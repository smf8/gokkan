package request

// Login represents Login request body.
type Login struct {
	Username string `json:"username" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

// Signup represents user signup request body.
type Signup struct {
	Username       string `json:"username" validate:"required,email,max=255"`
	Password       string `json:"password" validate:"required,min=8,max=255"`
	FullName       string `json:"full_name" validate:"max=255"`
	BillingAddress string `json:"billing_address" validate:"max=1000"`
}

// ChargeBalance represents balance increase request.
type ChargeBalance struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

// UpdateUser represents an update user request
type UpdateUser struct {
	Password       string `json:"password" validate:"max=255"`
	FullName       string `json:"full_name" validate:"max=255"`
	BillingAddress string `json:"billing_address" validate:"max=1000"`
}
