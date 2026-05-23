---
name: v9-middleware-deployer
description: 部署或对齐 V9/AnyBackup 的共享中间件栈。用于处理 `resource_01/k8s-yaml`、`v9-infra` Helm Chart、`k8s-cluster` 角色、`middleware` namespace、`local-path`、PostgreSQL、Redis、RabbitMQ、OpenSearch 的新增、迁移、对齐或验收。
---

# V9 Middleware Deployer

## Overview

把 `resource_01/k8s-yaml` 当作 V9 中间件拓扑的事实来源。目标不是机械照抄 YAML，而是把组件集合、命名、存储、镜像、资源和运行时连接口径稳定落到发布包里。

## Cluster Baseline

- 默认按标准 `k8s` 集群处理，不按 `k3s` 做前提假设。
- 不要引入 `k3s` 专属命令、服务名或目录约定，例如 `k3s kubectl`、`k3s ctr`、`/var/lib/rancher/k3s`，除非用户明确说明现场就是 `k3s`。
- 涉及镜像导入、节点运行时排障或离线分发时，先确认现场使用的容器运行时和镜像分发方式，再决定具体命令。

## Workflow

1. 先比对这几处，再决定改动范围：
   - `resource_01/k8s-yaml`
   - `anybackup-agent-release/helm-chart/anybackup-agent/charts/v9-infra/`
   - `anybackup-agent-release/ansible/roles/k8s-cluster/`
   - `anybackup-agent-release/ansible/group_vars/`
2. 固定边界：
   - `local-path` 属于集群级能力，放在 `k8s-cluster`，不要塞进 Helm。
   - `PostgreSQL`、`Redis`、`RabbitMQ`、`OpenSearch` 属于 `v9-infra`。
   - `namespace` 默认用 `middleware`。
3. 默认命名优先保持稳定：
   - `postgres`
   - `redis`
   - `rabbitmq`
   - `opensearch-cluster-master`
4. `PostgreSQL`、`Redis`、`RabbitMQ` 优先使用 `Deployment + PVC + Service + Secret/ConfigMap`。
5. `OpenSearch` 优先保留 vendored subchart，不要手写一套等价模板。
6. 同步修改运行时默认连接串、验证脚本和文档，避免“chart 已变，secret/README 还停留在旧口径”。

## Guardrails

- 不要把源 YAML 里的密码直接硬编码进交付物；全部保持参数化。
- 不要把 KWeaver 自己的 `MariaDB/Kafka` 依赖链并入这套 middleware，除非用户明确要求做架构收敛。
- 改完后清理旧口径残留，例如：
  - `v9-system`
  - `v9-infra-postgres`
  - `v9-infra-rabbitmq`
- 任何中间件拓扑变更后，都要同步检查：
  - `v9_runtime` 默认 URL
  - `cluster-smoke.sh`
  - `README.md`
  - `install.sh`

## Validation

- 先跑 `helm lint`。
- 再跑 `helm template`，确认四个组件都能渲染出来。
- 检查 `k8s-cluster` 是否能幂等安装 `local-path`。
- 检查运行时 secret 是否默认指向 `middleware` 下的服务名。
