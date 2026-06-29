package crawler

import (
	"errors"
	"net/url"
	"strings"
)

// ErrUnsupportedSource — хост ссылки не обслуживается ни одним адаптером.
var ErrUnsupportedSource = errors.New("unsupported product source")

// Crawler — единый контракт парсера карточки товара.
type Crawler interface {
	FetchProductByURL(productURL string) (*ProductInfo, error)
}

// Adapter — Crawler, привязанный к источнику (хосту).
type Adapter interface {
	Crawler
	Supports(host string) bool
}

// Dispatcher выбирает адаптер по хосту URL.
type Dispatcher struct {
	adapters []Adapter
}

func NewDispatcher(adapters ...Adapter) *Dispatcher {
	return &Dispatcher{adapters: adapters}
}

func (d *Dispatcher) FetchProductByURL(productURL string) (*ProductInfo, error) {
	u, err := url.Parse(productURL)
	if err != nil || u.Hostname() == "" {
		return nil, ErrUnsupportedSource
	}
	host := strings.TrimPrefix(strings.ToLower(u.Hostname()), "www.")
	for _, a := range d.adapters {
		if a.Supports(host) {
			return a.FetchProductByURL(productURL)
		}
	}
	return nil, ErrUnsupportedSource
}
