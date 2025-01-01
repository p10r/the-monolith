package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/p10r/pedro/pkg/l"
	"github.com/p10r/pedro/serve/domain"
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
	dl := log.With(slog.String("adapter", "discord"))

	return &Client{c, fullUrl, dl, retries}
}

func (c *Client) SendUpcomingMatches(
	_ context.Context,
	matches domain.Matches,
	now time.Time,
) error {
	msg := NewUpcomingMatchesMsg(matches, now)

	payload, err := json.Marshal(msg)
	if err != nil {
		c.log.Error(l.Error("cannot marshal discord message", err))
		return err
	}

	var res *http.Response
	for c.retries > 0 {
		res, err = http.Post(c.fullUrl, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			c.log.Error(l.Error("cannot send discord message", err))
			c.retries -= 1
		} else {
			break
		}
	}

	if err != nil {
		return err
	}

	if res == nil {
		c.log.Error(l.Error("discord res was nil", err))
		return fmt.Errorf("discord res was nil")
	}

	if res.StatusCode != http.StatusNoContent {
		c.log.Error(l.Error("req failed with status code", err))
		return fmt.Errorf("request failed with status code: %v", res.StatusCode)
	}

	//TODO send response for error handling
	return nil
}

func (c *Client) SendFinishedMatches(
	_ context.Context,
	matches domain.FinishedMatchesByLeague,
	now time.Time,
) error {
	msg := NewFinishedMatchesMsg(matches, now)

	payload, err := json.Marshal(msg)
	if err != nil {
		c.log.Error(l.Error("cannot marshal discord message", err))
		return err
	}

	var res *http.Response
	for c.retries > 0 {
		res, err = http.Post(c.fullUrl, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			c.log.Error(l.Error("cannot send discord message", err))
			c.retries -= 1
		} else {
			break
		}
	}

	if err != nil {
		return err
	}

	if res == nil {
		c.log.Error(l.Error("discord res was nil", err))
		return fmt.Errorf("discord res was nil")
	}

	if res.StatusCode != http.StatusNoContent {
		c.log.Error(l.Error("req failed with status code", err))
		return fmt.Errorf("request failed with status code: %v", res.StatusCode)
	}

	//TODO send response for error handling
	return nil
}
