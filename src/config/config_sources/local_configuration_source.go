package config_sources

import (
	"chat/config"
	"os"
	"path/filepath"
)

type LocalConfigurationSource struct{}

func (LocalConfigurationSource) Load() (config.Configuration, error) {
	var configs []config.Configuration

	// Start from the current working directory and work upwards.
	dir, err := os.Getwd()
	if err != nil {
		return config.Configuration{}, err
	}

	for {
		// Construct the path to the .chatconfig file in the current directory
		configPath := filepath.Join(dir, ".chatconfig")

		// Check if the .chatconfig file exists and is not a directory
		info, err := os.Stat(configPath)

		// If .chatconfig is found, load and collect its configuration
		if err == nil && !info.IsDir() {
			// Parse the .chatconfig file into a Config object
			cfg, err := config.ConfigurationFromFile(configPath)
			if err != nil {
				return config.Configuration{}, err
			}
			configs = append(configs, cfg)
			// On unexpected errors (other than file not found), abort
		} else if err != nil && !os.IsNotExist(err) {
			return config.Configuration{}, err
		}
		// Determine the parent directory of the current directory
		parent := filepath.Dir(dir)

		// If we have reached the root (no further parents), stop the loop
		if parent == dir {
			break
		}

		// Move up to the parent directory for the next iteration
		dir = parent
	}

	// Merge all the loaded configurations.
	localConfiguration := config.MergeConfigurations(configs)

	// Return all loaded configurations (nearest-first)
	return localConfiguration, nil
}
