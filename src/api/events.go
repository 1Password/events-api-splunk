package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type EventsAPI struct {
	client    *http.Client
	AuthToken string
	BaseUrl   string
}

const DefaultClientTimeout = 15 * time.Second
const DefaultUserAgent = "1Password Events API for Splunk / 1.3.0"

func NewEventsAPI(authToken string, url string) *EventsAPI {
	log.Println("New Events API")
	c := &http.Client{
		Timeout: DefaultClientTimeout,
	}
	client := &EventsAPI{
		client:    c,
		AuthToken: authToken,
		BaseUrl:   url,
	}
	return client
}

func (e *EventsAPI) request(ctx context.Context, method string, route string, body interface{}) (*http.Response, error) {
	var b io.Reader
	if body != nil {
		reqBody, err := json.Marshal(body)
		if err != nil {
			err := fmt.Errorf("could not marshal request: %w", err)
			panic(err)
		}
		b = bytes.NewReader(reqBody)
	}
	req, err := http.NewRequestWithContext(ctx, method, e.BaseUrl+route, b)
	if err != nil {
		err := fmt.Errorf("could not create new request: %w", err)
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", e.AuthToken))
	req.Header.Add("User-Agent", DefaultUserAgent)
	res, err := e.client.Do(req)
	if err != nil {
		err := fmt.Errorf("could not make request: %w", err)
		return nil, err
	}
	return res, nil
}

type CursorRequest struct {
	Cursor string `json:"cursor"`
}
type CursorResetRequest struct {
	Limit     int        `json:"limit"`
	StartTime *time.Time `json:"start_time,omitempty"`
}
