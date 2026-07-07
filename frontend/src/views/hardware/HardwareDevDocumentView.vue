<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1>开发文档</h1>
      </div>

      <button v-if="canUseAction('hardware-dev-doc:upload')" class="primary-btn" @click="openUploadDialog">
        上传文档
      </button>
    </div>

    <div class="filter-card">
      <input v-model="filters.keyword" placeholder="搜索项目名称 / 文件名 / 上传人" />
      <select v-model="filters.projectName">
        <option value="">全部项目</option>
        <option v-for="project in projectOptions" :key="project" :value="project">
          {{ project }}
        </option>
      </select>
      <button class="query-btn" @click="loadDocuments">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <div class="table-card">
      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>项目名称</th>
              <th>资料数量</th>
              <th>最近上传时间</th>
              <th>最近上传人</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>
          <tbody>
            <template v-for="project in groupedDocuments" :key="project.projectName">
              <tr class="project-row">
                <td>
                  <button class="project-name-btn" @click="toggleProject(project.projectName)">
                    <span class="expand-icon">{{ expandedProjects.includes(project.projectName) ? '▼' : '▶' }}</span>
                    {{ project.projectName }}
                  </button>
                </td>
                <td><span class="count-tag">{{ project.files.length }} 份文档</span></td>
                <td class="muted">{{ project.latestUploadTime || '-' }}</td>
                <td>{{ project.latestUploader || '-' }}</td>
                <td class="operation-col">
                  <div class="action-group">
                    <button class="text-btn blue" @click="toggleProject(project.projectName)">
                      {{ expandedProjects.includes(project.projectName) ? '收起' : '查看文档' }}
                    </button>
                  </div>
                </td>
              </tr>

              <tr v-if="expandedProjects.includes(project.projectName)" class="child-row">
                <td colspan="5">
                  <div class="child-table-wrapper">
                    <table class="child-table">
                      <thead>
                        <tr>
                          <th>资料名称</th>
                          <th>文件名称</th>
                          <th>上传人</th>
                          <th>上传时间</th>
                          <th>文件大小</th>
                          <th class="child-operation-col">操作</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="item in project.files" :key="item.id">
                          <td>
                            <button class="file-link" @click="viewDocument(item)">
                              {{ item.fileName }}
                            </button>
                          </td>
                          <td>
                            <div class="file-name" :title="item.fileName">{{ item.fileName }}</div>
                          </td>
                          <td>{{ item.uploader }}</td>
                          <td class="muted">{{ item.uploadTime || '-' }}</td>
                          <td class="muted">{{ item.fileSizeText || '-' }}</td>
                          <td class="child-operation-col">
                            <div class="action-group">
                              <button class="text-btn" @click="viewDocument(item)">查看</button>
                              <button class="text-btn blue" @click="downloadDocument(item)">下载</button>
                              <button v-if="canUseAction('hardware-dev-doc:delete')" class="text-btn red" @click="deleteDocument(item)">删除</button>
                            </div>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
      <div class="table-footer">共 {{ filteredDocuments.length }} 条开发文档，按项目归类展示</div>
    </div>

    <div v-if="showUploadDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>上传开发文档</h3>
          <button @click="showUploadDialog = false">×</button>
        </div>
        <div class="form-grid">
          <label>
            绑定项目
            <select v-model="uploadForm.projectName">
              <option value="">请选择项目</option>
              <option v-for="project in projectOptions" :key="project" :value="project">{{ project }}</option>
            </select>
          </label>
          <label>
            上传人
            <input v-model="currentUserName" disabled />
          </label>
          <label class="full-row">
            开发文档
            <input type="file" accept=".txt,.ini,.conf,.json,.xml,.yaml,.yml,.cfg,.xlsx,.xls,.doc,.docx,.zip,.pdf" @change="handleFileChange" />
          </label>
          <label class="full-row">
            文档说明
            <textarea v-model="uploadForm.remark" placeholder="填写开发文档说明"></textarea>
          </label>
        </div>
        <div class="dialog-footer">
          <button class="reset-btn" @click="showUploadDialog = false">取消</button>
          <button class="primary-btn" @click="uploadDocument">保存上传</button>
        </div>
      </div>
    </div>

    <div v-if="selectedDocument" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>开发文档详情</h3>
          <button @click="selectedDocument = null">×</button>
        </div>
        <div class="detail-card">
          <div><span>资料名称</span><strong>{{ selectedDocument.fileName }}</strong></div>
          <div><span>绑定项目</span><strong>{{ selectedDocument.projectName }}</strong></div>
          <div><span>上传人</span><strong>{{ selectedDocument.uploader }}</strong></div>
          <div><span>上传时间</span><strong>{{ selectedDocument.uploadTime || '-' }}</strong></div>
          <div><span>文件大小</span><strong>{{ selectedDocument.fileSizeText || '-' }}</strong></div>
        </div>
        <div class="remark-card">
          <span>文档说明</span>
          <p>{{ selectedDocument.remark || '暂无说明' }}</p>
        </div>
        <div class="dialog-footer">
          <button class="reset-btn" @click="downloadDocument(selectedDocument)">下载文件</button>
          <button class="primary-btn" @click="selectedDocument = null">关闭</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { canUseAction } from '@/utils/permission'
import { buildUploadFilePayload, downloadLocalFile, getFilePreviewUrl } from '@/utils/filePreview'
import { getCurrentUserParams } from '@/utils/currentUser'
import { getProjects } from '@/api/project'
import { createHardwareDevDocument, deleteHardwareDevDocument, getHardwareDevDocuments } from '@/api/hardwareDevDocument'

const currentUserName = ref(localStorage.getItem('realName') || localStorage.getItem('username') || '当前用户')
const filters = reactive({ keyword: '', projectName: '' })
const showUploadDialog = ref(false)
const selectedDocument = ref(null)
const expandedProjects = ref([])
const projectOptions = ref([])
const projectMap = ref({})
const documentList = ref([])
const uploadForm = reactive({
  projectName: '',
  fileId: 0,
  fileName: '',
  fileContentType: '',
  fileData: '',
  remark: ''
})

onMounted(async () => {
  await loadProjects()
  await loadDocuments()
})

function getResponseData(res) {
  return res?.data || res
}

function formatDate(value) {
  if (!value) return ''
  if (typeof value === 'string') return value.slice(0, 10)
  if (typeof value === 'object' && value.Valid !== false) return String(value.Time || value.time || '').slice(0, 10)
  return ''
}

function formatFileSize(value) {
  const size = Number(value || 0)
  if (!size) return '-'
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`
  return `${(size / 1024 / 1024).toFixed(1)} MB`
}

function findProjectName(projectId) {
  const found = Object.entries(projectMap.value).find(([, id]) => Number(id) === Number(projectId))
  return found ? found[0] : `项目ID-${projectId}`
}

async function loadProjects() {
  const res = await getProjects()
  const result = getResponseData(res)
  const list = result.data || []
  projectOptions.value = list.map(item => item.projectName)
  projectMap.value = Object.fromEntries(list.map(item => [item.projectName, item.id]))
}

async function loadDocuments() {
  try {
    const res = await getHardwareDevDocuments()
    const result = getResponseData(res)
    if (result.code !== 200) {
      alert(result.msg || '加载开发文档失败')
      return
    }
    documentList.value = (result.data || []).map(item => ({
      id: item.id,
      projectId: item.projectId,
      projectName: item.projectName || findProjectName(item.projectId),
      fileId: item.fileId || 0,
      fileName: item.fileName || (item.fileId ? `文件ID-${item.fileId}` : '暂无文件'),
      fileSize: item.fileSize || 0,
      fileSizeText: formatFileSize(item.fileSize || 0),
      fileUrl: item.fileUrl || getFilePreviewUrl(item.fileId),
      downloadUrl: item.downloadUrl || '',
      uploader: item.uploadUserName || item.uploader || '未知上传人',
      uploadTime: formatDate(item.uploadTime || item.createdAt),
      remark: item.remark || ''
    }))
  } catch (err) {
    console.error('加载开发文档失败：', err)
    alert(err.response?.data || '加载开发文档失败')
  }
}

const filteredDocuments = computed(() => {
  const keyword = filters.keyword.trim()
  return documentList.value.filter(item => {
    const keywordMatch = !keyword ||
      item.projectName.includes(keyword) ||
      item.fileName.includes(keyword) ||
      item.uploader.includes(keyword) ||
      item.remark.includes(keyword)
    const projectMatch = !filters.projectName || item.projectName === filters.projectName
    return keywordMatch && projectMatch
  })
})

const groupedDocuments = computed(() => {
  const map = new Map()
  filteredDocuments.value.forEach(item => {
    if (!map.has(item.projectName)) map.set(item.projectName, { projectName: item.projectName, files: [] })
    map.get(item.projectName).files.push(item)
  })
  return Array.from(map.values()).map(project => {
    const files = [...project.files].sort((a, b) => new Date(b.uploadTime) - new Date(a.uploadTime))
    return {
      projectName: project.projectName,
      files,
      latestUploadTime: files[0]?.uploadTime || '',
      latestUploader: files[0]?.uploader || ''
    }
  })
})

function resetFilters() {
  filters.keyword = ''
  filters.projectName = ''
}

function toggleProject(projectName) {
  expandedProjects.value = expandedProjects.value.includes(projectName)
    ? expandedProjects.value.filter(item => item !== projectName)
    : [...expandedProjects.value, projectName]
}

function openUploadDialog() {
  Object.assign(uploadForm, { projectName: '', fileId: 0, fileName: '', fileContentType: '', fileData: '', remark: '' })
  showUploadDialog.value = true
}

async function handleFileChange(event) {
  const file = event.target.files?.[0]
  if (!file) return
  Object.assign(uploadForm, await buildUploadFilePayload(file))
}

async function uploadDocument() {
  if (!uploadForm.projectName) {
    alert('请选择绑定项目')
    return
  }
  if (!uploadForm.fileId || !uploadForm.fileName) {
    alert('请上传开发文档文件')
    return
  }
  const projectId = projectMap.value[uploadForm.projectName]
  const user = getCurrentUserParams()
  try {
    const res = await createHardwareDevDocument({
      projectId,
      fileId: uploadForm.fileId,
      fileName: uploadForm.fileName,
      fileContentType: uploadForm.fileContentType,
      fileData: uploadForm.fileData,
      uploadUserId: user.userId || 0,
      uploadUserName: user.realName || user.username || currentUserName.value,
      remark: uploadForm.remark
    })
    const result = getResponseData(res)
    if (result.code !== 200) {
      alert(result.msg || '上传开发文档失败')
      return
    }
    if (!expandedProjects.value.includes(uploadForm.projectName)) expandedProjects.value.push(uploadForm.projectName)
    showUploadDialog.value = false
    await loadDocuments()
  } catch (err) {
    console.error('上传开发文档失败：', err)
    alert(err.response?.data || '上传开发文档失败')
  }
}

function viewDocument(item) {
  selectedDocument.value = item
}

function downloadDocument(item) {
  downloadLocalFile(item, item.fileName || '开发文档')
}

async function deleteDocument(item) {
  if (!confirm(`确认删除开发文档【${item.fileName}】吗？`)) return
  const res = await deleteHardwareDevDocument(item.id)
  const result = getResponseData(res)
  if (result.code !== 200) {
    alert(result.msg || '删除失败')
    return
  }
  selectedDocument.value = null
  await loadDocuments()
}
</script>

<style scoped>
@import '../requirement/customer-supplied-shared.css';
</style>
