//Формируем ответ клиенту

package endpoint

import (
	"net/http"

	"github.com/jedyEvgeny/time_tracker/pkg/logger"
)

type Endpoint struct {
	dec Decoder
	adr Adder
	edc EndpointCaller
}

func New(d Decoder, a Adder, c EndpointCaller) *Endpoint {
	return &Endpoint{
		dec: d,
		adr: a,
		edc: c,
	}
}

var (
	msgErrJSON = "не распознаны номер и серия пользователя"
)

func (e *Endpoint) StatusAdd(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Получили запрос от клиента")
	if r.Method != "POST" {
		w.Write([]byte("метод запроса должен быть POST"))
		logger.Log.Debug("метод запроса не POST, а", r.Method)
		return
	}
	serie, number, err := e.dec.DecodeJSON(r)
	if err != nil {
		http.Error(w, msgErrJSON, http.StatusBadRequest)
		logger.Log.Debug(msgErrJSON, err)
		return
	}
	resp, err := e.edc.CallEndpoint(serie, number)
	if err != nil {
		logger.Log.Debug("неудача при выполнении GET-запроса на эндпоинт /info", err)
	} else {
		defer resp.Body.Close()
		logger.Log.Debug("получен ответ от info эндпоинта: ", resp.StatusCode)
	}
	enrichedUserData, err := e.dec.EnrichUserData(resp, serie, number)
	if err != nil {
		return
	}
	err = e.adr.AddPerson(enrichedUserData)
	if err != nil {
		http.Error(w, "неудача при добавлении данных в БД", http.StatusInternalServerError)
		logger.Log.Debug("неудача при добавлении данных в БД", err)
		return
	}
	logger.Log.Info("Пользователь успешно добавлен в БД, паспорт:", serie, number)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Пользователь успешно добавлен в БД"))
}

func (e *Endpoint) StatusDel(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Получили запрос от клиента")
	if r.Method != "DELETE" {
		w.Write([]byte("метод запроса должен быть DELETE"))
		logger.Log.Debug("метод запроса не DELETE, а", r.Method)
		return
	}
	serie, number, err := e.dec.DecodeJSON(r)
	if err != nil {
		http.Error(w, msgErrJSON, http.StatusBadRequest)
		logger.Log.Debug(msgErrJSON, err)
		return
	}
	err = e.adr.DelPerson(serie, number)
	if err != nil {
		http.Error(w, "неудача при удалении данных из БД", http.StatusInternalServerError)
		logger.Log.Debug("неудача при удалении данных из БД", err)
		return
	}
	logger.Log.Info("Пользователь успешно удалён из БД, паспорт:", serie, number)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Пользователь успешно удалён из БД"))
}

func (e *Endpoint) StatusChange(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Получили запрос от клиента")
	if r.Method != "PATCH" {
		w.Write([]byte("метод запроса должен быть PATCH"))
		logger.Log.Debug("метод запроса не PATCH, а", r.Method)
		return
	}
	enrichedUserData, err := e.dec.ChangeUserData(r)
	if err != nil {
		http.Error(w, msgErrJSON, http.StatusBadRequest)
		return
	}
	err = e.adr.ChangePerson(enrichedUserData)
	if err != nil {
		http.Error(w, "неудача при изменении данных в БД", http.StatusInternalServerError)
		logger.Log.Debug("неудача при изменении данных в БД", err)
		return
	}
	logger.Log.Info("Информация о пользователе успешно изменена в БД, паспорт:", enrichedUserData.PassportSerie, enrichedUserData.PassportNumber)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Пользователь успешно обновлён в БД"))
}

func (e *Endpoint) StatusFilter(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Получили запрос от клиента")
	if r.Method != "GET" {
		w.Write([]byte("метод запроса должен быть GET"))
		logger.Log.Debug("метод запроса не GET, а", r.Method)
		return
	}
	enrichedUserData, err := e.dec.ChangeUserData(r)
	if err != nil {
		http.Error(w, "не распознано тело запроса", http.StatusBadRequest)
		return
	}
	sliceUsers, err := e.adr.GetUsersByFilter(enrichedUserData)
	if err != nil {
		http.Error(w, "неудача при получении данных из БД", http.StatusInternalServerError)
		logger.Log.Debug("неудача при получении данных из БД", err)
		return
	}
	logger.Log.Info("Перечень пользователей успешно получен из БД")
	response, err := e.dec.GetPudding(sliceUsers)
	if err != nil {
		http.Error(w, "не удалось преобразовать данные в JSON", http.StatusInternalServerError)
		logger.Log.Debug("не удалось преобразовать данные в JSON", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (e *Endpoint) StatusStart(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Получили запрос от клиента")
	if r.Method != "PUT" {
		w.Write([]byte("метод запроса должен быть PUT"))
		logger.Log.Debug("метод запроса должен быть PUT, а не", r.Method)
		return
	}
	userTask, err := e.dec.DecodeJSONTask(r)
	if err != nil {
		http.Error(w, "не распознан JSON в задаче", http.StatusBadRequest)
		logger.Log.Debug("не распознан JSON в задаче")
		http.Error(w, msgErrJSON, http.StatusBadRequest)
		return
	}
	startTimeTask, err := e.dec.Now()
	if err != nil {
		logger.Log.Debug("не удалось вычислить текущее время")
		http.Error(w, "не удалось вычислить текущее время", http.StatusInternalServerError)
		return
	}
	logger.Log.Info("Отсчёт времени начат по задаче: ", userTask.TaskName)

	err = e.adr.AddStartTask(userTask, startTimeTask)
	if err != nil {
		logger.Log.Debug("не удалось добавить старт задачи в БД")
		http.Error(w, "не удалось добавить старт задачи в БД", http.StatusInternalServerError)
		return
	}
	logger.Log.Info("Старт задачи добавлен в БД")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Старт задачи добавлен в БД"))
}

func (e *Endpoint) StatusFinish(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Получили запрос от клиента")
	if r.Method != "PUT" {
		w.Write([]byte("метод запроса должен быть PUT"))
		logger.Log.Debug("метод запроса не PUT, а", r.Method)
		return
	}
	userTask, err := e.dec.DecodeJSONTask(r)
	if err != nil {
		http.Error(w, "не распознан JSON в задаче", http.StatusBadRequest)
		logger.Log.Debug("не распознан JSON в задаче")
		return
	}
	finishTimeTask, err := e.dec.Now()
	if err != nil {
		logger.Log.Debug("не удалось вычислить текущее время")
		http.Error(w, "не удалось вычислить текущее время", http.StatusInternalServerError)
		return
	}
	logger.Log.Info("Вычислено время окончания задачи: ", userTask.TaskName)

	err = e.adr.AddFinishTask(userTask, finishTimeTask)
	if err != nil {
		http.Error(w, "не удалось добавить окончание задачи в БД", http.StatusInternalServerError)
		return
	}
	logger.Log.Info("Окончание задачи добавлено в БД")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Время окончания задачи добавлено в БД"))
}

func (e *Endpoint) StatusDur(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Получили запрос от клиента")
	if r.Method != "GET" {
		w.Write([]byte("метод запроса должен быть GET"))
		logger.Log.Debug("метод запроса не GET, а", r.Method)
		return
	}
	userTask, err := e.dec.DecodeJsonTaskDur(r)
	if err != nil {
		http.Error(w, "не распознан JSON в задаче", http.StatusBadRequest)
		logger.Log.Debug("не распознан JSON в задаче")
		return
	}
	userTaskDur, err := e.adr.FindTimeTask(userTask)
	if err != nil {
		logger.Log.Debug("не получено время старта и финиша задачи")
		http.Error(w, "не получено время старта и финиша задачи", http.StatusInternalServerError)
		return
	}
	logger.Log.Info("Получены данные по задачам пользователя из БД")
	response, err := e.dec.GetSortTasks(userTaskDur)
	if err != nil {
		logger.Log.Debug("не удалось преобразовать данные из БД в JSON")
		http.Error(w, "не удалось преобразовать данные из БД в JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
