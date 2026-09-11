package types

// CartItem represents a product within a user's cart.
type CartItem struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	Image     string `json:"image"`
	Qty       int    `json:"qty"`
}

// WishlistItem represents a product saved by the user.
type WishlistItem struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	Image     string `json:"image"`
}
