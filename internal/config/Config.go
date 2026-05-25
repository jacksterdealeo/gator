package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Reads the user's config file found in their home directory.
func Read() (*Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err = json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	if err := c.write(); err != nil {
		return err
	}
	return nil
}

func (c *Config) write() error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0666); err != nil {
		return err
	}
	return nil
}

// Gets the user's config file path.
func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(home, configFileName)
	return path, nil
}
