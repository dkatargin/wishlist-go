package postgres

import (
	"context"
	"errors"
	"time"
	"wishlist-go/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type wishlistModel struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OwnerID     int64     `gorm:"index;not null" json:"owner"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"not null" json:"description"`
	ShareCode   uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"share_code"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Owner accountModel `json:"-" gorm:"foreignKey:OwnerID"`
}

func (wishlistModel) TableName() string {
	return "wish_lists"
}

type wishlistRepo struct {
	db *gorm.DB
}

func NewWishlistRepository(db *gorm.DB) domain.WishlistRepository {
	return &wishlistRepo{db: db}
}

func (r *wishlistRepo) CreateWishlist(ctx context.Context, wishlist *domain.Wishlist) error {
	model := wlDomainToModel(wishlist)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}

	// Обновляем domain модель сгенерированными значениями
	wishlist.ID = model.ID
	wishlist.Name = model.Name
	wishlist.Description = &model.Description
	wishlist.OwnerID = model.OwnerID
	wishlist.Color = nil
	wishlist.ShareCode = model.ShareCode
	wishlist.CreatedAt = model.CreatedAt
	wishlist.UpdatedAt = model.UpdatedAt

	return nil
}

// GetWishlistByID получает вишлист по shareCode
func (r *wishlistRepo) GetWishlistByCode(ctx context.Context, shareCode uuid.UUID) (*domain.Wishlist, error) {
	var model wishlistModel

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).First(&model, "share_code = ?", shareCode.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrWishlistNotFound
		}
		return nil, err
	}

	return wlModelToDomain(&model), nil
}

// GetWishlistsByOwnerID получает все вишлисты пользователя
func (r *wishlistRepo) GetWishlistsByOwnerID(ctx context.Context, ownerID int64, offset int, limit int) ([]*domain.Wishlist, error) {
	var models []wishlistModel

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.db.WithContext(ctx).
		Where("owner_id = ?", ownerID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}

	wishlists := make([]*domain.Wishlist, 0, len(models))
	for i := range models {
		wishlists = append(wishlists, wlModelToDomain(&models[i]))
	}

	return wishlists, nil
}

// UpdateWishlist обновляет вишлист
func (r *wishlistRepo) UpdateWishlist(ctx context.Context, wishlist *domain.Wishlist) error {
	model := wlDomainToModel(wishlist)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updatedAt := time.Now()

	updates := map[string]interface{}{}
	if model.Name != "" {
		updates["name"] = model.Name
	}
	if model.Description != "" {
		updates["description"] = model.Description
	}
	updates["updated_at"] = updatedAt

	result := r.db.WithContext(ctx).
		Model(&wishlistModel{}).
		Where("share_code = ?", wishlist.ShareCode.String()).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrWishlistNotFound
	}

	wishlist.Name = model.Name
	wishlist.Description = &model.Description
	wishlist.UpdatedAt = updatedAt

	return nil
}

// DeleteWishlist удаляет вишлист
func (r *wishlistRepo) DeleteWishlist(ctx context.Context, shareCode uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := r.db.WithContext(ctx).Delete(&wishlistModel{}, "share_code = ?", shareCode.String())

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrWishlistNotFound
	}

	return nil
}

// domainToModel конвертирует domain модель в GORM модель
func wlDomainToModel(wl *domain.Wishlist) *wishlistModel {
	return &wishlistModel{
		ID:          wl.ID,
		OwnerID:     wl.OwnerID,
		ShareCode:   wl.ShareCode,
		Name:        wl.Name,
		Description: *wl.Description,
		CreatedAt:   wl.CreatedAt,
		UpdatedAt:   wl.UpdatedAt,
	}
}

// modelToDomain конвертирует GORM модель в domain модель
func wlModelToDomain(m *wishlistModel) *domain.Wishlist {

	return &domain.Wishlist{
		ID:          m.ID,
		OwnerID:     m.OwnerID,
		ShareCode:   m.ShareCode,
		Name:        m.Name,
		Description: &m.Description,
		Color:       nil,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
