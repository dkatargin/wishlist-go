package postgres

import (
	"context"
	"testing"

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

	// вставляем элемент напрямую моделью — изолируем баг колонки от конвертеров (P1).
	if err := db.Create(&wishItemModel{
		WishListCode:   listCode,
		OwnerID:        owner,
		Name:           "item",
		Priority:       1,
		Status:         "pending",
		MarketLink:     "http://market",
		MarketPicture:  "http://pic",
		MarketPrice:    10,
		MarketCurrency: "RUB",
		MarketQuantity: 1,
	}).Error; err != nil {
		t.Fatalf("insert wish item: %v", err)
	}

	items, err := NewWishItemRepository(db).GetWishItemsByWishlistID(listCode, 50, 0)
	if err != nil {
		t.Fatalf("GetWishItemsByWishlistID: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("ожидался 1 элемент по share-коду, получили %d", len(items))
	}
}
