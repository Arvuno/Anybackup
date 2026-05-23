import type { ConversationWorkspaceSelection, ConversationWorkspaceState, LocalDraftWorkspace } from "@/types/conversation"

const DRAFT_PREFIX = "agent_web_conversation_draft_"
const WORKSPACE_STATE_KEY = "agent_web_conversation_workspace"

interface PersistedConversationWorkspaceState {
  selectedWorkspace: ConversationWorkspaceSelection | null
  localDraftWorkspace: LocalDraftWorkspace | null
}

function getStorage(): Storage | null {
  if (typeof window === "undefined") return null
  return window.localStorage
}

export function conversationDraftKeyForConversation(conversationId: string): string {
  return `${DRAFT_PREFIX}conversation_${conversationId}`
}

export function conversationDraftKeyForLocalDraft(localDraftId: string): string {
  return `${DRAFT_PREFIX}local_${localDraftId}`
}

export function readConversationDraft(draftKey: string): string {
  const storage = getStorage()
  return storage?.getItem(draftKey) ?? ""
}

export function writeConversationDraft(draftKey: string, value: string): void {
  const storage = getStorage()
  if (!storage) return

  if (!value.trim()) {
    storage.removeItem(draftKey)
    return
  }

  storage.setItem(draftKey, value)
}

export function removeConversationDraft(draftKey: string): void {
  const storage = getStorage()
  storage?.removeItem(draftKey)
}

export function readConversationWorkspaceState(): ConversationWorkspaceState {
  const storage = getStorage()
  if (!storage) return { selectedWorkspace: null, localDraftWorkspace: null }

  const raw = storage.getItem(WORKSPACE_STATE_KEY)
  if (!raw) return { selectedWorkspace: null, localDraftWorkspace: null }

  try {
    const parsed = JSON.parse(raw) as PersistedConversationWorkspaceState
    return {
      selectedWorkspace: parsed.selectedWorkspace ?? null,
      localDraftWorkspace: parsed.localDraftWorkspace ?? null,
    }
  } catch {
    return { selectedWorkspace: null, localDraftWorkspace: null }
  }
}

export function writeConversationWorkspaceState(state: ConversationWorkspaceState): void {
  const storage = getStorage()
  if (!storage) return

  storage.setItem(
    WORKSPACE_STATE_KEY,
    JSON.stringify({
      selectedWorkspace: state.selectedWorkspace,
      localDraftWorkspace: state.localDraftWorkspace,
    } satisfies PersistedConversationWorkspaceState),
  )
}
