import type { ImportStatus } from './import'

export interface ReportSummary {
  receiptsCount: number
  totalSum: number
  cashSum: number
  cardSum: number
  isStub: boolean
  note?: string
}

export interface Report {
  importId: string
  status: ImportStatus
  data: ReportSummary
  createdAt?: string
  updatedAt?: string
  errorCode?: string
}

export interface IlliquidItem {
  sku: string
  name: string
  category?: string
  stock: number
  daysWithoutSale: number
  lastSaleAt?: string
  recommendation?: 'discount' | 'bundle' | 'writeoff' | 'monitor'
  recommendationNote?: string
}

export interface ReportInsightsResponse {
  importId: string
  generatedAt: string
  items: IlliquidItem[]
}
