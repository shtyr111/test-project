package handlers

import (
	"net/http"
	"os"
	"test-project/internal/config"

	log "github.com/sirupsen/logrus"
)

func RunHttpServer() {
	log.Info("Http server started on port 8080")

	mux := http.NewServeMux()
	mux.HandleFunc("/users", UsersHandler)

	e := http.ListenAndServe(":"+config.SERVER_CONFIG.Port, mux)

	if e != nil {
		log.Fatal("Произошла ошибка при старте сервера", e)
		os.Exit(1)
	}
}
