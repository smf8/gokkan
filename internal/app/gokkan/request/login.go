package request

// Login represents Login request body.
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

// Add request validation in the future
