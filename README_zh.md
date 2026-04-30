# Anybackup

<p align="center">
  <a href="https://github.com/anybackup-ai/Anybackup/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-SSPL--1.0-blue.svg" alt="License"></a>
  <a href="https://github.com/anybackup-ai/Anybackup/blob/main/VERSION.txt"><img src="https://img.shields.io/badge/version-9.0.0--alpha-orange.svg" alt="Version"></a>
</p>

<p align="center">
  <em>AI-Native 数据韧性平台 — 自主备份、自主恢复、自主优化。基于开源商业模式，总体拥有成本降低 35%。</em>
</p>

<p align="center">
  <a href="./README.md">English</a>
</p>

---

## 📖 目录

- [产品概述](#产品概述)
- [核心能力](#核心能力)
- [工作机制](#工作机制)
- [产品架构](#产品架构)
- [快速开始](#快速开始)
- [仓库结构](#仓库结构)
- [社区](#社区)
- [参与贡献](#参与贡献)
- [相关项目](#相关项目)
- [许可证与第三方声明](#许可证与第三方声明)

---

## 产品概述

作为更经济、智能的数据韧性平台，Anybackup V9 基于开源商业模式，为客户实现业务所需的数据韧性保障。Anybackup Agent 打造的 AI 备份管理员，实现自主备份、自主恢复、自主优化，总体拥有成本可降低 35%，告别被动响应。

Anybackup V9 以 AI 原生思维重塑，以智能体 Agent 为中心。整体架构为 **采集 · 执行 · 决策 三位一体**，各层职责分明：

- **Anybackup Agent**（决策） — 负责思考与规划的智能体：意图识别、方案生成、风险评估、审批控制、全程审计追踪。运行在 SaaS 侧。
- **Anybackup Foundation**（执行） — 负责执行的备份引擎：全量/增量/日志备份、保留策略、恢复编排。
- **Anybackup Client**（采集） — 按工作负载部署的采集代理：运行在被保护资产上。
---

## 核心能力

- **自然语言备份推荐** — 描述业务需求，系统输出 3 套按 RPO/RTO 排序的候选方案供你挑选
- **灵活的备份粒度** — 支持分钟级、天级、周级全量/增量/日志备份，适配不同工作负载
- **应用识别与路由** — 同一资源下存在多实例时，自动识别并引导选择目标实例后再推荐
- **候选方案闭环** — 从 3 套方案中选定一套，系统记录选择并关闭推荐流程
- **故障场景识别** — 识别实例崩溃、库级误删、日志损坏、勒索攻击等常见故障类型，恢复范围由管理员在对话中确认
- **多粒度恢复路径** — 支持实例级恢复、库级恢复、时间点恢复、日志文件恢复
- **Kubernetes 原生部署** — 所有 Agent 侧服务与基础设施通过 Helm Chart 部署
- **MySQL 首批支持** — XtraBackup、实时日志、多通道备份，更多工作负载陆续支持


---

## 工作机制

```
  用户开口                  Agent 推理               Foundation 执行
  ────────               ──────────────────        ──────────────────────
  "给我们的      →        意图识别               →  查询当前备份配置
   数据库推荐个           方案生成                   与 Foundation 能力
   备份方案"              风险评估
                               │
                     ┌─────────┴─────────┐
                     ▼                   ▼
               3 个候选方案         你选定一个
               按 RPO/RTO 排序      → Agent 记录选择，
               与保留周期对比          下发至 Foundation
```

**示例 — 备份推荐：**

```
用户:   "帮我们的金融数据库推荐个备份方案。"

Agent:  基于你的环境，三套候选方案：

        A — 常规保护
        周级全备 + 天级增量，保留 90 天。

        B — 低 RPO（★ 推荐）
        周级全备 + 天级增量 + 分钟级日志备份，保留 5 年。
        开启实时日志备份。

        C — 最简模式
        每日全备，保留 90 天。
```

**示例 — 实例崩溃恢复：**

```
用户:   "生产数据库崩了，起不来，所有库都访问不了。"

Agent:  实例级故障。

        恢复方案：将整个实例恢复到最近一次有效备份状态，
        恢复目标为原生产主机。
        ⚠️  此操作将覆盖主机当前数据。[ 确认执行 ]  [ 取消 ]
```

---

## 产品架构

```
  ┌───────────────────────────────────────────────────┐
  │         Anybackup Agent  ·  SaaS 侧             │
  │                                                   │
  │  意图识别 · 方案生成 · 风险评估                   │
  │  审批确认 · 人工审核 · 审计追踪                   │
  └──────────────────────▬───────────────────────────┘
                         │
             skill + foundation-cli
                         │
                         ▼
  ┌───────────────────────────────────────────────────┐
  │      Anybackup Foundation  ·  你的数据中心        │
  │                                                   │
  │  备份：全量 · 增量 · 日志 · 保留策略              │
  │  恢复：实例级 · 库级 · 时间点恢复                  │
  └──────────────────────▬───────────────────────────┘
                         │
                      接入采集
                         │
                         ▼
  ┌───────────────────────────────────────────────────┐
  │        Anybackup Client  ·  被保护资产侧          │
  │                                                   │
  └───────────────────────────────────────────────────┘
```

---

## 快速开始

### 前置条件

- **Foundation**：已在你的数据中心部署并可访问（部署细节请联系我们）
- **Kubernetes**：Agent 侧服务所需的 K8s 集群
- **Helm**：3.x 或更高版本
- **目标资产**：MySQL 实例运行中，可与 Foundation 互通

### 安装步骤

**第一步 — 部署基础设施服务**

```bash
helm install v9-infra ./deploy/helm/v9-infra
```

此命令将为 Agent 栈部署 Postgres、Redis、RabbitMQ 和 OpenSearch。

**第二步 — 部署 Agent 业务服务**

```bash
# 核心 Agent 引擎
helm install core-agent ./deploy/helm/core_agent

# 会话服务
helm install conversation ./deploy/helm/conversation

# 认证服务（Keycloak）
helm install auth ./deploy/helm/auth

# API 网关（Traefik）
helm install api-gateway ./deploy/helm/api_gateway

# Web 门户
helm install web ./deploy/helm/web
```

**第三步 — 安装 CLI**

```bash
cd CLI
go build -o foundation-cli ./cmd/foundation-cli
```

**第四步 — 检查连通性**

```bash
foundation-cli version
```

**第五步 — 开始对话**

打开 Agent Web 门户，用自然语言描述你的备份需求；或直接通过 CLI 调用：

```bash
anybackup chat "给我的生产数据库推荐一个备份方案"
```

详细配置选项参见 [`deploy/helm/`](./deploy/helm/) 下各 Chart 的 `values.yaml`。

---

## 仓库结构

```
Anybackup/
├── Agent/                     # AI 决策层
│   ├── portal/                #   Web 前端（React）
│   ├── service/               #   后端服务
│   │   ├── conversation/      #     会话编排
│   │   └── core_agent/        #     AI Agent 引擎
│   ├── skills/                #   Agent 技能与 AG-UI 设计模式
│   └── knowledge/             #   知识网络（备份与恢复 BKN）
├── CLI/                       # foundation-cli — 命令行工具
├── deploy/                    # Helm Chart 与部署配置
│   └── helm/
│       ├── v9-infra/          #   基础设施（Postgres、Redis、RabbitMQ、OpenSearch）
│       ├── api_gateway/       #   API 网关（Traefik）
│       ├── auth/              #   认证服务（Keycloak）
│       ├── conversation/      #   会话服务
│       ├── core_agent/        #   核心 Agent 服务
│       └── web/               #   Web SPA
├── LICENSE                    # SSPL-1.0
├── NOTICE                     # 版权声明
└── VERSION.txt                # 当前版本号
```

---

## 社区（筹备中）

- **问题反馈**：[GitHub Issues](https://github.com/anybackup-ai/Anybackup/issues)
- **讨论交流**：[GitHub Discussions](https://github.com/anybackup-ai/Anybackup/discussions)

---

## 参与贡献

我们欢迎社区贡献。在提交 Pull Request 之前：

1. 先通过 Issue 讨论你的改动或新功能。
2. 遵循仓库的编码规范。
3. 若依赖变更，同步更新相关文档和第三方声明。
4. 确保所有测试通过。

详细贡献指南见 [CONTRIBUTING.md](./CONTRIBUTING.md)（筹备中）。

---

## 相关项目

| 项目 | 关系 |
|---|---|
| **[Agent Portal](./Agent/portal/)** | 基于对话的备份工作流 Web 前端 |
| **[foundation-cli](./CLI/)** | Foundation 控制面命令行工具 |

---

## 许可证与第三方声明

- 源代码：[LICENSE](./LICENSE) (SSPL-1.0)，附 [NOTICE](./NOTICE)
- 第三方开源组件按分发单元声明：
  - [`CLI/THIRD_PARTY_NOTICES.md`](./CLI/THIRD_PARTY_NOTICES.md) — Go 依赖（CLI 二进制）
  - [`deploy/helm/v9-infra/THIRD_PARTY_NOTICES.md`](./deploy/helm/v9-infra/THIRD_PARTY_NOTICES.md) — 基础设施（PostgreSQL、RabbitMQ、Redis、OpenSearch）
  - [`deploy/helm/api_gateway/THIRD_PARTY_NOTICES.md`](./deploy/helm/api_gateway/THIRD_PARTY_NOTICES.md) — API 网关（Traefik）
  - [`deploy/helm/auth/THIRD_PARTY_NOTICES.md`](./deploy/helm/auth/THIRD_PARTY_NOTICES.md) — 认证服务（Keycloak）
  - [`deploy/helm/conversation/THIRD_PARTY_NOTICES.md`](./deploy/helm/conversation/THIRD_PARTY_NOTICES.md) — 会话服务
  - [`deploy/helm/core_agent/THIRD_PARTY_NOTICES.md`](./deploy/helm/core_agent/THIRD_PARTY_NOTICES.md) — 核心智能体服务
  - [`deploy/helm/web/THIRD_PARTY_NOTICES.md`](./deploy/helm/web/THIRD_PARTY_NOTICES.md) — Web SPA
  - [`Agent/portal/deploy/helm/agent-web/THIRD_PARTY_NOTICES.md`](./Agent/portal/deploy/helm/agent-web/THIRD_PARTY_NOTICES.md) — Agent Web 门户

引入新镜像或升级依赖时，须在同一提交中更新对应声明文件。
