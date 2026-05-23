---
name: design-doc-change-handler
description: 根据代码、Ansible、Helm chart、验证脚本和现有设计文档,更新 `docs/00_设计文档/` 下受影响的 Markdown 文档。Use when a design change, deployment flow change, middleware change, runtime dependency change, or documentation correction should land in one or more design docs.
---

# Design Doc Change Handler

把设计变更准确落到 `docs/00_设计文档/` 的具体文档里,优先维护现有文件,只在必要时新增设计文档。

## First Reads

- 先读 [../design-doc-orchestrator/references/source-map.md](../design-doc-orchestrator/references/source-map.md)。
- 再读 [references/update-checklist.md](references/update-checklist.md)。
- 如果需要新建文档,使用 [../design-doc-orchestrator/assets/design-doc-template.md](../design-doc-orchestrator/assets/design-doc-template.md) 起草。

## Workflow

### 1. Identify the changed facts

先确认设计结论来自哪里:

- `ansible/` 里的安装顺序、变量、前置检查。
- `helm-chart/` 里的组件、依赖、环境变量和 Service 命名。
- `scripts/validate_deploy.py` 与 `docs/操作文档/cluster-smoke.sh` 里的验证口径。
- 已有 `docs/00_设计文档/` 里的旧结论。

只在能找到代码或脚本事实支撑时,再把它写进正式设计文档。

### 2. Pick the landing document

按内容落位:

- 总体拓扑、部署总览、命名空间归属: 优先更新 `00_V9部署总体架构设计.md`。
- 部署时序、人工步骤、控制机/目标机职责: 优先更新部署时序类文档; 如果还没拆文档,就补进总体文档对应章节。
- 运行时依赖、环境变量、中间件依赖矩阵: 优先更新运行时依赖类文档; 如果还没拆文档,就补进总体文档或新增矩阵文档。
- 验证路径、回归策略、交付判断标准: 优先更新验证设计类文档。

如果一个变更跨了两类以上内容,优先维护主文档,再同步辅文档。

### 3. Edit with minimal disturbance

- 尽量保留已有标题、编号和主线。
- 只改受影响章节,避免整篇重写。
- 如果图已经过时,同步调用 `design-diagram-maintainer` 或至少标记需要补图。
- 如果旧章节已经失真,直接收敛到新结论,不要保留“旧结论 + 新补丁”的堆叠写法。

### 4. Sync nearby artifacts

修改正文后,同步检查:

- 图下说明文字是否还成立。
- 表格、清单、步骤编号是否需要跟着变。
- 同目录其他设计文档是否引用了旧名称或旧流程。
- 是否需要在文末补一条修改记录。

### 5. Validate when deployment meaning changed

如果改动触及以下内容,补跑:

```bash
python scripts/validate_deploy.py --mode regression
```

触发条件:

- `install.sh`
- `ansible/`
- `helm-chart/`
- `foundation/`
- `init-scripts/`
- `docs/操作文档/cluster-smoke.sh`
- `scripts/validate_deploy.py`

## Guardrails

- 不要把设计文档写成实现流水账; 重点写结构、边界、依赖和结论。
- 不要新增“临时说明”文档来回避更新正式设计文档。
- 不要把未验证的猜测写成确定结论。
- 不要遗漏人工步骤和自动化步骤之间的责任交接点。

## Done Criteria

- 受影响设计文档已更新到位。
- 图表、表格和文字至少在相邻范围内保持一致。
- 如果部署含义变了,已经运行或明确说明未运行验证脚本。
