// internal/infrastructure/crawler/httputil.go
package crawler

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/andybalholm/brotli"
)

// productUserAgent — честное имя нашего бота. Используем там, где источник НЕ фильтрует
// по UA (WB-CDN). Где честный UA не отдаёт данные — представляемся иначе вынужденно:
// Ozon — WhatsApp-превью-UA (его allow-list), Яндекс — браузерный UA (боту микроразметку не отдаёт).
const productUserAgent = "WishCraft/TelegramBot"

// safeImageURL пропускает только http(s)-ссылки (защита от javascript:/data: в og:image);
// иначе возвращает пустую строку. Фронт рендерит результат строго как <img src>.
func safeImageURL(raw string) string {
	if strings.HasPrefix(raw, "https://") || strings.HasPrefix(raw, "http://") {
		return raw
	}
	return ""
}

// readBody читает тело ответа с декомпрессией по Content-Encoding.
func readBody(resp *http.Response) (string, error) {
	var r io.Reader = resp.Body
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			return "", err
		}
		defer func() { _ = gz.Close() }()
		r = gz
	case "br":
		r = brotli.NewReader(resp.Body)
	case "deflate":
		fr := flate.NewReader(resp.Body)
		defer func() { _ = fr.Close() }()
		r = fr
	}
	b, err := io.ReadAll(r)
	return string(b), err
}

// httpGet делает GET с декомпрессией; headers перекрывают дефолты.
func httpGet(client *http.Client, rawURL string, headers map[string]string) (string, int, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := readBody(resp)
	return body, resp.StatusCode, err
}
