import { defineStore } from 'pinia'

interface UIState {
  sidebarOpen: boolean
  mobileSidebarOpen: boolean
}

export const useUIStore = defineStore('ui', {
  state: (): UIState => ({
    sidebarOpen: true,
    mobileSidebarOpen: false,
  }),

  actions: {
    toggleSidebar() {
      this.sidebarOpen = !this.sidebarOpen
    },
    openMobileSidebar() {
      this.mobileSidebarOpen = true
    },
    closeMobileSidebar() {
      this.mobileSidebarOpen = false
    },
  },
})
