---
name: anybackup-v9-architecture-charter
description: 维护 AnyBackup V9 的纲领性架构对齐规则。用于在设计、实现、评审、打包、部署或文档修订时，对齐 `E:\Code\aa-feat\docs\05-技术架构` 的目标架构，特别是 Kubernetes/VKE 部署基线、Anybackup Agent 与 Foundation 的系统边界、Agent 侧中间件职责、Core Agent Service 与 KWeaver Core 的职责分层，以及 `E:\Code\kweaver-core\deploy` 中的 KWeaver Core 上游部署资产。
---

# Anybackup V9 Architecture Charter

## Overview

这个 skill 是 AnyBackup V9 的“总纲”，先回答“这次改动应该对齐成什么样”，再把具体实现交给细分 skill。它的目标不是替代部署细节，而是防止实现、文档和上游组件接入把长期架构边界写乱。

## First Reads

- 先读 `references/source-map.md`，按任务类型选择最小必要文档集合。
- 如果任务涉及服务代码仓、镜像构建、Helm Chart 归位、版本号或团队交付物收口，优先读：
  - `E:\Code\script\anybackup-agent-release\docs\00_设计文档\01_V9服务仓库交付规范.md`
- 如果任务涉及总体边界、组件归位或术语统一，优先读：
  - `E:\Code\aa-feat\docs\05-技术架构\架构总览.md`
  - `E:\Code\aa-feat\docs\05-技术架构\架构设计\部署架构设计.md`
- 如果任务涉及 AI 协作链路、Core Agent Service 或 KWeaver 接入，额外读：
  - `E:\Code\aa-feat\docs\05-技术架构\核心智能体服务\Core Agent Service-架构设计.md`
  - `E:\Code\aa-feat\docs\05-技术架构\核心智能体服务\KWeaver Core SDK能力说明.md`
- 如果任务涉及 KWeaver Core 的安装、裁剪或嵌入式部署，额外读：
  - `E:\Code\kweaver-core\deploy\README.zh.md`
  - `E:\Code\kweaver-core\deploy\deploy.sh`
  - `E:\Code\kweaver-core\deploy\scripts\services\core.sh`
  - `E:\Code\kweaver-core\deploy\scripts\services\k8s.sh`

## Charter

1. `E:\Code\aa-feat\docs\05-技术架构` 是目标架构事实来源。
   - 当前仓库实现如果与这些文档不一致，优先收敛实现。
   - 如果短期做不到，明确标成“版本差异”或“过渡方案”，不要把临时实现混成长期架构基线。
2. `AnyBackup V9` 的整体边界固定为三段：
   - `Anybackup Agent`：SaaS 侧决策、协作与治理运行时
   - `Anybackup Foundation`：用户数据中心内的执行链路
   - `Anybackup Client`：被保护资产侧采集层
3. `KWeaver Core` 是 `AnyBackup V9` 的一个组件，不是重定义系统边界的独立产品。
   - 在 V9 里，它属于 SaaS 侧 AI 能力面。
   - 对外口径上，`Core Agent Service + KWeaver Core` 共同构成统一 AI 能力面。
4. 默认部署基线是标准 `k8s` / `VKE`。
   - 不把 `k3s`、单节点集群或某个现场运维习惯写成产品默认前提。
5. `Redis`、`PostgreSQL`、`RabbitMQ`、`OpenSearch` 只属于 `Anybackup Agent` 侧。
   - 不把这四件套写成 `Foundation` 侧基础设施。
6. 运行期协作链路固定为：
   - 正向：`5 个业务 Service -> RabbitMQ -> Core Agent Service -> KWeaver Core`
   - 反向：`KWeaver Core / Decision Agent -> Skill / CLI -> 业务 Service 或 Foundation 执行链路`
7. `Core Agent Service` 是内部事件消费、能力映射与 SDK 转发层。
   - 不是北向业务 Web API。
   - 不是 5 个业务 Service 之一。
   - 运行期不应让业务 Service 直接依赖 KWeaver SDK。
8. 最终收包交付必须遵守 `docs/00_设计文档/01_V9服务仓库交付规范.md`。
   - 一个可部署服务必须同时交付代码、`Dockerfile`、仓内 Helm Chart 和统一健康检查约定。
   - 只有 Chart、没有代码和 Dockerfile 的仓库，不视为可收货的最终服务仓。

## KWeaver Deployment Position

- `E:\Code\kweaver-core\deploy` 是 `KWeaver Core` 的上游部署入口，要把它当作组件部署事实来源。
- 上游 deploy 默认会处理：
  - 单节点 Kubernetes
  - `local-path`
  - `ingress-nginx`
  - MariaDB、Redis、Kafka、ZooKeeper、OpenSearch 等依赖
- 集成到 AnyBackup V9 时，不要无意识照搬这些默认值。
  - 先判断哪些能力由 V9 集群统一提供。
  - 再决定哪些安装步骤需要关闭、绕开、保留还是延后。
- 优先采用“渲染配置 + 调用上游入口”的方式对接 KWeaver 部署，不要随意 fork `E:\Code\kweaver-core\deploy` 的脚本逻辑，除非用户明确要求。
- 如果上游 deploy 的默认行为会覆盖 V9 的集群、存储、Ingress 或中间件边界，先收口边界，再落实现。

## Working Mode

1. 先判断任务类型：
   - 架构对齐
   - 部署/打包落地
   - KWeaver 集成
   - 文档回写
2. 按最小必要集合读取架构文档、ADR 和上游 deploy 资产。
3. 先检查是否触碰以下硬边界：
   - SaaS 与 Foundation 的部署边界
   - AI 决策与业务执行边界
   - Agent 侧中间件职责边界
   - `Core Agent Service` 与 `KWeaver Core` 的分层
   - KWeaver 上游 deploy 默认值与 V9 所有权边界
4. 再决定具体工作由谁承接：
   - 中间件：`v9-middleware-deployer`
   - KWeaver 集成：`kweaver-core-integrator`
   - 发布包：`v9-release-packager`
   - 验收：`v9-deployment-verifier`

## Conflict Handling

- 默认冲突优先级：
  1. 已确认 ADR
  2. 专项架构设计文档
  3. `架构总览.md`
  4. 当前仓库实现
- 如果 `E:\Code\kweaver-core\deploy` 与 V9 目标架构冲突，不要直接把上游实现当最终答案。
  - 先保留上游作为组件部署事实。
  - 再通过配置、编排顺序和边界裁剪把它纳入 V9。
- 只有当用户明确说“实现现状就是新基线”时，才把文档对齐方向反过来。

## Guardrails

- 不要把临时 MVP 做法写成长期架构原则。
- 不要把平台外运维工具链当成产品运行时组件。
- 不要把 Foundation 设备认证混进 `Auth Service`；Foundation 控制面与接入链路应按架构文档走 `Platform Governance` 与网关链路。
- 不要把 KWeaver 本地 CLI、TUI 或 SKILL 入口误写成 `Core Agent Service` 的运行期主依赖；服务内主依赖应保持在 SDK 与资源能力层。
- 不要只因为上游 deploy 自带某项依赖，就自动认定它属于 V9 的正式边界。

## Done Criteria

- 改动后的设计、脚本、chart、Ansible 或文档，与 `E:\Code\aa-feat\docs\05-技术架构` 一致，或者差异被明确标记为版本差异。
- `KWeaver Core` 被清晰表述为 `AnyBackup V9` 的组成组件，并有明确的部署归位方式。
- 系统边界、调用边界、中间件边界和部署边界没有互相打架。
- 实现说明、部署入口和文档术语使用的是同一套口径。
