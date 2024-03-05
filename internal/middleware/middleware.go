package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/vermacodes/email-alert-processor/internal/config"
)

func AuthMiddleware(appConfig *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Debug("AuthMiddleware")

		// if request is for /healthz, skip auth
		if c.Request.URL.Path == "/healthz" {
			c.Next()
			return
		}

		apiKey := c.GetHeader("ApiKey")
		if apiKey == "" {
			slog.Error("API Key not found")
			c.AbortWithStatusJSON(401, gin.H{"error": "API Key not found"})
			return
		}

		if apiKey != appConfig.APIKey {
			slog.Error("Invalid API Key")
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid API Key"})
			return
		}

		c.Next()
	}
}
