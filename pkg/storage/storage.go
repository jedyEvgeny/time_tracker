package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	config "github.com/jedyEvgeny/time_tracker/etc"
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
	log.Println("соединение PostgreSQL установлено")

	// Создаем таблицу в БД
	_, err = d.poolConnectionsDb.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            passport_number TEXT NOT NULL
        )
    `)
	if err != nil {
		log.Printf("не удалось создать таблицу: %v\n", err)
		return err
	}
	log.Println("Таблица в БД успешно создана")
	return nil
}

func (d *Database) AddPerson(u User) error {
	_, err := d.poolConnectionsDb.Exec(context.Background(), "INSERT INTO users (passport_number) VALUES ($1)", u.PassportNumber)
	return err
}
