package httpclient

import (
	"net/http"

	config "github.com/jedyEvgeny/time_tracker/etc"
	"github.com/jedyEvgeny/time_tracker/pkg/logger"
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
	baseUrl := "https://" + config.HTTPClientHost + "/info"
	infoEndpoint := baseUrl + "?passportSerie=" + serie + "&passportNumber=" + number
	logger.Log.Info("Сформирован get-запрос: ", infoEndpoint)
	resp, err := http.Get(infoEndpoint)
	return resp, err
}
