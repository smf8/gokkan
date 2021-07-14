package model

import "time"

const (
	// ReceiptStatusProcessing is for an in-process order
	ReceiptStatusProcessing = iota
	// ReceiptStatusDone is for a finished order
	ReceiptStatusDone
	// ReceiptStatusCanceled is for a canceled order
	ReceiptStatusCanceled
)

// ReceiptStatus is one of ReceiptStatusProcessing, ReceiptStatusDone, ReceiptStatusCanceled
type ReceiptStatus int

// Receipt is a receipt for an item. it has a `belong-to` relation with model.Item
type Receipt struct {
	ID             int           `json:"id"`
	ItemID         int           `json:"-"`
	Item           Item          `json:"item"`
	Quantity       int           `json:"quantity"`
	BuyerName      string        `json:"buyer_name"`
	BillingAddress string        `json:"billing_address"`
	Price          float64       `json:"price"`
	Date           time.Time     `json:"date"`
	TrackingCode   string        `json:"tracking_code"`
	Status         ReceiptStatus `json:"status"`
}
