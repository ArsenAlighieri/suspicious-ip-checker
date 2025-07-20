package kafka

import (
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"suspicious-ip-checker/config"
)

// ScanResult, Kafka'ya gönderilen verinin yapısını temsil eder.
type ScanResult struct {
	IP        string    `json:"ip"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
}

// SendToKafka, yapılandırılmış Kafka broker'ına bir ScanResult mesajı gönderir.
func SendToKafka(cfg *config.Config, logger *zap.Logger, result ScanResult) error {
	// Yeni bir Sarama senkronize üretici oluştur.
	// Bu üretici, mesaj başarıyla gönderilene veya bir hata oluşana kadar engellenir.
	producer, err := sarama.NewSyncProducer([]string{cfg.KafkaBroker}, nil)
	if err != nil {
		logger.Error("Kafka üreticisi oluşturulamadı", zap.Error(err))
		return err
	}
	// Fonksiyon çıkışında üreticinin kapatıldığından emin ol.
	defer producer.Close()

	// ScanResult yapısını bir JSON bayt dizisine dönüştür.
	// Bu, mesajı Kafka üzerinden iletim için hazırlar.
	messageBytes, err := json.Marshal(result)
	if err != nil {
		logger.Error("Kafka mesajı serileştirilemedi", zap.Error(err))
		return err
	}

	// Yeni bir Kafka üretici mesajı oluştur.
	// Mesaj, "ip_scan_result" konusuna JSON kodlu ScanResult değeri olarak gönderilir.
	msg := &sarama.ProducerMessage{
		Topic: "ip_scan_result",
		Value: sarama.ByteEncoder(messageBytes),
	}

	// Mesajı Kafka'ya gönder.
	// Bu işlem, mesajın depolandığı bölümü ve ofseti veya bir hata döndürür.
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		logger.Error("Kafka'ya mesaj gönderilemedi", zap.Error(err))
		return err
	}

	logger.Info("Mesaj Kafka'ya gönderildi",
		zap.String("topic", "ip_scan_result"),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
	)

	return nil
}
