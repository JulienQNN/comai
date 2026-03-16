package config

type Config struct {
	ProviderName       string `mapstructure:"provider"`
	ModelName          string `mapstructure:"model"`
	CommitMaxLength    string `mapstructure:"commit_max_length"`
	Language           string `mapstructure:"language"`
	CustomInstructions string `mapstructure:"custom_instructions"`
}
