package admin

import (
	"cmp"
	"context"
	"github.com/0fau/logs/pkg/database"
	"github.com/0fau/logs/pkg/database/sql"
	"github.com/0fau/logs/pkg/database/sql/structs"
	"github.com/0fau/logs/pkg/process"
	"github.com/0fau/logs/pkg/process/meter"
	"github.com/0fau/logs/pkg/s3"
	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/cockroachdb/errors"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc"
	"log"
	"net"
	"runtime/debug"
	"slices"
	"sync"
	"time"
)

var _ AdminServer = (*Server)(nil)

type Config struct {
	Address string
}

type Server struct {
	config *Config

	db        *database.DB
	s3        *s3.Client
	processor *process.Processor

	UnimplementedAdminServer
}

func NewServer(c *Config, db *database.DB, s3 *s3.Client, processor *process.Processor) *Server {
	return &Server{config: c, db: db, s3: s3}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		return errors.Wrap(err, "listening on endpoint")
	}

	grpcServer := grpc.NewServer()
	RegisterAdminServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		return errors.Wrap(err, "grpc serve")
	}
	return nil
}

func (s *Server) Process(ctx context.Context, req *ProcessRequest) (*ProcessResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic from processing %d: %v\n%s\n", req.Encounter, r, string(debug.Stack()))
		}
	}()

	raw, err := s.s3.FetchEncounter(ctx, req.Encounter)
	if err != nil {
		return nil, errors.Wrap(err, "fetching encounter")
	}

	var enc *meter.Encounter
	if err := json.Unmarshal(raw, &enc); err != nil {
		return nil, errors.Wrap(err, "unmarshalling encounter")
	}

	s.processor.Preprocess(enc)

	proc, err := s.processor.Process(enc)
	if err != nil {
		return nil, errors.Wrap(err, "processing encounter")
	}

	order := make([]string, 0, len(proc.Header.Players))
	for player := range proc.Header.Players {
		order = append(order, player)
	}
	slices.SortFunc(order, func(a, b string) int {
		return cmp.Compare(
			proc.Header.Players[b].Damage,
			proc.Header.Players[a].Damage,
		)
	})

	hash := proc.UniqueHash(order)

	start := time.UnixMilli(enc.FightStart).UTC()
	var date pgtype.Timestamp
	if err := date.Scan(start); err != nil {
		return nil, errors.Wrap(err, "scanning duration pgtype.Timstamp")
	}

	if err := crdbpgx.ExecuteTx(ctx, s.db.Pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		qtx := s.db.Queries.WithTx(tx)

		user, err := qtx.ProcessEncounter(ctx, sql.ProcessEncounterParams{
			ID:         req.Encounter,
			Header:     proc.Header,
			Data:       proc.Data,
			UniqueHash: hash,
		})
		if err != nil {
			return errors.Wrap(err, "saving encounter")
		}

		group, err := qtx.GetUniqueGroup(ctx, sql.GetUniqueGroupParams{
			UniqueHash: hash,
			Date:       date,
			Duration:   enc.Duration,
		})
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return errors.Wrap(err, "getting unique group")
		}

		if group == 0 {
			group = req.Encounter
		} else if group != req.Encounter {
			if err := qtx.UpsertEncounterGroup(ctx, sql.UpsertEncounterGroupParams{
				GroupID: group,
				Column2: user,
			}); err != nil {
				return errors.Wrap(err, "upserting encounter group")
			}
		}

		if err := qtx.UpdateUniqueGroup(ctx, sql.UpdateUniqueGroupParams{
			ID:          req.Encounter,
			UniqueGroup: group,
		}); err != nil {
			return errors.Wrap(err, "updating unique group")
		}

		for i, player := range order {
			header := proc.Header.Players[player]
			params := sql.InsertPlayerInternalParams{
				Encounter: req.Encounter,
				Class:     header.Class,
				Name:      player,
				Dead:      header.Dead,
				Place:     int32(i + 1),
				Data: structs.IndexedPlayerData{
					Damage: header.Damage,
					DPS:    header.DPS,
				},
				Dps: header.DPS,
			}
			if err := qtx.InsertPlayerInternal(ctx, params); err != nil {
				return errors.Wrap(err, "inserting player")
			}
		}

		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "executing transaction")
	}

	return &ProcessResponse{}, nil
}

func (s *Server) ProcessHash(ctx context.Context, req *ProcessHashRequest) (*ProcessHashResponse, error) {
	row, err := s.db.Queries.GetHeader(ctx, req.Encounter)
	if err != nil {
		return nil, errors.Wrap(err, "getting encounter header")
	}

	players := make([]string, 0, len(row.Header.Players))
	for name := range row.Header.Players {
		players = append(players, name)
	}

	enc := process.Encounter{Header: row.Header, Raw: &meter.Encounter{CurrentBossName: row.Boss, Difficulty: row.Difficulty}}
	hash := enc.UniqueHash(players)

	if err := crdbpgx.ExecuteTx(ctx, s.db.Pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		qtx := s.db.Queries.WithTx(tx)
		_, err := tx.Exec(ctx, "UPDATE encounters SET unique_hash = $1 WHERE id = $2", hash, req.Encounter)
		if err != nil {
			return errors.Wrap(err, "updating encounter unique hash")
		}

		group, err := qtx.GetUniqueGroup(ctx, sql.GetUniqueGroupParams{
			UniqueHash: hash,
			Date:       row.Date,
			Duration:   row.Duration,
		})
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return errors.Wrap(err, "getting unique group")
		}

		if group == 0 {
			group = req.Encounter
		} else if group != req.Encounter {
			if err := qtx.UpsertEncounterGroup(ctx, sql.UpsertEncounterGroupParams{
				GroupID: group,
				Column2: row.UploadedBy,
			}); err != nil {
				return errors.Wrap(err, "upserting encounter group")
			}
		}

		if err := qtx.UpdateUniqueGroup(ctx, sql.UpdateUniqueGroupParams{
			ID:          req.Encounter,
			UniqueGroup: group,
		}); err != nil {
			return errors.Wrap(err, "updating unique group")
		}

		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "executing transaction")
	}
	return &ProcessHashResponse{}, nil
}

func (s *Server) ProcessAll(ctx context.Context, req *ProcessAllRequest) (*ProcessAllResponse, error) {
	ids, err := s.db.Queries.ListEncounters(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "listing encounter ids")
	}

	sem := make(chan struct{}, 1)
	var wg sync.WaitGroup
	wg.Add(len(ids))

	for i := 0; i < len(ids); i++ {
		sem <- struct{}{}
		go func(enc int32) {
			defer wg.Done()
			//if _, err := s.Process(ctx, &ProcessRequest{Encounter: enc}); err != nil {
			//	log.Println(err)
			//}

			if _, err := s.ProcessHash(ctx, &ProcessHashRequest{Encounter: enc}); err != nil {
				log.Println(err)
			}

			<-sem
		}(ids[i])
	}

	wg.Wait()
	return &ProcessAllResponse{}, nil
}

func (s *Server) Role(ctx context.Context, req *RoleRequest) (*RoleResponse, error) {
	user, err := s.db.Queries.GetUser(ctx, req.Discord)
	if err != nil {
		return nil, errors.Wrap(err, "fetch user")
	}

	roles := user.Roles
	switch req.Action {
	case RoleRequest_Add:
		if slices.Contains(roles, req.Role) {
			return &RoleResponse{}, nil
		}
		roles = append(roles, req.Role)
	case RoleRequest_Remove:
		if !slices.Contains(roles, req.Role) {
			return &RoleResponse{}, nil
		}
		roles = slices.DeleteFunc(roles, func(role string) bool {
			return role == req.Role
		})
	}

	if err := s.db.Queries.SetUserRoles(ctx, sql.SetUserRolesParams{
		DiscordTag: req.Discord,
		Roles:      roles,
	}); err != nil {
		return nil, errors.Wrap(err, "setting roles")
	}

	return &RoleResponse{}, nil
}

func (s *Server) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	if err := s.s3.DeleteEncounter(ctx, req.Encounter); err != nil {
		return nil, errors.Wrap(err, "s3 delete")
	}

	if err := s.db.Queries.DeleteEncounter(ctx, req.Encounter); err != nil {
		return nil, errors.Wrap(err, "db delete")
	}

	return &DeleteResponse{}, nil
}
