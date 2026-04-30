import { beforeEach, describe, expect, it, vi } from "vitest"
import { storeSession } from "@/lib/session"
import {
  createUser,
  disableUser,
  enableUser,
  listUsers,
  resetUserPassword,
  updateUser,
} from "@/services/user-management-service"
import type { AuthSession } from "@/types/auth"

function jsonResponse(body: unknown, status = 200, headers?: HeadersInit): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { "Content-Type": "application/json", ...headers },
  })
}

function emptyResponse(status = 204, headers?: HeadersInit): Response {
  return new Response(null, { status, headers })
}

function keycloakUser(overrides: Record<string, unknown> = {}) {
  return {
    id: "user-001",
    username: "admin",
    firstName: "备份管理员",
    enabled: true,
    createdTimestamp: 1_713_600_000_000,
    attributes: {
      remark: ["系统内置管理员"],
    },
    ...overrides,
  }
}

function storeApiSession() {
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
  } as AuthSession)
}

describe("user management service", () => {
  const fetchMock = vi.fn<typeof fetch>()

  beforeEach(() => {
    vi.restoreAllMocks()
    fetchMock.mockReset()
    vi.stubGlobal("fetch", fetchMock)
    localStorage.clear()
    storeApiSession()
  })

  it("lists users through the Auth Service Keycloak users contract", async () => {
    fetchMock.mockResolvedValueOnce(jsonResponse([keycloakUser()]))

    const users = await listUsers()

    expect(users[0]).toMatchObject({ username: "admin", role: "backup_admin", status: "enabled" })
    expect(users[0]).not.toHaveProperty("password")

    const [url, init] = fetchMock.mock.calls[0]
    expect(url).toBe("/api/auth_service/v1/admin/realms/master/users?briefRepresentation=false")
    expect(init?.headers).toMatchObject({ Authorization: "Bearer access-token" })
  })

  it("creates a backup admin user and assigns the builtin admin role", async () => {
    fetchMock.mockResolvedValueOnce(
      emptyResponse(201, {
        Location: "/api/auth_service/v1/admin/realms/master/users/user-002",
      }),
    )
    fetchMock.mockResolvedValueOnce(
      jsonResponse(
        keycloakUser({
          id: "user-002",
          username: "operator",
          firstName: "备份操作员",
          attributes: { remark: ["MVP 用户"] },
        }),
      ),
    )
    fetchMock.mockResolvedValueOnce(
      jsonResponse({
        id: "role-001",
        name: "backup_admin",
      }),
    )
    fetchMock.mockResolvedValueOnce(emptyResponse(204))

    const user = await createUser({
      username: "operator",
      displayName: "备份操作员",
      password: "operator123",
      confirmPassword: "operator123",
      status: "enabled",
      remark: "MVP 用户",
    })

    expect(user.role).toBe("backup_admin")
    expect(user.username).toBe("operator")

    const [, createInit] = fetchMock.mock.calls[0]
    expect(createInit?.method).toBe("POST")
    expect(JSON.parse(String(createInit?.body))).toMatchObject({
      username: "operator",
      enabled: true,
      firstName: "备份操作员",
      attributes: { remark: ["MVP 用户"] },
      credentials: [{ type: "password", value: "operator123", temporary: false }],
    })

    const [roleUrl] = fetchMock.mock.calls[2]
    expect(roleUrl).toBe("/api/auth_service/v1/admin/realms/master/roles/backup_admin")

    const [mappingUrl, mappingInit] = fetchMock.mock.calls[3]
    expect(mappingUrl).toBe("/api/auth_service/v1/admin/realms/master/users/user-002/role-mappings/realm")
    expect(mappingInit?.method).toBe("POST")
    expect(JSON.parse(String(mappingInit?.body))).toEqual([
      {
        id: "role-001",
        name: "backup_admin",
      },
    ])
  })

  it("falls back to the legacy admin role name when the configured role is not found", async () => {
    fetchMock.mockResolvedValueOnce(
      emptyResponse(201, {
        Location: "/api/auth_service/v1/admin/realms/master/users/user-002",
      }),
    )
    fetchMock.mockResolvedValueOnce(
      jsonResponse(
        keycloakUser({
          id: "user-002",
          username: "operator",
          firstName: "澶囦唤鎿嶄綔鍛?",
        }),
      ),
    )
    fetchMock.mockResolvedValueOnce(jsonResponse({ message: "Role not found" }, 404))
    fetchMock.mockResolvedValueOnce(
      jsonResponse({
        id: "role-admin-001",
        name: "admin",
      }),
    )
    fetchMock.mockResolvedValueOnce(emptyResponse(204))

    await expect(
      createUser({
        username: "operator",
        displayName: "澶囦唤鎿嶄綔鍛?",
        password: "operator123",
        confirmPassword: "operator123",
        status: "enabled",
      }),
    ).resolves.toMatchObject({ username: "operator", role: "backup_admin" })

    const [configuredRoleUrl] = fetchMock.mock.calls[2]
    expect(configuredRoleUrl).toBe("/api/auth_service/v1/admin/realms/master/roles/backup_admin")

    const [fallbackRoleUrl] = fetchMock.mock.calls[3]
    expect(fallbackRoleUrl).toBe("/api/auth_service/v1/admin/realms/master/roles/admin")

    const [, mappingInit] = fetchMock.mock.calls[4]
    expect(JSON.parse(String(mappingInit?.body))).toEqual([
      {
        id: "role-admin-001",
        name: "admin",
      },
    ])
  })

  it("maps duplicate username errors", async () => {
    fetchMock.mockResolvedValueOnce(jsonResponse({ errorMessage: "User exists with same username" }, 409))

    await expect(
      createUser({
        username: "admin",
        displayName: "重复管理员",
        password: "admin5678",
        confirmPassword: "admin5678",
        status: "enabled",
      }),
    ).rejects.toMatchObject({ code: "USERNAME_EXISTS", message: "Username already exists." })
  })

  it("updates profile fields without changing username", async () => {
    fetchMock.mockResolvedValueOnce(jsonResponse(keycloakUser()))
    fetchMock.mockResolvedValueOnce(emptyResponse())
    fetchMock.mockResolvedValueOnce(
      jsonResponse(
        keycloakUser({
          firstName: "系统管理员",
          attributes: { remark: ["更新备注"] },
        }),
      ),
    )

    const updated = await updateUser("user-001", {
      displayName: "系统管理员",
      status: "enabled",
      remark: "更新备注",
    })

    expect(updated.username).toBe("admin")
    expect(updated.displayName).toBe("系统管理员")

    const [, init] = fetchMock.mock.calls[1]
    const body = JSON.parse(String(init?.body))
    expect(body.username).toBe("admin")
    expect(body.firstName).toBe("系统管理员")
    expect(body.attributes.remark).toEqual(["更新备注"])
  })

  it("does not allow disabling the current user", async () => {
    await expect(disableUser("user-001", "user-001")).rejects.toMatchObject({
      code: "CANNOT_DISABLE_SELF",
      message: "You cannot disable the currently signed-in user.",
    })
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it("disables and enables another user", async () => {
    fetchMock.mockResolvedValueOnce(jsonResponse(keycloakUser({ id: "user-002", username: "operator" })))
    fetchMock.mockResolvedValueOnce(emptyResponse())
    fetchMock.mockResolvedValueOnce(
      jsonResponse(keycloakUser({ id: "user-002", username: "operator", enabled: false })),
    )
    await expect(disableUser("user-002", "user-001")).resolves.toMatchObject({ status: "disabled" })

    fetchMock.mockResolvedValueOnce(
      jsonResponse(keycloakUser({ id: "user-002", username: "operator", enabled: false })),
    )
    fetchMock.mockResolvedValueOnce(emptyResponse())
    fetchMock.mockResolvedValueOnce(
      jsonResponse(keycloakUser({ id: "user-002", username: "operator", enabled: true })),
    )
    await expect(enableUser("user-002")).resolves.toMatchObject({ status: "enabled" })
  })

  it("resets password with policy validation", async () => {
    fetchMock.mockResolvedValueOnce(jsonResponse(keycloakUser({ id: "user-002", username: "operator" })))
    fetchMock.mockResolvedValueOnce(emptyResponse())
    fetchMock.mockResolvedValueOnce(jsonResponse(keycloakUser({ id: "user-002", username: "operator" })))

    await expect(
      resetUserPassword("user-002", {
        password: "newpass123",
        confirmPassword: "newpass123",
      }),
    ).resolves.toMatchObject({ id: "user-002" })

    const [, init] = fetchMock.mock.calls[1]
    expect(init?.method).toBe("PUT")
    expect(JSON.parse(String(init?.body))).toEqual({
      type: "password",
      value: "newpass123",
      temporary: false,
    })
  })
})
