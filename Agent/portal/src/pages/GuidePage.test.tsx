import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { MemoryRouter, useLocation } from "react-router-dom"
import { beforeEach, describe, expect, it, vi } from "vitest"
import { GuidePage } from "@/pages/GuidePage"
import { routes } from "@/config/routes"
import { isGuideDone, storeSession } from "@/lib/session"
import type { AuthSession, CurrentUser } from "@/types/auth"

const adminUser: CurrentUser = {
  id: "user-001",
  username: "admin",
  displayName: "备份管理员",
  role: "backup_admin",
  tenantId: "tenant-001",
}

const secondUser: CurrentUser = {
  id: "user-002",
  username: "operator",
  displayName: "恢复操作员",
  role: "backup_admin",
  tenantId: "tenant-001",
}

function seedSession(user: CurrentUser = adminUser) {
  const session: AuthSession = {
    accessToken: `token-${user.id}`,
    expiresAt: Date.now() + 60 * 60 * 1000,
    idleExpiresAt: Date.now() + 30 * 60 * 1000,
    user,
  }
  storeSession(session)
}

function CurrentPath() {
  const location = useLocation()
  return <span data-testid="current-path">{location.pathname}</span>
}

describe("GuidePage", () => {
  beforeEach(() => {
    localStorage.clear()
    sessionStorage.clear()
    seedSession()
    vi.useRealTimers()
  })

  it("renders the prototype guide landing content", () => {
    render(
      <MemoryRouter>
        <GuidePage />
      </MemoryRouter>,
    )

    expect(screen.getByRole("heading", { name: "Anybackup" })).toBeInTheDocument()
    expect(
      screen.getByRole("heading", { name: "Natural language driven, AI-generated plans, and human approval" }),
    ).toBeInTheDocument()
    expect(screen.getByText("Anybackup Agent")).toBeInTheDocument()
    expect(screen.getByText("Anybackup Foundation")).toBeInTheDocument()
    expect(screen.getByText("Anybackup Client")).toBeInTheDocument()
    expect(screen.getByRole("button", { name: /Get started/ })).toBeInTheDocument()
    expect(screen.getByLabelText("Do not show this next time")).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Skip guide and enter now →" })).toBeInTheDocument()
  })

  it("moves through the prototype workflow intro", async () => {
    const user = userEvent.setup()

    render(
      <MemoryRouter>
        <GuidePage />
      </MemoryRouter>,
    )

    await user.click(screen.getByRole("button", { name: /Get started/ }))

    expect(await screen.findByRole("heading", { name: "Workflow" })).toBeInTheDocument()
    expect(screen.getByText("Describe requirements naturally")).toBeInTheDocument()
    expect(screen.getByText("Agent analyzes environment")).toBeInTheDocument()
    expect(screen.getByRole("button", { name: /View demo/ })).toBeInTheDocument()
  })

  it("temporarily dismisses guide for the current session when skipping without opt-out", async () => {
    const user = userEvent.setup()

    render(
      <MemoryRouter initialEntries={[routes.guide]}>
        <GuidePage />
        <CurrentPath />
      </MemoryRouter>,
    )

    await user.click(screen.getByRole("button", { name: "Skip guide and enter now →" }))

    expect(isGuideDone()).toBe(true)
    expect(screen.getByTestId("current-path")).toHaveTextContent(routes.home)

    sessionStorage.clear()

    expect(isGuideDone()).toBe(false)
  })

  it("persists don't-show-again by current user only", async () => {
    const user = userEvent.setup()

    render(
      <MemoryRouter initialEntries={[routes.guide]}>
        <GuidePage />
        <CurrentPath />
      </MemoryRouter>,
    )

    await user.click(screen.getByLabelText("Do not show this next time"))
    await user.click(screen.getByRole("button", { name: "Skip guide and enter now →" }))

    expect(isGuideDone(adminUser)).toBe(true)
    expect(screen.getByTestId("current-path")).toHaveTextContent(routes.home)

    sessionStorage.clear()

    expect(isGuideDone(adminUser)).toBe(true)
    expect(isGuideDone(secondUser)).toBe(false)
  })
})
