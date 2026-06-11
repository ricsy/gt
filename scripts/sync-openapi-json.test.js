const assert = require("node:assert/strict");
const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");
const test = require("node:test");

const {
  formatOpenApiJson,
  parseArgs,
  syncOpenApiJson,
  writeIfChanged,
} = require("./sync-openapi-json");

test("parseArgs supports config, root, and output overrides", () => {
  const parsed = parseArgs(["--config", "custom.json", "--root", "repo", "--output", "spec.json"]);

  assert.deepEqual(parsed, {
    configPath: "custom.json",
    root: "repo",
    outputPath: "spec.json",
  });
});

test("formatOpenApiJson produces stable pretty JSON with trailing newline", () => {
  assert.equal(
    formatOpenApiJson({ swagger: "2.0", info: { title: "demo" } }),
    '{\n  "swagger": "2.0",\n  "info": {\n    "title": "demo"\n  }\n}\n',
  );
});

test("writeIfChanged only rewrites files when content changes", () => {
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), "sync-openapi-write-"));
  const filePath = path.join(tempDir, "api.json");

  assert.equal(writeIfChanged(filePath, "first\n"), true);
  assert.equal(writeIfChanged(filePath, "first\n"), false);
  assert.equal(writeIfChanged(filePath, "second\n"), true);
  assert.equal(fs.readFileSync(filePath, "utf8"), "second\n");
});

test("syncOpenApiJson writes the fetched spec to configured localPath", async () => {
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), "sync-openapi-run-"));
  const expected = { swagger: "2.0", paths: { "/v5/user": { get: {} } } };

  const result = await syncOpenApiJson({
    rootDir: tempDir,
    loadConfig: () => ({
      openapi: { localPath: "data/api-1.json" },
    }),
    fetchOpenApiJson: async () => expected,
  });

  assert.equal(result.updated, true);
  assert.equal(result.outputPath, path.join(tempDir, "data", "api-1.json"));
  assert.deepEqual(JSON.parse(fs.readFileSync(result.outputPath, "utf8")), expected);

  const second = await syncOpenApiJson({
    rootDir: tempDir,
    loadConfig: () => ({
      openapi: { localPath: "data/api-1.json" },
    }),
    fetchOpenApiJson: async () => expected,
  });

  assert.equal(second.updated, false);
});
