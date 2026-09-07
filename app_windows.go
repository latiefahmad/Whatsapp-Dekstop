//go:build windows

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
	"unsafe"

	"github.com/go-toast/toast"
	"github.com/jchv/go-webview2"
	"golang.org/x/sys/windows"
)

var (
	kernel32          = windows.NewLazySystemDLL("kernel32.dll")
	user32            = windows.NewLazySystemDLL("user32.dll")
	dwmapi            = windows.NewLazySystemDLL("dwmapi.dll")
	procCreateMutex   = kernel32.NewProc("CreateMutexW")
	procFindWindow    = user32.NewProc("FindWindowW")
	procSetFgWindow   = user32.NewProc("SetForegroundWindow")
	procShowNormal    = user32.NewProc("ShowWindow")
	procDwmSetAttr    = dwmapi.NewProc("DwmSetWindowAttribute")
	procGetWindowLong = user32.NewProc("GetWindowLongW")
	procSetWindowLong = user32.NewProc("SetWindowLongW")
	procSetWindowPos  = user32.NewProc("SetWindowPos")
	procGetWindowRect = user32.NewProc("GetWindowRect")
	procMoveWindow    = user32.NewProc("MoveWindow")

	isAlwaysOnTopWin = false
)

const (
	mutexName = "WhatsAppDesktopSingleInstanceMutex"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"

	// Win32 Window Styles for Dynamic Resizability
	GWL_STYLE        = 0xFFFFFFF0 // -16
	WS_THICKFRAME    = 0x00040000
	WS_MAXIMIZEBOX   = 0x00010000
	WS_MINIMIZEBOX   = 0x00020000
	SWP_FRAMECHANGED = 0x0020
	SWP_NOMOVE       = 0x0002
	SWP_NOSIZE       = 0x0001
	SWP_NOZORDER     = 0x0004

	// Window Z-Order constants for Always On Top
	HWND_TOPMOST   = ^uintptr(0) // -1
	HWND_NOTOPMOST = ^uintptr(1) // -2

	// DWM Window Attributes for Dark Theme
	DWMWA_USE_IMMERSIVE_DARK_MODE_BEFORE_20H1 = 19
	DWMWA_USE_IMMERSIVE_DARK_MODE             = 20
	DWMWA_CAPTION_COLOR                       = 35
	DWMWA_TEXT_COLOR                          = 36
)

func toggleAlwaysOnTop(hwnd uintptr) bool {
	isAlwaysOnTopWin = !isAlwaysOnTopWin
	target := uintptr(HWND_NOTOPMOST)
	if isAlwaysOnTopWin {
		target = uintptr(HWND_TOPMOST)
	}
	procSetWindowPos.Call(hwnd, target, 0, 0, 0, 0, uintptr(SWP_NOMOVE|SWP_NOSIZE))
	return isAlwaysOnTopWin
}

func toggleAutoStartWindows() bool {
	runKey := `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`
	valName := "WhatsAppDesktop"

	// Check if already configured
	err := exec.Command("reg", "query", runKey, "/v", valName).Run()
	if err == nil {
		// Key exists, remove it
		_ = exec.Command("reg", "delete", runKey, "/v", valName, "/f").Run()
		return false
	}

	// Key does not exist, add it
	execPath, err := os.Executable()
	if err != nil {
		return false
	}
	dataVal := fmt.Sprintf("\"%s\"", execPath)
	err = exec.Command("reg", "add", runKey, "/v", valName, "/t", "REG_SZ", "/d", dataVal, "/f").Run()
	return err == nil
}

type RECT struct {
	Left, Top, Right, Bottom int32
}

func setDarkWindowFrame(hwnd uintptr) {
	darkMode := int32(1)
	// Try standard DWMWA_USE_IMMERSIVE_DARK_MODE (Win10 20H1+ & Win11)
	procDwmSetAttr.Call(
		hwnd,
		uintptr(DWMWA_USE_IMMERSIVE_DARK_MODE),
		uintptr(unsafe.Pointer(&darkMode)),
		unsafe.Sizeof(darkMode),
	)
	// Try older Win10 build
	procDwmSetAttr.Call(
		hwnd,
		uintptr(DWMWA_USE_IMMERSIVE_DARK_MODE_BEFORE_20H1),
		uintptr(unsafe.Pointer(&darkMode)),
		unsafe.Sizeof(darkMode),
	)

	// Set dark caption color (COLORREF: 0x00111B21 WhatsApp Dark Header: RGB 17, 27, 33)
	captionColor := uint32(0x00211B11) // 0x00BBGGRR
	procDwmSetAttr.Call(
		hwnd,
		uintptr(DWMWA_CAPTION_COLOR),
		uintptr(unsafe.Pointer(&captionColor)),
		unsafe.Sizeof(captionColor),
	)

	// Set white caption text (RGB 255, 255, 255)
	textColor := uint32(0x00FFFFFF)
	procDwmSetAttr.Call(
		hwnd,
		uintptr(DWMWA_TEXT_COLOR),
		uintptr(unsafe.Pointer(&textColor)),
		unsafe.Sizeof(textColor),
	)
}

func checkSingleInstance() (uintptr, bool) {
	namePtr, _ := syscall.UTF16PtrFromString(mutexName)
	handle, _, err := procCreateMutex.Call(0, 1, uintptr(unsafe.Pointer(namePtr)))
	if err == windows.ERROR_ALREADY_EXISTS {
		titlePtr, _ := syscall.UTF16PtrFromString(windowTitle)
		hwnd, _, _ := procFindWindow.Call(0, uintptr(unsafe.Pointer(titlePtr)))
		if hwnd != 0 {
			procShowNormal.Call(hwnd, 9) // SW_RESTORE
			procSetFgWindow.Call(hwnd)
		}
		return handle, false
	}
	return handle, true
}

func getUserDataDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = os.Getenv("APPDATA")
		if configDir == "" {
			configDir = "."
		}
	}
	dir := filepath.Join(configDir, "WhatsAppDesktopLight", "UserData")
	_ = os.MkdirAll(dir, 0755)
	return dir
}

func showNativeNotification(title, message, iconPath string) {
	notification := toast.Notification{
		AppID:   "WhatsApp Desktop",
		Title:   title,
		Message: message,
		Icon:    iconPath,
	}
	_ = notification.Push()
}

func configureWindow(hwnd uintptr) {
	setDarkWindowFrame(hwnd)

	// Ensure sizing border and maximize/minimize buttons are enabled
	gwlStyle := uintptr(GWL_STYLE)
	style, _, _ := procGetWindowLong.Call(hwnd, gwlStyle)
	style |= WS_THICKFRAME | WS_MAXIMIZEBOX | WS_MINIMIZEBOX
	procSetWindowLong.Call(hwnd, gwlStyle, style)
	procSetWindowPos.Call(hwnd, 0, 0, 0, 0, 0, uintptr(SWP_NOMOVE|SWP_NOSIZE|SWP_NOZORDER|SWP_FRAMECHANGED))
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

func saveWindowState(dir string, hwnd uintptr) {
	var r RECT
	ret, _, _ := procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&r)))
	if ret == 0 {
		return
	}
	w := r.Right - r.Left
	h := r.Bottom - r.Top
	if w >= 450 && h >= 320 {
		state := WindowState{
			X:      float64(r.Left),
			Y:      float64(r.Top),
			Width:  float64(w),
			Height: float64(h),
		}
		data, err := json.MarshalIndent(state, "", "  ")
		if err == nil {
			_ = os.WriteFile(filepath.Join(dir, "window_state.json"), data, 0644)
		}
	}
}

func runApp() {
	cleanupOldWindowsBinary()
	_, isSingle := checkSingleInstance()
	if !isSingle {
		os.Exit(0)
	}
	userDataDir := getUserDataDir()
	executablePath, _ := os.Executable()
	iconFullPath := filepath.Join(filepath.Dir(executablePath), "icon.ico")

	opts := webview2.WebViewOptions{
		Window:    nil,
		Debug:     false,
		DataPath:  userDataDir,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  windowTitle,
			Width:  windowWidth,
			Height: windowHeight,
			IconId: 2,
			Center: true,
		},
	}

	w := webview2.NewWithOptions(opts)
	if w == nil {
		log.Fatalln("Gagal inisialisasi WebView2")
	}
	defer w.Destroy()

	hwnd := uintptr(w.Window())
	configureWindow(hwnd)

	w.SetTitle(windowTitle)
	w.SetSize(450, 320, webview2.HintMin)
	w.SetSize(windowWidth, windowHeight, webview2.HintNone)

	if state := loadWindowState(userDataDir); state != nil {
		procMoveWindow.Call(hwnd, uintptr(int32(state.X)), uintptr(int32(state.Y)), uintptr(int32(state.Width)), uintptr(int32(state.Height)), 1)
	}

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			saveWindowState(userDataDir, hwnd)
		}
	}()

	// Bind native notification bridge
	_ = w.Bind("sendNativeNotification", func(title, body string) {
		go showNativeNotification(title, body, iconFullPath)
	})

	// Bind external link handler to open links in default Windows browser
	_ = w.Bind("openExternalLink", func(rawURL string) {
		if strings.HasPrefix(rawURL, "http://") || strings.HasPrefix(rawURL, "https://") {
			go func() {
				_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", rawURL).Start()
			}()
		}
	})

	// Bind Always on Top toggle
	_ = w.Bind("toggleAlwaysOnTopNative", func() bool {
		return toggleAlwaysOnTop(hwnd)
	})

	// Bind Auto-Start toggle
	_ = w.Bind("toggleAutoStartNative", func() bool {
		return toggleAutoStartWindows()
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

	w.Init(getInitScript(userAgent))
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

	defer saveWindowState(userDataDir, hwnd)
	w.Run()
}
