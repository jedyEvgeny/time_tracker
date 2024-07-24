// Точка входа

package main

import (
	"log"

	"github.com/jedyEvgeny/time_tracker/internal/pkg/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal("обнаружена ошибка при старте сервиса: ", err)
	}
	err = a.Run()
	if err != nil {
		log.Fatal("не удалось прослушать порт", err)
	}
}
