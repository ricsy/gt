# gt

Gitee 命令行工具。

## 安装

### 使用 go install

```bash
go install github.com/ricsy/gt@latest
```

安装后需要将 `~/go/bin` 添加到 PATH：

```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc && source ~/.zshrc
```

验证安装：

```bash
gt --version
```

### 从源码编译

```bash
git clone https://github.com/ricsy/gt.git
cd gt
go build -o gt .
```

### 下载预编译版本

从 [GitHub Releases](https://github.com/ricsy/gt/releases) 下载对应平台的二进制文件。

## 配置

Gitee API Token 可从 [Gitee 个人设置](https://gitee.com/profile/personal_access_tokens) 生成。

设置环境变量：

```bash
export GITEE_TOKEN="your_token_here"
```

也可以基于仓库提供的 `.env.example` 创建本地环境文件，再通过 `--env-file` 或 `GT_ENV_FILE` 临时加载：

```bash
cp .env.example .env.local
gt --env-file .env.local repo list
```

或首次使用配置访问令牌：

```bash
gt auth login
```

## Codex 插件

仓库内置了一个 Codex 插件，插件目录位于 `plugins/gitee`。

从 GitHub 仓库来源安装：

```bash
codex plugin marketplace add ricsy/one-hive
codex plugin add gitee@one-hive
```

卸载插件：

```bash
codex plugin remove gitee@one-hive
```

查看已加载的 marketplace 和插件：

```bash
codex plugin marketplace list --json
codex plugin list --json
```
