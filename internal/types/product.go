// Package types defines shared data models used across the backend.
// Product is the canonical shape for a product record in the database and API responses.
package types

// Product represents a product row in the products table.
type Product struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Category       string  `json:"category"`
	Price          int     `json:"price"`
	CompareAtPrice int     `json:"compare_at_price,omitempty"`
	Rating         float64 `json:"rating"`
	ReviewCount    int     `json:"review_count"`
	SkinType       string  `json:"skin_type"`
	Badge          string  `json:"badge,omitempty"`
	ImageURL       string  `json:"image_url"`
	Description    string  `json:"description"`
	Volume         string  `json:"volume"`
}
