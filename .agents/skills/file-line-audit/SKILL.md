---
name: file-line-audit
description: 审查项目文件行数并输出超阈值文件列表。在需要识别超长文件、评估拆分优先级、按 .gitignore 规则统计行数时使用。
---

# 文件行数审查

## 触发场景

- 文件过长，需要定位拆分目标
- 需要按 `.gitignore` 规则做行数审查
- 需要输出超阈值文件清单

## 执行步骤

1. 在项目根目录执行：

```bash
node .agents/skills/file-line-audit/scripts/line-audit.mjs
```

2. 需要自定义阈值时，先创建 `.line-audit.json`：

```json
{
  "threshold": 400
}
```

3. 使用自定义配置运行：

```bash
node .agents/skills/file-line-audit/scripts/line-audit.mjs --config .line-audit.json
```

## 输出约定

- 仅输出“超出阈值文件列表”
- 列表项格式：`- 路径 行数`
