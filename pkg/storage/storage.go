package storage

import (
	"context"
	"errors"
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

func (d *Database) AddPerson(e EnrichedUser) error {
	err := d.checkPresonInDB(e)
	if err != nil {
		return err
	}
	log.Println("Приступили к добавлению информации в БД")
	query := `
	INSERT INTO users
	(passport_serie, passport_number, surname, name, patronymic, address)
	VALUES
	($1, $2, $3, $4, $5, $6)
	`
	_, err = d.poolConnectionsDb.Exec(context.Background(),
		query,
		e.PassportSerie,
		e.PassportNumber,
		e.Surname,
		e.Name,
		e.Patronymic,
		e.Address,
	)
	return err
}

func (d *Database) checkPresonInDB(e EnrichedUser) error {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE passport_serie = $1 AND passport_number = $2"
	err := d.poolConnectionsDb.QueryRow(context.Background(), query, e.PassportSerie, e.PassportNumber).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		log.Println("человек с данными паспортными данными уже есть в БД:", e.PassportSerie, e.PassportNumber)
		err = errors.New("паспортные данные уже есть в БД")
		return err
	}
	return nil
}
