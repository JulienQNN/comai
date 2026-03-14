package copilot

import copilot "github.com/github/copilot-sdk/go"

const defaultModel = "claude-sonnet-4.5"

type Client struct {
	model  string
	client *copilot.Client
}
