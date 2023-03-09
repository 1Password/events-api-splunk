package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/dimchansky/utfbom"
)

type Config struct {
	Url                   string
	AuthToken             string
	StartAt               time.Time
	Limit                 int
	SignInCursorFile      string
	ItemUsageCursorFile   string
	AuditEventsCursorFile string
}

type SplunkEnv struct {
	Home       string
	ConfigPath string
	Config     Config
}

// Gets configuration and normalizes values to EnvConfig.
// Note that the toml parsing library does not support BOM characters.
// LoadConfig must trim a BOM prefix before passing the config bytes to the parser.
func NewSplunkEnv(splunkHome string) (*SplunkEnv, error) {
	log.Println("New Config")

	sc := SplunkEnv{
		Home:       splunkHome,
		ConfigPath: path.Join(splunkHome, "/etc/apps/onepassword_events_api/local/events_reporting.conf"),
	}

	configFile, err := os.Open(sc.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	defer configFile.Close()

	br := utfbom.SkipOnly(configFile)

	if _, err := toml.DecodeReader(br, &sc); err != nil {
		return nil, fmt.Errorf("could not decode toml config file: %w", err)
	}

	return &sc, nil
}

func (e *SplunkEnv) UpdateConfig(newConfig Config) error {
	configTemp := e.ConfigPath + "_" + strconv.Itoa(os.Getpid()) // create temporary file
	configFile, err := os.Create(configTemp)
	if err != nil {
		return fmt.Errorf("could not create config file: %w", err)
	}
	defer configFile.Close()

	type LocalConfig struct {
		Config Config
	}
	lc := LocalConfig{
		Config: newConfig,
	}
	encoder := toml.NewEncoder(configFile)
	err = encoder.Encode(lc)
	if err != nil {
		return fmt.Errorf("could not write toml to file: %w", err)
	}

	err = os.Rename(configTemp, e.ConfigPath)
	if err != nil {
		return fmt.Errorf("could not overwrite previous config: %w", err)
	}
	return nil
}
