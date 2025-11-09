package postgres

import (
	"errors"
	"time"
	"wishlist-go/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type wishItemModel struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	WishListCode   uuid.UUID `gorm:"type:uuid; index;not null" json:"wishlist_code"`
	OwnerID        int64     `gorm:"index;not null" json:"owner_id"`
	Name           string    `gorm:"not null" json:"name"`
	Priority       int       `gorm:"not null" json:"priority"`
	Status         string    `gorm:"not null" json:"status"` // возможные значения: "pending", "reserved", "purchased"
	MarketLink     string    `gorm:"not null" json:"market_link"`
	MarketPicture  string    `gorm:"not null" json:"market_picture"`
	MarketPrice    float64   `gorm:"not null" json:"market_price"`
	MarketCurrency string    `gorm:"not null" json:"market_currency"`
	MarketQuantity int       `gorm:"not null" json:"market_quantity"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Owner    accountModel  `json:"-" gorm:"foreignKey:OwnerID"`
	WishList wishlistModel `json:"-" gorm:"foreignKey:WishListCode;references:ShareCode"`
}

func (wishItemModel) TableName() string {
	return "wish_items"
}

// wishItemRepo - реализация domain.WishItemRepository
type wishItemRepo struct {
	db *gorm.DB
}

// NewWishItemRepository создает новый экземпляр репозитория
func NewWishItemRepository(db *gorm.DB) domain.WishItemRepository {
	return &wishItemRepo{db: db}
}

// CreateWishItem создает новый элемент вишлиста
func (r *wishItemRepo) CreateWishItem(wishItem *domain.WishItem) error {
	model := domainWishItemToModel(wishItem)

	if err := r.db.Create(&model).Error; err != nil {
		return err
	}

	*wishItem = *modelToDomainWishItem(model)
	return nil
}

// GetWishItemByID получает элемент по ID
func (r *wishItemRepo) GetWishItemByID(id int64, wishlistCode uuid.UUID) (*domain.WishItem, error) {
	var model wishItemModel

	err := r.db.Where("id = ? AND wish_list_code = ?", id, wishlistCode.String()).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrWishItemNotFound
		}
		return nil, err
	}

	return modelToDomainWishItem(&model), nil
}

// GetWishItemsByWishlistID получает все элементы вишлиста
func (r *wishItemRepo) GetWishItemsByWishlistID(wishlistCode uuid.UUID, limit int, offset int) ([]*domain.WishItem, error) {
	var models []wishItemModel

	err := r.db.Where("wish_list_id = ?", wishlistCode.String()).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	items := make([]*domain.WishItem, 0, len(models))
	for _, model := range models {
		items = append(items, modelToDomainWishItem(&model))
	}

	return items, nil
}

// UpdateWishItem обновляет элемент вишлиста
func (r *wishItemRepo) UpdateWishItem(wishItem *domain.WishItem) error {
	model := domainWishItemToModel(wishItem)

	result := r.db.Model(&wishItemModel{}).
		Where("id = ?", wishItem.ID).
		Updates(model)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrWishItemNotFound
	}

	return nil
}

// DeleteWishItem удаляет элемент вишлиста
func (r *wishItemRepo) DeleteWishItem(id int64) error {
	result := r.db.Where("id = ?", id).Delete(&wishItemModel{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrWishItemNotFound
	}

	return nil
}

// domainWishItemToModel конвертирует domain модель в GORM модель
func domainWishItemToModel(item *domain.WishItem) *wishItemModel {
	return &wishItemModel{
		ID:            item.ID,
		WishListCode:  item.WishListCode,
		Name:          *item.Name,
		MarketLink:    item.MarketURL,
		MarketPicture: *item.MarketPictureURL,
		MarketPrice:   *item.MarketPrice,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
	}
}

// modelToDomainWishItem конвертирует GORM модель в domain модель
func modelToDomainWishItem(model *wishItemModel) *domain.WishItem {
	return &domain.WishItem{
		ID:               model.ID,
		WishListCode:     model.WishListCode,
		Name:             &model.Name,
		MarketURL:        model.MarketLink,
		MarketPictureURL: &model.MarketPicture,
		MarketPrice:      &model.MarketPrice,
		CreatedAt:        model.CreatedAt,
		UpdatedAt:        model.UpdatedAt,
	}
}
