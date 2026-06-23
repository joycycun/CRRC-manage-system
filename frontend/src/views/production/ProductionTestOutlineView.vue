<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1>生产测试大纲</h1>
      </div>

      <button
        v-if="canUpload"
        class="primary-btn"
        @click="openUploadDialog"
      >
        上传测试大纲
      </button>
    </div>

    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索产品型号 / 板卡型号 / 终端类型 / 文件名"
      />

      <select v-model="filters.hardwareId">
        <option value="">全部产品型号</option>
        <option
          v-for="item in hardwareOptions"
          :key="item.id"
          :value="String(item.id)"
        >
          {{ formatHardwareOption(item) }}
        </option>
      </select>

      <button class="query-btn" @click="loadOutlines">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <div class="table-card">
      <div class="table-card-header">
        <div>
          <h3>测试大纲列表</h3>
          <span>共 {{ filteredOutlines.length }} 条记录</span>
        </div>
      </div>

      <div class="table-wrapper">
        <table class="batch-table">
          <thead>
            <tr>
              <th>产品型号</th>
              <th>板卡型号</th>
              <th>终端类型</th>
              <th>测试大纲文件</th>
              <th>上传人</th>
              <th>上传时间</th>
              <th>备注</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>

          <tbody>
            <template v-for="item in filteredOutlines" :key="item.id">
              <tr>
                <td>
                  <span class="version-tag">{{ item.hardwareVersion || '-' }}</span>
                </td>
                <td>
                  <button
                    class="board-models-btn"
                    :class="{ active: expandedBoardRowId === item.id }"
                    :title="item.boardModels"
                    @click="toggleBoardModels(item)"
                  >
                    <span>{{ formatBoardModelsSummary(item.boardModels) }}</span>
                    <span class="expand-icon">{{ expandedBoardRowId === item.id ? '收起' : '展开' }}</span>
                  </button>
                </td>
                <td>{{ item.deviceType || '-' }}</td>
                <td>
                  <span class="file-name" :title="item.fileName">{{ item.fileName || '-' }}</span>
                </td>
                <td>{{ item.uploaderName || '-' }}</td>
                <td class="muted">{{ item.uploadTime || '-' }}</td>
                <td>
                  <span class="remark-text" :title="item.remark">{{ item.remark || '-' }}</span>
                </td>
                <td class="operation-col">
                  <button class="text-btn" @click="downloadFile(item)">下载</button>
                </td>
              </tr>

              <tr v-if="expandedBoardRowId === item.id" class="board-detail-row">
                <td colspan="8">
                  <div class="board-detail-panel">
                    <span
                      v-for="model in splitBoardModels(item.boardModels)"
                      :key="model"
                      class="board-chip"
                    >
                      {{ model }}
                    </span>
                    <span v-if="splitBoardModels(item.boardModels).length === 0" class="empty-board-model">
                      暂无板卡型号
                    </span>
                  </div>
                </td>
              </tr>
            </template>

            <tr v-if="filteredOutlines.length === 0">
              <td colspan="8" class="empty-table">暂无生产测试大纲</td>
            </tr>
          </tbody>
        </table>
      </div>

    </div>

    <div v-if="showUploadDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>上传生产测试大纲</h3>
          <button @click="closeUploadDialog">×</button>
        </div>

        <div class="upload-tip">
          <strong>上传规则：</strong>
          <p>
            产品型号来自硬件版本管理中已登记的数据。重复上传同一产品型号时会保留历史文件，
            生产人员页面只展示最新一份。
          </p>
        </div>

        <div class="form-grid">
          <label>
            产品型号
            <select v-model.number="uploadForm.hardwareId">
              <option :value="0">请选择硬件版本里的产品</option>
              <option
                v-for="item in hardwareOptions"
                :key="item.id"
                :value="item.id"
              >
                {{ formatHardwareOption(item) }}
              </option>
            </select>
          </label>

          <label class="full-row">
            板卡型号
            <input
              v-model="uploadForm.boardModels"
              placeholder="可填写多个，用逗号分隔，例如：PA-01, AMP-02"
            />
          </label>

          <label class="full-row">
            测试大纲文件
            <div class="file-upload-row">
              <input
                ref="fileInputRef"
                class="hidden-file-input"
                type="file"
                accept=".doc,.docx,.pdf,.xls,.xlsx"
                @change="handleFileChange"
              />
              <button class="outline-btn" type="button" @click="triggerFileInput">
                选择文件
              </button>
              <strong :title="uploadForm.fileName">
                {{ uploadForm.fileName || '未选择文件' }}
              </strong>
            </div>
          </label>

          <label class="full-row">
            备注
            <textarea v-model="uploadForm.remark" placeholder="填写版本说明或更新内容"></textarea>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="closeUploadDialog">取消</button>
          <button class="primary-btn" :disabled="submitting" @click="submitUpload">
            {{ submitting ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { getHardwareVersions } from '@/api/hardware'
import { createProductionTestOutline, getProductionTestOutlines } from '@/api/productionTestOutline'
import { buildUploadFilePayload, downloadLocalFile } from '@/utils/filePreview'
import { canUseAction } from '@/utils/permission'

const outlineList = ref([])
const hardwareOptions = ref([])
const showUploadDialog = ref(false)
const submitting = ref(false)
const fileInputRef = ref(null)
const expandedBoardRowId = ref(null)

const filters = reactive({
  keyword: '',
  hardwareId: ''
})

const uploadForm = reactive({
  hardwareId: 0,
  boardModels: '',
  fileId: 0,
  fileName: '',
  fileContentType: '',
  fileData: '',
  remark: ''
})

const canUpload = computed(() => canUseAction('production:outline:upload'))

const filteredOutlines = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase()
  return outlineList.value.filter(item => {
    const matchesHardware = !filters.hardwareId || String(item.hardwareId) === filters.hardwareId
    const text = [
      item.hardwareVersion,
      item.boardModels,
      item.deviceType,
      item.fileName,
      item.uploaderName,
      item.remark
    ].join(' ').toLowerCase()
    return matchesHardware && (!keyword || text.includes(keyword))
  })
})

onMounted(async () => {
  await Promise.all([loadHardwareOptions(), loadOutlines()])
})

function formatHardwareOption(item) {
  return [item.hardwareVersion, item.deviceType].filter(Boolean).join(' / ') || `硬件版本 ${item.id}`
}

async function loadHardwareOptions() {
  try {
    const res = await getHardwareVersions()
    hardwareOptions.value = res.data?.data || []
  } catch (err) {
    console.error('加载硬件版本失败：', err)
    alert('加载硬件版本失败')
  }
}

async function loadOutlines() {
  try {
    const res = await getProductionTestOutlines()
    outlineList.value = res.data?.data || []
  } catch (err) {
    console.error('加载生产测试大纲失败：', err)
    alert(err.response?.data || '加载生产测试大纲失败')
  }
}

function resetFilters() {
  filters.keyword = ''
  filters.hardwareId = ''
  expandedBoardRowId.value = null
  loadOutlines()
}

function openUploadDialog() {
  resetUploadForm()
  showUploadDialog.value = true
}

function closeUploadDialog() {
  showUploadDialog.value = false
}

function triggerFileInput() {
  fileInputRef.value?.click()
}

async function handleFileChange(event) {
  const file = event.target.files?.[0]
  if (!file) return
  const payload = await buildUploadFilePayload(file)
  Object.assign(uploadForm, {
    fileId: payload.fileId,
    fileName: payload.fileName,
    fileContentType: payload.fileContentType,
    fileData: payload.fileData
  })
}

function getCurrentUser() {
  try {
    return JSON.parse(localStorage.getItem('user') || '{}')
  } catch (err) {
    return {}
  }
}

async function submitUpload() {
  if (!uploadForm.hardwareId) {
    alert('请选择产品型号')
    return
  }
  if (!uploadForm.fileId || !uploadForm.fileName) {
    alert('请选择测试大纲文件')
    return
  }
  const boardModels = uploadForm.boardModels
    .split(/[，,]/)
    .map(item => item.trim())
    .filter(Boolean)
    .join(', ')

  if (!boardModels) {
    alert('请填写板卡型号，多个型号请用逗号分隔')
    return
  }

  submitting.value = true
  try {
    const user = getCurrentUser()
    await createProductionTestOutline({
      hardwareId: uploadForm.hardwareId,
      boardModels,
      fileId: uploadForm.fileId,
      fileName: uploadForm.fileName,
      fileContentType: uploadForm.fileContentType,
      fileData: uploadForm.fileData,
      uploaderId: user.id || 0,
      uploaderName: user.realName || user.username || '',
      remark: uploadForm.remark
    })
    alert('保存成功')
    closeUploadDialog()
    await loadOutlines()
  } catch (err) {
    console.error('保存生产测试大纲失败：', err)
    alert(err.response?.data || '保存生产测试大纲失败')
  } finally {
    submitting.value = false
  }
}

function resetUploadForm() {
  Object.assign(uploadForm, {
    hardwareId: 0,
    boardModels: '',
    fileId: 0,
    fileName: '',
    fileContentType: '',
    fileData: '',
    remark: ''
  })
  if (fileInputRef.value) {
    fileInputRef.value.value = ''
  }
}

function downloadFile(item) {
  downloadLocalFile(item, item.fileName || '生产测试大纲')
}

function splitBoardModels(value) {
  return String(value || '')
    .split(/[，,]/)
    .map(item => item.trim())
    .filter(Boolean)
}

function formatBoardModelsSummary(value) {
  const models = splitBoardModels(value)
  if (models.length === 0) return '-'
  if (models.length <= 2) return models.join(', ')
  return `${models.slice(0, 2).join(', ')} 等 ${models.length} 个`
}

function toggleBoardModels(item) {
  expandedBoardRowId.value = expandedBoardRowId.value === item.id ? null : item.id




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
.outline-btn {
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

.primary-btn:disabled {
  cursor: not-allowed;
  opacity: 0.65;
}

.reset-btn,
.outline-btn {
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

.form-grid textarea {
  min-height: 90px;
  padding: 10px 12px;
  resize: vertical;
}

.filter-card input::placeholder,
.form-grid textarea::placeholder {
  color: #64748b;
}

.table-card {
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 14px;
  overflow: hidden;
}

.table-card-header {
  padding: 16px 18px;
  border-bottom: 1px solid #1e293b;
}

.table-card-header h3 {
  margin: 0 0 4px;
  color: #f8fafc;
  font-size: 16px;
}

.table-card-header span,
.muted {
  color: #94a3b8;
  font-size: 12px;
}

.table-wrapper {
  width: 100%;
  overflow-x: auto;
  overflow-y: hidden;
  scrollbar-width: thin;
  scrollbar-color: #334155 #020617;
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

.batch-table {
  width: 100%;
  min-width: 1180px;
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

.batch-table tbody tr:hover {
  background: #1e293b80;
}

.batch-table th:nth-child(1),
.batch-table td:nth-child(1) {
  width: 180px;
}

.batch-table th:nth-child(2),
.batch-table td:nth-child(2) {
  width: 150px;
}

.batch-table th:nth-child(3),
.batch-table td:nth-child(3) {
  width: 150px;
}

.batch-table th:nth-child(4),
.batch-table td:nth-child(4) {
  width: 300px;
}

.batch-table th:nth-child(5),
.batch-table td:nth-child(5) {
  width: 120px;
}

.batch-table th:nth-child(6),
.batch-table td:nth-child(6) {
  width: 160px;
}

.batch-table th:nth-child(7),
.batch-table td:nth-child(7) {
  width: 190px;
}

.batch-table th:nth-child(8),
.batch-table td:nth-child(8) {
  width: 150px;
}

.operation-col {
  width: 150px;
}

.version-tag {
  display: inline-block;
  max-width: 100%;
  padding: 4px 9px;
  border-radius: 999px;
  background: #1d4ed833;
  color: #93c5fd;
  font-size: 12px;
  font-weight: 700;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;
}

.file-name,
.remark-text {
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;
}

.board-models-btn {
  width: 100%;
  min-width: 0;
  border: none;
  background: transparent;
  color: #93c5fd;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 0;
  font-size: 13px;
  cursor: pointer;
}

.board-models-btn span:first-child {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.board-models-btn:hover,
.board-models-btn.active {
  color: #bfdbfe;
}

.expand-icon {
  flex-shrink: 0;
  color: #64748b;
  font-size: 12px;
}

.board-detail-row td {
  padding: 0 16px 16px;
  background: #020617;
  white-space: normal;
  overflow: visible;
}

.board-detail-panel {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 12px;
  border: 1px solid #1e293b;
  border-radius: 10px;
  background: #0f172a;
}

.board-chip {
  display: inline-flex;
  align-items: center;
  max-width: 220px;
  min-height: 26px;
  padding: 4px 10px;
  border-radius: 999px;
  background: #0f766e33;
  color: #5eead4;
  font-size: 12px;
  font-weight: 700;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.empty-board-model {
  color: #64748b;
  font-size: 13px;
}

.empty-table {
  color: #64748b;
  text-align: center !important;
  padding: 38px 16px !important;
}

.table-footer {
  padding: 14px 16px;
  color: #64748b;
  font-size: 12px;
  border-top: 1px solid #1e293b;
  background: #020617;
}

.text-btn {
  border: none;
  background: transparent;
  color: #cbd5e1;
  padding: 0 8px 0 0;
  font-size: 13px;
  cursor: pointer;
}

.text-btn.blue {
  color: #60a5fa;
}

.text-btn:hover {
  color: #93c5fd;
  text-decoration: underline;
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
  cursor: pointer;
  font-size: 24px;
  color: #94a3b8;
}

.upload-tip {
  margin: 20px 20px 0;
  padding: 14px;
  background: #020617;
  border: 1px solid #1e293b;
  border-radius: 12px;
  color: #cbd5e1;
}

.upload-tip strong {
  display: block;
  margin-bottom: 8px;
  color: #f8fafc;
  font-size: 14px;
}

.upload-tip p {
  margin: 0;
  color: #94a3b8;
  font-size: 13px;
  line-height: 1.7;
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

.file-upload-row {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.file-upload-row strong {
  min-width: 0;
  overflow: hidden;
  color: #e2e8f0;
  font-size: 13px;
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.hidden-file-input {
  display: none;
}

.dialog-footer {
  padding: 16px 20px;
  border-top: 1px solid #1e293b;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

@media (max-width: 900px) {
  .filter-card {
    grid-template-columns: 1fr;
  }

  .form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
