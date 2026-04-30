function cleanBasePath(value: string): string {
  const normalized = value.trim() || "/api/conversation_service/v1"
  return normalized.endsWith("/") ? normalized.slice(0, -1) : normalized
}

export const conversationApiConfig = {
  basePath: cleanBasePath(import.meta.env.VITE_CONVERSATION_SERVICE_BASE_PATH ?? "/api/conversation_service/v1"),
}
