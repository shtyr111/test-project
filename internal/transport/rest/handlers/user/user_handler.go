package user

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"test-project/internal/models"
	"test-project/internal/service/user"

	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
)

type UserHandler struct {
	userService *user.UserService
}

func New(userService *user.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (u UserHandler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsersHandler(u, w, r)
	case http.MethodPost:
		postUsersHandler(u, w, r)
	case http.MethodPut:
		putUsersHandler(u, w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getUsersHandler(u UserHandler, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	strId := r.URL.Query().Get("id")
	id, err := uuid.Parse(strId)

	if err != nil {
		wrapErr := fmt.Errorf("Произошла ошибка при десериализации параметра id запроса: %w", err)
		handleError(wrapErr, w)

		return
	}

	user, err := u.userService.FindById(ctx, id)
	if err != nil {
		handleError(err, w)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func postUsersHandler(u UserHandler, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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

	updateUsers := u.userService.SaveUsers(ctx, users)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateUsers)
}

func putUsersHandler(u UserHandler, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		wrapErr := fmt.Errorf("Произошла ошибка при десериализации тела запроса: %w", err)
		handleError(wrapErr, w)

		return
	}

	var user models.User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		handleError(err, w)

		return
	}

	updateUser, err := u.userService.PutUser(ctx, &user)
	if err != nil {
		handleError(err, w)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateUser)
}

func handleError(err error, w http.ResponseWriter) {
	log.Error(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(models.NewErrorResponse(http.StatusInternalServerError, err.Error()))
	return
}
