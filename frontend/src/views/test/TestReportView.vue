<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>测试报告管理</h1>
      </div>
      <button class="primary-btn" @click="openUploadDialog">
        上传测试报告
      </button>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索项目名称 / 测试报告名称 / 上传人"
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

    <!-- 数据表格 -->
    <div class="table-card">
      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>测试报告</th>
              <th>绑定项目</th>
              <th>上传人</th>
              <th>上传时间</th>
              <th>审核状态</th>
              <th>审核人</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="item in filteredTestReportList" :key="item.id">
              <td>
                <button class="report-link" @click="viewTestReport(item)">
                  {{ item.reportName }}
                </button>
              </td>

              <td>
                <span class="project-tag">
                  {{ item.projectName }}
                </span>
              </td>

              <td>{{ item.uploader }}</td>

              <td class="muted">{{ item.uploadTime }}</td>

              <td>
                <span class="status-tag" :class="item.auditStatus">
                  {{ getAuditStatusText(item.auditStatus) }}
                </span>
              </td>

              <td>{{ item.auditor || '-' }}</td>

              <td class="operation-col">
                <div class="action-group">
                  <button class="text-btn" @click="viewTestReport(item)">
                    查看
                  </button>

                  <button class="text-btn blue" @click="downloadTestReport(item)">
                    下载
                  </button>

                  <button
                    v-if="item.auditStatus === 'draft' || item.auditStatus === 'rejected'"
                    class="text-btn blue"
                    @click="submitTestReport(item)"
                  >
                    提交
                  </button>

                  <button
                    v-if="item.auditStatus === 'submitted'"
                    class="text-btn green"
                    @click="auditTestReport(item)"
                  >
                    审核
                  </button>

                  <button
                    v-if="item.auditStatus === 'draft' || item.auditStatus === 'rejected'"
                    class="text-btn red"
                    @click="deleteTestReport(item)"
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
        共 {{ filteredTestReportList.length }} 条测试报告记录
      </div>
    </div>

    <!-- 上传测试报告弹窗 -->
    <div v-if="showUploadDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>上传测试报告</h3>
          <button @click="showUploadDialog = false">×</button>
        </div>

        <div class="form-grid">
          <label>
            绑定项目
            <select v-model="uploadForm.projectName">
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
            测试报告名称
            <input
              v-model="uploadForm.reportName"
              placeholder="例如：香港屯马项目内部初始测试报告"
            />
          </label>

          <label>
            当前上传人
            <input
              :value="currentUserName"
              disabled
            />
          </label>

          <label class="full-row">
            测试报告文件
            <input
              type="file"
              accept=".doc,.docx,.xls,.xlsx,.pdf,.zip"
              @change="handleFileChange"
            />
          </label>

          <label class="full-row">
            测试结论 / 说明
            <textarea
              v-model="uploadForm.remark"
              placeholder="例如：本轮测试覆盖广播、报警、SIP注册、音频播放、网络通信等功能，遗留问题已记录到问题闭环"
            ></textarea>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showUploadDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="uploadTestReport">
            保存上传
          </button>
        </div>
      </div>
    </div>

    <!-- 查看详情 / 审核弹窗 -->
    <div v-if="selectedTestReport" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>测试报告详情</h3>
          <button @click="selectedTestReport = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>测试报告名称</span>
            <strong>{{ selectedTestReport.reportName }}</strong>
          </div>

          <div>
            <span>绑定项目</span>
            <strong>{{ selectedTestReport.projectName }}</strong>
          </div>

          <div>
            <span>文件名称</span>
            <strong>{{ selectedTestReport.fileName || '-' }}</strong>
          </div>

          <div>
            <span>上传人</span>
            <strong>{{ selectedTestReport.uploader }}</strong>
          </div>

          <div>
            <span>上传时间</span>
            <strong>{{ selectedTestReport.uploadTime }}</strong>
          </div>

          <div>
            <span>审核状态</span>
            <strong>{{ getAuditStatusText(selectedTestReport.auditStatus) }}</strong>
          </div>

          <div>
            <span>审核人</span>
            <strong>{{ selectedTestReport.auditor || '-' }}</strong>
          </div>

          <div>
            <span>审核时间</span>
            <strong>{{ selectedTestReport.auditTime || '-' }}</strong>
          </div>

          <div>
            <span>文件查看</span>
            <button class="inline-link" @click="openTestReportFile(selectedTestReport)">
              点开查看文件
            </button>
          </div>
        </div>

        <div class="remark-card">
          <span>测试结论 / 说明</span>
          <p>{{ selectedTestReport.remark || '暂无说明' }}</p>
        </div>

        <div class="dialog-footer">
          <button
            v-if="selectedTestReport.auditStatus === 'submitted'"
            class="green-btn"
            @click="approveTestReport(selectedTestReport)"
          >
            审核通过
          </button>

          <button
            v-if="selectedTestReport.auditStatus === 'submitted'"
            class="red-btn"
            @click="rejectTestReport(selectedTestReport)"
          >
            审核驳回
          </button>

          <button class="reset-btn" @click="downloadTestReport(selectedTestReport)">
            下载文件
          </button>

          <button class="primary-btn" @click="selectedTestReport = null">
            关闭
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'

const currentUserName = ref(
  localStorage.getItem('username') ||
  localStorage.getItem('accountName') ||
  localStorage.getItem('realName') ||
  '当前用户'
)

const filters = reactive({
  keyword: '',
  projectName: '',
  auditStatus: ''
})

const showUploadDialog = ref(false)
const selectedTestReport = ref(null)

const projectOptions = [
  '香港屯马项目',
  '波尔图二期项目',
  '阿根廷有轨项目',
  '波哥大有轨项目',
  '迪拜项目'
]

const uploadForm = reactive({
  projectName: '',
  reportName: '',
  fileName: '',
  file: null,
  fileUrl: '',
  remark: ''
})

const testReportList = ref([
  {
    id: 1,
    projectName: '香港屯马项目',
    reportName: '香港屯马项目内部初始测试报告',
    fileName: '香港屯马内部初始测试报告_V1.0.docx',
    fileUrl: '',
    uploader: '寸诗睿',
    uploadTime: '2026-05-10',
    auditStatus: 'approved',
    auditor: '领导',
    auditTime: '2026-05-11',
    remark: '内部初始测试通过，广播、报警、SIP注册、网络通信功能测试正常。'
  },
  {
    id: 2,
    projectName: '波尔图二期项目',
    reportName: '波尔图二期广播系统测试报告',
    fileName: '波尔图二期测试报告_V0.9.docx',
    fileUrl: '',
    uploader: '寸诗睿',
    uploadTime: '2026-05-16',
    auditStatus: 'submitted',
    auditor: '',
    auditTime: '',
    remark: '待领导审核确认。'
  },
  {
    id: 3,
    projectName: '阿根廷有轨项目',
    reportName: '阿根廷DACU测试报告',
    fileName: '阿根廷DACU测试报告_草稿.xlsx',
    fileUrl: '',
    uploader: '寸诗睿',
    uploadTime: '2026-05-18',
    auditStatus: 'draft',
    auditor: '',
    auditTime: '',
    remark: '草稿，尚未提交审核。'
  },
  {
    id: 4,
    projectName: '波哥大有轨项目',
    reportName: '波哥大乘客报警器测试报告',
    fileName: '波哥大乘客报警器测试报告_V1.1.pdf',
    fileUrl: '',
    uploader: '寸诗睿',
    uploadTime: '2026-05-20',
    auditStatus: 'rejected',
    auditor: '领导',
    auditTime: '2026-05-21',
    remark: '测试结论不完整，需要补充异常断网、SIP重注册、报警恢复测试结果。'
  }
])

const filteredTestReportList = computed(() => {
  return testReportList.value.filter(item => {
    const keywordMatch =
      !filters.keyword ||
      item.projectName.includes(filters.keyword) ||
      item.reportName.includes(filters.keyword) ||
      item.uploader.includes(filters.keyword) ||
      item.fileName.includes(filters.keyword)

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

function openUploadDialog() {
  uploadForm.projectName = ''
  uploadForm.reportName = ''
  uploadForm.fileName = ''
  uploadForm.file = null
  uploadForm.fileUrl = ''
  uploadForm.remark = ''

  showUploadDialog.value = true
}

function handleFileChange(event) {
  const file = event.target.files[0]
  if (!file) return

  uploadForm.file = file
  uploadForm.fileName = file.name
  uploadForm.fileUrl = URL.createObjectURL(file)
}

function uploadTestReport() {
  if (!uploadForm.projectName) {
    alert('请选择测试报告绑定项目')
    return
  }

  if (!uploadForm.reportName) {
    alert('请输入测试报告名称')
    return
  }

  if (!uploadForm.file) {
    alert('请上传测试报告文件')
    return
  }

  testReportList.value.unshift({
    id: Date.now(),
    projectName: uploadForm.projectName,
    reportName: uploadForm.reportName,
    fileName: uploadForm.fileName,
    fileUrl: uploadForm.fileUrl,
    uploader: currentUserName.value,
    uploadTime: new Date().toISOString().slice(0, 10),
    auditStatus: 'draft',
    auditor: '',
    auditTime: '',
    remark: uploadForm.remark
  })

  showUploadDialog.value = false
}

function viewTestReport(item) {
  selectedTestReport.value = item
}

function openTestReportFile(item) {
  if (!item.fileUrl) {
    alert('当前是模拟数据，暂无可直接打开的原始文件')
    return
  }

  window.open(item.fileUrl, '_blank')
}

function downloadTestReport(item) {
  if (!item.fileUrl) {
    alert('当前是模拟数据，暂无可下载的原始文件')
    return
  }

  const link = document.createElement('a')
  link.href = item.fileUrl
  link.download = item.fileName || '测试报告文件'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

function submitTestReport(item) {
  item.auditStatus = 'submitted'
  item.auditor = ''
  item.auditTime = ''
  alert(`测试报告【${item.reportName}】已提交领导审核`)
}

function auditTestReport(item) {
  selectedTestReport.value = item
}

function approveTestReport(item) {
  item.auditStatus = 'approved'
  item.auditor = '领导'
  item.auditTime = new Date().toISOString().slice(0, 10)
  selectedTestReport.value = null
  alert(`测试报告【${item.reportName}】审核通过`)
}

function rejectTestReport(item) {
  item.auditStatus = 'rejected'
  item.auditor = '领导'
  item.auditTime = new Date().toISOString().slice(0, 10)
  selectedTestReport.value = null
  alert(`测试报告【${item.reportName}】已驳回`)
}

function deleteTestReport(item) {
  const ok = confirm(`确认删除测试报告【${item.reportName}】吗？`)
  if (!ok) return

  testReportList.value = testReportList.value.filter(
    record => record.id !== item.id
  )
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
  min-width: 1100px;
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
  width: 220px;
}

.table-card th:nth-child(2),
.table-card td:nth-child(2) {
  width: 180px;
}

.table-card th:nth-child(3),
.table-card td:nth-child(3) {
  width: 120px;
}

.table-card th:nth-child(4),
.table-card td:nth-child(4) {
  width: 130px;
}

.table-card th:nth-child(5),
.table-card td:nth-child(5) {
  width: 130px;
}

.table-card th:nth-child(6),
.table-card td:nth-child(6) {
  width: 120px;
}

.report-link {
  display: inline-block;
  max-width: 200px;
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

.report-link:hover {
  color: #93c5fd;
  text-decoration: underline;
}

.project-tag {
  display: inline-block;
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
  vertical-align: middle;
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

.status-tag.rejected {
  background: #dc262633;
  color: #f87171;
}

.muted {
  color: #94a3b8 !important;
}

.operation-col {
  width: 280px;
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

  .table-card table {
    min-width: 1100px;
  }

  .form-grid,
  .detail-card {
    grid-template-columns: 1fr;
  }
}
</style>