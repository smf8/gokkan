package profiling

import (
	"github.com/pyroscope-io/pyroscope/pkg/agent/profiler"
	"github.com/sirupsen/logrus"
	"github.com/smf8/gokkan/internal/app/gokkan/config"
)

// Start starts the pyroscope profiler.
func Start(cfg config.Pyroscope) (*profiler.Profiler, error) {
	config := profiler.Config{
		ApplicationName: config.Namespace,
		ServerAddress:   cfg.Server,
		SampleRate:      cfg.SampleRate,
		ProfileTypes:    cfg.Profiles,
	}

	if cfg.EnableLogs {
		config.Logger = logrus.StandardLogger()
	}

	return profiler.Start(config)
}
