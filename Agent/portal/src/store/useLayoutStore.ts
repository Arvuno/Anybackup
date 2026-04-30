import { create } from "zustand"

interface LayoutState {
  sidebarCollapsed: boolean
}

interface LayoutActions {
  toggleSidebar: () => void
}

export const useLayoutStore = create<LayoutState & LayoutActions>((set) => ({
  sidebarCollapsed: false,

  toggleSidebar: () => set((s) => ({ sidebarCollapsed: !s.sidebarCollapsed })),
}))
