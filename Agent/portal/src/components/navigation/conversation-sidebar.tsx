import { useEffect, useRef, useState } from "react"
import { useNavigate } from "react-router-dom"
import {
  ChevronLeft,
  ChevronRight,
  LoaderCircle,
  LogOut,
  MessageSquare,
  MessageSquarePlus,
  Search,
  Settings,
  User,
} from "lucide-react"
import { routes } from "@/config/routes"
import { useI18n } from "@/i18n"
import { cn } from "@/lib/cn"
import { useAuthStore } from "@/store/useAuthStore"
import { useConversationStore } from "@/store/useConversationStore"
import { useLayoutStore } from "@/store/useLayoutStore"
import type { ConversationSummary } from "@/types/conversation"

export function ConversationSidebar() {
  const { t } = useI18n()
  const navigate = useNavigate()
  const sidebarCollapsed = useLayoutStore((state) => state.sidebarCollapsed)
  const toggleSidebar = useLayoutStore((state) => state.toggleSidebar)
  const currentUser = useAuthStore((state) => state.currentUser)
  const logout = useAuthStore((state) => state.logout)
  const conversations = useConversationStore((state) => state.conversations)
  const query = useConversationStore((state) => state.query)
  const listLoading = useConversationStore((state) => state.listLoading)
  const selectedWorkspace = useConversationStore((state) => state.selectedWorkspace)
  const activateLocalDraftWorkspace = useConversationStore((state) => state.activateLocalDraftWorkspace)
  const selectConversation = useConversationStore((state) => state.selectConversation)
  const setSearchQuery = useConversationStore((state) => state.setSearchQuery)

  const [userMenuOpen, setUserMenuOpen] = useState(false)
  const menuRef = useRef<HTMLDivElement>(null)
  const displayName = currentUser?.displayName ?? t("sidebar.demoUser")
  const accountName = currentUser?.username ?? "backup_admin"
  const roleLabel = t("sidebar.backupAdmin")

  useEffect(() => {
    if (!userMenuOpen) return

    const handler = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setUserMenuOpen(false)
      }
    }

    document.addEventListener("mousedown", handler)
    return () => document.removeEventListener("mousedown", handler)
  }, [userMenuOpen])

  const handleNewConversation = () => {
    activateLocalDraftWorkspace()
    navigate(routes.home)
  }

  const handleSelectConversation = async (conversationId: string) => {
    navigate(routes.home)
    try {
      await selectConversation(conversationId)
    } catch {
      // The store already exposes the error to the main panel.
    }
  }

  const handleLogout = async () => {
    setUserMenuOpen(false)
    await logout()
    navigate(routes.login)
  }

  const renderHistoryButton = (conversation: ConversationSummary, collapsed = false) => {
    const active =
      selectedWorkspace?.kind === "conversation" && selectedWorkspace.conversationId === conversation.conversationId

    return (
      <button
        key={conversation.conversationId}
        type="button"
        onClick={() => void handleSelectConversation(conversation.conversationId)}
        title={collapsed ? conversation.title : undefined}
        className={cn(
          "focus-ring transition-fast",
          collapsed
            ? "flex h-10 w-10 items-center justify-center rounded-lg bg-white/85 text-xs font-semibold text-foreground shadow-[var(--shadow-xs)] hover:bg-white"
            : "flex h-9 w-full items-center gap-2 rounded-md px-2.5 text-left text-xs",
          active
            ? collapsed
              ? "bg-ai-surface text-ai shadow-card"
              : "bg-accent text-accent-foreground shadow-sm"
            : collapsed
              ? "border-transparent"
              : "text-foreground/70 hover:bg-accent/50",
        )}
      >
        {collapsed ? (
          conversation.title.slice(0, 1)
        ) : (
          <>
            <MessageSquare className="h-3.5 w-3.5 shrink-0 text-muted-foreground" />
            <span className="min-w-0 truncate font-medium">{conversation.title}</span>
          </>
        )}
      </button>
    )
  }

  const renderUserMenuCard = (className: string) => (
    <div
      className={cn(
        "z-[160] w-64 overflow-hidden rounded-xl border border-border/80 bg-card/95 shadow-[0_20px_48px_-28px_rgba(15,23,42,0.35)] backdrop-blur-sm animate-fade-in-scale",
        className,
      )}
    >
      <div className="border-b border-border/70 bg-gradient-to-br from-ai-surface via-white to-white/95 px-4 py-4">
        <div className="flex items-start gap-3">
          <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-gradient-to-br from-[hsl(var(--ai))] to-[hsl(var(--ai)/0.7)] text-primary-foreground shadow-sm">
            <User className="h-4 w-4" />
          </div>
          <div className="min-w-0 flex-1">
            <p className="truncate text-sm font-semibold text-foreground">{displayName}</p>
            <p className="mt-1 text-xs text-muted-foreground">{roleLabel}</p>
            <p className="mt-2 text-xs text-muted-foreground">
              {t("sidebar.loginAccount")}：{accountName}
            </p>
          </div>
        </div>
      </div>
      <div className="p-2">
        <button
          type="button"
          onClick={() => void handleLogout()}
          className="focus-ring flex w-full items-center gap-2.5 rounded-md px-3 py-2 text-xs font-medium text-destructive transition-colors hover:bg-destructive/5"
        >
          <LogOut className="h-3.5 w-3.5" />
          {t("sidebar.logout")}
        </button>
      </div>
    </div>
  )

  return (
    <aside
      className={cn(
        "flex h-screen shrink-0 flex-col border-r border-border/70 bg-[hsl(var(--sidebar))] transition-all duration-200 ease-in-out",
        sidebarCollapsed ? "w-14" : "w-[240px]",
      )}
    >
      <div className={cn("flex h-14 items-center px-3", sidebarCollapsed && "justify-center")}>
        <div className={cn("flex items-center gap-2.5 overflow-hidden", sidebarCollapsed && "justify-center")}>
          <img src="/images/logo-icon.png" alt="Anybackup" className="h-6 w-6 shrink-0 object-contain" />
          {!sidebarCollapsed ? (
            <div className="flex items-baseline gap-1.5 whitespace-nowrap">
              <span className="text-[15px] font-semibold tracking-tight text-foreground">Anybackup</span>
              <span className="text-[13px] font-semibold text-[hsl(var(--ai))]">Agent</span>
            </div>
          ) : null}
        </div>
      </div>

      {!sidebarCollapsed ? (
        <div className="px-3 pb-2 pt-2">
          <button
            type="button"
            onClick={handleNewConversation}
            title={t("sidebar.newConversation")}
            className={cn(
              "focus-ring flex h-7 w-full items-center gap-2 rounded-sm border border-border/60 bg-white/75 px-2.5 text-xs font-medium text-foreground transition-fast hover:bg-white hover:shadow-card",
              selectedWorkspace?.kind === "localDraft" && "bg-white shadow-sm",
            )}
          >
            <MessageSquarePlus className="h-3 w-3 shrink-0 text-ai" />
            <span>{t("sidebar.newConversation")}</span>
          </button>
        </div>
      ) : null}

      {!sidebarCollapsed ? (
        <div className="px-3 pb-1.5">
          <label className="relative flex h-7 items-center rounded-sm border border-input bg-background transition-fast focus-within:border-ai/30 focus-within:ring-2 focus-within:ring-ai/10">
            <Search className="absolute left-2 top-1/2 h-3 w-3 shrink-0 -translate-y-1/2 text-muted-foreground" />
            <input
              value={query}
              onChange={(event) => void setSearchQuery(event.target.value)}
              placeholder={t("sidebar.searchConversation")}
              className="min-w-0 flex-1 border-0 bg-transparent pl-6 pr-2 text-xs text-foreground outline-none placeholder:text-muted-foreground"
            />
          </label>
        </div>
      ) : null}

      <div className="flex min-h-0 flex-1 flex-col overflow-hidden">
        <div className="scrollbar-sidebar flex min-h-0 flex-1 flex-col overflow-y-auto px-2 pb-3">
          {sidebarCollapsed ? (
            <div className="flex flex-col items-center gap-2 pt-1">
              <button
                type="button"
                onClick={handleNewConversation}
                aria-label={t("sidebar.newConversation")}
                title={t("sidebar.newConversation")}
                className={cn(
                  "focus-ring flex h-10 w-10 items-center justify-center rounded-lg text-ai shadow-card transition-fast hover:bg-ai-surface",
                  selectedWorkspace?.kind === "localDraft" ? "bg-ai-surface" : "bg-white",
                )}
              >
                <MessageSquarePlus className="h-4 w-4" />
              </button>
              {conversations.map((conversation) => renderHistoryButton(conversation, true))}
            </div>
          ) : listLoading ? (
            <div className="mx-1 flex items-center gap-2 rounded-lg border border-border bg-white/85 px-3 py-2 text-xs text-muted-foreground shadow-card">
              <LoaderCircle className="h-4 w-4 animate-spin" />
              {t("sidebar.refreshingConversations")}
            </div>
          ) : conversations.length === 0 ? (
            <div className="mx-1 rounded-md border border-dashed border-border bg-white/60 px-3 py-4 text-xs text-muted-foreground">
              {query.trim() ? t("sidebar.noMatchedConversation") : t("sidebar.noConversationYet")}
            </div>
          ) : (
            <section>
              <div className="px-1 pb-2 text-xs font-semibold tracking-wider text-muted-foreground">
                {t("sidebar.historyConversation")}
                <span className="ml-1 font-normal text-muted-foreground/40">
                  ({conversations.length})
                </span>
              </div>
              <div className="space-y-0.5 px-1">{conversations.map((conversation) => renderHistoryButton(conversation))}</div>
            </section>
          )}
        </div>
      </div>

      <div className="flex flex-col gap-1 py-3">
        {sidebarCollapsed ? (
          <>
            <div className="relative px-2" ref={menuRef}>
              <button
                type="button"
                onClick={() => setUserMenuOpen((open) => !open)}
                title={displayName}
                className="focus-ring flex h-10 w-full items-center justify-center rounded-lg bg-white/90 transition-fast hover:bg-white"
              >
                <div className="flex h-7 w-7 items-center justify-center rounded-full bg-gradient-to-br from-[hsl(var(--ai))] to-[hsl(var(--ai)/0.7)] text-primary-foreground">
                  <User className="h-3.5 w-3.5" />
                </div>
              </button>

              {userMenuOpen ? renderUserMenuCard("absolute bottom-0 left-[calc(100%+8px)]") : null}
            </div>
            <div className="px-2">
              <button
                type="button"
                onClick={toggleSidebar}
                title={t("sidebar.expandSidebar")}
                className="focus-ring flex h-10 w-full items-center justify-center rounded-lg text-muted-foreground transition-all duration-200 hover:bg-white/60 hover:text-foreground"
              >
                <ChevronRight className="h-4 w-4" />
              </button>
            </div>
            <div className="px-2">
              <button
                type="button"
                onClick={() => navigate(routes.settings)}
                title={t("sidebar.settings")}
                className="focus-ring flex h-10 w-full items-center justify-center rounded-lg text-muted-foreground transition-all duration-200 hover:bg-white/60 hover:text-foreground"
              >
                <Settings className="h-4 w-4" />
              </button>
            </div>
          </>
        ) : (
          <div className="px-2">
            <div className="mb-1.5 flex items-center gap-1">
              <div className="relative min-w-0 flex-1" ref={menuRef}>
                <button
                  type="button"
                  onClick={() => setUserMenuOpen((open) => !open)}
                  className={cn(
                    "focus-ring flex h-9 w-full items-center gap-3 rounded-lg px-3 transition-all duration-200",
                    userMenuOpen ? "bg-white shadow-sm" : "hover:bg-white/60",
                  )}
                >
                  <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-gradient-to-br from-[hsl(var(--ai))] to-[hsl(var(--ai)/0.7)] text-primary-foreground">
                    <User className="h-3.5 w-3.5" />
                  </div>
                  <div className="min-w-0 flex-1 text-left">
                    <span className="block truncate text-sm font-medium text-foreground">{displayName}</span>
                    <span className="block truncate text-xs text-muted-foreground">{roleLabel}</span>
                  </div>
                </button>

                {userMenuOpen ? renderUserMenuCard("fixed bottom-4 left-[248px]") : null}
              </div>
              <button
                type="button"
                onClick={toggleSidebar}
                aria-label={t("sidebar.collapseSidebar")}
                className="focus-ring flex h-9 w-9 shrink-0 items-center justify-center rounded-lg text-muted-foreground transition-all duration-200 hover:bg-white/60 hover:text-foreground"
              >
                <ChevronLeft className="h-4 w-4" />
              </button>
            </div>

            <button
              type="button"
              onClick={() => navigate(routes.settings)}
              className="focus-ring flex h-9 w-full items-center gap-3 rounded-lg px-3 text-muted-foreground transition-all duration-200 hover:bg-white/60 hover:text-foreground"
            >
              <Settings className="h-4 w-4 shrink-0" />
              <span className="text-sm font-medium">{t("sidebar.settings")}</span>
            </button>
          </div>
        )}
      </div>
    </aside>
  )
}
