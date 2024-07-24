package service

import "time"

type TaskDuration struct {
	TaskName      string
	Hours         int
	Minutes       int
	TotalDuration time.Duration
}
