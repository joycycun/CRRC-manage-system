<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1>出库记录管理</h1>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="summary-grid two-col">
      <div class="summary-card green">
        <span>本月出库数量</span>
        <strong>{{ currentMonthOutboundCount }}</strong>
        <p>本月已进入出库记录的设备数量</p>
      </div>

      <div class="summary-card blue">
        <span>发货数量</span>
        <strong>{{ shippingCount }}</strong>
        <p>已从出库记录生成发货批次的设备数量</p>
      </div>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索 SN / MAC / 软件版本 / 硬件版本 / 操作人"
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

    <!-- 出库表格 -->
    <div class="table-card">
      <div class="table-card-header">
        <div>
          <h3>出库列表</h3>
          <span>已勾选 {{ selectedOutboundIds.length }} 条，共 {{ filteredOutboundList.length }} 条</span>
        </div>

        <button
          class="ship-action-btn"
          :disabled="selectedOutboundIds.length === 0"
          @click="goToShippingBatch"
        >
          生成发货批次
        </button>
      </div>

      <div class="table-wrapper">
        <table class="version-table">
          <thead>
            <tr>
              <th class="check-col">
                <input
                  type="checkbox"
                  :checked="isAllSelected"
                  @change="toggleSelectAll"
                />
              </th>
              <th>终端类型</th>
              <th>SN序列号</th>
              <th>MAC地址</th>
              <th>软件版本</th>
              <th>硬件版本</th>
              <th>入库时间</th>
              <th>出库时间</th>
              <th>操作人</th>
              <th>发货状态</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="item in filteredOutboundList" :key="item.outboundId">
              <td class="check-col">
                <input
                  type="checkbox"
                  :value="item.outboundId"
                  v-model="selectedOutboundIds"
                />
              </td>

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
              <td class="muted">{{ item.outboundTime }}</td>

              <td>
                <span class="normal-text" :title="item.operator">
                  {{ item.operator }}
                </span>
              </td>

              <td>
                <span
                  class="ship-status-tag"
                  :class="item.outboundStatus === 'shipping' ? 'shipping' : 'pending'"
                >
                  {{ getOutboundStatusText(item.outboundStatus) }}
                </span>
              </td>

              <td class="operation-col">
                <div class="action-group">
                  <button class="text-btn" @click="viewOutbound(item)">
                    查看
                  </button>

                  <button class="text-btn green" @click="returnToInventory(item)">
                    返回入库
                  </button>

                  <button class="text-btn red" @click="deleteOutbound(item)">
                    删除
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="table-footer">
        勾选需要发货的出库设备，点击“生成发货批次”后会自动跳转到发货管理页面。
      </div>
    </div>

    <!-- 查看详情弹窗 -->
    <div v-if="selectedOutbound" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>出库设备详情</h3>
          <button @click="selectedOutbound = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>终端类型</span>
            <strong>{{ selectedOutbound.deviceType }}</strong>
          </div>

          <div>
            <span>SN序列号</span>
            <strong>{{ selectedOutbound.sn }}</strong>
          </div>

          <div>
            <span>MAC地址</span>
            <strong>{{ selectedOutbound.macAddress }}</strong>
          </div>

          <div>
            <span>软件版本</span>
            <strong>{{ selectedOutbound.softwareVersion }}</strong>
          </div>

          <div>
            <span>硬件版本</span>
            <strong>{{ selectedOutbound.hardwareVersion }}</strong>
          </div>

          <div>
            <span>入库时间</span>
            <strong>{{ selectedOutbound.inTime }}</strong>
          </div>

          <div>
            <span>出库时间</span>
            <strong>{{ selectedOutbound.outboundTime }}</strong>
          </div>

          <div>
            <span>操作人</span>
            <strong>{{ selectedOutbound.operator }}</strong>
          </div>

          <div>
            <span>发货状态</span>
            <strong>{{ getOutboundStatusText(selectedOutbound.outboundStatus) }}</strong>
          </div>
        </div>

        <div class="remark-card">
          <span>出库说明</span>
          <p>{{ selectedOutbound.remark || '暂无说明' }}</p>
        </div>

        <div class="dialog-footer">
          <button class="green-btn" @click="returnToInventory(selectedOutbound)">
            返回入库
          </button>

          <button class="primary-btn" @click="selectedOutbound = null">
            关闭
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const STORAGE_INVENTORY_KEY = 'inventoryList'
const STORAGE_OUTBOUND_KEY = 'outInventoryList'
const STORAGE_PENDING_SHIPPING_KEY = 'pendingShippingDevices'

const filters = reactive({
  keyword: '',
  deviceType: ''
})

const selectedOutbound = ref(null)
const selectedOutboundIds = ref([])

const deviceTypeOptions = [
  '司机室控制盒',
  '解码板',
  '司机提醒单元',
  '噪声检测',
  '编码板',
  '功放板'
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

const inventoryList = ref(readStorageList(STORAGE_INVENTORY_KEY, []))
const outboundList = ref(readStorageList(STORAGE_OUTBOUND_KEY, []))

watch(
  inventoryList,
  value => {
    saveStorageList(STORAGE_INVENTORY_KEY, value)
  },
  { deep: true }
)

watch(
  outboundList,
  value => {
    saveStorageList(STORAGE_OUTBOUND_KEY, value)
  },
  { deep: true }
)

const filteredOutboundList = computed(() => {
  return outboundList.value.filter(item => {
    const keywordMatch =
      !filters.keyword ||
      item.sn.includes(filters.keyword) ||
      item.macAddress.includes(filters.keyword) ||
      item.softwareVersion.includes(filters.keyword) ||
      item.hardwareVersion.includes(filters.keyword) ||
      item.deviceType.includes(filters.keyword) ||
      item.operator.includes(filters.keyword)

    const deviceTypeMatch =
      !filters.deviceType || item.deviceType === filters.deviceType

    return keywordMatch && deviceTypeMatch
  })
})

const currentMonthOutboundCount = computed(() => {
  const currentMonth = new Date().toISOString().slice(0, 7)

  return outboundList.value.filter(item => {
    return item.outboundTime && item.outboundTime.startsWith(currentMonth)
  }).length
})

const shippingCount = computed(() => {
  return outboundList.value.filter(item => item.outboundStatus === 'shipping').length
})

const isAllSelected = computed(() => {
  return (
    filteredOutboundList.value.length > 0 &&
    filteredOutboundList.value.every(item =>
      selectedOutboundIds.value.includes(item.outboundId)
    )
  )
})

function getOutboundStatusText(status) {
  const map = {
    pending: '待生成发货',
    shipping: '已生成发货'
  }

  return map[status] || '待生成发货'
}

function resetFilters() {
  filters.keyword = ''
  filters.deviceType = ''
}

function toggleSelectAll(event) {
  if (event.target.checked) {
    selectedOutboundIds.value = filteredOutboundList.value.map(item => item.outboundId)
  } else {
    selectedOutboundIds.value = []
  }
}

function goToShippingBatch() {
  if (selectedOutboundIds.value.length === 0) {
    alert('请先勾选需要生成发货批次的出库设备')
    return
  }

  const selectedDevices = outboundList.value.filter(item =>
    selectedOutboundIds.value.includes(item.outboundId)
  )

  localStorage.setItem(
    STORAGE_PENDING_SHIPPING_KEY,
    JSON.stringify(selectedDevices)
  )

  outboundList.value = outboundList.value.map(item => {
    if (selectedOutboundIds.value.includes(item.outboundId)) {
      return {
        ...item,
        outboundStatus: 'shipping'
      }
    }

    return item
  })

  selectedOutboundIds.value = []

  alert('已选择出库设备，即将跳转到发货批次管理页面')

  router.push('/shipping/batch')
}

function viewOutbound(item) {
  selectedOutbound.value = item
}

function returnToInventory(item) {
  const ok = confirm(`确认将设备【${item.sn}】返回库存吗？`)
  if (!ok) return

  inventoryList.value.unshift({
    id: Date.now(),
    deviceType: item.deviceType,
    sn: item.sn,
    macAddress: item.macAddress,
    softwareVersion: item.softwareVersion,
    hardwareVersion: item.hardwareVersion,
    inventoryStatus: 'ready',
    inTime: item.inTime,
    updateTime: new Date().toISOString().slice(0, 10),
    remark: '该设备由出库记录返回库存。'
  })

  outboundList.value = outboundList.value.filter(
    record => record.outboundId !== item.outboundId
  )

  selectedOutboundIds.value = selectedOutboundIds.value.filter(
    id => id !== item.outboundId
  )

  selectedOutbound.value = null

  alert('设备已返回库存')
}

function deleteOutbound(item) {
  const ok = confirm(`确认删除出库记录【${item.sn}】吗？删除后不会返回库存。`)
  if (!ok) return

  outboundList.value = outboundList.value.filter(
    record => record.outboundId !== item.outboundId
  )

  selectedOutboundIds.value = selectedOutboundIds.value.filter(
    id => id !== item.outboundId
  )

  selectedOutbound.value = null
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

.summary-card.blue strong {
  color: #60a5fa;
}

.check-col {
  width: 56px;
  text-align: center !important;
}

.ship-action-btn {
  height: 34px;
  padding: 0 16px;
  border: none;
  border-radius: 8px;
  background: #2563eb;
  color: #fff;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 0 0 1px rgba(37, 99, 235, 0.35);
}

.ship-action-btn:hover {
  background: #1d4ed8;
}

.ship-action-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.primary-btn,
.query-btn,
.reset-btn,
.green-btn,
.red-btn {
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

.green-btn {
  border: 1px solid #166534;
  background: #052e16;
  color: #86efac;
}

.red-btn {
  border: 1px solid #7f1d1d;
  background: #450a0a;
  color: #fca5a5;
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
.filter-card select {
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
  min-width: 1420px;
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
  width: 250px;
  text-align: right !important;
}

.action-group {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 12px;
  flex-wrap: nowrap;
  white-space: nowrap;
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

.mac-text,
.normal-text {
  display: inline-block;
  max-width: 160px;
  color: #cbd5e1;
  font-size: 12px;
  font-family: Consolas, Monaco, monospace;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.normal-text {
  font-family: inherit;
  font-size: 13px;
}

.version-cell {
  overflow: hidden;
}

.muted {
  color: #94a3b8 !important;
}

.ship-status-tag {
  display: inline-block;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
}

.ship-status-tag.pending {
  background: #47556933;
  color: #94a3b8;
}

.ship-status-tag.shipping {
  background: #1d4ed833;
  color: #60a5fa;
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

.text-btn.green {
  color: #4ade80;
}

.text-btn.red {
  color: #f87171;
}

.table-footer {
  padding: 12px 16px;
  color: #64748b;
  font-size: 12px;
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

@media (max-width: 960px) {
  .filter-card {
    grid-template-columns: 1fr;
  }

  .summary-grid {
    grid-template-columns: 1fr;
  }

  .version-table {
    min-width: 1420px;
  }

  .detail-card {
    grid-template-columns: 1fr;
  }
}
</style>