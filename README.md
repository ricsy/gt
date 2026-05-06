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

或首次使用配置访问令牌：

```bash
gt auth login
```

## 使用

```bash
# 查看仓库
gt repo view owner/repo

# 管理 Issue
gt issue list --owner owner --repo repo
gt issue create --owner owner --repo repo --title "标题" --body "内容"

# 管理 PR
gt pr list --repo owner/repo
gt pr create --repo owner/repo --title "标题" --head branch --base main

# 管理 Release
gt release list --repo owner/repo
gt release create --repo owner/repo --name "v1.0.0" --body "更新内容" tag
```
