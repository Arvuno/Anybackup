import type { ReactNode } from "react"
import { Bot } from "lucide-react"
import { useI18n } from "@/i18n"

interface ConversationEmptyStateProps {
  composer: ReactNode
}

export function ConversationEmptyState({ composer }: ConversationEmptyStateProps) {
  const { t } = useI18n()
  return (
    <div className="relative grid min-h-0 flex-1 place-items-center overflow-hidden px-6">
      <div className="flex w-full max-w-2xl flex-col items-center justify-center gap-5">
        <div className="flex h-14 w-14 items-center justify-center rounded-xl bg-[hsl(var(--ai))] text-primary-foreground shadow-ai">
          <Bot className="h-7 w-7" />
        </div>

        <div className="text-center">
          <h1 className="text-lg font-bold text-foreground">{t("chat.empty.title")}</h1>
          <p className="mx-auto mt-1.5 max-w-sm text-sm leading-relaxed text-muted-foreground">
            {t("chat.empty.description")}
          </p>
        </div>

        <div className="w-full">{composer}</div>
      </div>
    </div>
  )
}
