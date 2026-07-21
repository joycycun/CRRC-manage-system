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

    <!-- 项目折叠列表 -->
    <div class="project-version-groups">
      <div
        v-for="group in groupedHardwareVersions"
        :key="group.projectName"
        class="project-version-group"
      >
        <button class="project-group-header" @click="toggleProjectGroup(group.projectName)">
          <span class="fold-icon">{{ isProjectCollapsed(group.projectName) ? '›' : '⌄' }}</span>
          <span class="project-group-title">{{ group.projectName }}</span>
          <span class="project-group-count">{{ group.items.length }} 个硬件版本</span>
        </button>

        <div v-show="!isProjectCollapsed(group.projectName)" class="table-card">
          <table>
            <thead>
              <tr>
                <th>硬件版本</th>
                <th>终端类型</th>
                <th>版本状态</th>
                <th>负责人</th>
                <th>更新时间</th>
                <th>硬件更改文档</th>
                <th class="operation-col">操作</th>
              </tr>
            </thead>

            <tbody>
              <tr v-for="item in group.items" :key="`${group.projectName}-${item.id}`">
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
                  <span class="status-tag" :class="item.status">
                    {{ getStatusText(item.status) }}
                  </span>
                </td>

                <td>{{ item.owner }}</td>

                <td class="muted">{{ item.updateTime }}</td>

                <td>
                  <span v-if="item.changeDocFileName" class="file-name">
                    {{ item.changeDocFileName }}
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

                    <button v-if="canUseAction('hardware:upload')" class="text-btn yellow" @click="openDocumentUploadDialog(item)">
                      上传文档
                    </button>

                    <button v-if="canUseAction('hardware:download')" class="text-btn green" @click="downloadDocument(item)">
                      下载
                    </button>

                    <button v-if="isSystemAdmin" class="text-btn danger" @click="deleteHardware(item)">
                      删除
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div v-if="groupedHardwareVersions.length === 0" class="empty-card">
        暂无硬件版本记录
      </div>

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
            硬件更改文档
            <input
              type="file"
              accept=".doc,.docx,.txt,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document,text/plain"
              @change="handleHardwareFileChange"
            />

            <span v-if="hardwareForm.changeDocFileName" class="selected-file">
              已选择：{{ hardwareForm.changeDocFileName }}
            </span>

            <span v-else class="file-tip">
              支持上传 Word 或 TXT 格式的硬件更改说明文档。
            </span>
          </label>

          <label class="full-row">
            版本说明
            <textarea
              v-model="hardwareForm.description"
              placeholder="例如：适配控制盒（主）新板卡，修改电源模块和音频接口"
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

    <!-- 单独上传硬件更改文档弹窗 -->
    <div v-if="showDocumentDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>上传硬件更改文档</h3>
          <button @click="showDocumentDialog = false">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>硬件版本</span>
            <strong>{{ currentDocumentHardware?.hardwareVersion }}</strong>
          </div>

          <div>
            <span>终端类型</span>
            <strong>{{ currentDocumentHardware?.deviceType }}</strong>
          </div>

          <div>
            <span>当前更改文档</span>
            <strong>{{ currentDocumentHardware?.changeDocFileName || '未上传' }}</strong>
          </div>

          <div>
            <span>更新时间</span>
            <strong>{{ currentDocumentHardware?.updateTime }}</strong>
          </div>
        </div>

        <div class="form-grid document-form">
          <label class="full-row">
            Word/TXT 文档
            <input
              type="file"
              accept=".doc,.docx,.txt,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document,text/plain"
              @change="handleDocumentFileChange"
            />
          </label>

          <label class="full-row">
            上传说明
            <textarea
              v-model="documentForm.remark"
              placeholder="例如：记录本版本硬件更改点、适配说明、注意事项"
            ></textarea>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showDocumentDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="saveDocumentFile">
            保存文档
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
            <span>硬件更改文档</span>
            <strong>{{ selectedHardware.changeDocFileName || '未上传' }}</strong>
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
          <button v-if="canUseAction('hardware:download')" class="reset-btn" @click="downloadDocument(selectedHardware)">
            下载文档
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
import { canUseAction, getStoredRoles } from '@/utils/permission'
import { buildUploadFilePayload, getFilePreviewUrl } from '@/utils/filePreview'
import { DEVICE_TYPE_OPTIONS } from '@/constants/deviceTypes'

import { getProjects } from '@/api/project'

import {
  getHardwareVersions,
  createHardwareVersion,
  updateHardwareVersion,
  uploadHardwareDocument,
  deleteHardwareVersion
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
const showDocumentDialog = ref(false)
const selectedHardware = ref(null)
const currentEditHardware = ref(null)
const currentDocumentHardware = ref(null)

const editMode = ref('create')

const deviceTypeOptions = DEVICE_TYPE_OPTIONS

const projectOptions = ref([])
const projectMap = ref({})
const collapsedProjects = reactive({})

const hardwareForm = reactive({
  hardwareVersion: '',
  deviceType: '',
  owner: '',
  status: 'trial',
  bindProjects: [],
  description: '',
  changeDocFile: null,
  changeDocFileId: 0,
  changeDocFileName: '',
  changeDocFileUrl: '',
  changeDocFileContentType: '',
  changeDocFileData: ''
})

const documentForm = reactive({
  file: null,
  fileId: 0,
  fileName: '',
  fileUrl: '',
  fileContentType: '',
  fileData: '',
  remark: ''
})

const hardwareVersionList = ref([])
const isSystemAdmin = computed(() => getStoredRoles().includes('system_admin'))

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
    changeDocFileId: item.changeDocFileId || item.changeDocFileID || 0,
    changeDocFileName:
      item.changeDocFileName ||
      item.fileName ||
      item.fileDisplayName ||
      (item.changeDocFileId || item.changeDocFileID ? `文件ID-${item.changeDocFileId || item.changeDocFileID}.docx` : ''),
    changeDocFileUrl: item.changeDocFileUrl || item.fileUrl || getFilePreviewUrl(item.changeDocFileId || item.changeDocFileID),
    changeDocDownloadUrl: item.changeDocDownloadUrl || '',
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
      (item.changeDocFileName && item.changeDocFileName.includes(filters.keyword))

    const deviceTypeMatch =
      !filters.deviceType || item.deviceType === filters.deviceType

    const statusMatch =
      !filters.status || item.status === filters.status

    return keywordMatch && deviceTypeMatch && statusMatch
  })
})

const groupedHardwareVersions = computed(() => {
  const groupMap = new Map()

  filteredHardwareList.value.forEach(item => {
    const projects = item.bindProjects.length > 0 ? item.bindProjects : ['未绑定项目']

    projects.forEach(projectName => {
      if (!groupMap.has(projectName)) {
        groupMap.set(projectName, [])
      }
      groupMap.get(projectName).push(item)
    })
  })

  return Array.from(groupMap.entries()).map(([projectName, items]) => ({
    projectName,
    items
  }))
})

function toggleProjectGroup(projectName) {
  collapsedProjects[projectName] = !collapsedProjects[projectName]
}

function isProjectCollapsed(projectName) {
  return Boolean(collapsedProjects[projectName])
}

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
  hardwareForm.changeDocFile = null
  hardwareForm.changeDocFileId = 0
  hardwareForm.changeDocFileName = ''
  hardwareForm.changeDocFileUrl = ''
  hardwareForm.changeDocFileContentType = ''
  hardwareForm.changeDocFileData = ''
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
  hardwareForm.changeDocFile = null
  hardwareForm.changeDocFileId = item.changeDocFileId || 0
  hardwareForm.changeDocFileName = item.changeDocFileName || ''
  hardwareForm.changeDocFileUrl = item.changeDocFileUrl || ''
  hardwareForm.changeDocFileContentType = ''
  hardwareForm.changeDocFileData = ''

  showEditDialog.value = true
}

async function handleHardwareFileChange(event) {
  const file = event.target.files[0]
  if (!file) return

  if (!isChangeDocumentFile(file)) {
    alert('只能上传 Word 或 TXT 格式的硬件更改文档')
    event.target.value = ''
    return
  }

  hardwareForm.changeDocFile = file
  hardwareForm.changeDocFileName = file.name
  try {
    const payload = await buildUploadFilePayload(file)
    hardwareForm.changeDocFileId = payload.fileId
    hardwareForm.changeDocFileUrl = payload.fileUrl
    hardwareForm.changeDocFileContentType = payload.fileContentType
    hardwareForm.changeDocFileData = payload.fileData
  } catch (err) {
    alert('读取硬件更改文档失败，请重新选择')
  }
}

function isChangeDocumentFile(file) {
  const name = (file?.name || '').toLowerCase()
  return name.endsWith('.doc') || name.endsWith('.docx') || name.endsWith('.txt')
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
    changeDocFileId: hardwareForm.changeDocFileId,
    changeDocFileName: hardwareForm.changeDocFileName || '',
    fileContentType: hardwareForm.changeDocFileContentType,
    fileData: hardwareForm.changeDocFileData,
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
    if (err.code === 'ECONNABORTED') {
      alert('保存硬件版本超时：文档较大或网络较慢，请稍后重试')
      return
    }
    alert('保存硬件版本失败，请检查后端接口')
  }
}

function viewHardware(item) {
  selectedHardware.value = item
}

function openDocumentUploadDialog(item) {
  currentDocumentHardware.value = item

  documentForm.file = null
  documentForm.fileId = 0
  documentForm.fileName = ''
  documentForm.fileUrl = ''
  documentForm.fileContentType = ''
  documentForm.fileData = ''
  documentForm.remark = ''

  showDocumentDialog.value = true
}

async function handleDocumentFileChange(event) {
  const file = event.target.files[0]
  if (!file) return

  if (!isChangeDocumentFile(file)) {
    alert('只能上传 Word 或 TXT 格式的硬件更改文档')
    event.target.value = ''
    return
  }

  documentForm.file = file
  documentForm.fileName = file.name
  try {
    const payload = await buildUploadFilePayload(file)
    documentForm.fileId = payload.fileId
    documentForm.fileUrl = payload.fileUrl
    documentForm.fileContentType = payload.fileContentType
    documentForm.fileData = payload.fileData
  } catch (err) {
    alert('读取硬件更改文档失败，请重新选择')
  }
}

async function saveDocumentFile() {
  if (!currentDocumentHardware.value) return

  if (!documentForm.file) {
    alert('请选择 Word 或 TXT 格式的硬件更改文档')
    return
  }

  const payload = {
    changeDocFileId: documentForm.fileId,
    changeDocFileName: documentForm.fileName,
    fileContentType: documentForm.fileContentType,
    fileData: documentForm.fileData,
    remark: documentForm.remark || ''
  }

  try {
    const res = await uploadHardwareDocument(currentDocumentHardware.value.id, payload)
    const result = getResponseData(res)

    console.log('上传硬件更改文档返回：', result)

    if (result.code === 200) {
      alert(`硬件版本【${currentDocumentHardware.value.hardwareVersion}】更改文档已上传`)
      showDocumentDialog.value = false
      await loadHardwareVersions()
    } else {
      alert(result.msg || '上传硬件更改文档失败')
    }
  } catch (err) {
    console.error('上传硬件更改文档失败：', err)
    if (err.code === 'ECONNABORTED') {
      alert('上传硬件更改文档超时：文件较大或网络较慢，请稍后重试')
      return
    }
    alert('上传硬件更改文档失败，请检查后端接口')
  }
}

function downloadDocument(item) {
  if (!item.changeDocFileUrl) {
    alert('当前文件暂无可下载内容')
    return
  }

  const link = document.createElement('a')
  link.href = item.changeDocFileUrl
  link.download = item.changeDocFileName || `${item.hardwareVersion}.docx`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

async function deleteHardware(item) {
  if (!item?.id) return
  if (!confirm(`确认删除硬件版本【${item.hardwareVersion}】吗？`)) return

  try {
    const res = await deleteHardwareVersion(item.id)
    const result = res?.data || res
    if (result.code !== 200) {
      alert(result.msg || '删除硬件版本失败')
      return
    }
    alert('删除硬件版本成功')
    await loadHardwareVersions()
  } catch (err) {
    console.error('删除硬件版本失败：', err)
    alert(err.response?.data || '删除硬件版本失败，请检查后端接口')
  }
}

function exportHardwareVersions() {
  const header = [
    '硬件版本',
    '终端类型',
    '绑定项目',
    '版本状态',
    '负责人',
    '更新时间',
    '硬件更改文档'
  ]

  const rows = hardwareVersionList.value.map(item => [
    item.hardwareVersion,
    item.deviceType,
    item.bindProjects.join('、'),
    getStatusText(item.status),
    item.owner,
    item.updateTime,
    item.changeDocFileName || '未上传'
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

.project-version-groups {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.project-version-group {
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 14px;
  overflow: hidden;
}

.project-group-header {
  width: 100%;
  height: 48px;
  padding: 0 16px;
  border: none;
  border-bottom: 1px solid #1e293b;
  background: #020617;
  color: #e2e8f0;
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  text-align: left;
}

.project-group-header:hover {
  background: #0b1120;
}

.fold-icon {
  width: 18px;
  color: #60a5fa;
  font-size: 18px;
  font-weight: 800;
  line-height: 1;
}

.project-group-title {
  flex: 1;
  font-size: 14px;
  font-weight: 800;
}

.project-group-count {
  color: #94a3b8;
  font-size: 12px;
}

.project-version-group .table-card {
  border: none;
  border-radius: 0;
}

.empty-card {
  padding: 28px 16px;
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 14px;
  color: #94a3b8;
  text-align: center;
  font-size: 13px;
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
  white-space: pre-wrap;
  text-align: left;
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

.text-btn.danger {
  color: #f87171;
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
  white-space: pre-wrap;
  text-align: left;
}

.document-form {
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
