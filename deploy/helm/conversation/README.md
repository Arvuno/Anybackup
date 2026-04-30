# Conversation Service Helm Chart

This local chart deploys the AnyBackup AI Conversation Service.

## Files

- `Chart.yaml`: chart metadata.
- `values.yaml`: image, runtime settings, Secret references, Service, and probes.
- `templates/deployment.yaml`: FastAPI/uvicorn Deployment.
- `templates/service.yaml`: internal ClusterIP Service.
- `templates/ingressroute.yaml`: service-owned Traefik route rules.
- `templates/migrations-job.yaml`: Helm hook Job that runs `alembic upgrade head`.
- `templates/secret.yaml`: optional Secret for local or test environments.

## Authentication Boundary

Conversation Service is a business service. It does not validate Bearer tokens, call Auth Service, or deploy an auth middleware.

This chart owns the Conversation Service route for `/api/conversation_service/v1`. The route references shared API Gateway middleware, which authenticates callers before forwarding and injects `X-User`. The `X-User` value must be a JSON object with at least `sub`, for example `{"sub":"user-001","roles":["backup_admin"]}`.

## Secrets

By default, this chart expects an existing Kubernetes Secret and does not store credentials in Git.

```bash
kubectl -n anybackup-ai create secret generic conversation-service-secrets \
  --from-literal=database-url='postgresql+asyncpg://conversation:<password>@postgres.middleware:5432/conversation' \
  --from-literal=rabbitmq-url='amqp://<user>:<password>@rabbitmq.middleware:5672/' \
  --from-literal=redis-url='redis://redis.middleware:6379/0' \
  --from-literal=core-agent-service-token='<replace-with-service-token>'
```

For non-production environments, the chart can create the Secret:

```bash
helm upgrade --install conversation-service src/helm/conversation_service_chart \
  --namespace anybackup-ai \
  --create-namespace \
  --set secrets.create=true \
  --set secrets.databaseUrl='postgresql+asyncpg://conversation:conversation@postgres.middleware:5432/conversation' \
  --set secrets.rabbitmqUrl='amqp://guest:guest@rabbitmq.middleware:5672/' \
  --set secrets.redisUrl='redis://redis.middleware:6379/0' \
  --set secrets.coreAgentServiceToken='<replace-with-service-token>'
```

## Install

```bash
helm upgrade --install conversation-service src/helm/conversation_service_chart \
  --namespace anybackup-ai \
  --create-namespace \
  -f src/helm/conversation_service_chart/values.yaml
```

## Verify

```bash
helm lint src/helm/conversation_service_chart
helm template conversation-service src/helm/conversation_service_chart
kubectl rollout status deployment/conversation-service-conversation-service -n anybackup-ai
```

Internal health checks use `/healthz` and `/readyz`. The default API gateway path is:

```text
/api/conversation_service/v1
```

## License and Third-Party Notices

- This chart's source is distributed under the repository root [LICENSE](../../../LICENSE) (SSPL-1.0) together with the root [NOTICE](../../../NOTICE).
- The Python base image and runtime dependencies (FastAPI, SQLAlchemy, aio-pika, asyncpg, etc.) bundled into the Conversation Service image shipped by this chart are declared uniformly in [`THIRD_PARTY_NOTICES.md`](./THIRD_PARTY_NOTICES.md).
- When upgrading the base image or any runtime dependency, the Inventory table of the above notice file MUST be updated in the same commit.
- PostgreSQL, RabbitMQ, and Redis are accessed via pure protocol-level invocations and do not trigger any third-party open-source component notice obligation; the Keycloak image is declared separately by `deploy/helm/auth`.
