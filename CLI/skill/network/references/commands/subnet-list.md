# `foundation-cli network subnet list`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli network subnet list --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> [--plane-type <1|2|3|4>]` |
| 方法 | `GET` |
| 路径 | `/clusters/v1/{storageSvcId}/subnet` |
| 风险 | `只读` |

## 必填输入

| 参数 | 必填 | 说明 |
|---|---|---|
| `--tenant-id` | 否 | 租户上下文 |
| `--endpoint` | 是 | Foundation 接口地址 |
| `--ak` | 是 | 访问密钥 |
| `--sk` | 是 | 访问密钥密文 |
| `--storage-svc-id` | 是 | 存储服务 ID |

## 可选输入

| 参数 | 默认值 | 说明 |
|---|---|---|
| `--plane-type` | `3` | 网络平面类型：`1` 管理平面，`2` 存储平面，`3` 备份平面，`4` 复制平面 |

## 参数映射

| CLI Flag | 位置 | 后端参数 |
|---|---|---|
| `--storage-svc-id` | 路径 | `{storageSvcId}` |
| `--plane-type` | Query | `planeType` |

## 返回结果关注点

常见字段：

| 字段 | 说明 |
|---|---|
| `responseData.data[].subnetId` | 子网 ID，后续查询节点列表时会用到 |
| `responseData.data[].subnetName` | 子网名称 |
| `responseData.data[].subnetType` | 子网类型 |
| `responseData.data[].isDefault` | 是否默认子网 |
| `responseData.data[].configProgress` | 子网配置进度 |
| `responseData.subnetMaxNum` | 当前服务允许的子网最大数量 |

## 示例

```bash
foundation-cli network subnet list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --storage-svc-id 1234567890abcdef1234567890abcdef \
  --plane-type 3
```

## 使用建议

1. 默认 `--plane-type 3`；仅在明确需要管理、存储或复制平面时再改成 `1`、`2`、`4`。
2. 先执行该命令确认目标 `subnetId`，再继续执行 `network subnet node list`。
3. 如果 agent 只需要做后续编排，优先提取 `subnetId`、`subnetName` 和 `isDefault`。

## 请求参数

请结合上文参数表传参；常见查询命令可按需传 `--index` 与 `--count` 控制分页。

## 返回案例

```json
{"status":"success","error":null,"responseData":null}
```
