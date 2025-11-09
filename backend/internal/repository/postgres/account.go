package postgres

import (
	"context"
	"errors"
	"time"
	"wishlist-go/internal/domain"

	"gorm.io/gorm"
)

// accountModel - GORM модель для таблицы accounts
type accountModel struct {
	ID        int64     `gorm:"primaryKey" json:"id"` // telegram-id
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (accountModel) TableName() string {
	return "accounts"
}

// accountRepo - реализация domain.AccountRepository
type accountRepo struct {
	db *gorm.DB
}

// CreateAccount создает новый аккаунт
func (r *accountRepo) CreateAccount(ctx context.Context, telegramId *int64) (int64, error) {
	model := &accountModel{
		ID: *telegramId,
	}
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return 0, err
	}
	return model.ID, nil
}

// GetAccountByID получает аккаунт по ID
func (r *accountRepo) GetAccountByID(id int64) (*domain.Account, error) {
	var model accountModel

	if err := r.db.First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrAccountNotFound
		}
		return nil, err
	}

	return modelToDomain(&model), nil
}

func (r *accountRepo) DeleteAccount(id int64) error {
	//TODO implement me
	panic("implement me")
}

// NewAccountRepository создает новый репозиторий для работы с аккаунтами
func NewAccountRepository(db *gorm.DB) domain.AccountRepository {
	return &accountRepo{db: db}
}

// modelToDomain конвертирует GORM модель в domain модель
func modelToDomain(m *accountModel) *domain.Account {
	return &domain.Account{
		ID:               m.ID,
		RegistrationType: "telegram",
		IsActive:         true,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}
