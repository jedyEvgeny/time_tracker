// Реализация бизнес-логики
package service

import (
	"encoding/json"
	"net/http"

	"github.com/jedyEvgeny/time_tracker/internal/database"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) DecodeJSON(r *http.Request) (database.User, error) {
	var userData database.User
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		return database.User{}, err
	}
	return userData, nil
}
