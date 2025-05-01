package config_sources

import (
	"chat/config"
	"os"
)

type EnvConfigurationSource struct{}

func (EnvConfigurationSource) Load() (config.Configuration, error) {
	configuration := config.New()
	configuration.Add(config.ConfigKeyApiKey, os.Getenv(config.ConfigKeyApiKey))
	configuration.Add(config.ConfigKeySystemPrompt, os.Getenv(config.ConfigKeySystemPrompt))
	configuration.Add(config.ConfigKeyModel, os.Getenv(config.ConfigKeyModel))

	return configuration, nil
}
