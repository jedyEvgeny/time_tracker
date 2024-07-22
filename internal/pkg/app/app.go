//Вызов задач

package app

import (
	"log"
	"net/http"
	"os"

	"github.com/jedyEvgeny/time_tracker/internal/app/endpoint"
	"github.com/jedyEvgeny/time_tracker/internal/service"
	"github.com/jedyEvgeny/time_tracker/pkg/httpclient"
	"github.com/jedyEvgeny/time_tracker/pkg/storage"
	"github.com/joho/godotenv"
)

type App struct {
	db *storage.Database
	e  *endpoint.Endpoint
	s  *service.Service
	c  *httpclient.Client
}

func New() (*App, error) {
	a := &App{}
	a.db = storage.New()
	a.s = service.New()
	a.c = httpclient.New()
	a.e = endpoint.New(a.s, a.db, a.c)
	err := a.db.Setup()
	if err != nil {
		return a, err
	}
	httpServerPath, err := findEnvironmentVariable("APP_HTTP_SERVER_PATH")
	if err != nil {
		return a, err
	}
	http.HandleFunc(httpServerPath, a.e.Status)
	return a, nil
}

func (a *App) Run() error {
	log.Println("Запускаем сервер")
	httpServerPort, err := findEnvironmentVariable("APP_HTTP_SERVER_PORT")
	if err != nil {
		return err
	}
	err = http.ListenAndServe(httpServerPort, nil)
	if err != nil {
		return err
	}
	log.Println("Сервер закончил работу")
	return nil
}

func findEnvironmentVariable(vrbl string) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("ошибка загрузки переменных окружения: %v\n", err)
		return "", err
	}
	value := os.Getenv(vrbl)
	return value, nil
}
