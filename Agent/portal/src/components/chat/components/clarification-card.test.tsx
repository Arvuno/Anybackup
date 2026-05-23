import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"
import { ClarificationCard } from "@/components/chat/components/clarification-card"

const payload = {
  clarificationId: "restore_window",
  prompt: "Please confirm the restore window.",
  options: [
    { label: "Latest safe point", value: "latest_safe_point" },
    { label: "Custom timestamp", value: "custom_timestamp" },
  ],
  inputConstraints: {
    required: true,
    allowFreeText: true,
    freeTextLabel: "Restore time",
    freeTextPlaceholder: "Enter a timestamp",
  },
}

describe("ClarificationCard", () => {
  it("renders the prompt, options, and free-text input", () => {
    render(<ClarificationCard payload={payload} submitting={false} onSubmit={vi.fn()} />)

    expect(screen.getByText("Please confirm the restore window.")).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Latest safe point" })).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Custom timestamp" })).toBeInTheDocument()
    expect(screen.getByLabelText("Restore time")).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Submit confirmation" })).toBeDisabled()
  })

  it("submits the selected option and free-text value", async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn()

    render(<ClarificationCard payload={payload} submitting={false} onSubmit={onSubmit} />)

    await user.click(screen.getByRole("button", { name: "Custom timestamp" }))
    await user.type(screen.getByLabelText("Restore time"), "2026-04-23 21:30:00")
    await user.click(screen.getByRole("button", { name: "Submit confirmation" }))

    expect(onSubmit).toHaveBeenCalledWith({
      clarificationId: "restore_window",
      selectedValue: "custom_timestamp",
      freeText: "2026-04-23 21:30:00",
    })
  })

  it("locks the card after a successful local submit", async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn()

    render(<ClarificationCard payload={payload} submitting={false} onSubmit={onSubmit} />)

    await user.click(screen.getByRole("button", { name: "Latest safe point" }))
    await user.click(screen.getByRole("button", { name: "Submit confirmation" }))

    expect(screen.getByRole("button", { name: "Latest safe point" })).toBeDisabled()
    expect(screen.getByRole("button", { name: "Custom timestamp" })).toBeDisabled()
    expect(screen.getByRole("button", { name: "Submitted" })).toBeDisabled()
    expect(screen.getAllByText("Submitted").length).toBeGreaterThan(0)
  })
})
