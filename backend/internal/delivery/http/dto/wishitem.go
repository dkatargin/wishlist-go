package dto

type CreateWishItemRequest struct {
	Name           string  `json:"name" binding:"required"`
	Priority       int     `json:"priority"`
	MarketLink     string  `json:"market_link"`
	MarketPicture  string  `json:"market_picture"`
	MarketPrice    float64 `json:"market_price"`
	MarketCurrency string  `json:"market_currency"`
	MarketQuantity int     `json:"market_quantity"`
}

type UpdateWishItemRequest struct {
	Name           *string  `json:"name" binding:"required"`
	Priority       *int     `json:"priority"`
	MarketLink     *string  `json:"market_link"`
	MarketPicture  *string  `json:"market_picture"`
	MarketPrice    *float64 `json:"market_price"`
	MarketCurrency *string  `json:"market_currency"`
	MarketQuantity *int     `json:"market_quantity"`
}
