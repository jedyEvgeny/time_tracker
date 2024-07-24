//Формируем ответ клиенту

package endpoint

import (
	"log"
	"net/http"
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
	log.Println("Получили запрос от клиента")
	if r.Method != "POST" {
		w.Write([]byte("метод запроса должен быть POST"))
		log.Println("метод запроса не POST, а", r.Method)
		return
	}
	serie, number, err := e.dec.DecodeJSON(r)
	if err != nil {
		http.Error(w, msgErrJSON, http.StatusBadRequest)
		log.Println(msgErrJSON, err)
		return
	}
	resp, err := e.edc.CallEndpoint(serie, number)
	if err != nil {
		http.Error(w, "неудача при выполнении GET-запроса на эндпоинт /info", http.StatusInternalServerError)
		log.Println("неудача при выполнении GET-запроса на эндпоинт /info", err)
	} else {
		defer resp.Body.Close()
		log.Printf("Получен ответ от info эндпоинта: %d\n", resp.StatusCode)
	}
	enrichedUserData, err := e.dec.EnrichUserData(resp, serie, number)
	if err != nil {
		return
	}
	err = e.adr.AddPerson(enrichedUserData)
	if err != nil {
		http.Error(w, "неудача при добавлении данных в БД", http.StatusInternalServerError)
		log.Println("неудача при добавлении данных в БД", err)
		return
	}
	log.Printf("Пользователь с паспортом серии %v и номером %v успешно добавлен в БД\n", serie, number)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Пользователь успешно добавлен в БД"))
}

func (e *Endpoint) StatusDel(w http.ResponseWriter, r *http.Request) {
	log.Println("Получили запрос от клиента")
	if r.Method != "DELETE" {
		w.Write([]byte("метод запроса должен быть DELETE"))
		log.Println("метод запроса не DELETE, а", r.Method)
		return
	}
	serie, number, err := e.dec.DecodeJSON(r)
	if err != nil {
		http.Error(w, msgErrJSON, http.StatusBadRequest)
		log.Println(msgErrJSON, err)
		return
	}
	err = e.adr.DelPerson(serie, number)
	if err != nil {
		http.Error(w, "неудача при удалении данных из БД", http.StatusInternalServerError)
		log.Println("неудача при удалении данных из БД", err)
		return
	}
	log.Printf("Пользователь с паспортом серии %v и номером %v успешно удалён из БД\n", serie, number)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Пользователь успешно удалён из БД"))
}

func (e *Endpoint) StatusChange(w http.ResponseWriter, r *http.Request) {
	log.Println("Получили запрос от клиента")
	if r.Method != "PATCH" {
		w.Write([]byte("метод запроса должен быть PATCH"))
		log.Println("метод запроса не PATCH, а", r.Method)
		return
	}
	enrichedUserData, err := e.dec.ChangeUserData(r)
	if err != nil {
		return
	}
	err = e.adr.ChangePerson(enrichedUserData)
	if err != nil {
		http.Error(w, "неудача при изменении данных в БД", http.StatusInternalServerError)
		log.Println("неудача при изменении данных в БД", err)
		return
	}
	log.Printf("Пользователь с паспортом серии %v и номером %v успешно обновлён в БД\n", enrichedUserData.PassportSerie, enrichedUserData.PassportNumber)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Пользователь успешно обновлён в БД"))
}

func (e *Endpoint) StatusFilter(w http.ResponseWriter, r *http.Request) {
	log.Println("Получили запрос от клиента")
	if r.Method != "GET" {
		w.Write([]byte("метод запроса должен быть GET"))
		log.Println("метод запроса не POST, а", r.Method)
		return
	}
	enrichedUserData, err := e.dec.ChangeUserData(r)
	if err != nil {
		return
	}
	sliceUsers, err := e.adr.GetUsersByFilter(enrichedUserData)
	if err != nil {
		http.Error(w, "неудача при получении данных из БД", http.StatusInternalServerError)
		log.Println("неудача при получении данных из БД", err)
		return
	}
	log.Println("Перечень пользователей успешно получен из БД")
	response, err := e.dec.GetPudding(sliceUsers)
	if err != nil {
		http.Error(w, "Не удалось преобразовать данные в JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (e *Endpoint) StatusStart(w http.ResponseWriter, r *http.Request) {
	log.Println("Получили запрос от клиента")
	if r.Method != "PUT" {
		w.Write([]byte("метод запроса должен быть PUT"))
		log.Printf("метод запроса не %v, а PUT\n", r.Method)
		return
	}
	userTask, err := e.dec.DecodeJSONTask(r)
	if err != nil {
		http.Error(w, "не распознан JSON в задаче", http.StatusBadRequest)
		return
	}
	startTimeTask, err := e.dec.Now()
	if err != nil {
		return
	}
	log.Printf("Старт задачи %v в : %v.\n", userTask.TaskName, startTimeTask)

	err = e.adr.AddStartTask(userTask, startTimeTask)
	if err != nil {
		return
	}
	log.Println("Старт задачи добавлен в БД")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Старт задачи добавлен в БД"))
}

func (e *Endpoint) StatusFinish(w http.ResponseWriter, r *http.Request) {
	log.Println("Получили запрос от клиента")
	if r.Method != "GET" {
		w.Write([]byte("метод запроса должен быть GET"))
		log.Println("метод запроса не PUT, а", r.Method)
		return
	}
	userTask, err := e.dec.DecodeJSONTask(r)
	if err != nil {
		http.Error(w, "не распознан JSON в задаче", http.StatusBadRequest)
		return
	}
	finishTimeTask, err := e.dec.Now()
	if err != nil {
		return
	}
	log.Printf("Окончание задачи %v в : %v.\n", userTask.TaskName, finishTimeTask)

	err = e.adr.AddFinishTask(userTask, finishTimeTask)
	if err != nil {
		return
	}
	log.Println("Окончание задачи добавлено в БД")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Время окончания задачи добавлено в БД"))
}

func (e *Endpoint) StatusDur(w http.ResponseWriter, r *http.Request) {
	log.Println("Получили запрос от клиента")
	if r.Method != "GET" {
		w.Write([]byte("метод запроса должен быть GET"))
		log.Println("метод запроса не GET, а", r.Method)
		return
	}
	userTask, err := e.dec.DecodeJsonTaskDur(r)
	if err != nil {
		http.Error(w, "не распознан JSON в задаче", http.StatusInternalServerError)
		return
	}
	userTaskDur, err := e.adr.FindTimeTask(userTask)
	if err != nil {
		return
	}
	log.Println("Получены данные по задачам пользователя из БД")
	response, err := e.dec.GetSortTasks(userTaskDur)
	if err != nil {
		log.Println("не удалось преобразовать данные из БД в JSON")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
