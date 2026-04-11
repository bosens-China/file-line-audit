package audit

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestResolveConfigPath(t *testing.T) {
	t.Parallel()

	if got := ResolveConfigPath([]string{"--config", "custom.json"}); got != "custom.json" {
		t.Fatalf("expected custom.json, got %q", got)
	}

	if got := ResolveConfigPath([]string{"-c", "alt.json"}); got != "alt.json" {
		t.Fatalf("expected alt.json, got %q", got)
	}

	if got := ResolveConfigPath(nil); got != DefaultConfigPath {
		t.Fatalf("expected %q, got %q", DefaultConfigPath, got)
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	config, err := LoadConfig(DefaultConfigPath, dir)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if config.Threshold != DefaultThreshold {
		t.Fatalf("expected threshold %d, got %d", DefaultThreshold, config.Threshold)
	}

	if !reflect.DeepEqual(config.Include, DefaultInclude) {
		t.Fatalf("unexpected default include patterns: %#v", config.Include)
	}

	if !reflect.DeepEqual(config.Exclude, DefaultExclude) {
		t.Fatalf("unexpected default exclude patterns: %#v", config.Exclude)
	}
}

func TestLoadConfigReadsCustomFields(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	configPath := filepath.Join(dir, DefaultConfigPath)
	content := []byte(`{
  "threshold": 128,
  "include": ["src/**/*.go", "apps/**/*.{ts,tsx}"],
  "exclude": ["dist/", "*.snap"]
}`)
	if err := os.WriteFile(configPath, content, 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	config, err := LoadConfig(DefaultConfigPath, dir)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if config.Threshold != 128 {
		t.Fatalf("expected threshold 128, got %d", config.Threshold)
	}

	wantInclude := []string{"src/**/*.go", "apps/**/*.{ts,tsx}"}
	if !reflect.DeepEqual(config.Include, wantInclude) {
		t.Fatalf("unexpected include patterns: %#v", config.Include)
	}

	wantExclude := []string{"dist/", "*.snap"}
	if !reflect.DeepEqual(config.Exclude, wantExclude) {
		t.Fatalf("unexpected exclude patterns: %#v", config.Exclude)
	}
}

func TestLoadConfigRejectsInvalidIncludePattern(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	configPath := filepath.Join(dir, DefaultConfigPath)
	content := []byte(`{
  "include": ["src/**["]
}`)
	if err := os.WriteFile(configPath, content, 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	if _, err := LoadConfig(DefaultConfigPath, dir); err == nil {
		t.Fatal("expected invalid include pattern error")
	}
}
