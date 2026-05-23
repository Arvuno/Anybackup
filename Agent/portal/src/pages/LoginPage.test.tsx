import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { MemoryRouter } from "react-router-dom"
import { beforeEach, describe, expect, it, vi } from "vitest"
import { clearStoredSession } from "@/lib/session"
import { LoginPage } from "@/pages/LoginPage"
import { useAuthStore } from "@/store/useAuthStore"

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { "Content-Type": "application/json" },
  })
}

function renderLogin() {
  return render(
    <MemoryRouter initialEntries={["/login"]}>
      <LoginPage />
    </MemoryRouter>,
  )
}

describe("LoginPage", () => {
  const fetchMock = vi.fn<typeof fetch>()

  beforeEach(() => {
    vi.restoreAllMocks()
    fetchMock.mockReset()
    vi.stubGlobal("fetch", fetchMock)
    localStorage.clear()
    clearStoredSession()
    useAuthStore.setState({
      currentUser: null,
      bootstrapped: false,
      loading: false,
      error: null,
    })
  })

  it("shows username required first", async () => {
    const user = userEvent.setup()
    renderLogin()

    await user.click(screen.getByRole("button", { name: "Sign in" }))

    expect(screen.getByText("Please enter username")).toBeInTheDocument()
  })

  it("shows password required when username exists", async () => {
    const user = userEvent.setup()
    renderLogin()

    await user.type(screen.getByLabelText("Username"), "admin")
    await user.click(screen.getByRole("button", { name: "Sign in" }))

    expect(screen.getByText("Please enter password")).toBeInTheDocument()
  })

  it("shows invalid credential message", async () => {
    fetchMock.mockResolvedValueOnce(jsonResponse({ error: "invalid_grant" }, 401))
    const user = userEvent.setup()
    renderLogin()

    await user.type(screen.getByLabelText("Username"), "admin")
    await user.type(screen.getByLabelText("Password"), "wrong1234")
    await user.click(screen.getByRole("button", { name: "Sign in" }))

    expect(await screen.findByText("Invalid username or password.")).toBeInTheDocument()
  })

  it("toggles password visibility", async () => {
    const user = userEvent.setup()
    renderLogin()

    const password = screen.getByLabelText("Password")
    expect(password).toHaveAttribute("type", "password")

    await user.click(screen.getByRole("button", { name: "Show password" }))

    expect(password).toHaveAttribute("type", "text")
  })
})
