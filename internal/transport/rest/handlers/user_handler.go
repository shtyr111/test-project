package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"test-project/internal/models"
	"test-project/internal/service"

	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
)

type UserHandler struct {
	userService *service.UserService
}

func New(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (u UserHandler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsersHandler(u, w, r)
	case http.MethodPost:
		postUsersHandler(u, w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getUsersHandler(u UserHandler, w http.ResponseWriter, r *http.Request) {
	log.Info("Получен запрос GET /users: ", r)

	strId := r.URL.Query().Get("id")
	id, err := uuid.Parse(strId)

	if err != nil {
		log.Fatalf("Ошибка парсинга: %v\n", err)
	}

	user, err := u.userService.FindById(id)
	if err != nil {
		log.Fatalf("Ошибка парсинга: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	log.Info("Отправлен ответ GET /users:", user)
}

func postUsersHandler(u UserHandler, w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("Ошибка чтения тела запроса:", err)
	}

	log.Info("Получен запрос POST /users. Хедеры: ", r.Header, " Тело: ", string(bodyBytes))

	var users []models.User
	e := json.Unmarshal(bodyBytes, &users)
	if e != nil {
		log.Error("Произошла ошибка при десериалазации тела: ", e.Error())
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	updateUsers := u.userService.SaveUsers(users)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateUsers)
	log.Info("Отправлен ответ POST /users:", updateUsers)
}
