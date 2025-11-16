package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"test-project/internal/config/application"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	log "github.com/sirupsen/logrus"
)

var POOL *pgxpool.Pool

type PostgresConfig struct {
	config *application.Config
	dns    string
	Pool   *pgxpool.Pool
}

func New(config *application.Config) *PostgresConfig {
	databaseConfig := application.DATABASE_CONFIG
	dns := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		databaseConfig.Host, databaseConfig.Port, databaseConfig.Username, databaseConfig.Password, databaseConfig.DBName)

	return &PostgresConfig{config: config, dns: dns}
}

func (p PostgresConfig) InitPool() (*pgxpool.Pool, error) {

	pool, err := pgxpool.New(context.Background(), p.dns)

	if err != nil {
		log.Fatal("Ошибка подключения:", err)
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("Ошибка соединения:", err)
		pool.Close()
		return nil, err
	}

	p.Pool = pool
	log.Info(p.Pool)
	return pool, nil
}

func (p PostgresConfig) ClosePool() {
	log.Info("here")
	if p.Pool != nil {
		p.Pool.Close()
	}
}

func (p PostgresConfig) InitializeMigration() error {
	db, err := sql.Open("postgres", p.dns)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver,
	)
	if err != nil {
		log.Fatal(err)
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
		return err
	}
	return nil
}
