// internal/infrastructure/crawler/ozon.go
package crawler

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// OzonClient берёт карточку Ozon через OG-разметку.
// Ozon отдаёт OG только на UA мессенджер-краулеров (allow-list по UA).
type OzonClient struct {
	client *http.Client
	uas    []string
}

func NewOzonClient() *OzonClient {
	return &OzonClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
			// При анти-бот-блоке Ozon зацикливает редиректы (?__rr=...); ограничиваем,
			// чтобы блок проваливался быстро (последний ответ), а не после 10 редиректов.
			CheckRedirect: func(_ *http.Request, via []*http.Request) error {
				if len(via) >= 3 {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
		// Только WhatsApp: единственный превью-UA, который Ozon реально пускает (allow-list).
		// facebook/twitter ловят редирект-цикл и не помогают — добавлять рабочие по факту.
		uas: []string{"WhatsApp/2.23"},
	}
}

func (c *OzonClient) Supports(host string) bool { return host == "ozon.ru" }

func (c *OzonClient) FetchProductByURL(productURL string) (*ProductInfo, error) {
	for _, ua := range c.uas {
		body, status, err := httpGet(c.client, productURL, map[string]string{"User-Agent": ua})
		if err != nil || status != http.StatusOK {
			continue
		}
		info := parseOzon(body, productURL)
		if info.Title != "" {
			return info, nil
		}
	}
	return nil, fmt.Errorf("ozon: карточка не получена для %s", productURL)
}

var (
	ozonTitleTag = regexp.MustCompile(`<meta\b[^>]*property="og:title"[^>]*>`)
	ozonImageTag = regexp.MustCompile(`<meta\b[^>]*property="og:image"[^>]*>`)
	ozonContent  = regexp.MustCompile(`content="([^"]*)"`)
	ozonPrice    = regexp.MustCompile(`"price"\s*:\s*"?([0-9.]+)"?\s*,\s*"priceCurrency"\s*:\s*"[A-Z]+"`)
)

// ozonTagContent находит <meta>-тег по precompiled-регэкспу и достаёт content,
// не завися от порядка атрибутов внутри тега.
func ozonTagContent(tagRe *regexp.Regexp, html string) string {
	tag := tagRe.FindString(html)
	if tag == "" {
		return ""
	}
	if m := ozonContent.FindStringSubmatch(tag); len(m) > 1 {
		return m[1]
	}
	return ""
}

// parseOzon достаёт имя/картинку/цену из HTML карточки Ozon.
func parseOzon(html, sourceURL string) *ProductInfo {
	p := &ProductInfo{URL: sourceURL}
	p.Title = strings.TrimSpace(ozonTagContent(ozonTitleTag, html))
	p.ImageURL = safeImageURL(ozonTagContent(ozonImageTag, html))
	if m := ozonPrice.FindStringSubmatch(html); len(m) > 1 {
		p.Price = m[1] + " ₽"
	}
	return p
}
