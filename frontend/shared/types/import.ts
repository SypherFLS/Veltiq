export type ImportStatus =
  | 'pending'
  | 'processing'
  | 'partial_failed'
  | 'done'
  | 'failed'

export interface ImportRecord {
  id: string
  status: ImportStatus
  fileName?: string
  createdAt?: string
  updatedAt?: string
  errorCode?: string
}

export interface ImportUploadResponse {
  id: string
  status: ImportStatus
}

export interface ImportStatusResponse {
  id: string
  status: ImportStatus
  errorCode?: string
  updatedAt?: string
}
