# `foundation-cli storage pool create`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli storage pool create --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> --body-file <path>` |
| 方法 | `POST` |
| 路径 | `/storageresmgm/v1/{storageSvcId}/pool` |
| 风险 | `写入` |

## 必填输入

| 参数 | 必填 | 说明 |
|---|---|---|
| `--tenant-id` | 是 | 租户上下文 |
| `--endpoint` | 是 | Foundation 接口地址 |
| `--ak` | 是 | 访问密钥 |
| `--sk` | 是 | 访问密钥密文 |
| `--storage-svc-id` | 是 | 存储服务 ID |
| `--data` / `--body-file` | 是 | 请求体 JSON |

## 请求体约束

1. 顶层必须包含 `name`、`type`、`redundancyPolicy`、`resource`。
2. `redundancyPolicy.policy` 必须存在。
3. `warnThreshold` 如果存在，取值范围必须是 `1..98`。
4. `safeThreshold` 如果存在，取值范围必须是 `2..99`。
5. 更深的资源类型组合规则仍由后端校验。

## 类型枚举

| `type` | 含义 | 常用资源字段 |
|---|---|---|
| `1` | SAN 存储池 | `resource[].device` |
| `2` | 本地块存储池 | `resource[].device` |
| `3` | 分布式存储池 | `resource[].device`，通常同时带 `readCacheDevices`、`writeCacheDevices`、`mainDevices` |
| `4` | NAS 存储池 | `resource[].mountPath` |
| `5` | 磁带存储池 | `resource[].tape` |
| `7` | 本地文件系统存储池 | `resource[].mountPath` |
| `8` | 备份对象存储池 | `resource[].objectStorage` |

## 冗余策略枚举

| `redundancyPolicy.policy` | 含义 | 常用附加字段 |
|---|---|---|
| `1` | 副本 | 可带 `replicationNum` |
| `2` | 纠删码 | 可带 `dataUnitsNum`、`parityUnitsNum` |

## 最小通用模板

```json
{
  "name": "pool-a",
  "type": 3,
  "redundancyPolicy": {
    "policy": 1
  },
  "resource": []
}
```

这个模板只适合说明 CLI 的最低校验门槛，不保证后端一定接受。真正调用时，必须根据 `type` 补齐对应的资源字段。

## 资源字段选择规则

1. `type=1`、`type=2`、`type=3` 时，优先准备 `resource[].device`。
2. `type=4`、`type=7` 时，优先准备 `resource[].mountPath`。
3. `type=5` 时，准备 `resource[].tape`。
4. `type=8` 时，准备 `resource[].objectStorage`。

## 常见示例

### 示例 1：分布式存储池

```json
{
  "name": "dist-pool-a",
  "type": 3,
  "redundancyPolicy": {
    "policy": 1,
    "replicationNum": 1
  },
  "warnThreshold": 90,
  "safeThreshold": 95,
  "resource": [
    {
      "device": [
        {
          "nodeID": "node-1",
          "readCacheDevices": [
            "read-cache-dev-1"
          ],
          "writeCacheDevices": [
            "write-cache-dev-1"
          ],
          "mainDevices": [
            "main-dev-1",
            "main-dev-2"
          ]
        }
      ]
    }
  ]
}
```

适用场景：

1. 先通过 `storage pool node list` 找节点。
2. 再通过 `storage pool node device list` 获取该节点可用设备。
3. 最后把设备 ID 填入 `readCacheDevices`、`writeCacheDevices`、`mainDevices`。

### 示例 2：NAS 存储池

```json
{
  "name": "nas-pool-a",
  "type": 4,
  "redundancyPolicy": {
    "policy": 1
  },
  "warnThreshold": 85,
  "safeThreshold": 95,
  "resource": [
    {
      "mountPath": [
        {
          "nodeID": "node-1",
          "mountPaths": [
            "/mnt/nas-share-a"
          ]
        }
      ]
    }
  ]
}
```

适用场景：

1. `mountPaths` 表达的是挂载路径列表。
2. `type=7` 的本地文件系统存储池可参考同一结构，只是资源来源不同。

### 示例 3：备份对象存储池

```json
{
  "name": "object-pool-a",
  "type": 8,
  "redundancyPolicy": {
    "policy": 1
  },
  "resource": [
    {
      "objectStorage": {
        "acc": {
          "warnThreshold": 90,
          "backup": [
            {
              "nodeID": "node-1",
              "mountPaths": [
                "/mnt/object-cache-a"
              ]
            }
          ]
        },
        "objectID": "object-storage-1",
        "bucketID": "bucket-1"
      }
    }
  ]
}
```

适用场景：

1. `objectID` 表示对象存储资源 ID。
2. `bucketID` 表示目标桶 ID。
3. `acc.backup[].mountPaths` 表示对象存储加速介质使用的挂载路径。

### 示例 4：创建本地块存储池（type=2）

```json
{
  "name": "sdd",
  "safeThreshold": 95,
  "warnThreshold": 90,
  "redundancyPolicy": {
    "policy": 1,
    "replicationNum": 1
  },
  "type": 2,
  "resource": {
    "device": [
      {
        "nodeID": "6852e9be0e8d4382bd4cd78a3ea0ee66",
        "mainDevices": [
          "ffa924c841e411f1a9de005056891898"
        ],
        "writeCacheDevices": []
      }
    ]
  }
}
```

## 推荐调用方式

请求体较长时，优先使用 `--body-file`：

```bash
foundation-cli storage pool create \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --storage-svc-id <storage-svc-id> \
  --body-file pool-create.json
```

该接口签名遵循 SDK demo 的 forward 规则：请求路径保留 `{storageSvcId}`，签名 path 使用 `/storageresmgm/v1/pool`。

## 请求参数

请结合上文参数表传参；常见查询命令可按需传 `--index` 与 `--count` 控制分页。

## Body 参数

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| data | object | 是 | 请求体对象，建议按“请求体示例”中的字段组织 |

## 枚举列表

- 枚举值请以上文参数说明为准；未列出的取值按后端接口定义传入。

## 请求体示例

```json
{}
```

## 返回案例

```json
{"status":"success","error":null,"responseData":null}
```

## 示例

```bash
foundation-cli storage pool create
```

## Agent 创建前置决策流程（必须）

在调用 `foundation-cli storage pool create` 之前，Agent 必须按以下顺序完成用户确认与参数选择：

1. 先调用 `storage service list`，让用户先选择存储服务，并确认 `storage-svc-id`。
2. 再让用户确认存储池类型（`type`）。
3. 再让用户确认冗余策略（`redundancyPolicy.policy` 及对应附加参数）。
4. 再让用户确认介质类型（对应 `resource` 结构：`device/mountPath/tape/objectStorage`）。
5. 先调用 `storage pool node list`，让用户选择节点（`nodeID`）。
6. 再调用 `storage pool node device list`，让用户选择设备 ID。
7. 最后组装 `resource` 并执行 `storage pool create`。

硬约束：

- 未完成“存储服务 + 类型 + 冗余策略 + 介质类型 + 节点 + 设备”六类确认时，禁止发起创建请求。
- 设备必须来自 `storage pool node device list` 的真实返回，禁止猜测设备 ID。
- 创建设备型存储池时，优先使用可创建状态的设备（详见 `pool-node-device-list.md` 的设备状态规则）。
- 创建失败时，必须回传后端错误并提示用户重新确认上述 6 类输入。
