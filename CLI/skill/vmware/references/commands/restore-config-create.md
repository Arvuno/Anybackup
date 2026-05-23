# `foundation-cli vmware restore-config create`

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli vmware restore-config create --object-id <id> --data '<json>'` |
| Method | `POST` |
| Path | `/backupmgm/v1/virtual/vmware/{objectId}/recovery` |
| Risk | `write` |

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
| objectId | `--object-id` | string | 是 | 保护对象 ID |
| body | `--data` | JSON string | 是 | 恢复配置请求体 |

## Body 参数（`--data`）

对象：`VMwareRecoveryConfigRequest`

| 字段 | 类型 | 必填 | 来源结构体 | 说明 |
|---|---|---|---|---|
| timePointId | string | 是 | VMwareRecoveryConfigRequest | 恢复时间点唯一标识，`len=32` |
| productionSystemId | string | 是 | VMwareRecoveryConfigRequest | 恢复目标生产系统 ID（`uri` 字段语义） |
| runImmediately | bool | 否 | VMwareRecoveryConfigRequest | 是否立即执行 |
| objectPassword | string | 否 | ObjectPasswordConfig | 对象密码（无密码可传空） |
| failureRetry | int32 | 否 | FailureRetryConfig | 失败重试开关 |
| failureRetryCount | int32 | 否 | FailureRetryConfig | 失败重试次数 |
| failureRetryInterval | int32 | 否 | FailureRetryConfig | 失败重试间隔 |
| encryptionTrans | int32 | 否 | EncryptionConfig | 传输加密开关 |
| clientId | string | 是 | VMwareRecoveryConfig | 执行客户端 ID |
| vmgns | string | 是 | VMwareRecoveryConfig | 恢复对象路径 |
| transportMode | string | 否 | VMwareRecoveryConfig | 传输模式（异构恢复可不传，非异构恢复建议传） |
| usetoNBD | bool | 否 | VMwareRecoveryConfig | 是否转 NBD |
| autoStartVm | bool | 否 | VMwareRecoveryConfig | 恢复后是否自动开机 |
| useOriginalMac | bool | 否 | VMwareRecoveryConfig | 是否使用原 MAC |
| hostname | string | 否 | VMwareRecoveryConfig | 恢复主机名 |
| location | string | 否 | VMwareRecoveryConfig | 恢复位置 |
| resource | string | 否 | VMwareRecoveryConfig | 恢复资源池 |
| overwriteVmPath | string | 否 | VMwareRecoveryConfig | 覆盖恢复目标 VM 路径 |
| overwriteVmUuid | string | 否 | VMwareRecoveryConfig | 覆盖恢复目标 VM UUID |
| overwriteVmName | string | 否 | VMwareRecoveryConfig | 覆盖恢复目标 VM 名称 |
| nameSuffix | string | 否 | VMwareRecoveryConfig | 名称后缀 |
| newName | string | 否 | VMwareRecoveryConfig | 新建恢复 VM 名称 |
| recoveryMode | int | 是 | VMwareRecoveryConfig | 恢复方式 |
| configMode | int | 否 | VMwareRecoveryConfig | 配置方式（新建恢复时常用） |
| overlayMode | int | 否 | VMwareRecoveryConfig | 覆盖方式（覆盖恢复时常用） |
| enableAutoRepair | int32 | 否 | VMwareRecoveryConfig | 异构恢复自动修复开关 |
| isFromDrm | bool | 否 | VMwareRecoveryConfig | 是否 DRM 恢复 |
| enableCustomCpuAndMem | bool | 否 | VMwareRecoveryConfig | 是否启用自定义 CPU/内存 |
| memorySize | int | 否 | VMwareRecoveryConfig | 内存大小（GiB） |
| numCoresPerSocket | int | 否 | VMwareRecoveryConfig | 每 CPU 插槽核心数 |
| cpuSlots | int | 否 | VMwareRecoveryConfig | CPU 插槽数 |
| disksConfig | object[] | 否 | VMwareRecoveryConfigRequest | 恢复磁盘配置 |
| vifsConfig | object[] | 否 | VMwareRecoveryConfigRequest | 恢复网卡配置 |
| metaData | object | 否 | VMwareRecoveryConfigRequest | 时间点元数据 |
| supportCrossPlatRecovery | bool | 否 | VMwareRecoveryConfigRequest | 是否支持异构恢复 |

`disksConfig[]`（`DiskConfig`）字段：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| id | string | 否 | 磁盘 ID |
| name | string | 否 | 磁盘名称 |
| storageName | string | 否 | 存储名称 |
| storageNameWithSize | string | 否 | 存储名称（含大小） |

`vifsConfig[]`（`VifConfig`）字段：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| name | string | 否 | 网卡名称 |
| mac | string | 否 | MAC 地址 |
| networkName | string | 否 | 网络名称 |

`metaData`（`TimePointMetadata`）字段：

| 字段路径 | 类型 | 必填 | 说明 |
|---|---|---|---|
| metaData.version | int | 否 | 元数据版本 |
| metaData.platformType | int | 否 | 平台类型 |
| metaData.platform | object | 否 | 平台信息 |
| metaData.cpuSlots | int | 否 | CPU 槽数 |
| metaData.cpuCores | int | 否 | 每 CPU 核心数 |
| metaData.memorySize | int64 | 否 | 内存大小（字节） |
| metaData.osSystem | string | 否 | 操作系统名称 |
| metaData.osBit | int | 否 | 操作系统位数 |
| metaData.bootType | string | 否 | 引导类型（uefi/bios） |
| metaData.vmTools | bool | 否 | 是否安装 VMTools |
| metaData.arch | string | 否 | 架构（x86/aarch） |
| metaData.vmId | string | 否 | VM ID |
| metaData.vmName | string | 否 | VM 名称 |
| metaData.vmFullName | string | 否 | VM 完整路径 |
| metaData.totalSize | int64 | 否 | 总占用空间（字节） |
| metaData.disk | object[] | 否 | 磁盘元数据 |
| metaData.nic | object[] | 否 | 网卡元数据 |

`metaData.platform`（`Platform`）字段：

| 字段路径 | 类型 | 必填 | 说明 |
|---|---|---|---|
| metaData.platform.platformVersion | string | 否 | 平台版本 |
| metaData.platform.platformId | string | 否 | 平台 ID |
| metaData.platform.platformIp | string | 否 | 平台 IP |
| metaData.platform.platformName | string | 否 | 平台名称 |

`metaData.disk[]`（`Disk`）字段：

| 字段路径 | 类型 | 必填 | 说明 |
|---|---|---|---|
| metaData.disk[].id | string | 否 | 磁盘 ID |
| metaData.disk[].name | string | 否 | 磁盘名称 |
| metaData.disk[].devName | string | 否 | 设备名称 |
| metaData.disk[].size | int64 | 否 | 大小（字节） |
| metaData.disk[].bootable | bool | 否 | 是否可启动 |
| metaData.disk[].format | string | 否 | 磁盘格式 |
| metaData.disk[].bus | string | 否 | 总线类型 |
| metaData.disk[].sequence | int64 | 否 | 磁盘序号 |
| metaData.disk[].storageName | string | 否 | 存储名称 |

`metaData.nic[]`（`NIC`）字段：

| 字段路径 | 类型 | 必填 | 说明 |
|---|---|---|---|
| metaData.nic[].id | string | 否 | 网卡 ID |
| metaData.nic[].name | string | 否 | 网卡名称 |
| metaData.nic[].ip | string | 否 | IP 地址 |
| metaData.nic[].macAddr | string | 否 | MAC 地址 |
| metaData.nic[].networkName | string | 否 | 网络名称 |

结构体关系（用于校准对象层级）：

| 层级 | 结构体 |
|---|---|
| 1 | VMwareRecoveryConfigRequest |
| 2 | ObjectPasswordConfig + FailureRetryConfig + VMwareRecoveryConfig + EncryptionConfig |
| 3 | DiskConfig / VifConfig / TimePointMetadata |
| 4 | Platform / Disk / NIC |

说明：

| 项 | 内容 |
|---|---|
| 已包含字段 | 当前文档已覆盖 `vm_recover.go` 请求体全部字段及其嵌套对象 |
| 一致性约束 | 命令路径使用 `--object-id`，参数模型包含 `productionSystemId`（`uri` 语义）；调用时建议两者指向同一恢复目标上下文 |

## 枚举说明

### `recoveryMode`

| 值 | 含义 |
|---|---|
| 1 | 新建恢复 |
| 2 | 覆盖恢复 |
| 3 | 异构恢复 |

### `configMode`

| 值 | 含义 |
|---|---|
| 1 | 原配置恢复 |
| 2 | 指定配置恢复 |

### `overlayMode`

| 值 | 含义 |
|---|---|
| 1 | 全量覆盖恢复 |

### `enableAutoRepair`

| 值 | 含义 |
|---|---|
| 0 | 关闭自动修复 |
| 1 | 开启自动修复 |

### `transportMode`

| 值 | 含义 |
|---|---|
| nbd | NBD 模式 |
| nbdssl | NBDSSL 模式 |
| hotadd | HotAdd 模式 |
| san | SAN 模式 |
| auto | 自动选择 |

### `encryptionTrans`

| 值 | 含义 |
|---|---|
| 0 | 关闭传输加密 |
| 1 | 开启传输加密 |

### `failureRetry`

| 值 | 含义 |
|---|---|
| 0 | 关闭失败重试 |
| 1 | 开启失败重试 |

### `metaData.platformType`

| 值 | 含义 |
|---|---|
| 1 | XenServer |
| 2 | InCloud Sphere |
| 3 | VMware |
| 4 | FusionCompute |
| 5 | RHEV |
| 6 | Hyper-V |
| 7 | H3CLOUD |
| 8 | CAS |
| 9 | 废弃 |
| 10 | 废弃 |
| 11 | SANGFOR HCI |
| 12 | Zstack |
| 13 | OracleVM |
| 14 | IROS |
| 15 | InCloudSphereKvm |
| 16 | SmartX |
| 17 | Nutanix |
| 18 | QingCloud |
| 19 | OpenStack |
| 20 | HCS |
| 21 | AliyunPrivate |
| 22 | Bingo |
| 30 | FusionOne |
| 31 | ProxmoxVE |

## 返回结果

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 执行状态 |
| error | object/null | 错误信息 |
| responseData | object | 配置创建结果 |

## 示例

```bash
foundation-cli vmware restore-config create \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --object-id <object-id> \
  --data '{"timePointId":"<time-point-id-32chars>","productionSystemId":"<production-system-id>","clientId":"<client-id>","vmgns":"<target-vm-path>","recoveryMode":1,"runImmediately":true}'
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli\ vmware\ restore-config\ create |
| 风险 | write |

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
