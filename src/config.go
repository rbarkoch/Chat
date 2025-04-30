package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// Config holds key/value pairs loaded from a file.
type Config struct {
	values map[string]string
}

// Reads a configuration file.
func LoadConfiguration(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := &Config{values: make(map[string]string)}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		// Strip comments starting with '#' or ';'
		if idx := strings.IndexAny(line, "#;"); idx != -1 {
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		cfg.values[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Get returns the value associated with the given key, or an empty string if not present.
func (c *Config) Get(key string) string {
	return c.values[key]
}

// Merges configurations. First has the highest priority.
func MergeConfigurations(cfgs []*Config) Config {
	merged := Config{
		values: make(map[string]string),
	}

	for _, cfg := range cfgs {
		for k, v := range cfg.values {
			merged.values[k] = v
		}
	}
	return merged
}

func LoadConfigurationsFromWorkingDirectory() ([]*Config, error) {
	var configs []*Config
	// Start from the current working directory and work upwards.
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for {
		// Construct the path to the .chatconfig file in the current directory
		configPath := filepath.Join(dir, ".chatconfig")

		// Check if the .chatconfig file exists and is not a directory
		info, err := os.Stat(configPath)

		// If .chatconfig is found, load and collect its configuration
		if err == nil && !info.IsDir() {
			// Parse the .chatconfig file into a Config object
			cfg, err := LoadConfiguration(configPath)
			if err != nil {
				return nil, err
			}
			configs = append(configs, cfg)
			// On unexpected errors (other than file not found), abort
		} else if err != nil && !os.IsNotExist(err) {
			return nil, err
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

	// Return all loaded configurations (nearest-first)
	return configs, nil
}
