package audit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAuditRespectsGitIgnoreIncludeAndExclude(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, ".gitignore"), "apps/site/ignored.py\nnode_modules/\n")
	mustWriteFile(t, filepath.Join(root, "src", "keep.ts"), "1\n2\n3\n4\n5\n")
	mustWriteFile(t, filepath.Join(root, "src", "binary.ts"), string([]byte{0, 1, 2, 3}))
	mustWriteFile(t, filepath.Join(root, "src", "note.md"), "1\n2\n3\n4\n5\n6\n")
	mustWriteFile(t, filepath.Join(root, "packages", "pkg", ".gitignore"), "vendor/\n")
	mustWriteFile(t, filepath.Join(root, "packages", "pkg", "index.js"), "1\n2\n3\n4\n5\n6\n")
	mustWriteFile(t, filepath.Join(root, "packages", "pkg", "vendor", "generated.js"), "1\n2\n3\n4\n5\n6\n7\n")
	mustWriteFile(t, filepath.Join(root, "apps", "site", "page.py"), "1\n2\n3\n4\n5\n6\n7\n")
	mustWriteFile(t, filepath.Join(root, "apps", "site", "ignored.py"), "1\n2\n3\n4\n5\n6\n7\n8\n")
	mustWriteFile(t, filepath.Join(root, "apps", "site", "generated", "report.py"), "1\n2\n3\n4\n5\n6\n7\n8\n9\n")
	mustWriteFile(t, filepath.Join(root, "node_modules", "lib", "index.js"), "1\n2\n3\n4\n5\n6\n")

	config := Config{
		Threshold: 5,
		Include: []string{
			"src/**/*.{ts,py}",
			"packages/**/*.{js,ts}",
			"apps/**/*.{py,ts}",
		},
		Exclude: []string{
			"apps/site/generated/",
			"*.snap",
		},
	}

	stats, err := Audit(root, config)
	if err != nil {
		t.Fatalf("Audit returned error: %v", err)
	}

	report := BuildReport(config, stats)
	if !strings.Contains(report, "- apps/site/page.py 7") {
		t.Fatalf("report missing apps/site/page.py: %s", report)
	}

	if !strings.Contains(report, "- packages/pkg/index.js 6") {
		t.Fatalf("report missing packages/pkg/index.js: %s", report)
	}

	if !strings.Contains(report, "- src/keep.ts 5") {
		t.Fatalf("report missing src/keep.ts: %s", report)
	}

	for _, unwanted := range []string{
		"src/binary.ts",
		"src/note.md",
		"packages/pkg/vendor/generated.js",
		"apps/site/ignored.py",
		"apps/site/generated/report.py",
		"node_modules/lib/index.js",
	} {
		if strings.Contains(report, unwanted) {
			t.Fatalf("report should not contain %q: %s", unwanted, report)
		}
	}
}

func TestRunReadsConfigFile(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, DefaultConfigPath), `{
  "threshold": 3,
  "include": ["src/**/*.py"],
  "exclude": ["src/generated/"]
}`)
	mustWriteFile(t, filepath.Join(root, "src", "app.py"), "1\n2\n3\n4\n")
	mustWriteFile(t, filepath.Join(root, "src", "generated", "skip.py"), "1\n2\n3\n4\n5\n")

	report, err := Run(Options{
		Args: []string{"--config", DefaultConfigPath},
		Cwd:  root,
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if !strings.Contains(report, "- src/app.py 4") {
		t.Fatalf("report missing src/app.py: %s", report)
	}

	if strings.Contains(report, "src/generated/skip.py") {
		t.Fatalf("report should not contain excluded file: %s", report)
	}
}

func mustWriteFile(t *testing.T, path string, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
}
