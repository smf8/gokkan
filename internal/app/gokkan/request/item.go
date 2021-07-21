package request

// Item represents a product create/update request.
type Item struct {
	Name       string  `json:"name" validate:"required"`
	CategoryID int     `json:"category_id" validate:"required, gt=0"`
	Price      float64 `json:"price" validate:"required, gt=0"`
}
