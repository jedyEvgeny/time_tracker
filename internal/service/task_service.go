package service

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jedyEvgeny/time_tracker/pkg/storage"
)

func (s *Service) DecodeJSONTask(r *http.Request) (storage.TaskEnrichedUser, error) {
	var userTask storage.TaskEnrichedUser
	log.Println("Приступили к декодированию входящего JSON задачи")
	err := json.NewDecoder(r.Body).Decode(&userTask)
	if err != nil {
		log.Println("не распознано наименование задачи")
		return userTask, err
	}
	log.Println("Закончили декодирование входящего JSON пользователя")
	err = checkJsonTask(userTask)
	if err != nil {
		return userTask, err
	}
	log.Println("Проверки корректности JSON выполнены успешно")
	return userTask, nil
}

func checkJsonTask(u storage.TaskEnrichedUser) error {
	if u.TaskName == "" {
		err := errors.New("нет наименования задачи")
		return err
	}
	if len(u.PassportSerie) != 4 {
		log.Println("ожидалось 4 символов в серии паспорта")
		err := errors.New("неверная длина серии паспорта, ожидается 4 символа")
		return err
	}
	if len(u.PassportNumber) != 6 {
		log.Println("ожидалось 6 символов в номере паспорта")
		err := errors.New("неверная длина номера паспорта, ожидается 6 символов")
		return err
	}
	elFirst, err := strconv.Atoi(u.PassportSerie)
	if err != nil {
		log.Println("в серии паспорта не только цифры")
		err := errors.New("в номере паспорта не только цифры")
		return err
	}
	elSecond, err := strconv.Atoi(u.PassportNumber)
	if err != nil {
		log.Println("во втором блоке тела запроса не только цифры")
		err := errors.New("во втором блоке тела запроса не только цифры")
		return err
	}
	if elFirst < 0 || elSecond < 0 {
		log.Println("цифровой блок со знаком минус не допустим")
		err := errors.New("цифровой блок со знаком минус не допустим")
		return err
	}
	return nil
}

func (s *Service) Now() (time.Time, error) {
	vladivostokLocation, err := time.LoadLocation("Asia/Vladivostok")
	if err != nil {
		log.Println("ошибка при загрузке локации:", err)
		return time.Time{}, err
	}
	return time.Now().In(vladivostokLocation), nil
}
