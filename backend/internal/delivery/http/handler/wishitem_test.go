package handler

import (
	"testing"
	"wishlist-go/internal/usecase/wishitem"
	"wishlist-go/internal/usecase/wishlist"
)

// Регрессионный тест на DI-дыру: конструктор обязан инжектить обе usecase-зависимости.
// Без wishlistUsecase каждый метод хендлера падает с nil-pointer на CheckAccess.
func TestNewWishItemHandler_InjectsBothUsecases(t *testing.T) {
	h := NewWishItemHandler(&wishitem.Service{}, &wishlist.Service{})

	if h.usecase == nil {
		t.Fatal("wishitem usecase must be injected")
	}
	if h.wishlistUsecase == nil {
		t.Fatal("wishlist usecase must be injected (иначе nil-deref в CheckAccess)")
	}
}
