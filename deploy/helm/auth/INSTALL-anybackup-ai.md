# Auth Service Manual Deployment

This guide installs the packaged Auth Service chart into namespace `anybackup-ai`.

## Assumptions

- Helm package: `auth-service-0.1.6.tgz`
- Namespace: `anybackup-ai`
- Release name: `auth-service`
- Keycloak image: `docker.aityp.com/image/docker.io/keycloak/keycloak:26.5.1`
- PostgreSQL is already deployed outside this chart.
- PostgreSQL host: `postgres.middleware`
- PostgreSQL port: `5432`
- PostgreSQL database: `keycloak`
- PostgreSQL user: `postgres`

Auth Service exposes only internal Kubernetes services. API Gateway is the only external entry and routes `/api/auth_service/v1` to Auth Service.

## 1. Check Cluster Access

```bash
kubectl config current-context
kubectl get nodes
helm version --short
```

## 2. Prepare Package Directory

```bash
mkdir -p /opt/anybackup-ai/helm
cd /opt/anybackup-ai/helm
ls -lh auth-service-0.1.6.tgz
helm lint ./auth-service-0.1.6.tgz
```

## 3. Create Namespace And Secret

```bash
kubectl create namespace anybackup-ai

kubectl -n anybackup-ai create secret generic auth-service-keycloak-secrets \
  --from-literal=admin-password='<replace-with-admin-password>' \
  --from-literal=database-password='<replace-with-keycloak-db-password>'
```

## 4. Create Install Values

```bash
cat > auth-service-values.yaml <<EOF
keycloak:
  replicaCount: 1

  image:
    registry: docker.aityp.com
    repository: image/docker.io/keycloak/keycloak
    tag: 26.5.1
    pullPolicy: IfNotPresent

  database:
    vendor: postgres
    host: postgres.middleware
    port: 5432
    name: keycloak
    schema: public
    username: postgres
    jdbcUrl: ""

  secrets:
    create: false
    name: auth-service-keycloak-secrets
    adminPasswordKey: admin-password
    databasePasswordKey: database-password

  bootstrapAdmin:
    username: admin

  realm:
    name: master
    accessTokenLifespan: 1800

  http:
    enabled: true
    port: 8080
    managementPort: 9000
    relativePath: /api/auth_service/v1

  proxy:
    headers: xforwarded

  features:
    enabled: []
    disabled:
      - admin
      - account

  health:
    enabled: true

  metrics:
    enabled: true

  service:
    type: ClusterIP
    ports:
      http: 80
      management: 9000
EOF
```

## 5. Dry Run And Install

```bash
helm upgrade --install auth-service ./auth-service-0.1.6.tgz \
  --namespace anybackup-ai \
  -f auth-service-values.yaml \
  --dry-run --debug

helm upgrade --install auth-service ./auth-service-0.1.6.tgz \
  --namespace anybackup-ai \
  -f auth-service-values.yaml
```

## 6. Check Deployment

```bash
helm status auth-service -n anybackup-ai
kubectl get pods -n anybackup-ai
kubectl rollout status deployment/auth-service-auth-service -n anybackup-ai
kubectl logs -n anybackup-ai deployment/auth-service-auth-service -f
kubectl get job -n anybackup-ai auth-service-auth-service-realm-config
```

## 7. Verify Internal Service

```bash
kubectl -n anybackup-ai port-forward svc/auth-service-auth-service 8080:80
curl -i http://127.0.0.1:8080/api/auth_service/v1/realms/master/.well-known/openid-configuration
curl -i http://127.0.0.1:8080/api/auth_service/v1/realms/master/protocol/openid-connect/token
```

A `400` or `405` response from the token endpoint is acceptable for this check because no grant request body was sent.

## 8. Verify API Gateway Access

API Gateway must be installed and configured with the Auth Service route enabled.

```bash
curl -i http://192.168.40.107/api/auth_service/v1/realms/master/.well-known/openid-configuration
curl -i http://192.168.40.107/api/auth_service/v1/realms/master/protocol/openid-connect/certs
curl -i http://192.168.40.107/api/auth_service/v1/realms/master/protocol/openid-connect/userinfo
```

Without a valid `Authorization: Bearer ...` header, `userinfo` should return `401`.

## 9. Uninstall

```bash
helm uninstall auth-service -n anybackup-ai
kubectl delete secret auth-service-keycloak-secrets -n anybackup-ai
```

Only delete the namespace if no other workloads use it:

```bash
kubectl delete namespace anybackup-ai
```

## Notes

- `KC_BOOTSTRAP_ADMIN_USERNAME` and `KC_BOOTSTRAP_ADMIN_PASSWORD` initialize the admin user only when Keycloak first creates its database state.
- If the PostgreSQL database already contains initialized Keycloak data, changing the Kubernetes Secret does not reset the existing admin password.
- Kubernetes image references must not include `https://`; use `docker.aityp.com/image/docker.io/keycloak/keycloak:26.5.1`.
- `/api/auth_service/v1` is Keycloak's configured HTTP relative path. External routing is owned by API Gateway.
