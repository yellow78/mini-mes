<template>
  <div class="spc-mini-chart" ref="chartRef" />
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'
import * as echarts from 'echarts/core'
import { LineChart } from 'echarts/charts'
import {
  GridComponent,
  MarkLineComponent,
  TooltipComponent,
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import type { SpcDataPoint } from '../../types/mes'

echarts.use([LineChart, GridComponent, MarkLineComponent, TooltipComponent, CanvasRenderer])

const props = defineProps<{
  data: SpcDataPoint[]
  parameter: string   // 'temperature' | 'pressure'
}>()

const chartRef = ref<HTMLElement | null>(null)
let chartInstance: echarts.ECharts | null = null

function buildOption() {
  const values = props.data.map(d => d.value)
  const ucl    = props.data[0]?.ucl ?? 0
  const lcl    = props.data[0]?.lcl ?? 0
  const labels = props.data.map((_, i) => i + 1)

  return {
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'axis',
      backgroundColor: '#1e293b',
      borderColor: '#334155',
      textStyle: { color: '#e2e8f0', fontSize: 12 },
      formatter: (params: any[]) => {
        const v = params[0]?.value ?? 0
        const isOut = v > ucl || v < lcl
        return `點位 ${params[0]?.axisValue}<br/>數值: <b style="color:${isOut ? '#ef4444' : '#22c55e'}">${v}</b>`
      },
    },
    grid: { top: 8, right: 8, bottom: 20, left: 40 },
    xAxis: {
      type: 'category',
      data: labels,
      axisLabel: { color: '#64748b', fontSize: 10 },
      axisLine: { lineStyle: { color: '#334155' } },
    },
    yAxis: {
      type: 'value',
      axisLabel: { color: '#64748b', fontSize: 10 },
      axisLine: { show: false },
      splitLine: { lineStyle: { color: '#1e293b' } },
    },
    series: [
      {
        type: 'line',
        data: values,
        smooth: true,
        lineStyle: { width: 1.5, color: '#3b82f6' },
        itemStyle: {
          color: (params: any) => {
            const v = params.value
            return v > ucl || v < lcl ? '#ef4444' : '#3b82f6'
          },
        },
        symbol: 'circle',
        symbolSize: 4,
        markLine: {
          silent: true,
          symbol: 'none',
          label: { fontSize: 10 },
          data: [
            { yAxis: ucl, lineStyle: { color: '#ef4444', type: 'dashed', width: 1 }, label: { formatter: `UCL ${ucl}`, color: '#ef4444' } },
            { yAxis: lcl, lineStyle: { color: '#f59e0b', type: 'dashed', width: 1 }, label: { formatter: `LCL ${lcl}`, color: '#f59e0b' } },
          ],
        },
      },
    ],
  }
}

function init() {
  if (!chartRef.value) return
  chartInstance = echarts.init(chartRef.value, 'dark')
  chartInstance.setOption(buildOption())
}

function resize() {
  chartInstance?.resize()
}

watch(() => props.data, () => {
  chartInstance?.setOption(buildOption())
}, { deep: true })

onMounted(() => {
  init()
  window.addEventListener('resize', resize)
})

onUnmounted(() => {
  chartInstance?.dispose()
  window.removeEventListener('resize', resize)
})
</script>

<style scoped>
.spc-mini-chart {
  width: 100%;
  height: 160px;
}
</style>
