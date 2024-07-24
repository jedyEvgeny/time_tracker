package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/jedyEvgeny/time_tracker/pkg/storage"
)

func (s *Service) DecodeJSONTask(r *http.Request) (storage.TaskEnrichedUser, error) {
	var userTask storage.TaskEnrichedUser
	log.Println("Приступили к декодированию входящего JSON задачи")
	err := json.NewDecoder(r.Body).Decode(&userTask)
	if err != nil {
		log.Println("не распознан JSON задачи")
		return userTask, err
	}
	log.Println("Закончили декодирование входящего JSON пользователя")
	if userTask.TaskName == "" {
		err := errors.New("нет наименования задачи")
		return userTask, err
	}
	err = checkJsonPassport(userTask)
	if err != nil {
		return userTask, err
	}
	log.Println("Проверки корректности JSON выполнены успешно")
	return userTask, nil
}

func checkJsonPassport(u storage.TaskEnrichedUser) error {
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

func (s *Service) DecodeJsonTaskDur(r *http.Request) (storage.TaskEnrichedUser, error) {
	var userTask storage.TaskEnrichedUser
	log.Println("Приступили к декодированию входящего JSON задачи")
	err := json.NewDecoder(r.Body).Decode(&userTask)
	if err != nil {
		log.Println("не распознан JSON задачи")
		return userTask, err
	}
	log.Println("Закончили декодировать входящий JSON пользователя")
	err = checkJsonPassport(userTask)
	if err != nil {
		return userTask, err
	}
	log.Println("Проверки корректности JSON выполнены успешно")
	return userTask, nil
}

func (s *Service) GetSortTasks(u []storage.UserTask) ([]byte, error) {
	taskDurationsSlice := countDuringEvenTask(u)
	sort.Slice(taskDurationsSlice, func(i, j int) bool {
		return taskDurationsSlice[i].TotalDuration > taskDurationsSlice[j].TotalDuration
	})
	response := getResponse(taskDurationsSlice)
	return response, nil
}

func countDuringEvenTask(u []storage.UserTask) []TaskDuration {
	var taskDurationsSlice []TaskDuration
	for _, task := range u {
		duration := task.EndTime.Sub(task.StartTime)
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60

		taskDurationsSlice = append(taskDurationsSlice, TaskDuration{
			TaskName:      task.TaskName,
			Hours:         hours,
			Minutes:       minutes,
			TotalDuration: duration,
		})
	}
	return taskDurationsSlice
}

func getResponse(taskDurationsSlice []TaskDuration) []byte {
	var response []byte
	for _, task := range taskDurationsSlice {
		taskInfo := fmt.Sprintf("задача '%s' занимает: %d часов - %d минут", task.TaskName, task.Hours, task.Minutes)
		response = append(response, taskInfo...)
		response = append(response, '\n')
	}
	return response
}
