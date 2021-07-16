package model

import "gorm.io/gorm"

// Category specifies an item category
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CategoryRepo defines allowed operations on Category
type CategoryRepo interface {
	Delete(name string) error
	Save(category *Category) error
	FindAll() ([]Category, error)
}

// make sure SQLCategoryRepo implements Category at compile time
var _ CategoryRepo = SQLCategoryRepo{}

// SQLCategoryRepo is SQL implementation of CategoryRepo
type SQLCategoryRepo struct {
	DB *gorm.DB
}

// Delete removes a category with given name.
func (c SQLCategoryRepo) Delete(name string) error {
	return c.DB.Where("name = ?", name).Delete(&Category{}).Error
}

// Save saves given category. if it contains an ID it's updated
func (c SQLCategoryRepo) Save(category *Category) error {
	return c.DB.Save(category).Error
}

// FindAll retrieves all categories.
func (c SQLCategoryRepo) FindAll() ([]Category, error) {
	var result []Category
	err := c.DB.Find(&result).Error

	return result, err
}
