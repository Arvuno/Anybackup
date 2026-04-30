import * as React from "react"
import { vi } from "vitest"

vi.mock("@testing-library/react", async (importOriginal) => {
  const actual = await importOriginal<typeof import("@testing-library/react")>()
  const { I18nProvider } = await import("@/i18n")

  return {
    ...actual,
    render: (ui: React.ReactElement, options?: Parameters<typeof actual.render>[1]) => {
      const wrap = options?.wrapper
      const Wrapper = wrap
        ? function CombinedWrapper({ children }: { children: React.ReactNode }) {
            return React.createElement(I18nProvider, null, React.createElement(wrap, null, children))
          }
        : function I18nOnly({ children }: { children: React.ReactNode }) {
            return React.createElement(I18nProvider, null, children)
          }

      return actual.render(ui, {
        ...options,
        wrapper: Wrapper,
      })
    },
  }
})
