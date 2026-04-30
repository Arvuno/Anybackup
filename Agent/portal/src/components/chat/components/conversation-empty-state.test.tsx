import { render, screen } from "@testing-library/react"
import { describe, expect, it } from "vitest"
import { ConversationEmptyState } from "@/components/chat/components/conversation-empty-state"

describe("ConversationEmptyState", () => {
  it("does not paint its own background over the content area", () => {
    const { container } = render(<ConversationEmptyState composer={<div>Composer slot</div>} />)

    const root = screen.getByText("Composer slot").closest(".relative.grid")

    expect(root).not.toHaveClass("bg-background")
    expect(root).not.toHaveClass("bg-gradient-hero")
    expect(container.querySelector(".bg-gradient-hero")).not.toBeInTheDocument()
  })
})
