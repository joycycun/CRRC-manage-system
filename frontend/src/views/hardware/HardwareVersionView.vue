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

        <button v-if="canUseAction('hardware:create')" class="primary-btn" @click="openCreateDialog">
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

      <button class="query-btn" @click="loadHardwareVersions">查询</button>
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
              <button class="version-link" @click="viewHardware(item)">
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

                <button v-if="canUseAction('hardware:update')" class="text-btn blue" @click="openEditDialog(item)">
                  修改
                </button>

                <button v-if="canUseAction('hardware:upload')" class="text-btn yellow" @click="openZipUploadDialog(item)">
                  上传ZIP
                </button>

                <button v-if="canUseAction('hardware:download')" class="text-btn green" @click="downloadZip(item)">
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
            版本状态
            <select v-model="hardwareForm.status">
              <option value="trial">试产</option>
              <option value="batch">批量</option>
              <option value="sample">样品</option>
            </select>
          </label>

          <label>
            当前负责人
            <input
              :value="editMode === 'create' ? currentUserName : hardwareForm.owner"
              disabled
            />
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
            硬件文件 ZIP
            <input
              type="file"
              accept=".zip,application/zip,application/x-zip-compressed"
              @change="handleHardwareFileChange"
            />

            <span v-if="hardwareForm.zipFileName" class="selected-file">
              已选择：{{ hardwareForm.zipFileName }}
            </span>

            <span v-else class="file-tip">
              可在新增硬件版本时直接上传硬件原理图、PCB资料、BOM清单、生产资料等 ZIP 压缩包。
            </span>
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

    <!-- 单独上传 ZIP 弹窗 -->
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
          <button v-if="canUseAction('hardware:download')" class="reset-btn" @click="downloadZip(selectedHardware)">
            下载 ZIP
          </button>

          <button class="primary-btn" @click="selectedHardware = null">
            关闭
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { canUseAction } from '@/utils/permission'

import { getProjects } from '@/api/project'

import {
  getHardwareVersions,
  createHardwareVersion,
  updateHardwareVersion,
  uploadHardwareZip
} from '@/api/hardware'

const currentUserName = ref(
  localStorage.getItem('username') ||
  localStorage.getItem('accountName') ||
  localStorage.getItem('realName') ||
  '当前用户'
)

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

const projectOptions = ref([])
const projectMap = ref({})

const hardwareForm = reactive({
  hardwareVersion: '',
  deviceType: '',
  owner: '',
  status: 'trial',
  bindProjects: [],
  description: '',
  zipFile: null,
  zipFileName: '',
  zipFileUrl: ''
})

const zipForm = reactive({
  file: null,
  fileName: '',
  fileUrl: '',
  remark: ''
})

const hardwareVersionList = ref([])

onMounted(async () => {
  await loadProjects()
  await loadHardwareVersions()
})

function getResponseData(res) {
  if (res && res.data) return res.data
  return res
}

function formatDate(value) {
  if (!value) return ''
  if (typeof value === 'string') return value.slice(0, 10)
  return value
}

function backendStatusToFrontend(status) {
  const map = {
    试产: 'trial',
    批量: 'batch',
    样品: 'sample',
    trial: 'trial',
    batch: 'batch',
    sample: 'sample'
  }

  return map[status] || status || 'trial'
}

function frontendStatusToBackend(status) {
  const map = {
    trial: '试产',
    batch: '批量',
    sample: '样品'
  }

  return map[status] || status || '试产'
}

function findProjectName(projectId) {
  const found = Object.entries(projectMap.value).find(([, id]) => Number(id) === Number(projectId))
  return found ? found[0] : `项目ID-${projectId}`
}

function getProjectIdsByNames(projectNames) {
  return projectNames
    .map(name => projectMap.value[name])
    .filter(id => id !== undefined && id !== null)
}

async function loadProjects() {
  try {
    const res = await getProjects()
    const result = getResponseData(res)

    console.log('项目列表返回：', result)

    if (result.code !== 200) {
      alert(result.msg || '加载项目失败')
      return
    }

    const list = result.data || []

    projectOptions.value = list.map(item => item.projectName)

    const map = {}
    list.forEach(item => {
      map[item.projectName] = item.id
    })

    projectMap.value = map
  } catch (err) {
    console.error('加载项目失败：', err)
    alert('加载项目失败')
  }
}

async function loadHardwareVersions() {
  try {
    const res = await getHardwareVersions()
    const result = getResponseData(res)

    console.log('硬件版本列表返回：', result)

    if (result.code !== 200) {
      alert(result.msg || '加载硬件版本失败')
      return
    }

    hardwareVersionList.value = (result.data || []).map(item => normalizeHardware(item))
  } catch (err) {
    console.error('加载硬件版本失败：', err)
    alert('加载硬件版本失败，请检查后端接口')
  }
}

function normalizeHardware(item) {
  let bindProjects = []

  if (Array.isArray(item.bindProjects)) {
    bindProjects = item.bindProjects
  } else if (Array.isArray(item.projectNames)) {
    bindProjects = item.projectNames
  } else if (item.projectName) {
    bindProjects = [item.projectName]
  } else if (item.projectId) {
    bindProjects = [findProjectName(item.projectId)]
  }

  return {
    id: item.id,
    hardwareVersion: item.hardwareVersion || '',
    deviceType: item.deviceType || '',
    bindProjects,
    status: backendStatusToFrontend(item.status),
    backendStatus: item.status || '',
    ownerId: item.ownerId || 1,
    owner: item.owner || item.ownerName || '未分配',
    ownerName: item.ownerName || item.owner || '未分配',
    updateTime: formatDate(item.updateTime || item.updatedAt || item.createdAt),
    zipFileId: item.zipFileId || item.zipFileID || 0,
    zipFileName:
      item.zipFileName ||
      item.fileName ||
      item.fileDisplayName ||
      (item.zipFileId || item.zipFileID ? `文件ID-${item.zipFileId || item.zipFileID}.zip` : ''),
    zipFileUrl: item.zipFileUrl || item.fileUrl || '',
    description: item.description || ''
  }
}

const filteredHardwareList = computed(() => {
  return hardwareVersionList.value.filter(item => {
    const keywordMatch =
      !filters.keyword ||
      item.hardwareVersion.includes(filters.keyword) ||
      item.deviceType.includes(filters.keyword) ||
      item.owner.includes(filters.keyword) ||
      item.bindProjects.some(project => project.includes(filters.keyword)) ||
      (item.zipFileName && item.zipFileName.includes(filters.keyword))

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

function resetHardwareForm() {
  hardwareForm.hardwareVersion = ''
  hardwareForm.deviceType = ''
  hardwareForm.owner = ''
  hardwareForm.status = 'trial'
  hardwareForm.bindProjects = []
  hardwareForm.description = ''
  hardwareForm.zipFile = null
  hardwareForm.zipFileName = ''
  hardwareForm.zipFileUrl = ''
}

function openCreateDialog() {
  editMode.value = 'create'
  currentEditHardware.value = null

  resetHardwareForm()

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
  hardwareForm.zipFile = null
  hardwareForm.zipFileName = item.zipFileName || ''
  hardwareForm.zipFileUrl = item.zipFileUrl || ''

  showEditDialog.value = true
}

function handleHardwareFileChange(event) {
  const file = event.target.files[0]
  if (!file) return

  const isZip = file.name.toLowerCase().endsWith('.zip')

  if (!isZip) {
    alert('只能上传 ZIP 压缩包文件')
    event.target.value = ''
    return
  }

  hardwareForm.zipFile = file
  hardwareForm.zipFileName = file.name
  hardwareForm.zipFileUrl = URL.createObjectURL(file)
}

async function saveHardwareVersion() {
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

  const projectIds = getProjectIdsByNames(hardwareForm.bindProjects)

  if (projectIds.length === 0) {
    alert('没有找到绑定项目ID，请重新选择项目')
    return
  }

  const payload = {
    hardwareVersion: hardwareForm.hardwareVersion,
    deviceType: hardwareForm.deviceType,
    projectId: projectIds[0],
    projectIds,
    bindProjects: hardwareForm.bindProjects,
    status: frontendStatusToBackend(hardwareForm.status),
    ownerId: 1,
    owner: editMode.value === 'create'
      ? currentUserName.value
      : (hardwareForm.owner || currentUserName.value),
    ownerName: editMode.value === 'create'
      ? currentUserName.value
      : (hardwareForm.owner || currentUserName.value),
    zipFileId: hardwareForm.zipFileName ? 1 : 0,
    zipFileName: hardwareForm.zipFileName || '',
    description: hardwareForm.description || ''
  }

  try {
    let res

    if (editMode.value === 'create') {
      res = await createHardwareVersion(payload)
    } else {
      res = await updateHardwareVersion(currentEditHardware.value.id, payload)
    }

    const result = getResponseData(res)

    console.log('保存硬件版本返回：', result)

    if (result.code === 200) {
      alert(editMode.value === 'create' ? '新增硬件版本成功' : '修改硬件版本成功')
      showEditDialog.value = false
      await loadHardwareVersions()
    } else {
      alert(result.msg || '保存硬件版本失败')
    }
  } catch (err) {
    console.error('保存硬件版本失败：', err)
    alert('保存硬件版本失败，请检查后端接口')
  }
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

async function saveZipFile() {
  if (!currentZipHardware.value) return

  if (!zipForm.file) {
    alert('请选择 ZIP 压缩包文件')
    return
  }

  const payload = {
    zipFileId: 1,
    zipFileName: zipForm.fileName,
    remark: zipForm.remark || ''
  }

  try {
    const res = await uploadHardwareZip(currentZipHardware.value.id, payload)
    const result = getResponseData(res)

    console.log('上传硬件 ZIP 返回：', result)

    if (result.code === 200) {
      alert(`硬件版本【${currentZipHardware.value.hardwareVersion}】ZIP 文件已上传`)
      showZipDialog.value = false
      await loadHardwareVersions()
    } else {
      alert(result.msg || '上传 ZIP 失败')
    }
  } catch (err) {
    console.error('上传硬件 ZIP 失败：', err)
    alert('上传 ZIP 失败，请检查后端接口')
  }
}

function downloadZip(item) {
  if (!item.zipFileUrl) {
    alert('当前还没有接真实文件下载，后面做 project_files 文件上传下载时再接')
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

.selected-file {
  margin-top: 6px;
  color: #5eead4;
  font-size: 12px;
}

.file-tip {
  margin-top: 6px;
  color: #64748b;
  font-size: 12px;
  line-height: 1.6;
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
  display: inline-block;
  max-width: 150px;
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
