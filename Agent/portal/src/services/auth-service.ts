import { authApiConfig } from "@/config/auth"
import {
  clearStoredSession,
  getStoredSession,
  isAccessTokenValid,
  isStoredSessionValid,
  storeSession,
  touchStoredSession,
} from "@/lib/session"
import { authHeader, realmPath, requestJson, requestNoContent } from "@/services/auth-api-client"
import type { AuthSession, CurrentUser, LoginRequest } from "@/types/auth"
import { ServiceError } from "@/types/auth"
import { translate } from "@/i18n/messages"

const IDLE_TTL_MS = 30 * 60 * 1000
const TOKEN_PATH = realmPath("/realms/{realm}/protocol/openid-connect/token")
const LOGOUT_PATH = realmPath("/realms/{realm}/protocol/openid-connect/logout")
const USERINFO_PATH = realmPath("/realms/{realm}/protocol/openid-connect/userinfo")
const REFRESH_SKEW_MS = 10 * 1000

interface KeycloakTokenResponse {
  access_token: string
  expires_in: number
  refresh_expires_in?: number
  refresh_token?: string
  token_type?: string
  id_token?: string
  scope?: string
}

interface KeycloakUserInfoResponse {
  sub?: string
  preferred_username?: string
  name?: string
  given_name?: string
  family_name?: string
  email?: string
  roles?: string[]
}

function tokenForm(values: Record<string, string | undefined>): URLSearchParams {
  const body = new URLSearchParams()
  for (const [key, value] of Object.entries(values)) {
    if (value) body.set(key, value)
  }
  return body
}

function tokenHeaders(): Record<string, string> {
  return { "Content-Type": "application/x-www-form-urlencoded" }
}

function toCurrentUser(userInfo: KeycloakUserInfoResponse): CurrentUser {
  const username = userInfo.preferred_username ?? userInfo.email ?? userInfo.sub ?? "unknown"
  const fullName = [userInfo.given_name, userInfo.family_name].filter(Boolean).join(" ").trim()
  return {
    id: userInfo.sub ?? username,
    username,
    displayName: userInfo.name || fullName || username,
    role: "backup_admin",
    tenantId: authApiConfig.realm,
  }
}

function createSession(token: KeycloakTokenResponse, user: CurrentUser): AuthSession {
  const now = Date.now()
  return {
    accessToken: token.access_token,
    refreshToken: token.refresh_token,
    tokenType: token.token_type,
    idToken: token.id_token,
    scope: token.scope,
    expiresAt: now + token.expires_in * 1000,
    refreshExpiresAt: token.refresh_expires_in ? now + token.refresh_expires_in * 1000 : undefined,
    idleExpiresAt: now + IDLE_TTL_MS,
    user,
  }
}

function mergeRefreshedSession(session: AuthSession, token: KeycloakTokenResponse): AuthSession {
  const now = Date.now()
  return {
    ...session,
    accessToken: token.access_token,
    refreshToken: token.refresh_token ?? session.refreshToken,
    tokenType: token.token_type ?? session.tokenType,
    idToken: token.id_token ?? session.idToken,
    scope: token.scope ?? session.scope,
    expiresAt: now + token.expires_in * 1000,
    refreshExpiresAt: token.refresh_expires_in ? now + token.refresh_expires_in * 1000 : session.refreshExpiresAt,
    idleExpiresAt: now + IDLE_TTL_MS,
  }
}

async function requestToken(body: URLSearchParams): Promise<KeycloakTokenResponse> {
  return requestJson<KeycloakTokenResponse>(
    TOKEN_PATH,
    {
      method: "POST",
      headers: tokenHeaders(),
      body,
    },
    translate("auth.authServiceUnavailable"),
  )
}

function normalizeLoginError(error: unknown): never {
  if (error instanceof ServiceError) {
    const lowerMessage = error.message.toLowerCase()
    if (error.status === 403 || lowerMessage.includes("disabled")) {
      throw new ServiceError("USER_DISABLED", translate("auth.userDisabled"), error.status, error.traceId)
    }
    if (error.status === 400 || error.status === 401) {
      throw new ServiceError("INVALID_CREDENTIALS", translate("auth.invalidCredentials"), error.status, error.traceId)
    }
  }
  throw error
}

async function fetchCurrentUser(accessToken: string): Promise<CurrentUser> {
  const userInfo = await requestJson<KeycloakUserInfoResponse>(
    USERINFO_PATH,
    {
      method: "GET",
      headers: authHeader(accessToken),
    },
    translate("auth.sessionValidationFailed"),
  )
  return toCurrentUser(userInfo)
}

async function refreshSession(session: AuthSession): Promise<AuthSession> {
  if (!session.refreshToken) {
    clearStoredSession()
    throw new ServiceError("UNAUTHORIZED", translate("auth.sessionExpired"), 401)
  }

  if (session.refreshExpiresAt && session.refreshExpiresAt <= Date.now()) {
    clearStoredSession()
    throw new ServiceError("SESSION_EXPIRED", translate("auth.sessionExpired"), 401)
  }

  try {
    const token = await requestToken(
      tokenForm({
        grant_type: "refresh_token",
        client_id: authApiConfig.clientId,
        refresh_token: session.refreshToken,
      }),
    )
    const refreshed = mergeRefreshedSession(session, token)
    storeSession(refreshed)
    return refreshed
  } catch (error) {
    clearStoredSession()
    if (error instanceof ServiceError) {
      throw new ServiceError("UNAUTHORIZED", translate("auth.sessionExpired"), error.status, error.traceId)
    }
    throw error
  }
}

export async function getAuthorizedSession(): Promise<AuthSession> {
  let session = touchStoredSession()
  if (!session || !isStoredSessionValid(session)) {
    clearStoredSession()
    throw new ServiceError("UNAUTHORIZED", translate("auth.sessionExpired"), 401)
  }

  if (!isAccessTokenValid(session) || session.expiresAt - Date.now() <= REFRESH_SKEW_MS) {
    session = await refreshSession(session)
  }

  return session
}

export async function login(request: LoginRequest): Promise<AuthSession> {
  try {
    const token = await requestToken(
      tokenForm({
        grant_type: "password",
        client_id: authApiConfig.clientId,
        username: request.username,
        password: request.password,
        scope: authApiConfig.scope,
      }),
    )
    const currentUser = await fetchCurrentUser(token.access_token)
    const session = createSession(token, currentUser)
    storeSession(session)
    return session
  } catch (error) {
    clearStoredSession()
    normalizeLoginError(error)
  }
}

export async function getCurrentUser(): Promise<CurrentUser> {
  const session = await getAuthorizedSession()
  const user = await fetchCurrentUser(session.accessToken)
  storeSession({ ...session, user })
  return user
}

export function getCurrentSession(): AuthSession | null {
  const session = getStoredSession()
  return isStoredSessionValid(session) ? session : null
}

export async function logoutCurrentUser(): Promise<void> {
  const session = getStoredSession()
  if (!session || !isStoredSessionValid(session)) {
    clearStoredSession()
    return
  }

  if (!session.refreshToken) {
    clearStoredSession()
    return
  }

  try {
    await requestNoContent(
      LOGOUT_PATH,
      {
        method: "POST",
        headers: tokenHeaders(),
        body: tokenForm({
          client_id: authApiConfig.clientId,
          refresh_token: session.refreshToken,
        }),
      },
      translate("auth.logoutFailed"),
    )
    clearStoredSession()
  } catch (error) {
    if (error instanceof ServiceError && (error.status === 400 || error.status === 401)) {
      clearStoredSession()
      return
    }
    throw error
  }
}
