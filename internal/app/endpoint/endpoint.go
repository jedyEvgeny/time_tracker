//Формируем ответ клиенту

package endpoint

import (
	"log"
	"net/http"
)

type Decoder interface {
	DecodeJSON(*http.Request) (string, string, error)
}

type Adder interface {
	AddPerson(string, string) error
}

type EndpointCaller interface {
	CallEndpoint(string, string) (*http.Response, error)
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
	serie, number, err := e.dec.DecodeJSON(r)
	if err != nil {
		http.Error(w, "нераспознан JSON", http.StatusBadRequest)
		log.Println("неудача в распознавании JSON", err)
		return
	}

	//        https://localhost/info?passportSerie=6776&passportNumber=614544
	//https://editor.swagger.io/info?passportSerie=1234&passportNumber=567890

	infoEndpoint := "https://localhost/info" + "?passportSerie=" + serie + "&passportNumber=" + number
	log.Println("Сформирован get-запрос: ", infoEndpoint)
	resp, err := http.Get(infoEndpoint)
	if err != nil {
		http.Error(w, "неудача при выполнении GET-запроса на эндпоинт /info", http.StatusInternalServerError)
		log.Println("неудача при выполнении GET-запроса на эндпоинт /info", err)
	} else {
		defer resp.Body.Close()
		log.Printf("Получен ответ от info эндпоинта: %d\n", resp.StatusCode)
	}

	err = e.adr.AddPerson(serie, number)
	if err != nil {
		http.Error(w, "неудача при добавлении данных в БД", http.StatusInternalServerError)
		log.Println("неудача при добавлении данных в БД", err)
		return
	}
	log.Printf("Пользователь с паспортом серии %v и номером %v успешно добавлен в БД\n", serie, number)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Пользователь успешно добавлен в БД"))
}
