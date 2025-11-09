package domain

import (
	"time"

	"github.com/google/uuid"
)

type WishItem struct {
	ID               int64     `json:"id"`
	WishListCode     uuid.UUID `json:"wishlist_id"`
	Name             *string   `json:"name"`
	MarketURL        string    `json:"market_url"`
	MarketPictureURL *string   `json:"market_picture_url"`
	MarketPrice      *float64  `json:"market_price"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WishItemRepository interface {
	CreateWishItem(wishItem *WishItem) error
	GetWishItemByID(id int64, wishlistCode uuid.UUID) (*WishItem, error)
	GetWishItemsByWishlistID(wishlistCode uuid.UUID, limit int, offset int) ([]*WishItem, error)
	UpdateWishItem(wishItem *WishItem) error
	DeleteWishItem(id int64) error
}
