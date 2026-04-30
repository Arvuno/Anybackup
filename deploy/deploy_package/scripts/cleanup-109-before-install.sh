#!/usr/bin/env bash
set -euo pipefail

DRY_RUN=true
YES=false
WAIT_TIMEOUT_SECONDS=600

APP_NAMESPACES=(
  anybackup-ai
  kweaver
  resource
  v9-system
  ingress-nginx
  middleware
)

INGRESS_CLASSES=(
  class-443
  traefik
  nginx
)

REMOTE_PATHS=(
  /opt/v9-alpha-deploy
  /opt/v9-sources
  /root/.kweaver
  /root/.kweaver-admin
)

usage() {
  cat <<'EOF'
Usage:
  bash scripts/cleanup-109-before-install.sh --dry-run
  bash scripts/cleanup-109-before-install.sh --yes

Purpose:
  Clean the 109 validation host before rerunning the integrated deployment.

This script removes:
  - Helm releases and Kubernetes resources in app namespaces:
    anybackup-ai, kweaver, resource, v9-system, ingress-nginx, middleware
  - IngressClass residue: class-443, traefik, nginx
  - PVs whose claimRef points to the removed app namespaces
  - V9/KWeaver working directories under /opt
  - root KWeaver CLI cached login state

This script intentionally keeps:
  - Kubernetes itself: kube-system, kube-flannel, local-path-storage, etc.
  - Container runtime and imported images
  - Foundation package: /backupsoft/Linux_el7_x64-latest.tar.gz
  - Foundation installation: /backupsoft/AnyBackupServer

Options:
  --dry-run                 Print actions without deleting anything. Default.
  --yes                     Execute cleanup.
  --wait-timeout SECONDS    Namespace/PV wait timeout. Default: 600.
  -h, --help                Show this help.
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --dry-run)
      DRY_RUN=true
      ;;
    --yes)
      DRY_RUN=false
      YES=true
      ;;
    --wait-timeout)
      shift
      if [[ $# -eq 0 ]]; then
        echo "ERROR: --wait-timeout requires a value" >&2
        exit 1
      fi
      WAIT_TIMEOUT_SECONDS="$1"
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "ERROR: unknown argument: $1" >&2
      usage
      exit 1
      ;;
  esac
  shift
done

if [[ "${DRY_RUN}" == "false" && "${YES}" != "true" ]]; then
  echo "ERROR: destructive cleanup requires --yes" >&2
  exit 1
fi

log() {
  printf '[cleanup] %s\n' "$*"
}

run() {
  if [[ "${DRY_RUN}" == "true" ]]; then
    printf '[dry-run] %q' "$1"
    shift || true
    for arg in "$@"; do
      printf ' %q' "$arg"
    done
    printf '\n'
  else
    "$@"
  fi
}

run_shell() {
  local script="$1"
  if [[ "${DRY_RUN}" == "true" ]]; then
    printf '[dry-run] bash -lc %q\n' "${script}"
  else
    bash -lc "${script}"
  fi
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "ERROR: required command not found: $1" >&2
    exit 1
  }
}

wait_for_namespace_absent() {
  local ns="$1"
  local elapsed=0
  while kubectl get namespace "${ns}" >/dev/null 2>&1; do
    if (( elapsed >= WAIT_TIMEOUT_SECONDS )); then
      echo "ERROR: namespace still exists after waiting: ${ns}" >&2
      return 1
    fi
    sleep 5
    elapsed=$((elapsed + 5))
  done
}

wait_for_pv_absent() {
  local pv="$1"
  local elapsed=0
  while kubectl get pv "${pv}" >/dev/null 2>&1; do
    if (( elapsed >= WAIT_TIMEOUT_SECONDS )); then
      echo "ERROR: persistent volume still exists after waiting: ${pv}" >&2
      return 1
    fi
    sleep 5
    elapsed=$((elapsed + 5))
  done
}

delete_helm_releases() {
  local ns="$1"
  kubectl get namespace "${ns}" >/dev/null 2>&1 || return 0

  local releases
  releases="$(helm list -n "${ns}" -q 2>/dev/null || true)"
  if [[ -z "${releases}" ]]; then
    return 0
  fi

  while IFS= read -r release; do
    [[ -n "${release}" ]] || continue
    log "Uninstall Helm release ${release} in ${ns}"
    run helm uninstall -n "${ns}" "${release}" --wait --timeout 10m
  done <<< "${releases}"
}

force_clear_namespace_finalizers() {
  local ns="$1"
  kubectl get namespace "${ns}" >/dev/null 2>&1 || return 0

  log "Clearing finalizers for resources in namespace ${ns} if needed"
  run_shell "kubectl api-resources --verbs=list --namespaced -o name \
    | while read -r resource; do \
        kubectl get \"\${resource}\" -n '${ns}' -o name 2>/dev/null \
          | while read -r item; do \
              kubectl patch \"\${item}\" -n '${ns}' --type=merge -p '{\"metadata\":{\"finalizers\":[]}}' >/dev/null 2>&1 || true; \
            done; \
      done"
}

collect_pvs_for_namespaces() {
  local ns_regex
  ns_regex="$(IFS='|'; echo "${APP_NAMESPACES[*]}")"
  kubectl get pv -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.claimRef.namespace}{"\n"}{end}' 2>/dev/null \
    | awk -F '\t' -v re="^(${ns_regex})$" '$2 ~ re {print $1}' \
    | sort -u
}

main() {
  require_cmd kubectl
  require_cmd helm
  require_cmd awk

  if [[ "${DRY_RUN}" == "true" ]]; then
    log "DRY RUN only. Re-run with --yes to execute cleanup."
  else
    log "Executing destructive cleanup."
  fi

  log "Target namespaces: ${APP_NAMESPACES[*]}"
  log "Keeping Kubernetes system namespaces and Foundation installation."

  local pvs_before=()
  mapfile -t pvs_before < <(collect_pvs_for_namespaces || true)

  for ns in "${APP_NAMESPACES[@]}"; do
    delete_helm_releases "${ns}"
  done

  for ns in "${APP_NAMESPACES[@]}"; do
    if kubectl get namespace "${ns}" >/dev/null 2>&1; then
      log "Delete namespace ${ns}"
      run kubectl delete namespace "${ns}" --ignore-not-found=true
    fi
  done

  if [[ "${DRY_RUN}" == "false" ]]; then
    for ns in "${APP_NAMESPACES[@]}"; do
      if ! wait_for_namespace_absent "${ns}"; then
        force_clear_namespace_finalizers "${ns}"
        wait_for_namespace_absent "${ns}"
      fi
    done
  fi

  for ingress_class in "${INGRESS_CLASSES[@]}"; do
    log "Delete ingress class ${ingress_class}"
    run kubectl delete ingressclass "${ingress_class}" --ignore-not-found=true
  done

  local pvs_after=()
  if [[ "${DRY_RUN}" == "true" ]]; then
    pvs_after=("${pvs_before[@]}")
  else
    mapfile -t pvs_after < <(collect_pvs_for_namespaces || true)
  fi

  local pvs=("${pvs_before[@]}" "${pvs_after[@]}")
  if (( ${#pvs[@]} > 0 )); then
    mapfile -t pvs < <(printf '%s\n' "${pvs[@]}" | awk 'NF' | sort -u)
  fi

  for pv in "${pvs[@]}"; do
    [[ -n "${pv}" ]] || continue
    if [[ "${DRY_RUN}" == "false" ]] && ! kubectl get pv "${pv}" >/dev/null 2>&1; then
      log "Skip PV ${pv}; already absent"
      continue
    fi
    log "Clear claimRef and delete PV ${pv}"
    if [[ "${DRY_RUN}" == "true" ]]; then
      run kubectl patch pv "${pv}" --type=merge --patch '{"spec":{"claimRef":null}}'
    else
      kubectl patch pv "${pv}" --type=merge --patch '{"spec":{"claimRef":null}}' >/dev/null 2>&1 || true
    fi
    run kubectl delete pv "${pv}" --ignore-not-found=true
  done

  if [[ "${DRY_RUN}" == "false" ]]; then
    for pv in "${pvs[@]}"; do
      [[ -n "${pv}" ]] || continue
      wait_for_pv_absent "${pv}"
    done
  fi

  for path in "${REMOTE_PATHS[@]}"; do
    log "Remove path ${path}"
    run rm -rf --one-file-system "${path}"
  done

  log "Post-cleanup snapshot:"
  run_shell "kubectl get pod -A | grep -E 'anybackup-ai|kweaver|resource|v9-system|ingress-nginx|middleware' || true"
  run_shell "helm list -A -a | grep -E 'anybackup-ai|kweaver|resource|v9-system|ingress-nginx|middleware|failed|pending' || true"

  log "Cleanup finished."
}

main "$@"
