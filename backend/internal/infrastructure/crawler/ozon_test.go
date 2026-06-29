// internal/infrastructure/crawler/ozon_test.go
package crawler

import "testing"

const ozonFixture = `<html><head>
<meta data-hid="property::og:title" property="og:title" content="8BitDo Геймпад M30 Bluetooth ">
<meta property="og:image" content="https://ir.ozone.ru/s3/x/c600/7433553559.jpg">
<meta property="og:description" content="Успейте купить">
</head><body>
<script>var x={"availability":"http://schema.org/InStock","price":"2803","priceCurrency":"RUB","sku":"1"};</script>
</body></html>`

func TestParseOzon(t *testing.T) {
	p := parseOzon(ozonFixture, "https://www.ozon.ru/product/x-1991629925/?z=1")
	if p.Title != "8BitDo Геймпад M30 Bluetooth" {
		t.Fatalf("title: %q", p.Title)
	}
	if p.ImageURL != "https://ir.ozone.ru/s3/x/c600/7433553559.jpg" {
		t.Fatalf("image: %q", p.ImageURL)
	}
	if p.Price != "2803 ₽" {
		t.Fatalf("price: %q", p.Price)
	}
	if p.URL != "https://www.ozon.ru/product/x-1991629925/?z=1" {
		t.Fatalf("url: %q", p.URL)
	}
}

func TestParseOzon_NoData(t *testing.T) {
	p := parseOzon(`<html><head></head><body>пусто</body></html>`, "https://www.ozon.ru/x")
	if p.Title != "" {
		t.Fatalf("ожидали пустой Title, got %q", p.Title)
	}
}

func TestParseOzon_RejectsUnsafeImage(t *testing.T) {
	html := `<meta property="og:title" content="X"><meta property="og:image" content="javascript:alert(1)">`
	if p := parseOzon(html, "https://www.ozon.ru/x"); p.ImageURL != "" {
		t.Fatalf("небезопасную картинку нужно отбросить, got %q", p.ImageURL)
	}
}

func TestOzonClient_Supports(t *testing.T) {
	var _ Adapter = NewOzonClient()
	if !NewOzonClient().Supports("ozon.ru") || NewOzonClient().Supports("wildberries.ru") {
		t.Fatal("Supports ozon.ru неверен")
	}
}

func TestParseOzon_AttrOrderIndependent(t *testing.T) {
	// content раньше property — HTML не гарантирует порядок атрибутов
	html := `<meta content="Товар X" property="og:title"><meta content="https://img/y.jpg" property="og:image">`
	p := parseOzon(html, "https://www.ozon.ru/x")
	if p.Title != "Товар X" {
		t.Fatalf("title при обратном порядке атрибутов: %q", p.Title)
	}
	if p.ImageURL != "https://img/y.jpg" {
		t.Fatalf("image при обратном порядке атрибутов: %q", p.ImageURL)
	}
}

func TestParseOzon_NumericPrice(t *testing.T) {
	html := `<meta property="og:title" content="X"><script>{"price":2803,"priceCurrency":"RUB"}</script>`
	if p := parseOzon(html, "https://www.ozon.ru/x"); p.Price != "2803 ₽" {
		t.Fatalf("числовая цена: %q", p.Price)
	}
}
