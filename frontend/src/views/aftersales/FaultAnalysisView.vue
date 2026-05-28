<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>故障分析方案管理</h1>
      </div>

      <div class="header-actions">
        <button class="primary-btn" @click="openCreateDialog">
          提交故障分析方案
        </button>
      </div>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索项目 / 方案名称 / 提交人 / 文件名"
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

      <select v-model="filters.auditStatus">
        <option value="">全部审核状态</option>
        <option value="draft">草稿</option>
        <option value="submitted">待审核</option>
        <option value="approved">审核通过</option>
        <option value="rejected">审核驳回</option>
      </select>

      <button class="query-btn">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <!-- 故障分析方案表格 -->
    <div class="table-card">
      <div class="table-card-header">
        <div>
          <h3>故障分析方案列表</h3>
          <span>共 {{ filteredAnalysisList.length }} 条故障分析方案</span>
        </div>
      </div>

      <div class="table-wrapper">
        <table class="version-table">
          <thead>
            <tr>
              <th>绑定项目</th>
              <th>方案名称</th>
              <th>方案文件</th>
              <th>提交人</th>
              <th>提交时间</th>
              <th>审核状态</th>
              <th>审核人</th>
              <th>审核时间</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="item in filteredAnalysisList" :key="item.id">
              <td>
                <span class="project-tag">{{ item.projectName }}</span>
              </td>

              <td>
                <button class="record-link" @click="viewAnalysis(item)">
                  {{ item.analysisName }}
                </button>
              </td>

              <td>
                <span class="file-tag" :title="item.fileName">
                  {{ item.fileName }}
                </span>
              </td>

              <td>{{ item.submitUser }}</td>

              <td class="muted">{{ item.submitTime }}</td>

              <td>
                <span class="status-tag" :class="item.auditStatus">
                  {{ getAuditStatusText(item.auditStatus) }}
                </span>
              </td>

              <td>{{ item.auditor || '-' }}</td>

              <td class="muted">{{ item.auditTime || '-' }}</td>

              <td class="operation-col">
                <div class="action-group">
                  <button class="text-btn" @click="viewAnalysis(item)">
                    查看
                  </button>

                  <button class="text-btn blue" @click="downloadFile(item)">
                    下载
                  </button>

                  <button
                    v-if="item.auditStatus === 'draft' || item.auditStatus === 'rejected'"
                    class="text-btn blue"
                    @click="submitAnalysis(item)"
                  >
                    提交
                  </button>

                  <button
                    v-if="item.auditStatus === 'submitted'"
                    class="text-btn green"
                    @click="auditAnalysis(item)"
                  >
                    审核
                  </button>

                  <button
                    v-if="item.auditStatus === 'draft' || item.auditStatus === 'rejected'"
                    class="text-btn red"
                    @click="deleteAnalysis(item)"
                  >
                    删除
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="table-footer">
      </div>
    </div>

    <!-- 新增故障分析方案弹窗 -->
    <div v-if="showCreateDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>提交故障分析方案</h3>
          <button @click="showCreateDialog = false">×</button>
        </div>

        <div class="form-grid">
          <label>
            绑定项目
            <select v-model="analysisForm.projectName">
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
            方案名称
            <input
              v-model="analysisForm.analysisName"
              placeholder="例如：香港屯马项目广播中断故障分析方案"
            />
          </label>

          <label>
            提交人
            <input
              v-model="analysisForm.submitUser"
              placeholder="请输入提交人"
            />
          </label>

          <label>
            文件名称
            <input
              v-model="analysisForm.fileName"
              placeholder="选择文件后自动填充"
              disabled
            />
          </label>

          <label class="full-row">
            故障分析方案文件
            <input
              type="file"
              accept=".doc,.docx,.pdf,.xls,.xlsx,.zip"
              @change="handleFileChange"
            />
          </label>

          <label class="full-row">
            故障分析说明
            <textarea
              v-model="analysisForm.remark"
              placeholder="例如：故障现象、原因分析、处理方案、验证结果、后续预防措施等"
            ></textarea>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showCreateDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="createAnalysis">
            保存
          </button>
        </div>
      </div>
    </div>

    <!-- 查看详情 / 审核弹窗 -->
    <div v-if="selectedAnalysis" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>故障分析方案详情</h3>
          <button @click="selectedAnalysis = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>绑定项目</span>
            <strong>{{ selectedAnalysis.projectName }}</strong>
          </div>

          <div>
            <span>方案名称</span>
            <strong>{{ selectedAnalysis.analysisName }}</strong>
          </div>

          <div>
            <span>方案文件</span>
            <strong>{{ selectedAnalysis.fileName }}</strong>
          </div>

          <div>
            <span>提交人</span>
            <strong>{{ selectedAnalysis.submitUser }}</strong>
          </div>

          <div>
            <span>提交时间</span>
            <strong>{{ selectedAnalysis.submitTime }}</strong>
          </div>

          <div>
            <span>审核状态</span>
            <strong>{{ getAuditStatusText(selectedAnalysis.auditStatus) }}</strong>
          </div>

          <div>
            <span>审核人</span>
            <strong>{{ selectedAnalysis.auditor || '-' }}</strong>
          </div>

          <div>
            <span>审核时间</span>
            <strong>{{ selectedAnalysis.auditTime || '-' }}</strong>
          </div>

          <div>
            <span>文件查看</span>
            <button class="inline-link" @click="openFile(selectedAnalysis)">
              点开查看文件
            </button>
          </div>
        </div>

        <div class="remark-card">
          <span>故障分析说明</span>
          <p>{{ selectedAnalysis.remark || '暂无说明' }}</p>
        </div>

        <div class="dialog-footer">
          <button
            v-if="selectedAnalysis.auditStatus === 'submitted'"
            class="green-btn"
            @click="approveAnalysis(selectedAnalysis)"
          >
            审核通过
          </button>

          <button
            v-if="selectedAnalysis.auditStatus === 'submitted'"
            class="red-btn"
            @click="rejectAnalysis(selectedAnalysis)"
          >
            审核驳回
          </button>

          <button class="reset-btn" @click="downloadFile(selectedAnalysis)">
            下载文件
          </button>

          <button class="primary-btn" @click="selectedAnalysis = null">
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
  projectName: '',
  auditStatus: ''
})

const showCreateDialog = ref(false)
const selectedAnalysis = ref(null)

const projectOptions = [
  '香港屯马项目',
  '波尔图二期项目',
  '阿根廷有轨项目',
  '波哥大有轨项目',
  '迪拜项目'
]

const analysisForm = reactive({
  projectName: '',
  analysisName: '',
  submitUser: '',
  fileName: '',
  file: null,
  fileUrl: '',
  remark: ''
})

const analysisList = ref([
  {
    id: 1,
    projectName: '香港屯马项目',
    analysisName: '香港屯马广播控制盒无音频输出故障分析方案',
    fileName: '香港屯马广播控制盒无音频输出故障分析方案.docx',
    fileUrl: '',
    submitUser: '售后人员',
    submitTime: '2026-05-10',
    auditStatus: 'approved',
    auditor: '领导',
    auditTime: '2026-05-11',
    remark: '分析广播控制盒上电后无音频输出问题，结论为音频服务启动时序异常，已提出软件优化方案。'
  },
  {
    id: 2,
    projectName: '波尔图二期项目',
    analysisName: '波尔图二期解码板SIP注册失败故障分析方案',
    fileName: '波尔图二期解码板SIP注册失败分析.pdf',
    fileUrl: '',
    submitUser: '研发人员',
    submitTime: '2026-05-16',
    auditStatus: 'submitted',
    auditor: '',
    auditTime: '',
    remark: '待领导审核确认。主要分析解码板在现场网络环境下SIP注册失败的问题。'
  },
  {
    id: 3,
    projectName: '阿根廷有轨项目',
    analysisName: '阿根廷司机提醒单元提醒音异常分析方案',
    fileName: '阿根廷司机提醒单元故障分析_草稿.docx',
    fileUrl: '',
    submitUser: '售后人员',
    submitTime: '2026-05-18',
    auditStatus: 'draft',
    auditor: '',
    auditTime: '',
    remark: '草稿，尚未提交审核。'
  },
  {
    id: 4,
    projectName: '迪拜项目',
    analysisName: '迪拜编码板音频中断故障分析方案',
    fileName: '迪拜编码板音频中断故障分析方案.xlsx',
    fileUrl: '',
    submitUser: '研发人员',
    submitTime: '2026-05-20',
    auditStatus: 'rejected',
    auditor: '领导',
    auditTime: '2026-05-21',
    remark: '分析结论不完整，需要补充RTP时间戳、buffer状态、网络抓包和现场日志。'
  }
])

const filteredAnalysisList = computed(() => {
  return analysisList.value.filter(item => {
    const keywordMatch =
      !filters.keyword ||
      item.projectName.includes(filters.keyword) ||
      item.analysisName.includes(filters.keyword) ||
      item.submitUser.includes(filters.keyword) ||
      item.fileName.includes(filters.keyword) ||
      item.remark.includes(filters.keyword)

    const projectMatch =
      !filters.projectName || item.projectName === filters.projectName

    const auditStatusMatch =
      !filters.auditStatus || item.auditStatus === filters.auditStatus

    return keywordMatch && projectMatch && auditStatusMatch
  })
})

function getAuditStatusText(status) {
  const map = {
    draft: '草稿',
    submitted: '待审核',
    approved: '审核通过',
    rejected: '审核驳回'
  }

  return map[status] || status
}

function resetFilters() {
  filters.keyword = ''
  filters.projectName = ''
  filters.auditStatus = ''
}

function openCreateDialog() {
  analysisForm.projectName = ''
  analysisForm.analysisName = ''
  analysisForm.submitUser = ''
  analysisForm.fileName = ''
  analysisForm.file = null
  analysisForm.fileUrl = ''
  analysisForm.remark = ''

  showCreateDialog.value = true
}

function handleFileChange(event) {
  const file = event.target.files[0]
  if (!file) return

  analysisForm.file = file
  analysisForm.fileName = file.name
  analysisForm.fileUrl = URL.createObjectURL(file)
}

function createAnalysis() {
  if (!analysisForm.projectName) {
    alert('请选择绑定项目')
    return
  }

  if (!analysisForm.analysisName) {
    alert('请输入方案名称')
    return
  }

  if (!analysisForm.submitUser) {
    alert('请输入提交人')
    return
  }

  if (!analysisForm.file) {
    alert('请上传故障分析方案文件')
    return
  }

  analysisList.value.unshift({
    id: Date.now(),
    projectName: analysisForm.projectName,
    analysisName: analysisForm.analysisName,
    fileName: analysisForm.fileName,
    fileUrl: analysisForm.fileUrl,
    submitUser: analysisForm.submitUser,
    submitTime: new Date().toISOString().slice(0, 10),
    auditStatus: 'draft',
    auditor: '',
    auditTime: '',
    remark: analysisForm.remark
  })

  showCreateDialog.value = false
}

function viewAnalysis(item) {
  selectedAnalysis.value = item
}

function auditAnalysis(item) {
  selectedAnalysis.value = item
}

function submitAnalysis(item) {
  item.auditStatus = 'submitted'
  item.auditor = ''
  item.auditTime = ''

  alert(`故障分析方案【${item.analysisName}】已提交领导审核`)
}

function approveAnalysis(item) {
  item.auditStatus = 'approved'
  item.auditor = '领导'
  item.auditTime = new Date().toISOString().slice(0, 10)
  selectedAnalysis.value = null

  alert(`故障分析方案【${item.analysisName}】审核通过`)
}

function rejectAnalysis(item) {
  item.auditStatus = 'rejected'
  item.auditor = '领导'
  item.auditTime = new Date().toISOString().slice(0, 10)
  selectedAnalysis.value = null

  alert(`故障分析方案【${item.analysisName}】已驳回`)
}

function openFile(item) {
  if (!item.fileUrl) {
    alert('当前是模拟数据，暂无可直接打开的原始文件')
    return
  }

  window.open(item.fileUrl, '_blank')
}

function downloadFile(item) {
  if (!item.fileUrl) {
    alert('当前是模拟数据，暂无可下载的原始文件')
    return
  }

  const link = document.createElement('a')
  link.href = item.fileUrl
  link.download = item.fileName || '故障分析方案文件'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

function deleteAnalysis(item) {
  const ok = confirm(`确认删除故障分析方案【${item.analysisName}】吗？`)
  if (!ok) return

  analysisList.value = analysisList.value.filter(record => record.id !== item.id)
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
  min-width: 1350px;
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

.project-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  background: #1d4ed833;
  color: #60a5fa;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
}

.file-tag {
  display: inline-block;
  width: 170px;
  max-width: 170px;
  padding: 4px 9px;
  border-radius: 999px;
  background: #33415566;
  color: #cbd5e1;
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;
}

.status-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
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

.inline-link {
  border: none;
  background: transparent;
  color: #60a5fa;
  cursor: pointer;
  padding: 0;
  font-size: 14px;
  text-align: left;
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

  .version-table {
    min-width: 1350px;
  }

  .form-grid,
  .detail-card {
    grid-template-columns: 1fr;
  }
}
</style>
