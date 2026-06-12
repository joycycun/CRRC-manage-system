<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>项目进度报告</h1>
      </div>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索项目名称 / 负责人 / 当前阶段"
      />

      <select v-model="filters.currentStage">
        <option value="">全部阶段</option>
        <option
          v-for="stage in stageOptions"
          :key="stage.value"
          :value="stage.value"
        >
          {{ stage.label }}
        </option>
      </select>

      <select v-model="filters.projectStatus">
        <option value="">全部项目状态</option>
        <option value="running">进行中</option>
        <option value="closed">已关闭</option>
      </select>

      <button class="query-btn" @click="loadProjectProgress">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <!-- 项目进度表格 -->
    <div class="table-card">
      <div class="table-card-header">
        <div>
          <h3>项目进度列表</h3>
          <span>共 {{ filteredProjectList.length }} 个项目</span>
        </div>
      </div>

      <div class="table-wrapper">
        <table class="version-table">
          <thead>
            <tr>
              <th>项目名称</th>
              <th>负责人</th>
              <th>当前阶段</th>
              <th>项目进度</th>
              <th>项目状态</th>
              <th>最后更新时间</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="item in filteredProjectList" :key="item.id">
              <td>
                <button class="record-link" @click="viewProject(item)">
                  {{ item.name }}
                </button>
              </td>

              <td>{{ item.owner }}</td>

              <td>
                <span class="stage-tag">
                  {{ getStageText(item.currentStage) }}
                </span>
              </td>

              <td>
                <div class="progress-cell">
                  <div class="progress-title">
                    <span>{{ getProjectProgress(item) }}%</span>
                  </div>

                  <div class="progress-bg">
                    <div
                      class="progress-bar"
                      :class="{ closed: item.projectStatus === 'closed' }"
                      :style="{ width: getProjectProgress(item) + '%' }"
                    ></div>
                  </div>
                </div>
              </td>

              <td>
                <span class="status-tag" :class="item.projectStatus">
                  {{ getProjectStatusText(item.projectStatus) }}
                </span>
              </td>

              <td class="muted">{{ item.updateTime }}</td>

              <td class="operation-col">
                <div class="action-group">
                  <button class="text-btn" @click="viewProject(item)">
                    查看
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="table-footer">
        每个阶段占比 10%，选择到当前阶段后，系统自动计算项目进度。
      </div>
    </div>

    <!-- 阶段进度总览 -->
    <div class="stage-overview-card">
      <div class="table-card-header">
        <div>
          <h3>项目阶段说明</h3>
          <span>共 {{ stageOptions.length }} 个阶段，每个阶段 10%</span>
        </div>
      </div>

      <div class="stage-grid">
        <div
          v-for="stage in stageOptions"
          :key="stage.value"
          class="stage-item"
        >
          <span>{{ stage.percent }}%</span>
          <strong>{{ stage.label }}</strong>
        </div>
      </div>
    </div>

    <!-- 设置阶段弹窗 -->
    <div v-if="showStageDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>设置项目当前阶段</h3>
          <button @click="showStageDialog = false">×</button>
        </div>

        <div v-if="editingProject" class="stage-setting">
          <div class="project-info-card">
            <span>项目名称</span>
            <strong>{{ editingProject.name }}</strong>
          </div>

          <label>
            当前阶段
            <select v-model="stageForm.currentStage">
              <option
                v-for="stage in stageOptions"
                :key="stage.value"
                :value="stage.value"
              >
                {{ stage.label }} / {{ stage.percent }}%
              </option>
            </select>
          </label>

          <div class="stage-preview">
            <span>设置后项目进度</span>
            <strong>{{ getStagePercent(stageForm.currentStage) }}%</strong>

            <div class="progress-bg">
              <div
                class="progress-bar"
                :style="{ width: getStagePercent(stageForm.currentStage) + '%' }"
              ></div>
            </div>
          </div>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showStageDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="saveProjectStage">
            保存
          </button>
        </div>
      </div>
    </div>

    <!-- 查看详情弹窗 -->
    <div v-if="selectedProject" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>项目进度详情</h3>
          <button @click="selectedProject = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>项目名称</span>
            <strong>{{ selectedProject.name }}</strong>
          </div>

          <div>
            <span>负责人</span>
            <strong>{{ selectedProject.owner }}</strong>
          </div>

          <div>
            <span>当前阶段</span>
            <strong>{{ getStageText(selectedProject.currentStage) }}</strong>
          </div>

          <div>
            <span>项目进度</span>
            <strong>{{ getProjectProgress(selectedProject) }}%</strong>
          </div>

          <div>
            <span>项目状态</span>
            <strong>{{ getProjectStatusText(selectedProject.projectStatus) }}</strong>
          </div>

          <div>
            <span>最后更新时间</span>
            <strong>{{ selectedProject.updateTime }}</strong>
          </div>
        </div>

        <div class="progress-detail">
          <h4>阶段完成情况</h4>

          <div class="timeline">
            <div
              v-for="stage in stageOptions"
              :key="stage.value"
              class="timeline-item"
              :class="{
                finished: stage.percent <= getProjectProgress(selectedProject),
                current: stage.value === selectedProject.currentStage
              }"
            >
              <div class="timeline-dot"></div>

              <div class="timeline-content">
                <strong>{{ stage.label }}</strong>
                <span>{{ stage.percent }}%</span>
              </div>
            </div>
          </div>
        </div>

        <div class="dialog-footer">
          <button class="primary-btn" @click="selectedProject = null">
            关闭
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { getProjectProgressReport } from '@/api/report'

const filters = reactive({
  keyword: '',
  currentStage: '',
  projectStatus: ''
})

const showStageDialog = ref(false)
const selectedProject = ref(null)
const editingProject = ref(null)

const stageForm = reactive({
  currentStage: ''
})

const stageOptions = [
  {
    value: '立项',
    label: '立项',
    percent: 10
  },
  {
    value: '需求书首次确认',
    label: '需求书首次确认',
    percent: 20
  },
  {
    value: '硬件检测',
    label: '硬件检测',
    percent: 30
  },
  {
    value: '软件研发',
    label: '软件研发',
    percent: 40
  },
  {
    value: '内部初始测试',
    label: '内部初始测试',
    percent: 50
  },
  {
    value: '内部初始测试是否问题闭环',
    label: '内部初始测试是否问题闭环',
    percent: 60
  },
  {
    value: '样机联调',
    label: '样机联调',
    percent: 70
  },
  {
    value: '出厂测试',
    label: '出厂测试',
    percent: 80
  },
  {
    value: '收货',
    label: '收货',
    percent: 90
  },
  {
    value: '已关闭',
    label: '项目关闭',
    percent: 100
  }
]

const projectList = ref([])

onMounted(() => {
  loadProjectProgress()
})

function getResponseData(res) {
  return res && res.data ? res.data : res
}

function normalizeProject(item) {
  return {
    id: item.id || item.projectId,
    name: item.name || item.projectName || '',
    owner: item.owner || '',
    currentStage: item.currentStage || item.stage || '',
    projectStatus: item.projectStatus || (item.status === '已关闭' ? 'closed' : 'running'),
    updateTime: item.updateTime || '',
    progress: Math.round(Number(item.progress || 0))
  }
}

async function loadProjectProgress() {
  try {
    const res = await getProjectProgressReport()
    const result = getResponseData(res)
    if (result.code !== 200) {
      alert(result.msg || '加载项目进度失败')
      return
    }
    projectList.value = (result.data || []).map(normalizeProject)
  } catch (err) {
    console.error('加载项目进度失败：', err)
    alert(err.response?.data || '加载项目进度失败')
  }
}

const filteredProjectList = computed(() => {
  return projectList.value.filter(item => {
    const keywordMatch =
      !filters.keyword ||
      item.name.includes(filters.keyword) ||
      item.owner.includes(filters.keyword) ||
      getStageText(item.currentStage).includes(filters.keyword)

    const stageMatch =
      !filters.currentStage || item.currentStage === filters.currentStage

    const statusMatch =
      !filters.projectStatus || item.projectStatus === filters.projectStatus

    return keywordMatch && stageMatch && statusMatch
  })
})

function getStageText(stageValue) {
  const stage = stageOptions.find(item => item.value === stageValue)
  return stage ? stage.label : stageValue
}

function getStagePercent(stageValue) {
  const stage = stageOptions.find(item => item.value === stageValue)
  return stage ? stage.percent : 0
}

function getProjectProgress(project) {
  return Math.round(Number(project.progress || getStagePercent(project.currentStage)))
}

function getProjectStatusText(status) {
  const map = {
    running: '进行中',
    closed: '已关闭'
  }

  return map[status] || status
}

function resetFilters() {
  filters.keyword = ''
  filters.currentStage = ''
  filters.projectStatus = ''
}

function viewProject(item) {
  selectedProject.value = item
}

function openStageDialog(item) {
  editingProject.value = item
  stageForm.currentStage = item.currentStage
  showStageDialog.value = true
}

function saveProjectStage() {
  if (!editingProject.value) return

  editingProject.value.currentStage = stageForm.currentStage
  editingProject.value.updateTime = new Date().toISOString().slice(0, 10)

  if (stageForm.currentStage === 'project_closed') {
    editingProject.value.projectStatus = 'closed'
  } else {
    editingProject.value.projectStatus = 'running'
  }

  showStageDialog.value = false

  if (selectedProject.value) {
    selectedProject.value = editingProject.value
  }
}

function closeProject(item) {
  const ok = confirm(`确认将项目【${item.name}】设置为关闭吗？`)
  if (!ok) return

  item.currentStage = 'project_closed'
  item.projectStatus = 'closed'
  item.updateTime = new Date().toISOString().slice(0, 10)

  if (selectedProject.value) {
    selectedProject.value = item
  }
}

function reopenProject(item) {
  item.projectStatus = 'running'

  if (item.currentStage === 'project_closed') {
    item.currentStage = 'receipt'
  }

  item.updateTime = new Date().toISOString().slice(0, 10)
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

.header-actions {
  display: flex;
  gap: 12px;
}

.primary-btn,
.query-btn,
.reset-btn,
.green-btn,
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

.green-btn {
  border: 1px solid #166534;
  background: #052e16;
  color: #86efac;
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
  grid-template-columns: 1.4fr 220px 180px 90px 90px;
  gap: 12px;
  margin-bottom: 20px;
}

.filter-card input,
.filter-card select,
.stage-setting select {
  border: 1px solid #334155;
  border-radius: 8px;
  background: #020617;
  color: #e2e8f0;
  padding: 0 12px;
  outline: none;
  height: 36px;
}

.filter-card input::placeholder {
  color: #64748b;
}

.table-card,
.stage-overview-card {
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 14px;
  overflow: hidden;
  margin-bottom: 20px;
}

.table-card-header {
  padding: 14px 16px;
  border-bottom: 1px solid #1e293b;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.table-card-header h3 {
  margin: 0;
  color: #f8fafc;
  font-size: 15px;
}

.table-card-header span {
  color: #64748b;
  font-size: 12px;
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

.version-table {
  width: 100%;
  min-width: 1150px;
  border-collapse: collapse;
  table-layout: fixed;
}

.version-table thead {
  background: #020617;
}

.version-table th {
  padding: 14px 16px;
  color: #94a3b8;
  font-size: 12px;
  font-weight: 600;
  text-align: left;
  border-bottom: 1px solid #1e293b;
  white-space: nowrap;
}

.version-table td {
  padding: 15px 16px;
  font-size: 13px;
  color: #e2e8f0;
  border-bottom: 1px solid #1e293b;
  vertical-align: middle;
}

.record-link {
  border: none;
  background: transparent;
  color: #60a5fa;
  font-size: 13px;
  font-weight: 800;
  cursor: pointer;
  padding: 0;
  text-align: left;
}

.record-link:hover {
  color: #93c5fd;
  text-decoration: underline;
}

.stage-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  background: #1d4ed833;
  color: #60a5fa;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
}

.status-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
}

.status-tag.running {
  background: #16a34a33;
  color: #4ade80;
}

.status-tag.closed {
  background: #47556933;
  color: #94a3b8;
}

.progress-cell {
  width: 100%;
}

.progress-title {
  display: flex;
  justify-content: flex-end;
  color: #60a5fa;
  font-size: 12px;
  font-weight: 700;
  margin-bottom: 6px;
}

.progress-bg {
  width: 100%;
  height: 8px;
  background: #1e293b;
  border-radius: 999px;
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  background: #3b82f6;
  border-radius: 999px;
}

.progress-bar.closed {
  background: #64748b;
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

.text-btn.red {
  color: #f87171;
}

.table-footer {
  padding: 12px 16px;
  color: #64748b;
  font-size: 12px;
}

.stage-grid {
  padding: 16px;
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 12px;
}

.stage-item {
  background: #020617;
  border: 1px solid #1e293b;
  border-radius: 10px;
  padding: 14px;
}

.stage-item span {
  display: block;
  color: #60a5fa;
  font-size: 13px;
  font-weight: 800;
  margin-bottom: 8px;
}

.stage-item strong {
  color: #f8fafc;
  font-size: 13px;
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
  width: 920px;
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

.stage-setting {
  padding: 20px;
  display: grid;
  gap: 16px;
}

.stage-setting label {
  display: flex;
  flex-direction: column;
  gap: 8px;
  color: #cbd5e1;
  font-size: 13px;
}

.project-info-card,
.stage-preview {
  background: #020617;
  border: 1px solid #1e293b;
  border-radius: 10px;
  padding: 14px;
}

.project-info-card span,
.stage-preview span {
  display: block;
  color: #64748b;
  font-size: 12px;
  margin-bottom: 8px;
}

.project-info-card strong,
.stage-preview strong {
  color: #f8fafc;
  font-size: 16px;
}

.stage-preview .progress-bg {
  margin-top: 12px;
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
  word-break: break-all;
}

.progress-detail {
  margin: 0 20px 20px;
  padding: 16px;
  background: #020617;
  border: 1px solid #1e293b;
  border-radius: 12px;
}

.progress-detail h4 {
  margin: 0 0 16px;
  font-size: 15px;
  color: #f8fafc;
}

.timeline {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 12px;
}

.timeline-item {
  position: relative;
  border: 1px solid #1e293b;
  border-radius: 10px;
  padding: 12px;
  background: #0f172a;
}

.timeline-item.finished {
  border-color: #2563eb;
  background: #1d4ed81f;
}

.timeline-item.current {
  border-color: #4ade80;
  background: #16a34a1f;
}

.timeline-dot {
  width: 9px;
  height: 9px;
  border-radius: 999px;
  background: #475569;
  margin-bottom: 8px;
}

.timeline-item.finished .timeline-dot {
  background: #60a5fa;
}

.timeline-item.current .timeline-dot {
  background: #4ade80;
}

.timeline-content strong {
  display: block;
  color: #f8fafc;
  font-size: 12px;
  margin-bottom: 6px;
}

.timeline-content span {
  color: #94a3b8;
  font-size: 12px;
}

@media (max-width: 1200px) {
  .stage-grid,
  .timeline {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 960px) {
  .filter-card {
    grid-template-columns: 1fr;
  }

  .version-table {
    min-width: 1150px;
  }

  .stage-grid,
  .timeline {
    grid-template-columns: 1fr;
  }

  .detail-card {
    grid-template-columns: 1fr;
  }
}
</style>
