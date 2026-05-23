import type { UserRole } from "@/types/auth"

export type UserStatus = "enabled" | "disabled"

export interface ManagedUser {
  id: string
  username: string
  displayName: string
  role: UserRole
  status: UserStatus
  remark?: string
  createdAt: string
  updatedAt: string
  lastLoginAt?: string
}

export interface CreateUserInput {
  username: string
  displayName: string
  password: string
  confirmPassword: string
  status: UserStatus
  remark?: string
}

export interface UpdateUserInput {
  displayName: string
  status: UserStatus
  remark?: string
}

export interface ResetPasswordInput {
  password: string
  confirmPassword: string
}
