package storage

import (
	"context"
	"errors"
	"log"
	"time"
)

func (d *Database) AddStartTask(t TaskEnrichedUser, startTime time.Time) error {
	log.Println("Проверка наличия пользователя в БД")
	idUser, err := d.searchIDUser(t)
	if err != nil {
		log.Printf("пользователь с серией паспорта %v и номером %v в БД не найден\n", t.PassportSerie, t.PassportNumber)
		return err
	}
	log.Printf("Пользователь с серией паспорта %v и номером %v найден в БД под ID: %v\n", t.PassportSerie, t.PassportNumber, idUser)
	query := `
	INSERT INTO user_tasks
	(user_id, task_name, start_time)
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

func (d *Database) AddFinishTask(t TaskEnrichedUser, endTime time.Time) error {
	log.Println("Проверка наличия пользователя в БД")
	idUser, err := d.searchIDUser(t)
	if err != nil {
		log.Printf("пользователь с серией паспорта %v и номером %v в БД не найден\n", t.PassportSerie, t.PassportNumber)
		return err
	}
	log.Printf("Пользователь с серией паспорта %v и номером %v найден в БД под ID: %v\n", t.PassportSerie, t.PassportNumber, idUser)

	ok := d.searchTaskName(t, idUser)
	if !ok {
		err := errors.New("задача не запущена")
		return err
	}
	log.Printf("Задача %v для пользователя %v запущена\n", t.TaskName, t.Name)
	query := `
	UPDATE user_tasks
	SET end_time = $1
	WHERE task_name = $2 AND user_id = $3
	`
	_, err = d.poolConnectionsDb.Exec(context.Background(),
		query,
		endTime,
		t.TaskName,
		idUser,
	)
	if err != nil {
		log.Printf("не удалось добавить окончание задачи %v для пользователя %v в БД\n", t.TaskName, t.Name)
		return err
	}
	log.Printf("Время окончания задачи для пользователя %v успешно создано в %v\n", t.Name, endTime)

	return nil
}

func (d *Database) searchTaskName(t TaskEnrichedUser, idUser int) bool {
	var taskStart time.Time
	query := `
	SELECT start_time FROM user_tasks
	WHERE task_name = $1 AND user_id = $2
	`
	err := d.poolConnectionsDb.QueryRow(context.Background(), query,
		t.TaskName,
		idUser,
	).Scan(&taskStart)
	if err != nil {
		log.Println("задача не начата")
		return false
	}
	return true
}
