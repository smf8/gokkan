package model

import "gorm.io/gorm"

// User represents a normal user in database.
type User struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	Password       string    `json:"password"`
	FullName       string    `json:"full_name"`
	BillingAddress string    `json:"billing_address"`
	Balance        float64   `json:"balance"`
	Receipts       []Receipt `json:"-"`
}

// UserRepo defins allowed operations on User.
type UserRepo interface {
	Save(user *User) error
	Find(username string) (User, error)
}

// SQLUserRepo is SQL implementation of UserRepo.
type SQLUserRepo struct {
	DB *gorm.DB
}

// Save saves user in database. if User.ID id not 0 it's record is updated.
func (u SQLUserRepo) Save(user *User) error {
	return u.DB.Save(user).Error
}

// Find finds user with given username.
func (u SQLUserRepo) Find(username string) (*User, error) {
	var user User

	err := u.DB.Where("username = ?", username).Find(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, err
}
