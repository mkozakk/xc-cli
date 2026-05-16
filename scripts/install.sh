#!/usr/bin/env bash
set -euo pipefail

REPO="mkozakk/xc-cli"
INSTALL_DIR="${XC_INSTALL_DIR:-$HOME/.local/bin}"
BINARY_NAME="xc"

arch="$(uname -m)"
case "$arch" in
    x86_64)  GOARCH="amd64" ;;
    aarch64) GOARCH="arm64" ;;
    *)
        echo "Error: unsupported architecture ($arch)" >&2
        exit 1
        ;;
esac

if [[ -z "${XC_VERSION:-}" ]]; then
    echo "Fetching latest release from GitHub..."
    LATEST_URL="https://api.github.com/repos/${REPO}/releases/latest"
    VERSION="$(curl -fsSL "$LATEST_URL" | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')"
else
    VERSION="${XC_VERSION#v}"
fi

echo "Installing xc v${VERSION}..."

TARBALL="xc_${VERSION}_linux_${GOARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/v${VERSION}/${TARBALL}"
CHECKSUM_URL="https://github.com/${REPO}/releases/download/v${VERSION}/xc_${VERSION}_checksums.txt"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Downloading binary..."
curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/$TARBALL"
curl -fsSL "$CHECKSUM_URL" -o "$TMP_DIR/checksums.txt"

echo "Verifying checksum..."
cd "$TMP_DIR"
grep "$TARBALL" checksums.txt | sha256sum --check --status
echo "✓ Checksum verified"

echo "Installing binary..."
tar xzf "$TARBALL" -C "$TMP_DIR"
mkdir -p "$INSTALL_DIR"
install -m 755 "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
echo "✓ Installed to $INSTALL_DIR/$BINARY_NAME"

add_to_path_if_needed() {
    local rc_file="$1"
    local path_line='export PATH="$HOME/.local/bin:$PATH"'

    if [[ ! -f "$rc_file" ]]; then
        return
    fi

    if grep -qF '.local/bin' "$rc_file"; then
        return
    fi

    printf '\nexport PATH="$HOME/.local/bin:$PATH"\n' >> "$rc_file"
}

setup_shell_hook() {
    local rc_file="$1"
    local shell_name="$2"
    local hook_line='eval "$(xc init '"$shell_name"')"'

    if [[ ! -f "$rc_file" ]]; then
        touch "$rc_file"
    fi

    if grep -qF 'xc init' "$rc_file"; then
        return 1
    fi

    printf '\neval "$(xc init %s)"\n' "$shell_name" >> "$rc_file"
    return 0
}

CURRENT_SHELL="$(basename "$SHELL")"

case "$CURRENT_SHELL" in
    bash)
        echo "Setting up bash hook..."
        add_to_path_if_needed "$HOME/.bashrc"
        if setup_shell_hook "$HOME/.bashrc" "bash"; then
            echo "✓ Added shell hook to ~/.bashrc"
        else
            echo "✓ Shell hook already present in ~/.bashrc"
        fi
        RC_FILE="$HOME/.bashrc"
        ;;
    zsh)
        echo "Setting up zsh hook..."
        add_to_path_if_needed "$HOME/.zshrc"
        if setup_shell_hook "$HOME/.zshrc" "zsh"; then
            echo "✓ Added shell hook to ~/.zshrc"
        else
            echo "✓ Shell hook already present in ~/.zshrc"
        fi
        RC_FILE="$HOME/.zshrc"
        ;;
    *)
        echo "Warning: shell '$CURRENT_SHELL' not automatically configured"
        echo "Add these lines to your shell rc file manually:"
        printf '\nexport PATH="$HOME/.local/bin:$PATH"\neval "$(%s init bash)"\n' "$BINARY_NAME"
        RC_FILE=""
        ;;
esac

echo ""
echo "Installation complete! 🎉"

if [[ -n "$RC_FILE" ]]; then
    echo ""
    echo "To start using xc right now, run:"
    printf '  source %s\n' "$RC_FILE"
    echo ""
    echo "Or just open a new terminal."
fi
