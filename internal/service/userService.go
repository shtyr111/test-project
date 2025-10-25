package service

import (
	"test-project/internal/models"
)

func SaveUsers(users []models.User) []models.User {
	for i := range users {
		users[i].SetNumber(i)
	}

	return users
}
