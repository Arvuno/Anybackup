# `foundation-cli network subnet node list`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli network subnet node list --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> --subnet-id <subnet-id>` |
| 方法 | `GET` |
| 路径 | `/clusters/v1/{storageSvcId}/subnet/nodes` |
| 风险 | `只读` |

## 必填输入

| 参数 | 必填 | 说明 |
|---|---|---|
| `--tenant-id` | 否 | 租户上下文 |
| `--endpoint` | 是 | Foundation 接口地址 |
| `--ak` | 是 | 访问密钥 |
| `--sk` | 是 | 访问密钥密文 |
| `--storage-svc-id` | 是 | 存储服务 ID |
| `--subnet-id` | 是 | 子网 ID |

## 可选输入

| 参数 | 默认值 | 说明 |
|---|---|---|
| `--index` | `0` | 分页起始位置 |
| `--count` | `30` | 分页数量 |
| `--plane-type` | `3` | 网络平面类型，当前网络平面接口默认使用 `3` |

## 查询参数映射

| CLI Flag | 后端参数 |
|---|---|
| `--index` | `index` |
| `--count` | `count` |
| `--plane-type` | `planeType` |
| `--subnet-id` | `subnetId` |

## 返回结果关注点

常见字段：

| 字段 | 说明 |
|---|---|
| `responseData.data[].nodeId` | 节点 ID，后续查询节点业务 IP 或设置业务 IP 时会用到 |
| `responseData.data[].hostName` | 节点主机名 |
| `responseData.data[].nodeIp` | 节点管理 IP |
| `responseData.data[].address` | 节点地址信息，部分环境可能为空 |
| `responseData.portDescData` | 端口说明映射，可辅助判断网络连通要求 |

## 示例

```bash
foundation-cli network subnet node list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --storage-svc-id 1234567890abcdef1234567890abcdef \
  --subnet-id fab4ec3015df11f1966b0050568952bd \
  --plane-type 3 \
  --index 0 \
  --count 30
```

## 使用建议

1. 默认 `--plane-type 3`，如果后端接口约束未变化，通常不建议覆盖。
2. 该命令通常紧跟在 `network subnet list` 之后执行，先拿到目标 `subnetId`。
3. 后续要查询节点可用业务 IP 时，直接复用返回里的 `nodeId` 执行 `network subnet node ip list`。

## 请求参数

请结合上文参数表传参；常见查询命令可按需传 `--index` 与 `--count` 控制分页。

## 返回案例

```json
{"status":"success","error":null,"responseData":null}
```
