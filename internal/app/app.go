package app

import (
	"test-project/internal/config/application"
	"test-project/internal/config/logger"
	"test-project/internal/config/postgres"
	"test-project/internal/repository"
	sendToOlbScheduler "test-project/internal/scheduler/send_to_olb"
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
	log.Info("Properties config: ", config.Properties)

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

	sendToOlbScheduler := sendToOlbScheduler.New(config.Properties.SendUserToOlbSchedulerCron, config.Properties.SendUserToOlbSchedulerSectionAdvisoryCron, service)
	sendToOlbScheduler.Start()

	server.RunHttpServer()
	defer postgresConfig.ClosePool()
}
