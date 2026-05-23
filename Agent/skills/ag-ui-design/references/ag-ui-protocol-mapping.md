# AG-UI 协议映射

## 设计边界

`ag-ui-design` 的输出仍然是 Markdown 设计约束，不直接生成 draft JSON。设计约束必须让后续 `ag-ui-response` 能明确生成 AG-UI 协议字段，避免只给“好看一点”的视觉描述。

每次设计都要说明：

- 要创建新的可见卡片，还是更新已有卡片。
- 要使用哪一种 `meta.intent`：`thought`、`tool_call`、`clarification`、`result` 或 `error`。
- 可见内容如何组成 `conversation.ui.layout-tree@1` 的 layout tree。
- 用户动作是否需要进入 `actions`，以及哪个 `action-row` 节点承载这些动作。
- 当前交互状态是否需要 `state_snapshot` 或 `state_delta`。

## 布局树约束

Markdown 的“布局约束”章节必须能映射到 `ACTIVITY_SNAPSHOT` 或 `ACTIVITY_DELTA`：

- 新建卡片：描述一个完整 activity，后续由 `ag-ui-response` 生成 `ACTIVITY_SNAPSHOT`，`activityType` 固定为 `conversation.ui.layout-tree`，`content.contract` 固定为 `conversation.ui.layout-tree@1`。
- 更新卡片：描述要更新的已有卡片、目标节点和变化字段，后续由 `ag-ui-response` 生成 `ACTIVITY_DELTA`，只能 patch 已存在 activity，不创建新的 activity。
- 根节点优先使用 `stack`；多个并列方案才使用 `grid`；单个方案、恢复点、告警或审查对象使用 `card`。
- 每个可见卡片必须定义稳定的业务块含义，便于后续生成 `message_id`、`block_id` 和更新路径。
- 不要把业务结果设计成 paragraph-only；结果、工具调用、错误、候选方案都必须包含 `kv-list`、`data-table`、`callout`、`metric-list`、`card` 或 `action-row` 中的结构化节点。

## 卡片内容模型

卡片内容按“标题区、摘要区、事实区、依据区、风险区、动作区、状态区”设计，不要求每张卡全部包含，但必须说明保留和省略的理由。

| 区域 | 推荐组件 | 内容要求 |
| --- | --- | --- |
| 标题区 | `heading`、`badge-row` | 业务对象、结论、状态徽标、推荐或待审核标记 |
| 摘要区 | `paragraph`、`callout` | 一到两句结论，不承载内部推理链 |
| 事实区 | `kv-list`、`data-table`、`metric-list` | 系统、数据库、时间点、数据量、RPO/RTO、地址、策略、命中数量 |
| 依据区 | `callout`、`data-table`、详情折叠区 | 只显示脱敏证据和可验证事实 |
| 风险区 | `callout`、`badge-row` | 影响范围、失败原因、冲突、风险等级 |
| 动作区 | `action-row` | 用户可执行动作，主动作唯一，文案使用业务语义 |
| 状态区 | `metric-list`、`badge-row`、`state_snapshot`/`state_delta` | 执行中、待确认、已完成、错误、选择状态、当前视图 |

## 组件选择矩阵

| 内容类型 | 首选组件 | 禁止做法 |
| --- | --- | --- |
| 单个事实摘要 | `kv-list` | 写成连续长段落 |
| 同构列表 | `data-table` | 用 bullet 堆叠大量字段 |
| 多方案比较 | `grid` + `card` | 把所有方案混在一张普通表里且没有动作 |
| 风险、空结果、失败原因 | `callout` | 只用红色文字描述 |
| 状态、优先级、信心值 | `badge-row`、`metric-list` | 只靠自然语言表达状态 |
| 用户确认或选择 | `action-row` + `actions` | 只在正文里说“请选择” |
| 趋势或报告附件 | `chart`、`attachment-list` | 粘贴原始日志或不可控链接 |

## 动作绑定规则

Markdown 约束中出现用户动作时，必须同时描述动作的协议意图：

- `action-row` 只展示动作按钮；真实动作定义由同一个 activity 的 `actions` 承载。
- 每个动作要有稳定业务语义，例如“确认继续”“调整参数”“审核通过”“选择该方案”“查看详情”。
- 主动作只能有一个，使用 primary 样式；其他动作使用 secondary 或 link 样式。
- 会继续对话的动作映射为 `submit_message`；查看受控引用映射为 `open_ref`；复制文本映射为 `copy_text`；外部跳转需要有明确业务理由才使用 `open_url`。
- 高风险动作必须要求确认文案，后续由 `ag-ui-response` 生成 `confirmText`。

## 状态映射规则

Markdown 约束必须说明状态是否改变，方便后续生成 `STATE_SNAPSHOT` 或 `STATE_DELTA`：

- 新建交互结果时，优先定义 `state_snapshot`，至少描述 `interaction.status`。
- 更新执行进度时，定义 `state_delta`，只描述变更路径和值。
- `interaction.status` 只能使用业务阶段语义：`thinking`、`clarifying`、`executing`、`completed`、`error` 或回到 `idle`。
- 需要用户选择时，说明 `selection.required`、默认推荐项、是否锁定选择，以及按钮选择后应写入的选择字段。
- 需要控制界面展开或激活区域时，说明 `view.activeBlockIds` 或同类视图状态。

## 序列和更新约束

- 设计新卡片时不需要生成真实 `message_id`，但必须说明卡片的业务块含义，便于 `ag-ui-response` 生成稳定标识。
- 更新已有卡片时必须明确“更新已有卡片”，并说明目标卡片和要变化的字段；后续 draft 必须携带已有 `message_id`。
- 每次可见输出都对应递增的 `sequence`；同一 `sequence` 的重复更新会覆盖该序号当前值。
- 思考条和业务卡片可以分开成不同 activity；思考条不得承载主要业务结果。

## 结果源一致性规则

- 终态 `result`、最终报告和终态 `error` 必须在 AI 已完成结果输出后再进入 AG-UI 设计和发布流程。
- 将 AI 已完成结果输出作为结果源文本，并在发布前冻结；后续 layout tree 只做结构化映射，不重新推断业务结论。
- 不得提前发布终态 `result`；如果最终结论尚未稳定，只能发布进度、工具状态、主动询问或阶段性结果。
- `ag-ui-response` 生成 draft 前必须做一致性校验：`content`、卡片标题、核心字段、风险和动作不得偏离结果源文本。
- 不得新增结果源文本之外的事实；如果 AG-UI 需要额外字段，先让 AI 补全结果源文本，再重新生成 Markdown 设计约束和 draft。

## 推理链边界

- 可见思考只允许写阶段摘要和下一步，使用 `meta.reasoningTraceId` 建立关联。
- 不得要求把 `reasoningTrace`、内部推理链、系统提示词或原始执行轨迹放入可见 layout tree。
- 需要解释方案依据时，写成“AI 生成思路摘要”或“决策依据摘要”，并限制为用户可验证事实和业务判断。
