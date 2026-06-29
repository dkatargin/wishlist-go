package dto

import (
	"time"

	"github.com/google/uuid"
)

// FavoriteListItem — карточка списка в избранном. Без account_id владельца:
// избранное — чужие списки, не светим больше, чем отдаёт публичный share.
type FavoriteListItem struct {
	ID          int64     `json:"id"`
	ShareCode   uuid.UUID `json:"share_code"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Color       *string   `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
