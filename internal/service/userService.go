package service

import (
	"test-project/internal/models"
	"test-project/internal/repository"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func New(userRepository *repository.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (u UserService) FindById(id uuid.UUID) (*models.User, error) {
	return u.userRepository.FindById(id)
}

func (u UserService) SaveUsers(users []models.User) []models.User {
	var newUsers []models.User
	for i := range users {
		users[i].SetStatus("NEW")
		user, err := u.userRepository.Insert(users[i])
		if err != nil {
		}

		newUsers = append(newUsers, *user)
	}

	return newUsers
}

func (u UserService) SendUsersWithStatusNewToInternalSystem(sectionNumber int) {
	log.Info("Старт выполнения задачи SendUsersWithStatusNewToInternalSystem")
	users, err := u.userRepository.FindAllWithStatusNew()
	if err != nil {
		log.Error("Error while fetching users")
		return
	}

	log.Info(users)
}
