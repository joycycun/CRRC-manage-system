<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>需求书管理</h1>
      </div>

      <button class="primary-btn" @click="openUploadDialog">
        上传需求书
      </button>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索项目名称 / 需求书名称 / 上传人"
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

      <select v-model="filters.status">
        <option value="">全部状态</option>
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
      <table>
        <thead>
          <tr>
            <th>需求书名称</th>
            <th>对应项目</th>
            <th>上传人</th>
            <th>上传时间</th>
            <th>审核状态</th>
            <th>审核人</th>
            <th class="operation-col">操作</th>
          </tr>
        </thead>

        <tbody>
          <tr v-for="item in filteredRequirementBooks" :key="item.id">
            <td>
              <div class="book-name">{{ item.bookName }}</div>
              <div class="file-name">{{ item.fileName }}</div>
            </td>

            <td>
              <span class="project-tag">{{ item.projectName }}</span>
            </td>

            <td>{{ item.uploader }}</td>

            <td class="muted">{{ item.uploadTime }}</td>

            <td>
              <span class="status-tag" :class="item.status">
                {{ getStatusText(item.status) }}
              </span>
            </td>

            <td>{{ item.auditor || '-' }}</td>

            <td class="operation-col">
              <div class="action-group">
                <button class="text-btn" @click="viewBook(item)">
                  查看
                </button>

                <button class="text-btn blue" @click="downloadBook(item)">
                  下载
                </button>

                <button
                  v-if="item.status === 'draft' || item.status === 'rejected'"
                  class="text-btn blue"
                  @click="submitBook(item)"
                >
                  提交
                </button>

                <button
                  v-if="item.status === 'submitted'"
                  class="text-btn green"
                  @click="auditBook(item)"
                >
                  审核
                </button>

                <button
                  v-if="item.status === 'draft' || item.status === 'rejected'"
                  class="text-btn red"
                  @click="deleteBook(item)"
                >
                  删除
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <div class="table-footer">
        共 {{ filteredRequirementBooks.length }} 条需求书记录
      </div>
    </div>

    <!-- 上传需求书弹窗 -->
    <div v-if="showUploadDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>上传需求书</h3>
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
            需求书名称
            <input
              v-model="uploadForm.bookName"
              placeholder="例如：香港屯马项目需求书"
            />
          </label>

          <label>
            文件名称
            <input
              v-model="uploadForm.fileName"
              placeholder="例如：需求书_V1.0.docx"
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
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showUploadDialog = false">
            取消
          </button>
          <button class="primary-btn" @click="uploadBook">
            保存上传
          </button>
        </div>
      </div>
    </div>
    <!-- 查看需求书详情弹窗 -->
    <div v-if="selectedBook" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>需求书详情</h3>
          <button @click="selectedBook = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>需求书名称</span>
            <strong>{{ selectedBook.bookName }}</strong>
          </div>

          <div>
            <span>对应项目</span>
            <strong>{{ selectedBook.projectName }}</strong>
          </div>

          <div>
            <span>文件名称</span>
            <strong>{{ selectedBook.fileName }}</strong>
          </div>

          <div>
            <span>上传人</span>
            <strong>{{ selectedBook.uploader }}</strong>
          </div>

          <div>
            <span>上传时间</span>
            <strong>{{ selectedBook.uploadTime }}</strong>
          </div>

          <div>
            <span>审核状态</span>
            <strong>{{ getStatusText(selectedBook.status) }}</strong>
          </div>

          <div>
            <span>审核人</span>
            <strong>{{ selectedBook.auditor || '-' }}</strong>
          </div>

          <div>
            <span>审核时间</span>
            <strong>{{ selectedBook.auditTime || '-' }}</strong>
          </div>
        </div>

        <div class="remark-card">
          <span>备注</span>
          <p>{{ selectedBook.remark || '暂无备注' }}</p>
        </div>

        <div class="dialog-footer">
          <button
            v-if="selectedBook.status === 'submitted'"
            class="green-btn"
            @click="approveBook(selectedBook)"
          >
            审核通过
          </button>

          <button
            v-if="selectedBook.status === 'submitted'"
            class="red-btn"
            @click="rejectBook(selectedBook)"
          >
            审核驳回
          </button>

          <button class="primary-btn" @click="selectedBook = null">
            关闭
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'

import { getProjects } from '@/api/project'

import {
  getRequirementBooks,
  createRequirementBook,
  submitRequirementBook,
  auditRequirementBook,
  deleteRequirementBook
} from '@/api/requirement'

const currentUserName = ref(
  localStorage.getItem('username') ||
  localStorage.getItem('accountName') ||
  localStorage.getItem('realName') ||
  '当前用户'
)

const filters = reactive({
  keyword: '',
  projectName: '',
  status: ''
})

const showUploadDialog = ref(false)
const selectedBook = ref(null)

const projectOptions = ref([])
const projectMap = ref({})

const uploadForm = reactive({
  projectName: '',
  bookName: '',
  fileName: '',
  file: null,
  fileUrl: '',
  remark: ''
})

const requirementBookList = ref([])

onMounted(async () => {
  await loadProjects()
  await loadRequirementBooks()
})

function getResponseData(res) {
  if (res && res.data) return res.data
  return res
}

function formatDate(value) {
  if (!value) return ''
  if (typeof value === 'string') return value.slice(0, 10)
  return value
}

function backendStatusToFrontend(status) {
  const map = {
    草稿: 'draft',
    未提交: 'draft',
    待审核: 'submitted',
    已通过: 'approved',
    审核通过: 'approved',
    已驳回: 'rejected',
    审核驳回: 'rejected'
  }

  return map[status] || status || 'draft'
}

function frontendStatusToBackend(status) {
  const map = {
    draft: '草稿',
    submitted: '待审核',
    approved: '已通过',
    rejected: '已驳回'
  }

  return map[status] || status
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

async function loadRequirementBooks() {
  try {
    const res = await getRequirementBooks()
    const result = getResponseData(res)

    console.log('需求书列表返回：', result)

    if (result.code !== 200) {
      alert(result.msg || '加载需求书失败')
      return
    }

    requirementBookList.value = (result.data || []).map(item => ({
      id: item.id,
      projectId: item.projectId,
      projectName: item.projectName || findProjectName(item.projectId),
      bookName: item.bookName || '',
      fileId: item.fileId || 0,
      fileName: item.fileName || uploadFileNameFromItem(item),
      fileUrl: item.fileUrl || '',
      uploader: item.submitUserName || item.uploader || '',
      uploadTime: formatDate(item.submitTime || item.uploadTime || item.createdAt),
      status: backendStatusToFrontend(item.status),
      auditor: item.auditUserName || item.auditor || '',
      auditTime: formatDate(item.auditTime),
      remark: item.remark || ''
    }))
  } catch (err) {
    console.error('加载需求书失败：', err)
    alert('加载需求书失败，请检查后端接口')
  }
}

function uploadFileNameFromItem(item) {
  if (item.fileName) return item.fileName
  if (item.fileDisplayName) return item.fileDisplayName
  if (item.fileId) return `文件ID-${item.fileId}`
  return '暂无文件'
}

const filteredRequirementBooks = computed(() => {
  return requirementBookList.value.filter(item => {
    const keywordMatch =
      !filters.keyword ||
      item.projectName.includes(filters.keyword) ||
      item.bookName.includes(filters.keyword) ||
      item.uploader.includes(filters.keyword)

    const projectMatch =
      !filters.projectName || item.projectName === filters.projectName

    const statusMatch =
      !filters.status || item.status === filters.status

    return keywordMatch && projectMatch && statusMatch
  })
})

function getStatusText(status) {
  const map = {
    draft: '草稿',
    submitted: '待审核',
    approved: '审核通过',
    rejected: '审核驳回'
  }

  return map[status] || status
}

function downloadBook(item) {
  if (!item.fileUrl) {
    alert('当前还没有接真实文件下载，后面做 project_files 文件上传下载时再接')
    return
  }

  const link = document.createElement('a')
  link.href = item.fileUrl
  link.download = item.fileName || '需求书.docx'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
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

function resetFilters() {
  filters.keyword = ''
  filters.projectName = ''
  filters.status = ''
}

function openUploadDialog() {
  uploadForm.projectName = ''
  uploadForm.bookName = ''
  uploadForm.fileName = ''
  uploadForm.file = null
  uploadForm.fileUrl = ''
  uploadForm.remark = ''
  showUploadDialog.value = true
}

async function uploadBook() {
  if (!uploadForm.projectName) {
    alert('请选择需求书对应项目')
    return
  }

  if (!uploadForm.bookName) {
    alert('请输入需求书名称')
    return
  }

  if (!uploadForm.file) {
    alert('请上传 Word 需求书文档')
    return
  }

  const projectId = projectMap.value[uploadForm.projectName]

  if (!projectId) {
    alert('没有找到对应项目ID，请重新选择项目')
    return
  }

  const payload = {
    projectId,
    bookName: uploadForm.bookName,
    fileId: 1,
    status: frontendStatusToBackend('draft'),
    submitUserId: 1,
    submitUserName: currentUserName.value,
    remark: uploadForm.remark || ''
  }

  try {
    const res = await createRequirementBook(payload)
    const result = getResponseData(res)

    console.log('新增需求书返回：', result)

    if (result.code === 200) {
      alert('上传需求书成功')
      showUploadDialog.value = false
      await loadRequirementBooks()
    } else {
      alert(result.msg || '上传需求书失败')
    }
  } catch (err) {
    console.error('上传需求书失败：', err)
    alert('上传需求书失败，请检查后端接口')
  }
}

function viewBook(item) {
  selectedBook.value = item
}

async function submitBook(item) {
  if (!item || !item.id) {
    alert('提交失败：没有拿到需求书ID')
    return
  }

  try {
    const res = await submitRequirementBook(item.id)
    const result = getResponseData(res)

    console.log('提交需求书返回：', result)

    if (result.code === 200) {
      alert(`需求书【${item.bookName}】已提交领导审核`)
      await loadRequirementBooks()
    } else {
      alert(result.msg || '提交失败')
    }
  } catch (err) {
    console.error('提交需求书失败：', err)
    alert('提交需求书失败')
  }
}

function auditBook(item) {
  selectedBook.value = item
}

async function approveBook(item) {
  if (!item || !item.id) {
    alert('审核失败：没有拿到需求书ID')
    return
  }

  try {
    const res = await auditRequirementBook(item.id, {
      auditUserId: 1,
      auditUserName: '领导',
      auditStatus: '已通过',
      rejectReason: ''
    })

    const result = getResponseData(res)

    console.log('审核通过返回：', result)

    if (result.code === 200) {
      alert(`需求书【${item.bookName}】审核通过`)
      selectedBook.value = null
      await loadRequirementBooks()
    } else {
      alert(result.msg || '审核失败')
    }
  } catch (err) {
    console.error('审核需求书失败：', err)
    alert('审核失败')
  }
}

async function rejectBook(item) {
  if (!item || !item.id) {
    alert('审核失败：没有拿到需求书ID')
    return
  }

  const reason = prompt('请输入驳回原因') || '需求描述不完整，请补充后重新提交'

  try {
    const res = await auditRequirementBook(item.id, {
      auditUserId: 1,
      auditUserName: '领导',
      auditStatus: '已驳回',
      rejectReason: reason
    })

    const result = getResponseData(res)

    console.log('审核驳回返回：', result)

    if (result.code === 200) {
      alert(`需求书【${item.bookName}】已驳回`)
      selectedBook.value = null
      await loadRequirementBooks()
    } else {
      alert(result.msg || '审核失败')
    }
  } catch (err) {
    console.error('驳回需求书失败：', err)
    alert('审核失败')
  }
}

async function deleteBook(item) {
  if (!item || !item.id) {
    alert('删除失败：没有拿到需求书ID')
    return
  }

  const ok = confirm(`确认删除需求书【${item.bookName}】吗？`)
  if (!ok) return

  try {
    const res = await deleteRequirementBook(item.id)
    const result = getResponseData(res)

    console.log('删除需求书返回：', result)

    if (result.code === 200) {
      alert('删除成功')
      selectedBook.value = null
      await loadRequirementBooks()
    } else {
      alert(result.msg || '删除失败')
    }
  } catch (err) {
    console.error('删除需求书失败：', err)
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

.book-name {
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
  width: 320px;
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

.action-group {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.action-group .text-btn {
  margin-left: 0;
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
  width: 680px;
  max-width: 100%;
  background: #0f172a;
  border: 1px solid #334155;
  border-radius: 16px;
  color: #f8fafc;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.45);
}

.large-dialog {
  width: 820px;
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
.text-btn.yellow {
  color: #fbbf24;
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
    min-width: 980px;
  }

  .form-grid,
  .detail-card {
    grid-template-columns: 1fr;
  }
}
</style>