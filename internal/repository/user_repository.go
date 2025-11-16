package repository

import (
	"context"
	"test-project/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"

	log "github.com/sirupsen/logrus"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func New(p *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: p}
}

func (u UserRepository) Insert(user models.User) (*models.User, error) {
	err := u.pool.QueryRow(context.Background(),
		"INSERT INTO users (name, age, status) VALUES ($1, $2, $3) RETURNING id, number", user.Name, user.Age, user.Status).Scan(&user.Id, &user.Number)

	if err != nil {
		log.Error(err)

		return nil, err
	}

	log.Info(user)
	return &user, nil
}

func (u UserRepository) FindById(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := u.pool.QueryRow(context.Background(),
		"SELECT id, name, number, age FROM users where id = $1", id).Scan(&user.Id, &user.Name, &user.Number, &user.Age)

	if err != nil {
		log.Fatalf("Row scan failed: %v\n", err)
	}

	return &user, nil
}

func (u UserRepository) FindAllWithStatusNew() ([]models.User, error) {
	var users []models.User

	rows, err := u.pool.Query(context.Background(),
		"SELECT * FROM users where status = $1", "NEW")

	if err != nil {
		log.Fatalf("Row scan failed: %v\n", err)

		return nil, err
	}

	for rows.Next() {
		var user models.User

		err1 := rows.Scan(&user.Id, &user.Name, &user.Age, &user.Number, &user.Status)

		if err1 != nil {
			log.Fatalf("Row scan failed: %v\n", err)

			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (u UserRepository) BeginTxWithAdvisoryLock(ctx context.Context, lockID int) (pgx.Tx, error) {
	conn, err := u.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		conn.Release()
		return nil, err
	}

	var gotLock bool
	err = tx.QueryRow(ctx, "SELECT pg_try_advisory_xact_lock($1)", lockID).Scan(&gotLock)
	if err != nil {
		tx.Rollback(ctx)
		conn.Release()
		return nil, err
	}

	if !gotLock {
		// Если lock не получен, откатываем и освобождаем соединение, возвращаем nil без ошибки
		tx.Rollback(ctx)
		conn.Release()
		return nil, nil
	}

	return tx, nil
}
