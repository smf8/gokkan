package request

// CreateItem represents a product create/update request.
type CreateItem struct {
	Name       string  `json:"name" validate:"required"`
	CategoryID int     `json:"category_id" validate:"required,gt=0"`
	Price      float64 `json:"price" validate:"required,gt=0"`
}

// ItemFilter represents a find item filter request.
type ItemFilter struct {
	CategoryID   *int     `query:"category_id"`
	MinPrice     *float64 `query:"min_price"`
	MaxPrice     *float64 `query:"max_price"`
	CreatedSince *int64   `query:"created_since"`
	SortByPrice  *bool    `query:"sort_by_price"`
	SortDesc     *bool    `query:"sort_desc"`
}
