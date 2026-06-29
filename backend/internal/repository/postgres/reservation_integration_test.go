package postgres

import (
	"context"
	"errors"
	"testing"
	"wishlist-go/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// seedWish создаёт владельца, список и желание (для FK), возвращает id желания и share-код.
func seedWish(t *testing.T, db *gorm.DB, ownerID int64, marketQuantity int) (wishID int64, shareCode uuid.UUID) {
	t.Helper()
	if _, err := NewAccountRepository(db).CreateAccount(context.Background(), &ownerID); err != nil {
		t.Fatalf("CreateAccount: %v", err)
	}
	shareCode = uuid.New()
	if err := db.Create(&wishlistModel{OwnerID: ownerID, Name: "l", ShareCode: shareCode}).Error; err != nil {
		t.Fatalf("wishlist: %v", err)
	}
	wi := &wishItemModel{
		WishListCode: shareCode, OwnerID: ownerID, Name: "gift",
		MarketLink: "m", MarketPicture: "p", MarketPrice: 1, MarketCurrency: "RUB", MarketQuantity: marketQuantity,
	}
	if err := db.Create(wi).Error; err != nil {
		t.Fatalf("wishitem: %v", err)
	}
	return wi.ID, shareCode
}

func TestReservationRepo_CRUDUniqueAndSum(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()
	repo := NewReservationRepository(db)
	wishID, _ := seedWish(t, db, 100, 3)

	var reserver int64 = 200
	res := &domain.Reservation{WishID: wishID, ReserverID: reserver, Quantity: 2}
	if err := repo.Create(ctx, res); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if res.ID == uuid.Nil {
		t.Fatal("Create должен проставить ID")
	}

	// unique(reserver, wish): повтор тем же дарителем отвергается БД и транслируется
	// репозиторием в доменную ошибку (чтобы гонка дублей давала 409, а не 500)
	if err := repo.Create(ctx, &domain.Reservation{WishID: wishID, ReserverID: reserver, Quantity: 1}); !errors.Is(err, domain.ErrReservationExists) {
		t.Fatalf("повторный резерв тем же дарителем: ожидался ErrReservationExists, got %v", err)
	}

	got, err := repo.GetByReserverAndWish(ctx, reserver, wishID)
	if err != nil || got.Quantity != 2 {
		t.Fatalf("GetByReserverAndWish: %+v, %v", got, err)
	}

	// второй даритель резервирует ещё 1 → сумма по желанию = 3
	if err := repo.Create(ctx, &domain.Reservation{WishID: wishID, ReserverID: 300, Quantity: 1}); err != nil {
		t.Fatalf("Create#2: %v", err)
	}
	if sum, err := repo.SumQuantityByWish(ctx, wishID); err != nil || sum != 3 {
		t.Fatalf("SumQuantityByWish = %d, %v; ожидали 3", sum, err)
	}

	if mine, err := repo.ListByReserver(ctx, reserver); err != nil || len(mine) != 1 {
		t.Fatalf("ListByReserver = %d, %v; ожидали 1", len(mine), err)
	}

	// Delete скоупится по дарителю: чужой не удаляет
	if err := repo.Delete(ctx, res.ID, 999); !errors.Is(err, domain.ErrReservationNotFound) {
		t.Fatalf("Delete чужим: ожидался ErrReservationNotFound, got %v", err)
	}
	if err := repo.Delete(ctx, res.ID, reserver); err != nil {
		t.Fatalf("Delete своим: %v", err)
	}
	if _, err := repo.GetByReserverAndWish(ctx, reserver, wishID); !errors.Is(err, domain.ErrReservationNotFound) {
		t.Fatalf("после удаления ожидался ErrReservationNotFound, got %v", err)
	}
}

func TestReservationRepo_ListByWishlistAndMarkDone(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()
	repo := NewReservationRepository(db)
	wishID, shareCode := seedWish(t, db, 100, 5)

	res := &domain.Reservation{WishID: wishID, ReserverID: 200, Quantity: 1}
	if err := repo.Create(ctx, res); err != nil {
		t.Fatalf("Create: %v", err)
	}

	// join reservations → wish_items по share-коду
	list, err := repo.ListByWishlist(ctx, shareCode)
	if err != nil || len(list) != 1 || list[0].WishID != wishID || list[0].ID != res.ID {
		t.Fatalf("ListByWishlist = %+v, %v", list, err)
	}

	// MarkDone скоупится по дарителю
	if err := repo.MarkDone(ctx, res.ID, 999); !errors.Is(err, domain.ErrReservationNotFound) {
		t.Fatalf("MarkDone чужим: ожидался ErrReservationNotFound, got %v", err)
	}
	if err := repo.MarkDone(ctx, res.ID, 200); err != nil {
		t.Fatalf("MarkDone своим: %v", err)
	}
	if got, _ := repo.GetByReserverAndWish(ctx, 200, wishID); !got.IsDone {
		t.Fatal("MarkDone должен выставить is_done=true")
	}
}

func TestReservationRepo_CascadeOnWishDelete(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()
	repo := NewReservationRepository(db)
	wishID, _ := seedWish(t, db, 100, 3)

	if err := repo.Create(ctx, &domain.Reservation{WishID: wishID, ReserverID: 200, Quantity: 1}); err != nil {
		t.Fatalf("Create: %v", err)
	}

	// hard-delete желания не должен падать на FK и должен каскадно снести его резервы
	// (иначе владелец не сможет удалить зарезервированное желание — получит 500).
	if err := db.Delete(&wishItemModel{}, wishID).Error; err != nil {
		t.Fatalf("удаление желания с резервом не должно падать на FK: %v", err)
	}
	if sum, err := repo.SumQuantityByWish(ctx, wishID); err != nil || sum != 0 {
		t.Fatalf("резервы желания должны исчезнуть каскадом: sum=%d, err=%v", sum, err)
	}
}
