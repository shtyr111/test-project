package app

import (
	"test-project/internal/config"
	"test-project/internal/transport/rest/handlers"

	log "github.com/sirupsen/logrus"
)

func Run() {
	config.StartLogger()
	config.LoadApplicationConfig()
	log.Info("Server config: ", config.SERVER_CONFIG)
	log.Info("Database config: ", config.DATABASE_CONFIG)
	handlers.RunHttpServer()
}
