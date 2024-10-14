package main

import (
	"book-reading-tracker/config"
	"book-reading-tracker/internal/app/routes"
	"book-reading-tracker/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	utils.InitLogger(cfg.LogLevel)

	utils.Logger.Info("Starting the application...")

	r := gin.Default()

	routes.SetupRoutes(r)

	utils.Logger.Infof("Server is running on %s", cfg.ServerAddress)
	if err := r.Run(cfg.ServerAddress); err != nil {
		utils.LogError("Failed to start the server", err)
		panic(err)
	}
}
