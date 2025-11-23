package repository

import (
	"context"
	"fmt"
	"test-project/internal/models"
	"time"

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

// Прокинуть контекст сверху
func (u UserRepository) Insert(user models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := u.pool.QueryRow(ctx,
		"INSERT INTO users (name, age, status) VALUES ($1, $2, $3) RETURNING id, number", user.Name, user.Age, user.Status).Scan(&user.Id, &user.Number)

	if err != nil {
		return nil, fmt.Errorf("Произошла ошибка при инсерт. Параметры user: %+v\n, ошибка: %w", user, err)
	}

	return &user, nil
}

func (u UserRepository) FindById(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := u.pool.QueryRow(context.Background(),
		"SELECT id, name, number, age, status FROM users where id = $1", id).Scan(&user.Id, &user.Name, &user.Number, &user.Age, &user.Status)

	if err != nil {
		return nil, fmt.Errorf("пользователь с id %x не найден: %w", id.String(), err)
	}

	return &user, nil
}

func (u UserRepository) Upsert(user *models.User) (*models.User, error) {
	_, err := u.pool.Exec(context.Background(),
		"INSERT INTO users (id, name, age, status) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET name = $2, age = $3, status = $4", user.Id, user.Name, user.Age, user.Status)

	if err != nil {
		return nil, fmt.Errorf("Произошла ошибка при upsert. Параметры user: %+v\n, ошибка: %w", user, err)
	}

	return user, nil
}

func (u UserRepository) UpdateStatusById(id uuid.UUID, status string) error {
	_, err := u.pool.Exec(context.Background(),
		"UPDATE users SET status = $1 WHERE id = $2", status, id)

	if err != nil {
		return fmt.Errorf("произошла ошибка при обновлении юзена с id %x: %w", id.String(), err)
	}

	return nil
}

func (u UserRepository) FindAllWithStatusNew(tx pgx.Tx) ([]models.User, error) {
	var users []models.User

	rows, err := tx.Query(context.Background(),
		"SELECT * FROM users where status = $1", "NEW")

	if err != nil {
		return nil, fmt.Errorf("Произошла ошибка при выполненении запроса на получение юзеров с статусом NEW: %w", err)
	}

	for rows.Next() {
		var user models.User

		err1 := rows.Scan(&user.Id, &user.Name, &user.Age, &user.Number, &user.Status)

		if err1 != nil {
			return nil, fmt.Errorf("Произошла ошибка при сканировании после запроса на получение юзеров с статусом NEW: %w", err)
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
		log.Info("Блокировка не была получена")
		// Если lock не получен, откатываем и освобождаем соединение, возвращаем nil без ошибки
		tx.Rollback(ctx)
		conn.Release()
		return nil, nil
	}

	log.Info("Блокировка успешно получена")
	return tx, nil
}
