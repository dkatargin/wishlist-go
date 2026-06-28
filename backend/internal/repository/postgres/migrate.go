package postgres

import "gorm.io/gorm"

// AutoMigrate прогоняет схему БД для всех GORM-моделей пакета.
// Модели неэкспортируемые, поэтому миграция живёт рядом с ними.
func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&accountModel{},
		&wishlistModel{},
		&wishItemModel{},
		&migrationModel{},
	); err != nil {
		return err
	}

	// Колонка status осталась от старой схемы (заменена на is_done); AutoMigrate её
	// не удаляет сам, а NOT NULL без default ломал бы INSERT по новой модели.
	if db.Migrator().HasColumn(&wishItemModel{}, "status") {
		if err := db.Migrator().DropColumn(&wishItemModel{}, "status"); err != nil {
			return err
		}
	}

	return nil
}
