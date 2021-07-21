package splunk

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type SplunkAPI struct {
	client    *http.Client
	SessionKey string
	BaseUrl   string
}

const DefaultClientTimeout = 15 * time.Second
const DefaultUserAgent = "1Password Insights / 1.5.0"

func NewSplunkAPI(sessionKey string) *SplunkAPI {
	log.Println("New Splunk API")

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	c := &http.Client{
		Timeout: DefaultClientTimeout,
		Transport: transport,
	}
	
	client := &SplunkAPI{
		client:    c,
		SessionKey: sessionKey,
		BaseUrl: "https://localhost:8089", // Probably shouldn't hard code
	}
	return client
}

func (e *SplunkAPI) request(ctx context.Context, method string, route string, body interface{}) (*http.Response, error) {
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
	req.Header.Add("Authorization", fmt.Sprintf("Splunk %s", e.SessionKey))
	req.Header.Add("Content-Type", "application/json")
	res, err := e.client.Do(req)
	if err != nil {
		err := fmt.Errorf("could not make request: %w", err)
		return nil, err
	}
	return res, nil
}
