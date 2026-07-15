<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1>板卡组成</h1>
      </div>
      <button v-if="canManage" class="primary-btn" @click="openEditDialog">
        新增板卡组成
      </button>
    </div>

    <div class="filter-card">
      <input v-model="filters.keyword" placeholder="搜索项目名称 / 产品名称 / 入库型号 / 出库型号" />
      <select v-model="filters.projectName">
        <option value="">全部项目</option>
        <option v-for="project in projectOptions" :key="project" :value="project">
          {{ project }}
        </option>
      </select>
      <button class="query-btn" @click="loadCompositions">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <div class="table-card">
      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>项目名称</th>
              <th>最近更新时间</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>
          <tbody>
            <template v-for="group in groupedCompositions" :key="group.projectName">
              <tr class="project-row">
                <td>
                  <button class="project-name-btn" @click="toggleProject(group.projectName)">
                    <span class="expand-icon">{{ expandedProjects.includes(group.projectName) ? '▼' : '▶' }}</span>
                    {{ group.projectName }}
                  </button>
                </td>
                <td class="muted">{{ group.latestTime || '-' }}</td>
                <td class="operation-col">
                  <button class="text-btn blue" @click="toggleProject(group.projectName)">
                    {{ expandedProjects.includes(group.projectName) ? '收起' : '查看组成' }}
                  </button>
                </td>
              </tr>
              <tr v-if="expandedProjects.includes(group.projectName)" class="child-panel-row">
                <td colspan="3">
                  <div class="composition-panel">
                    <div class="composition-header">
                      <span>产品名称</span>
                      <span>出库型号</span>
                      <span>入库型号</span>
                    </div>
                    <div
                      v-for="composition in group.compositions"
                      :key="composition.key"
                      class="composition-row"
                    >
                      <div>
                        <span class="product-tag">{{ composition.productName }}</span>
                      </div>
                      <div>
                        <span class="model-tag outbound">{{ composition.outboundModel }}</span>
                      </div>
                      <div class="inbound-model-list">
                        <span
                          v-for="item in composition.items"
                          :key="item.id"
                          class="model-tag inbound"
                        >
                          {{ item.inboundModel }}
                          <button v-if="canManage" class="chip-delete" @click="deleteComposition(item)">×</button>
                        </span>
                      </div>
                    </div>
                  </div>
                </td>
              </tr>
            </template>
            <tr v-if="groupedCompositions.length === 0">
              <td colspan="3" class="empty-table">暂无板卡组成数据</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="table-footer">共 {{ filteredCompositions.length }} 条板卡组成，按项目归类展示</div>
    </div>

    <div v-if="showEditDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>新增板卡组成</h3>
          <button @click="showEditDialog = false">×</button>
        </div>

        <div class="upload-tip">
          <strong>Excel 读取规则：</strong>
          <p>系统会从 Excel 第 1 行读取表头。请保证第 1 行包含：项目名称、产品名称、入库型号、出库型号。合并单元格下方空白内容会自动沿用上一行。</p>
        </div>

        <div class="form-grid">
          <label class="full-row">
            板卡组成 Excel
            <input
              ref="excelFileInput"
              type="file"
              accept=".xls,.xlsx"
              @change="handleExcelChange"
            />
            <button type="button" class="file-picker-btn" @click="excelFileInput?.click()">
              选择 Excel 文件
            </button>
            <span class="file-name">
              {{ excelFileName || '请选择 .xls / .xlsx 文件' }}
            </span>
          </label>
        </div>

        <div class="manual-card">
          <div class="manual-card-title">
            <strong>手动新增</strong>
            <span>产品名称和入库型号来自板卡入库，可直接选择或手动填写</span>
          </div>
          <div class="manual-grid">
            <label>
              项目名称
              <select v-model="manualForm.projectName">
                <option value="">请选择项目</option>
                <option v-for="project in projectOptions" :key="project" :value="project">{{ project }}</option>
              </select>
            </label>
            <label>
              产品名称
              <input v-model="manualForm.productName" list="boardProductNames" placeholder="产品名称" />
              <datalist id="boardProductNames">
                <option v-for="name in boardProductNames" :key="name" :value="name" />
              </datalist>
            </label>
            <label>
              入库型号
              <input v-model="manualForm.inboundModel" list="boardInboundModels" placeholder="入库型号" />
              <datalist id="boardInboundModels">
                <option v-for="model in boardInboundModels" :key="model" :value="model" />
              </datalist>
            </label>
            <label>
              出库型号
              <input v-model="manualForm.outboundModel" placeholder="生产烧录中的产品型号" />
            </label>
          </div>
          <div class="manual-actions">
            <button class="reset-btn" type="button" @click="resetManualForm">清空</button>
            <button class="query-btn" type="button" @click="addManualComposition">加入列表</button>
          </div>
        </div>

        <div v-if="previewList.length > 0" class="preview-card">
          <div class="preview-title">待保存 {{ previewList.length }} 条板卡组成</div>
          <div class="preview-wrapper">
            <table>
              <thead>
                <tr>
                  <th>项目名称</th>
                  <th>产品名称</th>
                  <th>入库型号</th>
                  <th>出库型号</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(item, index) in previewList" :key="index">
                  <td>{{ item.projectName }}</td>
                  <td>{{ item.productName }}</td>
                  <td>{{ item.inboundModel }}</td>
                  <td>{{ item.outboundModel }}</td>
                  <td><button class="text-btn red" @click="previewList.splice(index, 1)">移除</button></td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showEditDialog = false">取消</button>
          <button class="primary-btn" @click="saveCompositions">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import * as XLSX from 'xlsx'
import { getProjects } from '@/api/project'
import { getBoardInboundRecords } from '@/api/boardInbound'
import { createBoardCompositions, deleteBoardComposition, getBoardCompositions } from '@/api/boardComposition'
import { canUseAction } from '@/utils/permission'

const filters = reactive({ keyword: '', projectName: '' })
const compositionList = ref([])
const projectOptions = ref([])
const boardInboundList = ref([])
const expandedProjects = ref([])
const showEditDialog = ref(false)
const previewList = ref([])
const excelFileInput = ref(null)
const excelFileName = ref('')
const manualForm = reactive({
  projectName: '',
  productName: '',
  inboundModel: '',
  outboundModel: ''
})

const canManage = computed(() => canUseAction('board-composition:create') || canUseAction('board-composition:upload'))

const boardProductNames = computed(() => [...new Set(boardInboundList.value.map(item => item.productName).filter(Boolean))])
const boardInboundModels = computed(() => [...new Set(boardInboundList.value.map(item => item.productModel).filter(Boolean))])

const filteredCompositions = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase()
  return compositionList.value.filter(item => {
    const keywordMatch = !keyword || [
      item.projectName,
      item.productName,
      item.inboundModel,
      item.outboundModel
    ].some(value => String(value || '').toLowerCase().includes(keyword))
    const projectMatch = !filters.projectName || item.projectName === filters.projectName
    return keywordMatch && projectMatch
  })
})

const groupedCompositions = computed(() => {
  const map = new Map()
  filteredCompositions.value.forEach(item => {
    const projectName = item.projectName || '未绑定项目'
    if (!map.has(projectName)) {
      map.set(projectName, [])
    }
    map.get(projectName).push(item)
  })
  return Array.from(map.entries()).map(([projectName, items]) => {
    const compositionMap = new Map()
    items.forEach(item => {
      const key = `${item.productName}__${item.outboundModel}`
      if (!compositionMap.has(key)) {
        compositionMap.set(key, [])
      }
      compositionMap.get(key).push(item)
    })

    const compositions = Array.from(compositionMap.entries()).map(([key, groupItems]) => ({
      key,
      productName: groupItems[0]?.productName || '-',
      outboundModel: groupItems[0]?.outboundModel || '-',
      items: groupItems,
      latestTime: groupItems.map(item => item.updatedAt || item.createdAt).filter(Boolean).sort().reverse()[0] || ''
    }))

    return {
      projectName,
      items,
      compositions,
      latestTime: items.map(item => item.updatedAt || item.createdAt).filter(Boolean).sort().reverse()[0] || ''
    }
  })
})

onMounted(async () => {
  await Promise.all([loadProjects(), loadBoardInboundOptions(), loadCompositions()])
})

function getResponseData(res) {
  return res?.data || res
}

async function loadProjects() {
  const res = await getProjects()
  const result = getResponseData(res)
  projectOptions.value = (result.data || []).map(item => item.projectName || item.project_name).filter(Boolean)
}

async function loadBoardInboundOptions() {
  const res = await getBoardInboundRecords()
  const result = getResponseData(res)
  boardInboundList.value = (result.data || []).map(item => ({
    productName: item.productName || item.product_name || '',
    productModel: item.productModel || item.product_model || ''
  }))
}

async function loadCompositions() {
  const res = await getBoardCompositions()
  const result = getResponseData(res)
  compositionList.value = (result.data || []).map(item => ({
    id: item.id,
    projectId: item.projectId || item.project_id || 0,
    projectName: item.projectName || item.project_name || '',
    productName: item.productName || item.product_name || '',
    inboundModel: item.inboundModel || item.inbound_model || '',
    outboundModel: item.outboundModel || item.outbound_model || '',
    createdBy: item.createdBy || item.created_by || '',
    createdAt: formatDateTime(item.createdAt || item.created_at),
    updatedAt: formatDateTime(item.updatedAt || item.updated_at)
  }))
}

function formatDateTime(value) {
  if (!value) return ''
  return String(value).slice(0, 19).replace('T', ' ')
}

function resetFilters() {
  filters.keyword = ''
  filters.projectName = ''
}

function toggleProject(projectName) {
  if (expandedProjects.value.includes(projectName)) {
    expandedProjects.value = expandedProjects.value.filter(item => item !== projectName)
  } else {
    expandedProjects.value.push(projectName)
  }
}

function openEditDialog() {
  previewList.value = []
  excelFileName.value = ''
  if (excelFileInput.value) excelFileInput.value.value = ''
  resetManualForm()
  showEditDialog.value = true
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
  const normalized = {}
  Object.keys(row).forEach(key => {
    normalized[normalizeHeader(key)] = row[key]
  })
  return {
    projectName: getCellValue(normalized, ['项目名称']),
    productName: getCellValue(normalized, ['产品名称']),
    inboundModel: getCellValue(normalized, ['入库型号']),
    outboundModel: getCellValue(normalized, ['出库型号'])
  }
}

function fillMergedLikeRows(rows) {
  const last = {
    projectName: '',
    productName: '',
    outboundModel: ''
  }

  return rows.map(row => {
    const item = normalizeExcelRow(row)
    if (item.projectName) last.projectName = item.projectName
    if (item.productName) last.productName = item.productName
    if (item.outboundModel) last.outboundModel = item.outboundModel

    return {
      projectName: item.projectName || last.projectName,
      productName: item.productName || last.productName,
      inboundModel: item.inboundModel,
      outboundModel: item.outboundModel || last.outboundModel
    }
  })
}

function handleExcelChange(event) {
  const file = event.target.files?.[0]
  if (!file) return
  if (!file.name.endsWith('.xls') && !file.name.endsWith('.xlsx')) {
    alert('请上传 Excel 文件，格式为 .xls 或 .xlsx')
    event.target.value = ''
    excelFileName.value = ''
    return
  }
  excelFileName.value = file.name

  const reader = new FileReader()
  reader.onload = e => {
    try {
      const data = new Uint8Array(e.target.result)
      const workbook = XLSX.read(data, { type: 'array' })
      const worksheet = workbook.Sheets[workbook.SheetNames[0]]
      const rows = XLSX.utils.sheet_to_json(worksheet, { defval: '' })
      if (rows.length === 0) {
        alert('Excel 中未识别到有效数据，请确认第 1 行是否为表头')
        return
      }
      const headers = Object.keys(rows[0] || {}).map(key => normalizeHeader(key))
      const missing = ['项目名称', '产品名称', '入库型号', '出库型号'].filter(header => !headers.includes(header))
      if (missing.length > 0) {
        alert(`Excel 第 1 行缺少表头：${missing.join('、')}`)
        return
      }
      const list = fillMergedLikeRows(rows)
        .filter(item => item.projectName || item.productName || item.inboundModel || item.outboundModel)
      previewList.value = list
    } catch (err) {
      console.error(err)
      alert('Excel 解析失败，请检查文件格式或第 1 行表头是否正确')
    }
  }
  reader.readAsArrayBuffer(file)
}

function resetManualForm() {
  Object.assign(manualForm, {
    projectName: '',
    productName: '',
    inboundModel: '',
    outboundModel: ''
  })
}

function addManualComposition() {
  const item = {
    projectName: manualForm.projectName.trim(),
    productName: manualForm.productName.trim(),
    inboundModel: manualForm.inboundModel.trim(),
    outboundModel: manualForm.outboundModel.trim()
  }
  if (!item.projectName || !item.productName || !item.inboundModel || !item.outboundModel) {
    alert('请填写项目名称、产品名称、入库型号、出库型号')
    return
  }
  previewList.value.push(item)
  resetManualForm()
}

async function saveCompositions() {
  if (previewList.value.length === 0) {
    alert('请先上传 Excel 或手动加入板卡组成')
    return
  }
  const invalid = previewList.value.find(item => !item.projectName || !item.productName || !item.inboundModel || !item.outboundModel)
  if (invalid) {
    alert('保存失败：项目名称、产品名称、入库型号、出库型号不能为空')
    return
  }

  const res = await createBoardCompositions({ records: previewList.value })
  const result = getResponseData(res)
  if (result.code === 200) {
    alert(`保存成功，共保存 ${result.data?.count || previewList.value.length} 条`)
    showEditDialog.value = false
    await loadCompositions()
  } else {
    alert(result.msg || '保存失败')
  }
}

async function deleteComposition(item) {
  if (!confirm(`确认删除【${item.outboundModel} / ${item.inboundModel}】这条板卡组成吗？`)) return
  await deleteBoardComposition(item.id)
  await loadCompositions()
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
.upload-tip,
.manual-card {
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
.filter-card select,
.form-grid input,
.manual-grid input,
.manual-grid select {
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
.reset-btn {
  height: 38px;
  border: 0;
  border-radius: 6px;
  padding: 0 16px;
  color: #f8fafc;
  cursor: pointer;
}

.primary-btn,
.query-btn { background: #2563eb; }
.reset-btn { background: rgba(51, 65, 85, 0.92); }

.table-card {
  border-radius: 8px;
  overflow: hidden;
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

.project-row {
  background: rgba(15, 23, 42, 0.64);
}

.project-name-btn {
  border: 0;
  background: transparent;
  color: #bfdbfe;
  font-weight: 800;
  cursor: pointer;
}

.expand-icon {
  display: inline-block;
  width: 18px;
  color: #60a5fa;
}

.count-tag {
  display: inline-flex;
  padding: 4px 10px;
  border-radius: 999px;
  background: rgba(37, 99, 235, 0.16);
  color: #93c5fd;
  font-size: 12px;
  font-weight: 700;
}

.operation-col {
  width: 160px;
}

.child-panel-row > td {
  padding: 0;
  border-top: 1px solid rgba(96, 165, 250, 0.28);
  background: #020617;
}

.composition-panel {
  min-width: 840px;
  padding: 10px 12px 14px;
}

.composition-header,
.composition-row {
  display: grid;
  grid-template-columns: minmax(160px, 1fr) minmax(180px, 1.1fr) minmax(420px, 2.4fr);
  gap: 14px;
  align-items: flex-start;
}

.composition-header {
  padding: 10px 12px;
  border-radius: 6px;
  background: rgba(15, 23, 42, 0.94);
  color: #93c5fd;
  font-size: 12px;
  font-weight: 800;
}

.composition-row {
  margin-top: 8px;
  padding: 12px;
  border: 1px solid rgba(51, 65, 85, 0.76);
  border-radius: 6px;
  background: rgba(2, 6, 23, 0.72);
}

.composition-row:hover {
  background: rgba(30, 41, 59, 0.72);
}

.product-tag,
.model-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  max-width: 100%;
  min-height: 26px;
  padding: 4px 10px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 800;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-tag {
  background: rgba(37, 99, 235, 0.18);
  color: #bfdbfe;
}

.model-tag.inbound {
  background: rgba(20, 184, 166, 0.16);
  color: #99f6e4;
  max-width: none;
  height: auto;
  min-height: 28px;
  white-space: normal;
  overflow: visible;
  text-overflow: initial;
  word-break: break-word;
  line-height: 1.35;
}

.model-tag.outbound {
  background: rgba(168, 85, 247, 0.16);
  color: #ddd6fe;
}

.inbound-model-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: flex-start;
}

.chip-delete {
  flex: 0 0 auto;
  width: 18px;
  height: 18px;
  border: 0;
  border-radius: 999px;
  background: rgba(15, 23, 42, 0.72);
  color: #fca5a5;
  cursor: pointer;
  line-height: 18px;
  padding: 0;
}

.chip-delete:hover {
  background: rgba(127, 29, 29, 0.72);
  color: #fee2e2;
}

.action-group {
  display: flex;
  align-items: center;
  gap: 12px;
  justify-content: flex-start;
}

.text-btn {
  border: 0;
  background: transparent;
  color: #cbd5e1;
  cursor: pointer;
}

.text-btn.blue { color: #60a5fa; }
.text-btn.red { color: #f87171; }

.table-footer,
.muted,
.upload-tip p,
.manual-card-title span {
  color: #94a3b8;
}

.table-footer {
  padding: 14px 16px;
  font-size: 12px;
  background: #020617;
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
  width: min(980px, calc(100vw - 32px));
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
.form-grid,
.manual-card,
.preview-card {
  margin: 16px;
  border-radius: 8px;
  padding: 14px;
}

.form-grid,
.manual-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
}

.form-grid label,
.manual-grid label {
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

.manual-card-title,
.manual-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 14px;
}

.manual-actions {
  justify-content: flex-end;
  margin-top: 14px;
  margin-bottom: 0;
}

.preview-title {
  margin-bottom: 10px;
  font-weight: 700;
}
</style>
