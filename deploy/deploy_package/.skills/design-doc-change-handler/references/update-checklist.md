# Design Doc Update Checklist

## Change Classes

- 拓扑变化: 组件新增、删除、重命名,或连接关系变化。
- 部署变化: 安装顺序、初始化动作、前置条件、控制机职责变化。
- 运行时变化: 环境变量、依赖中间件、回调链路、Service 关系变化。
- 验证变化: smoke check、回归脚本、可交付判断标准变化。
- 文档结构变化: 一份文档承载不下,需要拆分或新增设计文档。

## Landing Guidance

- 先更新现有文档。
- 只有当现有文档的职责边界已被破坏时才新增。
- 新增文档时,标题要直接反映职责,并保持编号递增。

## Mandatory Sync Points

- 图文是否一致。
- 章节标题和文内术语是否一致。
- 组件名、命名空间名、Service 名和 chart/脚本是否一致。
- 如果新增了人工步骤,是否明确了执行者和前置条件。
- 如果删掉了步骤,是否同步移除了过期风险提示或验证步骤。

## Validation Follow-up

部署含义变化后优先运行:

```bash
python scripts/validate_deploy.py --mode regression
```

只有当发布口径或客户交付判断发生变化时,才额外考虑:

```bash
python scripts/validate_deploy.py --mode release
```
