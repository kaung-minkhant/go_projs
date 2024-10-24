package config

import (
	"encoding/json"
	"os"
)

var (
	Token    string
	Prefix   string
	AppId    string
	ServerId string
)

type config struct {
	BotToken      string `json:"BotToken"`
	BotPrefix     string `json:"BotPrefix"`
	ApplicationID string `json:"ApplicationID"`
	ServerId      string `json:"ServerId"`
}

func LoadConfig(filename string) error {
	body, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var configData config

	if err := json.Unmarshal(body, &configData); err != nil {
		return err
	}

	Token = configData.BotToken
	Prefix = configData.BotPrefix
	AppId = configData.ApplicationID
	ServerId = configData.ServerId

	return nil
}
