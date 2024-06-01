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
	retries int
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
	retries := 5
	l := log.With(slog.String("adapter", "discord"))

	return &Client{c, fullUrl, l, retries}
}

func (c *Client) SendMessage(
	_ context.Context,
	matches domain.Matches,
	now time.Time,
) error {
	msg := NewMessage(matches, now)

	payload, err := json.Marshal(msg)
	if err != nil {
		log.Fatal(err)
		return err
	}

	var res *http.Response
	for c.retries > 0 {
		res, err = http.Post(c.fullUrl, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			c.log.Error("cannot send discord message", slog.Any("error", err))
			c.retries -= 1
		} else {
			break
		}
	}

	if err != nil {
		return err
	}

	if res == nil {
		c.log.Error("discord res was nil", slog.Any("error", err))
		return fmt.Errorf("discord res was nil")
	}

	if res.StatusCode != http.StatusNoContent {
		c.log.Error("req failed with status code", slog.Any("error", err))
		log.Printf("Discord request failed with status code: %v\n", res.StatusCode)
		return fmt.Errorf("request failed with status code: %v", res.StatusCode)
	}

	//TODO send response for error handling
	return nil
}
