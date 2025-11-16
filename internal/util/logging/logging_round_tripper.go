package logging

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type LoggingRoundTripper struct {
	Rt http.RoundTripper
}

func (l *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Логируем метод, URL и заголовки запроса
	log.Printf(">> Request %s %s\nHeaders: %v", req.Method, req.URL.String(), req.Header)

	if req.Body != nil {
		bodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			return nil, err
		}
		// Лог тела запроса
		log.Printf(">> Request Body: %s", string(bodyBytes))
		// Восстанавливаем тело запроса для дальнейшей отправки
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	start := time.Now()
	// Выполнение запроса
	resp, err := l.Rt.RoundTrip(req)
	duration := time.Since(start)

	if err != nil {
		log.Printf("Request error: %v", err)
		return resp, err
	}

	// Копируем тело ответа для логирования
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return resp, err
	}
	// Восстанавливаем тело ответа для дальнейшей обработки
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBody))

	// Логируем статус, заголовки и тело ответа
	log.Printf("<< Response status: %s\nHeaders: %v\nBody: %s\nDuration: %v",
		resp.Status, resp.Header, string(respBody), duration)

	return resp, nil
}
