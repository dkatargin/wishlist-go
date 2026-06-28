package postgres

import "gorm.io/gorm"

// AutoMigrate прогоняет схему БД для всех GORM-моделей пакета.
// Модели неэкспортируемые, поэтому миграция живёт рядом с ними.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&accountModel{},
		&wishlistModel{},
		&wishItemModel{},
		&migrationModel{},
	)
}
