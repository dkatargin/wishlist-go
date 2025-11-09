package database

import (
	"fmt"
	"log"
	"wishlist-go/internal/infrastructure/config"
	"wishlist-go/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	db  *gorm.DB
	cfg *config.DB
}

func PostgresDSN(host string, port int, user string, password string, dbname string) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
}

func ConnectDB(cfg *config.DB) *gorm.DB {

	masterDSN := PostgresDSN(cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
	dbConnect, err := gorm.Open(postgres.Open(masterDSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = dbConnect.AutoMigrate(
		&models.Account{},
		&models.WishList{},
		&models.WishItem{},
		&models.WishReservation{},
		&models.Migration{},
		&models.WishItemDataRequest{},
	)
	if err != nil {
		log.Fatalf("auto migration failed: %w", err)
	}

	return dbConnect
}
