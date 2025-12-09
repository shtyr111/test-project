package user

import (
	"context"
	"errors"
	"sync"
	"test-project/internal/http_client/user_client"
	"test-project/internal/models"
	"test-project/internal/repository/postgres"
	redis2 "test-project/internal/repository/redis"
	"test-project/internal/service/websocket"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

type UserService struct {
	userRepository      *postgres.UserRepository
	internalClient      *user_client.InternalClient
	webSocketService    *websocket.WebsocketService
	redisUserRepository *redis2.RedisUserRepository
}

type UserRepositoryInterface interface {
	Insert(ctx context.Context, user models.User) (*models.User, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.User, error)
	Upsert(ctx context.Context, user *models.User) (*models.User, error)
}

var countSendUserWithStatusNewToInternalSystemAndSave int

func New(userRepository *postgres.UserRepository, internalClient *user_client.InternalClient, webSocketService *websocket.WebsocketService, redisUserRepository *redis2.RedisUserRepository) *UserService {
	return &UserService{userRepository: userRepository, internalClient: internalClient, webSocketService: webSocketService, redisUserRepository: redisUserRepository}
}

func (u UserService) FindById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := u.redisUserRepository.FindById(ctx, id)
	if err != nil {
		log.Warn("Ошибка получения юзера из редиса", err)
		return u.userRepository.FindById(ctx, id)
	}

	return user, nil
}

func (u UserService) SaveUsers(ctx context.Context, users []models.User) ([]models.User, error) {
	var newUsers []models.User
	for i := range users {
		users[i].SetStatus("NEW")
		user, err := u.userRepository.Insert(ctx, users[i])
		if err != nil {
			return nil, err
		}

		userSavedToRedis, err := u.redisUserRepository.Insert(ctx, *user)
		if err != nil {
			return nil, err
		}

		newUsers = append(newUsers, *userSavedToRedis)
	}

	return newUsers, nil
}

func (u UserService) PutUser(ctx context.Context, user *models.User) (*models.User, error) {
	oldUser, err := u.userRepository.FindById(ctx, user.Id)
	if err != nil {
		notification := models.Notification{Id: user.Id, Time: time.Now(), NewStatus: user.Status}
		u.webSocketService.SendToClient(&notification)

		return u.userRepository.Upsert(ctx, user)
	}

	notification := models.Notification{Id: user.Id, Time: time.Now(), NewStatus: user.Status, OldStatus: oldUser.Status}
	u.webSocketService.SendToClient(&notification)

	return u.userRepository.Upsert(ctx, user)
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
		log.Error(err)
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
	var mutex sync.Mutex

	for _, user := range users {
		wg.Add(1)
		go func() {
			defer wg.Done()
			u.sendUserWithStatusNewToInternalSystemAndSave(user, semaphore, &mutex)
		}()
	}

	wg.Wait()
}

func (u UserService) sendUserWithStatusNewToInternalSystemAndSave(user models.User, semaphore chan struct{}, mutex *sync.Mutex) {
	semaphore <- struct{}{}
	defer func() { <-semaphore }()

	mutex.Lock()
	countSendUserWithStatusNewToInternalSystemAndSave++
	log.Info("countSendUserWithStatusNewToInternalSystemAndSave: ", countSendUserWithStatusNewToInternalSystemAndSave)
	mutex.Unlock()

	internalResponse, err := u.internalClient.SendToInternal(user)
	if err != nil {
		log.Error(err)
		return
	}

	err = u.userRepository.UpdateStatusById(user.Id, internalResponse.Status)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Статус успешно обновлен")
}
