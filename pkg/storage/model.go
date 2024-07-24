package storage

import "time"

type User struct {
	PassportNumber string `json:"passportNumber"`
}

type EnrichedUser struct {
	ID             int    `json:"id"`
	PassportSerie  string `json:"passportSerie"`
	PassportNumber string `json:"passportNumber"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
}

type UserTask struct {
	ID        int       `json:"idTask"`
	UserID    int       `json:"userID"`
	TaskName  string    `json:"taskName"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Duration  int       `json:"duration"`
}

type TaskEnrichedUser struct {
	EnrichedUser
	UserTask
}
