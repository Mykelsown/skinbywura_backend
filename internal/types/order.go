package types

import "time"

// Order represents a purchased order created from the authenticated user's cart.
type Order struct {
	ID          int64       `json:"id"`
	Status      string      `json:"status"`
	Total       int         `json:"total"`
	CreatedAt   time.Time   `json:"created_at"`
	ShipName    string      `json:"ship_name"`
	ShipPhone   string      `json:"ship_phone"`
	ShipAddress string      `json:"ship_address"`
	ShipCity    string      `json:"ship_city"`
	ShipState   string      `json:"ship_state"`
	Notes       string      `json:"notes"`
	Items       []OrderItem `json:"items"`
}

// OrderItem stores the product snapshot as it existed at the time of purchase.
type OrderItem struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	Qty       int    `json:"qty"`
}

// CheckoutRequest carries the delivery details collected at checkout.
type CheckoutRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Notes   string `json:"notes"`
}
