package model

// User represents a normal user in database
type User struct {
	ID             int     `json:"id"`
	Username       string  `json:"username"`
	Password       string  `json:"password"`
	FullName       string  `json:"full_name"`
	BillingAddress string  `json:"billing_address"`
	Balance        float64 `json:"balance"`
}
