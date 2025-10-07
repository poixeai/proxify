#!/bin/bash
set -e

# === basic config ===
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RELEASES_DIR="${SCRIPT_DIR}/releases"
CURRENT_DIR="${SCRIPT_DIR}/current"

# === custom config ===
REPO="poixeai/proxify"            
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi
GITHUB_TOKEN="${GITHUB_TOKEN:?Environment variable GITHUB_TOKEN not set}"

APP_NAME="$1"
VERSION="$2"

if [ -z "$APP_NAME" ] || [ -z "$VERSION" ]; then
  echo "[ERROR] Usage: $0 <app_name> <version>"
  echo "Example: $0 ci-react-web-prod v1.2.3"
  exit 1
fi

echo "===== Step 1: Preparing directories ====="
mkdir -p "$RELEASES_DIR" "$CURRENT_DIR"
ZIP_FILE="${RELEASES_DIR}/${APP_NAME}-${VERSION}.zip"

echo "===== Step 2: Downloading release asset ====="

if [ ! -f "$ZIP_FILE" ]; then
  RELEASE_JSON=$(curl -s -H "Authorization: token ${GITHUB_TOKEN}" \
    "https://api.github.com/repos/${REPO}/releases/tags/${VERSION}")

  ASSET_URL=$(echo "$RELEASE_JSON" | jq -r ".assets[] | select(.name==\"${APP_NAME}.zip\") | .url")
  if [ -z "$ASSET_URL" ] || [ "$ASSET_URL" == "null" ]; then
    echo "[ERROR] Asset '${APP_NAME}.zip' not found in release ${VERSION}."
    exit 1
  fi

  curl -L -H "Authorization: token ${GITHUB_TOKEN}" \
       -H "Accept: application/octet-stream" \
       -o "$ZIP_FILE" "$ASSET_URL"

  echo "[INFO] Downloaded ${ZIP_FILE}"
else
  echo "[INFO] Using cached zip file: ${ZIP_FILE}"
fi

echo "===== Step 3: Clearing and deploying to current/ ====="
rm -rf "${CURRENT_DIR:?}/"*
unzip -o "$ZIP_FILE" -d "$CURRENT_DIR"
echo "[INFO] Unzipped contents to ${CURRENT_DIR}"

echo "===== Step 4: Writing package.json ====="
cat > "${CURRENT_DIR}/package.json" <<EOF
{
  "name": "${APP_NAME}",
  "version": "${VERSION}"
}
EOF
echo "[INFO] package.json written to ${CURRENT_DIR}/package.json"

echo "===== Deployment Complete ====="
