package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"suspicious-ip-checker/config"
)

// VTResponse, bir IP adresini sorgularken VirusTotal API'sinden gelen JSON yanıtının yapısını temsil eder.
// Analiz istatistiklerine erişmek için iç içe geçmiş yapılar içerir.
type VTResponse struct {
	Data struct {
		Attributes struct {
			LastAnalysisStats struct {
				Malicious  int `json:"malicious"`  // IP'yi kötü amaçlı olarak işaretleyen motor sayısı.
				Suspicious int `json:"suspicious"` // IP'yi şüpheli olarak işaretleyen motor sayısı.
				Harmless   int `json:"harmless"`   // IP'yi zararsız olarak işaretleyen motor sayısı.
			} `json:"last_analysis_stats"`
		} `json:"attributes"`
	} `json:"data"`
}

// CheckIP, verilen IP adresi için VirusTotal API'sini sorgular.

func CheckIP(ip string, cfg *config.Config) (string, error) {
	// Asılı kalan istekleri önlemek için zaman aşımı olan bir HTTP istemcisi oluştur.
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	// IP adresi araması için VirusTotal API URL'sini oluştur.
	url := fmt.Sprintf("https://www.virustotal.com/api/v3/ip_addresses/%s", ip)

	// Yeni bir HTTP GET isteği oluştur.
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// "x-apikey" başlığını yapılandırmadan alınan VirusTotal API anahtarıyla ayarla.
	req.Header.Set("x-apikey", cfg.VirusTotalApiKey)

	// HTTP isteğini yürüt.
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	// Okuduktan sonra yanıt gövdesinin kapatıldığından emin ol.
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("VirusTotal API hatası: %d - %s", resp.StatusCode, string(body))
	}

	// JSON yanıt gövdesini VTResponse yapısına dönüştür.
	// Dönüştürme başarısız olursa, bir hata döndür.
	var vtResp VTResponse
	if err := json.NewDecoder(resp.Body).Decode(&vtResp); err != nil {
		return "", err
	}

	// VirusTotal'dan gelen analiz sonuçlarını yorumla.
	stats := vtResp.Data.Attributes.LastAnalysisStats
	if stats.Malicious > 0 {
		return "malicious", nil
	} else if stats.Suspicious > 0 {
		return "suspicious", nil
	} else {
		return "clean", nil
	}
}
