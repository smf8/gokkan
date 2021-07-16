package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// ErrItemNotFound the purpose of this variable is pretty obvious.
var ErrItemNotFound = errors.New("item not found")

// Item represents an item in our website, it has a `belong-to` relation with model.Category
type Item struct {
	ID         int       `json:"id"`
	CategoryID int       `json:"category_id"`
	Category   Category  `json:"category"`
	Price      float64   `json:"price"`
	Remaining  int       `json:"remaining"`
	Sold       int       `json:"sold"`
	PhotoURL   string    `json:"photo_url"`
	CreatedAt  time.Time `json:"created_at"`
}

type findOption struct {
	minPrice        float64
	maxPrice        float64
	limit           int
	categoryID      int
	descendingOrder bool
	orderPrice      bool
	orderDate       bool
}

// ItemRepo defines allowed operations on an item in database.
type ItemRepo interface {
	Find(options ...ItemOption) ([]Item, error)
	Save(item *Item) error
}

var _ ItemRepo = SQLItemRepo{}

// SQLItemRepo is SQL implementation of ItemRepo.
type SQLItemRepo struct {
	DB *gorm.DB
}

// ItemOption is a function which accepts a findOption and changes something
// in it. read https://www.calhoun.io/using-functional-options-instead-of-method-chaining-in-go/
// to learn more about functional options.
type ItemOption func(option *findOption)

// we use something like functional options
//for providing dynamic selection cases

// WithPriceRange sets price minimum and maximum in select statement.
func WithPriceRange(minPrice, maxPrice float64) ItemOption {
	return func(option *findOption) {
		option.minPrice = minPrice
		option.maxPrice = maxPrice
	}
}

// WithLimit sets number of records to fetch.
func WithLimit(limit int) ItemOption {
	return func(option *findOption) {
		option.limit = limit
	}
}

// WithCategory specifies the category to select for
func WithCategory(categoryID int) ItemOption {
	return func(option *findOption) {
		option.categoryID = categoryID
	}
}

// WithDescendingOrder sets order of results to descending order
func WithDescendingOrder() ItemOption {
	return func(option *findOption) {
		option.descendingOrder = true
	}
}

// WithPriceOrder sets order by price
func WithPriceOrder() ItemOption {
	return func(option *findOption) {
		option.orderPrice = true
		// we set this in case of overriding
		option.orderDate = false
	}
}

// WithDateOrder sets order by date
func WithDateOrder() ItemOption {
	return func(option *findOption) {
		option.orderPrice = false
		option.orderDate = true
	}
}

func (o findOption) toSQL(tx *gorm.DB) *gorm.DB {
	if o.maxPrice != 0 {
		tx = tx.Where("price < ? and price > ?", o.maxPrice, o.minPrice)
	}

	if o.orderPrice {
		if o.descendingOrder {
			tx = tx.Order("price desc")
		} else {
			tx = tx.Order("price ")
		}
	} else {
		if o.descendingOrder {
			tx = tx.Order("created_at desc")
		} else {
			tx = tx.Or("created_at")
		}
	}

	if o.categoryID != 0 {
		tx = tx.Where("category_id = ?", o.categoryID)
	}

	if o.limit != 0 {
		tx = tx.Limit(o.limit)
	}

	return tx
}

// Save saves an item in database, if Item.ID is 0 it will
// insert it as a new record.
func (i SQLItemRepo) Save(item *Item) error {
	return i.DB.Save(item).Error
}

// Find creates a query to find items with given parametes.
// see ItemOption and findOption for details.
func (i SQLItemRepo) Find(options ...ItemOption) ([]Item, error) {
	// apply default item options
	defaultOptions := &findOption{}

	for _, option := range options {
		option(defaultOptions)
	}

	queryTx := defaultOptions.toSQL(i.DB)

	var result []Item

	if err := queryTx.Find(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrItemNotFound
		}

		return nil, err
	}

	return result, nil
}
