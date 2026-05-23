import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"
import { ChatMessageList } from "@/components/chat/components/chat-message-list"
import { WaitingMessage } from "@/components/chat/components/waiting-message"
import type { ConversationMessageSummary } from "@/types/conversation"

const plainMessages: ConversationMessageSummary[] = [
  {
    messageId: "plain-user-001",
    conversationId: "demo-plain",
    role: "user",
    contentType: "text",
    content: "Start a restore check.",
    createdAt: "2026-04-22T16:00:00+08:00",
    status: "responded",
  },
  {
    messageId: "plain-assistant-001",
    conversationId: "demo-plain",
    role: "assistant",
    contentType: "text",
    content: "I will verify the latest restore points first.",
    createdAt: "2026-04-22T16:00:20+08:00",
    status: "responded",
  },
]

const richMessages: ConversationMessageSummary[] = [
  {
    messageId: "rich-user-001",
    conversationId: "demo-rich",
    role: "user",
    contentType: "text",
    content: "Show me the restore options.",
    createdAt: "2026-04-22T16:00:00+08:00",
    status: "responded",
  },
  {
    messageId: "rich-assistant-001",
    conversationId: "demo-rich",
    role: "assistant",
    contentType: "rich_content",
    content: "I found three restore options.",
    createdAt: "2026-04-22T16:00:20+08:00",
    status: "responded",
    richPayload: {
      kind: "candidate_options",
      data: {
        reasoningTraceId: "trace-001",
        title: "Restore options",
        summary: "Three structured restore options are available.",
        actions: [
          { type: "confirm", label: "Confirm restore plan" },
          { type: "reject", label: "Reject plan" },
          {
            type: "revise",
            label: "Add constraints",
            inputLabel: "Add constraints",
            inputPlaceholder: "For example: generate the plan only.",
            submitLabel: "Submit constraints",
          },
        ],
        options: [
          {
            optionId: "option_a",
            title: "Option A: Restore by exporting the target table",
            recommended: true,
            summary: "Recommended option with the lowest rollback impact.",
            fields: [
              { key: "scope", label: "Restore scope", value: "Database-level" },
              { key: "rpo_rto", label: "RPO / RTO", value: "< 2 min / 1.5 h" },
            ],
            extra: {
              title: "Recommendation",
              content: "This path keeps the production database isolated while recovering the target table.",
            },
          },
        ],
      },
    },
  },
]

const clarificationRichMessages: ConversationMessageSummary[] = [
  {
    messageId: "clarify-user-001",
    conversationId: "demo-clarify",
    role: "user",
    contentType: "text",
    content: "Please clarify the restore window.",
    createdAt: "2026-04-22T16:05:00+08:00",
    status: "responded",
  },
  {
    messageId: "clarify-assistant-001",
    conversationId: "demo-clarify",
    role: "assistant",
    contentType: "clarification",
    content: "Please confirm the restore window.",
    createdAt: "2026-04-22T16:05:20+08:00",
    status: "responded",
    richPayload: {
      kind: "clarification",
      data: {
        clarificationId: "restore_window",
        prompt: "Please confirm the restore window.",
        options: [
          { label: "Latest safe point", value: "latest_safe_point" },
          { label: "Custom timestamp", value: "custom_timestamp" },
        ],
        inputConstraints: {
          required: true,
          allowFreeText: true,
          freeTextLabel: "Restore time",
          freeTextPlaceholder: "Enter a timestamp",
        },
      },
    },
  },
]

const layoutTreeMessages: ConversationMessageSummary[] = [
  {
    messageId: "layout-user-001",
    conversationId: "demo-layout-tree",
    role: "user",
    contentType: "text",
    content: "Show me the protocol-driven restore options.",
    createdAt: "2026-04-24T10:00:00+08:00",
    status: "responded",
  },
  {
    messageId: "layout-assistant-001",
    conversationId: "demo-layout-tree",
    role: "assistant",
    contentType: "rich_content",
    content: "I generated a layout-tree response.",
    createdAt: "2026-04-24T10:00:20+08:00",
    status: "responded",
    richPayload: {
      kind: "layout_tree",
      data: {
        activity: {
          contract: "conversation.ui.layout-tree@1",
          blockId: "candidate_compare_001",
          ui: {
            id: "root",
            type: "grid",
            props: {
              columns: 2,
              gap: "md",
            },
            children: [
              {
                id: "candidate_a",
                type: "card",
                props: {
                  gap: "md",
                },
                children: [
                  {
                    type: "heading",
                    props: {
                      level: 3,
                      text: "Plan A",
                    },
                  },
                  {
                    type: "action-row",
                    props: {
                      actionIds: ["choose_candidate_a"],
                    },
                  },
                ],
              },
              {
                id: "candidate_b",
                type: "card",
                props: {
                  gap: "md",
                },
                children: [
                  {
                    type: "heading",
                    props: {
                      level: 3,
                      text: "Plan B",
                    },
                  },
                  {
                    type: "action-row",
                    props: {
                      actionIds: ["choose_candidate_b"],
                    },
                  },
                ],
              },
            ],
          },
          actions: [
            {
              id: "choose_candidate_a",
              kind: "submit_message",
              label: "Confirm Plan A",
              payload: {
                type: "candidate_selection",
                candidate_option_id: "candidate_a",
                selection: "confirm",
              },
            },
            {
              id: "choose_candidate_b",
              kind: "submit_message",
              label: "Confirm Plan B",
              payload: {
                type: "candidate_selection",
                candidate_option_id: "candidate_b",
                selection: "confirm",
              },
            },
          ],
          meta: {
            intent: "result",
            reasoningTraceId: "trace-layout-001",
          },
        },
        stateSnapshot: {
          interaction: {
            status: "clarifying",
          },
          selection: {
            required: true,
            selectedCandidateOptionId: "candidate_a",
            selectionLocked: false,
          },
        },
      },
    } as unknown as ConversationMessageSummary["richPayload"],
  },
]

const compactLayoutTreeMessages: ConversationMessageSummary[] = [
  {
    messageId: "compact-layout-user-001",
    conversationId: "demo-compact-layout-tree",
    role: "user",
    contentType: "text",
    content: "Summarize the restore context.",
    createdAt: "2026-04-24T10:10:00+08:00",
    status: "responded",
  },
  {
    messageId: "compact-layout-assistant-001",
    conversationId: "demo-compact-layout-tree",
    role: "assistant",
    contentType: "rich_content",
    content: "I summarized the current restore context.",
    createdAt: "2026-04-24T10:10:20+08:00",
    status: "responded",
    richPayload: {
      kind: "layout_tree",
      data: {
        activity: {
          contract: "conversation.ui.layout-tree@1",
          blockId: "restore_summary_001",
          ui: {
            id: "root",
            type: "stack",
            props: {
              gap: "lg",
            },
            children: [
              {
                type: "heading",
                props: {
                  level: 2,
                  text: "Restore analysis",
                },
              },
              {
                type: "paragraph",
                props: {
                  text: "The target database and recovery window have been parsed from the request.",
                },
              },
              {
                type: "card",
                props: {
                  title: "Current understanding",
                },
                children: [
                  {
                    type: "kv-list",
                    props: {
                      items: [
                        { label: "Target asset", value: "orders-db" },
                        { label: "Requested window", value: "Last Friday afternoon" },
                      ],
                    },
                  },
                  {
                    type: "callout",
                    props: {
                      title: "Next action",
                      tone: "info",
                      text: "Search restore points and compare execution paths.",
                    },
                  },
                ],
              },
            ],
          },
          meta: {
            intent: "thought",
            terminal: false,
          },
        },
        stateSnapshot: {
          interaction: {
            status: "thinking",
          },
        },
      },
    } as unknown as ConversationMessageSummary["richPayload"],
  },
]

const clarificationLayoutTreeMessages: ConversationMessageSummary[] = [
  {
    messageId: "clarification-layout-user-001",
    conversationId: "demo-clarification-layout-tree",
    role: "user",
    contentType: "text",
    content: "Clarify restore window.",
    createdAt: "2026-04-24T10:15:00+08:00",
    status: "responded",
  },
  {
    messageId: "clarification-layout-assistant-001",
    conversationId: "demo-clarification-layout-tree",
    role: "assistant",
    contentType: "clarification",
    content: "Please confirm the target recovery time window.",
    createdAt: "2026-04-24T10:15:20+08:00",
    status: "responded",
    richPayload: {
      kind: "layout_tree",
      data: {
        activity: {
          contract: "conversation.ui.layout-tree@1",
          blockId: "clarification_001",
          ui: {
            id: "root",
            type: "stack",
            props: {
              gap: "lg",
            },
            children: [
              {
                type: "heading",
                props: {
                  level: 2,
                  text: "Clarification required",
                },
              },
              {
                type: "paragraph",
                props: {
                  text: "The request can proceed after the recovery window is explicitly confirmed.",
                },
              },
              {
                type: "action-row",
                props: {
                  actionIds: ["submit_latest_safe_point", "submit_specific_timestamp"],
                },
              },
            ],
          },
          actions: [
            {
              id: "submit_latest_safe_point",
              kind: "submit_message",
              label: "Use latest safe point",
              style: "primary",
              payload: {
                type: "clarification_response",
                selectedValue: "latest_safe_point",
                clarificationId: "recovery_window",
              },
            },
            {
              id: "submit_specific_timestamp",
              kind: "submit_message",
              label: "Provide specific timestamp",
              payload: {
                type: "clarification_response",
                freeText: "2026-04-24T15:30:00+08:00",
                clarificationId: "recovery_window",
              },
            },
          ],
          meta: {
            intent: "clarification",
            terminal: false,
          },
        },
        stateSnapshot: {
          interaction: {
            status: "clarifying",
          },
          selection: {
            required: true,
            selectedCandidateOptionId: null,
            selectionLocked: false,
          },
        },
      },
    } as unknown as ConversationMessageSummary["richPayload"],
  },
]

const unsupportedRichMessages: ConversationMessageSummary[] = [
  {
    messageId: "unsupported-rich-001",
    conversationId: "demo-unsupported-rich",
    role: "assistant",
    contentType: "rich_content",
    content: "Fallback text for an unsupported rich payload.",
    createdAt: "2026-04-22T16:01:00+08:00",
    status: "responded",
    richPayload: {
      kind: "unsupported_card",
      data: {
        title: "Unknown rich block",
      },
    } as unknown as ConversationMessageSummary["richPayload"],
  },
]

const turnAggregatedMessages: ConversationMessageSummary[] = [
  {
    messageId: "turn-user-001",
    conversationId: "demo-turn",
    turnId: "turn-001",
    role: "user",
    contentType: "text",
    content: "测试对话",
    createdAt: "2026-04-23T20:57:35+08:00",
    status: "persisted",
  },
  {
    messageId: "turn-thought-001",
    conversationId: "demo-turn",
    turnId: "turn-001",
    parentMessageId: "turn-user-001",
    role: "assistant",
    contentType: "rich_content",
    content: "正在理解问题并准备查询上下文。",
    createdAt: "2026-04-23T20:57:37+08:00",
    status: "responded",
    agUiSequence: 1,
    agUiEventName: "anybackup.conversation.thought.render",
    richPayload: {
      kind: "thought",
      data: {
        status: "running",
        summary: "正在提取用户目标和上下文线索。",
      },
    },
  },
  {
    messageId: "turn-result-001",
    conversationId: "demo-turn",
    turnId: "turn-001",
    parentMessageId: "turn-user-001",
    role: "assistant",
    contentType: "rich_content",
    content: "已完成会话的模拟分析，生成结果批次 mock-940。",
    createdAt: "2026-04-23T20:57:38+08:00",
    status: "responded",
    agUiSequence: 2,
    agUiEventName: "anybackup.conversation.result.render",
    richPayload: {
      kind: "result",
      data: {
        title: "模拟分析结果",
        summary: "已完成会话的模拟分析，生成结果批次 mock-940。",
        resultId: "mock-940",
        sourceMessage: "测试对话",
      },
    },
  },
  {
    messageId: "turn-chart-001",
    conversationId: "demo-turn",
    turnId: "turn-001",
    parentMessageId: "turn-user-001",
    role: "assistant",
    contentType: "rich_content",
    content: "已生成一份可视化分析。",
    createdAt: "2026-04-23T20:57:39+08:00",
    status: "responded",
    agUiSequence: 3,
    agUiEventName: "anybackup.conversation.chart.render",
    richPayload: {
      kind: "ag_ui",
      data: {
        eventName: "anybackup.conversation.chart.render",
        title: "趋势图",
        summary: "已生成一份可视化分析。",
      },
    },
  },
]

const mixedWidthAssistantTurnMessages: ConversationMessageSummary[] = [
  {
    messageId: "mixed-user-001",
    conversationId: "demo-mixed-width",
    turnId: "mixed-turn-001",
    role: "user",
    contentType: "text",
    content: "Assess restore options.",
    createdAt: "2026-04-24T10:20:00+08:00",
    status: "responded",
  },
  {
    messageId: "mixed-layout-summary-001",
    conversationId: "demo-mixed-width",
    turnId: "mixed-turn-001",
    parentMessageId: "mixed-user-001",
    role: "assistant",
    contentType: "rich_content",
    content: "I summarized the restore context.",
    createdAt: "2026-04-24T10:20:10+08:00",
    status: "responded",
    agUiSequence: 1,
    richPayload: {
      kind: "layout_tree",
      data: {
        activity: {
          contract: "conversation.ui.layout-tree@1",
          blockId: "restore_summary_002",
          ui: {
            id: "summary-root",
            type: "stack",
            children: [
              {
                type: "heading",
                props: {
                  level: 2,
                  text: "Restore-point search",
                },
              },
              {
                type: "paragraph",
                props: {
                  text: "The tool returned one recommended restore point and two fallback alternatives.",
                },
              },
            ],
          },
          meta: {
            intent: "tool_call",
            terminal: false,
          },
        },
        stateSnapshot: {
          interaction: {
            status: "executing",
          },
        },
      },
    } as unknown as ConversationMessageSummary["richPayload"],
  },
  {
    messageId: "mixed-layout-candidate-001",
    conversationId: "demo-mixed-width",
    turnId: "mixed-turn-001",
    parentMessageId: "mixed-user-001",
    role: "assistant",
    contentType: "rich_content",
    content: "Generated three recovery candidates.",
    createdAt: "2026-04-24T10:20:20+08:00",
    status: "responded",
    agUiSequence: 2,
    richPayload: {
      kind: "layout_tree",
      data: {
        activity: {
          contract: "conversation.ui.layout-tree@1",
          blockId: "candidate_compare_002",
          ui: {
            id: "candidate-root",
            type: "grid",
            props: {
              columns: 2,
              gap: "md",
            },
            children: [
              {
                id: "candidate_a",
                type: "card",
                children: [
                  {
                    type: "heading",
                    props: {
                      level: 3,
                      text: "Plan A",
                    },
                  },
                  {
                    type: "action-row",
                    props: {
                      actionIds: ["choose_candidate_a"],
                    },
                  },
                ],
              },
            ],
          },
          actions: [
            {
              id: "choose_candidate_a",
              kind: "submit_message",
              label: "Confirm Plan A",
              payload: {
                type: "candidate_selection",
                candidate_option_id: "candidate_a",
                selection: "confirm",
              },
            },
          ],
          meta: {
            intent: "result",
            reasoningTraceId: "trace-mixed-001",
          },
        },
        stateSnapshot: {
          interaction: {
            status: "clarifying",
          },
          selection: {
            required: true,
            selectedCandidateOptionId: "candidate_a",
            selectionLocked: false,
          },
        },
      },
    } as unknown as ConversationMessageSummary["richPayload"],
  },
]

describe("ChatMessageList", () => {
  it("renders assistant text messages with the standard AI bubble layout", () => {
    render(
      <ChatMessageList
        title="Plain text loop"
        summary="Assistant text response"
        messages={plainMessages}
        onCandidateSelection={vi.fn()}
      />,
    )

    const assistantText = screen.getByText("I will verify the latest restore points first.")
    const wrapper = assistantText.closest("[data-message-width]")
    const bubble = wrapper?.firstElementChild
    const userText = screen.getByText("Start a restore check.")
    const userArticle = userText.closest("article")
    const userWrapper = userArticle?.firstElementChild
    const userBubble = userWrapper?.firstElementChild

    expect(wrapper).toHaveClass("w-full")
    expect(wrapper).toHaveClass("max-w-3xl")
    expect(wrapper).not.toHaveClass("max-w-[min(100%,72rem)]")
    expect(bubble).toHaveClass("bubble-ai")
    expect(bubble).toHaveClass("max-w-[85%]")
    expect(userWrapper).toHaveClass("w-full")
    expect(userWrapper).toHaveClass("max-w-3xl")
    expect(userBubble).toHaveClass("bubble-user")
    expect(userBubble).toHaveClass("max-w-[85%]")
  })

  it("renders supported rich content messages with a wider layout", () => {
    render(
      <ChatMessageList
        title="Restore options"
        summary="Candidate cards"
        messages={richMessages}
        onCandidateSelection={vi.fn()}
      />,
    )

    const richHeading = screen.getByRole("heading", { name: "Restore options", level: 3 })
    const wrapper = richHeading.closest("[data-message-width]")
    const bubble = wrapper?.firstElementChild

    expect(wrapper).toHaveClass("w-full")
    expect(wrapper).toHaveClass("max-w-[min(100%,72rem)]")
    expect(wrapper).not.toHaveClass("max-w-[78%]")
    expect(bubble).toHaveClass("bg-transparent")
  })

  it("renders clarification cards inside the standard readable column", () => {
    render(
      <ChatMessageList
        title="Restore clarification"
        summary="Clarification cards"
        messages={clarificationRichMessages}
        onCandidateSelection={vi.fn()}
      />,
    )

    const prompt = screen.getByText("Please confirm the restore window.")
    const wrapper = prompt.closest("[data-message-width]")
    const bubble = wrapper?.firstElementChild

    expect(wrapper).toHaveClass("w-full")
    expect(wrapper).toHaveClass("max-w-3xl")
    expect(wrapper).not.toHaveClass("max-w-[min(100%,72rem)]")
    expect(bubble).toHaveClass("bg-transparent")
    expect(screen.getByRole("button", { name: "Latest safe point" })).toBeInTheDocument()
  })

  it("renders candidate compare layout-tree messages with the wider layout and routes submit actions", async () => {
    const user = userEvent.setup()
    const onCandidateSelection = vi.fn()

    render(
      <ChatMessageList
        title="Layout tree"
        summary="Protocol-driven cards"
        messages={layoutTreeMessages}
        onCandidateSelection={onCandidateSelection}
      />,
    )

    const planHeading = screen.getByRole("heading", { name: "Plan A", level: 3 })
    const wrapper = planHeading.closest("[data-message-width]")

    expect(wrapper).toHaveClass("w-full")
    expect(wrapper).toHaveClass("max-w-[min(100%,72rem)]")
    expect(screen.getByRole("button", { name: "Confirm Plan B" })).toBeInTheDocument()

    await user.click(screen.getByText("Plan B"))
    await user.click(screen.getByRole("button", { name: "Confirm Plan B" }))

    expect(onCandidateSelection).toHaveBeenCalledWith({
      type: "candidate_selection",
      messageId: "layout-assistant-001",
      reasoningTraceId: "trace-layout-001",
      candidateOptionId: "candidate_b",
      selection: "confirm",
      additionalConstraints: undefined,
    })
  })

  it("renders lightweight layout-tree panels inside the standard readable column", () => {
    render(
      <ChatMessageList
        title="Compact layout tree"
        summary="Lightweight structured panel"
        messages={compactLayoutTreeMessages}
        onCandidateSelection={vi.fn()}
      />,
    )

    const heading = screen.getByRole("heading", { name: "Restore analysis", level: 2 })
    const wrapper = heading.closest("[data-message-width]")
    const bubble = wrapper?.firstElementChild

    expect(wrapper).toHaveClass("w-full")
    expect(wrapper).toHaveClass("max-w-3xl")
    expect(wrapper).not.toHaveClass("max-w-[min(100%,72rem)]")
    expect(bubble).toHaveClass("bg-transparent")
  })

  it("keeps clarification layout-tree panels inside the standard readable column", () => {
    render(
      <ChatMessageList
        title="Clarification layout tree"
        summary="Clarification panel"
        messages={clarificationLayoutTreeMessages}
        onCandidateSelection={vi.fn()}
      />,
    )

    const heading = screen.getByRole("heading", { name: "Clarification required", level: 2 })
    const wrapper = heading.closest("[data-message-width]")
    const bubble = wrapper?.firstElementChild

    expect(wrapper).toHaveAttribute("data-message-width", "standard")
    expect(wrapper).toHaveClass("max-w-3xl")
    expect(wrapper).not.toHaveClass("max-w-[min(100%,72rem)]")
    expect(bubble).toHaveClass("bg-transparent")
  })

  it("keeps mixed assistant messages in the same turn on their own width tracks", () => {
    render(
      <ChatMessageList
        title="Mixed assistant widths"
        summary="Same turn, different message widths"
        messages={mixedWidthAssistantTurnMessages}
        onCandidateSelection={vi.fn()}
      />,
    )

    const summaryHeading = screen.getByRole("heading", { name: "Restore-point search", level: 2 })
    const summaryArticle = summaryHeading.closest("article")
    const summaryWrapper = summaryHeading.closest("[data-message-width]")
    const candidateHeading = screen.getByRole("heading", { name: "Plan A", level: 3 })
    const candidateArticle = candidateHeading.closest("article")
    const candidateWrapper = candidateHeading.closest("[data-message-width]")

    expect(summaryArticle).toBe(candidateArticle)
    expect(summaryWrapper).toHaveAttribute("data-message-width", "standard")
    expect(summaryWrapper).toHaveClass("max-w-3xl")
    expect(candidateWrapper).toHaveAttribute("data-message-width", "expanded")
    expect(candidateWrapper).toHaveClass("max-w-[min(100%,72rem)]")
  })

  it("falls back to the standard AI bubble for unsupported rich payloads", () => {
    render(
      <ChatMessageList
        title="Unsupported rich payload"
        summary="Fallback rendering"
        messages={unsupportedRichMessages}
        onCandidateSelection={vi.fn()}
      />,
    )

    const fallbackText = screen.getByText("Fallback text for an unsupported rich payload.")
    const wrapper = fallbackText.closest("[data-message-width]")
    const bubble = wrapper?.firstElementChild

    expect(wrapper).toHaveClass("w-full")
    expect(wrapper).toHaveClass("max-w-3xl")
    expect(wrapper).not.toHaveClass("max-w-[min(100%,72rem)]")
    expect(bubble).toHaveClass("bubble-ai")
    expect(bubble).not.toHaveClass("bg-transparent")
    expect(bubble).toHaveClass("max-w-[85%]")
  })

  it("aggregates thought, result, and future AG-UI messages into one assistant turn", () => {
    const { container } = render(
      <ChatMessageList
        title="Turn aggregation"
        messages={turnAggregatedMessages}
        onCandidateSelection={vi.fn()}
      />,
    )

    const articles = container.querySelectorAll("article")

    expect(articles).toHaveLength(2)
    expect(articles[0]).toHaveTextContent("测试对话")
    expect(articles[1]).toHaveTextContent("Thought")
    expect(articles[1]).toHaveTextContent("Done")
    expect(articles[1]).toHaveTextContent("正在提取用户目标和上下文线索。")
    expect(articles[1]).not.toHaveTextContent("模拟分析结果")
    expect(articles[1]).toHaveTextContent("已完成会话的模拟分析，生成结果批次 mock-940。")
    expect(articles[1]).toHaveTextContent("趋势图")
    expect(articles[1]).toHaveTextContent("已生成一份可视化分析。")
  })

  it("renders thought as a collapsible preface and result as plain assistant body", async () => {
    const user = userEvent.setup()

    render(
      <ChatMessageList
        title="Thought preface"
        messages={[
          {
            messageId: "preface-user-001",
            conversationId: "demo-preface",
            turnId: "turn-preface",
            role: "user",
            contentType: "text",
            content: "Plan the backup.",
            createdAt: "2026-04-23T21:30:00+08:00",
            status: "responded",
          },
          {
            messageId: "preface-thought-001",
            conversationId: "demo-preface",
            turnId: "turn-preface",
            parentMessageId: "preface-user-001",
            role: "assistant",
            contentType: "rich_content",
            content: "Collecting backup context.",
            createdAt: "2026-04-23T21:30:02+08:00",
            status: "responded",
            agUiSequence: 1,
            richPayload: {
              kind: "thought",
              data: {
                status: "completed",
                summary: "Collecting backup context.",
              },
            },
          },
          {
            messageId: "preface-result-001",
            conversationId: "demo-preface",
            turnId: "turn-preface",
            parentMessageId: "preface-user-001",
            role: "assistant",
            contentType: "rich_content",
            content: "Use replica backups for lower primary impact.",
            createdAt: "2026-04-23T21:30:03+08:00",
            status: "responded",
            agUiSequence: 2,
            richPayload: {
              kind: "result",
              data: {
                title: "Mock analysis result",
                summary: "Use replica backups for lower primary impact.",
              },
            },
          },
        ]}
        onCandidateSelection={vi.fn()}
      />,
    )

    const thoughtToggle = screen.getByRole("button", { name: /Thought Done/ })
    expect(thoughtToggle).toHaveAttribute("aria-expanded", "true")
    expect(screen.getByText("Collecting backup context.")).toBeInTheDocument()

    const resultText = screen.getByText("Use replica backups for lower primary impact.")
    expect(resultText.closest(".bubble-ai")).toBeInTheDocument()
    expect(screen.queryByText("Mock analysis result")).not.toBeInTheDocument()

    await user.click(thoughtToggle)

    expect(thoughtToggle).toHaveAttribute("aria-expanded", "false")
    expect(screen.queryByText("Collecting backup context.")).not.toBeInTheDocument()
    expect(screen.getByText("Use replica backups for lower primary impact.")).toBeInTheDocument()
  })

  it("does not crash when a message timestamp is invalid", () => {
    render(
      <ChatMessageList
        title="Invalid timestamp"
        messages={[
          {
            messageId: "invalid-time-001",
            conversationId: "demo-invalid-time",
            role: "assistant",
            contentType: "text",
            content: "This message has a malformed timestamp.",
            createdAt: "",
            status: "responded",
          },
        ]}
        onCandidateSelection={vi.fn()}
      />,
    )

    expect(screen.getByText("This message has a malformed timestamp.")).toBeInTheDocument()
    expect(screen.getByText("Just now")).toBeInTheDocument()
  })
  it("uses a user-facing thinking label for the waiting placeholder", () => {
    render(<WaitingMessage />)

    expect(screen.getByText("Thinking")).toBeInTheDocument()
    expect(screen.queryByText("正在处理这轮对话...")).not.toBeInTheDocument()
  })

  it("shows elapsed thinking time for the waiting placeholder", () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date("2026-04-26T02:03:20.000Z"))

    try {
      render(<WaitingMessage startedAtMs={new Date("2026-04-26T02:00:00.000Z").getTime()} />)

      expect(screen.getByText("Thought for 3m20s")).toBeInTheDocument()
      expect(screen.queryByText("Thinking")).not.toBeInTheDocument()
    } finally {
      vi.useRealTimers()
    }
  })
})
