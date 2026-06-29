package postgres

import (
	"context"
	"errors"
	"testing"
	"wishlist-go/internal/domain"

	"github.com/google/uuid"
)

func TestFavoriteRepo_AddIdempotentListRemove(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()
	repo := NewFavoriteRepository(db)

	// владелец + два его списка
	var owner int64 = 100
	if _, err := NewAccountRepository(db).CreateAccount(ctx, &owner); err != nil {
		t.Fatalf("CreateAccount owner: %v", err)
	}
	code1, code2 := uuid.New(), uuid.New()
	for name, code := range map[string]uuid.UUID{"A": code1, "B": code2} {
		if err := db.Create(&wishlistModel{OwnerID: owner, Name: name, Description: "d", ShareCode: code}).Error; err != nil {
			t.Fatalf("wishlist %s: %v", name, err)
		}
	}

	// тот, кто добавляет в избранное
	var acc int64 = 200
	if _, err := NewAccountRepository(db).CreateAccount(ctx, &acc); err != nil {
		t.Fatalf("CreateAccount acc: %v", err)
	}

	if err := repo.Add(ctx, &domain.Favorite{AccountID: acc, WishlistCode: code1}); err != nil {
		t.Fatalf("Add#1: %v", err)
	}
	if err := repo.Add(ctx, &domain.Favorite{AccountID: acc, WishlistCode: code2}); err != nil {
		t.Fatalf("Add#2: %v", err)
	}
	// идемпотентность: повтор того же списка тем же юзером — не ошибка и без дубля
	if err := repo.Add(ctx, &domain.Favorite{AccountID: acc, WishlistCode: code1}); err != nil {
		t.Fatalf("повторный Add не должен быть ошибкой: %v", err)
	}

	lists, err := repo.ListByAccount(ctx, acc)
	if err != nil || len(lists) != 2 {
		t.Fatalf("ListByAccount = %d, %v; ожидали 2 (без дубля)", len(lists), err)
	}
	// вернулись карточки (имя/share_code), а не голые id
	names := map[string]bool{}
	for _, w := range lists {
		names[w.Name] = true
		if w.ShareCode == uuid.Nil {
			t.Fatalf("share_code не заполнен: %+v", w)
		}
	}
	if !names["A"] || !names["B"] {
		t.Fatalf("ожидали карточки A и B: %+v", lists)
	}

	if err := repo.Remove(ctx, acc, code1); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if l, _ := repo.ListByAccount(ctx, acc); len(l) != 1 {
		t.Fatalf("после удаления ожидали 1, got %d", len(l))
	}

	// удаление отсутствующего → ErrFavoriteNotFound
	if err := repo.Remove(ctx, acc, uuid.New()); !errors.Is(err, domain.ErrFavoriteNotFound) {
		t.Fatalf("удаление отсутствующего: ожидался ErrFavoriteNotFound, got %v", err)
	}
}
