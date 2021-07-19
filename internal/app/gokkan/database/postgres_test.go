package database_test

import (
	"testing"

	"github.com/smf8/gokkan/internal/app/gokkan/config"
	"github.com/smf8/gokkan/internal/app/gokkan/database"
	"github.com/stretchr/testify/assert"
)

func TestConnectDatabase(t *testing.T) {
	t.Parallel()

	cfg := config.New()
	a := assert.New(t)
	db, err := database.New(cfg.Database)

	a.NoError(err)

	// test ping database
	database, err := db.DB()

	a.NoError(err)

	a.NoError(database.Ping())
}
