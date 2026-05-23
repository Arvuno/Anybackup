import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { MemoryRouter } from "react-router-dom"
import { beforeEach, describe, expect, it } from "vitest"
import { ConversationSidebar } from "@/components/navigation/conversation-sidebar"
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

describe("ConversationSidebar", () => {
  beforeEach(() => {
    seedUser()
    useLayoutStore.setState({ sidebarCollapsed: false })
    useConversationStore.setState({
      conversations: [
        {
          conversationId: "conv-001",
          title: "订单数据库恢复",
          updatedAt: "2026-04-22T16:00:00+08:00",
        },
      ],
      query: "",
      listLoading: false,
      selectedWorkspace: {
        kind: "localDraft",
        localDraftId: "local-001",
      },
      error: null,
    })
  })

  it("uses the prototype-like brand row with a compact standalone new-conversation button", () => {
    const { container } = render(
      <MemoryRouter>
        <ConversationSidebar />
      </MemoryRouter>,
    )

    const aside = container.querySelector("aside")
    const newConversationButton = screen.getByRole("button", { name: "New conversation" })
    const searchInput = screen.getByPlaceholderText("Search conversations")

    expect(aside).toHaveClass("w-[240px]")
    expect(screen.getByText("Anybackup")).toBeInTheDocument()
    expect(screen.getByText("Agent")).toBeInTheDocument()
    expect(screen.queryByText("工作台 / 对话")).not.toBeInTheDocument()
    expect(newConversationButton).toHaveClass("h-7", "rounded-sm", "text-xs")
    expect(searchInput).toHaveClass("text-xs")
    expect(searchInput.parentElement).toHaveClass("h-7", "rounded-sm")
  })

  it("keeps settings separate from the user menu and shows logout inside the menu", async () => {
    const user = userEvent.setup()

    render(
      <MemoryRouter>
        <ConversationSidebar />
      </MemoryRouter>,
    )

    expect(screen.getAllByRole("button", { name: "Settings" })).toHaveLength(1)

    await user.click(screen.getByRole("button", { name: /备份管理员/ }))

    expect(screen.getAllByRole("button", { name: "Settings" })).toHaveLength(1)
    expect(screen.getByRole("button", { name: "Sign out" })).toBeInTheDocument()
  })

  it("renders history items with the lighter list treatment from the prototype", () => {
    render(
      <MemoryRouter>
        <ConversationSidebar />
      </MemoryRouter>,
    )

    const historyButton = screen.getByRole("button", { name: "订单数据库恢复" })

    expect(historyButton).toHaveClass("h-9", "rounded-md", "text-xs")
    expect(historyButton).not.toHaveClass("rounded-2xl")
  })

  it("shows all history items when the sidebar is collapsed", () => {
    useLayoutStore.setState({ sidebarCollapsed: true })
    useConversationStore.setState({
      conversations: [
        {
          conversationId: "conv-001",
          title: "会话一",
          updatedAt: "2026-04-22T16:00:00+08:00",
        },
        {
          conversationId: "conv-002",
          title: "会话二",
          updatedAt: "2026-04-22T16:01:00+08:00",
        },
        {
          conversationId: "conv-003",
          title: "会话三",
          updatedAt: "2026-04-22T16:02:00+08:00",
        },
        {
          conversationId: "conv-004",
          title: "会话四",
          updatedAt: "2026-04-22T16:03:00+08:00",
        },
      ],
    })

    render(
      <MemoryRouter>
        <ConversationSidebar />
      </MemoryRouter>,
    )

    expect(screen.getByTitle("会话一")).toBeInTheDocument()
    expect(screen.getByTitle("会话二")).toBeInTheDocument()
    expect(screen.getByTitle("会话三")).toBeInTheDocument()
    expect(screen.getByTitle("会话四")).toBeInTheDocument()
  })
})
