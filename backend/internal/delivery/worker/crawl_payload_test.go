package worker

import (
	"testing"
	"wishlist-go/internal/infrastructure/crawler"

	"github.com/google/uuid"
)

func TestCrawlPayloadToWishItem_HappyPath(t *testing.T) {
	code := uuid.New()
	info := &crawler.ProductInfo{Title: "gift", Price: "12345 ₽", ImageURL: "http://pic", URL: "http://market/x"}
	payload := map[string]interface{}{
		"wish_list_code": code.String(),
		"owner_id":       float64(42),
	}

	item, err := crawlPayloadToWishItem(payload, info)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.WishListCode != code || item.OwnerID != 42 {
		t.Fatalf("ids неверны: %+v", item)
	}
	if item.Name == nil || *item.Name != "gift" {
		t.Fatalf("name неверно: %v", item.Name)
	}
	if item.MarketPrice == nil || *item.MarketPrice != 12345 {
		t.Fatalf("price неверна: %v", item.MarketPrice)
	}
	if item.MarketCurrency != "RUB" || item.MarketURL != "http://market/x" {
		t.Fatalf("market-поля неверны: %+v", item)
	}
}

func TestCrawlPayloadToWishItem_BadPayload(t *testing.T) {
	info := &crawler.ProductInfo{Title: "x"}

	if _, err := crawlPayloadToWishItem(map[string]interface{}{"wish_list_code": "not-uuid", "owner_id": float64(1)}, info); err == nil {
		t.Fatal("ожидалась ошибка на невалидный wish_list_code")
	}
	if _, err := crawlPayloadToWishItem(map[string]interface{}{"wish_list_code": uuid.New().String()}, info); err == nil {
		t.Fatal("ожидалась ошибка на отсутствующий owner_id")
	}
}
