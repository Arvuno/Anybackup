import { LANGUAGE_SWITCH_ENABLED, useI18n } from "@/i18n"

export function LanguageSwitch() {
  const { locale, setLocale, t } = useI18n()

  if (!LANGUAGE_SWITCH_ENABLED) return null

  return (
    <div className="inline-flex items-center gap-2 rounded-md border border-border bg-card p-1 text-xs">
      <button
        type="button"
        onClick={() => setLocale("en")}
        className={locale === "en" ? "rounded bg-secondary px-2 py-1 text-foreground" : "px-2 py-1 text-muted-foreground"}
      >
        EN
      </button>
      <button
        type="button"
        onClick={() => setLocale("zh-CN")}
        className={
          locale === "zh-CN" ? "rounded bg-secondary px-2 py-1 text-foreground" : "px-2 py-1 text-muted-foreground"
        }
      >
        {t("language.zhShort")}
      </button>
    </div>
  )
}
