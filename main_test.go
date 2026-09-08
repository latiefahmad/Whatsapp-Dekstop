package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestPDFPreviewOwnsBlobURLLifecycle(t *testing.T) {
	script := getInitScript("test-agent")

	checks := []string{
		"origCreateObjectURL(previewBlob)",
		"URL.revokeObjectURL(ownedBlobUrl)",
		"isRecentPDFIntent()",
		"blob.type === 'application/octet-stream'",
		"findDocumentDownloadControl(el)",
		"e.stopImmediatePropagation()",
		"extractDocumentName(el)",
	}
	for _, want := range checks {
		if !strings.Contains(script, want) {
			t.Errorf("PDF preview script is missing %q", want)
		}
	}
}

func TestHiddenWindowRequestsNativeMemoryRelease(t *testing.T) {
	script := getInitScript("test-agent")
	for _, want := range []string{"visibilitychange", "releaseMemoryNative", "lastClickedDocName = ''"} {
		if !strings.Contains(script, want) {
			t.Errorf("memory lifecycle script is missing %q", want)
		}
	}
}

func TestSettingsControlsRemainWired(t *testing.T) {
	script := getInitScript("test-agent")
	ids := []string{
		"wa-theme-btn-dark", "wa-theme-btn-light", "wa-theme-btn-system",
		"wa-action-toggle-priv", "wa-action-toggle-pin", "wa-action-toggle-mute", "wa-action-toggle-auto",
		"wa-btn-change-folder", "wa-btn-open-folder", "wa-btn-reset-folder",
		"wa-btn-check-updates-modal", "wa-btn-reload-modal", "wa-btn-hardref-modal", "wa-btn-onboard-modal",
	}
	for _, id := range ids {
		if strings.Count(script, `id="`+id+`"`) != 1 {
			t.Errorf("control %s must be rendered exactly once", id)
		}
		if !strings.Contains(script, "getElementById('"+id+"')") {
			t.Errorf("control %s has no event or state binding", id)
		}
	}
}

func TestOnboardingIsQuietAndReplayable(t *testing.T) {
	script := getOnboardingScript()
	for _, unwanted := range []string{"radial-gradient", "backdrop-filter", "wa-feat-card", "Hemat RAM ~90%"} {
		if strings.Contains(script, unwanted) {
			t.Errorf("onboarding still contains noisy pattern %q", unwanted)
		}
	}
	for _, want := range []string{"window.showOnboardingModal", "prefers-reduced-motion", "Lewati panduan"} {
		if !strings.Contains(script, want) {
			t.Errorf("onboarding is missing %q", want)
		}
	}
}

func TestInjectedJavaScriptParses(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node is not available")
	}
	cmd := exec.Command("node", "--check", "-")
	cmd.Stdin = strings.NewReader(getInitScript("test-agent"))
	if output, err := cmd.CombinedOutput(); err != nil {
		if strings.Contains(err.Error(), "operation not permitted") || strings.Contains(err.Error(), "permission denied") {
			t.Skipf("skipping node execution in sandboxed environment: %v", err)
		}
		t.Fatalf("injected JavaScript does not parse: %v\n%s", err, output)
	}
}

func TestDarwinMenuBridgeUsesStableAppWindow(t *testing.T) {
	source, err := os.ReadFile("app_darwin.go")
	if err != nil {
		t.Fatal(err)
	}
	script := string(source)
	if !strings.Contains(script, "static NSWindow* appWindow(void)") {
		t.Fatal("native menu must resolve the app's stored window")
	}

	start := strings.Index(script, "@implementation MenuBridge")
	if start < 0 {
		t.Fatal("MenuBridge implementation not found")
	}
	endOffset := strings.Index(script[start:], "@end")
	if endOffset < 0 {
		t.Fatal("MenuBridge implementation has no end")
	}
	bridge := script[start : start+endOffset]
	if strings.Contains(bridge, "[NSApp keyWindow] ?: [NSApp mainWindow]") {
		t.Fatal("menu actions must not depend on the transient key/main window")
	}
	for _, action := range []string{
		"menuSettings:", "menuCheckUpdates:", "menuOpenDownloads:",
		"menuTogglePrivacy:", "menuToggleAlwaysOnTop:", "menuToggleMuteAudio:",
		"menuReloadChat:", "menuHardRefresh:", "menuShowApp:",
		"menuSetThemeDark:", "menuSetThemeLight:", "menuSetThemeSystem:",
	} {
		if !strings.Contains(bridge, action) {
			t.Errorf("native menu action %s is missing", action)
		}
	}
}

func TestDarwinPDFUsesNativePDFKitPreview(t *testing.T) {
	source, err := os.ReadFile("app_darwin.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"-framework PDFKit", "#import <PDFKit/PDFKit.h>",
		"showNativePDFPreview", "showPDFPreviewNative",
	} {
		if !strings.Contains(string(source), want) {
			t.Errorf("native PDF preview is missing %q", want)
		}
	}

	script := getInitScript("test-agent")
	if !strings.Contains(script, "window.showPDFPreviewNative(savedPath)") {
		t.Fatal("PDF modal must hand saved PDFs to the native renderer on macOS")
	}
}

func TestWindowsProcessProtectionPermitsChildBreakaway(t *testing.T) {
	source, err := os.ReadFile("app_windows.go")
	if err != nil {
		t.Fatal(err)
	}
	content := string(source)

	// Must permit child breakaway so Chromium sandbox job objects don't fail
	if !strings.Contains(content, "JOB_OBJECT_LIMIT_SILENT_BREAKAWAY_OK") {
		t.Error("Windows Job Object must include JOB_OBJECT_LIMIT_SILENT_BREAKAWAY_OK")
	}

	// Must not kill own child processes at startup
	if strings.Contains(content, "Stop-Process") {
		t.Error("Windows startup must not invoke PowerShell Stop-Process which kills active webview instances")
	}

	// Must not pass dangerous flags that trigger black screen on Windows 10/11
	for _, dangerous := range []string{"--disable-gpu-shader-disk-cache", "CalculateNativeWinOcclusion", `--js-flags="`} {
		if strings.Contains(content, dangerous) {
			t.Errorf("Windows browser args contain dangerous flag %q known to cause black screen", dangerous)
		}
	}
}
