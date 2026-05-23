import { useEffect, useState } from "react"
import { useI18n } from "@/i18n"
import { Button } from "@/components/ui/button"
import { Modal } from "@/components/ui/modal"
import { TextField } from "@/components/ui/text-field"
import type { ManagedUser, ResetPasswordInput } from "@/types/user-management"

interface ResetPasswordModalProps {
  open: boolean
  user?: ManagedUser
  saving: boolean
  onClose: () => void
  onReset: (userId: string, input: ResetPasswordInput) => Promise<void>
}

export function ResetPasswordModal({ open, user, saving, onClose, onReset }: ResetPasswordModalProps) {
  const { t } = useI18n()
  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")

  useEffect(() => {
    if (!open) return
    setPassword("")
    setConfirmPassword("")
  }, [open])

  const submit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (!user) return
    await onReset(user.id, { password, confirmPassword })
    onClose()
  }

  return (
    <Modal
      open={open}
      title={t("resetPassword.title")}
      description={
        user ? t("resetPassword.description").replace("{{username}}", user.username) : undefined
      }
      onClose={onClose}
      footer={
        <>
          <Button type="button" variant="ghost" onClick={onClose}>
            {t("resetPassword.cancel")}
          </Button>
          <Button type="submit" form="reset-password-form" disabled={saving}>
            {t("resetPassword.submit")}
          </Button>
        </>
      }
    >
      <form id="reset-password-form" className="space-y-5" onSubmit={submit}>
        {user ? (
          <div className="rounded-lg border border-border/70 bg-muted/25 px-3 py-3">
            <p className="text-[11px] text-muted-foreground">{t("resetPassword.targetLabel")}</p>
            <p className="mt-1 text-sm font-medium text-foreground">{user.displayName || user.username}</p>
            <p className="text-xs text-muted-foreground">{user.username}</p>
          </div>
        ) : null}

        <div className="grid gap-4 sm:grid-cols-2">
          <TextField
            label={t("resetPassword.newPassword")}
            name="newPassword"
            type="password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
          />
          <TextField
            label={t("resetPassword.confirmNewPassword")}
            name="confirmNewPassword"
            type="password"
            value={confirmPassword}
            onChange={(event) => setConfirmPassword(event.target.value)}
          />
        </div>
      </form>
    </Modal>
  )
}
