package wishlist

import (
	"context"
	"errors"
	"testing"
	"wishlist-go/internal/domain"

	"github.com/google/uuid"
)

type mockWishlistRepo struct {
	getByCode func(uuid.UUID) (*domain.Wishlist, error)
}

func (m *mockWishlistRepo) CreateWishlist(context.Context, *domain.Wishlist) error { return nil }
func (m *mockWishlistRepo) GetWishlistByCode(_ context.Context, code uuid.UUID) (*domain.Wishlist, error) {
	return m.getByCode(code)
}
func (m *mockWishlistRepo) GetWishlistsByOwnerID(context.Context, int64, int, int) ([]*domain.Wishlist, error) {
	return nil, nil
}
func (m *mockWishlistRepo) UpdateWishlist(context.Context, *domain.Wishlist) error { return nil }
func (m *mockWishlistRepo) DeleteWishlist(context.Context, uuid.UUID) error        { return nil }

func TestGetDetail_OwnerFlag(t *testing.T) {
	code := uuid.New()
	repo := &mockWishlistRepo{getByCode: func(uuid.UUID) (*domain.Wishlist, error) {
		return &domain.Wishlist{ShareCode: code, OwnerID: 5}, nil
	}}
	svc := NewService(repo)

	if _, isOwner, err := svc.GetDetail(context.Background(), code, 5); err != nil || !isOwner {
		t.Fatalf("владелец: ожидался isOwner=true, err=nil; получили isOwner=%v err=%v", isOwner, err)
	}
	if _, isOwner, err := svc.GetDetail(context.Background(), code, 9); err != nil || isOwner {
		t.Fatalf("не-владелец: ожидался isOwner=false, err=nil; получили isOwner=%v err=%v", isOwner, err)
	}
}

func TestGetDetail_NotFoundPropagates(t *testing.T) {
	repo := &mockWishlistRepo{getByCode: func(uuid.UUID) (*domain.Wishlist, error) {
		return nil, domain.ErrWishlistNotFound
	}}
	svc := NewService(repo)

	if _, _, err := svc.GetDetail(context.Background(), uuid.New(), 1); !errors.Is(err, domain.ErrWishlistNotFound) {
		t.Fatalf("ожидался ErrWishlistNotFound, получили %v", err)
	}
}

func TestGetByCode_NotFoundPropagates(t *testing.T) {
	repo := &mockWishlistRepo{getByCode: func(uuid.UUID) (*domain.Wishlist, error) {
		return nil, domain.ErrWishlistNotFound
	}}
	svc := NewService(repo)

	if _, err := svc.GetByCode(context.Background(), uuid.New()); !errors.Is(err, domain.ErrWishlistNotFound) {
		t.Fatalf("ожидался ErrWishlistNotFound, получили %v", err)
	}
}
