# File Line Audit Skill

[中文说明](./README.zh-CN.md)

`file-line-audit` audits file line counts in a project.  
In a Git repository it follows Git's tracked and untracked file set with `.gitignore` rules applied, including nested `.gitignore` files. It skips binary files, counts `\n` characters in each text file, and prints only files that exceed the configured threshold.

## Install via npx skills

Install this skill from the local path:

```bash
npx skills add . --skill file-line-audit
```

## Configuration

`.agents/skills/file-line-audit/.line-audit.example.json` is the template config file.  
Copy it to `.line-audit.json` in your project root and adjust as needed:

- `threshold`: line-count threshold for reporting (default `400`)

## License

MIT
