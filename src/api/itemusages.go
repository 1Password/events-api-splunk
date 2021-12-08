package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ItemUsage struct {
	UUID        string          `json:"uuid"`
	Timestamp   FixedFormatTime `json:"timestamp"`
	UsedVersion uint32          `json:"used_version"`
	VaultUUID   string          `json:"vault_uuid"`
	ItemUUID    string          `json:"item_uuid"`
	User        ItemUsageUser   `json:"user"`
	Client      ItemUsageClient `json:"client"`
	Action      string          `json:"action"`
}

type ItemUsageUser struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ItemUsageClient struct {
	AppName         string `json:"app_name"`
	AppVersion      string `json:"app_version"`
	PlatformName    string `json:"platform_name"`
	PlatformVersion string `json:"platform_version"`
	OSName          string `json:"os_name"`
	OSVersion       string `json:"os_version"`
	IPAddress       string `json:"ip_address"`
}

type ItemUsageResponse struct {
	*CursorResponse
	Items []ItemUsage `json:"items"`
}

func (s *ItemUsageResponse) PrintEvents() error {
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

func (e *EventsAPI) ItemUsagesRequest(ctx context.Context, body interface{}) (*ItemUsageResponse, error) {
	res, err := e.request(ctx, "POST", "/api/v1/itemusages", body)
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

	itemUsagesRes := &ItemUsageResponse{}
	err = json.Unmarshal(resBody, itemUsagesRes)
	if err != nil {
		err := fmt.Errorf("could not unmarshal response: %s", string(resBody))
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err := fmt.Errorf("received a non 200 response: %v", string(resBody))
		return nil, err
	}

	return itemUsagesRes, nil
}
