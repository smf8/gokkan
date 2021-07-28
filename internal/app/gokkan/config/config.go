package config

import (
	"strings"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/pyroscope-io/pyroscope/pkg/agent/spy"
	"github.com/sirupsen/logrus"
)

const _Prefix = "GOKKAN_"

type (
	// Config represents a gokkan config instance.
	Config struct {
		Logger    Logger    `koanf:"logger"`
		Database  Database  `koanf:"database"`
		Server    Server    `koanf:"server"`
		Pyroscope Pyroscope `koanf:"pyroscope"`
	}

	// Server struct specifies echo server settings.
	Server struct {
		Timeout time.Duration `koanf:"timeout"`
		Secret  string        `koanf:"secret"`
		Port    int           `koanf:"port"`
	}

	// Logger represents logger(logrus) config information.
	Logger struct {
		Level   logrus.Level `koanf:"level"`
		Enabled bool         `koanf:"enabled"`
	}

	// Database is PostgreSQL configuration.
	Database struct {
		Host     string        `koanf:"host"`
		Port     string        `koanf:"port"`
		Name     string        `koanf:"name"`
		Username string        `koanf:"username"`
		Password string        `koanf:"password"`
		Timeout  time.Duration `koanf:"timeout"`
	}

	// Pyroscope represents configuations required to run pyroscope agent.
	Pyroscope struct {
		Server     string `koanf:"server"`
		SampleRate uint32 `koanf:"sample-rate"`
		EnableLogs bool   `koanf:"enable-logs"`
		Enable     bool   `koanf:"enable"`
		// possible values: "cpu, inuse_objects, alloc_objects, inuse_space, alloc_space"
		Profiles []spy.ProfileType `koanf:"profiles"`
	}
)

// New creates a new config instance with this order : default -> config.yml.
func New() Config {
	var instance Config

	k := koanf.New(".")

	if err := k.Load(structs.Provider(def, "koanf"), nil); err != nil {
		logrus.Fatalf("error loading default: %s", err)
	}

	if err := k.Load(file.Provider("config.yml"), yaml.Parser()); err != nil {
		logrus.Errorf("error loading file: %s", err)
	}

	if err := k.Load(env.Provider(_Prefix, ".", func(s string) string {
		parsedEnv := strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(s, _Prefix)), "__", "-")

		return strings.ReplaceAll(parsedEnv, "_", ".")
	}), nil); err != nil {
		logrus.Errorf("error loading environment variables: %s", err)
	}

	if err := k.Unmarshal("", &instance); err != nil {
		logrus.Fatalf("error unmarshalling config: %s", err)
	}

	logrus.Infof("following configuration is loaded:\n%+v", instance)

	return instance
}
