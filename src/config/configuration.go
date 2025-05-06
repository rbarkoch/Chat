package config

import (
	"bufio"
	"os"
	"strings"
)

const (
	ConfigKeyModel        = "CHAT_MODEL"
	ConfigKeyApiKey       = "CHAT_API_KEY"
	ConfigKeySystemPrompt = "CHAT_SYSTEM_PROMPT"
	ConfigWebSearch       = "CHAT_WEB_SEARCH"
)

// Configuration holds key/value pairs loaded from a file.
type Configuration struct {
	values map[string]string
}

// Constructs a new empty configuration.
func New() Configuration {
	return Configuration{
		values: make(map[string]string),
	}
}

// Adds a key to the configuration.
func (cfg *Configuration) Add(key string, value string) {
	cfg.values[key] = value
}

// Get returns the value associated with the given key, or an empty string if not present.
func (cfg *Configuration) Get(key string) string {
	return cfg.values[key]
}

// Reads a configuration file. Configuration files should be simple KEY=VALUE files.
func ConfigurationFromFile(path string) (Configuration, error) {
	file, err := os.Open(path)
	if err != nil {
		return Configuration{}, err
	}
	defer file.Close()

	cfg := New()
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
		return Configuration{}, err
	}
	return cfg, nil
}

// Merges configurations. First has the highest priority.
func MergeConfigurations(cfgs []Configuration) Configuration {
	merged := New()

	for _, cfg := range cfgs {
		for k, v := range cfg.values {

			_, exists := merged.values[k]
			if exists {
				// If we already have a key, skip because we have higher priority.
				continue
			}
			merged.values[k] = v
		}
	}
	return merged
}
