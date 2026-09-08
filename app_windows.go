//go:build windows

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/go-toast/toast"
	"github.com/jchv/go-webview2"
	"golang.org/x/sys/windows"
)

var (
	kernel32                     = windows.NewLazySystemDLL("kernel32.dll")
	user32                       = windows.NewLazySystemDLL("user32.dll")
	dwmapi                       = windows.NewLazySystemDLL("dwmapi.dll")
	procCreateMutex              = kernel32.NewProc("CreateMutexW")
	procFindWindow               = user32.NewProc("FindWindowW")
	procSetFgWindow              = user32.NewProc("SetForegroundWindow")
	procShowNormal               = user32.NewProc("ShowWindow")
	procDwmSetAttr               = dwmapi.NewProc("DwmSetWindowAttribute")
	procGetWindowLong            = user32.NewProc("GetWindowLongW")
	procSetWindowLong            = user32.NewProc("SetWindowLongW")
	procSetWindowPos             = user32.NewProc("SetWindowPos")
	procGetWindowRect            = user32.NewProc("GetWindowRect")
	procMoveWindow               = user32.NewProc("MoveWindow")
	procCreateJobObject          = kernel32.NewProc("CreateJobObjectW")
	procSetInformationJobObject  = kernel32.NewProc("SetInformationJobObject")
	procAssignProcessToJobObject = kernel32.NewProc("AssignProcessToJobObject")
	procSetProcessWorkingSetSize = kernel32.NewProc("SetProcessWorkingSetSize")
	procGetCurrentProcess        = kernel32.NewProc("GetCurrentProcess")

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
	valName := "WhatsAppDesk"

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

const (
	JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE = 0x00002000
	JobObjectExtendedLimitInformation  = 9
)

type IO_COUNTERS struct {
	ReadOperationCount  uint64
	WriteOperationCount uint64
	OtherOperationCount uint64
	ReadTransferCount   uint64
	WriteTransferCount  uint64
	OtherTransferCount  uint64
}

type JOBOBJECT_BASIC_LIMIT_INFORMATION struct {
	PerProcessUserTimeLimit int64
	PerJobUserTimeLimit     int64
	LimitFlags              uint32
	MinimumWorkingSetSize   uintptr
	MaximumWorkingSetSize   uintptr
	ActiveProcessLimit      uint32
	Affinity                uintptr
	PriorityClass           uint32
	SchedulingClass         uint32
}

type JOBOBJECT_EXTENDED_LIMIT_INFORMATION struct {
	BasicLimitInformation JOBOBJECT_BASIC_LIMIT_INFORMATION
	IoInfo                IO_COUNTERS
	ProcessMemoryLimit    uintptr
	JobMemoryLimit        uintptr
	PeakProcessMemoryUsed uintptr
	PeakJobMemoryUsed     uintptr
}

func initWindowsProcessProtection() {
	// 1. Assign process to Job Object with KILL_ON_JOB_CLOSE so Windows kernel
	// automatically kills all child msedgewebview2.exe processes on exit or crash.
	// This permanently prevents orphaned "WebView2 Manager" processes in RAM.
	job, _, _ := procCreateJobObject.Call(0, 0)
	if job != 0 {
		var info JOBOBJECT_EXTENDED_LIMIT_INFORMATION
		info.BasicLimitInformation.LimitFlags = JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
		procSetInformationJobObject.Call(
			job,
			uintptr(JobObjectExtendedLimitInformation),
			uintptr(unsafe.Pointer(&info)),
			uintptr(unsafe.Sizeof(info)),
		)
		curProc, _, _ := procGetCurrentProcess.Call()
		procAssignProcessToJobObject.Call(job, curProc)
	}

	// 2. Configure WebView2 / Chromium engine arguments for aggressive memory & cache limits:
	// - Limit V8 JS heap to 256MB and optimize for memory footprint
	// - Cap disk cache to 32MB and media cache to 16MB
	// - Disable background telemetry, component updates (which caused 8.8 MB/s disk I/O), and caching
	browserArgs := []string{
		"--disable-features=Translate,OptimizationHints,MediaRouter,CalculateNativeWinOcclusion,BackForwardCache,InterestFeedContentSuggestions",
		"--disk-cache-size=33554432",                                  // Cap disk cache to 32MB
		"--media-cache-size=16777216",                                 // Cap media cache to 16MB
		"--js-flags=\"--max-old-space-size=256 --optimize_for_size --expose-gc\"", // Cap V8 JS heap to 256MB & expose GC
		"--disable-gpu-shader-disk-cache",
		"--disable-background-networking",
		"--disable-component-update",
		"--disable-domain-reliability",
		"--disable-sync",
		"--renderer-process-limit=1",
		"--disable-site-isolation-trials",
		"--disable-speech-api",
	}
	_ = os.Setenv("WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS", strings.Join(browserArgs, " "))

	// 3. Clean up any leftover orphaned WebView2 instances from previous sessions
	go func() {
		psCmd := `Get-CimInstance Win32_Process -Filter "Name = 'msedgewebview2.exe'" | Where-Object { $_.CommandLine -like "*WhatsAppDesk*" } | ForEach-Object { Stop-Process -Id $_.ProcessId -Force }`
		_ = exec.Command("powershell", "-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden", "-Command", psCmd).Run()
	}()
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
	dir := filepath.Join(configDir, "WhatsAppDesk", "UserData")
	_ = os.MkdirAll(dir, 0755)
	return dir
}

func showNativeNotification(title, message, iconPath string) {
	notification := toast.Notification{
		AppID:   "WhatsApp Desk",
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

var lastSavedStateWin *WindowState

func saveWindowState(dir string, hwnd uintptr) {
	var r RECT
	ret, _, _ := procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&r)))
	if ret == 0 {
		return
	}
	w := r.Right - r.Left
	h := r.Bottom - r.Top
	if w >= 450 && h >= 320 {
		if lastSavedStateWin != nil &&
			lastSavedStateWin.X == float64(r.Left) &&
			lastSavedStateWin.Y == float64(r.Top) &&
			lastSavedStateWin.Width == float64(w) &&
			lastSavedStateWin.Height == float64(h) {
			return // Position and size unchanged, avoid disk write
		}
		state := WindowState{
			X:      float64(r.Left),
			Y:      float64(r.Top),
			Width:  float64(w),
			Height: float64(h),
		}
		data, err := json.MarshalIndent(state, "", "  ")
		if err == nil {
			_ = os.WriteFile(filepath.Join(dir, "window_state.json"), data, 0644)
			lastSavedStateWin = &state
		}
	}
}

func runApp() {
	cleanupOldWindowsBinary()
	initWindowsProcessProtection()
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
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			saveWindowState(userDataDir, hwnd)
		}
	}()

	// Periodic Working Set & Go runtime memory cleanup (every 60s)
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		curProc, _, _ := procGetCurrentProcess.Call()
		for range ticker.C {
			debug.FreeOSMemory()
			procSetProcessWorkingSetSize.Call(curProc, ^uintptr(0), ^uintptr(0))
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
