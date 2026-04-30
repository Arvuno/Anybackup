# `foundation-cli storage pool node device list`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli storage pool node device list --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> --node-id <node-id> [--types <type> ...]` |
| 方法 | `GET` |
| 路径 | `/storageresmgm/v1/{storageSvcId}/device` |
| 风险 | `只读` |

## 必填输入

| 参数 | 必填 | 说明 |
|---|---|---|
| `--tenant-id` | 是 | 租户上下文 |
| `--endpoint` | 是 | Foundation 接口地址 |
| `--ak` | 是 | 访问密钥 |
| `--sk` | 是 | 访问密钥密文 |
| `--storage-svc-id` | 是 | 存储服务 ID |
| `--node-id` | 是 | 节点 ID |

## 可选输入

| 参数 | 默认值 | 说明 |
|---|---|---|
| `--types` | `1,2,3,4,5,7` | 设备类型筛选，支持重复传参 |
| `--authorized` | `2` | 授权状态筛选 |

## 查询参数映射

| CLI Flag | 后端参数 |
|---|---|
| `--node-id` | `nodeId` |
| `--types` | `types` |
| `--authorized` | `authorized` |

## 参数说明

| 参数 | 典型含义 |
|---|---|
| `--types` | 设备类型筛选，支持重复传参；默认展开为 `1,2,3,4,5,7` |
| `--authorized` | 授权状态筛选，默认沿用当前接口主路径的 `2` |

## 返回结果关注点

常见字段：

| 字段 | 说明 |
|---|---|
| `responseData.data[].id` | 设备 ID，可作为 `pool create` 请求体中的设备标识 |
| `responseData.data[].name` | 设备名称 |
| `responseData.data[].ip` | 设备所属节点 IP |
| `responseData.data[].status` | 设备状态 |
| `responseData.data[].isUsed` | 设备是否已被占用 |

## 示例

```bash
foundation-cli storage pool node device list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --storage-svc-id 52247d9c121745eab3a7880817457df5 \
  --node-id 839b4b217e214c348152221288257773
```

如果需要覆盖默认筛选条件：

```bash
foundation-cli storage pool node device list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --storage-svc-id 52247d9c121745eab3a7880817457df5 \
  --node-id 839b4b217e214c348152221288257773 \
  --types 1 \
  --types 2 \
  --types 3 \
  --types 4 \
  --types 5 \
  --types 7 \
  --authorized 2
```

## 使用建议

1. 该命令通常紧跟在 `storage pool node list` 之后执行。
2. 返回中的设备 ID 可直接填入 `pool create` 请求体的 `mainDevices`、`readCacheDevices`、`writeCacheDevices`。

## 当前环境验证

- 已使用真实 `node-id=babd3c3fbfd34359800ecffcb8907402` 实测成功。
- 当前环境返回 2 个设备，包含 `vda`、`vdb`。
- 该接口签名遵循 SDK demo 的 forward 规则：请求路径保留 `{storageSvcId}`，签名 path 使用 `/storageresmgm/v1/device`。

## 请求参数

请结合上文参数表传参；常见查询命令可按需传 `--index` 与 `--count` 控制分页。

## 返回案例

```json
{
    "error": null,
    "status": "success",
    "responseData": {
        "totalNum": 2,
        "data": [
            {
                "id": "daa2b0e22f1811f1aa4600163e33df73",
                "sn": "3X0GXAC5XT554MMRTPCM",
                "wwn": "",
                "name": "vda",
                "type": 7,
                "status": 2,
                "totalSize": 214748364800,
                "node": {
                    "id": "babd3c3fbfd34359800ecffcb8907402",
                    "ip": "172.31.12.91"
                }
            },
            {
                "id": "daa4639c2f1811f1aa4600163e33df73",
                "sn": "3X0GXABA6Z590XPSQRY0",
                "wwn": "",
                "name": "vdb",
                "type": 7,
                "status": 2,
                "totalSize": 214748364800,
                "node": {
                    "id": "babd3c3fbfd34359800ecffcb8907402",
                    "ip": "172.31.12.91"
                },
                "pool": {
                    "id": "0c6022822f1d11f19b0600163e33df73",
                    "name": "block"
                }
            }
        ]
    }
}
```

## Agent 执行强约束（用于 storage pool create）

- 在将设备 ID 传给 `foundation-cli storage pool create` 之前，必须先筛选 `responseData.data[]`。
- 仅允许使用 `status = 1` 的设备。
- `status != 1` 的设备一律禁止进入 `mainDevices/readCacheDevices/writeCacheDevices`。
- 如果本次返回中没有任何 `status = 1` 的设备，必须停止创建流程并提示用户先处理设备状态。

快速判断示例：

- 可用于建池：`{"id":"dev-1","status":1}`
- 不可用于建池：`{"id":"dev-2","status":2}`
