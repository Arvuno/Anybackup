import { conversationApiConfig } from "@/config/conversation"
import { createId } from "@/lib/ids"
import { translate } from "@/i18n/messages"
import { getAuthorizedSession } from "@/services/auth-service"
import { ServiceError } from "@/types/auth"

type ConversationErrorEnvelope = {
  code?: string
  message?: string
  retryable?: boolean
}

type ConversationErrorBody = {
  error?: ConversationErrorEnvelope | string
  message?: string
  request_id?: string
  trace_id?: string
}

function toUrl(path: string): string {
  return `${conversationApiConfig.basePath}${path.startsWith("/") ? path : `/${path}`}`
}

async function readBody(response: Response): Promise<unknown> {
  if (response.status === 204) return null

  const text = await response.text()
  if (!text) return null

  try {
    return JSON.parse(text)
  } catch {
    return text
  }
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null
}

function errorMessage(body: unknown, fallback: string): string {
  if (!isRecord(body)) return fallback

  const nestedError = body.error
  if (isRecord(nestedError) && typeof nestedError.message === "string" && nestedError.message.trim()) {
    return nestedError.message
  }

  if (typeof body.message === "string" && body.message.trim()) {
    return body.message
  }

  if (typeof nestedError === "string" && nestedError.trim()) {
    return nestedError
  }

  return fallback
}

function errorCode(status: number, body: unknown): string {
  if (isRecord(body)) {
    const nestedError = body.error
    if (isRecord(nestedError) && typeof nestedError.code === "string" && nestedError.code.trim()) {
      return nestedError.code
    }
  }

  if (status === 400) return "BAD_REQUEST"
  if (status === 401) return "UNAUTHORIZED"
  if (status === 403) return "FORBIDDEN"
  if (status === 404) return "NOT_FOUND"
  if (status === 409) return "CONFLICT"
  if (status === 422) return "VALIDATION_FAILED"
  if (status >= 500) return "INTERNAL_ERROR"
  return "REQUEST_FAILED"
}

async function toServiceError(response: Response, fallback: string): Promise<ServiceError> {
  const body = (await readBody(response)) as ConversationErrorBody | string | null
  const traceId =
    response.headers.get("x-trace-id") ??
    response.headers.get("x-request-id") ??
    (isRecord(body) && typeof body.trace_id === "string" ? body.trace_id : undefined)

  return new ServiceError(errorCode(response.status, body), errorMessage(body, fallback), response.status, traceId)
}

async function withAuthorization(init: RequestInit): Promise<RequestInit> {
  const session = await getAuthorizedSession()
  const headers = new Headers(init.headers)

  headers.set("Authorization", `Bearer ${session.accessToken}`)

  if (!headers.has("X-Request-Id")) {
    headers.set("X-Request-Id", createId("req"))
  }

  return {
    ...init,
    headers,
  }
}

export function jsonHeaders(extra?: Record<string, string>): Record<string, string> {
  return {
    "Content-Type": "application/json",
    ...extra,
  }
}

export async function conversationApiFetch(
  path: string,
  init: RequestInit,
  fallbackError?: string,
): Promise<Response> {
  const fallback = fallbackError ?? translate("conversation.serviceUnavailable")
  let response: Response

  try {
    response = await fetch(toUrl(path), await withAuthorization(init))
  } catch {
    throw new ServiceError("NETWORK_UNAVAILABLE", translate("errors.networkUnavailable"))
  }

  if (!response.ok) {
    throw await toServiceError(response, fallback)
  }

  return response
}

export async function readJson<T>(response: Response): Promise<T> {
  return (await readBody(response)) as T
}

export async function requestJson<T>(path: string, init: RequestInit, fallbackError?: string): Promise<T> {
  const response = await conversationApiFetch(path, init, fallbackError)
  return readJson<T>(response)
}

export async function requestNoContent(path: string, init: RequestInit, fallbackError?: string): Promise<Response> {
  return conversationApiFetch(path, init, fallbackError)
}
