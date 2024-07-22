package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	config "github.com/jedyEvgeny/time_tracker/etc"

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

	err = useMigrations(dbUrl)
	if err != nil {
		return err
	}

	log.Println("Таблицы БД готовы к работе")

	//Создаём пул для добавления информации в БД
	d.poolConnectionsDb, err = pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Printf("не удаётся связаться с БД по пути: %v; %v\n", dbUrl, err)
		return err
	}
	log.Println("Соединение PostgreSQL установлено")

	return nil
}

func useMigrations(dbUrl string) error {
	pathToMigrations := "file://migrations"

	m, err := migrate.New(
		pathToMigrations,
		dbUrl)
	if err != nil {
		log.Fatalf("не удаётся создать миграцию: %v\n", err)
	}

	err = m.Up()
	if err != nil {
		log.Printf("Не требуется применять миграцию: %v\n", err)
		return nil
	}

	log.Println("Миграции успешно выполнены")
	return nil
}

func (d *Database) AddPerson(serie, number string) error {
	log.Println("Приступили к добавлению информации в БД")
	_, err := d.poolConnectionsDb.Exec(context.Background(), "INSERT INTO users (passport_serie, passport_number) VALUES ($1, $2)", serie, number)
	return err
}
