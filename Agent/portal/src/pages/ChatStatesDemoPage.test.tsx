import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it } from "vitest"
import { ChatStatesDemoPage } from "@/pages/ChatStatesDemoPage"

describe("ChatStatesDemoPage", () => {
  it("renders the three chat state demos", () => {
    render(<ChatStatesDemoPage />)

    expect(screen.getByRole("heading", { name: "Conversation state demo" })).toBeInTheDocument()
    expect(screen.getByRole("heading", { name: "Pattern 1 · Candidate confirmation card" })).toBeInTheDocument()
    expect(screen.getByRole("heading", { name: "Pattern 2 · Plain message" })).toBeInTheDocument()
    expect(screen.getByRole("heading", { name: "Pattern 3 · Waiting for reply" })).toBeInTheDocument()
  })

  it("lets the candidate demo grow to its full content height", () => {
    render(<ChatStatesDemoPage />)

    const candidateHeading = screen.getByRole("heading", { name: "Pattern 1 · Candidate confirmation card" })
    const section = candidateHeading.closest("section")

    expect(section).toBeInTheDocument()
    expect(section?.querySelector(".h-auto")).toBeInTheDocument()
  })

  it("updates the interaction summary when a candidate option is confirmed", async () => {
    const user = userEvent.setup()

    render(<ChatStatesDemoPage />)

    await user.click(screen.getByRole("button", { name: "Confirm and submit" }))

    expect(screen.getByText("Simulated confirm:option_a")).toBeInTheDocument()
  })
})
