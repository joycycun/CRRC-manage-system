<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>库存情况管理</h1>
      </div>

      <button v-if="canReturnToBoardInbound" class="return-btn" type="button" @click="openReturnDialog">
        退回板卡入库
      </button>
    </div>

    <!-- 统计卡片 -->
    <div class="summary-grid two-col">
      <div class="summary-card">
        <span>当前库存总数</span>
        <strong>{{ totalCount }}</strong>
        <p>统计已完成生产测试并进入库存的设备</p>
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
        placeholder="搜索产品名称 / 产品型号 / 产品编码 / 产品序列（SN） / MAC / 版本"
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

      <select v-model="filters.inventoryStatus">
        <option value="">全部库存状态</option>
        <option
          v-for="status in inventoryStatusOptions"
          :key="status"
          :value="status"
        >
          {{ getInventoryStatusText(status) }}
        </option>
      </select>

      <button class="query-btn" @click="loadInventory">查询</button>
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
              <th>产品名称</th>
              <th>产品型号</th>
              <th>产品编码</th>
              <th>产品序列（SN）</th>
              <th>MAC地址</th>
              <th>PCB二维码</th>
              <th>软件版本</th>
              <th>硬件版本</th>
              <th>库存状态</th>
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
                <span class="normal-text" :title="item.productName">
                  {{ item.productName }}
                </span>
              </td>

              <td class="product-model-cell">
                <span class="product-model-text" :title="item.productModel">
                  {{ item.productModel }}
                </span>
              </td>

              <td>
                <span class="normal-text" :title="item.productCode">
                  {{ item.productCode }}
                </span>
              </td>

              <td>
                <span class="sn-tag" :title="item.sn">
                  {{ item.sn }}
                </span>
              </td>

              <td>
                <span class="mac-text" :title="item.macAddress">
                  {{ isHandsetItem(item) ? '无需' : item.macAddress }}
                </span>
              </td>

              <td>
                <span class="normal-text" :title="item.pcbQrCode">
                  {{ isHandsetItem(item) ? '无需' : item.pcbQrCode }}
                </span>
              </td>

              <td class="version-cell">
                <span class="software-tag" :title="item.softwareVersion">
                  {{ isHandsetItem(item) ? '无需' : item.softwareVersion }}
                </span>
              </td>

              <td class="version-cell">
                <span class="hardware-tag" :title="item.hardwareVersion">
                  {{ isHandsetItem(item) ? '无需' : item.hardwareVersion }}
                </span>
              </td>

              <td>
                <span class="status-tag" :class="getInventoryStatusClass(item.inventoryStatus)">
                  {{ getInventoryStatusText(item.inventoryStatus) }}
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
              <td colspan="12" class="empty-table">
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
        库存页面仅展示已完成生产测试且尚未出库的设备。发货批次审核通过后，设备会进入出库记录并从库存情况中移除。
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
            <span>产品名称</span>
            <strong>{{ selectedInventory.productName }}</strong>
          </div>

          <div>
            <span>产品型号</span>
            <strong>{{ selectedInventory.productModel }}</strong>
          </div>

          <div>
            <span>产品编码</span>
            <strong>{{ selectedInventory.productCode }}</strong>
          </div>

          <div>
            <span>产品序列（SN）</span>
            <strong>{{ selectedInventory.sn }}</strong>
          </div>

          <div v-if="!isHandsetItem(selectedInventory)">
            <span>MAC地址</span>
            <strong>{{ selectedInventory.macAddress }}</strong>
          </div>

          <div v-if="!isHandsetItem(selectedInventory)">
            <span>PCB二维码</span>
            <strong>{{ selectedInventory.pcbQrCode }}</strong>
          </div>

          <div v-if="!isHandsetItem(selectedInventory)">
            <span>软件版本</span>
            <strong>{{ selectedInventory.softwareVersion }}</strong>
          </div>

          <div v-if="!isHandsetItem(selectedInventory)">
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

          <div>
            <span>废弃审核状态</span>
            <strong>{{ selectedInventory.scrapAuditStatus || '未申请' }}</strong>
          </div>

          <div>
            <span>废弃申请人</span>
            <strong>{{ selectedInventory.scrapRequestUserName || '-' }}</strong>
          </div>

          <div>
            <span>废弃申请时间</span>
            <strong>{{ selectedInventory.scrapRequestTime || '-' }}</strong>
          </div>

          <div>
            <span>废弃审核人</span>
            <strong>{{ selectedInventory.scrapAuditUserName || '-' }}</strong>
          </div>
        </div>

        <div class="remark-card">
          <span>库存说明</span>
          <p>{{ selectedInventory.remark || '暂无说明' }}</p>
        </div>

        <div v-if="selectedInventory.scrapRejectReason" class="remark-card">
          <span>废弃驳回原因</span>
          <p>{{ selectedInventory.scrapRejectReason }}</p>
        </div>

        <div v-if="canSubmitScrapRequest" class="scrap-card">
          <label>
            <input v-model="scrapChecked" type="checkbox" />
            <span>确认申请废弃该库存设备</span>
          </label>
          <button class="red-btn" :disabled="!scrapChecked" @click="submitScrapRequest">
            提交废弃申请
          </button>
        </div>

        <div class="dialog-footer">
          <button
            v-if="canAuditScrapRequest"
            class="red-btn"
            @click="auditScrapRequest('rejected')"
          >
            审核驳回
          </button>

          <button
            v-if="canAuditScrapRequest"
            class="green-btn"
            @click="auditScrapRequest('approved')"
          >
            审核通过
          </button>

          <button class="primary-btn" @click="selectedInventory = null">
            关闭
          </button>
        </div>
      </div>
    </div>

    <div v-if="showReturnDialog" class="dialog-mask" @click.self="closeReturnDialog">
      <div class="dialog return-dialog">
        <div class="dialog-header">
          <h3>按产品型号退回板卡入库</h3>
          <button @click="closeReturnDialog">×</button>
        </div>

        <div class="return-form">
          <label>
            <span>产品型号</span>
            <select v-model="returnForm.productModel">
              <option value="">请选择需要退回的产品型号</option>
              <option v-for="item in returnTypeOptions" :key="item.productModel" :value="item.productModel">
                {{ item.productName }} / {{ item.productModel }}（{{ item.count }} 件）
              </option>
            </select>
          </label>

          <div v-if="selectedReturnType" class="return-summary">
            <div>
              <span>产品名称</span>
              <strong>{{ selectedReturnType.productName }}</strong>
            </div>
            <div>
              <span>退回成品数量</span>
              <strong>{{ selectedReturnType.count }}</strong>
            </div>
          </div>

          <label>
            <span>退回原因</span>
            <textarea v-model.trim="returnForm.reason" rows="3" placeholder="请填写本次退回原因"></textarea>
          </label>

          <p class="return-warning">
            该型号当前可退回的库存将全部移除，并恢复其烧录时扣减的板卡入库数量。已绑定发货批次的数据不会执行退回。
          </p>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" type="button" @click="closeReturnDialog">取消</button>
          <button class="return-btn" type="button" :disabled="returningInventory" @click="submitReturnToBoardInbound">
            {{ returningInventory ? '退回中...' : '确认退回' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import {
  auditInventoryScrap,
  getInventory,
  returnInventoryTypeToBoardInbound,
  submitInventoryScrap
} from '@/api/inventory'
import { canUseAction, hasLeaderRole, hasRole } from '@/utils/permission'

const route = useRoute()
const filters = reactive({
  keyword: '',
  deviceType: '',
  inventoryStatus: ''
})

const selectedInventory = ref(null)
const scrapChecked = ref(false)
const showReturnDialog = ref(false)
const returningInventory = ref(false)
const returnForm = reactive({
  productModel: '',
  reason: ''
})

const currentPage = ref(1)
const pageSize = ref(10)

const pageSizeOptions = [10, 20, 50, 100]

const inventoryList = ref([])

onMounted(async () => {
  syncKeywordFromRoute()
  await loadInventory()
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
  if (res && res.data) return res.data
  return res
}
function getInventoryStatusText(status) {
  const map = {
    in_stock: '在库',
    reserved: '已预占',
    shipped: '已出库',
    repairing: '维修中',
    scrapped: '已报废',
    board_inbound: '板卡入库',

    在库: '在库',
    板卡入库: '板卡入库',
    已预占: '已预占',
    已出库: '已出库',
    维修中: '维修中',
    返厂: '返厂',
    更换: '更换',
    已废弃: '已废弃',
    已报废: '已报废'
  }

  return map[status] || status || '-'
}
function getInventoryStatusClass(status) {
  const map = {
    in_stock: 'in-stock',
    reserved: 'reserved',
    shipped: 'shipped',
    repairing: 'repairing',
    scrapped: 'scrapped',
    board_inbound: 'in-stock',

    在库: 'in-stock',
    板卡入库: 'in-stock',
    已预占: 'reserved',
    已出库: 'shipped',
    维修中: 'repairing',
    返厂: 'repairing',
    更换: 'repairing',
    已废弃: 'scrapped',
    已报废: 'scrapped'
  }

  return map[status] || 'in-stock'
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

    inventoryList.value = (result.data || []).map(item => normalizeInventory(item))
  } catch (err) {
    console.error('加载库存失败：', err)
    alert(err.response?.data || '加载库存失败，请检查 /api/inventory 接口')
  }
}

function normalizeInventory(item) {
  return {
    id: item.id,

    projectId: item.projectId || item.project_id || 0,

    deviceType: item.deviceType || item.device_type || '-',
    productName: item.productName || item.product_name || '-',
    productModel: item.productModel || item.product_model || '-',
    productCode: item.productCode || item.product_code || '-',

    sn: item.sn || '-',
    macAddress: item.macAddress || item.mac_address || '-',
    pcbQrCode: item.pcbQrCode || item.pcb_qr_code || '-',

    hardwareId: item.hardwareId || item.hardware_id || 0,
    hardwareVersion: item.hardwareVersion || item.hardware_version || '-',

    softwareId: item.softwareId || item.software_id || 0,
    softwareVersion: item.softwareVersion || item.software_version || '-',

    inventoryStatus: item.inventoryStatus || item.inventory_status || '',
    scrapAuditStatus: item.scrapAuditStatus || item.scrap_audit_status || '',
    scrapRequestUserId: item.scrapRequestUserId || item.scrap_request_user_id || 0,
    scrapRequestUserName: item.scrapRequestUserName || item.scrap_request_user_name || '',
    scrapRequestTime: formatDateTime(item.scrapRequestTime || item.scrap_request_time),
    scrapAuditUserId: item.scrapAuditUserId || item.scrap_audit_user_id || 0,
    scrapAuditUserName: item.scrapAuditUserName || item.scrap_audit_user_name || '',
    scrapAuditTime: formatDateTime(item.scrapAuditTime || item.scrap_audit_time),
    scrapRejectReason: item.scrapRejectReason || item.scrap_reject_reason || '',

    sourceBurnRecordId:
      item.sourceBurnRecordId ||
      item.source_burn_record_id ||
      0,

    factoryTestId:
      item.factoryTestId ||
      item.factory_test_id ||
      0,

    inTime: formatDateTime(item.inTime || item.in_time),
    updateTime: formatDateTime(item.updateTime || item.update_time),

    remark: item.remark || ''
  }
}

function isHandsetItem(item) {
  const productName = String(item?.productName || '').replace(/\s+/g, '')
  const productModel = String(item?.productModel || '').trim().toLowerCase()
  return productName.includes('手持话柄') || productModel === 'handheld mic-zycoo'
}

function formatDateTime(value) {
  if (!value) return ''

  if (typeof value === 'string') {
    return value.slice(0, 19).replace('T', ' ')
  }

  if (value.Time) {
    return String(value.Time).slice(0, 19).replace('T', ' ')
  }

  return String(value)
}

// 不按库存状态过滤；后端返回什么库存记录，前端就展示什么
const completedInventoryList = computed(() => {
  return inventoryList.value
})

const deviceTypeOptions = computed(() => {
  const types = completedInventoryList.value
    .map(item => item.deviceType)
    .filter(Boolean)
    .filter(type => type !== '-')

  return [...new Set(types)]
})

const inventoryStatusOptions = computed(() => {
  const statuses = completedInventoryList.value
    .map(item => item.inventoryStatus)
    .filter(Boolean)

  return [...new Set(statuses)]
})

const filteredInventoryList = computed(() => {
  return completedInventoryList.value.filter(item => {
    const keyword = filters.keyword.trim().toLowerCase()

    const keywordMatch =
      !keyword ||
      item.sn.toLowerCase().includes(keyword) ||
      item.macAddress.toLowerCase().includes(keyword) ||
      item.softwareVersion.toLowerCase().includes(keyword) ||
      item.hardwareVersion.toLowerCase().includes(keyword) ||
      item.deviceType.toLowerCase().includes(keyword) ||
      item.productName.toLowerCase().includes(keyword) ||
      item.productModel.toLowerCase().includes(keyword) ||
      item.productCode.toLowerCase().includes(keyword) ||
      item.pcbQrCode.toLowerCase().includes(keyword) ||
      item.remark.toLowerCase().includes(keyword)

    const deviceTypeMatch =
      !filters.deviceType || item.deviceType === filters.deviceType

    const inventoryStatusMatch =
      !filters.inventoryStatus || item.inventoryStatus === filters.inventoryStatus

    return keywordMatch && deviceTypeMatch && inventoryStatusMatch
  })
})

const totalCount = computed(() => completedInventoryList.value.length)

// 现在没有接出库表接口，先按 inventory_status = 已出库 且 updateTime 是当月统计
const currentMonthOutboundCount = computed(() => {
  const currentMonth = new Date().toISOString().slice(0, 7)

  return inventoryList.value.filter(item => {
    return (
      item.inventoryStatus === '已出库' &&
      item.updateTime &&
      item.updateTime.startsWith(currentMonth)
    )
  }).length
})

const deviceTypeSummary = computed(() => {
  return deviceTypeOptions.value.map(type => {
    return {
      deviceType: type,
      count: completedInventoryList.value.filter(item => item.deviceType === type).length
    }
  })
})

const canSubmitScrapRequest = computed(() => {
  if (!selectedInventory.value) return false
  if (!hasRole('production_staff') && !hasRole('system_admin') && !canUseAction('production:update')) return false
  if (selectedInventory.value.scrapAuditStatus === '待审核') return false
  return !['已废弃', '已报废', '已出库'].includes(selectedInventory.value.inventoryStatus)
})

const canAuditScrapRequest = computed(() => {
  if (!selectedInventory.value) return false
  if (selectedInventory.value.scrapAuditStatus !== '待审核') return false
  return hasLeaderRole() || hasRole('system_admin') || canUseAction('production:audit')
})

const canReturnToBoardInbound = computed(() => hasRole('system_admin'))

const returnTypeOptions = computed(() => {
  const grouped = new Map()
  inventoryList.value
    .filter(item => item.sourceBurnRecordId > 0)
    .filter(item => ['在库', '返厂', '更换'].includes(item.inventoryStatus))
    .forEach(item => {
      if (!item.productModel || item.productModel === '-') return
      if (!grouped.has(item.productModel)) {
        grouped.set(item.productModel, {
          productModel: item.productModel,
          productName: item.productName,
          count: 0
        })
      }
      grouped.get(item.productModel).count += 1
    })

  return Array.from(grouped.values()).sort((a, b) => a.productModel.localeCompare(b.productModel))
})

const selectedReturnType = computed(() => {
  return returnTypeOptions.value.find(item => item.productModel === returnForm.productModel) || null
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

function resetFilters() {
  filters.keyword = ''
  filters.deviceType = ''
  filters.inventoryStatus = ''
  currentPage.value = 1
}

function viewInventory(item) {
  selectedInventory.value = item
  scrapChecked.value = false
}

function openReturnDialog() {
  returnForm.productModel = ''
  returnForm.reason = ''
  showReturnDialog.value = true
}

function closeReturnDialog() {
  if (returningInventory.value) return
  showReturnDialog.value = false
}

async function submitReturnToBoardInbound() {
  if (!returnForm.productModel) {
    alert('请选择需要退回的产品型号')
    return
  }
  if (!returnForm.reason) {
    alert('请填写退回原因')
    return
  }

  const selected = selectedReturnType.value
  const count = selected?.count || 0
  if (!confirm(`确认将【${returnForm.productModel}】的 ${count} 件库存全部退回板卡入库吗？`)) {
    return
  }

  returningInventory.value = true
  try {
    const res = await returnInventoryTypeToBoardInbound({
      productModel: returnForm.productModel,
      reason: returnForm.reason
    })
    const result = getResponseData(res)
    if (result.code !== 200) {
      alert(result.msg || '退回板卡入库失败')
      return
    }
    alert(`退回成功：移除 ${result.data?.inventoryCount || count} 件成品库存，恢复 ${result.data?.restoredBoardQuantity || 0} 件板卡库存`)
    showReturnDialog.value = false
    await loadInventory()
  } catch (err) {
    console.error('退回板卡入库失败：', err)
    alert(err.response?.data || '退回板卡入库失败')
  } finally {
    returningInventory.value = false
  }
}

async function submitScrapRequest() {
  if (!selectedInventory.value || !scrapChecked.value) return

  if (!confirm('确认提交该库存设备的废弃申请吗？')) {
    return
  }

  try {
    const res = await submitInventoryScrap(selectedInventory.value.id)
    const result = getResponseData(res)
    if (result.code !== 200) {
      alert(result.msg || '提交废弃申请失败')
      return
    }
    alert('废弃申请已提交，等待领导审核')
    selectedInventory.value = null
    await loadInventory()
  } catch (err) {
    console.error('提交废弃申请失败：', err)
    alert(err.response?.data || '提交废弃申请失败')
  }
}

async function auditScrapRequest(status) {
  if (!selectedInventory.value) return

  let rejectReason = ''
  if (status === 'rejected') {
    rejectReason = window.prompt('请输入驳回原因') || ''
  }

  const message = status === 'approved' ? '确认审核通过该废弃申请吗？' : '确认驳回该废弃申请吗？'
  if (!confirm(message)) {
    return
  }

  try {
    const res = await auditInventoryScrap(selectedInventory.value.id, {
      auditStatus: status,
      rejectReason
    })
    const result = getResponseData(res)
    if (result.code !== 200) {
      alert(result.msg || '废弃审核失败')
      return
    }
    alert('废弃审核已处理')
    selectedInventory.value = null
    await loadInventory()
  } catch (err) {
    console.error('废弃审核失败：', err)
    alert(err.response?.data || '废弃审核失败')
  }
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
.return-btn,
.page-btn {
  height: 36px;
  padding: 0 16px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
}

.return-btn {
  border: 1px solid #f59e0b;
  background: #b45309;
  color: #fff;
}

.return-btn:hover:not(:disabled) {
  background: #92400e;
}

.return-btn:disabled {
  opacity: 0.55;
  cursor: not-allowed;
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

.green-btn {
  border: none;
  background: #16a34a;
  color: #fff;
}

.green-btn:hover {
  background: #15803d;
}

.red-btn {
  border: none;
  background: #dc2626;
  color: #fff;
}

.red-btn:hover:not(:disabled) {
  background: #b91c1c;
}

.red-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
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
  grid-template-columns: 1.4fr 200px 180px 90px 90px;
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

.return-dialog {
  width: min(620px, calc(100vw - 32px));
}

.return-form {
  display: grid;
  gap: 18px;
  padding: 20px;
}

.return-form label {
  display: grid;
  gap: 8px;
}

.return-form label > span,
.return-summary span {
  color: #94a3b8;
  font-size: 12px;
}

.return-form select,
.return-form textarea {
  width: 100%;
  box-sizing: border-box;
  border: 1px solid #334155;
  border-radius: 8px;
  background: #020617;
  color: #e2e8f0;
  padding: 10px 12px;
  outline: none;
}

.return-form textarea {
  resize: vertical;
  min-height: 84px;
  line-height: 1.55;
}

.return-summary {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.return-summary > div {
  min-width: 0;
  border: 1px solid #1e293b;
  border-radius: 8px;
  background: #020617;
  padding: 12px;
}

.return-summary strong {
  display: block;
  margin-top: 7px;
  color: #f8fafc;
  overflow-wrap: anywhere;
}

.return-warning {
  margin: 0;
  border-left: 3px solid #f59e0b;
  background: #78350f33;
  color: #fcd34d;
  padding: 10px 12px;
  font-size: 12px;
  line-height: 1.6;
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
  min-width: 1680px;
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

.normal-text {
  display: inline-block;
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-model-cell {
  white-space: normal !important;
  overflow: visible !important;
}

.product-model-text {
  display: inline-block;
  min-width: 220px;
  max-width: 320px;
  white-space: normal;
  overflow: visible;
  word-break: break-word;
  line-height: 1.35;
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

.status-tag {
  display: inline-block;
  min-width: 58px;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
  text-align: center;
}

.status-tag.in-stock {
  background: #16a34a33;
  color: #4ade80;
}

.status-tag.reserved {
  background: #eab30833;
  color: #fde047;
}

.status-tag.shipped {
  background: #2563eb33;
  color: #93c5fd;
}

.status-tag.repairing {
  background: #f9731633;
  color: #fdba74;
}

.status-tag.scrapped {
  background: #dc262633;
  color: #fca5a5;
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

.scrap-card {
  margin: 0 20px 20px;
  background: #020617;
  border: 1px solid #7f1d1d;
  border-radius: 10px;
  padding: 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.scrap-card label {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #fecaca;
  font-size: 13px;
  cursor: pointer;
}

.scrap-card input {
  width: 16px;
  height: 16px;
  accent-color: #dc2626;
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
