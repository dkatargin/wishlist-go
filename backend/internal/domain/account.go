package domain

import (
	"context"
	"time"
)

type Account struct {
	ID               int64     `json:"id"`
	RegistrationType string    `json:"registration_type"` // e.g., "telegram", "web"
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type AccountRepository interface {
	CreateAccount(ctx context.Context, telegramId *int64) (int64, error)
	GetAccountByID(id int64) (*Account, error)
	DeleteAccount(id int64) error
}
