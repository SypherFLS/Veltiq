export function useToast() {
  const store = useToastStore()
  return {
    info: store.info,
    success: store.success,
    error: store.error,
    warning: store.warning,
    dismiss: store.dismiss,
  }
}
