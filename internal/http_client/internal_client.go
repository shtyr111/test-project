package http_client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"test-project/internal/models"
	"test-project/internal/util/logging"

	log "github.com/sirupsen/logrus"
)

type InternalClient struct {
}

func New() *InternalClient {
	return &InternalClient{}
}

func (i InternalClient) SendToInternal(user models.User) (*models.InternalResponse, error) {
	requestBodyBytes, err := json.Marshal(user)
	if err != nil {
		log.Error("Произошла ошибка при парсинге", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", "http://localhost:8081/saveToInternalSystem", bytes.NewBuffer(requestBodyBytes))

	if err != nil {
		log.Error("Произошла ошибка создании запроса /saveToInternalSystem", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("id", strconv.Itoa(user.Number))

	req.URL.RawQuery = q.Encode()

	log.Info("URL: ", req.URL.String())

	client := &http.Client{
		Transport: &logging.LoggingRoundTripper{Rt: http.DefaultTransport},
	}

	log.Info("Производится запрос http://localhost:8081/saveToInternalSystem")
	resp, err := client.Do(req)

	if err != nil {
		log.Error("Произошла ошибка при запросе в /saveToInternalSystem", err)
		return nil, err
	}

	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Ошибка при чтении тела ответа /saveToInternalSystem", err)
		return nil, err
	}

	var response models.InternalResponse

	err = json.Unmarshal(respBodyBytes, &response)
	if err != nil {
		log.Error("Ошибка при парсинге тела ответа /saveToInternalSystem", err)
		return nil, err
	}

	return &response, nil
}
