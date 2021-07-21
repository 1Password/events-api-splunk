package actions

import (
	"context"
	"fmt"
	"log"

	"go.1password.io/eventsapi-splunk/splunk"
)

func GetEventsToken(ctx context.Context, splunkAPI *splunk.SplunkAPI) (string, error)  {
	log.Println("Calling Splunk API")
	
	splunkRes, err := splunkAPI.Passwords(ctx, "events_reporting_realm", "events_api_token")
	if err != nil {
		err := fmt.Errorf("call to splunk failed: %w", err)
		return "", err
	}

	return splunkRes.Entry.Content.Dict.Key[0].Text, nil
}
