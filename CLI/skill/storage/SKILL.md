---
name: foundation-cli-storage
description: 当用户需要通过 `foundation-cli storage` 执行存储服务、存储池、建池节点或节点设备相关命令时使用。
metadata:
  requires:
    bins: ["foundation-cli"]
  cliHelp: "foundation-cli storage --help"
---

# Foundation CLI Storage 技能

当任务属于 `foundation-cli storage` 域命令时，使用本技能。

本文件是入口，不是完整参考手册。
它的作用是帮助 agent：
1. 判断用户意图是否属于 `storage` 域。
2. 选择正确的 `storage` 命令。
3. 跳转到 `references/commands/*.md` 中对应的单命令文档。

## 适用场景

- 用户要获取存储服务列表。
- 用户要获取存储池列表。
- 用户要创建存储池。
- 用户要删除存储池。
- 用户要获取可用于创建存储池的节点列表。
- 用户要获取指定节点下可用于建池的设备列表。

## 不适用场景

- 用户要查询子网、子网节点或节点业务 IP，此时应优先切换到 `network` 域技能。
- 用户需要的是其他业务域对象能力，而不是存储服务或存储池相关操作。
- 用户只是要走通用 API 透传，且没有标准 `storage` 命令覆盖时，应回到共享技能重新判断是否改走 `api` 域。

## 标准流程

1. 先使用“快速判断规则”排除最常见的误选命令场景。
2. 再使用 [命令映射](./references/commands.md) 选择命令。
3. 立即打开 `./references/commands/` 下对应的单命令文档。
4. 基于该单命令文档构造最终 CLI。

## 快速判断规则

- 如果用户要查“存储服务列表”，使用 `storage service list`。
- 如果用户要查“存储池列表”，使用 `storage pool list`。
- 如果用户要“创建存储池”，使用 `storage pool create`。
- 如果用户要“删除存储池”，使用 `storage pool delete`。
- 如果用户要查“可建池节点”，使用 `storage pool node list`。
- 如果用户要查“某个节点下可建池设备”，使用 `storage pool node device list`。
- 除 `service list` 外，其余命令通常都依赖 `--storage-svc-id`。
- `pool delete` 额外依赖 `--pool-id`。
- `pool node device list` 额外依赖 `--node-id`。
- 当前环境实测：`storage service list`、`storage pool list`、`storage pool node list`、`storage pool node device list` 均已验证成功。
- `storageresmgm` 这组命令签名需遵循 SDK demo 的 forward 规则：请求路径保留 `{storageSvcId}`，但签名 path 要去掉该段。

## 写入与只读

写入命令：

- `storage pool create`
- `storage pool delete`

只读命令：

- `storage service list`
- `storage pool list`
- `storage pool node list`
- `storage pool node device list`

## 必要的全局上下文

除非调用环境已经隐式提供，否则先确认这些参数：

- `--tenant-id`
- `--endpoint`
- `--ak`
- `--sk`

## 常见误判

- 不要把 `storage` 域和 `network` 域混淆。
- 不要把“存储服务列表”和“存储池列表”混淆。
- 不要在没有 `--storage-svc-id` 的情况下执行池相关命令。
- 不要在没有 `--pool-id` 的情况下执行 `pool delete`。
- `pool create` 是写操作，且请求体较长时优先用 `--body-file`。

## 意图到命令

| 意图 | 推荐命令 | 下一步打开的文件 |
|---|---|---|
| 获取存储服务列表 | `foundation-cli storage service list --tenant-id <tenant-id>` | [service-list.md](./references/commands/service-list.md) |
| 获取存储池列表 | `foundation-cli storage pool list --tenant-id <tenant-id> --storage-svc-id <storage-svc-id>` | [pool-list.md](./references/commands/pool-list.md) |
| 创建存储池 | `foundation-cli storage pool create --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> --body-file <path>` | [pool-create.md](./references/commands/pool-create.md) |
| 删除存储池 | `foundation-cli storage pool delete --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> --pool-id <pool-id>` | [pool-delete.md](./references/commands/pool-delete.md) |
| 获取可建池节点列表 | `foundation-cli storage pool node list --tenant-id <tenant-id> --storage-svc-id <storage-svc-id>` | [pool-node-list.md](./references/commands/pool-node-list.md) |
| 获取节点设备列表 | `foundation-cli storage pool node device list --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> --node-id <node-id>` | [pool-node-device-list.md](./references/commands/pool-node-device-list.md) |

## 参考资料

- [命令映射](./references/commands.md)
