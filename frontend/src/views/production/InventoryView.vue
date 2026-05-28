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
        <p>当前仍在库存中的设备</p>
      </div>

      <div class="summary-card green">
        <span>当月出库数量</span>
        <strong>{{ currentMonthOutboundCount }}</strong>
        <p>本月已完成发货出库的设备数量</p>
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
        <span>按终端类型统计当前库存数量</span>
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
          <span>已勾选 {{ selectedIds.length }} 条</span>
        </div>

        <button
          class="ship-action-btn"
          :disabled="selectedIds.length === 0"
          @click="openActionDialog"
        >
          处理已选设备
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
            </tr>
          </thead>

          <tbody>
            <tr v-for="item in filteredInventoryList" :key="item.id">
              <td class="check-col">
                <input
                  type="checkbox"
                  :value="item.id"
                  v-model="selectedIds"
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
            </tr>
          </tbody>
        </table>
      </div>

      <div class="table-footer">
        共 {{ filteredInventoryList.length }} 条库存设备记录
      </div>
    </div>

    <!-- 出库列表 -->
    <div class="table-card outbound-card">
      <div class="table-card-header">
        <h3>出库列表</h3>
        <span>已出库 {{ outboundList.length }} 条</span>
      </div>

      <div class="table-wrapper">
        <table class="version-table outbound-table">
          <thead>
            <tr>
              <th>发货批次</th>
              <th>终端类型</th>
              <th>SN序列号</th>
              <th>MAC地址</th>
              <th>软件版本</th>
              <th>硬件版本</th>
              <th>入库时间</th>
              <th>出库时间</th>
              <th>操作人</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="item in outboundList" :key="item.outboundId">
              <td>
                <span class="batch-tag" :title="item.shippingBatchName">
                  {{ item.shippingBatchNo }}
                </span>
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
              <td>{{ item.operator }}</td>
              <td class="operation-col">
                <button class="ship-action-btn" @click="openReturnDialog(item)">
                  返回入库
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="table-footer">
        出库后的设备会从当前库存列表移除，并进入出库列表
      </div>
    </div>

    <!-- 确认发货弹窗 -->
   <div v-if="showActionDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>已选设备操作</h3>
          <button @click="showActionDialog = false">×</button>
        </div>

        <div class="ship-warning">
          <strong>当前已选择 {{ pendingShipList.length }} 台设备</strong>
          <p>可以在这里确认发货出库。确认后，设备会从当前库存列表移除，并自动进入出库列表。</p>
        </div>
        <div class="ship-batch-select">
          <label>
            关联发货批次
            <select v-model="selectedShippingBatchId">
              <option value="">请选择已有发货批次</option>
              <option
                v-for="batch in approvedShippingBatchList"
                :key="batch.id"
                :value="batch.id"
              >
                {{ batch.batchNo }} / {{ batch.batchName }}
              </option>
            </select>
          </label>

          <p>
            只能选择已经审核通过的发货批次。未审核通过的批次不能用于发货出库。
          </p>
        </div>

        <div class="ship-list">
          <div
            v-for="item in pendingShipList"
            :key="item.id"
            class="ship-item"
          >
            <span>{{ item.deviceType }}</span>
            <strong>{{ item.sn }}</strong>
            <em>{{ item.macAddress }}</em>
          </div>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showActionDialog = false">
            取消
          </button>

          <button class="ship-action-btn" @click="confirmShip">
            确认发货出库
          </button>
        </div>
      </div>
    </div>
    <!-- 返回入库确认弹窗 -->
    <div v-if="showReturnDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>确认返回入库</h3>
          <button @click="showReturnDialog = false">×</button>
        </div>

        <div class="ship-warning">
          <strong>确认将该设备返回当前库存？</strong>
          <p>确认后，该设备会从出库列表移除，并重新回到当前库存列表。</p>
        </div>

        <div v-if="pendingReturnItem" class="ship-list">
          <div class="ship-item">
            <span>{{ pendingReturnItem.deviceType }}</span>
            <strong>{{ pendingReturnItem.sn }}</strong>
            <em>{{ pendingReturnItem.macAddress }}</em>
          </div>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showReturnDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="confirmReturnInventory">
            确认返回入库
          </button>
        </div>
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
          <button class="green-btn" @click="openSingleShipDialog(selectedInventory)">
            发货出库
          </button>
          <button class="primary-btn" @click="selectedInventory = null">
            关闭
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'

const filters = reactive({
  keyword: '',
  deviceType: '',
  inventoryStatus: ''
})

const selectedInventory = ref(null)
const selectedIds = ref([])
const showActionDialog = ref(false)
const pendingShipIds = ref([])
const showReturnDialog = ref(false)
const pendingReturnItem = ref(null)
const selectedShippingBatchId = ref('')
const deviceTypeOptions = [
  '司机室控制盒',
  '解码板',
  '司机提醒单元',
  '噪声检测',
  '编码板',
  '功放板'
]
const shippingBatchList = ref([
  {
    id: 1,
    batchNo: '0000013289',
    batchName: '5月第一批司机室控制盒发货',
    auditStatus: 'approved'
  },
  {
    id: 2,
    batchNo: '0000013290',
    batchName: '5月第二批解码板发货',
    auditStatus: 'approved'
  },
  {
    id: 3,
    batchNo: '0000013291',
    batchName: '司机提醒单元发货批次',
    auditStatus: 'submitted'
  }
])


const approvedShippingBatchList = computed(() => {
  return shippingBatchList.value.filter(item => item.auditStatus === 'approved')
})
const inventoryList = ref([
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
    updateTime: '2026-05-12'
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
    updateTime: '2026-05-15'
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
    updateTime: '2026-05-14'
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
])

const outboundList = ref([])

const filteredInventoryList = computed(() => {
  return inventoryList.value.filter(item => {
    const keywordMatch =
      !filters.keyword ||
      item.sn.includes(filters.keyword) ||
      item.macAddress.includes(filters.keyword) ||
      item.softwareVersion.includes(filters.keyword) ||
      item.hardwareVersion.includes(filters.keyword) ||
      item.deviceType.includes(filters.keyword)

    const deviceTypeMatch =
      !filters.deviceType || item.deviceType === filters.deviceType

    const inventoryStatusMatch =
      !filters.inventoryStatus || item.inventoryStatus === filters.inventoryStatus

    return keywordMatch && deviceTypeMatch && inventoryStatusMatch
  })
})

const totalCount = computed(() => inventoryList.value.length)


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
      count: inventoryList.value.filter(item => item.deviceType === type).length
    }
  })
})

const isAllSelected = computed(() => {
  return (
    filteredInventoryList.value.length > 0 &&
    filteredInventoryList.value.every(item => selectedIds.value.includes(item.id))
  )
})

const pendingShipList = computed(() => {
  return inventoryList.value.filter(item => pendingShipIds.value.includes(item.id))
})

function getInventoryStatusText(status) {
  const map = {
    ready: '可发货',
    waiting_burn: '待烧录',
    repair: '返修库存',
    testing: '出厂测试中'
  }

  return map[status] || status
}

function resetFilters() {
  filters.keyword = ''
  filters.deviceType = ''
  filters.inventoryStatus = ''
}

function viewInventory(item) {
  selectedInventory.value = item
}

function toggleSelectAll(event) {
  if (event.target.checked) {
    selectedIds.value = filteredInventoryList.value.map(item => item.id)
  } else {
    selectedIds.value = []
  }
}

function openActionDialog() {
  pendingShipIds.value = [...selectedIds.value]
  selectedShippingBatchId.value = ''
  showActionDialog.value = true
}

function openSingleShipDialog(item) {
  pendingShipIds.value = [item.id]
  selectedShippingBatchId.value = ''
  showActionDialog.value = true
}

function confirmShip() {
  if (!selectedShippingBatchId.value) {
    alert('请选择发货批次')
    return
  }

  const selectedBatch = shippingBatchList.value.find(
    item => item.id === Number(selectedShippingBatchId.value)
  )

  if (!selectedBatch) {
    alert('发货批次不存在')
    return
  }

  if (selectedBatch.auditStatus !== 'approved') {
    alert('只能选择审核通过的发货批次')
    return
  }

  const today = new Date().toISOString().slice(0, 10)

  const shipItems = inventoryList.value.filter(item =>
    pendingShipIds.value.includes(item.id)
  )

  shipItems.forEach(item => {
    outboundList.value.unshift({
      ...item,
      outboundId: Date.now() + Math.random(),
      shippingBatchId: selectedBatch.id,
      shippingBatchNo: selectedBatch.batchNo,
      shippingBatchName: selectedBatch.batchName,
      outboundTime: today,
      operator: '当前用户'
    })
  })

  inventoryList.value = inventoryList.value.filter(
    item => !pendingShipIds.value.includes(item.id)
  )

  selectedIds.value = selectedIds.value.filter(
    id => !pendingShipIds.value.includes(id)
  )

  pendingShipIds.value = []
  selectedShippingBatchId.value = ''
  showActionDialog.value = false
  selectedInventory.value = null

  alert(`已确认发货，设备已关联发货批次【${selectedBatch.batchNo}】并自动出库`)
}
function openReturnDialog(item) {
  pendingReturnItem.value = item
  showReturnDialog.value = true
}
function confirmReturnInventory() {
  if (!pendingReturnItem.value) return

  const item = pendingReturnItem.value

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
    remark: '该设备由出库列表返回库存。'
  })

  outboundList.value = outboundList.value.filter(
    record => record.outboundId !== item.outboundId
  )

  pendingReturnItem.value = null
  showReturnDialog.value = false

  alert('设备已返回入库')
}
function markReady(item) {
  item.inventoryStatus = 'ready'
  item.updateTime = new Date().toISOString().slice(0, 10)

  if (selectedInventory.value) {
    selectedInventory.value = item
  }
}


</script>
<style scoped>
.page {
  width: 100%;
  min-height: 100%;
  color: #f8fafc;
}
.ship-batch-select {
  margin: 0 20px 20px;
  padding: 14px;
  border: 1px solid #1e293b;
  border-radius: 10px;
  background: #020617;
}

.ship-batch-select label {
  display: flex;
  flex-direction: column;
  gap: 8px;
  color: #cbd5e1;
  font-size: 13px;
}

.ship-batch-select select {
  height: 36px;
  border: 1px solid #334155;
  border-radius: 8px;
  background: #020617;
  color: #e2e8f0;
  padding: 0 12px;
  outline: none;
}

.ship-batch-select p {
  margin: 8px 0 0;
  color: #64748b;
  font-size: 12px;
}

.batch-tag {
  display: inline-block;
  max-width: 150px;
  padding: 4px 9px;
  border-radius: 999px;
  background: #1d4ed833;
  color: #60a5fa;
  font-size: 12px;
  font-weight: 700;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;
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

.page-header p {
  margin: 8px 0 0;
  color: #94a3b8;
  font-size: 14px;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 20px;
}

.summary-grid.two-col {
  grid-template-columns: repeat(2, 1fr);
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

.summary-card.yellow strong {
  color: #fbbf24;
}

.summary-card.red strong {
  color: #f87171;
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
.filter-card select,
.form-grid input,
.form-grid select,
.form-grid textarea {
  border: 1px solid #334155;
  border-radius: 8px;
  background: #020617;
  color: #e2e8f0;
  padding: 0 12px;
  outline: none;
}

.filter-card input,
.filter-card select,
.form-grid input,
.form-grid select {
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
  min-width: 1050px;
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

.file-name {
  margin-top: 4px;
  color: #64748b;
  font-size: 12px;
  line-height: 1.5;
  word-break: break-all;
}

.project-tag,
.device-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
}

.project-tag {
  background: #1d4ed833;
  color: #60a5fa;
}

.device-tag {
  background: #0f766e33;
  color: #5eead4;
}

.version-cell {
  overflow: hidden;
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

.sn-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  max-width: 180px;
}

.sn-tag {
  display: inline-flex;
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
  max-width: 130px;
  color: #cbd5e1;
  font-size: 12px;
  word-break: break-all;
  font-family: Consolas, Monaco, monospace;
}

.status-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
}

.status-tag.ready {
  background: #16a34a33;
  color: #4ade80;
}


.status-tag.testing {
  background: #1d4ed833;
  color: #60a5fa;
}

.status-tag.shipped {
  background: #47556933;
  color: #94a3b8;
}

.muted {
  color: #94a3b8 !important;
}


.text-btn.yellow {
  color: #fbbf24;
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

@media (max-width: 1200px) {
  .summary-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .type-grid {
    grid-template-columns: repeat(3, 1fr);
  }
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

@media (max-width: 960px) {
  .filter-card {
    grid-template-columns: 1fr;
  }

  .summary-grid,
  .type-grid {
    grid-template-columns: 1fr;
  }

  .version-table {
    min-width: 1050px;
  }

  .detail-card {
    grid-template-columns: 1fr;
  }
}
</style>