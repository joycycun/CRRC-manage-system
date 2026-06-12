<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1>发货批次管理</h1>
      </div>

      <button class="primary-btn" @click="openCreateDialog">
        新增发货批次
      </button>
    </div>

    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索批次号 / 快递单号 / SN / MAC / 上传人"
      />

      <select v-model="filters.auditStatus">
        <option value="">全部审核状态</option>
        <option value="draft">草稿</option>
        <option value="submitted">待审核</option>
        <option value="approved">审核通过</option>
        <option value="rejected">审核驳回</option>
      </select>

      <button class="query-btn" @click="loadShippingBatches">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <div class="table-card">
      <div class="table-wrapper">
        <table class="model-table">
          <thead>
            <tr>
              <th>发货批次号</th>
              <th>设备数量</th>
              <th>快递单号</th>
              <th>上传人</th>
              <th>上传时间</th>
              <th>审核状态</th>
              <th>审核人</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>

          <tbody>
            <tr
              v-for="item in filteredBatchList"
              :key="item.id"
              class="model-row"
            >
              <td>
                <button
                  class="model-name-btn"
                  @click="viewBatch(item)"
                >
                  {{ item.batchNo }}
                </button>
              </td>

              <td>
                <span class="count-tag">
                  {{ item.deviceCount }} 台
                </span>
              </td>

              <td>
                <span class="file-text" :title="item.expressNo">
                  {{ item.expressNo || '-' }}
                </span>
              </td>

              <td>
                <span class="normal-text" :title="item.uploader">
                  {{ item.uploader || '-' }}
                </span>
              </td>

              <td class="muted nowrap">
                {{ item.uploadTime || '-' }}
              </td>

              <td>
                <span class="status-tag" :class="item.auditStatus">
                  {{ getAuditStatusText(item.auditStatus) }}
                </span>
              </td>

              <td>
                {{ item.auditor || '-' }}
              </td>

              <td class="operation-col">
                <div class="action-group">
                  <button
                    class="text-btn blue"
                    @click="viewBatch(item)"
                  >
                    查看
                  </button>

                  <button
                    v-if="item.auditStatus === 'draft' || item.auditStatus === 'rejected'"
                    class="text-btn blue"
                    @click="submitBatch(item)"
                  >
                    提交
                  </button>

                  <button
                    v-if="item.auditStatus === 'submitted'"
                    class="text-btn green"
                    @click="auditBatch(item)"
                  >
                    审核
                  </button>

                  <button
                    v-if="item.auditStatus === 'draft' || item.auditStatus === 'rejected'"
                    class="text-btn red"
                    @click="deleteBatch(item)"
                  >
                    删除
                  </button>
                </div>
              </td>
            </tr>

            <tr v-if="filteredBatchList.length === 0">
              <td colspan="8" class="empty-table">
                暂无发货批次记录
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="table-footer">
        共 {{ filteredBatchList.length }} 条发货批次记录。审核通过后，批次内设备会自动进入出库记录。
      </div>
    </div>

    <div v-if="showCreateDialog" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>新增发货批次</h3>
          <button @click="showCreateDialog = false">×</button>
        </div>

        <div class="form-grid">
          <label>
            发货批次号
            <input v-model="batchForm.batchNo" placeholder="请输入发货批次号" />
          </label>

          <label>
            上传人
            <input v-model="batchForm.uploader" disabled />
          </label>

          <label>
            快递单号
            <input v-model="batchForm.expressNo" placeholder="请输入快递单号" />
          </label>

          <label>
            发货单文件
            <input type="file" accept=".xls,.xlsx,.csv,.pdf,.doc,.docx,.zip" @change="handleFileChange" />
          </label>

          <label class="full-row">
            发货说明
            <textarea v-model="batchForm.remark" placeholder="请输入发货说明"></textarea>
          </label>

          <label class="full-row">
            选择库存设备
            <div class="mac-select-panel">
              <input
                v-model="inventoryFilters.keyword"
                class="mac-search-input"
                placeholder="搜索 SN / MAC / 终端类型 / 软件版本 / 硬件版本"
              />

              <select v-model="inventoryFilters.deviceType" class="mac-search-input">
                <option value="">全部终端类型</option>
                <option v-for="type in deviceTypeOptions" :key="type" :value="type">
                  {{ type }}
                </option>
              </select>

              <div class="mac-panel-header">
                <span>
                  当前筛选 {{ filteredInventoryList.length }} 台，
                  显示第 {{ inventoryPageStartIndex }} - {{ inventoryPageEndIndex }} 台
                </span>
                <strong>已选择 {{ selectedInventoryIds.length }} 台</strong>
              </div>

              <div class="action-group">
                <button class="text-btn blue" @click="selectAllFilteredInventory">
                  全选当前筛选
                </button>
                <button class="text-btn red" @click="clearSelectedInventory">
                  清空选择
                </button>
              </div>

              <table class="mac-dialog-table">
                <thead>
                  <tr>
                    <th>
                      <input
                        type="checkbox"
                        :checked="isAllFilteredInventorySelected"
                        @change="toggleAllFilteredInventory"
                      />
                    </th>
                    <th>终端类型</th>
                    <th>SN序列号</th>
                    <th>MAC地址</th>
                    <th>软件版本</th>
                    <th>硬件版本</th>
                  </tr>
                </thead>

                <tbody>
                  <tr v-for="item in paginatedInventoryList" :key="item.id">
                    <td>
                      <input
                        v-model="selectedInventoryIds"
                        type="checkbox"
                        :value="item.id"
                      />
                    </td>
                    <td>{{ item.deviceType || '-' }}</td>
                    <td>{{ item.sn || '-' }}</td>
                    <td>{{ item.macAddress || '-' }}</td>
                    <td>{{ item.softwareVersion || '-' }}</td>
                    <td>{{ item.hardwareVersion || '-' }}</td>
                  </tr>

                  <tr v-if="paginatedInventoryList.length === 0">
                    <td colspan="6" class="empty-table">
                      暂无可选择的在库设备
                    </td>
                  </tr>
                </tbody>
              </table>

              <div class="pagination-bar">
                <div class="pagination-info">
                  第 {{ inventoryCurrentPage }} / {{ inventoryTotalPage }} 页
                </div>
                <div class="pagination-actions">
                  <button
                    class="reset-btn"
                    :disabled="inventoryCurrentPage === 1"
                    @click="goPrevInventoryPage"
                  >
                    上一页
                  </button>
                  <button
                    class="reset-btn"
                    :disabled="inventoryCurrentPage === inventoryTotalPage"
                    @click="goNextInventoryPage"
                  >
                    下一页
                  </button>
                </div>
              </div>
            </div>
          </label>

          <label class="full-row">
            已选设备
            <div class="mac-select-panel">
              <div v-if="selectedInventoryDevices.length === 0" class="empty-dialog-data">
                暂未选择库存设备
              </div>
              <div
                v-for="device in selectedInventoryPreview"
                v-else
                :key="device.id"
                class="mac-check-item"
              >
                <span>{{ device.deviceType || '-' }}</span>
                <em>{{ device.sn || '-' }}</em>
                <em>{{ device.macAddress || '-' }}</em>
                <em>{{ device.softwareVersion || '-' }}</em>
                <em>{{ device.hardwareVersion || '-' }}</em>
              </div>
              <div v-if="selectedInventoryDevices.length > selectedInventoryPreview.length" class="empty-dialog-data">
                还有 {{ selectedInventoryDevices.length - selectedInventoryPreview.length }} 台已选择设备未展开显示。
              </div>
            </div>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showCreateDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="createBatch">
            保存批次
          </button>
        </div>
      </div>
    </div>

    <div v-if="selectedBatch" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>发货批次详情</h3>
          <button @click="selectedBatch = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>发货批次号</span>
            <strong>{{ selectedBatch.batchNo || '-' }}</strong>
          </div>

          <div>
            <span>快递单号</span>
            <strong>{{ selectedBatch.expressNo || '-' }}</strong>
          </div>

          <div>
            <span>设备数量</span>
            <strong>{{ selectedBatch.deviceCount }} 台</strong>
          </div>

          <div>
            <span>上传人</span>
            <strong>{{ selectedBatch.uploader || '-' }}</strong>
          </div>

          <div>
            <span>上传时间</span>
            <strong>{{ selectedBatch.uploadTime || '-' }}</strong>
          </div>

          <div>
            <span>审核状态</span>
            <strong>{{ getAuditStatusText(selectedBatch.auditStatus) }}</strong>
          </div>

          <div>
            <span>审核人</span>
            <strong>{{ selectedBatch.auditor || '-' }}</strong>
          </div>

          <div>
            <span>审核时间</span>
            <strong>{{ selectedBatch.auditTime || '-' }}</strong>
          </div>
        </div>

        <div class="remark-card">
          <span>发货说明</span>
          <p>{{ selectedBatch.remark || selectedBatch.shippingDesc || '暂无说明' }}</p>
        </div>

        <div class="mac-dialog-table-wrapper">
          <input
            v-model="batchDeviceKeyword"
            class="mac-search-input"
            placeholder="搜索本批次设备 SN / MAC / 终端类型 / 软件版本 / 硬件版本"
          />

          <div class="mac-panel-header">
            <span>
              当前筛选 {{ selectedBatchFilteredDevices.length }} 台，
              显示第 {{ batchDevicePageStartIndex }} - {{ batchDevicePageEndIndex }} 台
            </span>
            <strong>共 {{ selectedBatch.deviceList.length }} 台</strong>
          </div>

          <table class="mac-dialog-table">
            <thead>
              <tr>
                <th>序号</th>
                <th>终端类型</th>
                <th>SN序列号</th>
                <th>MAC地址</th>
                <th>软件版本</th>
                <th>硬件版本</th>
              </tr>
            </thead>

            <tbody>
              <tr v-for="(device, index) in paginatedSelectedBatchDevices" :key="device.detailId || device.id">
                <td>{{ batchDevicePageStartIndex + index }}</td>
                <td>{{ device.deviceType || '-' }}</td>
                <td>{{ device.sn || '-' }}</td>
                <td>{{ device.macAddress || '-' }}</td>
                <td>{{ device.softwareVersion || '-' }}</td>
                <td>{{ device.hardwareVersion || '-' }}</td>
              </tr>

              <tr v-if="paginatedSelectedBatchDevices.length === 0">
                <td colspan="6" class="empty-table">
                  暂无设备明细
                </td>
              </tr>
            </tbody>
          </table>

          <div class="pagination-bar">
            <div class="pagination-info">
              第 {{ batchDeviceCurrentPage }} / {{ batchDeviceTotalPage }} 页
            </div>
            <div class="pagination-actions">
              <button
                class="reset-btn"
                :disabled="batchDeviceCurrentPage === 1"
                @click="goPrevBatchDevicePage"
              >
                上一页
              </button>
              <button
                class="reset-btn"
                :disabled="batchDeviceCurrentPage === batchDeviceTotalPage"
                @click="goNextBatchDevicePage"
              >
                下一页
              </button>
            </div>
          </div>
        </div>

        <div class="dialog-footer">
          <button
            v-if="selectedBatch.auditStatus === 'submitted'"
            class="green-btn"
            @click="approveBatch(selectedBatch)"
          >
            审核通过
          </button>

          <button
            v-if="selectedBatch.auditStatus === 'submitted'"
            class="red-btn"
            @click="rejectBatch(selectedBatch)"
          >
            审核驳回
          </button>

          <button class="reset-btn" @click="openFile(selectedBatch)">
            打开文件
          </button>

          <button class="reset-btn" @click="downloadFile(selectedBatch)">
            下载文件
          </button>

          <button class="primary-btn" @click="selectedBatch = null">
            关闭
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { getInventory } from '@/api/inventory'
import {
  getShippingBatches,
  createShippingBatch,
  submitShippingBatch,
  auditShippingBatch,
  deleteShippingBatch
} from '@/api/shippingBatch'

const filters = reactive({ keyword: '', auditStatus: '' })
const inventoryFilters = reactive({ keyword: '', deviceType: '' })
const macSelectForm = reactive({ macText: '', startMac: '', endMac: '' })

const showCreateDialog = ref(false)
const selectedBatch = ref(null)
const selectedInventoryIds = ref([])
const inventoryList = ref([])
const batchList = ref([])
const inventoryCurrentPage = ref(1)
const inventoryPageSize = 20
const batchDeviceKeyword = ref('')
const batchDeviceCurrentPage = ref(1)
const batchDevicePageSize = 30

const batchForm = reactive({
  batchNo: '',
  uploader: '',
  expressNo: '',
  fileName: '',
  file: null,
  fileUrl: '',
  remark: ''
})

onMounted(async () => {
  await loadInventory()
  await loadShippingBatches()
})

function getResponseData(res) {
  return res && res.data ? res.data : res
}

function formatDateTime(value) {
  if (!value) return ''
  if (typeof value === 'string') return value.slice(0, 19).replace('T', ' ')
  if (value.Time) return String(value.Time).slice(0, 19).replace('T', ' ')
  return String(value)
}

function normalizeInventoryStatus(status) {
  const map = {
    ready: '在库',
    locked: '锁定',
    outbound: '已出库',
    repair: '返修',
    在库: '在库',
    锁定: '锁定',
    已出库: '已出库',
    返修: '返修'
  }
  return map[status] || status
}

function normalizeInventory(item) {
  return {
    id: item.id,
    deviceType: item.deviceType || item.device_type || '',
    productName: item.productName || item.product_name || '',
    productModel: item.productModel || item.product_model || '',
    sn: item.sn || '',
    macAddress: item.macAddress || item.mac_address || '',
    hardwareVersion: item.hardwareVersion || item.hardware_version || '',
    softwareVersion: item.softwareVersion || item.software_version || '',
    inventoryStatus: normalizeInventoryStatus(item.inventoryStatus || item.inventory_status || ''),
    inTime: formatDateTime(item.inTime || item.in_time),
    updateTime: formatDateTime(item.updateTime || item.update_time),
    remark: item.remark || ''
  }
}

function normalizeAuditStatus(status) {
  const map = {
    草稿: 'draft',
    待审核: 'submitted',
    已通过: 'approved',
    已驳回: 'rejected',
    draft: 'draft',
    submitted: 'submitted',
    approved: 'approved',
    rejected: 'rejected'
  }
  return map[status] || status || 'draft'
}

function normalizeBatch(item) {
  const rawDeviceList = Array.isArray(item.deviceList)
    ? item.deviceList
    : Array.isArray(item.device_list)
      ? item.device_list
      : []

  const deviceList = rawDeviceList.map(device => ({
    id: device.id || device.inventoryDeviceId || device.inventory_device_id,
    detailId: device.detailId || device.detail_id || device.id,
    inventoryDeviceId: device.inventoryDeviceId || device.inventory_device_id || device.id,
    deviceType: device.deviceType || device.device_type || '',
    sn: device.sn || '',
    macAddress: device.macAddress || device.mac_address || '',
    hardwareVersion: device.hardwareVersion || device.hardware_version || '',
    softwareVersion: device.softwareVersion || device.software_version || ''
  }))

  const snList = Array.isArray(item.snList)
    ? item.snList
    : Array.isArray(item.sn_list)
      ? item.sn_list
      : deviceList.map(device => device.sn).filter(Boolean)

  const macList = Array.isArray(item.macList)
    ? item.macList
    : Array.isArray(item.mac_list)
      ? item.mac_list
      : deviceList.map(device => device.macAddress).filter(Boolean)

  return {
    id: item.id,
    batchNo: item.batchNo || item.batch_no || '',
    projectId: item.projectId || item.project_id || 0,
    expressNo: item.expressNo || item.express_no || '',
    deviceCount: item.deviceCount || item.device_count || deviceList.length || 0,
    fileId: item.fileId || item.file_id || 0,
    fileName: item.fileName || item.file_name || '',
    fileUrl: item.fileUrl || item.file_url || '',
    uploaderId: item.uploaderId || item.uploader_id || 0,
    uploader: item.uploader || item.uploaderName || item.uploader_name || '',
    uploaderName: item.uploaderName || item.uploader_name || item.uploader || '',
    uploadTime: formatDateTime(item.uploadTime || item.upload_time),
    auditStatus: normalizeAuditStatus(item.auditStatus || item.audit_status || 'draft'),
    auditorId: item.auditorId || item.auditor_id || 0,
    auditor: item.auditor || item.auditorName || item.auditor_name || '',
    auditorName: item.auditorName || item.auditor_name || item.auditor || '',
    auditTime: formatDateTime(item.auditTime || item.audit_time),
    remark: item.remark || '',
    rejectReason: item.rejectReason || item.reject_reason || '',
    shippingDesc: item.shippingDesc || item.shipping_desc || '',
    deviceList,
    snList,
    macList
  }
}

async function loadInventory() {
  try {
    const res = await getInventory()
    const result = getResponseData(res)
    console.log('库存接口返回：', result)

    if (result.code !== 200) {
      alert(result.msg || '加载库存失败')
      return
    }

    inventoryList.value = (result.data || []).map(normalizeInventory)
  } catch (err) {
    console.error('加载库存失败：', err)
    alert(err.response?.data || '加载库存失败，请检查 /api/inventory 接口')
  }
}

async function loadShippingBatches() {
  try {
    const res = await getShippingBatches()
    const result = getResponseData(res)
    console.log('发货批次返回：', result)

    if (result.code !== 200) {
      alert(result.msg || '加载发货批次失败')
      return
    }

    batchList.value = (result.data || []).map(normalizeBatch)
  } catch (err) {
    console.error('加载发货批次失败：', err)
    alert(err.response?.data || '加载发货批次失败，请检查 /api/shipping-batches 接口')
  }
}

const availableInventoryList = computed(() => {
  return inventoryList.value.filter(item => item.inventoryStatus === '在库')
})

const deviceTypeOptions = computed(() => {
  return [...new Set(availableInventoryList.value.map(item => item.deviceType).filter(Boolean))]
})

const filteredInventoryList = computed(() => {
  return availableInventoryList.value.filter(item => {
    const keyword = inventoryFilters.keyword.trim().toLowerCase()
    const keywordMatch =
      !keyword ||
      String(item.sn || '').toLowerCase().includes(keyword) ||
      String(item.macAddress || '').toLowerCase().includes(keyword) ||
      String(item.deviceType || '').toLowerCase().includes(keyword) ||
      String(item.softwareVersion || '').toLowerCase().includes(keyword) ||
      String(item.hardwareVersion || '').toLowerCase().includes(keyword)

    const deviceTypeMatch = !inventoryFilters.deviceType || item.deviceType === inventoryFilters.deviceType
    return keywordMatch && deviceTypeMatch
  })
})

const inventoryTotalPage = computed(() => Math.max(1, Math.ceil(filteredInventoryList.value.length / inventoryPageSize)))
const paginatedInventoryList = computed(() => {
  if (inventoryCurrentPage.value > inventoryTotalPage.value) inventoryCurrentPage.value = inventoryTotalPage.value
  const start = (inventoryCurrentPage.value - 1) * inventoryPageSize
  return filteredInventoryList.value.slice(start, start + inventoryPageSize)
})
const inventoryPageStartIndex = computed(() => filteredInventoryList.value.length === 0 ? 0 : (inventoryCurrentPage.value - 1) * inventoryPageSize + 1)
const inventoryPageEndIndex = computed(() => Math.min(inventoryCurrentPage.value * inventoryPageSize, filteredInventoryList.value.length))

const selectedInventoryDevices = computed(() => {
  return availableInventoryList.value.filter(item => selectedInventoryIds.value.includes(item.id))
})

const selectedInventoryPreview = computed(() => selectedInventoryDevices.value.slice(0, 30))

const isAllFilteredInventorySelected = computed(() => {
  return filteredInventoryList.value.length > 0 &&
    filteredInventoryList.value.every(item => selectedInventoryIds.value.includes(item.id))
})

const filteredBatchList = computed(() => {
  return batchList.value.filter(item => {
    const keyword = filters.keyword.trim().toLowerCase()
    const snMatch = item.snList && item.snList.some(sn => String(sn || '').toLowerCase().includes(keyword))
    const macMatch = item.deviceList && item.deviceList.some(device => String(device.macAddress || '').toLowerCase().includes(keyword))
    const keywordMatch =
      !keyword ||
      String(item.batchNo || '').toLowerCase().includes(keyword) ||
      String(item.uploader || '').toLowerCase().includes(keyword) ||
      String(item.fileName || '').toLowerCase().includes(keyword) ||
      String(item.expressNo || '').toLowerCase().includes(keyword) ||
      snMatch ||
      macMatch

    const auditStatusMatch = !filters.auditStatus || item.auditStatus === filters.auditStatus
    return keywordMatch && auditStatusMatch
  })
})

const selectedBatchFilteredDevices = computed(() => {
  if (!selectedBatch.value) return []
  const keyword = batchDeviceKeyword.value.trim().toLowerCase()
  return selectedBatch.value.deviceList.filter(item => {
    if (!keyword) return true
    return String(item.sn || '').toLowerCase().includes(keyword) ||
      String(item.macAddress || '').toLowerCase().includes(keyword) ||
      String(item.deviceType || '').toLowerCase().includes(keyword) ||
      String(item.softwareVersion || '').toLowerCase().includes(keyword) ||
      String(item.hardwareVersion || '').toLowerCase().includes(keyword)
  })
})

const batchDeviceTotalPage = computed(() => Math.max(1, Math.ceil(selectedBatchFilteredDevices.value.length / batchDevicePageSize)))
const paginatedSelectedBatchDevices = computed(() => {
  if (batchDeviceCurrentPage.value > batchDeviceTotalPage.value) batchDeviceCurrentPage.value = batchDeviceTotalPage.value
  const start = (batchDeviceCurrentPage.value - 1) * batchDevicePageSize
  return selectedBatchFilteredDevices.value.slice(start, start + batchDevicePageSize)
})
const batchDevicePageStartIndex = computed(() => selectedBatchFilteredDevices.value.length === 0 ? 0 : (batchDeviceCurrentPage.value - 1) * batchDevicePageSize + 1)
const batchDevicePageEndIndex = computed(() => Math.min(batchDeviceCurrentPage.value * batchDevicePageSize, selectedBatchFilteredDevices.value.length))

function getCurrentUserName() {
  return localStorage.getItem('username') || localStorage.getItem('accountName') || localStorage.getItem('realName') || '当前用户'
}

function getAuditStatusText(status) {
  const map = { draft: '草稿', submitted: '待审核', approved: '审核通过', rejected: '审核驳回' }
  return map[status] || status
}

function resetFilters() {
  filters.keyword = ''
  filters.auditStatus = ''
}

async function openCreateDialog() {
  batchForm.batchNo = ''
  batchForm.uploader = getCurrentUserName()
  batchForm.expressNo = ''
  batchForm.fileName = ''
  batchForm.file = null
  batchForm.fileUrl = ''
  batchForm.remark = ''

  inventoryFilters.keyword = ''
  inventoryFilters.deviceType = ''
  inventoryCurrentPage.value = 1
  macSelectForm.macText = ''
  macSelectForm.startMac = ''
  macSelectForm.endMac = ''
  selectedInventoryIds.value = []

  await loadInventory()
  showCreateDialog.value = true
}

function handleFileChange(event) {
  const file = event.target.files[0]
  if (!file) return
  batchForm.file = file
  batchForm.fileName = file.name
  batchForm.fileUrl = URL.createObjectURL(file)
}

function normalizeMac(mac) {
  return String(mac || '').trim().replace(/-/g, ':').toUpperCase()
}

function parseMacText(text) {
  return String(text || '')
    .split(/[\n\r,，;；\s]+/)
    .map(item => normalizeMac(item))
    .filter(Boolean)
}

function macToNumber(mac) {
  const normalized = normalizeMac(mac)
  const hex = normalized.replace(/:/g, '')
  if (!/^[0-9A-F]{12}$/.test(hex)) return null
  return parseInt(hex, 16)
}

function selectByMacList() {
  const macList = parseMacText(macSelectForm.macText)
  if (macList.length === 0) {
    alert('请先填写 MAC 地址')
    return
  }

  const matchedIds = availableInventoryList.value
    .filter(item => macList.includes(normalizeMac(item.macAddress)))
    .map(item => item.id)

  if (matchedIds.length === 0) {
    alert('库存中没有匹配到这些 MAC 地址')
    return
  }

  selectedInventoryIds.value = [...new Set([...selectedInventoryIds.value, ...matchedIds])]
  alert(`已按 MAC 列表选择 ${matchedIds.length} 台库存设备`)
}

function selectByMacRange() {
  const start = macToNumber(macSelectForm.startMac)
  const end = macToNumber(macSelectForm.endMac)
  if (start === null || end === null) {
    alert('请填写正确的起始 MAC 和结束 MAC')
    return
  }

  const min = Math.min(start, end)
  const max = Math.max(start, end)
  const matchedIds = availableInventoryList.value
    .filter(item => {
      const value = macToNumber(item.macAddress)
      return value !== null && value >= min && value <= max
    })
    .map(item => item.id)

  if (matchedIds.length === 0) {
    alert('库存中没有匹配到该 MAC 区间内的设备')
    return
  }

  selectedInventoryIds.value = [...new Set([...selectedInventoryIds.value, ...matchedIds])]
  alert(`已按 MAC 区间选择 ${matchedIds.length} 台库存设备`)
}

function clearMacText() {
  macSelectForm.macText = ''
}

function selectAllFilteredInventory() {
  const ids = filteredInventoryList.value.map(item => item.id)
  selectedInventoryIds.value = [...new Set([...selectedInventoryIds.value, ...ids])]
}

function toggleAllFilteredInventory(event) {
  const ids = filteredInventoryList.value.map(item => item.id)
  if (event.target.checked) {
    selectedInventoryIds.value = [...new Set([...selectedInventoryIds.value, ...ids])]
  } else {
    selectedInventoryIds.value = selectedInventoryIds.value.filter(id => !ids.includes(id))
  }
}

function clearSelectedInventory() {
  selectedInventoryIds.value = []
}

function goPrevInventoryPage() {
  if (inventoryCurrentPage.value > 1) inventoryCurrentPage.value -= 1
}

function goNextInventoryPage() {
  if (inventoryCurrentPage.value < inventoryTotalPage.value) inventoryCurrentPage.value += 1
}

async function createBatch() {
  if (!batchForm.batchNo) {
    alert('请输入发货批次号')
    return
  }
  if (selectedInventoryIds.value.length === 0) {
    alert('请先选择需要发货的库存设备')
    return
  }
  if (!batchForm.uploader) {
    alert('上传人读取失败')
    return
  }
  if (!batchForm.file) {
    alert('请上传发货单 Excel 文件')
    return
  }

  try {
    const res = await createShippingBatch({
      batchNo: batchForm.batchNo,
      expressNo: batchForm.expressNo,
      uploaderId: 1,
      uploaderName: batchForm.uploader,
      remark: batchForm.remark,
      shippingDesc: batchForm.remark,
      inventoryDeviceIds: selectedInventoryIds.value
    })
    const result = getResponseData(res)
    if (result.code !== 200) {
      alert(result.msg || '新增发货批次失败')
      return
    }

    selectedInventoryIds.value = []
    showCreateDialog.value = false
    await loadInventory()
    await loadShippingBatches()
    alert('发货批次已生成，审核通过后进入出库管理')
  } catch (err) {
    console.error('新增发货批次失败：', err)
    alert(err.response?.data || '新增发货批次失败，请检查后端接口')
  }
}

function viewBatch(item) {
  selectedBatch.value = item
  batchDeviceKeyword.value = ''
  batchDeviceCurrentPage.value = 1
}

function auditBatch(item) {
  selectedBatch.value = item
  batchDeviceKeyword.value = ''
  batchDeviceCurrentPage.value = 1
}

function goPrevBatchDevicePage() {
  if (batchDeviceCurrentPage.value > 1) batchDeviceCurrentPage.value -= 1
}

function goNextBatchDevicePage() {
  if (batchDeviceCurrentPage.value < batchDeviceTotalPage.value) batchDeviceCurrentPage.value += 1
}

async function submitBatch(item) {
  try {
    const res = await submitShippingBatch(item.id)
    const result = getResponseData(res)
    if (result.code !== 200) {
      alert(result.msg || '提交失败')
      return
    }
    await loadShippingBatches()
    alert(`发货批次【${item.batchNo}】已提交领导审核`)
  } catch (err) {
    console.error('提交失败：', err)
    alert(err.response?.data || '提交失败')
  }
}

async function approveBatch(item) {
  const ok = confirm(`确认发货批次【${item.batchNo}】审核通过吗？`)
  if (!ok) return

  try {
    const res = await auditShippingBatch(item.id, { auditStatus: '已通过', auditorId: 1, auditorName: '领导' })
    const result = getResponseData(res)
    if (result.code !== 200) {
      alert(result.msg || '审核失败')
      return
    }
    selectedBatch.value = null
    await loadInventory()
    await loadShippingBatches()
    alert(`发货批次【${item.batchNo}】审核通过，设备已进入出库管理`)
  } catch (err) {
    console.error('审核失败：', err)
    alert(err.response?.data || '审核失败')
  }
}

async function rejectBatch(item) {
  try {
    const res = await auditShippingBatch(item.id, { auditStatus: '已驳回', auditorId: 1, auditorName: '领导', rejectReason: '审核驳回' })
    const result = getResponseData(res)
    if (result.code !== 200) {
      alert(result.msg || '审核驳回失败')
      return
    }
    selectedBatch.value = null
    await loadInventory()
    await loadShippingBatches()
    alert(`发货批次【${item.batchNo}】已驳回`)
  } catch (err) {
    console.error('审核驳回失败：', err)
    alert(err.response?.data || '审核驳回失败')
  }
}

function openFile(item) {
  if (!item.fileUrl) {
    alert('当前表结构只有 file_id，没有 file_url，暂不支持直接打开文件')
    return
  }
  window.open(item.fileUrl, '_blank')
}

function downloadFile(item) {
  if (!item.fileUrl) {
    alert('当前表结构只有 file_id，没有 file_url，暂不支持下载文件')
    return
  }
  const link = document.createElement('a')
  link.href = item.fileUrl
  link.download = item.fileName || '发货单.xlsx'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

async function deleteBatch(item) {
  const ok = confirm(`确认删除发货批次【${item.batchNo}】吗？`)
  if (!ok) return

  try {
    const res = await deleteShippingBatch(item.id)
    const result = getResponseData(res)
    if (result.code !== 200) {
      alert(result.msg || '删除失败')
      return
    }
    await loadInventory()
    await loadShippingBatches()
    alert('删除成功')
  } catch (err) {
    console.error('删除失败：', err)
    alert(err.response?.data || '删除失败')
  }
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

.reset-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
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
  grid-template-columns: 1.4fr 220px 180px 90px 90px;
  gap: 12px;
  margin-bottom: 20px;
}

.filter-card input,
.filter-card select,
.form-grid input,
.form-grid select,
.form-grid textarea,
.mac-search-input {
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
.form-grid select,
.mac-search-input {
  height: 36px;
}

.form-grid textarea {
  min-height: 90px;
  padding: 10px 12px;
  resize: vertical;
}

.filter-card input::placeholder,
.form-grid input::placeholder,
.form-grid textarea::placeholder,
.mac-search-input::placeholder {
  color: #64748b;
}

.form-grid input[type="file"] {
  height: auto;
  padding: 8px 12px;
  cursor: pointer;
}

.form-grid input[type="file"]::file-selector-button {
  height: 28px;
  padding: 0 12px;
  margin-right: 12px;
  border: 1px solid #334155;
  border-radius: 6px;
  background: #1e293b;
  color: #cbd5e1;
  cursor: pointer;
}

.table-card {
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 14px;
  overflow: hidden;
}

.table-wrapper,
.mac-dialog-table-wrapper {
  width: 100%;
  overflow-x: auto;
  overflow-y: hidden;
}

.table-wrapper::-webkit-scrollbar,
.mac-dialog-table-wrapper::-webkit-scrollbar {
  height: 10px;
}

.table-wrapper::-webkit-scrollbar-track,
.mac-dialog-table-wrapper::-webkit-scrollbar-track {
  background: #020617;
  border-radius: 999px;
}

.table-wrapper::-webkit-scrollbar-thumb,
.mac-dialog-table-wrapper::-webkit-scrollbar-thumb {
  background: #334155;
  border-radius: 999px;
  border: 2px solid #020617;
}

.table-wrapper::-webkit-scrollbar-thumb:hover,
.mac-dialog-table-wrapper::-webkit-scrollbar-thumb:hover {
  background: #475569;
}

.table-wrapper::-webkit-scrollbar-button,
.mac-dialog-table-wrapper::-webkit-scrollbar-button {
  display: none;
}

.table-wrapper,
.mac-dialog-table-wrapper {
  scrollbar-width: thin;
  scrollbar-color: #334155 #020617;
}

.model-table {
  width: 100%;
  min-width: 1500px;
  border-collapse: collapse;
  table-layout: fixed;
}

.model-table thead {
  background: #020617;
}

.model-table th,
.model-table td {
  box-sizing: border-box;
  white-space: nowrap;
}

.model-table th {
  padding: 14px 16px;
  color: #94a3b8;
  font-size: 12px;
  font-weight: 600;
  text-align: left;
  border-bottom: 1px solid #1e293b;
}

.model-table td {
  padding: 15px 16px;
  font-size: 13px;
  color: #e2e8f0;
  border-bottom: 1px solid #1e293b;
  vertical-align: middle;
  overflow: hidden;
}

.model-table th:nth-child(1),
.model-table td:nth-child(1) {
  width: 220px;
}

.model-table th:nth-child(2),
.model-table td:nth-child(2) {
  width: 120px;
}

.model-table th:nth-child(3),
.model-table td:nth-child(3) {
  width: 280px;
}

.model-table th:nth-child(4),
.model-table td:nth-child(4) {
  width: 120px;
}

.model-table th:nth-child(5),
.model-table td:nth-child(5) {
  width: 140px;
}

.model-table th:nth-child(6),
.model-table td:nth-child(6) {
  width: 130px;
}

.model-table th:nth-child(7),
.model-table td:nth-child(7) {
  width: 110px;
}

.model-table th:nth-child(8),
.model-table td:nth-child(8) {
  width: 340px;
}

.model-row {
  background: #0f172a;
}

.model-row:hover {
  background: #1e293b80;
}

.model-name-btn {
  display: inline-block;
  max-width: 190px;
  border: none;
  background: transparent;
  color: #60a5fa;
  font-size: 14px;
  font-weight: 800;
  cursor: pointer;
  padding: 0;
  text-align: left;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;
}

.model-name-btn:hover {
  color: #93c5fd;
  text-decoration: underline;
}

.count-tag {
  display: inline-block;
  padding: 4px 9px;
  border-radius: 999px;
  background: #33415566;
  color: #cbd5e1;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
}

.mac-tag,
.sn-tag,
.code-text,
.file-text,
.normal-text {
  display: inline-block;
  max-width: 100%;
  color: #cbd5e1;
  font-size: 12px;
  font-family: Consolas, Monaco, monospace;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;
}

.mac-tag,
.sn-tag {
  padding: 3px 8px;
  border-radius: 999px;
  background: #33415566;
}

.normal-text {
  font-family: inherit;
  font-size: 13px;
}

.status-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
}

.status-tag.draft {
  background: #47556933;
  color: #94a3b8;
}

.status-tag.submitted {
  background: #d9770633;
  color: #fbbf24;
}

.status-tag.approved {
  background: #16a34a33;
  color: #4ade80;
}

.status-tag.rejected {
  background: #dc262633;
  color: #f87171;
}

.muted {
  color: #94a3b8 !important;
}

.nowrap {
  white-space: nowrap !important;
}

.operation-col {
  text-align: left !important;
}

.action-group {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  gap: 10px;
  flex-wrap: nowrap;
  white-space: nowrap;
  min-width: 300px;
}

.text-btn {
  border: none;
  background: transparent;
  color: #cbd5e1;
  font-size: 13px;
  cursor: pointer;
  white-space: nowrap;
  padding: 0;
}

.text-btn:hover {
  color: #fff;
}

.text-btn.blue {
  color: #60a5fa;
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

.pagination-bar {
  padding: 12px 0 0;
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

.upload-tip {
  margin: 18px 20px 0;
  padding: 12px;
  background: #020617;
  border: 1px solid #1e293b;
  border-radius: 10px;
}

.upload-tip strong {
  color: #f8fafc;
  font-size: 13px;
}

.upload-tip p {
  margin: 6px 0 0;
  color: #94a3b8;
  font-size: 12px;
  line-height: 1.6;
}

.mac-select-panel {
  min-height: 120px;
  max-height: 260px;
  overflow-y: auto;
  border: 1px solid #334155;
  border-radius: 10px;
  background: #020617;
  padding: 10px;
}

.mac-select-panel::-webkit-scrollbar {
  width: 8px;
}

.mac-select-panel::-webkit-scrollbar-track {
  background: #020617;
}

.mac-select-panel::-webkit-scrollbar-thumb {
  background: #334155;
  border-radius: 999px;
}

.mac-search-input {
  width: 100%;
  box-sizing: border-box;
  margin-bottom: 10px;
}

.mac-panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  color: #64748b;
  font-size: 12px;
}

.mac-panel-header strong {
  color: #5eead4;
  font-size: 12px;
}

.exclude-tag {
  display: inline-flex;
  justify-content: center;
  align-items: center;
  height: 22px;
  padding: 0 8px;
  border-radius: 999px;
  background: #7f1d1d55;
  color: #fca5a5;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
}

.empty-mac,
.empty-dialog-data {
  color: #64748b;
  font-size: 13px;
  padding: 8px 0;
}

.empty-dialog-data {
  padding: 0 20px 20px;
}

.mac-check-item {
  display: grid !important;
  grid-template-columns: 18px 180px 130px 1fr 70px;
  align-items: center;
  gap: 8px !important;
  padding: 8px;
  border-bottom: 1px solid #1e293b;
  color: #cbd5e1 !important;
}

.mac-check-item:last-child {
  border-bottom: none;
}

.mac-check-item input {
  width: 14px;
  height: 14px;
}

.mac-check-item span {
  font-family: Consolas, Monaco, monospace;
  color: #e2e8f0;
  font-size: 12px;
}

.mac-check-item em {
  color: #64748b;
  font-size: 12px;
  font-style: normal;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* MAC 弹窗 */
.model-summary-card {
  padding: 20px;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 14px;
}

.model-summary-card div {
  background: #020617;
  border: 1px solid #1e293b;
  border-radius: 10px;
  padding: 12px;
}

.model-summary-card span {
  display: block;
  color: #64748b;
  font-size: 12px;
  margin-bottom: 6px;
}

.model-summary-card strong {
  color: #f8fafc;
  font-size: 14px;
  word-break: break-all;
}

.mac-dialog-table-wrapper {
  padding: 0 20px 20px;
  box-sizing: border-box;
}

.mac-dialog-table {
  width: 100%;
  min-width: 760px;
  border-collapse: collapse;
  table-layout: fixed;
  border: 1px solid #1e293b;
  border-radius: 10px;
  overflow: hidden;
}

.mac-dialog-table thead {
  background: #020617;
}

.mac-dialog-table th,
.mac-dialog-table td {
  padding: 13px 14px;
  border-bottom: 1px solid #1e293b;
  color: #e2e8f0;
  font-size: 13px;
  text-align: left;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.mac-dialog-table th {
  color: #94a3b8;
  font-size: 12px;
  font-weight: 600;
}

.mac-dialog-table th:nth-child(1),
.mac-dialog-table td:nth-child(1) {
  width: 70px;
}

.mac-dialog-table th:nth-child(2),
.mac-dialog-table td:nth-child(2) {
  width: 160px;
}

.mac-dialog-table th:nth-child(3),
.mac-dialog-table td:nth-child(3) {
  width: 190px;
}

.mac-dialog-table th:nth-child(4),
.mac-dialog-table td:nth-child(4) {
  width: 180px;
}

.mac-dialog-table th:nth-child(5),
.mac-dialog-table td:nth-child(5) {
  width: 180px;
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
  max-height: 92vh;
  overflow-y: auto;
  background: #0f172a;
  border: 1px solid #334155;
  border-radius: 16px;
  color: #f8fafc;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.45);
}

.large-dialog {
  width: 900px;
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

.form-grid {
  padding: 20px;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.form-grid label {
  display: flex;
  flex-direction: column;
  gap: 8px;
  color: #cbd5e1;
  font-size: 13px;
}

.full-row {
  grid-column: 1 / -1;
}

.dialog-footer {
  padding: 16px 20px;
  border-top: 1px solid #1e293b;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
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

.inline-link {
  border: none;
  background: transparent;
  color: #60a5fa;
  cursor: pointer;
  padding: 0;
  font-size: 14px;
  text-align: left;
}

.inline-link:hover {
  color: #93c5fd;
  text-decoration: underline;
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

@media (max-width: 960px) {
  .filter-card {
    grid-template-columns: 1fr;
  }

  .model-table {
    min-width: 1280px;
  }

  .form-grid,
  .detail-card,
  .model-summary-card {
    grid-template-columns: 1fr;
  }

  .mac-check-item {
    grid-template-columns: 18px 1fr;
  }

  .mac-check-item em {
    grid-column: 2;
  }

  .mac-dialog-table {
    min-width: 760px;
  }
}
</style>
