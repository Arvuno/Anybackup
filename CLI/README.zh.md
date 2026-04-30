# foundation-cli

`foundation-cli` 是面向 Foundation 控制面的统一命令行工具，提供保护对象、策略、作业、客户端、存储、网络和应用专属能力的标准化入口，方便人工运维、脚本集成和 Agent 编排调用。

## 特性
- 统一入口：将分散在多个后端接口上的控制能力收敛到一个 CLI。
- AK/SK 认证：统一使用 `endpoint + tenant-id + ak + sk` 发起签名请求。
- 面向 Agent：返回结构稳定，适合脚本、自动化任务和上层 Agent 调用。
- 业务域清晰：按 `protect`、`mysql`、`vmware`、`host`、`job`、`client` 等域组织命令。
- 原始结果友好：默认保留后端 `status/error/responseData` 结构，便于继续编排。

## 功能列表

| 类别 | 能力 |
|---|---|
| 客户端管理 | 部署客户端代理、查看主机列表、部署作业、部署作业执行输出、查看客户端和代理列表 |
| 作业管理 | 停止作业、删除作业、查看作业列表、备份作业详情、作业执行输出、子作业列表 |
| MYSQL | MySQL 对象和数据源查询、配置备份配置、发起恢复和备份、备份配置/恢复时间点/恢复和备份作业配置详情查看、授权 |
| 策略管理 | 创建、删除备份策略、查看策略列表 |
| 存储池管理 | 创建、删除存储池、获取存储池列表 |
| 时间点管理 | 清理时间点、查看时间点列表 |

## 安装

### 环境要求
- Go 1.22 或更高版本
- 可访问目标 Foundation 控制面
- 有效的 `ak`、`sk` 与目标 `endpoint`（`tenant-id` 可选）

### 从源码构建
Windows:

```powershell
go build -o foundation-cli.exe .\cmd\foundation-cli
```

macOS / Linux:

```bash
go build -o foundation-cli ./cmd/foundation-cli
```

### 查看版本

```bash
foundation-cli version
```

当前源码默认支持的目标版本是 `9.0.9.0`。

## 快速开始

### 1. 查看根命令帮助
```bash
foundation-cli --help
```

### 2. 执行一个只读命令
```bash
foundation-cli job logs \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --job-id <job-id> \
  --index 0 \
  --count 30
```

### 3. 使用透传 API 调试未标准化能力
```bash
foundation-cli api \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --path "/job_center/v1/business_types"
```

## 全局参数

大多数业务命令都会继承这组远端访问参数：

| 参数 | 说明 |
|---|---|
| `--tenant-id` | 租户 ID（可选） |
| `--endpoint` | Foundation 控制面地址 |
| `--ak` | Access Key |
| `--sk` | Secret Key |
| `--target-version` | 目标版本，默认 `9.0.9.0` |

## 命令总览

当前源码中的一级命令域如下：

| 命令域 | 说明 |
|---|---|
| `protect` | 通用保护运营能力：发起备份 |
| `policy` | 创建、删除备份策略、查看策略列表 |
| `timepoint` | 清理时间点、查看时间点列表 |
| `job` | 作业列表、详情、日志、子作业列表等运行态观测能力 |
| `client` | 部署客户端代理、查看主机列表、部署作业、部署作业执行输出、查看客户端和代理列表 |
| `mysql` | MySQL 对象和数据源查询、配置备份配置、发起恢复、备份配置/恢复时间点/恢复和备份作业配置详情查看、授权 |
| `network` | 网络子网与节点查询能力 |
| `storage` | 创建、删除存储池、获取存储池列表 |

你可以继续用 `foundation-cli <domain> --help` 查看某个命令域的具体子命令。

## 常见使用场景

### 查询作业列表
```bash
foundation-cli job list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk>
```

### 查询子作业列表
```bash
foundation-cli job child list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --job-id <job-id> \
  --index 0 \
  --count 30
```

### 查询 MySQL 时间点
```bash
foundation-cli mysql timepoint list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --object-id <object-id>
```

### 查询存储服务
```bash
foundation-cli storage service list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk>
```

## 输出与错误处理
- 成功时通常返回 `status=success`。
- 失败时保留后端 `errorCode / errorArgs`。
- 命令退出码由 CLI 统一管理，适合脚本判断执行结果。

示例：

```json
{
  "status": "success",
  "error": null,
  "responseData": {
    "data": [],
    "totalNum": 0
  }
}
```

## 项目结构

```text
cmd/        Cobra 入口与根命令
internal/   业务命令、请求构造、签名与输出实现
skill/      面向 Agent 的命令路由与单命令文档
scripts/    校验与辅助脚本
tests/      单元测试
docs/       架构、设计和参考资料
```

## Agent / Skill 文档

如果你要让 Agent 更稳定地调用本工具，建议直接查看 `skill/` 目录下的入口文档，例如：

- [skill/job/SKILL.md](./skill/job/SKILL.md)
- [skill/mysql/SKILL.md](./skill/mysql/SKILL.md)
- [skill/protect/SKILL.md](./skill/protect/SKILL.md)

这些文档更适合 Agent 做命令选择、参数补齐和单命令跳转。

## 开发

### 运行测试

```bash
go test ./...
```

如果你只想验证某个命令域，推荐跑聚焦测试，例如：

```bash
go test ./tests/unit/business/domains -run '^TestJobReadMappings$' -count=1
```

### 文档校验

针对单命令文档，可用 `doc-checker` 脚本做结构检查：

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\doc-checker\references\check-command-doc.ps1 -Path .\skill\job\references\commands\child-list.md
```

## 参考资料

- [Foundation CLI 管理简报](./docs/references/2026-04-08-foundation-cli-management-brief.md)
