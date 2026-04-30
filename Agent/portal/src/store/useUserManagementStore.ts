import { create } from "zustand"
import {
  createUser,
  disableUser,
  enableUser,
  listUsers,
  resetUserPassword,
  updateUser,
} from "@/services/user-management-service"
import { ServiceError } from "@/types/auth"
import type {
  CreateUserInput,
  ManagedUser,
  ResetPasswordInput,
  UpdateUserInput,
} from "@/types/user-management"
import { translate } from "@/i18n/messages"

interface UserManagementState {
  users: ManagedUser[]
  loading: boolean
  saving: boolean
  message: string | null
  error: string | null
  feedbackAutoDismiss: boolean
}

interface UserManagementActions {
  loadUsers: () => Promise<void>
  create: (input: CreateUserInput) => Promise<void>
  update: (userId: string, input: UpdateUserInput) => Promise<void>
  enable: (userId: string) => Promise<void>
  disable: (userId: string, currentUserId: string) => Promise<void>
  resetPassword: (userId: string, input: ResetPasswordInput) => Promise<void>
  clearFeedback: () => void
}

function toMessage(error: unknown): string {
  if (error instanceof ServiceError) return error.message
  return translate("errors.genericUnavailable")
}

export const useUserManagementStore = create<UserManagementState & UserManagementActions>((set, get) => ({
  users: [],
  loading: false,
  saving: false,
  message: null,
  error: null,
  feedbackAutoDismiss: false,

  loadUsers: async () => {
    set({ loading: true, error: null, feedbackAutoDismiss: false })
    try {
      const users = await listUsers()
      set({ users, loading: false, feedbackAutoDismiss: false })
    } catch (error) {
      set({ error: toMessage(error), loading: false, feedbackAutoDismiss: false })
    }
  },

  create: async (input) => {
    set({ saving: true, error: null, message: null, feedbackAutoDismiss: false })
    try {
      await createUser(input)
      const users = await listUsers()
      set({ users, saving: false, message: translate("users.feedback.created"), feedbackAutoDismiss: true })
    } catch (error) {
      set({ saving: false, feedbackAutoDismiss: false })
      throw error
    }
  },

  update: async (userId, input) => {
    set({ saving: true, error: null, message: null, feedbackAutoDismiss: false })
    try {
      await updateUser(userId, input)
      const users = await listUsers()
      set({
        users,
        saving: false,
        message: translate("users.feedback.saved"),
        feedbackAutoDismiss: true,
      })
    } catch (error) {
      set({ saving: false, feedbackAutoDismiss: false })
      throw error
    }
  },

  enable: async (userId) => {
    set({ saving: true, error: null, message: null, feedbackAutoDismiss: false })
    try {
      await enableUser(userId)
      const users = await listUsers()
      set({ users, saving: false, message: translate("users.feedback.enabled"), feedbackAutoDismiss: true })
    } catch (error) {
      set({ saving: false, error: toMessage(error), feedbackAutoDismiss: true })
      throw error
    }
  },

  disable: async (userId, currentUserId) => {
    set({ saving: true, error: null, message: null, feedbackAutoDismiss: false })
    try {
      await disableUser(userId, currentUserId)
      const users = await listUsers()
      set({ users, saving: false, message: translate("users.feedback.disabled"), feedbackAutoDismiss: true })
    } catch (error) {
      set({ saving: false, error: toMessage(error), feedbackAutoDismiss: true })
      throw error
    }
  },

  resetPassword: async (userId, input) => {
    set({ saving: true, error: null, message: null, feedbackAutoDismiss: false })
    try {
      await resetUserPassword(userId, input)
      set({
        saving: false,
        message: translate("users.feedback.passwordReset"),
        feedbackAutoDismiss: true,
      })
    } catch (error) {
      set({ saving: false, error: toMessage(error), feedbackAutoDismiss: true })
      throw error
    }
  },

  clearFeedback: () => {
    const { message, error } = get()
    if (message || error) set({ message: null, error: null, feedbackAutoDismiss: false })
  },
}))
