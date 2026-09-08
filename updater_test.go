package main

import "testing"

func TestLinuxUpdaterSelectsPortableArchive(t *testing.T) {
	release := &GitHubRelease{Assets: []GitHubAsset{
		{Name: "WhatsApp-Desk-Linux-amd64.deb", BrowserDownloadURL: "https://example.test/app.deb"},
		{Name: "WhatsApp-Desk-Linux-x64.tar.gz", BrowserDownloadURL: "https://example.test/app.tar.gz"},
	}}
	asset := findAssetForOS(release, "linux")
	if asset == nil || asset.Name != "WhatsApp-Desk-Linux-x64.tar.gz" {
		t.Fatalf("Linux self-update must select tar.gz, got %#v", asset)
	}
	if got := updateDownloadExtension(asset.BrowserDownloadURL); got != ".tar.gz" {
		t.Fatalf("download extension = %q, want .tar.gz", got)
	}
}

func TestIsNewerVersion(t *testing.T) {
	cases := []struct {
		current string
		latest  string
		want    bool
	}{
		{"1.4.0", "1.4.0", false},
		{"1.4.0", "v1.4.0", false},
		{"1.4.0", "1.4.1", true},
		{"1.4.0", "v1.5.0", true},
		{"1.5.0", "v1.5.1", true},
		{"1.5.1", "v1.5.2", true},
		{"1.5.2", "v1.5.3", true},
		{"1.5.3", "1.5.3", false},
		{"1.5.2", "1.5.2", false},
		{"1.5.2", "1.5.1", false},
		{"1.5.1", "v1.5.1", false},
		{"1.5.1", "1.5.0", false},
		{"1.4.0", "2.0.0", true},
		{"1.4.0", "1.3.9", false},
		{"1.4.0", "1.4.0-beta", false},
	}

	for _, c := range cases {
		got := isNewerVersion(c.current, c.latest)
		if got != c.want {
			t.Errorf("isNewerVersion(%q, %q) = %v; want %v", c.current, c.latest, got, c.want)
		}
	}
}

func TestCheckForUpdateLive(t *testing.T) {
	// Version 1.0.0 should always detect an update on live repo
	infoOld, err := checkForUpdate("1.0.0")
	if err != nil {
		t.Skipf("skipping live network test: %v", err)
		return
	}
	if !infoOld.Available {
		t.Errorf("expected update for 1.0.0, got false")
	}
	if infoOld.LatestVersion == "" {
		t.Errorf("expected non-empty latest version")
	}
	if infoOld.DownloadURL == "" {
		t.Errorf("expected non-empty download URL")
	}
}
