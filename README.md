# httpx

[![Build and Test](https://github.com/JckHoe/http-cli/actions/workflows/build.yml/badge.svg)](https://github.com/JckHoe/http-cli/actions/workflows/build.yml)
[![Release](https://github.com/JckHoe/http-cli/actions/workflows/release.yml/badge.svg)](https://github.com/JckHoe/http-cli/actions/workflows/release.yml)

A powerful HTTP file runner CLI tool with an interactive TUI interface.

## Installation

### Quick Install (Recommended)

Install or update to the latest version:

```bash
curl -fsSL https://raw.githubusercontent.com/JckHoe/http-cli/main/install.sh | bash
```

### Installation Options

```bash
# Force update without prompts
curl -fsSL https://raw.githubusercontent.com/JckHoe/http-cli/main/install.sh | bash -s -- --force

# Install specific version
curl -fsSL https://raw.githubusercontent.com/JckHoe/http-cli/main/install.sh | bash -s -- --version v1.0.0

# Install to custom directory
curl -fsSL https://raw.githubusercontent.com/JckHoe/http-cli/main/install.sh | bash -s -- --dir ~/bin
```

### Manual Installation

Download the appropriate binary for your system from the [releases page](https://github.com/JckHoe/http-cli/releases):

- **macOS (ARM64)**: `httpx-darwin-arm64`
- **Linux (AMD64)**: `httpx-linux-amd64`

```bash
# Example for macOS ARM64
curl -L https://github.com/JckHoe/http-cli/releases/latest/download/httpx-darwin-arm64 -o httpx
chmod +x httpx
sudo mv httpx /usr/local/bin/
```

### Build from Source

```bash
git clone https://github.com/JckHoe/http-cli.git
cd http-cli
make build
sudo make install
```

## Usage

### Interactive TUI Mode

```bash
httpx tui examples/sample.http
```

### Run Specific Request

```bash
httpx run examples/sample.http --request 1
```

### Update to Latest Version

The installer script automatically checks for updates:

```bash
curl -fsSL https://raw.githubusercontent.com/JckHoe/http-cli/main/install.sh | bash
```

If an update is available, you'll be prompted to install it. Use `--force` to skip the prompt.

## Features

- Interactive TUI for browsing and executing HTTP requests
- Support for `.http` file format
- Environment variable support
- Request history and response viewing
- Cross-platform support (macOS ARM64, Linux AMD64)
- Automatic version updates

## Development

### Prerequisites

- Go 1.24 or higher
- Make

### Building

```bash
make build
```

### Running Tests

```bash
make test
```

### Local Development

```bash
make run ARGS="your arguments here"
```

## License

[Add your license here]

## Contributing

[Add contributing guidelines here]
