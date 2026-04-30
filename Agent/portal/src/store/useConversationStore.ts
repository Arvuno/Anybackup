import { create } from "zustand"
import { createId } from "@/lib/ids"
import { translate } from "@/i18n/messages"
import {
  conversationDraftKeyForConversation,
  conversationDraftKeyForLocalDraft,
  readConversationDraft,
  readConversationWorkspaceState,
  removeConversationDraft,
  writeConversationDraft,
  writeConversationWorkspaceState,
} from "@/lib/conversation-draft"
import {
  createConversation,
  getConversationDetail,
  getConversationMessages,
  listConversationEvents,
  listConversations,
  searchConversations,
  sendMessage,
} from "@/services/conversation-service"
import { ServiceError } from "@/types/auth"
import type {
  CandidateSelectionInput,
  ClarificationResponseInput,
  ConversationDetail,
  ConversationMessageSummary,
  ConversationStatusEvent,
  ConversationSummary,
  ConversationWorkspaceSelection,
  ConversationWorkspaceState,
  LocalDraftWorkspace,
  MessageAcceptedResult,
  UserMessageInput,
} from "@/types/conversation"

interface ConversationState {
  bootstrapped: boolean
  listLoading: boolean
  conversationLoading: boolean
  query: string
  error: string | null
  conversations: ConversationSummary[]
  selectedWorkspace: ConversationWorkspaceSelection | null
  localDraftWorkspace: LocalDraftWorkspace | null
  detailsByConversationId: Record<string, ConversationDetail>
  messagesByConversationId: Record<string, ConversationMessageSummary[]>
  draftsByKey: Record<string, string>
  pendingTurnByConversationId: Record<string, PendingTurnState>
  pendingTurnStartedAtMsByConversationId: Record<string, number>
  nextPollAfterMsByConversationId: Record<string, number>
  eventCursorByConversationId: Record<string, string | null>
  latestEventSequenceByConversationId: Record<string, number>
  appliedStatusEventIdsByConversationId: Record<string, Record<string, true>>
  submittingWorkspaceKey: string | null
}

interface ConversationActions {
  hydrate: () => Promise<void>
  activateLocalDraftWorkspace: () => void
  selectConversation: (conversationId: string) => Promise<void>
  setSearchQuery: (query: string) => Promise<void>
  setDraft: (value: string) => void
  submitComposerMessage: () => Promise<void>
  submitCandidateSelection: (input: CandidateSelectionInput) => Promise<void>
  submitClarificationResponse: (input: ClarificationResponseInput) => Promise<void>
  submitLayoutTreeUserMessage: (input: UserMessageInput) => Promise<void>
  showInteractionError: (message: string) => void
  clearError: () => void
}

type ConversationStore = ConversationState & ConversationActions

const persistedWorkspaceState = readConversationWorkspaceState()
const conversationPollHandles = new Map<string, ReturnType<typeof setTimeout>>()
const DEFAULT_CONVERSATION_POLL_DELAY_MS = 5000
const MIN_CONVERSATION_REPOLL_DELAY_MS = 5000
const MAX_CONVERSATION_REPOLL_DELAY_MS = 5000
const PENDING_TURN_TIMEOUT_MS = 15 * 60 * 1000
const ENABLE_MESSAGE_INTERACTION_STATE_RECONCILIATION = true

interface PendingTurnState {
  state: "idle" | "thinking" | "clarifying" | "submitting_selection" | "error"
  turnId?: string
  sourceMessageId?: string
}

function toMessage(error: unknown): string {
  if (error instanceof ServiceError) return error.message
  return translate("conversation.serviceUnavailable")
}

function createLocalDraftWorkspace(seedDraft = ""): LocalDraftWorkspace {
  const now = new Date().toISOString()
  return {
    localDraftId: createId("local_draft"),
    title: translate("conversation.newConversationTitle"),
    draft: seedDraft,
    createdAt: now,
    updatedAt: now,
  }
}

function loadDraftForSelection(
  selection: ConversationWorkspaceSelection | null,
  localDraftWorkspace: LocalDraftWorkspace | null,
  draftsByKey: Record<string, string>,
): Record<string, string> {
  if (!selection) return draftsByKey

  const draftKey =
    selection.kind === "conversation"
      ? conversationDraftKeyForConversation(selection.conversationId)
      : localDraftWorkspace
        ? conversationDraftKeyForLocalDraft(localDraftWorkspace.localDraftId)
        : null

  if (!draftKey || draftKey in draftsByKey) return draftsByKey

  return {
    ...draftsByKey,
    [draftKey]: readConversationDraft(draftKey),
  }
}

function persistWorkspaceState(state: ConversationWorkspaceState): void {
  writeConversationWorkspaceState(state)
}

async function loadConversationPayload(conversationId: string): Promise<{
  detail: ConversationDetail
  messages: ConversationMessageSummary[]
}> {
  const [detail, messages] = await Promise.all([
    getConversationDetail(conversationId),
    getConversationMessages(conversationId),
  ])

  return { detail, messages }
}

function resolveInteractionStateFromMessages(
  messages: ConversationMessageSummary[],
  activeTurnId?: string,
): {
  interactionState: ConversationDetail["interactionState"]
  turnId?: string
  sourceMessageId?: string
} | null {
  const candidateMessages = [...messages]
    .filter((message): message is ConversationMessageSummary & {
      richPayload: Extract<ConversationMessageSummary["richPayload"], { kind: "layout_tree" }>
    } => {
      if (activeTurnId && message.turnId && message.turnId !== activeTurnId) return false
      return message.richPayload?.kind === "layout_tree"
    })
    .reverse()

  for (const message of candidateMessages) {
    const interactionState = message.richPayload?.data.stateSnapshot?.interaction?.status
    if (!interactionState) continue

    return {
      interactionState,
      turnId: message.turnId,
      sourceMessageId: message.messageId,
    }
  }

  return null
}

function reconcileDetailWithMessages(
  detail: ConversationDetail,
  messages: ConversationMessageSummary[],
): ConversationDetail {
  if (!ENABLE_MESSAGE_INTERACTION_STATE_RECONCILIATION) return detail

  const messageState = resolveInteractionStateFromMessages(messages, detail.activeTurnId)
  if (!messageState?.interactionState) return detail

  const interactionState = messageState.interactionState

  return {
    ...detail,
    interactionState,
    activeTurnId: shouldClearActiveTurnId(interactionState) ? undefined : detail.activeTurnId ?? messageState.turnId,
  }
}

function replaceConversation(
  conversations: ConversationSummary[],
  conversation: ConversationSummary,
): ConversationSummary[] {
  return [conversation, ...conversations.filter((item) => item.conversationId !== conversation.conversationId)]
}

function isActiveInteractionState(state?: ConversationDetail["interactionState"]): boolean {
  return state === "thinking" || state === "executing" || state === "clarifying"
}

function pendingTurnFromInteractionState(input: {
  interactionState?: ConversationDetail["interactionState"]
  activeTurnId?: string | null
  turnId?: string
  sourceMessageId?: string
}): PendingTurnState {
  const turnId = input.activeTurnId ?? input.turnId

  if (isActiveInteractionState(input.interactionState) && input.interactionState !== "clarifying") {
    return {
      state: "thinking",
      ...(turnId ? { turnId } : {}),
    }
  }

  if (input.interactionState === "clarifying") {
    return {
      state: "clarifying",
      ...(turnId ? { turnId } : {}),
      ...(input.sourceMessageId ? { sourceMessageId: input.sourceMessageId } : {}),
    }
  }

  if (input.interactionState === "error") {
    return {
      state: "error",
      ...(turnId ? { turnId } : {}),
    }
  }

  return { state: "idle" }
}

function shouldContinuePolling(pendingTurn: PendingTurnState, hasMore = false): boolean {
  return hasMore || pendingTurn.state === "thinking"
}

function pendingTurnIdentity(pendingTurn?: PendingTurnState): string {
  if (!pendingTurn) return ""
  return `${pendingTurn.turnId ?? ""}:${pendingTurn.sourceMessageId ?? ""}`
}

function omitRecordKey<T>(record: Record<string, T>, key: string): Record<string, T> {
  if (!(key in record)) return record

  const { [key]: _removed, ...rest } = record
  return rest
}

function pendingTurnStartedAtMsByConversation(
  current: ConversationState,
  conversationId: string,
  pendingTurn: PendingTurnState,
  restoredStartedAtMs?: number,
): Record<string, number> {
  const startedAtByConversation = current.pendingTurnStartedAtMsByConversationId ?? {}

  if (pendingTurn.state !== "thinking") {
    return omitRecordKey(startedAtByConversation, conversationId)
  }

  const previousPendingTurn = current.pendingTurnByConversationId[conversationId]
  const previousStartedAtMs = startedAtByConversation[conversationId]

  if (
    previousPendingTurn?.state === "thinking" &&
    pendingTurnIdentity(previousPendingTurn) === pendingTurnIdentity(pendingTurn) &&
    typeof previousStartedAtMs === "number"
  ) {
    return startedAtByConversation
  }

  return {
    ...startedAtByConversation,
    [conversationId]: restoredStartedAtMs ?? Date.now(),
  }
}

function toTimestampMs(value?: string): number | undefined {
  if (!value) return undefined
  const timestamp = Date.parse(value)
  return Number.isFinite(timestamp) ? timestamp : undefined
}

function resolvePendingTurnStartedAtMsFromMessages(
  messages: ConversationMessageSummary[],
  turnId?: string,
): number | undefined {
  if (!turnId) return undefined

  const turnMessages = messages.filter((message) => message.turnId === turnId)
  const userTurnMessages = turnMessages.filter((message) => message.role === "user")
  const candidateMessages = userTurnMessages.length > 0 ? userTurnMessages : turnMessages

  const timestamps = candidateMessages
    .map((message) => toTimestampMs(message.createdAt))
    .filter((timestamp): timestamp is number => typeof timestamp === "number")

  if (timestamps.length === 0) return undefined
  return Math.min(...timestamps)
}

function resolveInteractionStateFromMessage(
  message?: ConversationMessageSummary,
): ConversationDetail["interactionState"] | undefined {
  if (message?.richPayload?.kind === "layout_tree") {
    return message.richPayload.data.stateSnapshot?.interaction?.status
  }

  return undefined
}

function shouldClearActiveTurnId(interactionState?: ConversationDetail["interactionState"]): boolean {
  return interactionState === "idle" || interactionState === "completed" || interactionState === "error"
}

function normalizeConversationRepollDelay(delayMs: number): number {
  return Number.isFinite(delayMs) && delayMs > 0
    ? Math.min(Math.max(delayMs, MIN_CONVERSATION_REPOLL_DELAY_MS), MAX_CONVERSATION_REPOLL_DELAY_MS)
    : DEFAULT_CONVERSATION_POLL_DELAY_MS
}

function normalizeAcceptedPollDelay(delayMs: number): number {
  if (!Number.isFinite(delayMs) || delayMs <= 0) return 0
  return normalizeConversationRepollDelay(delayMs)
}

function upsertMessage(
  messages: ConversationMessageSummary[],
  incoming: ConversationMessageSummary,
): ConversationMessageSummary[] {
  const existingIndex = messages.findIndex((message) => message.messageId === incoming.messageId)
  if (existingIndex < 0) {
    return [...messages, incoming]
  }

  return messages.map((message, index) =>
    index === existingIndex
      ? {
          ...message,
          ...incoming,
        }
      : message,
  )
}

function workspaceRequestKey(selection: ConversationWorkspaceSelection): string {
  return selection.kind === "conversation"
    ? `conversation:${selection.conversationId}`
    : `local:${selection.localDraftId}`
}

export const useConversationStore = create<ConversationStore>((set, get) => {
  const applyConversationStatusEvent = (conversationId: string, event: ConversationStatusEvent): void => {
    set((current) => {
      const applied = current.appliedStatusEventIdsByConversationId[conversationId] ?? {}
      if (applied[event.statusEventId]) return current

      const latestSequence = current.latestEventSequenceByConversationId[conversationId] ?? 0
      if (event.sequence <= latestSequence) {
        return {
          ...current,
          appliedStatusEventIdsByConversationId: {
            ...current.appliedStatusEventIdsByConversationId,
            [conversationId]: {
              ...applied,
              [event.statusEventId]: true,
            },
          },
        }
      }

      const currentMessages = current.messagesByConversationId[conversationId] ?? []
      const nextMessages = event.message ? upsertMessage(currentMessages, event.message) : currentMessages
      const currentDetail = current.detailsByConversationId[conversationId]
      const messageInteractionState = resolveInteractionStateFromMessage(event.message)
      const activeTurnId =
        shouldClearActiveTurnId(messageInteractionState)
          ? undefined
          : event.activeTurnId === null
            ? undefined
            : event.activeTurnId ?? currentDetail?.activeTurnId
      const interactionState = messageInteractionState ?? event.interactionState ?? currentDetail?.interactionState
      const nextPendingTurn =
        event.eventType === "interaction.status_changed" || event.interactionState || messageInteractionState
          ? pendingTurnFromInteractionState({
              interactionState,
              activeTurnId,
              turnId: event.turnId,
              sourceMessageId: event.messageId,
            })
          : current.pendingTurnByConversationId[conversationId] ?? { state: "idle" }
      const nextPendingTurnStartedAtMsByConversationId = pendingTurnStartedAtMsByConversation(
        current,
        conversationId,
        nextPendingTurn,
      )

      return {
        ...current,
        detailsByConversationId: currentDetail
          ? {
              ...current.detailsByConversationId,
              [conversationId]: {
                ...currentDetail,
                ...(interactionState ? { interactionState } : {}),
                activeTurnId,
              },
            }
          : current.detailsByConversationId,
        messagesByConversationId: {
          ...current.messagesByConversationId,
          [conversationId]: nextMessages,
        },
        pendingTurnByConversationId: {
          ...current.pendingTurnByConversationId,
          [conversationId]: nextPendingTurn,
        },
        pendingTurnStartedAtMsByConversationId: nextPendingTurnStartedAtMsByConversationId,
        latestEventSequenceByConversationId: {
          ...current.latestEventSequenceByConversationId,
          [conversationId]: event.sequence,
        },
        appliedStatusEventIdsByConversationId: {
          ...current.appliedStatusEventIdsByConversationId,
          [conversationId]: {
            ...applied,
            [event.statusEventId]: true,
          },
        },
      }
    })
  }

  const applyAcceptedResponse = (accepted: MessageAcceptedResult): PendingTurnState => {
    const conversationId = accepted.conversation.conversationId
    const nextPollDelayMs = normalizeAcceptedPollDelay(accepted.nextPollAfterMs)
    const pendingTurn = pendingTurnFromInteractionState({
      interactionState: accepted.statusEvent.interactionState ?? accepted.conversation.interactionState,
      activeTurnId: accepted.conversation.activeTurnId,
      turnId: accepted.statusEvent.turnId ?? accepted.message.turnId,
      sourceMessageId: accepted.statusEvent.messageId ?? accepted.message.messageId,
    })

    set((current) => ({
      detailsByConversationId: {
        ...current.detailsByConversationId,
        [conversationId]: accepted.conversation,
      },
      messagesByConversationId: {
        ...current.messagesByConversationId,
        [conversationId]: upsertMessage(current.messagesByConversationId[conversationId] ?? [], accepted.message),
      },
      conversations: replaceConversation(current.conversations, accepted.conversation),
      pendingTurnByConversationId: {
        ...current.pendingTurnByConversationId,
        [conversationId]: pendingTurn,
      },
      pendingTurnStartedAtMsByConversationId: pendingTurnStartedAtMsByConversation(
        current,
        conversationId,
        pendingTurn,
      ),
      nextPollAfterMsByConversationId: {
        ...current.nextPollAfterMsByConversationId,
        [conversationId]: nextPollDelayMs,
      },
    }))

    applyConversationStatusEvent(conversationId, accepted.statusEvent)
    return pendingTurn
  }

  const stopEventRefresh = (conversationId?: string | null): void => {
    if (!conversationId) return
    const existingHandle = conversationPollHandles.get(conversationId)
    if (existingHandle) {
      clearTimeout(existingHandle)
      conversationPollHandles.delete(conversationId)
    }
  }

  const isSelectedConversation = (conversationId: string): boolean => {
    const selection = get().selectedWorkspace
    return selection?.kind === "conversation" && selection.conversationId === conversationId
  }

  const pendingTurnRemainingMs = (conversationId: string): number | undefined => {
    const state = get()
    const pendingTurn = state.pendingTurnByConversationId[conversationId]
    const startedAtMs = state.pendingTurnStartedAtMsByConversationId[conversationId]

    if (pendingTurn?.state !== "thinking" || typeof startedAtMs !== "number") return undefined

    return Math.max(PENDING_TURN_TIMEOUT_MS - (Date.now() - startedAtMs), 0)
  }

  const markPendingTurnTimedOut = (conversationId: string): void => {
    set((current) => {
      const pendingTurn = current.pendingTurnByConversationId[conversationId]
      if (pendingTurn?.state !== "thinking") return current

      const currentDetail = current.detailsByConversationId[conversationId]

      return {
        ...current,
        error: translate("conversation.pendingTurnTimeout"),
        detailsByConversationId: currentDetail
          ? {
              ...current.detailsByConversationId,
              [conversationId]: {
                ...currentDetail,
                interactionState: "error",
                activeTurnId: undefined,
              },
            }
          : current.detailsByConversationId,
        pendingTurnByConversationId: {
          ...current.pendingTurnByConversationId,
          [conversationId]: {
            state: "error",
            ...(pendingTurn.turnId ? { turnId: pendingTurn.turnId } : {}),
            ...(pendingTurn.sourceMessageId ? { sourceMessageId: pendingTurn.sourceMessageId } : {}),
          },
        },
        pendingTurnStartedAtMsByConversationId: omitRecordKey(
          current.pendingTurnStartedAtMsByConversationId ?? {},
          conversationId,
        ),
        submittingWorkspaceKey:
          current.submittingWorkspaceKey === `conversation:${conversationId}` ? null : current.submittingWorkspaceKey,
      }
    })
  }

  const scheduleEventRefresh = (conversationId: string, delayMs: number): void => {
    stopEventRefresh(conversationId)
    const remainingMs = pendingTurnRemainingMs(conversationId)

    if (remainingMs !== undefined && remainingMs <= 0) {
      markPendingTurnTimedOut(conversationId)
      return
    }

    const handle = setTimeout(async () => {
      conversationPollHandles.delete(conversationId)
      if (!isSelectedConversation(conversationId)) return
      const remainingMs = pendingTurnRemainingMs(conversationId)

      if (remainingMs !== undefined && remainingMs <= 0) {
        markPendingTurnTimedOut(conversationId)
        return
      }

      try {
        const cursor = get().eventCursorByConversationId[conversationId] ?? null
        const previousLatestSequence = get().latestEventSequenceByConversationId[conversationId] ?? 0
        const result = await listConversationEvents(conversationId, { cursor })

        for (const event of result.events) {
          applyConversationStatusEvent(conversationId, event)
        }

        let pendingTurn = get().pendingTurnByConversationId[conversationId] ?? { state: "idle" }

        if (result.interactionState && !result.hasMore && pendingTurn.state === "thinking") {
          const detail = get().detailsByConversationId[conversationId]
          pendingTurn = pendingTurnFromInteractionState({
            interactionState: result.interactionState,
            activeTurnId: detail?.activeTurnId,
            turnId: pendingTurn.turnId,
            sourceMessageId: pendingTurn.sourceMessageId,
          })
        }

        let reconciledDetail: ConversationDetail | null = null

        if (
          !result.hasMore &&
          pendingTurn.state === "thinking" &&
          result.events.length === 0 &&
          result.latestSequence <= previousLatestSequence
        ) {
          try {
            const detail = await getConversationDetail(conversationId)
            const reconciledPendingTurn = pendingTurnFromInteractionState({
              interactionState: detail.interactionState,
              activeTurnId: detail.activeTurnId,
              turnId: pendingTurn.turnId,
              sourceMessageId: pendingTurn.sourceMessageId,
            })

            if (reconciledPendingTurn.state !== pendingTurn.state) {
              pendingTurn = reconciledPendingTurn
              reconciledDetail = detail
            }
          } catch {
            // Keep polling if the control-state reconciliation request fails.
          }
        }

        const nextPollDelayMs = normalizeConversationRepollDelay(result.recommendedPollIntervalMs)

        set((current) => ({
          detailsByConversationId: reconciledDetail
            ? {
                ...current.detailsByConversationId,
                [conversationId]: reconciledDetail,
              }
            : current.detailsByConversationId,
          eventCursorByConversationId: {
            ...current.eventCursorByConversationId,
            [conversationId]: result.nextCursor ?? null,
          },
          nextPollAfterMsByConversationId: {
            ...current.nextPollAfterMsByConversationId,
            [conversationId]: nextPollDelayMs,
          },
          pendingTurnByConversationId: {
            ...current.pendingTurnByConversationId,
            [conversationId]: pendingTurn,
          },
          pendingTurnStartedAtMsByConversationId: pendingTurnStartedAtMsByConversation(
            current,
            conversationId,
            pendingTurn,
          ),
          submittingWorkspaceKey:
            current.submittingWorkspaceKey === `conversation:${conversationId}` ? null : current.submittingWorkspaceKey,
        }))

        if (
          isSelectedConversation(conversationId) &&
          shouldContinuePolling(get().pendingTurnByConversationId[conversationId] ?? pendingTurn, result.hasMore)
        ) {
          scheduleEventRefresh(conversationId, result.hasMore ? 0 : nextPollDelayMs)
        }
      } catch (error) {
        set((current) => ({
          error: toMessage(error),
          pendingTurnByConversationId: {
            ...current.pendingTurnByConversationId,
            [conversationId]: {
              state: "error",
              ...(current.pendingTurnByConversationId[conversationId]?.turnId
                ? { turnId: current.pendingTurnByConversationId[conversationId].turnId }
                : {}),
            },
          },
          pendingTurnStartedAtMsByConversationId: omitRecordKey(
            current.pendingTurnStartedAtMsByConversationId ?? {},
            conversationId,
          ),
          submittingWorkspaceKey:
            current.submittingWorkspaceKey === `conversation:${conversationId}` ? null : current.submittingWorkspaceKey,
        }))
      }
    }, Math.max(0, remainingMs === undefined ? delayMs : Math.min(delayMs, remainingMs)))

    conversationPollHandles.set(conversationId, handle)
  }

  return {
    bootstrapped: false,
    listLoading: false,
    conversationLoading: false,
    query: "",
    error: null,
    conversations: [],
    selectedWorkspace: persistedWorkspaceState.selectedWorkspace,
    localDraftWorkspace: persistedWorkspaceState.localDraftWorkspace,
    detailsByConversationId: {},
    messagesByConversationId: {},
    pendingTurnByConversationId: {},
    pendingTurnStartedAtMsByConversationId: {},
    nextPollAfterMsByConversationId: {},
    eventCursorByConversationId: {},
    latestEventSequenceByConversationId: {},
    appliedStatusEventIdsByConversationId: {},
    submittingWorkspaceKey: null,
    draftsByKey: loadDraftForSelection(
      persistedWorkspaceState.selectedWorkspace,
      persistedWorkspaceState.localDraftWorkspace,
      {},
    ),

    hydrate: async () => {
      set({ listLoading: true, error: null })

      try {
        const conversations = await listConversations()
        const state = get()

      let localDraftWorkspace = state.localDraftWorkspace
      let selectedWorkspace = state.selectedWorkspace
      const previouslySelectedConversationId =
        state.selectedWorkspace?.kind === "conversation" ? state.selectedWorkspace.conversationId : undefined
      const selectedConversationId =
        selectedWorkspace?.kind === "conversation" ? selectedWorkspace.conversationId : null

      if (selectedConversationId && !conversations.some((conversation) => conversation.conversationId === selectedConversationId)) {
        selectedWorkspace = null
      }

      if (selectedWorkspace?.kind === "localDraft") {
        if (!localDraftWorkspace || localDraftWorkspace.localDraftId !== selectedWorkspace.localDraftId) {
          localDraftWorkspace = createLocalDraftWorkspace()
          selectedWorkspace = {
            kind: "localDraft",
            localDraftId: localDraftWorkspace.localDraftId,
          }
        }
      }

      if (!selectedWorkspace) {
        localDraftWorkspace = localDraftWorkspace ?? createLocalDraftWorkspace()
        selectedWorkspace = {
          kind: "localDraft",
          localDraftId: localDraftWorkspace.localDraftId,
        }
      }

      const resolvedSelectedWorkspace = selectedWorkspace as ConversationWorkspaceSelection
      const draftsByKey = loadDraftForSelection(resolvedSelectedWorkspace, localDraftWorkspace, state.draftsByKey)

      if (
        previouslySelectedConversationId &&
        (resolvedSelectedWorkspace.kind !== "conversation" ||
          resolvedSelectedWorkspace.conversationId !== previouslySelectedConversationId)
      ) {
        stopEventRefresh(previouslySelectedConversationId)
      }

      set({
        conversations,
        selectedWorkspace: resolvedSelectedWorkspace,
        localDraftWorkspace,
        draftsByKey,
        listLoading: false,
        bootstrapped: true,
      })

      persistWorkspaceState({ selectedWorkspace: resolvedSelectedWorkspace, localDraftWorkspace })

        if (resolvedSelectedWorkspace.kind === "conversation") {
          set({ conversationLoading: true })
          const { detail, messages } = await loadConversationPayload(resolvedSelectedWorkspace.conversationId)
          const reconciledDetail = reconcileDetailWithMessages(detail, messages)
          const pendingTurn = pendingTurnFromInteractionState({
            interactionState: reconciledDetail.interactionState,
            activeTurnId: reconciledDetail.activeTurnId,
          })
          const restoredStartedAtMs = resolvePendingTurnStartedAtMsFromMessages(
            messages,
            pendingTurn.turnId ?? reconciledDetail.activeTurnId,
          )

          set((current) => ({
            conversationLoading: false,
            detailsByConversationId: {
              ...current.detailsByConversationId,
              [detail.conversationId]: reconciledDetail,
            },
            messagesByConversationId: {
              ...current.messagesByConversationId,
              [detail.conversationId]: messages,
            },
            pendingTurnByConversationId: {
              ...current.pendingTurnByConversationId,
              [detail.conversationId]: pendingTurn,
            },
            pendingTurnStartedAtMsByConversationId: pendingTurnStartedAtMsByConversation(
              current,
              detail.conversationId,
              pendingTurn,
              restoredStartedAtMs,
            ),
          }))

          if (shouldContinuePolling(pendingTurn)) {
            scheduleEventRefresh(detail.conversationId, DEFAULT_CONVERSATION_POLL_DELAY_MS)
          }
        }
      } catch (error) {
        const localDraftWorkspace = get().localDraftWorkspace ?? createLocalDraftWorkspace()
        const currentSelection = get().selectedWorkspace
        const previouslySelectedConversationId =
          currentSelection?.kind === "conversation" ? currentSelection.conversationId : undefined
        const selectedWorkspace: ConversationWorkspaceSelection = {
          kind: "localDraft",
          localDraftId: localDraftWorkspace.localDraftId,
        }

        stopEventRefresh(previouslySelectedConversationId)

        set((current) => ({
          error: toMessage(error),
          listLoading: false,
          bootstrapped: true,
          selectedWorkspace,
          localDraftWorkspace,
          draftsByKey: loadDraftForSelection(selectedWorkspace, localDraftWorkspace, current.draftsByKey),
        }))

        persistWorkspaceState({ selectedWorkspace, localDraftWorkspace })
      }
    },

    activateLocalDraftWorkspace: () => {
      const state = get()
      const previouslySelectedConversationId =
        state.selectedWorkspace?.kind === "conversation" ? state.selectedWorkspace.conversationId : undefined
      const localDraftWorkspace = state.localDraftWorkspace ?? createLocalDraftWorkspace()
      const selectedWorkspace: ConversationWorkspaceSelection = {
        kind: "localDraft",
        localDraftId: localDraftWorkspace.localDraftId,
      }
      const draftsByKey = loadDraftForSelection(selectedWorkspace, localDraftWorkspace, state.draftsByKey)

      stopEventRefresh(previouslySelectedConversationId)

      set({
        selectedWorkspace,
        localDraftWorkspace,
        draftsByKey,
        error: null,
      })

      persistWorkspaceState({ selectedWorkspace, localDraftWorkspace })
    },

    selectConversation: async (conversationId) => {
      const state = get()
      const previouslySelectedConversationId =
        state.selectedWorkspace?.kind === "conversation" ? state.selectedWorkspace.conversationId : undefined
      const selectedWorkspace: ConversationWorkspaceSelection = {
        kind: "conversation",
        conversationId,
      }
      const draftsByKey = loadDraftForSelection(selectedWorkspace, state.localDraftWorkspace, state.draftsByKey)

      stopEventRefresh(previouslySelectedConversationId)

      set({
        selectedWorkspace,
        draftsByKey,
        conversationLoading: true,
        error: null,
      })

      persistWorkspaceState({ selectedWorkspace, localDraftWorkspace: state.localDraftWorkspace })

      try {
        const { detail, messages } = await loadConversationPayload(conversationId)
        const reconciledDetail = reconcileDetailWithMessages(detail, messages)
        const pendingTurn = pendingTurnFromInteractionState({
          interactionState: reconciledDetail.interactionState,
          activeTurnId: reconciledDetail.activeTurnId,
        })
        const restoredStartedAtMs = resolvePendingTurnStartedAtMsFromMessages(
          messages,
          pendingTurn.turnId ?? reconciledDetail.activeTurnId,
        )

        set((current) => ({
          conversationLoading: false,
          detailsByConversationId: {
            ...current.detailsByConversationId,
            [conversationId]: reconciledDetail,
          },
          messagesByConversationId: {
            ...current.messagesByConversationId,
            [conversationId]: messages,
          },
          pendingTurnByConversationId: {
            ...current.pendingTurnByConversationId,
            [conversationId]: pendingTurn,
          },
          pendingTurnStartedAtMsByConversationId: pendingTurnStartedAtMsByConversation(
            current,
            conversationId,
            pendingTurn,
            restoredStartedAtMs,
          ),
        }))

        if (shouldContinuePolling(pendingTurn)) {
          scheduleEventRefresh(conversationId, DEFAULT_CONVERSATION_POLL_DELAY_MS)
        }
      } catch (error) {
        set({ conversationLoading: false, error: toMessage(error) })
        throw error
      }
    },

    setSearchQuery: async (query) => {
      set({ query, listLoading: true, error: null })

      try {
        const conversations = query.trim() ? await searchConversations(query) : await listConversations()
        set({ conversations, listLoading: false })
      } catch (error) {
        set({ listLoading: false, error: toMessage(error) })
      }
    },

    setDraft: (value) => {
      const state = get()
      const selection = state.selectedWorkspace
      if (!selection) return

      const draftKey =
        selection.kind === "conversation"
          ? conversationDraftKeyForConversation(selection.conversationId)
          : conversationDraftKeyForLocalDraft(selection.localDraftId)

      writeConversationDraft(draftKey, value)

      set((current) => {
        const nextDraftsByKey = {
          ...current.draftsByKey,
          [draftKey]: value,
        }

        if (selection.kind === "localDraft" && current.localDraftWorkspace) {
          const updatedAt = new Date().toISOString()
          const localDraftWorkspace = {
            ...current.localDraftWorkspace,
            draft: value,
            updatedAt,
          }

          persistWorkspaceState({
            selectedWorkspace: current.selectedWorkspace,
            localDraftWorkspace,
          })

          return {
            draftsByKey: nextDraftsByKey,
            localDraftWorkspace,
          }
        }

        return {
          draftsByKey: nextDraftsByKey,
        }
      })
    },

    submitComposerMessage: async () => {
      const state = get()
      const selection = state.selectedWorkspace
      if (!selection) return

      const draftKey =
        selection.kind === "conversation"
          ? conversationDraftKeyForConversation(selection.conversationId)
          : conversationDraftKeyForLocalDraft(selection.localDraftId)

      const content = (state.draftsByKey[draftKey] ?? "").trim()
      if (!content) return

      set({
        error: null,
        submittingWorkspaceKey: workspaceRequestKey(selection),
      })

      try {
        if (selection.kind === "localDraft") {
          const accepted = await createConversation({
            initialMessage: {
              type: "user_message",
              content,
            },
          })

          const pendingTurn = applyAcceptedResponse(accepted)
          const nextLocalDraftWorkspace = createLocalDraftWorkspace()
          const nextSelection: ConversationWorkspaceSelection = {
            kind: "conversation",
            conversationId: accepted.conversation.conversationId,
          }

          removeConversationDraft(draftKey)

          set((current) => ({
            selectedWorkspace: nextSelection,
            localDraftWorkspace: nextLocalDraftWorkspace,
            draftsByKey: {
              ...current.draftsByKey,
              [draftKey]: "",
              [conversationDraftKeyForLocalDraft(nextLocalDraftWorkspace.localDraftId)]: "",
            },
            submittingWorkspaceKey: null,
            error: null,
          }))

          persistWorkspaceState({
            selectedWorkspace: nextSelection,
            localDraftWorkspace: nextLocalDraftWorkspace,
          })

          if (shouldContinuePolling(pendingTurn)) {
            scheduleEventRefresh(
              accepted.conversation.conversationId,
              normalizeAcceptedPollDelay(accepted.nextPollAfterMs),
            )
          }
          return
        }

        const accepted = await sendMessage(selection.conversationId, {
          type: "user_message",
          content,
        })

        removeConversationDraft(draftKey)
        const pendingTurn = applyAcceptedResponse(accepted)

        set((current) => ({
          draftsByKey: {
            ...current.draftsByKey,
            [draftKey]: "",
          },
          submittingWorkspaceKey: null,
          error: null,
        }))

        if (shouldContinuePolling(pendingTurn)) {
          scheduleEventRefresh(selection.conversationId, normalizeAcceptedPollDelay(accepted.nextPollAfterMs))
        }
      } catch (error) {
        set((current) => ({
          error: toMessage(error),
          submittingWorkspaceKey: null,
          pendingTurnByConversationId:
            selection.kind === "conversation"
              ? {
                  ...current.pendingTurnByConversationId,
                  [selection.conversationId]: { state: "error" },
                }
              : current.pendingTurnByConversationId,
          pendingTurnStartedAtMsByConversationId:
            selection.kind === "conversation"
              ? omitRecordKey(current.pendingTurnStartedAtMsByConversationId ?? {}, selection.conversationId)
              : current.pendingTurnStartedAtMsByConversationId,
        }))
      }
    },

    submitCandidateSelection: async (input) => {
      const state = get()
      const selection = state.selectedWorkspace
      if (!selection || selection.kind !== "conversation") return

      set((current) => ({
        error: null,
        pendingTurnByConversationId: {
          ...current.pendingTurnByConversationId,
          [selection.conversationId]: {
            state: "submitting_selection",
            sourceMessageId: input.messageId,
          },
        },
        pendingTurnStartedAtMsByConversationId: omitRecordKey(
          current.pendingTurnStartedAtMsByConversationId ?? {},
          selection.conversationId,
        ),
      }))

      try {
        const accepted = await sendMessage(selection.conversationId, input)
        const pendingTurn = applyAcceptedResponse(accepted)

        set({ error: null })

        if (shouldContinuePolling(pendingTurn)) {
          scheduleEventRefresh(selection.conversationId, normalizeAcceptedPollDelay(accepted.nextPollAfterMs))
        }
      } catch (error) {
        set((current) => ({
          error: toMessage(error),
          pendingTurnByConversationId: {
            ...current.pendingTurnByConversationId,
            [selection.conversationId]: {
              state: "clarifying",
              sourceMessageId: input.messageId,
            },
          },
          pendingTurnStartedAtMsByConversationId: omitRecordKey(
            current.pendingTurnStartedAtMsByConversationId ?? {},
            selection.conversationId,
          ),
        }))
      }
    },

    submitClarificationResponse: async (input) => {
      const state = get()
      const selection = state.selectedWorkspace
      if (!selection || selection.kind !== "conversation") return

      set((current) => ({
        error: null,
        pendingTurnByConversationId: {
          ...current.pendingTurnByConversationId,
          [selection.conversationId]: {
            state: "submitting_selection",
            sourceMessageId: input.messageId,
          },
        },
        pendingTurnStartedAtMsByConversationId: omitRecordKey(
          current.pendingTurnStartedAtMsByConversationId ?? {},
          selection.conversationId,
        ),
      }))

      try {
        const accepted = await sendMessage(selection.conversationId, input)
        const pendingTurn = applyAcceptedResponse(accepted)

        set({ error: null })

        if (shouldContinuePolling(pendingTurn)) {
          scheduleEventRefresh(selection.conversationId, normalizeAcceptedPollDelay(accepted.nextPollAfterMs))
        }
      } catch (error) {
        set((current) => ({
          error: toMessage(error),
          pendingTurnByConversationId: {
            ...current.pendingTurnByConversationId,
            [selection.conversationId]: {
              state: "clarifying",
              sourceMessageId: input.messageId,
            },
          },
          pendingTurnStartedAtMsByConversationId: omitRecordKey(
            current.pendingTurnStartedAtMsByConversationId ?? {},
            selection.conversationId,
          ),
        }))
      }
    },

    submitLayoutTreeUserMessage: async (input) => {
      const state = get()
      const selection = state.selectedWorkspace
      if (!selection || selection.kind !== "conversation") return

      const content = input.content.trim()
      if (!content) return

      set({
        error: null,
        submittingWorkspaceKey: workspaceRequestKey(selection),
      })

      try {
        const accepted = await sendMessage(selection.conversationId, {
          type: "user_message",
          content,
        })
        const pendingTurn = applyAcceptedResponse(accepted)

        set({
          error: null,
          submittingWorkspaceKey: null,
        })

        if (shouldContinuePolling(pendingTurn)) {
          scheduleEventRefresh(selection.conversationId, normalizeAcceptedPollDelay(accepted.nextPollAfterMs))
        }
      } catch (error) {
        set((current) => ({
          error: toMessage(error),
          submittingWorkspaceKey: null,
          pendingTurnByConversationId: {
            ...current.pendingTurnByConversationId,
            [selection.conversationId]: { state: "error" },
          },
          pendingTurnStartedAtMsByConversationId: omitRecordKey(
            current.pendingTurnStartedAtMsByConversationId ?? {},
            selection.conversationId,
          ),
        }))
      }
    },

    showInteractionError: (message) => set({ error: message }),

    clearError: () => set({ error: null }),
  }
})
