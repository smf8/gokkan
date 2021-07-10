package config

//Namespace is the name for application instance
const Namespace = "Gokkan"

//nolint:lll,gochecknoglobals,gomnd
var def = Config{
	Logger: Logger{
		Level:   5,
		Enabled: true,
	},
}
