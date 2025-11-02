package db

import (
	"fmt"
	"log"
	"wishlist-api/internal/config"
	"wishlist-api/internal/db/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ORM *gorm.DB

func dsnFromConfig() string {
	log.Println(config.Config)

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Config.Database.Host, config.Config.Database.Port, config.Config.Database.User,
		config.Config.Database.Password, config.Config.Database.Name,
	)
}

func ConnectDB() (err error) {
	if ORM != nil {
		log.Println("ORM is already initialized")
		return nil
	}
	masterDSN := dsnFromConfig()
	db, err := gorm.Open(postgres.Open(masterDSN), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(
		&models.Account{},
		&models.WishList{},
		&models.WishItem{},
		&models.WishReservation{},
		&models.Migration{},
		&models.WishItemDataRequest{},
	)
	if err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}

	ORM = db
	return nil
}
