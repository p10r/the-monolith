package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/p10r/pedro/serve/domain"
	"log"
	"log/slog"
	"net"
	"net/http"
	"time"
)

type Client struct {
	http    *http.Client
	fullUrl string
	log     *slog.Logger
}

func NewClient(fullUrl string, log *slog.Logger) *Client {
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

	l := log.With(slog.String("adapter", "discord"))
	return &Client{c, fullUrl, l}
}

func (c Client) SendMessage(_ context.Context, matches domain.Matches, now time.Time) error {
	msg := NewMessage(matches, now)

	payload, err := json.Marshal(msg)
	if err != nil {
		log.Fatal(err)
		return err
	}

	res, err := http.Post(c.fullUrl, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
		return err
	}

	if res.StatusCode != http.StatusNoContent {
		log.Printf("Discord request failed with status code: %v\n", res.StatusCode)
		return fmt.Errorf("request failed with status code: %v", res.StatusCode)
	}

	//TODO send response for error handling
	return nil
}
