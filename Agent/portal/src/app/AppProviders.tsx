import { useEffect, type ReactNode } from "react"
import { useAuthStore } from "@/store/useAuthStore"
import { I18nProvider } from "@/i18n"

function AuthBootstrap({ children }: { children: ReactNode }) {
  const bootstrap = useAuthStore((state) => state.bootstrap)

  useEffect(() => {
    void bootstrap()
  }, [bootstrap])

  return <>{children}</>
}

export function AppProviders({ children }: { children: ReactNode }) {
  return (
    <I18nProvider>
      <AuthBootstrap>{children}</AuthBootstrap>
    </I18nProvider>
  )
}
