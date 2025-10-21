package models

type Account struct {
	ID        int64 `gorm:"primaryKey" json:"id"` // telegram-id
	CreatedAt int64 `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64 `gorm:"autoUpdateTime" json:"updated_at"`
}
