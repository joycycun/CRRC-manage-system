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

      <button class="query-btn" @click="loadHardwareTests">查询</button>
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
import { computed, onMounted, reactive, ref } from 'vue'

import { getProjects } from '@/api/project'

import {
  getHardwareVersions,
  getHardwareTests,
  createHardwareTest,
  submitHardwareTest,
  auditHardwareTest,
  deleteHardwareTest
} from '@/api/hardware'

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

const projectOptions = ref([])
const projectMap = ref({})

const hardwareVersionOptions = ref([])
const hardwareVersionMap = ref({})

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

const hardwareTestList = ref([])

onMounted(async () => {
  await loadProjects()
  await loadHardwareVersions()
  await loadHardwareTests()
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
function backendAuditStatusToFrontend(status) {
  const map = {
    草稿: 'draft',
    未提交: 'draft',
    待审核: 'submitted',
    已通过: 'approved',
    审核通过: 'approved',
    已驳回: 'rejected',
    审核驳回: 'rejected',
    draft: 'draft',
    submitted: 'submitted',
    approved: 'approved',
    rejected: 'rejected'
  }

  return map[status] || status || 'draft'
}

function findProjectName(projectId) {
  const found = Object.entries(projectMap.value).find(([, id]) => Number(id) === Number(projectId))
  return found ? found[0] : `项目ID-${projectId}`
}

function findHardwareVersion(hardwareId) {
  const found = Object.entries(hardwareVersionMap.value).find(([, id]) => Number(id) === Number(hardwareId))
  return found ? found[0] : `硬件ID-${hardwareId}`
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

async function loadHardwareVersions() {
  try {
    const res = await getHardwareVersions()
    const result = getResponseData(res)

    if (result.code !== 200) {
      alert(result.msg || '加载硬件版本失败')
      return
    }

    const list = result.data || []

    hardwareVersionOptions.value = list.map(item => item.hardwareVersion)

    const map = {}
    list.forEach(item => {
      map[item.hardwareVersion] = item.id
    })

    hardwareVersionMap.value = map
  } catch (err) {
    console.error('加载硬件版本失败：', err)
    alert('加载硬件版本失败')
  }
}

async function loadHardwareTests() {
  try {
    const res = await getHardwareTests()
    const result = getResponseData(res)

    console.log('硬件测试列表返回：', result)

    if (result.code !== 200) {
      alert(result.msg || '加载硬件测试记录失败')
      return
    }

    hardwareTestList.value = (result.data || []).map(item => normalizeHardwareTest(item))
  } catch (err) {
    console.error('加载硬件测试记录失败：', err)
    alert('加载硬件测试记录失败，请检查后端接口')
  }
}

function normalizeHardwareTest(item) {
  return {
    id: item.id,
    projectId: item.projectId || 0,
    hardwareId: item.hardwareId || 0,
    projectName: item.projectName || findProjectName(item.projectId),
    recordName: item.recordName || item.testName || item.hardwareTestName || '',
    hardwareVersion: item.hardwareVersion || findHardwareVersion(item.hardwareId),
    deviceType: item.deviceType || '',
    fileId: item.fileId || 0,
    fileName:
      item.fileName ||
      item.fileDisplayName ||
      (item.fileId ? `文件ID-${item.fileId}.docx` : '暂无文件'),
    fileUrl: item.fileUrl || '',
    uploaderId: item.uploaderId || 0,
    uploader:
  item.uploaderName ||
  item.uploadUserName ||
  item.submitUserName ||
  item.creatorName ||
  item.uploader ||
  currentUserName.value ||
  '当前用户',
    uploadTime: formatDate(item.uploadTime || item.createdAt),
    auditStatus: backendAuditStatusToFrontend(item.auditStatus || item.status),
    auditor: item.auditorName || item.auditUserName || item.auditor || '',
    auditTime: formatDate(item.auditTime),
    rejectReason: item.rejectReason || '',
    remark: item.remark || ''
  }
}

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

async function uploadTest() {
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

  const projectId = projectMap.value[uploadForm.projectName]
  const hardwareId = hardwareVersionMap.value[uploadForm.hardwareVersion]

  if (!projectId) {
    alert('没有找到项目ID，请重新选择项目')
    return
  }

  if (!hardwareId) {
    alert('没有找到硬件版本ID，请重新选择硬件版本')
    return
  }

const payload = {
  projectId,
  hardwareId,
  deviceType: uploadForm.deviceType,
  testName: uploadForm.recordName,
  recordName: uploadForm.recordName,
  hardwareVersion: uploadForm.hardwareVersion,
  fileId: 1,
  fileName: uploadForm.fileName,

  // 上传人
  uploaderId: 1,
  uploaderName: currentUserName.value,
  uploadUserId: 1,
  uploadUserName: currentUserName.value,

  auditStatus: '草稿',
  status: '草稿',
  remark: uploadForm.remark || ''
}

  try {
    const res = await createHardwareTest(payload)
    const result = getResponseData(res)

    console.log('新增硬件测试记录返回：', result)

    if (result.code === 200) {
      alert('上传硬件测试记录成功')
      showUploadDialog.value = false
      await loadHardwareTests()
    } else {
      alert(result.msg || '上传硬件测试记录失败')
    }
  } catch (err) {
    console.error('上传硬件测试记录失败：', err)
    alert('上传硬件测试记录失败，请检查后端接口')
  }
}

function viewTest(item) {
  selectedTest.value = item
}

function downloadTest(item) {
  if (!item.fileUrl) {
    alert('当前还没有接真实文件下载，后面做 project_files 文件上传下载时再接')
    return
  }

  const link = document.createElement('a')
  link.href = item.fileUrl
  link.download = item.fileName || '硬件测试记录.docx'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

async function submitTest(item) {
  if (!item || !item.id) {
    alert('提交失败：没有拿到硬件测试记录ID')
    return
  }

  try {
    const res = await submitHardwareTest(item.id)
    const result = getResponseData(res)

    console.log('提交硬件测试记录返回：', result)

    if (result.code === 200) {
      alert(`硬件测试记录【${item.recordName}】已提交领导审核`)
      await loadHardwareTests()
    } else {
      alert(result.msg || '提交失败')
    }
  } catch (err) {
    console.error('提交硬件测试记录失败：', err)
    alert('提交失败。如果这里是 404，说明后端还没有 /api/hardware-tests/{id}/submit 接口')
  }
}

function auditTest(item) {
  selectedTest.value = item
}

async function approveTest(item) {
  if (!item || !item.id) {
    alert('审核失败：没有拿到硬件测试记录ID')
    return
  }

  try {
    const res = await auditHardwareTest(item.id, {
      auditorId: 1,
      auditorName: '领导',
      auditUserId: 1,
      auditUserName: '领导',
      auditStatus: '已通过',
      rejectReason: ''
    })

    const result = getResponseData(res)

    console.log('硬件测试审核通过返回：', result)

    if (result.code === 200) {
      alert(`硬件测试记录【${item.recordName}】审核通过`)
      selectedTest.value = null
      await loadHardwareTests()
    } else {
      alert(result.msg || '审核失败')
    }
  } catch (err) {
    console.error('硬件测试审核失败：', err)
    alert('审核失败，请检查后端接口')
  }
}

function openRejectDialog(item) {
  currentRejectTest.value = item
  rejectForm.reason = item.rejectReason || ''
  showRejectDialog.value = true
}

async function confirmRejectTest() {
  if (!currentRejectTest.value) return

  if (!rejectForm.reason) {
    alert('请填写驳回原因')
    return
  }

  try {
    const res = await auditHardwareTest(currentRejectTest.value.id, {
      auditorId: 1,
      auditorName: '领导',
      auditUserId: 1,
      auditUserName: '领导',
      auditStatus: '已驳回',
      rejectReason: rejectForm.reason
    })

    const result = getResponseData(res)

    console.log('硬件测试审核驳回返回：', result)

    if (result.code === 200) {
      alert(`硬件测试记录【${currentRejectTest.value.recordName}】已驳回`)
      showRejectDialog.value = false
      selectedTest.value = null
      await loadHardwareTests()
    } else {
      alert(result.msg || '驳回失败')
    }
  } catch (err) {
    console.error('硬件测试驳回失败：', err)
    alert('驳回失败，请检查后端接口')
  }
}

function viewRejectReason(item) {
  selectedRejectReason.value = item.rejectReason || '暂无驳回原因'
}

async function deleteTest(item) {
  if (!item || !item.id) {
    alert('删除失败：没有拿到硬件测试记录ID')
    return
  }

  const ok = confirm(`确认删除硬件测试记录【${item.recordName}】吗？`)
  if (!ok) return

  try {
    const res = await deleteHardwareTest(item.id)
    const result = getResponseData(res)

    console.log('删除硬件测试记录返回：', result)

    if (result.code === 200) {
      alert('删除成功')
      selectedTest.value = null
      await loadHardwareTests()
    } else {
      alert(result.msg || '删除失败')
    }
  } catch (err) {
    console.error('删除硬件测试记录失败：', err)
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