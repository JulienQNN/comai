# comai

Generate git commit messages with a local AI model.

## Installation

### With a Package Manager (Recommended)

**Homebrew** (macOS / Linux)
```bash
brew tap JulienQNN/tap
brew install comai
```

**Winget** (Windows)
```bash
winget install JulienQNN.comai
```

**Chocolatey** (Windows)
```bash
choco install comai
```

### Manual install

Download the latest binary or package from the [releases page](https://github.com/JulienQNN/comai/releases).

**apt** (Debian / Ubuntu)
```bash
sudo dpkg -i comai_*_linux_amd64.deb
```

**rpm** (Fedora / RHEL)
```bash
sudo rpm -i comai_*_linux_amd64.rpm
```

**pacman** (Arch Linux)
```bash
sudo pacman -U comai_*_linux_amd64.pkg.tar.zst
```

## Usage

```bash
# Configure comai
comai init

# Generate a commit message from your staged changes
git add .
comai generate
```

## Requirements

- [Ollama](https://ollama.com) running locally
- At least one model pulled, e.g. `ollama pull qwen3:8b`
