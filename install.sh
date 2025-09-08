#!/bin/bash

set -e

# Configuration
REPO="JckHoe/http-cli"
BINARY_NAME="httpx"
INSTALL_DIR="${HTTPX_INSTALL_DIR:-/usr/local/bin}"
GITHUB_API="https://api.github.com"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
print_error() {
    echo -e "${RED}Error: $1${NC}" >&2
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_info() {
    echo -e "${BLUE}→ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

detect_os() {
    OS="$(uname -s)"
    case "${OS}" in
        Linux*)     OS_TYPE="linux";;
        Darwin*)    OS_TYPE="darwin";;
        *)          print_error "Unsupported OS: ${OS}"; exit 1;;
    esac
}

detect_arch() {
    ARCH="$(uname -m)"
    case "${ARCH}" in
        x86_64|amd64)   ARCH_TYPE="amd64";;
        arm64|aarch64)  ARCH_TYPE="arm64";;
        *)              print_error "Unsupported architecture: ${ARCH}"; exit 1;;
    esac
}

get_current_version() {
    if command -v ${BINARY_NAME} &> /dev/null; then
        # Try to get version from binary
        local version=$(${BINARY_NAME} --version 2>/dev/null | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' | head -1)
        if [ -z "$version" ]; then
            # If --version doesn't work, check if binary exists
            if [ -f "${INSTALL_DIR}/${BINARY_NAME}" ]; then
                echo "unknown"
            else
                echo "none"
            fi
        else
            echo "$version"
        fi
    else
        echo "none"
    fi
}

get_latest_release() {
    local api_url="${GITHUB_API}/repos/${REPO}/releases/latest"
    
    if command -v curl &> /dev/null; then
        curl -s "$api_url" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command -v wget &> /dev/null; then
        wget -qO- "$api_url" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    else
        print_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi
}

download_file() {
    local url=$1
    local output=$2
    
    if command -v curl &> /dev/null; then
        curl -sL "$url" -o "$output"
    elif command -v wget &> /dev/null; then
        wget -q "$url" -O "$output"
    else
        print_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi
}

verify_checksum() {
    local file=$1
    local checksum_url=$2
    local temp_checksum="/tmp/${BINARY_NAME}.sha256"
    
    print_info "Verifying checksum..."
    download_file "$checksum_url" "$temp_checksum"
    
    if command -v sha256sum &> /dev/null; then
        # Linux
        local expected=$(cat "$temp_checksum" | awk '{print $1}')
        local actual=$(sha256sum "$file" | awk '{print $1}')
    elif command -v shasum &> /dev/null; then
        # macOS
        local expected=$(cat "$temp_checksum" | awk '{print $1}')
        local actual=$(shasum -a 256 "$file" | awk '{print $1}')
    else
        print_warning "Cannot verify checksum (no sha256sum or shasum found)"
        rm -f "$temp_checksum"
        return 0
    fi
    
    if [ "$expected" = "$actual" ]; then
        print_success "Checksum verified"
        rm -f "$temp_checksum"
        return 0
    else
        print_error "Checksum verification failed"
        rm -f "$temp_checksum"
        return 1
    fi
}

install_binary() {
    local version=$1
    local force=${2:-false}
    
    # Determine binary name based on OS and architecture
    local binary_suffix="${OS_TYPE}-${ARCH_TYPE}"
    
    # Special case: Linux only supports amd64 for now
    if [ "$OS_TYPE" = "linux" ] && [ "$ARCH_TYPE" != "amd64" ]; then
        print_error "Linux ${ARCH_TYPE} is not supported yet"
        exit 1
    fi
    
    # Special case: macOS only supports arm64 for now
    if [ "$OS_TYPE" = "darwin" ] && [ "$ARCH_TYPE" != "arm64" ]; then
        print_error "macOS ${ARCH_TYPE} is not supported yet"
        exit 1
    fi
    
    local download_url="https://github.com/${REPO}/releases/download/${version}/${BINARY_NAME}-${binary_suffix}"
    local checksum_url="${download_url}.sha256"
    local temp_file="/tmp/${BINARY_NAME}-${version}"
    
    print_info "Downloading ${BINARY_NAME} ${version} for ${OS_TYPE}/${ARCH_TYPE}..."
    download_file "$download_url" "$temp_file"
    
    if [ ! -f "$temp_file" ]; then
        print_error "Failed to download binary"
        exit 1
    fi
    
    # Verify checksum
    if ! verify_checksum "$temp_file" "$checksum_url"; then
        rm -f "$temp_file"
        exit 1
    fi
    
    # Check if we need sudo for installation
    if [ -w "$INSTALL_DIR" ]; then
        print_info "Installing to ${INSTALL_DIR}/${BINARY_NAME}..."
        mv "$temp_file" "${INSTALL_DIR}/${BINARY_NAME}"
        chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    else
        print_info "Installing to ${INSTALL_DIR}/${BINARY_NAME} (requires sudo)..."
        sudo mv "$temp_file" "${INSTALL_DIR}/${BINARY_NAME}"
        sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    fi
    
    print_success "${BINARY_NAME} ${version} installed successfully!"
}

prompt_update() {
    local current=$1
    local latest=$2
    
    echo ""
    print_warning "Update available: ${current} → ${latest}"
    echo -n "Would you like to update? (y/n): "
    read -r response
    case "$response" in
        [yY][eE][sS]|[yY])
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

show_usage() {
    cat << EOF
httpx Installer

Usage: $0 [OPTIONS]

OPTIONS:
    --force             Force installation without prompts
    --version VERSION   Install a specific version (e.g., --version v1.0.0)
    --dir DIRECTORY     Install to a custom directory (default: /usr/local/bin)
    --help              Show this help message

EXAMPLES:
    # Install or update to latest version
    $0

    # Force update without prompts
    $0 --force

    # Install specific version
    $0 --version v1.2.3

    # Install to custom directory
    $0 --dir ~/bin

EOF
}

# Main script
main() {
    local force=false
    local specific_version=""
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --force)
                force=true
                shift
                ;;
            --version)
                specific_version="$2"
                shift 2
                ;;
            --dir)
                INSTALL_DIR="$2"
                shift 2
                ;;
            --help)
                show_usage
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
    
    echo ""
    echo "========================================="
    echo "     httpx Installer"
    echo "========================================="
    echo ""
    
    # Detect OS and architecture
    detect_os
    detect_arch
    
    print_info "Detected: ${OS_TYPE}/${ARCH_TYPE}"
    print_info "Install directory: ${INSTALL_DIR}"
    
    # Get version information
    local current_version=$(get_current_version)
    local target_version="${specific_version}"
    
    if [ -z "$target_version" ]; then
        print_info "Checking for latest version..."
        target_version=$(get_latest_release)
        
        if [ -z "$target_version" ]; then
            print_error "Failed to get latest release information"
            exit 1
        fi
    fi
    
    print_info "Latest version: ${target_version}"
    
    # Check if update is needed
    if [ "$current_version" = "none" ]; then
        print_info "No existing installation found"
        install_binary "$target_version" "$force"
    elif [ "$current_version" = "unknown" ]; then
        print_warning "Current version unknown"
        if [ "$force" = true ] || prompt_update "$current_version" "$target_version"; then
            install_binary "$target_version" "$force"
        else
            print_info "Installation cancelled"
            exit 0
        fi
    elif [ "$current_version" = "$target_version" ]; then
        print_success "Already up to date (${current_version})"
        if [ "$force" = true ]; then
            print_info "Force reinstalling..."
            install_binary "$target_version" "$force"
        fi
    else
        print_info "Current version: ${current_version}"
        if [ "$force" = true ] || prompt_update "$current_version" "$target_version"; then
            install_binary "$target_version" "$force"
        else
            print_info "Update cancelled"
            exit 0
        fi
    fi
    
    # Verify installation
    if command -v ${BINARY_NAME} &> /dev/null; then
        print_success "Installation complete!"
        echo ""
        echo "Run '${BINARY_NAME} --help' to get started"
    else
        print_warning "${BINARY_NAME} installed but not in PATH"
        print_info "Add ${INSTALL_DIR} to your PATH or specify the full path to use ${BINARY_NAME}"
    fi
}

# Run main function
main "$@"
