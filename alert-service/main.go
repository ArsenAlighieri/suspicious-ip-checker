package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"suspicious-ip-checker/alert-service/config"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	kafkaTopic = "ip_scan_result"
)

// ScanResult yapısı, üreticinin mesaj formatıyla eşleşir.
type ScanResult struct {
	IP        string    `json:"ip"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
}

func main() {
	// JSON çıktısı için Zap günlükçüsünü başlat.
	zapCfg := zap.NewProductionConfig()
	zapCfg.Encoding = "json"
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := zapCfg.Build()
	if err != nil {
		log.Fatalf("Günlükleyici oluşturulamadı: %v", err)
	}
	defer logger.Sync() // Çıkışta arabelleğe alınmış günlükleri temizle.

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Yapılandırma yüklenemedi", zap.Error(err))
	}

	logger.Info("Uyarı Servisi başlatılıyor...")

	// Kafka tüketici yapılandırması.
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest // En eski mevcut ofsetten tüketmeye başla.

	// Yeni bir tüketici oluştur.
	master, err := sarama.NewConsumer([]string{cfg.KafkaBroker}, saramaConfig)
	if err != nil {
		logger.Fatal("Kafka tüketicisi başlatılamadı", zap.Error(err))
	}
	defer func() {
		if err := master.Close(); err != nil {
			logger.Error("Kafka tüketicisi kapatılamadı", zap.Error(err))
		}
	}()

	// Konudan mesajları tüket.
	consumer, err := master.ConsumePartition(kafkaTopic, 0, sarama.OffsetOldest)
	if err != nil {
		logger.Fatal("Bölüm tüketilemedi", zap.Error(err))
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Mesajları tüketmek için ana döngü.
	for {
		select {
		case msg := <-consumer.Messages():
			var result ScanResult
			if err := json.Unmarshal(msg.Value, &result); err != nil {
				logger.Error("Kafka mesajı ayrıştırılamadı", zap.Error(err), zap.ByteString("message", msg.Value))
				continue
			}
			// Mesajı JSON formatında günlüğe kaydet.
			logger.Info("IP tarama sonucu alındı",
				zap.String("ip", result.IP),
				zap.String("status", result.Status),
				zap.Time("timestamp", result.Timestamp),
				zap.String("source", result.Source),
				zap.String("topic", msg.Topic),
				zap.Int32("partition", msg.Partition),
				zap.Int64("offset", msg.Offset),
			)
		case err := <-consumer.Errors():
			logger.Error("Kafka tüketici hatası", zap.Error(err))
		case <-signals:
			logger.Info("Uyarı Servisi kapatılıyor...")
			return
		}
	}
}
