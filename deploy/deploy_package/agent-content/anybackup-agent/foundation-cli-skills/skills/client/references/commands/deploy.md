# `foundation-cli client deploy`

## 命令摘要

客户端部署复合命令：先创建 `deploy-config`，再将创建结果中的配置 ID 注入到部署请求体 `hostList`，最后发起部署作业。

| 项 | 值 |
|---|---|
| CLI | `foundation-cli client deploy --tenant-id <tenant-id> --os <linux|windows|unix> --data '<json>'` |
| Method | `POST`（复合：`/deploy/v1/hostConfig/{OS}` -> `/deploy/v1/job/config`） |
| Risk | `write` |

## 请求参数

### 共享参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户标识 |
| endpoint | `--endpoint` | string | 是 | Foundation Console 地址 |
| ak | `--ak` | string | 是 | Access Key |
| sk | `--sk` | string | 是 | Secret Key |
| targetVersion | `--target-version` | string | 否 | 目标版本，默认 `8.0.9.0` |

### 命令参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| os | `--os` | string | 是 | `linux` / `windows` / `unix` |
| data | `--data` | json string | 是 | 复合请求体 |

## Body 参数

### `/deploy/v1/job/config`（部署阶段）

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| name / jobName | string | 是（二选一） | 部署任务名称 |
| LinuxInstallPath | string | 否 | Linux 安装路径 |
| WindowsInstallPath | string | 否 | Windows 安装路径 |
| concurrent | int | 否 | 并发数 |
| uninstallType | bool | 否 | 是否卸载已安装客户端 |
| hostList | string[] | 是 | 部署目标列表（复合命令中会被 deploy-config 返回 ID 覆盖） |
| runners | object[] | 否 | Runner 列表 |
| runners[].runnerType | string | 否 | Runner 类型，枚举见下方“runnerType 枚举（部署）”；必须由用户明确确认后输入 |

### `/deploy/v1/hostConfig/{OS}`（创建配置阶段）

### `--os linux`

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| name | string | 是 | 配置名称 |
| hostList | string[] | 是 | 主机列表 |
| port | int | 是 | SSH 端口 |
| administrator | bool | 是 | 是否管理员模式 |
| account.rootPassword | string | 是 | Root 用户密码 |

Linux 主机配置推荐请求体（已验证）：

```json
{
  "name": "ceshi",
  "hostList": ["115.190.186.186"],
  "port": 22,
  "administrator": true,
  "account": {
    "rootPassword": "12345678"
  }
}
```

### `--os windows`

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| name | string | 是 | 配置名称 |
| hostList | string[] | 是 | 主机列表 |
| account.user | string | 是 | 登录用户名 |
| account.password | string | 是 | 登录密码 |

### `--os unix`

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| name | string | 是 | 配置名称 |
| hostList | string[] | 是 | 主机列表 |
| port | int | 是 | SSH 端口 |
| administrator | bool | 是 | 是否管理员模式 |
| userType | int | 条件必填 | `administrator=false` 时可传 |
| account.rootPassword | string | 条件必填 | `administrator=true` 时必填；`administrator=false` 时部分场景也会携带 |
| account.user | string | 条件必填 | `administrator=false` 时必填 |
| account.password | string | 条件必填 | `administrator=false` 时必填 |

## 枚举列表

- `--os`: `linux` / `windows` / `unix`
- `hostList` 在复合流程中会被 `deploy-config create` 返回的配置 ID 覆盖

### runnerType 枚举（部署）

来源：[`runner_deployment.json`](./runner_deployment.json)（以该文件为准）。

| runnerType | 中文名 |
|---|---|
| `Gen` | 通用代理 |
| `Basic` | 基础代理 |
| `Machine` | 机器系统保护代理 |
| `File` | 文件/对象代理 |
| `Volume` | 卷代理 |
| `ProgramExecutor` | 脚本代理 |
| `MySQL` | MySQL 代理 |
| `Oracle` | Oracle 代理 |
| `VMware` | VMware 代理 |
| `HyperV` | Hyper-V 代理 |

说明：
- 以上为部署场景常用枚举。
- `runners[].runnerType` 必须由用户确认输入，不允许 Agent 默认猜测或自动替换。
- 若用户未明确 runnerType，先让用户从枚举中确认后再执行部署命令。

## 执行流程

1. `GET /deploy/v1/hostConfig/nameExist?name=<name>`
2. `POST /deploy/v1/hostConfig/{Linux|Windows|Unix}`
3. 从创建响应读取配置 ID 数组
4. 将配置 ID 注入部署请求体 `hostList`
5. `GET /deploy/v1/job/config/nameExist?name=<name>`
6. `POST /deploy/v1/job/config`

## 返回案例

```json
{
  "error": null,
  "status": "success",
  "responseData": {
    "jobId": "sample-job-id"
  }
}
```

## 示例

```bash
foundation-cli client deploy \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --os linux \
  --data '{"name":"deploy-job","hostList":["115.190.186.186"],"port":22,"administrator":true,"account":{"rootPassword":"12345678"},"runners":[{"runnerType":"Gen"}]}'
```
