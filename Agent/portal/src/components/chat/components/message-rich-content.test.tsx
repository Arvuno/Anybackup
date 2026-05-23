import { fireEvent, render, screen } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"
import { MessageRichContent } from "@/components/chat/components/message-rich-content"
import type { ConversationMessageSummary } from "@/types/conversation"

describe("MessageRichContent", () => {
  it("renders markdown rich payloads inside the standard markdown document wrapper", () => {
    const message: ConversationMessageSummary = {
      messageId: "msg_markdown_001",
      conversationId: "conv_markdown_001",
      role: "assistant",
      contentType: "rich_content",
      content: "markdown summary",
      createdAt: "2026-04-28T20:00:00.000Z",
      status: "responded",
      richPayload: {
        kind: "markdown",
        data: {
          markdown:
            "# 备份方案推荐\n\n| 属性 | 值 |\n| --- | --- |\n| 生产资源名称 | **mysql3306_U0HYTDM3RENXS4EJ** |",
        },
      },
    }

    const { container } = render(
      <MessageRichContent
        message={message}
        submitting={false}
        onCandidateSelection={vi.fn()}
      />,
    )

    expect(container.querySelector(".markdown-body")).toBeInTheDocument()
    expect(screen.getByRole("heading", { level: 1, name: "备份方案推荐" })).toBeInTheDocument()
    expect(screen.getByRole("table")).toBeInTheDocument()
    expect(screen.getByText("mysql3306_U0HYTDM3RENXS4EJ", { selector: "strong" })).toBeInTheDocument()
    expect(screen.queryByText("**mysql3306_U0HYTDM3RENXS4EJ**")).not.toBeInTheDocument()
  })

  it("passes camelCase layout-tree clarification payload fields to clarification responses", () => {
    const onClarificationResponse = vi.fn()
    const message: ConversationMessageSummary = {
      messageId: "173612696141303809",
      conversationId: "173403886646726656",
      role: "assistant",
      contentType: "clarification",
      content: "Please confirm the target recovery time window.",
      createdAt: "2026-04-25T01:54:53.660000Z",
      status: "responded",
      richPayload: {
        kind: "layout_tree",
        data: {
          activity: {
            contract: "conversation.ui.layout-tree@1",
            blockId: "clarification_946",
            ui: {
              type: "stack",
              children: [
                {
                  type: "action-row",
                  props: {
                    actionIds: ["submit_latest_safe_point"],
                  },
                },
              ],
            },
            meta: {
              intent: "clarification",
              terminal: false,
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
            ],
          },
        },
      },
    }

    render(
      <MessageRichContent
        message={message}
        submitting={false}
        onCandidateSelection={vi.fn()}
        onClarificationResponse={onClarificationResponse}
      />,
    )

    fireEvent.click(screen.getByRole("button", { name: "Use latest safe point" }))

    expect(onClarificationResponse).toHaveBeenCalledWith({
      type: "clarification_response",
      messageId: "173612696141303809",
      clarificationId: "recovery_window",
      selectedValue: "latest_safe_point",
      freeText: undefined,
    })
  })
})
