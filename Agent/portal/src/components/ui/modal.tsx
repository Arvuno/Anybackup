import type { ReactNode } from "react"
import { X } from "lucide-react"
import { cn } from "@/lib/cn"
import { useI18n } from "@/i18n"

interface ModalProps {
  open: boolean
  title: string
  description?: string
  children: ReactNode
  footer?: ReactNode
  onClose: () => void
  className?: string
}

export function Modal({ open, title, description, children, footer, onClose, className }: ModalProps) {
  const { t } = useI18n()
  if (!open) return null

  return (
    <div className="fixed inset-0 z-[220] flex items-center justify-center bg-foreground/18 p-4 backdrop-blur-[2px]">
      <section
        role="dialog"
        aria-modal="true"
        aria-labelledby="modal-title"
        aria-describedby={description ? "modal-description" : undefined}
        className={cn(
          "w-full max-w-lg rounded-xl border border-border/80 bg-card/95 shadow-modal",
          className,
        )}
      >
        <header className="flex items-start justify-between gap-4 border-b border-border/70 px-5 py-4">
          <div className="min-w-0">
            <h2 id="modal-title" className="text-[15px] font-semibold tracking-tight text-foreground">
              {title}
            </h2>
            {description ? (
              <p id="modal-description" className="mt-1 text-sm leading-6 text-muted-foreground">
                {description}
              </p>
            ) : null}
          </div>
          <button
            type="button"
            aria-label={t("modal.close")}
            onClick={onClose}
            className="rounded-md p-1 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
          >
            <X className="h-4 w-4" />
          </button>
        </header>
        <div className="px-5 py-4">{children}</div>
        {footer ? (
          <footer className="flex flex-col-reverse justify-end gap-2 border-t border-border/70 bg-muted/15 px-5 py-4 sm:flex-row">
            {footer}
          </footer>
        ) : null}
      </section>
    </div>
  )
}
