package main

import (
	"collectionsservice/internal/config"
	"collectionsservice/internal/database"
	"collectionsservice/internal/grpc"
	"collectionsservice/internal/repository"
	"collectionsservice/internal/service"

	"github.com/rs/zerolog/log"
)

func main() {
	config.LoadEnv()

	db, err := database.ConnectToDatabase()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the database")
	}

	repo := repository.NewCollectionRepository(db)
	ser := service.NewCollectionService(repo)

	grpc.StartGRPCServer(ser)
}
