package storage

import (
	"context"
	"log"
	"time"
)

func (d *Database) AddStartTask(t TaskEnrichedUser, startTime time.Time) error {
	log.Println("Проверка наличия пользователя в БД")
	idUser, err := d.searchIDUser(t)
	log.Printf("Пользователь с серией паспорта %v и номером %v найден в БД под ID: %v\n", t.PassportSerie, t.PassportNumber, idUser)
	if err != nil {
		log.Printf("пользователь с серией паспорта %v и номером %v в БД не найден\n", t.PassportSerie, t.PassportNumber)
		return err
	}
	log.Println("Пользователь найден в БД")
	query := `
	INSERT INTO user_tasks
	(user_id, task_name, start_time, end_time, duration)
	VALUES
	($1, $2, $3)
	`
	_, err = d.poolConnectionsDb.Exec(context.Background(),
		query,
		idUser,
		t.TaskName,
		startTime,
	)

	if err != nil {
		log.Printf("не удалось добавить задачу %v для пользователя %v в БД.\n", t.TaskName, t.Name)
		return err
	}
	log.Printf("Время начала задачи для пользователя %v успешно создано в %v\n", t.Name, startTime)

	return nil
}

func (d *Database) searchIDUser(t TaskEnrichedUser) (int, error) {
	var userId int
	query := `
	SELECT id FROM users
	WHERE passport_serie = $1 AND
	passport_number = $2
	`
	err := d.poolConnectionsDb.QueryRow(context.Background(), query,
		t.PassportSerie,
		t.PassportNumber,
	).Scan(&userId)
	if err != nil {
		log.Println("пользователь не найден")
		return -1, err
	}
	return userId, nil
}
