package crawler

import "testing"

func TestYaMarketClient_Supports(t *testing.T) {
	c := NewYaMarketClient()
	yes := []string{"market.yandex.ru", "yandex.ru"}
	no := []string{"ozon.ru", "wildberries.ru", "example.com"}
	for _, h := range yes {
		if !c.Supports(h) {
			t.Fatalf("должен поддерживать %q", h)
		}
	}
	for _, h := range no {
		if c.Supports(h) {
			t.Fatalf("не должен поддерживать %q", h)
		}
	}
	// проверка, что тип удовлетворяет Adapter
	var _ Adapter = NewYaMarketClient()
}
