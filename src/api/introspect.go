package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		err := fmt.Errorf("could not make EventAPIRequest: %w", err)
		return nil, err
	}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err := fmt.Errorf("could not read response: %w", err)
		return nil, err
	}
	res.Body.Close()

	introspectRes := &IntrospectResponse{}
	err = json.Unmarshal(resBody, introspectRes)
	if err != nil {
		err := fmt.Errorf("could not unmarshal response: %s", string(resBody))
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err := fmt.Errorf("received a non 200 response: %v", string(resBody))
		return nil, err
	}

	return introspectRes, nil
}
