import ReactMarkdown from "react-markdown"
import remarkGfm from "remark-gfm"
import { AgUiLayoutTreeDemo } from "@/components/chat/components/ag-ui-layout-tree-demo"
import { ClarificationCard } from "@/components/chat/components/clarification-card"
import { CandidateOptionsCard } from "@/components/chat/components/candidate-options-card"
import type {
  CandidateSelectionInput,
  ClarificationResponseInput,
  ConversationMessageSummary,
  ConversationRichPayload,
  UserMessageInput,
} from "@/types/conversation"

interface MessageRichContentProps {
  message: ConversationMessageSummary
  submitting: boolean
  onCandidateSelection: (input: CandidateSelectionInput) => void
  onClarificationResponse?: (input: ClarificationResponseInput) => void
  onUserMessageAction?: (input: UserMessageInput) => void
  onOpenReference?: (refId: string) => void
}

type RenderableRichMessage = ConversationMessageSummary & {
  richPayload: Extract<ConversationRichPayload, { kind: "markdown" | "candidate_options" | "clarification" | "layout_tree" }>
}

export function canRenderRichContent(message: ConversationMessageSummary): message is RenderableRichMessage {
  return (
    message.richPayload?.kind === "markdown" ||
    message.richPayload?.kind === "candidate_options" ||
    message.richPayload?.kind === "clarification" ||
    message.richPayload?.kind === "layout_tree"
  )
}

function MarkdownContent({ text }: { text: string }) {
  return (
    <div className="markdown-rich-content overflow-x-auto rounded-xl border border-border/60 bg-background/80 px-4 py-3">
      <div className="markdown-body">
        <ReactMarkdown
          remarkPlugins={[remarkGfm]}
          skipHtml
          components={{
            a: ({ children, href }) => (
              <a href={href} target="_blank" rel="noreferrer">
                {children}
              </a>
            ),
          }}
        >
          {text}
        </ReactMarkdown>
      </div>
    </div>
  )
}

export function MessageRichContent({
  message,
  submitting,
  onCandidateSelection,
  onClarificationResponse,
  onUserMessageAction,
  onOpenReference,
}: MessageRichContentProps) {
  if (canRenderRichContent(message)) {
    const richPayload = message.richPayload

    if (richPayload.kind === "markdown") {
      return <MarkdownContent text={richPayload.data.markdown} />
    }

    if (richPayload.kind === "layout_tree") {
      return (
        <AgUiLayoutTreeDemo
          activity={richPayload.data.activity}
          stateSnapshot={richPayload.data.stateSnapshot}
          showMeta={false}
          onSubmitMessage={(actionPayload, actionEvent) => {
            const type =
              typeof actionPayload.type === "string"
                ? actionPayload.type
                : typeof actionPayload.action === "string" &&
                    (typeof actionPayload.candidate_option_id === "string" ||
                      typeof actionPayload.candidateOptionId === "string")
                  ? "candidate_selection"
                  : undefined

            if (type === "candidate_selection") {
              const reasoningTraceId =
                (typeof actionPayload.reasoning_trace_id === "string" ? actionPayload.reasoning_trace_id : undefined) ??
                (typeof actionPayload.reasoningTraceId === "string" ? actionPayload.reasoningTraceId : undefined) ??
                richPayload.data.activity.meta?.reasoningTraceId
              const candidateOptionId =
                (typeof actionPayload.candidate_option_id === "string"
                  ? actionPayload.candidate_option_id
                  : undefined) ??
                (typeof actionPayload.candidateOptionId === "string" ? actionPayload.candidateOptionId : undefined)
              const selection =
                actionPayload.selection === "confirm" ||
                actionPayload.selection === "reject" ||
                actionPayload.selection === "revise"
                  ? actionPayload.selection
                  : actionPayload.action === "confirm" ||
                      actionPayload.action === "reject" ||
                      actionPayload.action === "revise"
                    ? actionPayload.action
                  : undefined

              if (!reasoningTraceId || !candidateOptionId || !selection) return

              onCandidateSelection({
                type: "candidate_selection",
                messageId: message.messageId,
                reasoningTraceId,
                candidateOptionId,
                selection,
                additionalConstraints:
                  typeof actionPayload.additional_constraints === "string"
                    ? actionPayload.additional_constraints
                    : undefined,
              })
              return
            }

            if (type === "clarification_response") {
              onClarificationResponse?.({
                type: "clarification_response",
                messageId: message.messageId,
                clarificationId:
                  (typeof actionPayload.clarification_id === "string" ? actionPayload.clarification_id : undefined) ??
                  (typeof actionPayload.clarificationId === "string" ? actionPayload.clarificationId : undefined),
                selectedValue:
                  (typeof actionPayload.selected_value === "string" ? actionPayload.selected_value : undefined) ??
                  (typeof actionPayload.selectedValue === "string" ? actionPayload.selectedValue : undefined),
                freeText:
                  (typeof actionPayload.free_text === "string" ? actionPayload.free_text : undefined) ??
                  (typeof actionPayload.freeText === "string" ? actionPayload.freeText : undefined),
              })
              return
            }

            const fallbackContent =
              typeof actionPayload.content === "string"
                ? actionPayload.content
                : typeof actionPayload.text === "string"
                  ? actionPayload.text
                  : actionEvent.label

            if (!fallbackContent) return

            onUserMessageAction?.({
              type: "user_message",
              content: fallbackContent,
            })
          }}
          onAction={(event) => {
            if (event.kind !== "open_ref") return

            const refId =
              typeof event.payload?.ref_id === "string"
                ? event.payload.ref_id
                : typeof event.payload?.refId === "string"
                  ? event.payload.refId
                  : undefined

            if (!refId) return
            onOpenReference?.(refId)
          }}
        />
      )
    }

    if (richPayload.kind === "clarification") {
      return (
        <ClarificationCard
          payload={richPayload.data}
          submitting={submitting}
          onSubmit={(response) =>
            onClarificationResponse?.({
              type: "clarification_response",
              messageId: message.messageId,
              clarificationId: response.clarificationId,
              selectedValue: response.selectedValue,
              freeText: response.freeText,
            })
          }
        />
      )
    }

    if (richPayload.kind === "candidate_options") {
      return (
        <CandidateOptionsCard
          payload={richPayload.data}
          submitting={submitting}
          onSelect={(selection) =>
            onCandidateSelection({
              type: "candidate_selection",
              messageId: message.messageId,
              reasoningTraceId: selection.reasoningTraceId,
              candidateOptionId: selection.candidateOptionId,
              selection: selection.selection,
              additionalConstraints: selection.additionalConstraints,
            })
          }
        />
      )
    }
  }

  return <p className="whitespace-pre-wrap break-words">{message.content}</p>
}
