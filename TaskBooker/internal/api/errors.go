package api

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	log.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
