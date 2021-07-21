package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

const (
	// ReceiptStatusProcessing is for an in-process order.
	ReceiptStatusProcessing = iota
	// ReceiptStatusDone is for a finished order.
	ReceiptStatusDone
	// ReceiptStatusCanceled is for a canceled order.
	ReceiptStatusCanceled
)

// ReceiptStatus is one of ReceiptStatusProcessing, ReceiptStatusDone, ReceiptStatusCanceled.
type ReceiptStatus int

// Receipt is a receipt for an item. it has a `belong-to` relation with model.Item.
type Receipt struct {
	ID             int           `json:"id"`
	ItemID         int           `json:"-"`
	Item           Item          `json:"item,omitempty"`
	Quantity       int           `json:"quantity"`
	UserID         int           `json:"user_id"`
	BuyerName      string        `json:"buyer_name"`
	BillingAddress string        `json:"billing_address"`
	Price          float64       `json:"price"`
	Date           time.Time     `json:"date"`
	TrackingCode   string        `json:"tracking_code"`
	Status         ReceiptStatus `json:"status"`
}

// ReceiptRepo defines allowed operations on receipts.
type ReceiptRepo interface {
	Save(receipt *Receipt) error
	UpdateStatus(receiptID int, status ReceiptStatus) error
	FindForUser(userID int) ([]Receipt, error)
	FindAll() ([]Receipt, error)
}

var _ ReceiptRepo = SQLReceiptRepo{}

// SQLReceiptRepo is the SQL implementation of ReceiptRepo.
type SQLReceiptRepo struct {
	DB *gorm.DB
}

// Save saves given receipt in database.
func (r SQLReceiptRepo) Save(receipt *Receipt) error {
	return r.DB.Save(receipt).Error
}

// UpdateStatus changes the status of a specific receipt.
func (r SQLReceiptRepo) UpdateStatus(receiptID int, status ReceiptStatus) error {
	receipt := Receipt{ID: receiptID}
	if err := r.DB.Find(&receipt).Error; err != nil {
		return err
	}

	receipt.Status = status

	return r.DB.Save(&receipt).Error
}

// FindForUser returns receipts for given user.
func (r SQLReceiptRepo) FindForUser(userID int) ([]Receipt, error) {
	var result []Receipt

	err := r.DB.Joins("Item").Where("user_id = ?", userID).Find(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}

		return nil, err
	}

	if len(result) == 0 {
		return nil, ErrRecordNotFound
	}

	return result, nil
}

// FindAll retrieves all receipts from Database.
func (r SQLReceiptRepo) FindAll() ([]Receipt, error) {
	var result []Receipt

	err := r.DB.Joins("Item").Find(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}

		return nil, err
	}

	if len(result) == 0 {
		return nil, ErrRecordNotFound
	}

	return result, nil
}
