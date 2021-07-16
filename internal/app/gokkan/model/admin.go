package model

import (
	"errors"

	"gorm.io/gorm"
)

// ErrUserNotFound tells that user with given data was not found.
var ErrUserNotFound = errors.New("user not found")

// Admin is an admin user. and it's a read-only entity.
type Admin struct {
	ID       int    `json:"id" gorm:"->"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// AdminRepo is an interface for managing admin model.
type AdminRepo interface {
	Find(username string) (*Admin, error)
}

// make sure SQLAdminRepo implements AdminRepo at compile time
var _ AdminRepo = SQLAdminRepo{}

// SQLAdminRepo is SQL implementation of AdminRepo.
type SQLAdminRepo struct {
	DB *gorm.DB
}

// Find finds admin with given username.
func (a SQLAdminRepo) Find(username string) (*Admin, error) {
	var admin *Admin
	if err := a.DB.Where("username=?", username).
		First(admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return admin, nil
}
