# File Line Audit

[English](./README.md)

一个专门为 Agent Skills 设计的工具，用于审计仓库中超过 400 行（默认）的源代码文件。你可以通过创建本地配置文件来灵活自定义审查阈值和扫描模式。

## 安装

将此技能添加至你的 Agent：

```bash
npx skills add bosens-China/file-line-audit
```

## 自定义配置

如果需要修改默认的 400 行阈值或扫描范围，请参考 [示例配置](./.agents/skills/file-line-audit/assets/default-config.json) 在项目根目录创建 `.line-audit.json` 文件。

## 开源协议

MIT
