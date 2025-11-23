package rest

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
	"test-project/internal/config/application"
	"test-project/internal/transport/rest/handlers/user"
	"test-project/internal/transport/rest/handlers/websocket"
	"time"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	userHandler *user.UserHandler
	wsHandler   *websocket.WebsocketHandler
}

func New(userHandler *user.UserHandler, wsHandler *websocket.WebsocketHandler) *Server {
	return &Server{userHandler: userHandler, wsHandler: wsHandler}
}

func (s Server) RunHttpServer() {
	log.Info("Http server started on port 8080")

	mux := http.NewServeMux()
	mux.HandleFunc("/users", s.userHandler.UsersHandler)
	mux.HandleFunc("/ws", s.wsHandler.WebSocketHandler)

	//loggerMux := loggingMiddleware(mux)

	e := http.ListenAndServe(strings.Join([]string{":", application.SERVER_CONFIG.Port}, ""), mux)

	if e != nil {
		log.Fatal("Произошла ошибка при старте сервера", e)
		os.Exit(1)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Логируем URL, метод и параметры запроса
		log.Printf("Request: %s %s %s", r.Method, r.URL.String(), r.RemoteAddr)

		// Логируем заголовки запроса
		//for name, values := range r.Header {
		//	for _, value := range values {
		//		log.Printf("Header: %s=%s", name, value)
		//	}
		//}

		// Читаем тело запроса
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // восстанавливаем тело для обработки дальше
			log.Printf("Request Body: %s", string(bodyBytes))
		}

		// Создаем ResponseWriter для перехвата ответа
		rw := &responseWriter{ResponseWriter: w, body: &bytes.Buffer{}}

		next.ServeHTTP(rw, r) // вызываем следующий обработчик

		// Логируем статус код и тело ответа
		log.Printf("Response status: %d", rw.statusCode)
		log.Printf("Response body: %s", rw.body.String())

		log.Printf("Request processed in %s", time.Since(start))
	})
}

// Обертка над ResponseWriter для перехвата кода статуса и тела
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}
