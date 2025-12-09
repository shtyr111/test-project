package app

import (
	"test-project/internal/config/application"
	"test-project/internal/config/logger"
	"test-project/internal/config/postgres"
	"test-project/internal/config/redis"
	"test-project/internal/http_client/user_client"
	postgres2 "test-project/internal/repository/postgres"
	redis2 "test-project/internal/repository/redis"
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

	redisConn := redis.CreateNewConnection(&config)

	log.Info("Redis config: ", redisConn)

	client := user_client.New()
	repository := postgres2.New(pool)
	redisUserRepository := redis2.NewRedisUserRepository(redisConn)
	webSocketService := websocket.NewWebsocketService()
	userService := user2.New(repository, client, webSocketService, redisUserRepository)
	handler := user.New(userService)
	wsHandler := ws.NewWebsocketHandler(webSocketService)

	server := rest.New(handler, wsHandler)

	sendToOlbScheduler := sendToOlbScheduler.New(config.Properties.SendUserToOlbSchedulerCron, config.Properties.SendUserToOlbSchedulerSectionAdvisoryCron,
		config.Properties.SendUserToOlbSchedulerParallelCurrencySend, userService)
	sendToOlbScheduler.Start()

	server.RunHttpServer()
}
