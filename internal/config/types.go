package config

type Config struct {
	ModelName       string `mapstructure:"model"`
	Language        string `mapstructure:"language"`
	CommitMaxLength int    `mapstructure:"commit_max_length"`
	PromptAddition  string `mapstructure:"prompt_addition"`
}
