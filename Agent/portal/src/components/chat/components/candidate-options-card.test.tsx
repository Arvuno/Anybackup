import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"
import { CandidateOptionsCard } from "@/components/chat/components/candidate-options-card"

const payload = {
  reasoningTraceId: "trace_restore_001",
  title: "Two restore candidates are ready",
  summary: "The restore plan assessment has finished.",
  actions: [
    { type: "confirm" as const, label: "Confirm restore plan" },
    { type: "reject" as const, label: "Reject plan" },
    {
      type: "revise" as const,
      label: "Add constraints",
      inputLabel: "Additional constraints",
      inputPlaceholder: "For example: generate the plan only.",
      submitLabel: "Submit constraints",
    },
  ],
  options: [
    {
      optionId: "option_a",
      title: "Option A: Restore by exporting the target table",
      recommended: true,
      summary: "Recommended option",
      fields: [
        { key: "restore_scope", label: "Restore scope", value: "Database-level" },
        { key: "rpo_rto", label: "RPO / RTO", value: "< 2 min / 1.5 h" },
        { key: "target", label: "Target", value: "Isolated MySQL restore environment" },
        { key: "coverage", label: "Coverage", value: "order_details only" },
      ],
      extra: {
        title: "Recommendation",
        content: "This path keeps the production database isolated while recovering the target table.",
      },
    },
    {
      optionId: "option_b",
      title: "Option B: Restore the full database in production",
      fields: [
        { key: "restore_scope", label: "Restore scope", value: "Database-level" },
        { key: "rpo_rto", label: "RPO / RTO", value: "< 2 min / 4 h" },
        { key: "target", label: "Target", value: "Production OrderDB" },
        { key: "coverage", label: "Coverage", value: "Entire OrderDB" },
      ],
      extra: {
        title: "Recommendation",
        content: "This path overwrites other healthy tables in production and carries higher risk.",
      },
    },
  ],
}

describe("CandidateOptionsCard", () => {
  it("renders structured candidate options and emits confirm selection", async () => {
    const user = userEvent.setup()
    const onSelect = vi.fn()

    render(<CandidateOptionsCard payload={payload} submitting={false} onSelect={onSelect} />)

    expect(screen.getByText("Option A: Restore by exporting the target table")).toBeInTheDocument()
    expect(screen.getAllByText("Restore scope").length).toBeGreaterThan(0)
    expect(screen.getAllByText("Database-level").length).toBeGreaterThan(0)
    expect(screen.getAllByText("Recommendation")).toHaveLength(2)

    await user.click(screen.getByRole("button", { name: "Confirm restore plan" }))

    expect(onSelect).toHaveBeenCalledWith({
      reasoningTraceId: "trace_restore_001",
      candidateOptionId: "option_a",
      selection: "confirm",
    })
  })

  it("shows only the confirm action for the selected option", () => {
    const onSelect = vi.fn()

    render(<CandidateOptionsCard payload={payload} submitting={false} onSelect={onSelect} />)

    expect(screen.getByRole("button", { name: "Confirm restore plan" })).toBeInTheDocument()
    expect(screen.queryByRole("button", { name: "Reject plan" })).not.toBeInTheDocument()
    expect(screen.queryByRole("button", { name: "Add constraints" })).not.toBeInTheDocument()
    expect(screen.queryByRole("textbox", { name: "Additional constraints" })).not.toBeInTheDocument()
  })

  it("uses a wrapping grid layout so up to three options can sit on one row", () => {
    const onSelect = vi.fn()

    render(<CandidateOptionsCard payload={payload} submitting={false} onSelect={onSelect} />)

    expect(screen.getByTestId("candidate-options-grid")).toHaveClass(
      "[grid-template-columns:repeat(auto-fit,minmax(min(100%,16rem),1fr))]",
    )
  })

  it("allows switching to another option before confirming", async () => {
    const user = userEvent.setup()
    const onSelect = vi.fn()

    render(<CandidateOptionsCard payload={payload} submitting={false} onSelect={onSelect} />)

    await user.click(screen.getByRole("button", { name: /Option B: Restore the full database in production/ }))
    await user.click(screen.getByRole("button", { name: "Confirm restore plan" }))

    expect(onSelect).toHaveBeenCalledWith({
      reasoningTraceId: "trace_restore_001",
      candidateOptionId: "option_b",
      selection: "confirm",
    })
  })

  it("locks the card buttons after confirming a selection", async () => {
    const user = userEvent.setup()
    const onSelect = vi.fn()

    render(<CandidateOptionsCard payload={payload} submitting={false} onSelect={onSelect} />)

    const optionAButton = screen.getByRole("button", { name: /Option A: Restore by exporting the target table/i })
    const optionBButton = screen.getByRole("button", { name: /Option B: Restore the full database in production/i })
    const confirmButton = screen.getByRole("button", { name: "Confirm restore plan" })

    await user.click(confirmButton)

    expect(onSelect).toHaveBeenCalledWith({
      reasoningTraceId: "trace_restore_001",
      candidateOptionId: "option_a",
      selection: "confirm",
    })
    expect(optionAButton).toBeDisabled()
    expect(optionBButton).toBeDisabled()
    expect(screen.getByRole("button", { name: "Confirm restore plan" })).toBeDisabled()
  })
})
