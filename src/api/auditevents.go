package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type AuditEvent struct {
	UUID       string    `json:"uuid"`
	Timestamp  time.Time `json:"timestamp"`
	ActorUUID  string    `json:"actor_uuid"`
	Action     string    `json:"action"`
	ObjectType string    `json:"object_type"`
	ObjectUUID string    `json:"object_uuid"`
	AuxID      int64     `json:"aux_id,omitempty"`
	AuxUUID    string    `json:"aux_uuid,omitempty"`
	AuxInfo    string    `json:"aux_info,omitempty"`

	Session  AuditEventSession   `json:"session"`
	Location *AuditEventLocation `json:"location,omitempty"`
}

type AuditEventSession struct {
	UUID       string    `json:"uuid"`
	LoginTime  time.Time `json:"login_time"`
	DeviceUUID string    `json:"device_uuid"`
	IP         string    `json:"ip"`
}

type AuditEventLocation struct {
	Country   string  `json:"country,omitempty"`
	Region    string  `json:"region,omitempty"`
	City      string  `json:"city,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

type AuditEventsResponse struct {
	*CursorResponse
	AuditEvents []AuditEvent `json:"items"`
}

func (s *AuditEventsResponse) PrintEvents() error {
	for i, v := range s.AuditEvents {
		raw, err := json.Marshal(v)
		if err != nil {
			err := fmt.Errorf("could not marshal event: %d, error: %s", i, err)
			return err
		}
		fmt.Println(string(raw))
	}
	return nil
}

func (e *EventsAPI) AuditEventsRequest(ctx context.Context, body interface{}) (*AuditEventsResponse, error) {
	res, err := e.request(ctx, "POST", "/api/v1/auditevents", body)
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

	auditEventsRes := &AuditEventsResponse{}
	err = json.Unmarshal(resBody, auditEventsRes)
	if err != nil {
		err := fmt.Errorf("could not unmarshal response: %s", string(resBody))
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err := fmt.Errorf("received a non 200 response: %v", string(resBody))
		return nil, err
	}

	return auditEventsRes, nil
}
