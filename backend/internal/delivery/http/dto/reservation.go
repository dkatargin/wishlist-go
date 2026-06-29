package dto

import (
	"time"

	"github.com/google/uuid"
)

// ReserveRequest — тело POST .../wishes/:wishId/reserve.
type ReserveRequest struct {
	Quantity    int  `json:"quantity"`
	IsAnonymous bool `json:"is_anonymous"`
}

// ReservationResponse — резерв в ответе («мои резервы», результат создания).
type ReservationResponse struct {
	ID          uuid.UUID `json:"id"`
	WishID      int64     `json:"wish_id"`
	Quantity    int       `json:"quantity"`
	IsAnonymous bool      `json:"is_anonymous"`
	IsDone      bool      `json:"is_done"`
	CreatedAt   time.Time `json:"created_at"`
}
