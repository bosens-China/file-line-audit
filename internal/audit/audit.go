package audit

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Options struct {
	Args []string
	Cwd  string
}

type FileStat struct {
	Path  string
	Lines int
}

func Run(options Options) (string, error) {
	cwd := options.Cwd
	if cwd == "" {
		cwd = "."
	}

	configPath := ResolveConfigPath(options.Args)
	config, err := LoadConfig(configPath, cwd)
	if err != nil {
		return "", err
	}

	stats, err := Audit(cwd, config)
	if err != nil {
		return "", err
	}

	return BuildReport(config, stats), nil
}

func Audit(root string, config Config) ([]FileStat, error) {
	matcher, err := newPathMatcher(root, config)
	if err != nil {
		return nil, err
	}

	stats := make([]FileStat, 0)
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if path == root {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		relPath = filepath.ToSlash(relPath)
		if matcher.shouldSkip(relPath, entry.IsDir()) {
			if entry.IsDir() {
				return fs.SkipDir
			}

			return nil
		}

		if entry.IsDir() {
			return nil
		}

		shouldAudit, err := matcher.shouldAudit(relPath)
		if err != nil {
			return err
		}

		if !shouldAudit {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file %q: %w", relPath, err)
		}

		if isBinary(content) {
			return nil
		}

		lines := countLines(content)
		if lines < config.Threshold {
			return nil
		}

		stats = append(stats, FileStat{
			Path:  relPath,
			Lines: lines,
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Lines == stats[j].Lines {
			return stats[i].Path < stats[j].Path
		}

		return stats[i].Lines > stats[j].Lines
	})

	return stats, nil
}

func BuildReport(config Config, stats []FileStat) string {
	lines := []string{
		"# File Line Audit",
		"",
		fmt.Sprintf("## Files Over Threshold (>= %d lines)", config.Threshold),
		"",
	}

	if len(stats) == 0 {
		lines = append(lines, "- none")
		return strings.Join(lines, "\n")
	}

	for _, stat := range stats {
		lines = append(lines, fmt.Sprintf("- %s %d", stat.Path, stat.Lines))
	}

	return strings.Join(lines, "\n")
}

func isBinary(content []byte) bool {
	if len(content) == 0 {
		return false
	}

	sampleSize := len(content)
	if sampleSize > 8000 {
		sampleSize = 8000
	}

	suspicious := 0
	for _, value := range content[:sampleSize] {
		if value == 0 {
			return true
		}

		isControl := (value >= 1 && value <= 8) || value == 11 || value == 12 || (value >= 14 && value <= 31)
		if isControl {
			suspicious++
		}
	}

	return float64(suspicious)/float64(sampleSize) > 0.3
}

func countLines(content []byte) int {
	return bytes.Count(content, []byte{'\n'})
}
