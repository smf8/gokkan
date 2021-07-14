package model

// Admin is an admin user. and it's a read-only entity
type Admin struct {
	ID       int    `json:"id" gorm:"->"`
	Username string `json:"username"`
	Password string `json:"password"`
}
