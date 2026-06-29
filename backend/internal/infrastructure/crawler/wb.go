// internal/infrastructure/crawler/wb.go
package crawler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// WbClient берёт имя+картинку WB с контентного CDN (basket-NN.wbbasket.ru).
// Цена WB — вне объёма: живая цена за анти-бот-токеном API.
type WbClient struct {
	client *http.Client
}

func NewWbClient() *WbClient {
	return &WbClient{client: &http.Client{Timeout: 20 * time.Second}}
}

func (c *WbClient) Supports(host string) bool { return host == "wildberries.ru" }

const wbMaxBasket = 40

var wbNmRe = regexp.MustCompile(`/catalog/(\d+)`)

func parseWbNm(rawURL string) (int64, error) {
	m := wbNmRe.FindStringSubmatch(rawURL)
	if len(m) < 2 {
		return 0, fmt.Errorf("wb: артикул не найден в %s", rawURL)
	}
	return strconv.ParseInt(m[1], 10, 64)
}

// volToBasket — таблица соответствия vol → номер basket-шарда WB.
// ОБНОВЛЯТЬ при росте WB (добавлении новых basket-ов); добор сканом ниже страхует.
func volToBasket(vol int64) int {
	bounds := []struct {
		max    int64
		basket int
	}{
		{143, 1}, {287, 2}, {431, 3}, {719, 4}, {1007, 5}, {1061, 6}, {1115, 7},
		{1169, 8}, {1313, 9}, {1601, 10}, {1655, 11}, {1919, 12}, {2045, 13},
		{2189, 14}, {2405, 15}, {2621, 16}, {2837, 17}, {3053, 18}, {3269, 19},
		{3485, 20}, {3701, 21}, {3917, 22}, {4133, 23}, {4349, 24}, {4565, 25},
		{4877, 26}, {5189, 27}, {5501, 28}, {5813, 29}, {6125, 30},
	}
	for _, b := range bounds {
		if vol <= b.max {
			return b.basket
		}
	}
	return 31
}

func wbCardURL(basket int, vol, part, nm int64) string {
	return fmt.Sprintf("https://basket-%02d.wbbasket.ru/vol%d/part%d/%d/info/ru/card.json", basket, vol, part, nm)
}

func wbImageURL(basket int, vol, part, nm int64) string {
	return fmt.Sprintf("https://basket-%02d.wbbasket.ru/vol%d/part%d/%d/images/big/1.webp", basket, vol, part, nm)
}

func parseWbCardName(jsonBody string) string {
	var card struct {
		ImtName string `json:"imt_name"`
	}
	_ = json.Unmarshal([]byte(jsonBody), &card)
	return card.ImtName
}

func (c *WbClient) FetchProductByURL(productURL string) (*ProductInfo, error) {
	nm, err := parseWbNm(productURL)
	if err != nil {
		return nil, err
	}
	vol, part := nm/100000, nm/1000

	basket, body, err := c.fetchCard(vol, part, nm)
	if err != nil {
		return nil, err
	}
	name := parseWbCardName(body)
	if name == "" {
		return nil, fmt.Errorf("wb: imt_name пуст для nm=%d", nm)
	}
	return &ProductInfo{
		Title:    name,
		ImageURL: wbImageURL(basket, vol, part, nm),
		URL:      productURL,
	}, nil
}

// wbCardStatus классифицирует статус ответа basket-CDN: ok=карточка найдена;
// cont=искать в следующем basket (404). Иначе (429/5xx/блок) — стоп скана.
func wbCardStatus(status int) (ok, cont bool) {
	switch status {
	case http.StatusOK:
		return true, false
	case http.StatusNotFound:
		return false, true
	default:
		return false, false
	}
}

// fetchCard пробует вычисленный basket, затем сканирует известный диапазон.
// Прерываем скан при сетевой ошибке ИЛИ не-404 статусе (429/5xx/блок), чтобы не
// слать до wbMaxBasket запросов при rate-limit/недоступности CDN.
func (c *WbClient) fetchCard(vol, part, nm int64) (int, string, error) {
	try := func(b int) (string, bool, error) {
		body, status, err := httpGet(c.client, wbCardURL(b, vol, part, nm), map[string]string{"User-Agent": productUserAgent})
		if err != nil {
			return "", false, err
		}
		ok, cont := wbCardStatus(status)
		if ok {
			return body, true, nil
		}
		if cont {
			return "", false, nil
		}
		return "", false, fmt.Errorf("wb: basket-CDN вернул %d", status)
	}

	computed := volToBasket(vol)
	if body, ok, err := try(computed); err != nil {
		return 0, "", err
	} else if ok {
		return computed, body, nil
	}

	for b := 1; b <= wbMaxBasket; b++ {
		if b == computed {
			continue
		}
		if body, ok, err := try(b); err != nil {
			return 0, "", err
		} else if ok {
			return b, body, nil
		}
	}
	return 0, "", fmt.Errorf("wb: card.json не найден для nm=%d", nm)
}
