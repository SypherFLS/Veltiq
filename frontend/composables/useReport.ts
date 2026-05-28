import { useQuery } from '@tanstack/vue-query'
import type { MaybeRef } from 'vue'
import type { Report, ReportInsightsResponse } from '~/shared/types/report'

export function useReport(id: MaybeRef<string>) {
  const { $api } = useNuxtApp()
  const idRef = computed(() => unref(id))

  return useQuery<Report>({
    queryKey: ['report', idRef],
    queryFn: () => $api.imports.report(idRef.value),
    enabled: computed(() => Boolean(idRef.value)),
  })
}

export function useReportInsights(id: MaybeRef<string>) {
  const { $api } = useNuxtApp()
  const idRef = computed(() => unref(id))

  return useQuery<ReportInsightsResponse>({
    queryKey: ['report-insights', idRef],
    queryFn: () => $api.imports.insights(idRef.value),
    enabled: computed(() => Boolean(idRef.value)),
    retry: 0,
  })
}
