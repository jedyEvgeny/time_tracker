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

	err = executeMigrations(dbUrl)
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

func executeMigrations(dbUrl string) error {
	pathToMigrations := "file://migrations"

	m, err := migrate.New(
		pathToMigrations,
		dbUrl)
	if err != nil {
		log.Printf("не удаётся создать миграцию: %v\n", err)
		return err
	}
	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("В базе данных нет новых миграций для применения")
		return nil
	}
	if err != nil {
		log.Println("ошибка при выполнении миграции:", err)
		return err
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

func (d *Database) GetUsersByFilter(filter EnrichedUser) ([]EnrichedUser, error) {
	var usersSlice []EnrichedUser

	query := "SELECT id, passport_serie, passport_number, surname, name, patronymic, address FROM users WHERE " +
		"($1 = '' OR passport_serie = $1) AND " +
		"($2 = '' OR passport_number = $2) AND " +
		"($3 = '' OR surname = $3) AND " +
		"($4 = '' OR name = $4) AND " +
		"($5 = '' OR patronymic = $5) AND " +
		"($6 = '' OR address = $6)"

	rows, err := d.poolConnectionsDb.Query(context.Background(), query,
		filter.PassportSerie,
		filter.PassportNumber,
		filter.Surname,
		filter.Name,
		filter.Patronymic,
		filter.Address,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user EnrichedUser
		err := rows.Scan(
			&user.ID,
			&user.PassportSerie,
			&user.PassportNumber,
			&user.Surname,
			&user.Name,
			&user.Patronymic,
			&user.Address,
		)
		if err != nil {
			return nil, err
		}
		usersSlice = append(usersSlice, user)
	}

	return usersSlice, nil
}
