#!/usr/bin/env node

import fs from "node:fs/promises";
import path from "node:path";
import process from "node:process";
import { spawnSync } from "node:child_process";
import { pathToFileURL } from "node:url";

const DEFAULT_CONFIG_PATH = ".line-audit.json";
const DEFAULT_CONFIG = {
  threshold: 400,
};

export async function runLineAudit(options = {}) {
  const args = options.args ?? process.argv.slice(2);
  const cwd = options.cwd ?? process.cwd();

  const configPath = resolveConfigPath(args);
  const config = await loadConfig(configPath, cwd);
  const fileCandidates = await listCandidates(cwd);

  const stats = [];
  let binarySkipped = 0;

  for (const relPath of fileCandidates) {
    const absPath = path.join(cwd, relPath);
    let buffer;
    try {
      buffer = await fs.readFile(absPath);
    } catch {
      continue;
    }

    if (isBinaryBuffer(buffer)) {
      binarySkipped += 1;
      continue;
    }

    const lines = countPhysicalLines(buffer);
    stats.push({ path: relPath, lines });
  }

  stats.sort((a, b) => b.lines - a.lines || a.path.localeCompare(b.path));
  const totalLines = stats.reduce((sum, item) => sum + item.lines, 0);
  const overThreshold = stats.filter((item) => item.lines >= config.threshold);

  return buildMarkdownReport({
    cwd,
    config,
    stats,
    totalLines,
    binarySkipped,
    overThreshold,
  });
}

export function resolveConfigPath(args) {
  const flagIndex = args.findIndex((arg) => arg === "--config" || arg === "-c");
  if (flagIndex >= 0 && args[flagIndex + 1]) {
    return args[flagIndex + 1];
  }
  return DEFAULT_CONFIG_PATH;
}

export async function loadConfig(configPath, cwd = process.cwd()) {
  let merged = { ...DEFAULT_CONFIG };
  try {
    const raw = await fs.readFile(path.resolve(cwd, configPath), "utf8");
    const parsed = JSON.parse(raw);
    merged = {
      ...merged,
      ...parsed,
    };
  } catch (error) {
    if (error && error.code !== "ENOENT") {
      throw new Error(
        `Config parse failed at "${configPath}": ${error.message}`,
      );
    }
  }

  const threshold = Number(merged.threshold);

  return {
    threshold:
      Number.isFinite(threshold) && threshold > 0 ? Math.floor(threshold) : 400,
  };
}

export async function listCandidates(cwd) {
  if (isGitRepo(cwd)) {
    const result = spawnSync(
      "git",
      ["ls-files", "--cached", "--others", "--exclude-standard", "-z"],
      { cwd, encoding: "utf8" },
    );
    if (result.status === 0) {
      return result.stdout
        .split("\0")
        .map((item) => item.trim())
        .filter(Boolean)
        .filter((item) => !isLikelyDirectoryPath(item));
    }
  }

  const rgFiles = listByRipgrep(cwd);
  if (rgFiles.length > 0) {
    return rgFiles;
  }

  const out = [];
  await walkDir(cwd, "", out);
  return out;
}

function listByRipgrep(cwd) {
  const result = spawnSync("rg", ["--files", "-0"], {
    cwd,
    encoding: "utf8",
  });
  if (result.status !== 0) {
    return [];
  }
  return result.stdout
    .split("\0")
    .map((item) => item.trim())
    .filter(Boolean)
    .filter((item) => !isLikelyDirectoryPath(item));
}

function isGitRepo(cwd) {
  const check = spawnSync("git", ["rev-parse", "--is-inside-work-tree"], {
    cwd,
    encoding: "utf8",
  });
  return check.status === 0 && check.stdout.trim() === "true";
}

async function walkDir(root, relDir, out) {
  const absDir = path.join(root, relDir);
  let entries;
  try {
    entries = await fs.readdir(absDir, { withFileTypes: true });
  } catch {
    return;
  }

  for (const entry of entries) {
    const relPath = relDir ? path.posix.join(relDir, entry.name) : entry.name;
    if (entry.isDirectory()) {
      if (shouldSkipDir(entry.name)) {
        continue;
      }
      await walkDir(root, relPath, out);
      continue;
    }
    if (entry.isFile()) {
      out.push(relPath);
    }
  }
}

function shouldSkipDir(name) {
  return (
    name === ".git" ||
    name === "node_modules" ||
    name === ".idea" ||
    name === ".vscode" ||
    name === ".cursor"
  );
}

function isLikelyDirectoryPath(relPath) {
  return relPath.endsWith("/");
}

export function isBinaryBuffer(buffer) {
  if (buffer.length === 0) {
    return false;
  }

  const sampleSize = Math.min(buffer.length, 8000);
  let suspicious = 0;

  for (let i = 0; i < sampleSize; i += 1) {
    const byte = buffer[i];
    if (byte === 0) {
      return true;
    }
    const isControl =
      (byte >= 1 && byte <= 8) ||
      byte === 11 ||
      byte === 12 ||
      (byte >= 14 && byte <= 31);
    if (isControl) {
      suspicious += 1;
    }
  }

  return suspicious / sampleSize > 0.3;
}

export function countPhysicalLines(buffer) {
  let lines = 0;
  for (let i = 0; i < buffer.length; i += 1) {
    if (buffer[i] === 10) {
      lines += 1;
    }
  }
  return lines;
}

export function buildMarkdownReport({ config, overThreshold }) {
  const lines = [];
  lines.push("# 文件行数审查结果");
  lines.push("");
  lines.push(`## 超出阈值文件（>= ${config.threshold} 行）`);
  lines.push("");
  if (overThreshold.length === 0) {
    lines.push("- 无");
  } else {
    for (const item of overThreshold) {
      lines.push(`- ${item.path} ${item.lines}行`);
    }
  }

  return lines.join("\n");
}

async function main() {
  const report = await runLineAudit();
  process.stdout.write(report + "\n");
}

function isDirectExecution() {
  if (!process.argv[1]) {
    return false;
  }
  return import.meta.url === pathToFileURL(process.argv[1]).href;
}

if (isDirectExecution()) {
  main().catch((error) => {
    process.stderr.write(`line-audit failed: ${error.message}\n`);
    process.exit(1);
  });
}
