# comai
![](https://github.com/JulienQNN/comai/demo.gif)

Generate AI-powered git commit messages automatically from staged changes. Supports **Ollama** (local) and **GitHub Copilot**.

## Installation

**Homebrew** (macOS / Linux)
```bash
brew tap JulienQNN/tap
brew install comai
```

**From Source**
```bash
git clone https://github.com/JulienQNN/comai.git
cd comai && go build -o comai main.go && sudo mv comai /usr/local/bin/
```

## Quick Start

```bash
# 1. Configure comai (choose provider: Ollama or Copilot)
comai init

# 2. Stage your changes
git add .

# 3. Generate commit message
comai generate

# 4. Review, edit if needed, confirm
```

## Commands

### `comai init`
Initialize with interactive setup. Choose provider, model, language, max length.

**Flags:**
- `--global, -g`: Save to `~/.comai.yaml` instead of `.comai.yaml`

### `comai generate`
Generate commit from staged changes.

**Flags:**
- `--global, -g`: Use `~/.comai.yaml`
- `--verbose, -v`: Show config details
- `--date "yesterday"`: Set custom commit date (natural language or ISO format)
- `--date-interactive, -D`: Choose date after generation

## Providers

### Ollama (Local)
- Runs on your machine
- Requires [Ollama](https://ollama.com) installed and running
- Models: `qwen3.5:9b` (recommended), `qwen3.5:4b`, `mistral:7b`, etc.
- Setup: `ollama serve` in one terminal, then `ollama pull qwen3.5:9b`

### Copilot
- Uses GitHub Copilot API
- Requires [GitHub CLI](https://github.com/cli/cli) + Copilot subscription

## Configuration
Use `comai init` or edit `.comai.yaml` (local) or `~/.comai.yaml` (global):

```yaml
provider: ollama              # or "copilot"
model: qwen3.5:9b            # Provider-specific
language: en
commit_max_length: 50         # Characters
custom_instructions: ""       # Optional: override system prompt
```

## Features

- **Intelligent filtering**: Ignores lock files, binaries, auto-generated content
- **Custom prompts**: Override system message with your own instructions
- **Date handling**: Use natural language ("yesterday", "last week") or ISO format
- **Streaming**: Real-time token display while generating
- **Editing**: Review and edit generated message before committing

## Requirements

- **Git**: With configured `user.name` and `user.email`
- **Copilot path**: GitHub Copilot CLI + active subscription

## License

See [LICENSE](LICENSE) | [Contributing](CONTRIBUTING.md)
