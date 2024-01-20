package admin

import (
	"bytes"
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"image/png"
	"log"
)

func (s *Server) RunOperation(ctx context.Context, req *RunOperationRequest) (*RunOperationResponse, error) {
	if err := s.fetchDiscordAvatar(ctx); err != nil {
		return nil, err
	}

	return &RunOperationResponse{}, nil
}

func (s *Server) fetchDiscordAvatar(ctx context.Context) error {
	rows, err := s.conn.Pool.Query(ctx, "SELECT id, discord_id FROM users")
	if err != nil {
		return err
	}

	var (
		ids  []pgtype.UUID
		dids []string
	)
	for rows.Next() {
		var (
			id  pgtype.UUID
			did string
		)
		if err := rows.Scan(&id, &did); err != nil {
			return err
		}
		ids = append(ids, id)
		dids = append(dids, did)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for i, did := range dids {
		u, err := s.dg.User(did)
		if err != nil {
			return err
		}

		if u.Avatar == "" {
			continue
		}

		var avatar pgtype.Text
		if err := avatar.Scan(u.Avatar); err != nil {
			return err
		}

		_, err = s.conn.Pool.Exec(ctx, "UPDATE users SET avatar = $2 WHERE discord_id = $1", did, avatar)
		if err != nil {
			return err
		}

		img, err := s.dg.UserAvatarDecode(u)
		if err != nil {
			return err
		}

		buf := new(bytes.Buffer)
		if err := png.Encode(buf, img); err != nil {
			return err
		}

		uuid, err := ids[i].Value()
		if err != nil {
			return err
		}

		if err := s.s3.SaveAvatar(ctx, uuid.(string), buf.Bytes()); err != nil {
			return err
		}
		log.Printf("Saved %s's avatar\n", u.Username)
	}

	return nil
}
