# Design Doc Source Map

按最小必要原则读取事实来源,不要一次把整仓库所有部署和文档文件都塞进上下文。

## Target Documents

当前目录:

- `docs/00_设计文档/00_V9部署总体架构设计.md`
  - 适合承接总体部署拓扑、命名空间归属、部署链路总览和人工步骤边界。

后续如需拆分,优先沿下面的落点扩展,而不是重命名现有 `00_` 文档:

- `docs/00_设计文档/01_部署时序与执行责任设计.md`
  - 控制机、客户目标机、Ansible、Helm、人工初始化步骤的时序与职责。
- `docs/00_设计文档/02_运行时组件与依赖矩阵.md`
  - 服务到 PostgreSQL / RabbitMQ / Redis / OpenSearch 的依赖矩阵和关键环境变量。
- `docs/00_设计文档/03_命名空间与资源边界设计.md`
  - `v9-system`、`v9`、`kweaver`、`resource` 等命名空间和资源归属。
- `docs/00_设计文档/04_部署验证与回归策略设计.md`
  - 验证路径、smoke check、回归脚本和“可交付”判断标准。

## Source of Truth

### Repository entrypoints

- `README.md`
- `install.sh`
- `scripts/validate_deploy.py`

### Ansible

- `ansible/site.yml`
- `ansible/group_vars/all.yml`
- `ansible/inventory.ini`
- `ansible/roles/internal/preflight/tasks/main.yml`
- `ansible/roles/internal/v9_infra/tasks/main.yml`
- `ansible/roles/deploy-services/tasks/main.yml`
- `ansible/roles/internal/release/tasks/main.yml`
- `ansible/roles/internal/kweaver_core/tasks/main.yml`
- `ansible/roles/internal/manual_init_gate/tasks/main.yml`
- `ansible/roles/internal/build_import/tasks/main.yml`
- `ansible/roles/internal/verify/tasks/main.yml`

### Helm charts

- `helm-chart/anybackup-agent/Chart.yaml`
- `helm-chart/anybackup-agent/charts/v9-infra/values.yaml`
- `helm-chart/anybackup-agent/charts/v9-infra/templates/postgres.yaml`
- `helm-chart/anybackup-agent/charts/v9-infra/templates/rabbitmq.yaml`
- `helm-chart/anybackup-agent/charts/v9-infra/templates/redis.yaml`
- `helm-chart/anybackup-agent/charts/v9-infra/templates/opensearch.yaml`
- `helm-chart/anybackup-agent/charts/core-agent-service/templates/deployment.yaml`
- `helm-chart/anybackup-agent/charts/v9-services/templates/deployment.yaml`

### Runtime and verification companions

- `config/kweaver-config.yaml.j2`
- `docs/操作文档/cluster-smoke.sh`

## Decision Hints

- 变更了中间件数量、镜像、PVC、Service 名称:
  - 先读 `v9-infra` chart 和 `ansible/roles/internal/v9_infra/tasks/main.yml`
  - 再判断是否要同步总体架构图和依赖矩阵。
- 变更了环境变量、服务依赖或 Core Agent 连接:
  - 先读 `core-agent-service` 和 `v9-services` deployment 模板
  - 再判断是否要同步运行时依赖说明。
- 变更了安装顺序、人工步骤、前置检查:
  - 先读 `install.sh`、`site.yml`、`preflight`、`manual_init_gate`
  - 再判断是否要同步部署时序或验证策略说明。
- 变更了验证脚本或“是否可部署”的判断:
  - 先读 `scripts/validate_deploy.py` 和 `docs/操作文档/cluster-smoke.sh`
  - 再同步验证设计文档。

## Writing Rules

- 先更新已有设计文档,再考虑新增。
- 新增文档前,先确认是否只是当前文档的一个新章节。
- 若代码事实与旧文档冲突,优先修正文档,不要替旧文档找借口。
