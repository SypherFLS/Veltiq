import dayjs from 'dayjs'
import 'dayjs/locale/ru'
import relativeTime from 'dayjs/plugin/relativeTime'

dayjs.extend(relativeTime)
dayjs.locale('ru')

const moneyFormatter = new Intl.NumberFormat('ru-RU', {
  style: 'currency',
  currency: 'RUB',
  maximumFractionDigits: 0,
})

const numberFormatter = new Intl.NumberFormat('ru-RU')

export function formatMoney(value: number | null | undefined): string {
  if (value == null || Number.isNaN(value)) return '—'
  return moneyFormatter.format(value)
}

export function formatNumber(value: number | null | undefined): string {
  if (value == null || Number.isNaN(value)) return '—'
  return numberFormatter.format(value)
}

export function formatDate(
  value: string | Date | null | undefined,
  format = 'DD.MM.YYYY',
): string {
  if (!value) return '—'
  return dayjs(value).format(format)
}

export function formatDateTime(value: string | Date | null | undefined): string {
  return formatDate(value, 'DD.MM.YYYY HH:mm')
}

export function formatRelative(
  value: string | Date | null | undefined,
): string {
  if (!value) return '—'
  return dayjs(value).fromNow()
}
