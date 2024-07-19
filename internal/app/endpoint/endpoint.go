//Формируем ответ клиенту

package endpoint

import (
	"log"
	"net/http"

	"github.com/jedyEvgeny/time_tracker/internal/database"
)

type Decoder interface {
	DecodeJSON(*http.Request) (database.User, error)
}

type Adder interface {
	AddPerson(database.User) error
}

type Endpoint struct {
	dcr Decoder
	adr Adder
}

func New(d Decoder, a Adder) *Endpoint {
	return &Endpoint{
		dcr: d,
		adr: a,
	}
}

func (e *Endpoint) Status(w http.ResponseWriter, r *http.Request) {
	userData, err := e.dcr.DecodeJSON(r)
	if err != nil {
		http.Error(w, "нераспознан JSON", http.StatusBadRequest)
		return
	}

	err = e.adr.AddPerson(userData)
	if err != nil {
		http.Error(w, "неудача при добавлении данных в БД", http.StatusInternalServerError)
		return
	}
	log.Printf("Пользователь с паспортом %v успешно добавлен в БД\n", userData.PassportNumber)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Пользователь успешно добавлен в БД"))
}
