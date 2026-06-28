package postgres

import (
	"context"
	"errors"
	"testing"
	"wishlist-go/internal/domain"

	"github.com/google/uuid"
)

func TestWishItemRepo_GetByWishlist_FiltersByListCode(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	// владелец + список нужны для FK
	var owner int64 = 777
	if _, err := NewAccountRepository(db).CreateAccount(ctx, &owner); err != nil {
		t.Fatalf("CreateAccount: %v", err)
	}

	listCode := uuid.New()
	if err := db.Create(&wishlistModel{OwnerID: owner, Name: "list", ShareCode: listCode}).Error; err != nil {
		t.Fatalf("insert wishlist: %v", err)
	}

	// вставляем элемент напрямую моделью — изолируем баг колонки от конвертеров.
	if err := db.Create(&wishItemModel{
		WishListCode:   listCode,
		OwnerID:        owner,
		Name:           "item",
		Priority:       1,
		MarketLink:     "http://market",
		MarketPicture:  "http://pic",
		MarketPrice:    10,
		MarketCurrency: "RUB",
		MarketQuantity: 1,
	}).Error; err != nil {
		t.Fatalf("insert wish item: %v", err)
	}

	items, err := NewWishItemRepository(db).GetWishItemsByWishlistID(listCode, 50, 0, false)
	if err != nil {
		t.Fatalf("GetWishItemsByWishlistID: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("ожидался 1 элемент по share-коду, получили %d", len(items))
	}
}

func TestWishItemRepo_CreatePersistsFieldsAndSortsByPriority(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	var owner int64 = 555
	if _, err := NewAccountRepository(db).CreateAccount(ctx, &owner); err != nil {
		t.Fatalf("CreateAccount: %v", err)
	}
	listCode := uuid.New()
	if err := db.Create(&wishlistModel{OwnerID: owner, Name: "list", ShareCode: listCode}).Error; err != nil {
		t.Fatalf("insert wishlist: %v", err)
	}

	repo := NewWishItemRepository(db)
	mk := func(name string, prio int) *domain.WishItem {
		n, pic, price, qty := name, "http://pic", 1.0, 2
		return &domain.WishItem{
			WishListCode: listCode, OwnerID: owner, Name: &n, Priority: prio,
			MarketURL: "http://m", MarketPictureURL: &pic, MarketPrice: &price,
			MarketCurrency: "RUB", MarketQuantity: &qty,
		}
	}
	if err := repo.CreateWishItem(mk("low", 1)); err != nil {
		t.Fatalf("create low: %v", err)
	}
	if err := repo.CreateWishItem(mk("high", 9)); err != nil {
		t.Fatalf("create high: %v", err)
	}

	items, err := repo.GetWishItemsByWishlistID(listCode, 50, 0, false)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("ожидалось 2 элемента, получили %d", len(items))
	}
	if items[0].Priority != 9 || items[1].Priority != 1 {
		t.Fatalf("сортировка priority DESC нарушена: [%d,%d]", items[0].Priority, items[1].Priority)
	}
	if items[0].MarketCurrency != "RUB" || items[0].MarketQuantity == nil || *items[0].MarketQuantity != 2 {
		t.Fatalf("новые поля не сохранились: %+v", items[0])
	}
}

func TestWishItemRepo_UpdatePartial_PersistsZeroValues(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	var owner int64 = 333
	if _, err := NewAccountRepository(db).CreateAccount(ctx, &owner); err != nil {
		t.Fatalf("CreateAccount: %v", err)
	}
	listCode := uuid.New()
	if err := db.Create(&wishlistModel{OwnerID: owner, Name: "l", ShareCode: listCode}).Error; err != nil {
		t.Fatalf("insert wishlist: %v", err)
	}

	repo := NewWishItemRepository(db)
	n, pic, price, qty := "x", "p", 1.0, 1
	item := &domain.WishItem{
		WishListCode: listCode, OwnerID: owner, Name: &n, Priority: 5, IsDone: true,
		MarketURL: "m", MarketPictureURL: &pic, MarketPrice: &price, MarketCurrency: "RUB", MarketQuantity: &qty,
	}
	if err := repo.CreateWishItem(item); err != nil {
		t.Fatalf("create: %v", err)
	}

	// PATCH: нулевые значения должны сохраниться (чего .Updates(struct) не умеет).
	no, zero := false, 0
	if err := repo.UpdateWishItem(item.ID, listCode, domain.WishItemUpdate{IsDone: &no, Priority: &zero}); err != nil {
		t.Fatalf("update: %v", err)
	}

	got, err := repo.GetWishItemByID(item.ID, listCode)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.IsDone {
		t.Fatalf("is_done не сброшен в false")
	}
	if got.Priority != 0 {
		t.Fatalf("priority не сброшен в 0: %d", got.Priority)
	}
}

func TestWishItemRepo_UpdateScopedToList(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	var owner int64 = 444
	if _, err := NewAccountRepository(db).CreateAccount(ctx, &owner); err != nil {
		t.Fatalf("CreateAccount: %v", err)
	}
	listA, listB := uuid.New(), uuid.New()
	if err := db.Create(&wishlistModel{OwnerID: owner, Name: "la", ShareCode: listA}).Error; err != nil {
		t.Fatalf("wishlist A: %v", err)
	}
	if err := db.Create(&wishlistModel{OwnerID: owner, Name: "lb", ShareCode: listB}).Error; err != nil {
		t.Fatalf("wishlist B: %v", err)
	}

	repo := NewWishItemRepository(db)
	n, pic, price, qty := "x", "p", 1.0, 1
	item := &domain.WishItem{
		WishListCode: listA, OwnerID: owner, Name: &n,
		MarketURL: "m", MarketPictureURL: &pic, MarketPrice: &price, MarketCurrency: "RUB", MarketQuantity: &qty,
	}
	if err := repo.CreateWishItem(item); err != nil {
		t.Fatalf("create: %v", err)
	}

	// item в списке A; апдейт с кодом списка B → не найдено (скоуп по wish_list_code).
	hacked := "hacked"
	err := repo.UpdateWishItem(item.ID, listB, domain.WishItemUpdate{Name: &hacked})
	if !errors.Is(err, domain.ErrWishItemNotFound) {
		t.Fatalf("ожидался ErrWishItemNotFound при чужом списке, получили %v", err)
	}
}

func TestWishItemRepo_GetByWishlist_OnlyActiveFiltersDone(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	var owner int64 = 666
	if _, err := NewAccountRepository(db).CreateAccount(ctx, &owner); err != nil {
		t.Fatalf("CreateAccount: %v", err)
	}
	listCode := uuid.New()
	if err := db.Create(&wishlistModel{OwnerID: owner, Name: "l", ShareCode: listCode}).Error; err != nil {
		t.Fatalf("insert wishlist: %v", err)
	}
	for _, it := range []wishItemModel{
		{WishListCode: listCode, OwnerID: owner, Name: "done", IsDone: true, MarketLink: "m", MarketPicture: "p", MarketPrice: 1, MarketCurrency: "RUB", MarketQuantity: 1},
		{WishListCode: listCode, OwnerID: owner, Name: "active", IsDone: false, MarketLink: "m", MarketPicture: "p", MarketPrice: 1, MarketCurrency: "RUB", MarketQuantity: 1},
	} {
		if err := db.Create(&it).Error; err != nil {
			t.Fatalf("insert item: %v", err)
		}
	}

	repo := NewWishItemRepository(db)

	active, err := repo.GetWishItemsByWishlistID(listCode, 50, 0, true)
	if err != nil {
		t.Fatalf("list active: %v", err)
	}
	if len(active) != 1 || active[0].Name == nil || *active[0].Name != "active" {
		t.Fatalf("onlyActive=true: ожидался 1 невыполненный 'active', получили %d", len(active))
	}

	all, err := repo.GetWishItemsByWishlistID(listCode, 50, 0, false)
	if err != nil {
		t.Fatalf("list all: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("onlyActive=false: ожидалось 2, получили %d", len(all))
	}
}
