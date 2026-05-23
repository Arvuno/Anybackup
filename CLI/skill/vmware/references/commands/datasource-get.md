# `foundation-cli vmware datasource get`

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli vmware datasource get --production-system-id <id>` |
| Method | `GET` |
| Path | `/backupmgm/v1/virtual/vmware/{productionSystemId}/sub_objects` |
| Risk | `read-only` |

## 请求参数

### 共享参数

| 字段 | CLI Flag | 必填 | 说明 |
|---|---|---|---|
| tenantId | `--tenant-id` | 否 | 租户标识 |
| endpoint | `--endpoint` | 是 | Foundation 服务地址 |
| ak | `--ak` | 是 | Access Key |
| sk | `--sk` | 是 | Secret Key |
| targetVersion | `--target-version` | 否 | 目标版本，默认 `9.0.9.0` |

### 命令参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| productionSystemId | `--production-system-id` | string | 是 | 生产系统 ID |

## 返回参数

CLI 外层结构：

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 执行状态 |
| error | object/null | 错误信息 |
| responseData | object | 业务返回体 |

`responseData`（`VMwareDatasourceListResponse`）：

| 字段 | 类型 | 说明 |
|---|---|---|
| totalNum | int64 | 总数 |
| data | object[] | 数据源子对象列表 |

`responseData.data[]`（`VMwareDatasource`）字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| productionResourceId | string | 生产系统下资源唯一标识 |
| id | string | 对象 ID |
| name | string | 对象名称 |
| type | int | 对象类型 |
| displayType | string | 展示类型/路径分类 |
| path | string | 对象路径 |
| expandable | bool | 是否可展开 |
| checkable | bool | 是否可选择 |

## 枚举说明

| 字段 | 值 | 说明 |
|---|---|---|
| type | 后端动态枚举 | `params/vmware/vmware_sub_object_list.go` 仅定义为 `int`，未给出固定值映射。建议联调时按接口实际返回字典展示。 |

## 示例

```bash
foundation-cli vmware datasource get \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --production-system-id <production-system-id>
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli\ vmware\ datasource\ get |
| 风险 | read-only |

## 返回案例

```json
{"status":"success","error":null,"responseData":null}
```
