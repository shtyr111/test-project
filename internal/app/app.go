package app

import (
	"test-project/internal/config/application"
	"test-project/internal/config/logger"
	"test-project/internal/config/postgres"
	"test-project/internal/http_client/user_client"
	"test-project/internal/repository"
	sendToOlbScheduler "test-project/internal/scheduler/send_to_olb"
	user2 "test-project/internal/service/user"
	"test-project/internal/service/websocket"
	"test-project/internal/transport/rest"
	"test-project/internal/transport/rest/handlers/user"
	ws "test-project/internal/transport/rest/handlers/websocket"

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
	pool, err := postgresConfig.InitPool()
	if err == nil {
	}
	defer postgresConfig.ClosePool()

	client := user_client.New()
	repository := repository.New(pool)
	webSocketService := websocket.NewWebsocketService()
	userService := user2.New(repository, client, webSocketService)
	handler := user.New(userService)
	wsHandler := ws.NewWebsocketHandler(webSocketService)

	server := rest.New(handler, wsHandler)

	sendToOlbScheduler := sendToOlbScheduler.New(config.Properties.SendUserToOlbSchedulerCron, config.Properties.SendUserToOlbSchedulerSectionAdvisoryCron,
		config.Properties.SendUserToOlbSchedulerParallelCurrencySend, userService)
	sendToOlbScheduler.Start()

	server.RunHttpServer()
}
