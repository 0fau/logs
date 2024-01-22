package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (s *Server) avatarHandler(c *gin.Context) {
	if c.Param("user") == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	image, err := s.s3.FetchAvatar(context.Background(), c.Param("user"))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Header("cache-control", "public, max-age=14400")
	c.Header("content-length", strconv.Itoa(len(image)))

	c.Data(http.StatusOK, "image/png", image)
}

func (s *Server) thumbnailHandler(c *gin.Context) {
	if c.Param("log") == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	image, err := s.s3.FetchImage(context.Background(), "thumbnail/"+c.Param("log"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Header("cache-control", "public, max-age=14400")
	c.Header("content-length", strconv.Itoa(len(image)))

	c.Data(http.StatusOK, "image/png", image)
}
