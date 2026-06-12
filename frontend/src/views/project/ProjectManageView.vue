<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>项目立项管理</h1>
      </div>

      <button class="primary-btn" @click="openCreateDialog">
        新增项目
      </button>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索项目名称 / 负责人"
      />

      <select v-model="filters.status">
        <option value="">全部状态</option>
        <option value="draft">草稿</option>
        <option value="submitted">待审核</option>
        <option value="approved">已立项</option>
        <option value="running">进行中</option>
        <option value="closed">已关闭</option>
        <option value="archived">已归档</option>
        <option value="rejected">已驳回</option>
      </select>

      <select v-model="filters.stage">
        <option value="">全部阶段</option>
        <option
          v-for="stage in projectStages"
          :key="stage"
          :value="stage"
        >
          {{ stage }}
        </option>
      </select>

      <button class="query-btn" @click="loadProjects">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <!-- 数据表格 -->
    <div class="table-card">
      <table>
        <thead>
          <tr>
            <th>项目名称</th>
            <th>当前阶段</th>
            <th>项目负责人</th>
            <th>项目状态</th>
            <th>创建时间</th>
            <th class="operation-col">操作</th>
          </tr>
        </thead>

        <tbody>
          <tr v-for="item in filteredProjects" :key="item.id">
            <td>
              <div class="project-name">{{ item.projectName }}</div>
              <div class="project-code">{{ item.projectCode }}</div>
            </td>

            <td>
              <span class="stage-tag" :class="getStageClass(item.stage)">
                {{ item.stage }}
              </span>
            </td>

            <td>{{ item.owner }}</td>

            <td>
              <span class="status-tag" :class="item.status">
                {{ getStatusText(item.status) }}
              </span>
            </td>

            <td class="muted">
              {{ item.createTime }}
            </td>

            <td class="operation-col">
              <button class="text-btn" @click="viewProject(item)">
                查看
              </button>

              <button
                v-if="item.status === 'draft'"
                class="text-btn blue"
                @click="submitProject(item)"
              >
                提交立项
              </button>

              <button
                v-if="item.status === 'submitted'"
                class="text-btn green"
                @click="auditProject(item)"
              >
                审核
              </button>

              <button
                v-if="item.status === 'running' || item.status === 'approved'"
                class="text-btn yellow"
                @click="viewStage(item)"
              >
                阶段
              </button>

              <button
                v-if="item.status === 'running' || item.status === 'approved'"
                class="text-btn red"
                @click="closeProject(item)"
              >
                关闭
              </button>

              <button
                v-if="item.status === 'closed'"
                class="text-btn gray"
                @click="archiveProject(item)"
              >
                归档
              </button>

              <button
                class="text-btn red"
                @click="deleteProject(item)"
              >
                删除
              </button>
            </td>
          </tr>
        </tbody>
      </table>

      <div class="table-footer">
        共 {{ filteredProjects.length }} 条项目记录
      </div>
    </div>

    <!-- 新增项目弹窗 -->
    <div v-if="showCreateDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>新增项目</h3>
          <button @click="showCreateDialog = false">×</button>
        </div>

        <div class="form-grid">
          <label>
            项目名称
            <input v-model="projectForm.projectName" placeholder="请输入项目名称" />
          </label>

          <label>
            项目编号
            <input v-model="projectForm.projectCode" placeholder="例如：CRRC-2026-001" />
          </label>

          <label>
            项目负责人
            <input v-model="projectForm.owner" placeholder="请输入负责人" />
          </label>

          <label>
            当前阶段
            <select v-model="projectForm.stage">
              <option
                v-for="stage in projectStages"
                :key="stage"
                :value="stage"
              >
                {{ stage }}
              </option>
            </select>
          </label>

          <label>
            项目状态
            <select v-model="projectForm.status">
              <option value="draft">草稿</option>
              <option value="submitted">待审核</option>
              <option value="approved">已立项</option>
              <option value="running">进行中</option>
            </select>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showCreateDialog = false">取消</button>
          <button class="primary-btn" @click="createProject">保存</button>
        </div>
      </div>
    </div>

    <!-- 查看详情 / 阶段弹窗 -->
    <div v-if="selectedProject" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>项目详情</h3>
          <button @click="selectedProject = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>项目名称</span>
            <strong>{{ selectedProject.projectName }}</strong>
          </div>

          <div>
            <span>项目编号</span>
            <strong>{{ selectedProject.projectCode }}</strong>
          </div>

          <div>
            <span>负责人</span>
            <strong>{{ selectedProject.owner }}</strong>
          </div>

          <div>
            <span>当前阶段</span>
            <strong>{{ selectedProject.stage }}</strong>
          </div>

          <div>
            <span>项目状态</span>
            <strong>{{ getStatusText(selectedProject.status) }}</strong>
          </div>

          <div>
            <span>创建时间</span>
            <strong>{{ selectedProject.createTime }}</strong>
          </div>
        </div>

        <div class="stage-section">
          <h4>项目阶段流程</h4>

          <div class="stage-list">
            <div
              v-for="stage in projectStages"
              :key="stage"
              class="stage-item"
              :class="{
                active: stage === selectedProject.stage,
                passed: isStagePassed(stage, selectedProject.stage)
              }"
              @click="changeProjectStage(stage)"
            >
              <div class="stage-dot"></div>
              <span>{{ stage }}</span>
            </div>
          </div>
        </div>

        <div class="dialog-footer">
          <button
            v-if="selectedProject.status !== 'closed'"
            class="red-btn"
            @click="closeProject(selectedProject)"
          >
            手动关闭项目
          </button>

          <button class="primary-btn" @click="selectedProject = null">
            关闭弹窗
          </button>


        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'

import {
  getProjects,
  createProject as createProjectApi,
  updateProject as updateProjectApi,
  submitProject as submitProjectApi,
  auditProject as auditProjectApi,
  archiveProject as archiveProjectApi,
  closeProject as closeProjectApi,
  deleteProject as deleteProjectApi
} from '@/api/project'

const filters = reactive({
  keyword: '',
  status: '',
  stage: ''
})

const loading = ref(false)
const showCreateDialog = ref(false)
const selectedProject = ref(null)

const projectForm = reactive({
  projectName: '',
  projectCode: '',
  ownerId: 1,
  owner: '',
  stage: '立项',
  status: 'draft',
  remark: ''
})

const projectStages = [
  '立项',
  '需求书首次确认',
  '硬件检测',
  '软件研发',
  '内部初始测试',
  '内部初始测试是否问题闭环',
  '样机联调',
  '出厂测试',
  '收货'
]

const projectList = ref([])

onMounted(() => {
  loadProjects()
})

function getResponseData(res) {
  if (res && res.data) return res.data
  return res
}

/**
 * 后端项目状态 -> 前端页面状态
 */
function backendStatusToFrontend(status, auditStatus) {
  if (status === '归档') return 'archived'
  if (status === '已关闭') return 'closed'
  if (status === '已完成') return 'finished'

  if (auditStatus === '未提交' || !auditStatus) return 'draft'
  if (auditStatus === '待审核') return 'submitted'
  if (auditStatus === '已驳回') return 'rejected'

  if (auditStatus === '已通过') {
    if (status === '进行中') return 'running'
    return 'approved'
  }

  const statusMap = {
    draft: 'draft',
    submitted: 'submitted',
    approved: 'approved',
    running: 'running',
    closed: 'closed',
    archived: 'archived',
    rejected: 'rejected',
    立项中: 'draft',
    进行中: 'running',
    已关闭: 'closed',
    归档: 'archived'
  }

  return statusMap[status] || 'draft'
}

/**
 * 前端页面状态 -> 后端项目状态
 */
function frontendStatusToBackend(status) {
  const map = {
    draft: '立项中',
    submitted: '立项中',
    approved: '进行中',
    running: '进行中',
    closed: '已关闭',
    archived: '归档',
    rejected: '立项中'
  }

  return map[status] || status || '立项中'
}

function normalizeProject(item) {
  return {
    id: item.id,
    projectName: item.projectName || '',
    projectCode: item.projectCode || '未填写',
    stage: item.stage || '立项',
    ownerId: item.ownerId || 1,
    owner: item.owner || item.ownerName || '未分配',
    ownerName: item.ownerName || item.owner || '未分配',
    status: backendStatusToFrontend(item.status, item.auditStatus),
    backendStatus: item.status || '',
    auditStatus: item.auditStatus || '',
    auditUserName: item.auditUserName || '',
    auditTime: formatDate(item.auditTime),
    submitTime: formatDate(item.submitTime),
    createTime: formatDate(item.createTime || item.createdAt),
    archiveTime: formatDate(item.archiveTime),
    closeTime: formatDate(item.closeTime),
    remark: item.remark || ''
  }
}

function formatDate(value) {
  if (!value) return ''
  if (typeof value === 'string') return value.slice(0, 10)
  return value
}

async function loadProjects() {
  loading.value = true

  try {
    const res = await getProjects()
    const result = getResponseData(res)

    console.log('项目列表返回：', result)

    if (result.code === 200) {
      projectList.value = (result.data || []).map(normalizeProject)
    } else {
      alert(result.msg || '查询项目失败')
    }
  } catch (err) {
    console.error('查询项目失败：', err)
    alert('查询项目失败，请检查后端是否启动')
  } finally {
    loading.value = false
  }
}

const filteredProjects = computed(() => {
  return projectList.value.filter(item => {
    const keywordMatch =
      !filters.keyword ||
      item.projectName.includes(filters.keyword) ||
      item.owner.includes(filters.keyword)

    const statusMatch = !filters.status || item.status === filters.status
    const stageMatch = !filters.stage || item.stage === filters.stage

    return keywordMatch && statusMatch && stageMatch
  })
})

function getStatusText(status) {
  const statusMap = {
    draft: '草稿',
    submitted: '待审核',
    approved: '已立项',
    running: '进行中',
    finished: '已完成',
    closed: '已关闭',
    archived: '已归档',
    rejected: '已驳回'
  }

  return statusMap[status] || status
}

function getStageClass(stage) {
  const stageMap = {
    立项: 'stage-project',
    需求书首次确认: 'stage-requirement',
    硬件检测: 'stage-hardware',
    软件研发: 'stage-software',
    内部初始测试: 'stage-test',
    内部初始测试是否问题闭环: 'stage-issue',
    样机联调: 'stage-debug',
    出厂测试: 'stage-factory',
    收货: 'stage-receipt'
  }

  return stageMap[stage] || 'stage-default'
}

function resetFilters() {
  filters.keyword = ''
  filters.status = ''
  filters.stage = ''
}

function openCreateDialog() {
  projectForm.projectName = ''
  projectForm.projectCode = ''
  projectForm.ownerId = 1
  projectForm.owner = ''
  projectForm.stage = '立项'
  projectForm.status = 'draft'
  projectForm.remark = ''

  showCreateDialog.value = true
}

async function createProject() {
  if (!projectForm.projectName) {
    alert('请输入项目名称')
    return
  }

  const payload = {
    projectName: projectForm.projectName,
    projectCode: projectForm.projectCode || '未填写',
    ownerId: projectForm.ownerId || 1,
    owner: projectForm.owner || '未分配',
    ownerName: projectForm.owner || '未分配',
    stage: projectForm.stage,
    status: frontendStatusToBackend(projectForm.status),
    remark: projectForm.remark || ''
  }

  try {
    const res = await createProjectApi(payload)
    const result = getResponseData(res)

    console.log('新增项目返回：', result)

    if (result.code === 200) {
      alert('新增项目成功')
      showCreateDialog.value = false
      await loadProjects()
    } else {
      alert(result.msg || '新增项目失败')
    }
  } catch (err) {
    console.error('新增项目失败：', err)
    alert('新增项目失败，请检查后端接口')
  }
}

function viewProject(item) {
  selectedProject.value = item
}

function viewStage(item) {
  selectedProject.value = item
}

async function changeProjectStage(stage) {
  if (!selectedProject.value) return

  try {
    const row = selectedProject.value

    const res = await updateProjectApi(row.id, {
      projectName: row.projectName,
      projectCode: row.projectCode,
      ownerId: row.ownerId || 1,
      owner: row.owner,
      ownerName: row.owner,
      stage,
      status: frontendStatusToBackend(row.status),
      remark: row.remark || ''
    })

    const result = getResponseData(res)

    if (result.code === 200) {
      alert('阶段修改成功')
      selectedProject.value.stage = stage
      await loadProjects()
    } else {
      alert(result.msg || '阶段修改失败')
    }
  } catch (err) {
    console.error('阶段修改失败：', err)
    alert('阶段修改失败')
  }
}

function isStagePassed(stage, currentStage) {
  const currentIndex = projectStages.indexOf(currentStage)
  const stageIndex = projectStages.indexOf(stage)

  return stageIndex !== -1 && currentIndex !== -1 && stageIndex < currentIndex
}

/**
 * 提交立项：真正调用后端
 */
async function submitProject(item) {
  try {
    const res = await submitProjectApi(item.id)
    const result = getResponseData(res)

    if (result.code === 200) {
      alert(`项目【${item.projectName}】已提交立项审核`)
      await loadProjects()
    } else {
      alert(result.msg || '提交失败')
    }
  } catch (err) {
    console.error('提交项目失败：', err)
    alert('提交项目失败')
  }
}

/**
 * 审核项目：这里用 confirm 简化
 * 点确定 = 审核通过
 * 点取消 = 驳回
 */
async function auditProject(item) {
  const pass = confirm(`是否审核通过项目【${item.projectName}】？\n确定=通过，取消=驳回`)

  let rejectReason = ''

  if (!pass) {
    rejectReason = prompt('请输入驳回原因') || '未填写驳回原因'
  }

  try {
    const res = await auditProjectApi(item.id, {
      auditUserId: 1,
      auditUserName: '领导',
      auditStatus: pass ? '已通过' : '已驳回',
      rejectReason
    })

    const result = getResponseData(res)

    if (result.code === 200) {
      alert(pass ? '审核通过' : '已驳回')
      await loadProjects()
    } else {
      alert(result.msg || '审核失败')
    }
  } catch (err) {
    console.error('审核项目失败：', err)
    alert('审核项目失败')
  }
}

/**
 * 关闭项目：真正调用后端
 */
async function closeProject(item) {
  const ok = confirm(`确认关闭项目【${item.projectName}】吗？`)
  if (!ok) return

  try {
    const res = await closeProjectApi(item.id)
    const result = getResponseData(res)

    if (result.code === 200) {
      alert(`项目【${item.projectName}】已关闭`)
      selectedProject.value = null
      await loadProjects()
    } else {
      alert(result.msg || '关闭失败')
    }
  } catch (err) {
    console.error('关闭项目失败：', err)
    alert('关闭项目失败')
  }
}

/**
 * 归档项目：真正调用后端
 */
async function archiveProject(item) {
  const ok = confirm(`确认归档项目【${item.projectName}】吗？`)
  if (!ok) return

  try {
    const res = await archiveProjectApi(item.id)
    const result = getResponseData(res)

    if (result.code === 200) {
      alert(`项目【${item.projectName}】已归档`)
      await loadProjects()
    } else {
      alert(result.msg || '归档失败')
    }
  } catch (err) {
    console.error('归档项目失败：', err)
    alert('归档项目失败')
  }
}

/**
 * 删除项目：模板里目前没有按钮，后面需要时可以绑定
 */
async function deleteProject(item) {
  console.log('点击删除按钮：', item)

  const ok = confirm(`确认删除项目【${item.projectName}】吗？`)
  if (!ok) return

  try {
    const res = await deleteProjectApi(item.id)
    const result = getResponseData(res)

    console.log('删除项目返回：', result)

    if (result.code === 200) {
      alert('删除成功')
      await loadProjects()
    } else {
      alert(result.msg || '删除失败')
    }
  } catch (err) {
    console.error('删除项目失败：', err)
    alert('删除项目失败，请检查后端接口')
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

.page-header p {
  margin: 8px 0 0;
  color: #94a3b8;
  font-size: 14px;
}

.primary-btn,
.query-btn,
.reset-btn,
.red-btn {
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

.red-btn {
  border: 1px solid #7f1d1d;
  background: #450a0a;
  color: #fca5a5;
}

.filter-card {
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 14px;
  padding: 16px;
  display: grid;
  grid-template-columns: 1.4fr 180px 240px 90px 90px;
  gap: 12px;
  margin-bottom: 20px;
}

.filter-card input,
.filter-card select,
.form-grid input,
.form-grid select {
  height: 36px;
  border: 1px solid #334155;
  border-radius: 8px;
  background: #020617;
  color: #e2e8f0;
  padding: 0 12px;
  outline: none;
}

.filter-card input::placeholder {
  color: #64748b;
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

.project-name {
  color: #f8fafc;
  font-weight: 700;
}

.project-code {
  margin-top: 4px;
  color: #64748b;
  font-size: 12px;
}

.stage-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
}

.stage-project {
  background: #1d4ed833;
  color: #60a5fa;
}

.stage-requirement {
  background: #0f766e33;
  color: #5eead4;
}

.stage-hardware {
  background: #9333ea33;
  color: #c084fc;
}

.stage-software {
  background: #2563eb33;
  color: #93c5fd;
}

.stage-test {
  background: #d9770633;
  color: #fbbf24;
}

.stage-issue {
  background: #dc262633;
  color: #f87171;
}

.stage-debug {
  background: #7c3aed33;
  color: #c4b5fd;
}

.stage-factory {
  background: #16a34a33;
  color: #4ade80;
}

.stage-receipt {
  background: #05966933;
  color: #34d399;
}

.stage-default {
  background: #47556933;
  color: #94a3b8;
}

.status-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
}

.status-tag.draft {
  background: #47556933;
  color: #94a3b8;
}

.status-tag.submitted {
  background: #d9770633;
  color: #fbbf24;
}

.status-tag.approved {
  background: #16a34a33;
  color: #4ade80;
}

.status-tag.running {
  background: #1d4ed833;
  color: #60a5fa;
}

.status-tag.closed {
  background: #64748b33;
  color: #cbd5e1;
}

.status-tag.archived {
  background: #64748b33;
  color: #94a3b8;
}

.status-tag.rejected {
  background: #dc262633;
  color: #f87171;
}

.muted {
  color: #94a3b8 !important;
}

.operation-col {
  width: 300px;
  text-align: right !important;
}

.text-btn {
  border: none;
  background: transparent;
  color: #cbd5e1;
  font-size: 13px;
  cursor: pointer;
  margin-left: 10px;
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

.text-btn.gray {
  color: #94a3b8;
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
  width: 620px;
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

.stage-section {
  padding: 0 20px 20px;
}

.stage-section h4 {
  margin: 0 0 14px;
  font-size: 15px;
}

.stage-list {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.stage-item {
  background: #020617;
  border: 1px solid #1e293b;
  border-radius: 10px;
  padding: 12px 8px;
  text-align: center;
  color: #64748b;
  font-size: 12px;
  cursor: pointer;
}

.stage-item:hover {
  border-color: #334155;
  color: #cbd5e1;
}

.stage-item.passed {
  border-color: #14532d;
  background: #052e1633;
  color: #4ade80;
}

.stage-item.active {
  border-color: #2563eb;
  color: #60a5fa;
  background: #1d4ed833;
}

.stage-dot {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  background: #334155;
  margin: 0 auto 8px;
}

.stage-item.passed .stage-dot {
  background: #22c55e;
}

.stage-item.active .stage-dot {
  background: #3b82f6;
}

/* 小屏幕适配 */
@media (max-width: 960px) {
  .filter-card {
    grid-template-columns: 1fr;
  }

  .table-card {
    overflow-x: auto;
  }

  .table-card table {
    min-width: 980px;
  }

  .form-grid,
  .detail-card {
    grid-template-columns: 1fr;
  }

  .stage-list {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>