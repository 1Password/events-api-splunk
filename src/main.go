package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"go.1password.io/eventsapi-splunk/actions"
	events "go.1password.io/eventsapi-splunk/api"
	"go.1password.io/eventsapi-splunk/config"
	"go.1password.io/eventsapi-splunk/splunk"
	"go.1password.io/eventsapi-splunk/utils"
)

var EventBuildType string // Injected at build time so we can make multiple apps

func main() {
	log.Println("Booting...")
	if EventBuildType == "" {
		err := fmt.Errorf("missing EventBuildType flag")
		panic(err)
	}

	splunkHome := os.Getenv("SPLUNK_HOME")
	if splunkHome == "" {
		err := fmt.Errorf("SPLUNK_HOME environment variable must be set")
		panic(err)
	}

	splunkEnv, err := config.NewSplunkEnv(splunkHome)
	if err != nil {
		err := fmt.Errorf("could not create new splunk env: %w", err)
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)
	splunkSession, _, err := reader.ReadLine()
	if err != nil {
		err := fmt.Errorf("could not read session: %w", err)
		panic(err)
	}

	splunkAPI := splunk.NewSplunkAPI(string(splunkSession))

	// Versions less than 1.5.0 of the Events API stored the token on disk
	// If we find it, move it to the splunk storage/passwords service
	var eventsToken string
	if splunkEnv.Config.AuthToken != "" {
		eventsToken = splunkEnv.Config.AuthToken
		err := actions.CreateEventsToken(context.TODO(), splunkAPI, eventsToken)
		if err != nil {
			err := fmt.Errorf("could not backup token: %w", err)
			panic(err)
		}
		splunkEnv.Config.AuthToken = "" // Remove token on disk
		err = splunkEnv.UpdateConfig(splunkEnv.Config)
		if err != nil {
			err := fmt.Errorf("could not remove auth token: %w", err)
			panic(err)
		}
	} else {
		eventsToken, err = actions.GetEventsToken(context.TODO(), splunkAPI)
		if err != nil {
			err := fmt.Errorf("could not get token: %w", err)
			panic(err)
		}
	}

	jwt, err := utils.ParseJWTClaims(eventsToken)
	if err != nil {
		err := fmt.Errorf("could not parse jwt: %w", err)
		panic(err)
	}

	url, err := jwt.GetEventsURL()
	// The config url will be used if the token was generated before
	// this update and does not contain a url
	if err == nil {
		splunkEnv.Config.Url = url
	}

	eventsAPI := events.NewEventsAPI(eventsToken, url)

	if jwt.Features.Contains(utils.SignInAttemptsFeatureScope) && EventBuildType == utils.SignInAttemptsFeatureScope {
		cursorFile := path.Join(splunkEnv.Home, splunkEnv.Config.SignInCursorFile)
		actions.StartSignIns(cursorFile, splunkEnv.Config.Limit, &splunkEnv.Config.StartAt, eventsAPI)
	} else if jwt.Features.Contains(utils.ItemUsageFeatureScope) && EventBuildType == utils.ItemUsageFeatureScope {
		cursorFile := path.Join(splunkEnv.Home, splunkEnv.Config.ItemUsageCursorFile)
		actions.StartItemUsages(cursorFile, splunkEnv.Config.Limit, &splunkEnv.Config.StartAt, eventsAPI)
	}
}
