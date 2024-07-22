// Реализация бизнес-логики
package service

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jedyEvgeny/time_tracker/pkg/storage"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) DecodeJSON(r *http.Request) (string, string, error) {
	var userData storage.User
	log.Println("Приступили к декодированию входящего JSON пользователя")
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		return "", "", err
	}
	log.Println("Закончили декодирование входящего JSON пользователя")
	serie, number, err := checkJSON(userData)
	if err != nil {
		return "", "", err
	}
	log.Println("Проверки корректности JSON выполнены успешно")
	return serie, number, nil
}

func checkJSON(u storage.User) (string, string, error) {
	if len(u.PassportNumber) != 11 {
		log.Println("ожидалось 11 символов в запросе")
		err := errors.New("неверная длина строки, ожидается 11 символов")
		return "", "", err
	}
	parts := strings.Split(u.PassportNumber, " ")
	if len(parts) != 2 {
		log.Println("ожидался один пробел в запросе")
		err := errors.New("ожидался один пробел в запросе")
		return "", "", err
	}
	if len(parts[0]) != 4 && len(parts[1]) != 6 {
		log.Println("неверная длина блока серии или номера паспорта")
		err := errors.New("неверная длина блока серии или номера паспорта")
		return "", "", err
	}
	elFirst, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Println("в первом блоке тела запроса не только цифры")
		err := errors.New("в первом блоке тела запроса не только цифры")
		return "", "", err
	}
	elSecond, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Println("во втором блоке тела запроса не только цифры")
		err := errors.New("во втором блоке тела запроса не только цифры")
		return "", "", err
	}
	if elFirst < 0 || elSecond < 0 {
		log.Println("цифровой блок со знаком минус не допустим")
		err := errors.New("цифровой блок со знаком минус не допустим")
		return "", "", err
	}
	return parts[0], parts[1], nil
}

func (s *Service) EnrichUserData(r *http.Response, serie, number string) (storage.EnrichedUser, error) {
	var userData storage.EnrichedUser
	userData.PassportNumber = number
	userData.PassportSerie = serie
	if r == nil {
		log.Println("Обогощение данных со стороннего API не выполнено")
		return userData, nil
	}
	log.Println("Приступили к декодированию входящего JSON стороннего API")
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		log.Println("декодирование входящего JSON стороннего API не выполнено")
		return userData, nil
	}
	log.Println("Выполнено обогощение данных на стороннем API")
	return storage.EnrichedUser{}, nil
}
