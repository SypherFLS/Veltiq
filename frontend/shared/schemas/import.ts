import { z } from 'zod'

export const MAX_IMPORT_FILE_BYTES = 25 * 1024 * 1024

export const ACCEPTED_IMPORT_MIME = [
  'text/csv',
  'application/vnd.ms-excel',
  'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
  'application/json',
  'text/xml',
  'application/xml',
]

export const importFileSchema = z
  .instanceof(File, { message: 'Выберите файл' })
  .refine((f) => f.size > 0, 'Файл пустой')
  .refine(
    (f) => f.size <= MAX_IMPORT_FILE_BYTES,
    `Максимальный размер ${(MAX_IMPORT_FILE_BYTES / 1024 / 1024).toFixed(0)} МБ`,
  )
