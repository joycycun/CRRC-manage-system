<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>生产烧录记录管理</h1>
      </div>

      <button class="primary-btn" @click="openUploadDialog">
        上传烧录记录 Excel
      </button>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索生产批次号 / 产品名称 / 产品型号 / 产品编码 / 序列号 / MAC地址 / 硬件版本 / 软件版本 / PCB二维码 / 上传人 / 备注"
      />

      <select v-model="filters.batchNo">
        <option value="">全部生产批次</option>
        <option
          v-for="batch in batchNoOptions"
          :key="batch"
          :value="batch"
        >
          {{ batch }}
        </option>
      </select>

      <button class="query-btn">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <!-- 按生产批次展开 -->
    <div class="table-card">
      <div class="table-wrapper">
        <table class="batch-table">
          <thead>
            <tr>
              <th>生产批次号</th>
              <th>产品数量</th>
              <th>来源文件</th>
              <th>上传人</th>
              <th>上传时间</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>

          <tbody>
            <template
              v-for="batch in visibleBatchGroupList"
              :key="batch.batchNo"
            >
              <!-- 第一层：生产批次 -->
              <tr class="batch-row">
                <td>
                  <button
                    class="batch-name-btn"
                    @click="toggleBatch(batch.batchNo)"
                  >
                    <span class="expand-icon">
                      {{ expandedBatchNo === batch.batchNo ? '▼' : '▶' }}
                    </span>
                    {{ batch.batchNo }}
                  </button>
                </td>

                <td>
                  <span class="count-tag">
                    {{ batch.products.length }} 条记录
                  </span>
                </td>

                <td>
                  <span class="file-text" :title="batch.fileName">
                    {{ batch.fileName }}
                  </span>
                </td>

                <td>
                  <span class="normal-text" :title="batch.uploader">
                    {{ batch.uploader || '-' }}
                  </span>
                </td>

                <td class="muted nowrap">
                  {{ batch.uploadTime || '-' }}
                </td>

                <td class="operation-col">
                  <div class="action-group">
                    <button
                      class="text-btn blue"
                      @click="toggleBatch(batch.batchNo)"
                    >
                      {{ expandedBatchNo === batch.batchNo ? '收起' : '查看产品' }}
                    </button>

                    <button
                      class="text-btn red"
                      @click="deleteBatch(batch)"
                    >
                      删除批次
                    </button>
                  </div>
                </td>
              </tr>

              <!-- 第二层：批次下的产品明细 -->
              <tr
                v-if="expandedBatchNo === batch.batchNo"
                class="child-row"
              >
                <td colspan="6">
                  <div class="child-table-wrapper">
                    <table class="child-table">
                      <thead>
                        <tr>
                          <th>产品名称</th>
                          <th>产品型号</th>
                          <th>产品编码</th>
                          <th>序列号</th>
                          <th>MAC地址</th>
                          <th>硬件版本</th>
                          <th>软件版本</th>
                          <th>PCB二维码</th>
                          <th>备注</th>
                          <th>来源文件</th>
                          <th>上传人</th>
                          <th>上传时间</th>
                          <th class="child-operation-col">操作</th>
                        </tr>
                      </thead>

                      <tbody>
                        <tr
                          v-for="item in batch.products"
                          :key="item.id"
                        >
                          <td>
                            <span class="product-tag" :title="item.productName">
                              {{ item.productName }}
                            </span>
                          </td>

                          <td>
                            <span class="model-text" :title="item.productModel">
                              {{ item.productModel }}
                            </span>
                          </td>

                          <td>
                            <span class="code-text" :title="item.productCode">
                              {{ item.productCode }}
                            </span>
                          </td>

                          <td>
                            <span class="sn-tag" :title="item.serialNumber">
                              {{ item.serialNumber }}
                            </span>
                          </td>

                          <td>
                            <span class="mac-text" :title="item.macAddress">
                              {{ item.macAddress }}
                            </span>
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

                          <td>
                            <span class="pcb-text" :title="item.pcbQrCode">
                              {{ item.pcbQrCode }}
                            </span>
                          </td>

                          <td>
                            <span class="remark-text" :title="item.note">
                              {{ item.note }}
                            </span>
                          </td>

                          <td>
                            <span class="file-text" :title="item.fileName">
                              {{ item.fileName }}
                            </span>
                          </td>

                          <td>
                            <span class="normal-text" :title="item.uploader">
                              {{ item.uploader }}
                            </span>
                          </td>

                          <td class="muted nowrap">
                            {{ item.uploadTime }}
                          </td>

                          <td class="child-operation-col">
                            <div class="action-group">
                              <button class="text-btn" @click="viewBurnRecord(item)">
                                查看
                              </button>

                              <button class="text-btn blue" @click="downloadBurnRecord(item)">
                                下载
                              </button>

                              <button class="text-btn red" @click="deleteBurnRecord(item)">
                                删除
                              </button>
                            </div>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>

      <div class="table-footer">
        共 {{ filteredBurnRecordList.length }} 条烧录记录，按 {{ filteredBatchGroupList.length }} 个生产批次归类展示。
        展开某一批次后，其余批次会自动隐藏。
      </div>
    </div>

    <!-- 上传烧录记录弹窗 -->
    <div v-if="showUploadDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>上传烧录记录 Excel</h3>
          <button @click="showUploadDialog = false">×</button>
        </div>

        <div class="upload-tip">
          <strong>Excel 读取规则：</strong>
          <p>
            系统会从 Excel 第 3 行开始读取表头。请保证第 3 行包含：
            产品名称、产品型号、产品编码、序列号、MAC地址、硬件版本、软件版本、PCB二维码、备注。
          </p>
          <p>
            如果 Excel 没有“生产批次号”这一列，系统会使用下面填写的生产批次号；
            如果这里也不填写，则从文件名中自动提取 8 位数字，例如 20260402。
          </p>
        </div>

        <div class="form-grid">
          <label>
            生产批次号
            <input
              v-model="uploadForm.batchNo"
              placeholder="例如：20260402"
            />
          </label>

          <label>
            当前上传人
            <input
              :value="currentUserName"
              disabled
            />
          </label>

          <label class="full-row">
            文件名称
            <input
              v-model="uploadForm.fileName"
              placeholder="选择 Excel 后自动填充"
              disabled
            />
          </label>

          <label class="full-row">
            烧录记录 Excel
            <input
              type="file"
              accept=".xls,.xlsx"
              @change="handleExcelFileChange"
            />
          </label>

          <label class="full-row">
            烧录说明
            <textarea
              v-model="uploadForm.remark"
              placeholder="例如：本次上传为某批次产品烧录记录，系统会自动读取 Excel 第 3 行表头后的所有产品信息"
            ></textarea>
          </label>
        </div>

        <div v-if="excelPreviewList.length > 0" class="preview-card">
          <div class="preview-title">
            已识别 {{ previewRecordCount }} 条产品记录
          </div>

          <div class="preview-wrapper">
            <table>
              <thead>
                <tr>
                  <th>生产批次号</th>
                  <th>产品名称</th>
                  <th>产品型号</th>
                  <th>产品编码</th>
                  <th>序列号</th>
                  <th>MAC地址</th>
                  <th>硬件版本</th>
                  <th>软件版本</th>
                  <th>PCB二维码</th>
                  <th>备注</th>
                </tr>
              </thead>

              <tbody>
                <tr
                  v-for="(item, index) in excelPreviewList.slice(0, 8)"
                  :key="index"
                >
                  <td>{{ getFinalBatchNo(item.batchNo) }}</td>
                  <td>{{ item.productName }}</td>
                  <td>{{ item.productModel }}</td>
                  <td>{{ item.productCode }}</td>
                  <td>
                    {{ item.serialNumbers.length > 0 ? item.serialNumbers.join('、') : '-' }}
                  </td>
                  <td>{{ item.macAddress }}</td>
                  <td>{{ item.hardwareVersion }}</td>
                  <td>{{ item.softwareVersion }}</td>
                  <td>{{ item.pcbQrCode }}</td>
                  <td>{{ item.note }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <p v-if="excelPreviewList.length > 8" class="preview-more">
            仅预览前 8 行，保存后会导入全部记录。
          </p>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showUploadDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="saveExcelBurnRecords">
            保存导入
          </button>
        </div>
      </div>
    </div>

    <!-- 查看烧录记录详情弹窗 -->
    <div v-if="selectedBurnRecord" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>烧录记录详情</h3>
          <button @click="selectedBurnRecord = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>生产批次号</span>
            <strong>{{ selectedBurnRecord.batchNo }}</strong>
          </div>

          <div>
            <span>产品名称</span>
            <strong>{{ selectedBurnRecord.productName }}</strong>
          </div>

          <div>
            <span>产品型号</span>
            <strong>{{ selectedBurnRecord.productModel }}</strong>
          </div>

          <div>
            <span>产品编码</span>
            <strong>{{ selectedBurnRecord.productCode }}</strong>
          </div>

          <div>
            <span>序列号</span>
            <strong>{{ selectedBurnRecord.serialNumber }}</strong>
          </div>

          <div>
            <span>MAC地址</span>
            <strong>{{ selectedBurnRecord.macAddress }}</strong>
          </div>

          <div>
            <span>硬件版本</span>
            <strong>{{ selectedBurnRecord.hardwareVersion }}</strong>
          </div>

          <div>
            <span>软件版本</span>
            <strong>{{ selectedBurnRecord.softwareVersion }}</strong>
          </div>

          <div>
            <span>PCB二维码</span>
            <strong>{{ selectedBurnRecord.pcbQrCode }}</strong>
          </div>

          <div>
            <span>备注</span>
            <strong>{{ selectedBurnRecord.note }}</strong>
          </div>

          <div>
            <span>来源文件</span>
            <strong>{{ selectedBurnRecord.fileName }}</strong>
          </div>

          <div>
            <span>上传人</span>
            <strong>{{ selectedBurnRecord.uploader }}</strong>
          </div>

          <div>
            <span>上传时间</span>
            <strong>{{ selectedBurnRecord.uploadTime }}</strong>
          </div>

          <div>
            <span>文件查看</span>
            <button class="inline-link" @click="openBurnRecordFile(selectedBurnRecord)">
              点开查看文件
            </button>
          </div>
        </div>

        <div class="remark-card">
          <span>烧录说明</span>
          <p>{{ selectedBurnRecord.importRemark || '暂无说明' }}</p>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="downloadBurnRecord(selectedBurnRecord)">
            下载源文件
          </button>

          <button class="primary-btn" @click="selectedBurnRecord = null">
            关闭
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import * as XLSX from 'xlsx'

const currentUserName = ref(
  localStorage.getItem('username') ||
  localStorage.getItem('accountName') ||
  localStorage.getItem('realName') ||
  '当前用户'
)

const filters = reactive({
  keyword: '',
  batchNo: ''
})

const showUploadDialog = ref(false)
const selectedBurnRecord = ref(null)
const excelPreviewList = ref([])
const expandedBatchNo = ref('')

const uploadForm = reactive({
  batchNo: '',
  fileName: '',
  file: null,
  fileUrl: '',
  remark: ''
})

const burnRecordList = ref([
  {
    id: 1,
    batchNo: '20260402',
    productName: '控制盒',
    productModel: 'DACU-SIP-Porto',
    productCode: 'M00004665176',
    serialNumber: 'ZDAPO2300001',
    macAddress: '68:69:2E:DD:17:03',
    hardwareVersion: '-',
    softwareVersion: '-',
    pcbQrCode: '-',
    note: '波尔图',
    fileName: '生产烧录记录_20260402.xlsx',
    fileUrl: '',
    uploader: '生产人员',
    uploadTime: '2026-05-10',
    importRemark: '模拟数据，按 Excel 第 3 行表头结构生成。'
  },
  {
    id: 2,
    batchNo: '20260402',
    productName: '乘客报警器',
    productModel: 'PECU-BOG-Metro-PCBA',
    productCode: 'M00004665176',
    serialNumber: 'ZPEBOG230000X',
    macAddress: '-',
    hardwareVersion: '-',
    softwareVersion: '-',
    pcbQrCode: '-',
    note: '模板',
    fileName: '生产烧录记录_20260402.xlsx',
    fileUrl: '',
    uploader: '生产人员',
    uploadTime: '2026-05-10',
    importRemark: '模拟数据，按 Excel 第 3 行表头结构生成。'
  },
  {
    id: 3,
    batchNo: '20260402',
    productName: '编码板',
    productModel: 'ACSU-BOG-Encode-PCBA',
    productCode: 'M00004665172',
    serialNumber: 'ZENBOG230000X',
    macAddress: '-',
    hardwareVersion: '-',
    softwareVersion: '-',
    pcbQrCode: '-',
    note: '模板',
    fileName: '生产烧录记录_20260402.xlsx',
    fileUrl: '',
    uploader: '生产人员',
    uploadTime: '2026-05-10',
    importRemark: '模拟数据，按 Excel 第 3 行表头结构生成。'
  }
])

const batchNoOptions = computed(() => {
  const batches = burnRecordList.value
    .map(item => item.batchNo)
    .filter(Boolean)

  return [...new Set(batches)]
})

const filteredBurnRecordList = computed(() => {
  return burnRecordList.value.filter(item => {
    const keyword = filters.keyword.trim()

    const keywordMatch =
      !keyword ||
      item.batchNo.includes(keyword) ||
      item.productName.includes(keyword) ||
      item.productModel.includes(keyword) ||
      item.productCode.includes(keyword) ||
      item.serialNumber.includes(keyword) ||
      item.macAddress.includes(keyword) ||
      item.hardwareVersion.includes(keyword) ||
      item.softwareVersion.includes(keyword) ||
      item.pcbQrCode.includes(keyword) ||
      item.note.includes(keyword) ||
      item.uploader.includes(keyword) ||
      item.fileName.includes(keyword)

    const batchMatch =
      !filters.batchNo || item.batchNo === filters.batchNo

    return keywordMatch && batchMatch
  })
})

const filteredBatchGroupList = computed(() => {
  const map = new Map()

  filteredBurnRecordList.value.forEach(item => {
    const batchNo = item.batchNo || '未填写批次'

    if (!map.has(batchNo)) {
      map.set(batchNo, {
        batchNo,
        products: []
      })
    }

    map.get(batchNo).products.push(item)
  })

  return Array.from(map.values()).map(batch => {
    const sortedProducts = [...batch.products].sort((a, b) => {
      return String(a.serialNumber).localeCompare(String(b.serialNumber))
    })

    const fileNames = [
      ...new Set(sortedProducts.map(item => item.fileName).filter(Boolean))
    ]

    const uploaders = [
      ...new Set(sortedProducts.map(item => item.uploader).filter(Boolean))
    ]

    const uploadTimes = [
      ...new Set(sortedProducts.map(item => item.uploadTime).filter(Boolean))
    ]

    return {
      batchNo: batch.batchNo,
      products: sortedProducts,
      fileName: fileNames.length > 1 ? `${fileNames.length} 个来源文件` : fileNames[0] || '-',
      uploader: uploaders.length > 1 ? `${uploaders.length} 个上传人` : uploaders[0] || '-',
      uploadTime: uploadTimes.length > 1 ? `${uploadTimes.length} 个上传时间` : uploadTimes[0] || '-'
    }
  })
})

const visibleBatchGroupList = computed(() => {
  if (!expandedBatchNo.value) {
    return filteredBatchGroupList.value
  }

  return filteredBatchGroupList.value.filter(
    batch => batch.batchNo === expandedBatchNo.value
  )
})

const previewRecordCount = computed(() => {
  return excelPreviewList.value.reduce((count, item) => {
    const snCount = item.serialNumbers.length > 0 ? item.serialNumbers.length : 1
    return count + snCount
  }, 0)
})

function resetFilters() {
  filters.keyword = ''
  filters.batchNo = ''
  expandedBatchNo.value = ''
}

function toggleBatch(batchNo) {
  if (expandedBatchNo.value === batchNo) {
    expandedBatchNo.value = ''
  } else {
    expandedBatchNo.value = batchNo
  }
}

function openUploadDialog() {
  uploadForm.batchNo = ''
  uploadForm.fileName = ''
  uploadForm.file = null
  uploadForm.fileUrl = ''
  uploadForm.remark = ''
  excelPreviewList.value = []
  showUploadDialog.value = true
}

function normalizeHeader(value) {
  return String(value || '')
    .trim()
    .replace(/\s+/g, '')
    .replace(/：/g, ':')
}

function getCellValue(row, headerNames) {
  for (const name of headerNames) {
    if (row[name] !== undefined && row[name] !== null && row[name] !== '') {
      return String(row[name]).trim()
    }
  }

  return ''
}

function splitSerialNumbers(value) {
  return String(value || '')
    .split(/[\n\r,，、;；\s]+/)
    .map(item => item.trim())
    .filter(Boolean)
}

function extractBatchNoFromFileName(fileName) {
  const match = String(fileName || '').match(/\d{8}/)
  return match ? match[0] : ''
}

function getFinalBatchNo(excelBatchNo) {
  return (
    excelBatchNo ||
    uploadForm.batchNo ||
    extractBatchNoFromFileName(uploadForm.fileName) ||
    '未填写批次'
  )
}

function normalizeExcelRow(row) {
  const normalizedRow = {}

  Object.keys(row).forEach(key => {
    normalizedRow[normalizeHeader(key)] = row[key]
  })

  const serialNumberText = getCellValue(normalizedRow, [
    '序列号',
    'SN',
    'SN序列号',
    'SerialNumber'
  ])

  return {
    batchNo: getCellValue(normalizedRow, [
      '生产批次号',
      '批次号',
      '生产批次',
      '批次',
      'BatchNo',
      'batchNo'
    ]),
    productName: getCellValue(normalizedRow, [
      '产品名称',
      '产品名',
      '品名'
    ]),
    productModel: getCellValue(normalizedRow, [
      '产品型号',
      '型号',
      '产品类型',
      '终端类型'
    ]),
    productCode: getCellValue(normalizedRow, [
      '产品编码',
      '产品编号',
      '编码'
    ]),
    serialNumberText,
    serialNumbers: splitSerialNumbers(serialNumberText),
    macAddress: getCellValue(normalizedRow, [
      'MAC地址',
      'MAC',
      'Mac地址',
      'mac地址'
    ]),
    hardwareVersion: getCellValue(normalizedRow, [
      '硬件版本',
      '硬件版本号',
      'HardwareVersion'
    ]),
    softwareVersion: getCellValue(normalizedRow, [
      '软件版本',
      '软件版本号',
      'SoftwareVersion'
    ]),
    pcbQrCode: getCellValue(normalizedRow, [
      'PCB二维码',
      'PCB码',
      'PCBQRCode',
      'PCBQR'
    ]),
    note: getCellValue(normalizedRow, [
      '备注',
      '说明',
      'Remark',
      'Note'
    ])
  }
}

function handleExcelFileChange(event) {
  const file = event.target.files[0]
  if (!file) return

  const isExcel =
    file.name.endsWith('.xls') ||
    file.name.endsWith('.xlsx')

  if (!isExcel) {
    alert('请上传 Excel 文件，格式为 .xls 或 .xlsx')
    event.target.value = ''
    return
  }

  uploadForm.file = file
  uploadForm.fileName = file.name
  uploadForm.fileUrl = URL.createObjectURL(file)

  const extractedBatchNo = extractBatchNoFromFileName(file.name)
  if (!uploadForm.batchNo && extractedBatchNo) {
    uploadForm.batchNo = extractedBatchNo
  }

  const reader = new FileReader()

  reader.onload = e => {
    try {
      const data = new Uint8Array(e.target.result)
      const workbook = XLSX.read(data, { type: 'array' })
      const firstSheetName = workbook.SheetNames[0]
      const worksheet = workbook.Sheets[firstSheetName]

      const rows = XLSX.utils.sheet_to_json(worksheet, {
        defval: '',
        range: 2
      })

      const parsedRows = rows
        .map(row => normalizeExcelRow(row))
        .filter(row => {
          return (
            row.batchNo ||
            row.productName ||
            row.productModel ||
            row.productCode ||
            row.serialNumbers.length > 0 ||
            row.macAddress ||
            row.hardwareVersion ||
            row.softwareVersion ||
            row.pcbQrCode ||
            row.note
          )
        })

      if (parsedRows.length === 0) {
        alert('Excel 中未识别到有效数据，请确认第 3 行是否包含：产品名称、产品型号、产品编码、序列号、MAC地址、硬件版本、软件版本、PCB二维码、备注')
        excelPreviewList.value = []
        return
      }

      excelPreviewList.value = parsedRows
    } catch (error) {
      console.error(error)
      alert('Excel 解析失败，请检查文件格式或第 3 行表头是否正确')
      excelPreviewList.value = []
    }
  }

  reader.readAsArrayBuffer(file)
}

function saveExcelBurnRecords() {
  if (!uploadForm.file) {
    alert('请先上传 Excel 烧录记录文件')
    return
  }

  if (excelPreviewList.value.length === 0) {
    alert('当前 Excel 没有可导入的数据')
    return
  }

  const now = new Date().toISOString().slice(0, 10)
  const uploader = currentUserName.value

  const records = excelPreviewList.value.flatMap((item, rowIndex) => {
    const serialNumbers =
      item.serialNumbers.length > 0
        ? item.serialNumbers
        : ['-']

    return serialNumbers.map((serialNumber, snIndex) => {
      return {
        id: Date.now() + rowIndex * 1000 + snIndex,
        batchNo: getFinalBatchNo(item.batchNo),
        productName: item.productName || '-',
        productModel: item.productModel || '-',
        productCode: item.productCode || '-',
        serialNumber,
        macAddress: item.macAddress || '-',
        hardwareVersion: item.hardwareVersion || '-',
        softwareVersion: item.softwareVersion || '-',
        pcbQrCode: item.pcbQrCode || '-',
        note: item.note || '-',
        fileName: uploadForm.fileName,
        fileUrl: uploadForm.fileUrl,
        uploader,
        uploadTime: now,
        importRemark: uploadForm.remark
      }
    })
  })

  burnRecordList.value.unshift(...records)

  const firstBatchNo = records[0]?.batchNo
  if (firstBatchNo) {
    expandedBatchNo.value = firstBatchNo
  }

  showUploadDialog.value = false
}

function viewBurnRecord(item) {
  selectedBurnRecord.value = item
}

function openBurnRecordFile(item) {
  if (!item.fileUrl) {
    alert('当前是模拟数据，暂无可直接打开的原始文件')
    return
  }

  window.open(item.fileUrl, '_blank')
}

function downloadBurnRecord(item) {
  if (!item.fileUrl) {
    alert('当前是模拟数据，暂无可下载的原始文件')
    return
  }

  const link = document.createElement('a')
  link.href = item.fileUrl
  link.download = item.fileName || '烧录记录文件.xlsx'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

function deleteBurnRecord(item) {
  const ok = confirm(`确认删除序列号【${item.serialNumber}】的烧录记录吗？`)
  if (!ok) return

  burnRecordList.value = burnRecordList.value.filter(
    record => record.id !== item.id
  )
}

function deleteBatch(batch) {
  const ok = confirm(`确认删除生产批次【${batch.batchNo}】下的全部 ${batch.products.length} 条烧录记录吗？`)
  if (!ok) return

  const ids = batch.products.map(item => item.id)

  burnRecordList.value = burnRecordList.value.filter(
    record => !ids.includes(record.id)
  )

  if (expandedBatchNo.value === batch.batchNo) {
    expandedBatchNo.value = ''
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
  grid-template-columns: 1.4fr 180px 90px 90px;
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

.form-grid textarea {
  min-height: 90px;
  padding: 10px 12px;
  resize: vertical;
}

.filter-card input::placeholder,
.form-grid input::placeholder,
.form-grid textarea::placeholder {
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
.child-table-wrapper,
.preview-wrapper {
  width: 100%;
  overflow-x: auto;
  overflow-y: hidden;
}

.table-wrapper::-webkit-scrollbar,
.child-table-wrapper::-webkit-scrollbar,
.preview-wrapper::-webkit-scrollbar {
  height: 10px;
}

.table-wrapper::-webkit-scrollbar-track,
.child-table-wrapper::-webkit-scrollbar-track,
.preview-wrapper::-webkit-scrollbar-track {
  background: #020617;
  border-radius: 999px;
}

.table-wrapper::-webkit-scrollbar-thumb,
.child-table-wrapper::-webkit-scrollbar-thumb,
.preview-wrapper::-webkit-scrollbar-thumb {
  background: #334155;
  border-radius: 999px;
  border: 2px solid #020617;
}

.table-wrapper::-webkit-scrollbar-thumb:hover,
.child-table-wrapper::-webkit-scrollbar-thumb:hover,
.preview-wrapper::-webkit-scrollbar-thumb:hover {
  background: #475569;
}

.table-wrapper::-webkit-scrollbar-button,
.child-table-wrapper::-webkit-scrollbar-button,
.preview-wrapper::-webkit-scrollbar-button {
  display: none;
}

.table-wrapper,
.child-table-wrapper,
.preview-wrapper {
  scrollbar-width: thin;
  scrollbar-color: #334155 #020617;
}

.batch-table {
  width: 100%;
  min-width: 1050px;
  border-collapse: collapse;
  table-layout: fixed;
}

.batch-table thead {
  background: #020617;
}

.batch-table th,
.batch-table td {
  box-sizing: border-box;
  white-space: nowrap;
}

.batch-table th {
  padding: 14px 16px;
  color: #94a3b8;
  font-size: 12px;
  font-weight: 600;
  text-align: left;
  border-bottom: 1px solid #1e293b;
}

.batch-table td {
  padding: 15px 16px;
  font-size: 13px;
  color: #e2e8f0;
  border-bottom: 1px solid #1e293b;
  vertical-align: middle;
  overflow: hidden;
}

.batch-table th:nth-child(1),
.batch-table td:nth-child(1) {
  width: 180px;
}

.batch-table th:nth-child(2),
.batch-table td:nth-child(2) {
  width: 120px;
}

.batch-table th:nth-child(3),
.batch-table td:nth-child(3) {
  width: 260px;
}

.batch-table th:nth-child(4),
.batch-table td:nth-child(4) {
  width: 120px;
}

.batch-table th:nth-child(5),
.batch-table td:nth-child(5) {
  width: 140px;
}

.batch-table th:nth-child(6),
.batch-table td:nth-child(6) {
  width: 230px;
}

.batch-row {
  background: #0f172a;
}

.batch-row:hover {
  background: #1e293b80;
}

.batch-name-btn {
  display: inline-block;
  max-width: 150px;
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

.batch-name-btn:hover {
  color: #93c5fd;
  text-decoration: underline;
}

.expand-icon {
  display: inline-block;
  width: 18px;
  color: #94a3b8;
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

.child-row td {
  padding: 0;
  background: #020617;
}

.child-table-wrapper {
  padding: 12px 16px 16px;
  box-sizing: border-box;
}

.child-table {
  width: 100%;
  min-width: 2050px;
  border-collapse: collapse;
  table-layout: fixed;
  border: 1px solid #1e293b;
  border-radius: 10px;
  overflow: hidden;
}

.child-table thead {
  background: #0f172a;
}

.child-table th,
.child-table td {
  box-sizing: border-box;
  white-space: nowrap;
}

.child-table th {
  padding: 12px 14px;
  color: #94a3b8;
  font-size: 12px;
  font-weight: 600;
  border-bottom: 1px solid #1e293b;
}

.child-table td {
  padding: 13px 14px;
  border-bottom: 1px solid #1e293b;
  color: #e2e8f0;
  font-size: 13px;
  vertical-align: middle;
  overflow: hidden;
}

.child-table tbody tr:hover {
  background: #1e293b80;
}

.child-table th:nth-child(1),
.child-table td:nth-child(1) {
  width: 130px;
}

.child-table th:nth-child(2),
.child-table td:nth-child(2) {
  width: 190px;
}

.child-table th:nth-child(3),
.child-table td:nth-child(3) {
  width: 160px;
}

.child-table th:nth-child(4),
.child-table td:nth-child(4) {
  width: 200px;
}

.child-table th:nth-child(5),
.child-table td:nth-child(5) {
  width: 170px;
}

.child-table th:nth-child(6),
.child-table td:nth-child(6) {
  width: 210px;
}

.child-table th:nth-child(7),
.child-table td:nth-child(7) {
  width: 210px;
}

.child-table th:nth-child(8),
.child-table td:nth-child(8) {
  width: 180px;
}

.child-table th:nth-child(9),
.child-table td:nth-child(9) {
  width: 140px;
}

.child-table th:nth-child(10),
.child-table td:nth-child(10) {
  width: 220px;
}

.child-table th:nth-child(11),
.child-table td:nth-child(11) {
  width: 110px;
}

.child-table th:nth-child(12),
.child-table td:nth-child(12) {
  width: 130px;
}

.child-table th:nth-child(13),
.child-table td:nth-child(13) {
  width: 210px;
}

.product-tag,
.software-tag,
.hardware-tag {
  display: inline-block;
  max-width: 190px;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  vertical-align: middle;
  box-sizing: border-box;
}

.product-tag {
  background: #0f766e33;
  color: #5eead4;
}

.software-tag {
  background: #16a34a33;
  color: #4ade80;
}

.hardware-tag {
  background: #9333ea33;
  color: #c084fc;
}

.model-text,
.code-text,
.sn-tag,
.mac-text,
.file-text,
.normal-text,
.pcb-text,
.remark-text {
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

.normal-text {
  font-family: inherit;
  font-size: 13px;
}

.sn-tag {
  padding: 3px 8px;
  border-radius: 999px;
  background: #33415566;
}

.version-cell {
  overflow: hidden;
}

.nowrap {
  white-space: nowrap !important;
}

.muted {
  color: #94a3b8 !important;
}

.operation-col,
.child-operation-col {
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

.text-btn.red {
  color: #f87171;
}

.table-footer {
  padding: 12px 16px;
  color: #64748b;
  font-size: 12px;
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
  width: 860px;
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
  width: 960px;
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

.preview-card {
  margin: 0 20px 20px;
  padding: 14px;
  background: #020617;
  border: 1px solid #1e293b;
  border-radius: 12px;
}

.preview-title {
  color: #f8fafc;
  font-size: 14px;
  font-weight: 700;
  margin-bottom: 12px;
}

.preview-wrapper {
  width: 100%;
  overflow-x: auto;
}

.preview-wrapper table {
  width: 100%;
  min-width: 1400px;
  border-collapse: collapse;
}

.preview-wrapper th,
.preview-wrapper td {
  padding: 10px 12px;
  border-bottom: 1px solid #1e293b;
  color: #cbd5e1;
  font-size: 12px;
  text-align: left;
  white-space: nowrap;
}

.preview-wrapper th {
  color: #94a3b8;
  background: #0f172a;
}

.preview-more {
  margin: 10px 0 0;
  color: #64748b;
  font-size: 12px;
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

  .batch-table {
    min-width: 1050px;
  }

  .child-table {
    min-width: 2050px;
  }

  .preview-wrapper table {
    min-width: 1400px;
  }

  .form-grid,
  .detail-card {
    grid-template-columns: 1fr;
  }
}
</style>