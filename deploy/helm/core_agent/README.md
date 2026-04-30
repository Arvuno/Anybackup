# core_agent_service_chart

## 仓库说明

这是 `core_agent_service` 的 Helm 部署仓库，用于承载当前服务的 Kubernetes Chart 内容。
当前 Chart 按“MQ worker”形态设计，默认不对外暴露 Service 或 Ingress。

## Chart 目录结构

本仓库根目录就是 Chart 根目录，直接包含以下内容：

- `Chart.yaml`
- `values.yaml`
- `templates/`

不再保留 `deploy/helm/core-agent-service` 这类嵌套层级。

## values 覆盖项

部署时可通过 `values.yaml` 或安装命令中的覆盖参数，调整以下配置：

- 镜像参数
- Secret 参数
- KWeaver 参数
- `.env` 兜底挂载参数
- 节点本地 `~/.kweaver` 挂载参数

## 关键说明

- 当前服务通过 RabbitMQ 与外部系统交互，默认不需要 Service 和 Ingress。
- PostgreSQL、RabbitMQ 和 KWeaver 凭据通过 `secrets` 配置和 `secretKeyRef` 注入。
- KWeaver 鉴权优先使用用户名密码，其次使用显式 token，最后才回退节点本地 `~/.kweaver` 挂载。
- 如已配置用户名密码，则不再需要挂载 `/root/.kweaver`。
- 如需挂载节点本地 `~/.kweaver`，请在 values 中使用绝对路径配置 `kweaverHostMount.hostPath`。
- `kweaverHostMount.enabled` 默认关闭，只有在需要复用节点本地 ConfigAuth 登录态时才开启。

## 当前状态

当前仓库已完成 Helm Chart 迁移与 MQ worker 形态收敛，但尚未进行 Kubernetes 联调验证。

## 部署脚本

仓库提供 Bash 部署脚本：

- `scripts/deploy.sh`

可通过以下命令查看帮助：

```bash
bash scripts/deploy.sh --help
```

脚本入口约定：

- 通过 `--image` 传完整镜像地址，例如 `registry.example.com/core-agent-service:107`
- `--database-url` 默认值为 `postgresql+asyncpg://conversation:conversation@postgres.middleware:5432/conversation`
- `--rabbitmq-url` 默认值为 `amqp://guest:guest@rabbitmq.middleware:5672/`
- `--rabbitmq-exchange` 默认值为 `conversation.message.events`
- `--rabbitmq-queue` 默认值为 `core_agent.message.events`
- `--rabbitmq-status-exchange` 默认值为 `core_agent.run_status.events`
- `--rabbitmq-status-routing-key` 默认值为空字符串，表示按 `core_agent.run.accepted.v1`、`core_agent.run.completed.v1`、`core_agent.run.failed.v1` 三类状态 routing key 发布
- 如需显式用户名密码鉴权，可额外传入 `--kweaver-username` 与 `--kweaver-password`
- `--kweaver-tls-insecure` 默认值为 `true`
- `--kweaver-probe-on-startup` 默认值为 `true`
- `--kweaver-chat-timeout` 默认不限制超时；如需限制可传入秒数
- `config.kweaverStreamProgressInterval` 默认值为 `100`，仅影响 DEBUG 级别的流式进度采样日志
- `config.coreAgentLogLevel` 调试阶段默认值为 `DEBUG`，生产环境建议覆盖为 `INFO`
- `config.kweaverStreamTraceEnabled` 调试阶段默认值为 `true`，完整流式事件写入独立 trace 文件；生产环境建议覆盖为 `false`
- `config.kweaverStreamTraceDir` 默认值为 `/tmp/core-agent-service-trace-log`
- 当 `config.kweaverStreamTraceEnabled=true` 时，Pod 会把节点本地 hostPath `/tmp/core-agent-service-trace-log` 挂载到容器内同路径，trace 文件会写入该目录

## License and Third-Party Notices

- This chart's source is distributed under the repository root [LICENSE](../../../LICENSE) (SSPL-1.0) together with the root [NOTICE](../../../NOTICE).
- The Python base image and runtime dependencies (FastAPI, SQLAlchemy, psycopg, aio-pika, kweaver-sdk, etc.) bundled into the Core Agent Service image shipped by this chart are declared uniformly in [`THIRD_PARTY_NOTICES.md`](./THIRD_PARTY_NOTICES.md).
- **LGPL notice**: the image contains `psycopg[binary]` distributed under LGPL-3.0-or-later. Any redistribution of images built from this chart MUST carry the full LGPL license text and the upstream copyright notice. See §4 of the notice file for details.
- The license of `kweaver-sdk` is pending upstream confirmation; it will be back-filled into §3.5 of the notice file.
- When upgrading the base image or any runtime dependency, the Inventory table of the above notice file MUST be updated in the same commit.
- PostgreSQL and RabbitMQ are accessed via pure protocol-level invocations and do not trigger any third-party open-source component notice obligation.
