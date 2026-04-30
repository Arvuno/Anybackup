import { beforeEach, describe, expect, it, vi } from "vitest"
import { clearStoredSession, getStoredSession, isAuthenticated, storeSession } from "@/lib/session"
import { getCurrentUser, login, logoutCurrentUser } from "@/services/auth-service"
import type { AuthSession } from "@/types/auth"

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { "Content-Type": "application/json" },
  })
}

function emptyResponse(status = 204): Response {
  return new Response(null, { status })
}

function tokenResponse() {
  return {
    access_token: "access-token",
    expires_in: 300,
    refresh_expires_in: 1800,
    refresh_token: "refresh-token",
    token_type: "Bearer",
    scope: "openid profile",
  }
}

function userInfoResponse() {
  return {
    sub: "user-001",
    preferred_username: "admin",
    name: "备份管理员",
    roles: ["backup_admin"],
  }
}

function storeApiSession(overrides: Partial<AuthSession> = {}) {
  const now = Date.now()
  storeSession({
    accessToken: "access-token",
    refreshToken: "refresh-token",
    tokenType: "Bearer",
    expiresAt: now + 300_000,
    refreshExpiresAt: now + 1_800_000,
    idleExpiresAt: now + 1_800_000,
    user: {
      id: "user-001",
      username: "admin",
      displayName: "备份管理员",
      role: "backup_admin",
      tenantId: "master",
    },
    ...overrides,
  } as AuthSession)
}

describe("auth service", () => {
  const fetchMock = vi.fn<typeof fetch>()

  beforeEach(() => {
    vi.restoreAllMocks()
    fetchMock.mockReset()
    vi.stubGlobal("fetch", fetchMock)
    localStorage.clear()
  })

  it("logs in through the Auth Service token and userinfo contracts", async () => {
    fetchMock.mockResolvedValueOnce(jsonResponse(tokenResponse()))
    fetchMock.mockResolvedValueOnce(jsonResponse(userInfoResponse()))

    const session = await login({ username: "admin", password: "admin1234" })

    expect(session.user.displayName).toBe("备份管理员")
    expect(session.user.role).toBe("backup_admin")
    expect(session.accessToken).toBe("access-token")
    expect(session.refreshToken).toBe("refresh-token")
    expect(isAuthenticated()).toBe(true)
    expect(getStoredSession()?.user.username).toBe("admin")

    const [tokenUrl, tokenInit] = fetchMock.mock.calls[0]
    expect(tokenUrl).toBe("/api/auth_service/v1/realms/master/protocol/openid-connect/token")
    expect(tokenInit?.method).toBe("POST")
    expect(tokenInit?.headers).toEqual({ "Content-Type": "application/x-www-form-urlencoded" })
    const tokenBody = tokenInit?.body as URLSearchParams
    expect(tokenBody.get("grant_type")).toBe("password")
    expect(tokenBody.get("client_id")).toBe("admin-cli")
    expect(tokenBody.get("scope")).toBe("openid")
    expect(tokenBody.get("username")).toBe("admin")

    const [, userInfoInit] = fetchMock.mock.calls[1]
    expect(userInfoInit?.headers).toMatchObject({ Authorization: "Bearer access-token" })
  })

  it("rejects invalid credentials", async () => {
    fetchMock.mockResolvedValueOnce(jsonResponse({ error: "invalid_grant" }, 401))

    await expect(login({ username: "admin", password: "wrong1234" })).rejects.toMatchObject({
      code: "INVALID_CREDENTIALS",
      message: "Invalid username or password.",
    })
  })

  it("returns current user from a valid session", async () => {
    storeApiSession()
    fetchMock.mockResolvedValueOnce(jsonResponse(userInfoResponse()))

    await expect(getCurrentUser()).resolves.toMatchObject({ username: "admin" })

    const [url, init] = fetchMock.mock.calls[0]
    expect(url).toBe("/api/auth_service/v1/realms/master/protocol/openid-connect/userinfo")
    expect(init?.headers).toMatchObject({ Authorization: "Bearer access-token" })
  })

  it("clears session on logout", async () => {
    storeApiSession()
    fetchMock.mockResolvedValueOnce(emptyResponse())

    await logoutCurrentUser()

    expect(getStoredSession()).toBeNull()
    expect(isAuthenticated()).toBe(false)

    const [url, init] = fetchMock.mock.calls[0]
    expect(url).toBe("/api/auth_service/v1/realms/master/protocol/openid-connect/logout")
    expect(init?.method).toBe("POST")
    const body = init?.body as URLSearchParams
    expect(body.get("client_id")).toBe("admin-cli")
    expect(body.get("refresh_token")).toBe("refresh-token")
  })

  it("throws session error when stored session is missing", async () => {
    clearStoredSession()
    await expect(getCurrentUser()).rejects.toMatchObject({ code: "UNAUTHORIZED" })
  })
})
