package api

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
)

func (s *Server) userHandler(c *gin.Context) {
	ctx := context.Background()
	user, err := s.conn.Queries.GetUser(ctx, c.Param("user"))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			log.Println(errors.Wrap(err, "get user by name"))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	uuid, _ := user.ID.Value()
	c.JSON(http.StatusOK, ReturnedUser{
		ID:       uuid.(string),
		Username: c.Param("user"),
	})
}
