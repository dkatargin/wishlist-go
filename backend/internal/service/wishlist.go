package service

import (
	"wishlist-api/internal/db"
	"wishlist-api/internal/db/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WishlistService struct {
	orm *gorm.DB
}

func NewWishlistService() *WishlistService {
	return &WishlistService{orm: db.ORM}
}

type WishlistInsert struct {
	Name        *string
	Description *string
	Owner       *int64
}

func (s *WishlistService) GetAllByOwner(ownerTelegramId int64, limit int, offset int) ([]models.WishList, error) {
	var wishlists []models.WishList
	err := db.ORM.Model(&models.WishList{}).Debug().Where("owner_id = ?", ownerTelegramId).Limit(limit).Offset(offset).Find(&wishlists).Error
	return wishlists, err
}

func (s *WishlistService) Get(uuid uuid.UUID) (*models.WishList, error) {
	var wishlist *models.WishList
	err := db.ORM.Model(&models.WishList{}).Where("share_code = ?", uuid).First(&wishlist).Error
	return wishlist, err
}

func (s *WishlistService) Create(insert *WishlistInsert) (*models.WishList, error) {
	// генерируем уникальный share_code
	shareCode := uuid.New()
	err := db.ORM.Model(&models.WishList{}).Create(&models.WishList{
		Name:        *insert.Name,
		Description: *insert.Description,
		OwnerID:     *insert.Owner,
		ShareCode:   shareCode,
	}).Error
	if err != nil {
		return nil, err
	}
	return s.Get(shareCode)
}

func (s *WishlistService) Update(shareCode uuid.UUID, patch WishlistInsert) (*models.WishList, error) {
	updates := make(map[string]interface{})

	if patch.Name != nil {
		updates["name"] = patch.Name
	}
	if patch.Description != nil {
		updates["description"] = patch.Description
	}

	if len(updates) == 0 {
		return s.Get(shareCode)
	}

	err := db.ORM.Model(&models.WishList{}).Where("share_code = ?", shareCode).Updates(updates).Error
	if err != nil {
		return nil, err
	}
	return s.Get(shareCode)
}

func (s *WishlistService) Delete(shareCode uuid.UUID) error {
	return db.ORM.Model(&models.WishList{}).Where("share_code = ?", shareCode).Delete(&models.WishList{}).Error
}

func (s *WishlistService) CheckAccess(ownerTelegramId int64, shareCode uuid.UUID) bool {
	var wishlistExists bool
	err := db.ORM.Model(&models.WishList{}).
		Select("count(*) > 0").
		Where("owner_id = ? AND share_code = ?", ownerTelegramId, shareCode).
		Find(&wishlistExists).Error
	if err != nil {
		return false
	}
	return wishlistExists
}
