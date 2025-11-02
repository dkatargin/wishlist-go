package service

import (
	"wishlist-go/internal/db"
	"wishlist-go/internal/db/models"

	"github.com/google/uuid"
)

type WishItemService struct {
	WishList       uuid.UUID
	Owner          *int64
	Name           *string
	Priority       *int
	Status         *string
	MarketLink     *string
	MarketPicture  *string
	MarketPrice    *float64
	MarketCurrency *string
	MarketQuantity *int
}

func NewWishItemService() *WishItemService {
	return &WishItemService{}
}

func (s *WishItemService) GetAll(limit int, offset int) ([]models.WishItem, error) {
	var wishItems []models.WishItem
	err := db.ORM.Model(&models.WishItem{}).Where("wish_list_code = ?", s.WishList).Limit(limit).Offset(offset).Find(&wishItems).Error
	return wishItems, err
}

func (s *WishItemService) Get(id int64) (*models.WishItem, error) {
	var wishItem *models.WishItem
	err := db.ORM.Model(&models.WishItem{}).Where("id = ? AND wish_list_code = ?", id, s.WishList).First(&wishItem).Error
	return wishItem, err
}

func (s *WishItemService) Create() (*models.WishItem, error) {
	wishItem := &models.WishItem{
		WishListCode:   s.WishList,
		OwnerID:        *s.Owner,
		Name:           *s.Name,
		Priority:       *s.Priority,
		Status:         *s.Status,
		MarketLink:     *s.MarketLink,
		MarketPicture:  *s.MarketPicture,
		MarketPrice:    *s.MarketPrice,
		MarketCurrency: *s.MarketCurrency,
		MarketQuantity: *s.MarketQuantity,
	}

	err := db.ORM.Model(&models.WishItem{}).Create(wishItem).Error
	if err != nil {
		return nil, err
	}

	return s.Get(wishItem.ID)
}

func (s *WishItemService) Update(id int64) (*models.WishItem, error) {
	updates := make(map[string]interface{})

	if s.Name != nil {
		updates["name"] = *s.Name
	}
	if s.Priority != nil {
		updates["priority"] = *s.Priority
	}
	if s.Status != nil {
		updates["status"] = *s.Status
	}
	if s.MarketLink != nil {
		updates["market_link"] = *s.MarketLink
	}
	if s.MarketPicture != nil {
		updates["market_picture"] = *s.MarketPicture
	}
	if s.MarketPrice != nil {
		updates["market_price"] = *s.MarketPrice
	}
	if s.MarketCurrency != nil {
		updates["market_currency"] = *s.MarketCurrency
	}
	if s.MarketQuantity != nil {
		updates["market_quantity"] = *s.MarketQuantity
	}

	if len(updates) == 0 {
		return s.Get(id)
	}

	err := db.ORM.Model(&models.WishItem{}).Where("id = ? AND wish_list_code = ?", id, s.WishList).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	return s.Get(id)

}

func (s *WishItemService) Delete(id int64) error {
	return db.ORM.Model(&models.WishItem{}).Where("id = ? AND wish_list_code = ?", id, s.WishList).Delete(&models.WishItem{}).Error
}
