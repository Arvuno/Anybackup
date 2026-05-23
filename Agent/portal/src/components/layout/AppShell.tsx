import { useEffect } from "react"
import { Outlet } from "react-router-dom"
import { ConversationSidebar } from "@/components/navigation/conversation-sidebar"
import { useConversationStore } from "@/store/useConversationStore"

export function AppShell() {
  const hydrate = useConversationStore((state) => state.hydrate)

  useEffect(() => {
    void hydrate()
  }, [hydrate])

  return (
    <div className="flex h-screen w-screen overflow-hidden bg-background">
      <ConversationSidebar />
      <main className="flex min-h-0 flex-1 overflow-hidden bg-background">
        <div className="flex min-h-0 min-w-0 flex-1 flex-col overflow-hidden">
          <Outlet />
        </div>
      </main>
    </div>
  )
}
