import test from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import { spawnSync } from "node:child_process";

import {
  buildMarkdownReport,
  countPhysicalLines,
  isBinaryBuffer,
  loadConfig,
  resolveConfigPath,
} from "../scripts/line-audit.mjs";

test("resolveConfigPath returns passed config path", () => {
  assert.equal(resolveConfigPath(["--config", "custom.json"]), "custom.json");
  assert.equal(resolveConfigPath(["-c", "x.json"]), "x.json");
});

test("resolveConfigPath falls back to default path", () => {
  assert.equal(resolveConfigPath([]), ".line-audit.json");
});

test("countPhysicalLines counts newline characters", () => {
  assert.equal(countPhysicalLines(Buffer.from("")), 0);
  assert.equal(countPhysicalLines(Buffer.from("one")), 0);
  assert.equal(countPhysicalLines(Buffer.from("one\ntwo\nthree")), 2);
  assert.equal(countPhysicalLines(Buffer.from("one\ntwo\n")), 2);
});

test("isBinaryBuffer identifies binary-like content", () => {
  assert.equal(isBinaryBuffer(Buffer.from("hello\nworld\n")), false);
  assert.equal(isBinaryBuffer(Buffer.from([0, 1, 2, 3])), true);
});

test("loadConfig reads and normalizes config values", async () => {
  const tempDir = await fs.mkdtemp(path.join(os.tmpdir(), "line-audit-test-"));
  const configPath = path.join(tempDir, ".line-audit.json");
  await fs.writeFile(
    configPath,
    JSON.stringify({ threshold: "520.8" }),
    "utf8",
  );

  const cfg = await loadConfig(".line-audit.json", tempDir);
  assert.deepEqual(cfg, {
    threshold: 520,
  });
});

test("buildMarkdownReport only lists exceeded files as bullet lines", () => {
  const report = buildMarkdownReport({
    config: { threshold: 400 },
    overThreshold: [{ path: "a.js", lines: 500 }],
  });

  assert.match(report, /文件行数审查结果/);
  assert.match(report, /超出阈值文件（>= 400 行）/);
  assert.match(report, /- a\.js 500行/);
  assert.doesNotMatch(report, /\| File \| Lines \| Ratio \|/);
});

test("runLineAudit returns only files over threshold", async () => {
  const tempDir = await fs.mkdtemp(path.join(os.tmpdir(), "line-audit-e2e-"));
  await fs.writeFile(path.join(tempDir, "small.js"), "a\nb\n", "utf8");
  await fs.writeFile(path.join(tempDir, "large.js"), "1\n2\n3\n4\n5\n6\n", "utf8");
  await fs.writeFile(
    path.join(tempDir, ".line-audit.json"),
    JSON.stringify({ threshold: 5 }),
    "utf8",
  );

  const { runLineAudit } = await import("../scripts/line-audit.mjs");
  const report = await runLineAudit({ cwd: tempDir, args: [] });

  assert.match(report, /超出阈值文件（>= 5 行）/);
  assert.match(report, /- large\.js 6行/);
  assert.doesNotMatch(report, /- small\.js/);
});

test("runLineAudit respects gitignore rules and skips binary files", async () => {
  const tempDir = await fs.mkdtemp(path.join(os.tmpdir(), "line-audit-git-"));

  const gitInit = spawnSync("git", ["init"], {
    cwd: tempDir,
    encoding: "utf8",
  });
  assert.equal(gitInit.status, 0, gitInit.stderr);

  await fs.mkdir(path.join(tempDir, "nested"), { recursive: true });
  await fs.writeFile(path.join(tempDir, ".gitignore"), "ignored-root.js\n", "utf8");
  await fs.writeFile(
    path.join(tempDir, "nested", ".gitignore"),
    "ignored-nested.js\n",
    "utf8",
  );
  await fs.writeFile(path.join(tempDir, ".hidden.js"), "1\n2\n3\n", "utf8");
  await fs.writeFile(path.join(tempDir, "ignored-root.js"), "1\n2\n3\n4\n", "utf8");
  await fs.writeFile(path.join(tempDir, "nested", "ignored-nested.js"), "1\n2\n3\n4\n", "utf8");
  await fs.writeFile(path.join(tempDir, "nested", "kept.js"), "1\n2\n3\n4\n", "utf8");
  await fs.writeFile(
    path.join(tempDir, "binary.bin"),
    Buffer.from([0, 1, 2, 3, 4, 5]),
  );
  await fs.writeFile(
    path.join(tempDir, ".line-audit.json"),
    JSON.stringify({ threshold: 3 }),
    "utf8",
  );

  const { runLineAudit } = await import("../scripts/line-audit.mjs");
  const report = await runLineAudit({ cwd: tempDir, args: [] });

  assert.match(report, /- \.hidden\.js 3行/);
  assert.match(report, /- nested\/kept\.js 4行/);
  assert.doesNotMatch(report, /ignored-root\.js/);
  assert.doesNotMatch(report, /ignored-nested\.js/);
  assert.doesNotMatch(report, /binary\.bin/);
});
