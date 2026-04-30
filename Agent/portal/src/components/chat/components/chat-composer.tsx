import { useEffect, useRef } from "react"
import { SendHorizontal } from "lucide-react"
import { Button } from "@/components/ui/button"
import { useI18n } from "@/i18n"
import { cn } from "@/lib/cn"

interface ChatComposerProps {
  value: string
  onChange: (value: string) => void
  onSubmit?: () => void
  disabled?: boolean
  centered?: boolean
  pending?: boolean
  showHint?: boolean
}

export function ChatComposer({
  value,
  onChange,
  onSubmit,
  disabled = false,
  centered = false,
  pending = false,
  showHint = !centered,
}: ChatComposerProps) {
  const { t } = useI18n()
  const textareaRef = useRef<HTMLTextAreaElement>(null)
  const submitDisabled = disabled || pending || !value.trim()

  useEffect(() => {
    const element = textareaRef.current
    if (!element) return

    element.style.height = "0px"
    element.style.height = `${Math.min(element.scrollHeight, 160)}px`
  }, [value])

  return (
    <div className={cn("shrink-0 px-4", centered ? "bg-transparent pb-0" : "bg-background pb-5")}>
      <div className={cn("mx-auto", centered ? "max-w-2xl" : "max-w-3xl")}>
        <div
          className={cn(
            "flex items-center gap-2.5 rounded-xl border border-border/60 px-3 py-1.5 transition-fast focus-within:border-ai/40 focus-within:ring-2 focus-within:ring-ai/20",
            centered
              ? "bg-white/92 shadow-[0_16px_40px_-20px_rgba(16,24,40,0.24)] backdrop-blur-sm"
              : "bg-card/90 shadow-card",
          )}
        >
          <textarea
            ref={textareaRef}
            aria-label={t("chat.composer.inputLabel")}
            className={cn(
              "max-h-40 flex-1 resize-none border-0 bg-transparent px-1 text-sm leading-5 text-foreground placeholder:text-muted-foreground/70 focus:outline-none",
              centered ? "min-h-9 py-[7px]" : "min-h-9 py-[7px]",
            )}
            placeholder={t("chat.composer.placeholder")}
            rows={1}
            value={value}
            onChange={(event) => onChange(event.target.value)}
            onKeyDown={(event) => {
              if (event.key === "Enter" && !event.shiftKey) {
                event.preventDefault()
                if (!submitDisabled) onSubmit?.()
              }
            }}
          />
          <Button
            type="button"
            aria-label={t("chat.composer.send")}
            title={t("chat.composer.send")}
            variant="ai"
            size="icon"
            disabled={submitDisabled}
            className="h-8 w-8 shrink-0 rounded-lg focus-visible:ring-ai/30"
            onClick={() => onSubmit?.()}
          >
            <SendHorizontal className="h-3.5 w-3.5" />
          </Button>
        </div>
        {showHint ? (
          <p className="mt-1.5 text-center text-[10px] text-muted-foreground/40">
            {t("chat.composer.hint")}
          </p>
        ) : null}
      </div>
    </div>
  )
}
