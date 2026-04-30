function cleanBasePath(value: string): string {
  const normalized = value.trim() || "/api/auth_service/v1"
  return normalized.endsWith("/") ? normalized.slice(0, -1) : normalized
}

export const authApiConfig = {
  basePath: cleanBasePath(import.meta.env.VITE_AUTH_SERVICE_BASE_PATH ?? "/api/auth_service/v1"),
  realm: import.meta.env.VITE_AUTH_REALM ?? "master",
  clientId: import.meta.env.VITE_AUTH_CLIENT_ID ?? "admin-cli",
  scope: import.meta.env.VITE_AUTH_SCOPE ?? "openid",
  adminRoleName: import.meta.env.VITE_AUTH_ADMIN_ROLE_NAME ?? "backup_admin",
}
