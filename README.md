# WhatsApp Webview Desktop

Lightweight desktop wrapper for [WhatsApp Web](https://web.whatsapp.com), built with Go, Microsoft Edge WebView2 (Windows), and WebKit WKWebView (macOS).

## Features

- Native desktop window instead of heavy Electron (WebView2 on Windows, WKWebView on macOS).
- Persistent WhatsApp session across app restarts.
- Cookies, LocalStorage, IndexedDB, and service-worker data stored in a dedicated profile.
- Native notification bridge (Windows Toast and macOS Notification Center).
- Camera and microphone access for WhatsApp voice and video calls.
- Single-instance protection.
- High-DPI display support.
- Small native executable.

## Requirements

### Windows
- Windows 10 or newer.
- Microsoft Edge WebView2 Runtime.

### macOS
- macOS 11.0 (Big Sur) or newer (Apple Silicon and Intel supported).

## Session Data

Profile data is stored at:
- **Windows**: `%APPDATA%\WhatsAppDesktopLight\UserData`
- **macOS**: `~/Library/Application Support/WhatsAppDesktopLight/UserData`

Do not delete this folder if the existing login session must remain available. Closing the app does not clear session data.

## Build From Source

### macOS
Run the build script to compile the binary and generate `WhatsApp.app`:

```bash
chmod +x build_mac.sh
./build_mac.sh
```

To run:
```bash
open WhatsApp.app
```

### Windows
Install Go and a Windows C compiler, then run:

```powershell
go mod download
go build -ldflags="-H windowsgui -s -w" -o WhatsApp.exe .
```

## Project Files

- `main.go`: Shared entry point, user-agent spoofing, and notification polyfill scripts.
- `app_darwin.go`: macOS implementation with WebKit (WKWebView), AppleScript notifications, and file lock.
- `app_windows.go`: Windows implementation with Edge WebView2, DWM dark mode, and toast notifications.
- `build_mac.sh`: macOS build and .app packaging script.
- `app.manifest`: Windows DPI and application manifest.
- `resource.rc`: Windows icon and manifest resource definitions.
- `icon.ico`: Application icon.

## Privacy

This app loads WhatsApp Web directly. Chat data and authentication state are handled by WhatsApp Web and stored locally in the WebView2 profile above. This project is not affiliated with WhatsApp or Meta.

## License

No license has been declared yet.
