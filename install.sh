#!/bin/bash
set -e

REPO="hulk510/readme-gen"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Detect OS and Arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest version
VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
  echo "Failed to get latest version"
  exit 1
fi

echo "Installing readme-gen ${VERSION}..."

# Download
URL="https://github.com/${REPO}/releases/download/${VERSION}/readme-gen_${OS}_${ARCH}.tar.gz"
TMP_DIR=$(mktemp -d)
curl -sL "$URL" -o "${TMP_DIR}/readme-gen.tar.gz"

# Extract and install
tar -xzf "${TMP_DIR}/readme-gen.tar.gz" -C "$TMP_DIR"

if [ -w "$INSTALL_DIR" ]; then
  mv "${TMP_DIR}/readme-gen" "$INSTALL_DIR/"
else
  sudo mv "${TMP_DIR}/readme-gen" "$INSTALL_DIR/"
fi

rm -rf "$TMP_DIR"

echo "✨ readme-gen installed to ${INSTALL_DIR}/readme-gen"
echo "Run 'readme-gen --help' to get started"
