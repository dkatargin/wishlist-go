//go:build live

// Живой smoke-тест парсеров по реальным ссылкам. По умолчанию НЕ компилируется
// (build-тег live), чтобы обычный `go test ./...` оставался герметичным.
//
// Запуск (лучше с прод-egress, напр. dedicated-сервера):
//
//	go test -tags live -run TestLive_Sources -v ./internal/infrastructure/crawler/
//
// Проверяет, что диспетчер по хосту достаёт хотя бы название для каждого источника,
// и логирует имя/цену/картинку для глазной проверки. Имена/цены меняются на сайтах —
// тест НЕ ассертит конкретные значения, только непустое имя.
package crawler

import "testing"

func TestLive_Sources(t *testing.T) {
	d := NewDispatcher(NewYaMarketClient(), NewOzonClient(), NewWbClient())

	cases := []struct {
		name, url string
		wantImage bool // у кого ждём ещё и картинку
	}{
		{"yandex", "https://market.yandex.ru/card/dok-stantsiya-ugreen-cm198-dlya-hddssd-2535-c-usb-30/5089436739", true},
		{"ozon", "https://www.ozon.ru/product/8bitdo-geympad-m30-bluetooth-bluetooth-provodnoy-belyy-1991629925/", true},
		{"wb", "https://www.wildberries.ru/catalog/420939815/detail.aspx", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p, err := d.FetchProductByURL(c.url)
			if err != nil {
				t.Fatalf("%s: %v", c.name, err)
			}
			if p.Title == "" {
				t.Fatalf("%s: пустое название: %+v", c.name, p)
			}
			if c.wantImage && p.ImageURL == "" {
				t.Errorf("%s: ожидали картинку, пусто", c.name)
			}
			t.Logf("%s ok name=%q price=%q image=%q", c.name, p.Title, p.Price, p.ImageURL)
		})
	}
}
