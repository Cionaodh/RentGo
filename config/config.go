package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		App  App
		HTTP HTTP
		Log  Log
		PG   PG
		HDFS HDFSConfig
	}

	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	HTTP struct {
		Port           string `env:"HTTP_PORT,required"`
		UsePreforkMode bool   `env:"HTTP_USE_PREFORK_MODE" envDefault:"false"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	PG struct {
		PoolMax int    `env:"PG_POOL_MAX,required"`
		Conn    string `env:"POSTGRES_CONN,required"`
	}

	HDFSConfig struct {
		Addr    string `env:"HDFS_ADDR"    envDefault:"namenode:9000"`
		Enabled bool   `env:"HDFS_ENABLED" envDefault:"false"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
