#!/usr/bin/env node

const { spawnSync } = require("node:child_process");
const fs = require("node:fs");
const path = require("node:path");

const HTTP_METHODS = new Set(["get", "post", "put", "patch", "delete"]);
const DEFAULT_CONFIG_PATH = "scripts/api-doc-diff.config.json";

function createDefaultConfig() {
  return {
    openapi: {
      localPath: "data/api-1.json",
      basePathPrefix: "/v5",
      live: {
        pageUrl: "https://help.gitee.com/openapi/v5",
        downloadSelector: ".download-button",
        directUrl: "https://gitee.com/sdk/typescript-sdk-v5/raw/main/openapi-spec.json",
      },
      ignoredParameters: ["access_token"],
    },
    go: {
      extractor: "scripts/go-ast-extract.go",
      endpointFile: "pkg/api/endpoint.go",
      endpointGroupType: "EndpointGroup",
      responseGlobs: ["pkg/api/response/*.go"],
      apiGlobs: ["pkg/api/*.go"],
      endpointGlobs: ["pkg/api/endpoint.go", "pkg/api/openapi_coverage.go"],
      coverageEndpointFile: "pkg/api/openapi_coverage.go",
      compatibilityOptionsFile: "pkg/api/response/openapi_compatibility.go",
      compatibilityOptionsType: "OpenAPICompatibilityOptions",
      requestStructSuffixes: ["Options", "Request"],
      parameterTags: ["url", "form", "query", "json"],
      parameterNameCase: "snake_case",
      directCallMethod: "Do",
      ignoreDirectCallFiles: ["pkg/api/client.go"],
    },
  };
}

function loadConfig(configPath) {
  let loaded = {};
  if (configPath && fs.existsSync(configPath)) {
    try {
      loaded = JSON.parse(fs.readFileSync(configPath, "utf8"));
    } catch (error) {
      throw new UsageError(`Invalid config JSON at ${configPath}: ${error.message}`);
    }
  } else if (configPath) {
    throw new UsageError(`Config file not found: ${configPath}`);
  }
  const config = mergeConfig(createDefaultConfig(), loaded);
  validateConfig(config);
  return config;
}

function mergeConfig(base, override) {
  if (!override || typeof override !== "object" || Array.isArray(override)) {
    return base;
  }
  const result = { ...base };
  for (const [key, value] of Object.entries(override)) {
    if (value && typeof value === "object" && !Array.isArray(value)) {
      result[key] = mergeConfig(base[key] || {}, value);
    } else {
      result[key] = value;
    }
  }
  return result;
}

function validateConfig(config) {
  const requiredStrings = [
    ["openapi.localPath", config.openapi?.localPath],
    ["openapi.basePathPrefix", config.openapi?.basePathPrefix],
    ["openapi.live.pageUrl", config.openapi?.live?.pageUrl],
    ["go.extractor", config.go?.extractor],
    ["go.endpointGroupType", config.go?.endpointGroupType],
    ["go.directCallMethod", config.go?.directCallMethod],
  ];
  for (const [name, value] of requiredStrings) {
    if (!value || typeof value !== "string") {
      throw new UsageError(`Config value '${name}' must be a non-empty string.`);
    }
  }
  const requiredArrays = [
    ["go.responseGlobs", config.go?.responseGlobs],
    ["go.apiGlobs", config.go?.apiGlobs],
    ["go.endpointGlobs", config.go?.endpointGlobs],
    ["go.requestStructSuffixes", config.go?.requestStructSuffixes],
    ["go.parameterTags", config.go?.parameterTags],
  ];
  for (const [name, value] of requiredArrays) {
    if (!Array.isArray(value) || value.length === 0) {
      throw new UsageError(`Config value '${name}' must be a non-empty array.`);
    }
  }
}

class UsageError extends Error {
  constructor(message) {
    super(message);
    this.name = "UsageError";
  }
}

class DownloadError extends Error {
  constructor(message, cause) {
    super(message);
    this.name = "DownloadError";
    this.cause = cause;
  }
}

function normalizeOpenApiPath(openApiPath, config = createDefaultConfig()) {
  const params = [];
  const prefix = escapeRegExp(config.openapi.basePathPrefix || "");
  const prefixPattern = prefix ? new RegExp(`^${prefix}(?=/|$)`) : null;
  let normalizedPath = prefixPattern ? openApiPath.replace(prefixPattern, "") : openApiPath;
  normalizedPath = normalizedPath.replace(/\{([^}]+)\}/g, (_match, name) => {
    params.push(name);
    return `{param${params.length}}`;
  });
  return { normalizedPath: normalizedPath || "/", params };
}

function normalizeGoPath(goPath) {
  const params = [];
  const normalizedPath = goPath.replace(/%[-+#0-9.]*[sdv]/g, () => {
    params.push(`arg${params.length + 1}`);
    return `{param${params.length}}`;
  });
  return { normalizedPath: normalizedPath || "/", params };
}

function collectOpenApiOperations(openApi, config = createDefaultConfig()) {
  if (!openApi || typeof openApi !== "object" || !openApi.paths) {
    throw new UsageError("OpenAPI document must contain a top-level paths object.");
  }

  const operations = [];
  for (const originalPath of Object.keys(openApi.paths).sort()) {
    const pathItem = openApi.paths[originalPath] || {};
    const commonParams = Array.isArray(pathItem.parameters) ? pathItem.parameters : [];
    for (const methodName of Object.keys(pathItem).sort()) {
      if (!HTTP_METHODS.has(methodName.toLowerCase())) {
        continue;
      }
      const operation = pathItem[methodName] || {};
      const { normalizedPath, params: pathParams } = normalizeOpenApiPath(originalPath, config);
      const parameters = collectOperationParameters(openApi, config, commonParams, operation.parameters || []);
      operations.push({
        method: methodName.toUpperCase(),
        normalizedPath,
        originalPath,
        parameters,
        pathParams,
      });
    }
  }
  return operations;
}

function collectOperationParameters(openApi, config, ...parameterLists) {
  const byKey = new Map();
  const ignoredParameters = new Set(config.openapi.ignoredParameters || []);
  for (const list of parameterLists) {
    for (const parameter of list) {
      if (!parameter || !parameter.name) {
        continue;
      }
      if (ignoredParameters.has(parameter.name)) {
        continue;
      }
      if (parameter.in === "body" && parameter.schema) {
        for (const name of schemaPropertyNames(openApi, parameter.schema)) {
          if (ignoredParameters.has(name)) {
            continue;
          }
          byKey.set(`body:${name}`, { in: "body", name });
        }
        continue;
      }
      const location = parameter.in || "unknown";
      byKey.set(`${location}:${parameter.name}`, { in: location, name: parameter.name });
    }
  }
  return [...byKey.values()].sort(compareParameter);
}

function schemaPropertyNames(openApi, schema, seen = new Set()) {
  if (!schema || typeof schema !== "object") {
    return [];
  }
  if (schema.$ref) {
    const ref = schema.$ref.replace(/^#\//, "").split("/");
    let target = openApi;
    for (const part of ref) {
      target = target && target[part];
    }
    if (!target || seen.has(schema.$ref)) {
      return [];
    }
    seen.add(schema.$ref);
    return schemaPropertyNames(openApi, target, seen);
  }
  if (schema.properties && typeof schema.properties === "object") {
    return Object.keys(schema.properties).sort();
  }
  if (schema.items) {
    return schemaPropertyNames(openApi, schema.items, seen);
  }
  return [];
}

function compareApiDocs(openApi, goMetadata, config = createDefaultConfig()) {
  const openApiOperations = collectOpenApiOperations(openApi, config);
  const goEndpoints = (goMetadata.endpoints || []).map((endpoint) => ({
    ...endpoint,
    method: String(endpoint.method || "").toUpperCase(),
    normalizedPath: normalizeGoPath(endpoint.path || "").normalizedPath,
  }));
  const goParameterNames = collectGoParameterNames(goMetadata.requestStructs || {});

  const findings = {
    missing: [],
    extra: [],
    different: [],
    unknown: [],
  };

  for (const operation of openApiOperations) {
    const samePath = goEndpoints.filter((endpoint) => endpoint.normalizedPath === operation.normalizedPath);
    if (samePath.length === 0) {
      findings.missing.push({
        code: "endpoint.missing",
        message: `${operation.method} ${operation.originalPath} is missing from pkg/api/endpoint.go`,
        path: operation.normalizedPath,
        method: operation.method,
      });
      continue;
    }

    const sameMethod = samePath.find((endpoint) => endpoint.method === operation.method);
    if (!sameMethod) {
      findings.different.push({
        code: "endpoint.method",
        message: `${operation.originalPath} has OpenAPI method ${operation.method}, Go methods: ${samePath.map((endpoint) => endpoint.method).sort().join(", ")}`,
        path: operation.normalizedPath,
        method: operation.method,
      });
    }

    for (const parameter of operation.parameters) {
      if (parameter.in === "path") {
        continue;
      }
      if (!goParameterNames.has(parameter.name)) {
        findings.missing.push({
          code: `parameter.${parameter.in}`,
          message: `${operation.method} ${operation.originalPath} documents ${parameter.in} parameter '${parameter.name}' with no matching request/options field`,
          path: operation.normalizedPath,
          method: operation.method,
          parameter: parameter.name,
        });
      }
    }
  }

  const openApiEndpointKeys = new Set(openApiOperations.map((operation) => `${operation.method} ${operation.normalizedPath}`));
  const openApiPaths = new Set(openApiOperations.map((operation) => operation.normalizedPath));
  for (const endpoint of goEndpoints) {
    const key = `${endpoint.method} ${endpoint.normalizedPath}`;
    if (!openApiEndpointKeys.has(key) && !openApiPaths.has(endpoint.normalizedPath)) {
      findings.extra.push({
        code: "endpoint.extra",
        message: `${endpoint.group}.${endpoint.action} defines ${endpoint.method} ${endpoint.path} not found in OpenAPI document`,
        path: endpoint.normalizedPath,
        method: endpoint.method,
      });
    }
  }

  for (const call of goMetadata.directCalls || []) {
    findings.unknown.push({
      code: "go.direct_call",
      message: `${call.file}:${call.line} calls Client.Do directly (${call.method || "unknown method"}, ${call.path || "unknown path"}) and cannot be confidently mapped`,
      file: call.file,
      line: call.line,
    });
  }

  for (const key of Object.keys(findings)) {
    findings[key].sort(compareFinding);
  }
  return findings;
}

function collectGoParameterNames(requestStructs) {
  const names = new Set();
  for (const requestStruct of Object.values(requestStructs)) {
    for (const field of requestStruct.fields || []) {
      if (field.paramName) {
        names.add(field.paramName);
      }
    }
  }
  return names;
}

function hasFindings(findings) {
  return Object.values(findings).some((items) => items.length > 0);
}

function formatFindings(findings) {
  const sections = [
    ["missing", "Missing"],
    ["different", "Different"],
    ["extra", "Extra"],
    ["unknown", "Unknown"],
  ];
  const lines = [];
  for (const [key, title] of sections) {
    const items = findings[key] || [];
    if (items.length === 0) {
      continue;
    }
    lines.push(`${title} (${items.length})`);
    for (const item of items) {
      lines.push(`  - [${item.code}] ${item.message}`);
    }
  }
  return lines.join("\n");
}

function formatSummary(findings) {
  return `missing=${findings.missing.length} different=${findings.different.length} extra=${findings.extra.length} unknown=${findings.unknown.length}`;
}

function loadOpenApiJson(filePath) {
  let text;
  try {
    text = fs.readFileSync(filePath, "utf8");
  } catch (error) {
    throw new UsageError(`Unable to read OpenAPI document at ${filePath}: ${error.message}`);
  }
  try {
    const parsed = JSON.parse(text);
    if (!parsed.paths || typeof parsed.paths !== "object") {
      throw new Error("missing top-level paths object");
    }
    return parsed;
  } catch (error) {
    throw new UsageError(`Invalid OpenAPI JSON at ${filePath}: ${error.message}`);
  }
}

async function fetchOfficialOpenApiJson(config = createDefaultConfig()) {
  const pageUrl = config.openapi.live.pageUrl;
  const downloadUrl = await getLiveOpenApiUrlFromNetwork(config);

  let downloadResponse;
  try {
    downloadResponse = await fetch(downloadUrl);
  } catch (error) {
    throw new DownloadError(`Unable to download official OpenAPI JSON ${downloadUrl}: ${error.message}`, error);
  }
  if (!downloadResponse.ok) {
    throw new DownloadError(`Official OpenAPI download returned HTTP ${downloadResponse.status}: ${downloadUrl}`);
  }
  try {
    const json = await downloadResponse.json();
    if (!json.paths || typeof json.paths !== "object") {
      throw new Error("missing top-level paths object");
    }
    return json;
  } catch (error) {
    throw new DownloadError(`Official OpenAPI download did not return valid OpenAPI JSON: ${error.message}`, error);
  }
}

async function getLiveOpenApiUrlFromNetwork(config = createDefaultConfig()) {
  if (config.openapi.live.directUrl) {
    return config.openapi.live.directUrl;
  }

  const pageUrl = config.openapi.live.pageUrl;
  let pageResponse;
  try {
    pageResponse = await fetch(pageUrl);
  } catch (error) {
    throw new DownloadError(`Unable to fetch official OpenAPI page ${pageUrl}: ${error.message}`, error);
  }
  if (!pageResponse.ok) {
    throw new DownloadError(`Official OpenAPI page returned HTTP ${pageResponse.status}: ${pageUrl}`);
  }
  const html = await pageResponse.text();
  const downloadUrl = getLiveOpenApiUrl(html, config);
  if (!downloadUrl) {
    throw new DownloadError(`Unable to find ${config.openapi.live.downloadSelector} target on ${pageUrl}`);
  }
  return downloadUrl;
}

function getLiveOpenApiUrl(html, config = createDefaultConfig()) {
  if (config.openapi.live.directUrl) {
    return config.openapi.live.directUrl;
  }
  return resolveDownloadButtonUrl(html, config.openapi.live.pageUrl, config.openapi.live.downloadSelector);
}

function resolveDownloadButtonUrl(html, baseUrl, selector = ".download-button") {
  const className = selector.startsWith(".") ? selector.slice(1) : selector;
  const buttonPattern = new RegExp(`<a\\b[^>]*class=["'][^"']*\\b${escapeRegExp(className)}\\b[^"']*["'][^>]*>`, "i");
  const button = html.match(buttonPattern);
  const tag = button ? button[0] : "";
  const href = tag.match(/\bhref=["']([^"']+)["']/i);
  if (!href) {
    return "";
  }
  return new URL(href[1], baseUrl).toString();
}

function extractGoMetadata(rootDir, configPath, config = createDefaultConfig()) {
  const result = spawnSync("go", ["run", config.go.extractor, "--root", rootDir, "--config", configPath], {
    cwd: rootDir,
    encoding: "utf8",
  });
  if (result.error) {
    throw new UsageError(`Unable to run Go AST extractor: ${result.error.message}`);
  }
  if (result.status !== 0) {
    throw new UsageError(`Go AST extractor failed: ${(result.stderr || result.stdout).trim()}`);
  }
  try {
    return JSON.parse(result.stdout);
  } catch (error) {
    throw new UsageError(`Go AST extractor returned invalid JSON: ${error.message}`);
  }
}

function parseArgs(argv) {
  const options = {
    configPath: DEFAULT_CONFIG_PATH,
    openapiPath: "",
    live: false,
    json: false,
    root: process.cwd(),
  };
  for (let i = 0; i < argv.length; i++) {
    const arg = argv[i];
    if (arg === "--help" || arg === "-h") {
      options.help = true;
    } else if (arg === "--live") {
      options.live = true;
    } else if (arg === "--json") {
      options.json = true;
    } else if (arg === "--config") {
      options.configPath = argv[++i];
      if (!options.configPath) {
        throw new UsageError("--config requires a file path");
      }
    } else if (arg === "--openapi") {
      options.openapiPath = argv[++i];
      if (!options.openapiPath) {
        throw new UsageError("--openapi requires a file path");
      }
    } else if (arg === "--root") {
      options.root = argv[++i];
      if (!options.root) {
        throw new UsageError("--root requires a directory path");
      }
    } else {
      throw new UsageError(`Unknown argument: ${arg}`);
    }
  }
  return options;
}

function helpText() {
  return [
    "Usage: node scripts/compare-api-docs.js [options]",
    "",
    "Options:",
    "  --config <path>   Comparison config JSON (default: scripts/api-doc-diff.config.json)",
    "  --openapi <path>  Override the configured local OpenAPI JSON document",
    "  --live            Fetch the configured official document via direct URL or download selector",
    "  --json            Print findings as JSON",
    "  --root <path>     Repository root (default: current working directory)",
    "  -h, --help        Show this help",
  ].join("\n");
}

async function main(argv = process.argv.slice(2)) {
  const options = parseArgs(argv);
  if (options.help) {
    console.log(helpText());
    return 0;
  }

  const rootDir = path.resolve(options.root);
  const configPath = path.resolve(rootDir, options.configPath);
  const config = loadConfig(configPath);
  const openApiPath = options.openapiPath || config.openapi.localPath;
  const openApi = options.live
    ? await fetchOfficialOpenApiJson(config)
    : loadOpenApiJson(path.resolve(rootDir, openApiPath));
  const goMetadata = extractGoMetadata(rootDir, configPath, config);
  const findings = compareApiDocs(openApi, goMetadata, config);

  if (options.json) {
    console.log(JSON.stringify({ summary: formatSummary(findings), findings }, null, 2));
  } else if (hasFindings(findings)) {
    console.log(`Gitee OpenAPI comparison found differences: ${formatSummary(findings)}`);
    console.log(formatFindings(findings));
  } else {
    console.log("Gitee OpenAPI comparison passed: no supported endpoint or parameter differences found.");
  }

  return hasFindings(findings) ? 1 : 0;
}

function compareParameter(left, right) {
  return `${left.in}:${left.name}`.localeCompare(`${right.in}:${right.name}`);
}

function compareFinding(left, right) {
  return `${left.code}:${left.path || ""}:${left.method || ""}:${left.message}`.localeCompare(
    `${right.code}:${right.path || ""}:${right.method || ""}:${right.message}`,
  );
}

function escapeRegExp(value) {
  return String(value).replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}

if (require.main === module) {
  main().then(
    (code) => {
      process.exitCode = code;
    },
    (error) => {
      const prefix = error instanceof DownloadError ? "Download error" : "Error";
      console.error(`${prefix}: ${error.message}`);
      process.exitCode = 2;
    },
  );
}

module.exports = {
  DEFAULT_CONFIG_PATH,
  DownloadError,
  UsageError,
  collectOpenApiOperations,
  compareApiDocs,
  createDefaultConfig,
  extractGoMetadata,
  fetchOfficialOpenApiJson,
  formatFindings,
  formatSummary,
  getLiveOpenApiUrl,
  loadOpenApiJson,
  loadConfig,
  normalizeGoPath,
  normalizeOpenApiPath,
  parseArgs,
  resolveDownloadButtonUrl,
};
