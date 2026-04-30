# Review Checklist

## Severity

- `P1`
  - 会导致错误部署理解、错误组件边界、错误依赖关系或错误验证口径。
- `P2`
  - 不至于直接误导实现,但会造成明显歧义、遗漏关键步骤或图文冲突。
- `P3`
  - 元数据、标题、编号、轻微命名不一致或局部表达问题。

## Review Dimensions

- 结构是否清楚。
- 边界是否清楚。
- 命名是否和代码一致。
- 图文是否一致。
- 部署前提和验证前提是否一致。
- 是否保留了仍然重要的人工步骤说明。

## Required Output Shape

- Findings first
- Open questions or assumptions
- Brief pass/fix summary

## Common Risks

- 文档仍然沿用旧的中间件清单。
- Service 名和 URL 示例与实际 chart 不一致。
- 把控制机安装工具写成客户目标机前提。
- 图已经更新,正文却还在描述旧流程。
- 回归验证脚本已经变化,设计文档还停留在旧验证口径。
