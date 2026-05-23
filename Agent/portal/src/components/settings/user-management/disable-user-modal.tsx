import { AlertTriangle } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Modal } from "@/components/ui/modal"
import { useI18n } from "@/i18n"
import type { ManagedUser } from "@/types/user-management"

interface DisableUserModalProps {
  open: boolean
  user?: ManagedUser
  saving: boolean
  onClose: () => void
  onConfirm: (userId: string) => Promise<void>
}

export function DisableUserModal({ open, user, saving, onClose, onConfirm }: DisableUserModalProps) {
  const { t } = useI18n()

  return (
    <Modal
      open={open}
      title={t("disableUser.title")}
      description={t("disableUser.description")}
      onClose={onClose}
      footer={
        <>
          <Button type="button" variant="ghost" onClick={onClose}>
            {t("disableUser.cancel")}
          </Button>
          <Button
            type="button"
            variant="destructive"
            disabled={saving || !user}
            onClick={async () => {
              if (!user) return
              await onConfirm(user.id)
              onClose()
            }}
          >
            {t("disableUser.confirm")}
          </Button>
        </>
      }
    >
      <div className="space-y-4">
        <div className="rounded-lg border border-destructive/15 bg-destructive/5 px-4 py-3">
          <div className="flex gap-3">
            <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-destructive/10 text-destructive">
              <AlertTriangle className="h-4 w-4" />
            </div>
            <div>
              <p className="text-sm font-medium text-foreground">{t("disableUser.warningTitle")}</p>
              <p className="mt-1 text-xs leading-5 text-muted-foreground">{t("disableUser.warningBody")}</p>
            </div>
          </div>
        </div>

        <div className="rounded-lg border border-border/70 bg-muted/25 px-4 py-3 text-sm text-foreground">
          {user ? t("disableUser.preview").replace("{{username}}", user.username) : t("disableUser.selectUser")}
        </div>
      </div>
    </Modal>
  )
}
