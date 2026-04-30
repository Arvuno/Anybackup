---
name: sandbox-cli-overlay
description: 处理 `foundation-cli` 注入 KWeaver sandbox 镜像的完整流程。用于查询基础镜像、生成覆盖 Dockerfile、同 tag 构建、镜像离线导出或分发到目标 Kubernetes 集群、重启 sandbox pod、以及验证镜像里是否包含 CLI。
---

# Sandbox CLI Overlay

## Overview

这个 skill 只解决一件事：让 KWeaver sandbox 运行的镜像里稳定带上 `foundation-cli`，并且尽量不改数据库里的镜像引用。

## Default Workflow

1. 先查当前 sandbox 模板实际使用的基础镜像。
2. 查不到时，才退回显式变量，例如 `sandbox_base_image`。
3. 基于当前基础镜像生成 Dockerfile：
   - `FROM <base-image>`
   - `COPY foundation-cli /usr/local/bin/foundation-cli`
4. 优先采用“同名同 tag 覆盖”策略。
5. 构建后执行镜像分发：
   - 需要离线交付时，先 `docker save`
   - 再按现场 `k8s` 集群实际方式导入节点运行时或推送到私有镜像仓库
6. 重启 sandbox 相关 pod，让新镜像真正被拉起。

## Guardrails

- 优先保留数据库中的镜像引用不变。
- 不要在没记录原镜像地址的情况下直接覆盖。
- `foundation-cli` 的来源默认是发布包里的 `bin/foundation-cli`。
- 如果现场明确要求新 tag 策略，再改数据库或 runtime 配置；否则默认不要主动换 tag。
- 默认按标准 `k8s` 处理，不要先入为主使用 `k3s ctr images import`；只有用户明确说明现场是 `k3s`，才用 `k3s` 专属命令。

## Validation

- 记录原始基础镜像。
- 验证导入后的镜像 tag 与模板 tag 保持一致。
- 验证镜像已经按现场口径进入目标 `k8s` 集群可见的运行时或仓库。
- 进入容器或用 `docker run` 校验：
  - `command -v foundation-cli`
- 如果有 smoke 脚本，确保 sandbox 镜像检查也被同步更新。
