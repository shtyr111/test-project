package repository

import (
	"context"
	"test-project/internal/models"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgxpool"

	log "github.com/sirupsen/logrus"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func New(p *pgxpool.Pool) *UserRepository {
	log.Info(p)
	return &UserRepository{pool: p}
}

func (u UserRepository) Insert(user models.User) (*models.User, error) {
	log.Info(u.pool)
	_ = u.pool.QueryRow(context.Background(),
		"INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id, number", user.Name, user.Age).Scan(&user.Id, &user.Number)

	log.Info(user)
	return &user, nil
}

func (u UserRepository) FindById(id uuid.UUID) (*models.User, error) {
	log.Info(u.pool)

	var user models.User
	err := u.pool.QueryRow(context.Background(),
		"SELECT id, name, number, age FROM users where id = $1", id).Scan(&user.Id, &user.Name, &user.Number, &user.Age)

	if err != nil {
		log.Fatalf("Row scan failed: %v\n", err)
	}

	return &user, nil
}
