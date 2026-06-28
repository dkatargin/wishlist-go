package handler

import (
	"testing"
	"wishlist-go/internal/delivery/http/dto"
	"wishlist-go/internal/usecase/reservation"
)

func TestApplyReservationView(t *testing.T) {
	mq := 5
	agg := reservation.WishReservationView{ReservedQuantity: 3, Reservers: []int64{7}}

	t.Run("владелец не видит резервы вообще", func(t *testing.T) {
		si := dto.SharedWishItem{MarketQuantity: &mq}
		applyReservationView(&si, true /*isOwner*/, true /*authenticated*/, agg)
		if si.ReservedQuantity != 0 || si.Remaining != nil || len(si.Reservers) != 0 {
			t.Fatalf("владельцу резервы должны быть скрыты: %+v", si)
		}
	})

	t.Run("неаутентифицированный видит агрегат, но НЕ дарителей", func(t *testing.T) {
		si := dto.SharedWishItem{MarketQuantity: &mq}
		applyReservationView(&si, false, false, agg)
		if si.ReservedQuantity != 3 || si.Remaining == nil || *si.Remaining != 2 {
			t.Fatalf("агрегат количества должен быть виден: %+v", si)
		}
		if si.Reservers == nil || len(si.Reservers) != 0 {
			t.Fatalf("состав дарителей анониму не раскрывается и не должен быть null: %v", si.Reservers)
		}
	})

	t.Run("аутентифицированный гость видит не-анонимных дарителей", func(t *testing.T) {
		si := dto.SharedWishItem{MarketQuantity: &mq}
		applyReservationView(&si, false, true, agg)
		if si.ReservedQuantity != 3 || si.Remaining == nil || *si.Remaining != 2 {
			t.Fatalf("агрегат количества должен быть виден: %+v", si)
		}
		if len(si.Reservers) != 1 || si.Reservers[0] != 7 {
			t.Fatalf("гость видит не-анонимных дарителей: %v", si.Reservers)
		}
	})
}

func TestRemainingQuantity(t *testing.T) {
	ptr := func(n int) *int { return &n }
	cases := []struct {
		name     string
		mq       *int
		reserved int
		want     *int
	}{
		{"market_quantity не задан (nil) → nil", nil, 2, nil},
		{"market_quantity = 0 (не указано) → nil", ptr(0), 1, nil},
		{"остаток = market_quantity - reserved", ptr(5), 2, ptr(3)},
		{"перебронь не уходит в минус → 0", ptr(2), 5, ptr(0)},
		{"ничего не зарезервировано → весь market_quantity", ptr(4), 0, ptr(4)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := remainingQuantity(tc.mq, tc.reserved)
			if (got == nil) != (tc.want == nil) {
				t.Fatalf("nil-несовпадение: got %v, want %v", got, tc.want)
			}
			if got != nil && *got != *tc.want {
				t.Fatalf("got %d, want %d", *got, *tc.want)
			}
		})
	}
}
