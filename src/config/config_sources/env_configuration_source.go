package config_sources

import (
	"chat/config"
	"os"
)

type EnvConfigurationSource struct{}

func (EnvConfigurationSource) Load() (config.Configuration, error) {
	configuration := config.New()
	apiKey := os.Getenv(config.ConfigKeyApiKey)
	if apiKey != "" {
		configuration.Add(config.ConfigKeyApiKey, os.Getenv(config.ConfigKeyApiKey))
	}

	systemPrompt := os.Getenv(config.ConfigKeySystemPrompt)
	if systemPrompt != "" {
		configuration.Add(config.ConfigKeySystemPrompt, os.Getenv(config.ConfigKeySystemPrompt))
	}

	model := os.Getenv(config.ConfigKeyModel)
	if model != "" {
		configuration.Add(config.ConfigKeyModel, os.Getenv(config.ConfigKeyModel))
	}

	return configuration, nil
}
