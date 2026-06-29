package favorite

import (
	"context"
	"errors"
	"testing"
	"wishlist-go/internal/domain"

	"github.com/google/uuid"
)

type mockFavRepo struct {
	added     *domain.Favorite
	removed   bool
	lists     []*domain.Wishlist
	removeErr error
}

func (m *mockFavRepo) Add(_ context.Context, f *domain.Favorite) error { m.added = f; return nil }
func (m *mockFavRepo) Remove(_ context.Context, _ int64, _ uuid.UUID) error {
	m.removed = true
	return m.removeErr
}
func (m *mockFavRepo) ListByAccount(context.Context, int64) ([]*domain.Wishlist, error) {
	return m.lists, nil
}

type mockWlRepo struct {
	wl  *domain.Wishlist
	err error
}

func (m *mockWlRepo) CreateWishlist(context.Context, *domain.Wishlist) error { return nil }
func (m *mockWlRepo) GetWishlistByCode(context.Context, uuid.UUID) (*domain.Wishlist, error) {
	return m.wl, m.err
}
func (m *mockWlRepo) GetWishlistsByOwnerID(context.Context, int64, int, int) ([]*domain.Wishlist, error) {
	return nil, nil
}
func (m *mockWlRepo) UpdateWishlist(context.Context, *domain.Wishlist) error { return nil }
func (m *mockWlRepo) DeleteWishlist(context.Context, uuid.UUID) error        { return nil }

func TestAdd_WishlistMustExist(t *testing.T) {
	svc := NewService(&mockFavRepo{}, &mockWlRepo{err: domain.ErrWishlistNotFound})
	if err := svc.Add(context.Background(), 1, uuid.New()); !errors.Is(err, domain.ErrWishlistNotFound) {
		t.Fatalf("добавление несуществующего списка: ожидался ErrWishlistNotFound, got %v", err)
	}
}

func TestAdd_HappyPath(t *testing.T) {
	fav := &mockFavRepo{}
	code := uuid.New()
	svc := NewService(fav, &mockWlRepo{wl: &domain.Wishlist{ShareCode: code}})

	if err := svc.Add(context.Background(), 42, code); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fav.added == nil || fav.added.AccountID != 42 || fav.added.WishlistCode != code {
		t.Fatalf("Favorite сохранён неверно: %+v", fav.added)
	}
}

func TestRemoveAndList_Delegate(t *testing.T) {
	fav := &mockFavRepo{lists: []*domain.Wishlist{{Name: "x"}}}
	svc := NewService(fav, &mockWlRepo{})

	if err := svc.Remove(context.Background(), 1, uuid.New()); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if !fav.removed {
		t.Fatal("Remove должен делегироваться репозиторию")
	}
	got, err := svc.List(context.Background(), 1)
	if err != nil || len(got) != 1 {
		t.Fatalf("List = %d, %v; ожидали 1", len(got), err)
	}
}
