package httpclient

import (
	"log"
	"net/http"

	config "github.com/jedyEvgeny/time_tracker/etc"
)

type Client struct{}

func New() *Client {
	return &Client{}
}

func (c *Client) CallEndpoint(serie, number string) (*http.Response, error) {
	config, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	log.Println(config.HTTPClientHost)
	baseUrl := "https://" + config.HTTPClientHost + "/info"
	infoEndpoint := baseUrl + "?passportSerie=" + serie + "&passportNumber=" + number
	log.Println("Сформирован get-запрос: ", infoEndpoint)
	resp, err := http.Get(infoEndpoint)
	return resp, err
}
