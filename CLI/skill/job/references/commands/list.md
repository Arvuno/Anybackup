# `foundation-cli job list`

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | `foundation-cli job list --tenant-id <tenant-id> [filters]` |
| 方法 | `GET` |
| 路径 | `/job_center/v1/jobs` |
| 风险 | `只读` |

## 请求参数

### 共享参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户标识；可由环境变量 FOUNDATION_TENANT_ID 提供 |
| endpoint | `--endpoint` | string | 是 | Foundation console 基础 URL |
| ak | `--ak` | string | 是 | 访问密钥（AK） |
| sk | `--sk` | string | 是 | 访问密钥密文（SK） |
| targetVersion | `--target-version` | string | 否 | 目标版本，默认 `9.0.9.0` |

### 命令参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| index | `--index` | int | 否 | 分页起始位置，默认 `0` |
| count | `--count` | int | 否 | 分页数量，默认 `10` |
| statuses | `--statuses` | int64[] | 否 | 作业状态列表，重复 flag 传入；取值见下方“作业状态枚举” |
| objectId | `--object-id` | string | 否 | 保护对象 ID |
| objectName | `--object-name` | string | 否 | 保护对象名称 |
| appTypes | `--app-types` | string[] | 否 | 应用类型列表，重复 flag 传入 |
| businessTypes | `--business-types` | int64[] | 否 | 业务类型列表，重复 flag 传入 |
| startTime | `--start-time` | int64 | 否 | 起始时间（毫秒） |
| endTime | `--end-time` | int64 | 否 | 结束时间（毫秒） |
| sort | `--sort` | string | 否 | 排序字段 |

## 返回结果

## 作业状态枚举（status / --statuses）

| 状态码 | 枚举名 | 含义 |
|---|---|---|
| `0` | `ExecStatusProto_Not_Start` | 未启动 |
| `10` | `ExecStatusProto_Approving` | 审批中 |
| `100` | `ExecStatusProto_Queuing` | 排队中 |
| `200` | `ExecStatusProto_Ready` | 准备中 |
| `300` | `ExecStatusProto_Running` | 正在执行 |
| `310` | `ExecStatusProto_Retrying` | 重试中 |
| `400` | `ExecStatusProto_Stopping` | 正在停止 |
| `500` | `ExecStatusProto_Abnormal` | 异常 |
| `600` | `ExecStatusProto_Failed` | 失败 |
| `700` | `ExecStatusProto_Canceled` | 已取消 |
| `800` | `ExecStatusProto_Success` | 成功 |
| `900` | `ExecStatusProto_WarnSuss` | 成功但有警告 |
| `1000` | `ExecStatusProto_PartSuss` | 部分成功 |
| `1100` | `ExecStatusProto_Stopped` | 已停止 |

| 字段 | 类型 | 说明 |
|---|---|---|
| totalNum | int64 | 作业总数 |
| data | Job[] | 作业列表 |

### `data[]` 核心字段

| 字段 | 类型 | 说明 |
|---|---|---|
| jobId | string | 作业 ID |
| objectId | string | 保护对象 ID |
| objectName | string | 保护对象名称 |
| status | int32 | 作业状态 |
| appType | string | 应用类型 |
| operationType | int32 | 一级业务类型 |
| businessType | int32 | 业务类型 |
| retriedCount | int64 | 已重试次数 |
| startTime | int64 | 开始时间 |
| endTime | int64 | 结束时间，未结束时可能为 `-1` |
| runTime | int64 | 运行时长 |
| queueTime | int64 | 排队时长 |
| strategyId | string | 策略 ID |
| strategyName | string | 策略名称 |
| strategyType | int | 策略类型 |
| clientId | string | 客户端 ID |
| clientName | string | 客户端名称 |
| clientIp | string | 客户端 IP |
| userId | string | 用户 ID |
| userName | string | 用户名称 |
| isAutoExecute | bool | 是否自动执行 |
| isExist | bool | 目标对象是否仍存在 |
| isolation | bool | 是否孤立对象 |
| genMode | int32 | 来源类型 |
| remark | string | 作业描述 |
| tenantName | string | 租户名称 |
| speed | string | 作业速度 |
| completedData | int64 | 已完成数据量 |
| storedSize | int64 | 已存储数据量（可选） |
| sendSize | int64 | 已发送数据量 |
| storageServiceName | string | 存储服务名称 |
| storagePoolName | string | 存储池名称 |
| targetDestinationIp | string | 目标端 IP（可选） |
| deduplication | int | 去重配置 |

## 返回案例

```json
{
	"status": "success",
	"error": null,
	"responseData": {
		"totalNum": 5,
		"data": [{
			"jobId": "f12f1d3061581124b50400163e33df73",
			"objectId": "3086b6b2db7b67950be03a56081e41a1",
			"objectName": "test",
			"status": 700,
			"appType": "eso_backupengine_fileengine",
			"operationType": 1,
			"businessType": 1,
			"retriedCount": 0,
			"startTime": 1775193783000,
			"endTime": 1775195235000,
			"runTime": 1452000,
			"queueTime": 0,
			"strategyId": "",
			"strategyName": "",
			"strategyType": 0,
			"clientId": "5c872f826b2f03101b6578947ee13d2c",
			"clientName": "anybackup",
			"clientIp": "172.31.12.91",
			"userId": "cdec38822f1811f19fd400163e33df73",
			"userName": "admin",
			"isAutoExecute": false,
			"isExist": true,
			"isolation": false,
			"genMode": 0,
			"remark": "",
			"tenantName": "",
			"speed": "0",
			"completedData": 0,
			"storedSize": 0,
			"sendSize": 0,
			"storageServiceName": "StorageService_172.31.12.91",
			"storagePoolName": "block",
			"deduplication": 0
		}]
	}
}
```

## 示例

```bash
foundation-cli job list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --index 0 \
  --count 20
```
