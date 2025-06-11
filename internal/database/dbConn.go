package database

import (
	"collectionsservice/internal/config"
	"collectionsservice/internal/models"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDatabase() (*gorm.DB, error) {
	dsn := config.GetPostgresDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("Error connecting to the database")
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get sql instance")
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		log.Error().Err(err).Msg("Database ping failed")
		return nil, err
	}

	if err := db.AutoMigrate(&models.Collection{}, &models.Request{}); err != nil {
		log.Error().Err(err).Msg("Failed auto-migrating tables")
		return nil, err
	}

	return db, nil
}
