package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"suspicious-ip-checker/api"
	"suspicious-ip-checker/config"
	"suspicious-ip-checker/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Yapılandırma yüklenemedi: %v", err)
	}

	log := logger.NewLogger(cfg.LogLevel)

	app := fiber.New()

	api.RegisterRoutes(app, log, cfg)

	log.Info("API başlatıldı", zap.String("port", cfg.ServerPort))
	if err := app.Listen(":" + cfg.ServerPort); err != nil {
		log.Fatal("API başlatılamadı", zap.Error(err))
	}
}
