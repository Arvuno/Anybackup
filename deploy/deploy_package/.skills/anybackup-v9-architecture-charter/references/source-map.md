# AnyBackup V9 Source Map

按最小必要原则读取外部资料，不要一次性把整套文档全量塞进上下文。

## Overall Alignment

优先用于判断长期目标架构、系统边界和术语口径：

- `E:\Code\aa-feat\docs\05-技术架构\架构总览.md`
- `E:\Code\aa-feat\docs\05-技术架构\架构设计\总体架构设计.md`
- `E:\Code\aa-feat\docs\05-技术架构\架构设计\部署架构设计.md`

## Delivery Contract

优先用于判断“最终打包人到底按什么格式收货”：

- `E:\Code\script\anybackup-agent-release\docs\00_设计文档\01_V9服务仓库交付规范.md`

重点检查这些关键词：

- `Dockerfile`
- `helm/<service-name>-chart`
- `8080`
- `/health`
- `/ready`
- `chart-only`

## ADR Baseline

优先用于判断“哪些边界已经拍板”：

- `E:\Code\aa-feat\docs\05-技术架构\架构决策\ADR-13-部署架构采用KubernetesVKE.md`
  - Kubernetes/VKE 是 SaaS 侧正式部署基线
- `E:\Code\aa-feat\docs\05-技术架构\架构决策\ADR-15-SaaS与Foundation混合部署架构.md`
  - SaaS 与 Foundation 的混合部署、反向通道、Gateway/Connection Manager 边界
- `E:\Code\aa-feat\docs\05-技术架构\架构决策\ADR-19-全局中间件选型.md`
  - `Redis / PostgreSQL / RabbitMQ / OpenSearch` 仅属于 Agent 侧

## AI And Core Boundary

优先用于判断 `Core Agent Service` 与 `KWeaver Core` 的职责分层：

- `E:\Code\aa-feat\docs\05-技术架构\核心智能体服务\Core Agent Service-架构设计.md`
- `E:\Code\aa-feat\docs\05-技术架构\核心智能体服务\KWeaver Core SDK能力说明.md`

重点检查这些关键词：

- `Core Agent Service`
- `Kweaver_Core`
- `RabbitMQ`
- `Skill / CLI`
- `SDK`
- `success / failed`

## Upstream KWeaver Deploy

优先用于判断上游 `KWeaver Core` 组件的真实部署假设，而不是直接拿来定义 V9 架构：

- `E:\Code\kweaver-core\deploy\README.zh.md`
- `E:\Code\kweaver-core\deploy\deploy.sh`
- `E:\Code\kweaver-core\deploy\scripts\services\core.sh`
- `E:\Code\kweaver-core\deploy\scripts\services\k8s.sh`

重点检查这些关键词：

- `kweaver-core install`
- `AUTO_INSTALL_LOCALPV`
- `AUTO_INSTALL_INGRESS_NGINX`
- `local-path`
- `ingress-nginx`
- `single-node Kubernetes`

## Decision Hints

- 如果任务是“这套东西是不是符合目标架构”，先读 `Overall Alignment` + `ADR Baseline`。
- 如果任务是“Core Agent / KWeaver 应该怎么分层”，再读 `AI And Core Boundary`。
- 如果任务是“怎么把上游 KWeaver 部署嵌进 V9”，再读 `Upstream KWeaver Deploy`。
- 如果实现和文档冲突，先查 ADR，再决定是收敛实现还是显式记录版本差异。
