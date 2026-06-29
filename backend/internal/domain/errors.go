package domain

import "errors"

var (
	ErrAccountNotFound  = errors.New("account not found")
	ErrWishlistNotFound = errors.New("wishlist not found")
	ErrWishItemNotFound = errors.New("wish not found")

	ErrReservationNotFound     = errors.New("reservation not found")
	ErrSelfReservation         = errors.New("cannot reserve item from own wishlist")
	ErrReservationExists       = errors.New("reservation already exists")
	ErrReserveQuantityExceeded = errors.New("reserve quantity exceeds available")
)
