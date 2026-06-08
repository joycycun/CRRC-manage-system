<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>库存情况管理</h1>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="summary-grid two-col">
      <div class="summary-card">
        <span>当前库存总数</span>
        <strong>{{ totalCount }}</strong>
        <p>仅统计已完成烧录、出厂测试且当前仍在库存中的设备</p>
      </div>

      <div class="summary-card green">
        <span>当月出库数量</span>
        <strong>{{ currentMonthOutboundCount }}</strong>
        <p>本月已确认出库的设备数量</p>
      </div>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索 SN / MAC / 软件版本 / 硬件版本"
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

      <button class="query-btn">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <!-- 终端类型库存统计 -->
    <div class="type-card">
      <div class="section-title">
        <h3>终端类型库存</h3>
        <span>按终端类型统计已完成烧录、出厂测试后的当前库存数量</span>
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

    <!-- 当前库存表格 -->
    <div class="table-card">
      <div class="table-card-header">
        <div>
          <h3>当前库存列表</h3>
          <span>
            共 {{ filteredInventoryList.length }} 条，当前第 {{ currentPage }} / {{ totalPage }} 页
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
              <th>终端类型</th>
              <th>SN序列号</th>
              <th>MAC地址</th>
              <th>软件版本</th>
              <th>硬件版本</th>
              <th>入库时间</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="item in paginatedInventoryList" :key="item.id">
              <td>
                <span class="device-tag">{{ item.deviceType }}</span>
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

              <td class="muted">{{ item.inTime }}</td>

              <td class="operation-col">
                <button class="text-btn blue" @click="viewInventory(item)">
                  查看
                </button>
              </td>
            </tr>

            <tr v-if="paginatedInventoryList.length === 0">
              <td colspan="7" class="empty-table">
                暂无符合条件的库存设备
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="pagination-bar">
        <div class="pagination-info">
          当前显示第 {{ pageStartIndex }} - {{ pageEndIndex }} 条，
          共 {{ filteredInventoryList.length }} 条库存设备记录
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
        库存页面仅展示已完成烧录、出厂测试并处于可发货库存状态的设备。出库操作请在出库记录页面中处理。
      </div>
    </div>

    <!-- 查看库存详情弹窗 -->
    <div v-if="selectedInventory" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>库存设备详情</h3>
          <button @click="selectedInventory = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>终端类型</span>
            <strong>{{ selectedInventory.deviceType }}</strong>
          </div>

          <div>
            <span>SN序列号</span>
            <strong>{{ selectedInventory.sn }}</strong>
          </div>

          <div>
            <span>MAC地址</span>
            <strong>{{ selectedInventory.macAddress }}</strong>
          </div>

          <div>
            <span>软件版本</span>
            <strong>{{ selectedInventory.softwareVersion }}</strong>
          </div>

          <div>
            <span>硬件版本</span>
            <strong>{{ selectedInventory.hardwareVersion }}</strong>
          </div>

          <div>
            <span>入库时间</span>
            <strong>{{ selectedInventory.inTime }}</strong>
          </div>

          <div>
            <span>库存状态</span>
            <strong>{{ getInventoryStatusText(selectedInventory.inventoryStatus) }}</strong>
          </div>

          <div>
            <span>最后更新时间</span>
            <strong>{{ selectedInventory.updateTime }}</strong>
          </div>
        </div>

        <div class="remark-card">
          <span>库存说明</span>
          <p>{{ selectedInventory.remark || '暂无说明' }}</p>
        </div>

        <div class="dialog-footer">
          <button class="primary-btn" @click="selectedInventory = null">
            关闭
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'

const STORAGE_INVENTORY_KEY = 'inventoryList'
const STORAGE_OUTBOUND_KEY = 'outInventoryList'

const filters = reactive({
  keyword: '',
  deviceType: ''
})

const selectedInventory = ref(null)

const currentPage = ref(1)
const pageSize = ref(10)

const pageSizeOptions = [10, 20, 50, 100]

const deviceTypeOptions = [
  '司机室控制盒',
  '解码板',
  '司机提醒单元',
  '噪声检测',
  '编码板',
  '功放板'
]

const defaultInventoryList = [
  {
    id: 1,
    deviceType: '司机室控制盒',
    sn: 'DCCU-202605100001',
    macAddress: '00:11:22:33:44:01',
    softwareVersion: 'SW-DCCU.V1.2.0',
    hardwareVersion: 'HD-DCCU.V1.1.0',
    inventoryStatus: 'ready',
    inTime: '2026-05-10',
    updateTime: '2026-05-11',
    remark: '已完成烧录和出厂测试，可发货。'
  },
  {
    id: 2,
    deviceType: '司机室控制盒',
    sn: 'DCCU-202605100002',
    macAddress: '00:11:22:33:44:02',
    softwareVersion: 'SW-DCCU.V1.2.0',
    hardwareVersion: 'HD-DCCU.V1.1.0',
    inventoryStatus: 'ready',
    inTime: '2026-05-10',
    updateTime: '2026-05-11',
    remark: '已完成烧录和出厂测试，可发货。'
  },
  {
    id: 3,
    deviceType: '司机室控制盒',
    sn: 'DCCU-202605100003',
    macAddress: '00:11:22:33:44:03',
    softwareVersion: 'SW-DCCU.V1.2.0',
    hardwareVersion: 'HD-DCCU.V1.1.0',
    inventoryStatus: 'ready',
    inTime: '2026-05-10',
    updateTime: '2026-05-10'
  },
  {
    id: 4,
    deviceType: '解码板',
    sn: 'DEC-202605110001',
    macAddress: '00:11:22:55:66:01',
    softwareVersion: 'SW-DECODER.V1.0.3',
    hardwareVersion: 'HD-DECODER.V1.1.0',
    inventoryStatus: 'ready',
    inTime: '2026-05-11',
    updateTime: '2026-05-12'
  },
  {
    id: 5,
    deviceType: '解码板',
    sn: 'DEC-202605110002',
    macAddress: '00:11:22:55:66:02',
    softwareVersion: 'SW-DECODER.V1.0.3',
    hardwareVersion: 'HD-DECODER.V1.1.0',
    inventoryStatus: 'ready',
    inTime: '2026-05-11',
    updateTime: '2026-05-12'
  },
  {
    id: 6,
    deviceType: '解码板',
    sn: 'DEC-202605110003',
    macAddress: '00:11:22:55:66:03',
    softwareVersion: 'SW-DECODER.V1.0.2',
    hardwareVersion: 'HD-DECODER.V1.0.0',
    inventoryStatus: 'ready',
    inTime: '2026-05-11',
    updateTime: '2026-05-13'
  },
  {
    id: 7,
    deviceType: '司机提醒单元',
    sn: 'DRU-202605120001',
    macAddress: '00:11:22:77:88:01',
    softwareVersion: 'SW-DRU.V1.0.0',
    hardwareVersion: 'HD-DRU.V1.0.0',
    inventoryStatus: 'ready',
    inTime: '2026-05-12',
    updateTime: '2026-05-13'
  },
  {
    id: 8,
    deviceType: '司机提醒单元',
    sn: 'DRU-202605120002',
    macAddress: '00:11:22:77:88:02',
    softwareVersion: 'SW-DRU.V1.0.0',
    hardwareVersion: 'HD-DRU.V1.0.0',
    inventoryStatus: 'waiting_burn',
    inTime: '2026-05-12',
    updateTime: '2026-05-12',
    remark: '未完成烧录，不在库存页面展示。'
  },
  {
    id: 9,
    deviceType: '噪声检测',
    sn: 'NOISE-202605130001',
    macAddress: '00:11:22:99:AA:01',
    softwareVersion: 'SW-NOISE.V1.1.0',
    hardwareVersion: 'HD-NOISE.V1.0.1',
    inventoryStatus: 'ready',
    inTime: '2026-05-13',
    updateTime: '2026-05-14'
  },
  {
    id: 10,
    deviceType: '噪声检测',
    sn: 'NOISE-202605130002',
    macAddress: '00:11:22:99:AA:02',
    softwareVersion: 'SW-NOISE.V1.1.0',
    hardwareVersion: 'HD-NOISE.V1.0.1',
    inventoryStatus: 'repair',
    inTime: '2026-05-13',
    updateTime: '2026-05-15',
    remark: '返修库存，不在当前库存页面展示。'
  },
  {
    id: 11,
    deviceType: '编码板',
    sn: 'ENC-202605140001',
    macAddress: '00:11:22:BB:CC:01',
    softwareVersion: 'SW-ENCODER.V1.3.0',
    hardwareVersion: 'HD-ENCODER.V1.2.0',
    inventoryStatus: 'ready',
    inTime: '2026-05-14',
    updateTime: '2026-05-15'
  },
  {
    id: 12,
    deviceType: '编码板',
    sn: 'ENC-202605140002',
    macAddress: '00:11:22:BB:CC:02',
    softwareVersion: 'SW-ENCODER.V1.3.0',
    hardwareVersion: 'HD-ENCODER.V1.2.0',
    inventoryStatus: 'repair',
    inTime: '2026-05-14',
    updateTime: '2026-05-14',
    remark: '返修库存，不在当前库存页面展示。'
  },
  {
    id: 13,
    deviceType: '功放板',
    sn: 'AMP-202605150001',
    macAddress: '00:11:22:DD:EE:01',
    softwareVersion: 'SW-AMP.V1.0.5',
    hardwareVersion: 'HD-AMP.V1.0.2',
    inventoryStatus: 'ready',
    inTime: '2026-05-15',
    updateTime: '2026-05-16'
  },
  {
    id: 14,
    deviceType: '功放板',
    sn: 'AMP-202605150002',
    macAddress: '00:11:22:DD:EE:02',
    softwareVersion: 'SW-AMP.V1.0.5',
    hardwareVersion: 'HD-AMP.V1.0.2',
    inventoryStatus: 'ready',
    inTime: '2026-05-15',
    updateTime: '2026-05-16'
  }
]

function readStorageList(key, fallback) {
  try {
    const raw = localStorage.getItem(key)
    if (!raw) return fallback
    return JSON.parse(raw)
  } catch (error) {
    return fallback
  }
}

function saveStorageList(key, value) {
  localStorage.setItem(key, JSON.stringify(value))
}

const inventoryList = ref(readStorageList(STORAGE_INVENTORY_KEY, defaultInventoryList))
const outboundList = ref(readStorageList(STORAGE_OUTBOUND_KEY, []))

watch(
  inventoryList,
  value => {
    saveStorageList(STORAGE_INVENTORY_KEY, value)
  },
  { deep: true }
)

const completedInventoryList = computed(() => {
  return inventoryList.value.filter(item => {
    return item.inventoryStatus === 'ready'
  })
})

const filteredInventoryList = computed(() => {
  return completedInventoryList.value.filter(item => {
    const keyword = filters.keyword.trim()

    const keywordMatch =
      !keyword ||
      item.sn.includes(keyword) ||
      item.macAddress.includes(keyword) ||
      item.softwareVersion.includes(keyword) ||
      item.hardwareVersion.includes(keyword) ||
      item.deviceType.includes(keyword)

    const deviceTypeMatch =
      !filters.deviceType || item.deviceType === filters.deviceType

    return keywordMatch && deviceTypeMatch
  })
})

const totalCount = computed(() => completedInventoryList.value.length)

const currentMonthOutboundCount = computed(() => {
  const currentMonth = new Date().toISOString().slice(0, 7)

  return outboundList.value.filter(item => {
    return item.outboundTime && item.outboundTime.startsWith(currentMonth)
  }).length
})

const deviceTypeSummary = computed(() => {
  return deviceTypeOptions.map(type => {
    return {
      deviceType: type,
      count: completedInventoryList.value.filter(item => item.deviceType === type).length
    }
  })
})

const totalPage = computed(() => {
  return Math.max(1, Math.ceil(filteredInventoryList.value.length / pageSize.value))
})

const paginatedInventoryList = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value

  return filteredInventoryList.value.slice(start, end)
})

const pageStartIndex = computed(() => {
  if (filteredInventoryList.value.length === 0) {
    return 0
  }

  return (currentPage.value - 1) * pageSize.value + 1
})

const pageEndIndex = computed(() => {
  return Math.min(currentPage.value * pageSize.value, filteredInventoryList.value.length)
})

watch(
  () => [filters.keyword, filters.deviceType],
  () => {
    currentPage.value = 1
  }
)

watch(pageSize, () => {
  currentPage.value = 1
})

watch(totalPage, value => {
  if (currentPage.value > value) {
    currentPage.value = value
  }
})

function getInventoryStatusText(status) {
  const map = {
    ready: '已完成烧录和出厂测试',
    waiting_burn: '待烧录',
    repair: '返修库存',
    testing: '出厂测试中'
  }

  return map[status] || status
}

function resetFilters() {
  filters.keyword = ''
  filters.deviceType = ''
  currentPage.value = 1
}

function viewInventory(item) {
  selectedInventory.value = item
}

function goFirstPage() {
  currentPage.value = 1
}

function goPrevPage() {
  if (currentPage.value > 1) {
    currentPage.value -= 1
  }
}

function goNextPage() {
  if (currentPage.value < totalPage.value) {
    currentPage.value += 1
  }
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