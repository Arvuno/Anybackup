# `foundation-cli client deploy-config create`

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | `foundation-cli client deploy-config create --tenant-id <tenant-id> --os <linux|windows|unix> --data '<json>'` |
| Method | `POST` |
| Path | `/deploy/v1/hostConfig/{Linux|Windows|Unix}` |
| Risk | `write` |

## 请求参数

### 共享参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户标识 |
| endpoint | `--endpoint` | string | 是 | Foundation Console 地址 |
| ak | `--ak` | string | 是 | Access Key |
| sk | `--sk` | string | 是 | Secret Key |
| targetVersion | `--target-version` | string | 否 | 目标版本，默认 `9.0.9.0` |

### 业务参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| os | `--os` | string | 是 | `linux` / `windows` / `unix` |
| data | `--data` | json string | 是 | JSON 请求体 |

## Body 参数

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

- `--os=linux` -> `/deploy/v1/hostConfig/Linux`
- `--os=windows` -> `/deploy/v1/hostConfig/Windows`
- `--os=unix` -> `/deploy/v1/hostConfig/Unix`
- `userType` 常见值：`1` root 类型，`2` 普通用户类型

## 返回案例

```json
{
  "error": null,
  "status": "success",
  "responseData": [
    "dcd029549e4811f1a97d0050568943e2"
  ]
}
```

## 示例

```bash
foundation-cli client deploy-config create \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --os linux \
  --data '{"name":"ceshi","hostList":["115.190.186.186"],"port":22,"administrator":true,"account":{"rootPassword":"12345678"}}'
```
