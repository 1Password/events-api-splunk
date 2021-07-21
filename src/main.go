package main

import (
	"bufio"
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
	"go.1password.io/eventsapi-splunk/splunk"
	"go.1password.io/eventsapi-splunk/utils"
)

const ItemUsageFeatureScope = "itemusages"
const SignInAttemptsFeatureScope = "signinattempts"

var EventBuildType string // Injected at build time so we can make multiple apps
type EnvConfig struct {
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

	configPath := splunkHome + "/etc/apps/onepassword_events_api/local/events_reporting.conf"
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

	return &config, nil
}

func main() {
	log.Println("Booting...")
	if EventBuildType == "" {
		err := fmt.Errorf("missing EventBuildType flag")
		panic(err)
	}

	env, err := loadConfig()
	if err != nil {
		err := fmt.Errorf("could not start: %w", err)
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)
	splunkSession, _, err := reader.ReadLine()
	if err != nil {
		err := fmt.Errorf("could not read session: %w", err)
		panic(err)
	}

	splunkAPI := splunk.NewSplunkAPI(string(splunkSession))
	eventsToken, err := actions.GetEventsToken(context.TODO(), splunkAPI)
	if err != nil {
		err := fmt.Errorf("could not get token: %w", err)
		panic(err)
	}
	
	jwt, err := utils.ParseJWTClaims(eventsToken)
	if err != nil {
		err := fmt.Errorf("could not parse jwt: %w", err)
		panic(err)
	}

	url, err := jwt.GetEventsURL()
	if err != nil {
		err := fmt.Errorf("could not get url from token: %w", err)
		panic(err)
	}

	eventsAPI := events.NewEventsAPI(eventsToken, url)
	eventsRes, err := eventsAPI.Introspect(context.TODO())
	if err != nil {
		err := fmt.Errorf("introspect request failed: %w", err)
		panic(err)
	}

	if utils.ContainsString(SignInAttemptsFeatureScope, eventsRes.Features) && EventBuildType == SignInAttemptsFeatureScope {
		actions.StartSignIns(env.SignInCursorFile, env.Limit, &env.StartAt, eventsAPI)
	} else if utils.ContainsString(ItemUsageFeatureScope, eventsRes.Features) && EventBuildType == ItemUsageFeatureScope {
		actions.StartItemUsages(env.ItemUsageCursorFile, env.Limit, &env.StartAt, eventsAPI)
	}
}
