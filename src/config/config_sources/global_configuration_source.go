package config_sources

import (
	"chat/config"
	"fmt"
	"os"
	"path"
)

type GlobalConfigurationSource struct{}

func (GlobalConfigurationSource) Load() (config.Configuration, error) {
	// Find the users home directory.
	userHomePath, err := os.UserHomeDir()
	if err != nil {
		return config.Configuration{}, fmt.Errorf("error finding user home directory\n%w", err)
	}

	// Find a configuration file in the users home directory.
	homeConfigPath := path.Join(userHomePath, ".chatconfig")
	_, err = os.Stat(homeConfigPath)
	if os.IsNotExist(err) {
		// If the user does not have a home directory defined, just return an empty configuration.
		return config.New(), nil
	} else if err != nil {
		return config.Configuration{}, fmt.Errorf("failed while checking for users home directory\n%w", err)
	}

	// Load configuration file from users home directory.
	homeConfig, err := config.ConfigurationFromFile(homeConfigPath)
	if err != nil {
		return config.Configuration{}, fmt.Errorf("error loading configuration file from users home directory\n%w", err)
	}

	return homeConfig, nil
}
