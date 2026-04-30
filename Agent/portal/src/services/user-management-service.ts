import {
  validatePasswordConfirmation,
  validatePasswordPolicy,
  validateRequired,
} from "@/lib/password-policy"
import { authApiConfig } from "@/config/auth"
import { translate } from "@/i18n/messages"
import { authHeader, readJson, realmPath, requestJson, requestNoContent } from "@/services/auth-api-client"
import { getAuthorizedSession } from "@/services/auth-service"
import { ServiceError } from "@/types/auth"
import type {
  CreateUserInput,
  ManagedUser,
  ResetPasswordInput,
  UpdateUserInput,
  UserStatus,
} from "@/types/user-management"

interface KeycloakUserRepresentation {
  id?: string
  username?: string
  firstName?: string
  lastName?: string
  email?: string
  enabled?: boolean
  createdTimestamp?: number
  requiredActions?: string[]
  attributes?: Record<string, string[] | string | undefined>
  credentials?: KeycloakCredentialRepresentation[]
}

interface KeycloakCredentialRepresentation {
  type: "password"
  value: string
  temporary?: boolean
}

interface KeycloakRoleRepresentation {
  id?: string
  name?: string
}

interface ListUsersQuery {
  username?: string
  exact?: boolean
}

const USERS_PATH = realmPath("/admin/realms/{realm}/users")
const ROLES_PATH = realmPath("/admin/realms/{realm}/roles")
const BUILTIN_ADMIN_ROLE_FALLBACK_NAMES = ["admin", "backup_admin"] as const

function userPath(userId: string): string {
  return `${USERS_PATH}/${encodeURIComponent(userId)}`
}

function resetPasswordPath(userId: string): string {
  return `${userPath(userId)}/reset-password`
}

function roleMappingsPath(userId: string): string {
  return `${userPath(userId)}/role-mappings/realm`
}

function rolePath(roleName: string): string {
  return `${ROLES_PATH}/${encodeURIComponent(roleName)}`
}

function jsonHeaders(accessToken: string): Record<string, string> {
  return {
    ...authHeader(accessToken),
    "Content-Type": "application/json",
  }
}

function firstAttribute(attributes: KeycloakUserRepresentation["attributes"], key: string): string | undefined {
  const value = attributes?.[key]
  if (Array.isArray(value)) return value[0]
  return value
}

function toIso(timestamp?: number): string {
  return timestamp ? new Date(timestamp).toISOString() : ""
}

function displayNameOf(user: KeycloakUserRepresentation): string {
  const fromAttribute = firstAttribute(user.attributes, "displayName")
  const fullName = [user.firstName, user.lastName].filter(Boolean).join(" ").trim()
  return fromAttribute || fullName || user.username || user.id || translate("users.unnamed")
}

function toManagedUser(user: KeycloakUserRepresentation): ManagedUser {
  const createdAt = toIso(user.createdTimestamp)
  return {
    id: user.id ?? user.username ?? "",
    username: user.username ?? "",
    displayName: displayNameOf(user),
    role: "backup_admin",
    status: user.enabled === false ? "disabled" : "enabled",
    remark: firstAttribute(user.attributes, "remark"),
    createdAt,
    updatedAt: firstAttribute(user.attributes, "updatedAt") ?? createdAt,
    lastLoginAt: firstAttribute(user.attributes, "lastLoginAt"),
  }
}

function withRemark(
  attributes: KeycloakUserRepresentation["attributes"] = {},
  remark?: string,
): KeycloakUserRepresentation["attributes"] {
  const next = { ...attributes }
  const normalized = remark?.trim()

  if (normalized) {
    next.remark = [normalized]
  } else {
    delete next.remark
  }

  return next
}

function createUserRepresentation(input: CreateUserInput): KeycloakUserRepresentation {
  return {
    username: input.username.trim(),
    enabled: input.status === "enabled",
    firstName: input.displayName.trim(),
    attributes: withRemark({}, input.remark),
    credentials: [
      {
        type: "password",
        value: input.password,
        temporary: false,
      },
    ],
  }
}

function updateUserRepresentation(
  user: KeycloakUserRepresentation,
  input: UpdateUserInput,
): KeycloakUserRepresentation {
  return {
    ...user,
    firstName: input.displayName.trim(),
    enabled: input.status === "enabled",
    attributes: withRemark(user.attributes, input.remark),
  }
}

function enabledUserRepresentation(
  user: KeycloakUserRepresentation,
  status: UserStatus,
): KeycloakUserRepresentation {
  return {
    ...user,
    enabled: status === "enabled",
  }
}

function createdUserId(response: Response): string | null {
  const location = response.headers.get("Location")
  if (!location) return null

  const match = location.match(/\/users\/([^/?#]+)$/)
  return match ? decodeURIComponent(match[1]) : null
}

function assertCreateInput(input: CreateUserInput): void {
  const usernameError = validateRequired(input.username, "username", translate("validation.requiredUsername"))
  if (usernameError) throw new ServiceError("VALIDATION_ERROR", usernameError.message, 400)

  const displayNameError = validateRequired(
    input.displayName,
    "displayName",
    translate("validation.requiredDisplayName"),
  )
  if (displayNameError) throw new ServiceError("VALIDATION_ERROR", displayNameError.message, 400)

  const passwordRequiredError = validateRequired(input.password, "password", translate("validation.requiredPassword"))
  if (passwordRequiredError) throw new ServiceError("VALIDATION_ERROR", passwordRequiredError.message, 400)

  const confirmPasswordRequiredError = validateRequired(
    input.confirmPassword,
    "confirmPassword",
    translate("validation.confirmPasswordRequired"),
  )
  if (confirmPasswordRequiredError) {
    throw new ServiceError("VALIDATION_ERROR", confirmPasswordRequiredError.message, 400)
  }

  const passwordError = validatePasswordPolicy(input.password, input.username)
  if (passwordError) throw new ServiceError("PASSWORD_POLICY_VIOLATION", passwordError, 400)

  const confirmationError = validatePasswordConfirmation(input.password, input.confirmPassword)
  if (confirmationError) throw new ServiceError("PASSWORD_CONFIRMATION_MISMATCH", confirmationError, 400)
}

function assertUpdateInput(input: UpdateUserInput): void {
  const displayNameError = validateRequired(
    input.displayName,
    "displayName",
    translate("validation.requiredDisplayName"),
  )
  if (displayNameError) throw new ServiceError("VALIDATION_ERROR", displayNameError.message, 400)
}

function assertResetPasswordInput(username: string, input: ResetPasswordInput): void {
  const passwordRequiredError = validateRequired(input.password, "password", translate("validation.requiredPassword"))
  if (passwordRequiredError) throw new ServiceError("VALIDATION_ERROR", passwordRequiredError.message, 400)

  const confirmPasswordRequiredError = validateRequired(
    input.confirmPassword,
    "confirmPassword",
    translate("validation.confirmPasswordRequired"),
  )
  if (confirmPasswordRequiredError) {
    throw new ServiceError("VALIDATION_ERROR", confirmPasswordRequiredError.message, 400)
  }

  const passwordError = validatePasswordPolicy(input.password, username)
  if (passwordError) throw new ServiceError("PASSWORD_POLICY_VIOLATION", passwordError, 400)

  const confirmationError = validatePasswordConfirmation(input.password, input.confirmPassword)
  if (confirmationError) throw new ServiceError("PASSWORD_CONFIRMATION_MISMATCH", confirmationError, 400)
}

async function authorizedHeaders(): Promise<Record<string, string>> {
  const session = await getAuthorizedSession()
  return authHeader(session.accessToken)
}

async function authorizedJsonHeaders(): Promise<Record<string, string>> {
  const session = await getAuthorizedSession()
  return jsonHeaders(session.accessToken)
}

function listUsersPath(query?: ListUsersQuery): string {
  const params = new URLSearchParams()
  params.set("briefRepresentation", "false")

  if (query?.username) params.set("username", query.username)
  if (query?.exact !== undefined) params.set("exact", String(query.exact))

  return `${USERS_PATH}?${params.toString()}`
}

export async function listUsers(query?: ListUsersQuery): Promise<ManagedUser[]> {
  const users = await requestJson<KeycloakUserRepresentation[]>(
    listUsersPath(query),
    {
      method: "GET",
      headers: await authorizedHeaders(),
    },
    translate("users.listLoadFailed"),
  )

  return users.map(toManagedUser)
}

async function getUserRepresentation(userId: string): Promise<KeycloakUserRepresentation> {
  return requestJson<KeycloakUserRepresentation>(
    userPath(userId),
    {
      method: "GET",
      headers: await authorizedHeaders(),
    },
    translate("users.detailLoadFailed"),
  )
}

function builtinAdminRoleNames(): string[] {
  return [...new Set([authApiConfig.adminRoleName, ...BUILTIN_ADMIN_ROLE_FALLBACK_NAMES].map((name) => name.trim()).filter(Boolean))]
}

async function getRoleRepresentation(roleName: string): Promise<KeycloakRoleRepresentation> {
  return requestJson<KeycloakRoleRepresentation>(
    rolePath(roleName),
    {
      method: "GET",
      headers: await authorizedHeaders(),
    },
    translate("users.roleLoadFailed"),
  )
}

async function resolveBuiltinAdminRole(): Promise<Required<KeycloakRoleRepresentation>> {
  let lastTraceId: string | undefined

  for (const roleName of builtinAdminRoleNames()) {
    try {
      const role = await getRoleRepresentation(roleName)
      if (role.id && role.name) {
        return { id: role.id, name: role.name }
      }
    } catch (error) {
      if (error instanceof ServiceError && error.status === 404) {
        lastTraceId = error.traceId
        continue
      }
      throw error
    }
  }

  throw new ServiceError(
    "ADMIN_ROLE_NOT_FOUND",
    translate("users.adminRoleNotFound").replace("{{roles}}", builtinAdminRoleNames().join(" / ")),
    404,
    lastTraceId,
  )
}

async function assignBuiltinAdminRole(userId: string): Promise<void> {
  const role = await resolveBuiltinAdminRole()

  await requestNoContent(
    roleMappingsPath(userId),
    {
      method: "POST",
      headers: await authorizedJsonHeaders(),
      body: JSON.stringify([
        {
          id: role.id,
          name: role.name,
        },
      ]),
    },
    translate("users.assignRoleFailed"),
  )
}

export async function getUser(userId: string): Promise<ManagedUser> {
  return toManagedUser(await getUserRepresentation(userId))
}

async function getCreatedUser(response: Response, username: string): Promise<ManagedUser> {
  const userId = createdUserId(response)
  if (userId) return getUser(userId)

  const [user] = await listUsers({ username, exact: true })
  if (!user) {
    throw new ServiceError(
      "USER_NOT_FOUND",
      translate("users.createReadFailed"),
      404,
    )
  }

  return user
}

export async function createUser(input: CreateUserInput): Promise<ManagedUser> {
  assertCreateInput(input)

  try {
    const response = await requestNoContent(
      USERS_PATH,
      {
        method: "POST",
        headers: await authorizedJsonHeaders(),
        body: JSON.stringify(createUserRepresentation(input)),
      },
      translate("users.createFailed"),
    )

    if (response.status !== 204) {
      const contentType = response.headers.get("Content-Type") ?? ""
      if (contentType.includes("application/json")) {
        const body = await readJson<KeycloakUserRepresentation>(response)
        if (body && body.id) {
          await assignBuiltinAdminRole(body.id)
          return toManagedUser(body)
        }
      }
    }

    const createdUser = await getCreatedUser(response, input.username.trim())

    try {
      await assignBuiltinAdminRole(createdUser.id)
    } catch (error) {
      if (error instanceof ServiceError) {
        throw new ServiceError(
          "ADMIN_ROLE_ASSIGNMENT_FAILED",
          translate("users.assignAfterCreateFailed"),
          error.status,
          error.traceId,
        )
      }
      throw error
    }

    return createdUser
  } catch (error) {
    if (error instanceof ServiceError && error.status === 409) {
      throw new ServiceError(
        "USERNAME_EXISTS",
        translate("users.usernameExists"),
        error.status,
        error.traceId,
      )
    }
    throw error
  }
}

export async function updateUser(userId: string, input: UpdateUserInput): Promise<ManagedUser> {
  assertUpdateInput(input)
  const headers = await authorizedJsonHeaders()
  const user = await getUserRepresentation(userId)

  await requestNoContent(
    userPath(userId),
    {
      method: "PUT",
      headers,
      body: JSON.stringify(updateUserRepresentation(user, input)),
    },
    translate("users.saveFailed"),
  )

  return getUser(userId)
}

async function setUserStatus(userId: string, status: UserStatus): Promise<ManagedUser> {
  const headers = await authorizedJsonHeaders()
  const user = await getUserRepresentation(userId)

  await requestNoContent(
    userPath(userId),
    {
      method: "PUT",
      headers,
      body: JSON.stringify(enabledUserRepresentation(user, status)),
    },
    translate("users.statusSaveFailed"),
  )

  return getUser(userId)
}

export async function enableUser(userId: string): Promise<ManagedUser> {
  return setUserStatus(userId, "enabled")
}

export async function disableUser(userId: string, currentUserId: string): Promise<ManagedUser> {
  if (userId === currentUserId) {
    throw new ServiceError(
      "CANNOT_DISABLE_SELF",
      translate("users.cannotDisableSelf"),
      400,
    )
  }

  return setUserStatus(userId, "disabled")
}

export async function resetUserPassword(userId: string, input: ResetPasswordInput): Promise<ManagedUser> {
  const headers = await authorizedJsonHeaders()
  const user = await getUserRepresentation(userId)
  assertResetPasswordInput(user.username ?? "", input)

  const credential: KeycloakCredentialRepresentation = {
    type: "password",
    value: input.password,
    temporary: false,
  }

  await requestNoContent(
    resetPasswordPath(userId),
    {
      method: "PUT",
      headers,
      body: JSON.stringify(credential),
    },
    translate("users.resetPasswordFailed"),
  )

  return getUser(userId)
}
