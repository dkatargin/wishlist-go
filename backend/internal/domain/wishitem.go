package domain

import (
	"time"

	"github.com/google/uuid"
)

type WishItem struct {
	ID               int64     `json:"id"`
	WishListCode     uuid.UUID `json:"wishlist_id"`
	OwnerID          int64     `json:"owner_id"`
	Name             *string   `json:"name"`
	Priority         int       `json:"priority"`
	IsDone           bool      `json:"is_done"`
	MarketURL        string    `json:"market_url"`
	MarketPictureURL *string   `json:"market_picture_url"`
	MarketPrice      *float64  `json:"market_price"`
	MarketCurrency   string    `json:"market_currency"`
	MarketQuantity   *int      `json:"market_quantity"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WishItemUpdate — частичное обновление (PATCH): применяются только не-nil поля.
// Доменное представление, без знания имён колонок БД (их знает только repository).
type WishItemUpdate struct {
	Name             *string
	Priority         *int
	IsDone           *bool
	MarketURL        *string
	MarketPictureURL *string
	MarketPrice      *float64
	MarketCurrency   *string
	MarketQuantity   *int
}

type WishItemRepository interface {
	CreateWishItem(wishItem *WishItem) error
	GetWishItemByID(id int64, wishlistCode uuid.UUID) (*WishItem, error)
	GetWishItemsByWishlistID(wishlistCode uuid.UUID, limit int, offset int, onlyActive bool) ([]*WishItem, error)
	UpdateWishItem(id int64, shareCode uuid.UUID, upd WishItemUpdate) error
	DeleteWishItem(id int64) error
}
