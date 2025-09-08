# Release Notes

## v1.0.0 - Initial Release

### Features
- **Interactive TUI Mode**: Browse and execute HTTP requests with an intuitive terminal interface
- **.http File Support**: Full support for standard HTTP file format
- **Environment Variables**: Dynamic variable substitution in requests
- **Cross-Platform**: Binaries available for macOS (ARM64) and Linux (AMD64)
- **Request Management**: 
  - Execute specific requests by index
  - View request history
  - Response viewing with syntax highlighting
- **Easy Installation**: One-line installation script with automatic updates

### Installation
```bash
# Quick install
curl -fsSL https://raw.githubusercontent.com/JckHoe/http-cli/main/install.sh | bash

# Manual download (macOS ARM64)
curl -L https://github.com/JckHoe/http-cli/releases/latest/download/httpx-darwin-arm64 -o httpx
chmod +x httpx
sudo mv httpx /usr/local/bin/
```

### Usage Examples
```bash
# Interactive TUI mode
httpx tui examples/sample.http

# Run specific request
httpx run examples/sample.http --request 1
```

### Development
- Built with Go 1.24
- Includes comprehensive test suite
- GitHub Actions CI/CD pipeline for automated builds and releases

### Binary Artifacts
- `httpx-darwin-arm64` - macOS ARM64 (Apple Silicon)
- `httpx-linux-amd64` - Linux AMD64

### Contributors
- Initial implementation and architecture

---

For bug reports and feature requests, please visit: https://github.com/JckHoe/http-cli/issues