package request

// Login represents Login request body.
type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	IsAdmin  bool   `json:"is_admin"`
}

// Signup represents user signup request body.
type Signup struct {
	Username       string `json:"username" validate:"required"`
	Password       string `json:"password" validate:"required"`
	FullName       string `json:"full_name"`
	BillingAddress string `json:"billing_address"`
}
