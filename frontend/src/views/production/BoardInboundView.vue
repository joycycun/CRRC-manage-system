<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1>板卡入库</h1>
      </div>

      <button v-if="canImport" class="primary-btn" @click="openUploadDialog">
        上传板卡入库 Excel
      </button>
    </div>

    <div class="summary-grid two-col">
      <div class="summary-card">
        <span>板卡入库总数</span>
        <strong>{{ totalCount }}</strong>
        <p>统计通过板卡入库 Excel 导入的库存记录</p>
      </div>

      <div class="summary-card green">
        <span>本月入库数量</span>
        <strong>{{ currentMonthInboundCount }}</strong>
        <p>按入库时间统计当前月份的板卡入库记录</p>
      </div>
    </div>

    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索产品名称 / 产品型号 / 产品编码 / 序列号 / MAC地址 / PCB二维码 / 备注"
      />

      <select v-model="filters.productModel">
        <option value="">全部产品型号</option>
        <option
          v-for="model in productModelOptions"
          :key="model"
          :value="model"
        >
          {{ model }}
        </option>
      </select>

      <button class="query-btn" @click="loadBoardInboundRecords">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <div class="type-card">
      <div class="section-title">
        <h3>产品型号入库统计</h3>
        <span>按产品型号统计当前板卡入库数量</span>
      </div>

      <div class="type-grid">
        <div
          v-for="item in productModelSummary"
          :key="item.productModel"
          class="type-item"
        >
          <span>{{ item.productModel }}</span>
          <strong>{{ item.count }}</strong>
        </div>

        <div v-if="productModelSummary.length === 0" class="type-item">
          <span>暂无数据</span>
          <strong>0</strong>
        </div>
      </div>
    </div>

    <div class="table-card">
      <div class="table-card-header">
        <div>
          <h3>板卡入库记录</h3>
          <span>共 {{ filteredList.length }} 条，当前第 {{ currentPage }} / {{ totalPage }} 页</span>
        </div>

        <div class="page-size-control">
          <span>每页</span>
          <select v-model.number="pageSize">
            <option v-for="size in pageSizeOptions" :key="size" :value="size">
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
              <th>产品名称</th>
              <th>产品型号</th>
              <th>产品编码</th>
              <th>序列号</th>
              <th>MAC地址</th>
              <th>PCB二维码</th>
              <th>入库时间</th>
              <th>来源文件</th>
              <th>备注</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="item in paginatedList" :key="item.id">
              <td><span class="product-tag" :title="item.productName">{{ item.productName }}</span></td>
              <td><span class="model-text" :title="item.productModel">{{ item.productModel }}</span></td>
              <td><span class="code-text" :title="item.productCode">{{ item.productCode }}</span></td>
              <td><span class="sn-tag" :title="item.sn">{{ item.sn }}</span></td>
              <td><span class="mac-text" :title="item.macAddress">{{ item.macAddress }}</span></td>
              <td><span class="pcb-text" :title="item.pcbQrCode">{{ item.pcbQrCode }}</span></td>
              <td class="muted nowrap">{{ item.inTime || '-' }}</td>
              <td>
                <button
                  v-if="item.sourceFileId"
                  class="inline-link"
                  @click="downloadFile(item)"
                >
                  {{ item.sourceFileName || '下载源文件' }}
                </button>
                <span v-else>-</span>
              </td>
              <td><span class="remark-text" :title="item.remark">{{ item.remark || '-' }}</span></td>
            </tr>

            <tr v-if="paginatedList.length === 0">
              <td colspan="9" class="empty-table">暂无板卡入库记录</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="pagination-bar">
        <div class="pagination-info">
          当前显示第 {{ pageStartIndex }} - {{ pageEndIndex }} 条，共 {{ filteredList.length }} 条
        </div>
        <div class="pagination-actions">
          <button class="page-btn" :disabled="currentPage === 1" @click="currentPage = 1">首页</button>
          <button class="page-btn" :disabled="currentPage === 1" @click="currentPage -= 1">上一页</button>
          <span class="page-number">{{ currentPage }} / {{ totalPage }}</span>
          <button class="page-btn" :disabled="currentPage === totalPage" @click="currentPage += 1">下一页</button>
          <button class="page-btn" :disabled="currentPage === totalPage" @click="currentPage = totalPage">末页</button>
        </div>
      </div>
    </div>

    <div v-if="showUploadDialog" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>上传板卡入库 Excel</h3>
          <button @click="showUploadDialog = false">×</button>
        </div>

        <div class="upload-tip">
          <strong>Excel 读取规则：</strong>
          <p>
            系统会从 Excel 第 3 行开始读取表头。请保证第 3 行包含：
            产品名称、产品型号、产品编码、序列号、MAC地址、PCB二维码、备注。
          </p>
        </div>

        <div class="form-grid">
          <label class="full-row">
            文件名称
            <input v-model="uploadForm.fileName" disabled placeholder="选择 Excel 后自动填充" />
          </label>

          <label class="full-row">
            板卡入库记录 Excel
            <input
              ref="excelFileInput"
              type="file"
              accept=".xls,.xlsx"
              @change="handleExcelFileChange"
            />
            <button type="button" class="file-picker-btn" @click="excelFileInput?.click()">
              选择 Excel 文件
            </button>
            <span class="file-name">
              {{ uploadForm.fileName || '请选择 .xls / .xlsx 文件' }}
            </span>
          </label>
        </div>

        <div v-if="excelPreviewList.length > 0" class="preview-card">
          <div class="preview-title">
            已识别 {{ excelPreviewList.length }} 条板卡入库记录
          </div>

          <div class="preview-wrapper">
            <table>
              <thead>
                <tr>
                  <th>产品名称</th>
                  <th>产品型号</th>
                  <th>产品编码</th>
                  <th>序列号</th>
                  <th>MAC地址</th>
                  <th>PCB二维码</th>
                  <th>备注</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(item, index) in excelPreviewList.slice(0, 8)" :key="index">
                  <td>{{ item.productName }}</td>
                  <td>{{ item.productModel }}</td>
                  <td>{{ item.productCode }}</td>
                  <td>{{ item.sn }}</td>
                  <td>{{ item.macAddress }}</td>
                  <td>{{ item.pcbQrCode }}</td>
                  <td>{{ item.remark }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <p v-if="excelPreviewList.length > 8" class="preview-more">
            仅预览前 8 行，保存后会导入全部记录。
          </p>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showUploadDialog = false">取消</button>
          <button class="primary-btn" @click="saveBoardInboundRecords">保存导入</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import * as XLSX from 'xlsx'
import { buildUploadFilePayload, downloadLocalFile, getFileDownloadUrl } from '@/utils/filePreview'
import { getBoardInboundRecords, importBoardInboundRecords } from '@/api/boardInbound'
import { canUseAction } from '@/utils/permission'

const filters = reactive({
  keyword: '',
  productModel: ''
})
const boardInboundList = ref([])
const excelPreviewList = ref([])
const showUploadDialog = ref(false)
const excelFileInput = ref(null)
const currentPage = ref(1)
const pageSize = ref(10)
const pageSizeOptions = [10, 20, 50, 100]
const canImport = computed(() => canUseAction('board-inbound:import'))

const uploadForm = reactive({
  file: null,
  fileId: 0,
  fileName: '',
  fileContentType: '',
  fileData: ''
})

onMounted(async () => {
  await loadBoardInboundRecords()
})

function getResponseData(res) {
  return res?.data || res
}

function normalizeHeader(value) {
  return String(value || '').trim().replace(/\s+/g, '')
}

function getCellValue(row, names) {
  for (const name of names) {
    if (row[name] !== undefined && row[name] !== null && row[name] !== '') {
      return String(row[name]).trim()
    }
  }
  return ''
}

function normalizeExcelRow(row) {
  const normalizedRow = {}
  Object.keys(row).forEach(key => {
    normalizedRow[normalizeHeader(key)] = row[key]
  })

  return {
    productName: getCellValue(normalizedRow, ['产品名称']),
    productModel: getCellValue(normalizedRow, ['产品型号']),
    productCode: getCellValue(normalizedRow, ['产品编码']),
    sn: getCellValue(normalizedRow, ['序列号', 'SN', 'SN序列号']),
    macAddress: getCellValue(normalizedRow, ['MAC地址', 'MAC']),
    pcbQrCode: getCellValue(normalizedRow, ['PCB二维码', 'PCB码', 'PCBQRCode']),
    remark: getCellValue(normalizedRow, ['备注', '说明'])
  }
}

function validateRequiredHeaders(row) {
  const headers = Object.keys(row || {}).map(key => normalizeHeader(key))
  const required = ['产品名称', '产品型号', '产品编码', '序列号', 'MAC地址', 'PCB二维码', '备注']
  return required.filter(header => !headers.includes(header))
}

function formatDuplicateMessages(duplicates) {
  if (duplicates.length === 0) return ''
  const lines = duplicates
    .slice(0, 20)
    .map(item => `${item.type}【${item.value}】：第 ${item.rowNumber} 行重复第 ${item.firstRowNumber} 行`)
  const more = duplicates.length > 20 ? `\n还有 ${duplicates.length - 20} 条重复未显示` : ''
  return lines.join('\n') + more
}

async function loadBoardInboundRecords() {
  try {
    const res = await getBoardInboundRecords()
    const result = getResponseData(res)
    if (result.code !== 200) {
      alert(result.msg || '加载板卡入库记录失败')
      return
    }
    boardInboundList.value = (result.data || []).map(item => ({
      id: item.id,
      productName: item.productName || item.product_name || '-',
      productModel: item.productModel || item.product_model || '-',
      productCode: item.productCode || item.product_code || '-',
      sn: item.sn || '-',
      macAddress: item.macAddress || item.mac_address || '-',
      pcbQrCode: item.pcbQrCode || item.pcb_qr_code || '-',
      remark: item.remark || '',
      sourceFileId: item.sourceFileId || item.source_file_id || 0,
      sourceFileName: item.sourceFileName || item.source_file_name || '',
      inTime: formatDateTime(item.inTime || item.in_time),
      updateTime: formatDateTime(item.updateTime || item.update_time)
    }))
  } catch (err) {
    console.error('加载板卡入库记录失败：', err)
    alert(err.response?.data || '加载板卡入库记录失败')
  }
}

function formatDateTime(value) {
  if (!value) return ''
  if (typeof value === 'string') return value.slice(0, 19).replace('T', ' ')
  if (value.Time) return String(value.Time).slice(0, 19).replace('T', ' ')
  return String(value)
}

function openUploadDialog() {
  uploadForm.file = null
  uploadForm.fileId = 0
  uploadForm.fileName = ''
  uploadForm.fileContentType = ''
  uploadForm.fileData = ''
  excelPreviewList.value = []
  if (excelFileInput.value) {
    excelFileInput.value.value = ''
  }
  showUploadDialog.value = true
}

async function handleExcelFileChange(event) {
  const file = event.target.files[0]
  if (!file) return

  const isExcel = file.name.endsWith('.xls') || file.name.endsWith('.xlsx')
  if (!isExcel) {
    alert('请上传 Excel 文件，格式为 .xls 或 .xlsx')
    event.target.value = ''
    return
  }

  uploadForm.file = file
  uploadForm.fileName = file.name

  try {
    const payload = await buildUploadFilePayload(file)
    uploadForm.fileId = payload.fileId
    uploadForm.fileContentType = payload.fileContentType
    uploadForm.fileData = payload.fileData
  } catch (err) {
    alert('读取板卡入库 Excel 文件失败，请重新选择')
    return
  }

  const reader = new FileReader()
  reader.onload = e => {
    try {
      const data = new Uint8Array(e.target.result)
      const workbook = XLSX.read(data, { type: 'array' })
      const worksheet = workbook.Sheets[workbook.SheetNames[0]]
      const rows = XLSX.utils.sheet_to_json(worksheet, { defval: '', range: 2 })

      if (rows.length === 0) {
        alert('Excel 中未识别到有效数据，请确认第 3 行是否为表头')
        excelPreviewList.value = []
        return
      }

      const missingHeaders = validateRequiredHeaders(rows[0])
      if (missingHeaders.length > 0) {
        alert(`Excel 第 3 行缺少表头：${missingHeaders.join('、')}`)
        excelPreviewList.value = []
        return
      }

      const seenSN = new Map()
      const seenMac = new Map()
      const duplicateMessages = []
      const parsedRows = rows
        .map((row, index) => ({
          ...normalizeExcelRow(row),
          rowNumber: index + 4
        }))
        .filter(row => row.productName || row.productModel || row.productCode || row.sn || row.macAddress || row.pcbQrCode || row.remark)
        .filter(row => {
          const sn = row.sn.toUpperCase()
          const mac = row.macAddress.toUpperCase()
          let isDuplicate = false
          if (sn && seenSN.has(sn)) {
            duplicateMessages.push({
              type: '序列号',
              value: row.sn,
              rowNumber: row.rowNumber,
              firstRowNumber: seenSN.get(sn)
            })
            isDuplicate = true
          }
          if (mac && seenMac.has(mac)) {
            duplicateMessages.push({
              type: 'MAC地址',
              value: row.macAddress,
              rowNumber: row.rowNumber,
              firstRowNumber: seenMac.get(mac)
            })
            isDuplicate = true
          }
          if (isDuplicate) {
            return false
          }
          if (sn) seenSN.set(sn, row.rowNumber)
          if (mac) seenMac.set(mac, row.rowNumber)
          return true
        })

      if (parsedRows.length === 0) {
        alert('Excel 中未识别到有效数据，请确认第 3 行表头和下面的数据是否正确')
        excelPreviewList.value = []
        return
      }

      excelPreviewList.value = parsedRows
      if (duplicateMessages.length > 0) {
        alert(`已自动忽略 Excel 内重复的序列号或 MAC：${duplicateMessages.length} 处\n${formatDuplicateMessages(duplicateMessages)}`)
      }
    } catch (error) {
      console.error(error)
      alert('Excel 解析失败，请检查文件格式或第 3 行表头是否正确')
      excelPreviewList.value = []
    }
  }
  reader.readAsArrayBuffer(file)
}

async function saveBoardInboundRecords() {
  if (!uploadForm.file) {
    alert('请先上传板卡入库 Excel 文件')
    return
  }
  if (excelPreviewList.value.length === 0) {
    alert('当前 Excel 没有可导入的数据')
    return
  }

  const invalid = excelPreviewList.value.find(item => !item.sn || !item.macAddress)
  if (invalid) {
    alert('导入失败：序列号和 MAC地址不能为空')
    return
  }

  try {
    const res = await importBoardInboundRecords({
      sourceFileId: uploadForm.fileId,
      fileName: uploadForm.fileName,
      fileContentType: uploadForm.fileContentType,
      fileData: uploadForm.fileData,
      records: excelPreviewList.value
    })
    const result = getResponseData(res)
    if (result.code === 200) {
      alert(`导入成功，共导入 ${result.data?.count || excelPreviewList.value.length} 条记录`)
      showUploadDialog.value = false
      await loadBoardInboundRecords()
    } else {
      alert(result.msg || '导入失败')
    }
  } catch (err) {
    console.error('导入板卡入库失败：', err)
    alert(err.response?.data || '导入板卡入库失败')
  }
}

const filteredList = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase()
  return boardInboundList.value.filter(item => {
    const keywordMatch =
      !keyword ||
      [
        item.productName,
        item.productModel,
        item.productCode,
        item.sn,
        item.macAddress,
        item.pcbQrCode,
        item.remark
      ].some(value => String(value || '').toLowerCase().includes(keyword))

    const productModelMatch =
      !filters.productModel || item.productModel === filters.productModel

    return keywordMatch && productModelMatch
  })
})

const totalCount = computed(() => boardInboundList.value.length)

const currentMonthInboundCount = computed(() => {
  const currentMonth = new Date().toISOString().slice(0, 7)
  return boardInboundList.value.filter(item => item.inTime && item.inTime.startsWith(currentMonth)).length
})

const productModelOptions = computed(() => {
  const models = boardInboundList.value
    .map(item => item.productModel)
    .filter(Boolean)
    .filter(model => model !== '-')
  return [...new Set(models)]
})

const productModelSummary = computed(() => {
  return productModelOptions.value.map(model => ({
    productModel: model,
    count: boardInboundList.value.filter(item => item.productModel === model).length
  }))
})

const totalPage = computed(() => Math.max(1, Math.ceil(filteredList.value.length / pageSize.value)))
const paginatedList = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredList.value.slice(start, start + pageSize.value)
})
const pageStartIndex = computed(() => (filteredList.value.length === 0 ? 0 : (currentPage.value - 1) * pageSize.value + 1))
const pageEndIndex = computed(() => Math.min(currentPage.value * pageSize.value, filteredList.value.length))

watch(() => [filters.keyword, filters.productModel], () => {
  currentPage.value = 1
})
watch(pageSize, () => {
  currentPage.value = 1
})
watch(totalPage, value => {
  if (currentPage.value > value) currentPage.value = value
})

function resetFilters() {
  filters.keyword = ''
  filters.productModel = ''
  currentPage.value = 1
}

function downloadFile(item) {
  downloadLocalFile(
    {
      fileId: item.sourceFileId,
      fileName: item.sourceFileName,
      downloadUrl: getFileDownloadUrl(item.sourceFileId)
    },
    item.sourceFileName || '板卡入库记录.xlsx'
  )
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

.filter-card,
.table-card,
.preview-card,
.upload-tip {
  border: 1px solid rgba(148, 163, 184, 0.18);
  background: rgba(15, 23, 42, 0.72);
  box-shadow: 0 16px 32px rgba(2, 6, 23, 0.22);
}

.filter-card {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 18px;
  padding: 16px;
  border-radius: 8px;
}

.filter-card input,
.form-grid input,
.page-size-control select {
  height: 38px;
  border: 1px solid rgba(148, 163, 184, 0.28);
  border-radius: 6px;
  background: rgba(15, 23, 42, 0.92);
  color: #f8fafc;
  padding: 0 12px;
  outline: none;
}

.filter-card input {
  min-width: 420px;
}

.primary-btn,
.query-btn,
.reset-btn,
.page-btn {
  height: 38px;
  border: 0;
  border-radius: 6px;
  padding: 0 16px;
  color: #f8fafc;
  cursor: pointer;
}

.primary-btn { background: #2563eb; }
.query-btn { background: #1d4ed8; }
.reset-btn,
.page-btn { background: rgba(51, 65, 85, 0.92); }
.page-btn:disabled { opacity: 0.45; cursor: not-allowed; }

.table-card {
  border-radius: 8px;
  overflow: hidden;
}

.table-card-header,
.pagination-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  padding: 16px;
}

.table-card-header h3 {
  margin: 0 0 4px;
  font-size: 17px;
}

.table-card-header span,
.pagination-info,
.upload-tip p,
.preview-more {
  color: #94a3b8;
}

.page-size-control,
.pagination-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.table-wrapper,
.preview-wrapper {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  table-layout: fixed;
}

th,
td {
  border-top: 1px solid rgba(148, 163, 184, 0.14);
  padding: 12px;
  text-align: left;
  font-size: 13px;
  vertical-align: middle;
}

th {
  color: #cbd5e1;
  background: rgba(15, 23, 42, 0.88);
  font-weight: 700;
}

td {
  color: #e2e8f0;
}

.product-tag,
.model-text,
.code-text,
.sn-tag,
.mac-text,
.pcb-text,
.remark-text {
  display: block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.sn-tag {
  color: #bfdbfe;
  font-weight: 700;
}

.mac-text {
  color: #fde68a;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}

.pcb-text {
  color: #c4b5fd;
}

.muted {
  color: #94a3b8;
}

.nowrap {
  white-space: nowrap;
}

.inline-link {
  border: 0;
  background: transparent;
  color: #60a5fa;
  cursor: pointer;
  padding: 0;
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.empty-table {
  text-align: center;
  color: #94a3b8;
  padding: 30px;
}

.dialog-mask {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(2, 6, 23, 0.72);
  z-index: 1000;
}

.dialog {
  width: min(920px, calc(100vw - 32px));
  max-height: calc(100vh - 48px);
  overflow: auto;
  border: 1px solid rgba(148, 163, 184, 0.22);
  border-radius: 8px;
  background: #0f172a;
  box-shadow: 0 24px 60px rgba(2, 6, 23, 0.42);
}

.dialog-header,
.dialog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border-bottom: 1px solid rgba(148, 163, 184, 0.14);
}

.dialog-footer {
  border-top: 1px solid rgba(148, 163, 184, 0.14);
  border-bottom: 0;
  justify-content: flex-end;
}

.dialog-header h3 {
  margin: 0;
}

.dialog-header button {
  border: 0;
  background: transparent;
  color: #cbd5e1;
  font-size: 22px;
  cursor: pointer;
}

.upload-tip,
.preview-card,
.form-grid {
  margin: 16px;
  border-radius: 8px;
  padding: 14px;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
  padding: 0;
}

.form-grid label {
  display: flex;
  flex-direction: column;
  gap: 8px;
  color: #cbd5e1;
  font-size: 13px;
}

.form-grid input[type='file'] {
  display: none;
}

.file-picker-btn {
  width: fit-content;
  height: 34px;
  padding: 0 14px;
  border: 1px solid #3b82f6;
  border-radius: 8px;
  background: #1d4ed833;
  color: #bfdbfe;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
}

.file-picker-btn:hover {
  background: #2563eb55;
  color: #eff6ff;
}

.file-name {
  min-height: 18px;
  color: #94a3b8;
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.full-row {
  grid-column: 1 / -1;
}

.preview-title {
  margin-bottom: 10px;
  font-weight: 700;
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
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.type-item strong {
  color: #f8fafc;
  font-size: 22px;
}

.primary-btn,
.query-btn,
.reset-btn,
.page-btn {
  height: 36px;
  padding: 0 16px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
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

.filter-card {
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 14px;
  box-shadow: none;
  padding: 16px;
  display: grid;
  grid-template-columns: 1.4fr 220px 90px 90px;
  gap: 12px;
  margin-bottom: 20px;
}

.filter-card input,
.filter-card select,
.page-size-control select,
.form-grid input {
  border: 1px solid #334155;
  border-radius: 8px;
  background: #020617;
  color: #e2e8f0;
  padding: 0 12px;
  outline: none;
  height: 36px;
}

.filter-card input {
  min-width: 0;
}

.filter-card input::placeholder {
  color: #64748b;
}

.table-card,
.preview-card,
.upload-tip {
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 14px;
  box-shadow: none;
}

.table-card {
  overflow: hidden;
}

.table-card-header {
  padding: 14px 16px;
  border-bottom: 1px solid #1e293b;
}

.table-card-header h3 {
  font-size: 16px;
}

.table-card-header span,
.pagination-info {
  color: #64748b;
  font-size: 12px;
}

.page-size-control span {
  color: #94a3b8;
  font-size: 13px;
}

.page-size-control select {
  width: 90px;
}

.table-wrapper,
.preview-wrapper {
  width: 100%;
  overflow-x: auto;
  overflow-y: hidden;
  scrollbar-width: thin;
  scrollbar-color: #334155 #020617;
}

.table-wrapper::-webkit-scrollbar,
.preview-wrapper::-webkit-scrollbar {
  height: 10px;
}

.table-wrapper::-webkit-scrollbar-track,
.preview-wrapper::-webkit-scrollbar-track {
  background: #020617;
  border-radius: 999px;
}

.table-wrapper::-webkit-scrollbar-thumb,
.preview-wrapper::-webkit-scrollbar-thumb {
  background: #334155;
  border-radius: 999px;
  border: 2px solid #020617;
}

.version-table {
  min-width: 1320px;
  table-layout: fixed;
}

.version-table thead,
.preview-wrapper thead {
  background: #020617;
}

.version-table th,
.version-table td,
.preview-wrapper th,
.preview-wrapper td {
  box-sizing: border-box;
  white-space: nowrap;
  border-top: 0;
  border-bottom: 1px solid #1e293b;
}

.version-table th,
.preview-wrapper th {
  padding: 14px 16px;
  color: #94a3b8;
  font-size: 12px;
  font-weight: 600;
}

.version-table td,
.preview-wrapper td {
  padding: 15px 16px;
  color: #e2e8f0;
  overflow: hidden;
}

.product-tag,
.model-text,
.code-text,
.remark-text {
  color: #e2e8f0;
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
}

.mac-text {
  max-width: 160px;
  color: #cbd5e1;
  font-size: 12px;
  font-family: Consolas, Monaco, monospace;
}

.pcb-text {
  color: #cbd5e1;
}

.empty-table {
  color: #64748b !important;
  padding: 28px 16px !important;
}

.pagination-bar {
  padding: 14px 16px;
  border-top: 1px solid #1e293b;
}

.page-number {
  min-width: 70px;
  text-align: center;
  color: #cbd5e1;
  font-size: 13px;
}

.dialog-mask {
  z-index: 999;
  padding: 20px;
}

.dialog {
  width: 760px;
  max-width: 100%;
  max-height: calc(100vh - 40px);
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
}

.dialog-header h3 {
  font-size: 18px;
}

.dialog-header button {
  border: none;
  background: transparent;
  color: #94a3b8;
  font-size: 24px;
}

.upload-tip,
.preview-card,
.form-grid {
  margin: 20px;
  border-radius: 10px;
}

.upload-tip {
  padding: 14px;
}

.form-grid {
  padding: 0;
}

.dialog-footer {
  padding: 16px 20px;
  border-top: 1px solid #1e293b;
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
}
</style>
