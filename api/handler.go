package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"suspicious-ip-checker/config"
	"suspicious-ip-checker/kafka"
	"suspicious-ip-checker/service"
)

// RegisterRoutes, Fiber uygulaması için API uç noktalarını ayarlar.
// Fiber uygulama örneğini, bir Zap günlükçüsünü ve uygulama yapılandırmasını alır.
func RegisterRoutes(app *fiber.App, logger *zap.Logger, cfg *config.Config) {
	// /submit-ip adresinde bir POST uç noktası tanımlar.
	// Bu uç nokta, bir IP adresi almak, VirusTotal ile kontrol etmek,
	// sonucu Kafka'ya göndermek ve bir yanıt döndürmekten sorumludur.
	app.Post("/submit-ip", func(c *fiber.Ctx) error {
		// Gelen JSON istek gövdesi için yapıyı tanımlar.
		type request struct {
			IP string `json:"ip"` // Kontrol edilecek IP adresi.
		}
		var body request

		// İstek gövdesini 'body' yapısına ayrıştırır.
		// Ayrıştırma başarısız olursa, bir hata günlüğe kaydeder ve Bad Request durumu döndürür.
		if err := c.BodyParser(&body); err != nil {
			logger.Error("Geçersiz istek gövdesi", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Geçersiz istek",
			})
		}

		// Sağlanan IP adresini kontrol etmek için VirusTotal servisini çağırır.
		// Servis bir durum (örn. "malicious", "suspicious", "clean") ve varsa bir hata döndürür.
		status, err := service.CheckIP(body.IP, cfg)
		if err != nil {
			logger.Error("VirusTotal API hatası", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "IP kontrol edilemedi",
			})
		}

		// Kafka'ya gönderilecek tarama sonucunu hazırlar.
		// Bu, IP'yi, belirlenen durumunu, mevcut zaman damgasını ve kaynak hizmeti içerir.
		result := kafka.ScanResult{
			IP:        body.IP,
			Status:    status,
			Timestamp: time.Now(),
			Source:    "ip-submission-service",
		}

		// Tarama sonucunu Kafka konusuna gönderir.
		// Gönderme başarısız olursa, bir hata günlüğe kaydeder ve Internal Server Error durumu döndürür.
		if err := kafka.SendToKafka(cfg, logger, result); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Kafka'ya mesaj gönderilemedi",
			})
		}

		// Başarılı IP kontrolünü ve Kafka mesajı gönderimini günlüğe kaydet.
		logger.Info("IP kontrol edildi ve Kafka'ya gönderildi",
			zap.String("ip", body.IP),
			zap.String("status", status),
		)

		// Kontrol edilen IP'yi ve durumunu belirten bir JSON yanıtı döndür.
		return c.JSON(fiber.Map{
			"ip":     body.IP,
			"status": status,
		})
	})
}
