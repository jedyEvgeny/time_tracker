package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

type Database struct {
	poolConnectionsDb *pgxpool.Pool
}

func New() *Database {
	return &Database{}
}

func (d *Database) Setup() error {
	scheme, user, password, host, port, dbName, sslMode, err := findEnvironmentVariables()
	if err != nil {
		return err
	}
	dbUrl := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		scheme,
		user,
		password,
		host,
		port,
		dbName,
		sslMode,
	)
	log.Println("URL для подключения к БД:", dbUrl)

	// Создаём пул соединений
	d.poolConnectionsDb, err = pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Printf("нет возможности связаться с БД: %v\n", err)
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
	log.Println("Таблица успешно создана")
	return nil
}

func findEnvironmentVariables() (string, string, string, string, string, string, string, error) {
	err := godotenv.Load("etc/.env")
	if err != nil {
		log.Printf("ошибка загрузки переменных окружения: %v\n", err)
		return "", "", "", "", "", "", "", err
	}
	sheme := os.Getenv("APP_DATABASE_SCHEME")
	user := os.Getenv("APP_STORAGE_POSTGRES_USER")
	password := os.Getenv("APP_STORAGE_POSTGRES_PASSWORD")
	host := os.Getenv("APP_STORAGE_POSTGRES_HOST")
	port := os.Getenv("APP_STORAGE_POSTGRES_PORT")
	dbName := os.Getenv("APP_STORAGE_POSTGRES_DBNAME")
	sslMode := os.Getenv("APP_STORAGE_POSTGRES_SSLMODE")
	return sheme, user, password, host, port, dbName, sslMode, nil
}

func (d *Database) AddPerson(u User) error {
	_, err := d.poolConnectionsDb.Exec(context.Background(), "INSERT INTO users (passport_number) VALUES ($1)", u.PassportNumber)
	return err
}
