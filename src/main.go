package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"

	"go.1password.io/eventsapi-splunk/actions"
	events "go.1password.io/eventsapi-splunk/api"
	"go.1password.io/eventsapi-splunk/config"
	"go.1password.io/eventsapi-splunk/splunk"
	"go.1password.io/eventsapi-splunk/utils"
)

var EventBuildType string // Injected at build time so we can make multiple apps

func main() {
	err := run()
	if err != nil {
		slog.Error("unexpected error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	// Set up slog logger with INFO level by default, DEBUG if LOG_LEVEL=debug
	logLevel := slog.LevelInfo
	if strings.ToLower(os.Getenv("LOG_LEVEL")) == "debug" {
		logLevel = slog.LevelDebug
	}
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, opts))
	slog.SetDefault(logger)

	slog.Info("Booting...")
	if EventBuildType == "" {
		return fmt.Errorf("missing EventBuildType flag")
	}

	splunkHome := os.Getenv("SPLUNK_HOME")
	if splunkHome == "" {
		return fmt.Errorf("SPLUNK_HOME environment variable must be set")
	}

	splunkEnv, err := config.NewSplunkEnv(splunkHome)
	if err != nil {
		return fmt.Errorf("could not create new splunk env: %w", err)
	}

	reader := bufio.NewReader(os.Stdin)
	splunkSession, _, err := reader.ReadLine()
	if err != nil {
		return fmt.Errorf("could not read session: %w", err)
	}

	splunkAPI := splunk.NewSplunkAPI(string(splunkSession))

	// Versions less than 1.5.0 of the Events API stored the token on disk
	// If we find it, move it to the splunk storage/passwords service
	var eventsToken string
	if splunkEnv.Config.AuthToken != "" {
		eventsToken = splunkEnv.Config.AuthToken
		err := actions.CreateEventsToken(context.TODO(), splunkAPI, eventsToken)
		if err != nil {
			return fmt.Errorf("could not backup token: %w", err)
		}
		splunkEnv.Config.AuthToken = "" // Remove token on disk
		err = splunkEnv.UpdateConfig(splunkEnv.Config)
		if err != nil {
			return fmt.Errorf("could not remove auth token: %w", err)
		}
	} else {
		eventsToken, err = actions.GetEventsToken(context.TODO(), splunkAPI)
		if err != nil {
			return fmt.Errorf("could not get token: %w", err)
		}
	}

	jwt, err := utils.ParseJWTClaims(eventsToken)
	if err != nil {
		return fmt.Errorf("could not parse jwt: %w", err)
	}

	url, err := jwt.GetEventsURL()
	// The config url will be used if the token was generated before
	// this update and does not contain a url
	if err == nil {
		splunkEnv.Config.Url = url
	}

	eventsAPI := events.NewEventsAPI(eventsToken, url)

	slog.Info("Starting", "eventType", EventBuildType, "limit", splunkEnv.Config.Limit, "startAt", splunkEnv.Config.StartAt)

	switch {
	case jwt.Features.Contains(utils.SignInAttemptsFeatureScope) && EventBuildType == utils.SignInAttemptsFeatureScope:
		cursorFile := path.Join(splunkEnv.Home, splunkEnv.Config.SignInCursorFile)
		return actions.StartSignIns(cursorFile, splunkEnv.Config.Limit, &splunkEnv.Config.StartAt, eventsAPI)
	case jwt.Features.Contains(utils.ItemUsageFeatureScope) && EventBuildType == utils.ItemUsageFeatureScope:
		cursorFile := path.Join(splunkEnv.Home, splunkEnv.Config.ItemUsageCursorFile)
		return actions.StartItemUsages(cursorFile, splunkEnv.Config.Limit, &splunkEnv.Config.StartAt, eventsAPI)
	case jwt.Features.Contains(utils.AuditEventsFeatureScope) && EventBuildType == utils.AuditEventsFeatureScope:
		cursorFile := path.Join(splunkEnv.Home, splunkEnv.Config.AuditEventsCursorFile)
		return actions.StartAuditEvents(cursorFile, splunkEnv.Config.Limit, &splunkEnv.Config.StartAt, eventsAPI)
	}
	return fmt.Errorf("JWT features do not match build type")
}
