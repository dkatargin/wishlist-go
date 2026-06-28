package wishitem

import (
	"context"
	"testing"
	"wishlist-go/internal/domain"

	"github.com/google/uuid"
)

type mockWishItemRepo struct {
	created *domain.WishItem
}

func (m *mockWishItemRepo) CreateWishItem(wi *domain.WishItem) error { m.created = wi; return nil }
func (m *mockWishItemRepo) GetWishItemByID(int64, uuid.UUID) (*domain.WishItem, error) {
	return nil, nil
}
func (m *mockWishItemRepo) GetWishItemsByWishlistID(uuid.UUID, int, int) ([]*domain.WishItem, error) {
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
