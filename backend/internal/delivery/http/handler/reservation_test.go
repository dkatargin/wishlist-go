package handler

import (
	"encoding/json"
	"testing"
	"wishlist-go/internal/domain"
	"wishlist-go/internal/usecase/reservation"
)

// Регрессия на DI: конструктор обязан инжектить usecase, иначе nil-deref в каждом методе.
func TestNewReservationHandler_InjectsUsecase(t *testing.T) {
	h := NewReservationHandler(&reservation.Service{})
	if h.reservationUC == nil {
		t.Fatal("reservation usecase must be injected")
	}
}

// Страж сюрприза: вид владельца (GET /list/:listId/wishes и .../:wishId) отдаёт domain.WishItem.
// Фиксируем точный набор JSON-ключей — ЛЮБОЕ новое поле (reserved_quantity, booked, …) уронит
// тест и заставит осознанно проверить, не утекают ли данные о резервах владельцу.
func TestOwnerWishItemResponseHidesReservations(t *testing.T) {
	qty := 3
	j, err := json.Marshal(domain.WishItem{ID: 1, OwnerID: 100, MarketQuantity: &qty})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(j, &fields); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	allowed := map[string]bool{
		"id": true, "wishlist_id": true, "owner_id": true, "name": true,
		"priority": true, "is_done": true, "market_url": true, "market_picture_url": true,
		"market_price": true, "market_currency": true, "market_quantity": true,
		"created_at": true, "updated_at": true,
	}
	for k := range fields {
		if !allowed[k] {
			t.Fatalf("в owner-представлении WishItem новое поле %q — проверь, не утечка ли это резерва", k)
		}
	}
}
