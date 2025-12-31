package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

const (
	// Tor proxy adresleri (varsayılan: 9050, Tor Browser: 9150)
	torProxyAddr = "127.0.0.1:9150"

	// Timeout değerleri
	httpTimeout = 30 * time.Second

	// Çıktı dizini
	outputDir = "scraped_data"

	// Log dizini
	logDir = "logs"

	// Log dosyası
	logFile = "logs/scan_report.log"
)

// ANSI renk kodları (Windows 10+ destekler)
const (
	colorReset  = "\033[0m"
	colorPurple = "\033[95m"
	colorCyan   = "\033[96m"
	colorGreen  = "\033[92m"
	colorYellow = "\033[93m"
	colorRed    = "\033[91m"
	colorWhite  = "\033[97m"
)

// Forum yapılandırması
type Forum struct {
	Name string
	URLs []string
}

// Forum listesi
var forums = []Forum{
	{Name: "GhostHub", URLs: []string{"http://aniozgjggq2pzxznogrlpoioks7iu3emj6bwebz3yptl4pkoukzd6kid.onion"}},
	{Name: "Darkzone", URLs: []string{"http://cashvosowt4iblo6levginhhquzi5iot5wedgtfyxxs5bkyk33roybqd.onion"}},
	{Name: "DeepWeb Question and Answers", URLs: []string{"http://fullzcxfok5643yy667mn75cqlcazzvr6oytdziffhku66k7wbnghxqd.onion"}},
	{Name: "Respostas Ocultas", URLs: []string{"http://exiliokpa5hknpknevkccokdrkmw2nibcrl3ehjwkr2nj55ff5xfuwad.onion"}},
	{Name: "Out3r Space", URLs: []string{"http://forumup4mrybxqweysxevr3nfivmphrwq5tngxbhixwppzji34y454id.onion"}},
	{Name: "The Tor Forum", URLs: []string{"http://forumup4mrybxqweysxevr3nfivmphrwq5tngxbhixwppzji34y454id.onion"}},
	{Name: "DarkWeb Forums", URLs: []string{"http://dreamr6y7e7ilemgws2ccuegvwdge3sfdm2ijh4s4hc2eggz3cvgl2ad.onion"}},
	{Name: "Suprbay", URLs: []string{"http://darkzqtmbdeauwq5mzcmgeeuhet42fhfjj4p5wbak3ofx2yqgecoeqyd.onion"}},
	{Name: "Hidden Answers", URLs: []string{"http://bags3opronrsfn4umq5w7zaaqillxfq4ml3vrqyb47vid6ujbywarqyd.onion"}},
	{Name: "FrenchPool", URLs: []string{"http://24ov5b5kawzsevv2ayltf52shlutzsewcjligevx35izqahqivdujryd.onion"}},
	{Name: "Wall of Shame", URLs: []string{"http://apwalloxvnx4jzy2q5lomljv2jfoxtlsv2acafrklmi2hy36yi32stad.onion/"}},
}

// Log tipi
type logLevel string

const (
	INFO logLevel = "INFO"
	ERR  logLevel = "ERR"
	WARN logLevel = "WARN"
)

// Logger yapısı
type Logger struct {
	logFile *os.File
}

// Yeni logger oluştur
func NewLogger() (*Logger, error) {
	// Log dizinini oluştur
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("log dizini oluşturulamadı: %v", err)
	}

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("log dosyası oluşturulamadı: %v", err)
	}
	return &Logger{logFile: file}, nil
}

// Log yaz
func (l *Logger) Log(level logLevel, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message)

	// Konsola yazdır (renksiz)
	fmt.Print(logMsg)

	// Dosyaya yaz
	if l.logFile != nil {
		l.logFile.WriteString(logMsg)
	}
}

// Logger'ı kapat
func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// Renkli çıktı fonksiyonları
func printColor(color, text string) {
	fmt.Print(color + text + colorReset)
}

func printInfo(text string) {
	printColor(colorCyan, "i "+text+"\n")
}

func printSuccess(text string) {
	printColor(colorGreen, "✓ "+text+"\n")
}

func printError(text string) {
	printColor(colorRed, "✗ "+text+"\n")
}

// Menüyü göster
func showMenu() {
	fmt.Println()
	printColor(colorPurple, "=== Dark Web Forum Scraper ===\n")
	fmt.Println()

	for i, forum := range forums {
		fmt.Printf("  %d.  %s\n", i+1, forum.Name)
	}

	fmt.Printf("  12. Scrape all forums\n")
	fmt.Printf("  0.  Exit\n")
	fmt.Println()
}

// Tor proxy üzerinden HTTP client oluştur
func createTorClient() (*http.Client, error) {
	// SOCKS5 proxy dialer oluştur
	dialer, err := proxy.SOCKS5("tcp", torProxyAddr, nil, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("SOCKS5 proxy oluşturulamadı: %v", err)
	}

	// HTTP transport yapılandırması
	transport := &http.Transport{
		Dial: dialer.Dial,
	}

	// HTTP client oluştur
	client := &http.Client{
		Transport: transport,
		Timeout:   httpTimeout,
	}

	return client, nil
}

// URL'den veri çek
func scrapeURL(client *http.Client, url string, logger *Logger) error {
	logger.Log(INFO, fmt.Sprintf("Taranıyor: %s", url))

	// HTTP isteği oluştur
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Log(ERR, fmt.Sprintf("İstek oluşturulamadı (%s): %v", url, err))
		return err
	}

	// User-Agent header'ı ekle (OpSec için)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")

	// İsteği gönder
	resp, err := client.Do(req)
	if err != nil {
		logger.Log(ERR, fmt.Sprintf("İstek başarısız (%s): %v", url, err))
		return err
	}
	defer resp.Body.Close()

	// HTTP durum kodunu kontrol et
	if resp.StatusCode != http.StatusOK {
		logger.Log(ERR, fmt.Sprintf("HTTP hatası (%s): %d %s", url, resp.StatusCode, resp.Status))
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// Response body'yi oku
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log(ERR, fmt.Sprintf("Response okunamadı (%s): %v", url, err))
		return err
	}

	// Dosya adı oluştur (URL'den güvenli bir dosya adı)
	safeFilename := sanitizeFilename(url)
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.html", safeFilename, timestamp)

	// Çıktı dizinini oluştur
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		logger.Log(ERR, fmt.Sprintf("Çıktı dizini oluşturulamadı: %v", err))
		return err
	}

	// Dosyaya yaz
	filePath := filepath.Join(outputDir, filename)
	if err := os.WriteFile(filePath, body, 0644); err != nil {
		logger.Log(ERR, fmt.Sprintf("Dosya yazılamadı (%s): %v", filePath, err))
		return err
	}

	logger.Log(INFO, fmt.Sprintf("Başarılı: %s -> %s (%d bytes)", url, filePath, len(body)))
	return nil
}

// Dosya adı için güvenli string oluştur
func sanitizeFilename(url string) string {
	// URL'deki özel karakterleri temizle
	safe := strings.ReplaceAll(url, "://", "_")
	safe = strings.ReplaceAll(safe, "/", "_")
	safe = strings.ReplaceAll(safe, ":", "_")
	safe = strings.ReplaceAll(safe, ".", "_")
	safe = strings.ReplaceAll(safe, "?", "_")
	safe = strings.ReplaceAll(safe, "&", "_")
	safe = strings.ReplaceAll(safe, "=", "_")

	// Maksimum uzunluk sınırla
	if len(safe) > 100 {
		safe = safe[:100]
	}

	return safe
}

// Forum tarama işlemi
func scrapeForum(client *http.Client, forum Forum, logger *Logger) {
	if len(forum.URLs) == 0 {
		printInfo(fmt.Sprintf("%s için URL bulunamadı (boş liste)", forum.Name))
		logger.Log(WARN, fmt.Sprintf("%s için URL listesi boş", forum.Name))
		return
	}

	printInfo(fmt.Sprintf("%s taraması başlatılıyor (%d URL)...", forum.Name, len(forum.URLs)))
	logger.Log(INFO, fmt.Sprintf("=== %s Taraması Başlatıldı ===", forum.Name))

	var successCount, failCount int

	for i, url := range forum.URLs {
		logger.Log(INFO, fmt.Sprintf("İlerleme [%s]: %d/%d", forum.Name, i+1, len(forum.URLs)))

		// Hata toleransı: Bir URL başarısız olsa bile devam et
		if err := scrapeURL(client, url, logger); err != nil {
			failCount++
			logger.Log(ERR, fmt.Sprintf("Tarama başarısız [%s] (%s): %v", forum.Name, url, err))
			continue
		}
		successCount++

		// Rate limiting
		time.Sleep(1 * time.Second)
	}

	logger.Log(INFO, fmt.Sprintf("=== %s Taraması Tamamlandı ===", forum.Name))
	logger.Log(INFO, fmt.Sprintf("%s - Başarılı: %d, Başarısız: %d", forum.Name, successCount, failCount))
	printSuccess(fmt.Sprintf("%s tamamlandı (Başarılı: %d, Başarısız: %d)", forum.Name, successCount, failCount))
}

// Tor bağlantısını test et
func testTorConnection(client *http.Client) bool {
	testURL := "https://check.torproject.org"

	req, err := http.NewRequest("GET", testURL, nil)
	if err != nil {
		printError("Tor bağlantı testi oluşturulamadı")
		return false
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		printError(fmt.Sprintf("Tor bağlantı testi başarısız: %v", err))
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	// Basit bir kontrol
	if strings.Contains(string(body), "Congratulations") || strings.Contains(string(body), "torproject") {
		return true
	}

	return false
}

func main() {
	// Log dizinini oluştur
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Log dizini oluşturulamadı: %v\n", err)
		os.Exit(1)
	}

	// Logger oluştur
	logger, err := NewLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Logger oluşturulamadı: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	// Tor bağlantı kontrolü
	printInfo(fmt.Sprintf("Tor bağlantısı kontrol ediliyor (Port: %s)...", strings.Split(torProxyAddr, ":")[1]))
	client, err := createTorClient()
	if err != nil {
		printError(fmt.Sprintf("Tor client oluşturulamadı: %v", err))
		printError("Tor servisinin çalıştığından emin olun")
		logger.Log(ERR, fmt.Sprintf("Tor client oluşturulamadı: %v", err))
		os.Exit(1)
	}

	// Tor bağlantı testi
	if testTorConnection(client) {
		printSuccess("Tor bağlantısı başarılı!")
		logger.Log(INFO, "Tor bağlantısı doğrulandı")
	} else {
		printError("Tor bağlantı testi başarısız, ancak devam ediliyor...")
		logger.Log(WARN, "Tor bağlantı testi belirsiz")
	}

	// Scan raporu oluşturuluyor mesajı
	printInfo("Scan raporu oluşturuluyor...")
	printSuccess(fmt.Sprintf("Scan raporu oluşturuldu: %s", logFile))
	logger.Log(INFO, "=== Dark Web Forum Scraper Başlatıldı ===")

	// Ana döngü
	for {
		showMenu()

		fmt.Print("Select an option: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			printError("Giriş okunamadı")
			continue
		}

		// Satır sonu karakterlerini temizle
		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil {
			printError("Geçersiz seçim")
			continue
		}

		// Seçim işleme
		if choice == 0 {
			printInfo("Çıkılıyor...")
			logger.Log(INFO, "Program sonlandırıldı")
			break
		} else if choice == 12 {
			// Tüm forumları tara
			printInfo("Tüm forumlar taranacak...")
			logger.Log(INFO, "=== Tüm Forumlar Taraması Başlatıldı ===")

			for _, forum := range forums {
				scrapeForum(client, forum, logger)
				time.Sleep(2 * time.Second) // Forumlar arası bekleme
			}

			logger.Log(INFO, "=== Tüm Forumlar Taraması Tamamlandı ===")
			printSuccess("Tüm forumlar taraması tamamlandı!")
		} else if choice >= 1 && choice <= len(forums) {
			// Seçilen forumu tara
			forum := forums[choice-1]
			scrapeForum(client, forum, logger)
		} else {
			printError("Geçersiz seçim")
		}

		fmt.Println()
		fmt.Print("Devam etmek için Enter'a basın...")
		reader.ReadString('\n')
	}
}
