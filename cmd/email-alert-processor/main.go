package main

import (
	"log/slog"

	"github.com/vermacodes/email-alert-processor/internal/config"
	"github.com/vermacodes/email-alert-processor/internal/handler"
	"github.com/vermacodes/email-alert-processor/internal/middleware"
	"github.com/vermacodes/email-alert-processor/internal/router"
	"github.com/vermacodes/email-alert-processor/internal/service"
)

func main() {
	appConfig, err := config.NewConfig()
	if err != nil {
		slog.Error("Error reading config: ", err)
		panic(err)
	}

	alertService := service.NewAlertService()

	router := router.NewDefaultRouter()
	router.Use(middleware.AuthMiddleware(appConfig))

	handler.NewEmailAlertProcessingHandler(router.Group("/"), alertService)

	port := appConfig.Port
	router.Run(":" + port)
}
