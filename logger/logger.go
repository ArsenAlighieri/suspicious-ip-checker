package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger, yeni bir Zap günlükleyici örneği başlatır ve döndürür.
// Günlükleme seviyesi, sağlanan 'level' dizesiyle belirlenir.
// Günlükleyici, günlükleri JSON biçiminde çıktı verecek şekilde yapılandırılmıştır.
func NewLogger(level string) *zap.Logger {
	// Zap günlükleyici için yeni bir üretim yapılandırması oluştur.
	// Bu yapılandırma, üretim ortamlarında yüksek performanslı günlükleme için optimize edilmiştir.
	cfg := zap.NewProductionConfig()

	// Yapılandırılmış günlükleme için kodlamayı JSON olarak ayarla.
	cfg.Encoding = "json"

	// Zaman damgaları için ISO8601 biçimini kullanmak üzere zaman kodlayıcısını yapılandır.
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Giriş dizesine göre günlükleme seviyesini ayarla.
	// Bu, günlüklerin ayrıntı düzeyini dinamik olarak kontrol etmeyi sağlar.
	switch level {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		// Bilinmeyen veya boş bir seviye sağlanırsa varsayılan olarak bilgi seviyesine ayarla.
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Yapılandırmadan günlükleyici örneğini oluştur.
	// Bu bağlamda basitlik için Build() tarafından döndürülen hatayı yoksaymak için alt çizgi kullanılır.
	logger, _ := cfg.Build()
	return logger
}
