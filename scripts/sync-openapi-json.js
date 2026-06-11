#!/usr/bin/env node

const fs = require("node:fs");
const path = require("node:path");

const {
  DEFAULT_CONFIG_PATH,
  DownloadError,
  fetchOfficialOpenApiJson,
  loadConfig,
  UsageError,
} = require("./compare-api-docs");

function parseArgs(argv) {
  const options = {
    configPath: DEFAULT_CONFIG_PATH,
    root: process.cwd(),
  };

  for (let i = 0; i < argv.length; i++) {
    const arg = argv[i];
    if (arg === "--config") {
      options.configPath = argv[++i];
      if (!options.configPath) {
        throw new UsageError("--config requires a file path");
      }
    } else if (arg === "--root") {
      options.root = argv[++i];
      if (!options.root) {
        throw new UsageError("--root requires a directory path");
      }
    } else if (arg === "--output") {
      options.outputPath = argv[++i];
      if (!options.outputPath) {
        throw new UsageError("--output requires a file path");
      }
    } else if (arg === "--help" || arg === "-h") {
      options.help = true;
    } else {
      throw new UsageError(`Unknown argument: ${arg}`);
    }
  }

  return options;
}

function helpText() {
  return [
    "Usage: node scripts/sync-openapi-json.js [options]",
    "",
    "Options:",
    `  --config <path>   Comparison config JSON (default: ${DEFAULT_CONFIG_PATH})`,
    "  --root <path>     Repository root (default: current working directory)",
    "  --output <path>   Override sync target path (default: config openapi.localPath)",
    "  -h, --help        Show this help",
  ].join("\n");
}

function formatOpenApiJson(openApi) {
  return `${JSON.stringify(openApi, null, 2)}\n`;
}

function writeIfChanged(filePath, content) {
  fs.mkdirSync(path.dirname(filePath), { recursive: true });

  const previous = fs.existsSync(filePath) ? fs.readFileSync(filePath, "utf8") : null;
  if (previous === content) {
    return false;
  }

  fs.writeFileSync(filePath, content);
  return true;
}

async function syncOpenApiJson(options = {}) {
  const rootDir = path.resolve(options.rootDir || process.cwd());
  const configPath = path.resolve(rootDir, options.configPath || DEFAULT_CONFIG_PATH);
  const config = (options.loadConfig || loadConfig)(configPath);
  const outputPath = path.resolve(rootDir, options.outputPath || config.openapi.localPath);
  const openApi = await (options.fetchOpenApiJson || fetchOfficialOpenApiJson)(config);
  const content = formatOpenApiJson(openApi);
  const updated = writeIfChanged(outputPath, content);

  return {
    updated,
    outputPath,
    content,
  };
}

async function main(argv = process.argv.slice(2)) {
  const options = parseArgs(argv);
  if (options.help) {
    console.log(helpText());
    return 0;
  }

  const result = await syncOpenApiJson({
    rootDir: options.root,
    configPath: options.configPath,
    outputPath: options.outputPath,
  });

  if (result.updated) {
    console.log(`Synced OpenAPI JSON to ${path.relative(process.cwd(), result.outputPath)}`);
  } else {
    console.log(`OpenAPI JSON already up to date: ${path.relative(process.cwd(), result.outputPath)}`);
  }

  return 0;
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
  formatOpenApiJson,
  helpText,
  parseArgs,
  syncOpenApiJson,
  writeIfChanged,
};
