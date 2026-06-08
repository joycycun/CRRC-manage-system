<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>发货批次管理</h1>
      </div>

      <button class="primary-btn" @click="openCreateDialog">
        新增发货批次
      </button>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索批次号 / SN序列号 / MAC地址 / 上传人 / 发货单文件 / 快递单号"
      />

      <select v-model="filters.auditStatus">
        <option value="">全部审核状态</option>
        <option value="draft">草稿</option>
        <option value="submitted">待审核</option>
        <option value="approved">审核通过</option>
        <option value="rejected">审核驳回</option>
      </select>

      <button class="query-btn">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <!-- 发货批次表格 -->
    <div class="table-card">
      <div class="table-card-header">
        <div>
          <h3>发货批次列表</h3>
          <span>共 {{ filteredBatchList.length }} 条发货批次</span>
        </div>
      </div>

      <div class="table-wrapper">
        <table class="version-table">
          <thead>
            <tr>
              <th>发货批次</th>
              <th>发货单文件</th>
              <th>设备数量</th>
              <th>上传人</th>
              <th>上传时间</th>
              <th>审核状态</th>
              <th>审核人</th>
              <th>快递单号</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="item in filteredBatchList" :key="item.id">
              <td>
                <button class="record-link" @click="viewBatch(item)">
                  {{ item.batchNo }}
                </button>
              </td>

              <td>
                <span class="file-tag" :title="item.fileName">
                  {{ item.fileName || '-' }}
                </span>
              </td>

              <td>
                <span class="count-tag">
                  {{ item.deviceCount }} 台
                </span>
              </td>

              <td>{{ item.uploader }}</td>

              <td class="muted">{{ item.uploadTime }}</td>

              <td>
                <span class="status-tag" :class="item.auditStatus">
                  {{ getAuditStatusText(item.auditStatus) }}
                </span>
              </td>

              <td>{{ item.auditor || '-' }}</td>

              <td>
                <span class="express-tag" :title="item.expressNo">
                  {{ item.expressNo || '-' }}
                </span>
              </td>

              <td class="operation-col">
                <div class="action-group">
                  <button class="text-btn" @click="viewBatch(item)">
                    查看
                  </button>

                  <button class="text-btn blue" @click="downloadFile(item)">
                    下载
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
              <td colspan="9" class="empty-table">
                暂无发货批次记录
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="table-footer">
        新增发货批次时，可直接从当前库存中选择设备；保存后设备会自动进入出库记录页面。
      </div>
    </div>

    <!-- 新增发货批次弹窗 -->
    <div v-if="showCreateDialog" class="dialog-mask">
      <div class="dialog create-dialog">
        <div class="dialog-header">
          <h3>新增发货批次</h3>
          <button @click="showCreateDialog = false">×</button>
        </div>

        <div class="form-grid">
          <label>
            发货批次号
            <input
              v-model="batchForm.batchNo"
              placeholder="例如：0000013289"
            />
          </label>

          <label>
            设备数量
            <input
              :value="selectedInventoryDevices.length"
              type="number"
              disabled
            />
          </label>

          <label>
            上传人
            <input
              v-model="batchForm.uploader"
              disabled
            />
          </label>

          <label>
            快递单号
            <input
              v-model="batchForm.expressNo"
              placeholder="例如：SF1234567890"
            />
          </label>

          <label class="full-row">
            从库存选择设备
            <div class="inventory-select-panel">
              <div class="inventory-select-tools">
                <input
                  v-model="inventoryFilters.keyword"
                  placeholder="按 SN / MAC / 终端类型 / 软件版本 / 硬件版本模糊筛选"
                />

                <select v-model="inventoryFilters.deviceType">
                  <option value="">全部终端类型</option>
                  <option
                    v-for="type in deviceTypeOptions"
                    :key="type"
                    :value="type"
                  >
                    {{ type }}
                  </option>
                </select>
              </div>

              <div class="mac-tool-card">
                <div class="mac-tool-title">
                  <strong>按 MAC 列表选择</strong>
                  <span>支持换行、逗号、空格分隔多个 MAC</span>
                </div>

                <textarea
                  v-model="macSelectForm.macText"
                  placeholder="例如：
00:11:22:33:44:01
00:11:22:33:44:02"
                ></textarea>

                <div class="mac-tool-actions">
                  <button class="reset-btn small-btn" @click="selectByMacList">
                    按列表选择
                  </button>
                  <button class="reset-btn small-btn" @click="clearMacText">
                    清空列表
                  </button>
                </div>
              </div>

              <div class="mac-tool-card">
                <div class="mac-tool-title">
                  <strong>按 MAC 区间选择</strong>
                  <span>适合同一批连续 MAC，例如 00:11:22:33:44:01 到 00:11:22:33:44:10</span>
                </div>

                <div class="range-grid">
                  <input
                    v-model="macSelectForm.startMac"
                    placeholder="起始 MAC"
                  />
                  <input
                    v-model="macSelectForm.endMac"
                    placeholder="结束 MAC"
                  />
                  <button class="reset-btn small-btn" @click="selectByMacRange">
                    按区间选择
                  </button>
                </div>
              </div>

              <div class="inventory-panel-header">
                <div>
                  <strong>当前库存设备</strong>
                  <span>
                    已选择 {{ selectedInventoryIds.length }} 台 /
                    当前筛选 {{ filteredInventoryList.length }} 台
                  </span>
                </div>

                <div class="inventory-header-actions">
                  <button class="reset-btn small-btn" @click="selectAllFilteredInventory">
                    选择当前筛选
                  </button>

                  <button class="reset-btn small-btn" @click="clearSelectedInventory">
                    清空选择
                  </button>
                </div>
              </div>

              <div class="inventory-table-wrapper">
                <table class="inventory-table">
                  <thead>
                    <tr>
                      <th class="check-col">
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
                    <tr
                      v-for="device in filteredInventoryList"
                      :key="device.id"
                    >
                      <td class="check-col">
                        <input
                          v-model="selectedInventoryIds"
                          type="checkbox"
                          :value="device.id"
                        />
                      </td>

                      <td>
                        <span class="device-tag">
                          {{ device.deviceType }}
                        </span>
                      </td>

                      <td>
                        <span class="sn-tag" :title="device.sn">
                          {{ device.sn }}
                        </span>
                      </td>

                      <td>
                        <span class="mac-text" :title="device.macAddress">
                          {{ device.macAddress }}
                        </span>
                      </td>

                      <td>
                        <span class="software-tag" :title="device.softwareVersion">
                          {{ device.softwareVersion }}
                        </span>
                      </td>

                      <td>
                        <span class="hardware-tag" :title="device.hardwareVersion">
                          {{ device.hardwareVersion }}
                        </span>
                      </td>
                    </tr>

                    <tr v-if="filteredInventoryList.length === 0">
                      <td colspan="6" class="empty-table">
                        当前库存中暂无符合条件的设备
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </label>

          <label class="full-row">
            已选设备
            <div class="selected-device-panel">
              <div v-if="selectedInventoryDevices.length === 0" class="empty-device">
                暂未选择设备。请在上方库存列表中勾选设备，或通过 MAC 列表 / MAC 区间自动选择。
              </div>

              <div
                v-for="device in selectedInventoryDevices"
                v-else
                :key="device.id"
                class="selected-device-item"
              >
                <span>{{ device.deviceType }}</span>
                <strong>{{ device.sn }}</strong>
                <em>{{ device.macAddress }}</em>
              </div>
            </div>
          </label>

          <label class="full-row">
            发货单 Excel 文件
            <input
              type="file"
              accept=".xls,.xlsx"
              @change="handleFileChange"
            />
          </label>

          <label class="full-row">
            文件名称
            <input
              v-model="batchForm.fileName"
              placeholder="选择文件后自动填充"
              disabled
            />
          </label>

          <label class="full-row">
            发货说明
            <textarea
              v-model="batchForm.remark"
              placeholder="例如：本批次发货单包含SN、MAC、终端类型、软件版本、硬件版本等信息"
            ></textarea>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showCreateDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="createBatch">
            保存
          </button>
        </div>
      </div>
    </div>

    <!-- 查看详情 / 审核弹窗 -->
    <div v-if="selectedBatch" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>发货批次详情</h3>
          <button @click="selectedBatch = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>发货批次号</span>
            <strong>{{ selectedBatch.batchNo }}</strong>
          </div>

          <div>
            <span>快递单号</span>
            <strong>{{ selectedBatch.expressNo || '-' }}</strong>
          </div>

          <div>
            <span>发货单文件</span>
            <strong>{{ selectedBatch.fileName || '-' }}</strong>
          </div>

          <div>
            <span>设备数量</span>
            <strong>{{ selectedBatch.deviceCount }} 台</strong>
          </div>

          <div>
            <span>上传人</span>
            <strong>{{ selectedBatch.uploader }}</strong>
          </div>

          <div>
            <span>上传时间</span>
            <strong>{{ selectedBatch.uploadTime }}</strong>
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

          <div>
            <span>文件查看</span>
            <button class="inline-link" @click="openFile(selectedBatch)">
              点开查看文件
            </button>
          </div>
        </div>

        <div class="remark-card">
          <span>发货设备</span>
          <div class="shipping-device-list">
            <div
              v-for="device in selectedBatch.deviceList"
              :key="device.sn"
              class="shipping-device-item"
            >
              <span>{{ device.deviceType }}</span>
              <strong>{{ device.sn }}</strong>
              <em>{{ device.macAddress }}</em>
            </div>

            <p v-if="!selectedBatch.deviceList || selectedBatch.deviceList.length === 0">
              暂无设备明细
            </p>
          </div>
        </div>

        <div class="remark-card">
          <span>发货说明</span>
          <p>{{ selectedBatch.remark || '暂无说明' }}</p>
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
import { computed, reactive, ref, watch } from 'vue'

const STORAGE_INVENTORY_KEY = 'inventoryList'
const STORAGE_OUTBOUND_KEY = 'outInventoryList'
const STORAGE_SHIPPING_BATCH_KEY = 'shippingBatchList'

const filters = reactive({
  keyword: '',
  auditStatus: ''
})

const inventoryFilters = reactive({
  keyword: '',
  deviceType: ''
})

const macSelectForm = reactive({
  macText: '',
  startMac: '',
  endMac: ''
})

const showCreateDialog = ref(false)
const selectedBatch = ref(null)
const selectedInventoryIds = ref([])

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
    inventoryStatus: 'ready',
    inTime: '2026-05-12',
    updateTime: '2026-05-13'
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
    inventoryStatus: 'ready',
    inTime: '2026-05-13',
    updateTime: '2026-05-14'
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
    inventoryStatus: 'ready',
    inTime: '2026-05-14',
    updateTime: '2026-05-15'
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

const defaultBatchList = [
  {
    id: 1,
    batchNo: '0000013289',
    deviceCount: 3,
    expressNo: 'SF202605100001',
    snList: [
      'DCCU-202605100001',
      'DCCU-202605100002',
      'DCCU-202605100003'
    ],
    deviceList: [],
    fileName: '5月第一批司机室控制盒发货单.xlsx',
    fileUrl: '',
    uploader: '袁晓兰',
    uploadTime: '2026-05-10',
    auditStatus: 'approved',
    auditor: '领导',
    auditTime: '2026-05-11',
    remark: '本批次包含司机室控制盒，发货单记录SN、MAC、软件版本、硬件版本。'
  },
  {
    id: 2,
    batchNo: '0000013290',
    deviceCount: 3,
    expressNo: 'SF202605160001',
    snList: [
      'DEC-202605110001',
      'DEC-202605110002',
      'DEC-202605110003'
    ],
    deviceList: [],
    fileName: '5月第二批解码板发货单.xlsx',
    fileUrl: '',
    uploader: '袁晓兰',
    uploadTime: '2026-05-16',
    auditStatus: 'submitted',
    auditor: '',
    auditTime: '',
    remark: '待领导审核确认。'
  }
]

const batchForm = reactive({
  batchNo: '',
  uploader: '',
  expressNo: '',
  fileName: '',
  file: null,
  fileUrl: '',
  remark: ''
})

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
const batchList = ref(readStorageList(STORAGE_SHIPPING_BATCH_KEY, defaultBatchList))

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

watch(
  batchList,
  value => {
    saveStorageList(STORAGE_SHIPPING_BATCH_KEY, value)
  },
  { deep: true }
)

const availableInventoryList = computed(() => {
  return inventoryList.value.filter(item => item.inventoryStatus === 'ready')
})

const filteredInventoryList = computed(() => {
  return availableInventoryList.value.filter(item => {
    const keyword = inventoryFilters.keyword.trim()

    const keywordMatch =
      !keyword ||
      item.sn.includes(keyword) ||
      item.macAddress.includes(keyword) ||
      item.deviceType.includes(keyword) ||
      item.softwareVersion.includes(keyword) ||
      item.hardwareVersion.includes(keyword)

    const deviceTypeMatch =
      !inventoryFilters.deviceType || item.deviceType === inventoryFilters.deviceType

    return keywordMatch && deviceTypeMatch
  })
})

const selectedInventoryDevices = computed(() => {
  return availableInventoryList.value.filter(item =>
    selectedInventoryIds.value.includes(item.id)
  )
})

const isAllFilteredInventorySelected = computed(() => {
  return (
    filteredInventoryList.value.length > 0 &&
    filteredInventoryList.value.every(item =>
      selectedInventoryIds.value.includes(item.id)
    )
  )
})

const filteredBatchList = computed(() => {
  return batchList.value.filter(item => {
    const keyword = filters.keyword.trim()

    const snMatch =
      item.snList &&
      item.snList.some(sn => sn.includes(keyword))

    const macMatch =
      item.deviceList &&
      item.deviceList.some(device => device.macAddress && device.macAddress.includes(keyword))

    const keywordMatch =
      !keyword ||
      item.batchNo.includes(keyword) ||
      item.uploader.includes(keyword) ||
      item.fileName.includes(keyword) ||
      item.expressNo.includes(keyword) ||
      snMatch ||
      macMatch

    const auditStatusMatch =
      !filters.auditStatus || item.auditStatus === filters.auditStatus

    return keywordMatch && auditStatusMatch
  })
})

function getCurrentUserName() {
  return (
    localStorage.getItem('username') ||
    localStorage.getItem('accountName') ||
    localStorage.getItem('realName') ||
    '当前用户'
  )
}

function getAuditStatusText(status) {
  const map = {
    draft: '草稿',
    submitted: '待审核',
    approved: '审核通过',
    rejected: '审核驳回'
  }

  return map[status] || status
}

function resetFilters() {
  filters.keyword = ''
  filters.auditStatus = ''
}

function openCreateDialog() {
  batchForm.batchNo = ''
  batchForm.uploader = getCurrentUserName()
  batchForm.expressNo = ''
  batchForm.fileName = ''
  batchForm.file = null
  batchForm.fileUrl = ''
  batchForm.remark = ''

  inventoryFilters.keyword = ''
  inventoryFilters.deviceType = ''
  macSelectForm.macText = ''
  macSelectForm.startMac = ''
  macSelectForm.endMac = ''
  selectedInventoryIds.value = []

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
  return String(mac || '')
    .trim()
    .replace(/-/g, ':')
    .toUpperCase()
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

  if (!/^[0-9A-F]{12}$/.test(hex)) {
    return null
  }

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

  selectedInventoryIds.value = [
    ...new Set([
      ...selectedInventoryIds.value,
      ...matchedIds
    ])
  ]

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

  selectedInventoryIds.value = [
    ...new Set([
      ...selectedInventoryIds.value,
      ...matchedIds
    ])
  ]

  alert(`已按 MAC 区间选择 ${matchedIds.length} 台库存设备`)
}

function clearMacText() {
  macSelectForm.macText = ''
}

function selectAllFilteredInventory() {
  const ids = filteredInventoryList.value.map(item => item.id)

  selectedInventoryIds.value = [
    ...new Set([
      ...selectedInventoryIds.value,
      ...ids
    ])
  ]
}

function toggleAllFilteredInventory(event) {
  const ids = filteredInventoryList.value.map(item => item.id)

  if (event.target.checked) {
    selectedInventoryIds.value = [
      ...new Set([
        ...selectedInventoryIds.value,
        ...ids
      ])
    ]
  } else {
    selectedInventoryIds.value = selectedInventoryIds.value.filter(
      id => !ids.includes(id)
    )
  }
}

function clearSelectedInventory() {
  selectedInventoryIds.value = []
}

function createBatch() {
  if (!batchForm.batchNo) {
    alert('请输入发货批次号')
    return
  }

  if (selectedInventoryDevices.value.length === 0) {
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

  const today = new Date().toISOString().slice(0, 10)

  // 这里保存的就是你在新增批次弹窗里选中的 MAC / SN
  const devices = selectedInventoryDevices.value.map(item => {
    return {
      id: item.id,
      deviceType: item.deviceType,
      sn: item.sn,
      macAddress: item.macAddress,
      softwareVersion: item.softwareVersion,
      hardwareVersion: item.hardwareVersion,
      inventoryStatus: item.inventoryStatus,
      inTime: item.inTime,
      updateTime: item.updateTime,
      remark: item.remark || ''
    }
  })

  batchList.value.unshift({
    id: Date.now(),
    batchNo: batchForm.batchNo,
    deviceCount: devices.length,
    expressNo: batchForm.expressNo,
    snList: devices.map(item => item.sn),
    macList: devices.map(item => item.macAddress),

    // 重点：选中的 MAC / SN 都存在这里
    deviceList: devices,

    fileName: batchForm.fileName,
    fileUrl: batchForm.fileUrl,
    uploader: batchForm.uploader,
    uploadTime: today,

    auditStatus: 'draft',
    auditor: '',
    auditTime: '',

    // 还没有审核通过，所以还不是出库
    shippingStatus: 'not_shipped',
    shippingTime: '',

    remark: batchForm.remark
  })

  selectedInventoryIds.value = []
  showCreateDialog.value = false

  alert('发货批次已生成，审核通过后才会进入出库管理')
}

function viewBatch(item) {
  selectedBatch.value = item
}

function auditBatch(item) {
  selectedBatch.value = item
}

function submitBatch(item) {
  item.auditStatus = 'submitted'
  item.auditor = ''
  item.auditTime = ''
  alert(`发货批次【${item.batchNo}】已提交领导审核`)
}

function approveBatch(item) {
  const ok = confirm(
    `确认发货批次【${item.batchNo}】审核通过吗？审核通过后，该批次选择的 MAC / SN 会自动进入出库管理。`
  )

  if (!ok) return

  const today = new Date().toISOString().slice(0, 10)

  item.auditStatus = 'approved'
  item.auditor = '领导'
  item.auditTime = today
  item.shippingStatus = 'shipped'
  item.shippingTime = today

  const devices = item.deviceList || []

  if (devices.length === 0) {
    alert('当前发货批次没有选择任何 MAC / SN，无法生成出库记录')
    return
  }

  const existingOutboundMacSet = new Set(
    outboundList.value.map(record => record.macAddress)
  )

  const outboundItems = devices
    .filter(device => !existingOutboundMacSet.has(device.macAddress))
    .map((device, index) => {
      return {
        outboundId: Date.now() + index,

        // 这些就是新增批次时选择的 MAC / SN
        deviceType: device.deviceType,
        sn: device.sn,
        macAddress: device.macAddress,
        softwareVersion: device.softwareVersion,
        hardwareVersion: device.hardwareVersion,
        inTime: device.inTime,

        outboundTime: today,
        operator: item.uploader,
        outboundStatus: 'shipped',

        shippingBatchNo: item.batchNo,
        expressNo: item.expressNo,

        remark: `由发货批次【${item.batchNo}】审核通过后自动生成出库记录。`
      }
    })

  if (outboundItems.length === 0) {
    alert('该批次中的设备已经存在于出库管理中，不再重复生成')
    selectedBatch.value = null
    return
  }

  // 重点：写入出库管理页面读取的 outInventoryList
  outboundList.value.unshift(...outboundItems)

  // 审核通过后，从库存中移除这些设备
  const outboundIds = devices.map(device => device.id)

  inventoryList.value = inventoryList.value.filter(
    inventory => !outboundIds.includes(inventory.id)
  )

  selectedBatch.value = null

  alert(
    `发货批次【${item.batchNo}】审核通过，${outboundItems.length} 台设备已进入出库管理`
  )
}

function rejectBatch(item) {
  item.auditStatus = 'rejected'
  item.auditor = '领导'
  item.auditTime = new Date().toISOString().slice(0, 10)
  selectedBatch.value = null
  alert(`发货批次【${item.batchNo}】已驳回`)
}

function openFile(item) {
  if (!item.fileUrl) {
    alert('当前是模拟数据，暂无可直接打开的原始文件')
    return
  }

  window.open(item.fileUrl, '_blank')
}

function downloadFile(item) {
  if (!item.fileUrl) {
    alert('当前是模拟数据，暂无可下载的原始文件')
    return
  }

  const link = document.createElement('a')
  link.href = item.fileUrl
  link.download = item.fileName || '发货单.xlsx'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

function deleteBatch(item) {
  const ok = confirm(`确认删除发货批次【${item.batchNo}】吗？`)
  if (!ok) return

  batchList.value = batchList.value.filter(record => record.id !== item.id)
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

.small-btn {
  height: 30px;
  padding: 0 12px;
  font-size: 12px;
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
  grid-template-columns: 1.4fr 180px 90px 90px;
  gap: 12px;
  margin-bottom: 20px;
}

.filter-card input,
.filter-card select,
.form-grid input,
.form-grid select,
.form-grid textarea,
.inventory-select-tools input,
.inventory-select-tools select,
.range-grid input,
.mac-tool-card textarea {
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
.inventory-select-tools input,
.inventory-select-tools select,
.range-grid input {
  height: 36px;
}

.form-grid textarea,
.mac-tool-card textarea {
  min-height: 90px;
  padding: 10px 12px;
  resize: vertical;
}

.filter-card input::placeholder,
.form-grid input::placeholder,
.form-grid textarea::placeholder,
.mac-tool-card textarea::placeholder,
.range-grid input::placeholder,
.inventory-select-tools input::placeholder {
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

.table-wrapper,
.inventory-table-wrapper {
  width: 100%;
  overflow-x: auto;
  overflow-y: hidden;
}

.table-wrapper::-webkit-scrollbar,
.inventory-table-wrapper::-webkit-scrollbar,
.selected-device-panel::-webkit-scrollbar {
  height: 10px;
  width: 8px;
}

.table-wrapper::-webkit-scrollbar-track,
.inventory-table-wrapper::-webkit-scrollbar-track,
.selected-device-panel::-webkit-scrollbar-track {
  background: #020617;
  border-radius: 999px;
}

.table-wrapper::-webkit-scrollbar-thumb,
.inventory-table-wrapper::-webkit-scrollbar-thumb,
.selected-device-panel::-webkit-scrollbar-thumb {
  background: #334155;
  border-radius: 999px;
  border: 2px solid #020617;
}

.table-wrapper::-webkit-scrollbar-thumb:hover,
.inventory-table-wrapper::-webkit-scrollbar-thumb:hover,
.selected-device-panel::-webkit-scrollbar-thumb:hover {
  background: #475569;
}

.version-table {
  width: 100%;
  min-width: 1320px;
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

.record-link {
  display: inline-block;
  max-width: 130px;
  border: none;
  background: transparent;
  color: #60a5fa;
  font-size: 13px;
  font-weight: 800;
  cursor: pointer;
  padding: 0;
  text-align: left;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.record-link:hover {
  color: #93c5fd;
  text-decoration: underline;
}

.file-tag,
.express-tag {
  display: inline-block;
  width: 150px;
  max-width: 150px;
  padding: 4px 9px;
  border-radius: 999px;
  background: #33415566;
  color: #cbd5e1;
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;
}

.express-tag {
  background: #1d4ed833;
  color: #60a5fa;
  font-family: Consolas, Monaco, monospace;
}

.count-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  background: #0f766e33;
  color: #5eead4;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
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

.operation-col {
  width: 300px;
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

.empty-table {
  text-align: center;
  color: #64748b !important;
  padding: 28px 16px !important;
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
  max-height: 92vh;
  overflow-y: auto;
  background: #0f172a;
  border: 1px solid #334155;
  border-radius: 16px;
  color: #f8fafc;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.45);
}

.create-dialog {
  width: 980px;
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

.inventory-select-panel {
  background: #020617;
  border: 1px solid #334155;
  border-radius: 12px;
  padding: 12px;
}

.inventory-select-tools {
  display: grid;
  grid-template-columns: 1fr 180px;
  gap: 12px;
  margin-bottom: 12px;
}

.mac-tool-card {
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 10px;
  padding: 12px;
  margin-bottom: 12px;
}

.mac-tool-title {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 10px;
}

.mac-tool-title strong {
  color: #f8fafc;
  font-size: 13px;
}

.mac-tool-title span {
  color: #64748b;
  font-size: 12px;
}

.mac-tool-actions {
  display: flex;
  gap: 10px;
  margin-top: 10px;
}

.range-grid {
  display: grid;
  grid-template-columns: 1fr 1fr 120px;
  gap: 10px;
}

.inventory-panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin: 12px 0;
}

.inventory-panel-header strong {
  color: #f8fafc;
  font-size: 13px;
  margin-right: 10px;
}

.inventory-panel-header span {
  color: #64748b;
  font-size: 12px;
}

.inventory-header-actions {
  display: flex;
  gap: 10px;
}

.inventory-table {
  width: 100%;
  min-width: 900px;
  border-collapse: collapse;
  table-layout: fixed;
  border: 1px solid #1e293b;
}

.inventory-table thead {
  background: #0f172a;
}

.inventory-table th,
.inventory-table td {
  padding: 12px 14px;
  border-bottom: 1px solid #1e293b;
  color: #e2e8f0;
  font-size: 12px;
  text-align: left;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.inventory-table th {
  color: #94a3b8;
}

.check-col {
  width: 48px;
  text-align: center !important;
}

.selected-device-panel {
  min-height: 100px;
  max-height: 220px;
  overflow-y: auto;
  border: 1px solid #334155;
  border-radius: 10px;
  background: #020617;
  padding: 10px;
}

.selected-device-item,
.shipping-device-item {
  display: grid;
  grid-template-columns: 140px 1fr 180px;
  align-items: center;
  gap: 10px;
  padding: 8px;
  border-bottom: 1px solid #1e293b;
}

.selected-device-item:last-child,
.shipping-device-item:last-child {
  border-bottom: none;
}

.selected-device-item span,
.shipping-device-item span {
  color: #5eead4;
  font-size: 12px;
}

.selected-device-item strong,
.selected-device-item em,
.shipping-device-item strong,
.shipping-device-item em {
  color: #cbd5e1;
  font-size: 12px;
  font-family: Consolas, Monaco, monospace;
  font-style: normal;
}

.empty-device {
  color: #64748b;
  font-size: 13px;
  padding: 10px 0;
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
  max-width: 170px;
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
  max-width: 180px;
  padding: 4px 9px;
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

.shipping-device-list {
  display: flex;
  flex-direction: column;
}

@media (max-width: 960px) {
  .filter-card,
  .form-grid,
  .inventory-select-tools,
  .range-grid,
  .detail-card {
    grid-template-columns: 1fr;
  }

  .version-table {
    min-width: 1320px;
  }

  .selected-device-item,
  .shipping-device-item {
    grid-template-columns: 1fr;
  }

  .inventory-panel-header {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>