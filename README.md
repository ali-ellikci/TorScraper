# TOR Scraper - Otomatize Tor Ağı Tarama Aracı

## 📋 Proje Amacı

Siber tehdit aktörleri izlerini kaybettirmek için Tor ağını kullanmaktadır. Tekil analizler manuel yapılabilse de, yüzlerce .onion adresini (sızıntı siteleri, forumlar, marketler) düzenli olarak taramak insan gücüyle imkansızdır.

Bu proje; **Go (Golang) dilini kullanarak**, toplu hedef listesini (YAML) işleyebilen, trafiği Tor ağı üzerinden anonim olarak yönlendiren ve elde edilen istihbaratı raporlayan bir otomasyon aracı sunar.

### Hedefler
- ✅ CTI süreçlerindeki **Collection** (Toplama) yetkinliği kazanma
- ✅ **Automation** (Otomasyon) yetkinliği kazanma
- ✅ Go'nun Goroutine'leri ile paralel işleme performansı

---

## 🏗️ Proje Mimarisi

Proje 4 ana modülden oluşmaktadır:

### 1. **Dosya Okuma Modülü (Input Handler)** - `internal/input/`
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

## ⚡ Performans Özellikleri

### Goroutines ile Hızlandırma (İsteğe Bağlı)
Projenin basit sürümü sırayla tarama yaparken, ileri kullanıcılar **Goroutines** kullanarak paralelleştirme yapabilir:

```go
// Sırayla tarama (Temel)
for _, target := range targets {
    scanTarget(target)
}

// Paralel tarama (İleri - Goroutines)
for _, target := range targets {
    go scanTarget(target)
}
```

**Not:** Goroutines kullanırken rate limiting ve bağlantı yönetimi kritiktir!

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

## 👨‍💻 Geliştirici Notları

> "Go dili, modern bulut ve ağ araçlarının dilidir. Bu projede Python yerine Go kullanmamızın sebebi, ileride binlerce siteyi aynı anda taramak istediğinizde Go'nun 'Goroutines' yapısının size sağlayacağı performansı şimdiden hissetmenizdir. Bu ödevde basit bir döngü kullanabilirsiniz, ancak meraklıları 'goroutine' ile taramayı hızlandırmayı deneyebilir!"

---

**Proje Tamamlanma Tarihi:** Aralık 2025
