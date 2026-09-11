package types

import "time"

// Order represents a purchased order created from the authenticated user's cart.
type Order struct {
	ID        int64       `json:"id"`
	Status    string      `json:"status"`
	Total     int         `json:"total"`
	CreatedAt time.Time   `json:"created_at"`
	Items     []OrderItem `json:"items"`
}

// OrderItem stores the product snapshot as it existed at the time of purchase.
type OrderItem struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	Qty       int    `json:"qty"`
}
