package app

import (
	"test-project/internal/config/application"
	"test-project/internal/config/logger"
	"test-project/internal/config/postgres"
	"test-project/internal/repository"
	"test-project/internal/service"
	"test-project/internal/transport/rest"
	"test-project/internal/transport/rest/handlers"

	log "github.com/sirupsen/logrus"
)

func Run() {
	logger.StartLogger()
	config := application.LoadApplicationConfig()

	log.Info("Server config: ", config.Server)
	log.Info("Database config: ", config.Database)

	postgresConfig := postgres.New(&config)

	err := postgresConfig.InitializeMigration()
	if err != nil {
		return
	}
	pool, err1 := postgresConfig.InitPool()
	if err1 != nil {
	}

	repository := repository.New(pool)
	service := service.New(repository)
	handler := handlers.New(service)
	server := rest.New(handler)

	server.RunHttpServer()
	defer postgresConfig.ClosePool()
}
