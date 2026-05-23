---
name: v9-release-packager
description: 维护 AnyBackup Agent/V9 的发布包和服务发布链路。用于整理 `anybackup-agent-release` 目录、同步 Ansible/Helm/脚本/文档、执行业务服务镜像 build-import-release 流程、以及保持交付物与部署现实一致。
---

# V9 Release Packager

## Overview

这个 skill 关注发布包视角，而不是单个 chart 视角。它负责把中间件、KWeaver、业务服务、离线镜像、文档和脚本收拢成一套能交付、能复跑、边界清楚的目录。

## Core Responsibilities

1. 维护发布根目录结构，例如：
   - `ansible/`
   - `helm-chart/`
   - `images/`
   - `bin/`
   - `config/`
   - `docs/`
   - `init-scripts/`
2. 保持三层部署资产的职责清晰：
   - `deploy.sh` 是上游黑盒
   - `Helm chart` 是自研部署模板
   - `Ansible` 是总编排
3. 处理业务服务发布链路：
   - `git clone/pull`
   - `docker build`
   - `docker save`
   - 按目标 `k8s` 集群口径分发镜像（节点运行时导入或私有仓库推送）
   - `helm upgrade --install`
4. 保持发布顺序稳定：
   - middleware
   - KWeaver
   - sandbox overlay
   - manual init
   - core-agent
   - business services
5. 同步更新交付入口：
   - `install.sh`
   - `images/load-images.sh`
   - `README.md`
   - `docs/*`

## Guardrails

- 不要把上游 `E:\\Code\\kweaver-core\\deploy` 复制进发布包再改。
- 不要用裸 `kubectl apply -f` 替代已经确定用 Helm 管理的服务。
- 默认以标准 `k8s` 为交付前提，不要把 `k3s` 专属命令写进通用脚本或文档，除非环境已被明确标注为 `k3s`。
- 不要伪造真实资产：
  - 离线镜像 tar
  - Foundation 安装包
  - Skill zip
  - 真实 init SQL
- 如果只是占位，就明确标成占位。

## Validation

- 检查 chart、Ansible、脚本、文档是否说的是同一套拓扑。
- 检查 `images/load-images.sh` 能覆盖新增目录层级。
- 检查业务服务发布仍然遵守：
  - 先 `core-agent-service`
  - 再 `v9-services`
- 改完后至少跑一次：
  - `helm lint`
  - `helm template`
  - 能跑的话再补 `ansible-playbook --syntax-check`
