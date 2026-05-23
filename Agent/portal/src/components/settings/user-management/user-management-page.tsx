import { useEffect, useMemo, useState } from "react"
import { Pencil, Plus, RotateCcw, UserRoundCheck, UserRoundX } from "lucide-react"
import { Button } from "@/components/ui/button"
import { UserFormModal } from "@/components/settings/user-management/user-form-modal"
import { ResetPasswordModal } from "@/components/settings/user-management/reset-password-modal"
import { DisableUserModal } from "@/components/settings/user-management/disable-user-modal"
import { UserStatusBadge } from "@/components/settings/user-management/user-status-badge"
import { useAuthStore } from "@/store/useAuthStore"
import { useUserManagementStore } from "@/store/useUserManagementStore"
import { useI18n } from "@/i18n"

const FEEDBACK_DISMISS_DELAY_MS = 4000

export function UserManagementPageContent() {
  const { t } = useI18n()
  const currentUser = useAuthStore((state) => state.currentUser)
  const users = useUserManagementStore((state) => state.users)
  const loading = useUserManagementStore((state) => state.loading)
  const saving = useUserManagementStore((state) => state.saving)
  const message = useUserManagementStore((state) => state.message)
  const error = useUserManagementStore((state) => state.error)
  const feedbackAutoDismiss = useUserManagementStore((state) => state.feedbackAutoDismiss)
  const loadUsers = useUserManagementStore((state) => state.loadUsers)
  const create = useUserManagementStore((state) => state.create)
  const update = useUserManagementStore((state) => state.update)
  const enable = useUserManagementStore((state) => state.enable)
  const disable = useUserManagementStore((state) => state.disable)
  const resetPassword = useUserManagementStore((state) => state.resetPassword)
  const clearFeedback = useUserManagementStore((state) => state.clearFeedback)

  const [formMode, setFormMode] = useState<"create" | "edit" | null>(null)
  const [selectedUserId, setSelectedUserId] = useState<string | null>(null)
  const [resetUserId, setResetUserId] = useState<string | null>(null)
  const [disableUserId, setDisableUserId] = useState<string | null>(null)

  useEffect(() => {
    void loadUsers()
  }, [loadUsers])

  useEffect(() => {
    if (!feedbackAutoDismiss || (!message && !error)) return

    const timeoutId = window.setTimeout(() => {
      clearFeedback()
    }, FEEDBACK_DISMISS_DELAY_MS)

    return () => window.clearTimeout(timeoutId)
  }, [clearFeedback, error, feedbackAutoDismiss, message])

  const selectedUser = useMemo(
    () => users.find((user) => user.id === selectedUserId),
    [selectedUserId, users],
  )
  const resetUser = useMemo(() => users.find((user) => user.id === resetUserId), [resetUserId, users])
  const disablingUser = useMemo(() => users.find((user) => user.id === disableUserId), [disableUserId, users])
  const enabledCount = useMemo(() => users.filter((user) => user.status === "enabled").length, [users])

  const onDisableClick = (userId: string) => {
    if (userId === currentUser?.id) {
      void disable(userId, currentUser.id).catch(() => undefined)
      return
    }
    setDisableUserId(userId)
  }

  return (
    <div className="h-full min-h-0 flex-1 overflow-auto p-6">
      <div className="mx-auto flex max-w-6xl flex-col gap-5">
        <header className="rounded-xl border border-border/70 bg-card/90 px-5 py-4 shadow-card">
          <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
            <div>
              <h1 className="text-[22px] font-semibold tracking-tight text-foreground">{t("usersPage.title")}</h1>
              <p className="mt-1 text-sm text-muted-foreground">{t("usersPage.subtitle")}</p>
            </div>

            <div className="flex items-center gap-2 sm:pt-1">
              <div className="hidden rounded-full border border-border/70 bg-white/80 px-3 py-1 text-[11px] font-medium text-muted-foreground sm:inline-flex">
                {t("usersPage.enabledCount")
                  .replace("{enabled}", String(enabledCount))
                  .replace("{total}", String(users.length || 0))}
              </div>
              <Button
                type="button"
                size="sm"
                className="gap-2 rounded-md"
                onClick={() => {
                  setSelectedUserId(null)
                  setFormMode("create")
                }}
              >
                <Plus className="h-3.5 w-3.5" />
                {t("usersPage.createUser")}
              </Button>
            </div>
          </div>
        </header>

        {message ? (
          <div
            role="status"
            aria-live="polite"
            className="rounded-lg border border-success/15 bg-success-surface/85 px-4 py-3 text-sm text-success shadow-[var(--shadow-xs)]"
          >
            {message}
          </div>
        ) : null}
        {error ? (
          <div
            role="alert"
            className="rounded-lg border border-destructive/15 bg-destructive/5 px-4 py-3 text-sm text-destructive shadow-[var(--shadow-xs)]"
          >
            {error}
          </div>
        ) : null}

        <section className="overflow-hidden rounded-xl border border-border/80 bg-card/95 shadow-card">
          <div className="flex flex-col gap-3 border-b border-border/70 bg-gradient-to-r from-ai-surface/80 via-white to-white px-4 py-3 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <h2 className="text-sm font-semibold text-foreground">{t("usersPage.listTitle")}</h2>
              <p className="mt-1 text-xs text-muted-foreground">{t("usersPage.listHint")}</p>
            </div>
            <div className="inline-flex w-fit rounded-full border border-border/70 bg-white/85 px-3 py-1 text-[11px] font-medium text-muted-foreground">
              {t("usersPage.userCount").replace("{count}", String(users.length))}
            </div>
          </div>

          <div className="overflow-x-auto">
            <div className="min-w-[900px]">
              <div className="grid grid-cols-[1.1fr_1.2fr_0.9fr_0.8fr_1.5fr] items-center gap-3 border-b border-border bg-muted/30 px-4 py-3 text-center text-[11px] font-medium text-muted-foreground">
                <span>{t("usersPage.colUsername")}</span>
                <span>{t("usersPage.colDisplayName")}</span>
                <span>{t("usersPage.colRole")}</span>
                <span>{t("usersPage.colStatus")}</span>
                <span>{t("usersPage.colActions")}</span>
              </div>

              {loading ? (
                <div className="px-4 py-12 text-center text-sm text-muted-foreground">{t("usersPage.loading")}</div>
              ) : users.length === 0 ? (
                <div className="px-4 py-12 text-center text-sm text-muted-foreground">{t("usersPage.empty")}</div>
              ) : (
                users.map((user) => {
                  const isCurrentUser = user.id === currentUser?.id

                  return (
                    <div
                      key={user.id}
                      className="grid grid-cols-[1.1fr_1.2fr_0.9fr_0.8fr_1.5fr] items-center gap-3 border-b border-border/70 px-4 py-3 text-center transition-colors last:border-b-0 hover:bg-muted/15"
                    >
                      <div className="min-w-0">
                        <span className="block truncate text-sm font-medium text-foreground">{user.username}</span>
                        {isCurrentUser ? (
                          <span className="mt-1 inline-flex rounded-full border border-ai/15 bg-ai-surface px-2 py-0.5 text-[10px] font-medium text-ai">
                            {t("usersPage.currentAccount")}
                          </span>
                        ) : null}
                      </div>

                      <span className="block min-w-0 truncate text-sm text-foreground">{user.displayName}</span>
                      <span className="text-sm text-muted-foreground">{t("sidebar.backupAdmin")}</span>
                      <div className="flex justify-center">
                        <UserStatusBadge status={user.status} />
                      </div>
                      <div className="flex flex-wrap justify-center gap-2">
                        <Button
                          type="button"
                          variant="ghost"
                          size="xs"
                          className="gap-1.5 rounded-sm"
                          onClick={() => {
                            setSelectedUserId(user.id)
                            setFormMode("edit")
                          }}
                        >
                          <Pencil className="h-3.5 w-3.5" />
                          {t("usersPage.edit")}
                        </Button>
                        {user.status === "enabled" ? (
                          <Button
                            type="button"
                            variant="ghost"
                            size="xs"
                            className="gap-1.5 rounded-sm text-destructive hover:text-destructive"
                            aria-label={`${t("usersPage.disable")} ${user.username}`}
                            onClick={() => onDisableClick(user.id)}
                          >
                            <UserRoundX className="h-3.5 w-3.5" />
                            {t("usersPage.disable")}
                          </Button>
                        ) : (
                          <Button
                            type="button"
                            variant="ghost"
                            size="xs"
                            className="gap-1.5 rounded-sm"
                            onClick={() => void enable(user.id).catch(() => undefined)}
                          >
                            <UserRoundCheck className="h-3.5 w-3.5" />
                            {t("usersPage.enable")}
                          </Button>
                        )}
                        <Button
                          type="button"
                          variant="ghost"
                          size="xs"
                          className="gap-1.5 rounded-sm"
                          onClick={() => setResetUserId(user.id)}
                        >
                          <RotateCcw className="h-3.5 w-3.5" />
                          {t("usersPage.resetPassword")}
                        </Button>
                      </div>
                    </div>
                  )
                })
              )}
            </div>
          </div>
        </section>
      </div>

      <UserFormModal
        open={formMode !== null}
        mode={formMode ?? "create"}
        user={selectedUser}
        saving={saving}
        onClose={() => setFormMode(null)}
        onCreate={create}
        onUpdate={update}
      />
      <ResetPasswordModal
        open={Boolean(resetUserId)}
        user={resetUser}
        saving={saving}
        onClose={() => setResetUserId(null)}
        onReset={resetPassword}
      />
      <DisableUserModal
        open={Boolean(disableUserId)}
        user={disablingUser}
        saving={saving}
        onClose={() => setDisableUserId(null)}
        onConfirm={(userId) => disable(userId, currentUser?.id ?? "")}
      />
    </div>
  )
}
