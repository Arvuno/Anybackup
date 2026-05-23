import { ChatPanel } from "@/components/chat/ChatPanel"

export function WorkspaceHome() {
  return (
    <div className="flex h-full min-h-0 flex-1 flex-col overflow-hidden bg-gradient-hero">
      <div className="flex min-h-0 flex-1 flex-col px-4 pb-4 pt-3 md:px-8">
        <ChatPanel fillHeight className="min-h-[280px]" />
      </div>
    </div>
  )
}
