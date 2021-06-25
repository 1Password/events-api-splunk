package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/dimchansky/utfbom"
	"go.1password.io/eventsapi-splunk/actions"
	events "go.1password.io/eventsapi-splunk/api"
	"go.1password.io/eventsapi-splunk/utils"
)

const ItemUsageFeatureScope = "itemusages"
const SignInAttemptsFeatureScope = "signinattempts"

var EventBuildType string // Injected at build time so we can make multiple apps
type EnvConfig struct {
	Url                 string
	AuthToken           string
	StartAt             time.Time
	Limit               int
	SignInCursorFile    string
	ItemUsageCursorFile string
}

// Gets environment variables and normalizes values to EnvConfig.
// Note that the toml parsing library does not support BOM characters.
// LoadConfig must trim a BOM prefix before passing the config bytes to the parser.
func loadConfig() (*EnvConfig, error) {
	log.Println("Loading config")
	splunkHome := os.Getenv("SPLUNK_HOME")
	if splunkHome == "" {
		return nil, fmt.Errorf("SPLUNK_HOME environment variable must be set")
	}

	configPath := splunkHome + "/etc/apps/op_events_reporting/local/events_reporting.conf"
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	defer configFile.Close()

	br := utfbom.SkipOnly(configFile)

	type LocalConfig struct {
		Config EnvConfig
	}
	var lc LocalConfig
	if _, err := toml.DecodeReader(br, &lc); err != nil {
		return nil, fmt.Errorf("could not decode toml config file: %w", err)
	}
	config := lc.Config
	config.ItemUsageCursorFile = path.Join(splunkHome, config.ItemUsageCursorFile)
	config.SignInCursorFile = path.Join(splunkHome, config.SignInCursorFile)

	jwt, err := utils.ParseJWTClaims(config.AuthToken)
	if err != nil {
		return nil, err
	}

	url := jwt.GetEventsURL()
	// The config url will be used if the token was generated before
	// this update and does not contain a url
	if url != "" {
		config.Url = url
	}

	return &config, nil
}

func main() {
	log.Println("Booting...")

	env, err := loadConfig()
	if err != nil {
		err := fmt.Errorf("could not start: %w", err)
		panic(err)
	}

	if EventBuildType == "" {
		err := fmt.Errorf("missing EventBuildType flag")
		panic(err)
	}

	eventsAPI := events.NewEventsAPI(env.AuthToken, env.Url)
	res, err := eventsAPI.Introspect(context.TODO())
	if err != nil {
		err := fmt.Errorf("introspect request failed: %w", err)
		panic(err)
	}

	if utils.ContainsString(SignInAttemptsFeatureScope, res.Features) && EventBuildType == SignInAttemptsFeatureScope {
		actions.StartSignIns(env.SignInCursorFile, env.Limit, &env.StartAt, eventsAPI)
	} else if utils.ContainsString(ItemUsageFeatureScope, res.Features) && EventBuildType == ItemUsageFeatureScope {
		actions.StartItemUsages(env.ItemUsageCursorFile, env.Limit, &env.StartAt, eventsAPI)
	}
}
