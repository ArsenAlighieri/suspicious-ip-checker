package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"suspicious-ip-checker/api"
	"suspicious-ip-checker/config"
	"suspicious-ip-checker/logger"
)

// main fonksiyonu, IP Gönderim Servisi'nin giriş noktasıdır.
func main() {
	// Uygulama yapılandırmasını config.yaml dosyasından ve ortam değişkenlerinden yükle.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Yapılandırma yüklenemedi: %v", err)
	}

	// Yapılandırılan günlük seviyesiyle Zap günlükçüsünü başlat.
	// Bu günlükçü, yapılandırılmış günlükleme için uygulama genelinde kullanılacaktır.
	log := logger.NewLogger(cfg.LogLevel)

	// Yeni bir Fiber web uygulama örneği oluştur.
	// Fiber, Go için hızlı, görüşsüz ve minimalist bir web çerçevesidir.
	app := fiber.New()

	// Fiber uygulaması için API rotalarını kaydet.
	// Bu, /submit-ip uç noktasını ve işleyicisini ayarlar.
	api.RegisterRoutes(app, log, cfg)

	// Fiber API sunucusunu başlat ve yapılandırılan portta gelen istekleri dinle.
	log.Info("API başlatıldı", zap.String("port", cfg.ServerPort))
	if err := app.Listen(":" + cfg.ServerPort); err != nil {
		log.Fatal("API başlatılamadı", zap.Error(err))
	}
}
