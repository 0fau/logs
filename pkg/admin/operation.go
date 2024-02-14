package admin

import (
	"bytes"
	"context"
	"fmt"
	"github.com/0fau/logs/pkg/database/sql/structs"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"image/png"
	"log"
	"sync"
)

func (s *Server) RunOperation(ctx context.Context, req *RunOperationRequest) (*RunOperationResponse, error) {
	s.populatePlayerData()
	return &RunOperationResponse{}, nil
}

func (s *Server) populatePlayerData() error {
	for {
		rows, err := s.conn.Pool.Query(context.Background(), "SELECT id, header FROM encounters WHERE version = 0 ORDER BY id DESC LIMIT 50")
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				log.Println(err)
			}
			break
		}

		has := false

		var wg sync.WaitGroup
		for rows.Next() {
			has = true
			var id int32
			var header structs.EncounterHeader
			if err := rows.Scan(&id, &header); err != nil {
				return err
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				crdbpgx.ExecuteTx(context.Background(), s.conn.Pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
					for name, player := range header.Players {
						if _, err := tx.Exec(
							context.Background(),
							"UPDATE players SET gear_score = $3 WHERE encounter = $1 AND name = $2",
							id, name, player.GearScore,
						); err != nil {
							return err
						}
					}

					_, err = tx.Exec(context.Background(), "UPDATE encounters SET version = 1 WHERE id = $1", id)
					return err
				})
			}()
		}
		if err := rows.Err(); err != nil {
			return err
		}

		if !has {
			break
		}

		wg.Wait()
	}
	return nil
}

func (s *Server) generateLogThumbnail() error {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
		chromedp.WindowSize(1920, 1080),
	)

	rows, err := s.conn.Pool.Query(context.Background(), "SELECT id FROM encounters WHERE thumbnail = false ORDER BY id DESC LIMIT 10000")
	if err != nil {
		return err
	}

	ids := []int32{}
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	sem := make(chan struct{}, 10)
	var wg sync.WaitGroup
	wg.Add(len(ids))
	for _, id := range ids {
		sem <- struct{}{}
		go func(enc int32) {
			defer wg.Done()

			allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
			defer cancel()

			ctx, cancel := chromedp.NewContext(
				allocCtx,
			)
			defer cancel()

			var buf []byte
			var width, height int64
			if err := chromedp.Run(ctx,
				chromedp.ActionFunc(func(ctx context.Context) error {
					return chromedp.EmulateViewport(width, height, func(sdmop *emulation.SetDeviceMetricsOverrideParams, steep *emulation.SetTouchEmulationEnabledParams) {
						sdmop.DeviceScaleFactor = 3
					}).Do(ctx)
				}), elementScreenshot(`http://logsbyfaust.logsbyfaust.svc.cluster.local:3000/screenshot/log/`+fmt.Sprintf("%d", enc), `.screenshot`, &buf)); err != nil {
				log.Fatal(err)
			}

			if err := s.s3.SaveImage(ctx, "thumbnail/"+fmt.Sprintf("%d", enc), buf); err != nil {
				log.Println(err)
			}

			if _, err := s.conn.Pool.Exec(context.Background(), "UPDATE encounters SET thumbnail = true WHERE id = $1", enc); err != nil {
				log.Println(err)
			} else {
				log.Printf("Saved encounter %d thumbnail\n", enc)
			}
			<-sem
		}(id)
	}
	wg.Wait()

	return nil
}

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
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
