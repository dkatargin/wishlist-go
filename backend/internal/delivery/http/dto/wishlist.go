package dto

import (
	"time"
	"wishlist-go/internal/domain"

	"github.com/google/uuid"
)

type CreateWishlistRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Color       *string `json:"color"`
}

type UpdateWishlistRequest struct {
	Name        *string `json:"name" binding:"required"`
	Description *string `json:"description" binding:"required"`
	Color       *string `json:"color"`
}

type CreateWishlistResponse struct {
	domain.Wishlist
}

type UserWishlistsResponse struct {
	domain.Wishlist
}

// WishlistDetailResponse — деталь списка без OwnerID (account_id не светим не-владельцу).
type WishlistDetailResponse struct {
	ShareCode   uuid.UUID `json:"share_code"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Color       *string   `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsOwner     bool      `json:"is_owner"`
}

// SharedWishItem — публично-безопасное представление желания (без OwnerID).
type SharedWishItem struct {
	ID               int64    `json:"id"`
	Name             *string  `json:"name"`
	Priority         int      `json:"priority"`
	MarketURL        string   `json:"market_url"`
	MarketPictureURL *string  `json:"market_picture_url"`
	MarketPrice      *float64 `json:"market_price"`
	MarketCurrency   string   `json:"market_currency"`
	MarketQuantity   *int     `json:"market_quantity"`
}

// SharedWishlistResponse — гостевой просмотр шаренного списка (без OwnerID).
type SharedWishlistResponse struct {
	ShareCode   uuid.UUID        `json:"share_code"`
	Name        string           `json:"name"`
	Description *string          `json:"description"`
	Color       *string          `json:"color"`
	Items       []SharedWishItem `json:"items"`
}
