//Формируем ответ клиенту

package endpoint

import (
	"log"
	"net/http"

	"github.com/jedyEvgeny/time_tracker/pkg/storage"
)

type Decoder interface {
	DecodeJSON(*http.Request) (storage.User, error)
}

type Adder interface {
	AddPerson(storage.User) error
}

type Endpoint struct {
	dec Decoder
	adr Adder
}

func New(d Decoder, a Adder) *Endpoint {
	return &Endpoint{
		dec: d,
		adr: a,
	}
}

func (e *Endpoint) Status(w http.ResponseWriter, r *http.Request) {
	log.Println("Получили запрос от клиента")
	userData, err := e.dec.DecodeJSON(r)
	if err != nil {
		http.Error(w, "нераспознан JSON", http.StatusBadRequest)
		log.Println("неудача в распозновании JSON")
		return
	}

	err = e.adr.AddPerson(userData)
	if err != nil {
		http.Error(w, "неудача при добавлении данных в БД", http.StatusInternalServerError)
		log.Println("неудача при добавлении данных в БД")
		return
	}
	log.Printf("Пользователь с паспортом %v успешно добавлен в БД\n", userData.PassportNumber)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Пользователь успешно добавлен в БД"))
}
