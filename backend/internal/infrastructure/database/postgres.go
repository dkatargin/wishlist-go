package database

import (
	"fmt"
	"log"
	"wishlist-go/internal/infrastructure/config"
	pgrepo "wishlist-go/internal/repository/postgres"

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
	// TranslateError: GORM приводит ошибки драйвера к доменным (ErrDuplicatedKey,
	// ErrForeignKeyViolated), чтобы репозитории не парсили SQLSTATE руками.
	dbConnect, err := gorm.Open(postgres.Open(masterDSN), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatal(err)
	}

	if err := pgrepo.AutoMigrate(dbConnect); err != nil {
		log.Fatalf("auto migration failed: %v", err)
	}

	return dbConnect
}
