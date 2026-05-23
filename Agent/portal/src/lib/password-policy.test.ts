import { describe, expect, it } from "vitest"
import {
  validateLoginFields,
  validatePasswordConfirmation,
  validatePasswordPolicy,
} from "@/lib/password-policy"

describe("password policy", () => {
  it("returns username error before password error", () => {
    expect(validateLoginFields("", "")).toEqual({ field: "username", message: "Please enter username" })
  })

  it("returns password error when username is present", () => {
    expect(validateLoginFields("admin", "")).toEqual({ field: "password", message: "Please enter password" })
  })

  it("accepts valid login fields", () => {
    expect(validateLoginFields("admin", "admin1234")).toBeNull()
  })

  it("requires at least 8 characters", () => {
    expect(validatePasswordPolicy("abc12", "admin")).toBe("Password must be at least 8 characters")
  })

  it("requires letters and numbers", () => {
    expect(validatePasswordPolicy("abcdefgh", "admin")).toBe("Password must contain at least letters and numbers")
  })

  it("rejects password equal to username", () => {
    expect(validatePasswordPolicy("admin123", "admin123")).toBe("Password cannot be exactly the same as username")
  })

  it("accepts a valid password", () => {
    expect(validatePasswordPolicy("admin1234", "admin")).toBeNull()
  })

  it("checks confirmation mismatch", () => {
    expect(validatePasswordConfirmation("admin1234", "admin12345")).toBe("The two passwords do not match")
  })
})
