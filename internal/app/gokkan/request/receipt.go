package request

// BuyItem WTF should I write here you damn linter?.
type BuyItem struct {
	ItemID   int `json:"item_id" validate:"required,gt=0"`
	Quantity int `json:"quantity" validate:"required,gt=0"`
}

// UpdateReceipt represents receipt update request.
type UpdateReceipt struct {
	ReceiptID int `json:"receipt_id" validate:"required,gt=0"`
	Status    int `json:"status" validate:"required,min=0,max=2"`
}
