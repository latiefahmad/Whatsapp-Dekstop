package main

import "testing"

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
	// 1. Current version 1.4.0 should not see an update if 1.4.0 is latest
	info, err := checkForUpdate("1.4.0")
	if err != nil {
		t.Skipf("skipping live network test: %v", err)
		return
	}
	if info.Available {
		t.Errorf("expected no update for 1.4.0, got available: %v (latest: %s)", info.Available, info.LatestVersion)
	}

	// 2. Old version 1.3.0 should detect 1.4.0
	infoOld, err := checkForUpdate("1.3.0")
	if err != nil {
		t.Fatalf("checkForUpdate failed: %v", err)
	}
	if !infoOld.Available {
		t.Errorf("expected update for 1.3.0, got false")
	}
	if infoOld.LatestVersion != "1.4.0" {
		t.Errorf("expected latest version 1.4.0, got %s", infoOld.LatestVersion)
	}
	if infoOld.DownloadURL == "" {
		t.Errorf("expected non-empty download URL")
	}
}
