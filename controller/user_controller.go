package controller

import (
	"logger-example/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (uc *UserController) GetUser(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("id")

	start := time.Now()
	defer func() {
		latency := time.Since(start)
		logger.Info(ctx, "Request latency", zap.String("path", c.FullPath()), zap.String("latency", latency.String()))
	}()

	logger.Info(ctx, "Request started", zap.String("user_id", userID))

	time.Sleep(50 * time.Millisecond) // simulate DB query

	logger.Info(ctx, "Database query finished", zap.String("user_id", userID))

	if userID == "0" {
		logger.Error(ctx, http.ErrNoCookie, "User not found", zap.String("user_id", userID))
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":   userID,
		"name": "Jean Bonilha",
	})

	logger.Info(ctx, "Request completed successfully", zap.String("user_id", userID))
}
