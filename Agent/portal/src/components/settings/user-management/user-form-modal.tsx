import { useEffect, useState } from "react"
import { ToggleLeft, ToggleRight } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Modal } from "@/components/ui/modal"
import { TextField } from "@/components/ui/text-field"
import {
  validatePasswordConfirmation,
  validatePasswordPolicy,
  validateRequired,
} from "@/lib/password-policy"
import { useI18n } from "@/i18n"
import { translate } from "@/i18n/messages"
import { ServiceError } from "@/types/auth"
import type { CreateUserInput, ManagedUser, UpdateUserInput, UserStatus } from "@/types/user-management"

type UserFormField = "username" | "displayName" | "password" | "confirmPassword"
type UserFormFieldErrors = Partial<Record<UserFormField, string>>

interface UserFormModalProps {
  open: boolean
  mode: "create" | "edit"
  user?: ManagedUser
  saving: boolean
  onClose: () => void
  onCreate: (input: CreateUserInput) => Promise<void>
  onUpdate: (userId: string, input: UpdateUserInput) => Promise<void>
}

function hasFieldErrors(errors: UserFormFieldErrors): boolean {
  return Object.values(errors).some(Boolean)
}

function validateCreateForm(input: CreateUserInput): UserFormFieldErrors {
  const errors: UserFormFieldErrors = {}

  const usernameError = validateRequired(input.username, "username", translate("validation.requiredUsername"))
  if (usernameError) errors.username = usernameError.message

  const displayNameError = validateRequired(input.displayName, "displayName", translate("validation.requiredDisplayName"))
  if (displayNameError) errors.displayName = displayNameError.message

  const passwordRequiredError = validateRequired(input.password, "password", translate("validation.requiredPassword"))
  if (passwordRequiredError) errors.password = passwordRequiredError.message

  const confirmPasswordRequiredError = validateRequired(
    input.confirmPassword,
    "confirmPassword",
    translate("validation.confirmPasswordRequired"),
  )
  if (confirmPasswordRequiredError) errors.confirmPassword = confirmPasswordRequiredError.message

  if (!errors.password) {
    const passwordPolicyError = validatePasswordPolicy(input.password, input.username)
    if (passwordPolicyError) errors.password = passwordPolicyError
  }

  if (!errors.confirmPassword && !errors.password) {
    const confirmPasswordError = validatePasswordConfirmation(input.password, input.confirmPassword)
    if (confirmPasswordError) errors.confirmPassword = confirmPasswordError
  }

  return errors
}

function validateUpdateForm(input: UpdateUserInput): UserFormFieldErrors {
  const errors: UserFormFieldErrors = {}
  const displayNameError = validateRequired(input.displayName, "displayName", translate("validation.requiredDisplayName"))
  if (displayNameError) errors.displayName = displayNameError.message
  return errors
}

function mapServiceError(error: unknown): { fieldErrors: UserFormFieldErrors; formError: string | null } {
  if (error instanceof ServiceError) {
    if (error.code === "USERNAME_EXISTS") {
      return { fieldErrors: { username: error.message }, formError: null }
    }
    if (error.code === "PASSWORD_POLICY_VIOLATION") {
      return { fieldErrors: { password: error.message }, formError: null }
    }
    if (error.code === "PASSWORD_CONFIRMATION_MISMATCH") {
      return { fieldErrors: { confirmPassword: error.message }, formError: null }
    }
    return { fieldErrors: {}, formError: error.message }
  }

  return {
    fieldErrors: {},
    formError: translate("userForm.formError"),
  }
}

export function UserFormModal({ open, mode, user, saving, onClose, onCreate, onUpdate }: UserFormModalProps) {
  const { t } = useI18n()
  const [username, setUsername] = useState("")
  const [displayName, setDisplayName] = useState("")
  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const [status, setStatus] = useState<UserStatus>("enabled")
  const [remark, setRemark] = useState("")
  const [fieldErrors, setFieldErrors] = useState<UserFormFieldErrors>({})
  const [formError, setFormError] = useState<string | null>(null)

  useEffect(() => {
    if (!open) return

    setUsername(user?.username ?? "")
    setDisplayName(user?.displayName ?? "")
    setPassword("")
    setConfirmPassword("")
    setStatus(user?.status ?? "enabled")
    setRemark(user?.remark ?? "")
    setFieldErrors({})
    setFormError(null)
  }, [open, user])

  const clearFieldError = (field: UserFormField) => {
    setFieldErrors((current) => {
      if (!current[field]) return current
      return { ...current, [field]: undefined }
    })
    setFormError(null)
  }

  const submit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setFormError(null)

    if (mode === "create") {
      const input: CreateUserInput = { username, displayName, password, confirmPassword, status, remark }
      const nextErrors = validateCreateForm(input)
      if (hasFieldErrors(nextErrors)) {
        setFieldErrors(nextErrors)
        return
      }

      setFieldErrors({})
      try {
        await onCreate(input)
        onClose()
      } catch (error) {
        const nextState = mapServiceError(error)
        setFieldErrors(nextState.fieldErrors)
        setFormError(nextState.formError)
      }
      return
    }

    if (!user) return

    const input: UpdateUserInput = { displayName, status, remark }
    const nextErrors = validateUpdateForm(input)
    if (hasFieldErrors(nextErrors)) {
      setFieldErrors(nextErrors)
      return
    }

    setFieldErrors({})
    try {
      await onUpdate(user.id, input)
      onClose()
    } catch (error) {
      const nextState = mapServiceError(error)
      setFieldErrors(nextState.fieldErrors)
      setFormError(nextState.formError)
    }
  }

  return (
    <Modal
      open={open}
      title={mode === "create" ? t("userForm.titleCreate") : t("userForm.titleEdit")}
      description={mode === "create" ? t("userForm.descCreate") : t("userForm.descEdit")}
      onClose={onClose}
      footer={
        <>
          <Button type="button" variant="ghost" onClick={onClose}>
            {t("userForm.cancel")}
          </Button>
          <Button type="submit" form="user-form" disabled={saving}>
            {mode === "create" ? t("userForm.create") : t("userForm.save")}
          </Button>
        </>
      }
    >
      <form id="user-form" className="space-y-5" onSubmit={submit}>
        {formError ? (
          <div
            role="alert"
            className="rounded-lg border border-destructive/15 bg-destructive/5 px-3 py-2 text-sm text-destructive"
          >
            {formError}
          </div>
        ) : null}

        <div className="grid gap-4 sm:grid-cols-2">
          <TextField
            label={t("userForm.labelUsername")}
            name="username"
            required
            error={fieldErrors.username}
            value={username}
            onChange={(event) => {
              setUsername(event.target.value)
              clearFieldError("username")
            }}
            disabled={mode === "edit"}
          />
          <TextField
            label={t("userForm.labelDisplayName")}
            name="displayName"
            required
            error={fieldErrors.displayName}
            value={displayName}
            onChange={(event) => {
              setDisplayName(event.target.value)
              clearFieldError("displayName")
            }}
          />
        </div>

        {mode === "create" ? (
          <div className="grid gap-4 sm:grid-cols-2">
            <TextField
              label={t("userForm.labelPassword")}
              name="password"
              type="password"
              required
              error={fieldErrors.password}
              value={password}
              onChange={(event) => {
                setPassword(event.target.value)
                clearFieldError("password")
              }}
            />
            <TextField
              label={t("userForm.labelConfirmPassword")}
              name="confirmPassword"
              type="password"
              required
              error={fieldErrors.confirmPassword}
              value={confirmPassword}
              onChange={(event) => {
                setConfirmPassword(event.target.value)
                clearFieldError("confirmPassword")
              }}
            />
          </div>
        ) : null}

        <div className="space-y-1.5">
          <div>
            <span className="block text-xs font-medium text-foreground">{t("userForm.status")}</span>
            <p className="mt-1 text-[11px] text-muted-foreground">{t("userForm.statusHint")}</p>
          </div>
          <div className="rounded-lg border border-border/70 bg-card px-3 py-3 shadow-[var(--shadow-xs)]">
            <div className="flex items-center gap-3">
              <div className="min-w-0 flex-1">
                <p className="text-sm font-medium text-foreground">
                  {status === "enabled" ? t("userForm.statusEnabled") : t("userForm.statusDisabled")}
                </p>
              </div>
              <button
                type="button"
                role="switch"
                aria-checked={status === "enabled"}
                aria-label={status === "enabled" ? t("userForm.enableToDisableHint") : t("userForm.disableToEnableHint")}
                onClick={() => {
                  setStatus((current) => (current === "enabled" ? "disabled" : "enabled"))
                  setFormError(null)
                }}
                className="focus-ring shrink-0 rounded-md transition-fast"
              >
                {status === "enabled" ? (
                  <ToggleRight className="h-8 w-8 text-ai" />
                ) : (
                  <ToggleLeft className="h-8 w-8 text-muted-foreground/40" />
                )}
              </button>
            </div>
          </div>
        </div>
      </form>
    </Modal>
  )
}
