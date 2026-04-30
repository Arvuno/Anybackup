# Domain Entry Skill Template

下面模板用于创建或重构某个业务域的入口型 `SKILL.md`。

将占位符替换为实际 domain 信息：

- `{{domain}}`
- `{{domain_title}}`
- `{{domain_intents}}`
- `{{not_applicable_cases}}`
- `{{quick_rules}}`
- `{{write_commands}}`
- `{{read_commands}}`
- `{{global_context}}`
- `{{misroutes}}`
- `{{intent_rows}}`

```md
---
name: foundation-cli-{{domain}}
description: 当用户需要通过 `foundation-cli {{domain}}` 执行 {{domain_title}} 相关命令时使用。
metadata:
  requires:
    bins: ["foundation-cli"]
  cliHelp: "foundation-cli {{domain}} --help"
---

# Foundation CLI {{domain_title}} 技能

当任务属于 `foundation-cli {{domain}}` 域命令时，使用本技能。

本文件是入口，不是完整参考手册。
它的作用是帮助 agent：
1. 判断用户意图是否属于 `{{domain}}` 域。
2. 选择正确的 `{{domain}}` 命令。
3. 跳转到 `references/commands/*.md` 中对应的单命令文档。

## 适用场景

{{domain_intents}}

## 不适用场景

{{not_applicable_cases}}

## 标准流程

1. 先使用“快速判断规则”排除最常见的误选命令场景。
2. 再使用 [命令映射](./references/commands.md) 选择命令。
3. 立即打开 `./references/commands/` 下对应的单命令文档。
4. 基于该单命令文档构造最终 CLI。

不要把本文件当成查询完整请求体、示例或排错细节的地方。

## 快速判断规则

{{quick_rules}}

## 写入与只读

写入命令：

{{write_commands}}

只读命令：

{{read_commands}}

对于写入命令，优先使用 `--data` 传递结构化载荷。

## 必要的全局上下文

除非调用环境已经隐式提供，否则先确认这些参数：

{{global_context}}

## 常见误判

{{misroutes}}

## 意图到命令

| 意图 | 推荐命令 | 下一步打开的文件 |
|---|---|---|
{{intent_rows}}

## 参考资料

- [命令映射](./references/commands.md)
```

## 使用建议

1. 保持入口页短小，优先承担分流职责。
2. 如果某条命令需要大量细节，放到单命令文档，不要塞回入口页。
3. 如果某组命令极易混淆，把这组对比写进“快速判断规则”而不是只写在大表里。
