# TOR Scraper - Otomatize Tor Ağı Tarama Aracı

## 📋 Proje Amacı

Siber tehdit aktörleri izlerini kaybettirmek için Tor ağını kullanmaktadır. Tekil analizler manuel yapılabilse de, yüzlerce .onion adresini (sızıntı siteleri, forumlar, marketler) düzenli olarak taramak insan gücüyle imkansızdır.

Bu proje; **Go (Golang) dilini kullanarak**, toplu hedef listesini (YAML) işleyebilen, trafiği Tor ağı üzerinden anonim olarak yönlendiren ve elde edilen istihbaratı raporlayan bir otomasyon aracı sunar.

### Hedefler
- ✅ CTI süreçlerindeki **Collection** (Toplama) yetkinliği kazanma
- ✅ **Automation** (Otomasyon) yetkinliği kazanma
- ✅ Go'nun Goroutine'leri ile paralel işleme performansı
- ✅ **IP sızıntısını önlemek** için özel HTTP Transport/Client kullanımı
- ✅ Kapsamlı loglama ve JSON raporlama

---

## 🏗️ Proje Mimarisi

Proje 4 ana modülden oluşmaktadır:

### 1. **Dosya Okuma Modülü (Input Handler)** - `internal/input/`
- YAML formatında .onion adreslerini okur
- Her URL'i temizler (whitespace trimming)
- Hedef listesini döndürür

### 2. **Tor Proxy Yönetimi (Go Proxy Client)** - `internal/tor/`
- `net/http` kütüphanesini SOCKS5 proxy'sine yönlendirir (127.0.0.1:9050)
- **IP sızıntısını önlemek** için özel `http.Transport` ve `http.Client` yapılandırması
- TOR ağı üzerinden anonim bağlantı sağlar

### 3. **Tarama Modülü (Scanner)** - `internal/scanner/`
- **HTTP Client**: SOCKS5 proxy üzerinden HTML içeriğini çeker
- **IP Verification**: check.torproject.org adresi kontrol ederek TOR IP'sini doğrular
- **Chromedp**: TOR proxy üzerinden tarayıcı ile ekran görüntüsü alır
- Hata yönetimi: Dead site'ler programı durdurmaz, loglayıp devam eder

### 4. **Veri Kayıt Modülü (Output Writer)** - `internal/output/`
- **HTML Dosyaları**: `output/html/` dizinine URL adı + timestamp ile kaydedilir
- **Screenshot'lar**: `output/screenshots/` dizinine PNG formatında kaydedilir
- **JSON Rapor**: `output/scan_report_*.json` - Detaylı tarama sonuçları
- **Log Dosyası**: `output/scan_report_*.log` - Tüm işlemlerin kaydı

### 5. **Logger Modülü** - `internal/logger/`
- Console ve dosyaya eş zamanlı loglama
- Timestamp ile her log kaydı
- INFO, ERROR, SUCCESS, WARN seviyeleri

---

## 💻 Kurulum ve Çalıştırma

### Ön Koşullar
1. **Go** (1.19+) yüklü olmalı
2. **Tor Service** çalışır durumda olmalı (SOCKS5: 127.0.0.1:9050)
3. **Chromium/Chrome** tarayıcı yüklü olmalı (Chromedp için)

### Windows'ta TOR Kurulumu
```bash
# Option 1: Tor Browser kullanın (en kolay)
# https://www.torproject.org/download/

# Option 2: Tor Service kurulum
# https://www.torproject.org/download/#windows
```

### Go Projesini Kurma
```bash
# Repository'i clone et
git clone https://github.com/ali-ellikci/TorScraper.git
cd TorScraper

# Bağımlılıkları indir
go mod download
go mod tidy

# Projeyi çalıştır
go run .\cmd\tor-scraper\main.go

# Veya derle
go build -o TorScraper.exe .\cmd\tor-scraper\main.go
.\TorScraper.exe
```

---

## 📊 Çıktılar

### Başarılı çalıştırma sonrasında oluşan dosyalar:

```
output/
├── screenshots/          # PNG ekran görüntüleri
│   ├── bestteermb42clir_20251229_075511.png
│   ├── dreadytofatropt_20251229_075512.png
│   └── ...
├── html/                 # HTML dosyaları
│   ├── bestteermb42clir_20251229_075511.html
│   ├── dreadytofatropt_20251229_075512.html
│   └── ...
├── scan_report_20251229_075511.log    # Log dosyası
└── scan_report_20251229_075511.json   # JSON rapor
```

### JSON Rapor Formatı
```json
{
  "start_time": "2025-12-29T07:55:07.123456Z",
  "end_time": "2025-12-29T07:55:20.654321Z",
  "total_targets": 9,
  "success_count": 1,
  "fail_count": 8,
  "records": [
    {
      "url": "https://www.google.com/",
      "status": "SUCCESS",
      "status_code": 200,
      "ip_address": "{\"ip\":\"1.2.3.4\",\"is_tor\":true}",
      "timestamp": "2025-12-29T07:55:15.123456Z",
      "screenshot_path": "output/screenshots/www.google.com_20251229_075515.png",
      "html_path": "output/html/www.google.com_20251229_075515.html"
    },
    {
      "url": "http://bestteermb42clir6ux7xm76d4jjodh3fpahjqgbddbmfrgp4skg2wqd.onion/",
      "status": "FAILED",
      "timestamp": "2025-12-29T07:55:08.654321Z",
      "error": "[FAILED] failed to scan target: page load error net::ERR_NAME_NOT_RESOLVED"
    }
  ]
}
```

---

## 🔧 Yapı ve Dosya Tasnifi

```
TorScraper/
├── cmd/
│   └── tor-scraper/
│       └── main.go              # Ana program giriş noktası
├── internal/
│   ├── input/
│   │   └── reader.go            # YAML dosya okuma
│   ├── logger/
│   │   └── logger.go            # Loglama sistemi
│   ├── output/
│   │   ├── writer.go            # HTML/Screenshot kaydetme
│   │   ├── report.go            # JSON rapor oluşturma
│   │   └── screen_report.go     # Screen rapor (opsiyonel)
│   ├── scanner/
│   │   └── scanner.go           # Tarama motoru (HTTP + Chromedp)
│   └── tor/
│       └── client.go            # TOR SOCKS5 client yapılandırması
├── configs/
│   └── targets.yaml             # Taranacak .onion adresleri
├── output/                      # Çıktı dosyaları (otomatik oluşturulur)
├── go.mod                       # Go modülü tanımı
├── go.sum                       # Bağımlılık haritası
└── README.md                    # Bu dosya
```

---

## 🔐 Güvenlik Özellikleri

### 1. **IP Sızıntısı Koruma**
- Özel `http.Transport` yapılandırması
- SOCKS5 proxy aracılığıyla tüm trafiğin yönlendirilmesi
- IP verification: `check.torproject.org` kontrolü

### 2. **Hata Yönetimi**
- Dead site'ler (ERR_NAME_NOT_RESOLVED) programı durdurmaz
- Timeout yönetimi (30 saniye)
- Her hata kaydedilir ve rapora eklenir

### 3. **Veri Güvenliği**
- Tüm çıktılar `output/` dizininde merkezi yönetim
- JSON rapor ile yapılandırılmış veri depolama
- Timestamp ile dosya çakışmalarını önleme

---

## 📝 Kullanım Örneği

### 1. Hedef Dosyası Hazırlama
`configs/targets.yaml`:
```yaml
http://bestteermb42clir6ux7xm76d4jjodh3fpahjqgbddbmfrgp4skg2wqd.onion/
https://dreadytofatroptsdj6io7l3xptbet6onoyno2yv7jicoxknyazubrad.onion/
https://www.google.com/
```

### 2. TOR Servisini Başlat
```bash
# Windows: Tor Browser'ı çalıştırın
# veya Tor Service kurulu ise:
# Net Start Tor
```

### 3. Programı Çalıştır
```bash
go run .\cmd\tor-scraper\main.go
```

### 4. Sonuçları İnceле
```bash
# Log dosyasını oku
type output\scan_report_*.log

# JSON raporunu oku
type output\scan_report_*.json

# Screenshot'ları görüntüle
dir output\screenshots\

# HTML dosyalarını kontrol et
dir output\html\
```

---

## 🎯 Beklenen Terminal Çıktısı

```
2025/12/29 07:55:07 [INFO] Starting TOR Scraper with 9 targets
2025/12/29 07:55:07 [INFO] [1/9] Scanning: http://bestteermb42clir6ux7xm76d4jjodh3fpahjqgbddbmfrgp4skg2wqd.onion/
2025/12/29 07:55:08 [ERR] http://bestteermb42clir6ux7xm76d4jjodh3fpahjqgbddbmfrgp4skg2wqd.onion/ -> [FAILED] failed to scan target...
2025/12/29 07:55:08 [INFO] [9/9] Scanning: https://www.google.com/
2025/12/29 07:55:15 [INFO] Using TOR IP: {"ip":"1.2.3.4","is_tor":true}
2025/12/29 07:55:15 [SUCCESS] Screenshot saved: output\screenshots\www.google.com_20251229_075515.png
2025/12/29 07:55:15 [SUCCESS] HTML saved: output\html\www.google.com_20251229_075515.html
2025/12/29 07:55:15 [SUCCESS] https://www.google.com/ (Status: 200, IP: {"ip":"1.2.3.4","is_tor":true})
2025/12/29 07:55:16 [INFO] Report saved: output/scan_report_20251229_075515.json
2025/12/29 07:55:16 [INFO] ========================================
2025/12/29 07:55:16 [INFO] Total: 9, Success: 1, Failed: 8
2025/12/29 07:55:16 [INFO] Screenshots: output/screenshots/
2025/12/29 07:55:16 [INFO] HTML files: output/html/
2025/12/29 07:55:16 [INFO] Log file: output/scan_report_*.log
2025/12/29 07:55:16 [INFO] JSON Report: output/scan_report_20251229_075515.json
```

---

## 🛠️ Gerekli Kütüphaneler

```go
require (
	github.com/chromedp/chromedp v0.14.2
	golang.org/x/net v0.48.0  // SOCKS5 proxy desteği
)
```

---


## 📄 Lisans

Bu proje eğitim amaçlı geliştirilmiştir.

---

## ✍️ Yazar

**Siber Tehdit İstihbaratı (CTI) Projesi**  
Eğitim Amacı: Tor Ağı Üzerinde Otomatize Veri Toplama

- Komut satırından YAML dosyası okunur
- Her satır temizlenir (whitespace trimming)
- .onion adresleri listelenip işleme hazırlanır

### 2. **Tor Proxy Yönetimi (Go Proxy Client)** - `internal/tor/`
- Go'nun `net/http` kütüphanesi, yerel Tor servisine (127.0.0.1:9050/9150) yönlendirilir
- Özel `http.Transport` ve `http.Client` yapılandırması ile IP sızıntısı önlenir
- SOCKS5 proxy üzerinden anonim istek gönderimi

### 3. **İstek ve Hata Yönetimi** - `internal/scanner/`
- Çalışmayan/kapanmış siteler programı durdurmaz
- Hatalar loglanır, tarama devam eder
- Timeout ve connection error yönetimi

### 4. **Veri Kayıt (Output Writer)** - `internal/output/`
- Başarılı isteklerden dönen HTML verisi kaydedilir
- URL adına veya tarih damgasına göre ayrı dosyalar
- Yapılandırılmış JSON/LOG formatında raporlama

---

## 📦 Kurulum

### Gereksinimler
- **Go 1.16+**
- **Tor Service** (arka planda çalışır durumda)
- Windows/Linux/macOS

### Adım 1: Tor Servisini Başlatın
```bash
# Windows
tor.exe

# Linux
sudo service tor start

# macOS
brew services start tor
```

SOCKS5 proxy'nin `127.0.0.1:9050` veya `127.0.0.1:9150` portunda çalıştığını doğrulayın.

### Adım 2: Proje Dosyalarını Hazırlayın

Hedef adresleri içeren `targets.yaml` dosyasını oluşturun:

```yaml
http://hss3d3eo7oxabjjx.onion
http://darkweblink.onion
http://example.onion
```

### Adım 3: Go Modüllerini İndirin
```bash
go mod download
```

---

## 🚀 Kullanım

### Temel Kullanım
```bash
go run ./cmd/tor-scraper/main.go
```

### Derlenmiş Binary ile Çalıştırma
```bash
# Derleme
go build -o tor-scraper.exe ./cmd/tor-scraper/main.go

# Çalıştırma
./tor-scraper.exe
```

### Örnek Komut Satırı Argümanları
```bash
go run ./cmd/tor-scraper/main.go -targets targets.yaml -output ./results/ -timeout 30
```

---

## 📊 Beklenen Çıktılar

Proje tamamlandığında aşağıdaki somut çıktılar elde edilecektir:

### 1. **Otomatize Tarama Aracı**
- Yüzlerce linki tek komutla tarayabilen derlenmiş Go binary dosyası

### 2. **Toplu Veri Seti**
```
output/
├── hss3d3eo7oxabjjx.onion_2025-12-28.html
├── darkweblink.onion_2025-12-28.html
└── example.onion_2025-12-28.json
```

### 3. **Durum Raporu**
```
[INFO] Scanner başlatıldı: 5 hedef bulundu
[INFO] Scanning: http://hss3d3eo7oxabjjx.onion -> SUCCESS (200)
[ERR] Scanning: http://deadsite.onion -> TIMEOUT (30s)
[INFO] Scanning: http://example.onion -> SUCCESS (200)
═══════════════════════════════════════
Tarama Tamamlandı:
- Başarılı: 2
- Başarısız: 1
- Toplam: 3
```

---

## 🛠️ Kullanılan Teknolojiler

### Programlama Dili
- **Go (Golang)** - Performans ve concurrency avantajları

### Kritik Kütüphaneler
```go
import (
    "net/http"                    // HTTP istekleri
    "golang.org/x/net/proxy"      // SOCKS5 proxy desteği
    "os"                          // Dosya işlemleri
    "bufio"                       // Dosya okuma/yazma
)
```

### Ağ Altyapısı
- **Tor Service** - Anonim ağ bağlantısı
- **SOCKS5 Proxy** - 127.0.0.1:9050 (varsayılan)

---

## 📁 Proje Yapısı

```
TorScraper/
├── README.md                    # Proje dokumentasyonu
├── cmd/
│   └── tor-scraper/
│       └── main.go              # Ana entry point
├── configs/                     # Konfigürasyon dosyaları
├── internal/
│   ├── input/
│   │   └── reader.go            # YAML dosyası okuma
│   ├── logger/                  # Loglama modülü
│   ├── output/                  # Veri kayıt modülü
│   ├── scanner/                 # Tarama motoru
│   └── tor/
│       └── client.go            # Tor proxy client
├── output/                      # Tarama sonuçları
└── targets.yaml                 # Hedef .onion adresleri
```



---

## 🔒 Güvenlik Notları

### IP Sızıntısı Önleme
- ✅ Tüm HTTP istekleri SOCKS5 proxy üzerinden yönlendirilir
- ✅ DNS leak'leri önlemek için özel transport yapılandırması yapılır
- ✅ User-Agent spoofing ile kimlik gizleme

### Tor Bağlantısı Doğrulama
```bash
# Tor üzerinden çalıştığını doğrulayın
curl --socks5 127.0.0.1:9050 https://check.torproject.org
```

---

## 📝 Örnek Çalıştırma Senaryosu

### 1. targets.yaml Hazırlama
```yaml
http://3g2upl4pq3khfchsl.onion
http://thehiddenwiki.onion
http://6nhmgdpnywnfwzqq.onion
```

### 2. Programı Çalıştırma
```bash
$ go run ./cmd/tor-scraper/main.go

[INFO] Tor Proxy bağlantısı kuruldu: 127.0.0.1:9050
[INFO] Scanner başlatıldı: 3 hedef bulundu
[INFO] Scanning: http://3g2upl4pq3khfchsl.onion -> SUCCESS (200)
[ERR] Scanning: http://thehiddenwiki.onion -> TIMEOUT
[INFO] Scanning: http://6nhmgdpnywnfwzqq.onion -> SUCCESS (200)
═════════════════════════════════════════════════════
Tarama Tamamlandı:
- Başarılı: 2
- Başarısız: 1
- Süre: 45 saniye
```

### 3. Sonuçları Kontrol Etme
```bash
$ ls -la output/
3g2upl4pq3khfchsl.onion_2025-12-28.html
6nhmgdpnywnfwzqq.onion_2025-12-28.html
scan_report.log
```

---

## 🐛 Hata Giderme

### "Connection refused" Hatası
```
❌ Hata: dial tcp 127.0.0.1:9050: connection refused
✅ Çözüm: Tor servisinin çalıştığını doğrulayın
```

### "Timeout" Hatası
```
❌ Hata: context deadline exceeded (timeout 30s)
✅ Çözüm: Timeout süresini artırın veya site gerçekten offline olabilir
```

### DNS Leak Şüphesi
```bash
$ curl --socks5 127.0.0.1:9050 https://check.torproject.org
```

---

## 📚 Kaynaklar

- [Go Official Documentation](https://golang.org/doc)
- [Tor Project - SOCKS Protocol](https://www.torproject.org)
- [golang.org/x/net/proxy](https://pkg.go.dev/golang.org/x/net/proxy)

---

## 📄 Lisans

Bu proje eğitim amaçlıdır.


---

**Proje Tamamlanma Tarihi:** Aralık 2025
