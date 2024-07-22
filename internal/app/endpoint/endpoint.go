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
}

type Adder interface {
	AddPerson(storage.EnrichedUser) error
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

func (e *Endpoint) Status(w http.ResponseWriter, r *http.Request) {
	log.Println("Получили запрос от клиента")
	serie, number, err := e.dec.DecodeJSON(r)
	if err != nil {
		http.Error(w, "нераспознан JSON", http.StatusBadRequest)
		log.Println("неудача в распознавании JSON", err)
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
