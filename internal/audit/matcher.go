package audit

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/git-pkgs/gitignore"
)

type pathMatcher struct {
	include []string
	ignore  *gitignore.Matcher
}

func newPathMatcher(root string, config Config) (*pathMatcher, error) {
	matcher := gitignore.NewFromDirectory(root)

	if len(config.Exclude) > 0 {
		matcher.AddPatterns([]byte(strings.Join(config.Exclude, "\n")), "")
	}

	if matcherErrs := matcher.Errors(); len(matcherErrs) > 0 {
		return nil, fmt.Errorf("load ignore patterns: %w", matcherErrs[0])
	}

	return &pathMatcher{
		include: append([]string(nil), config.Include...),
		ignore:  matcher,
	}, nil
}

func (m *pathMatcher) shouldSkip(relPath string, isDir bool) bool {
	if relPath == "" {
		return false
	}

	return m.ignore.MatchPath(filepath.ToSlash(relPath), isDir)
}

func (m *pathMatcher) shouldAudit(relPath string) (bool, error) {
	slashPath := filepath.ToSlash(relPath)

	for _, pattern := range m.include {
		matched, err := doublestar.Match(pattern, slashPath)
		if err != nil {
			return false, fmt.Errorf("match include pattern %q: %w", pattern, err)
		}

		if matched {
			return true, nil
		}
	}

	return false, nil
}
