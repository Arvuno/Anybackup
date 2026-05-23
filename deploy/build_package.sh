#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
CONFIG_FILE="${SCRIPT_DIR}/build_config.json"
TEMP_DIR=""
SERVICE_VERSIONS=""

cleanup_on_exit() {
    if [ -n "$TEMP_DIR" ] && [ -d "$TEMP_DIR" ]; then
        echo "[INFO] 清理临时目录: ${TEMP_DIR}"
        rm -rf "${TEMP_DIR}"
    fi
}

trap cleanup_on_exit EXIT

log_info() {
    echo "[INFO] $1"
}

log_error() {
    echo "[ERROR] $1" >&2
}

read_config() {
    if [ ! -f "$CONFIG_FILE" ]; then
        log_error "配置文件不存在: $CONFIG_FILE"
        exit 1
    fi
    
    PACKAGE_NAME=$(jq -r '.package_name' "$CONFIG_FILE")
    PACKAGE_VERSION=$(jq -r '.package_version' "$CONFIG_FILE")
    PLATFORM=$(jq -r '.platform' "$CONFIG_FILE")
    BUILD_TARGET=$(jq -r '.build_target' "$CONFIG_FILE")
    LANGUAGE=$(jq -r '.language' "$CONFIG_FILE")
    PRODUCT_TYPE=$(jq -r '.product_type' "$CONFIG_FILE")
    BUILD_NUMBER=$(jq -r '.build_number' "$CONFIG_FILE")
    PACKAGE_BRANCH=$(jq -r '.package_branch' "$CONFIG_FILE")
    
    if [ "$LANGUAGE" = "null" ] || [ -z "$LANGUAGE" ]; then
        LANGUAGE="zh_CN"
    fi
    if [ "$PRODUCT_TYPE" = "null" ] || [ -z "$PRODUCT_TYPE" ]; then
        PRODUCT_TYPE="ABNormal"
    fi
    if [ "$BUILD_NUMBER" = "null" ] || [ -z "$BUILD_NUMBER" ]; then
        BUILD_NUMBER="1"
    fi
    if [ "$PACKAGE_BRANCH" = "null" ] || [ -z "$PACKAGE_BRANCH" ]; then
        PACKAGE_BRANCH="AB9.0.0"
    fi
    
    CONSOLE_PY_DIR=$(jq -r '.directories.console_py' "$CONFIG_FILE")
    CONSOLE_GO_DIR=$(jq -r '.directories.console_go' "$CONFIG_FILE")
    OTHER_SERVICES_DIR=$(jq -r '.directories.other_services' "$CONFIG_FILE")
    CLIENT_DIR=$(jq -r '.directories.client' "$CONFIG_FILE")
    STORAGE_DIR=$(jq -r '.directories.storage' "$CONFIG_FILE")
    
    if [ "$CONSOLE_PY_DIR" = "null" ] || [ -z "$CONSOLE_PY_DIR" ]; then
        CONSOLE_PY_DIR="${PROJECT_ROOT}/../consolePy"
    fi
    if [ "$CONSOLE_GO_DIR" = "null" ] || [ -z "$CONSOLE_GO_DIR" ]; then
        CONSOLE_GO_DIR="${PROJECT_ROOT}/../consoleGo"
    fi
    if [ "$OTHER_SERVICES_DIR" = "null" ] || [ -z "$OTHER_SERVICES_DIR" ]; then
        OTHER_SERVICES_DIR="${PROJECT_ROOT}/../otherServices"
    fi
    if [ "$CLIENT_DIR" = "null" ] || [ -z "$CLIENT_DIR" ]; then
        CLIENT_DIR="${PROJECT_ROOT}/../client"
    fi
    if [ "$STORAGE_DIR" = "null" ] || [ -z "$STORAGE_DIR" ]; then
        STORAGE_DIR="${PROJECT_ROOT}/../proplyd"
    fi
    
    GO_SERVICE_BASE_PATH=$(jq -r '.go_services_config.service_base_path' "$CONFIG_FILE")
    GO_PACKAGE_SUBPATH=$(jq -r '.go_services_config.package_subpath' "$CONFIG_FILE")
    GO_PACKAGE_NAME=$(jq -r '.go_services_config.package_name' "$CONFIG_FILE")
    
    OTHER_SERVICE_BASE_PATH=$(jq -r '.other_services_config.base_path' "$CONFIG_FILE")
    OTHER_BINARY_DIR=$(jq -r '.other_services_config.binary_dir' "$CONFIG_FILE")
    OTHER_PACKAGE_NAME=$(jq -r '.other_services_config.package_name' "$CONFIG_FILE")
    
    CLIENT_APP_BASE_PATH=$(jq -r '.client_services_config.app_base_path' "$CONFIG_FILE")
    CLIENT_PACKAGE_SUBPATH=$(jq -r '.client_services_config.package_subpath' "$CONFIG_FILE")
    CLIENT_TARGET_DIR=$(jq -r '.client_services_config.target_dir' "$CONFIG_FILE")
    
    STORAGE_PACKAGE_SUBPATH=$(jq -r '.storage_config.package_subpath' "$CONFIG_FILE")
    STORAGE_PACKAGE_NAME=$(jq -r '.storage_config.package_name' "$CONFIG_FILE")
    
    GO_SERVICES_FILTER=$(jq -r '.go_services | @json' "$CONFIG_FILE")
    OTHER_SERVICES_FILTER=$(jq -r '.other_services | @json' "$CONFIG_FILE")
    CLIENT_SERVICES_FILTER=$(jq -r '.client_services | @json' "$CONFIG_FILE")
    
    log_info "配置加载完成:"
    log_info "  包名: ${PACKAGE_NAME}"
    log_info "  版本: ${PACKAGE_VERSION}"
    log_info "  平台: ${PLATFORM}"
    log_info "  构建目标: ${BUILD_TARGET}"
    log_info "  语言: ${LANGUAGE}"
    log_info "  产品类型: ${PRODUCT_TYPE}"
    log_info "  构建号: ${BUILD_NUMBER}"
    log_info "  分支: ${PACKAGE_BRANCH}"
    log_info "  Py服务目录: ${CONSOLE_PY_DIR}"
    log_info "  Go服务目录: ${CONSOLE_GO_DIR}"
    log_info "  其他服务目录: ${OTHER_SERVICES_DIR}"
    log_info "  客户端目录: ${CLIENT_DIR}"
    log_info "  存储服务目录: ${STORAGE_DIR}"
}

is_service_allowed() {
    local service_name="$1"
    local filter_json="$2"
    
    if [ "$filter_json" = "null" ] || [ "$filter_json" = "[]" ]; then
        return 0
    fi
    
    local found=$(echo "$filter_json" | jq -r --arg name "$service_name" 'index($name) != null')
    
    if [ "$found" = "true" ]; then
        return 0
    else
        return 1
    fi
}

should_strip_root_dir() {
    local content_dir="$1"
    
    local items=($(ls -A "$content_dir"))
    local item_count=${#items[@]}
    
    if [ "$item_count" -eq 1 ] && [ -d "$content_dir/${items[0]}" ]; then
        local inner_dir="$content_dir/${items[0]}"
        local inner_items=($(ls -A "$inner_dir"))
        
        local has_etc=false
        local has_service=false
        
        for item in "${inner_items[@]}"; do
            if [ "$item" = "etc" ] && [ -d "$inner_dir/etc" ]; then
                has_etc=true
            elif [ -d "$inner_dir/$item" ] || [ -f "$inner_dir/$item" ]; then
                if [[ "$item" != "etc" ]]; then
                    has_service=true
                fi
            fi
        done
        
        if [ "$has_etc" = true ] && [ "$has_service" = true ]; then
            return 0
        fi
    fi
    
    return 1
}

extract_service_version() {
    local tar_file="$1"
    local service_name="$2"
    
    local tar_basename=$(basename "$tar_file" .tar.gz)
    
    echo "${tar_basename}"
}

extract_and_merge() {
    local tar_file="$1"
    local target_dir="$2"
    local service_name="$3"
    
    local tar_version=$(extract_service_version "$tar_file" "$service_name")
    if [ -n "$SERVICE_VERSIONS" ]; then
        SERVICE_VERSIONS="${SERVICE_VERSIONS}"$'\n'"${tar_version}"
    else
        SERVICE_VERSIONS="${tar_version}"
    fi
    
    local extract_temp=$(mktemp -d -t "extract_${service_name}_XXXXXX")
    
    log_info "  解压到临时目录: ${extract_temp}"
    tar -xzf "$tar_file" -C "$extract_temp"
    
    local content_dir="$extract_temp"
    
    if should_strip_root_dir "$extract_temp"; then
        local subdirs=($(ls -A "$extract_temp"))
        content_dir="$extract_temp/${subdirs[0]}"
        log_info "  检测到服务包结构，移除外层目录: ${subdirs[0]}"
    else
        log_info "  保持原有目录结构"
    fi
    
    log_info "  合并内容到目标目录..."
    
    for item in "$content_dir"/*; do
        if [ ! -e "$item" ]; then
            continue
        fi
        
        local item_name=$(basename "$item")
        local target_item="${target_dir}/${item_name}"
        
        if [ -d "$item" ]; then
            if [ -d "$target_item" ]; then
                log_info "    合并目录: ${item_name}/"
                cp -rf "$item"/* "$target_item"/ 2>/dev/null || true
            else
                log_info "    拷贝目录: ${item_name}/"
                cp -r "$item" "$target_item"
            fi
        else
            if [ -f "$target_item" ]; then
                log_info "    覆盖文件: ${item_name}"
            else
                log_info "    拷贝文件: ${item_name}"
            fi
            cp -f "$item" "$target_item"
        fi
    done
    
    rm -rf "$extract_temp"
}

set_executable_permissions() {
    log_info "设置可执行权限..."
    
    local sh_count=0
    local venv_bin_count=0
    
    for sh_file in $(find "${PACKAGE_DIR}" -name "*.sh" -type f 2>/dev/null); do
        chmod +x "$sh_file"
        sh_count=$((sh_count + 1))
    done
    log_info "  设置 *.sh 可执行权限: ${sh_count} 个文件"
    
    local venv_bin_dir="${PACKAGE_DIR}/virtualenv3/bin"
    if [ -d "$venv_bin_dir" ]; then
        for bin_file in "$venv_bin_dir"/*; do
            if [ -f "$bin_file" ]; then
                chmod +x "$bin_file"
                venv_bin_count=$((venv_bin_count + 1))
            fi
        done
        log_info "  设置 virtualenv3/bin 可执行权限: ${venv_bin_count} 个文件"
    fi
}

generate_version_details() {
    log_info "生成 VersionDetails 文件..."
    
    local build_date=$(date +%Y%m%d)
    local build_target_lower=$(echo "$BUILD_TARGET" | tr '[:upper:]' '[:lower:]')
    
    local full_package_name="${PACKAGE_NAME}-${PLATFORM}-${PACKAGE_VERSION}-${build_date}-${build_target_lower}-${LANGUAGE}-${PRODUCT_TYPE}-${BUILD_NUMBER}"
    
    local product_version=$(echo "$PACKAGE_VERSION" | cut -d'.' -f1,2)
    
    local version_file="${PACKAGE_DIR}/VersionDetails"
    
    {
        echo "Package Name=${full_package_name}"
        echo "Product Package Name=${full_package_name}"
        echo "Package Branch=${PACKAGE_BRANCH}"
        echo "Package Version=${PACKAGE_VERSION}"
        echo "Product Version=${product_version}"
        if [ -n "$SERVICE_VERSIONS" ]; then
            echo "$SERVICE_VERSIONS"
        fi
    } > "$version_file"
    
    log_info "  VersionDetails 已生成: ${version_file}"
}

setup_build_dir() {
    SETUP_DIR="${SCRIPT_DIR}/setup"
    
    if [ ! -d "$SETUP_DIR" ]; then
        log_error "setup目录不存在: ${SETUP_DIR}"
        exit 1
    fi
    
    TEMP_DIR=$(mktemp -d -t "${PACKAGE_NAME}_XXXXXX")
    log_info "创建临时目录: ${TEMP_DIR}"
    
    log_info "拷贝setup目录到临时目录..."
    cp -r "${SETUP_DIR}" "${TEMP_DIR}/"
    
    PACKAGE_DIR="${TEMP_DIR}/${PACKAGE_NAME}"
    
    log_info "重命名setup目录为${PACKAGE_NAME}"
    mv "${TEMP_DIR}/setup" "${PACKAGE_DIR}"
    
    log_info "构建目录准备完成: ${PACKAGE_DIR}"
}

copy_py_services() {
    log_info "开始拷贝Python服务..."
    
    PY_SERVICES_COUNT=$(jq '.py_services | length' "$CONFIG_FILE")
    
    if [ "$PY_SERVICES_COUNT" -eq 0 ]; then
        log_info "没有配置Python服务，跳过"
        return 0
    fi
    
    for ((i=0; i<PY_SERVICES_COUNT; i++)); do
        SERVICE_NAME=$(jq -r ".py_services[$i].name" "$CONFIG_FILE")
        SOURCE_PATH=$(jq -r ".py_services[$i].source_path" "$CONFIG_FILE")
        
        SERVICE_DIR="${CONSOLE_PY_DIR}/${SOURCE_PATH}"
        
        log_info "处理Python服务: ${SERVICE_NAME}"
        log_info "  源目录: ${SERVICE_DIR}"
        
        if [ ! -d "$SERVICE_DIR" ]; then
            log_error "服务目录不存在: ${SERVICE_DIR}"
            continue
        fi
        
        TAR_FILES=$(find "$SERVICE_DIR" -maxdepth 1 -name "*.tar.gz" -type f)
        
        if [ -z "$TAR_FILES" ]; then
            log_error "未找到压缩包文件: ${SERVICE_DIR}"
            continue
        fi
        
        for TAR_FILE in $TAR_FILES; do
            log_info "  解压: $(basename "$TAR_FILE")"
            extract_and_merge "$TAR_FILE" "${PACKAGE_DIR}" "$SERVICE_NAME"
        done
    done
    
    log_info "Python服务拷贝完成"
}

copy_go_services() {
    log_info "开始拷贝Go服务..."
    
    GO_SERVICE_DIR="${CONSOLE_GO_DIR}/${GO_SERVICE_BASE_PATH}"
    
    if [ ! -d "$GO_SERVICE_DIR" ]; then
        log_error "Go服务目录不存在: ${GO_SERVICE_DIR}"
        return 0
    fi
    
    MISSING_PACKAGES=()
    SUCCESS_COUNT=0
    
    for SERVICE_DIR in "$GO_SERVICE_DIR"/*; do
        if [ ! -d "$SERVICE_DIR" ]; then
            continue
        fi
        
        SERVICE_NAME=$(basename "$SERVICE_DIR")
        SERVICE_PACKAGE_DIR="${SERVICE_DIR}/${GO_PACKAGE_SUBPATH}"
        TAR_FILE="${SERVICE_PACKAGE_DIR}/${GO_PACKAGE_NAME}"
        
        log_info "处理Go服务目录: ${SERVICE_NAME}"
        log_info "  源目录: ${SERVICE_PACKAGE_DIR}"
        
        if [ ! -f "$TAR_FILE" ]; then
            log_error "  未找到包文件: ${TAR_FILE}"
            MISSING_PACKAGES+=("$SERVICE_NAME")
            continue
        fi
        
        log_info "  解压: ${GO_PACKAGE_NAME}"
        extract_and_merge "$TAR_FILE" "${PACKAGE_DIR}" "$SERVICE_NAME"
        
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    done
    
    log_info "Go服务拷贝完成: 成功 ${SUCCESS_COUNT} 个"
    
    if [ ${#MISSING_PACKAGES[@]} -gt 0 ]; then
        log_info "以下Go服务目录缺少包文件:"
        for SERVICE in "${MISSING_PACKAGES[@]}"; do
            log_info "  - ${SERVICE}"
        done
    fi
}

copy_other_services() {
    log_info "开始拷贝其他服务..."
    
    OTHER_SERVICE_DIR="${OTHER_SERVICES_DIR}/${OTHER_SERVICE_BASE_PATH}"
    
    if [ ! -d "$OTHER_SERVICE_DIR" ]; then
        log_error "其他服务目录不存在: ${OTHER_SERVICE_DIR}"
        return 0
    fi
    
    MISSING_PACKAGES=()
    SUCCESS_COUNT=0
    
    for SERVICE_DIR in "$OTHER_SERVICE_DIR"/*; do
        if [ ! -d "$SERVICE_DIR" ]; then
            continue
        fi
        
        SERVICE_NAME=$(basename "$SERVICE_DIR")
        BINARY_DIR="${SERVICE_DIR}/${OTHER_BINARY_DIR}"
        TAR_FILE="${BINARY_DIR}/${OTHER_PACKAGE_NAME}"
        
        log_info "处理其他服务目录: ${SERVICE_NAME}"
        log_info "  源目录: ${BINARY_DIR}"
        
        if [ ! -f "$TAR_FILE" ]; then
            log_error "  未找到包文件: ${TAR_FILE}"
            MISSING_PACKAGES+=("$SERVICE_NAME")
            continue
        fi
        
        log_info "  解压: ${OTHER_PACKAGE_NAME}"
        extract_and_merge "$TAR_FILE" "${PACKAGE_DIR}" "$SERVICE_NAME"
        
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    done
    
    log_info "其他服务拷贝完成: 成功 ${SUCCESS_COUNT} 个"
    
    if [ ${#MISSING_PACKAGES[@]} -gt 0 ]; then
        log_info "以下其他服务目录缺少包文件:"
        for SERVICE in "${MISSING_PACKAGES[@]}"; do
            log_info "  - ${SERVICE}"
        done
    fi
}

copy_client() {
    log_info "开始拷贝客户端..."
    
    CLIENT_APP_DIR="${CLIENT_DIR}/${CLIENT_APP_BASE_PATH}"
    
    if [ ! -d "$CLIENT_APP_DIR" ]; then
        log_error "客户端目录不存在: ${CLIENT_APP_DIR}"
        return 0
    fi
    
    TARGET_DIR="${PACKAGE_DIR}/${CLIENT_TARGET_DIR}"
    mkdir -p "${TARGET_DIR}"
    
    MISSING_PACKAGES=()
    SUCCESS_COUNT=0
    SKIPPED_COUNT=0
    
    for CLIENT_SERVICE_DIR in "$CLIENT_APP_DIR"/*; do
        if [ ! -d "$CLIENT_SERVICE_DIR" ]; then
            continue
        fi
        
        SERVICE_NAME=$(basename "$CLIENT_SERVICE_DIR")
        
        if ! is_service_allowed "$SERVICE_NAME" "$CLIENT_SERVICES_FILTER"; then
            log_info "跳过客户端: ${SERVICE_NAME} (不在配置列表中)"
            SKIPPED_COUNT=$((SKIPPED_COUNT + 1))
            continue
        fi
        
        PACKAGE_DIR_PATH="${CLIENT_SERVICE_DIR}/${CLIENT_PACKAGE_SUBPATH}"
        
        log_info "处理客户端: ${SERVICE_NAME}"
        log_info "  源目录: ${PACKAGE_DIR_PATH}"
        
        if [ ! -d "$PACKAGE_DIR_PATH" ]; then
            log_error "  包目录不存在: ${PACKAGE_DIR_PATH}"
            MISSING_PACKAGES+=("$SERVICE_NAME")
            continue
        fi
        
        TAR_FILES=$(find "$PACKAGE_DIR_PATH" -maxdepth 1 -name "*.tar.gz" -type f)
        
        if [ -z "$TAR_FILES" ]; then
            log_error "  未找到压缩包文件"
            MISSING_PACKAGES+=("$SERVICE_NAME")
            continue
        fi
        
        for TAR_FILE in $TAR_FILES; do
            ORIGINAL_NAME=$(basename "$TAR_FILE")
            BASE_NAME="${ORIGINAL_NAME%.tar.gz}"
            
            SPECIAL_NAME=$(jq -r ".client_services_config.special_names.${SERVICE_NAME}" "$CONFIG_FILE")
            
            if [ "$SPECIAL_NAME" != "null" ] && [ -n "$SPECIAL_NAME" ]; then
                NEW_NAME="${SPECIAL_NAME}-${PLATFORM}-latest.tar.gz"
            else
                NEW_NAME="${BASE_NAME}-${PLATFORM}-latest.tar.gz"
            fi
            
            log_info "  拷贝: ${ORIGINAL_NAME} -> ${NEW_NAME}"
            cp "$TAR_FILE" "${TARGET_DIR}/${NEW_NAME}"
            
            SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
        done
    done
    
    log_info "客户端拷贝完成: 成功 ${SUCCESS_COUNT} 个, 跳过 ${SKIPPED_COUNT} 个"
    
    if [ ${#MISSING_PACKAGES[@]} -gt 0 ]; then
        log_info "以下客户端缺少包文件:"
        for SERVICE in "${MISSING_PACKAGES[@]}"; do
            log_info "  - ${SERVICE}"
        done
    fi
}

copy_storage() {
    log_info "开始拷贝存储服务..."
    
    if [ "$STORAGE_PACKAGE_SUBPATH" = "null" ] || [ -z "$STORAGE_PACKAGE_SUBPATH" ]; then
        log_info "未配置存储服务路径，跳过"
        return 0
    fi
    
    STORAGE_PACKAGE_DIR="${STORAGE_DIR}/${STORAGE_PACKAGE_SUBPATH}"
    TAR_FILE="${STORAGE_PACKAGE_DIR}/${STORAGE_PACKAGE_NAME}"
    
    log_info "  源目录: ${STORAGE_PACKAGE_DIR}"
    log_info "  包文件: ${STORAGE_PACKAGE_NAME}"
    
    if [ ! -f "$TAR_FILE" ]; then
        log_error "  未找到存储服务包文件: ${TAR_FILE}"
        return 0
    fi
    
    log_info "  解压存储服务包..."
    extract_and_merge "$TAR_FILE" "${PACKAGE_DIR}" "storage"
    
    log_info "存储服务拷贝完成"
}

set_ownership() {
    log_info "设置属组为 2048:2048..."
    chown -R 2048:2048 "${PACKAGE_DIR}"
}

create_package() {
    log_info "开始创建最终安装包..."
    
    PACKAGES_DIR="${SCRIPT_DIR}/packages"
    mkdir -p "${PACKAGES_DIR}"
    
    FINAL_TAR_NAME="${PACKAGE_NAME}-${PACKAGE_VERSION}-${PLATFORM}.tar.gz"
    FINAL_TAR_PATH="${PACKAGES_DIR}/${FINAL_TAR_NAME}"
    
    cd "${TEMP_DIR}"
    
    log_info "创建压缩包: ${FINAL_TAR_NAME}"
    tar -czf "${FINAL_TAR_PATH}" "${PACKAGE_NAME}"
    
    log_info "安装包创建完成: ${FINAL_TAR_PATH}"
    
    PACKAGE_SIZE=$(du -h "${FINAL_TAR_PATH}" | cut -f1)
    log_info "安装包大小: ${PACKAGE_SIZE}"
}

cleanup() {
    log_info "清理临时文件..."
    
    if [ -n "$TEMP_DIR" ] && [ -d "$TEMP_DIR" ]; then
        log_info "删除临时目录: ${TEMP_DIR}"
        rm -rf "${TEMP_DIR}"
    fi
}

main() {
    log_info "========================================="
    log_info "开始构建整包: ${PACKAGE_NAME}"
    log_info "========================================="
    
    read_config
    setup_build_dir
    
    copy_py_services
    copy_go_services
    copy_other_services
    copy_client
    copy_storage
    
    generate_version_details
    set_executable_permissions
    set_ownership
    create_package
    cleanup
    
    log_info "========================================="
    log_info "整包构建完成！"
    log_info "========================================="
}

main "$@"
