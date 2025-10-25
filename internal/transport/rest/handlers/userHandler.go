package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"test-project/internal/models"
	"test-project/internal/service"

	log "github.com/sirupsen/logrus"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsersHandler(w, r)
	case http.MethodPost:
		postUsersHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Получен запрос GET /users: ", r)

	users := []models.User{{Name: "Alice", Age: 1}, {Name: "Bob", Age: 2}}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	log.Info("Отправлен ответ GET /users:", users)
}

func postUsersHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	log.Info("Получен запрос POST /users. Хедеры: ", r.Header, " Тело: ", string(bodyBytes))

	var users []models.User
	e := json.Unmarshal(bodyBytes, &users)
	if e != nil {
		if err != nil {
			log.Println("Ошибка чтения тела запроса:", err)
		}

		log.Error("Произошла ошибка при десериалазации тела: ", e.Error())
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	updateUsers := service.SaveUsers(users)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateUsers)
	log.Info("Отправлен ответ POST /users:", updateUsers)
}
