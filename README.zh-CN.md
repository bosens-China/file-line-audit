# File Line Audit

一个用于快速识别仓库中超长源码文件的 Agent Skill 和 Go 命令行工具。

## 安装

```bash
npx skills add bosens-China/file-line-audit
```

## 使用

运行对应平台的预编译二进制文件：

```bash
./.agents/skills/file-line-audit/scripts/line-audit-<target>
```

如果需要自定义阈值或扫描模式，请参考示例配置：

```bash
cp .agents/skills/file-line-audit/.line-audit.example.json .line-audit.json
```

## 开源协议

MIT
