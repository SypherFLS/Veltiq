import { defineStore } from 'pinia'

export type ToastKind = 'info' | 'success' | 'error' | 'warning'

export interface ToastItem {
  id: number
  kind: ToastKind
  title: string
  description?: string
  timeoutMs: number
}

interface ToastState {
  items: ToastItem[]
  nextId: number
}

interface PushOptions {
  description?: string
  timeoutMs?: number
}

export const useToastStore = defineStore('toast', {
  state: (): ToastState => ({
    items: [],
    nextId: 1,
  }),

  actions: {
    push(kind: ToastKind, title: string, opts: PushOptions = {}) {
      const id = this.nextId++
      const item: ToastItem = {
        id,
        kind,
        title,
        description: opts.description,
        timeoutMs: opts.timeoutMs ?? 5000,
      }
      this.items.push(item)
      if (item.timeoutMs > 0 && import.meta.client) {
        window.setTimeout(() => this.dismiss(id), item.timeoutMs)
      }
      return id
    },

    info(title: string, opts?: PushOptions) {
      return this.push('info', title, opts)
    },
    success(title: string, opts?: PushOptions) {
      return this.push('success', title, opts)
    },
    error(title: string, opts?: PushOptions) {
      return this.push('error', title, opts)
    },
    warning(title: string, opts?: PushOptions) {
      return this.push('warning', title, opts)
    },

    dismiss(id: number) {
      this.items = this.items.filter((t) => t.id !== id)
    },

    clear() {
      this.items = []
    },
  },
})
