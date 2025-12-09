package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"test-project/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisUserRepository struct {
	redisClient *redis.Client
}

func NewRedisUserRepository(redisClient *redis.Client) *RedisUserRepository {
	return &RedisUserRepository{redisClient: redisClient}
}

func (r RedisUserRepository) Insert(ctx context.Context, user models.User) (*models.User, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userJSON, err := json.Marshal(user)

	if err != nil {
		return nil, fmt.Errorf("Произошла ошибка при сериализации. Параметры user: %+v\n, ошибка: %w", user, err)
	}

	err = r.redisClient.HSet(ctxWithTimeout, "users", user.Id.String(), string(userJSON)).Err()

	if err != nil {
		return nil, fmt.Errorf("Произошла ошибка при инсерт в redis. Параметры user: %+v\n, ошибка: %w", user, err)
	}

	return &user, nil
}

func (r RedisUserRepository) FindById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	str, err := r.redisClient.HGet(ctxWithTimeout, "users", id.String()).Result()
	var user models.User

	if err != nil {
		return nil, fmt.Errorf("пользователь с id %x не найден: %w", id.String(), err)
	}

	err = json.Unmarshal([]byte(str), &user)
	if err != nil {
		return nil, fmt.Errorf("Ошибка десериализации юзера с id %x: %w", id.String(), err)
	}

	return &user, nil
}
