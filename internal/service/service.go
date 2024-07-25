// Реализация бизнес-логики
package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/jedyEvgeny/time_tracker/pkg/logger"
	"github.com/jedyEvgeny/time_tracker/pkg/storage"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) DecodeJSON(r *http.Request) (string, string, error) {
	var userData storage.User
	logger.Log.Info("Приступили к декодированию входящего JSON пользователя")
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		return "", "", err
	}
	logger.Log.Info("Закончили декодирование входящего JSON пользователя")
	serie, number, err := checkJSON(userData)
	if err != nil {
		return "", "", err
	}
	logger.Log.Info("Проверки корректности JSON выполнены успешно")
	return serie, number, nil
}

func checkJSON(u storage.User) (string, string, error) {
	if len(u.PassportNumber) != 11 {
		logger.Log.Debug("ожидалось 11 символов в запросе")
		err := errors.New("неверная длина строки, ожидается 11 символов")
		return "", "", err
	}
	parts := strings.Split(u.PassportNumber, " ")
	if len(parts) != 2 {
		logger.Log.Debug("ожидался один пробел в запросе")
		err := errors.New("ожидался один пробел в запросе")
		return "", "", err
	}
	if len(parts[0]) != 4 && len(parts[1]) != 6 {
		logger.Log.Debug("неверная длина блока серии или номера паспорта")
		err := errors.New("неверная длина блока серии или номера паспорта")
		return "", "", err
	}
	elFirst, err := strconv.Atoi(parts[0])
	if err != nil {
		logger.Log.Debug("в первом блоке тела запроса не только цифры")
		err := errors.New("в первом блоке тела запроса не только цифры")
		return "", "", err
	}
	elSecond, err := strconv.Atoi(parts[1])
	if err != nil {
		logger.Log.Debug("во втором блоке тела запроса не только цифры")
		err := errors.New("во втором блоке тела запроса не только цифры")
		return "", "", err
	}
	if elFirst < 0 || elSecond < 0 {
		logger.Log.Debug("цифровой блок со знаком минус не допустим")
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
		logger.Log.Debug("Обогощение данных со стороннего API не выполнено")
		return userData, nil
	}
	logger.Log.Info("Приступили к декодированию входящего JSON стороннего API")
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		logger.Log.Debug("декодирование входящего JSON стороннего API не выполнено")
		return userData, err
	}
	logger.Log.Info("Выполнено обогощение данных на стороннем API")
	return userData, nil
}

func (s *Service) ChangeUserData(req *http.Request) (storage.EnrichedUser, error) {
	var userData storage.EnrichedUser
	logger.Log.Info("Приступили к обогощению данных из JSON")
	err := json.NewDecoder(req.Body).Decode(&userData)
	if err != nil {
		logger.Log.Debug("декодирование обогощённого JSON не выполнено")
		return userData, err
	}
	logger.Log.Info("Выполнено обогощение данных из JSON")
	return userData, nil
}

func (s *Service) GetPudding(users []storage.EnrichedUser) ([]byte, error) {
	logger.Log.Info("приступили к пагинации")
	amountElemOnPage := 2 // Количество позиций в ответе клиенту на одной странице
	startPageDisplay := 1 // Отображаем ответ клиенту с первой страницы
	// В дальнейшем возможна реализация обработчика для обновления информации на новой странице
	// При получении get-запроса другого номера страницы
	startIdx := (startPageDisplay - 1) * amountElemOnPage
	endIdx := startIdx + amountElemOnPage
	if startIdx >= len(users) {
		return nil, errors.New("вне диапазона")
	}
	if endIdx > len(users) {
		endIdx = len(users)
	}
	sortUsersByID(users)
	slicedUsers := users[startIdx:endIdx]

	response, err := json.Marshal(slicedUsers)
	if err != nil {
		logger.Log.Debug("не удалось преобразовать данные из БД в JSON")
		return nil, err
	}
	return response, nil
}

func sortUsersByID(users []storage.EnrichedUser) {
	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})
}
