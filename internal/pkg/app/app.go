//Вызов задач

package app

import (
	"net/http"
	"os"

	"github.com/jedyEvgeny/time_tracker/internal/app/endpoint"
	routes "github.com/jedyEvgeny/time_tracker/internal/delivery/http"
	"github.com/jedyEvgeny/time_tracker/internal/service"
	"github.com/jedyEvgeny/time_tracker/pkg/httpclient"
	"github.com/jedyEvgeny/time_tracker/pkg/logger"
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
	logger.Run()
	a := &App{}
	a.db = storage.New()
	a.s = service.New()
	a.c = httpclient.New()
	a.e = endpoint.New(a.s, a.db, a.c)
	err := a.db.Setup()
	if err != nil {
		return a, err
	}
	routes.SetupRoutes(a.e)
	return a, nil
}

func (a *App) Run() error {
	logger.Log.Info("Запускаем сервер")
	httpServerPort, err := findEnvironmentVariable("APP_HTTP_SERVER_PORT")
	if err != nil {
		return err
	}
	err = http.ListenAndServe(httpServerPort, nil)
	if err != nil {
		return err
	}
	logger.Log.Info("Сервер закончил работу")
	return nil
}

func findEnvironmentVariable(vrbl string) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Log.Debug("ошибка загрузки переменных окружения: ", err)
		return "", err
	}
	value := os.Getenv(vrbl)
	return value, nil
}
