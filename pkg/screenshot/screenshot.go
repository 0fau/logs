package screenshot

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"github.com/cockroachdb/errors"

	"github.com/0fau/logs/pkg/database"
	"github.com/0fau/logs/pkg/s3"
)

type Config struct {
	FrontendURL string
	DatabaseURL string

	S3 *s3.Config
}

type Server struct {
	config *Config

	conn *database.DB
	s3   *s3.Client
}

func NewServer(c *Config) *Server {
	return &Server{config: c}
}

func (s *Server) Run(ctx context.Context) error {
	conn, err := database.Connect(ctx, s.config.DatabaseURL, "logs_screenshot", false)
	if err != nil {
		return errors.Wrap(err, "connecting to database")
	}
	s.conn = conn

	s.s3, err = s3.NewClient(s.config.S3)
	if err != nil {
		return errors.Wrap(err, "creating minio s3 client")
	}

	log.Println("Screenshot Service has started.")

	return s.Poll(ctx)
}

func (s *Server) Poll(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		time.Sleep(time.Millisecond * 500)

		ids, err := s.conn.Queries.GetNoThumbnailLogs(ctx)
		if len(ids) == 0 || err != nil {
			continue
		}

		sem := make(chan struct{}, 10)
		var wg sync.WaitGroup
		wg.Add(len(ids))

		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		for _, id := range ids {
			sem <- struct{}{}
			go func(enc int32) {
				defer wg.Done()
				s.Screenshot(ctx, enc)
				<-sem
			}(id)
		}
		wg.Wait()
		cancel()
	}
}

func (s *Server) Screenshot(ctx context.Context, enc int32) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
		chromedp.WindowSize(1920, 1080),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	chromedpCtx, cancel := chromedp.NewContext(
		allocCtx,
	)
	defer cancel()

	var buf []byte
	var width, height int64
	if err := chromedp.Run(
		chromedpCtx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			return chromedp.EmulateViewport(width, height, func(sdmop *emulation.SetDeviceMetricsOverrideParams, steep *emulation.SetTouchEmulationEnabledParams) {
				sdmop.DeviceScaleFactor = 4
			}).Do(ctx)
		}), elementScreenshot(fmt.Sprintf("%s/screenshot/log/%d", s.config.FrontendURL, enc), `.screenshot`, &buf)); err != nil {
		log.Println(errors.Wrapf(err, "Failed to screenshot encounter %d", enc))
		return
	}

	if err := s.s3.SaveImage(chromedpCtx, "thumbnail/"+fmt.Sprintf("%d", enc), buf); err != nil {
		log.Println(err)
		return
	}

	if err := s.conn.Queries.MarkThumbnail(ctx, enc); err != nil {
		log.Println(errors.Wrapf(err, "Failed to save encounter %d thumbnail", enc))
	} else {
		log.Printf("Saved encounter %d thumbnail\n", enc)
	}
}

func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}
