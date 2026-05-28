<script setup lang="ts">
import { Bar } from 'vue-chartjs'
import {
  BarElement,
  CategoryScale,
  Chart as ChartJS,
  Legend,
  LinearScale,
  Tooltip,
  type ChartOptions,
} from 'chart.js'
import { formatMoney } from '~/shared/utils/format'

ChartJS.register(CategoryScale, LinearScale, BarElement, Tooltip, Legend)

interface Props {
  cashSum: number
  cardSum: number
  totalSum: number
  loading?: boolean
}

const props = defineProps<Props>()

const chartData = computed(() => ({
  labels: ['Наличные', 'Безнал', 'Всего'],
  datasets: [
    {
      label: 'Сумма продаж',
      data: [props.cashSum, props.cardSum, props.totalSum],
      backgroundColor: ['#34d399', '#60a5fa', '#1c5ef5'],
      borderRadius: 8,
      borderSkipped: false,
      maxBarThickness: 64,
    },
  ],
}))

const chartOptions = computed<ChartOptions<'bar'>>(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
      callbacks: {
        label: (ctx) => formatMoney(ctx.parsed.y),
      },
    },
  },
  scales: {
    x: {
      grid: { display: false },
      ticks: { color: '#6b7280' },
    },
    y: {
      grid: { color: '#eef0f4' },
      ticks: {
        color: '#6b7280',
        callback: (value) => formatMoney(Number(value)),
      },
    },
  },
}))
</script>

<template>
  <div class="h-72">
    <div v-if="loading" class="flex h-full items-center justify-center">
      <VSkeleton height="100%" rounded="lg" />
    </div>
    <Bar
      v-else
      :data="chartData"
      :options="chartOptions"
    />
  </div>
</template>
