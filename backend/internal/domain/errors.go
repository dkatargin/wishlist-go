package domain

import "errors"

var (
	ErrAccountNotFound  = errors.New("account not found")
	ErrWishlistNotFound = errors.New("wishlist not found")
	ErrWishItemNotFound = errors.New("wish not found")
)
