<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>硬件版本管理</h1>
      </div>

      <div class="header-actions">
        <button class="secondary-btn" @click="exportHardwareVersions">
          导出硬件版本
        </button>

        <button class="primary-btn" @click="openCreateDialog">
          新增硬件版本
        </button>
      </div>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索硬件版本 / 终端类型 / 绑定项目"
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

      <select v-model="filters.status">
        <option value="">全部状态</option>
        <option value="trial">试产</option>
        <option value="batch">批量</option>
        <option value="sample">样品</option>
      </select>

      <button class="query-btn">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <!-- 数据表格 -->
    <div class="table-card">
      <table>
        <thead>
          <tr>
            <th>硬件版本</th>
            <th>终端类型</th>
            <th>绑定项目</th>
            <th>版本状态</th>
            <th>负责人</th>
            <th>更新时间</th>
            <th>ZIP 文件</th>
            <th class="operation-col">操作</th>
          </tr>
        </thead>

        <tbody>
          <tr v-for="item in filteredHardwareList" :key="item.id">
            <td>
              <button class="version-link" @click="openZipUploadDialog(item)">
                {{ item.hardwareVersion }}
              </button>
              <div class="version-desc">{{ item.description }}</div>
            </td>

            <td>
              <span class="device-tag">
                {{ item.deviceType }}
              </span>
            </td>

            <td>
              <div class="project-list">
                <span
                  v-for="project in item.bindProjects"
                  :key="project"
                  class="project-tag"
                >
                  {{ project }}
                </span>
              </div>
            </td>

            <td>
              <span class="status-tag" :class="item.status">
                {{ getStatusText(item.status) }}
              </span>
            </td>

            <td>{{ item.owner }}</td>

            <td class="muted">{{ item.updateTime }}</td>

            <td>
              <span v-if="item.zipFileName" class="file-name">
                {{ item.zipFileName }}
              </span>
              <span v-else class="muted">未上传</span>
            </td>

            <td class="operation-col">
              <div class="action-group">
                <button class="text-btn" @click="viewHardware(item)">
                  查看
                </button>

                <button class="text-btn blue" @click="openEditDialog(item)">
                  修改
                </button>

                <button class="text-btn yellow" @click="openZipUploadDialog(item)">
                  上传ZIP
                </button>

                <button class="text-btn green" @click="downloadZip(item)">
                  下载
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <div class="table-footer">
        共 {{ filteredHardwareList.length }} 条硬件版本记录
      </div>
    </div>

    <!-- 新增 / 修改硬件版本弹窗 -->
    <div v-if="showEditDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>{{ editMode === 'create' ? '新增硬件版本' : '修改硬件版本' }}</h3>
          <button @click="showEditDialog = false">×</button>
        </div>

        <div class="form-grid">
          <label>
            硬件版本号
            <input
              v-model="hardwareForm.hardwareVersion"
              placeholder="例如：HW_V2.1.0"
            />
          </label>

          <label>
            终端类型
            <select v-model="hardwareForm.deviceType">
              <option value="">请选择终端类型</option>
              <option
                v-for="type in deviceTypeOptions"
                :key="type"
                :value="type"
              >
                {{ type }}
              </option>
            </select>
          </label>

          <label>
            负责人
            <input
              v-model="hardwareForm.owner"
              placeholder="请输入硬件负责人"
            />
          </label>

          <label>
            版本状态
            <select v-model="hardwareForm.status">
              <option value="trial">试产</option>
              <option value="batch">批量</option>
              <option value="sample">样品</option>
            </select>
          </label>

          <label class="full-row">
            绑定项目
            <div class="checkbox-list">
              <label
                v-for="project in projectOptions"
                :key="project"
                class="checkbox-item"
              >
                <input
                  type="checkbox"
                  :value="project"
                  v-model="hardwareForm.bindProjects"
                />
                <span>{{ project }}</span>
              </label>
            </div>
          </label>

          <label class="full-row">
            版本说明
            <textarea
              v-model="hardwareForm.description"
              placeholder="例如：适配广播控制盒新板卡，修改电源模块和音频接口"
            ></textarea>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showEditDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="saveHardwareVersion">
            保存
          </button>
        </div>
      </div>
    </div>

    <!-- 上传 ZIP 弹窗 -->
    <div v-if="showZipDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>上传硬件版本 ZIP 文件</h3>
          <button @click="showZipDialog = false">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>硬件版本</span>
            <strong>{{ currentZipHardware?.hardwareVersion }}</strong>
          </div>

          <div>
            <span>终端类型</span>
            <strong>{{ currentZipHardware?.deviceType }}</strong>
          </div>

          <div>
            <span>当前 ZIP 文件</span>
            <strong>{{ currentZipHardware?.zipFileName || '未上传' }}</strong>
          </div>

          <div>
            <span>更新时间</span>
            <strong>{{ currentZipHardware?.updateTime }}</strong>
          </div>
        </div>

        <div class="form-grid zip-form">
          <label class="full-row">
            ZIP 压缩包
            <input
              type="file"
              accept=".zip,application/zip,application/x-zip-compressed"
              @change="handleZipFileChange"
            />
          </label>

          <label class="full-row">
            上传说明
            <textarea
              v-model="zipForm.remark"
              placeholder="例如：上传硬件原理图、PCB资料、BOM清单、生产资料等"
            ></textarea>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showZipDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="saveZipFile">
            保存 ZIP
          </button>
        </div>
      </div>
    </div>

    <!-- 查看详情弹窗 -->
    <div v-if="selectedHardware" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>硬件版本详情</h3>
          <button @click="selectedHardware = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>硬件版本</span>
            <strong>{{ selectedHardware.hardwareVersion }}</strong>
          </div>

          <div>
            <span>终端类型</span>
            <strong>{{ selectedHardware.deviceType }}</strong>
          </div>

          <div>
            <span>版本状态</span>
            <strong>{{ getStatusText(selectedHardware.status) }}</strong>
          </div>

          <div>
            <span>负责人</span>
            <strong>{{ selectedHardware.owner }}</strong>
          </div>

          <div>
            <span>更新时间</span>
            <strong>{{ selectedHardware.updateTime }}</strong>
          </div>

          <div>
            <span>ZIP 文件</span>
            <strong>{{ selectedHardware.zipFileName || '未上传' }}</strong>
          </div>
        </div>

        <div class="relation-card">
          <span>项目绑定关系</span>

          <div class="relation-list">
            <span
              v-for="project in selectedHardware.bindProjects"
              :key="project"
              class="project-tag"
            >
              {{ project }}
            </span>
          </div>
        </div>

        <div class="remark-card">
          <span>版本说明</span>
          <p>{{ selectedHardware.description || '暂无说明' }}</p>
        </div>

        <div class="dialog-footer">
          <button class="primary-btn" @click="selectedHardware = null">
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
  status: ''
})

const showEditDialog = ref(false)
const showZipDialog = ref(false)
const selectedHardware = ref(null)
const currentEditHardware = ref(null)
const currentZipHardware = ref(null)

const editMode = ref('create')

const deviceTypeOptions = [
  '广播控制盒',
  '客室解码板',
  '编码板',
  '乘客报警器',
  '司机室广播控制盒',
  '解码板',
  '功放板',
  '噪声检测器'
]

const projectOptions = [
  '香港屯马项目',
  '波尔图二期项目',
  '阿根廷有轨项目',
  '波哥大有轨项目',
  '成都项目'
]

const hardwareForm = reactive({
  hardwareVersion: '',
  deviceType: '',
  owner: '',
  status: 'trial',
  bindProjects: [],
  description: ''
})

const zipForm = reactive({
  file: null,
  fileName: '',
  fileUrl: '',
  remark: ''
})

const hardwareVersionList = ref([
  {
    id: 1,
    hardwareVersion: 'HD-CRRC-HKTM.01.V1.1.0',
    deviceType: '广播控制盒',
    bindProjects: ['香港屯马项目', '波尔图二期项目'],
    status: 'batch',
    owner: '王宇',
    updateTime: '2026-05-10',
    zipFileName: '',
    zipFileUrl: '',
    description: '广播控制盒硬件版本，适配司机人工广播和乘客报警功能'
  },
  {
    id: 2,
    hardwareVersion: 'HD-CRRC-AGTB-04.T1.1.0',
    deviceType: '客室解码板',
    bindProjects: ['阿根廷有轨项目'],
    status: 'trial',
    owner: '王宇',
    updateTime: '2026-05-16',
    zipFileName: '',
    zipFileUrl: '',
    description: '客室解码板硬件试产版本，用于车厢广播解码'
  },
  {
    id: 3,
    hardwareVersion: 'HD-CRRC-BOGT-03.T1.1.0',
    deviceType: '乘客报警器',
    bindProjects: ['波哥大有轨项目'],
    status: 'batch',
    owner: '郑宇',
    updateTime: '2026-05-18',
    zipFileName: '',
    zipFileUrl: '',
    description: '乘客报警器硬件批量版本'
  },
  {
    id: 4,
    hardwareVersion: 'HD-CRRC-DUBAI-05.S1.1.0',
    deviceType: '编码板',
    bindProjects: ['迪拜项目'],
    status: 'sample',
    owner: '郑宇',
    updateTime: '2026-05-20',
    zipFileName: '',
    zipFileUrl: '',
    description: '编码板硬件样品版本'
  }
])

const filteredHardwareList = computed(() => {
  return hardwareVersionList.value.filter(item => {
    const keywordMatch =
      !filters.keyword ||
      item.hardwareVersion.includes(filters.keyword) ||
      item.deviceType.includes(filters.keyword) ||
      item.owner.includes(filters.keyword) ||
      item.bindProjects.some(project => project.includes(filters.keyword))

    const deviceTypeMatch =
      !filters.deviceType || item.deviceType === filters.deviceType

    const statusMatch =
      !filters.status || item.status === filters.status

    return keywordMatch && deviceTypeMatch && statusMatch
  })
})

function getStatusText(status) {
  const map = {
    trial: '试产',
    batch: '批量',
    sample: '样品'
  }

  return map[status] || status
}

function resetFilters() {
  filters.keyword = ''
  filters.deviceType = ''
  filters.status = ''
}

function openCreateDialog() {
  editMode.value = 'create'
  currentEditHardware.value = null

  hardwareForm.hardwareVersion = ''
  hardwareForm.deviceType = ''
  hardwareForm.owner = ''
  hardwareForm.status = 'trial'
  hardwareForm.bindProjects = []
  hardwareForm.description = ''

  showEditDialog.value = true
}

function openEditDialog(item) {
  editMode.value = 'edit'
  currentEditHardware.value = item

  hardwareForm.hardwareVersion = item.hardwareVersion
  hardwareForm.deviceType = item.deviceType
  hardwareForm.owner = item.owner
  hardwareForm.status = item.status
  hardwareForm.bindProjects = [...item.bindProjects]
  hardwareForm.description = item.description

  showEditDialog.value = true
}

function saveHardwareVersion() {
  if (!hardwareForm.hardwareVersion) {
    alert('请输入硬件版本号')
    return
  }

  if (!hardwareForm.deviceType) {
    alert('请选择终端类型')
    return
  }

  if (hardwareForm.bindProjects.length === 0) {
    alert('请至少选择一个绑定项目')
    return
  }

  if (editMode.value === 'create') {
    hardwareVersionList.value.unshift({
      id: Date.now(),
      hardwareVersion: hardwareForm.hardwareVersion,
      deviceType: hardwareForm.deviceType,
      bindProjects: [...hardwareForm.bindProjects],
      status: hardwareForm.status,
      owner: hardwareForm.owner || '未填写',
      updateTime: new Date().toISOString().slice(0, 10),
      zipFileName: '',
      zipFileUrl: '',
      description: hardwareForm.description
    })
  } else {
    currentEditHardware.value.hardwareVersion = hardwareForm.hardwareVersion
    currentEditHardware.value.deviceType = hardwareForm.deviceType
    currentEditHardware.value.owner = hardwareForm.owner || '未填写'
    currentEditHardware.value.status = hardwareForm.status
    currentEditHardware.value.bindProjects = [...hardwareForm.bindProjects]
    currentEditHardware.value.description = hardwareForm.description
    currentEditHardware.value.updateTime = new Date().toISOString().slice(0, 10)
  }

  showEditDialog.value = false
}

function viewHardware(item) {
  selectedHardware.value = item
}

function openZipUploadDialog(item) {
  currentZipHardware.value = item

  zipForm.file = null
  zipForm.fileName = ''
  zipForm.fileUrl = ''
  zipForm.remark = ''

  showZipDialog.value = true
}

function handleZipFileChange(event) {
  const file = event.target.files[0]
  if (!file) return

  const isZip = file.name.toLowerCase().endsWith('.zip')

  if (!isZip) {
    alert('只能上传 ZIP 压缩包文件')
    event.target.value = ''
    return
  }

  zipForm.file = file
  zipForm.fileName = file.name
  zipForm.fileUrl = URL.createObjectURL(file)
}

function saveZipFile() {
  if (!currentZipHardware.value) return

  if (!zipForm.file) {
    alert('请选择 ZIP 压缩包文件')
    return
  }

  currentZipHardware.value.zipFileName = zipForm.fileName
  currentZipHardware.value.zipFileUrl = zipForm.fileUrl
  currentZipHardware.value.updateTime = new Date().toISOString().slice(0, 10)

  if (zipForm.remark) {
    currentZipHardware.value.description = zipForm.remark
  }

  showZipDialog.value = false

  alert(`硬件版本【${currentZipHardware.value.hardwareVersion}】ZIP 文件已上传`)
}

function downloadZip(item) {
  if (!item.zipFileUrl) {
    alert('当前硬件版本暂未上传 ZIP 文件')
    return
  }

  const link = document.createElement('a')
  link.href = item.zipFileUrl
  link.download = item.zipFileName || `${item.hardwareVersion}.zip`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

function exportHardwareVersions() {
  const header = [
    '硬件版本',
    '终端类型',
    '绑定项目',
    '版本状态',
    '负责人',
    '更新时间',
    'ZIP文件'
  ]

  const rows = hardwareVersionList.value.map(item => [
    item.hardwareVersion,
    item.deviceType,
    item.bindProjects.join('、'),
    getStatusText(item.status),
    item.owner,
    item.updateTime,
    item.zipFileName || '未上传'
  ])

  const csvContent = [header, ...rows]
    .map(row => row.map(cell => `"${cell}"`).join(','))
    .join('\n')

  const blob = new Blob(['\uFEFF' + csvContent], {
    type: 'text/csv;charset=utf-8;'
  })

  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = '硬件版本列表.csv'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
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

.page-header p {
  margin: 8px 0 0;
  color: #94a3b8;
  font-size: 14px;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.primary-btn,
.secondary-btn,
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

.secondary-btn,
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
  grid-template-columns: 1.4fr 220px 180px 90px 90px;
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
  min-height: 84px;
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

.table-card table {
  width: 100%;
  border-collapse: collapse;
  table-layout: fixed;
}

.table-card thead {
  background: #020617;
}

.table-card th {
  padding: 14px 16px;
  border-bottom: 1px solid #1e293b;
  color: #94a3b8;
  font-size: 12px;
  font-weight: 600;
  text-align: left;
  white-space: nowrap;
}

.table-card td {
  padding: 15px 16px;
  border-bottom: 1px solid #1e293b;
  color: #e2e8f0;
  font-size: 13px;
  vertical-align: middle;
}

.version-link {
  border: none;
  background: transparent;
  color: #60a5fa;
  font-size: 13px;
  font-weight: 800;
  cursor: pointer;
  padding: 0;
}

.version-link:hover {
  color: #93c5fd;
  text-decoration: underline;
}

.version-desc {
  margin-top: 4px;
  color: #64748b;
  font-size: 12px;
  word-break: break-all;
}

.device-tag,
.project-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
}

.device-tag {
  background: #1d4ed833;
  color: #60a5fa;
}

.project-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.project-tag {
  background: #0f766e33;
  color: #5eead4;
}

.status-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
}

.status-tag.trial {
  background: #d9770633;
  color: #fbbf24;
}

.status-tag.batch {
  background: #16a34a33;
  color: #4ade80;
}

.status-tag.sample {
  background: #1d4ed833;
  color: #60a5fa;
}

.file-name {
  color: #cbd5e1;
  font-size: 12px;
  word-break: break-all;
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
  flex-wrap: wrap;
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

.text-btn.green {
  color: #4ade80;
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
  width: 720px;
  max-width: 100%;
  background: #0f172a;
  border: 1px solid #334155;
  border-radius: 16px;
  color: #f8fafc;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.45);
}

.large-dialog {
  width: 860px;
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

.checkbox-list {
  background: #020617;
  border: 1px solid #334155;
  border-radius: 8px;
  padding: 12px;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
}

.checkbox-item {
  display: flex !important;
  flex-direction: row !important;
  align-items: center;
  gap: 6px !important;
  font-size: 13px;
  color: #cbd5e1;
}

.checkbox-item input {
  width: 14px;
  height: 14px;
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
.remark-card,
.relation-card {
  background: #020617;
  border: 1px solid #1e293b;
  border-radius: 10px;
  padding: 12px;
}

.detail-card span,
.remark-card span,
.relation-card span:first-child {
  display: block;
  color: #64748b;
  font-size: 12px;
  margin-bottom: 6px;
}

.detail-card strong {
  color: #f8fafc;
  font-size: 14px;
}

.relation-card,
.remark-card {
  margin: 0 20px 20px;
}

.relation-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.remark-card p {
  margin: 0;
  color: #cbd5e1;
  font-size: 13px;
  line-height: 1.6;
}

.zip-form {
  padding-top: 0;
}

@media (max-width: 960px) {
  .filter-card {
    grid-template-columns: 1fr;
  }

  .table-card {
    overflow-x: auto;
  }

  .table-card table {
    min-width: 1200px;
  }

  .form-grid,
  .detail-card {
    grid-template-columns: 1fr;
  }

  .checkbox-list {
    grid-template-columns: 1fr;
  }
}
</style>