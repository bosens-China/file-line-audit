---
name: file-line-audit
description: Audit oversized source files in a repository. Apply repository .gitignore rules plus extra exclude rules from config, filter candidate files with include glob patterns, count physical lines, and output only files at or above the threshold. Use this skill when the user asks to find long files, review file-length distribution, identify split/refactor candidates, or analyze technical debt.
license: MIT
compatibility: Requires a bundled binary from scripts/ matching the host OS and CPU architecture. Current packaged targets are Windows amd64, Linux amd64, Linux arm64, macOS amd64, and macOS arm64.
allowed-tools: run_shell_command
metadata:
  author: github.com/bosens-China
  version: "0.2.1"
---

# File Line Audit

Use this skill to audit file line counts in a project and return only files that exceed a configured threshold.

When using this skill:
- Run the packaged binary from `scripts/` that matches the current OS and CPU architecture.
- Always run from the target repository root so `.gitignore` and `.line-audit.json` are resolved correctly.
- Respect repository `.gitignore` rules first, then apply the skill config `exclude` rules on top.
- Use `include` patterns to limit scan scope.
- Output only the files whose physical line count is greater than or equal to the configured threshold.

## Installation

```bash
npx skills add bosens-China/file-line-audit
```

## Binary Selection

Choose the packaged executable under `scripts/`:

- Windows amd64: `scripts/line-audit-windows-amd64.exe`
- Linux amd64: `scripts/line-audit-linux-amd64`
- Linux arm64: `scripts/line-audit-linux-arm64`
- macOS amd64: `scripts/line-audit-darwin-amd64`
- macOS arm64: `scripts/line-audit-darwin-arm64`

## Steps

1. Confirm you are in the repository root.
2. Check for local configuration:
   - If the repository has `.line-audit.json`, use it.
   - If not, use the built-in example: `.agents/skills/file-line-audit/assets/default-config.json`.
3. Select the correct binary from `scripts/` for the current platform.
4. On Linux/macOS, ensure the binary is executable: `chmod +x <binary_path>`.
5. Run the binary with `--config` pointing to the selected configuration.
6. Return only the over-threshold file list to the user.

## Commands

### Run with auto-detected or example config

On Unix-like systems:
```bash
chmod +x .agents/skills/file-line-audit/scripts/line-audit-<target>
.agents/skills/file-line-audit/scripts/line-audit-<target> --config .agents/skills/file-line-audit/assets/default-config.json
```

On Windows:
```powershell
& .\.agents\skills\file-line-audit\scripts\line-audit-windows-amd64.exe --config .agents/skills/file-line-audit/assets/default-config.json
```

### Create a local config (Optional)
Only do this if the user specifically requests a persistent local configuration.
```bash
cp .agents/skills/file-line-audit/assets/default-config.json .line-audit.json
```

## Config Rules

- `threshold`: minimum line count to report, default `400`
- `include`: glob patterns for candidate files
- `exclude`: extra ignore rules in `.gitignore` syntax, applied after repository `.gitignore`

Default config intent:
- Focus on common source roots such as `src`, `app`, `apps`, `pkg`, `packages`, `lib`, `libs`, `cmd`, `internal`, `backend`, `frontend`, `client`, `server`, `service`, `services`, `api`, and `web`
- Cover mainstream source extensions including JavaScript, TypeScript, Vue, Python, Go, Rust, Java, Kotlin, Scala, Ruby, PHP, C#, Swift, and Dart

## Output Format

Return the tool output directly. The expected format is:

```text
# File Line Audit

## Files Over Threshold (>= 400 lines)

- src/example.ts 512
- apps/web/pages/home.tsx 438
```

## Limitations

- **Physical Lines Only**: The tool counts raw newlines and does not distinguish between code, comments, or blank lines.
- **Binary Files**: Automatically skipped.
- **Git Context**: Relies on `git` being available in the environment to resolve `.gitignore` rules effectively.
- **Performance**: Optimized for source code; avoid running on directories containing large data files or build artifacts not covered by `.gitignore`.

## Notes

- Repository `.gitignore` is always active even if config does not mention it.
- Nested `.gitignore` files are respected.
