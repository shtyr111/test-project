package handlers

import (
	"encoding/json"
	"fmt"
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
	strId := r.URL.Query().Get("id")
	id, err := uuid.Parse(strId)

	if err != nil {
		wrapErr := fmt.Errorf("Произошла ошибка при десериализации параметра id запроса: %w", err)
		handleError(wrapErr, w)

		return
	}

	user, err := u.userService.FindById(id)
	if err != nil {
		handleError(err, w)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func postUsersHandler(u UserHandler, w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		wrapErr := fmt.Errorf("Произошла ошибка при десериализации тела запроса: %w", err)
		handleError(wrapErr, w)

		return
	}

	var users []models.User
	err = json.Unmarshal(bodyBytes, &users)
	if err != nil {
		handleError(err, w)

		return
	}

	updateUsers := u.userService.SaveUsers(users)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateUsers)
}

func handleError(err error, w http.ResponseWriter) {
	log.Error(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(models.NewErrorResponse(http.StatusInternalServerError, err.Error()))
	return
}
