<template>
  <div class="chart-card">
    <div class="chart-header">
      <h3>问题出现频率</h3>
      <p>按严重程度统计当前问题分布</p>
    </div>

    <div class="chart-area">
      <div ref="chartRef" class="chart-box"></div>
    </div>

    <div class="legend-list">
      <div class="legend-item" v-for="item in legendItems" :key="item.level">
        <span class="dot" :class="item.className"></span>
        <span>{{ item.name }}（{{ item.value }}）</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import * as echarts from 'echarts'
import { getIssueStatisticsReport } from '@/api/report'

const chartRef = ref(null)
let chart = null

const levelConfig = [
  { level: 'serious', name: '严重', color: '#ef4444', className: 'red' },
  { level: 'normal', name: '一般', color: '#f59e0b', className: 'yellow' },
  { level: 'minor', name: '轻微', color: '#3b82f6', className: 'blue' },
  { level: 'suggestion', name: '建议', color: '#64748b', className: 'gray' }
]

const issueList = ref([])

const legendItems = computed(() => {
  return levelConfig.map(config => ({
    ...config,
    value: issueList.value.filter(item => item.issueLevel === config.level).length
  }))
})

const chartData = computed(() => {
  return legendItems.value.map(item => ({
    value: item.value,
    name: item.name,
    itemStyle: { color: item.color }
  }))
})

async function loadIssueStatistics() {
  try {
    const res = await getIssueStatisticsReport()
    const result = res?.data || res
    if (result.code !== 200) return
    issueList.value = (result.data || []).map(item => ({
      issueLevel: item.issueLevel || 'normal'
    }))
  } catch (err) {
    console.error('加载问题图表失败：', err)
  }
}

function initChart() {
  if (!chartRef.value) return

  chart = echarts.init(chartRef.value)

  chart.setOption({
    tooltip: {
      trigger: 'item',
      backgroundColor: '#020617',
      borderColor: '#334155',
      borderWidth: 1,
      padding: [8, 12],
      textStyle: {
        color: '#f8fafc',
        fontSize: 12
      },
      formatter: '{b}<br/>数量：{c}<br/>占比：{d}%'
    },

    series: [
      {
        name: '问题级别',
        type: 'pie',

        radius: ['42%', '62%'],
        center: ['50%', '50%'],

        itemStyle: {
          borderRadius: 8,
          borderColor: '#0f172a',
          borderWidth: 4
        },

        // 普通状态不显示文字
        label: {
          show: false
        },

        // 普通状态不显示引导线
        labelLine: {
          show: false
        },

        // 鼠标悬浮时，只放大扇区，不显示外圈文字
        emphasis: {
          scale: true,
          scaleSize: 6,
          label: {
            show: false
          },
          labelLine: {
            show: false
          }
        },

        data: chartData.value
      }
    ]
  })

  chart.resize()
}

function resizeChart() {
  if (chart) {
    chart.resize()
  }
}

onMounted(async () => {
  await loadIssueStatistics()
  await nextTick()
  initChart()
  window.addEventListener('resize', resizeChart)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', resizeChart)

  if (chart) {
    chart.dispose()
    chart = null
  }
})
</script>

<style scoped>
.chart-card {
  width: 100%;
  height: 420px;
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 14px;
  padding: 20px;
  color: #f8fafc;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.chart-header {
  height: 52px;
  flex-shrink: 0;
}

.chart-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
}

.chart-header p {
  margin: 6px 0 0;
  color: #94a3b8;
  font-size: 13px;
}

.chart-area {
  height: 280px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chart-box {
  width: 100%;
  height: 280px;
}

.legend-list {
  height: 48px;
  flex-shrink: 0;
  display: grid;
  grid-template-columns: repeat(4, auto);
  justify-content: center;
  align-items: center;
  column-gap: 28px;
  row-gap: 8px;
  padding-top: 8px;
}

.legend-item {
  display: flex;
  align-items: center;
  color: #94a3b8;
  font-size: 12px;
  white-space: nowrap;
}

.dot {
  width: 9px;
  height: 9px;
  border-radius: 999px;
  margin-right: 7px;
}

.red {
  background: #ef4444;
}

.yellow {
  background: #f59e0b;
}

.blue {
  background: #3b82f6;
}

.gray {
  background: #64748b;
}

@media (max-width: 768px) {
  .chart-card {
    height: 380px;
  }

  .chart-area {
    height: 230px;
  }

  .chart-box {
    height: 230px;
  }

  .legend-list {
    grid-template-columns: repeat(2, auto);
    height: 70px;
  }
}
</style>
