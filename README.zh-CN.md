# File Line Audit Skill

`file-line-audit` 是一个用于审查项目文件长度分布的 Agent Skill。  
在 Git 仓库中，它会按 Git 的 tracked/untracked 文件集合并结合 `.gitignore` 规则（含子目录规则）取文件，跳过二进制文件，按文本文件中的 `\n` 个数计数，并仅输出“超出阈值文件列表”。

## 通过 npx skills 添加

从本地路径安装该 Skill：

```bash
npx skills add . --skill file-line-audit
```

## 配置说明

`.agents/skills/file-line-audit/.line-audit.example.json` 是示例配置文件模板。  
你可以复制为项目根目录的 `.line-audit.json` 并手动调整：

- `threshold`：超出提醒阈值（默认 `400`）

## 协议

MIT
