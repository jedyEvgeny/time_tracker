package httpclient

import (
	"log"
	"net/http"
)

type Client struct{}

func New() *Client {
	return &Client{}
}

func (c *Client) CallEndpoint(serie, number string) (*http.Response, error) {
	infoEndpoint := "https://localhost/info" + "?passportSerie=" + serie + "&passportNumber=" + number
	log.Println("Сформирован get-запрос: ", infoEndpoint)
	resp, err := http.Get(infoEndpoint)
	return resp, err
}
