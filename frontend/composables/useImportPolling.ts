import { useQuery } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type { ImportStatusResponse } from '~/shared/types/import'

export function useImportPolling(id: MaybeRef<string>, intervalMs = 3000) {
  const { $api } = useNuxtApp()
  const idRef = computed(() => unref(id))

  return useQuery<ImportStatusResponse>({
    queryKey: ['import', idRef],
    queryFn: () => $api.imports.status(idRef.value),
    enabled: computed(() => Boolean(idRef.value)),
    refetchInterval: (query) => {
      const status = query.state.data?.status
      if (status === 'done' || status === 'failed' || status === 'partial_failed') {
        return false
      }
      return intervalMs
    },
    refetchIntervalInBackground: false,
  })
}
