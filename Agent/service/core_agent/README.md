# core_agent_service

`core_agent_service` 是面向 RabbitMQ 的最小中转版服务。

## 文档导航

主文档：

- [需求文档](docs/01-需求文档.md)
- [技术选型文档](docs/02-技术选型文档.md)
- [架构设计文档](docs/03-架构设计文档.md)
- [实现设计文档](docs/04-实现设计文档.md)
- [测试设计文档](docs/05-测试设计文档.md)
- [开发与演进说明](docs/06-开发与演进说明.md)
- [系统架构与交互图](docs/07-系统架构与交互图.md)

附录文档：

- [上游 MQ 契约](docs/contracts/upstream-mq-contract.md)
- [下游 SDK 中转契约](docs/contracts/downstream-sdk-relay-contract.md)
- [最小存储模型](docs/storage/minimal-schema.md)
- [本地开发说明](docs/runbooks/local-dev.md)
- [Helm 部署说明](docs/runbooks/helm-deploy.md)

补充业务参考文档：

- [服务交互定义](docs/服务交互定义.md)
- [核心智能体MQ消息契约](docs/核心智能体MQ消息契约.md)
- [消息格式定义](docs/消息格式定义.md)

## 服务职责

- 消费 `conversation.message.sent.v1`
- 消费 `conversation.message.cancel_requested.v1`
- 维护 `conversation_id -> decision_conversation_id` 会话映射
- 通过 `kweaver-sdk` 转发用户消息到 `KWeaver Core Decision Agent`
- 发布 `core_agent.run.accepted.v1`
- 发布 `core_agent.run.completed.v1`
- 发布 `core_agent.run.failed.v1`

## 服务边界

当前版本只做最小消息中转，职责边界以最小 MQ 契约、最小存储模型和最小 SDK 中转边界为准。

## 语言约定

- 项目文档与代码注释使用中文
- 代码标识符、程序日志与其他程序运行时输出使用英文
- 新增内容避免中英混写，除非引用外部协议字段或第三方固定名称

## 关键环境变量

必需：

- `DATABASE_URL`，PostgreSQL 必须使用 `postgresql+psycopg://` 方言
- `RABBITMQ_URL`
- `RABBITMQ_EXCHANGE`
- `RABBITMQ_QUEUE`
- `RABBITMQ_STATUS_EXCHANGE`
- `RABBITMQ_STATUS_ROUTING_KEY`
- `KWEAVER_BASE_URL`
- `KWEAVER_DECISION_AGENT_ID`

常用可选项：

- `KWEAVER_USERNAME`
- `KWEAVER_PASSWORD`
- `KWEAVER_TOKEN`
- `KWEAVER_BUSINESS_DOMAIN`
- `KWEAVER_TIMEOUT`
- `KWEAVER_TLS_INSECURE`
- `KWEAVER_PROBE_ON_STARTUP`
- `CORE_AGENT_LOG_LEVEL`

说明：

- `KWEAVER_PROBE_ON_STARTUP` 默认是 `true`
- `KWEAVER_TLS_INSECURE` 默认是 `true`
- MQ 默认命名为：上行 Exchange `conversation.message.events`，上行队列 `core_agent.message.events`，状态 Exchange `core_agent.run_status.events`，状态队列 `conversation.core_agent.run_status`
- KWeaver 鉴权优先级为：`KWEAVER_USERNAME` + `KWEAVER_PASSWORD`，其次 `KWEAVER_TOKEN`，最后回退到本机 `~/.kweaver` 登录态
- 如使用用户名密码模式，当前实现依赖 `kweaver-sdk>=0.6.10` 提供的 `HttpSigninAuth`

## 启动方式

推荐直接启动：

```bash
python -m core_agent_service
```

如仓库提供批处理脚本，也可使用：

```bat
start_core_agent_service.bat
```

## 部署产物

- 服务代码位于当前仓库
- 容器镜像构建文件位于仓库根目录：`Dockerfile`
- 运行与部署说明文档位于：`docs/runbooks/`
- Helm Chart 的唯一事实源位于独立仓库：`core_agent_service_chart`

当前仓库负责服务代码、镜像文件和运行文档，不再承载 Helm Chart 副本；部署时应以独立仓库 `core_agent_service_chart` 中的 Chart 内容为准。

## 测试方式

运行全部单元测试：

```bash
pytest -q
```

只验证中文文档约束：

```bash
pytest -q tests/unit/test_docs_contract.py
```

## 依赖约束

必须通过以下方式安装 SDK：

```bash
pip install kweaver-sdk
```

不得把工作区中的本地 `kweaver-sdk` 源码目录当作运行时依赖。

## License and Third-Party Notices

- This service's source code is distributed under the repository root [LICENSE](../../../LICENSE) (SSPL-1.0) together with the root [NOTICE](../../../NOTICE).
- When this service is packaged as a container image and shipped via its Helm chart, the Python base image (`python:3.12-slim`, Debian base layer) and runtime dependencies declared in `pyproject.toml` (including `psycopg[binary]` under LGPL-3.0-or-later) are declared uniformly on the chart side: [`deploy/helm/core_agent/THIRD_PARTY_NOTICES.md`](../../../deploy/helm/core_agent/THIRD_PARTY_NOTICES.md).
- The license of `kweaver-sdk` is pending confirmation with the upstream project. When replacing or upgrading any Python runtime dependency, the above chart-side notice file MUST be updated in the same commit.
