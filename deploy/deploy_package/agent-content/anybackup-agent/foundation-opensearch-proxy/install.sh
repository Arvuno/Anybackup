#!/usr/bin/env bash
set -euo pipefail

NAMESPACE="${FOUNDATION_OPENSEARCH_PROXY_NAMESPACE:-v9-system}"
NAME="${FOUNDATION_OPENSEARCH_PROXY_NAME:-v9-infra-foundation-opensearch-proxy}"
SERVICE_NAME="${FOUNDATION_OPENSEARCH_PROXY_SERVICE_NAME:-$NAME}"
IMAGE="${FOUNDATION_OPENSEARCH_PROXY_IMAGE:-swr.cn-east-3.myhuaweicloud.com/kweaver-ai/ingress-nginx/controller:v1.14.1}"
SERVICE_PORT="${FOUNDATION_OPENSEARCH_PROXY_SERVICE_PORT:-9896}"
UPSTREAM_HOST="${FOUNDATION_OPENSEARCH_PROXY_UPSTREAM_HOST:-192.168.40.109}"
UPSTREAM_PORT="${FOUNDATION_OPENSEARCH_PROXY_UPSTREAM_PORT:-9895}"
MANIFEST_FILE="${FOUNDATION_OPENSEARCH_PROXY_MANIFEST_FILE:-/opt/v9-alpha-deploy/generated/foundation-opensearch-http-proxy.yaml}"

usage() {
  cat <<'USAGE'
Usage:
  install.sh [options]

Options:
  --namespace NS       Kubernetes namespace. Default: v9-system
  --name NAME          Deployment/ConfigMap name. Default: v9-infra-foundation-opensearch-proxy
  --service-name NAME  Service name. Default: same as --name
  --image IMAGE        Nginx-capable image. Default: KWeaver ingress-nginx controller image
  --service-port PORT  ClusterIP HTTP port. Default: 9896
  --upstream-host HOST Foundation OpenSearch HTTPS host. Default: 192.168.40.109
  --upstream-port PORT Foundation OpenSearch HTTPS port. Default: 9895
  --manifest-file FILE Rendered manifest path.
  -h, --help           Show this help.

This deploys an internal ClusterIP Nginx HTTP proxy:
  http://<service>.<namespace>.svc.cluster.local:<service-port>
    -> https://<upstream-host>:<upstream-port>

The proxy does not store OpenSearch credentials. It forwards Authorization and
other request headers unchanged; Vega still owns username/password config.
USAGE
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --namespace)
      NAMESPACE="${2:-}"
      shift 2
      ;;
    --name)
      NAME="${2:-}"
      shift 2
      ;;
    --service-name)
      SERVICE_NAME="${2:-}"
      shift 2
      ;;
    --image)
      IMAGE="${2:-}"
      shift 2
      ;;
    --service-port)
      SERVICE_PORT="${2:-}"
      shift 2
      ;;
    --upstream-host)
      UPSTREAM_HOST="${2:-}"
      shift 2
      ;;
    --upstream-port)
      UPSTREAM_PORT="${2:-}"
      shift 2
      ;;
    --manifest-file)
      MANIFEST_FILE="${2:-}"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      usage >&2
      exit 2
      ;;
  esac
done

for value in "$SERVICE_PORT" "$UPSTREAM_PORT"; do
  [[ "$value" =~ ^[0-9]+$ ]] || {
    echo "Port must be numeric: $value" >&2
    exit 2
  }
done

command -v kubectl >/dev/null 2>&1 || {
  echo "kubectl is required" >&2
  exit 1
}

mkdir -p "$(dirname "$MANIFEST_FILE")"

cat >"$MANIFEST_FILE" <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: ${NAME}
  namespace: ${NAMESPACE}
  labels:
    app.kubernetes.io/name: ${NAME}
    app.kubernetes.io/part-of: v9-infra
data:
  nginx.conf: |
    worker_processes 1;
    pid /tmp/nginx.pid;

    events {
      worker_connections 1024;
    }

    http {
      access_log /dev/stdout;
      error_log /dev/stderr info;
      client_body_temp_path /tmp/client_temp;
      proxy_temp_path /tmp/proxy_temp;

      server {
        listen ${SERVICE_PORT};

        location / {
          proxy_pass https://${UPSTREAM_HOST}:${UPSTREAM_PORT};
          proxy_ssl_server_name on;
          proxy_ssl_verify off;
          proxy_set_header Host ${UPSTREAM_HOST}:${UPSTREAM_PORT};
          proxy_set_header Authorization \$http_authorization;
          proxy_set_header X-Real-IP \$remote_addr;
          proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
          proxy_set_header X-Forwarded-Proto http;
          proxy_http_version 1.1;
          proxy_connect_timeout 10s;
          proxy_send_timeout 60s;
          proxy_read_timeout 60s;
          proxy_buffering off;
        }
      }
    }
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${NAME}
  namespace: ${NAMESPACE}
  labels:
    app.kubernetes.io/name: ${NAME}
    app.kubernetes.io/part-of: v9-infra
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: ${NAME}
  template:
    metadata:
      labels:
        app.kubernetes.io/name: ${NAME}
        app.kubernetes.io/part-of: v9-infra
    spec:
      containers:
        - name: nginx
          image: ${IMAGE}
          imagePullPolicy: IfNotPresent
          command:
            - /usr/local/nginx/sbin/nginx
          args:
            - -c
            - /etc/nginx/nginx.conf
            - -g
            - daemon off;
          ports:
            - name: http
              containerPort: ${SERVICE_PORT}
          volumeMounts:
            - name: config
              mountPath: /etc/nginx/nginx.conf
              subPath: nginx.conf
      volumes:
        - name: config
          configMap:
            name: ${NAME}
---
apiVersion: v1
kind: Service
metadata:
  name: ${SERVICE_NAME}
  namespace: ${NAMESPACE}
  labels:
    app.kubernetes.io/name: ${NAME}
    app.kubernetes.io/part-of: v9-infra
spec:
  type: ClusterIP
  selector:
    app.kubernetes.io/name: ${NAME}
  ports:
    - name: http
      port: ${SERVICE_PORT}
      targetPort: http
EOF

kubectl get namespace "$NAMESPACE" >/dev/null
kubectl apply -f "$MANIFEST_FILE"
kubectl rollout restart "deployment/${NAME}" -n "$NAMESPACE"
kubectl rollout status "deployment/${NAME}" -n "$NAMESPACE" --timeout=120s

echo "Foundation OpenSearch HTTP proxy is ready:"
echo "  http://${SERVICE_NAME}.${NAMESPACE}.svc.cluster.local:${SERVICE_PORT}"
echo "  upstream: https://${UPSTREAM_HOST}:${UPSTREAM_PORT}"
echo "  manifest: ${MANIFEST_FILE}"
