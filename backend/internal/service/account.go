package service

import (
	"wishlist-go/internal/db"
	"wishlist-go/internal/db/models"
)

type AccountService struct{}

func NewAccountService() *AccountService {
	return &AccountService{}
}

func (s *AccountService) Get(telegramId int64) (*models.Account, error) {
	var account *models.Account
	err := db.ORM.Model(&models.Account{}).Where("id = ?", telegramId).First(&account).Error
	return account, err
}

func (s *AccountService) Create(telegramId int64) (*models.Account, error) {
	err := db.ORM.Model(&models.Account{}).Create(&models.Account{ID: telegramId}).Error
	if err != nil {
		return nil, err
	}
	return s.Get(telegramId)
}

func (s *AccountService) Delete(telegramId int64) error {
	return db.ORM.Model(&models.Account{}).Delete(&models.Account{ID: telegramId}).Error
}
