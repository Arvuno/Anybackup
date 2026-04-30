import { authApiConfig } from "@/config/auth"
import { ServiceError } from "@/types/auth"

type ErrorBody = {
  error?: string
  error_description?: string
  errorMessage?: string
  message?: string
  traceId?: string
}

function toUrl(path: string): string {
  return `${authApiConfig.basePath}${path.startsWith("/") ? path : `/${path}`}`
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

  const keys: Array<keyof ErrorBody> = ["error_description", "errorMessage", "message", "error"]
  for (const key of keys) {
    const value = body[key]
    if (typeof value === "string" && value.trim()) return value
  }

  return fallback
}

function errorCode(status: number): string {
  if (status === 400) return "VALIDATION_ERROR"
  if (status === 401) return "UNAUTHORIZED"
  if (status === 403) return "FORBIDDEN"
  if (status === 404) return "NOT_FOUND"
  if (status === 409) return "CONFLICT"
  if (status >= 500) return "SERVICE_UNAVAILABLE"
  return "REQUEST_FAILED"
}

async function toServiceError(response: Response, fallback: string): Promise<ServiceError> {
  const body = await readBody(response)
  const traceId =
    response.headers.get("x-trace-id") ??
    response.headers.get("x-request-id") ??
    (isRecord(body) && typeof body.traceId === "string" ? body.traceId : undefined)
  return new ServiceError(errorCode(response.status), errorMessage(body, fallback), response.status, traceId)
}

export async function authApiFetch(path: string, init: RequestInit, fallbackError = "服务暂时不可用，请稍后重试"): Promise<Response> {
  let response: Response
  try {
    response = await fetch(toUrl(path), init)
  } catch {
    throw new ServiceError("NETWORK_UNAVAILABLE", "网络连接不可用，请检查网络后重试")
  }

  if (!response.ok) {
    throw await toServiceError(response, fallbackError)
  }

  return response
}

export async function readJson<T>(response: Response): Promise<T> {
  return (await readBody(response)) as T
}

export async function requestJson<T>(path: string, init: RequestInit, fallbackError?: string): Promise<T> {
  const response = await authApiFetch(path, init, fallbackError)
  return readJson<T>(response)
}

export async function requestNoContent(path: string, init: RequestInit, fallbackError?: string): Promise<Response> {
  return authApiFetch(path, init, fallbackError)
}

export function realmPath(path: string): string {
  return path.replace("{realm}", encodeURIComponent(authApiConfig.realm))
}

export function authHeader(accessToken: string): Record<string, string> {
  return { Authorization: `Bearer ${accessToken}` }
}
