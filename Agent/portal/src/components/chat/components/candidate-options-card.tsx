import { useEffect, useMemo, useState } from "react"
import { Check, CheckCircle2, Circle, Sparkles } from "lucide-react"
import { Button } from "@/components/ui/button"
import { useI18n } from "@/i18n"
import { cn } from "@/lib/cn"
import type { CandidateOptionsRichPayload } from "@/types/conversation"

interface CandidateOptionsCardProps {
  payload: CandidateOptionsRichPayload
  submitting: boolean
  onSelect: (input: {
    reasoningTraceId: string
    candidateOptionId: string
    selection: "confirm" | "reject" | "revise"
    additionalConstraints?: string
  }) => void
}

function findAction(payload: CandidateOptionsRichPayload, type: "confirm" | "reject" | "revise") {
  return payload.actions.find((action) => action.type === type)
}

export function CandidateOptionsCard({ payload, submitting, onSelect }: CandidateOptionsCardProps) {
  const { t } = useI18n()
  const initialOptionId =
    payload.selectedOptionId ??
    payload.options.find((option) => option.recommended)?.optionId ??
    payload.options[0]?.optionId ??
    null

  const [selectedOptionId, setSelectedOptionId] = useState<string | null>(initialOptionId)
  const [optimisticLockedOptionId, setOptimisticLockedOptionId] = useState<string | null>(payload.selectedOptionId ?? null)

  const confirmAction = useMemo(() => findAction(payload, "confirm"), [payload])
  const lockedOptionId = optimisticLockedOptionId ?? payload.selectedOptionId ?? null
  const selectionLocked = payload.selectionLocked === true || lockedOptionId !== null
  const activeOptionId = lockedOptionId ?? selectedOptionId

  useEffect(() => {
    setSelectedOptionId(initialOptionId)
    setOptimisticLockedOptionId(payload.selectedOptionId ?? null)
  }, [initialOptionId, payload.reasoningTraceId, payload.selectedOptionId])

  const selectedOption = useMemo(
    () => payload.options.find((option) => option.optionId === activeOptionId) ?? payload.options[0],
    [activeOptionId, payload.options],
  )

  if (!selectedOption) {
    return (
      <section className="rounded-xl border border-border/70 bg-card p-3 text-sm text-muted-foreground shadow-card">
        {payload.summary ?? payload.title}
      </section>
    )
  }

  return (
    <section className="rounded-2xl border border-border/70 bg-card p-3 shadow-card sm:p-4">
      <header className="space-y-1.5">
        <p className="text-[10px] font-semibold uppercase tracking-[0.16em] text-ai/80">Candidate Options</p>
        <div className="space-y-1">
          <h3 className="text-sm font-semibold leading-6 text-foreground sm:text-[15px]">{payload.title}</h3>
          {payload.summary ? <p className="text-xs leading-5 text-muted-foreground">{payload.summary}</p> : null}
        </div>
      </header>

      <div
        data-testid="candidate-options-grid"
        className="mt-3 grid gap-3.5 [grid-template-columns:repeat(auto-fit,minmax(min(100%,16rem),1fr))]"
      >
        {payload.options.map((option) => {
          const isActive = option.optionId === selectedOption.optionId
          const disableOptionButton = submitting || selectionLocked

          return (
            <article
              key={option.optionId}
              className={cn(
                "flex h-full min-w-0 flex-col overflow-hidden rounded-2xl border bg-background/80 shadow-sm transition-all duration-200 ease-smooth",
                isActive
                  ? "border-ai/60 shadow-[0_12px_32px_-24px_rgba(10,132,255,0.55)] ring-1 ring-ai/25"
                  : "border-border/70",
                selectionLocked && !isActive ? "opacity-70" : "hover:border-ai/25 hover:bg-ai/5",
              )}
            >
              <button
                type="button"
                aria-pressed={isActive}
                disabled={disableOptionButton}
                className="flex h-full w-full min-w-0 flex-col text-left disabled:cursor-not-allowed disabled:opacity-90"
                onClick={() => setSelectedOptionId(option.optionId)}
              >
                <div className="flex items-start gap-2.5 border-b border-border/60 px-3.5 py-3">
                  <div className="pt-0.5 text-ai">
                    {isActive ? (
                      <CheckCircle2 className="h-[18px] w-[18px]" aria-hidden="true" />
                    ) : (
                      <Circle className="h-[18px] w-[18px] text-muted-foreground/60" aria-hidden="true" />
                    )}
                  </div>
                  <div className="min-w-0 flex-1 space-y-1">
                    <div className="flex flex-wrap items-center gap-2">
                      <p className="text-sm font-semibold leading-5 text-foreground">{option.title}</p>
                      {option.recommended ? (
                        <span className="rounded-full bg-success/12 px-2 py-0.5 text-[10px] font-medium text-success">
                          {t("candidate.recommended")}
                        </span>
                      ) : null}
                      {selectionLocked && isActive ? (
                        <span className="rounded-full bg-muted px-2 py-0.5 text-[10px] font-medium text-muted-foreground">
                          {t("candidate.confirmed")}
                        </span>
                      ) : null}
                    </div>
                    {option.summary ? <p className="text-[12px] leading-5 text-muted-foreground">{option.summary}</p> : null}
                  </div>
                </div>

                <div className="space-y-3 px-3.5 py-3">
                  <div className="grid gap-2.5 sm:grid-cols-2">
                    {option.fields.map((field) => (
                      <div
                        key={field.key}
                        className="min-w-0 rounded-xl border border-border/50 bg-muted/30 px-3 py-2.5 shadow-[inset_0_1px_0_rgba(255,255,255,0.4)]"
                      >
                        <p className="text-[11px] font-medium leading-4 text-muted-foreground">{field.label}</p>
                        <p className="mt-1 break-words text-sm font-medium leading-5 text-foreground">{field.value}</p>
                      </div>
                    ))}
                  </div>

                  {option.extra ? (
                    <div className="rounded-xl border border-ai/15 bg-ai/5 px-3 py-2.5">
                      {option.extra.title ? (
                        <p className="text-[10px] font-semibold uppercase tracking-[0.1em] text-ai">{option.extra.title}</p>
                      ) : null}
                      <div className="mt-1 flex items-start gap-2">
                        <Sparkles className="mt-0.5 h-4 w-4 shrink-0 text-ai" aria-hidden="true" />
                        <p className="break-words text-xs leading-5 text-muted-foreground">{option.extra.content}</p>
                      </div>
                    </div>
                  ) : null}
                </div>
              </button>

              {isActive && confirmAction ? (
                <div className="border-t border-border/60 px-3.5 pb-3.5 pt-3">
                  <Button
                    type="button"
                    variant={selectionLocked ? "secondary" : "ai"}
                    className="h-9 w-full gap-2 text-sm"
                    disabled={submitting || selectionLocked}
                    onClick={() => {
                      setOptimisticLockedOptionId(option.optionId)
                      onSelect({
                        reasoningTraceId: payload.reasoningTraceId,
                        candidateOptionId: option.optionId,
                        selection: "confirm",
                      })
                    }}
                  >
                    <Check className="h-4 w-4" aria-hidden="true" />
                    {confirmAction.label}
                  </Button>
                </div>
              ) : null}
            </article>
          )
        })}
      </div>
    </section>
  )
}
