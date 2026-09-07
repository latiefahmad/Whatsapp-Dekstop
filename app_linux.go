//go:build linux

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	webview "github.com/webview/webview_go"
)

const (
	userAgentLinux = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"
)

var (
	isAlwaysOnTopLinux = false
)

func checkSingleInstance() (*os.File, bool) {
	dataDir := getUserDataDir()
	lockPath := filepath.Join(dataDir, "app.lock")
	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, true
	}
	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		// Already running
		return nil, false
	}
	return file, true
}

func getUserDataDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		home, _ := os.UserHomeDir()
		configDir = filepath.Join(home, ".config")
	}
	dir := filepath.Join(configDir, "WhatsAppDesktopLight")
	_ = os.MkdirAll(dir, 0755)
	return dir
}

func showNativeNotification(title, message string) {
	_ = exec.Command("notify-send", "-a", "WhatsApp Desktop Light", title, message).Run()
}

func toggleAlwaysOnTopLinux() bool {
	isAlwaysOnTopLinux = !isAlwaysOnTopLinux
	state := "remove"
	if isAlwaysOnTopLinux {
		state = "add"
	}
	if path, err := exec.LookPath("wmctrl"); err == nil && path != "" {
		_ = exec.Command("wmctrl", "-r", windowTitle, "-b", fmt.Sprintf("%s,above", state)).Run()
	}
	return isAlwaysOnTopLinux
}

func getAutoStartDesktopPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "autostart", "whatsapp-desktop-light.desktop")
}

func toggleAutoStartLinux() bool {
	p := getAutoStartDesktopPath()
	if p == "" {
		return false
	}
	if _, err := os.Stat(p); err == nil {
		_ = os.Remove(p)
		return false
	}

	execPath, err := os.Executable()
	if err != nil {
		return false
	}

	_ = os.MkdirAll(filepath.Dir(p), 0755)
	desktopContent := fmt.Sprintf(`[Desktop Entry]
Type=Application
Version=1.0
Name=WhatsApp Desktop Light
Comment=Lightweight WhatsApp Desktop Client
Exec=%s
Icon=whatsapp-desktop-light
Terminal=false
Categories=Network;InstantMessaging;
StartupNotify=true
`, execPath)

	err = os.WriteFile(p, []byte(desktopContent), 0644)
	return err == nil
}

func loadWindowState(dir string) *WindowState {
	path := filepath.Join(dir, "window_state.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var state WindowState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil
	}
	if state.Width < 450 || state.Height < 320 {
		return nil
	}
	return &state
}

func saveWindowState(dir string, width, height int) {
	if width >= 450 && height >= 320 {
		state := WindowState{
			Width:  float64(width),
			Height: float64(height),
		}
		data, err := json.MarshalIndent(state, "", "  ")
		if err == nil {
			_ = os.WriteFile(filepath.Join(dir, "window_state.json"), data, 0644)
		}
	}
}

func runApp() {
	_, isSingle := checkSingleInstance()
	if !isSingle {
		fmt.Println("WhatsApp Desktop Light is already running.")
		os.Exit(0)
	}

	userDataDir := getUserDataDir()

	w := webview.New(false)
	if w == nil {
		log.Fatalln("Gagal inisialisasi WebKitGTK Webview")
	}
	defer w.Destroy()

	w.SetTitle(windowTitle)
	w.SetSize(windowWidth, windowHeight, webview.HintNone)

	// Bind native notification bridge
	_ = w.Bind("sendNativeNotification", func(title, body string) {
		go showNativeNotification(title, body)
	})

	// Bind external link handler (xdg-open)
	_ = w.Bind("openExternalLink", func(rawURL string) {
		if strings.HasPrefix(rawURL, "http://") || strings.HasPrefix(rawURL, "https://") {
			go func() {
				_ = exec.Command("xdg-open", rawURL).Start()
			}()
		}
	})

	// Bind dock badge (no-op on linux)
	_ = w.Bind("updateDockBadge", func(badge string) {
		// Not supported on standard Linux window managers
	})

	// Bind Always on Top toggle
	_ = w.Bind("toggleAlwaysOnTopNative", func() bool {
		return toggleAlwaysOnTopLinux()
	})

	// Bind Auto-Start toggle
	_ = w.Bind("toggleAutoStartNative", func() bool {
		return toggleAutoStartLinux()
	})

	// Bind in-app auto updater
	_ = w.Bind("checkForUpdateNative", func(manual bool) UpdateInfo {
		info, err := checkForUpdate(appVersion)
		if err != nil {
			return UpdateInfo{CurrentVersion: appVersion}
		}
		return *info
	})

	_ = w.Bind("startUpdateNative", func(downloadURL string) {
		go func() {
			_ = executeUpdate(w, downloadURL)
		}()
	})

	// Bind download and settings handlers
	_ = w.Bind("saveDownloadedFileNative", func(filename, dataURI string) string {
		path, err := saveDownloadedFile(filename, dataURI)
		if err != nil {
			return ""
		}
		return path
	})

	_ = w.Bind("getDownloadDirNative", func() string {
		s := loadSettings()
		return s.DownloadDir
	})

	_ = w.Bind("chooseDownloadDirNative", func() string {
		selected, err := chooseFolderDialog()
		if err != nil || selected == "" {
			return ""
		}
		s := loadSettings()
		s.DownloadDir = selected
		_ = saveSettings(s)
		return selected
	})

	_ = w.Bind("openDownloadDirNative", func() bool {
		s := loadSettings()
		_ = openFolderInFileManager(s.DownloadDir)
		return true
	})

	_ = w.Bind("resetDownloadDirNative", func() string {
		s := loadSettings()
		s.DownloadDir = getDefaultDownloadDir()
		_ = saveSettings(s)
		return s.DownloadDir
	})

	_ = w.Bind("getAppThemeNative", func() string {
		s := loadSettings()
		return s.Theme
	})

	_ = w.Bind("setAppThemeNative", func(theme string) string {
		return saveTheme(theme)
	})

	w.Init(getInitScript(userAgentLinux))
	w.Navigate(appURL)

	// Check for updates in the background after startup
	go func() {
		time.Sleep(5 * time.Second)
		info, err := checkForUpdate(appVersion)
		if err == nil && info != nil && info.Available {
			w.Dispatch(func() {
				script := fmt.Sprintf("if (window.showUpdateBanner) { window.showUpdateBanner(%q, %q, %q); }",
					info.LatestVersion, info.ReleaseTitle, info.DownloadURL)
				w.Eval(script)
			})
		}
	}()

	defer saveWindowState(userDataDir, windowWidth, windowHeight)
	w.Run()
}
