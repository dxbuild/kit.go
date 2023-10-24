package rt

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type Config struct {
	log zerolog.Logger
}

func (c Config) Get(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	fmt.Printf("key: %s, value: %s, exists: %t\n", key, value, exists)
	if !exists {
		return fallback
	}
	return value
}

func (c Config) GetDuration(key string, fallback time.Duration) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		c.log.Error().Err(err).Msgf("invalid duration for %s: %s", key, value)
		return fallback
	}
	return duration
}

func (c Config) GetInt(key string, fallback int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		c.log.Error().Err(err).Msgf("invalid int for %s: %s", key, value)
		return fallback
	}
	return i
}

func ProvideConfig(log zerolog.Logger) (*Config, error) {
	log.Debug().Msg("ProvideConfig")
	name := ".env"
	env := os.Getenv("ENV")
	if env != "" {
		name = ".env." + env
	}
	log.Info().Msgf("loading config from %s", name)
	err := godotenv.Load(name)
	if err != nil {
		log.Info().Msg("no .env file found")
	}
	return &Config{
		log: log,
	}, nil
}
