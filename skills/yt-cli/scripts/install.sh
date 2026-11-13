#!/usr/bin/env bash
set -euo pipefail

# Downloads the latest release of yt-cli from GitHub,
# verifies its SHA-256 checksum, and extracts it to the current directory.

# Detect OS
OS_UNAME="$(uname -s)"
case "$OS_UNAME" in
  Linux) OS="Linux" ;;
  Darwin) OS="Darwin" ;;
  *) echo "Unsupported OS: $OS_UNAME"; exit 1 ;;
esac

# Detect Architecture
ARCH_UNAME="$(uname -m)"
case "$ARCH_UNAME" in
  x86_64|amd64) ARCH="x86_64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported Architecture: $ARCH_UNAME"; exit 1 ;;
esac

REPO="ghchinoy/yt-cli"

# Resolve latest release version tag
LATEST_URL=$(curl -sLI -o /dev/null -w '%{url_effective}' "https://github.com/${REPO}/releases/latest")
TAG="${LATEST_URL##*/}"
VERSION="${TAG#v}"

FILENAME="yt-cli_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${TAG}/${FILENAME}"
CHECKSUMS_FILENAME="yt-cli_${VERSION}_checksums.txt"
CHECKSUMS_URL="https://github.com/${REPO}/releases/download/${TAG}/${CHECKSUMS_FILENAME}"

TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Downloading yt-cli from: ${URL}"
curl -sL "$URL" -o "${TMP_DIR}/${FILENAME}"

echo "Downloading checksums..."
curl -sL "$CHECKSUMS_URL" -o "${TMP_DIR}/checksums.txt"

echo "Verifying checksum..."
cd "$TMP_DIR"
if ! grep "${FILENAME}" checksums.txt > expected_checksum.txt; then
  echo "Error: Could not find checksum for ${FILENAME} in checksums.txt"
  exit 1
fi

if command -v sha256sum >/dev/null 2>&1; then
  sha256sum -c expected_checksum.txt
elif command -v shasum >/dev/null 2>&1; then
  shasum -a 256 -c expected_checksum.txt
else
  echo "Warning: No sha256sum or shasum command found. Skipping checksum validation."
fi
cd - > /dev/null

echo "Extracting..."
tar -xzf "${TMP_DIR}/${FILENAME}" yt-cli

chmod +x yt-cli
echo "Success! yt-cli is now installed in the current directory."
