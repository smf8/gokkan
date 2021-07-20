package response

// User is the user json response.
type User struct {
	ID             int     `json:"id"`
	Username       string  `json:"username"`
	FullName       string  `json:"full_name,omitempty"`
	BillingAddress string  `json:"billing_address,omitempty"`
	Balance        float64 `json:"balance,omitempty"`
	IsAdmin        bool    `json:"is_admin"`
	Token          string  `json:"token"`
}
