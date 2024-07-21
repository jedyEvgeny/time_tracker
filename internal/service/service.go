// Реализация бизнес-логики
package service

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jedyEvgeny/time_tracker/pkg/storage"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) DecodeJSON(r *http.Request) (storage.User, error) {
	var userData storage.User
	log.Println("Приступили к декодированию входящего JSON пользователя")
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		return storage.User{}, err
	}
	log.Println("Закончили декодирование входящего JSON пользователя")
	return userData, nil
}
