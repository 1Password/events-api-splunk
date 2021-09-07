package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SignInAttempt struct {
	UUID        string                  `json:"uuid"`
	SessionUUID string                  `json:"session_uuid"`
	Timestamp   FixedFormatTime         `json:"timestamp"`
	Country     string                  `json:"country"`
	Category    string                  `json:"category"`
	Type        string                  `json:"type"`
	Details     *SignInAttemptDetails   `json:"details"`
	TargetUser  SignInAttemptTargetUser `json:"target_user"`
	Client      SignInAttemptClient     `json:"client"`
}

type SignInAttemptDetails struct {
	Value string `json:"value"`
}
type SignInAttemptTargetUser struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
type SignInAttemptClient struct {
	AppName         string `json:"app_name"`
	AppVersion      string `json:"app_version"`
	PlatformName    string `json:"platform_name"`
	PlatformVersion string `json:"platform_version"`
	OSName          string `json:"os_name"`
	OSVersion       string `json:"os_version"`
	IPAddress       string `json:"ip_address"`
}

type CursorResponse struct {
	Cursor  string `json:"cursor"`
	HasMore bool   `json:"has_more"`
}

type SignInAttemptResponse struct {
	*CursorResponse
	Items []SignInAttempt `json:"items"`
}

func (s *SignInAttemptResponse) PrintEvents() error {
	for i, v := range s.Items {
		raw, err := json.Marshal(v)
		if err != nil {
			err := fmt.Errorf("could not marshal event: %d, error: %s", i, err)
			return err
		}
		fmt.Println(string(raw))
	}
	return nil
}

func (e *EventsAPI) SignInAttemptsRequest(ctx context.Context, body interface{}) (*SignInAttemptResponse, error) {
	res, err := e.request(ctx, "POST", "/api/v1/signinattempts", body)
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

	attemptsRes := &SignInAttemptResponse{}
	err = json.Unmarshal(resBody, attemptsRes)
	if err != nil {
		err := fmt.Errorf("could not unmarshal response: %s", string(resBody))
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err := fmt.Errorf("received a non 200 response: %v", string(resBody))
		return nil, err
	}

	return attemptsRes, nil
}
