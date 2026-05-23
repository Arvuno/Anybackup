import { translate } from "@/i18n/messages"

export interface FieldError {
  field: string
  message: string
}

export function validateLoginFields(username: string, password: string): FieldError | null {
  if (!username.trim()) return { field: "username", message: translate("validation.requiredUsername") }
  if (!password.trim()) return { field: "password", message: translate("validation.requiredPassword") }
  return null
}

export function validateRequired(value: string, field: string, message: string): FieldError | null {
  return value.trim() ? null : { field, message }
}

export function validatePasswordPolicy(password: string, username: string): string | null {
  if (password.length < 8) return translate("validation.passwordMinLength")
  if (password === username) return translate("validation.passwordSameAsUsername")
  if (!/[A-Za-z]/.test(password) || !/\d/.test(password)) return translate("validation.passwordNeedsLetterAndDigit")
  return null
}

export function validatePasswordConfirmation(password: string, confirmPassword: string): string | null {
  return password === confirmPassword ? null : translate("validation.passwordConfirmationMismatch")
}
