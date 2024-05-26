package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func (s *Server) LimitLog(c *gin.Context) error {
	ctx := context.Background()
	auth := c.GetHeader("Authorization")
	if auth != "" {
		if tok, ok := strings.CutPrefix(auth, "Bot: "); ok {
			_, err := s.conn.Queries.GetAPIToken(ctx, tok)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					c.AbortWithStatus(http.StatusUnauthorized)
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}
				return err
			}

		} else {
			c.AbortWithStatus(http.StatusBadRequest)
			return errors.New("invalid authorization header")
		}
	}

	if _, err := c.Cookie("sessions"); err != nil {
		sesh := sessions.Default(c)
		if val := sesh.Get("user"); val != nil {
			u := val.(*SessionUser)
			if val := sesh.Get("roles"); val != nil {
				roles := val.([]string)
				if hasRoles(roles, "alpha") {
					if err := s.Limit("User: "+u.ID, 300, time.Hour*24); err != nil {
						c.AbortWithStatus(http.StatusTooManyRequests)
						return err
					}
				}
			}
		}
	}
	return s.Limit(c.ClientIP(), 100, time.Hour*24)
}

func (s *Server) Limit(key string, limit int, dur time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC).Truncate(dur)
	end := date.Add(dur)

	key = "Rate Limit: [" + key + "] (" + end.Format("2006-01-02 15:04:05") + ")"
	res, err := s.redis.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return errors.Wrap(err, "fetching key from redis")
	}

	count := 0
	if res != "" {
		count, err = strconv.Atoi(res)
		if err != nil {
			return errors.Wrap(err, "converting string to int")
		}
	}

	if count >= limit {
		return errors.New("rate limit exceeded")
	}

	pipe := s.redis.TxPipeline()
	pipe.Incr(ctx, key)
	if count == 0 {
		pipe.Expire(ctx, key, dur)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return errors.Wrap(err, "executing pipeline")
	}

	return nil
}
