import { Bot } from "lucide-react"
import { useI18n } from "@/i18n"

export function ChatPanelHeader() {
  const { t } = useI18n()
  return (
    <div className="flex items-center gap-2 border-b border-border/40 bg-gradient-ai-subtle px-4 py-3">
      <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-[hsl(var(--ai))] shadow-ai">
        <Bot className="h-4 w-4 text-primary-foreground" />
      </div>
      <div>
        <h2 className="text-sm font-semibold text-foreground">{t("chat.panelHeader.title")}</h2>
        <p className="text-[10px] text-muted-foreground">{t("chat.panelHeader.subtitle")}</p>
      </div>
    </div>
  )
}
