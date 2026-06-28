package handler

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"wishlist-go/internal/delivery/http/dto"
	"wishlist-go/internal/delivery/http/middleware"
	"wishlist-go/internal/domain"
	"wishlist-go/internal/usecase/reservation"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReservationHandler — резервирование подарков залогиненным дарителем (под Telegram-auth).
type ReservationHandler struct {
	reservationUC *reservation.Service
}

func NewReservationHandler(resUC *reservation.Service) *ReservationHandler {
	return &ReservationHandler{reservationUC: resUC}
}

// Reserve — POST /share/:shareCode/wishes/:wishId/reserve.
func (h *ReservationHandler) Reserve(c *gin.Context) {
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	shareCode, err := uuid.Parse(c.Param("shareCode"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid share code"})
		return
	}
	wishID, err := strconv.ParseInt(c.Param("wishId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wish id"})
		return
	}

	// тело необязательно: пустой POST = quantity по умолчанию, не анонимно
	var req dto.ReserveRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reserverID := auth.(*middleware.TelegramAuthData).User.ID
	res, err := h.reservationUC.Reserve(c.Request.Context(), reserverID, shareCode, wishID, req.Quantity, req.IsAnonymous)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWishItemNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "wish item not found"})
		case errors.Is(err, domain.ErrSelfReservation):
			c.JSON(http.StatusForbidden, gin.H{"error": "cannot reserve item from own wishlist"})
		case errors.Is(err, domain.ErrReservationExists):
			c.JSON(http.StatusConflict, gin.H{"error": "reservation already exists"})
		case errors.Is(err, domain.ErrReserveQuantityExceeded):
			c.JSON(http.StatusConflict, gin.H{"error": "reserve quantity exceeds available"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reserve"})
		}
		return
	}

	c.JSON(http.StatusCreated, toReservationResponse(res))
}

// Cancel — DELETE /share/:shareCode/wishes/:wishId/reserve (снимает только свой резерв).
func (h *ReservationHandler) Cancel(c *gin.Context) {
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if _, err := uuid.Parse(c.Param("shareCode")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid share code"})
		return
	}
	wishID, err := strconv.ParseInt(c.Param("wishId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wish id"})
		return
	}

	reserverID := auth.(*middleware.TelegramAuthData).User.ID
	if err := h.reservationUC.CancelByWish(c.Request.Context(), reserverID, wishID); err != nil {
		if errors.Is(err, domain.ErrReservationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "reservation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel reservation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// MarkPurchased — POST /reservations/:reservationId/purchased (только автор резерва).
func (h *ReservationHandler) MarkPurchased(c *gin.Context) {
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	reservationID, err := uuid.Parse(c.Param("reservationId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reservation id"})
		return
	}

	reserverID := auth.(*middleware.TelegramAuthData).User.ID
	if err := h.reservationUC.MarkPurchased(c.Request.Context(), reserverID, reservationID); err != nil {
		// чужой резерв не находится по скоупу reserver_id → ErrReservationNotFound
		if errors.Is(err, domain.ErrReservationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "reservation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark purchased"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// ListMine — GET /reservations (резервы текущего пользователя по всем спискам).
func (h *ReservationHandler) ListMine(c *gin.Context) {
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	reserverID := auth.(*middleware.TelegramAuthData).User.ID
	list, err := h.reservationUC.ListMine(c.Request.Context(), reserverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch reservations"})
		return
	}

	out := make([]dto.ReservationResponse, 0, len(list))
	for _, r := range list {
		out = append(out, toReservationResponse(r))
	}
	c.JSON(http.StatusOK, gin.H{"reservations": out})
}

func toReservationResponse(r *domain.Reservation) dto.ReservationResponse {
	return dto.ReservationResponse{
		ID:          r.ID,
		WishID:      r.WishID,
		Quantity:    r.Quantity,
		IsAnonymous: r.IsAnonymous,
		IsDone:      r.IsDone,
		CreatedAt:   r.CreatedAt,
	}
}
