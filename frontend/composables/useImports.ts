import { useQuery } from '@tanstack/vue-query'
import type { Paginated } from '~/shared/types/api'
import type { ImportRecord } from '~/shared/types/import'

export function useImportsList(limit = 20) {
  const { $api } = useNuxtApp()

  return useQuery<Paginated<ImportRecord>>({
    queryKey: ['imports', { limit }],
    queryFn: () => $api.imports.list({ limit }),
    retry: 0,
    placeholderData: { items: [], total: 0, cursor: null },
  })
}
