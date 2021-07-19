package database

import (
	"fmt"

	"github.com/smf8/gokkan/internal/app/gokkan/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DSN creates a Data Source Name for connecting to postgresql.
func DSN(cfg config.Database) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Name)
}

// New creates a new database connection for given user.
func New(cfg config.Database) (*gorm.DB, error) {
	// it's better to try and connect to database with retry.
	// https://github.com/avast/retry-go might be a good option.
	db, err := gorm.Open(postgres.Open(DSN(cfg)))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
