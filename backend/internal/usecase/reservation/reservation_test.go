package reservation

import (
	"context"
	"errors"
	"testing"
	"wishlist-go/internal/domain"

	"github.com/google/uuid"
)

type mockResRepo struct {
	created         *domain.Reservation
	existing        *domain.Reservation
	sum             int
	byList          []*domain.Reservation
	deletedID       uuid.UUID
	deletedReserver int64
}

func (m *mockResRepo) Create(_ context.Context, r *domain.Reservation) error {
	m.created = r
	return nil
}
func (m *mockResRepo) Delete(_ context.Context, id uuid.UUID, reserverID int64) error {
	m.deletedID = id
	m.deletedReserver = reserverID
	return nil
}
func (m *mockResRepo) MarkDone(context.Context, uuid.UUID, int64) error { return nil }
func (m *mockResRepo) GetByReserverAndWish(context.Context, int64, int64) (*domain.Reservation, error) {
	if m.existing != nil {
		return m.existing, nil
	}
	return nil, domain.ErrReservationNotFound
}
func (m *mockResRepo) SumQuantityByWish(context.Context, int64) (int, error) { return m.sum, nil }
func (m *mockResRepo) ListByWishlist(context.Context, uuid.UUID) ([]*domain.Reservation, error) {
	return m.byList, nil
}
func (m *mockResRepo) ListByReserver(context.Context, int64) ([]*domain.Reservation, error) {
	return m.byList, nil
}

type mockWishRepo struct {
	wish    *domain.WishItem
	wishErr error
}

func (m *mockWishRepo) CreateWishItem(*domain.WishItem) error { return nil }
func (m *mockWishRepo) GetWishItemByID(int64, uuid.UUID) (*domain.WishItem, error) {
	return m.wish, m.wishErr
}
func (m *mockWishRepo) GetWishItemsByWishlistID(uuid.UUID, int, int, bool) ([]*domain.WishItem, error) {
	return nil, nil
}
func (m *mockWishRepo) UpdateWishItem(int64, uuid.UUID, domain.WishItemUpdate) error { return nil }
func (m *mockWishRepo) DeleteWishItem(int64) error                                   { return nil }

func wishOwned(ownerID int64, qty int) *domain.WishItem {
	q := qty
	return &domain.WishItem{ID: 1, OwnerID: ownerID, MarketQuantity: &q}
}

func TestReserve_SelfReservationForbidden(t *testing.T) {
	svc := NewService(&mockResRepo{}, &mockWishRepo{wish: wishOwned(5, 0)})
	if _, err := svc.Reserve(context.Background(), 5, uuid.New(), 1, 1, false); !errors.Is(err, domain.ErrSelfReservation) {
		t.Fatalf("резерв на своём списке: ожидался ErrSelfReservation, got %v", err)
	}
}

func TestReserve_WishNotFound(t *testing.T) {
	svc := NewService(&mockResRepo{}, &mockWishRepo{wishErr: domain.ErrWishItemNotFound})
	if _, err := svc.Reserve(context.Background(), 9, uuid.New(), 1, 1, false); !errors.Is(err, domain.ErrWishItemNotFound) {
		t.Fatalf("ожидался ErrWishItemNotFound, got %v", err)
	}
}

func TestReserve_AlreadyReserved(t *testing.T) {
	svc := NewService(&mockResRepo{existing: &domain.Reservation{}}, &mockWishRepo{wish: wishOwned(5, 0)})
	if _, err := svc.Reserve(context.Background(), 9, uuid.New(), 1, 1, false); !errors.Is(err, domain.ErrReservationExists) {
		t.Fatalf("ожидался ErrReservationExists, got %v", err)
	}
}

func TestReserve_QuantityExceeded(t *testing.T) {
	// market_quantity=3, уже зарезервировано 2, просим ещё 2 → 4 > 3
	svc := NewService(&mockResRepo{sum: 2}, &mockWishRepo{wish: wishOwned(5, 3)})
	if _, err := svc.Reserve(context.Background(), 9, uuid.New(), 1, 2, false); !errors.Is(err, domain.ErrReserveQuantityExceeded) {
		t.Fatalf("ожидался ErrReserveQuantityExceeded, got %v", err)
	}
}

func TestReserve_HappyPath(t *testing.T) {
	res := &mockResRepo{sum: 1}
	svc := NewService(res, &mockWishRepo{wish: wishOwned(5, 3)})
	out, err := svc.Reserve(context.Background(), 9, uuid.New(), 1, 1, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.created == nil || res.created.ReserverID != 9 || res.created.Quantity != 1 || !res.created.IsAnonymous {
		t.Fatalf("резерв сохранён неверно: %+v", res.created)
	}
	if out.WishID != 1 {
		t.Fatalf("вернулся неверный резерв: %+v", out)
	}
}

func TestAggregate_AnonymousHiddenButCounted(t *testing.T) {
	res := &mockResRepo{byList: []*domain.Reservation{
		{WishID: 1, ReserverID: 7, Quantity: 2, IsAnonymous: false},
		{WishID: 1, ReserverID: 8, Quantity: 1, IsAnonymous: true},
	}}
	svc := NewService(res, &mockWishRepo{})

	agg, err := svc.AggregateForSharedList(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v := agg[1]
	if v.ReservedQuantity != 3 {
		t.Fatalf("анонимный резерв должен учитываться в сумме: %d", v.ReservedQuantity)
	}
	if len(v.Reservers) != 1 || v.Reservers[0] != 7 {
		t.Fatalf("в reservers должен быть только не-анонимный (7): %v", v.Reservers)
	}
}

func TestCancelByWish_DeletesOwnReservation(t *testing.T) {
	existing := &domain.Reservation{ID: uuid.New(), WishID: 1, ReserverID: 9}
	res := &mockResRepo{existing: existing}
	svc := NewService(res, &mockWishRepo{})

	if err := svc.CancelByWish(context.Background(), 9, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.deletedID != existing.ID || res.deletedReserver != 9 {
		t.Fatalf("должен удалить найденный резерв своего дарителя: id=%v reserver=%d", res.deletedID, res.deletedReserver)
	}
}

func TestCancelByWish_NotFound(t *testing.T) {
	svc := NewService(&mockResRepo{}, &mockWishRepo{})
	if err := svc.CancelByWish(context.Background(), 9, 1); !errors.Is(err, domain.ErrReservationNotFound) {
		t.Fatalf("ожидался ErrReservationNotFound, got %v", err)
	}
}
