#!/bin/sh
set -e

REPO="tferazzi/ktidy"
BIN="ktidy"
INSTALL_DIR="/usr/local/bin"

# Detect OS
case "$(uname -s)" in
  Darwin) OS="darwin" ;;
  Linux)  OS="linux" ;;
  *)      echo "Unsupported OS: $(uname -s)" >&2; exit 1 ;;
esac

# Detect arch
case "$(uname -m)" in
  x86_64)        ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)             echo "Unsupported arch: $(uname -m)" >&2; exit 1 ;;
esac

# Resolve latest version
VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep '"tag_name"' | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')
[ -n "$VERSION" ] || { echo "Failed to resolve latest version" >&2; exit 1; }

# Goreleaser strips the leading 'v' from the archive name
VER="${VERSION#v}"
ARCHIVE="${BIN}_${VER}_${OS}_${ARCH}.tar.gz"
BASE_URL="https://github.com/${REPO}/releases/download/${VERSION}"
CHECKSUMS="${BIN}_${VER}_checksums.txt"

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

echo "Downloading ${BIN} ${VERSION} (${OS}/${ARCH})..."
curl -fsSL "${BASE_URL}/${ARCHIVE}"   -o "${TMP}/${ARCHIVE}"
curl -fsSL "${BASE_URL}/${CHECKSUMS}" -o "${TMP}/${CHECKSUMS}"

# Verify checksum
cd "$TMP"
if command -v sha256sum >/dev/null 2>&1; then
  grep "${ARCHIVE}" "${CHECKSUMS}" | sha256sum -c -
elif command -v shasum >/dev/null 2>&1; then
  grep "${ARCHIVE}" "${CHECKSUMS}" | shasum -a 256 -c -
else
  echo "Warning: no sha256 tool found, skipping checksum verification" >&2
fi

tar -xzf "${ARCHIVE}" "${BIN}"

if [ -w "$INSTALL_DIR" ]; then
  mv "${BIN}" "${INSTALL_DIR}/"
else
  INSTALL_DIR="${HOME}/.local/bin"
  mkdir -p "$INSTALL_DIR"
  mv "${BIN}" "${INSTALL_DIR}/"
  echo "Note: installed to ${INSTALL_DIR} — ensure it is on your PATH"
fi

echo "${BIN} ${VERSION} installed → ${INSTALL_DIR}/${BIN}"
