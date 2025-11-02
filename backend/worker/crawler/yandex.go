package crawler

import (
	"compress/flate"
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"
	"time"

	"github.com/andybalholm/brotli"
)

// YaMarketClient клиент
type YaMarketClient struct {
	httpClient *http.Client
	userAgents []string
	referer    string
}

// NewYaMarketClient создаёт новый экземпляр клиента
func NewYaMarketClient() *YaMarketClient {
	jar, _ := cookiejar.New(nil)

	client := &YaMarketClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     jar, // Важно: сохраняем cookies между запросами
			Transport: &http.Transport{
				MaxIdleConns:        10,
				MaxIdleConnsPerHost: 5,
				IdleConnTimeout:     60 * time.Second,
				TLSClientConfig: &tls.Config{
					MinVersion: tls.VersionTLS12,
				},
				DisableCompression: false,
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return nil // Следуем редиректам
			},
		},
		userAgents: []string{
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_7_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.1 Safari/605.1.15",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:131.0) Gecko/20100101 Firefox/131.0",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
		},
		referer: "https://yandex.ru/",
	}

	return client
}

// randomDelay добавляет случайную задержку
func (c *YaMarketClient) randomDelay() {
	delay := time.Duration(1000+rand.Intn(2000)) * time.Millisecond // 1-3 секунды
	time.Sleep(delay)
}

// getRandomUserAgent возвращает случайный User-Agent
func (c *YaMarketClient) getRandomUserAgent() string {
	return c.userAgents[rand.Intn(len(c.userAgents))]
}

// prepareRequest подготавливает HTTP запрос с необходимыми заголовками
func (c *YaMarketClient) prepareRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Заголовки, имитирующие браузер
	req.Header.Set("User-Agent", c.getRandomUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("DNT", "1")
	req.Header.Set("sec-ch-ua", `"Chromium";v="130", "Google Chrome";v="130", "Not?A_Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)

	// Добавляем referer для последующих запросов
	if c.referer != "" {
		req.Header.Set("Referer", c.referer)
		req.Header.Set("Sec-Fetch-Site", "same-origin")
	}

	return req, nil
}

// FetchProductByURL извлекает информацию о товаре по прямой ссылке
func (c *YaMarketClient) FetchProductByURL(productURL string) (*ProductInfo, error) {
	if !strings.Contains(productURL, "market.yandex.ru") && !strings.Contains(productURL, "yandex.ru/products") {
		return nil, fmt.Errorf("invalid Yandex Market URL")
	}

	// Случайная задержка перед запросом
	c.randomDelay()

	req, err := c.prepareRequest("GET", productURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product page: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("Error closing response body: %v", cerr)
		}
	}()

	// Сохраняем URL для последующих запросов
	c.referer = productURL

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Unexpected status code %d for URL %s, response: %s", resp.StatusCode, productURL, string(body[:min(500, len(body))]))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var reader io.Reader = resp.Body
	// Проверяем кодировку ответа
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzReader.Close()
		reader = gzReader
	case "br":
		reader = brotli.NewReader(resp.Body)
	case "deflate":
		reader = flate.NewReader(resp.Body)
	}
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return c.parseProductPage(string(body), productURL)
}

// ProductInfo содержит информацию о товаре
type ProductInfo struct {
	Title       string
	Price       string
	ImageURL    string
	Description string
	URL         string
}

// parseProductPage парсит страницу товара и извлекает информацию из микроразметки
func (c *YaMarketClient) parseProductPage(html, productURL string) (*ProductInfo, error) {
	product := &ProductInfo{
		URL: productURL,
	}

	// Ищем комментарий с путём содержащим "productPageMicromarkup" и следующий за ним script
	blockPattern := regexp.MustCompile(`<!--BEGIN \[@marketfront/JsonLd\] [^>]*productPageMicromarkup[^>]*--><script type="application/ld\+json">(.*?)</script>`)
	blockMatch := blockPattern.FindStringSubmatch(html)

	if len(blockMatch) > 1 {
		jsonData := blockMatch[1]

		// Извлекаем name
		namePattern := regexp.MustCompile(`"name"\s*:\s*"([^"]*)"`)
		if match := namePattern.FindStringSubmatch(jsonData); len(match) > 1 {
			product.Title = match[1]
		}

		// Извлекаем image
		imagePattern := regexp.MustCompile(`"image"\s*:\s*"([^"]*)"`)
		if match := imagePattern.FindStringSubmatch(jsonData); len(match) > 1 {
			product.ImageURL = match[1]
		}

		// Извлекаем offers.price
		pricePattern := regexp.MustCompile(`"offers"\s*:\s*{[^}]*"price"\s*:\s*"([^"]*)"`)
		if match := pricePattern.FindStringSubmatch(jsonData); len(match) > 1 {
			product.Price = match[1] + " ₽"
		}

		// Извлекаем прямой url
		urlPattern := regexp.MustCompile(`"url"\s*:\s*"([^"]*)"`)
		if match := urlPattern.FindStringSubmatch(jsonData); len(match) > 1 {
			product.URL = match[1]
		}

		// Извлекаем description (опционально)
		//descPattern := regexp.MustCompile(`"description"\s*:\s*"([^"]*)"`)
		//if match := descPattern.FindStringSubmatch(jsonData); len(match) > 1 {
		//	product.Description = match[1]
		//}
	}

	// Проверяем, что название извлечено
	if product.Title == "" {
		log.Printf("No product found in %v", html)
		return nil, fmt.Errorf("failed to extract product data from micromarkup")
	}

	return product, nil
}
