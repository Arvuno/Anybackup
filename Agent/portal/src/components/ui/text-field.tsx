import { forwardRef, type InputHTMLAttributes } from "react"
import { cn } from "@/lib/cn"

interface TextFieldProps extends InputHTMLAttributes<HTMLInputElement> {
  label: string
  error?: string
}

export const TextField = forwardRef<HTMLInputElement, TextFieldProps>(
  ({ id, label, error, className, required, ...props }, ref) => {
    const inputId = id ?? props.name
    const errorId = error && inputId ? `${inputId}-error` : undefined

    return (
      <div className="space-y-1.5">
        <label className="block text-xs font-medium text-foreground" htmlFor={inputId}>
          <span>{label}</span>
          {required ? (
            <span aria-hidden="true" className="ml-0.5 text-destructive">
              *
            </span>
          ) : null}
        </label>
        <input
          id={inputId}
          ref={ref}
          required={required}
          aria-invalid={Boolean(error)}
          aria-describedby={errorId}
          className={cn(
            "h-9 w-full rounded-md border bg-background px-3 text-sm text-foreground placeholder:text-muted-foreground transition-fast focus:border-ai/40 focus:outline-none focus:ring-2 focus:ring-ai/30 disabled:cursor-not-allowed disabled:bg-muted disabled:text-muted-foreground",
            error ? "border-destructive/60 bg-destructive/5" : "border-input",
            className,
          )}
          {...props}
        />
        {error ? (
          <p id={errorId} className="text-xs text-destructive">
            {error}
          </p>
        ) : null}
      </div>
    )
  },
)

TextField.displayName = "TextField"
