<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>软件版本管理</h1>
      </div>

      <button v-if="canUseAction('software:create')" class="primary-btn" @click="openCreateDialog">
        新增软件版本
      </button>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索软件版本 / 软件描述 / 项目 / 终端 / 硬件版本"
      />

      <select v-model="filters.projectName">
        <option value="">全部项目</option>
        <option
          v-for="project in projectOptions"
          :key="project"
          :value="project"
        >
          {{ project }}
        </option>
      </select>

      <select v-model="filters.deviceType">
        <option value="">全部终端</option>
        <option
          v-for="type in deviceTypeOptions"
          :key="type"
          :value="type"
        >
          {{ type }}
        </option>
      </select>

      <button class="query-btn" @click="loadSoftwareVersions">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <!-- 数据表格 -->
    <div class="table-card">
      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>软件版本</th>
              <th>软件描述</th>
              <th>适配项目</th>
              <th>适配终端</th>
              <th>适配硬件版本</th>
              <th>负责人</th>
              <th>发布日期</th>
              <th>状态</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="item in filteredSoftwareList" :key="item.id">
              <td>
                <button class="version-link" @click="openDownloadPage(item)">
                  {{ item.softwareVersion }}
                </button>
              </td>

              <td>
                <div class="software-desc" :title="item.businessDesc">
                  {{ item.businessDesc }}
                </div>
              </td>

              <td>
                <div class="project-list">
                  <span
                    v-for="project in getSoftwareProjects(item)"
                    :key="project"
                    class="project-tag"
                  >
                    {{ project }}
                  </span>
                </div>
              </td>

              <td>
                <span class="device-tag">
                  {{ item.deviceType }}
                </span>
              </td>

              <td class="hardware-version-cell">
                <span class="hardware-tag" :title="item.hardwareVersion">
                  {{ item.hardwareVersion }}
                </span>
              </td>

              <td>{{ item.owner }}</td>

              <td class="muted">
                {{ item.releaseDate || '-' }}
              </td>

              <td>
                <span class="status-tag" :class="item.softwareStatus">
                  {{ getSoftwareStatusText(item.softwareStatus) }}
                </span>
              </td>

              <td class="operation-col">
                <div class="action-group">
                  <button class="text-btn" @click="viewSoftware(item)">
                    查看
                  </button>

                  <button v-if="canUseAction('software:update')" class="text-btn blue" @click="openEditDialog(item)">
                    修改
                  </button>

                  <button
                    v-if="canUseAction('software:release') && item.softwareStatus !== 'released'"
                    class="text-btn green"
                    @click="releaseSoftware(item)"
                  >
                    发布
                  </button>

                  <button v-if="canUseAction('software:download')" class="text-btn green" @click="openDownloadPage(item)">
                    下载
                  </button>

                  <button v-if="canUseAction('software:delete')" class="text-btn red" @click="deleteSoftware(item)">
                    删除
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="table-footer">
        共 {{ filteredSoftwareList.length }} 条软件版本记录
      </div>
    </div>

    <!-- 新增 / 修改软件版本弹窗 -->
    <div v-if="showEditDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>{{ editMode === 'create' ? '新增软件版本' : '修改软件版本' }}</h3>
          <button @click="showEditDialog = false">×</button>
        </div>

        <div class="form-grid">
          <label>
            软件版本号
            <input
              v-model="softwareForm.softwareVersion"
              placeholder="例如：SW-CRRC-ARG-DACU.V1.0.0"
            />
          </label>

          <label>
            适配终端
            <select v-model="softwareForm.deviceType" @change="onDeviceTypeChange">
              <option value="">请选择终端</option>
              <option
                v-for="type in availableDeviceTypes"
                :key="type"
                :value="type"
              >
                {{ type }}
              </option>
            </select>
          </label>

          <label>
            适配硬件版本
            <select v-model="softwareForm.hardwareVersion">
              <option value="">请选择硬件版本</option>
              <option
                v-for="item in availableHardwareVersions"
                :key="item.id"
                :value="item.hardwareVersion"
              >
                {{ item.hardwareVersion }} / {{ item.deviceType }}
              </option>
            </select>
          </label>

          <label>
            当前负责人
            <input
              :value="editMode === 'create' ? currentUserName : softwareForm.owner"
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
                  v-model="softwareForm.projectNames"
                  @change="onProjectSelectionChange"
                />
                <span>{{ project }}</span>
              </label>
            </div>
          </label>

          <label class="full-row">
            下载网页地址
            <input
              v-model="softwareForm.downloadUrl"
              placeholder="例如：http://bc.zycoo.com:8050/files/Argentina/DACU/"
            />
          </label>

          <label class="full-row">
            软件描述
            <textarea
              v-model="softwareForm.businessDesc"
              placeholder="例如：实现阿根廷项目 DACU 广播控制盒的人工广播、OCC广播、PAD广播、紧急广播等业务功能"
            ></textarea>
          </label>

          <label class="full-row">
            版本说明
            <textarea
              v-model="softwareForm.description"
              placeholder="例如：修复音频播放异常、优化SIP注册流程、适配指定硬件版本"
            ></textarea>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showEditDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="saveSoftwareVersion">
            保存
          </button>
        </div>
      </div>
    </div>

    <!-- 查看详情弹窗 -->
    <div v-if="selectedSoftware" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>软件版本详情</h3>
          <button @click="selectedSoftware = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>软件版本</span>
            <strong>{{ selectedSoftware.softwareVersion }}</strong>
          </div>

          <div>
            <span>适配项目</span>
            <strong>{{ getSoftwareProjects(selectedSoftware).join('、') }}</strong>
          </div>

          <div>
            <span>适配终端</span>
            <strong>{{ selectedSoftware.deviceType }}</strong>
          </div>

          <div>
            <span>适配硬件版本</span>
            <strong>{{ selectedSoftware.hardwareVersion }}</strong>
          </div>

          <div>
            <span>负责人</span>
            <strong>{{ selectedSoftware.owner }}</strong>
          </div>

          <div>
            <span>发布日期</span>
            <strong>{{ selectedSoftware.releaseDate || '-' }}</strong>
          </div>

          <div>
            <span>状态</span>
            <strong>{{ getSoftwareStatusText(selectedSoftware.softwareStatus) }}</strong>
          </div>

          <div>
            <span>下载网页</span>
            <button class="inline-link" @click="openDownloadPage(selectedSoftware)">
              打开下载网页
            </button>
          </div>
        </div>

        <div class="remark-card">
          <span>软件描述</span>
          <p>{{ selectedSoftware.businessDesc || '暂无说明' }}</p>
        </div>

        <div class="remark-card">
          <span>版本说明</span>
          <p>{{ selectedSoftware.description || '暂无说明' }}</p>
        </div>

        <div class="dialog-footer">
          <button class="primary-btn" @click="selectedSoftware = null">
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
import { getHardwareVersions } from '@/api/hardware'

import {
  getSoftwareVersions,
  createSoftwareVersion,
  updateSoftwareVersion,
  releaseSoftwareVersion,
  discardSoftwareVersion,
  deleteSoftwareVersion
} from '@/api/software'

const currentUserName = ref(
  localStorage.getItem('username') ||
  localStorage.getItem('accountName') ||
  localStorage.getItem('realName') ||
  '当前用户'
)

const filters = reactive({
  keyword: '',
  projectName: '',
  deviceType: ''
})

const showEditDialog = ref(false)
const selectedSoftware = ref(null)
const currentEditSoftware = ref(null)
const editMode = ref('create')

const projectOptions = ref([])
const projectMap = ref({})

const deviceTypeOptions = [
  '广播控制盒',
  '客室解码板',
  '编码板',
  '乘客报警器',
  '司机室话筒',
  '功放模块',
  '司机室广播控制盒',
  '解码板',
  '功放板',
  '噪声检测器'
]

const hardwareVersionList = ref([])

const softwareForm = reactive({
  softwareVersion: '',
  projectNames: [],
  deviceType: '',
  hardwareVersion: '',
  owner: '',
  downloadUrl: '',
  businessDesc: '',
  description: ''
})

const softwareVersionList = ref([])

onMounted(async () => {
  await loadProjects()
  await loadHardwareVersions()
  await loadSoftwareVersions()
})

function getResponseData(res) {
  if (res && res.data) return res.data
  return res
}

function formatDate(value) {
  if (!value) return ''

  if (typeof value === 'string') {
    return value.slice(0, 10)
  }

  if (value.Time) {
    return String(value.Time).slice(0, 10)
  }

  return value
}

function backendSoftwareStatusToFrontend(status) {
  const map = {
    草稿: 'draft',
    已发布: 'released',
    已废弃: 'discarded',
    draft: 'draft',
    released: 'released',
    discarded: 'discarded'
  }

  return map[status] || status || 'draft'
}

function frontendSoftwareStatusToBackend(status) {
  const map = {
    draft: '草稿',
    released: '已发布',
    discarded: '已废弃'
  }

  return map[status] || status || '草稿'
}

function getSoftwareStatusText(status) {
  const map = {
    draft: '草稿',
    released: '已发布',
    discarded: '已废弃'
  }

  return map[status] || status
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

    if (result.code !== 200) {
      alert(result.msg || '加载硬件版本失败')
      return
    }

    hardwareVersionList.value = (result.data || []).map(item => normalizeHardware(item))
  } catch (err) {
    console.error('加载硬件版本失败：', err)
    alert('加载硬件版本失败')
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
    owner: item.owner || item.ownerName || '',
    updateTime: formatDate(item.updatedAt || item.createdAt),
    description: item.description || ''
  }
}

async function loadSoftwareVersions() {
  try {
    const res = await getSoftwareVersions()
    const result = getResponseData(res)

    console.log('软件版本列表返回：', result)

    if (result.code !== 200) {
      alert(result.msg || '加载软件版本失败')
      return
    }

    softwareVersionList.value = (result.data || []).map(item => normalizeSoftware(item))
  } catch (err) {
    console.error('加载软件版本失败：', err)
    alert('加载软件版本失败，请检查后端接口')
  }
}

function normalizeSoftware(item) {
  let projectNames = []

  if (Array.isArray(item.projectNames)) {
    projectNames = item.projectNames
  } else if (Array.isArray(item.bindProjects)) {
    projectNames = item.bindProjects
  } else if (item.projectName) {
    projectNames = [item.projectName]
  } else if (item.projectId) {
    projectNames = [findProjectName(item.projectId)]
  }

  return {
    id: item.id,
    projectId: item.projectId || 0,
    projectNames,
    softwareVersion: item.softwareVersion || '',
    deviceType: item.deviceType || '',
    hardwareId: item.hardwareId || 0,
    hardwareVersion: item.hardwareVersion || '',
    ownerId: item.ownerId || 1,
    owner: item.owner || item.ownerName || '未分配',
    releaseDate: formatDate(item.releaseDate || item.createdAt),
    downloadUrl: item.downloadUrl || '',
    businessDesc: item.businessDesc || '',
    description: item.description || '',
    softwareStatus: backendSoftwareStatusToFrontend(item.softwareStatus || item.status)
  }
}

const availableHardwareVersions = computed(() => {
  return hardwareVersionList.value.filter(item => {
    const projectMatch =
      softwareForm.projectNames.length === 0 ||
      item.bindProjects.length === 0 ||
      item.bindProjects.some(project =>
        softwareForm.projectNames.includes(project)
      )

    const deviceMatch =
      !softwareForm.deviceType ||
      item.deviceType === softwareForm.deviceType

    return projectMatch && deviceMatch
  })
})

const availableDeviceTypes = computed(() => {
  if (softwareForm.projectNames.length === 0) {
    return deviceTypeOptions
  }

  const types = hardwareVersionList.value
    .filter(item => {
      return (
        item.bindProjects.length === 0 ||
        item.bindProjects.some(project =>
          softwareForm.projectNames.includes(project)
        )
      )
    })
    .map(item => item.deviceType)

  return [...new Set([...types, ...deviceTypeOptions])]
})

const filteredSoftwareList = computed(() => {
  return softwareVersionList.value.filter(item => {
    const itemProjects = getSoftwareProjects(item)

    const keywordMatch =
      !filters.keyword ||
      item.softwareVersion.includes(filters.keyword) ||
      itemProjects.some(project => project.includes(filters.keyword)) ||
      item.deviceType.includes(filters.keyword) ||
      item.hardwareVersion.includes(filters.keyword) ||
      item.businessDesc.includes(filters.keyword) ||
      item.description.includes(filters.keyword)

    const projectMatch =
      !filters.projectName || itemProjects.includes(filters.projectName)

    const deviceTypeMatch =
      !filters.deviceType || item.deviceType === filters.deviceType

    return keywordMatch && projectMatch && deviceTypeMatch
  })
})

function getSoftwareProjects(item) {
  if (Array.isArray(item.projectNames) && item.projectNames.length > 0) {
    return item.projectNames
  }

  if (item.projectName) {
    return [item.projectName]
  }

  return []
}

function resetFilters() {
  filters.keyword = ''
  filters.projectName = ''
  filters.deviceType = ''
}

function onProjectSelectionChange() {
  const currentDeviceValid = availableDeviceTypes.value.includes(softwareForm.deviceType)

  if (!currentDeviceValid) {
    softwareForm.deviceType = ''
    softwareForm.hardwareVersion = ''
  } else {
    const currentHardwareValid = availableHardwareVersions.value.some(
      item => item.hardwareVersion === softwareForm.hardwareVersion
    )

    if (!currentHardwareValid) {
      softwareForm.hardwareVersion = ''
    }
  }
}

function onDeviceTypeChange() {
  const currentHardwareValid = availableHardwareVersions.value.some(
    item => item.hardwareVersion === softwareForm.hardwareVersion
  )

  if (!currentHardwareValid) {
    softwareForm.hardwareVersion = ''
  }
}

function resetSoftwareForm() {
  softwareForm.softwareVersion = ''
  softwareForm.projectNames = []
  softwareForm.deviceType = ''
  softwareForm.hardwareVersion = ''
  softwareForm.owner = ''
  softwareForm.downloadUrl = ''
  softwareForm.businessDesc = ''
  softwareForm.description = ''
}

function openCreateDialog() {
  editMode.value = 'create'
  currentEditSoftware.value = null

  resetSoftwareForm()

  showEditDialog.value = true
}

function openEditDialog(item) {
  editMode.value = 'edit'
  currentEditSoftware.value = item

  softwareForm.softwareVersion = item.softwareVersion
  softwareForm.projectNames = [...getSoftwareProjects(item)]
  softwareForm.deviceType = item.deviceType
  softwareForm.hardwareVersion = item.hardwareVersion
  softwareForm.owner = item.owner || currentUserName.value
  softwareForm.downloadUrl = item.downloadUrl
  softwareForm.businessDesc = item.businessDesc
  softwareForm.description = item.description

  showEditDialog.value = true
}

async function saveSoftwareVersion() {
  if (!softwareForm.softwareVersion) {
    alert('请输入软件版本号')
    return
  }

  if (softwareForm.projectNames.length === 0) {
    alert('请至少选择一个绑定项目')
    return
  }

  if (!softwareForm.deviceType) {
    alert('请选择适配终端')
    return
  }

  if (!softwareForm.hardwareVersion) {
    alert('请选择适配硬件版本')
    return
  }

  if (!softwareForm.businessDesc) {
    alert('请填写软件描述')
    return
  }

  if (!softwareForm.downloadUrl) {
    alert('请填写软件版本下载网页地址')
    return
  }

  const projectIds = getProjectIdsByNames(softwareForm.projectNames)
  const hardware = hardwareVersionList.value.find(
    item => item.hardwareVersion === softwareForm.hardwareVersion
  )

  if (projectIds.length === 0) {
    alert('没有找到项目ID，请重新选择项目')
    return
  }

  const payload = {
    projectId: projectIds[0],
    projectIds,
    projectNames: softwareForm.projectNames,
    softwareVersion: softwareForm.softwareVersion,
    deviceType: softwareForm.deviceType,
    hardwareId: hardware?.id || 0,
    hardwareVersion: softwareForm.hardwareVersion,
    ownerId: 1,
    owner: editMode.value === 'create'
      ? currentUserName.value
      : (softwareForm.owner || currentUserName.value),
    ownerName: editMode.value === 'create'
      ? currentUserName.value
      : (softwareForm.owner || currentUserName.value),
    releaseDate: new Date().toISOString().slice(0, 10),
    downloadUrl: softwareForm.downloadUrl,
    businessDesc: softwareForm.businessDesc,
    description: softwareForm.description,
    softwareStatus: frontendSoftwareStatusToBackend('draft'),
    status: frontendSoftwareStatusToBackend('draft')
  }

  try {
    let res

    if (editMode.value === 'create') {
      res = await createSoftwareVersion(payload)
    } else {
      res = await updateSoftwareVersion(currentEditSoftware.value.id, payload)
    }

    const result = getResponseData(res)

    console.log('保存软件版本返回：', result)

    if (result.code === 200) {
      alert(editMode.value === 'create' ? '新增软件版本成功' : '修改软件版本成功')
      showEditDialog.value = false
      await loadSoftwareVersions()
    } else {
      alert(result.msg || '保存软件版本失败')
    }
  } catch (err) {
    console.error('保存软件版本失败：', err)
    alert('保存软件版本失败，请检查后端接口')
  }
}

function viewSoftware(item) {
  selectedSoftware.value = item
}

function openDownloadPage(item) {
  if (!item.downloadUrl) {
    alert('当前软件版本没有配置下载网页地址')
    return
  }

  window.open(item.downloadUrl, '_blank')
}

async function releaseSoftware(item) {
  if (!item || !item.id) {
    alert('发布失败：没有拿到软件版本ID')
    return
  }

  const ok = confirm(`确认发布软件版本【${item.softwareVersion}】吗？`)
  if (!ok) return

  try {
    const res = await releaseSoftwareVersion(item.id)
    const result = getResponseData(res)

    console.log('发布软件版本返回：', result)

    if (result.code === 200) {
      alert('发布成功')
      await loadSoftwareVersions()
    } else {
      alert(result.msg || '发布失败')
    }
  } catch (err) {
    console.error('发布软件版本失败：', err)
    alert('发布失败，请检查后端接口')
  }
}

async function discardSoftware(item) {
  if (!item || !item.id) {
    alert('废弃失败：没有拿到软件版本ID')
    return
  }

  const ok = confirm(`确认废弃软件版本【${item.softwareVersion}】吗？`)
  if (!ok) return

  try {
    const res = await discardSoftwareVersion(item.id)
    const result = getResponseData(res)

    console.log('废弃软件版本返回：', result)

    if (result.code === 200) {
      alert('废弃成功')
      await loadSoftwareVersions()
    } else {
      alert(result.msg || '废弃失败')
    }
  } catch (err) {
    console.error('废弃软件版本失败：', err)
    alert('废弃失败，请检查后端接口')
  }
}

async function deleteSoftware(item) {
  if (!item || !item.id) {
    alert('删除失败：没有拿到软件版本ID')
    return
  }

  const ok = confirm(`确认删除软件版本【${item.softwareVersion}】吗？`)
  if (!ok) return

  try {
    const res = await deleteSoftwareVersion(item.id)
    const result = getResponseData(res)

    console.log('删除软件版本返回：', result)

    if (result.code === 200) {
      alert('删除成功')
      selectedSoftware.value = null
      await loadSoftwareVersions()
    } else {
      alert(result.msg || '删除失败')
    }
  } catch (err) {
    console.error('删除软件版本失败：', err)
    alert('删除失败，请检查后端接口')
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
  grid-template-columns: 1.4fr 200px 180px 90px 90px;
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

.table-card {
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 14px;
  overflow: hidden;
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

.table-card table {
  width: 100%;
  min-width: 1380px;
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

.table-card th:nth-child(1),
.table-card td:nth-child(1) {
  width: 190px;
}

.table-card th:nth-child(2),
.table-card td:nth-child(2) {
  width: 280px;
}

.table-card th:nth-child(3),
.table-card td:nth-child(3) {
  width: 180px;
}

.table-card th:nth-child(4),
.table-card td:nth-child(4) {
  width: 140px;
}

.table-card th:nth-child(5),
.table-card td:nth-child(5) {
  width: 210px;
}

.table-card th:nth-child(6),
.table-card td:nth-child(6) {
  width: 100px;
}

.table-card th:nth-child(7),
.table-card td:nth-child(7) {
  width: 110px;
}

.table-card th:nth-child(8),
.table-card td:nth-child(8) {
  width: 90px;
}

.version-link {
  border: none;
  background: transparent;
  color: #60a5fa;
  font-size: 13px;
  font-weight: 800;
  cursor: pointer;
  padding: 0;
  text-align: left;
}

.version-link:hover {
  color: #93c5fd;
  text-decoration: underline;
}

.software-desc {
  color: #cbd5e1;
  font-size: 12px;
  line-height: 1.5;
  max-width: 260px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.project-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.project-tag,
.device-tag,
.hardware-tag,
.status-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
}

.project-tag {
  background: #1d4ed833;
  color: #60a5fa;
}

.device-tag {
  background: #0f766e33;
  color: #5eead4;
}

.hardware-tag {
  background: #9333ea33;
  color: #c084fc;
  display: block;
  width: 100%;
  max-width: 170px;
  box-sizing: border-box;
}

.status-tag.draft {
  background: #47556933;
  color: #94a3b8;
}

.status-tag.released {
  background: #16a34a33;
  color: #4ade80;
}

.status-tag.discarded {
  background: #dc262633;
  color: #f87171;
}

.hardware-version-cell {
  overflow: hidden;
  max-width: 0;
}

.hardware-version-cell .hardware-tag {
  display: block;
  width: 100%;
  max-width: 170px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.muted {
  color: #94a3b8 !important;
}

.operation-col {
  width: 360px;
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

.text-btn.red {
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
  width: 760px;
  max-width: 100%;
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
}

.inline-link {
  border: none;
  background: transparent;
  color: #60a5fa;
  cursor: pointer;
  padding: 0;
  font-size: 14px;
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

  .table-card table {
    min-width: 1380px;
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
