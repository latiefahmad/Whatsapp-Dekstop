//go:build darwin

package main

/*
#cgo darwin CFLAGS: -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa -framework WebKit

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>
#include <stdlib.h>

@interface WhatsAppWindowDelegate : NSObject <NSWindowDelegate>
@end

@implementation WhatsAppWindowDelegate
- (BOOL)windowShouldClose:(NSWindow *)sender {
    [sender orderOut:nil];
    return NO;
}
@end

@interface WhatsAppAppDelegate : NSObject <NSApplicationDelegate>
@property (assign) NSWindow *window;
@end

@implementation WhatsAppAppDelegate
- (BOOL)applicationShouldHandleReopen:(NSApplication *)sender hasVisibleWindows:(BOOL)flag {
    if (self.window) {
        [self.window makeKeyAndOrderFront:nil];
        [NSApp activateIgnoringOtherApps:YES];
    }
    return YES;
}
@end

@interface WhatsAppUIDelegate : NSObject <WKUIDelegate>
@end

@implementation WhatsAppUIDelegate
- (void)webView:(WKWebView *)webView requestMediaCapturePermissionForOrigin:(WKSecurityOrigin *)origin initiatedByFrame:(WKFrameInfo *)frame type:(WKMediaCaptureType)type decisionHandler:(void (^)(WKPermissionDecision decision))decisionHandler {
    decisionHandler(WKPermissionDecisionGrant);
}
@end

static WhatsAppAppDelegate* g_appDelegate = nil;
static WhatsAppWindowDelegate* g_windowDelegate = nil;
static WhatsAppUIDelegate* g_uiDelegate = nil;

static void setWKWebViewUserAgentAndMedia(void* nsWindowPtr, const char* uaStr) {
    @autoreleasepool {
        NSWindow* win = (__bridge NSWindow*)nsWindowPtr;
        NSView* contentView = [win contentView];
        if ([contentView isKindOfClass:[WKWebView class]]) {
            WKWebView* wv = (WKWebView*)contentView;
            [wv setAutoresizingMask:(NSViewWidthSizable | NSViewHeightSizable)];
            NSString* ua = [NSString stringWithUTF8String:uaStr];
            [wv setCustomUserAgent:ua];

            g_uiDelegate = [[WhatsAppUIDelegate alloc] init];
            [wv setUIDelegate:g_uiDelegate];
        }
    }
}

static void configureWindowBehavior(void* nsWindowPtr) {
    @autoreleasepool {
        NSWindow* win = (__bridge NSWindow*)nsWindowPtr;

        // Ensure window is fully resizable, minimizable, and supports fullscreen
        NSWindowStyleMask mask = [win styleMask];
        mask |= (NSWindowStyleMaskResizable | NSWindowStyleMaskTitled | NSWindowStyleMaskClosable | NSWindowStyleMaskMiniaturizable);
        [win setStyleMask:mask];

        [win setCollectionBehavior:(NSWindowCollectionBehaviorFullScreenPrimary | NSWindowCollectionBehaviorDefault)];
        [win setShowsResizeIndicator:YES];

        // Minimum bounds: allow shrinking down dynamically to compact window
        [win setMinSize:NSMakeSize(450, 320)];
        [win setContentMinSize:NSMakeSize(450, 320)];
        [win setMaxSize:NSMakeSize(FLT_MAX, FLT_MAX)];
        [win setContentMaxSize:NSMakeSize(FLT_MAX, FLT_MAX)];

        win.appearance = [NSAppearance appearanceNamed:NSAppearanceNameDarkAqua];
        win.titlebarAppearsTransparent = YES;

        NSView* contentView = [win contentView];
        if (contentView) {
            [contentView setAutoresizingMask:(NSViewWidthSizable | NSViewHeightSizable)];
        }

        g_windowDelegate = [[WhatsAppWindowDelegate alloc] init];
        [win setDelegate:g_windowDelegate];

        g_appDelegate = [[WhatsAppAppDelegate alloc] init];
        g_appDelegate.window = win;
        [NSApp setDelegate:g_appDelegate];
    }
}

static BOOL g_isAlwaysOnTop = NO;

static BOOL toggleAlwaysOnTop(void* nsWindowPtr) {
    @autoreleasepool {
        NSWindow* win = (__bridge NSWindow*)nsWindowPtr;
        g_isAlwaysOnTop = !g_isAlwaysOnTop;
        if (g_isAlwaysOnTop) {
            [win setLevel:NSFloatingWindowLevel];
        } else {
            [win setLevel:NSNormalWindowLevel];
        }
        return g_isAlwaysOnTop;
    }
}

static void setDockBadge(const char* labelStr) {
    @autoreleasepool {
        NSString* label = (labelStr && strlen(labelStr) > 0) ? [NSString stringWithUTF8String:labelStr] : nil;
        [[NSApp dockTile] setBadgeLabel:label];
    }
}

static void getWindowFrame(void* nsWindowPtr, double* x, double* y, double* w, double* h) {
    @autoreleasepool {
        NSWindow* win = (__bridge NSWindow*)nsWindowPtr;
        NSRect frame = [win frame];
        *x = (double)frame.origin.x;
        *y = (double)frame.origin.y;
        *w = (double)frame.size.width;
        *h = (double)frame.size.height;
    }
}

static void setWindowFrame(void* nsWindowPtr, double x, double y, double w, double h) {
    @autoreleasepool {
        NSWindow* win = (__bridge NSWindow*)nsWindowPtr;
        NSRect rect = NSMakeRect((CGFloat)x, (CGFloat)y, (CGFloat)w, (CGFloat)h);
        [win setFrame:rect display:YES animate:NO];
    }
}

static void setupMacOSMenuBar(void) {
    @autoreleasepool {
        NSMenu* mainMenu = [[NSMenu alloc] init];

        // App Menu
        NSMenuItem* appMenuItem = [[NSMenuItem alloc] init];
        NSMenu* appMenu = [[NSMenu alloc] initWithTitle:@"WhatsApp"];
        [appMenu addItemWithTitle:@"About WhatsApp Desktop" action:@selector(orderFrontStandardAboutPanel:) keyEquivalent:@""];
        [appMenu addItem:[NSMenuItem separatorItem]];
        [appMenu addItemWithTitle:@"Hide WhatsApp Desktop" action:@selector(hide:) keyEquivalent:@"h"];
        NSMenuItem* hideOthers = [appMenu addItemWithTitle:@"Hide Others" action:@selector(hideOtherApplications:) keyEquivalent:@"h"];
        [hideOthers setKeyEquivalentModifierMask:(NSEventModifierFlagOption | NSEventModifierFlagCommand)];
        [appMenu addItemWithTitle:@"Show All" action:@selector(unhideAllApplications:) keyEquivalent:@""];
        [appMenu addItem:[NSMenuItem separatorItem]];
        [appMenu addItemWithTitle:@"Quit WhatsApp Desktop" action:@selector(terminate:) keyEquivalent:@"q"];
        [appMenuItem setSubmenu:appMenu];
        [mainMenu addItem:appMenuItem];

        // Edit Menu (Essential for Cmd+C, Cmd+V, Cmd+X, Cmd+A, Cmd+Z)
        NSMenuItem* editMenuItem = [[NSMenuItem alloc] init];
        NSMenu* editMenu = [[NSMenu alloc] initWithTitle:@"Edit"];
        [editMenu addItemWithTitle:@"Undo" action:@selector(undo:) keyEquivalent:@"z"];
        NSMenuItem* redo = [editMenu addItemWithTitle:@"Redo" action:@selector(redo:) keyEquivalent:@"Z"];
        [redo setKeyEquivalentModifierMask:(NSEventModifierFlagShift | NSEventModifierFlagCommand)];
        [editMenu addItem:[NSMenuItem separatorItem]];
        [editMenu addItemWithTitle:@"Cut" action:@selector(cut:) keyEquivalent:@"x"];
        [editMenu addItemWithTitle:@"Copy" action:@selector(copy:) keyEquivalent:@"c"];
        [editMenu addItemWithTitle:@"Paste" action:@selector(paste:) keyEquivalent:@"v"];
        [editMenu addItemWithTitle:@"Select All" action:@selector(selectAll:) keyEquivalent:@"a"];
        [editMenuItem setSubmenu:editMenu];
        [mainMenu addItem:editMenuItem];

        // Window Menu
        NSMenuItem* windowMenuItem = [[NSMenuItem alloc] init];
        NSMenu* windowMenu = [[NSMenu alloc] initWithTitle:@"Window"];
        [windowMenu addItemWithTitle:@"Minimize" action:@selector(performMiniaturize:) keyEquivalent:@"m"];
        [windowMenu addItemWithTitle:@"Zoom" action:@selector(performZoom:) keyEquivalent:@""];
        [windowMenu addItem:[NSMenuItem separatorItem]];
        [windowMenu addItemWithTitle:@"Close Window" action:@selector(performClose:) keyEquivalent:@"w"];
        [windowMenuItem setSubmenu:windowMenu];
        [mainMenu addItem:windowMenuItem];

        [NSApp setMainMenu:mainMenu];
        [NSApp setWindowsMenu:windowMenu];
    }
}
*/
import "C"
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

	"github.com/webview/webview_go"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"

func getUserDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	dir := filepath.Join(home, "Library", "Application Support", "WhatsAppDesktopLight", "UserData")
	_ = os.MkdirAll(dir, 0755)
	return dir
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

func saveWindowState(dir string, win unsafe.Pointer) {
	if win == nil {
		return
	}
	var x, y, w, h C.double
	C.getWindowFrame(win, &x, &y, &w, &h)
	if w >= 450 && h >= 320 {
		state := WindowState{
			X:      float64(x),
			Y:      float64(y),
			Width:  float64(w),
			Height: float64(h),
		}
		data, err := json.MarshalIndent(state, "", "  ")
		if err == nil {
			_ = os.WriteFile(filepath.Join(dir, "window_state.json"), data, 0644)
		}
	}
}

func checkSingleInstance() (*os.File, bool) {
	dir := getUserDataDir()
	lockPath := filepath.Join(dir, "whatsapp.lock")
	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, true
	}
	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		_ = exec.Command("osascript", "-e", `tell application "WhatsApp" to activate`).Run()
		return nil, false
	}
	return file, true
}

func getAppBundlePath() string {
	execPath, err := os.Executable()
	if err != nil {
		return "/Applications/WhatsApp Web.app"
	}
	if idx := strings.Index(execPath, ".app"); idx != -1 {
		return execPath[:idx+4]
	}
	if _, err := os.Stat("/Applications/WhatsApp Web.app"); err == nil {
		return "/Applications/WhatsApp Web.app"
	}
	return execPath
}

func getLaunchAgentPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, "Library", "LaunchAgents", "com.whatsapp.desktoplight.plist")
}

func toggleAutoStartMac() bool {
	plistPath := getLaunchAgentPath()
	if plistPath == "" {
		return false
	}
	if _, err := os.Stat(plistPath); err == nil {
		_ = os.Remove(plistPath)
		return false
	}

	appPath := getAppBundlePath()
	plistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.whatsapp.desktoplight</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/bin/open</string>
        <string>%s</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
</dict>
</plist>`, appPath)

	_ = os.MkdirAll(filepath.Dir(plistPath), 0755)
	if err := os.WriteFile(plistPath, []byte(plistContent), 0644); err != nil {
		return false
	}
	return true
}

func showNativeNotification(title, message string) {
	script := fmt.Sprintf(`display notification %q with title %q`, message, title)
	_ = exec.Command("osascript", "-e", script).Run()
}

func runApp() {
	lockFile, isSingle := checkSingleInstance()
	if !isSingle {
		os.Exit(0)
	}
	if lockFile != nil {
		defer lockFile.Close()
	}

	userDataDir := getUserDataDir()

	w := webview.New(false)
	if w == nil {
		log.Fatalln("Gagal inisialisasi WebKit WebView")
	}
	defer w.Destroy()

	// 1. Setup standard macOS menu bar (Cmd+C, Cmd+V, Cmd+X, Cmd+A, Cmd+Z, Cmd+Q)
	C.setupMacOSMenuBar()

	// 2. Configure window behavior: dark title bar, close-to-hide, and dock click reopen
	C.configureWindowBehavior(w.Window())

	// 3. Set native WebKit customUserAgent to Google Chrome & auto-grant media capture
	cua := C.CString(userAgent)
	C.setWKWebViewUserAgentAndMedia(w.Window(), cua)
	C.free(unsafe.Pointer(cua))

	w.SetTitle(windowTitle)

	// 4. Ensure window is initialized with HintNone (resizable), then restore saved state
	w.SetSize(windowWidth, windowHeight, webview.HintNone)
	if state := loadWindowState(userDataDir); state != nil {
		C.setWindowFrame(w.Window(), C.double(state.X), C.double(state.Y), C.double(state.Width), C.double(state.Height))
	}

	// 5. Save window state periodically
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			saveWindowState(userDataDir, w.Window())
		}
	}()

	// 6. Bind native notification bridge
	_ = w.Bind("sendNativeNotification", func(title, body string) {
		go showNativeNotification(title, body)
	})

	// 7. Bind external link handler to open links in macOS default browser
	_ = w.Bind("openExternalLink", func(rawURL string) {
		if strings.HasPrefix(rawURL, "http://") || strings.HasPrefix(rawURL, "https://") {
			_ = exec.Command("open", rawURL).Start()
		}
	})

	// 8. Bind dock badge unread counter
	_ = w.Bind("updateDockBadge", func(badge string) {
		cstr := C.CString(badge)
		defer C.free(unsafe.Pointer(cstr))
		C.setDockBadge(cstr)
	})

	// 9. Bind Always on Top toggle
	_ = w.Bind("toggleAlwaysOnTopNative", func() bool {
		return bool(C.toggleAlwaysOnTop(w.Window()))
	})

	// 10. Bind Auto-Start toggle
	_ = w.Bind("toggleAutoStartNative", func() bool {
		return toggleAutoStartMac()
	})

	// 11. Bind in-app auto updater
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

	// 13. Bind download and settings handlers
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

	w.Init(getInitScript(userAgent))
	w.Navigate(appURL)

	// 12. Check for updates in the background after startup
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

	defer saveWindowState(userDataDir, w.Window())
	w.Run()
}
