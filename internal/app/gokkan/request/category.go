package request

// CreateCategory represents a create category request.
type CreateCategory struct {
	Name string `json:"name"`
}

// DeleteCategory represents a delete category request.
// ID should be greater than 1 in order to prevent
// deleting default category.
type DeleteCategory struct {
	ID int `param:"id" validate:"required,gt=1"`
}
