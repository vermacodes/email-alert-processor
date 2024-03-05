package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewDefaultRouter() *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Authorization", "Content-Type", "ApiKey"}

	router.Use(cors.New(config))

	router.Group("/")

	return router
}
