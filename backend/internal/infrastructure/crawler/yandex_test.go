package crawler

import "testing"

// Реальный Я.Маркет отдаёт price ЧИСЛОМ без кавычек ("price":5304), а не строкой.
// Парсер обязан извлекать оба варианта.
func TestParseProductPage_ExtractsNumericPrice(t *testing.T) {
	html := `<!--BEGIN [@marketfront/JsonLd] productPageMicromarkup -->` +
		`<script type="application/ld+json">` +
		`{"@type":"Product","name":"Док-станция UGREEN CM198","image":"https://avatars.mds.yandex.net/x/orig",` +
		`"offers":{"@type":"Offer","availability":"https://schema.org/InStock","price":5304,"priceCurrency":"RUB","url":"https://market.yandex.ru/card/x/5089436739"},` +
		`"url":"https://market.yandex.ru/card/x/5089436739"}` +
		`</script>`

	c := NewYaMarketClient()
	p, err := c.parseProductPage(html, "https://market.yandex.ru/card/x/5089436739?utm=1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Title != "Док-станция UGREEN CM198" {
		t.Fatalf("title: %q", p.Title)
	}
	if p.ImageURL != "https://avatars.mds.yandex.net/x/orig" {
		t.Fatalf("image: %q", p.ImageURL)
	}
	if p.Price != "5304 ₽" {
		t.Fatalf("price: %q, ожидали \"5304 ₽\"", p.Price)
	}
}

// Обратная совместимость: цена строкой ("price":"5304") тоже извлекается.
func TestParseProductPage_ExtractsQuotedPrice(t *testing.T) {
	html := `<!--BEGIN [@marketfront/JsonLd] productPageMicromarkup -->` +
		`<script type="application/ld+json">` +
		`{"@type":"Product","name":"Товар","image":"https://img/x",` +
		`"offers":{"@type":"Offer","price":"1234","priceCurrency":"RUB"}}` +
		`</script>`

	c := NewYaMarketClient()
	p, err := c.parseProductPage(html, "https://market.yandex.ru/card/x/1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Price != "1234 ₽" {
		t.Fatalf("price: %q, ожидали \"1234 ₽\"", p.Price)
	}
}
