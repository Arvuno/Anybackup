import type { UserStatus } from "@/types/user-management"
import { useI18n } from "@/i18n"
import { cn } from "@/lib/cn"

export function UserStatusBadge({ status }: { status: UserStatus }) {
  const { t } = useI18n()
  const enabled = status === "enabled"

  return (
    <span
      className={cn(
        "inline-flex items-center gap-1 rounded-md px-2 py-0.5 text-[11px] font-medium",
        enabled ? "bg-success-surface text-success" : "bg-muted text-muted-foreground",
      )}
    >
      <span className={enabled ? "status-online" : "status-pending"} />
      {enabled ? t("userForm.statusEnabled") : t("userForm.statusDisabled")}
    </span>
  )
}
