package config

import (
	"github.com/spf13/viper"
)

// Config, uygulamanın yapılandırma parametrelerini tutar.
type Config struct {
	ServerPort       string // HTTP sunucusunun dinleyeceği port.
	VirusTotalApiKey string // VirusTotal servisi için API anahtarı.
	KafkaBroker      string // Kafka broker'ının adresi.
	LogLevel         string // Günlükleme seviyesi (örn. "debug", "info", "error").
}

// LoadConfig, yapılandırmayı config.yaml dosyasından ve ortam değişkenlerinden okur.
// Bir Config struct işaretçisi ve yükleme başarısız olursa bir hata döndürür.
func LoadConfig() (*Config, error) {
	// Yapılandırma dosyasının adını (uzantısız) ayarla.
	viper.SetConfigName("config")
	// Yapılandırma dosyasının türünü ayarla.
	viper.SetConfigType("yaml")
	// Yapılandırma dosyasının arama yoluna mevcut dizini ekle.
	viper.AddConfigPath(".")

	// Ortam değişkenlerinin otomatik olarak okunmasını etkinleştir.
	// Ortam değişkenleri, eşleşmeleri halinde yapılandırma dosyasındaki değerleri geçersiz kılar.
	viper.AutomaticEnv() // .env dosyalarını destekler (başka yollarla yüklenirse) veya doğrudan ortam değişkenlerini.

	// KAFKA_BROKER ortam değişkenini "kafka.broker" yapılandırma anahtarına bağla.
	// Bu, KAFKA_BROKER ortam değişkeninin Kafka broker adresi için config.yaml'deki değerden
	// öncelikli olmasını sağlar.
	viper.BindEnv("kafka.broker", "KAFKA_BROKER")

	// VIRUSTOTAL_API_KEY ortam değişkenini "virustotal.api_key" yapılandırma anahtarına bağla.
	// Bu, VIRUSTOTAL_API_KEY ortam değişkeninin VirusTotal API anahtarı için config.yaml'deki değerden
	// öncelikli olmasını sağlar.
	viper.BindEnv("virustotal.api_key", "VIRUSTOTAL_API_KEY")

	// Yapılandırma dosyasını oku.
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Viper'dan okunan değerlerle Config struct'ını doldur.
	cfg := &Config{
		ServerPort:       viper.GetString("server.port"),
		VirusTotalApiKey: viper.GetString("virustotal.api_key"),
		KafkaBroker:      viper.GetString("kafka.broker"),
		LogLevel:         viper.GetString("log.level"),
	}

	return cfg, nil
}
