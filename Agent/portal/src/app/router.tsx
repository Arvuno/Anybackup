import { createBrowserRouter, Navigate, useLocation } from "react-router-dom"
import { AppShell } from "@/components/layout/AppShell"
import { ChatStatesDemoPage } from "@/pages/ChatStatesDemoPage"
import { HomePage } from "@/pages/HomePage"
import { SettingsPage } from "@/pages/SettingsPage"
import { UserManagementPage } from "@/pages/UserManagementPage"
import { LoginPage } from "@/pages/LoginPage"
import { GuidePage } from "@/pages/GuidePage"
import { isAuthenticated, isGuideDone } from "@/lib/session"
import { routes } from "@/config/routes"
import { routerBasenameFromBasePath } from "@/config/base-path"

function loginRedirectFor(pathname: string, search: string) {
  const returnTo = encodeURIComponent(`${pathname}${search}`)
  return `${routes.login}?returnTo=${returnTo}`
}

function ProtectedLayout() {
  const location = useLocation()

  if (!isAuthenticated()) return <Navigate to={loginRedirectFor(location.pathname, location.search)} replace />
  if (!isGuideDone()) return <Navigate to={routes.guide} replace />
  return <AppShell />
}

function GuideRoute() {
  if (!isAuthenticated()) return <Navigate to={routes.login} replace />
  return <GuidePage />
}

export const router = createBrowserRouter(
  [
    { path: routes.login, element: <LoginPage /> },
    { path: routes.guide, element: <GuideRoute /> },
    {
      path: "/",
      element: <ProtectedLayout />,
      children: [
        { index: true, element: <HomePage /> },
        { path: "chat-demo", element: <ChatStatesDemoPage /> },
        { path: "settings", element: <SettingsPage /> },
        { path: "settings/users", element: <UserManagementPage /> },
      ],
    },
  ],
  { basename: routerBasenameFromBasePath(import.meta.env.BASE_URL) },
)
