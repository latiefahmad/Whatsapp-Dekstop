#!/usr/bin/env bash
set -e

APP_NAME="WhatsApp"
BUNDLE_DIR="${APP_NAME}.app"
CONTENTS_DIR="${BUNDLE_DIR}/Contents"
MACOS_DIR="${CONTENTS_DIR}/MacOS"
RESOURCES_DIR="${CONTENTS_DIR}/Resources"

echo "Building universal binary for macOS (arm64 + x86_64)..."
LIPO="/usr/bin/lipo"
if [ -x "$LIPO" ] && \
   CGO_ENABLED=1 GOARCH=arm64 CGO_CFLAGS="-arch arm64" CGO_LDFLAGS="-arch arm64" GOTOOLCHAIN=local go build -ldflags="-s -w" -o "${APP_NAME}_arm64" . 2>/dev/null && \
   CGO_ENABLED=1 GOARCH=amd64 CGO_CFLAGS="-arch x86_64" CGO_LDFLAGS="-arch x86_64" GOTOOLCHAIN=local go build -ldflags="-s -w" -o "${APP_NAME}_amd64" . 2>/dev/null; then
    "$LIPO" -create -output "${APP_NAME}_mac" "${APP_NAME}_arm64" "${APP_NAME}_amd64"
    rm -f "${APP_NAME}_arm64" "${APP_NAME}_amd64"
    echo "Successfully built universal binary (Apple Silicon + Intel)"
else
    echo "Falling back to native host architecture build..."
    GOTOOLCHAIN=local go build -ldflags="-s -w" -o "${APP_NAME}_mac" .
fi

echo "Packaging ${BUNDLE_DIR}..."
rm -rf "${BUNDLE_DIR}"
mkdir -p "${MACOS_DIR}" "${RESOURCES_DIR}"

if [ ! -f "AppIcon.icns" ]; then
    echo "Generating AppIcon.icns..."
    python3 generate_icns.py
fi
cp "AppIcon.icns" "${RESOURCES_DIR}/AppIcon.icns"

cp "${APP_NAME}_mac" "${MACOS_DIR}/${APP_NAME}"
chmod +x "${MACOS_DIR}/${APP_NAME}"

cat << 'EOF' > "${CONTENTS_DIR}/Info.plist"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>WhatsApp</string>
    <key>CFBundleIdentifier</key>
    <string>com.whatsapp.desktop.light</string>
    <key>CFBundleName</key>
    <string>WhatsApp Desktop</string>
    <key>CFBundleIconFile</key>
    <string>AppIcon</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>1.3.1</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>NSRequiresAquaSystemAppearance</key>
    <false/>
    <key>NSCameraUsageDescription</key>
    <string>WhatsApp requires camera access for video calls.</string>
    <key>NSMicrophoneUsageDescription</key>
    <string>WhatsApp requires microphone access for voice and video calls.</string>
</dict>
</plist>
EOF

if command -v hdiutil >/dev/null 2>&1; then
    echo "Creating DMG installer (WhatsApp-macOS-Universal.dmg)..."
    DMG_DIR="dmg_temp"
    rm -rf "${DMG_DIR}" "WhatsApp-macOS-Universal.dmg"
    mkdir -p "${DMG_DIR}"
    cp -R "${BUNDLE_DIR}" "${DMG_DIR}/"
    ln -s /Applications "${DMG_DIR}/Applications"
    if [ -f "dmg_background.png" ]; then
        mkdir -p "${DMG_DIR}/.background"
        cp "dmg_background.png" "${DMG_DIR}/.background/background.png"
    fi
    hdiutil create -volname "WhatsApp Desktop" -srcfolder "${DMG_DIR}" -ov -format UDZO "WhatsApp-macOS-Universal.dmg" || true
    rm -rf "${DMG_DIR}"
fi

echo "Done! Built ${APP_NAME}_mac, ${BUNDLE_DIR}, and DMG installer."
