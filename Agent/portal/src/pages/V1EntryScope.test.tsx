import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { MemoryRouter } from "react-router-dom"
import { beforeEach, describe, expect, it, vi } from "vitest"
import { AppShell } from "@/components/layout/AppShell"
import { HomePage } from "@/pages/HomePage"
import { GuidePage } from "@/pages/GuidePage"
import { useAuthStore } from "@/store/useAuthStore"
import { useConversationStore } from "@/store/useConversationStore"
import { useLayoutStore } from "@/store/useLayoutStore"

function seedUser() {
  useAuthStore.setState({
    currentUser: {
      id: "user-001",
      username: "admin",
      displayName: "备份管理员",
      role: "backup_admin",
      tenantId: "tenant-001",
    },
    bootstrapped: true,
    loading: false,
    error: null,
  })
}

describe("V1 entry scope", () => {
  beforeEach(() => {
    localStorage.clear()
    sessionStorage.clear()
    if (!Element.prototype.scrollTo) {
      Element.prototype.scrollTo = vi.fn()
    }
    useLayoutStore.setState({ sidebarCollapsed: false })
    seedUser()
    useConversationStore.setState({
      bootstrapped: true,
      listLoading: false,
      conversationLoading: false,
      query: "",
      error: null,
      conversations: [],
      selectedWorkspace: { kind: "localDraft", localDraftId: "v1-entry-draft" },
      localDraftWorkspace: {
        localDraftId: "v1-entry-draft",
        title: "",
        draft: "",
        createdAt: "2026-04-22T10:00:00.000Z",
        updatedAt: "2026-04-22T10:00:00.000Z",
      },
      detailsByConversationId: {},
      messagesByConversationId: {},
      draftsByKey: {},
      pendingTurnByConversationId: {},
      pendingTurnStartedAtMsByConversationId: {},
      nextPollAfterMsByConversationId: {},
      eventCursorByConversationId: {},
      latestEventSequenceByConversationId: {},
      appliedStatusEventIdsByConversationId: {},
      submittingWorkspaceKey: null,
    })
  })

  it("renders only chat entry on the home page", () => {
    render(
      <MemoryRouter>
        <HomePage />
      </MemoryRouter>,
    )

    expect(screen.getByRole("heading", { name: "Start a new conversation here", level: 1 })).toBeInTheDocument()
    expect(screen.getByPlaceholderText("Tell me what you want to back up, restore, or check...")).toBeInTheDocument()
    expect(screen.queryByText("待审批")).not.toBeInTheDocument()
    expect(screen.queryByText("执行")).not.toBeInTheDocument()
    expect(screen.queryByText("CLI")).not.toBeInTheDocument()
    expect(screen.queryByText("冷却")).not.toBeInTheDocument()
  })

  it("uses V9 semantic color tokens for the chat entry title and send action", () => {
    render(
      <MemoryRouter>
        <HomePage />
      </MemoryRouter>,
    )

    const heading = screen.getByRole("heading", { name: "Start a new conversation here", level: 1 })
    const sendButton = screen.getByRole("button", { name: "Send message" })

    expect(heading).toHaveClass("text-lg", "font-bold", "text-foreground")
    expect(sendButton.className).toContain("bg-gradient-ai")
  })

  it("renders prototype guide entry content", () => {
    render(
      <MemoryRouter>
        <GuidePage />
      </MemoryRouter>,
    )

    expect(screen.getByRole("heading", { name: "Anybackup" })).toBeInTheDocument()
    expect(screen.getByRole("heading", { name: "Natural language driven, AI-generated plans, and human approval" })).toBeInTheDocument()
    expect(screen.getByRole("button", { name: /Get started/ })).toBeInTheDocument()
    expect(screen.getByLabelText("Do not show this next time")).toBeInTheDocument()
  })

  it("keeps settings separate while leaving logout in the user menu", async () => {
    const user = userEvent.setup()
    render(
      <MemoryRouter>
        <AppShell />
      </MemoryRouter>,
    )

    expect(screen.queryByRole("link", { name: /新手引导/ })).not.toBeInTheDocument()
    expect(screen.getAllByRole("button", { name: "Settings" })).toHaveLength(1)

    await user.click(screen.getByRole("button", { name: /备份管理员/ }))

    expect(screen.queryByRole("button", { name: "新手引导" })).not.toBeInTheDocument()
    expect(screen.getAllByRole("button", { name: "Settings" })).toHaveLength(1)
    expect(screen.getByRole("button", { name: "Sign out" })).toBeInTheDocument()
  })

  it("logs out through the auth store", async () => {
    const logout = vi.fn().mockResolvedValue(undefined)
    useAuthStore.setState({ logout })
    const user = userEvent.setup()

    render(
      <MemoryRouter>
        <AppShell />
      </MemoryRouter>,
    )

    await user.click(screen.getByRole("button", { name: /备份管理员/ }))
    await user.click(screen.getByRole("button", { name: "Sign out" }))

    expect(logout).toHaveBeenCalledTimes(1)
  })

  it("logs out directly even when there is an unsent draft", async () => {
    const logout = vi.fn().mockResolvedValue(undefined)
    localStorage.setItem("agent_web_draft_test", "还没发送")
    useAuthStore.setState({ logout })
    const user = userEvent.setup()

    render(
      <MemoryRouter>
        <AppShell />
      </MemoryRouter>,
    )

    await user.click(screen.getByRole("button", { name: /备份管理员/ }))
    await user.click(screen.getByRole("button", { name: "Sign out" }))

    expect(logout).toHaveBeenCalledTimes(1)
  })
})
