// Package types defines shared data models used across the backend.
// Product is the canonical shape for a product record in the database and API responses.
package types

// Product represents a product row in the products table.
type Product struct {
	ID             int     `json:"id"`                         // ID uniquely identifies the product.
	Name           string  `json:"name"`                       // Name is the product's display name.
	Category       string  `json:"category"`                   // Category groups related products.
	Price          int     `json:"price"`                      // Price is the current product price.
	CompareAtPrice int     `json:"compare_at_price,omitempty"` // CompareAtPrice is the original price when discounted.
	Rating         float64 `json:"rating"`                     // Rating is the product's average review score.
	ReviewCount    int     `json:"review_count"`               // ReviewCount is the number of product reviews.
	SkinType       string  `json:"skin_type"`                  // SkinType identifies the intended skin type.
	Badge          string  `json:"badge,omitempty"`            // Badge is an optional merchandising label.
	ImageURL       string  `json:"image_url"`                  // ImageURL points to the product image.
	Description    string  `json:"description"`                // Description explains the product.
	Volume         string  `json:"volume"`                     // Volume is the package size.
}
