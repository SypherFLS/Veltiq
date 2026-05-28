import type { ApiClient } from './client'
import type {
  ImportRecord,
  ImportStatusResponse,
  ImportUploadResponse,
} from '~/shared/types/import'
import type { Paginated } from '~/shared/types/api'
import type { Report, ReportInsightsResponse } from '~/shared/types/report'

export interface UploadParams {
  file: File
  onProgress?: (percent: number) => void
  signal?: AbortSignal
}

export interface ListParams {
  limit?: number
  cursor?: string | null
}

export function createImportsApi(client: ApiClient, baseURL: string) {
  return {
    list(params: ListParams = {}) {
      return client.get<Paginated<ImportRecord>>('/api/v1/imports', {
        query: { limit: params.limit ?? 20, cursor: params.cursor ?? undefined },
      })
    },

    upload({ file, onProgress, signal }: UploadParams) {
      const form = new FormData()
      form.append('file', file)

      if (!onProgress) {
        return client.post<ImportUploadResponse>('/api/v1/imports', form, {
          signal,
        })
      }

      return new Promise<ImportUploadResponse>((resolve, reject) => {
        const xhr = new XMLHttpRequest()
        xhr.open('POST', `${baseURL}/api/v1/imports`)
        xhr.withCredentials = true

        xhr.upload.onprogress = (event) => {
          if (event.lengthComputable) {
            onProgress(Math.round((event.loaded / event.total) * 100))
          }
        }

        xhr.onload = () => {
          if (xhr.status >= 200 && xhr.status < 300) {
            try {
              resolve(JSON.parse(xhr.responseText) as ImportUploadResponse)
            } catch (err) {
              reject(err)
            }
          } else {
            reject({
              response: {
                status: xhr.status,
                _data: safeJson(xhr.responseText),
              },
            })
          }
        }

        xhr.onerror = () =>
          reject({ response: { status: xhr.status }, message: 'Network error' })

        if (signal) {
          signal.addEventListener('abort', () => xhr.abort(), { once: true })
        }

        xhr.send(form)
      })
    },

    status(id: string) {
      return client.get<ImportStatusResponse>(`/api/v1/imports/${id}/status`)
    },

    report(id: string) {
      return client.get<Report>(`/api/v1/imports/${id}/report`)
    },

    insights(id: string) {
      return client.get<ReportInsightsResponse>(
        `/api/v1/imports/${id}/insights`,
      )
    },
  }
}

function safeJson(text: string): unknown {
  try {
    return JSON.parse(text)
  } catch {
    return { message: text }
  }
}

export type ImportsApi = ReturnType<typeof createImportsApi>
