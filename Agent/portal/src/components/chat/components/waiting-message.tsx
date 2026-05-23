import { useEffect, useState } from "react"
import { LoaderCircle } from "lucide-react"
import { useI18n } from "@/i18n"

interface WaitingMessageProps {
  label?: string
  startedAtMs?: number
}

export function WaitingMessage({ label, startedAtMs }: WaitingMessageProps) {
  const { t } = useI18n()
  const [nowMs, setNowMs] = useState(() => Date.now())

  const thinkingLabel = label ?? t("chat.waiting.thinking")

  const formatElapsedThinkingTime = (elapsedMs: number): string => {
    const totalSeconds = Math.max(0, Math.floor(elapsedMs / 1000))
    const minutes = Math.floor(totalSeconds / 60)
    const seconds = totalSeconds % 60

    return t("chat.waiting.duration")
      .replace("{minutes}", String(minutes))
      .replace("{seconds}", String(seconds).padStart(2, "0"))
  }

  const resolvedLabel =
    typeof startedAtMs === "number" ? formatElapsedThinkingTime(nowMs - startedAtMs) : thinkingLabel

  useEffect(() => {
    if (typeof startedAtMs !== "number") return undefined

    const intervalId = window.setInterval(() => {
      setNowMs(Date.now())
    }, 1000)

    return () => window.clearInterval(intervalId)
  }, [startedAtMs])

  return (
    <article className="flex animate-fade-in-up justify-start">
      <div className="max-w-[78%]">
        <div className="bubble-ai flex items-center gap-3 border border-border/70 bg-card px-4 py-3 text-sm text-muted-foreground shadow-card">
          <LoaderCircle className="h-4 w-4 animate-spin" />
          <span>{resolvedLabel}</span>
        </div>
      </div>
    </article>
  )
}
