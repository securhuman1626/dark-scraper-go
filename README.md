# Dark Web Forum Scraper

Tor ağı üzerinden .onion forumlarını otomatik olarak tarayan, menü tabanlı Go uygulaması. Interaktif menü ile istediğiniz forumu seçerek tarama yapabilirsiniz.

## 🎯 Özellikler

- ✅ **Menü Tabanlı Arayüz**: Interaktif menü ile forum seçimi
- ✅ **Tor SOCKS5 Proxy Desteği**: 127.0.0.1:9150 (Tor Browser) veya 9050 (Standart Tor)
- ✅ **11 Farklı Forum**: Hazır yapılandırılmış forum listesi
- ✅ **Hata Toleransı**: Bir site başarısız olsa bile program devam eder
- ✅ **Renkli Konsol Çıktıları**: Windows terminal için renkli mesajlar
- ✅ **Otomatik HTML Kaydetme**: Tüm veriler `scraped_data/` dizinine kaydedilir
- ✅ **Detaylı Loglama**: Tüm işlemler `logs/scan_report.log` dosyasına kaydedilir
- ✅ **User-Agent ve OpSec**: Tarayıcı gibi görünecek şekilde header'lar
- ✅ **Tor Bağlantı Testi**: Başlangıçta otomatik Tor bağlantı kontrolü

## 📋 Gereksinimler

1. **Go (Golang)**: v1.18 veya üzeri
   - İndirme: https://golang.org/dl/

2. **Tor Browser veya Tor Servisi**: 
   - **Tor Browser** (Önerilen): Port 9150 kullanır
   - **Standart Tor Servisi**: Port 9050 kullanır

## 🚀 Kurulum

1. **Projeyi klonlayın veya indirin**

2. **Bağımlılıkları yükleyin:**
```bash
go mod download
```

3. **Tor Browser'ı başlatın** (Tor Browser kullanıyorsanız)
   - Tor Browser'ı açın ve bağlantının hazır olmasını bekleyin

## 💻 Kullanım

### Hızlı Başlangıç

1. **Programı çalıştırın:**
```bash
go run main.go
```

Veya derlenmiş binary ile:
```bash
# Derle
go build -o scraper.exe

# Çalıştır
./scraper.exe
```

2. **Menüden seçim yapın:**
   - Program başladığında Tor bağlantısını kontrol eder
   - Menü otomatik olarak gösterilir
   - İstediğiniz forumun numarasını girin

### Menü Seçenekleri

Program başlatıldığında aşağıdaki menü görüntülenir:

```
=== Dark Web Forum Scraper ===

  1.  GhostHub
  2.  Darkzone
  3.  DeepWeb Question and Answers
  4.  Respostas Ocultas
  5.  Out3r Space
  6.  The Tor Forum
  7.  DarkWeb Forums
  8.  Suprbay
  9.  Hidden Answers
  10. FrenchPool
  11. Wall of Shame
  12. Scrape all forums
  0.  Exit

Select an option:
```

**Seçimler:**
- **1-11**: Belirli bir forumu taramak için
- **12**: Tüm forumları sırayla taramak için
- **0**: Programdan çıkmak için

### Örnek Kullanım Senaryoları

#### Senaryo 1: Tek bir forumu taramak
```
Select an option: 1
```
GhostHub forumu taranacaktır.

#### Senaryo 2: Tüm forumları taramak
```
Select an option: 12
```
Tüm 11 forum sırayla taranacaktır.

#### Senaryo 3: Programdan çıkmak
```
Select an option: 0
```

## 📁 Çıktılar ve Dosyalar

### Çıktı Dizinleri

- **`scraped_data/`**: Taranan tüm HTML dosyaları burada saklanır
  - Dosya adı formatı: `{URL}_{timestamp}.html`
  - Örnek: `http___forum_onion_20250101_143022.html`

- **`logs/`**: Log dosyaları burada saklanır
  - **`scan_report.log`**: Detaylı tarama raporu
    - Başarılı/başarısız tüm işlemler
    - Timestamp'li log kayıtları
    - Hata mesajları ve istatistikler

### Log Formatı

```
[2025-01-01 14:30:22] [INFO] === GhostHub Taraması Başlatıldı ===
[2025-01-01 14:30:23] [INFO] Taranıyor: http://example.onion
[2025-01-01 14:30:25] [INFO] Başarılı: http://example.onion -> scraped_data/...html (15234 bytes)
[2025-01-01 14:30:26] [INFO] === GhostHub Taraması Tamamlandı ===
```

## ⚙️ Yapılandırma

### Tor Portu Değiştirme

`main.go` dosyasında `torProxyAddr` değişkenini düzenleyin:

```go
const (
    // Tor Browser için (varsayılan)
    torProxyAddr = "127.0.0.1:9150"
    
    // Standart Tor servisi için
    // torProxyAddr = "127.0.0.1:9050"
)
```

### Timeout Süresini Ayarlama

```go
const (
    httpTimeout = 30 * time.Second  // İstek timeout süresi
)
```

### Forum Linklerini Düzenleme

`main.go` dosyasında `forums` dizisini düzenleyin:

```go
var forums = []Forum{
    {Name: "GhostHub", URLs: []string{
        "http://example1.onion",
        "http://example2.onion",
        // Daha fazla URL ekleyebilirsiniz
    }},
    // ...
}
```

**Not**: Her forum için birden fazla URL ekleyebilirsiniz.

## 🔒 Güvenlik Notları

⚠️ **ÖNEMLİ**: Bu araç sadece **eğitim** ve **yasal CTI (Cyber Threat Intelligence)** amaçları için kullanılmalıdır.

### Güvenlik Önerileri

- ✅ Tor Browser veya Tor servisinin çalıştığından emin olun
- ✅ IP sızıntısı riskini azaltmak için VPN kullanmayı düşünün
- ✅ User-Agent ve header'lar gerçekçi görünecek şekilde ayarlanmıştır
- ✅ Tüm trafik SOCKS5 proxy üzerinden yönlendirilir
- ❌ Yasadışı aktiviteler için kullanmayın
- ❌ Kişisel verileri toplamayın veya kötüye kullanmayın

## 🐛 Sorun Giderme

### "Tor client oluşturulamadı" Hatası

**Çözüm:**
1. Tor Browser'ın açık ve bağlantının hazır olduğundan emin olun
2. Port numarasını kontrol edin:
   - Tor Browser: 9150
   - Standart Tor: 9050
3. `main.go` dosyasındaki `torProxyAddr` değerini kontrol edin

### "Tor bağlantı testi başarısız" Uyarısı

**Çözüm:**
- Bu uyarı görünse bile program çalışmaya devam edebilir
- Tor servisinin gerçekten çalıştığından emin olun
- İnternet bağlantınızı kontrol edin

### "Connection refused" Hatası

**Çözüm:**
- Tor servisi başlatılmamış olabilir
- Firewall Tor portunu engelliyor olabilir
- Windows Firewall'da portları kontrol edin

### Menü Görünmüyor / Renkler Çalışmıyor

**Çözüm:**
- Windows 10+ kullanıyorsanız ANSI renk desteği otomatik aktif olmalı
- Eski Windows sürümlerinde renkler görünmeyebilir ama program çalışır
- Terminal'iniz UTF-8 karakterleri desteklemelidir

### Dosyalar Kaydedilmiyor

**Çözüm:**
- `scraped_data/` dizini için yazma izinlerini kontrol edin
- Disk alanını kontrol edin
- Antivirus yazılımı dosyaları engelliyor olabilir

## 📊 İstatistikler ve Raporlama

Her tarama sonunda:

- ✅ Başarılı tarama sayısı
- ❌ Başarısız tarama sayısı
- 📁 Kaydedilen dosya sayısı
- 📝 Detaylı log kayıtları

Bu bilgiler hem konsola hem de `logs/scan_report.log` dosyasına yazılır.

## 🔄 Program Akışı

1. **Başlangıç**: Tor bağlantısı kontrol edilir
2. **Menü Gösterimi**: Forum listesi görüntülenir
3. **Seçim**: Kullanıcı bir seçenek girer
4. **Tarama**: Seçilen forum(lar) taranır
5. **Kayıt**: HTML dosyaları kaydedilir
6. **Rapor**: Sonuçlar loglanır ve gösterilir
7. **Döngü**: Menü tekrar gösterilir (Exit seçilene kadar)

## 📝 Notlar

- Her URL taraması arasında 1 saniye bekleme süresi vardır (rate limiting)
- Forumlar arası bekleme süresi 2 saniyedir
- HTTP 200 OK dışındaki yanıtlar hata olarak işaretlenir
- Program hata olsa bile çökmez, loglar ve devam eder

## 📄 Lisans

Bu proje eğitim amaçlıdır. Yasal sorumluluk kullanıcıya aittir.

## 🤝 Katkıda Bulunma

1. Forum linklerini güncelleyebilirsiniz
2. Yeni forumlar ekleyebilirsiniz
3. Hata raporları için issue açabilirsiniz

---

**İyi taramalar! 🚀**
