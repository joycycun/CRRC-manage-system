<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>需求变更管理</h1>
      </div>

      <button class="primary-btn" @click="openUploadDialog">
        上传需求变更
      </button>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索项目名称 / 变更名称 / 上传人"
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

      <select v-model="filters.closeStatus">
        <option value="">全部关闭状态</option>
        <option value="open">未关闭</option>
        <option value="closed">已关闭</option>
      </select>

      <button class="query-btn">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <!-- 数据表格 -->
    <div class="table-card">
      <table>
        <thead>
          <tr>
            <th>需求变更名称</th>
            <th>对应项目</th>
            <th>上传人</th>
            <th>上传时间</th>
            <th>审核状态</th>
            <th>关闭状态</th>
            <th>审核人</th>
            <th class="operation-col">操作</th>
          </tr>
        </thead>

        <tbody>
          <tr v-for="item in filteredChangeList" :key="item.id">
            <td>
              <div class="change-name">{{ item.changeName }}</div>
              <div class="file-name">{{ item.fileName }}</div>
            </td>

            <td>
              <span class="project-tag">{{ item.projectName }}</span>
            </td>

            <td>{{ item.uploader }}</td>

            <td class="muted">{{ item.uploadTime }}</td>

            <td>
              <span class="status-tag" :class="item.auditStatus">
                {{ getAuditStatusText(item.auditStatus) }}
              </span>
            </td>

            <td>
              <span class="close-tag" :class="item.closeStatus">
                {{ getCloseStatusText(item.closeStatus) }}
              </span>
            </td>

            <td>{{ item.auditor || '-' }}</td>

            <td class="operation-col">
              <div class="action-group">
                <button class="text-btn" @click="viewChange(item)">
                  查看
                </button>

                <button class="text-btn blue" @click="downloadChange(item)">
                  下载
                </button>

                <button
                  v-if="item.auditStatus === 'draft' || item.auditStatus === 'rejected'"
                  class="text-btn blue"
                  @click="submitChange(item)"
                >
                  提交
                </button>

                <button
                  v-if="item.auditStatus === 'submitted'"
                  class="text-btn green"
                  @click="auditChange(item)"
                >
                  审核
                </button>

                <button
                  v-if="item.auditStatus === 'approved' && item.closeStatus === 'open'"
                  class="text-btn purple"
                  @click="closeChange(item)"
                >
                  关闭
                </button>

                <button
                  v-if="item.auditStatus === 'draft' || item.auditStatus === 'rejected'"
                  class="text-btn red"
                  @click="deleteChange(item)"
                >
                  删除
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <div class="table-footer">
        共 {{ filteredChangeList.length }} 条需求变更记录
      </div>
    </div>

    <!-- 上传需求变更弹窗 -->
    <div v-if="showUploadDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>上传需求变更</h3>
          <button @click="showUploadDialog = false">×</button>
        </div>

        <div class="form-grid">
          <label>
            对应项目
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
            需求变更名称
            <input
              v-model="uploadForm.changeName"
              placeholder="例如：香港屯马项目需求变更"
            />
          </label>

          <label>
            上传人
            <input
              v-model="uploadForm.uploader"
              placeholder="请输入上传人"
            />
          </label>

          <label>
            文件名称
            <input
              v-model="uploadForm.fileName"
              placeholder="例如：需求变更_V1.0.docx"
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
            变更说明
            <textarea
              v-model="uploadForm.remark"
              placeholder="请输入本次需求变更的原因、影响范围、涉及终端类型等"
            ></textarea>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showUploadDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="uploadChange">
            保存上传
          </button>
        </div>
      </div>
    </div>

    <!-- 查看详情 / 审核弹窗 -->
    <div v-if="selectedChange" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>需求变更详情</h3>
          <button @click="selectedChange = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>需求变更名称</span>
            <strong>{{ selectedChange.changeName }}</strong>
          </div>

          <div>
            <span>对应项目</span>
            <strong>{{ selectedChange.projectName }}</strong>
          </div>

          <div>
            <span>文件名称</span>
            <strong>{{ selectedChange.fileName }}</strong>
          </div>

          <div>
            <span>上传人</span>
            <strong>{{ selectedChange.uploader }}</strong>
          </div>

          <div>
            <span>上传时间</span>
            <strong>{{ selectedChange.uploadTime }}</strong>
          </div>

          <div>
            <span>审核状态</span>
            <strong>{{ getAuditStatusText(selectedChange.auditStatus) }}</strong>
          </div>

          <div>
            <span>关闭状态</span>
            <strong>{{ getCloseStatusText(selectedChange.closeStatus) }}</strong>
          </div>

          <div>
            <span>审核人</span>
            <strong>{{ selectedChange.auditor || '-' }}</strong>
          </div>

          <div>
            <span>审核时间</span>
            <strong>{{ selectedChange.auditTime || '-' }}</strong>
          </div>

          <div>
            <span>关闭人</span>
            <strong>{{ selectedChange.closeUser || '-' }}</strong>
          </div>
        </div>

        <div class="remark-card">
          <span>变更说明</span>
          <p>{{ selectedChange.remark || '暂无说明' }}</p>
        </div>

        <div class="dialog-footer">
          <button
            v-if="selectedChange.auditStatus === 'submitted'"
            class="green-btn"
            @click="approveChange(selectedChange)"
          >
            审核通过
          </button>

          <button
            v-if="selectedChange.auditStatus === 'submitted'"
            class="red-btn"
            @click="rejectChange(selectedChange)"
          >
            审核驳回
          </button>

          <button
            v-if="selectedChange.auditStatus === 'approved' && selectedChange.closeStatus === 'open'"
            class="purple-btn"
            @click="closeChange(selectedChange)"
          >
            研发关闭
          </button>

          <button class="primary-btn" @click="selectedChange = null">
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
  auditStatus: '',
  closeStatus: ''
})

const showUploadDialog = ref(false)


const selectedChange = ref(null)


const projectOptions = [
  '香港屯马项目',
  '波尔图二期项目',
  '阿根廷有轨项目',
  '波哥大有轨项目',
  '成都项目'
]

const uploadForm = reactive({
  projectName: '',
  changeName: '',
  uploader: '',
  fileName: '',
  file: null,
  fileUrl: '',
  remark: ''
})


const changeList = ref([
  {
    id: 1,
    projectName: '香港屯马项目',
    changeName: '香港屯马项目需求变更',
    fileName: '香港屯马需求变更_V1.0.docx',
    fileUrl: '',
    uploader: '寸诗睿',
    uploadTime: '2026-05-10',
    auditStatus: 'approved',
    closeStatus: 'closed',
    auditor: '领导',
    auditTime: '2026-05-11',
    closeUser: '研发',
    closeTime: '2026-05-12',
    remark: '已完成变更研发处理，关闭需求变更'
  },
  {
    id: 2,
    projectName: '波尔图二期项目',
    changeName: '波尔图二期广播逻辑变更',
    fileName: '波尔图需求变更_V0.9.docx',
    fileUrl: '',
    uploader: '寸诗睿',
    uploadTime: '2026-05-16',
    auditStatus: 'submitted',
    closeStatus: 'open',
    auditor: '',
    auditTime: '',
    closeUser: '',
    closeTime: '',
    remark: '待领导审核确认'
  },
  {
    id: 3,
    projectName: '阿根廷有轨项目',
    changeName: '阿根廷报警器需求变更',
    fileName: '阿根廷需求变更_草稿.docx',
    fileUrl: '',
    uploader: '寸诗睿',
    uploadTime: '2026-05-18',
    auditStatus: 'draft',
    closeStatus: 'open',
    auditor: '',
    auditTime: '',
    closeUser: '',
    closeTime: '',
    remark: '草稿，尚未提交审核'
  },
  {
    id: 4,
    projectName: '波哥大有轨项目',
    changeName: '波哥大编码板需求变更',
    fileName: '波哥大需求变更_V1.1.docx',
    fileUrl: '',
    uploader: '寸诗睿',
    uploadTime: '2026-05-20',
    auditStatus: 'rejected',
    closeStatus: 'open',
    auditor: '领导',
    auditTime: '2026-05-21',
    closeUser: '',
    closeTime: '',
    remark: '变更影响范围不清晰，需要补充涉及终端类型'
  }
])

const filteredChangeList = computed(() => {
  return changeList.value.filter(item => {
    const keywordMatch =
      !filters.keyword ||
      item.projectName.includes(filters.keyword) ||
      item.changeName.includes(filters.keyword) ||
      item.uploader.includes(filters.keyword)

    const projectMatch =
      !filters.projectName || item.projectName === filters.projectName

    const auditStatusMatch =
      !filters.auditStatus || item.auditStatus === filters.auditStatus

    const closeStatusMatch =
      !filters.closeStatus || item.closeStatus === filters.closeStatus

    return keywordMatch && projectMatch && auditStatusMatch && closeStatusMatch
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

function getCloseStatusText(status) {
  const map = {
    open: '未关闭',
    closed: '已关闭'
  }

  return map[status] || status
}

function resetFilters() {
  filters.keyword = ''
  filters.projectName = ''
  filters.auditStatus = ''
  filters.closeStatus = ''
}

function openUploadDialog() {
  uploadForm.projectName = ''
  uploadForm.changeName = ''
  uploadForm.uploader = ''
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

function uploadChange() {
  if (!uploadForm.projectName) {
    alert('请选择需求变更对应项目')
    return
  }

  if (!uploadForm.changeName) {
    alert('请输入需求变更名称')
    return
  }

  if (!uploadForm.file) {
    alert('请上传 Word 需求变更文档')
    return
  }

  changeList.value.unshift({
    id: Date.now(),
    projectName: uploadForm.projectName,
    changeName: uploadForm.changeName,
    fileName: uploadForm.fileName,
    fileUrl: uploadForm.fileUrl,
    uploader: uploadForm.uploader || '当前用户',
    uploadTime: new Date().toISOString().slice(0, 10),
    auditStatus: 'draft',
    closeStatus: 'open',
    auditor: '',
    auditTime: '',
    closeUser: '',
    closeTime: '',
    remark: uploadForm.remark
  })

  showUploadDialog.value = false
}

function viewChange(item) {
  selectedChange.value = item
}

function downloadChange(item) {
  if (!item.fileUrl) {
    alert('当前是模拟数据，暂无可下载的 Word 文件')
    return
  }

  const link = document.createElement('a')
  link.href = item.fileUrl
  link.download = item.fileName || '需求变更.docx'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

function submitChange(item) {
  item.auditStatus = 'submitted'
  item.auditor = ''
  item.auditTime = ''
  alert(`需求变更【${item.changeName}】已提交领导审核`)
}

function auditChange(item) {
  selectedChange.value = item
}

function approveChange(item) {
  item.auditStatus = 'approved'
  item.auditor = '领导'
  item.auditTime = new Date().toISOString().slice(0, 10)
  selectedChange.value = null
  alert(`需求变更【${item.changeName}】审核通过`)
}

function rejectChange(item) {
  item.auditStatus = 'rejected'
  item.auditor = '领导'
  item.auditTime = new Date().toISOString().slice(0, 10)
  selectedChange.value = null
  alert(`需求变更【${item.changeName}】已驳回`)
}

function closeChange(item) {
  item.closeStatus = 'closed'
  item.closeUser = '研发'
  item.closeTime = new Date().toISOString().slice(0, 10)
  selectedChange.value = null
  alert(`研发已关闭需求变更【${item.changeName}】`)
}


function deleteChange(item) {
  const ok = confirm(`确认删除需求变更【${item.changeName}】吗？`)
  if (!ok) return

  changeList.value = changeList.value.filter(change => change.id !== item.id)
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
.red-btn,
.purple-btn {
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

.purple-btn {
  border: 1px solid #6d28d9;
  background: #2e1065;
  color: #c4b5fd;
}

.filter-card {
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 14px;
  padding: 16px;
  display: grid;
  grid-template-columns: 1.4fr 220px 180px 180px 90px 90px;
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

.change-name {
  color: #f8fafc;
  font-weight: 700;
}

.file-name {
  margin-top: 4px;
  color: #64748b;
  font-size: 12px;
  word-break: break-all;
}

.project-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  background: #1d4ed833;
  color: #60a5fa;
  font-size: 12px;
  font-weight: 700;
}

.status-tag,
.close-tag {
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

.close-tag.open {
  background: #1d4ed833;
  color: #60a5fa;
}

.close-tag.closed {
  background: #64748b33;
  color: #cbd5e1;
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

.text-btn.purple {
  color: #c4b5fd;
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
  width: 680px;
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
}
</style>