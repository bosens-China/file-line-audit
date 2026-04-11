# File Line Audit

[简体中文](./README.zh-CN.md)

An Agent Skill designed to audit source files that exceed 400 lines (default) in any repository. You can customize the threshold and scan patterns by creating a local configuration file.

## Installation

Install and add this skill to your agent:

```bash
npx skills add bosens-China/file-line-audit
```

## Customization

To override the default 400-line threshold or change the scan scope, create a `.line-audit.json` file in your repository root based on the [example config](./.agents/skills/file-line-audit/assets/default-config.json).

## License

MIT
