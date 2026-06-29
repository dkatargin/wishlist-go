package crawler

import (
	"errors"
	"testing"
)

type stubAdapter struct {
	host string
	info *ProductInfo
}

func (s stubAdapter) Supports(host string) bool { return host == s.host }
func (s stubAdapter) FetchProductByURL(string) (*ProductInfo, error) {
	return s.info, nil
}

func TestDispatcher_RoutesByHost(t *testing.T) {
	ya := stubAdapter{host: "market.yandex.ru", info: &ProductInfo{Title: "ya"}}
	oz := stubAdapter{host: "ozon.ru", info: &ProductInfo{Title: "oz"}}
	d := NewDispatcher(ya, oz)

	cases := map[string]string{
		"https://market.yandex.ru/card/x/1":         "ya",
		"https://www.ozon.ru/product/x-1991629925/": "oz",
		"https://OZON.ru/product/x/":                "oz",
	}
	for rawURL, want := range cases {
		info, err := d.FetchProductByURL(rawURL)
		if err != nil || info == nil || info.Title != want {
			t.Fatalf("%s -> %+v, %v; ожидали %q", rawURL, info, err, want)
		}
	}
}

func TestDispatcher_UnknownHost(t *testing.T) {
	d := NewDispatcher(stubAdapter{host: "ozon.ru"})
	if _, err := d.FetchProductByURL("https://example.com/x"); !errors.Is(err, ErrUnsupportedSource) {
		t.Fatalf("ожидался ErrUnsupportedSource, got %v", err)
	}
	if _, err := d.FetchProductByURL("::not a url::"); !errors.Is(err, ErrUnsupportedSource) {
		t.Fatalf("битый URL: ожидался ErrUnsupportedSource, got %v", err)
	}
}
