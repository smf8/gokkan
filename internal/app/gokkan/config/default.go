package config

import "time"

// Namespace is the name for application instance.
const Namespace = "Gokkan"

//nolint:lll,gomnd,gochecknoglobals
var def = Config{
	Logger: Logger{
		Level:   5,
		Enabled: true,
	},
	Database: Database{
		Host:     "localhost",
		Port:     "5432",
		Name:     "gokkan",
		Username: "gokkan",
		Password: "1",
		Timeout:  5 * time.Second,
	},
	Server: Server{
		Timeout: 5 * time.Second,
		Secret:  "super_top_classified_secret",
		Port:    8080,
	},
}
