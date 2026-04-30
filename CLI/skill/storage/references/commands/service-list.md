# `foundation-cli storage service list`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli storage service list --tenant-id <tenant-id>` |
| 方法 | `GET` |
| 路径 | `/mstoragesvcmgm/v1/storage-svc?onlyStorage=true` |
| 风险 | `只读` |

## 必填输入

| 参数 | 必填 | 说明 |
|---|---|---|
| `--tenant-id` | 否 | 租户上下文 |
| `--endpoint` | 是 | Foundation 接口地址 |
| `--ak` | 是 | 访问密钥 |
| `--sk` | 是 | 访问密钥密文 |

## 说明

1. 该命令由 CLI 固定追加 `onlyStorage=true`。
2. 默认输出保持后端原始 JSON。

## 查询参数映射

该命令不额外暴露业务查询参数，CLI 固定向后端发送：

| 后端参数 | 固定值 | 说明 |
|---|---|---|
| `onlyStorage` | `true` | 仅返回存储服务 |

## 返回结果关注点

常见字段：

| 字段 | 说明 |
|---|---|
| `responseData.data[].id` | 存储服务 ID，后续 `pool` 相关命令会用到 |
| `responseData.data[].name` | 存储服务名称 |
| `responseData.data[].vip` | 存储服务 VIP |
| `responseData.data[].status` | 存储服务状态 |
| `responseData.data[].type` | 存储服务类型 |

## 当前环境验证

- 已使用当前测试环境实测成功。
- 可拿到 `storage-svc-id=4235757746e447d18d4d7a8b3f5edc54`，后续池相关命令可复用该 ID 继续排查。

## 返回结果

CLI 标准返回结构：

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 执行状态 |
| error | object/null | 错误信息（失败时） |
| responseData | object | 接口业务返回 |
| responseData.totalNum | int | 存储服务总数 |
| responseData.data | object[] | 存储服务列表 |

## 返回案例

```json
{
  "error": null,
  "responseData": {
    "data": [
      {
        "id": "52247d9c121745eab3a7880817457df5",
        "name": "StorageService_192.168.1.1",
        "vip": "192.168.1.1",
        "status": 1,
        "dec": "",
        "type": 2
      },
      {
        "id": "9fc189d638944a57ac042874a21d94b7",
        "name": "StorageService_192.168.1.2",
        "vip": "192.168.1.2",
        "status": 1,
        "dec": "",
        "type": 2
      }
    ],
    "totalNum": 2
  },
  "status": "success"
}
```

## 示例

```bash
foundation-cli storage service list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk>
```

## 使用建议

1. 先执行该命令，拿到 `responseData.data[].id`。
2. 再把该 ID 作为 `--storage-svc-id` 传给 `pool list`、`pool create`、`pool delete`、`pool node list`、`pool node device list`。

## 请求参数

请结合上文参数表传参；常见查询命令可按需传 `--index` 与 `--count` 控制分页。
