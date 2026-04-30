# `foundation-cli mysql datasource recovery`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli mysql datasource recovery --tenant-id <tenant-id> --data-set-id <data-set-id> --storage-service-id <id> --timestamp <ts> [--restore-gran <0|1|2|3>] [--request-id <id>] [--full-path <path>] [--index <n>] [--count <n>]` |
| 方法 | `GET` |
| 路径 | `/backupmgm/v1/mysql/recovery_datasource` |
| 风险 | `只读` |

## 请求参数

### 公共参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户标识；可由环境变量 FOUNDATION_TENANT_ID 提供 |
| endpoint | `--endpoint` | string | 是 | Foundation 基础地址 |
| ak | `--ak` | string | 是 | 访问密钥 |
| sk | `--sk` | string | 是 | 访问密钥密文 |
| targetVersion | `--target-version` | string | 否 | 默认值为 `8.0.9.0` |

### 命令参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| dataSetId | `--data-set-id` | string | 是 | 恢复数据集 ID |
| storageServiceId | `--storage-service-id` | string | 是 | 存储服务 ID |
| timestamp | `--timestamp` | string | 是 | 时间点时间戳 |
| restoreGran | `--restore-gran` | string | 否 | 允许值： `0/1/2/3` |
| requestId | `--request-id` | string | 否 | 请求 ID；为空时可省略 |
| fullPath | `--full-path` | string | 否 | 查询路径，默认 `/` |
| index | `--index` | int | 否 | 分页起始，默认 `0` |
| count | `--count` | int | 否 | 分页条数，默认 `100`，范围 `1~100` |

## 响应参数

### 顶层结构

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 接口执行状态 |
| error | object/null | 错误信息 |
| responseData | object | 恢复时间点数据源响应体 |

### `responseData`

| 字段 | 类型 | 说明 |
|---|---|---|
| requestId | string | 请求 ID |
| totalNum | int | 总记录数 |
| remainAmount | int | 剩余记录数 |
| nextIndex | int | 下一页索引 |
| data | object[] | 恢复时间点数据源列表 |

### `data[]`

| 字段 | 类型 | 说明 |
|---|---|---|
| attributes | int | 属性位掩码 |
| checkAble | bool | 节点是否可选 |
| checked | bool | 节点是否已选中 |
| displayFlag | bool | 展示标记 |
| dispPath | string | 展示路径 |
| expandedFlag | bool | 节点是否已展开 |
| extAttributes | int | 扩展属性 |
| fullPath | string | 完整路径 |
| metadata | string | 元数据内容 |
| nodeType | int | 节点类型 |
| objectSize | int | 对象大小 |
| disableEdit | bool | 是否禁止编辑 |

## 返回案例

说明：
1. 当前仓库未附带该命令的真实接口成功响应。
2. 下例仅用于说明常见返回结构层级，字段名和取值以后端当前版本实际返回为准。
3. 建议后续补充 1 份真实成功响应和 1 份典型失败响应。

```json
{
  "error": null,
  "status": "success",
  "responseData": {
    "requestId": "sample-request-id",
    "totalNum": 1,
    "remainAmount": 0,
    "nextIndex": 1,
    "data": [
      {
        "attributes": 0,
        "checkAble": true,
        "checked": false,
        "displayFlag": true,
        "dispPath": "db1/table1",
        "expandedFlag": false,
        "extAttributes": 0,
        "fullPath": "/db1/table1",
        "metadata": "",
        "nodeType": 2,
        "objectSize": 1048576,
        "disableEdit": false
      }
    ]
  }
}
```

## 示例
```bash
foundation-cli mysql datasource recovery \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --storage-service-id <storage-service-id> \
  --data-set-id <data-set-id> \
  --timestamp <timestamp> \
  --restore-gran 2 \
  --full-path / \
  --index 0 \
  --count 100
```



