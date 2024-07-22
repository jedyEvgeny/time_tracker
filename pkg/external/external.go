package external

import (
	"log"
	"net/http"
)

type External struct{}

func New() *External {
	return &External{}
}

func (e *External) CallEndpoint(serie, number, string) (*http.Response, error) {
	infoEndpoint := "https://localhost/info" + "?passportSerie=" + serie + "&passportNumber=" + number
	log.Println("Сформирован get-запрос: ", infoEndpoint)
	resp, err := http.Get(infoEndpoint)
	return resp, err
}
