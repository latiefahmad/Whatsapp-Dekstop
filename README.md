# WhatsApp Desk

[![Release](https://img.shields.io/github/v/release/vianziro/Whatsapp-Dekstop?label=Latest%20Release&color=00a884)](https://github.com/vianziro/Whatsapp-Dekstop/releases/latest)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Windows%20%7C%20Linux-blue)](https://github.com/vianziro/Whatsapp-Dekstop/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

An ultra-lightweight, blazing-fast, and memory-efficient WhatsApp Desktop client for **macOS**, **Windows**, and **Linux**, built with Go.

Unlike official WhatsApp Desktop which bundles a heavy Electron runtime, WhatsApp Desk leverages native OS webview engines (Apple WebKit on macOS, Microsoft Edge WebView2 on Windows, and WebKitGTK on Linux). This saves gigabytes of disk space and hundreds of megabytes of RAM.

---

## 📥 Download (Release v1.5.4)

Ready-to-use binaries and installers are available directly from the official release links below:

| Operating System | Installer / Binary | Architecture | Format | Size |
| :--- | :--- | :--- | :--- | :--- |
| **macOS (DMG)** | [📥 **Download WhatsApp-Desk-macOS-Universal.dmg**](https://github.com/vianziro/Whatsapp-Dekstop/releases/latest/download/WhatsApp-Desk-macOS-Universal.dmg) | Universal (Apple Silicon & Intel) | Apple Disk Image (.dmg) | ~6.5 MB |
| **macOS (ZIP)** | [📥 **Download WhatsApp-Desk-macOS-Universal.zip**](https://github.com/vianziro/Whatsapp-Dekstop/releases/latest/download/WhatsApp-Desk-macOS-Universal.zip) | Universal (Apple Silicon & Intel) | Portable App (.zip) | ~5.6 MB |
| **Windows (ZIP)** | [📥 **Download WhatsApp-Desk-Windows-x64.zip**](https://github.com/vianziro/Whatsapp-Dekstop/releases/latest/download/WhatsApp-Desk-Windows-x64.zip) | 64-bit (x64) | Complete Portable Bundle (.zip) | ~5.3 MB |
| **Windows (EXE)** | [📥 **Download WhatsApp.exe**](https://github.com/vianziro/Whatsapp-Dekstop/releases/latest/download/WhatsApp.exe) | 64-bit (x64) | Standalone Executable (.exe) | ~6.4 MB |
| **Linux (DEB)** | [📥 **Download WhatsApp-Desk-Linux-amd64.deb**](https://github.com/vianziro/Whatsapp-Dekstop/releases/latest/download/WhatsApp-Desk-Linux-amd64.deb) | 64-bit (x86_64) | Debian / Ubuntu Package (.deb) | ~2.2 MB |
| **Linux (Tarball)** | [📥 **Download WhatsApp-Desk-Linux-x64.tar.gz**](https://github.com/vianziro/Whatsapp-Dekstop/releases/latest/download/WhatsApp-Desk-Linux-x64.tar.gz) | 64-bit (x86_64) | Portable Tarball (.tar.gz) | ~2.8 MB |

> Full release history, release notes, and older binary archives are accessible on the [GitHub Releases](https://github.com/vianziro/Whatsapp-Dekstop/releases) page.

---

## 🚀 Resource Utilization Comparison

| Metric | WhatsApp Desk | Official WhatsApp (Electron) |
| :--- | :--- | :--- |
| **Binary / App Size** | **~2.5 MB (macOS) / ~7 MB (Win)** | 470 MB – 900 MB+ |
| **Physical RAM Usage (RSS)** | **~35 MB (Core) / ~80–120 MB (Active Chat)** | 600 MB – 1.5 GB+ |
| **Rendering Engine** | Native OS (WebKit / WebView2 / WebKitGTK) | Bundled Chromium + Node.js |
| **Cold Startup Time** | **< 0.5s (Instant)** | 3 – 6 seconds |
| **CPU & Battery Drain** | Extremely Low (Hardware Accelerated) | High resource consumption |
| **Login Session** | Securely saved in isolated local profile | Saved in local profile |
| **Chat Encryption** | Official WhatsApp End-to-End (Signal Protocol) | Official WhatsApp End-to-End (Signal Protocol) |

---

## ✨ Features

### Cross-Platform Features (Windows, macOS & Linux)
- **Persistent Login Session**: Session state (cookies, local storage, indexedDB) is preserved in a dedicated, isolated profile directory. No need to re-scan QR codes each time you launch the app.
- **📄 Multi-Format In-App Document Preview**: Instant, dark-themed preview modal for **PDF**, **Excel** (`.xlsx`, `.xls`, `.csv`), **Word** (`.docx`, `.doc`), and **Plain Text** (`.txt`) documents directly inside the app. Includes quick actions to open files in default system applications (Apple Preview, Numbers, Excel, Pages, Word) or save them locally.
- **🛠️ Media Viewer & Quick Dismiss**: Enhanced WhatsApp Web media viewer interaction with immediate exit on ✕ (Close) button clicks or the `Escape` key, completely avoiding stuck loading spinners or hijacked event loops.
- **🎛️ Quick Control Center Toolbar**: Modern in-flow toolbar button next to the chat search bar to toggle Privacy Mode, Always on Top, Audio Mute, and open Settings with a single click—no need to memorize keyboard shortcuts.
- **🎨 Dynamic Theme Synchronization**: Supports Auto (system-matching), Dark, and Light themes seamlessly synchronized with WhatsApp Web.
- **🔒 Privacy Mode (Anti-Peeking)**: Press `Ctrl + Shift + P` (Windows/Linux) or `Cmd + Shift + P` (macOS) to blur messages, media, and contact names. Simply hover over any message to temporarily read it.
- **📌 Always on Top (Window Pinning)**: Press `Ctrl + Shift + T` / `Cmd + Shift + T` to keep WhatsApp pinned above all other windows during multitasking.
- **🔇 Instant Audio Mute**: Press `Ctrl + Shift + M` / `Cmd + Shift + M` to quickly mute or unmute chat notification sounds.
- **🚀 Auto-Start on System Boot**: Press `Ctrl + Shift + S` / `Cmd + Shift + S` to configure the app to launch automatically when your computer boots up.
- **🔄 Chat Reload & Hard Refresh**: Press `F5` / `Ctrl/Cmd + R` to reload chat view, or `Ctrl/Cmd + Shift + R` to hard refresh and purge interface caches.
- **⚡ In-App Auto-Updater**: Automatically checks for new GitHub Releases in the background with manual shortcut `Ctrl/Cmd + Shift + U`. Update the application directly from the chat window with one click without manually downloading installers.
- **💾 Permanent & Custom Downloads Directory**: Chat documents, files, and images are saved permanently to disk (defaults to `Downloads/WhatsApp Downloads`). The path can be customized via the Settings window (`Ctrl/Cmd + ,`) with a native folder dialog and duplicate filename auto-incrementing.
- **🔔 Minimalist Floating HUD Feedback**: Every shortcut action displays a clean, unobtrusive toast indicator inside the interface.
- **🎙️ Full Media & Call Support**: Complete support for voice notes recording, audio/video playback, and microphone/camera hardware access for voice and video calling.
- **🌐 External Link Interception**: Web links clicked inside chat open automatically in your default system browser (Chrome, Safari, Edge, Firefox, etc.) without disrupting the WhatsApp window.
- **🔍 Display Zoom Control**: Easily adjust text and interface scaling with `Ctrl/Cmd +`, `Ctrl/Cmd -`, and `Ctrl/Cmd 0`.
- **🛡️ Single-Instance Guard**: Prevents accidental multiple instances of the app from running concurrently.

### Windows System Integration
- **Windows 10 & 11 Full Compatibility (v1.5.3 Fix)**: Chromium process protection utilizing Windows Job Object with *Silent Breakaway* permission (`JOB_OBJECT_LIMIT_SILENT_BREAKAWAY_OK`), ensuring clean, instant startup without black screens across all Windows 10 and 11 devices.
- **Microsoft Edge WebView2**: Powered by the Evergreen WebView2 runtime for native DirectX hardware acceleration and minimal RAM consumption.
- **Windows Toast Notification Focus**: Clicking a native Windows toast notification instantly brings the WhatsApp window to the foreground.
- **Dark Mode Title Bar**: Seamless Win32 Desktop Window Manager (DWM) integration providing a native dark window frame matching WhatsApp Web.
- **Isolated User Data Profile**: Session data is stored in `%APPDATA%\WhatsAppDesk\UserData\`.
- **Fully Portable**: Available as a standalone executable (`WhatsApp.exe`) and portable zip archive (`WhatsApp-Desk-Windows-x64.zip`).

### macOS System Integration
- **Apple WebKit (Cocoa) & PDFKit**: Native WKWebView with Cocoa PDFKit integration for high-speed document rendering.
- **Universal Binary**: A single universal binary natively supporting both Apple Silicon (M1/M2/M3/M4) and Intel (x86_64) Macs.
- **Dock Unread Badge**: Displays live unread message counts directly on the application Dock icon.
- **Close-to-Hide Window Behavior**: Closing the window with `Cmd + W` or the red title bar button hides the window to the background; background notifications stay active, and the window pops up instantly upon clicking the Dock icon. Full exit is performed via `Cmd + Q`.
- **Complete Cocoa Menu Bar**: Standard macOS shortcuts such as `Cmd + C` (Copy), `Cmd + V` (Paste), `Cmd + X` (Cut), `Cmd + A` (Select All), and `Cmd + Z` (Undo).
- **Persistent Window Frame Memory**: Remembers your preferred window size and screen coordinates across restarts.

---

## 📖 Installation Guide & Security Warnings

Because this project is distributed independently and open-source (without an annual $99/year Apple Developer ID certificate or expensive Microsoft EV Code Signing certificate), your operating system might display an initial security prompt.

Here are quick, standard instructions for each operating system:

---

### 🍏 macOS Installation

#### Steps:
1. Download [**WhatsApp-Desk-macOS-Universal.dmg**](https://github.com/vianziro/Whatsapp-Dekstop/releases/latest/download/WhatsApp-Desk-macOS-Universal.dmg).
2. Double-click the `.dmg` file to mount it.
3. Drag and drop the **WhatsApp Desk** icon into your **Applications** folder.
4. Launch WhatsApp Desk from Launchpad, Spotlight (`Cmd + Space`), or Applications.

#### ⚠️ How to Resolve: *"WhatsApp Desk is damaged and can't be opened"* or *Unidentified Developer*
If downloaded via Safari or Chrome, macOS Gatekeeper automatically applies an extended quarantine attribute. This is standard macOS protection for third-party open-source applications outside the Mac App Store.

Choose either simple method below to launch the app:

- **Method 1: Right-Click Open (Easiest, No Terminal Required)**:
  1. Open your **Applications** folder in Finder.
  2. **Right-Click (or hold `Control` and click)** on **WhatsApp Desk**.
  3. Select **Open** from the context menu.
  4. In the security dialog, click **Open**.
  5. The application will launch immediately and will open normally with a regular left-click thereafter.

- **Method 2: Via Terminal (Remove Quarantine Flag)**:
  Open **Terminal** and run:
  ```bash
  xattr -cr "/Applications/WhatsApp Desk.app"
  ```
  *(Or if the app is still in your Downloads directory: `xattr -cr ~/Downloads/"WhatsApp Desk.app"`)*

*Note: Requires macOS 11.0 (Big Sur) or higher. Native Universal binary for both Apple Silicon (M1/M2/M3/M4) and Intel.*

---

### 🪟 Windows Installation

#### Steps:
1. Download [**WhatsApp.exe**](https://github.com/vianziro/Whatsapp-Dekstop/releases/latest/download/WhatsApp.exe) or [**WhatsApp-Desk-Windows-x64.zip**](https://github.com/vianziro/Whatsapp-Dekstop/releases/latest/download/WhatsApp-Desk-Windows-x64.zip).
2. Place the file in your preferred folder (e.g., `C:\Apps\WhatsAppDesk` or Desktop).
3. Create a shortcut to your Desktop or Start Menu if desired.
4. Double-click `WhatsApp.exe` to run (portable, no installer needed).

#### ⚠️ Bypassing Windows SmartScreen (*"Windows protected your PC"*):
1. When the blue SmartScreen prompt appears, click **"More info"**.
2. Click **"Run anyway"**.
3. The app will launch immediately.

*Note: Requires Microsoft Edge WebView2 Runtime (pre-installed by default on Windows 11 and up-to-date Windows 10).*

---

### 🐧 Linux Installation

#### Option A: Debian / Ubuntu (.deb package)
```bash
sudo dpkg -i WhatsApp-Desk-Linux-amd64.deb
sudo apt-get install -f # resolve any missing dependencies if needed
```
The app will be installed and appear in your desktop Application Menu / App Launcher.

#### Option B: Portable Tarball (.tar.gz)
```bash
tar -xzf WhatsApp-Desk-Linux-x64.tar.gz
cd dist_linux
chmod +x whatsapp-desk
./whatsapp-desk
```

---

## ⌨️ Keyboard Shortcuts

| macOS | Windows / Linux | Action |
| :--- | :--- | :--- |
| `Cmd + Shift + P` | `Ctrl + Shift + P` | Toggle Privacy Mode (blur messages & media in public) |
| `Cmd + Shift + T` | `Ctrl + Shift + T` | Toggle Always on Top (pin window above other apps) |
| `Cmd + Shift + M` | `Ctrl + Shift + M` | Toggle Mute chat notification audio |
| `Cmd + Shift + S` | `Ctrl + Shift + S` | Toggle Auto-Start on system boot |
| `Cmd + R` / `F5` | `Ctrl + R` / `F5` | Reload WhatsApp conversation view |
| `Cmd + Shift + R` | `Ctrl + Shift + R` | Hard refresh (reload and flush UI cache) |
| `Cmd + Shift + U` | `Ctrl + Shift + U` | Check for and install updates (In-App Auto-Updater) |
| `Cmd + ,` | `Ctrl + ,` | Open Settings (Downloads location & shortcuts info) |
| `Cmd + Shift + D` | `Ctrl + Shift + D` | Open downloads storage folder in Finder / Explorer |
| `Cmd + Shift + H` | `Ctrl + Shift + H` | Re-open welcome guide & shortcuts overview |
| `Cmd + +` / `Cmd + =` | `Ctrl + +` / `Ctrl + =` | Zoom in interface |
| `Cmd + -` | `Ctrl + -` | Zoom out interface |
| `Cmd + 0` | `Ctrl + 0` | Reset interface zoom to default (100%) |
| `Cmd + C` / `Cmd + V` | `Ctrl + C` / `Ctrl + V` | Copy / Paste text or media |
| `Cmd + W` | `Alt + F4` | Hide window to background (app remains active) |
| `Cmd + Q` | `Alt + F4` | Quit application completely |

---

## 🔒 Security & Data Privacy

- **Direct to Meta (Zero Middleman)**: This application connects directly to official WhatsApp Web servers (`https://web.whatsapp.com`). There are no relay servers, third-party proxies, analytics trackers, or intermediate backend services.
- **Official End-to-End Encryption**: All chats and voice/video calls remain encrypted end-to-end using WhatsApp's native Signal Protocol implementation inside the browser engine.
- **Isolated Local Storage**: User credentials, session tokens, and local cache are stored securely in your OS user directory:
  - **macOS**: `~/Library/Application Support/WhatsAppDesk/UserData/`
  - **Windows**: `%APPDATA%\WhatsAppDesk\UserData\`
  - **Linux**: `~/.config/whatsapp-desk/`
- **Zero Ads & Telemetry**: 100% free of advertisements, data tracking, or telemetry overhead.

---

## 📋 Changelog

- **v1.5.4 (Latest)**:
  - 📄 **Reliable macOS PDF Preview**: Routes WhatsApp document downloads into the native PDFKit viewer instead of the unsupported WKWebView PDF plugin path.
  - ⬇️ **Updated Viewer Toolbar Detection**: Detects the currently visible download control without depending on legacy WhatsApp `data-testid` attributes.
  - ↩️ **Close Returns to Chat**: Closing the native PDF preview also dismisses the WhatsApp document viewer and restores the active chat.
  - 🧠 **Safer Viewer Lifecycle**: Prevents repeated synthetic Escape events and releases temporary document URLs after preview.

- **v1.5.3 (Deprecated)**:
  - 📄 **Multi-Format In-App Document Preview**: Instant preview modal for **PDF**, **Excel** (`.xlsx`, `.xls`, `.csv`), **Word** (`.docx`, `.doc`), and **Text** (`.txt`) directly in chat with native app opening support.
  - 🛠️ **Media Viewer & ✕ Close Button Fix**: Resolved an issue where clicking the close button or toolbar buttons in WhatsApp Web's media viewer failed to return to the active chat screen.
  - 🪟 **Windows 10 & 11 Black Screen Fix**: Configured Windows Job Object with *Silent Breakaway* permission (`JOB_OBJECT_LIMIT_SILENT_BREAKAWAY_OK`), eliminated startup process termination conflicts, and optimized Chromium command-line flags.
  - 🍎 **macOS Menu Stability**: Bound native menu actions to a persistent window reference (`appWindow()`) and protected WebKit evaluation against unhandled Promise rejections.
- **v1.5.2**:
  - 📄 **In-App Document Preview**: In-app PDF viewer with dark blurred backdrop.
  - 🔄 **Stuck Spinner Auto-Dismiss**: Automatically dismissed stuck media viewer spinners caused by cross-origin Adobe Acrobat Web SDK limitations.
  - 📂 **Quick Action Controls**: Instant "Open in System App" (Apple Preview / Default Viewer) and save file buttons.
- **v1.5.1**:
  - 🔔 **Notification Window Focusing**: Clicking native OS push notifications automatically brings the chat window to the foreground.
  - 🌓 **Automatic Theme Synchronization**: Dynamic Dark, Light, and System theme synchronization matching OS appearance settings.
  - 🚀 **GPU Smooth Scrolling**: Eliminated aggressive memory cleanup during scrolling.
- **v1.5.0**:
  - 🔄 **In-App Auto-Updater**: One-click background updates via GitHub Releases.
  - 📁 **Permanent Downloads Manager**: Auto-saves received files to `Downloads/WhatsApp Downloads` with custom native folder picker.
  - 🎛️ **Control Center Toolbar**: Quick access toolbar in chat header for privacy mode, always-on-top, and audio mute.
- **v1.4.0**:
  - 🍏 **Universal Binary for macOS**: Native support for Apple Silicon and Intel Macs.
  - 📐 **Responsive Resizing**: Compact window resizing with minimum width constraints.
  - 📘 **Onboarding Guide**: Clean intro modal explaining key features and shortcuts.

---

## 📦 Official Binary Releases

For optimal performance and ease of use, download the precompiled official binaries from the [Download Table](#-download-release-v154) above. Release artifacts are built and package-verified for macOS (Apple Silicon & Intel), Windows 10/11, and Linux.

---

## License

This project is licensed under the [MIT License](LICENSE).

---

## Disclaimer

This project is an independent open-source software and is not affiliated, associated, authorized, endorsed by, or in any way officially connected with WhatsApp or Meta Platforms, Inc. WhatsApp is a registered trademark of Meta Platforms, Inc.
