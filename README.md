# File Line Audit

An Agent Skill and Go-powered CLI to quickly identify oversized source files in any repository.

## Installation

```bash
npx skills add bosens-China/file-line-audit
```

## Usage

Run the pre-packaged binary for your platform:

```bash
./.agents/skills/file-line-audit/scripts/line-audit-<target>
```

To use a custom threshold or patterns, copy the example config:

```bash
cp .agents/skills/file-line-audit/assets/line-audit.example.json .line-audit.json
```

## License

MIT
