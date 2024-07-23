//Формируем ответ клиенту

package endpoint

import (
	"log"
	"net/http"

	"github.com/jedyEvgeny/time_tracker/pkg/storage"
)

type Decoder interface {
	DecodeJSON(*http.Request) (string, string, error)
	EnrichUserData(*http.Response, string, string) (storage.EnrichedUser, error)
	ChangeUserData(*http.Request) (storage.EnrichedUser, error)
}

type Adder interface {
	AddPerson(storage.EnrichedUser) error
	DelPerson(string, string) error
	ChangePerson(storage.EnrichedUser) error
}

type EndpointCaller interface {
	CallEndpoint(string, string) (*http.Response, error)
}

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
		log.Println("метод запроса не POST, а", r.Method)
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
		log.Println("метод запроса не POST, а", r.Method)
		return
	}
	// serie, number, err := e.dec.DecodeJSON(r)
	// if err != nil {
	// 	http.Error(w, msgErrJSON, http.StatusBadRequest)
	// 	log.Println(msgErrJSON, err)
	// 	return
	// }

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
