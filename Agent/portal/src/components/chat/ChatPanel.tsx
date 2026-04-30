import { AlertCircle, LoaderCircle, X } from "lucide-react"
import { ChatComposer } from "@/components/chat/components/chat-composer"
import { ChatMessageList } from "@/components/chat/components/chat-message-list"
import { ConversationEmptyState } from "@/components/chat/components/conversation-empty-state"
import { cn } from "@/lib/cn"
import { conversationDraftKeyForConversation, conversationDraftKeyForLocalDraft } from "@/lib/conversation-draft"
import { useConversationStore } from "@/store/useConversationStore"
import { useI18n } from "@/i18n"

interface ChatPanelProps {
  fillHeight?: boolean
  className?: string
}

export function ChatPanel({ fillHeight, className }: ChatPanelProps) {
  const { t } = useI18n()
  const selectedWorkspace = useConversationStore((state) => state.selectedWorkspace)
  const detailsByConversationId = useConversationStore((state) => state.detailsByConversationId)
  const messagesByConversationId = useConversationStore((state) => state.messagesByConversationId)
  const draftsByKey = useConversationStore((state) => state.draftsByKey)
  const pendingTurnByConversationId = useConversationStore((state) => state.pendingTurnByConversationId)
  const pendingTurnStartedAtMsByConversationId = useConversationStore(
    (state) => state.pendingTurnStartedAtMsByConversationId,
  )
  const submittingWorkspaceKey = useConversationStore((state) => state.submittingWorkspaceKey)
  const bootstrapped = useConversationStore((state) => state.bootstrapped)
  const listLoading = useConversationStore((state) => state.listLoading)
  const conversationLoading = useConversationStore((state) => state.conversationLoading)
  const error = useConversationStore((state) => state.error)
  const setDraft = useConversationStore((state) => state.setDraft)
  const submitComposerMessage = useConversationStore((state) => state.submitComposerMessage)
  const submitCandidateSelection = useConversationStore((state) => state.submitCandidateSelection)
  const submitClarificationResponse = useConversationStore((state) => state.submitClarificationResponse)
  const submitLayoutTreeUserMessage = useConversationStore((state) => state.submitLayoutTreeUserMessage)
  const showInteractionError = useConversationStore((state) => state.showInteractionError)
  const clearError = useConversationStore((state) => state.clearError)

  const isLocalDraft = selectedWorkspace?.kind === "localDraft"
  const selectedConversationId = selectedWorkspace?.kind === "conversation" ? selectedWorkspace.conversationId : null
  const selectedDetail = selectedConversationId ? detailsByConversationId[selectedConversationId] : null
  const selectedMessages = selectedConversationId ? messagesByConversationId[selectedConversationId] ?? [] : []
  const pendingTurn = selectedConversationId ? pendingTurnByConversationId[selectedConversationId] : undefined
  const pendingTurnStartedAtMs = selectedConversationId
    ? pendingTurnStartedAtMsByConversationId[selectedConversationId]
    : undefined

  const activeDraftKey =
    selectedWorkspace?.kind === "conversation"
      ? conversationDraftKeyForConversation(selectedWorkspace.conversationId)
      : selectedWorkspace?.kind === "localDraft"
        ? conversationDraftKeyForLocalDraft(selectedWorkspace.localDraftId)
        : null

  const activeDraft = activeDraftKey ? draftsByKey[activeDraftKey] ?? "" : ""
  const isSubmittingCurrentWorkspace =
    selectedWorkspace &&
    submittingWorkspaceKey ===
      (selectedWorkspace.kind === "conversation"
        ? `conversation:${selectedWorkspace.conversationId}`
        : `local:${selectedWorkspace.localDraftId}`)

  const composer = (
    <ChatComposer
      value={activeDraft}
      onChange={setDraft}
      onSubmit={submitComposerMessage}
      centered={isLocalDraft}
      showHint={!isLocalDraft}
      pending={isSubmittingCurrentWorkspace || pendingTurn?.state === "thinking"}
      disabled={conversationLoading}
    />
  )

  return (
    <section
      className={cn(
        "relative flex flex-col bg-transparent",
        fillHeight ? "h-full min-h-0 flex-1 overflow-hidden" : "min-h-[32rem]",
        className,
      )}
    >
      {error ? (
        <div className="pointer-events-none absolute inset-x-0 top-4 z-20 flex justify-center px-4">
          <div className="pointer-events-auto flex w-full max-w-xl items-center gap-3 rounded-xl border border-destructive/20 bg-destructive/[0.05] px-4 py-3 text-sm text-destructive shadow-lg shadow-destructive/10 backdrop-blur-sm">
            <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-destructive/10">
              <AlertCircle className="h-4 w-4" />
            </div>
            <span className="min-w-0 flex-1 leading-6">{error}</span>
            <button
              type="button"
              onClick={clearError}
              aria-label={t("chat.closeHint")}
              className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-destructive/80 transition-colors hover:bg-destructive/10 hover:text-destructive"
            >
              <X className="h-4 w-4" />
            </button>
          </div>
        </div>
      ) : null}

      {!bootstrapped ? (
        <div className="flex min-h-0 flex-1 items-center justify-center px-4">
          <div className="flex items-center gap-3 rounded-xl border border-border bg-card px-4 py-3 text-sm text-muted-foreground shadow-card">
            <LoaderCircle className="h-4 w-4 animate-spin" />
            {t("chat.preparingWorkspace")}
          </div>
        </div>
      ) : isLocalDraft ? (
        <ConversationEmptyState composer={composer} />
      ) : selectedConversationId ? (
        conversationLoading || !selectedDetail ? (
          <div className="flex min-h-0 flex-1 items-center justify-center px-4">
            <div className="flex items-center gap-3 rounded-xl border border-border bg-card px-4 py-3 text-sm text-muted-foreground shadow-card">
              <LoaderCircle className="h-4 w-4 animate-spin" />
              {t("chat.restoringConversation")}
            </div>
          </div>
        ) : (
          <>
            <ChatMessageList
              title={selectedDetail.title}
              summary={selectedDetail.displaySummary}
              showHeader={false}
              messages={selectedMessages}
              submittingSelectionMessageId={
                pendingTurn?.state === "submitting_selection" ? pendingTurn.sourceMessageId : undefined
              }
              isAwaitingReply={pendingTurn?.state === "thinking"}
              awaitingReplyStartedAtMs={pendingTurnStartedAtMs}
              onCandidateSelection={submitCandidateSelection}
              onClarificationResponse={submitClarificationResponse}
              onUserMessageAction={submitLayoutTreeUserMessage}
              onOpenReference={(refId) => {
                showInteractionError(`${t("chat.refNotInMessagePrefix")} ${refId}`)
              }}
            />
            {composer}
          </>
        )
      ) : listLoading ? (
        <div className="flex min-h-0 flex-1 items-center justify-center px-4">
          <div className="flex items-center gap-3 rounded-xl border border-border bg-card px-4 py-3 text-sm text-muted-foreground shadow-card">
            <LoaderCircle className="h-4 w-4 animate-spin" />
            {t("chat.refreshingList")}
          </div>
        </div>
      ) : (
        <div className="flex min-h-0 flex-1 items-center justify-center px-4">
          <div className="rounded-xl border border-dashed border-border bg-card/70 px-5 py-4 text-sm text-muted-foreground">
            {t("chat.conversationNotLoaded")}
          </div>
        </div>
      )}
    </section>
  )
}
