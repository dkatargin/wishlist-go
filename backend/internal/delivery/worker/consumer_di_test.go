package worker

import (
	"testing"
	"wishlist-go/internal/infrastructure/crawler"
)

// stubCrawler — заглушка для проверки DI-шва. FetchProductByURL не вызывается
// в этом тесте, поэтому возврат (nil, nil) допустим и контракт не нарушает.
type stubCrawler struct{}

func (stubCrawler) FetchProductByURL(string) (*crawler.ProductInfo, error) { return nil, nil }

// Регресс: consumer должен зависеть от интерфейса crawler.Crawler,
// чтобы диспетчер с любыми адаптерами подставлялся без правок воркера.
func TestNewConsumer_AcceptsCrawlerInterface(t *testing.T) {
	var c crawler.Crawler = stubCrawler{}
	cons := NewConsumer(nil, c, nil)
	if cons.crawler == nil {
		t.Fatal("crawler не проинжектен")
	}
}
