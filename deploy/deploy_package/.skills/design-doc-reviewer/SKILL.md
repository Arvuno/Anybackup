---
name: design-doc-reviewer
description: 审查 `docs/00_设计文档/` 中设计文档的结构、边界、命名、图文一致性和部署事实是否成立。Use when the user asks for review/审查, or after creating or updating architecture, deployment, runtime dependency, or validation design docs.
---

# Design Doc Reviewer

优先找事实偏差、边界混乱、图文不一致和遗漏的验证风险,而不是帮作者重写一遍文档。

## First Reads

- 先读 [../design-doc-orchestrator/references/source-map.md](../design-doc-orchestrator/references/source-map.md)。
- 再读 [references/review-checklist.md](references/review-checklist.md)。
- 然后读取被 review 的设计文档和它引用的图表。

## Review Order

### 1. Verify the scope

先判断文档有没有回答清楚:

- 这份设计文档负责解释什么。
- 它不负责解释什么。
- 它和同目录其他设计文档怎么分工。

### 2. Check factual grounding

重点核对:

- 组件数量、名称、命名空间是否与 `ansible/` 和 `helm-chart/` 一致。
- 部署顺序、前置条件、人工步骤是否与脚本一致。
- 运行时依赖、环境变量、中间件归属是否与模板一致。
- 验证口径是否与 `scripts/validate_deploy.py` 和 `cluster-smoke.sh` 一致。

### 3. Check diagram-text consistency

- 图里有但文里没解释的节点或箭头。
- 文里写了但图里没有体现的关键关系。
- 图文对同一个组件使用了不同名称。

### 4. Check maintainability

- 章节主线是否稳定。
- 有没有把实现细节堆成流水账。
- 有没有把临时方案伪装成长期设计。

## Output Format

输出必须遵循这个顺序:

1. Findings first
   - 按 `P1`、`P2`、`P3` 排序。
   - 每条都说明受影响的文档和问题原因。
2. Open questions or assumptions
   - 只保留会影响设计判断的关键前提。
3. Brief summary
   - 说明是否建议直接通过、修改后通过,或暂不通过。

## Guardrails

- 不要把风格偏好当成缺陷。
- 不要在没有代码或脚本证据时判定“文档错误”。
- 不要以“建议补充更多背景”替代指出真正的问题。
- 如果没有发现问题,要明确说“未发现问题”,并补一句残余风险或未验证项。
