package config

type Config struct {
	ProviderName       string `mapstructure:"provider"`
	ModelName          string `mapstructure:"model"`
	Language           string `mapstructure:"language"`
	MaxLength          string `mapstructure:"max_length"`
	CustomInstructions string `mapstructure:"custom_instructions"`
}
