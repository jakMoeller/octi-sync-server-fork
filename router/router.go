package router

import (
	"context"
	"fmt"
	"net/http"
	"octi-sync-server/config"
	"octi-sync-server/middleware/logging"
	requestmiddleware "octi-sync-server/middleware/request"
	v1 "octi-sync-server/router/v1"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

// New generates the router used in the HTTP Server.
func New(ctx context.Context, config *config.Config) http.Handler {
	router := gin.New()

	router.Use(requestmiddleware.LimitHandler(requestmiddleware.DefaultLimit()))

	// Global Middleware
	router.Use(
		ginzap.RecoveryWithZap(config.Logger, true),
		logging.RequestLogging(config.Logger),
	)

	router.Use(requestmiddleware.BodySizeLimiter(config.Server.MaxRequestBodySize, gin.H{
		"msg": fmt.Sprintf("request too large, maximum allowed is %v bytes", config.Server.MaxRequestBodySize),
	}))

	v1.New(ctx, router, config)

	router.GET("/health", healthCheck(config))

	return router
}

func healthCheck(_ *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"health": "up"})
	}
}
