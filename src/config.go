package main

import (
	"encoding/json"
	"errors"
	"os"
)

const (
	configFilePath = "termpad-config.json"
)

// TODO: Application version specific version migration
// Structure representig the configuration properties insinde the termpad-config.json file
type Config struct {
	HistoryConfiguration  HistoryConfig  `json:"history-configuration"`
	KeybindsConfiguration KeybindsConfig `json:"keybinds-configuration"`
	CursorConfiguration   CursorConfig   `json:"cursor-configuration"`
	TextConfiguration     TextConfig     `json:"text-configuration"`
}

// Config structure initialization function. The function is retriving the config file or creating a default one if not present
func (config *Config) Init() error {
	var configFileExists bool
	if _, err := os.Stat(configFilePath); err == nil {
		configFileExists = true
	} else if errors.Is(err, os.ErrNotExist) {
		configFileExists = false
	} else {
		return errors.New("config: can not determine if the config file is accesable")
	}

	// NOTE: Retrieve config from existing file
	if configFileExists {
		configFileData, err := os.ReadFile(configFilePath)
		if err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(string(configFileData)), &config); err != nil {
			return err
		}

		return nil
	}

	// NOTE: Config file not found, creating config file with defaut values
	config.HistoryConfiguration = CreateDefaultHistoryConfig()
	config.KeybindsConfiguration = CreateDefaultKeybindsConfig()
	config.CursorConfiguration = CreateDefaultCursorConfig()
	config.TextConfiguration = CreateDefaultTextConfig()

	jsonConfig, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return err
	}

	configFile, err := os.Create(configFilePath)
	if err != nil {
		return err
	}

	if _, err := configFile.Write(jsonConfig); err != nil {
		if fileErr := configFile.Close(); fileErr != nil {
			return fileErr
		}

		return err
	}

	if err := configFile.Close(); err != nil {
		return err
	}

	return nil
}
