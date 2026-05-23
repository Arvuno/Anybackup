import type { AuthSession, CurrentUser } from "@/types/auth"

const AUTH_KEY = "agent_web_auth_session"
const GUIDE_KEY_PREFIX = "agent_web_guide_state"
const GUIDE_DISMISSED_PREFIX = "agent_web_guide_dismissed"

interface GuideState {
  completed?: boolean
  dontShowAgain?: boolean
}

export function getStoredSession(): AuthSession | null {
  const raw = localStorage.getItem(AUTH_KEY)
  if (!raw) return null

  try {
    return JSON.parse(raw) as AuthSession
  } catch {
    localStorage.removeItem(AUTH_KEY)
    return null
  }
}

export function storeSession(session: AuthSession): void {
  localStorage.setItem(AUTH_KEY, JSON.stringify(session))
}

export function clearStoredSession(): void {
  localStorage.removeItem(AUTH_KEY)
}

export function isStoredSessionValid(session: AuthSession | null = getStoredSession()): boolean {
  if (!session) return false
  const now = Date.now()
  return (session.refreshExpiresAt ?? session.expiresAt) > now && session.idleExpiresAt > now
}

export function isAccessTokenValid(session: AuthSession | null = getStoredSession()): boolean {
  if (!session) return false
  return session.expiresAt > Date.now()
}

export function isAuthenticated(): boolean {
  return isStoredSessionValid()
}

export function touchStoredSession(): AuthSession | null {
  const session = getStoredSession()
  if (!session || !isStoredSessionValid(session)) return null

  const touched: AuthSession = {
    ...session,
    idleExpiresAt: Date.now() + 30 * 60 * 1000,
  }
  storeSession(touched)
  return touched
}

export function logout(): void {
  clearStoredSession()
}

function guideUserKey(user?: CurrentUser): string | null {
  const resolvedUser = user ?? getStoredSession()?.user
  if (!resolvedUser) return null
  return `${resolvedUser.tenantId ?? "default"}:${resolvedUser.id || resolvedUser.username}`
}

function guideStateKey(user?: CurrentUser): string | null {
  const key = guideUserKey(user)
  return key ? `${GUIDE_KEY_PREFIX}:${key}` : null
}

function guideDismissedKey(user?: CurrentUser): string | null {
  const key = guideUserKey(user)
  return key ? `${GUIDE_DISMISSED_PREFIX}:${key}` : null
}

function readGuideState(user?: CurrentUser): GuideState {
  const key = guideStateKey(user)
  if (!key) return {}

  const raw = localStorage.getItem(key)
  if (!raw) return {}

  try {
    return JSON.parse(raw) as GuideState
  } catch {
    localStorage.removeItem(key)
    return {}
  }
}

function writeGuideState(state: GuideState, user?: CurrentUser): void {
  const key = guideStateKey(user)
  if (!key) return
  localStorage.setItem(key, JSON.stringify(state))
}

function clearGuideDismissal(user?: CurrentUser): void {
  const key = guideDismissedKey(user)
  if (key) sessionStorage.removeItem(key)
}

export function isGuideDone(user?: CurrentUser): boolean {
  const state = readGuideState(user)
  const dismissedKey = guideDismissedKey(user)
  const dismissedForSession = dismissedKey ? sessionStorage.getItem(dismissedKey) === "1" : false
  return state.completed === true || state.dontShowAgain === true || dismissedForSession
}

export function markGuideDone(user?: CurrentUser): void {
  writeGuideState({ ...readGuideState(user), completed: true }, user)
  clearGuideDismissal(user)
}

export function markGuideDontShowAgain(user?: CurrentUser): void {
  writeGuideState({ ...readGuideState(user), dontShowAgain: true }, user)
  clearGuideDismissal(user)
}

export function dismissGuideForNow(user?: CurrentUser): void {
  const key = guideDismissedKey(user)
  if (key) sessionStorage.setItem(key, "1")
}
