package query

import (
	"context"
	"github.com/0fau/logs/pkg/database"
	"github.com/0fau/logs/pkg/database/sql/structs"
	"github.com/0fau/logs/pkg/process"
	sq "github.com/Masterminds/squirrel"
	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type Raid struct {
	Gates        []int    `json:"gates"`
	Difficulties []string `json:"difficulties"`
}

type Selections struct {
	Raids     map[string]Raid `json:"raids"`
	Guardians []string        `json:"guardians"`
	Trials    []string        `json:"trials"`
	Classes   []string        `json:"classes"`
}

type Params struct {
	User      string
	Order     string
	Scope     string
	PastID    *int32
	PastField *int64

	Selections Selections
}

type User struct {
	DiscordTag string
	Username   pgtype.Text
}

type Encounter struct {
	Uploader User

	ID          int32
	Difficulty  string
	UploadedBy  pgtype.UUID
	UploadedAt  pgtype.Timestamp
	Settings    structs.EncounterSettings
	Tags        []string
	Header      structs.EncounterHeader
	Boss        string
	Date        pgtype.Timestamp
	Duration    int32
	LocalPlayer string
	Place       int32
}

func Query(db *database.DB, params *Params) ([]Encounter, error) {
	player := "e.local_player"

	focused := len(params.Selections.Classes) > 0 || params.Order == "performance"
	if focused {
		player = "p.name"
	}

	selects := []string{
		"u.discord_tag",
		"u.username",
		"e.id",
		"e.difficulty",
		"e.uploaded_by",
		"e.uploaded_at",
		"e.settings",
		"e.tags",
		"e.header",
		"e.boss",
		"e.date",
		"e.duration",
		player,
	}
	if focused {
		selects = append(selects, "p.place")
	}

	q := sq.Select(selects...)

	if focused {
		q = q.From("players p").
			Join("encounters e ON e.id = encounter")
	} else {
		q = q.From("encounters e")
	}
	q = q.Join("users u ON u.id = e.uploaded_by")

	var uuid pgtype.UUID
	if params.User != "" {
		if err := uuid.Scan(params.User); err != nil {
			return nil, errors.Wrap(err, "scanning user uuid")
		}
	}

	switch params.Scope {
	case "friends":
		q = q.Where("?::UUID = ANY (u.friends)", uuid)
	case "roster":
		q = q.Where(sq.Eq{"u.id": uuid})
	}

	selection := sq.Or{}
	for name, raid := range params.Selections.Raids {
		if len(raid.Gates) == 0 && len(raid.Difficulties) == 0 {
			continue
		}

		bosses := make([]string, len(raid.Gates))
		if len(raid.Gates) == 0 && len(raid.Difficulties) > 0 {
			gates := process.Raids[name]
			for _, bs := range gates {
				bosses = append(bosses, bs...)
			}
		} else {
			for _, gate := range raid.Gates {
				l := len(process.Raids[name])
				if gate < 0 || gate > l {
					continue
				}

				bosses = append(bosses, process.Raids[name][gate-1]...)
			}
		}

		if len(raid.Difficulties) > 0 {
			selection = append(selection, sq.And{
				sq.Eq{"e.boss": bosses},
				sq.Eq{"e.difficulty": raid.Difficulties},
			})
		} else {
			selection = append(selection, sq.Eq{
				"e.boss": bosses,
			})
		}
	}

	if len(params.Selections.Guardians) > 0 {
		selection = append(selection, sq.And{
			sq.Eq{
				"e.boss": params.Selections.Guardians,
			},
			sq.Eq{
				"e.difficulty": "Normal",
			},
		})
	}

	if len(params.Selections.Trials) > 0 {
		selection = append(selection, sq.And{
			sq.Eq{
				"e.boss": params.Selections.Trials,
			},
			sq.Eq{
				"e.difficulty": "Trial",
			},
		})
	}

	if len(selection) == 1 {
		q = q.Where(selection[0])
	} else if len(selection) > 1 {
		q = q.Where(selection)
	}

	if len(params.Selections.Classes) > 0 {
		q = q.Where(sq.Eq{
			"p.class": params.Selections.Classes,
		})
	}

	if focused && params.Scope != "arkesia" {
		q = q.Where("p.name = e.local_player")
	}

	switch params.Order {
	case "", "recent clear":
		if params.PastField != nil && params.PastID != nil {
			date := pgtype.Timestamp{
				Time:  time.UnixMilli(*params.PastField).UTC(),
				Valid: true,
			}

			q = q.Where(sq.Or{
				sq.Lt{"e.date": date},
				sq.And{
					sq.Eq{"e.date": date},
					sq.Gt{"e.id": *params.PastID},
				},
			})
		}
		q = q.OrderBy("e.date DESC", "e.id ASC")
	case "recent log":
		if params.PastID != nil {
			q = q.Where(sq.Lt{"e.id": *params.PastID})
		}
		q = q.OrderBy("e.id DESC")
	case "raid duration":
		if params.PastField != nil && params.PastID != nil {
			q = q.Where(sq.Or{
				sq.Gt{"e.duration": *params.PastField},
				sq.And{
					sq.Eq{"e.duration": *params.PastField},
					sq.Gt{"e.id": *params.PastID},
				},
			})
		}
		q = q.OrderBy("e.duration ASC", "e.id ASC")
	case "performance":
		if params.PastField != nil && params.PastID != nil {
			q = q.Where(sq.Or{
				sq.Lt{"((p.data->>'dps')::BIGINT)": *params.PastField},
				sq.And{
					sq.Eq{"((p.data->>'dps')::BIGINT)": *params.PastField},
					sq.Gt{"p.place": *params.PastID},
				},
			})
		}
		q = q.OrderBy("((p.data->>'dps')::BIGINT) DESC", "p.place ASC")
	}

	q = q.Limit(6).PlaceholderFormat(sq.Dollar)

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "to sql")
	}
	rows, err := db.Pool.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}
	defer rows.Close()

	var encs []Encounter
	for rows.Next() {
		var enc Encounter
		scan := []interface{}{
			&enc.Uploader.DiscordTag,
			&enc.Uploader.Username,
			&enc.ID,
			&enc.Difficulty,
			&enc.UploadedBy,
			&enc.UploadedAt,
			&enc.Settings,
			&enc.Tags,
			&enc.Header,
			&enc.Boss,
			&enc.Date,
			&enc.Duration,
			&enc.LocalPlayer,
		}
		if focused {
			scan = append(scan, &enc.Place)
		}

		if err := rows.Scan(scan...); err != nil {
			return nil, errors.Wrap(err, "scan err")
		}
		encs = append(encs, enc)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows err")
	}

	return encs, nil
}
