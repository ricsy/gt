# AGENTS.md

本文档用于向 Claude Code（claude.ai/code）说明在本仓库中工作的约束与方式。

## API 第一手资料优先

**当任务涉及新增功能、查看已有 API、核对接口定义、确认请求参数/响应结构、排查 API 行为时，必须优先查看
`swagger-gitee` MCP 或 `data/api-1.json`，不要先从互联网搜索开始。**

这是本仓库里关于 Gitee Open API v5 的第一手资料来源，优先级高于普通网页搜索、二手博客、论坛回答和印象式判断。
只有在这两处资料无法覆盖问题，或者需要核对上游官网补充说明时，才允许再扩大到外部搜索。

## 工具优先级

**代码操作优先使用 Goland MCP（`mcp__goland__*`）。** 读取、搜索、替换、跳转时优先走 MCP。只有在
MCP 不足以完成任务时，才使用直接文件工具（Read、Edit、Write）。

## API 参考

**MCP `swagger-gitee`**：直接查询 Gitee Open API v5 接口定义。

- `mcp__swagger-gitee__list_api_groups`：列出所有 API 分组
- `mcp__swagger-gitee__search_apis`：按关键字搜索 API
- `mcp__swagger-gitee__get_api_detail`：获取完整 API 定义（路径、方法、参数、Schema）

`data/api-1.json` 是同一份规范的本地副本（来源于 https://help.gitee.com/openapi/v5）。

## 项目概览

`gt` 是一个面向 Gitee（gitee.com）的 CLI 工具。它通过 Gitee API v5 管理仓库、Issue、PR、Release
和组织等能力。

## 常用命令

```bash
go build                          # 构建二进制（输出 gt/gt.exe）
go test ./...                     # 运行全部测试
go test -v ./internal/cmd         # 运行 cmd 包测试
go test -v ./pkg/api              # 运行 API 包测试
golangci-lint run ./...           # 运行全部 lint
pnpm run test:integration         # 运行集成测试（需要先构建 CLI）
pnpm run release:local            # 本地发布构建
```

## 架构

```text
internal/cmd/    # CLI 命令实现（Cobra 命令）
pkg/api/         # Gitee API Client（HTTP 调用与响应类型）
pkg/auth/        # 认证（Token 存储与缓存）
pkg/config/      # 配置文件与 hosts 文件管理
pkg/util/        # 共享工具
```

### API Client（`pkg/api/client.go`）

`Client` 结构体封装了对 Gitee API v5 的 HTTP 调用。所有 API 路径都必须通过
`config.ApiUrl(host)` 生成，不允许硬编码。认证头格式为 `Authorization: token <token>`。

每个领域（issue、pr、repo、org、release）都在 `pkg/api/` 下拥有独立文件。

### 认证流程（`pkg/auth/auth.go`）

`GetToken(host)` 会先读取 `GITEE_TOKEN` 环境变量；如果没有，再回退到 `config.LoadHosts()`
（YAML 配置文件）。结果会使用双检锁策略做内存缓存。

### 仓库解析（`internal/cmd/repo_helper.go`）

获取已认证 API Client 必须使用 `getClient()`。解析仓库必须使用 `resolveRepoFlag(flag)`
，其解析顺序为：命令 flag -> `GT_REPO` 环境变量 -> `default_repo` 配置。解析 `owner/repo`
字符串时使用 `ResolveRepo(owner/repo)`。

## 关键模式

- 所有命令处理函数都使用 `getClient()` + `resolveRepoFlag()`，不要直接调用
  `auth.GetToken` + `api.NewClient`
- URL 构造统一通过 `config.ApiUrl(host)`，不要硬编码 `https://gitee.com`
- Issue / PR 状态过滤统一使用 `api.StateOpen` 常量（值为 `"open"`）

## 示例数据规则

- 不要在产品代码、测试、夹具或示例中写死开发者个人仓库、组织或 namespace。
- 当必须使用稳定示例 owner 或 namespace 时，优先使用官方公开标识，例如 `gitee`。
- 如果仓库自身身份信息必须出现（例如 `go.mod`、import path、release URL、license
  元数据），应将其视为项目元数据，而不是可复用示例数据；不要把这些值继续复制到新的夹具或示例里。

## 变更管理规则

- 任何新增功能或行为变更，在视为完成之前，都必须先做一轮对仓库文档与插件工作流的影响评估。
- 如果变更影响命令用法、认证流程、配置、测试、发布流程，或 agent / plugin
  行为，必须在同一次交付中同步更新相关文档、示例、夹具和插件工作流说明。
- 不要把“代码改完”视为“工作完成”；在结束前必须显式检查 `README`、`AGENTS.md`
  、测试夹具、集成流程，以及插件相关的配置或使用说明是否需要更新。
- 对本仓库而言，新增或变更功能时，必须显式检查 `plugins/gitee` 以及相关插件工作流文档、示例是否需要同步更新。

## API 实现检查清单

实现或更新 API 接口时，必须对照 `data/api-1.json` 或 `swagger-gitee` MCP 进行核验：

### 1. Endpoint 定义（`pkg/api/endpoint.go`）

- [ ] 路径与 swagger 一致，例如：`/v5/repos/{owner}/{repo}/issues`
- [ ] HTTP 方法与 swagger 一致：GET、POST、PATCH、DELETE、PUT
- [ ] `EndpointGroup` 中的 action 名称正确

### 2. 请求参数（`pkg/api/response/{domain}.go`）

对 swagger 中每个 query/path/form 参数：

- [ ] 参数名完全一致（区分大小写）
- [ ] 参数类型一致：string、int、bool、array
- [ ] 必填 / 可选定义正确
- [ ] 如果 swagger 定义了枚举值，需正确校验

### 3. 请求体（`pkg/api/response/{domain}.go`）

对于 POST / PATCH / PUT 请求：

- [ ] 所有表单字段都在请求结构体中定义
- [ ] `json` tag 与 swagger 字段名完全一致
- [ ] 可选字段正确使用 `omitempty`
- [ ] 字段类型与 swagger 一致

### 4. 响应结构（`pkg/api/response/{domain}.go`）

对 swagger 响应定义中的每个字段：

- [ ] 字段名完全一致（区分大小写）
- [ ] `json` tag 与 swagger 一致
- [ ] 字段类型正确：string、int、bool、array、object
- [ ] 嵌套对象有正确的结构体定义
- [ ] 所有可选字段使用了合适的指针类型或 `omitempty`

### 5. 验证步骤

```bash
# 1. 搜索目标接口
mcp__swagger-gitee__search_apis tag:"Issues"

# 2. 获取完整 API 定义
mcp__swagger-gitee__get_api_detail --path "/v5/repos/{owner}/{repo}/issues" --method "get"

# 3. 获取响应 schema
mcp__swagger-gitee__get_schema --ref "Issue"

# 4. 验证构建通过
go build ./...

# 5. 验证 lint 通过
golangci-lint run ./...
```

## 新模块开发工作流

当实现新的 API 模块（例如 milestone、label）时，遵循以下目录结构：

### 文件位置

| 用途            | 文件路径                                                         |
|---------------|--------------------------------------------------------------|
| Endpoint 定义   | `pkg/api/endpoint.go`：向 `EndpointGroup` 结构体和 `var XXX` 中补充定义 |
| 请求 / 响应类型     | `pkg/api/response/{domain}.go`：定义 API 响应与参数结构                |
| API Client 方法 | `pkg/api/{domain}.go`：通过 `DoFromEndpoint` 实现类型别名与调用          |
| CLI 命令        | `internal/cmd/{domain}.go`：Cobra 命令实现                        |
| CLI 测试        | `internal/cmd/{domain}_test.go`：命令与 flag 的单元测试               |

### 开发步骤

1. **向 `pkg/api/endpoint.go` 添加 endpoint**
    - 如有必要，向 `EndpointGroup` 结构体增加新字段（通常很少，绝大多数可复用
      List/Get/Create/Update/Delete）
    - 增加 `var XXX = EndpointGroup{...}`，补齐对应路径

2. **在 `pkg/api/response/{domain}.go` 定义类型**
    - 响应结构体（例如 `Milestone`、`License`）
    - 参数结构体（例如 `ListXxxOptions`、`CreateXxxOptions`）
    - `json:"field_name"` 必须与 swagger 完全一致

3. **在 `pkg/api/{domain}.go` 实现 API 方法**
    - 使用类型别名，例如：`type ListXxxOptions = response.ListXxxOptions`
    - 所有 API 调用统一使用 `DoFromEndpoint()`
    - 返回类型统一来自 `response` 包

4. **在 `internal/cmd/{domain}.go` 实现 CLI**
    - 参考现有模式（例如 `issue.go`、`webhook.go`）
    - 使用 `newXxxCmd()` 模式创建命令
    - 合理组织 flags，并通过 `resolveRepoFlag()` 处理仓库解析

5. **在 `internal/cmd/{domain}_test.go` 添加测试**
    - 测试命令结构（`Use`、子命令数量等）
    - 测试 flags 是否存在
    - 可以参考 `webhook_test.go` 的测试模式

6. **对照上面的 API 实现检查清单完成验证**

## 常见错误

- **不要**在 `pkg/api/{domain}.go` 中定义类型；类型应该放在
  `pkg/api/response/{domain}.go`
- **不要**直接使用 `Do()`；统一通过带 endpoint 定义的 `DoFromEndpoint()`
- **不要**手写 query string；统一使用 `util.BuildQuery()`
- **不要**在 `os.Exit()` 前自己打印错误；Cobra 会统一处理错误输出
