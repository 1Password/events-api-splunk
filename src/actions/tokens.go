package actions

import (
	"context"
	"fmt"

	"go.1password.io/eventsapi-splunk/splunk"
)

func GetEventsToken(ctx context.Context, splunkAPI *splunk.SplunkAPI) (string, error)  {
	splunkRes, err := splunkAPI.GetPasswords(ctx, "events_api_token", "events_reporting_realm",)
	if err != nil {
		err := fmt.Errorf("call to splunk failed: %w", err)
		return "", err
	}
	if len(splunkRes.Entry.Content.Dict.Key) == 0 {
		err := fmt.Errorf("splunk response is missing credentials")
		return "", err
	}

	return splunkRes.Entry.Content.Dict.Key[0].Text, nil
}

func CreateEventsToken(ctx context.Context, splunkAPI *splunk.SplunkAPI, authToken string) error {
	err := splunkAPI.CreatePassword(ctx, "events_api_token", authToken, "events_reporting_realm")
	if err != nil {
		err := fmt.Errorf("call to splunk failed: %w", err)
		return err
	}
	return nil
}