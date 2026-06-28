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
	IsDone         bool      `gorm:"not null;default:false" json:"is_done"`
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

	err := r.db.Where("wish_list_code = ?", wishlistCode.String()).
		Order("priority DESC, created_at DESC").
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
func (r *wishItemRepo) UpdateWishItem(id int64, shareCode uuid.UUID, upd domain.WishItemUpdate) error {
	// Перевод доменного partial → колонки БД живёт ТОЛЬКО здесь.
	updates := map[string]any{}
	if upd.Name != nil {
		updates["name"] = *upd.Name
	}
	if upd.Priority != nil {
		updates["priority"] = *upd.Priority
	}
	if upd.IsDone != nil {
		updates["is_done"] = *upd.IsDone
	}
	if upd.MarketURL != nil {
		updates["market_link"] = *upd.MarketURL
	}
	if upd.MarketPictureURL != nil {
		updates["market_picture"] = *upd.MarketPictureURL
	}
	if upd.MarketPrice != nil {
		updates["market_price"] = *upd.MarketPrice
	}
	if upd.MarketCurrency != nil {
		updates["market_currency"] = *upd.MarketCurrency
	}
	if upd.MarketQuantity != nil {
		updates["market_quantity"] = *upd.MarketQuantity
	}
	if len(updates) == 0 {
		return nil
	}

	result := r.db.Model(&wishItemModel{}).
		Where("id = ? AND wish_list_code = ?", id, shareCode.String()).
		Updates(updates)

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
	m := &wishItemModel{
		ID:             item.ID,
		WishListCode:   item.WishListCode,
		OwnerID:        item.OwnerID,
		Priority:       item.Priority,
		IsDone:         item.IsDone,
		MarketLink:     item.MarketURL,
		MarketCurrency: item.MarketCurrency,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}
	if item.Name != nil {
		m.Name = *item.Name
	}
	if item.MarketPictureURL != nil {
		m.MarketPicture = *item.MarketPictureURL
	}
	if item.MarketPrice != nil {
		m.MarketPrice = *item.MarketPrice
	}
	if item.MarketQuantity != nil {
		m.MarketQuantity = *item.MarketQuantity
	}
	return m
}

// modelToDomainWishItem конвертирует GORM модель в domain модель
func modelToDomainWishItem(model *wishItemModel) *domain.WishItem {
	name := model.Name
	pic := model.MarketPicture
	price := model.MarketPrice
	qty := model.MarketQuantity
	return &domain.WishItem{
		ID:               model.ID,
		WishListCode:     model.WishListCode,
		OwnerID:          model.OwnerID,
		Name:             &name,
		Priority:         model.Priority,
		IsDone:           model.IsDone,
		MarketURL:        model.MarketLink,
		MarketPictureURL: &pic,
		MarketPrice:      &price,
		MarketCurrency:   model.MarketCurrency,
		MarketQuantity:   &qty,
		CreatedAt:        model.CreatedAt,
		UpdatedAt:        model.UpdatedAt,
	}
}
