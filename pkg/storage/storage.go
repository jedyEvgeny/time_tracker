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
		log.Printf("миграции не требуются: %v\n", err)
		return nil
	}
	log.Println("Миграции успешно выполнены")
	return nil
}

func (d *Database) AddPerson(e EnrichedUser) error {
	serie := e.PassportSerie
	number := e.PassportNumber
	err := d.checkPresonInDB(serie, number)
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

func (d *Database) checkPresonInDB(serie, number string) error {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE passport_serie = $1 AND passport_number = $2"
	err := d.poolConnectionsDb.QueryRow(context.Background(), query, serie, number).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		log.Println("человек с данным паспортом есть в БД:", serie, number)
		err = errors.New("паспортные данные есть в БД")
		return err
	}
	return nil
}

func (d *Database) DelPerson(serie, number string) error {
	err := d.checkPresonInDB(serie, number)
	if err == nil {
		err := errors.New("отсутствует пользователь для удаления из БД")
		log.Println("отсутствует пользователь для удаления из БД")
		return err
	}
	log.Println("Приступили к удалению информации из БД")
	query := "DELETE FROM users WHERE passport_serie = $1 AND passport_number = $2"

	_, err = d.poolConnectionsDb.Exec(context.Background(),
		query,
		serie,
		number,
	)
	return err
}

func (d *Database) ChangePerson(e EnrichedUser) error {
	serie := e.PassportSerie
	number := e.PassportNumber
	err := d.checkPresonInDB(serie, number)
	if err == nil {
		err := errors.New("отсутствует пользователь для удаления из БД")
		log.Println("отсутствует пользователь для удаления из БД")
		return err
	}
	log.Println("Приступили к изменению информации в БД")
	query := `
	UPDATE users
	SET surname = $1, name = $2, patronymic = $3, address = $4
	WHERE passport_serie = $5 AND passport_number = $6
	`
	_, err = d.poolConnectionsDb.Exec(context.Background(),
		query,
		e.Surname,
		e.Name,
		e.Patronymic,
		e.Address,
		e.PassportSerie,
		e.PassportNumber,
	)
	return err
}
