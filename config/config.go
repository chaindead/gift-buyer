package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

// Config holds all configuration values for the bot
type Config struct {
	// Notifications about new gifts
	Admin string `env:"TG_ADMIN,required"`

	// Telegram API configuration
	AppID   int32  `env:"TG_APP_ID,required"`
	AppHash string `env:"TG_API_HASH,required"`

	// Optional session string for user authentication
	Session string `env:"TG_SESSION"`

	// Check gift every diration
	PollInterval time.Duration `env:"POLL_INTERVAL" envDefault:"1s"`
}

// LoadConfig parses environment variables into Config struct
func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
