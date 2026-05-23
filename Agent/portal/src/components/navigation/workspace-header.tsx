import { useLocation, useNavigate } from "react-router-dom"
import { routes } from "@/config/routes"
import { useI18n } from "@/i18n"
import { ChevronRight } from "lucide-react"

export function WorkspaceHeader() {
  const { t } = useI18n()
  const location = useLocation()
  const navigate = useNavigate()

  const breadcrumbLabel: Record<string, string> = {
    [routes.home]: t("workspace.breadcrumb.home"),
    [routes.settings]: t("workspace.breadcrumb.settings"),
    [routes.users]: t("workspace.breadcrumb.users"),
  }

  const crumb = breadcrumbLabel[location.pathname] ?? t("workspace.breadcrumb.home")

  return (
    <header className="flex h-12 shrink-0 items-center px-5" style={{ backgroundColor: "#F3F5F7" }}>
      <nav className="flex items-center gap-2 text-sm">
        <button
          type="button"
          onClick={() => navigate(routes.home)}
          className="text-muted-foreground transition-colors hover:text-primary"
          title={t("workspace.navTitle")}
        >
          {t("workspace.navTitle")}
        </button>
        <ChevronRight className="h-4 w-4 text-muted-foreground/40" />
        <span className="font-medium text-foreground">{crumb}</span>
      </nav>
    </header>
  )
}
