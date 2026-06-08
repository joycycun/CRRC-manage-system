<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>硬件测试管理</h1>
      </div>

      <button class="primary-btn" @click="openUploadDialog">
        上传测试记录
      </button>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索项目名称 / 测试记录名称 / 上传人 / 终端类型"
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
              <th>终端类型</th>
              <th>测试记录名称</th>
              <th>绑定项目</th>
              <th>硬件版本</th>
              <th>上传人</th>
              <th>上传时间</th>
              <th>审核状态</th>
              <th>审核人</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="item in filteredTestList" :key="item.id">
              <td>
                <span class="device-tag">{{ item.deviceType }}</span>
              </td>

              <td>
                <div class="record-name">{{ item.recordName }}</div>
                <div class="file-name">{{ item.fileName }}</div>
              </td>

              <td>
                <span class="project-tag">{{ item.projectName }}</span>
              </td>

              <td>
                <span class="version-tag">{{ item.hardwareVersion }}</span>
              </td>

              <td>{{ item.uploader }}</td>

              <td class="muted">{{ item.uploadTime }}</td>

              <td>
                <div class="audit-cell">
                  <span class="status-tag" :class="item.auditStatus">
                    {{ getAuditStatusText(item.auditStatus) }}
                  </span>

                  <button
                    v-if="item.auditStatus === 'rejected'"
                    class="reason-btn"
                    @click="viewRejectReason(item)"
                  >
                    原因
                  </button>
                </div>
              </td>

              <td>{{ item.auditor || '-' }}</td>

              <td class="operation-col">
                <div class="action-group">
                  <button class="text-btn" @click="viewTest(item)">
                    查看
                  </button>

                  <button class="text-btn blue" @click="downloadTest(item)">
                    下载
                  </button>

                  <button
                    v-if="item.auditStatus === 'draft' || item.auditStatus === 'rejected'"
                    class="text-btn blue"
                    @click="submitTest(item)"
                  >
                    提交
                  </button>

                  <button
                    v-if="item.auditStatus === 'submitted'"
                    class="text-btn green"
                    @click="auditTest(item)"
                  >
                    审核
                  </button>

                  <button
                    v-if="item.auditStatus === 'draft' || item.auditStatus === 'rejected'"
                    class="text-btn red"
                    @click="deleteTest(item)"
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
        共 {{ filteredTestList.length }} 条硬件测试记录
      </div>
    </div>

    <!-- 上传硬件测试记录弹窗 -->
    <div v-if="showUploadDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>上传硬件测试记录</h3>
          <button @click="showUploadDialog = false">×</button>
        </div>

        <div class="form-grid">
          <label>
            终端类型
            <select v-model="uploadForm.deviceType">
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
            测试记录名称
            <input
              v-model="uploadForm.recordName"
              placeholder="例如：香港屯马硬件测试记录"
            />
          </label>

          <label>
            硬件版本
            <select v-model="uploadForm.hardwareVersion">
              <option value="">请选择硬件版本</option>
              <option
                v-for="version in hardwareVersionOptions"
                :key="version"
                :value="version"
              >
                {{ version }}
              </option>
            </select>
          </label>

          <label>
            文件名称
            <input
              v-model="uploadForm.fileName"
              placeholder="选择 Word 文件后自动填充"
              disabled
            />
          </label>

          <label class="full-row">
            Word 文档
            <input
              type="file"
              accept=".doc,.docx,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document"
              @change="handleFileChange"
            />
          </label>

          <label class="full-row">
            测试说明
            <textarea
              v-model="uploadForm.remark"
              placeholder="请输入硬件测试内容、测试结论、问题说明等"
            ></textarea>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showUploadDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="uploadTest">
            保存上传
          </button>
        </div>
      </div>
    </div>

    <!-- 查看详情 / 审核弹窗 -->
    <div v-if="selectedTest" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>硬件测试记录详情</h3>
          <button @click="selectedTest = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>测试记录名称</span>
            <strong>{{ selectedTest.recordName }}</strong>
          </div>

          <div>
            <span>绑定项目</span>
            <strong>{{ selectedTest.projectName }}</strong>
          </div>

          <div>
            <span>终端类型</span>
            <strong>{{ selectedTest.deviceType }}</strong>
          </div>

          <div>
            <span>硬件版本</span>
            <strong>{{ selectedTest.hardwareVersion }}</strong>
          </div>

          <div>
            <span>文件名称</span>
            <strong>{{ selectedTest.fileName }}</strong>
          </div>

          <div>
            <span>上传人</span>
            <strong>{{ selectedTest.uploader }}</strong>
          </div>

          <div>
            <span>上传时间</span>
            <strong>{{ selectedTest.uploadTime }}</strong>
          </div>

          <div>
            <span>审核状态</span>
            <strong>{{ getAuditStatusText(selectedTest.auditStatus) }}</strong>
          </div>

          <div>
            <span>审核人</span>
            <strong>{{ selectedTest.auditor || '-' }}</strong>
          </div>

          <div>
            <span>审核时间</span>
            <strong>{{ selectedTest.auditTime || '-' }}</strong>
          </div>

          <div v-if="selectedTest.auditStatus === 'rejected'">
            <span>驳回原因</span>
            <button class="inline-link" @click="viewRejectReason(selectedTest)">
              查看驳回原因
            </button>
          </div>
        </div>

        <div class="remark-card">
          <span>测试说明</span>
          <p>{{ selectedTest.remark || '暂无说明' }}</p>
        </div>

        <div class="dialog-footer">
          <button
            v-if="selectedTest.auditStatus === 'submitted'"
            class="green-btn"
            @click="approveTest(selectedTest)"
          >
            审核通过
          </button>

          <button
            v-if="selectedTest.auditStatus === 'submitted'"
            class="red-btn"
            @click="openRejectDialog(selectedTest)"
          >
            审核驳回
          </button>

          <button class="primary-btn" @click="selectedTest = null">
            关闭
          </button>
        </div>
      </div>
    </div>

    <!-- 审核驳回原因填写弹窗 -->
    <div v-if="showRejectDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>填写驳回原因</h3>
          <button @click="showRejectDialog = false">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>测试记录名称</span>
            <strong>{{ currentRejectTest?.recordName }}</strong>
          </div>

          <div>
            <span>绑定项目</span>
            <strong>{{ currentRejectTest?.projectName }}</strong>
          </div>

          <div>
            <span>终端类型</span>
            <strong>{{ currentRejectTest?.deviceType }}</strong>
          </div>

          <div>
            <span>硬件版本</span>
            <strong>{{ currentRejectTest?.hardwareVersion }}</strong>
          </div>
        </div>

        <div class="form-grid">
          <label class="full-row">
            驳回原因
            <textarea
              v-model="rejectForm.reason"
              placeholder="请填写驳回原因，例如：测试结论不完整、缺少测试数据、文件版本错误、测试项未覆盖等"
            ></textarea>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showRejectDialog = false">
            取消
          </button>

          <button class="red-btn" @click="confirmRejectTest">
            确认驳回
          </button>
        </div>
      </div>
    </div>

    <!-- 查看驳回原因弹窗 -->
    <div v-if="selectedRejectReason" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>驳回原因</h3>
          <button @click="selectedRejectReason = null">×</button>
        </div>

        <div class="remark-card reject-reason-card">
          <span>详细原因</span>
          <p>{{ selectedRejectReason }}</p>
        </div>

        <div class="dialog-footer">
          <button class="primary-btn" @click="selectedRejectReason = null">
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
const showRejectDialog = ref(false)

const selectedTest = ref(null)
const currentRejectTest = ref(null)
const selectedRejectReason = ref(null)

const projectOptions = [
  '香港屯马项目',
  '波尔图二期项目',
  '阿根廷有轨项目',
  '波哥大有轨项目',
  '成都项目'
]

const hardwareVersionOptions = [
  'HD-CRRC-HKTM.01.V1.1.0',
  'HD-CRRC-AGTB-04.T1.1.0',
  'HD-CRRC-BOGT-03.T1.1.0',
  'HD-CRRC-DUBAI-05.S1.1.0'
]

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

const uploadForm = reactive({
  projectName: '',
  recordName: '',
  hardwareVersion: '',
  deviceType: '',
  fileName: '',
  file: null,
  fileUrl: '',
  remark: ''
})

const rejectForm = reactive({
  reason: ''
})

const hardwareTestList = ref([
  {
    id: 1,
    projectName: '香港屯马项目',
    recordName: '香港屯马广播控制盒硬件测试记录',
    hardwareVersion: 'HD-CRRC-HKTM',
    deviceType: '广播控制盒',
    fileName: '香港屯马硬件测试记录_V1.0.docx',
    fileUrl: '',
    uploader: '郑宇',
    uploadTime: '2026-05-10',
    auditStatus: 'approved',
    auditor: '领导',
    auditTime: '2026-05-11',
    rejectReason: '',
    remark: '硬件功能测试通过，音频输入输出、按键、网络通信测试正常'
  },
  {
    id: 2,
    projectName: '波尔图二期项目',
    recordName: '波尔图客室解码板硬件测试记录',
    hardwareVersion: 'HD-CRRC-POR2',
    deviceType: '客室解码板',
    fileName: '波尔图硬件测试记录_V0.9.docx',
    fileUrl: '',
    uploader: '郑宇',
    uploadTime: '2026-05-16',
    auditStatus: 'submitted',
    auditor: '',
    auditTime: '',
    rejectReason: '',
    remark: '待领导审核确认'
  },
  {
    id: 3,
    projectName: '阿根廷有轨项目',
    recordName: '阿根廷乘客报警器硬件测试记录',
    hardwareVersion: 'HD-CRRC-AGTB',
    deviceType: '乘客报警器',
    fileName: '阿根廷硬件测试记录_草稿.docx',
    fileUrl: '',
    uploader: '郑宇',
    uploadTime: '2026-05-18',
    auditStatus: 'draft',
    auditor: '',
    auditTime: '',
    rejectReason: '',
    remark: '草稿，尚未提交审核'
  },
  {
    id: 4,
    projectName: '波哥大有轨项目',
    recordName: '波哥大编码板硬件测试记录',
    hardwareVersion: 'HD-CRRC-BOGT',
    deviceType: '编码板',
    fileName: '波哥大硬件测试记录_V1.1.docx',
    fileUrl: '',
    uploader: '郑宇',
    uploadTime: '2026-05-20',
    auditStatus: 'rejected',
    auditor: '领导',
    auditTime: '2026-05-21',
    rejectReason: '测试结论不完整，需要补充接口测试结果、测试截图和最终结论。',
    remark: '测试结论不完整，需要补充接口测试结果'
  }
])

const filteredTestList = computed(() => {
  return hardwareTestList.value.filter(item => {
    const keywordMatch =
      !filters.keyword ||
      item.projectName.includes(filters.keyword) ||
      item.recordName.includes(filters.keyword) ||
      item.uploader.includes(filters.keyword) ||
      item.hardwareVersion.includes(filters.keyword) ||
      item.deviceType.includes(filters.keyword) ||
      (item.rejectReason && item.rejectReason.includes(filters.keyword))

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
  uploadForm.recordName = ''
  uploadForm.hardwareVersion = ''
  uploadForm.deviceType = ''
  uploadForm.fileName = ''
  uploadForm.file = null
  uploadForm.fileUrl = ''
  uploadForm.remark = ''
  showUploadDialog.value = true
}

function handleFileChange(event) {
  const file = event.target.files[0]
  if (!file) return

  const isWord =
    file.name.endsWith('.doc') ||
    file.name.endsWith('.docx')

  if (!isWord) {
    alert('只能上传 Word 文档，格式为 .doc 或 .docx')
    event.target.value = ''
    return
  }

  uploadForm.file = file
  uploadForm.fileName = file.name
  uploadForm.fileUrl = URL.createObjectURL(file)
}

function uploadTest() {
  if (!uploadForm.deviceType) {
    alert('请选择终端类型')
    return
  }

  if (!uploadForm.projectName) {
    alert('请选择硬件测试记录绑定项目')
    return
  }

  if (!uploadForm.recordName) {
    alert('请输入硬件测试记录名称')
    return
  }

  if (!uploadForm.hardwareVersion) {
    alert('请选择硬件版本')
    return
  }

  if (!uploadForm.file) {
    alert('请上传 Word 硬件测试记录文档')
    return
  }

  hardwareTestList.value.unshift({
    id: Date.now(),
    projectName: uploadForm.projectName,
    recordName: uploadForm.recordName,
    hardwareVersion: uploadForm.hardwareVersion,
    deviceType: uploadForm.deviceType,
    fileName: uploadForm.fileName,
    fileUrl: uploadForm.fileUrl,
    uploader: currentUserName.value,
    uploadTime: new Date().toISOString().slice(0, 10),
    auditStatus: 'draft',
    auditor: '',
    auditTime: '',
    rejectReason: '',
    remark: uploadForm.remark
  })

  showUploadDialog.value = false
}

function viewTest(item) {
  selectedTest.value = item
}

function downloadTest(item) {
  if (!item.fileUrl) {
    alert('当前是模拟数据，暂无可下载的 Word 文件')
    return
  }

  const link = document.createElement('a')
  link.href = item.fileUrl
  link.download = item.fileName || '硬件测试记录.docx'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

function submitTest(item) {
  item.auditStatus = 'submitted'
  item.auditor = ''
  item.auditTime = ''
  item.rejectReason = ''
  alert(`硬件测试记录【${item.recordName}】已提交领导审核`)
}

function auditTest(item) {
  selectedTest.value = item
}

function approveTest(item) {
  item.auditStatus = 'approved'
  item.auditor = '领导'
  item.auditTime = new Date().toISOString().slice(0, 10)
  item.rejectReason = ''
  selectedTest.value = null
  alert(`硬件测试记录【${item.recordName}】审核通过`)
}

function openRejectDialog(item) {
  currentRejectTest.value = item
  rejectForm.reason = item.rejectReason || ''
  showRejectDialog.value = true
}

function confirmRejectTest() {
  if (!currentRejectTest.value) return

  if (!rejectForm.reason) {
    alert('请填写驳回原因')
    return
  }

  currentRejectTest.value.auditStatus = 'rejected'
  currentRejectTest.value.auditor = '领导'
  currentRejectTest.value.auditTime = new Date().toISOString().slice(0, 10)
  currentRejectTest.value.rejectReason = rejectForm.reason

  showRejectDialog.value = false
  selectedTest.value = null

  alert(`硬件测试记录【${currentRejectTest.value.recordName}】已驳回`)
}

function viewRejectReason(item) {
  selectedRejectReason.value = item.rejectReason || '暂无驳回原因'
}

function deleteTest(item) {
  const ok = confirm(`确认删除硬件测试记录【${item.recordName}】吗？`)
  if (!ok) return

  hardwareTestList.value = hardwareTestList.value.filter(
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

.page-header p {
  margin: 8px 0 0;
  color: #94a3b8;
  font-size: 14px;
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
  min-width: 1280px;
  border-collapse: collapse;
  table-layout: fixed;
}

.table-card th:nth-child(7),
.table-card td:nth-child(7) {
  width: 150px;
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

.record-name {
  color: #f8fafc;
  font-weight: 700;
}

.file-name {
  margin-top: 4px;
  color: #64748b;
  font-size: 12px;
  word-break: break-all;
}

.project-tag,
.version-tag,
.device-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
}

.project-tag {
  background: #1d4ed833;
  color: #60a5fa;
}

.version-tag {
  background: #9333ea33;
  color: #c084fc;
}

.device-tag {
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

.audit-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 120px;
  white-space: nowrap;
}

.reason-btn {
  height: 22px;
  padding: 0 8px;
  border: 1px solid #7f1d1d;
  border-radius: 999px;
  background: #450a0a;
  color: #fca5a5;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  white-space: nowrap;
}

.reason-btn:hover {
  background: #7f1d1d;
  color: #fee2e2;
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

.reject-reason-card {
  margin-top: 20px;
}

.reject-reason-card p {
  white-space: pre-wrap;
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
    min-width: 1200px;
  }

  .form-grid,
  .detail-card {
    grid-template-columns: 1fr;
  }
}
</style>