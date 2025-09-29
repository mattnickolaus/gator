package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home + "/" + configFileName, nil
}

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func write(cfg Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// 0644 is the file permissions
	err = os.WriteFile(configFilePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	err := write(*cfg)
	if err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
