package service

import (
	"context"
	"errors"
	"sync"
	"test-project/internal/http_client"
	"test-project/internal/models"
	"test-project/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

type UserService struct {
	userRepository *repository.UserRepository
	internalClient *http_client.InternalClient
}

func New(userRepository *repository.UserRepository, internalClient *http_client.InternalClient) *UserService {
	return &UserService{userRepository: userRepository, internalClient: internalClient}
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

func (u UserService) FindAndSendUsersWithStatusNewToInternalSystem(sectionNumber int, parallelCurrencySend int) {
	log.Info("Старт выполнения задачи FindAndSendUsersWithStatusNewToInternalSystem")
	ctx := context.Background()
	tx, err := u.userRepository.BeginTxWithAdvisoryLock(ctx, sectionNumber)

	if err != nil {
		log.Error("Не удалось получить блокировку", err)
		tx.Rollback(ctx)
	}

	if tx == nil {
		return
	}

	defer func() {
		err := tx.Rollback(ctx) // откатится, если commit не был вызван
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error("Ошибка при откате: %v", err)
		}
	}()

	users, err := u.userRepository.FindAllWithStatusNew(tx)
	if err != nil {
		log.Error("Error while fetching users", err)
		return
	}

	if len(users) != 0 {
		u.parallelSendUsersWithStatusNewToInternalSystemAndSave(users, parallelCurrencySend)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("Ошибка при коммите транзакции: %v", err)
	}
}

func (u UserService) parallelSendUsersWithStatusNewToInternalSystemAndSave(users []models.User, parallelCurrencySend int) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, parallelCurrencySend)

	for _, user := range users {
		wg.Add(1)
		go u.sendUserWithStatusNewToInternalSystemAndSave(user, semaphore, &wg)
	}

	wg.Wait()
}

func (u UserService) sendUserWithStatusNewToInternalSystemAndSave(user models.User, semaphore chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	semaphore <- struct{}{}
	defer func() { <-semaphore }()

	internalResponse, err := u.internalClient.SendToInternal(user)
	if err != nil {
		return
	}

	err = u.userRepository.UpdateStatusById(user.Id, internalResponse.Status)
	if err != nil {
		return
	}

	log.Info("Статус успешно обновлен")
}
