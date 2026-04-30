# `foundation-cli storage pool node list`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli storage pool node list --tenant-id <tenant-id> --storage-svc-id <storage-svc-id>` |
| 方法 | `GET` |
| 路径 | `/storageresmgm/v1/{storageSvcId}/pool/node` |
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
| `--count` | `100` | 分页数量 |

## 查询参数映射

| CLI Flag | 后端参数 |
|---|---|
| `--index` | `index` |
| `--count` | `count` |

## 返回结果关注点

常见字段：

| 字段 | 说明 |
|---|---|
| `responseData.data[].id` | 节点 ID，后续查询设备时需要 |
| `responseData.data[].name` | 节点名称 |
| `responseData.data[].ip` | 节点 IP |
| `responseData.data[].status` | 节点状态 |
| `responseData.data[].isUsed` | 是否已被使用 |

## 当前环境验证

- 已使用 `storage-svc-id=4235757746e447d18d4d7a8b3f5edc54` 实测成功。
- 当前环境返回节点 `babd3c3fbfd34359800ecffcb8907402`，名称 `anybackup`，IP `172.31.12.91`。
- 该接口签名遵循 SDK demo 的 forward 规则：请求路径保留 `{storageSvcId}`，签名 path 使用 `/storageresmgm/v1/pool/node`。

## 示例

```bash
foundation-cli storage pool node list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --storage-svc-id 52247d9c121745eab3a7880817457df5 \
  --index 0 \
  --count 100
```

## 使用建议

1. 该命令通常用于 `pool create` 前的资源发现。
2. 先从返回里挑出目标 `id`，再传给 `storage pool node device list --node-id`。

## 请求参数

请结合上文参数表传参；常见查询命令可按需传 `--index` 与 `--count` 控制分页。

## 返回案例

```json
{
    "error": null,
    "status": "success",
    "responseData": {
        "totalNum": 1,
        "data": [
            {
                "id": "babd3c3fbfd34359800ecffcb8907402",
                "name": "anybackup",
                "ip": "172.31.12.91",
                "status": 1,
                "isUsed": false
            }
        ]
    }
}
```
