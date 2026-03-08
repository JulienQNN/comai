package config

type Config struct {
	ProviderName       string `mapstructure:"provider"`
	ModelName          string `mapstructure:"model"`
	Language           string `mapstructure:"language"`
	CustomInstructions string `mapstructure:"custom_instructions"`
}
