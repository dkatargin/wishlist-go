package postgres

import (
	"testing"
	"wishlist-go/internal/domain"

	"github.com/google/uuid"
)

func TestWishItemConverters_RoundTripFull(t *testing.T) {
	name := "gift"
	pic := "http://pic"
	price := 99.5
	qty := 3
	in := &domain.WishItem{
		ID:               7,
		WishListCode:     uuid.New(),
		OwnerID:          42,
		Priority:         5,
		IsDone:           true,
		Name:             &name,
		MarketURL:        "http://market",
		MarketPictureURL: &pic,
		MarketPrice:      &price,
		MarketCurrency:   "USD",
		MarketQuantity:   &qty,
	}

	out := modelToDomainWishItem(domainWishItemToModel(in))

	if out.OwnerID != 42 || out.Priority != 5 || !out.IsDone {
		t.Fatalf("scalar fields lost: %+v", out)
	}
	if out.MarketCurrency != "USD" {
		t.Fatalf("currency lost: %q", out.MarketCurrency)
	}
	if out.MarketQuantity == nil || *out.MarketQuantity != 3 {
		t.Fatalf("quantity lost: %v", out.MarketQuantity)
	}
	if out.Name == nil || *out.Name != "gift" {
		t.Fatalf("name lost: %v", out.Name)
	}
}

func TestDomainWishItemToModel_NilSafe(t *testing.T) {
	// Все указатели nil — конвертер не должен паниковать.
	in := &domain.WishItem{WishListCode: uuid.New(), OwnerID: 1}
	_ = domainWishItemToModel(in)
}
