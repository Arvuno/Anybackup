# Storage 命令映射

本文件只负责帮助 agent 选择正确的 `foundation-cli storage` 命令，并跳转到对应的单命令文档。

## 先区分读写

| 类型 | 命令 |
|---|---|
| `写入` | `storage pool create`、`storage pool delete` |
| `只读` | `storage service list`、`storage pool list`、`storage pool node list`、`storage pool node device list` |

## 意图到命令

| 规范路径 | 用户常见意图 | 推荐命令 | 关键要求 | 单命令文档 |
|---|---|---|---|---|
| `storage/service/list` | 获取存储服务列表 | `foundation-cli storage service list --tenant-id <tenant-id>` | 当前环境已验证成功 | [service-list.md](./commands/service-list.md) |
| `storage/pool/list` | 获取存储池列表 | `foundation-cli storage pool list --tenant-id <tenant-id> --storage-svc-id <storage-svc-id>` | 当前环境已验证成功；签名 path 使用 `/storageresmgm/v1/pool` | [pool-list.md](./commands/pool-list.md) |
| `storage/pool/create` | 创建存储池 | `foundation-cli storage pool create --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> --body-file <path>` | 优先使用 `--body-file`；签名 path 使用 `/storageresmgm/v1/pool` | [pool-create.md](./commands/pool-create.md) |
| `storage/pool/delete` | 删除存储池 | `foundation-cli storage pool delete --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> --pool-id <pool-id>` | 必填 `--pool-id`；签名 path 使用 `/storageresmgm/v1/pool` | [pool-delete.md](./commands/pool-delete.md) |
| `storage/pool/node/list` | 获取可建池节点列表 | `foundation-cli storage pool node list --tenant-id <tenant-id> --storage-svc-id <storage-svc-id>` | 当前环境已验证成功；签名 path 使用 `/storageresmgm/v1/pool/node` | [pool-node-list.md](./commands/pool-node-list.md) |
| `storage/pool/node/device/list` | 获取节点设备列表 | `foundation-cli storage pool node device list --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> --node-id <node-id> [--types <type> ...]` | `--types` 支持重复传参；当前环境已验证成功；签名 path 使用 `/storageresmgm/v1/device` | [pool-node-device-list.md](./commands/pool-node-device-list.md) |

## 最容易漏掉的参数

| 命令 | 易漏参数 | 说明 |
|---|---|---|
| `storage pool list` | `--storage-svc-id` | 池查询依赖存储服务上下文 |
| `storage pool create` | `--storage-svc-id` | 创建池时必须锁定目标存储服务 |
| `storage pool delete` | `--pool-id` | 删除需要明确目标池 |
| `storage pool node device list` | `--node-id` | 设备查询依赖节点上下文 |
