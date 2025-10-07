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
NEED_ENV_FILE=true

APP_NAME="$1"
VERSION="$2"

if [ -z "$APP_NAME" ] || [ -z "$VERSION" ]; then
  echo "[ERROR] Usage: $0 <app_name> <version>"
  echo "Example: $0 ci-react-web-prod v1.2.3"
  exit 1
fi

if [ "$NEED_ENV_FILE" = true ]; then
  echo "===== Step 0: Checking .env file ====="
  ENV_FILE="${SCRIPT_DIR}/.env"
  if [ ! -f "$ENV_FILE" ]; then
    echo "[ERROR] .env file is required but not found at $ENV_FILE"
    exit 1
  else
    echo "[INFO] .env file exists: $ENV_FILE"
  fi
fi

echo "===== Step 1: Preparing directories ====="

if [ ! -d "$RELEASES_DIR" ]; then
  echo "[INFO] Creating directory: $RELEASES_DIR"
  mkdir -p "$RELEASES_DIR"
else
  echo "[INFO] Directory exists: $RELEASES_DIR"
fi

if [ ! -d "$CURRENT_DIR" ]; then
  echo "[INFO] Creating directory: $CURRENT_DIR"
  mkdir -p "$CURRENT_DIR"
else
  echo "[INFO] Directory exists: $CURRENT_DIR"
fi

echo "===== Step 2: Writing pm2.config.js ====="

cat > "${SCRIPT_DIR}/pm2.config.js" <<EOF
const path = require("path");

const APP_NAME = "${APP_NAME}";

module.exports = {
  apps: [
    {
      name: APP_NAME,
      cwd: __dirname,
      script: path.join("./current", APP_NAME),
      watch: false,
      autorestart: true,
      max_restarts: 10
    }
  ]
};
EOF

echo "[INFO] pm2.config.js has been written successfully."

echo "===== Step 3: Writing package.json ====="

cat > "${SCRIPT_DIR}/package.json" <<EOF
{
  "name": "${APP_NAME}",
  "version": "${VERSION}"
}
EOF

echo "[INFO] package.json has been written successfully."

echo "===== Step 4: Download and deploy binary ====="

TARGET_FILE="${RELEASES_DIR}/${APP_NAME}-${VERSION}"

if [ -f "$TARGET_FILE" ]; then
  echo "[INFO] Binary already exists locally: $TARGET_FILE"
else
  echo "[INFO] Querying GitHub API for release asset..."

  RELEASE_JSON=$(curl -s -H "Authorization: token ${GITHUB_TOKEN}" \
    "https://api.github.com/repos/${REPO}/releases/tags/${VERSION}")

  if [ -z "$RELEASE_JSON" ] || [ "$(echo "$RELEASE_JSON" | jq -r '.message')" = "Not Found" ]; then
    echo "[ERROR] Could not find release tag '${VERSION}'"
    exit 1
  fi

  ASSET_ID=$(echo "$RELEASE_JSON" | jq ".assets[] | select(.name==\"${APP_NAME}\") | .id")

  if [ -z "$ASSET_ID" ]; then
    echo "[ERROR] Could not find asset named '${APP_NAME}' in release '${VERSION}'"
    exit 1
  fi

  echo "[INFO] Found asset ID: $ASSET_ID"
  echo "[INFO] Downloading asset..."

  curl -L \
    -H "Authorization: token ${GITHUB_TOKEN}" \
    -H "Accept: application/octet-stream" \
    -o "$TARGET_FILE" \
    "https://api.github.com/repos/${REPO}/releases/assets/${ASSET_ID}"

  chmod +x "$TARGET_FILE"

  FILE_TYPE=$(file "$TARGET_FILE")
  echo "[INFO] Downloaded file type: $FILE_TYPE"

  if [[ "$FILE_TYPE" != *"executable"* ]]; then
    echo "[ERROR] Downloaded file is not an executable!"
    head "$TARGET_FILE"
    exit 1
  fi
fi

echo "===== Step 5: Updating symlink ====="
ln -sfn "$TARGET_FILE" "${CURRENT_DIR}/${APP_NAME}"
echo "[INFO] Symlink updated: ${CURRENT_DIR}/${APP_NAME} -> $TARGET_FILE"

echo "===== Step 6: Restarting PM2 process ====="

if pm2 describe "${APP_NAME}" > /dev/null; then
  echo "[INFO] Reloading PM2 process..."
  pm2 reload "${SCRIPT_DIR}/pm2.config.js"
else
  echo "[INFO] Starting PM2 process..."
  pm2 start "${SCRIPT_DIR}/pm2.config.js"
fi

echo "[SUCCESS] Switched to version $VERSION."

pm2 list