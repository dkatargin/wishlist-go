package postgres

import "testing"

// Переход Status(string)→IsDone(bool): AutoMigrate сам НЕ дропает старую NOT NULL
// колонку status, из-за чего INSERT на уже существующей БД падал бы. Проверяем,
// что AutoMigrate её убирает.
func TestAutoMigrate_DropsLegacyStatusColumn(t *testing.T) {
	db := newTestDB(t) // уже прогнал AutoMigrate; в актуальной схеме status нет

	// Эмулируем старую схему: NOT NULL колонка status без поля в модели.
	if err := db.Exec(`ALTER TABLE wish_items ADD COLUMN status varchar NOT NULL DEFAULT 'pending'`).Error; err != nil {
		t.Fatalf("add legacy status column: %v", err)
	}
	if !db.Migrator().HasColumn(&wishItemModel{}, "status") {
		t.Fatal("precondition: колонка status должна существовать")
	}

	if err := AutoMigrate(db); err != nil {
		t.Fatalf("AutoMigrate: %v", err)
	}

	if db.Migrator().HasColumn(&wishItemModel{}, "status") {
		t.Fatal("колонка status должна быть удалена AutoMigrate")
	}
}
