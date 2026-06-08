<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>测试用例管理</h1>
      </div>

      <button class="primary-btn" @click="openUploadDialog">
        上传测试用例
      </button>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索项目名称 / 上传人 / 测试用例文件 / 测试报告文件"
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
              <th>绑定项目</th>
              <th>上传人</th>
              <th>上传时间</th>
              <th>审核状态</th>
              <th>审核人</th>
              <th class="report-operation-col">报告上传/下载</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="item in filteredTestCaseList" :key="item.id">
              <td>
                <span class="project-tag" :title="item.projectName">
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

              <td class="report-operation-col">
                <div class="report-action-group">
                  <button
                    class="text-btn green"
                    @click="openReportUploadDialog(item)"
                  >
                    {{ item.reportFileName ? '更新报告' : '上传报告' }}
                  </button>

                  <button
                    v-if="item.reportFileName"
                    class="text-btn blue"
                    @click="downloadTestReport(item)"
                  >
                    下载报告
                  </button>

                  <span v-else class="empty-report-text">
                    未上传
                  </span>
                </div>
              </td>

              <td class="operation-col">
                <div class="action-group">
                  <button class="text-btn" @click="viewTestCase(item)">
                    查看
                  </button>

                  <button class="text-btn blue" @click="downloadTestCase(item)">
                    下载用例
                  </button>

                  <button
                    v-if="item.auditStatus === 'draft' || item.auditStatus === 'rejected'"
                    class="text-btn blue"
                    @click="submitTestCase(item)"
                  >
                    提交
                  </button>

                  <button
                    v-if="item.auditStatus === 'submitted'"
                    class="text-btn green"
                    @click="auditTestCase(item)"
                  >
                    审核
                  </button>

                  <button
                    v-if="item.auditStatus === 'draft' || item.auditStatus === 'rejected'"
                    class="text-btn red"
                    @click="deleteTestCase(item)"
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
        共 {{ filteredTestCaseList.length }} 条测试用例记录
      </div>
    </div>

    <!-- 上传测试用例弹窗 -->
    <div v-if="showUploadDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>上传测试用例</h3>
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
            测试用例名称
            <input
              v-model="uploadForm.caseName"
              placeholder="例如：香港屯马项目系统测试用例"
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
            测试用例文件
            <input
              type="file"
              accept=".doc,.docx,.xls,.xlsx,.pdf,.zip"
              @change="handleFileChange"
            />
          </label>

          <label class="full-row">
            测试用例说明
            <textarea
              v-model="uploadForm.remark"
              placeholder="例如：包含人工广播、乘客报警、SIP注册、音频播放、网络通信等测试项"
            ></textarea>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showUploadDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="uploadTestCase">
            保存上传
          </button>
        </div>
      </div>
    </div>

    <!-- 上传测试报告弹窗 -->
    <div v-if="showReportUploadDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>上传对应测试报告</h3>
          <button @click="showReportUploadDialog = false">×</button>
        </div>

        <div class="detail-tip">
          <div>
            <span>绑定项目</span>
            <strong>{{ currentReportCase?.projectName }}</strong>
          </div>

          <div>
            <span>测试用例文件</span>
            <strong>{{ currentReportCase?.fileName || '-' }}</strong>
          </div>
        </div>

        <div class="form-grid">
          <label>
            测试报告名称
            <input
              v-model="reportForm.reportName"
              placeholder="例如：香港屯马项目系统测试报告"
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
              @change="handleReportFileChange"
            />
          </label>

          <label class="full-row">
            测试报告说明
            <textarea
              v-model="reportForm.reportRemark"
              placeholder="例如：根据该测试用例执行测试，广播、报警、SIP注册、音频播放等功能测试通过"
            ></textarea>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showReportUploadDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="uploadTestReport">
            保存测试报告
          </button>
        </div>
      </div>
    </div>

    <!-- 查看详情 / 审核弹窗 -->
    <div v-if="selectedTestCase" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>测试用例详情</h3>
          <button @click="selectedTestCase = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>绑定项目</span>
            <strong>{{ selectedTestCase.projectName }}</strong>
          </div>

          <div>
            <span>测试用例文件</span>
            <strong>{{ selectedTestCase.fileName || '-' }}</strong>
          </div>

          <div>
            <span>上传人</span>
            <strong>{{ selectedTestCase.uploader }}</strong>
          </div>

          <div>
            <span>上传时间</span>
            <strong>{{ selectedTestCase.uploadTime }}</strong>
          </div>

          <div>
            <span>审核状态</span>
            <strong>{{ getAuditStatusText(selectedTestCase.auditStatus) }}</strong>
          </div>

          <div>
            <span>审核人</span>
            <strong>{{ selectedTestCase.auditor || '-' }}</strong>
          </div>

          <div>
            <span>审核时间</span>
            <strong>{{ selectedTestCase.auditTime || '-' }}</strong>
          </div>

          <div>
            <span>用例文件查看</span>
            <button class="inline-link" @click="openTestCaseFile(selectedTestCase)">
              点开查看测试用例
            </button>
          </div>
        </div>

        <div class="remark-card">
          <span>测试用例说明</span>
          <p>{{ selectedTestCase.remark || '暂无说明' }}</p>
        </div>

        <div class="report-detail-card">
          <div class="report-detail-header">
            <div>
              <span>对应测试报告</span>
              <strong>
                {{ selectedTestCase.reportFileName ? '已上传' : '未上传' }}
              </strong>
            </div>

            <button
              class="primary-btn small-btn"
              @click="openReportUploadDialog(selectedTestCase)"
            >
              {{ selectedTestCase.reportFileName ? '更新报告' : '上传报告' }}
            </button>
          </div>

          <div v-if="selectedTestCase.reportFileName" class="report-info-grid">
            <div>
              <span>测试报告文件</span>
              <strong>{{ selectedTestCase.reportFileName || '-' }}</strong>
            </div>

            <div>
              <span>报告上传人</span>
              <strong>{{ selectedTestCase.reportUploader || '-' }}</strong>
            </div>

            <div>
              <span>报告上传时间</span>
              <strong>{{ selectedTestCase.reportUploadTime || '-' }}</strong>
            </div>

            <div>
              <span>报告文件查看</span>
              <button class="inline-link" @click="openTestReportFile(selectedTestCase)">
                点开查看测试报告
              </button>
            </div>

            <div class="full-row">
              <span>报告说明</span>
              <p>{{ selectedTestCase.reportRemark || '暂无说明' }}</p>
            </div>
          </div>

          <div v-else class="empty-report-detail">
            当前测试用例暂未上传测试报告。
          </div>
        </div>

        <div class="dialog-footer">
          <button
            v-if="selectedTestCase.auditStatus === 'submitted'"
            class="green-btn"
            @click="approveTestCase(selectedTestCase)"
          >
            审核通过
          </button>

          <button
            v-if="selectedTestCase.auditStatus === 'submitted'"
            class="red-btn"
            @click="rejectTestCase(selectedTestCase)"
          >
            审核驳回
          </button>

          <button class="reset-btn" @click="downloadTestCase(selectedTestCase)">
            下载用例
          </button>

          <button
            v-if="selectedTestCase.reportFileName"
            class="reset-btn"
            @click="downloadTestReport(selectedTestCase)"
          >
            下载报告
          </button>

          <button class="primary-btn" @click="selectedTestCase = null">
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
const showReportUploadDialog = ref(false)

const selectedTestCase = ref(null)
const currentReportCase = ref(null)

const projectOptions = [
  '香港屯马项目',
  '波尔图二期项目',
  '阿根廷有轨项目',
  '波哥大有轨项目',
  '迪拜项目'
]

const uploadForm = reactive({
  projectName: '',
  caseName: '',
  fileName: '',
  file: null,
  fileUrl: '',
  remark: ''
})

const reportForm = reactive({
  reportName: '',
  reportFileName: '',
  reportFile: null,
  reportFileUrl: '',
  reportRemark: ''
})

const testCaseList = ref([
  {
    id: 1,
    projectName: '香港屯马项目',
    caseName: '香港屯马项目系统测试用例',
    fileName: '香港屯马系统测试用例_V1.0.xlsx',
    fileUrl: '',
    uploader: '寸诗睿',
    uploadTime: '2026-05-10',
    auditStatus: 'approved',
    auditor: '领导',
    auditTime: '2026-05-11',
    remark: '覆盖人工广播、乘客报警、客室广播、SIP注册、网络通信等测试项',
    reportName: '香港屯马项目系统测试报告',
    reportFileName: '香港屯马项目系统测试报告_V1.0.docx',
    reportFileUrl: '',
    reportUploader: '寸诗睿',
    reportUploadTime: '2026-05-12',
    reportRemark: '根据系统测试用例执行测试，主要功能测试通过。'
  },
  {
    id: 2,
    projectName: '波尔图二期项目',
    caseName: '波尔图二期广播系统测试用例',
    fileName: '波尔图二期测试用例_V0.9.docx',
    fileUrl: '',
    uploader: '寸诗睿',
    uploadTime: '2026-05-16',
    auditStatus: 'submitted',
    auditor: '',
    auditTime: '',
    remark: '待领导审核确认',
    reportName: '',
    reportFileName: '',
    reportFileUrl: '',
    reportUploader: '',
    reportUploadTime: '',
    reportRemark: ''
  },
  {
    id: 3,
    projectName: '阿根廷有轨项目',
    caseName: '阿根廷DACU测试用例',
    fileName: '阿根廷DACU测试用例_草稿.xlsx',
    fileUrl: '',
    uploader: '寸诗睿',
    uploadTime: '2026-05-18',
    auditStatus: 'draft',
    auditor: '',
    auditTime: '',
    remark: '草稿，尚未提交审核',
    reportName: '',
    reportFileName: '',
    reportFileUrl: '',
    reportUploader: '',
    reportUploadTime: '',
    reportRemark: ''
  },
  {
    id: 4,
    projectName: '波哥大有轨项目',
    caseName: '波哥大乘客报警器测试用例',
    fileName: '波哥大乘客报警器测试用例_V1.1.docx',
    fileUrl: '',
    uploader: '寸诗睿',
    uploadTime: '2026-05-20',
    auditStatus: 'rejected',
    auditor: '领导',
    auditTime: '2026-05-21',
    remark: '测试项不完整，需要补充异常断网、SIP重注册、报警恢复测试',
    reportName: '波哥大乘客报警器测试报告',
    reportFileName: '波哥大乘客报警器测试报告_V1.1.pdf',
    reportFileUrl: '',
    reportUploader: '寸诗睿',
    reportUploadTime: '2026-05-22',
    reportRemark: '报告已上传，但测试项不完整，需要补充后重新提交测试用例。'
  }
])

const filteredTestCaseList = computed(() => {
  return testCaseList.value.filter(item => {
    const keyword = filters.keyword.trim()

    const keywordMatch =
      !keyword ||
      item.projectName.includes(keyword) ||
      item.caseName.includes(keyword) ||
      item.uploader.includes(keyword) ||
      item.fileName.includes(keyword) ||
      (item.reportName && item.reportName.includes(keyword)) ||
      (item.reportFileName && item.reportFileName.includes(keyword)) ||
      (item.reportUploader && item.reportUploader.includes(keyword))

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
  uploadForm.caseName = ''
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

  if (!uploadForm.caseName) {
    uploadForm.caseName = file.name.replace(/\.[^/.]+$/, '')
  }
}

function uploadTestCase() {
  if (!uploadForm.projectName) {
    alert('请选择测试用例绑定项目')
    return
  }

  if (!uploadForm.caseName) {
    alert('请输入测试用例名称')
    return
  }

  if (!uploadForm.file) {
    alert('请上传测试用例文件')
    return
  }

  testCaseList.value.unshift({
    id: Date.now(),
    projectName: uploadForm.projectName,
    caseName: uploadForm.caseName,
    fileName: uploadForm.fileName,
    fileUrl: uploadForm.fileUrl,
    uploader: currentUserName.value,
    uploadTime: new Date().toISOString().slice(0, 10),
    auditStatus: 'draft',
    auditor: '',
    auditTime: '',
    remark: uploadForm.remark,
    reportName: '',
    reportFileName: '',
    reportFileUrl: '',
    reportUploader: '',
    reportUploadTime: '',
    reportRemark: ''
  })

  showUploadDialog.value = false
}

function viewTestCase(item) {
  selectedTestCase.value = item
}

function openTestCaseFile(item) {
  if (!item.fileUrl) {
    alert('当前是模拟数据，暂无可直接打开的原始测试用例文件')
    return
  }

  window.open(item.fileUrl, '_blank')
}

function downloadTestCase(item) {
  if (!item.fileUrl) {
    alert('当前是模拟数据，暂无可下载的原始测试用例文件')
    return
  }

  const link = document.createElement('a')
  link.href = item.fileUrl
  link.download = item.fileName || '测试用例文件'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

function openReportUploadDialog(item) {
  currentReportCase.value = item

  reportForm.reportName = item.reportName || `${item.projectName}测试报告`
  reportForm.reportFileName = ''
  reportForm.reportFile = null
  reportForm.reportFileUrl = ''
  reportForm.reportRemark = item.reportRemark || ''

  showReportUploadDialog.value = true
}

function handleReportFileChange(event) {
  const file = event.target.files[0]
  if (!file) return

  reportForm.reportFile = file
  reportForm.reportFileName = file.name
  reportForm.reportFileUrl = URL.createObjectURL(file)

  if (!reportForm.reportName) {
    reportForm.reportName = file.name.replace(/\.[^/.]+$/, '')
  }
}

function uploadTestReport() {
  if (!currentReportCase.value) return

  if (!reportForm.reportName) {
    alert('请输入测试报告名称')
    return
  }

  if (!reportForm.reportFile) {
    alert('请上传测试报告文件')
    return
  }

  currentReportCase.value.reportName = reportForm.reportName
  currentReportCase.value.reportFileName = reportForm.reportFileName
  currentReportCase.value.reportFileUrl = reportForm.reportFileUrl
  currentReportCase.value.reportUploader = currentUserName.value
  currentReportCase.value.reportUploadTime = new Date().toISOString().slice(0, 10)
  currentReportCase.value.reportRemark = reportForm.reportRemark

  showReportUploadDialog.value = false
  alert(`项目【${currentReportCase.value.projectName}】对应测试报告已上传`)
}

function openTestReportFile(item) {
  if (!item.reportFileUrl) {
    alert('当前是模拟数据，暂无可直接打开的原始测试报告文件')
    return
  }

  window.open(item.reportFileUrl, '_blank')
}

function downloadTestReport(item) {
  if (!item.reportFileName) {
    alert('当前测试用例还没有上传测试报告')
    return
  }

  if (!item.reportFileUrl) {
    alert('当前是模拟数据，暂无可下载的原始测试报告文件')
    return
  }

  const link = document.createElement('a')
  link.href = item.reportFileUrl
  link.download = item.reportFileName || '测试报告文件'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

function submitTestCase(item) {
  item.auditStatus = 'submitted'
  item.auditor = ''
  item.auditTime = ''
  alert(`测试用例【${item.fileName || item.caseName}】已提交领导审核`)
}

function auditTestCase(item) {
  selectedTestCase.value = item
}

function approveTestCase(item) {
  item.auditStatus = 'approved'
  item.auditor = '领导'
  item.auditTime = new Date().toISOString().slice(0, 10)
  selectedTestCase.value = null
  alert(`测试用例【${item.fileName || item.caseName}】审核通过`)
}

function rejectTestCase(item) {
  item.auditStatus = 'rejected'
  item.auditor = '领导'
  item.auditTime = new Date().toISOString().slice(0, 10)
  selectedTestCase.value = null
  alert(`测试用例【${item.fileName || item.caseName}】已驳回`)
}

function deleteTestCase(item) {
  const ok = confirm(`确认删除测试用例【${item.fileName || item.caseName}】吗？`)
  if (!ok) return

  testCaseList.value = testCaseList.value.filter(
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

.small-btn {
  height: 30px;
  padding: 0 12px;
  font-size: 12px;
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
  min-width: 1120px;
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
  width: 110px;
}

.table-card th:nth-child(3),
.table-card td:nth-child(3) {
  width: 120px;
}

.table-card th:nth-child(4),
.table-card td:nth-child(4) {
  width: 120px;
}

.table-card th:nth-child(5),
.table-card td:nth-child(5) {
  width: 100px;
}

.project-tag {
  display: inline-block;
  max-width: 160px;
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

.empty-report-text {
  color: #64748b;
  font-size: 12px;
  white-space: nowrap;
}

.muted {
  color: #94a3b8 !important;
}

.report-operation-col {
  width: 220px;
  text-align: left !important;
}

.operation-col {
  width: 310px;
  text-align: right !important;
}

.action-group {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.report-action-group {
  display: flex;
  justify-content: flex-start;
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
  flex-wrap: wrap;
}

.detail-card {
  padding: 20px;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 14px;
}

.detail-card div,
.remark-card,
.report-detail-card,
.detail-tip {
  background: #020617;
  border: 1px solid #1e293b;
  border-radius: 10px;
  padding: 12px;
}

.detail-card span,
.remark-card span,
.report-detail-card span,
.detail-tip span {
  display: block;
  color: #64748b;
  font-size: 12px;
  margin-bottom: 6px;
}

.detail-card strong,
.detail-tip strong {
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

.detail-tip {
  margin: 20px 20px 0;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.report-detail-card {
  margin: 0 20px 20px;
}

.report-detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 14px;
}

.report-detail-header strong {
  color: #5eead4;
  font-size: 15px;
}

.report-info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.report-info-grid div {
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 10px;
  padding: 12px;
}

.report-info-grid strong {
  color: #f8fafc;
  font-size: 14px;
  word-break: break-all;
}

.report-info-grid p {
  margin: 0;
  color: #cbd5e1;
  font-size: 13px;
  line-height: 1.6;
}

.empty-report-detail {
  padding: 12px;
  color: #64748b;
  font-size: 13px;
  background: #0f172a;
  border: 1px dashed #334155;
  border-radius: 10px;
}

@media (max-width: 960px) {
  .filter-card {
    grid-template-columns: 1fr;
  }

  .table-card table {
    min-width: 1120px;
  }

  .form-grid,
  .detail-card,
  .detail-tip,
  .report-info-grid {
    grid-template-columns: 1fr;
  }

  .report-operation-col {
    width: 220px;
  }

  .operation-col {
    width: 310px;
  }
}
</style>