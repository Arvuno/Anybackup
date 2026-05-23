---
name: kweaver-core-integrator
description: 将 KWeaver Core 作为黑盒组件接入 V9 发布流。用于处理上游 `deploy.sh` 的调用方式、KWeaver 配置渲染、裁剪安装、Ingress/LocalPV 禁用、manual init gate、以及 KWeaver 与 V9 之间的边界管理。
---

# KWeaver Core Integrator

## Overview

把 KWeaver 当成上游黑盒，不 fork，不改 `deploy.sh`。这个 skill 的核心是把它稳定嵌入 V9 的部署顺序，同时避免和 V9 自己的中间件、Ingress、StorageClass 发生重复安装或边界混乱。

## Integration Rules

1. 始终优先走“渲染配置 + 传参调用”，不要改上游脚本源码。
2. 默认调用方式是：
   - `./deploy.sh kweaver-core install`
   - `--minimum`
   - `--config=<rendered-config>`
3. 默认环境变量要显式带上：
   - `AUTO_INSTALL_INGRESS_NGINX=false`
   - `AUTO_INSTALL_LOCALPV=false`
4. KWeaver 配置模板的改动，只放在我们自己的发布包里，例如：
   - `config/kweaver-config.yaml.j2`
   - `ansible/group_vars/all.yml`
5. 保持 KWeaver 自己的依赖链现状，除非用户明确要求收敛：
   - `MariaDB`
   - `Kafka`
   - `KWeaver` 自己识别的 `Redis/OpenSearch`

## Default Sequence

1. 先完成 V9 侧基础中间件和 `local-path`。
2. 再部署 KWeaver Core。
3. 再做 sandbox CLI overlay。
4. 再进入手工初始化闸门：
   - `KWEAVER_BASE_URL`
   - `KWEAVER_TOKEN`
   - `preset_agents`
   - `backup_agent_id`
5. 最后再发布 `core-agent-service` 和业务服务。

## Guardrails

- 不要把 KWeaver 侧的部署资产复制进自研 chart。
- 不要默认暴露 KWeaver 外部入口。
- 不要在未确认的情况下把 V9 middleware URL 回灌给 KWeaver。
- 任何顺序变更，都要先检查 `manual_init_gate` 是否还成立。

## Validation

- 检查 `kweaver` 和 `resource` namespace 的 pod 就绪。
- 检查 `deploy.sh` 调用参数是否只来自我们的 Ansible，而不是脚本改造。
- 检查不存在重复安装的 `ingress-nginx` 和 `local-path`。
