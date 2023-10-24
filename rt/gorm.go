package rt

import (
	"os"

	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ProvideGorm(cfg *Config, log zerolog.Logger) (*gorm.DB, error) {
	err := os.MkdirAll("data", os.ModePerm)
	if err != nil {
		log.Error().Err(err).Msg("failed to create data directory")
		return nil, err
	}
	name := cfg.Get("DB_NAME", "hub.db")
	db, err := gorm.Open(sqlite.Open("data/"+name), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
