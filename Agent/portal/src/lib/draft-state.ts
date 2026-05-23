const DRAFT_PREFIXES = ["agent_web_draft_", "agent_web_conversation_draft_"]

export function hasUnsentDrafts(): boolean {
  for (let index = 0; index < localStorage.length; index += 1) {
    const key = localStorage.key(index)
    if (!key || !DRAFT_PREFIXES.some((prefix) => key.startsWith(prefix))) continue

    const value = localStorage.getItem(key)
    if (value && value.trim().length > 0) return true
  }

  return false
}
