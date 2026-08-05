package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"

	"alwis.dev/selectify/internal/program"
)

type DatabaseConfiguration struct {
	ApplicationName string `envconfig:"-"` // Set from program.AppPrefix, not from env

	Host     string `envconfig:"host" required:"true"`
	Port     string `envconfig:"port" default:"5432"`
	Database string `envconfig:"database" required:"true"`
	DatabaseConnectionSettings

	ReadWrite struct {
		User     string `envconfig:"user" required:"true"`
		Password string `envconfig:"password" required:"true"`
	} `envconfig:"rw"`

	ReadOnly struct {
		User     string `envconfig:"user" required:"true"`
		Password string `envconfig:"password" required:"true"`
	} `envconfig:"ro"`
}

type DatabaseConnectionSettings struct {
	MaxIdle     int           `envconfig:"max_idle" default:"15"`
	MaxOpen     int           `envconfig:"max_open" default:"25"`
	MaxLifetime time.Duration `envconfig:"max_lifetime" default:"10m"`
}

// ReadWriteDbUrl builds the PostgreSQL connection string for read-write database
func (c *DatabaseConfiguration) ReadWriteDbUrl() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		c.Host,
		c.Port,
		c.Database,
		c.ReadWrite.User,
		c.ReadWrite.Password,
	)
}

// ReadOnlyDbUrl builds the PostgreSQL connection string for read-only database
func (c *DatabaseConfiguration) ReadOnlyDbUrl() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		c.Host,
		c.Port,
		c.Database,
		c.ReadOnly.User,
		c.ReadOnly.Password,
	)
}

func LoadDatabaseConfig() (*DatabaseConfiguration, error) {
	var wrapper struct {
		DB DatabaseConfiguration `envconfig:"db"`
	}

	prefix := program.AppPrefix
	if prefix == "" {
		prefix = "API"
	}

	err := envconfig.Process(prefix, &wrapper)
	if err != nil {
		return nil, fmt.Errorf("failed to load database configuration: %w", err)
	}

	wrapper.DB.ApplicationName = program.AppPrefix

	return &wrapper.DB, nil
}
