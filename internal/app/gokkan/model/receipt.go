package model

import "time"

const (
	ReceiptStatusProcessing = iota
	ReceiptStatusDone
	ReceiptStatusCanceled
)

type ReceiptStatus int

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
