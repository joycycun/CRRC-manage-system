<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>客供资料</h1>
      </div>

      <button v-if="canUseAction('customer:upload')" class="primary-btn" @click="openUploadDialog">
        上传资料
      </button>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索项目名称 /  上传人 / 文件名"
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

      <button class="query-btn" @click="loadCustomerSuppliedFiles">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <!-- 项目资料列表 -->
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
            <template
              v-for="project in filteredProjectGroupList"
              :key="project.projectName"
            >
              <!-- 第一层：只展示项目名称 -->
              <tr class="project-row">
                <td>
                  <button
                    class="project-name-btn"
                    @click="toggleProject(project.projectName)"
                  >
                    <span class="expand-icon">
                      {{ expandedProjects.includes(project.projectName) ? '▼' : '▶' }}
                    </span>
                    {{ project.projectName }}
                  </button>
                </td>

                <td>
                  <span class="count-tag">
                    {{ project.files.length }} 份资料
                  </span>
                </td>

                <td class="muted">
                  {{ project.latestUploadTime || '-' }}
                </td>

                <td>
                  {{ project.latestUploader || '-' }}
                </td>

                <td class="operation-col">
                  <div class="action-group">
                    <button
                      class="text-btn blue"
                      @click="toggleProject(project.projectName)"
                    >
                      {{ expandedProjects.includes(project.projectName) ? '收起' : '查看资料' }}
                    </button>
                  </div>
                </td>
              </tr>

              <!-- 第二层：点击项目后展示资料 -->
              <tr
                v-if="expandedProjects.includes(project.projectName)"
                class="child-row"
              >
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
                        <tr
                          v-for="item in project.files"
                          :key="item.id"
                        >
                          <td>
                            <button class="file-link" @click="viewConfig(item)">
                              {{ item.configName }}
                            </button>
                          </td>

                          <td>
                            <div class="file-name">{{ item.fileName }}</div>
                          </td>

                          <td>{{ item.uploader }}</td>

                          <td class="muted">{{ item.uploadTime }}</td>

                          <td class="muted">{{ item.fileSize || '-' }}</td>

                          <td class="child-operation-col">
                            <div class="action-group">
                              <button class="text-btn" @click="viewConfig(item)">
                                查看
                              </button>

                              <button class="text-btn blue" @click="downloadConfig(item)">
                                下载
                              </button>

                              <button v-if="canUseAction('customer:delete')" class="text-btn red" @click="deleteConfig(item)">
                                删除
                              </button>
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

      <div class="table-footer">
        共 {{ filteredConfigList.length }} 条资料记录，按项目归类展示
      </div>
    </div>

    <!-- 上传资料弹窗 -->
    <div v-if="showUploadDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>上传资料</h3>
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
            资料名称
            <input
              v-model="uploadForm.configName"
              placeholder="例如：阿根廷项目SIP账号资料"
            />
          </label>

          <label>
            上传人
            <input
              v-model="currentUserName"
              disabled
            />
          </label>

          <label>
            文件名称
            <input
              v-model="uploadForm.fileName"
              placeholder="选择文件后自动填充"
              disabled
            />
          </label>

          <label class="full-row">
            资料文件
            <input
              type="file"
              accept=".txt,.ini,.conf,.json,.xml,.yaml,.yml,.cfg,.xlsx,.xls,.doc,.docx,.zip,.pdf"
              @change="handleFileChange"
            />
          </label>

          <label class="full-row">
            资料说明
            <textarea
              v-model="uploadForm.description"
              placeholder="例如：包含SIP账号、服务器IP、组播地址、协议参数、终端编号等资料内容"
            ></textarea>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showUploadDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="uploadConfig">
            保存上传
          </button>
        </div>
      </div>
    </div>

    <!-- 查看资料详情弹窗 -->
    <div v-if="selectedConfig" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>资料详情</h3>
          <button @click="selectedConfig = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>资料名称</span>
            <strong>{{ selectedConfig.configName }}</strong>
          </div>

          <div>
            <span>绑定项目</span>
            <strong>{{ selectedConfig.projectName }}</strong>
          </div>

          <div>
            <span>文件名称</span>
            <strong>{{ selectedConfig.fileName }}</strong>
          </div>

          <div>
            <span>上传人</span>
            <strong>{{ selectedConfig.uploader }}</strong>
          </div>

          <div>
            <span>上传时间</span>
            <strong>{{ selectedConfig.uploadTime }}</strong>
          </div>

          <div>
            <span>文件大小</span>
            <strong>{{ selectedConfig.fileSize || '-' }}</strong>
          </div>

          <div>
            <span>文件预览</span>
            <button class="inline-link" @click="openConfigFile(selectedConfig)">
              点开查看文件
            </button>
          </div>
        </div>

        <div class="remark-card">
          <span>资料说明</span>
          <p>{{ selectedConfig.description || '暂无说明' }}</p>
        </div>

        <div class="remark-card">
          <span>资料内容示例</span>
          <pre>{{ selectedConfig.previewContent || '当前文件暂无可预览内容，可点击下载查看原文件。' }}</pre>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="downloadConfig(selectedConfig)">
            下载文件
          </button>

          <button class="primary-btn" @click="openConfigFile(selectedConfig)">
            点开查看
          </button>

          <button class="primary-btn" @click="selectedConfig = null">
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
  getCustomerSuppliedFiles,
  createCustomerSuppliedFile,
  deleteCustomerSuppliedFile
} from '@/api/requirement'

const currentUserName = ref(
  localStorage.getItem('username') ||
  localStorage.getItem('accountName') ||
  localStorage.getItem('realName') ||
  '当前用户'
)

const filters = reactive({
  keyword: '',
  projectName: ''
})

const showUploadDialog = ref(false)
const selectedConfig = ref(null)
const expandedProjects = ref([])

const projectOptions = ref([])
const projectMap = ref({})

const uploadForm = reactive({
  projectName: '',
  configName: '',
  fileName: '',
  fileSize: '',
  file: null,
  fileUrl: '',
  description: '',
  previewContent: ''
})

const configFileList = ref([])

onMounted(async () => {
  await loadProjects()
  await loadCustomerSuppliedFiles()
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

function findProjectName(projectId) {
  const found = Object.entries(projectMap.value).find(([, id]) => Number(id) === Number(projectId))
  return found ? found[0] : `项目ID-${projectId}`
}

async function loadProjects() {
  try {
    const res = await getProjects()
    const result = getResponseData(res)

    console.log('项目列表返回：', result)

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

async function loadCustomerSuppliedFiles() {
  try {
    const res = await getCustomerSuppliedFiles()
    const result = getResponseData(res)

    console.log('客供资料列表返回：', result)

    if (result.code !== 200) {
      alert(result.msg || '加载客供资料失败')
      return
    }

    configFileList.value = (result.data || []).map(item => normalizeCustomerSupply(item))
  } catch (err) {
    console.error('加载客供资料失败：', err)
    alert('加载客供资料失败，请检查后端接口')
  }
}

function normalizeCustomerSupply(item) {
  const projectName = item.projectName || findProjectName(item.projectId)

  return {
    id: item.id,
    projectId: item.projectId,
    projectName,
    configName:
      item.materialName ||
      item.configName ||
      item.docName ||
      item.fileDisplayName ||
      '未命名资料',
    fileId: item.fileId || 0,
    fileName:
      item.fileDisplayName ||
      item.fileName ||
      (item.fileId ? `文件ID-${item.fileId}` : '暂无文件'),
    fileSize: item.fileSize || '-',
    fileUrl: item.fileUrl || '',
    uploaderId: item.uploadUserId || item.uploaderId || 0,
    uploader:
      item.uploadUserName ||
      item.uploader ||
      item.uploaderName ||
      '未知上传人',
    uploadTime: formatDate(item.uploadTime || item.createdAt),
    description:
      item.materialDesc ||
      item.description ||
      item.remark ||
      '',
    previewContent:
      item.previewContent ||
      '当前文件暂无可预览内容，可点击下载查看原文件。',
    remark: item.remark || ''
  }
}

const filteredConfigList = computed(() => {
  return configFileList.value.filter(item => {
    const keywordMatch =
      !filters.keyword ||
      item.projectName.includes(filters.keyword) ||
      item.configName.includes(filters.keyword) ||
      item.uploader.includes(filters.keyword) ||
      item.fileName.includes(filters.keyword) ||
      item.description.includes(filters.keyword)

    const projectMatch =
      !filters.projectName || item.projectName === filters.projectName

    return keywordMatch && projectMatch
  })
})

const filteredProjectGroupList = computed(() => {
  const map = new Map()

  filteredConfigList.value.forEach(item => {
    if (!map.has(item.projectName)) {
      map.set(item.projectName, {
        projectName: item.projectName,
        files: []
      })
    }

    map.get(item.projectName).files.push(item)
  })

  return Array.from(map.values()).map(project => {
    const sortedFiles = [...project.files].sort((a, b) => {
      return new Date(b.uploadTime) - new Date(a.uploadTime)
    })

    return {
      projectName: project.projectName,
      files: sortedFiles,
      latestUploadTime: sortedFiles[0]?.uploadTime || '',
      latestUploader: sortedFiles[0]?.uploader || ''
    }
  })
})

function resetFilters() {
  filters.keyword = ''
  filters.projectName = ''
}

function toggleProject(projectName) {
  if (expandedProjects.value.includes(projectName)) {
    expandedProjects.value = expandedProjects.value.filter(item => item !== projectName)
  } else {
    expandedProjects.value.push(projectName)
  }
}

function openUploadDialog() {
  uploadForm.projectName = ''
  uploadForm.configName = ''
  uploadForm.fileName = ''
  uploadForm.fileSize = ''
  uploadForm.file = null
  uploadForm.fileUrl = ''
  uploadForm.description = ''
  uploadForm.previewContent = ''

  showUploadDialog.value = true
}

function handleFileChange(event) {
  const file = event.target.files[0]
  if (!file) return

  uploadForm.file = file
  uploadForm.fileName = file.name
  uploadForm.fileSize = formatFileSize(file.size)
  uploadForm.fileUrl = URL.createObjectURL(file)

  const isTextFile =
    file.name.endsWith('.txt') ||
    file.name.endsWith('.ini') ||
    file.name.endsWith('.conf') ||
    file.name.endsWith('.json') ||
    file.name.endsWith('.xml') ||
    file.name.endsWith('.yaml') ||
    file.name.endsWith('.yml') ||
    file.name.endsWith('.cfg')

  if (isTextFile) {
    const reader = new FileReader()

    reader.onload = () => {
      uploadForm.previewContent = String(reader.result || '')
    }

    reader.readAsText(file)
  } else {
    uploadForm.previewContent = '当前文件不是纯文本资料，建议下载后查看原文件。'
  }
}

async function uploadConfig() {
  if (!uploadForm.projectName) {
    alert('请选择资料绑定项目')
    return
  }

  if (!uploadForm.configName) {
    alert('请输入资料名称')
    return
  }

  if (!uploadForm.file) {
    alert('请上传资料文件')
    return
  }

  const projectId = projectMap.value[uploadForm.projectName]

  if (!projectId) {
    alert('没有找到对应项目ID，请重新选择项目')
    return
  }

  const payload = {
    projectId,
    fileId: 1,
    materialName: uploadForm.configName,
    fileDisplayName: uploadForm.fileName,
    materialDesc: uploadForm.description || '',
    uploadUserId: 1,
    uploadUserName: currentUserName.value,
    remark: uploadForm.description || ''
  }

  try {
    const res = await createCustomerSuppliedFile(payload)
    const result = getResponseData(res)

    console.log('新增客供资料返回：', result)

    if (result.code === 200) {
      alert('上传客供资料成功')

      if (!expandedProjects.value.includes(uploadForm.projectName)) {
        expandedProjects.value.push(uploadForm.projectName)
      }

      showUploadDialog.value = false
      await loadCustomerSuppliedFiles()
    } else {
      alert(result.msg || '上传客供资料失败')
    }
  } catch (err) {
    console.error('上传客供资料失败：', err)
    alert('上传客供资料失败，请检查后端接口')
  }
}

function viewConfig(item) {
  selectedConfig.value = item
}

function openConfigFile(item) {
  if (!item.fileUrl) {
    selectedConfig.value = item
    alert('当前还没有接真实文件预览，后面做 project_files 文件上传下载时再接')
    return
  }

  window.open(item.fileUrl, '_blank')
}

function downloadConfig(item) {
  if (!item.fileUrl) {
    alert('当前还没有接真实文件下载，后面做 project_files 文件上传下载时再接')
    return
  }

  const link = document.createElement('a')
  link.href = item.fileUrl
  link.download = item.fileName || '资料文件'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

async function deleteConfig(item) {
  console.log('点击删除客供资料：', item)

  if (!item || !item.id) {
    alert('删除失败：没有拿到资料ID')
    return
  }

  const ok = confirm(`确认删除资料【${item.configName}】吗？`)
  if (!ok) return

  try {
    const res = await deleteCustomerSuppliedFile(item.id)
    const result = getResponseData(res)

    console.log('删除客供资料返回：', result)

    if (result.code === 200) {
      alert('删除成功')
      selectedConfig.value = null
      await loadCustomerSuppliedFiles()
    } else {
      alert(result.msg || '删除失败')
    }
  } catch (err) {
    console.error('删除客供资料失败：', err)
    alert('删除失败，请检查后端接口')
  }
}

function formatFileSize(size) {
  if (!size && size !== 0) return '-'

  if (size < 1024) {
    return `${size} B`
  }

  if (size < 1024 * 1024) {
    return `${(size / 1024).toFixed(1)} KB`
  }

  return `${(size / 1024 / 1024).toFixed(1)} MB`
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
  grid-template-columns: 1.4fr 220px 90px 90px;
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

.table-wrapper::-webkit-scrollbar,
.child-table-wrapper::-webkit-scrollbar {
  height: 10px;
}

.table-wrapper::-webkit-scrollbar-track,
.child-table-wrapper::-webkit-scrollbar-track {
  background: #020617;
  border-radius: 999px;
}

.table-wrapper::-webkit-scrollbar-thumb,
.child-table-wrapper::-webkit-scrollbar-thumb {
  background: #334155;
  border-radius: 999px;
  border: 2px solid #020617;
}

.table-wrapper::-webkit-scrollbar-thumb:hover,
.child-table-wrapper::-webkit-scrollbar-thumb:hover {
  background: #475569;
}

.table-wrapper::-webkit-scrollbar-button,
.child-table-wrapper::-webkit-scrollbar-button {
  display: none;
}

.table-wrapper,
.child-table-wrapper {
  scrollbar-width: thin;
  scrollbar-color: #334155 #020617;
}

.table-card table {
  width: 100%;
  min-width: 950px;
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

.project-row {
  background: #0f172a;
}

.project-row:hover {
  background: #1e293b80;
}

.project-name-btn {
  border: none;
  background: transparent;
  color: #60a5fa;
  font-size: 14px;
  font-weight: 800;
  cursor: pointer;
  padding: 0;
  text-align: left;
}

.project-name-btn:hover {
  color: #93c5fd;
  text-decoration: underline;
}

.expand-icon {
  display: inline-block;
  width: 18px;
  color: #94a3b8;
}

.count-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  background: #33415566;
  color: #cbd5e1;
  font-size: 12px;
  font-weight: 700;
}

.child-row td {
  padding: 0;
  background: #020617;
}

.child-table-wrapper {
  width: 100%;
  overflow-x: auto;
  overflow-y: hidden;
  padding: 12px 16px 16px;
  box-sizing: border-box;
}

.child-table {
  width: 100%;
  min-width: 850px;
  border-collapse: collapse;
  table-layout: fixed;
  border: 1px solid #1e293b;
  border-radius: 10px;
  overflow: hidden;
}

.child-table thead {
  background: #0f172a;
}

.child-table th {
  padding: 12px 14px;
  color: #94a3b8;
  font-size: 12px;
  font-weight: 600;
  border-bottom: 1px solid #1e293b;
}

.child-table td {
  padding: 13px 14px;
  border-bottom: 1px solid #1e293b;
  color: #e2e8f0;
  font-size: 13px;
}

.child-table tbody tr:hover {
  background: #1e293b80;
}

.file-link {
  border: none;
  background: transparent;
  color: #60a5fa;
  font-size: 13px;
  font-weight: 800;
  cursor: pointer;
  padding: 0;
  text-align: left;
}

.file-link:hover {
  color: #93c5fd;
  text-decoration: underline;
}

.file-name {
  color: #cbd5e1;
  font-size: 12px;
  word-break: break-all;
}

.muted {
  color: #94a3b8 !important;
}

.operation-col {
  width: 180px;
  text-align: right !important;
}

.child-operation-col {
  width: 220px;
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

.remark-card pre {
  margin: 0;
  color: #cbd5e1;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
  font-family: Consolas, Monaco, monospace;
}

@media (max-width: 960px) {
  .filter-card {
    grid-template-columns: 1fr;
  }

  .table-card table {
    min-width: 950px;
  }

  .child-table {
    min-width: 850px;
  }

  .form-grid,
  .detail-card {
    grid-template-columns: 1fr;
  }
}
</style>
