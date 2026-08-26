package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, configFileName), nil
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	jsonData, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config

	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing report file: %w", err)
	}
	return nil
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName

	err := write(*c)
	if err != nil {
		return err
	}

	return nil
}
