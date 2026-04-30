import { Link } from "react-router-dom"
import { Users } from "lucide-react"
import { LanguageSwitch } from "@/components/settings/language-switch"
import { routes } from "@/config/routes"
import { useI18n } from "@/i18n"

export function SettingsPage() {
  const { t } = useI18n()

  return (
    <div className="h-full min-h-0 flex-1 overflow-auto p-6">
      <div className="mx-auto max-w-5xl space-y-5">
        <header className="flex items-start justify-between gap-4">
          <div>
            <h1 className="text-h2 font-semibold text-foreground">{t("settings.title")}</h1>
            <p className="mt-1 text-sm text-muted-foreground">{t("settings.description")}</p>
          </div>
          <LanguageSwitch />
        </header>

        <section className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <Link
            to={routes.users}
            className="group rounded-lg border border-border bg-card p-5 shadow-card transition-smooth hover:-translate-y-0.5 hover:shadow-card-hover"
          >
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-ai-surface text-ai">
              <Users className="h-5 w-5" />
            </div>
            <h2 className="mt-4 text-base font-semibold text-foreground">{t("settings.userManagementTitle")}</h2>
            <p className="mt-2 text-sm leading-6 text-muted-foreground">{t("settings.userManagementDescription")}</p>
          </Link>
        </section>
      </div>
    </div>
  )
}
