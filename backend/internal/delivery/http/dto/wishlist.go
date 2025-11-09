package dto

import "wishlist-go/internal/domain"

type CreateWishlistRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Color       *string `json:"color"`
}

type UpdateWishlistRequest struct {
	Name        *string `json:"name" binding:"required"`
	Description *string `json:"description" binding:"required"`
	Color       *string `json:"color"`
}

type CreateWishlistResponse struct {
	domain.Wishlist
}

type UserWishlistsResponse struct {
	domain.Wishlist
}
