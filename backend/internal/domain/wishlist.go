package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Wishlist struct {
	ID          int64     `json:"id"`
	OwnerID     int64     `json:"account_id"`
	ShareCode   uuid.UUID `json:"share_code"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Color       *string   `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WishlistRepository interface {
	CreateWishlist(ctx context.Context, wishlist *Wishlist) error
	GetWishlistByCode(ctx context.Context, shareCode uuid.UUID) (*Wishlist, error)
	GetWishlistsByOwnerID(ctx context.Context, ownerID int64, offset int, limit int) ([]*Wishlist, error)
	UpdateWishlist(ctx context.Context, wishlist *Wishlist) error
	DeleteWishlist(ctx context.Context, shareCode uuid.UUID) error
}
