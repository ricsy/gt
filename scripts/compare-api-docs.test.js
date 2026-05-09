const assert = require("node:assert/strict");
const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");
const test = require("node:test");

const {
  collectOpenApiOperations,
  compareApiDocs,
  createDefaultConfig,
  formatFindings,
  getLiveOpenApiUrl,
  loadConfig,
  normalizeGoPath,
  normalizeOpenApiPath,
} = require("./compare-api-docs");

test("normalizes OpenAPI and Go endpoint templates into comparable paths", () => {
  const config = createDefaultConfig();

  assert.deepEqual(normalizeOpenApiPath("/v5/repos/{owner}/{repo}/issues", config), {
    normalizedPath: "/repos/{param1}/{param2}/issues",
    params: ["owner", "repo"],
  });

  assert.deepEqual(normalizeGoPath("/repos/%s/%s/issues"), {
    normalizedPath: "/repos/{param1}/{param2}/issues",
    params: ["arg1", "arg2"],
  });
});

test("collects OpenAPI operations with sorted parameters", () => {
  const config = createDefaultConfig();
  const spec = {
    paths: {
      "/v5/repos/{owner}/{repo}/issues": {
        get: {
          parameters: [
            { name: "repo", in: "path" },
            { name: "owner", in: "path" },
            { name: "state", in: "query" },
          ],
        },
      },
    },
  };

  assert.deepEqual(collectOpenApiOperations(spec, config), [
    {
      method: "GET",
      normalizedPath: "/repos/{param1}/{param2}/issues",
      originalPath: "/v5/repos/{owner}/{repo}/issues",
      parameters: [
        { in: "path", name: "owner" },
        { in: "path", name: "repo" },
        { in: "query", name: "state" },
      ],
      pathParams: ["owner", "repo"],
    },
  ]);
});

test("filters configured OpenAPI parameters that are handled outside request structs", () => {
  const config = createDefaultConfig();
  const spec = {
    paths: {
      "/v5/user": {
        get: {
          parameters: [
            { name: "access_token", in: "query" },
            { name: "page", in: "query" },
          ],
        },
      },
    },
  };

  assert.deepEqual(collectOpenApiOperations(spec, config), [
    {
      method: "GET",
      normalizedPath: "/user",
      originalPath: "/v5/user",
      parameters: [{ in: "query", name: "page" }],
      pathParams: [],
    },
  ]);
});

test("compares endpoints and parameters deterministically", () => {
  const config = createDefaultConfig();
  const openApi = {
    paths: {
      "/v5/repos/{owner}/{repo}/issues": {
        get: {
          parameters: [
            { name: "owner", in: "path" },
            { name: "repo", in: "path" },
            { name: "state", in: "query" },
          ],
        },
      },
    },
  };
  const go = {
    endpoints: [
      {
        group: "Issues",
        action: "List",
        method: "POST",
        path: "/repos/%s/%s/issues",
      },
      {
        group: "Extra",
        action: "List",
        method: "GET",
        path: "/extra",
      },
    ],
    requestStructs: {
      ListIssuesOptions: {
        fields: [{ name: "Page", paramName: "page" }],
      },
    },
    directCalls: [{ file: "pkg/api/issue.go", line: 10, method: "GET", path: "path" }],
  };

  const findings = compareApiDocs(openApi, go, config);

  assert.deepEqual(
    findings.different.map((finding) => finding.code),
    ["endpoint.method"],
  );
  assert.deepEqual(
    findings.missing.map((finding) => finding.code),
    ["parameter.query"],
  );
  assert.deepEqual(
    findings.extra.map((finding) => finding.code),
    ["endpoint.extra"],
  );
  assert.deepEqual(
    findings.unknown.map((finding) => finding.code),
    ["go.direct_call"],
  );

  assert.match(formatFindings(findings), /endpoint\.method/);
  assert.match(formatFindings(findings), /parameter\.query/);
});

test("default config keeps repository-specific values outside comparison logic", () => {
  const config = createDefaultConfig();

  assert.equal(config.openapi.localPath, "data/api-1.json");
  assert.equal(config.openapi.basePathPrefix, "/v5");
  assert.equal(config.openapi.live.pageUrl, "https://help.gitee.com/openapi/v5");
  assert.equal(config.openapi.live.downloadSelector, ".download-button");
  assert.equal(config.openapi.live.directUrl, "https://gitee.com/sdk/typescript-sdk-v5/raw/main/openapi-spec.json");
  assert.deepEqual(config.openapi.ignoredParameters, ["access_token"]);
  assert.equal(config.go.endpointFile, "pkg/api/endpoint.go");
  assert.deepEqual(config.go.endpointGlobs, ["pkg/api/endpoint.go", "pkg/api/openapi_coverage.go"]);
  assert.equal(config.go.coverageEndpointFile, "pkg/api/openapi_coverage.go");
  assert.equal(config.go.compatibilityOptionsFile, "pkg/api/response/openapi_compatibility.go");
  assert.equal(config.go.compatibilityOptionsType, "OpenAPICompatibilityOptions");
  assert.equal(config.go.endpointGroupType, "EndpointGroup");
  assert.deepEqual(config.go.requestStructSuffixes, ["Options", "Request"]);
  assert.equal(config.go.parameterNameCase, "snake_case");
});

test("loadConfig merges project config over generic defaults", () => {
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), "api-doc-diff-"));
  const configPath = path.join(tempDir, "config.json");
  fs.writeFileSync(
    configPath,
    JSON.stringify({
      openapi: { basePathPrefix: "/api" },
      go: {
        endpointFile: "src/endpoints.go",
        endpointGroupType: "RouteGroup",
      },
    }),
  );

  const config = loadConfig(configPath);

  assert.equal(config.openapi.localPath, "data/api-1.json");
  assert.equal(config.openapi.basePathPrefix, "/api");
  assert.equal(config.go.endpointFile, "src/endpoints.go");
  assert.equal(config.go.endpointGroupType, "RouteGroup");
  assert.deepEqual(config.go.requestStructSuffixes, ["Options", "Request"]);
});

test("live OpenAPI URL prefers configured direct URL over page scraping", () => {
  const config = createDefaultConfig();

  assert.equal(getLiveOpenApiUrl("<html></html>", config), config.openapi.live.directUrl);
});
