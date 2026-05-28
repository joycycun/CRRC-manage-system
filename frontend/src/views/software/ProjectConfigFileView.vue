<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>配置文件管理</h1>
      </div>

      <button class="primary-btn" @click="openUploadDialog">
        上传配置文件
      </button>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索项目名称 / 配置文件名称 / 上传人"
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

      <select v-model="filters.configType">
        <option value="">全部配置类型</option>
        <option
          v-for="type in configTypeOptions"
          :key="type"
          :value="type"
        >
          {{ type }}
        </option>
      </select>

      <button class="query-btn">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <!-- 数据表格 -->
    <div class="table-card">
      <table>
        <thead>
          <tr>
            <th>配置文件</th>
            <th>绑定项目</th>
            <th>配置类型</th>
            <th>上传人</th>
            <th>上传时间</th>
            <th>文件大小</th>
            <th class="operation-col">操作</th>
          </tr>
        </thead>

        <tbody>
          <tr v-for="item in filteredConfigList" :key="item.id">
            <td>
              <button class="file-link" @click="viewConfig(item)">
                {{ item.configName }}
              </button>
              <div class="file-name">{{ item.fileName }}</div>
            </td>

            <td>
              <span class="project-tag">
                {{ item.projectName }}
              </span>
            </td>

            <td>
              <span class="config-tag" :class="getConfigTypeClass(item.configType)">
                {{ item.configType }}
              </span>
            </td>

            <td>{{ item.uploader }}</td>

            <td class="muted">{{ item.uploadTime }}</td>

            <td class="muted">{{ item.fileSize || '-' }}</td>

            <td class="operation-col">
              <div class="action-group">
                <button class="text-btn" @click="viewConfig(item)">
                  查看
                </button>

                <button class="text-btn blue" @click="downloadConfig(item)">
                  下载
                </button>

                <button class="text-btn red" @click="deleteConfig(item)">
                  删除
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <div class="table-footer">
        共 {{ filteredConfigList.length }} 条配置文件记录
      </div>
    </div>

    <!-- 上传配置文件弹窗 -->
    <div v-if="showUploadDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>上传配置文件</h3>
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
            配置类型
            <select v-model="uploadForm.configType">
              <option value="">请选择配置类型</option>
              <option
                v-for="type in configTypeOptions"
                :key="type"
                :value="type"
              >
                {{ type }}
              </option>
            </select>
          </label>

          <label>
            配置名称
            <input
              v-model="uploadForm.configName"
              placeholder="例如：阿根廷项目SIP账号配置"
            />
          </label>

          <label>
            上传人
            <input
              v-model="uploadForm.uploader"
              placeholder="请输入上传人"
            />
          </label>

          <label class="full-row">
            配置文件
            <input
              type="file"
              accept=".txt,.ini,.conf,.json,.xml,.yaml,.yml,.cfg,.xlsx,.xls,.doc,.docx,.zip"
              @change="handleFileChange"
            />
          </label>

          <label class="full-row">
            文件名称
            <input
              v-model="uploadForm.fileName"
              placeholder="选择文件后自动填充"
              disabled
            />
          </label>

          <label class="full-row">
            配置说明
            <textarea
              v-model="uploadForm.description"
              placeholder="例如：包含SIP账号、服务器IP、组播地址、协议参数、终端编号等配置内容"
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

    <!-- 查看配置文件详情弹窗 -->
    <div v-if="selectedConfig" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>配置文件详情</h3>
          <button @click="selectedConfig = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>配置名称</span>
            <strong>{{ selectedConfig.configName }}</strong>
          </div>

          <div>
            <span>绑定项目</span>
            <strong>{{ selectedConfig.projectName }}</strong>
          </div>

          <div>
            <span>配置类型</span>
            <strong>{{ selectedConfig.configType }}</strong>
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
          <span>配置说明</span>
          <p>{{ selectedConfig.description || '暂无说明' }}</p>
        </div>

        <div class="remark-card">
          <span>配置内容示例</span>
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
import { computed, reactive, ref } from 'vue'

const filters = reactive({
  keyword: '',
  projectName: '',
  configType: ''
})

const showUploadDialog = ref(false)
const selectedConfig = ref(null)

const projectOptions = [
  '香港屯马项目',
  '波尔图二期项目',
  '阿根廷有轨项目',
  '波哥大有轨项目',
  '迪拜项目'
]

const configTypeOptions = [
  'SIP账号设置',
  'IP地址设置',
  '协议配置',
  '服务器配置',
  '音频参数配置',
  '终端编号配置',
  '其他配置'
]

const uploadForm = reactive({
  projectName: '',
  configType: '',
  configName: '',
  uploader: '',
  fileName: '',
  fileSize: '',
  file: null,
  fileUrl: '',
  description: '',
  previewContent: ''
})

const configFileList = ref([
  {
    id: 1,
    projectName: '阿根廷有轨项目',
    configType: 'SIP账号设置',
    configName: '阿根廷项目DACU SIP账号配置',
    fileName: 'Argentina_DACU_sip_accounts.ini',
    fileSize: '12 KB',
    fileUrl: '',
    uploader: '卢进',
    uploadTime: '2026-05-10',
    description: '用于配置阿根廷项目DACU广播控制盒的SIP账号、认证密码、注册服务器等参数。',
    previewContent:
`[sip_account]
terminal=DACU
server=10.0.11.11
backup_server=10.0.11.81
username=1001
password=******
transport=udp`
  },
  {
    id: 2,
    projectName: '波尔图二期项目',
    configType: 'IP地址设置',
    configName: '波尔图二期终端IP地址规划',
    fileName: 'Porto_IP_Address_Config.xlsx',
    fileSize: '28 KB',
    fileUrl: '',
    uploader: '寸诗睿',
    uploadTime: '2026-05-16',
    description: '用于记录波尔图二期项目各终端IP地址、网关、子网掩码和服务器地址。',
    previewContent:
`终端类型：广播控制盒
设备IP：10.0.24.41
主服务器：10.0.11.11
备用服务器：10.0.11.81
网关：10.0.24.1`
  },
  {
    id: 3,
    projectName: '香港屯马项目',
    configType: '协议配置',
    configName: '香港屯马PA协议配置',
    fileName: 'HKTM_PA_protocol_config.json',
    fileSize: '18 KB',
    fileUrl: '',
    uploader: '王宇',
    uploadTime: '2026-05-18',
    description: '用于配置香港屯马项目PA广播协议、心跳周期、报警上报协议和联动参数。',
    previewContent:
`{
  "heartbeat_interval": 30,
  "pa_protocol": "udp",
  "alarm_report": true,
  "multicast_enabled": false
}`
  },
  {
    id: 4,
    projectName: '迪拜项目',
    configType: '音频参数配置',
    configName: '迪拜项目编码板音频参数配置',
    fileName: 'Dubai_ECU_audio_config.conf',
    fileSize: '9 KB',
    fileUrl: '',
    uploader: '郑宇',
    uploadTime: '2026-05-20',
    description: '用于配置迪拜项目编码板采样率、编码格式、音频通道和码率。',
    previewContent:
`sample_rate=48000
codec=opus
channels=2
bitrate=64000
rtp_payload_type=111`
  }
])

const filteredConfigList = computed(() => {
  return configFileList.value.filter(item => {
    const keywordMatch =
      !filters.keyword ||
      item.projectName.includes(filters.keyword) ||
      item.configName.includes(filters.keyword) ||
      item.configType.includes(filters.keyword) ||
      item.uploader.includes(filters.keyword) ||
      item.fileName.includes(filters.keyword)

    const projectMatch =
      !filters.projectName || item.projectName === filters.projectName

    const configTypeMatch =
      !filters.configType || item.configType === filters.configType

    return keywordMatch && projectMatch && configTypeMatch
  })
})

function resetFilters() {
  filters.keyword = ''
  filters.projectName = ''
  filters.configType = ''
}

function openUploadDialog() {
  uploadForm.projectName = ''
  uploadForm.configType = ''
  uploadForm.configName = ''
  uploadForm.uploader = ''
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
    uploadForm.previewContent = '当前文件不是纯文本配置，建议下载后查看原文件。'
  }
}

function uploadConfig() {
  if (!uploadForm.projectName) {
    alert('请选择配置文件绑定项目')
    return
  }

  if (!uploadForm.configType) {
    alert('请选择配置类型')
    return
  }

  if (!uploadForm.configName) {
    alert('请输入配置名称')
    return
  }

  if (!uploadForm.file) {
    alert('请上传配置文件')
    return
  }

  configFileList.value.unshift({
    id: Date.now(),
    projectName: uploadForm.projectName,
    configType: uploadForm.configType,
    configName: uploadForm.configName,
    fileName: uploadForm.fileName,
    fileSize: uploadForm.fileSize,
    fileUrl: uploadForm.fileUrl,
    uploader: uploadForm.uploader || '当前用户',
    uploadTime: new Date().toISOString().slice(0, 10),
    description: uploadForm.description,
    previewContent: uploadForm.previewContent
  })

  showUploadDialog.value = false
}

function viewConfig(item) {
  selectedConfig.value = item
}

function openConfigFile(item) {
  if (!item.fileUrl) {
    selectedConfig.value = item
    return
  }

  window.open(item.fileUrl, '_blank')
}

function downloadConfig(item) {
  if (!item.fileUrl) {
    alert('当前是模拟数据，暂无可下载的原始配置文件')
    return
  }

  const link = document.createElement('a')
  link.href = item.fileUrl
  link.download = item.fileName || '配置文件'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

function deleteConfig(item) {
  const ok = confirm(`确认删除配置文件【${item.configName}】吗？`)
  if (!ok) return

  configFileList.value = configFileList.value.filter(
    config => config.id !== item.id
  )
}

function formatFileSize(size) {
  if (size < 1024) {
    return `${size} B`
  }

  if (size < 1024 * 1024) {
    return `${(size / 1024).toFixed(1)} KB`
  }

  return `${(size / 1024 / 1024).toFixed(1)} MB`
}

function getConfigTypeClass(type) {
  const map = {
    SIP账号设置: 'sip',
    IP地址设置: 'ip',
    协议配置: 'protocol',
    服务器配置: 'server',
    音频参数配置: 'audio',
    终端编号配置: 'terminal',
    其他配置: 'other'
  }

  return map[type] || 'other'
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
  grid-template-columns: 1.4fr 220px 200px 90px 90px;
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
  margin-top: 4px;
  color: #64748b;
  font-size: 12px;
  word-break: break-all;
}

.project-tag,
.config-tag {
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

.config-tag.sip {
  background: #16a34a33;
  color: #4ade80;
}

.config-tag.ip {
  background: #1d4ed833;
  color: #60a5fa;
}

.config-tag.protocol {
  background: #9333ea33;
  color: #c084fc;
}

.config-tag.server {
  background: #d9770633;
  color: #fbbf24;
}

.config-tag.audio {
  background: #0f766e33;
  color: #5eead4;
}

.config-tag.terminal {
  background: #be123c33;
  color: #fb7185;
}

.config-tag.other {
  background: #47556933;
  color: #94a3b8;
}

.muted {
  color: #94a3b8 !important;
}

.operation-col {
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

  .table-card {
    overflow-x: auto;
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