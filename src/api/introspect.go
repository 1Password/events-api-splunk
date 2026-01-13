package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type IntrospectResponse struct {
	UUID     string    `json:"UUID"`
	IssuedAt time.Time `json:"IssuedAt"`
	Features []string  `json:"Features"`
}

func (e *EventsAPI) Introspect(ctx context.Context) (*IntrospectResponse, error) {
	res, err := e.request(ctx, "GET", "/api/auth/introspect", nil)
	if err != nil {
		return nil, fmt.Errorf("could not make EventAPIRequest: %w", err)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %w", err)
	}
	res.Body.Close()

	introspectRes := &IntrospectResponse{}
	err = json.Unmarshal(resBody, introspectRes)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response: %s", string(resBody))
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received a non 200 response: %v", string(resBody))
	}

	return introspectRes, nil
}
