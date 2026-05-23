export type UserRole = "backup_admin"

export interface CurrentUser {
  id: string
  username: string
  displayName: string
  role: UserRole
  tenantId?: string
}

export interface AuthSession {
  accessToken: string
  refreshToken?: string
  tokenType?: string
  refreshExpiresAt?: number
  scope?: string
  idToken?: string
  expiresAt: number
  idleExpiresAt: number
  user: CurrentUser
}

export interface LoginRequest {
  username: string
  password: string
}

export type AuthErrorCode =
  | "INVALID_CREDENTIALS"
  | "USER_DISABLED"
  | "UNAUTHORIZED"
  | "TOKEN_EXPIRED"
  | "SESSION_EXPIRED"
  | "NETWORK_UNAVAILABLE"
  | "TIMEOUT"
  | "SERVICE_UNAVAILABLE"

export class ServiceError extends Error {
  code: string
  status?: number
  traceId?: string

  constructor(code: string, message: string, status?: number, traceId?: string) {
    super(message)
    this.name = "ServiceError"
    this.code = code
    this.status = status
    this.traceId = traceId
  }
}

export function isSessionError(code: string): boolean {
  return code === "UNAUTHORIZED" || code === "TOKEN_EXPIRED" || code === "SESSION_EXPIRED"
}
