# `foundation-cli network subnet node ip set`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli network subnet node ip set --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> --subnet-id <subnet-id> --node-id <node-id> --ip <ip>` |
| 方法 | `POST` |
| 路径 | `/clusters/v1/{storageSvcId}/subnet/nodes/node_ip_addresses` |
| 风险 | `写入` |

## 必填输入

| 参数 | 必填 | 说明 |
|---|---|---|
| `--tenant-id` | 否 | 租户上下文 |
| `--endpoint` | 是 | Foundation 接口地址 |
| `--ak` | 是 | 访问密钥 |
| `--sk` | 是 | 访问密钥密文 |
| `--storage-svc-id` | 是 | 存储服务 ID |
| `--subnet-id` | 是 | 子网 ID |
| `--node-id` | 是 | 节点 ID |
| `--ip` | 是 | 节点业务 IP，必须是标准 IPv4 或 IPv6 |

## 请求体

CLI 会自动构造如下 JSON：

```json
{
  "nodeId": "839b4b217e214c348152221288257773",
  "ip": "10.4.111.55",
  "subnetId": "fab4ec3015df11f1966b0050568952bd"
}
```

## 参数映射

| CLI Flag | 位置 | 后端参数 |
|---|---|---|
| `--storage-svc-id` | 路径 | `{storageSvcId}` |
| `--subnet-id` | Body | `subnetId` |
| `--node-id` | Body | `nodeId` |
| `--ip` | Body | `ip` |

## 返回结果关注点

1. 该接口通常以 HTTP 状态码和统一响应包中的 `status` 字段判断是否成功。
2. 如果后端返回空数据体，CLI 仍以请求成功响应为准，不额外补造业务字段。

## 示例

```bash
foundation-cli network subnet node ip set \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --storage-svc-id 1234567890abcdef1234567890abcdef \
  --subnet-id fab4ec3015df11f1966b0050568952bd \
  --node-id 839b4b217e214c348152221288257773 \
  --ip 10.4.111.55
```

## 使用建议

1. 本命令不使用 `--data` 或 `--body-file`，直接通过结构化 flag 构造请求体。
2. CLI 会先校验 `--subnet-id`、`--node-id` 非空，并校验 `--ip` 格式。
3. 写入前建议先执行 `network subnet node ip list`，确认目标 IP 存在且可用。

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
