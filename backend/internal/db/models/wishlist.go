package models

import "github.com/google/uuid"

type WishList struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OwnerID     int64     `gorm:"index;not null" json:"owner"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"not null" json:"description"`
	ShareCode   uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"share_code"`
	CreatedAt   int64     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   int64     `gorm:"autoUpdateTime" json:"updated_at"`

	Owner Account `json:"-" gorm:"foreignKey:OwnerID"`
}

type WishItem struct {
	ID             int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	WishListID     int64   `gorm:"index;not null" json:"wishlist_id"`
	OwnerID        int64   `gorm:"index;not null" json:"owner_id"`
	Name           string  `gorm:"not null" json:"name"`
	Priority       int     `gorm:"not null" json:"priority"`
	Status         string  `gorm:"not null" json:"status"` // возможные значения: "pending", "reserved", "purchased"
	MarketLink     string  `gorm:"not null" json:"market_link"`
	MarketPicture  string  `gorm:"not null" json:"market_picture"`
	MarketPrice    float64 `gorm:"not null" json:"market_price"`
	MarketCurrency string  `gorm:"not null" json:"market_currency"`
	MarketQuantity int     `gorm:"not null" json:"market_quantity"`
	CreatedAt      int64   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      int64   `gorm:"autoUpdateTime" json:"updated_at"`

	Owner    Account  `json:"-" gorm:"foreignKey:OwnerID"`
	WishList WishList `json:"-" gorm:"foreignKey:WishListID"`
}

type WishReservation struct {
	ID         int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	WishID     int64 `gorm:"index;not null" json:"wish_id"`
	ReserverID int64 `gorm:"index;not null" json:"reserver"`
	CreatedAt  int64 `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  int64 `gorm:"autoUpdateTime" json:"updated_at"`

	Wish     WishItem `json:"-" gorm:"foreignKey:WishID"`
	Reserver Account  `json:"-" gorm:"foreignKey:ReserverID"`
}
