package handler

import (
	"errors"
	"net/http"
	"wishlist-go/internal/delivery/http/dto"
	"wishlist-go/internal/delivery/http/middleware"
	"wishlist-go/internal/domain"
	"wishlist-go/internal/usecase/favorite"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// FavoriteHandler — избранные (сохранённые) чужие списки текущего пользователя.
type FavoriteHandler struct {
	usecase *favorite.Service
}

func NewFavoriteHandler(uc *favorite.Service) *FavoriteHandler {
	return &FavoriteHandler{usecase: uc}
}

// Add — POST /wishlist/:id/favorite (:id = ShareCode).
func (h *FavoriteHandler) Add(c *gin.Context) {
	accountID, ok := requireAccount(c)
	if !ok {
		return
	}
	shareCode, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wishlist id"})
		return
	}

	if err := h.usecase.Add(c.Request.Context(), accountID, shareCode); err != nil {
		if errors.Is(err, domain.ErrWishlistNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "wishlist not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add favorite"})
		return
	}
	// Состояние списка фронт перечитывает сам (FetchFavorites) — не делаем лишний
	// SELECT и не рискуем ложным 500 после успешной записи.
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// Remove — DELETE /wishlist/:id/favorite.
func (h *FavoriteHandler) Remove(c *gin.Context) {
	accountID, ok := requireAccount(c)
	if !ok {
		return
	}
	shareCode, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wishlist id"})
		return
	}

	if err := h.usecase.Remove(c.Request.Context(), accountID, shareCode); err != nil {
		if errors.Is(err, domain.ErrFavoriteNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "favorite not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove favorite"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// List — GET /favorites.
func (h *FavoriteHandler) List(c *gin.Context) {
	accountID, ok := requireAccount(c)
	if !ok {
		return
	}
	h.respondList(c, accountID)
}

// respondList отдаёт текущее избранное пользователя как []FavoriteListItem.
func (h *FavoriteHandler) respondList(c *gin.Context, accountID int64) {
	lists, err := h.usecase.List(c.Request.Context(), accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch favorites"})
		return
	}
	out := make([]dto.FavoriteListItem, 0, len(lists))
	for _, w := range lists {
		out = append(out, dto.FavoriteListItem{
			ID:          w.ID,
			ShareCode:   w.ShareCode,
			Name:        w.Name,
			Description: w.Description,
			Color:       w.Color,
			CreatedAt:   w.CreatedAt,
			UpdatedAt:   w.UpdatedAt,
		})
	}
	c.JSON(http.StatusOK, out)
}

// requireAccount достаёт id аккаунта из telegram_auth; пишет 401 и возвращает false, если нет.
func requireAccount(c *gin.Context) (int64, bool) {
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return 0, false
	}
	return auth.(*middleware.TelegramAuthData).User.ID, true
}
