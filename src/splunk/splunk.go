package splunk

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type SplunkAPI struct {
	client     *http.Client
	SessionKey string
	BaseUrl    string
}

const DefaultClientTimeout = 15 * time.Second

func NewSplunkAPI(sessionKey string) *SplunkAPI {
	log.Println("New Splunk API")

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	c := &http.Client{
		Timeout:   DefaultClientTimeout,
		Transport: transport,
	}

	client := &SplunkAPI{
		client:     c,
		SessionKey: sessionKey,
		BaseUrl:    "https://localhost:8089", // Probably shouldn't hard code
	}
	return client
}

func (e *SplunkAPI) request(ctx context.Context, method string, route string, data url.Values) (*http.Response, error) {
	log.Printf("Calling Splunk API: %s", route)

	var b io.Reader
	contentLength := 0
	if data != nil {
		b = strings.NewReader(data.Encode())
		contentLength = len(data.Encode())
	}
	req, err := http.NewRequestWithContext(ctx, method, e.BaseUrl+route, b)
	if err != nil {
		err := fmt.Errorf("could not create new request: %w", err)
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Splunk %s", e.SessionKey))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(contentLength))

	res, err := e.client.Do(req)
	if err != nil {
		err := fmt.Errorf("could not make request: %w", err)
		return nil, err
	}
	return res, nil
}
