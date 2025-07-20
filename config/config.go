package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort       string
	VirusTotalApiKey string
	KafkaBroker      string
	LogLevel         string
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
	viper.BindEnv("kafka.broker", "KAFKA_BROKER")
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
