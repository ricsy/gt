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
