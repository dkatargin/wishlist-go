package postgres

import (
	"context"
	"time"
	"wishlist-go/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type favoriteModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	AccountID    int64     `gorm:"not null;uniqueIndex:uq_favorite_account_wishlist"`
	WishlistCode uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uq_favorite_account_wishlist"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`

	// FK на share_code списка; при удалении списка его избранное чистится каскадом.
	Wishlist wishlistModel `json:"-" gorm:"foreignKey:WishlistCode;references:ShareCode;constraint:OnDelete:CASCADE"`
}

func (favoriteModel) TableName() string { return "favorites" }

type favoriteRepo struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) domain.FavoriteRepository {
	return &favoriteRepo{db: db}
}

// Add идемпотентен: дубль по unique(account, wishlist) тихо игнорируется.
func (r *favoriteRepo) Add(ctx context.Context, f *domain.Favorite) error {
	m := &favoriteModel{
		ID:           f.ID,
		AccountID:    f.AccountID,
		WishlistCode: f.WishlistCode,
		CreatedAt:    f.CreatedAt,
	}
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(m).Error
}

func (r *favoriteRepo) Remove(ctx context.Context, accountID int64, wishlistCode uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("account_id = ? AND wishlist_code = ?", accountID, wishlistCode).
		Delete(&favoriteModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrFavoriteNotFound
	}
	return nil
}

// ListByAccount джойнит избранное к спискам и возвращает сами карточки (новые сверху).
func (r *favoriteRepo) ListByAccount(ctx context.Context, accountID int64) ([]*domain.Wishlist, error) {
	var models []wishlistModel
	err := r.db.WithContext(ctx).
		Joins("JOIN favorites ON favorites.wishlist_code = wish_lists.share_code").
		Where("favorites.account_id = ?", accountID).
		Order("favorites.created_at DESC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	out := make([]*domain.Wishlist, 0, len(models))
	for i := range models {
		out = append(out, wlModelToDomain(&models[i]))
	}
	return out, nil
}
