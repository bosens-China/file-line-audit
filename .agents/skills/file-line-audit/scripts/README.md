# Packaged Binaries

This directory contains packaged `file-line-audit` executables for Agent Skills clients.

Target naming convention:

- `line-audit-windows-amd64.exe`
- `line-audit-linux-amd64`
- `line-audit-linux-arm64`
- `line-audit-darwin-amd64`
- `line-audit-darwin-arm64`

Rebuild from the repository root with Go:

```bash
go test ./...
./scripts/build-skill.sh
```

On Windows, you can use:

```powershell
./scripts/build-skill.ps1
```
