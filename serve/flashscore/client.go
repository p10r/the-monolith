package flashscore

import (
	"errors"
	"fmt"
	"github.com/p10r/monolith/pkg/l"
	"github.com/p10r/monolith/serve/domain"
	"log/slog"
	"net"
	"net/http"
	"time"
)

type Client struct {
	http    *http.Client
	baseUri string
	apiKey  string
	log     *slog.Logger
	retries int
}

func NewClient(baseUri, apiKey string, log *slog.Logger) *Client {
	c := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	fsl := log.With(slog.String("adapter", "flashscore"))
	r := 5
	return &Client{c, baseUri, apiKey, fsl, r}
}

func (c *Client) GetUpcomingMatches() (domain.Matches, error) {
	url := c.baseUri + "/v1/events/list?locale=en_GB&timezone=-4&sport_id=12&indent_days=0"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return domain.Matches{}, fmt.Errorf("err: %w, uri: %s", err, url)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-RapidAPI-Host", "flashscore.p.rapidapi.com")
	req.Header.Add("X-RapidAPI-Key", c.apiKey)

	retries := c.retries
	var res *http.Response
	for retries > 0 {
		res, err = c.http.Do(req)
		if err != nil {
			c.log.Error(l.Error("Error executing GET request", err))
			retries -= 1
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	if res == nil {
		c.log.Error(
			"flashscore res was nil",
			slog.Any("error", err),
			slog.Any("retries_left", retries),
		)
		return domain.Matches{}, fmt.Errorf("res was nil")
	}

	if res.StatusCode == http.StatusForbidden {
		c.log.Error("Forbidden - wrong API key?")
		return domain.Matches{}, err
	}

	if res.StatusCode != http.StatusOK {
		//Todo have Adapter logger that has "data" field
		c.log.Error(fmt.Sprintf("Request failed: %v, req: %v", res.StatusCode, req))
		err := fmt.Errorf("req failed: %v, body: %v", res.StatusCode, res.Body)
		return domain.Matches{}, err
	}

	if res.Body == nil {
		return domain.Matches{}, errors.New("no response body")
	}
	defer res.Body.Close()

	response, err := NewResponse(res.Body)
	if res.Body == nil {
		return domain.Matches{}, errors.New("could not parse JSON")
	}

	return response.ToMatches(), err
}
