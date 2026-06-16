<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>版本状态矩阵</h1>
      </div>

      <div class="header-actions">
        <button class="reset-btn" @click="exportVersionMatrix">
          导出 Excel
        </button>
      </div>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索项目 / 终端类型 / 硬件版本 / 软件版本"
      />

      <select v-model="filters.projectName">
        <option value="">全部项目</option>
        <option
          v-for="project in projectOptions"
          :key="project"
          :value="project"
        >
          {{ project }}
        </option>
      </select>

      <select v-model="filters.deviceType">
        <option value="">全部终端类型</option>
        <option
          v-for="type in deviceTypeOptions"
          :key="type"
          :value="type"
        >
          {{ type }}
        </option>
      </select>

      <button class="query-btn" @click="loadVersionMatrix">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <!-- 版本矩阵表格 -->
    <div class="table-card">
      <div class="table-card-header">
        <div>
          <h3>版本状态矩阵列表</h3>
          <span>共 {{ filteredMatrixList.length }} 条版本矩阵记录</span>
        </div>
      </div>

      <div class="table-wrapper">
        <table class="version-table">
          <thead>
          <tr>
            <th>终端类型</th>
            <th>项目名称</th>
            <th>硬件版本</th>
            <th>软件版本</th>
            <th>最近更新时间</th>
          </tr>
          </thead>

          <tbody>
            <tr v-for="item in filteredMatrixList" :key="item.id">
              <td>
                <span class="device-tag">
                  {{ item.deviceType }}
                </span>
              </td>

              <td>
                <span class="project-name">{{ item.projectName }}</span>
              </td>


              <td class="version-cell">
                <span class="hardware-tag" :title="item.hardwareVersion">
                  {{ item.hardwareVersion }}
                </span>
              </td>

              <td class="version-cell">
                <span class="software-tag" :title="item.softwareVersion">
                  {{ item.softwareVersion }}
                </span>
              </td>

              <td class="muted">{{ item.updateTime }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="table-footer">
      </div>
    </div>

  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getVersionMatrixReport } from '@/api/report'

const route = useRoute()
const filters = reactive({
  keyword: '',
  projectName: '',
  deviceType: ''
})

const matrixList = ref([])

onMounted(() => {
  syncKeywordFromRoute()
  loadVersionMatrix()
})

watch(
  () => route.query.keyword,
  () => {
    syncKeywordFromRoute()
  }
)

function syncKeywordFromRoute() {
  filters.keyword = String(route.query.keyword || '')
}

function getResponseData(res) {
  return res && res.data ? res.data : res
}

function normalizeMatrix(item) {
  return {
    id: item.id,
    projectName: item.projectName || '',
    deviceType: item.deviceType || '',
    hardwareVersion: item.hardwareVersion || '',
    softwareVersion: item.softwareVersion || '',
    updateTime: item.updateTime || '',
    owner: item.owner || item.ownerName || '',
    remark: item.remark || '',
    status: item.status || ''
  }
}

async function loadVersionMatrix() {
  try {
    const res = await getVersionMatrixReport()
    const result = getResponseData(res)
    if (result.code !== 200) {
      alert(result.msg || '加载版本矩阵失败')
      return
    }
    matrixList.value = (result.data || []).map(normalizeMatrix)
  } catch (err) {
    console.error('加载版本矩阵失败：', err)
    alert(err.response?.data || '加载版本矩阵失败')
  }
}

const projectOptions = computed(() => [...new Set(matrixList.value.map(item => item.projectName).filter(Boolean))])
const deviceTypeOptions = computed(() => [...new Set(matrixList.value.map(item => item.deviceType).filter(Boolean))])

const filteredMatrixList = computed(() => {
  return matrixList.value.filter(item => {
    const keywordMatch =
      !filters.keyword ||
      item.projectName.includes(filters.keyword) ||
      item.deviceType.includes(filters.keyword) ||
      item.hardwareVersion.includes(filters.keyword) ||
      item.softwareVersion.includes(filters.keyword) ||
      item.owner.includes(filters.keyword)

    const projectMatch =
      !filters.projectName || item.projectName === filters.projectName

    const deviceTypeMatch =
      !filters.deviceType || item.deviceType === filters.deviceType

    return keywordMatch && projectMatch && deviceTypeMatch
  })
})

function resetFilters() {
  filters.keyword = ''
  filters.projectName = ''
  filters.deviceType = ''
}

function exportVersionMatrix() {
  const header = [
    '终端类型', 
    '项目名称',
    '硬件版本',
    '软件版本',
    '最近更新时间',
    '负责人',
    '版本说明'
  ]

  const rows = filteredMatrixList.value.map(item => [
    item.deviceType,
    item.projectName,
    item.hardwareVersion,
    item.softwareVersion,
    item.updateTime,
    item.owner,
    item.remark || ''
  ])

  const csvContent = [header, ...rows]
    .map(row => row.map(cell => `"${String(cell).replace(/"/g, '""')}"`).join(','))
    .join('\n')

  const blob = new Blob(['\uFEFF' + csvContent], {
    type: 'text/csv;charset=utf-8;'
  })

  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = '版本状态矩阵.csv'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}
</script>

<style scoped>
.page {
  width: 100%;
  min-height: 100%;
  color: #f8fafc;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 22px;
}

.page-header h1 {
  margin: 0;
  font-size: 26px;
  font-weight: 800;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.primary-btn,
.query-btn,
.reset-btn {
  height: 36px;
  padding: 0 16px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
}

.primary-btn,
.query-btn {
  border: none;
  background: #2563eb;
  color: #fff;
}

.primary-btn:hover,
.query-btn:hover {
  background: #1d4ed8;
}

.reset-btn {
  border: 1px solid #334155;
  background: #1e293b;
  color: #cbd5e1;
}

.filter-card {
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 14px;
  padding: 16px;
  display: grid;
  grid-template-columns: 1.4fr 220px 220px 90px 90px;
  gap: 12px;
  margin-bottom: 20px;
}

.filter-card input,
.filter-card select {
  height: 36px;
  border: 1px solid #334155;
  border-radius: 8px;
  background: #020617;
  color: #e2e8f0;
  padding: 0 12px;
  outline: none;
}

.filter-card input::placeholder {
  color: #64748b;
}

.table-card {
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 14px;
  overflow: hidden;
}

.table-card-header {
  padding: 14px 16px;
  border-bottom: 1px solid #1e293b;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.table-card-header h3 {
  margin: 0;
  color: #f8fafc;
  font-size: 15px;
}

.table-card-header span {
  color: #64748b;
  font-size: 12px;
}

.table-wrapper {
  width: 100%;
  overflow-x: auto;
  overflow-y: hidden;
}

.table-wrapper::-webkit-scrollbar {
  height: 10px;
}

.table-wrapper::-webkit-scrollbar-track {
  background: #020617;
  border-radius: 999px;
}

.table-wrapper::-webkit-scrollbar-thumb {
  background: #334155;
  border-radius: 999px;
  border: 2px solid #020617;
}

.table-wrapper::-webkit-scrollbar-thumb:hover {
  background: #475569;
}

.table-wrapper::-webkit-scrollbar-button {
  display: none;
}

.table-wrapper {
  scrollbar-width: thin;
  scrollbar-color: #334155 #020617;
}

.version-table {
  width: 100%;
  min-width: 1100px;
  border-collapse: collapse;
  table-layout: fixed;
}

.version-table thead {
  background: #020617;
}

.version-table th {
  padding: 14px 16px;
  color: #94a3b8;
  font-size: 12px;
  font-weight: 600;
  text-align: left;
  border-bottom: 1px solid #1e293b;
  white-space: nowrap;
}

.version-table td {
  padding: 15px 16px;
  font-size: 13px;
  color: #e2e8f0;
  border-bottom: 1px solid #1e293b;
  vertical-align: middle;
}

.version-table tbody tr {
  transition: background 0.2s;
}

.version-table tbody tr:hover {
  background: #1e293b80;
}

.record-link {
  border: none;
  background: transparent;
  color: #60a5fa;
  font-size: 13px;
  font-weight: 800;
  cursor: pointer;
  padding: 0;
  text-align: left;
}

.record-link:hover {
  color: #93c5fd;
  text-decoration: underline;
}

.device-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  background: #0f766e33;
  color: #5eead4;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
}

.version-cell {
  overflow: hidden;
}

.hardware-tag,
.software-tag {
  display: inline-block;
  max-width: 180px;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  vertical-align: middle;
}

.hardware-tag {
  background: #9333ea33;
  color: #c084fc;
}

.software-tag {
  background: #16a34a33;
  color: #4ade80;
}

.muted {
  color: #94a3b8 !important;
}

.operation-col {
  width: 120px;
  text-align: right !important;
}

.action-group {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 12px;
}

.text-btn {
  border: none;
  background: transparent;
  color: #60a5fa;
  font-size: 13px;
  cursor: pointer;
  white-space: nowrap;
}

.text-btn:hover {
  color: #93c5fd;
}

.table-footer {
  padding: 12px 16px;
  color: #64748b;
  font-size: 12px;
}

/* 弹窗 */
.dialog-mask {
  position: fixed;
  inset: 0;
  background: rgba(2, 6, 23, 0.72);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 999;
  padding: 20px;
}

.dialog {
  width: 760px;
  max-width: 100%;
  background: #0f172a;
  border: 1px solid #334155;
  border-radius: 16px;
  color: #f8fafc;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.45);
}

.large-dialog {
  width: 920px;
}

.dialog-header {
  padding: 18px 20px;
  border-bottom: 1px solid #1e293b;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.dialog-header h3 {
  margin: 0;
  font-size: 18px;
}

.dialog-header button {
  border: none;
  background: transparent;
  color: #94a3b8;
  font-size: 24px;
  cursor: pointer;
}

.detail-card {
  padding: 20px;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 14px;
}

.detail-card div,
.remark-card {
  background: #020617;
  border: 1px solid #1e293b;
  border-radius: 10px;
  padding: 12px;
}

.detail-card span,
.remark-card span {
  display: block;
  color: #64748b;
  font-size: 12px;
  margin-bottom: 6px;
}

.detail-card strong {
  color: #f8fafc;
  font-size: 14px;
  word-break: break-all;
}

.remark-card {
  margin: 0 20px 20px;
}

.remark-card p {
  margin: 0;
  color: #cbd5e1;
  font-size: 13px;
  line-height: 1.6;
}

.dialog-footer {
  padding: 16px 20px;
  border-top: 1px solid #1e293b;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

@media (max-width: 960px) {
  .filter-card {
    grid-template-columns: 1fr;
  }

  .version-table {
    min-width: 1100px;
  }

  .detail-card {
    grid-template-columns: 1fr;
  }
}
</style>
