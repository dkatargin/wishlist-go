package handler

import (
	"testing"
	"wishlist-go/internal/usecase/favorite"
)

// Регрессия на DI: конструктор обязан инжектить usecase, иначе nil-deref в каждом методе.
func TestNewFavoriteHandler_InjectsUsecase(t *testing.T) {
	h := NewFavoriteHandler(&favorite.Service{})
	if h.usecase == nil {
		t.Fatal("favorite usecase must be injected")
	}
}
