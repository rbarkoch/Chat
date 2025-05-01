package config

type ConfigurationSource interface {
	Load() (Configuration, error)
}
