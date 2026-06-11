#!/usr/bin/env node
const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

function getVersion() {
  const args = process.argv.slice(2).filter(arg => !arg.startsWith('--'));
  if (args.length === 0) {
    console.error('Usage: node release.js <version>');
    process.exit(1);
  }
  return args[0];
}

function updatePackageJson(version) {
  const pkgPath = path.join(__dirname, '..', 'package.json');
  const pkg = JSON.parse(fs.readFileSync(pkgPath, 'utf8'));
  pkg.version = version;
  fs.writeFileSync(pkgPath, JSON.stringify(pkg, null, 2) + '\n');
  console.log(`Updated package.json to ${version}`);
}

function updateRootGo(version) {
  const rootGoPath = path.join(__dirname, '..', 'internal', 'cmd', 'root.go');
  let content = fs.readFileSync(rootGoPath, 'utf8');
  content = content.replace(/var version = "[^"]*"/, `var version = "${version}"`);
  fs.writeFileSync(rootGoPath, content);
  console.log(`Updated root.go to ${version}`);
}

function updatePluginManifest(version) {
  const manifestPath = path.join(__dirname, '..', 'plugins', 'gitee', '.codex-plugin', 'plugin.json');
  const manifest = JSON.parse(fs.readFileSync(manifestPath, 'utf8'));
  manifest.version = version;
  fs.writeFileSync(manifestPath, JSON.stringify(manifest, null, 2) + '\n');
  console.log(`Updated plugin.json to ${version}`);
}

function gitCommit(version) {
  const status = execSync('git status --porcelain').toString().trim();
  if (status) {
    execSync('git add package.json internal/cmd/root.go plugins/gitee/.codex-plugin/plugin.json', { stdio: 'inherit' });
    execSync(`git commit -m "release: bump version to ${version}"`, { stdio: 'inherit' });
    console.log(`Committed version bump to ${version}`);
  } else {
    console.log('No changes to commit, using existing commit');
  }
}

function gitTag(version) {
  const tag = version.startsWith('v') ? version : `v${version}`;
  execSync(`git tag -a ${tag} -m "Release ${tag}"`, { stdio: 'inherit' });
  execSync(`git push origin ${tag}`, { stdio: 'inherit' });
  console.log(`Pushed tag ${tag}`);
}

function run() {
  const version = getVersion();
  updatePackageJson(version);
  updateRootGo(version);
  updatePluginManifest(version);
  gitCommit(version);
  gitTag(version);
  console.log(`\nRelease ${version} created successfully!`);
}

run();
