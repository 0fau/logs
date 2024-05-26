package query

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/0fau/logs/pkg/database"
	"github.com/0fau/logs/pkg/process"
	structs2 "github.com/0fau/logs/pkg/process/structs"

	sq "github.com/Masterminds/squirrel"
	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5/pgtype"
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

	Search string `json:"search"`
}

type Params struct {
	User      pgtype.UUID
	Order     string
	Scope     string
	GearScore string
	PastID    *int32
	PastField *int64
	PastPlace *int32

	Selections Selections
	Privileged bool
	Admin      bool
}

type User struct {
	DiscordTag    string
	Username      pgtype.Text
	LogVisibility *structs2.EncounterVisibility
}

type Encounter struct {
	Uploader User

	ID            int32
	Difficulty    string
	UploadedBy    pgtype.UUID
	UploadedAt    pgtype.Timestamp
	Visibility    *structs2.EncounterVisibility
	Settings      structs2.EncounterSettings
	Tags          []string
	Header        structs2.EncounterHeader
	Boss          string
	Date          pgtype.Timestamp
	Duration      int32
	FocusedPlayer string
	LocalPlayer   string
	Place         int32
}

func Query(db *database.DB, params *Params) ([]Encounter, error) {
	if params.Selections.Search != "" && !params.User.Valid {
		return nil, nil
	}

	focused := len(params.Selections.Classes) > 0 ||
		params.Order == "performance" ||
		params.Selections.Search != "" ||
		params.GearScore != ""

	selects := []string{
		"u.discord_tag",
		"u.username",
		"u.log_visibility",
		"e.id",
		"e.difficulty",
		"e.uploaded_by",
		"e.uploaded_at",
		"e.visibility",
		"e.settings",
		"e.tags",
		"e.header",
		"e.boss",
		"e.date",
		"e.duration",
		"e.local_player",
	}
	if focused {
		selects = append(selects, "p.name", "p.place")
	}

	q := sq.Select(selects...)

	if focused {
		q = q.From("players p").
			Join("encounters e ON e.id = encounter")
	} else {
		q = q.From("encounters e")
	}
	q = q.Join("users u ON u.id = e.uploaded_by")

	if params.Scope == "friends" {
		q = q.Join("friends f ON f.user1 = u.id AND f.user2 = ?", params.User)
	} else {
		q = q.LeftJoin("grouped_encounters g ON e.unique_group = g.group_id")
	}

	bossColumn, difficultyColumn := "e.boss", "e.difficulty"
	if focused && !slices.Contains([]string{"duration", "recent log"}, params.Order) {
		bossColumn = "p.boss"
		difficultyColumn = "p.difficulty"
	}

	difficulties := map[string]struct{}{}
	for _, raid := range params.Selections.Raids {
		if len(raid.Gates) > 0 && len(raid.Difficulties) == 0 {
			difficulties[""] = struct{}{}
		}

		for _, difficulty := range raid.Difficulties {
			difficulties[difficulty] = struct{}{}
		}
	}

	selection := sq.Or{}
	if len(difficulties) == 1 {
		var bosses []string
		for name, raid := range params.Selections.Raids {
			if len(raid.Gates) == 0 && len(raid.Difficulties) == 0 {
				continue
			}

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
		}

		difficulty := ""
		for d := range difficulties {
			difficulty = d
			break
		}

		var cmp interface{}
		if len(bosses) == 1 {
			cmp = bosses[0]
		} else {
			cmp = bosses
		}

		if difficulty == "" {
			selection = append(selection, sq.Eq{bossColumn: cmp})
		} else {
			selection = append(selection, sq.And{
				sq.Eq{
					bossColumn: cmp,
				},
				sq.Eq{
					difficultyColumn: difficulty,
				},
			})
		}
	} else {
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
					sq.Eq{bossColumn: bosses},
					sq.Eq{difficultyColumn: raid.Difficulties},
				})
			} else {
				selection = append(selection, sq.Eq{
					bossColumn: bosses,
				})
			}
		}
	}

	if len(params.Selections.Guardians) > 0 {
		selection = append(selection, sq.And{
			sq.Eq{
				bossColumn: params.Selections.Guardians,
			},
			sq.Eq{
				difficultyColumn: "Normal",
			},
		})
	}

	if len(params.Selections.Trials) > 0 {
		selection = append(selection, sq.And{
			sq.Eq{
				bossColumn: params.Selections.Trials,
			},
			sq.Eq{
				difficultyColumn: "Trial",
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

	if focused {
		if params.Selections.Search != "" {
			q = q.Where(sq.Eq{"p.name": params.Selections.Search})
		} else if params.Scope != "arkesia" {
			q = q.Where("p.name = e.local_player")
		}
	}

	switch params.Order {
	case "", "recent clear":
		if params.PastField != nil && params.PastID != nil && (!focused || focused && params.PastPlace != nil) {
			date := pgtype.Timestamp{
				Time:  time.UnixMilli(*params.PastField).UTC(),
				Valid: true,
			}

			or := sq.Or{
				sq.Lt{"e.date": date},
				sq.And{
					sq.Eq{"e.date": date},
					sq.Gt{"e.id": *params.PastID},
				},
			}

			if focused {
				or = append(or, sq.And{
					sq.Eq{"e.date": date},
					sq.Eq{"e.id": *params.PastID},
					sq.Gt{"p.place": *params.PastPlace},
				})
			}

			q = q.Where(or)
		}
		q = q.OrderBy("e.date DESC", "e.id ASC")
	case "recent log":
		if params.PastID != nil && (!focused || focused && params.PastPlace != nil) {
			pastID := sq.Lt{"e.id": *params.PastID}

			if focused {
				q = q.Where(sq.Or{
					pastID,
					sq.And{
						sq.Eq{"e.id": *params.PastID},
						sq.Gt{"p.place": *params.PastPlace},
					},
				})
			} else {
				q = q.Where(pastID)
			}
		}
		q = q.OrderBy("e.id DESC")
	case "raid duration":
		if params.PastField != nil && params.PastID != nil && (!focused || focused && params.PastPlace != nil) {
			or := sq.Or{
				sq.Gt{"e.duration": *params.PastField},
				sq.And{
					sq.Eq{"e.duration": *params.PastField},
					sq.Gt{"e.id": *params.PastID},
				},
			}

			if focused {
				or = append(or, sq.And{
					sq.Eq{"e.duration": *params.PastField},
					sq.Eq{"e.id": *params.PastID},
					sq.Gt{"p.place": *params.PastPlace},
				})
			}

			q = q.Where(or)
		}
		q = q.OrderBy("e.duration ASC", "e.id ASC")
	case "performance":
		if params.PastField != nil && params.PastID != nil && params.PastPlace != nil {
			q = q.Where(sq.Or{
				sq.Lt{"p.dps": *params.PastField},
				sq.And{
					sq.Eq{"p.dps": *params.PastField},
					sq.Gt{"p.place": *params.PastID},
				},
				sq.And{
					sq.Eq{"p.dps": *params.PastField},
					sq.Eq{"p.place": *params.PastID},
					sq.Gt{"p.place": *params.PastPlace},
				},
			})
		}
		q = q.OrderBy("p.dps DESC")
	}

	if focused {
		q = q.OrderBy("p.place ASC")
	}

	switch params.Scope {
	case "arkesia":
		if params.User.Valid {
			q = q.Where(sq.Or{
				sq.Eq{"g.uploaders": nil},
				sq.Eq{"e.uploaded_by": params.User},
				sq.And{
					sq.Expr("e.id = e.unique_group"),
					sq.Expr("NOT (?::UUID = ANY (g.uploaders))", params.User),
				},
			})
		} else {
			q = q.Where(sq.Or{
				sq.Eq{"g.uploaders": nil},
				sq.Expr("e.id = e.unique_group"),
			})
		}
	}

	if params.Scope == "roster" {
		q = q.Where(sq.Eq{"u.id": params.User})
	} else if params.Scope == "arkesia" && params.Selections.Search != "" && !params.Privileged {
		q = q.Where(sq.Or{
			sq.Eq{"u.id": params.User},
			sq.And{
				sq.NotEq{"e.visibility": nil},
				sq.Or{
					sq.Eq{"(e.visibility->'names')::INT": structs2.ShowNames},
					sq.And{
						sq.Eq{"(e.visibility->'names')::INT": structs2.ShowSelf},
						sq.Eq{"e.local_player": "p.name"},
						sq.Expr("e.local_player = p.name"),
					},
				},
			},
			sq.And{
				sq.Eq{"e.visibility": nil},
				sq.NotEq{"u.log_visibility": nil},
				sq.Or{
					sq.Eq{"(u.log_visibility->'names')::INT": structs2.ShowNames},
					sq.And{
						sq.Eq{"(u.log_visibility->'names')::INT": structs2.ShowSelf},
						sq.Expr("e.local_player = p.name"),
					},
				},
			},
		})
	}

	if params.GearScore != "" && slices.Contains([]string{
		"1540-1559", "1560-1579", "1580-1599", "1600-1609", "1610-1619", "1620+",
	}, params.GearScore) {
		if params.GearScore == "1620+" {
			q = q.Where(sq.GtOrEq{"p.gear_score": 1620})
		} else {
			var min, max int
			_, err := fmt.Sscanf(params.GearScore, "%d-%d", &min, &max)
			if err == nil {
				q = q.Where(sq.And{
					sq.GtOrEq{"p.gear_score": min},
					sq.Lt{"p.gear_score": max + 1},
				})
			}
		}
	}

	//if !params.Admin {
	//	q = q.Where(sq.Or{
	//		sq.Eq{"e.private": false},
	//		sq.Eq{"e.uploaded_by": params.User},
	//	})
	//}

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
			&enc.Uploader.LogVisibility,
			&enc.ID,
			&enc.Difficulty,
			&enc.UploadedBy,
			&enc.UploadedAt,
			&enc.Visibility,
			&enc.Settings,
			&enc.Tags,
			&enc.Header,
			&enc.Boss,
			&enc.Date,
			&enc.Duration,
			&enc.LocalPlayer,
		}
		if focused {
			scan = append(scan, &enc.FocusedPlayer, &enc.Place)
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
