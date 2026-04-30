# Conversation Service

Conversation Service 是 Anybackup 9 的会话域入口服务，负责对话流管理、消息时间线持久化、上下文保管和提醒协调。

## 正式文档

产品需求、架构设计、接口契约和验收标准的**唯一事实源**在根 `docs/`：

| 文档 | 用途 |
|------|------|
| [版本索引](../../../../docs/06-版本内容/9.0.0-alpha1/版本索引.md) | alpha1 版本总入口 |
| [Conversation Service 专项设计](../../../../docs/06-版本内容/9.0.0-alpha1/03-版本架构/服务方案-Conversation%20Service专项设计.md) | 会话生命周期、北向南向规则、提醒机制、数据模型 |
| [服务能力边界与接口契约](../../../../docs/06-版本内容/9.0.0-alpha1/03-版本架构/服务方案-服务能力边界与接口契约.md) | 四个 Service 能力边界、API 最小集合、完成定义 |
| [会话领域模型](../../../../docs/06-版本内容/9.0.0-alpha1/03-版本架构/数据状态-会话领域模型.md) | 领域对象、状态机、轮次闭环 |
| [开发任务拆解](../../../../docs/06-版本内容/9.0.0-alpha1/03-版本架构/工程交付-开发任务拆解与完成定义.md) | WP-00 ~ WP-08 开发工作包 |
| [测试策略与质量门禁](../../../../docs/06-版本内容/9.0.0-alpha1/03-版本架构/工程交付-测试策略与质量门禁.md) | 测试分层与质量门禁 |

本 README **只**提供本地运行、测试和 OpenAPI 调试说明，不定义产品范围、版本验收或架构事实。

## 技术栈

| 组件 | 版本 |
|------|------|
| Python | 3.12+ |
| FastAPI | 0.128+ |
| SQLAlchemy | 2.0+ (async) |
| Alembic | 1.17+ |
| dependency-injector | 4.48+ |
| RabbitMQ | 3.12+ (aio-pika) |
| Redis | 7+ |
| PostgreSQL | 15+ (asyncpg) |
| structlog | 25+ |
| OpenTelemetry | 1.39+ |

## 本地环境

### 前置条件

- Python 3.12
- uv (Python 包管理)
- Docker (用于本地 PostgreSQL / RabbitMQ / Redis)

### 依赖服务

```powershell
# PostgreSQL
docker run --name conversation-postgres -e POSTGRES_USER=conversation -e POSTGRES_PASSWORD=conversation -e POSTGRES_DB=conversation -p 5432:5432 -d postgres:16

# RabbitMQ (带管理界面)
docker run --name conversation-rabbitmq -p 5672:5672 -p 15672:15672 -d rabbitmq:4.2-management

# Redis
docker run --name conversation-redis -p 6379:6379 -d redis:7
```

### 环境变量

```text
CONVERSATION_SERVICE_ENV=local
DATABASE_URL=postgresql+asyncpg://conversation:conversation@localhost:5432/conversation
REDIS_URL=redis://localhost:6379/0
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
SNOWFLAKE_NODE_ID=1
SNOWFLAKE_EPOCH_MS=1735689600000
AUTH_CONTEXT_SOURCE=nginx_ingress_x_user
AUTH_USER_HEADER_NAME=X-User
```

本地调试用 `X-User` 最小样例：

```json
{"sub":"user-001","preferred_username":"backup-admin","name":"Backup Admin","email":"backup-admin@example.com","email_verified":true,"roles":[]}
```

## 快速启动

```powershell
uv sync
uv run alembic upgrade head
uv run uvicorn app.main:app --reload --host 127.0.0.1 --port 8000
```

OpenAPI 文档：

- Swagger UI: http://127.0.0.1:8000/docs
- ReDoc: http://127.0.0.1:8000/redoc
- OpenAPI Schema: http://127.0.0.1:8000/openapi.json

## 常用命令

| 命令 | 用途 |
|------|------|
| `uv sync` | 安装依赖 |
| `uv run pytest tests/unit` | 单元测试 |
| `uv run pytest tests/integration` | 集成测试 |
| `uv run pytest tests/contract` | 契约测试 |
| `uv run pytest` | 全量测试 |
| `uv run ruff check .` | Lint |
| `uv run ruff format --check .` | 格式检查 |
| `uv run mypy app` | 类型检查 |
| `uv run alembic revision --autogenerate -m "message"` | 生成迁移 |
| `uv run alembic upgrade head` | 升级数据库 |

## 发布前检查

```powershell
uv lock --check
uv run ruff check .
uv run ruff format --check .
uv run mypy app
uv run pytest
uv run alembic upgrade head
```

## 健康检查

| Endpoint | 用途 | 通过标准 |
|----------|------|---------|
| `/healthz` | 进程存活 | 不依赖外部服务 |
| `/readyz` | 服务就绪 | DB 连接、迁移版本、MQ/Redis 连接可用 |
| `/metrics` | Prometheus 指标 | API/outbox/MQ/writeback 基础指标 |

## 运行模式

| 模式 | 说明 |
|------|------|
| API 进程 | REST API、写回 API、健康检查和 OpenAPI |
| Outbox publisher | 扫描 outbox 并发布 RabbitMQ，执行前获取 Redis 全局锁 |
| MQ consumer | 消费 Core Agent 状态事件 |
| Retention job | 自动归档和过期标记，可由 CronJob 触发 |
| Context merge worker | 合并 context delta，生成最新 context snapshot |

## License and Third-Party Notices

- This service's source code is distributed under the repository root [LICENSE](../../../LICENSE) (SSPL-1.0) together with the root [NOTICE](../../../NOTICE).
- When this service is packaged as a container image and shipped via its Helm chart, the Python base image (`python:3.12-slim`, Debian base layer) and runtime dependencies declared in `pyproject.toml` are declared uniformly on the chart side: [`deploy/helm/conversation/THIRD_PARTY_NOTICES.md`](../../../deploy/helm/conversation/THIRD_PARTY_NOTICES.md).
- When adding or upgrading any Python runtime dependency, the Inventory table of the above chart-side notice file MUST be updated in the same commit.
- Access to PostgreSQL, RabbitMQ, Redis, and Keycloak is a pure protocol-level invocation (asyncpg / aio-pika / redis-py / OIDC) and does not trigger any third-party open-source component notice obligation.
