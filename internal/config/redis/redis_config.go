package redis

import (
	"context"
	"test-project/internal/config/application"

	"github.com/redis/go-redis/v9"

	log "github.com/sirupsen/logrus"
)

type RedisConfig struct {
	config *application.Config
}

func CreateNewConnection(config *application.Config) *redis.Client {
	conn := redis.NewClient(&redis.Options{
		Addr:         config.Redis.Address, // хост:порт из docker-compose
		Username:     config.Redis.Username,
		Password:     config.Redis.Password,     // REDIS_PASSWORD из environment
		DB:           config.Redis.DB,           // база по умолчани
		MinIdleConns: config.Redis.Pool.MinSize, // минимальное количество подключений
		PoolSize:     config.Redis.Pool.MaxSize, // максимально количество подключений
	})

	pong, err := conn.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("❌ Redis error:", err)
	}
	log.Info("✅ Redis OK:", pong)

	return conn
}
