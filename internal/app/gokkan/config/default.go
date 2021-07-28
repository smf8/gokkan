package config

import (
	"time"

	"github.com/pyroscope-io/pyroscope/pkg/agent/profiler"
	"github.com/pyroscope-io/pyroscope/pkg/agent/spy"
)

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
	Pyroscope: Pyroscope{
		Enable:     true,
		Server:     "http://localhost:4040",
		SampleRate: 100,
		EnableLogs: true,
		Profiles: []spy.ProfileType{
			profiler.ProfileCPU,
			profiler.ProfileAllocObjects,
			profiler.ProfileAllocSpace,
			profiler.ProfileInuseObjects,
			profiler.ProfileInuseSpace,
		},
	},
}
