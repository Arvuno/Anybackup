#!/usr/bin/env bash

SCRIPT_LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PACKAGE_ROOT="${PACKAGE_ROOT:-$(cd "${SCRIPT_LIB_DIR}/../.." && pwd)}"

REMOTE_DEPLOY_ROOT="${REMOTE_DEPLOY_ROOT:-/opt/v9-alpha-deploy}"
REMOTE_SOURCE_ROOT="${REMOTE_SOURCE_ROOT:-/opt/v9-sources}"
GENERATED_ROOT="${GENERATED_ROOT:-${REMOTE_DEPLOY_ROOT}/generated}"
REMOTE_CHART_ROOT="${REMOTE_CHART_ROOT:-${REMOTE_DEPLOY_ROOT}/charts}"
OFFLINE_ROOT="${OFFLINE_ROOT:-${REMOTE_DEPLOY_ROOT}/offline}"

PROTON_REGISTRY_PORT="${PROTON_REGISTRY_PORT:-5000}"
PROTON_CHARTMUSEUM_PORT="${PROTON_CHARTMUSEUM_PORT:-8081}"
PROTON_PORT_WAIT_ATTEMPTS="${PROTON_PORT_WAIT_ATTEMPTS:-60}"
PROTON_PORT_WAIT_SECONDS="${PROTON_PORT_WAIT_SECONDS:-2}"
INGRESS_HTTP_NODE_PORT="${INGRESS_HTTP_NODE_PORT:-30080}"
INGRESS_HTTPS_NODE_PORT="${INGRESS_HTTPS_NODE_PORT:-30443}"
INGRESS_CLASS_NAME="${INGRESS_CLASS_NAME:-class-443}"

V9_INFRA_NAMESPACE="${V9_INFRA_NAMESPACE:-v9-system}"
V9_INFRA_RELEASE="${V9_INFRA_RELEASE:-v9-infra}"
V9_POSTGRES_SECRET_NAME="${V9_POSTGRES_SECRET_NAME:-v9-postgres-auth}"
V9_RABBITMQ_SECRET_NAME="${V9_RABBITMQ_SECRET_NAME:-v9-rabbitmq-auth}"
V9_REDIS_SECRET_NAME="${V9_REDIS_SECRET_NAME:-v9-redis-auth}"
V9_POSTGRES_USERNAME="${V9_POSTGRES_USERNAME:-v9}"
V9_POSTGRES_PASSWORD="${V9_POSTGRES_PASSWORD:-CHANGE_ME}"
V9_POSTGRES_DATABASE="${V9_POSTGRES_DATABASE:-v9}"
V9_RABBITMQ_USERNAME="${V9_RABBITMQ_USERNAME:-v9}"
V9_RABBITMQ_PASSWORD="${V9_RABBITMQ_PASSWORD:-CHANGE_ME}"
V9_RABBITMQ_ERLANG_COOKIE="${V9_RABBITMQ_ERLANG_COOKIE:-CHANGE_ME}"
V9_REDIS_PASSWORD="${V9_REDIS_PASSWORD:-CHANGE_ME}"

KWEAVER_VERSION="${KWEAVER_VERSION:-0.6.0}"
KWEAVER_NAMESPACE="${KWEAVER_NAMESPACE:-kweaver}"
KWEAVER_RESOURCE_NAMESPACE="${KWEAVER_RESOURCE_NAMESPACE:-resource}"
KWEAVER_CONFIG_DIR="${KWEAVER_CONFIG_DIR:-${REMOTE_DEPLOY_ROOT}/kweaver-config}"
KWEAVER_CONFIG_PATH="${KWEAVER_CONFIG_PATH:-${KWEAVER_CONFIG_DIR}/config.yaml}"
KWEAVER_REPO_URL="${KWEAVER_REPO_URL:-https://github.com/kweaver-ai/kweaver-core.git}"
KWEAVER_GIT_REF="${KWEAVER_GIT_REF:-release/0.6.0}"
KWEAVER_REPO_DIR="${KWEAVER_REPO_DIR:-${REMOTE_SOURCE_ROOT}/kweaver-core}"
KWEAVER_REMOTE_DEPLOY_DIR="${KWEAVER_REMOTE_DEPLOY_DIR:-${KWEAVER_REPO_DIR}/deploy}"
KWEAVER_STATE_CONFIGMAP="${KWEAVER_STATE_CONFIGMAP:-kweaver-online-install-state}"

PROTON_PACKAGE_LOCAL="${PROTON_PACKAGE_LOCAL:-${PACKAGE_ROOT}/offline-packages/proton/proton-offline-package-24768920921-amd64.tar}"
PROTON_CLI_LOCAL="${PROTON_CLI_LOCAL:-${PACKAGE_ROOT}/bin/proton-cli/proton-cli/bin/proton-cli}"
PROTON_CLI_REMOTE="${PROTON_CLI_REMOTE:-/usr/local/bin/proton-cli}"
PROTON_OFFLINE_DIR="${PROTON_OFFLINE_DIR:-${OFFLINE_ROOT}/proton}"
PROTON_PACKAGE_REMOTE="${PROTON_PACKAGE_REMOTE:-${PROTON_OFFLINE_DIR}/$(basename "${PROTON_PACKAGE_LOCAL}")}"
PROTON_EXTRACTED_DIR="${PROTON_EXTRACTED_DIR:-${PROTON_OFFLINE_DIR}}"
PROTON_SERVICE_IMAGES_DIR="${PROTON_SERVICE_IMAGES_DIR:-${PROTON_EXTRACTED_DIR}/service-package/images}"
PLATFORM_IMAGES_PACKAGE_LOCAL="${PLATFORM_IMAGES_PACKAGE_LOCAL:-${PACKAGE_ROOT}/offline-packages/platform-images/platform-images-0.6.0-amd64.tar}"
PLATFORM_IMAGES_DIR="${PLATFORM_IMAGES_DIR:-${OFFLINE_ROOT}/platform-images}"
PLATFORM_IMAGES_PACKAGE_REMOTE="${PLATFORM_IMAGES_PACKAGE_REMOTE:-${PLATFORM_IMAGES_DIR}/$(basename "${PLATFORM_IMAGES_PACKAGE_LOCAL}")}"

KWEAVER_PACKAGE_LOCAL="${KWEAVER_PACKAGE_LOCAL:-${PACKAGE_ROOT}/offline-packages/kweaver-core/${KWEAVER_VERSION}/kweaver-core-${KWEAVER_VERSION}-offline-package.tar}"
KWEAVER_DEPLOY_BUNDLE_LOCAL="${KWEAVER_DEPLOY_BUNDLE_LOCAL:-${PACKAGE_ROOT}/offline-packages/kweaver-core/${KWEAVER_VERSION}/deploy}"
KWEAVER_CORE_MANIFEST_LOCAL="${KWEAVER_CORE_MANIFEST_LOCAL:-${PACKAGE_ROOT}/offline-packages/kweaver-core/${KWEAVER_VERSION}/manifests/kweaver-core.yaml}"
KWEAVER_ISF_MANIFEST_LOCAL="${KWEAVER_ISF_MANIFEST_LOCAL:-${PACKAGE_ROOT}/offline-packages/kweaver-core/${KWEAVER_VERSION}/manifests/isf.yaml}"
KWEAVER_OFFLINE_DIR="${KWEAVER_OFFLINE_DIR:-${OFFLINE_ROOT}/kweaver-core/${KWEAVER_VERSION}}"
KWEAVER_PACKAGE_REMOTE="${KWEAVER_PACKAGE_REMOTE:-${KWEAVER_OFFLINE_DIR}/$(basename "${KWEAVER_PACKAGE_LOCAL}")}"
KWEAVER_MANIFESTS_REMOTE_DIR="${KWEAVER_MANIFESTS_REMOTE_DIR:-${KWEAVER_OFFLINE_DIR}/manifests}"

V9_INFRA_CHART_LOCAL="${V9_INFRA_CHART_LOCAL:-${PACKAGE_ROOT}/helm-chart/anybackup-agent/charts/v9-infra}"
V9_INFRA_CHART_REMOTE="${V9_INFRA_CHART_REMOTE:-${REMOTE_CHART_ROOT}/v9-infra}"

DEFAULT_VERIFY_TIMEOUT="${DEFAULT_VERIFY_TIMEOUT:-600s}"
HELM_WAIT_TIMEOUT="${HELM_WAIT_TIMEOUT:-30m}"

if command -v hostname >/dev/null 2>&1; then
  DETECTED_PUBLIC_IP="$(hostname -I 2>/dev/null | awk '{print $1}' || true)"
fi
if [[ -z "${DETECTED_PUBLIC_IP:-}" ]] && command -v ip >/dev/null 2>&1; then
  DETECTED_PUBLIC_IP="$(
    ip route get 1 2>/dev/null \
      | awk 'NR == 1 { for (i = 1; i <= NF; i++) if ($i == "src") { print $(i + 1); exit } }' \
      || true
  )"
fi
PUBLIC_IP="${PUBLIC_IP:-${DETECTED_PUBLIC_IP:-127.0.0.1}}"

log_info() {
  printf '[INFO] %s\n' "$*"
}

log_warn() {
  printf '[WARN] %s\n' "$*" >&2
}

log_error() {
  printf '[ERROR] %s\n' "$*" >&2
}

fail() {
  log_error "$*"
  exit 1
}

require_root() {
  if [[ "$(id -u)" -ne 0 ]]; then
    fail "This script must run as root."
  fi
}

require_cmd() {
  local cmd
  for cmd in "$@"; do
    if ! command -v "${cmd}" >/dev/null 2>&1; then
      fail "Required command not found: ${cmd}"
    fi
  done
}

require_file() {
  local path="$1"
  [[ -e "${path}" ]] || fail "Required path not found: ${path}"
}

ensure_dir() {
  mkdir -p "$1"
}

bool_true() {
  case "${1:-false}" in
    true|TRUE|1|yes|YES|on|ON)
      return 0
      ;;
    *)
      return 1
      ;;
  esac
}

base64_nowrap() {
  if base64 --help 2>/dev/null | grep -q -- '-w'; then
    printf '%s' "$1" | base64 -w 0
  else
    printf '%s' "$1" | base64 | tr -d '\n'
  fi
}

copy_dir_clean() {
  local src="$1"
  local dest="$2"
  rm -rf "${dest}"
  mkdir -p "$(dirname "${dest}")"
  cp -R "${src}" "${dest}"
}

copy_file_into_dir() {
  local src="$1"
  local dest_dir="$2"
  mkdir -p "${dest_dir}"
  cp -f "${src}" "${dest_dir}/"
}

render_v9_infra_values() {
  cat > "${GENERATED_ROOT}/v9-infra-values.yaml" <<EOF
global:
  imagePullPolicy: IfNotPresent
  defaultStorageClass: "local-path"

namespace: ${V9_INFRA_NAMESPACE}
postgres:
  image: docker.m.daocloud.io/library/postgres:17
  storage: 50Gi
  storageClass: "local-path"
  secretName: ${V9_POSTGRES_SECRET_NAME}
  serviceName: v9-postgres
rabbitmq:
  image: docker.m.daocloud.io/library/rabbitmq:3-management
  storage: 10Gi
  storageClass: "local-path"
  secretName: ${V9_RABBITMQ_SECRET_NAME}
  serviceName: v9-rabbitmq
redis:
  image: docker.m.daocloud.io/library/redis:7
  storage: 10Gi
  storageClass: "local-path"
  secretName: ${V9_REDIS_SECRET_NAME}
  serviceName: v9-redis
opensearch:
  image: docker.m.daocloud.io/opensearchproject/opensearch:3.6.0
  initImage: docker.m.daocloud.io/library/busybox:latest
  storage: 50Gi
  storageClass: "local-path"
  serviceName: v9-opensearch
  protocol: "http"
EOF
}

render_v9_infra_secrets() {
  cat > "${GENERATED_ROOT}/v9-postgres-secret.yaml" <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: ${V9_POSTGRES_SECRET_NAME}
  namespace: ${V9_INFRA_NAMESPACE}
type: Opaque
data:
  POSTGRES_USER: $(base64_nowrap "${V9_POSTGRES_USERNAME}")
  POSTGRES_PASSWORD: $(base64_nowrap "${V9_POSTGRES_PASSWORD}")
  POSTGRES_DB: $(base64_nowrap "${V9_POSTGRES_DATABASE}")
EOF

  cat > "${GENERATED_ROOT}/v9-rabbitmq-secret.yaml" <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: ${V9_RABBITMQ_SECRET_NAME}
  namespace: ${V9_INFRA_NAMESPACE}
type: Opaque
data:
  RABBITMQ_ERLANG_COOKIE: $(base64_nowrap "${V9_RABBITMQ_ERLANG_COOKIE}")
  RABBITMQ_DEFAULT_USER: $(base64_nowrap "${V9_RABBITMQ_USERNAME}")
  RABBITMQ_DEFAULT_PASS: $(base64_nowrap "${V9_RABBITMQ_PASSWORD}")
EOF

  cat > "${GENERATED_ROOT}/v9-redis-secret.yaml" <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: ${V9_REDIS_SECRET_NAME}
  namespace: ${V9_INFRA_NAMESPACE}
type: Opaque
data:
  REDIS_PASSWORD: $(base64_nowrap "${V9_REDIS_PASSWORD}")
EOF
}

ensure_package_assets_present() {
  require_file "${V9_INFRA_CHART_LOCAL}"
}

prepare_runtime_dirs() {
  ensure_dir "${GENERATED_ROOT}"
  ensure_dir "${REMOTE_CHART_ROOT}"
  ensure_dir "${REMOTE_SOURCE_ROOT}"
  ensure_dir "${KWEAVER_CONFIG_DIR}"
}

stage_runtime_assets() {
  copy_dir_clean "${V9_INFRA_CHART_LOCAL}" "${V9_INFRA_CHART_REMOTE}"
}

deploy_v9_infra() {
  log_info "Deploying v9_infra."
  render_v9_infra_values
  render_v9_infra_secrets
  kubectl create namespace "${V9_INFRA_NAMESPACE}" --dry-run=client -o yaml | kubectl apply -f -
  kubectl apply -f "${GENERATED_ROOT}/v9-postgres-secret.yaml"
  kubectl apply -f "${GENERATED_ROOT}/v9-rabbitmq-secret.yaml"
  kubectl apply -f "${GENERATED_ROOT}/v9-redis-secret.yaml"
  helm upgrade --install "${V9_INFRA_RELEASE}" "${V9_INFRA_CHART_REMOTE}" \
    --namespace "${V9_INFRA_NAMESPACE}" \
    --create-namespace \
    --values "${GENERATED_ROOT}/v9-infra-values.yaml" \
    --wait \
    --timeout "${HELM_WAIT_TIMEOUT}"
  kubectl rollout status statefulset/v9-infra-postgres -n "${V9_INFRA_NAMESPACE}" --timeout="${DEFAULT_VERIFY_TIMEOUT}"
  kubectl rollout status statefulset/v9-infra-rabbitmq -n "${V9_INFRA_NAMESPACE}" --timeout="${DEFAULT_VERIFY_TIMEOUT}"
  kubectl rollout status statefulset/v9-infra-redis -n "${V9_INFRA_NAMESPACE}" --timeout="${DEFAULT_VERIFY_TIMEOUT}"
  kubectl rollout status statefulset/v9-infra-opensearch -n "${V9_INFRA_NAMESPACE}" --timeout="${DEFAULT_VERIFY_TIMEOUT}"
}

install_proton_base() {
  local rpm_dir="${PROTON_EXTRACTED_DIR}/repos/Packages"
  local rpm_paths=(
    "${rpm_dir}/ecms-1.1.8-120.1.el7.x86_64.rpm"
    "${rpm_dir}/haproxy-2.5.6.x86_64.rpm"
    "${rpm_dir}/proton-cr-1.2.5-87.el7.x86_64.rpm"
    "${rpm_dir}/proton-cr-chartmuseum-0.15.0.x86_64.rpm"
    "${rpm_dir}/proton-cr-registry-2.7.1.x86_64.rpm"
  )
  local rpm_name
  local missing_rpms=()

  log_info "Installing Proton base components."
  if [[ ! -f "${PROTON_EXTRACTED_DIR}/install.sh" ]]; then
    prepare_proton_extract_dir
    tar -xf "${PROTON_PACKAGE_REMOTE}" -C "${PROTON_EXTRACTED_DIR}"
  fi

  for rpm_name in ecms haproxy proton-cr proton-cr-chartmuseum proton-cr-registry; do
    if ! rpm -q "${rpm_name}" >/dev/null 2>&1; then
      missing_rpms+=("${rpm_name}")
    fi
  done

  if [[ "${#missing_rpms[@]}" -gt 0 ]]; then
    yum install -y --nogpgcheck "${rpm_paths[@]}"
  else
    log_info "Proton base RPMs already installed."
  fi

  systemctl enable --now ecms proton-cr
  ensure_containerd_registry_override
  verify_proton_ports
  import_offline_platform_images
}

ensure_containerd_registry_override() {
  local config_changed="false"
  local host_dir="/etc/containerd/certs.d/${PUBLIC_IP}:${PROTON_REGISTRY_PORT}"
  local hosts_toml="${host_dir}/hosts.toml"

  mkdir -p "/etc/containerd/certs.d"
  if grep -q 'config_path = ""' /etc/containerd/config.toml 2>/dev/null; then
    python3 - <<'PY'
from pathlib import Path
path = Path("/etc/containerd/config.toml")
text = path.read_text()
text = text.replace('config_path = ""', "config_path = '/etc/containerd/certs.d'")
path.write_text(text)
PY
    config_changed="true"
  fi

  mkdir -p "${host_dir}"
  cat > "${hosts_toml}" <<EOF
server = "http://${PUBLIC_IP}:${PROTON_REGISTRY_PORT}"

[host."http://${PUBLIC_IP}:${PROTON_REGISTRY_PORT}"]
  capabilities = ["pull", "resolve"]
  skip_verify = true
EOF

  if [[ "${config_changed}" == "true" ]]; then
    systemctl restart containerd
    return 0
  fi

  systemctl restart containerd
}

proton_ports_listening() {
  ss -ltn | awk -v registry_port="${PROTON_REGISTRY_PORT}" \
    -v chartmuseum_port="${PROTON_CHARTMUSEUM_PORT}" '
    $4 ~ ":" registry_port "$" { registry = 1 }
    $4 ~ ":" chartmuseum_port "$" { chartmuseum = 1 }
    END {
      if (registry && chartmuseum) {
        exit 0
      }
      exit 1
    }
  '
}

verify_proton_ports() {
  local attempt

  for ((attempt = 1; attempt <= PROTON_PORT_WAIT_ATTEMPTS; attempt++)); do
    if proton_ports_listening; then
      log_info "Proton registry and ChartMuseum listeners are ready."
      return 0
    fi

    if [[ "${attempt}" -lt "${PROTON_PORT_WAIT_ATTEMPTS}" ]]; then
      log_info "Waiting for Proton registry and ChartMuseum listeners (${attempt}/${PROTON_PORT_WAIT_ATTEMPTS})."
      sleep "${PROTON_PORT_WAIT_SECONDS}"
    fi
  done

  fail "Proton registry or ChartMuseum listener is missing."
}

push_oci_images_to_registry() {
  local package_path="$1"
  local workdir="${REMOTE_DEPLOY_ROOT}/tmp-push-images"

  require_file "${package_path}"
  mkdir -p "${workdir}"
  "${PROTON_CLI_REMOTE}" push-images \
    --package "${package_path}" \
    --registry "${PUBLIC_IP}:${PROTON_REGISTRY_PORT}" \
    --workdir "${workdir}"
}

import_offline_platform_images() {
  log_info "Importing Proton service-package images into the local registry."
  push_oci_images_to_registry "${PROTON_SERVICE_IMAGES_DIR}"

  log_info "Importing KWeaver platform dependency images into the local registry."
  push_oci_images_to_registry "${PLATFORM_IMAGES_PACKAGE_REMOTE}"
}

prepare_proton_extract_dir() {
  local package_name
  package_name="$(basename "${PROTON_PACKAGE_REMOTE}")"

  if [[ "${PROTON_PACKAGE_REMOTE}" == "${PROTON_EXTRACTED_DIR}/"* ]]; then
    mkdir -p "${PROTON_EXTRACTED_DIR}"
    find "${PROTON_EXTRACTED_DIR}" -mindepth 1 -maxdepth 1 ! -name "${package_name}" -exec rm -rf {} +
  else
    rm -rf "${PROTON_EXTRACTED_DIR}"
    mkdir -p "${PROTON_EXTRACTED_DIR}"
  fi
}

run_kweaver_platform_install() {
  local component="$1"
  IMAGE_REGISTRY="${PUBLIC_IP}:${PROTON_REGISTRY_PORT}" \
  MARIADB_IMAGE="${PUBLIC_IP}:${PROTON_REGISTRY_PORT}/mariadb:11.4.7" \
  REDIS_IMAGE_REGISTRY="${PUBLIC_IP}:${PROTON_REGISTRY_PORT}" \
  ZOOKEEPER_IMAGE_REGISTRY="${PUBLIC_IP}:${PROTON_REGISTRY_PORT}" \
  KAFKA_IMAGE="${PUBLIC_IP}:${PROTON_REGISTRY_PORT}/bitnami/kafka:3.9.0-debian-12-r10" \
  OPENSEARCH_IMAGE="${PUBLIC_IP}:${PROTON_REGISTRY_PORT}/opensearchproject/opensearch:2.19.4" \
  OPENSEARCH_INIT_IMAGE="${PUBLIC_IP}:${PROTON_REGISTRY_PORT}/busybox:1.36.1" \
  INGRESS_NGINX_CONTROLLER_IMAGE="${PUBLIC_IP}:${PROTON_REGISTRY_PORT}/ingress-nginx-controller:v1.15.0" \
  INGRESS_NGINX_WEBHOOK_CERTGEN_IMAGE="${PUBLIC_IP}:${PROTON_REGISTRY_PORT}/ingress-nginx-kube-webhook-certgen:v1.6.8" \
  CONF_DIR="${KWEAVER_CONFIG_DIR}" \
  CONFIG_YAML_PATH="${KWEAVER_CONFIG_PATH}" \
  CONFIRM_ACCESS_ADDRESS="false" \
  KWEAVER_ACCESS_ADDRESS="https://${PUBLIC_IP}:${INGRESS_HTTPS_NODE_PORT}" \
  INGRESS_NGINX_HOSTNETWORK="false" \
  INGRESS_NGINX_CLASS="${INGRESS_CLASS_NAME}" \
  INGRESS_NGINX_HTTP_PORT="${INGRESS_HTTP_NODE_PORT}" \
  INGRESS_NGINX_HTTPS_PORT="${INGRESS_HTTPS_NODE_PORT}" \
  AUTO_INSTALL_INGRESS_NGINX="true" \
    "${KWEAVER_REMOTE_DEPLOY_DIR}/deploy.sh" "${component}" install
}

generate_kweaver_config() {
  IMAGE_REGISTRY="${PUBLIC_IP}:${PROTON_REGISTRY_PORT}" \
  CONF_DIR="${KWEAVER_CONFIG_DIR}" \
  CONFIG_YAML_PATH="${KWEAVER_CONFIG_PATH}" \
  CONFIRM_ACCESS_ADDRESS="false" \
  KWEAVER_ACCESS_ADDRESS="https://${PUBLIC_IP}:${INGRESS_HTTPS_NODE_PORT}" \
  INGRESS_NGINX_HOSTNETWORK="false" \
  INGRESS_NGINX_CLASS="${INGRESS_CLASS_NAME}" \
  INGRESS_NGINX_HTTP_PORT="${INGRESS_HTTP_NODE_PORT}" \
  INGRESS_NGINX_HTTPS_PORT="${INGRESS_HTTPS_NODE_PORT}" \
    "${KWEAVER_REMOTE_DEPLOY_DIR}/deploy.sh" config generate
}

deploy_kweaver_platform() {
  local component
  log_info "Deploying KWeaver platform dependencies."
  chmod +x "${KWEAVER_REMOTE_DEPLOY_DIR}/deploy.sh"
  for component in ingress-nginx mariadb redis zookeeper kafka opensearch; do
    run_kweaver_platform_install "${component}"
  done
  generate_kweaver_config
}

checkout_kweaver_source() {
  log_info "Cloning or updating KWeaver Core source from ${KWEAVER_REPO_URL}."
  mkdir -p "$(dirname "${KWEAVER_REPO_DIR}")"

  if [[ -d "${KWEAVER_REPO_DIR}/.git" ]]; then
    git -C "${KWEAVER_REPO_DIR}" remote set-url origin "${KWEAVER_REPO_URL}"
    git -C "${KWEAVER_REPO_DIR}" fetch --prune --tags origin
  else
    rm -rf "${KWEAVER_REPO_DIR}"
    git clone "${KWEAVER_REPO_URL}" "${KWEAVER_REPO_DIR}"
    git -C "${KWEAVER_REPO_DIR}" fetch --prune --tags origin
  fi

  checkout_kweaver_ref
}

checkout_kweaver_ref() {
  local ref
  local candidates=(
    "${KWEAVER_GIT_REF}"
    "release/${KWEAVER_VERSION}"
    "v${KWEAVER_VERSION}"
    "${KWEAVER_VERSION}"
  )

  for ref in "${candidates[@]}"; do
    [[ -n "${ref}" ]] || continue
    if git -C "${KWEAVER_REPO_DIR}" rev-parse --verify --quiet "origin/${ref}^{commit}" >/dev/null; then
      git -C "${KWEAVER_REPO_DIR}" checkout --force -B "deploy-${KWEAVER_VERSION}" "origin/${ref}"
      log_info "Checked out KWeaver source ref origin/${ref}."
      return 0
    fi
    if git -C "${KWEAVER_REPO_DIR}" rev-parse --verify --quiet "${ref}^{commit}" >/dev/null; then
      git -C "${KWEAVER_REPO_DIR}" checkout --force "${ref}"
      log_info "Checked out KWeaver source ref ${ref}."
      return 0
    fi
  done

  fail "Unable to find KWeaver source ref. Tried: ${candidates[*]}"
}

kweaver_source_commit() {
  git -C "${KWEAVER_REPO_DIR}" rev-parse HEAD 2>/dev/null || true
}

run_kweaver_online_install() {
  log_info "Installing KWeaver Core online through official atomic command."
  require_file "${KWEAVER_REMOTE_DEPLOY_DIR}/deploy.sh"
  chmod +x "${KWEAVER_REMOTE_DEPLOY_DIR}/deploy.sh"

  (
    cd "${KWEAVER_REMOTE_DEPLOY_DIR}"
    CONF_DIR="${KWEAVER_CONFIG_DIR}" \
    CONFIG_YAML_PATH="${KWEAVER_CONFIG_PATH}" \
    CONFIRM_ACCESS_ADDRESS="false" \
    KWEAVER_ACCESS_ADDRESS="https://${PUBLIC_IP}:${INGRESS_HTTPS_NODE_PORT}" \
    INGRESS_NGINX_HOSTNETWORK="false" \
    INGRESS_NGINX_CLASS="${INGRESS_CLASS_NAME}" \
    INGRESS_NGINX_HTTP_PORT="${INGRESS_HTTP_NODE_PORT}" \
    INGRESS_NGINX_HTTPS_PORT="${INGRESS_HTTPS_NODE_PORT}" \
    AUTO_INSTALL_INGRESS_NGINX="true" \
      ./deploy.sh kweaver-core install
  )
}

persist_kweaver_state() {
  local source_commit
  source_commit="$(kweaver_source_commit)"

  kubectl create namespace "${KWEAVER_NAMESPACE}" --dry-run=client -o yaml | kubectl apply -f - >/dev/null
  cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: ${KWEAVER_STATE_CONFIGMAP}
  namespace: ${KWEAVER_NAMESPACE}
data:
  install_mode: online
  target_version: "${KWEAVER_VERSION}"
  repo_url: "${KWEAVER_REPO_URL}"
  git_ref: "${KWEAVER_GIT_REF}"
  source_commit: "${source_commit}"
  install_command: "./deploy.sh kweaver-core install"
  success_source: official_command_exit_code
  last_success_time: "$(date -Iseconds)"
EOF
}

install_kweaver_core() {
  checkout_kweaver_source
  run_kweaver_online_install
  persist_kweaver_state
}

wait_running_pods_ready() {
  local namespace="$1"
  local timeout="${2:-${DEFAULT_VERIFY_TIMEOUT}}"
  local pods
  pods="$(kubectl get pods -n "${namespace}" --field-selector=status.phase=Running -o jsonpath='{range .items[*]}pod/{.metadata.name} {end}' 2>/dev/null || true)"
  if [[ -n "${pods}" ]]; then
    kubectl wait -n "${namespace}" --for=condition=Ready ${pods} --timeout="${timeout}"
  fi
}

verify_kweaver_core_only() {
  log_info "Verifying v9_infra, ingress-nginx, and KWeaver Core."
  kubectl rollout status statefulset/v9-infra-postgres -n "${V9_INFRA_NAMESPACE}" --timeout="${DEFAULT_VERIFY_TIMEOUT}"
  kubectl rollout status statefulset/v9-infra-rabbitmq -n "${V9_INFRA_NAMESPACE}" --timeout="${DEFAULT_VERIFY_TIMEOUT}"
  kubectl rollout status statefulset/v9-infra-redis -n "${V9_INFRA_NAMESPACE}" --timeout="${DEFAULT_VERIFY_TIMEOUT}"
  kubectl rollout status statefulset/v9-infra-opensearch -n "${V9_INFRA_NAMESPACE}" --timeout="${DEFAULT_VERIFY_TIMEOUT}"
  kubectl rollout status deployment/ingress-nginx-controller -n ingress-nginx --timeout="${DEFAULT_VERIFY_TIMEOUT}"
  kubectl rollout status statefulset/mariadb-proton-mariadb -n "${KWEAVER_RESOURCE_NAMESPACE}" --timeout="${DEFAULT_VERIFY_TIMEOUT}"
  kubectl rollout status statefulset/redis-proton-redis -n "${KWEAVER_RESOURCE_NAMESPACE}" --timeout="${DEFAULT_VERIFY_TIMEOUT}"
  kubectl rollout status statefulset/zookeeper -n "${KWEAVER_RESOURCE_NAMESPACE}" --timeout="${DEFAULT_VERIFY_TIMEOUT}"
  kubectl rollout status statefulset/kafka-broker -n "${KWEAVER_RESOURCE_NAMESPACE}" --timeout="${DEFAULT_VERIFY_TIMEOUT}"
  kubectl rollout status statefulset/kafka-controller -n "${KWEAVER_RESOURCE_NAMESPACE}" --timeout="${DEFAULT_VERIFY_TIMEOUT}"
  kubectl rollout status statefulset/opensearch-cluster-master -n "${KWEAVER_RESOURCE_NAMESPACE}" --timeout="${DEFAULT_VERIFY_TIMEOUT}"
  wait_running_pods_ready "${KWEAVER_NAMESPACE}" "${DEFAULT_VERIFY_TIMEOUT}"
  kubectl get ingressclass "${INGRESS_CLASS_NAME}" >/dev/null
  kubectl get ingress ingress-informationsecurityfabric -n "${KWEAVER_NAMESPACE}" >/dev/null
  kubectl get ingress rule-isfweb -n "${KWEAVER_NAMESPACE}" >/dev/null

  printf '\n[SUMMARY] v9-system\n'
  kubectl get pods,svc -n "${V9_INFRA_NAMESPACE}"
  printf '\n[SUMMARY] ingress-nginx\n'
  kubectl get all -n ingress-nginx
  printf '\n[SUMMARY] resource\n'
  kubectl get pods,svc -n "${KWEAVER_RESOURCE_NAMESPACE}"
  printf '\n[SUMMARY] kweaver\n'
  kubectl get pods,svc,ingress -n "${KWEAVER_NAMESPACE}"
  printf '\n[SUMMARY] helm releases\n'
  helm list -n "${V9_INFRA_NAMESPACE}" || true
  helm list -n "${KWEAVER_RESOURCE_NAMESPACE}" || true
  helm list -n "${KWEAVER_NAMESPACE}" || true
  printf '\n[SUMMARY] access\n'
  printf 'https://%s:%s\n' "${PUBLIC_IP}" "${INGRESS_HTTPS_NODE_PORT}"
}

delete_namespace_if_exists() {
  local namespace="$1"
  if kubectl get namespace "${namespace}" >/dev/null 2>&1; then
    log_info "Deleting namespace ${namespace}."
    kubectl delete namespace "${namespace}" --ignore-not-found --wait=false >/dev/null || true
  fi
}

wait_namespace_gone() {
  local namespace="$1"
  local attempt
  for attempt in $(seq 1 60); do
    if ! kubectl get namespace "${namespace}" >/dev/null 2>&1; then
      return 0
    fi
    sleep 5
  done

  log_warn "Namespace ${namespace} is still terminating. Clearing finalizers."
  kubectl patch namespace "${namespace}" --type=json -p='[{"op":"replace","path":"/spec/finalizers","value":[]}]' >/dev/null 2>&1 || true

  for attempt in $(seq 1 24); do
    if ! kubectl get namespace "${namespace}" >/dev/null 2>&1; then
      return 0
    fi
    sleep 5
  done

  fail "Namespace ${namespace} could not be removed cleanly."
}

delete_bound_pvs_for_validation_namespaces() {
  local target_namespaces
  local pv_json
  target_namespaces="middleware,${V9_INFRA_NAMESPACE},ingress-nginx,${KWEAVER_RESOURCE_NAMESPACE},${KWEAVER_NAMESPACE},v9,anybackup-ai"
  pv_json="$(kubectl get pv -o json 2>/dev/null || printf '{"items":[]}\n')"
  TARGET_NAMESPACES="${target_namespaces}" PV_JSON="${pv_json}" python3 - <<'PY' | while read -r pv_name; do
import json
import os

namespaces = set(filter(None, os.environ.get("TARGET_NAMESPACES", "").split(",")))
data = json.loads(os.environ.get("PV_JSON", '{"items":[]}'))
for item in data.get("items", []):
    claim_ref = (item.get("spec") or {}).get("claimRef") or {}
    if claim_ref.get("namespace") in namespaces:
        print(item["metadata"]["name"])
PY
      [[ -n "${pv_name}" ]] || continue
      kubectl delete pv "${pv_name}" --ignore-not-found >/dev/null || true
    done
}

delete_validation_pv() {
  local pv_name="$1"
  local attempt

  kubectl patch pv "${pv_name}" --type=merge -p '{"spec":{"claimRef":null}}' >/dev/null 2>&1 || true
  kubectl delete pv "${pv_name}" --ignore-not-found >/dev/null || true

  for attempt in $(seq 1 60); do
    if ! kubectl get pv "${pv_name}" >/dev/null 2>&1; then
      return 0
    fi
    sleep 2
  done

  fail "PersistentVolume ${pv_name} could not be removed cleanly."
}

cleanup_kweaver_core_validation() {
  local namespace
  local ingress_class
  local service_name
  local rpm_name

  for namespace in middleware "${V9_INFRA_NAMESPACE}" ingress-nginx "${KWEAVER_RESOURCE_NAMESPACE}" "${KWEAVER_NAMESPACE}" v9 anybackup-ai; do
    delete_namespace_if_exists "${namespace}"
  done

  for namespace in middleware "${V9_INFRA_NAMESPACE}" ingress-nginx "${KWEAVER_RESOURCE_NAMESPACE}" "${KWEAVER_NAMESPACE}" v9 anybackup-ai; do
    if kubectl get namespace "${namespace}" >/dev/null 2>&1; then
      wait_namespace_gone "${namespace}"
    fi
  done

  for ingress_class in "${INGRESS_CLASS_NAME}" traefik; do
    kubectl delete ingressclass "${ingress_class}" --ignore-not-found >/dev/null || true
  done

  delete_bound_pvs_for_validation_namespaces
  delete_validation_pv minio-pv

  for service_name in ecms proton-cr; do
    systemctl disable --now "${service_name}" >/dev/null 2>&1 || true
  done

  for rpm_name in ecms haproxy proton-cr proton-cr-chartmuseum proton-cr-registry; do
    if rpm -q "${rpm_name}" >/dev/null 2>&1; then
      yum remove -y "${rpm_name}" >/dev/null || true
    fi
  done

  rm -f "${PROTON_CLI_REMOTE}"
  rm -rf "/etc/containerd/certs.d/${PUBLIC_IP}:${PROTON_REGISTRY_PORT}"
  rm -rf "${REMOTE_DEPLOY_ROOT}" /opt/v9-sources "${KWEAVER_CONFIG_DIR}"
  if command -v crictl >/dev/null 2>&1; then
    crictl rmi --prune >/dev/null 2>&1 || true
  fi
  if command -v ctr >/dev/null 2>&1; then
    ctr -n k8s.io images prune >/dev/null 2>&1 || true
  fi
  systemctl restart containerd >/dev/null 2>&1 || true
}
