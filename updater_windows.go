//go:build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func applyUpdateWindows(newExePath string) error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	oldExePath := execPath + ".old"
	_ = os.Remove(oldExePath)

	// In Windows, a running executable can be renamed!
	if err := os.Rename(execPath, oldExePath); err != nil {
		// Fallback to background script if direct rename is not permitted
		script := fmt.Sprintf(`timeout /t 1 /nobreak > NUL & move /y "%s" "%s" & start "" "%s"`,
			newExePath, execPath, execPath)
		cmd := exec.Command("cmd.exe", "/c", script)
		if err := cmd.Start(); err != nil {
			return err
		}
		os.Exit(0)
		return nil
	}

	// Move new executable into place
	if err := os.Rename(newExePath, execPath); err != nil {
		// Rollback rename if moving new file fails
		_ = os.Rename(oldExePath, execPath)
		return err
	}

	// Launch newly replaced executable
	cmd := exec.Command(execPath)
	if err := cmd.Start(); err != nil {
		return err
	}

	os.Exit(0)
	return nil
}

func applyUpdate(downloadedFile string) error {
	return applyUpdateWindows(downloadedFile)
}

func cleanupOldWindowsBinary() {
	execPath, err := os.Executable()
	if err == nil {
		oldPath := execPath + ".old"
		if _, err := os.Stat(oldPath); err == nil {
			_ = os.Remove(oldPath)
		}
	}
}
