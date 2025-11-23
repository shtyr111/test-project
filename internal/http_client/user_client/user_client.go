package user_client

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		return nil, fmt.Errorf("Произошла ошибка при сериализации тела запроса /saveToInternalSystem %+v\n: %w", user, err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8081/saveToInternalSystem", bytes.NewBuffer(requestBodyBytes))

	if err != nil {
		return nil, fmt.Errorf("произошла ошибка при создании запроса /saveToInternalSystem: %w", err)
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
		return nil, fmt.Errorf("произошла ошибка при запросе в /saveToInternalSystem: %w", err)
	}

	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("произошла ошибка при чтении тела ответа h/saveToInternalSystem: %w", err)
	}

	var response models.InternalResponse

	err = json.Unmarshal(respBodyBytes, &response)
	if err != nil {
		return nil, fmt.Errorf("произошла ошибка при десериализации тела ответа h/saveToInternalSystem: %w", err)
	}

	return &response, nil
}
