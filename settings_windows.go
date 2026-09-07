//go:build windows

package main

import (
	"os/exec"
	"strings"
)

func chooseFolderDialog() (string, error) {
	script := `Add-Type -AssemblyName System.Windows.Forms; $f = New-Object System.Windows.Forms.FolderBrowserDialog; $f.Description = 'Pilih Folder Penyimpanan File WhatsApp'; if ($f.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) { Write-Output $f.SelectedPath }`
	out, err := exec.Command("powershell", "-NoProfile", "-Command", script).Output()
	if err != nil {
		return "", err
	}
	res := strings.TrimSpace(string(out))
	return res, nil
}
