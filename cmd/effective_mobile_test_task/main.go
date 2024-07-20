// Точка входа

package main

import (
	"log"

	"github.com/jedyEvgeny/time_tracker/internal/pkg/app"
)

func main() {
	a, _ := app.New()
	err := a.Run()
	if err != nil {
		log.Fatal("не удалось прослушать порт\n", err)
	}
}
