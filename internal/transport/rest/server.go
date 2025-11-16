package rest

import (
	"net/http"
	"os"
	"strings"
	"test-project/internal/config/application"
	"test-project/internal/transport/rest/handlers"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	userHandler *handlers.UserHandler
}

func New(userHandler *handlers.UserHandler) *Server {
	return &Server{userHandler: userHandler}
}

func (s Server) RunHttpServer() {
	log.Info("Http server started on port 8080")

	mux := http.NewServeMux()
	mux.HandleFunc("/users", s.userHandler.UsersHandler)

	e := http.ListenAndServe(strings.Join([]string{":", application.SERVER_CONFIG.Port}, ""), mux)

	if e != nil {
		log.Fatal("Произошла ошибка при старте сервера", e)
		os.Exit(1)
	}
}
