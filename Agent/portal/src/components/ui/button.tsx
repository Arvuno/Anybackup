import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "@/lib/cn"

const buttonVariants = cva(
  "inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-all duration-200 ease-smooth focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 active:scale-[0.98]",
  {
    variants: {
      variant: {
        default:
          "bg-primary text-primary-foreground shadow-sm hover:bg-primary/90 hover:-translate-y-[1px] hover:shadow-md",
        secondary: "bg-secondary text-secondary-foreground hover:bg-secondary/80 hover:-translate-y-[1px]",
        outline:
          "border border-input bg-card text-foreground hover:bg-accent hover:text-accent-foreground hover:-translate-y-[1px] hover:shadow-sm",
        ghost: "text-foreground hover:bg-accent hover:text-accent-foreground",
        destructive:
          "bg-destructive text-destructive-foreground hover:bg-destructive/90 hover:-translate-y-[1px] hover:shadow-md",
        link: "text-primary underline-offset-4 hover:underline",
        ai: "bg-gradient-ai text-primary-foreground shadow-ai hover:opacity-95 hover:-translate-y-[1px] hover:shadow-ai-lg",
        success: "bg-success text-success-foreground hover:bg-success/90 hover:-translate-y-[1px] hover:shadow-sm",
        warning: "bg-warning text-warning-foreground hover:bg-warning/90 hover:-translate-y-[1px] hover:shadow-sm",
      },
      size: {
        default: "h-9 px-4 py-2",
        sm: "h-8 rounded-md px-3 text-xs",
        lg: "h-11 rounded-md px-6",
        xs: "h-7 rounded-md px-2.5 text-xs",
        icon: "h-9 w-9",
      },
    },
    defaultVariants: { variant: "default", size: "default" },
  },
)

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, ...props }, ref) => (
    <button className={cn(buttonVariants({ variant, size, className }))} ref={ref} {...props} />
  ),
)
Button.displayName = "Button"

export { Button, buttonVariants }
