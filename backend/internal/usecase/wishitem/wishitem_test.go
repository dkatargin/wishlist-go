package wishitem

import (
	"context"
	"errors"
	"testing"
	"wishlist-go/internal/domain"

	"github.com/google/uuid"
)

type mockWishItemRepo struct {
	created    *domain.WishItem
	lastActive bool
}

func (m *mockWishItemRepo) CreateWishItem(wi *domain.WishItem) error { m.created = wi; return nil }
func (m *mockWishItemRepo) GetWishItemByID(int64, uuid.UUID) (*domain.WishItem, error) {
	return nil, nil
}
func (m *mockWishItemRepo) GetWishItemsByWishlistID(_ uuid.UUID, _ int, _ int, onlyActive bool) ([]*domain.WishItem, error) {
	m.lastActive = onlyActive
	return nil, nil
}
func (m *mockWishItemRepo) UpdateWishItem(int64, uuid.UUID, domain.WishItemUpdate) error {
	return nil
}
func (m *mockWishItemRepo) DeleteWishItem(int64) error { return nil }

func TestCreateWishItem_DefaultsCurrencyAndPersistsOwner(t *testing.T) {
	repo := &mockWishItemRepo{}
	svc := NewService(repo, nil, nil) // wishlistRepo и mq не нужны на этом пути

	name := "gift"
	item := &domain.WishItem{WishListCode: uuid.New(), OwnerID: 1, Priority: 3, Name: &name}

	out, err := svc.CreateWishItem(context.Background(), item)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.MarketCurrency != "RUB" {
		t.Fatalf("expected default currency RUB, got %q", out.MarketCurrency)
	}
	if repo.created == nil || repo.created.OwnerID != 1 || repo.created.Priority != 3 {
		t.Fatalf("item not persisted with owner/priority: %+v", repo.created)
	}
}

type mockPublisher struct {
	msgType string
	payload map[string]interface{}
	err     error
}

func (m *mockPublisher) PublishMessage(msgType string, payload map[string]interface{}) error {
	m.msgType = msgType
	m.payload = payload
	return m.err
}

func TestRequestCrawl_PublishesCrawlProduct(t *testing.T) {
	pub := &mockPublisher{}
	svc := NewService(nil, nil, pub)
	code := uuid.New()
	url := "https://market.yandex.ru/product/123"

	if err := svc.RequestCrawl(context.Background(), code, 42, url); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pub.msgType != "crawl_product" {
		t.Fatalf("ожидался тип crawl_product, получили %q", pub.msgType)
	}
	if pub.payload["product_url"] != url ||
		pub.payload["wish_list_code"] != code.String() ||
		pub.payload["owner_id"] != int64(42) {
		t.Fatalf("неверный payload: %v", pub.payload)
	}
}

func TestRequestCrawl_RejectsBadURL(t *testing.T) {
	pub := &mockPublisher{}
	svc := NewService(nil, nil, pub)

	for _, bad := range []string{
		"http://market.yandex.ru/x",           // не https
		"https://market.yandex.ru.attacker/x", // обход подстроки в host
		"https://evil.test/market.yandex.ru",  // обход подстроки в path
		"://broken",
	} {
		if err := svc.RequestCrawl(context.Background(), uuid.New(), 1, bad); !errors.Is(err, ErrInvalidMarketURL) {
			t.Fatalf("URL %q: ожидался ErrInvalidMarketURL, получили %v", bad, err)
		}
	}
	if pub.msgType != "" {
		t.Fatal("при невалидном URL ничего не должно публиковаться")
	}
}

func TestRequestCrawl_PropagatesPublisherError(t *testing.T) {
	boom := errors.New("publish failed")
	svc := NewService(nil, nil, &mockPublisher{err: boom})

	if err := svc.RequestCrawl(context.Background(), uuid.New(), 1, "https://market.yandex.ru/x"); !errors.Is(err, boom) {
		t.Fatalf("ожидалась ошибка publisher, получили %v", err)
	}
}

func TestVisibleVsAll_PassesOnlyActiveFlag(t *testing.T) {
	repo := &mockWishItemRepo{}
	svc := NewService(repo, nil, nil)

	_, _ = svc.GetVisibleByWishlist(context.Background(), uuid.New(), 10, 0)
	if !repo.lastActive {
		t.Fatal("GetVisibleByWishlist должен прокидывать onlyActive=true")
	}
	_, _ = svc.GetWishItemsByWishlist(context.Background(), uuid.New(), 10, 0)
	if repo.lastActive {
		t.Fatal("GetWishItemsByWishlist должен прокидывать onlyActive=false")
	}
}
