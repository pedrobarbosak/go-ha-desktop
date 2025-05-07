package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pedrobarbosak/go-errors"
)

type Config struct {
	Version       int           `json:"version"`
	URls          []string      `json:"urls"`
	AccessToken   string        `json:"accessToken"`
	PinnedDevices []string      `json:"pinnedDevices"`
	ScanInterval  time.Duration `json:"scanInterval"`
}

var mt sync.Mutex
var configDir string

const filePermissions = 0744

func Default() *Config {
	return &Config{
		Version:      1,
		ScanInterval: 60,
	}
}

func setup() error {
	mt.Lock()
	defer mt.Unlock()

	if configDir != "" {
		return nil
	}

	userDir, err := os.UserConfigDir()
	if err != nil {
		return errors.New("failed to get config dir:", err)
	}

	appDir := filepath.Join(userDir, "ha-desktop")
	if _, err = os.Stat(appDir); os.IsNotExist(err) {
		if err = os.MkdirAll(appDir, filePermissions); err != nil {
			return errors.New("failed to create config dir:", err)
		}
	}

	configDir = filepath.Join(appDir, "config.json")

	return nil
}

func Load() (*Config, error) {
	if err := setup(); err != nil {
		return nil, err
	}

	_, err := os.Stat(configDir)
	if err != nil && os.IsNotExist(err) {
		cfg := Default()
		if err = Save(cfg); err != nil {
			return nil, err
		}

		return cfg, nil
	}

	data, err := os.ReadFile(configDir)
	if err != nil {
		return nil, errors.New("failed to read config from file:", err)
	}

	cfg := Default()
	if err = json.Unmarshal(data, cfg); err != nil {
		return nil, errors.New("failed to unmarshal config:", err)
	}

	return cfg, nil
}

func Save(config *Config) error {
	if err := setup(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return errors.New("failed to save config:", err)
	}

	if err = os.WriteFile(configDir, data, filePermissions); err != nil {
		return errors.New("failed to save config:", err)
	}

	return nil
}
