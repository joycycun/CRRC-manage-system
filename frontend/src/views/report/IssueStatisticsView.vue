<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>问题统计分析</h1>
      </div>

      <div class="header-actions">
        <button class="reset-btn" @click="exportIssueStatistics">
          导出问题统计
        </button>
      </div>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索项目 / 问题描述 / 负责人"
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

      <select v-model="filters.issueLevel">
        <option value="">全部问题类型</option>
        <option value="serious">严重</option>
        <option value="normal">一般</option>
        <option value="minor">轻微</option>
        <option value="suggestion">建议</option>
      </select>

      <select v-model="filters.issueStatus">
        <option value="">全部问题状态</option>
        <option value="open">未处理</option>
        <option value="processing">处理中</option>
        <option value="closed">已闭环</option>
      </select>

      <button class="query-btn" @click="loadIssueStatistics">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <!-- 问题数量统计 -->
    <div class="summary-panel">
      <div class="summary-card">
        <span>问题总数</span>
        <strong>{{ totalIssueCount }}</strong>
        <p>当前系统记录的问题总数量</p>
      </div>

      <div class="summary-card red">
        <span>严重问题</span>
        <strong>{{ seriousCount }}</strong>
        <p>高频或影响项目交付的问题</p>
      </div>

      <div class="summary-card yellow">
        <span>处理中</span>
        <strong>{{ processingCount }}</strong>
        <p>当前仍在处理中的问题</p>
      </div>

      <div class="summary-card green">
        <span>已闭环</span>
        <strong>{{ closedCount }}</strong>
        <p>已经完成处理和确认的问题</p>
      </div>
    </div>

    <!-- 问题类型和项目分布 -->
    <div class="chart-grid">
      <div class="stat-card">
        <div class="table-card-header">
          <div>
            <h3>问题类型分布</h3>
            <span>按严重、一般、轻微、建议分类</span>
          </div>
        </div>

        <div class="type-list">
          <div
            v-for="item in issueLevelSummary"
            :key="item.level"
            class="type-item"
          >
            <div class="type-info">
              <span class="level-tag" :class="item.level">
                {{ getIssueLevelText(item.level) }}
              </span>
              <strong>{{ item.count }} 个</strong>
            </div>

            <div class="progress-bg">
              <div
                class="progress-bar"
                :class="item.level"
                :style="{ width: getPercent(item.count, totalIssueCount) + '%' }"
              ></div>
            </div>

            <p>{{ getPercent(item.count, totalIssueCount) }}%</p>
          </div>
        </div>
      </div>

      <div class="stat-card">
        <div class="table-card-header">
          <div>
            <h3>不同项目问题分布</h3>
            <span>查看每个项目的问题数量</span>
          </div>
        </div>

        <div class="project-distribution">
          <div
            v-for="item in projectIssueSummary"
            :key="item.projectName"
            class="project-stat-item"
          >
            <div class="project-stat-title">
              <strong>{{ item.projectName }}</strong>
              <span>{{ item.count }} 个</span>
            </div>

            <div class="progress-bg">
              <div
                class="progress-bar blue"
                :style="{ width: getPercent(item.count, maxProjectIssueCount) + '%' }"
              ></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 问题列表 -->
    <div class="table-card">
      <div class="table-card-header">
        <div>
          <h3>问题统计列表</h3>
          <span>共 {{ filteredIssueList.length }} 条问题记录</span>
        </div>
      </div>

      <div class="table-wrapper">
        <table class="version-table">
          <thead>
            <tr>
              <th>项目名称</th>
              <th>问题描述</th>
              <th>问题类型</th>
              <th>问题状态</th>
              <th>问题负责人</th>
              <th>出现次数</th>
              <th>重新打开</th>
              <th>最近更新时间</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="item in filteredIssueList" :key="item.id">
              <td>
                <span class="project-tag">
                  {{ item.projectName }}
                </span>
              </td>

              <td class="issue-title-cell">
                <button
                  class="record-link"
                  :title="item.issueTitle"
                  @click="viewIssue(item)"
                >
                  {{ item.issueTitle }}
                </button>
              </td>

              <td>
                <span class="level-tag" :class="item.issueLevel">
                  {{ getIssueLevelText(item.issueLevel) }}
                </span>
              </td>

              <td>
                <span class="status-tag" :class="item.issueStatus">
                  {{ getIssueStatusText(item.issueStatus) }}
                </span>
              </td>

              <td>{{ item.owner }}</td>

              <td>
                <span class="count-tag">
                  {{ item.frequency }} 次
                </span>
              </td>

              <td>
                <span class="count-tag">
                  {{ item.reopenCount }} 次
                </span>
              </td>

              <td class="muted">{{ item.updateTime }}</td>

              <td class="operation-col">
                <div class="action-group">
                  <button class="text-btn" @click="viewIssue(item)">
                    查看
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="table-footer">
        问题类型按照出现频率和影响程度划分为：严重、一般、轻微、建议。
      </div>
    </div>

    <!-- 查看问题详情弹窗 -->
    <div v-if="selectedIssue" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>问题详情</h3>
          <button @click="selectedIssue = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>项目名称</span>
            <strong>{{ selectedIssue.projectName }}</strong>
          </div>

          <div>
            <span>问题类型</span>
            <strong>{{ getIssueLevelText(selectedIssue.issueLevel) }}</strong>
          </div>

          <div>
            <span>问题状态</span>
            <strong>{{ getIssueStatusText(selectedIssue.issueStatus) }}</strong>
          </div>

          <div>
            <span>问题负责人</span>
            <strong>{{ selectedIssue.owner }}</strong>
          </div>

          <div>
            <span>出现次数</span>
            <strong>{{ selectedIssue.frequency }} 次</strong>
          </div>

          <div>
            <span>重新打开次数</span>
            <strong>{{ selectedIssue.reopenCount }} 次</strong>
          </div>

          <div>
            <span>最近更新时间</span>
            <strong>{{ selectedIssue.updateTime }}</strong>
          </div>
        </div>

        <div class="remark-card">
          <span>问题描述</span>
          <p>{{ selectedIssue.issueTitle }}</p>
        </div>

        <div class="remark-card">
          <span>处理说明</span>
          <p>{{ selectedIssue.remark || '暂无处理说明' }}</p>
        </div>

        <div class="dialog-footer">
          <button class="primary-btn" @click="selectedIssue = null">
            关闭
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { getIssueStatisticsReport } from '@/api/report'

const filters = reactive({
  keyword: '',
  projectName: '',
  issueLevel: '',
  issueStatus: ''
})

const selectedIssue = ref(null)

const issueList = ref([])

onMounted(() => {
  loadIssueStatistics()
})

function getResponseData(res) {
  return res && res.data ? res.data : res
}

function normalizeIssue(item) {
  return {
    id: item.id,
    projectName: item.projectName || '',
    issueTitle: item.issueTitle || '',
    issueLevel: item.issueLevel || 'normal',
    issueStatus: item.issueStatus || 'open',
    owner: item.owner || item.ownerName || '',
    frequency: Number(item.frequency || 1),
    reopenCount: Number(item.reopenCount || item.reopen_count || 0),
    updateTime: item.updateTime || '',
    remark: item.remark || '',
    issueSource: item.issueSource || ''
  }
}

async function loadIssueStatistics() {
  try {
    const res = await getIssueStatisticsReport()
    const result = getResponseData(res)
    if (result.code !== 200) {
      alert(result.msg || '加载问题统计失败')
      return
    }
    issueList.value = (result.data || []).map(normalizeIssue)
  } catch (err) {
    console.error('加载问题统计失败：', err)
    alert(err.response?.data || '加载问题统计失败')
  }
}

const projectOptions = computed(() => [...new Set(issueList.value.map(item => item.projectName).filter(Boolean))])

const filteredIssueList = computed(() => {
  return issueList.value.filter(item => {
    const keywordMatch =
      !filters.keyword ||
      item.projectName.includes(filters.keyword) ||
      item.issueTitle.includes(filters.keyword) ||
      item.owner.includes(filters.keyword) ||
      item.remark.includes(filters.keyword)

    const projectMatch =
      !filters.projectName || item.projectName === filters.projectName

    const levelMatch =
      !filters.issueLevel || item.issueLevel === filters.issueLevel

    const statusMatch =
      !filters.issueStatus || item.issueStatus === filters.issueStatus

    return keywordMatch && projectMatch && levelMatch && statusMatch
  })
})

const totalIssueCount = computed(() => filteredIssueList.value.length)

const seriousCount = computed(() => {
  return filteredIssueList.value.filter(item => item.issueLevel === 'serious').length
})

const processingCount = computed(() => {
  return filteredIssueList.value.filter(item => item.issueStatus === 'processing').length
})

const closedCount = computed(() => {
  return filteredIssueList.value.filter(item => item.issueStatus === 'closed').length
})

const issueLevelSummary = computed(() => {
  const levels = ['serious', 'normal', 'minor', 'suggestion']

  return levels.map(level => {
    return {
      level,
      count: filteredIssueList.value.filter(item => item.issueLevel === level).length
    }
  })
})

const projectIssueSummary = computed(() => {
  return projectOptions.value
    .map(projectName => {
      return {
        projectName,
        count: filteredIssueList.value.filter(item => item.projectName === projectName).length
      }
    })
    .filter(item => item.count > 0)
})

const maxProjectIssueCount = computed(() => {
  if (projectIssueSummary.value.length === 0) return 1

  return Math.max(...projectIssueSummary.value.map(item => item.count))
})

function getIssueLevelText(level) {
  const map = {
    serious: '严重',
    normal: '一般',
    minor: '轻微',
    suggestion: '建议'
  }

  return map[level] || level
}

function getIssueStatusText(status) {
  const map = {
    open: '未处理',
    processing: '处理中',
    closed: '已闭环'
  }

  return map[status] || status
}

function getPercent(count, total) {
  if (!total) return 0
  return Math.round((count / total) * 100)
}

function resetFilters() {
  filters.keyword = ''
  filters.projectName = ''
  filters.issueLevel = ''
  filters.issueStatus = ''
}

function viewIssue(item) {
  selectedIssue.value = item
}

function exportIssueStatistics() {
  const header = [
    '项目名称',
    '问题描述',
    '问题类型',
    '问题状态',
    '问题负责人',
    '出现次数',
    '重新打开次数',
    '最近更新时间',
    '处理说明'
  ]

  const rows = filteredIssueList.value.map(item => [
    item.projectName,
    item.issueTitle,
    getIssueLevelText(item.issueLevel),
    getIssueStatusText(item.issueStatus),
    item.owner,
    item.frequency,
    item.reopenCount,
    item.updateTime,
    item.remark || ''
  ])

  const csvContent = [header, ...rows]
    .map(row => row.map(cell => `"${String(cell).replace(/"/g, '""')}"`).join(','))
    .join('\n')

  const blob = new Blob(['\uFEFF' + csvContent], {
    type: 'text/csv;charset=utf-8;'
  })

  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = '问题统计分析.csv'
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

.header-actions {
  display: flex;
  gap: 12px;
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
  grid-template-columns: 1.4fr 180px 160px 160px 90px 90px;
  gap: 12px;
  margin-bottom: 20px;
}

.filter-card input,
.filter-card select {
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

.summary-panel {
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 14px;
  padding: 16px;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 14px;
  margin-bottom: 20px;
}

.summary-card {
  background: #020617;
  border: 1px solid #1e293b;
  border-radius: 12px;
  padding: 16px;
}

.summary-card span {
  color: #94a3b8;
  font-size: 13px;
}

.summary-card strong {
  display: block;
  margin-top: 8px;
  color: #f8fafc;
  font-size: 28px;
  font-weight: 800;
}

.summary-card p {
  margin: 6px 0 0;
  color: #64748b;
  font-size: 12px;
}

.summary-card.red strong {
  color: #f87171;
}

.summary-card.yellow strong {
  color: #fbbf24;
}

.summary-card.green strong {
  color: #4ade80;
}

.chart-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 18px;
  margin-bottom: 20px;
}

.stat-card,
.table-card {
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 14px;
  overflow: hidden;
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

.type-list,
.project-distribution {
  padding: 16px;
  display: grid;
  gap: 14px;
}

.type-item,
.project-stat-item {
  background: #020617;
  border: 1px solid #1e293b;
  border-radius: 10px;
  padding: 14px;
}

.type-info,
.project-stat-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.type-info strong,
.project-stat-title strong {
  color: #f8fafc;
  font-size: 13px;
}

.project-stat-title span,
.type-item p {
  color: #94a3b8;
  font-size: 12px;
}

.type-item p {
  margin: 8px 0 0;
  text-align: right;
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
  border-radius: 999px;
}

.progress-bar.serious {
  background: #ef4444;
}

.progress-bar.normal {
  background: #f59e0b;
}

.progress-bar.minor {
  background: #3b82f6;
}

.progress-bar.suggestion {
  background: #22c55e;
}

.progress-bar.blue {
  background: #3b82f6;
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
  min-width: 1380px;
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

.version-table th:nth-child(1),
.version-table td:nth-child(1) {
  width: 180px;
}

.version-table th:nth-child(2),
.version-table td:nth-child(2) {
  width: 300px;
}

.version-table th:nth-child(3),
.version-table td:nth-child(3),
.version-table th:nth-child(4),
.version-table td:nth-child(4) {
  width: 110px;
}

.version-table th:nth-child(5),
.version-table td:nth-child(5) {
  width: 130px;
}

.version-table th:nth-child(6),
.version-table td:nth-child(6),
.version-table th:nth-child(7),
.version-table td:nth-child(7) {
  width: 110px;
}

.version-table th:nth-child(8),
.version-table td:nth-child(8) {
  width: 170px;
}

.issue-title-cell {
  overflow: hidden;
}

.project-tag {
  display: inline-flex;
  max-width: 150px;
  padding: 4px 9px;
  border-radius: 999px;
  background: #1d4ed833;
  color: #60a5fa;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.record-link {
  display: block;
  width: 100%;
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
}

.record-link:hover {
  color: #93c5fd;
  text-decoration: underline;
}

.level-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
}

.level-tag.serious {
  background: #dc262633;
  color: #f87171;
}

.level-tag.normal {
  background: #d9770633;
  color: #fbbf24;
}

.level-tag.minor {
  background: #1d4ed833;
  color: #60a5fa;
}

.level-tag.suggestion {
  background: #16a34a33;
  color: #4ade80;
}

.status-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
}

.status-tag.open {
  background: #dc262633;
  color: #f87171;
}

.status-tag.processing {
  background: #d9770633;
  color: #fbbf24;
}

.status-tag.closed {
  background: #16a34a33;
  color: #4ade80;
}

.count-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  background: #33415566;
  color: #cbd5e1;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
}

.muted {
  color: #94a3b8 !important;
}

.operation-col {
  width: 120px;
  text-align: right !important;
}

.action-group {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 12px;
}

.text-btn {
  border: none;
  background: transparent;
  color: #60a5fa;
  font-size: 13px;
  cursor: pointer;
  white-space: nowrap;
}

.text-btn:hover {
  color: #93c5fd;
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

.remark-card {
  margin: 0 20px 20px;
}

.remark-card p {
  margin: 0;
  color: #cbd5e1;
  font-size: 13px;
  line-height: 1.6;
}

.dialog-footer {
  padding: 16px 20px;
  border-top: 1px solid #1e293b;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

@media (max-width: 1200px) {
  .summary-panel,
  .chart-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .filter-card {
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 960px) {
  .summary-panel,
  .chart-grid,
  .filter-card {
    grid-template-columns: 1fr;
  }

  .version-table {
    min-width: 1250px;
  }

  .detail-card {
    grid-template-columns: 1fr;
  }
}
</style>
