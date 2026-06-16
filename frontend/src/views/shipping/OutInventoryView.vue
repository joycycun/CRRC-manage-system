<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1>出库记录管理</h1>
      </div>
    </div>

    <div class="summary-grid two-col">
      <div class="summary-card">
        <span>累计出库数量</span>
        <strong>{{ outboundCount }}</strong>
        <p>发货批次审核通过后自动生成的出库设备数量</p>
      </div>

      <div class="summary-card green">
        <span>当月出库数量</span>
        <strong>{{ currentMonthOutboundCount }}</strong>
        <p>本月已确认出库的设备数量</p>
      </div>
    </div>

    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索批次号 / 快递单号 / SN / MAC / 软件版本 / 硬件版本 / 出库人"
      />

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

      <button class="query-btn" @click="loadOutboundRecords">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <div class="type-card">
      <div class="section-title">
        <h3>终端类型出库统计</h3>
        <span>按终端类型统计已出库设备数量</span>
      </div>

      <div class="type-grid">
        <div
          v-for="item in deviceTypeSummary"
          :key="item.deviceType"
          class="type-item"
        >
          <span>{{ item.deviceType }}</span>
          <strong>{{ item.count }}</strong>
        </div>
      </div>
    </div>

    <div class="table-card">
      <div class="table-card-header">
        <div>
          <h3>出库记录列表</h3>
          <span>
            共 {{ filteredOutboundList.length }} 条，当前第 {{ currentPage }} / {{ totalPage }} 页
          </span>
        </div>

        <div class="page-size-control">
          <span>每页</span>
          <select v-model.number="pageSize">
            <option
              v-for="size in pageSizeOptions"
              :key="size"
              :value="size"
            >
              {{ size }}
            </option>
          </select>
          <span>条</span>
        </div>
      </div>

      <div class="table-wrapper">
        <table class="version-table">
          <thead>
            <tr>
              <th>发货批次</th>
              <th>终端类型</th>
              <th>SN序列号</th>
              <th>MAC地址</th>
              <th>软件版本</th>
              <th>硬件版本</th>
              <th>出库时间</th>
              <th>出库人</th>
              <th>快递单号</th>
              <th>状态</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="item in paginatedOutboundList" :key="item.outboundId">
              <td>
                <span class="batch-tag" :title="item.shippingBatchNo">
                  {{ item.shippingBatchNo || '-' }}
                </span>
              </td>

              <td>
                <span class="device-tag">{{ item.deviceType || '-' }}</span>
              </td>

              <td>
                <span class="sn-tag" :title="item.sn">
                  {{ item.sn }}
                </span>
              </td>

              <td>
                <span class="mac-text" :title="item.macAddress">
                  {{ item.macAddress }}
                </span>
              </td>

              <td class="version-cell">
                <span class="software-tag" :title="item.softwareVersion">
                  {{ item.softwareVersion }}
                </span>
              </td>

              <td class="version-cell">
                <span class="hardware-tag" :title="item.hardwareVersion">
                  {{ item.hardwareVersion }}
                </span>
              </td>

              <td class="muted">{{ item.outboundTime || '-' }}</td>

              <td>{{ item.operator || '-' }}</td>

              <td>
                <span class="normal-text" :title="item.expressNo">
                  {{ item.expressNo || '-' }}
                </span>
              </td>

              <td>
                <span class="status-tag" :class="getOutboundStatusClass(item.outboundStatus)">
                  {{ getOutboundStatusText(item.outboundStatus) }}
                </span>
              </td>

              <td class="operation-col">
                <button class="text-btn blue" @click="viewOutbound(item)">
                  查看
                </button>
              </td>
            </tr>

            <tr v-if="paginatedOutboundList.length === 0">
              <td colspan="11" class="empty-table">
                暂无符合条件的出库记录
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="pagination-bar">
        <div class="pagination-info">
          当前显示第 {{ pageStartIndex }} - {{ pageEndIndex }} 条，
          共 {{ filteredOutboundList.length }} 条出库记录
        </div>

        <div class="pagination-actions">
          <button
            class="page-btn"
            :disabled="currentPage === 1"
            @click="goFirstPage"
          >
            首页
          </button>

          <button
            class="page-btn"
            :disabled="currentPage === 1"
            @click="goPrevPage"
          >
            上一页
          </button>

          <span class="page-number">
            {{ currentPage }} / {{ totalPage }}
          </span>

          <button
            class="page-btn"
            :disabled="currentPage === totalPage"
            @click="goNextPage"
          >
            下一页
          </button>

          <button
            class="page-btn"
            :disabled="currentPage === totalPage"
            @click="goLastPage"
          >
            末页
          </button>
        </div>
      </div>

      <div class="table-footer">
        出库记录由发货批次审核通过后自动生成，不支持在此页面直接新增或删除。
      </div>
    </div>

    <div v-if="selectedOutbound" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>出库记录详情</h3>
          <button @click="selectedOutbound = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>发货批次</span>
            <strong>{{ selectedOutbound.shippingBatchNo || '-' }}</strong>
          </div>

          <div>
            <span>快递单号</span>
            <strong>{{ selectedOutbound.expressNo || '-' }}</strong>
          </div>

          <div>
            <span>终端类型</span>
            <strong>{{ selectedOutbound.deviceType || '-' }}</strong>
          </div>

          <div>
            <span>SN序列号</span>
            <strong>{{ selectedOutbound.sn || '-' }}</strong>
          </div>

          <div>
            <span>MAC地址</span>
            <strong>{{ selectedOutbound.macAddress || '-' }}</strong>
          </div>

          <div>
            <span>软件版本</span>
            <strong>{{ selectedOutbound.softwareVersion || '-' }}</strong>
          </div>

          <div>
            <span>硬件版本</span>
            <strong>{{ selectedOutbound.hardwareVersion || '-' }}</strong>
          </div>

          <div>
            <span>入库时间</span>
            <strong>{{ selectedOutbound.inTime || '-' }}</strong>
          </div>

          <div>
            <span>出库时间</span>
            <strong>{{ selectedOutbound.outboundTime || '-' }}</strong>
          </div>

          <div>
            <span>出库人</span>
            <strong>{{ selectedOutbound.operator || '-' }}</strong>
          </div>

          <div>
            <span>出库状态</span>
            <strong>{{ getOutboundStatusText(selectedOutbound.outboundStatus) }}</strong>
          </div>
        </div>

        <div class="remark-card">
          <span>出库说明</span>
          <p>{{ selectedOutbound.remark || '暂无说明' }}</p>
        </div>

        <div class="dialog-footer">
          <button class="primary-btn" @click="selectedOutbound = null">
            关闭
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getOutboundRecords } from '@/api/outbound'

const route = useRoute()
const filters = reactive({ keyword: '', deviceType: '' })
const selectedOutbound = ref(null)
const outboundList = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const pageSizeOptions = [10, 20, 50, 100]

onMounted(async () => {
  syncKeywordFromRoute()
  await loadOutboundRecords()
})

watch(
  () => route.query.keyword,
  () => {
    syncKeywordFromRoute()
    currentPage.value = 1
  }
)

function syncKeywordFromRoute() {
  filters.keyword = String(route.query.keyword || '')
}

function getResponseData(res) {
  return res && res.data ? res.data : res
}

function formatDateTime(value) {
  if (!value) return ''
  if (typeof value === 'string') return value.slice(0, 19).replace('T', ' ')
  if (value.Time) return String(value.Time).slice(0, 19).replace('T', ' ')
  return String(value)
}

function normalizeOutbound(item) {
  const status = item.outboundStatus || item.outbound_status || item.status || '已出库'

  return {
    outboundId: item.outboundId || item.outbound_id || item.detailId || item.detail_id || item.id,
    detailId: item.detailId || item.detail_id || item.id || item.outboundId || item.outbound_id,
    batchId: item.batchId || item.batch_id || 0,
    inventoryDeviceId: item.inventoryDeviceId || item.inventory_device_id || 0,
    deviceType: item.deviceType || item.device_type || '',
    sn: item.sn || '',
    macAddress: item.macAddress || item.mac_address || '',
    softwareVersion: item.softwareVersion || item.software_version || '',
    hardwareVersion: item.hardwareVersion || item.hardware_version || '',
    inTime: formatDateTime(item.inTime || item.in_time),
    outboundTime: formatDateTime(item.outboundTime || item.outbound_time),
    operator: item.operator || item.outboundUserName || item.outbound_user_name || item.uploaderName || item.uploader_name || '',
    shippingBatchNo: item.shippingBatchNo || item.shipping_batch_no || item.batchNo || item.batch_no || '',
    expressNo: item.expressNo || item.express_no || '',
    outboundStatus: normalizeOutboundStatus(status),
    remark: item.remark || ''
  }
}

function normalizeOutboundStatus(status) {
  const map = {
    shipped: '已出库',
    shipping: '已出库',
    pending: '待出库',
    outbound: '已出库',
    returned: '已退回',
    已出库: '已出库',
    待出库: '待出库',
    已退回: '已退回'
  }
  return map[status] || status || '已出库'
}

async function loadOutboundRecords() {
  try {
    const res = await getOutboundRecords()
    const result = getResponseData(res)
    console.log('出库记录返回：', result)

    if (result.code !== 200) {
      alert(result.msg || '加载出库记录失败')
      return
    }

    outboundList.value = (result.data || []).map(normalizeOutbound)
  } catch (err) {
    console.error('加载出库记录失败：', err)
    alert(err.response?.data || '加载出库记录失败，请检查 /api/outbound-records 接口')
  }
}

const completedOutboundList = computed(() => {
  return outboundList.value.filter(item => item.outboundStatus !== '已退回')
})

const deviceTypeOptions = computed(() => {
  return [...new Set(completedOutboundList.value.map(item => item.deviceType).filter(Boolean))]
})

const filteredOutboundList = computed(() => {
  return completedOutboundList.value.filter(item => {
    const keyword = filters.keyword.trim().toLowerCase()
    const keywordMatch =
      !keyword ||
      String(item.sn || '').toLowerCase().includes(keyword) ||
      String(item.macAddress || '').toLowerCase().includes(keyword) ||
      String(item.softwareVersion || '').toLowerCase().includes(keyword) ||
      String(item.hardwareVersion || '').toLowerCase().includes(keyword) ||
      String(item.operator || '').toLowerCase().includes(keyword) ||
      String(item.shippingBatchNo || '').toLowerCase().includes(keyword) ||
      String(item.expressNo || '').toLowerCase().includes(keyword)

    const deviceTypeMatch = !filters.deviceType || item.deviceType === filters.deviceType
    return keywordMatch && deviceTypeMatch
  })
})

const outboundCount = computed(() => completedOutboundList.value.length)

const currentMonthOutboundCount = computed(() => {
  const currentMonth = new Date().toISOString().slice(0, 7)
  return completedOutboundList.value.filter(item => item.outboundTime && item.outboundTime.startsWith(currentMonth)).length
})

const deviceTypeSummary = computed(() => {
  const map = new Map()
  completedOutboundList.value.forEach(item => {
    const type = item.deviceType || '未填写'
    map.set(type, (map.get(type) || 0) + 1)
  })
  return [...map.entries()].map(([deviceType, count]) => ({ deviceType, count }))
})

const totalPage = computed(() => Math.max(1, Math.ceil(filteredOutboundList.value.length / pageSize.value)))
const paginatedOutboundList = computed(() => {
  if (currentPage.value > totalPage.value) currentPage.value = totalPage.value
  const start = (currentPage.value - 1) * pageSize.value
  return filteredOutboundList.value.slice(start, start + pageSize.value)
})
const pageStartIndex = computed(() => filteredOutboundList.value.length === 0 ? 0 : (currentPage.value - 1) * pageSize.value + 1)
const pageEndIndex = computed(() => Math.min(currentPage.value * pageSize.value, filteredOutboundList.value.length))

function getOutboundStatusText(status) {
  const map = { 已出库: '已出库', 待出库: '待出库', 已退回: '已退回' }
  return map[status] || status || '已出库'
}

function getOutboundStatusClass(status) {
  const map = { 已出库: 'shipped', 待出库: 'pending', 已退回: 'returned' }
  return map[status] || 'shipped'
}

function resetFilters() {
  filters.keyword = ''
  filters.deviceType = ''
  currentPage.value = 1
}

function viewOutbound(item) {
  selectedOutbound.value = item
}

function goFirstPage() {
  currentPage.value = 1
}

function goPrevPage() {
  if (currentPage.value > 1) currentPage.value -= 1
}

function goNextPage() {
  if (currentPage.value < totalPage.value) currentPage.value += 1
}

function goLastPage() {
  currentPage.value = totalPage.value
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

.summary-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 20px;
}

.summary-card {
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 14px;
  padding: 18px;
}

.summary-card span {
  color: #94a3b8;
  font-size: 13px;
}

.summary-card strong {
  display: block;
  margin-top: 10px;
  color: #f8fafc;
  font-size: 30px;
  font-weight: 800;
}

.summary-card p {
  margin: 8px 0 0;
  color: #64748b;
  font-size: 12px;
}

.summary-card.green strong {
  color: #4ade80;
}

.type-card {
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 14px;
  padding: 18px;
  margin-bottom: 20px;
}

.section-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-title h3 {
  margin: 0;
  font-size: 16px;
  color: #f8fafc;
}

.section-title span {
  color: #64748b;
  font-size: 12px;
}

.type-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 12px;
}

.type-item {
  background: #020617;
  border: 1px solid #1e293b;
  border-radius: 10px;
  padding: 14px;
}

.type-item span {
  display: block;
  color: #94a3b8;
  font-size: 12px;
  margin-bottom: 8px;
}

.type-item strong {
  color: #f8fafc;
  font-size: 22px;
}

.primary-btn,
.query-btn,
.reset-btn,
.green-btn,
.red-btn,
.page-btn {
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

.reset-btn,
.page-btn {
  border: 1px solid #334155;
  background: #1e293b;
  color: #cbd5e1;
}

.page-btn:hover:not(:disabled) {
  background: #334155;
  color: #fff;
}

.page-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.filter-card {
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 14px;
  padding: 16px;
  display: grid;
  grid-template-columns: 1.4fr 220px 90px 90px;
  gap: 12px;
  margin-bottom: 20px;
}

.filter-card input,
.filter-card select,
.page-size-control select {
  border: 1px solid #334155;
  border-radius: 8px;
  background: #020617;
  color: #e2e8f0;
  padding: 0 12px;
  outline: none;
  height: 36px;
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
  font-size: 16px;
}

.table-card-header span {
  color: #64748b;
  font-size: 12px;
}

.page-size-control {
  display: flex;
  align-items: center;
  gap: 8px;
}

.page-size-control span {
  color: #94a3b8;
  font-size: 13px;
}

.page-size-control select {
  width: 90px;
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
  min-width: 1060px;
  border-collapse: collapse;
  table-layout: fixed;
}

.version-table thead {
  background: #020617;
}

.version-table th,
.version-table td {
  box-sizing: border-box;
  white-space: nowrap;
}

.version-table th {
  padding: 14px 16px;
  color: #94a3b8;
  font-size: 12px;
  font-weight: 600;
  text-align: left;
  border-bottom: 1px solid #1e293b;
}

.version-table td {
  padding: 15px 16px;
  font-size: 13px;
  color: #e2e8f0;
  border-bottom: 1px solid #1e293b;
  vertical-align: middle;
  overflow: hidden;
}

.operation-col {
  width: 120px;
  text-align: right !important;
}

.device-tag {
  display: inline-block;
  max-width: 130px;
  padding: 4px 9px;
  border-radius: 999px;
  background: #0f766e33;
  color: #5eead4;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.software-tag,
.hardware-tag {
  display: inline-block;
  max-width: 230px;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  vertical-align: middle;
}

.software-tag {
  background: #16a34a33;
  color: #4ade80;
}

.hardware-tag {
  background: #9333ea33;
  color: #c084fc;
}

.sn-tag {
  display: inline-block;
  max-width: 170px;
  padding: 3px 8px;
  border-radius: 999px;
  background: #33415566;
  color: #cbd5e1;
  font-size: 12px;
  font-family: Consolas, Monaco, monospace;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.mac-text {
  display: inline-block;
  max-width: 160px;
  color: #cbd5e1;
  font-size: 12px;
  font-family: Consolas, Monaco, monospace;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.batch-tag,
.normal-text {
  display: inline-block;
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.batch-tag {
  padding: 3px 8px;
  border-radius: 999px;
  background: #1d4ed833;
  color: #93c5fd;
  font-size: 12px;
  font-weight: 700;
}

.normal-text {
  color: #cbd5e1;
}

.version-cell {
  overflow: hidden;
}

.muted {
  color: #94a3b8 !important;
}

.text-btn {
  border: none;
  background: transparent;
  color: #cbd5e1;
  font-size: 13px;
  cursor: pointer;
  white-space: nowrap;
}

.text-btn:hover {
  color: #fff;
}

.text-btn.blue {
  color: #60a5fa;
}

.status-tag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 64px;
  height: 24px;
  padding: 0 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
}

.status-tag.shipped {
  background: #16a34a33;
  color: #4ade80;
}

.status-tag.pending {
  background: #f59e0b33;
  color: #fbbf24;
}

.status-tag.returned {
  background: #64748b33;
  color: #cbd5e1;
}

.empty-table {
  text-align: center;
  color: #64748b !important;
  padding: 28px 16px !important;
}

.pagination-bar {
  padding: 14px 16px;
  border-top: 1px solid #1e293b;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.pagination-info {
  color: #64748b;
  font-size: 12px;
}

.pagination-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.page-number {
  min-width: 70px;
  text-align: center;
  color: #cbd5e1;
  font-size: 13px;
}

.table-footer {
  padding: 12px 16px;
  color: #64748b;
  font-size: 12px;
  border-top: 1px solid #1e293b;
}

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

@media (max-width: 1200px) {
  .type-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 960px) {
  .filter-card {
    grid-template-columns: 1fr;
  }

  .summary-grid,
  .type-grid {
    grid-template-columns: 1fr;
  }

  .version-table {
    min-width: 1060px;
  }

  .table-card-header,
  .pagination-bar {
    align-items: flex-start;
    flex-direction: column;
  }

  .pagination-actions {
    flex-wrap: wrap;
  }

  .detail-card {
    grid-template-columns: 1fr;
  }
}
</style>
