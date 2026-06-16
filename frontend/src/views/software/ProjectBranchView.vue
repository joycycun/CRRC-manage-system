<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>项目分支管理</h1>
      </div>

      <button v-if="canUseAction('branch:create')" class="primary-btn" @click="openCreateDialog">
        新增项目分支
      </button>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索项目名称 / 分支名称 / 负责人 / Clone地址"
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

      <button class="query-btn" @click="loadBranches">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <!-- 数据表格 -->
    <div class="table-card">
      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>项目分支</th>
              <th>绑定项目</th>
              <th>终端类型</th>
              <th>软件负责人</th>
              <th>创建时间</th>
              <th>Clone 地址</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="item in filteredBranchList" :key="item.id">
              <td>
                <button class="branch-link" @click="openCloneUrl(item)">
                  {{ item.branchName }}
                </button>
              </td>

              <td>
                <span class="project-tag">
                  {{ item.projectName }}
                </span>
              </td>

              <td>
                <span class="device-tag">
                  {{ item.deviceType }}
                </span>
              </td>

              <td>
                <span class="owner-text" :title="item.owner">
                  {{ item.owner }}
                </span>
              </td>

              <td class="muted">
                {{ item.createTime }}
              </td>

              <td>
                <div class="clone-url" :title="item.cloneUrl">
                  {{ item.cloneUrl }}
                </div>
              </td>

              <td class="operation-col">
                <div class="action-group">
                  <button class="text-btn" @click="viewBranch(item)">
                    查看
                  </button>

                  <button v-if="canUseAction('branch:update')" class="text-btn blue" @click="openEditDialog(item)">
                    修改
                  </button>

                  <button v-if="canUseAction('branch:download')" class="text-btn green" @click="openCloneUrl(item)">
                    打开
                  </button>

                  <button v-if="canUseAction('branch:download')" class="text-btn yellow" @click="copyCloneUrl(item)">
                    复制
                  </button>

                  <button v-if="canUseAction('branch:delete')" class="text-btn red" @click="deleteBranchRecord(item)">
                    删除
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="table-footer">
        共 {{ filteredBranchList.length }} 条项目分支记录
      </div>
    </div>

    <!-- 新增 / 修改项目分支弹窗 -->
    <div v-if="showEditDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>{{ editMode === 'create' ? '新增项目分支' : '修改项目分支' }}</h3>
          <button @click="showEditDialog = false">×</button>
        </div>

        <div class="form-grid">
          <label>
            绑定项目
            <select v-model="branchForm.projectName">
              <option value="">请选择项目</option>
              <option
                v-for="project in projectOptions"
                :key="project"
                :value="project"
              >
                {{ project }}
              </option>
            </select>
          </label>

          <label>
            终端类型
            <select v-model="branchForm.deviceType">
              <option value="">请选择终端</option>
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
            分支名称
            <input
              v-model="branchForm.branchName"
              placeholder="例如：X10-Series_PA_Intercom"
            />
          </label>

          <label>
            创建时间
            <input
              v-model="branchForm.createTime"
              type="date"
            />
          </label>

          <label class="full-row">
            软件负责人
            <input
              :value="editMode === 'create' ? currentUserName : branchForm.owner"
              disabled
            />
          </label>

          <label class="full-row">
            Clone 地址
            <input
              v-model="branchForm.cloneUrl"
              placeholder="例如：http://bc.zycoo.com:3000/speaker/X10-Series_PA_Intercom.git"
            />
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showEditDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="saveBranch">
            保存
          </button>
        </div>
      </div>
    </div>

    <!-- 查看详情弹窗 -->
    <div v-if="selectedBranch" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>项目分支详情</h3>
          <button @click="selectedBranch = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>分支名称</span>
            <strong>{{ selectedBranch.branchName }}</strong>
          </div>

          <div>
            <span>绑定项目</span>
            <strong>{{ selectedBranch.projectName }}</strong>
          </div>

          <div>
            <span>终端类型</span>
            <strong>{{ selectedBranch.deviceType }}</strong>
          </div>

          <div>
            <span>软件负责人</span>
            <strong>{{ selectedBranch.owner }}</strong>
          </div>

          <div>
            <span>创建时间</span>
            <strong>{{ selectedBranch.createTime }}</strong>
          </div>

          <div class="full-detail-row">
            <span>Clone 地址</span>
            <button v-if="canUseAction('branch:download')" class="inline-link" @click="openCloneUrl(selectedBranch)">
              {{ selectedBranch.cloneUrl }}
            </button>
          </div>
        </div>

        <div class="dialog-footer">
          <button v-if="canUseAction('branch:download')" class="reset-btn" @click="copyCloneUrl(selectedBranch)">
            复制 Clone 地址
          </button>

          <button v-if="canUseAction('branch:download')" class="primary-btn" @click="openCloneUrl(selectedBranch)">
            打开 Clone 地址
          </button>

          <button class="primary-btn" @click="selectedBranch = null">
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
  getBranches,
  createBranch,
  updateBranch,
  deleteBranch
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
const selectedBranch = ref(null)
const currentEditBranch = ref(null)
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

const branchForm = reactive({
  projectName: '',
  deviceType: '',
  branchName: '',
  owner: '',
  createTime: '',
  cloneUrl: ''
})

const branchList = ref([])

onMounted(async () => {
  await loadProjects()
  await loadBranches()
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

function findProjectName(projectId) {
  const found = Object.entries(projectMap.value).find(([, id]) => Number(id) === Number(projectId))
  return found ? found[0] : `项目ID-${projectId}`
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

async function loadBranches() {
  try {
    const res = await getBranches()
    const result = getResponseData(res)

    console.log('项目分支列表返回：', result)

    if (result.code !== 200) {
      alert(result.msg || '加载项目分支失败')
      return
    }

    branchList.value = (result.data || []).map(item => normalizeBranch(item))
  } catch (err) {
    console.error('加载项目分支失败：', err)
    alert('加载项目分支失败，请检查后端接口')
  }
}

function normalizeBranch(item) {
  return {
    id: item.id,
    projectId: item.projectId || 0,
    projectName: item.projectName || findProjectName(item.projectId),
    deviceType: item.deviceType || '',
    ownerId: item.ownerId || 1,
    owner: item.owner || item.ownerName || item.responsibleUser || '未分配',
    createTime: formatDate(item.createTime || item.createdAt),
    branchName: item.branchName || item.repoName || item.name || '',
    cloneUrl: item.cloneUrl || item.repoUrl || item.gitUrl || item.repositoryUrl || ''
    
  }
}

const filteredBranchList = computed(() => {
  return branchList.value.filter(item => {
    const keywordMatch =
      !filters.keyword ||
      item.projectName.includes(filters.keyword) ||
      item.deviceType.includes(filters.keyword) ||
      item.branchName.includes(filters.keyword) ||
      item.owner.includes(filters.keyword) ||
      item.cloneUrl.includes(filters.keyword)

    const projectMatch =
      !filters.projectName || item.projectName === filters.projectName

    const deviceTypeMatch =
      !filters.deviceType || item.deviceType === filters.deviceType

    return keywordMatch && projectMatch && deviceTypeMatch
  })
})

function resetFilters() {
  filters.keyword = ''
  filters.projectName = ''
  filters.deviceType = ''
}

function openCreateDialog() {
  editMode.value = 'create'
  currentEditBranch.value = null

  branchForm.projectName = ''
  branchForm.deviceType = ''
  branchForm.branchName = ''
  branchForm.owner = currentUserName.value
  branchForm.createTime = new Date().toISOString().slice(0, 10)
  branchForm.cloneUrl = ''

  showEditDialog.value = true
}

function openEditDialog(item) {
  editMode.value = 'edit'
  currentEditBranch.value = item

  branchForm.projectName = item.projectName
  branchForm.deviceType = item.deviceType
  branchForm.branchName = item.branchName
  branchForm.owner = item.owner || currentUserName.value
  branchForm.createTime = item.createTime
  branchForm.cloneUrl = item.cloneUrl

  showEditDialog.value = true
}

async function saveBranch() {
  if (!branchForm.projectName) {
    alert('请选择绑定项目')
    return
  }

  if (!branchForm.deviceType) {
    alert('请选择终端类型')
    return
  }

  if (!branchForm.branchName) {
    alert('请输入分支名称')
    return
  }

  if (!branchForm.cloneUrl) {
    alert('请输入 Clone 地址')
    return
  }

  const projectId = projectMap.value[branchForm.projectName]

  if (!projectId) {
    alert('没有找到项目ID，请重新选择项目')
    return
  }

  const payload = {
    projectId,
    projectName: branchForm.projectName,
    deviceType: branchForm.deviceType,
    branchName: branchForm.branchName,
    ownerId: 1,
    owner: editMode.value === 'create'
      ? currentUserName.value
      : (branchForm.owner || currentUserName.value),
    ownerName: editMode.value === 'create'
      ? currentUserName.value
      : (branchForm.owner || currentUserName.value),
    createTime: branchForm.createTime || new Date().toISOString().slice(0, 10),
    cloneUrl: branchForm.cloneUrl
  }

  try {
    let res

    if (editMode.value === 'create') {
      res = await createBranch(payload)
    } else {
      res = await updateBranch(currentEditBranch.value.id, payload)
    }

    const result = getResponseData(res)

    console.log('保存项目分支返回：', result)

    if (result.code === 200) {
      alert(editMode.value === 'create' ? '新增项目分支成功' : '修改项目分支成功')
      showEditDialog.value = false
      await loadBranches()
    } else {
      alert(result.msg || '保存项目分支失败')
    }
  } catch (err) {
    console.error('保存项目分支失败：', err)
    alert('保存项目分支失败，请检查后端接口')
  }
}

function viewBranch(item) {
  selectedBranch.value = item
}

function openCloneUrl(item) {
  if (!item.cloneUrl) {
    alert('当前项目分支没有配置 Clone 地址')
    return
  }

  window.open(item.cloneUrl, '_blank')
}

async function copyCloneUrl(item) {
  if (!item.cloneUrl) {
    alert('当前项目分支没有配置 Clone 地址')
    return
  }

  try {
    await navigator.clipboard.writeText(item.cloneUrl)
    alert('Clone 地址已复制')
  } catch (error) {
    const input = document.createElement('input')
    input.value = item.cloneUrl
    document.body.appendChild(input)
    input.select()
    document.execCommand('copy')
    document.body.removeChild(input)
    alert('Clone 地址已复制')
  }
}

async function deleteBranchRecord(item) {
  if (!item || !item.id) {
    alert('删除失败：没有拿到项目分支ID')
    return
  }

  const ok = confirm(`确认删除项目分支【${item.branchName}】吗？`)
  if (!ok) return

  try {
    const res = await deleteBranch(item.id)
    const result = getResponseData(res)

    console.log('删除项目分支返回：', result)

    if (result.code === 200) {
      alert('删除成功')
      selectedBranch.value = null
      await loadBranches()
    } else {
      alert(result.msg || '删除失败')
    }
  } catch (err) {
    console.error('删除项目分支失败：', err)
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
.form-grid select {
  border: 1px solid #334155;
  border-radius: 8px;
  background: #020617;
  color: #e2e8f0;
  padding: 0 12px;
  outline: none;
  height: 36px;
}

.form-grid input:disabled {
  color: #94a3b8;
  background: #0f172a;
  cursor: not-allowed;
}

.filter-card input::placeholder,
.form-grid input::placeholder {
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
  min-width: 1150px;
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
  width: 180px;
}

.table-card th:nth-child(2),
.table-card td:nth-child(2) {
  width: 160px;
}

.table-card th:nth-child(3),
.table-card td:nth-child(3) {
  width: 140px;
}

.table-card th:nth-child(4),
.table-card td:nth-child(4) {
  width: 120px;
}

.table-card th:nth-child(5),
.table-card td:nth-child(5) {
  width: 120px;
}

.table-card th:nth-child(6),
.table-card td:nth-child(6) {
  width: 280px;
}

.branch-link {
  display: inline-block;
  max-width: 160px;
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
  vertical-align: middle;
}

.branch-link:hover {
  color: #93c5fd;
  text-decoration: underline;
}

.project-tag,
.device-tag {
  display: inline-block;
  max-width: 130px;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  vertical-align: middle;
}

.project-tag {
  background: #1d4ed833;
  color: #60a5fa;
}

.device-tag {
  background: #0f766e33;
  color: #5eead4;
}

.owner-text {
  display: inline-block;
  max-width: 100px;
  color: #e2e8f0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;
}

.clone-url {
  display: block;
  max-width: 260px;
  color: #94a3b8;
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.muted {
  color: #94a3b8 !important;
}

.operation-col {
  width: 320px;
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

.detail-card div {
  background: #020617;
  border: 1px solid #1e293b;
  border-radius: 10px;
  padding: 12px;
}

.full-detail-row {
  grid-column: 1 / -1;
}

.detail-card span {
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
  word-break: break-all;
  text-align: left;
}

.inline-link:hover {
  color: #93c5fd;
  text-decoration: underline;
}

@media (max-width: 960px) {
  .filter-card {
    grid-template-columns: 1fr;
  }

  .table-card table {
    min-width: 1150px;
  }

  .form-grid,
  .detail-card {
    grid-template-columns: 1fr;
  }s

  .full-detail-row {
    grid-column: auto;
  }
}
</style>
