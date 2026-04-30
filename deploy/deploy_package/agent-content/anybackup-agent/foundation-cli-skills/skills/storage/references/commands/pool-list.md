# `foundation-cli storage pool list`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli storage pool list --tenant-id <tenant-id> --storage-svc-id <storage-svc-id>` |
| 方法 | `GET` |
| 路径 | `/storageresmgm/v1/{storageSvcId}/pool` |
| 风险 | `只读` |

## 必填输入

| 参数 | 必填 | 说明 |
|---|---|---|
| `--tenant-id` | 是 | 租户上下文 |
| `--endpoint` | 是 | Foundation 接口地址 |
| `--ak` | 是 | 访问密钥 |
| `--sk` | 是 | 访问密钥密文 |
| `--storage-svc-id` | 是 | 存储服务 ID |

## 可选输入

| 参数 | 默认值 | 说明 |
|---|---|---|
| `--index` | `0` | 分页起始位置 |
| `--count` | `30` | 分页数量 |

## 查询参数映射

| CLI Flag | 后端参数 |
|---|---|
| `--index` | `index` |
| `--count` | `count` |

## 返回结果关注点

常见字段：

| 字段 | 说明 |
|---|---|
| `responseData.data[].id` | 存储池 ID，后续删除等动作会用到 |
| `responseData.data[].name` | 存储池名称 |
| `responseData.data[].type` | 存储池类型 |
| `responseData.data[].status` | 存储池状态 |
| `responseData.data[].totalSize` | 总容量 |
| `responseData.data[].usedSize` | 已使用容量 |
| `responseData.data[].availableSize` | 可用容量 |
| `responseData.data[].warnThreshold` | 告警阈值 |
| `responseData.data[].safeThreshold` | 安全阈值 |

## 当前环境验证

- 已使用 `storage-svc-id=4235757746e447d18d4d7a8b3f5edc54` 实测成功。
- 当前环境返回了存储池 `block`，`poolId=0c6022822f1d11f19b0600163e33df73`。
- 该接口签名遵循 SDK demo 的 forward 规则：请求路径保留 `{storageSvcId}`，签名 path 使用 `/storageresmgm/v1/pool`。

## 返回结果

CLI 标准返回结构：

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 执行状态 |
| error | object/null | 错误信息（失败时） |
| responseData | object | 接口业务返回 |
| responseData.totalNum | int | 存储池总数 |
| responseData.data | object[] | 存储池列表 |

## 返回案例

```json
{
  "error": null,
  "status": "success",
  "responseData": {
    "totalNum": 1,
    "data": [
      {
        "id": "e177c3fe311f11f1a1220050568943e2",
        "name": "130sdb",
        "type": 2,
        "status": 1,
        "redundancyPolicy": {
          "policy": 1,
          "replicationNum": 1,
          "dataUnitsNum": -1,
          "parityUnitsNum": -1
        },
        "totalSize": 53687091200,
        "usedSize": 214958080,
        "availableSize": 53472133120,
        "warnThreshold": 90,
        "safeThreshold": 95,
        "accNodeCapacityStatus": {
          "normalNum": 1,
          "warnNum": 0,
          "fullNum": 0
        },
        "version": 1,
        "mainMediaType": 1,
        "partSize": -1
      }
    ]
  }
}
```

## 示例

```bash
foundation-cli storage pool list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --storage-svc-id 52247d9c121745eab3a7880817457df5 \
  --index 0 \
  --count 30
```

## 使用建议

1. 删除存储池前，先用该命令确认目标 `id`。
2. 如果只是给 agent 做后续编排，优先关注 `id`、`name`、`status` 和容量字段。

## 请求参数

请结合上文参数表传参；常见查询命令可按需传 `--index` 与 `--count` 控制分页。
