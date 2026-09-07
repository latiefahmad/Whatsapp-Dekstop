//go:build darwin

package main

import (
	"os/exec"
	"strings"
)

func chooseFolderDialog() (string, error) {
	script := `set selFolder to choose folder with prompt "Pilih Folder Penyimpanan File WhatsApp" default location (path to downloads folder)
return POSIX path of selFolder`
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return "", err
	}
	res := strings.TrimSpace(string(out))
	res = strings.TrimSuffix(res, "/")
	return res, nil
}
