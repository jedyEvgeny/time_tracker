package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	config "github.com/jedyEvgeny/time_tracker/etc"

	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Database struct {
	poolConnectionsDb *pgxpool.Pool
}

func New() *Database {
	return &Database{}
}

func (d *Database) Setup() error {
	config, err := config.NewConfig()
	if err != nil {
		return err
	}
	dbUrl := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		config.DatabaseScheme,
		config.StoragePostgresUser,
		config.StoragePostgresPassword,
		config.StoragePostgresHost,
		config.StoragePostgresPort,
		config.StoragePostgresDBName,
		config.StoragePostgresSSLMode,
	)
	log.Println("URL для подключения к БД:", dbUrl)

	// Создаём пул соединений
	d.poolConnectionsDb, err = pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Printf("не удаётся связаться с БД по пути: %v; %v\n", dbUrl, err)
		return err
	}
	log.Println("Соединение PostgreSQL установлено")

	pathToMigrations := "file://pkg/storage"

	m, err := migrate.New(
		pathToMigrations,
		dbUrl)

	if err != nil {
		log.Fatalf("не удаётся создать миграцию: %v", err)
	}

	if err := m.Up(); err != nil {
		log.Fatalf("не удалось применить миграцию: %v", err)
	}

	log.Println("Миграции при создании таблиц БД успешно применены")
	return nil
}

func (d *Database) AddPerson(u User) error {
	_, err := d.poolConnectionsDb.Exec(context.Background(), "INSERT INTO users (passport_number) VALUES ($1)", u.PassportNumber)
	return err
}
