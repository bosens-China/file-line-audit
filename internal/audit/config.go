package audit

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

const (
	DefaultConfigPath = ".line-audit.json"
	DefaultThreshold  = 400
)

var DefaultInclude = []string{
	"src/**/*.{js,ts,vue,jsx,tsx,mjs,cjs,cts,mts,py,go,rs,java,kt,kts,scala,rb,php,cs,swift,dart}",
	"app/**/*.{js,ts,vue,jsx,tsx,mjs,cjs,cts,mts,py,go,rs,java,kt,kts,scala,rb,php,cs,swift,dart}",
	"apps/**/*.{js,ts,vue,jsx,tsx,mjs,cjs,cts,mts,py,go,rs,java,kt,kts,scala,rb,php,cs,swift,dart}",
	"pkg/**/*.{js,ts,vue,jsx,tsx,mjs,cjs,cts,mts,py,go,rs,java,kt,kts,scala,rb,php,cs,swift,dart}",
	"packages/**/*.{js,ts,vue,jsx,tsx,mjs,cjs,cts,mts,py,go,rs,java,kt,kts,scala,rb,php,cs,swift,dart}",
	"lib/**/*.{js,ts,vue,jsx,tsx,mjs,cjs,cts,mts,py,go,rs,java,kt,kts,scala,rb,php,cs,swift,dart}",
	"libs/**/*.{js,ts,vue,jsx,tsx,mjs,cjs,cts,mts,py,go,rs,java,kt,kts,scala,rb,php,cs,swift,dart}",
	"cmd/**/*.{go,rs,py,js,ts}",
	"internal/**/*.{go,rs,py,js,ts}",
	"backend/**/*.{js,ts,vue,jsx,tsx,mjs,cjs,cts,mts,py,go,rs,java,kt,kts,scala,rb,php,cs,swift,dart}",
	"frontend/**/*.{js,ts,vue,jsx,tsx,mjs,cjs,cts,mts,py,go,rs,java,kt,kts,scala,rb,php,cs,swift,dart}",
	"client/**/*.{js,ts,vue,jsx,tsx,mjs,cjs,cts,mts,py,go,rs,java,kt,kts,scala,rb,php,cs,swift,dart}",
	"server/**/*.{js,ts,vue,jsx,tsx,mjs,cjs,cts,mts,py,go,rs,java,kt,kts,scala,rb,php,cs,swift,dart}",
	"service/**/*.{js,ts,vue,jsx,tsx,mjs,cjs,cts,mts,py,go,rs,java,kt,kts,scala,rb,php,cs,swift,dart}",
	"services/**/*.{js,ts,vue,jsx,tsx,mjs,cjs,cts,mts,py,go,rs,java,kt,kts,scala,rb,php,cs,swift,dart}",
	"api/**/*.{js,ts,vue,jsx,tsx,mjs,cjs,cts,mts,py,go,rs,java,kt,kts,scala,rb,php,cs,swift,dart}",
	"web/**/*.{js,ts,vue,jsx,tsx,mjs,cjs,cts,mts,py,go,rs,java,kt,kts,scala,rb,php,cs,swift,dart}",
}

var DefaultExclude = []string{
	"node_modules/",
	"dist/",
	"build/",
	"coverage/",
	".next/",
	".nuxt/",
	".output/",
	".svelte-kit/",
	"storybook-static/",
	"target/",
	"out/",
	"vendor/",
	"tmp/",
	"temp/",
}

type Config struct {
	Threshold int      `json:"threshold"`
	Include   []string `json:"include"`
	Exclude   []string `json:"exclude"`
}

type rawConfig struct {
	Threshold int      `json:"threshold"`
	Include   []string `json:"include"`
	Exclude   []string `json:"exclude"`
}

func ResolveConfigPath(args []string) string {
	for index := 0; index < len(args); index++ {
		switch args[index] {
		case "--config", "-c":
			if index+1 < len(args) {
				return args[index+1]
			}
		}
	}

	return DefaultConfigPath
}

func LoadConfig(configPath string, cwd string) (Config, error) {
	config := Config{
		Threshold: DefaultThreshold,
		Include:   append([]string(nil), DefaultInclude...),
		Exclude:   append([]string(nil), DefaultExclude...),
	}

	absPath := filepath.Join(cwd, configPath)
	content, err := os.ReadFile(absPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return config, nil
		}

		return Config{}, fmt.Errorf("read config %q: %w", configPath, err)
	}

	var raw rawConfig
	if err := json.Unmarshal(content, &raw); err != nil {
		return Config{}, fmt.Errorf("parse config %q: %w", configPath, err)
	}

	if raw.Threshold > 0 {
		config.Threshold = raw.Threshold
	}

	if len(raw.Include) > 0 {
		config.Include = normalizePatterns(raw.Include)
	}

	if len(raw.Exclude) > 0 {
		config.Exclude = normalizePatterns(raw.Exclude)
	}

	if len(config.Include) == 0 {
		return Config{}, errors.New("config include cannot be empty")
	}

	for _, pattern := range config.Include {
		if !doublestar.ValidatePattern(pattern) {
			return Config{}, fmt.Errorf("invalid include pattern %q", pattern)
		}
	}

	return config, nil
}

func normalizePatterns(patterns []string) []string {
	normalized := make([]string, 0, len(patterns))
	for _, pattern := range patterns {
		trimmed := strings.TrimSpace(pattern)
		if trimmed == "" {
			continue
		}

		normalized = append(normalized, filepath.ToSlash(trimmed))
	}

	return normalized
}
