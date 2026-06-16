<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>需求变更管理</h1>
      </div>

      <button v-if="canUseAction('requirement:upload')" class="primary-btn" @click="openUploadDialog">
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

      <button class="query-btn" @click="loadRequirementChanges">查询</button>
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
                  v-if="canUseAction('requirement:submit') && (item.auditStatus === 'draft' || item.auditStatus === 'rejected')"
                  class="text-btn blue"
                  @click="submitChange(item)"
                >
                  提交
                </button>

                <button
                  v-if="canUseAction('requirement:audit') && item.auditStatus === 'submitted'"
                  class="text-btn green"
                  @click="auditChange(item)"
                >
                  审核
                </button>

                <button
                  v-if="canUseAction('requirement:close') && item.auditStatus === 'approved' && item.closeStatus === 'open'"
                  class="text-btn purple"
                  @click="closeChange(item)"
                >
                  关闭
                </button>

                <button
                  v-if="canUseAction('requirement:delete') && (item.auditStatus === 'draft' || item.auditStatus === 'rejected')"
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
              placeholder="例如：波哥大有轨需求变更V1.0"
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
            v-if="canUseAction('requirement:audit') && selectedChange.auditStatus === 'submitted'"
            class="green-btn"
            @click="approveChange(selectedChange)"
          >
            审核通过
          </button>

          <button
            v-if="canUseAction('requirement:audit') && selectedChange.auditStatus === 'submitted'"
            class="red-btn"
            @click="rejectChange(selectedChange)"
          >
            审核驳回
          </button>

          <button
            v-if="canUseAction('requirement:close') && selectedChange.auditStatus === 'approved' && selectedChange.closeStatus === 'open'"
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
import { computed, onMounted, reactive, ref } from 'vue'
import { canUseAction } from '@/utils/permission'

import { getProjects } from '@/api/project'

import {
  getRequirementChanges,
  createRequirementChange,
  submitRequirementChange,
  auditRequirementChange,
  closeRequirementChange,
  deleteRequirementChange
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
  auditStatus: '',
  closeStatus: ''
})

const showUploadDialog = ref(false)
const selectedChange = ref(null)

const projectOptions = ref([])
const projectMap = ref({})

const uploadForm = reactive({
  projectName: '',
  changeName: '',
  changeType: '功能变更',
  fileName: '',
  file: null,
  fileUrl: '',
  remark: ''
})

const requirementChangeList = ref([])

onMounted(async () => {
  await loadProjects()
  await loadRequirementChanges()
})

function getResponseData(res) {
  if (res && res.data) return res.data
  return res
}

function formatDate(value) {
  if (!value) return ''

  // 后端返回普通字符串
  if (typeof value === 'string') {
    return value.slice(0, 10)
  }

  // 后端返回 Go 的 sql.NullTime：{ Time: "...", Valid: true }
  if (typeof value === 'object') {
    if (value.Valid === false) return ''

    const timeValue = value.Time || value.time
    if (timeValue) {
      return String(timeValue).slice(0, 10)
    }
  }

  return ''
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

function backendCloseStatusToFrontend(status) {
  const map = {
    未关闭: 'open',
    待关闭: 'open',
    已关闭: 'closed',
    已闭环: 'closed'
  }

  return map[status] || status || 'open'
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

async function loadRequirementChanges() {
  try {
    const res = await getRequirementChanges()
    const result = getResponseData(res)

    console.log('需求变更列表返回：', result)

    if (result.code !== 200) {
      alert(result.msg || '加载需求变更失败')
      return
    }

    requirementChangeList.value = (result.data || []).map(item => {
      const changeName = item.changeName || item.changeTitle || ''
      const submitUser = item.submitUserName || item.submitUser || item.uploader || ''
      const uploadTime = formatDate(item.submitTime || item.uploadTime || item.createdAt)

      return {
        id: item.id,
        projectId: item.projectId,
        projectName: item.projectName || findProjectName(item.projectId),

        changeName,
        changeTitle: changeName,
        changeType: item.changeType || '',

        fileId: item.fileId || 0,
        fileName:
          item.fileName ||
          item.fileDisplayName ||
          (item.fileId ? `文件ID-${item.fileId}` : '暂无文件'),
        fileUrl: item.fileUrl || '',

        submitUserId: item.submitUserId || 0,
        submitUserName: submitUser,
        submitUser,
        uploader: submitUser,

        submitTime: uploadTime,
        uploadTime,

        status: backendStatusToFrontend(item.status),
        auditStatus: backendStatusToFrontend(item.status),

        closeStatus: backendCloseStatusToFrontend(item.closeStatus),

        auditUserName: item.auditUserName || '',
        auditor: item.auditUserName || item.auditor || '',
        auditTime: formatDate(item.auditTime),

        closeUserName: item.closeUserName || '',
        closeUser: item.closeUserName || item.closeUser || '',
        closeTime: formatDate(item.closeTime),

        rejectReason: item.rejectReason || '',
        remark: item.remark || ''
      }
    })
  } catch (err) {
    console.error('加载需求变更失败：', err)
    alert('加载需求变更失败，请检查后端接口')
  }
}

const filteredChangeList = computed(() => {
  const list = requirementChangeList.value || []

  return list.filter(item => {
    const keyword = filters.keyword || ''

    const keywordMatch =
      !keyword ||
      String(item.projectName || '').includes(keyword) ||
      String(item.changeName || '').includes(keyword) ||
      String(item.changeType || '').includes(keyword) ||
      String(item.uploader || '').includes(keyword)

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
  uploadForm.changeType = '功能变更'
  uploadForm.fileName = ''
  uploadForm.file = null
  uploadForm.fileUrl = ''
  uploadForm.remark = ''

  showUploadDialog.value = true
}

function handleFileChange(event) {
  const file = event.target.files[0]
  if (!file) return

  const isAllowed =
    file.name.endsWith('.doc') ||
    file.name.endsWith('.docx') ||
    file.name.endsWith('.pdf') ||
    file.name.endsWith('.xlsx') ||
    file.name.endsWith('.xls')

  if (!isAllowed) {
    alert('只能上传 doc / docx / pdf / xls / xlsx 文件')
    event.target.value = ''
    return
  }

  uploadForm.file = file
  uploadForm.fileName = file.name
  uploadForm.fileUrl = URL.createObjectURL(file)
}

async function uploadChange() {
  if (!uploadForm.projectName) {
    alert('请选择项目')
    return
  }

  if (!uploadForm.changeName) {
    alert('请输入变更标题')
    return
  }

  if (!uploadForm.file) {
    alert('请上传需求变更文件')
    return
  }

  const projectId = projectMap.value[uploadForm.projectName]

  if (!projectId) {
    alert('没有找到对应项目ID，请重新选择项目')
    return
  }

  const payload = {
    projectId,
    changeTitle: uploadForm.changeName,
    changeName: uploadForm.changeName,
    changeType: uploadForm.changeType,
    fileId: 1,
    status: frontendStatusToBackend('draft'),
    closeStatus: '未关闭',
    submitUserId: 1,
    submitUserName: currentUserName.value,
    remark: uploadForm.remark || ''
  }

  try {
    const res = await createRequirementChange(payload)
    const result = getResponseData(res)

    console.log('新增需求变更返回：', result)

    if (result.code === 200) {
      alert('新增需求变更成功')
      showUploadDialog.value = false
      await loadRequirementChanges()
    } else {
      alert(result.msg || '新增需求变更失败')
    }
  } catch (err) {
    console.error('新增需求变更失败：', err)
    alert('新增需求变更失败，请检查后端接口')
  }
}

function viewChange(item) {
  selectedChange.value = item
}

function downloadChange(item) {
  if (!item.fileUrl) {
    alert('当前还没有接真实文件下载，后面做 project_files 文件上传下载时再接')
    return
  }

  const link = document.createElement('a')
  link.href = item.fileUrl
  link.download = item.fileName || '需求变更文件'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

async function submitChange(item) {
  if (!item || !item.id) {
    alert('提交失败：没有拿到需求变更ID')
    return
  }

  try {
    const res = await submitRequirementChange(item.id)
    const result = getResponseData(res)

    console.log('提交需求变更返回：', result)

    if (result.code === 200) {
      alert(`需求变更【${item.changeName}】已提交审核`)
      await loadRequirementChanges()
    } else {
      alert(result.msg || '提交失败')
    }
  } catch (err) {
    console.error('提交需求变更失败：', err)
    alert('提交需求变更失败')
  }
}

function auditChange(item) {
  selectedChange.value = item
}

async function approveChange(item) {
  if (!item || !item.id) {
    alert('审核失败：没有拿到需求变更ID')
    return
  }

  try {
    const res = await auditRequirementChange(item.id, {
      auditUserId: 1,
      auditUserName: '领导',
      auditStatus: '已通过',
      rejectReason: ''
    })

    const result = getResponseData(res)

    console.log('需求变更审核通过返回：', result)

    if (result.code === 200) {
      alert(`需求变更【${item.changeName}】审核通过`)
      selectedChange.value = null
      await loadRequirementChanges()
    } else {
      alert(result.msg || '审核失败')
    }
  } catch (err) {
    console.error('审核需求变更失败：', err)
    alert('审核失败')
  }
}

async function rejectChange(item) {
  if (!item || !item.id) {
    alert('审核失败：没有拿到需求变更ID')
    return
  }

  const reason = prompt('请输入驳回原因') || '变更说明不完整，请补充后重新提交'

  try {
    const res = await auditRequirementChange(item.id, {
      auditUserId: 1,
      auditUserName: '领导',
      auditStatus: '已驳回',
      rejectReason: reason
    })

    const result = getResponseData(res)

    console.log('需求变更审核驳回返回：', result)

    if (result.code === 200) {
      alert(`需求变更【${item.changeName}】已驳回`)
      selectedChange.value = null
      await loadRequirementChanges()
    } else {
      alert(result.msg || '审核失败')
    }
  } catch (err) {
    console.error('驳回需求变更失败：', err)
    alert('审核失败')
  }
}

async function closeChange(item) {
  if (!item || !item.id) {
    alert('关闭失败：没有拿到需求变更ID')
    return
  }

  const ok = confirm(`确认关闭需求变更【${item.changeName}】吗？`)
  if (!ok) return

  try {
    const res = await closeRequirementChange(item.id, {
      closeUserId: 1,
      closeUserName: currentUserName.value
    })

    const result = getResponseData(res)

    console.log('关闭需求变更返回：', result)

    if (result.code === 200) {
      alert('需求变更已关闭')
      selectedChange.value = null
      await loadRequirementChanges()
    } else {
      alert(result.msg || '关闭失败')
    }
  } catch (err) {
    console.error('关闭需求变更失败：', err)
    alert('关闭失败，请检查后端接口')
  }
}

async function deleteChange(item) {
  if (!item || !item.id) {
    alert('删除失败：没有拿到需求变更ID')
    return
  }

  const ok = confirm(`确认删除需求变更【${item.changeName}】吗？`)
  if (!ok) return

  try {
    const res = await deleteRequirementChange(item.id)
    const result = getResponseData(res)

    console.log('删除需求变更返回：', result)

    if (result.code === 200) {
      alert('删除成功')
      selectedChange.value = null
      await loadRequirementChanges()
    } else {
      alert(result.msg || '删除失败')
    }
  } catch (err) {
    console.error('删除需求变更失败：', err)
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
