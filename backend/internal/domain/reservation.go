package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Reservation — бронь желания дарителем (НЕ владельцем списка).
type Reservation struct {
	ID          uuid.UUID
	WishID      int64 // id желания (WishItem.ID)
	ReserverID  int64 // account дарителя
	Quantity    int
	IsAnonymous bool
	IsDone      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ReservationRepository interface {
	Create(ctx context.Context, r *Reservation) error
	// Delete/MarkDone скоупятся по reserverID — менять может только автор брони.
	Delete(ctx context.Context, id uuid.UUID, reserverID int64) error
	MarkDone(ctx context.Context, id uuid.UUID, reserverID int64) error
	GetByReserverAndWish(ctx context.Context, reserverID, wishID int64) (*Reservation, error)
	SumQuantityByWish(ctx context.Context, wishID int64) (int, error)
	ListByWishlist(ctx context.Context, shareCode uuid.UUID) ([]*Reservation, error)
	ListByReserver(ctx context.Context, reserverID int64) ([]*Reservation, error)
}
