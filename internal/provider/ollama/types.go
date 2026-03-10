package ollama

import "github.com/ollama/ollama/api"

type Client struct {
	model  string
	client *api.Client
}
