<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>生产烧录记录管理</h1>
      </div>

      <button class="primary-btn" @click="openUploadDialog">
        上传烧录记录
      </button>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索项目名称 / SN序列号 / MAC地址 / 上传人"
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

      <select v-model="filters.deviceType">
        <option value="">全部终端</option>
        <option
          v-for="type in deviceTypeOptions"
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
      <div class="table-wrapper">
        <table class="version-table">
          <thead>
            <tr>
              <th>烧录记录</th>
              <th>绑定项目</th>
              <th>终端类型</th>
              <th>SN序列号</th>
              <th>MAC地址</th>
              <th>软件版本</th>
              <th>硬件版本</th>
              <th>上传人</th>
              <th>上传时间</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="item in filteredBurnRecordList" :key="item.id">
              <td>
                <button class="record-link" @click="viewBurnRecord(item)">
                  {{ item.recordName }}
                </button>
                <div class="file-name">{{ item.fileName }}</div>
              </td>

              <td>
                <span class="project-tag">{{ item.projectName }}</span>
              </td>

              <td>
                <span class="device-tag">{{ item.deviceType }}</span>
              </td>

              <td>
                <div class="sn-list">
                  <span
                    v-for="sn in item.snList"
                    :key="sn"
                    class="sn-tag"
                    :title="sn"
                  >
                    {{ sn }}
                  </span>
                </div>
              </td>

              <td>
                <span class="mac-text" :title="item.macAddress">
                  {{ item.macAddress }}
                </span>
              </td>

              <td class="version-cell">
                <span class="software-tag" :title="item.softwareVersion">
                  {{ item.softwareVersion }}
                </span>
              </td>

              <td class="version-cell">
                <span class="hardware-tag" :title="item.hardwareVersion">
                  {{ item.hardwareVersion }}
                </span>
              </td>

              <td>{{ item.uploader }}</td>

              <td class="muted">{{ item.uploadTime }}</td>

              <td class="operation-col">
                <div class="action-group">
                  <button class="text-btn" @click="viewBurnRecord(item)">
                    查看
                  </button>

                  <button class="text-btn blue" @click="downloadBurnRecord(item)">
                    下载
                  </button>

                  <button class="text-btn red" @click="deleteBurnRecord(item)">
                    删除
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

  <div class="table-footer">
    共 {{ filteredBurnRecordList.length }} 条烧录记录
  </div>
</div>

    <!-- 上传烧录记录弹窗 -->
    <div v-if="showUploadDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>上传烧录记录</h3>
          <button @click="showUploadDialog = false">×</button>
        </div>

        <div class="form-grid">
          <label>
            绑定项目
            <select v-model="uploadForm.projectName" @change="onProjectChange">
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
            终端类型
            <select v-model="uploadForm.deviceType">
              <option value="">请选择终端</option>
              <option
                v-for="type in availableDeviceTypes"
                :key="type"
                :value="type"
              >
                {{ type }}
              </option>
            </select>
          </label>

          <label class="full-row">
            SN序列号
            <textarea
              v-model="uploadForm.snText"
              placeholder="请输入该终端对应的SN序列号，一行一个，例如：
          HKTM-DACU-202605100001
          HKTM-DACU-202605100002
          HKTM-DACU-202605100003"
            ></textarea>
          </label>

          <label>
            终端MAC地址
            <input
              v-model="uploadForm.macAddress"
              placeholder="例如：00:11:22:33:44:55"
            />
          </label>

          <label>
            实际烧录软件版本
            <select v-model="uploadForm.softwareVersion">
              <option value="">请选择软件版本</option>
              <option
                v-for="item in availableSoftwareVersions"
                :key="item.id"
                :value="item.softwareVersion"
              >
                {{ item.softwareVersion }} / {{ item.deviceType }}
              </option>
            </select>
          </label>

          <label>
            实际硬件版本
            <select v-model="uploadForm.hardwareVersion">
              <option value="">请选择硬件版本</option>
              <option
                v-for="item in availableHardwareVersions"
                :key="item.id"
                :value="item.hardwareVersion"
              >
                {{ item.hardwareVersion }} / {{ item.deviceType }}
              </option>
            </select>
          </label>

          <label>
            烧录记录名称
            <input
              v-model="uploadForm.recordName"
              placeholder="例如：香港屯马广播控制盒烧录记录"
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
            烧录记录文件
            <input
              type="file"
              accept=".doc,.docx,.xls,.xlsx,.pdf,.txt,.csv,.zip"
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
            烧录说明
            <textarea
              v-model="uploadForm.remark"
              placeholder="例如：记录本台设备烧录时间、烧录工具、烧录人员、烧录结果、异常说明等"
            ></textarea>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showUploadDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="uploadBurnRecord">
            保存上传
          </button>
        </div>
      </div>
    </div>

    <!-- 查看烧录记录详情弹窗 -->
    <div v-if="selectedBurnRecord" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>烧录记录详情</h3>
          <button @click="selectedBurnRecord = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>烧录记录名称</span>
            <strong>{{ selectedBurnRecord.recordName }}</strong>
          </div>

          <div>
            <span>绑定项目</span>
            <strong>{{ selectedBurnRecord.projectName }}</strong>
          </div>

          <div>
            <span>终端类型</span>
            <strong>{{ selectedBurnRecord.deviceType }}</strong>
          </div>


          <div>
            <span>SN序列号</span>
            <strong>{{ selectedBurnRecord.snList.join('、') }}</strong>
          </div>

          <div>
            <span>实际烧录软件版本</span>
            <strong>{{ selectedBurnRecord.softwareVersion }}</strong>
          </div>

          <div>
            <span>实际硬件版本</span>
            <strong>{{ selectedBurnRecord.hardwareVersion }}</strong>
          </div>

          <div>
            <span>文件名称</span>
            <strong>{{ selectedBurnRecord.fileName }}</strong>
          </div>

          <div>
            <span>上传人</span>
            <strong>{{ selectedBurnRecord.uploader }}</strong>
          </div>

          <div>
            <span>上传时间</span>
            <strong>{{ selectedBurnRecord.uploadTime }}</strong>
          </div>

          <div>
            <span>文件查看</span>
            <button class="inline-link" @click="openBurnRecordFile(selectedBurnRecord)">
              点开查看文件
            </button>
          </div>
        </div>

        <div class="remark-card">
          <span>烧录说明</span>
          <p>{{ selectedBurnRecord.remark || '暂无说明' }}</p>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="downloadBurnRecord(selectedBurnRecord)">
            下载文件
          </button>

          <button class="primary-btn" @click="selectedBurnRecord = null">
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
  deviceType: ''
})

const showUploadDialog = ref(false)
const selectedBurnRecord = ref(null)

const projectOptions = [
  '香港屯马项目',
  '波尔图二期项目',
  '阿根廷有轨项目',
  '波哥大有轨项目',
  '迪拜项目'
]

const deviceTypeOptions = [
  '广播控制盒',
  '客室解码板',
  '编码板',
  '乘客报警器',
  '司机室话筒',
  '功放模块'
]

const hardwareVersionList = ref([
  {
    id: 1,
    hardwareVersion: 'HD-CRRC-HKTM.01.V1.1.0',
    deviceType: '广播控制盒',
    bindProjects: ['香港屯马项目', '波尔图二期项目'],
    status: 'released',
    owner: '王宇',
    updateTime: '2026-05-10'
  },
  {
    id: 2,
    hardwareVersion: 'HD-CRRC-AGTB-04.T1.1.0',
    deviceType: '客室解码板',
    bindProjects: ['阿根廷有轨项目'],
    status: 'testing',
    owner: '王宇',
    updateTime: '2026-05-16'
  },
  {
    id: 3,
    hardwareVersion: 'HD-CRRC-BOGT-03.T1.1.0',
    deviceType: '乘客报警器',
    bindProjects: ['波哥大有轨项目'],
    status: 'released',
    owner: '郑宇',
    updateTime: '2026-05-18'
  },
  {
    id: 4,
    hardwareVersion: 'HD-CRRC-DUBAI-05.S1.1.0',
    deviceType: '编码板',
    bindProjects: ['迪拜项目'],
    status: 'draft',
    owner: '郑宇',
    updateTime: '2026-05-20'
  }
])

const softwareVersionList = ref([
  {
    id: 1,
    softwareVersion: 'SW-CRRC-ARG-DACU.V1.0.0',
    projectName: '阿根廷有轨项目',
    deviceType: '广播控制盒',
    hardwareVersion: 'HD-CRRC-HKTM.01.V1.1.0',
    status: 'released',
    owner: '卢进'
  },
  {
    id: 2,
    softwareVersion: 'SW-CRRC-HKTM-PACU.V1.2.0',
    projectName: '香港屯马项目',
    deviceType: '广播控制盒',
    hardwareVersion: 'HD-CRRC-HKTM.01.V1.1.0',
    status: 'testing',
    owner: '卢进'
  },
  {
    id: 3,
    softwareVersion: 'SW-CRRC-BOGT-PECU.V1.0.1',
    projectName: '波哥大有轨项目',
    deviceType: '乘客报警器',
    hardwareVersion: 'HD-CRRC-BOGT-03.T1.1.0',
    status: 'developing',
    owner: '寸诗睿'
  },
  {
    id: 4,
    softwareVersion: 'SW-CRRC-DUBAI-ECU.V0.9.0',
    projectName: '迪拜项目',
    deviceType: '编码板',
    hardwareVersion: 'HD-CRRC-DUBAI-05.S1.1.0',
    status: 'deprecated',
    owner: '寸诗睿'
  }
])

const uploadForm = reactive({
  projectName: '',
  deviceType: '',
  snText: '',
  macAddress: '',
  softwareVersion: '',
  hardwareVersion: '',
  recordName: '',
  uploader: '',
  fileName: '',
  file: null,
  fileUrl: '',
  remark: ''
})

const burnRecordList = ref([
  {
    id: 1,
    projectName: '香港屯马项目',
    deviceType: '广播控制盒',
    snList: ['HKTM-DACU-202605100001'],
    macAddress: '00:11:22:33:44:55',
    softwareVersion: 'SW-CRRC-HKTM-PACU.V1.2.0',
    hardwareVersion: 'HD-CRRC-HKTM.01.V1.1.0',
    recordName: '香港屯马广播控制盒烧录记录',
    fileName: '香港屯马广播控制盒烧录记录_001.xlsx',
    fileUrl: '',
    uploader: '生产人员',
    uploadTime: '2026-05-10',
    remark: '该设备已完成软件烧录、硬件版本确认、MAC地址登记和SN绑定。'
  },
  {
    id: 2,
    projectName: '阿根廷有轨项目',
    deviceType: '广播控制盒',
    snList: ['HKTM-DACU-202605100001',
    'HKTM-DACU-202605100002'
    ],
    macAddress: '00:11:22:AA:BB:08',
    softwareVersion: 'SW-CRRC-ARG-DACU.V1.0.0',
    hardwareVersion: 'HD-CRRC-HKTM.01.V1.1.0',
    recordName: '阿根廷DACU烧录记录',
    fileName: '阿根廷DACU烧录记录_008.xlsx',
    fileUrl: '',
    uploader: '生产人员',
    uploadTime: '2026-05-16',
    remark: '阿根廷DACU设备烧录完成，等待出厂测试。'
  },
  {
    id: 3,
    projectName: '波哥大有轨项目',
    deviceType: '乘客报警器',
    snList: ['BOGT-PECU-202605180012'],
    macAddress: '00:11:22:CC:DD:12',
    softwareVersion: 'SW-CRRC-BOGT-PECU.V1.0.1',
    hardwareVersion: 'HD-CRRC-BOGT-03.T1.1.0',
    recordName: '波哥大乘客报警器烧录记录',
    fileName: '波哥大乘客报警器烧录记录_012.xlsx',
    fileUrl: '',
    uploader: '生产人员',
    uploadTime: '2026-05-18',
    remark: '乘客报警器烧录完成，记录SN与MAC地址。'
  },
  {
    id: 4,
    projectName: '迪拜项目',
    deviceType: '编码板',
    snList: ['DUBAI-ECU-202605200003'],
    macAddress: '00:11:22:EE:FF:03',
    softwareVersion: 'SW-CRRC-DUBAI-ECU.V0.9.0',
    hardwareVersion: 'HD-CRRC-DUBAI-05.S1.1.0',
    recordName: '迪拜编码板烧录记录',
    fileName: '迪拜编码板烧录记录_003.xlsx',
    fileUrl: '',
    uploader: '生产人员',
    uploadTime: '2026-05-20',
    remark: '编码板烧录完成，软件版本为旧测试版本。'
  }
])

const availableDeviceTypes = computed(() => {
  if (!uploadForm.projectName) {
    return deviceTypeOptions
  }

  const softwareTypes = softwareVersionList.value
    .filter(item => item.projectName === uploadForm.projectName)
    .map(item => item.deviceType)

  const hardwareTypes = hardwareVersionList.value
    .filter(item => item.bindProjects.includes(uploadForm.projectName))
    .map(item => item.deviceType)

  return [...new Set([...softwareTypes, ...hardwareTypes])]
})

const availableSoftwareVersions = computed(() => {
  return softwareVersionList.value.filter(item => {
    const projectMatch =
      !uploadForm.projectName || item.projectName === uploadForm.projectName

    const deviceMatch =
      !uploadForm.deviceType || item.deviceType === uploadForm.deviceType

    return projectMatch && deviceMatch
  })
})

const availableHardwareVersions = computed(() => {
  return hardwareVersionList.value.filter(item => {
    const projectMatch =
      !uploadForm.projectName || item.bindProjects.includes(uploadForm.projectName)

    const deviceMatch =
      !uploadForm.deviceType || item.deviceType === uploadForm.deviceType

    return projectMatch && deviceMatch
  })
})

const filteredBurnRecordList = computed(() => {
  return burnRecordList.value.filter(item => {
    const keywordMatch =
      !filters.keyword ||
      item.projectName.includes(filters.keyword) ||
      item.deviceType.includes(filters.keyword) ||
      item.snList.some(sn => sn.includes(filters.keyword)) ||
      item.macAddress.includes(filters.keyword) ||
      item.softwareVersion.includes(filters.keyword) ||
      item.hardwareVersion.includes(filters.keyword) ||
      item.uploader.includes(filters.keyword) ||
      item.recordName.includes(filters.keyword)

    const projectMatch =
      !filters.projectName || item.projectName === filters.projectName

    const deviceTypeMatch =
      !filters.deviceType || item.deviceType === filters.deviceType

    return keywordMatch && projectMatch && deviceTypeMatch
  })
})

function resetFilters() {
  filters.keyword = ''
  filters.projectName = ''
  filters.deviceType = ''
}

function openUploadDialog() {
  uploadForm.projectName = ''
  uploadForm.deviceType = ''
  uploadForm.snText = ''
  uploadForm.macAddress = ''
  uploadForm.softwareVersion = ''
  uploadForm.hardwareVersion = ''
  uploadForm.recordName = ''
  uploadForm.uploader = ''
  uploadForm.fileName = ''
  uploadForm.file = null
  uploadForm.fileUrl = ''
  uploadForm.remark = ''

  showUploadDialog.value = true
}

function onProjectChange() {
  uploadForm.deviceType = ''
  uploadForm.softwareVersion = ''
  uploadForm.hardwareVersion = ''
}

function handleFileChange(event) {
  const file = event.target.files[0]
  if (!file) return

  uploadForm.file = file
  uploadForm.fileName = file.name
  uploadForm.fileUrl = URL.createObjectURL(file)
}

function uploadBurnRecord() {
  if (!uploadForm.projectName) {
    alert('请选择烧录记录绑定项目')
    return
  }

  if (!uploadForm.deviceType) {
    alert('请选择终端类型')
    return
  }

const snList = uploadForm.snText
  .split('\n')
  .map(sn => sn.trim())
  .filter(Boolean)

if (snList.length === 0) {
  alert('请至少输入一个SN序列号')
  return
}

  if (!uploadForm.macAddress) {
    alert('请输入终端MAC地址')
    return
  }

  if (!uploadForm.softwareVersion) {
    alert('请选择实际烧录的软件版本')
    return
  }

  if (!uploadForm.hardwareVersion) {
    alert('请选择实际硬件版本')
    return
  }

  if (!uploadForm.recordName) {
    alert('请输入烧录记录名称')
    return
  }

  if (!uploadForm.file) {
    alert('请上传烧录记录文件')
    return
  }

  burnRecordList.value.unshift({
    id: Date.now(),
    projectName: uploadForm.projectName,
    deviceType: uploadForm.deviceType,
    snList,
    macAddress: uploadForm.macAddress,
    softwareVersion: uploadForm.softwareVersion,
    hardwareVersion: uploadForm.hardwareVersion,
    recordName: uploadForm.recordName,
    fileName: uploadForm.fileName,
    fileUrl: uploadForm.fileUrl,
    uploader: uploadForm.uploader || '当前用户',
    uploadTime: new Date().toISOString().slice(0, 10),
    remark: uploadForm.remark
  })

  showUploadDialog.value = false
}

function viewBurnRecord(item) {
  selectedBurnRecord.value = item
}

function openBurnRecordFile(item) {
  if (!item.fileUrl) {
    alert('当前是模拟数据，暂无可直接打开的原始文件')
    return
  }

  window.open(item.fileUrl, '_blank')
}

function downloadBurnRecord(item) {
  if (!item.fileUrl) {
    alert('当前是模拟数据，暂无可下载的原始文件')
    return
  }

  const link = document.createElement('a')
  link.href = item.fileUrl
  link.download = item.fileName || '烧录记录文件'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

function deleteBurnRecord(item) {
  const ok = confirm(`确认删除烧录记录【${item.recordName}】吗？`)
  if (!ok) return

  burnRecordList.value = burnRecordList.value.filter(
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
  overflow-x: auto;
  overflow-y: hidden;
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

/* Chrome / Edge / Safari 滚动条整体 */
.table-wrapper::-webkit-scrollbar {
  height: 10px;
}

/* 滚动条轨道 */
.table-wrapper::-webkit-scrollbar-track {
  background: #020617;
  border-radius: 999px;
}

/* 可以拖动的滑块 */
.table-wrapper::-webkit-scrollbar-thumb {
  background: #334155;
  border-radius: 999px;
  border: 2px solid #020617;
}

/* 鼠标放上去时 */
.table-wrapper::-webkit-scrollbar-thumb:hover {
  background: #475569;
}

/* 左右两边的小箭头区域 */
.table-wrapper::-webkit-scrollbar-button {
  display: none;
}

/* Firefox */
.table-wrapper {
  scrollbar-width: thin;
  scrollbar-color: #334155 #020617;
}

.version-table {
  width: 100%;
  min-width: 1700px;
  border-collapse: collapse;
  table-layout: fixed;
}

.version-table thead {
  background: #020617;
}

.version-table th {
  padding: 14px 16px;
  color: #94a3b8;
  font-size: 12px;
  font-weight: 600;
  text-align: left;
  border-bottom: 1px solid #1e293b;
  white-space: nowrap;
}

.version-table td {
  padding: 15px 16px;
  font-size: 13px;
  color: #e2e8f0;
  border-bottom: 1px solid #1e293b;
  vertical-align: middle;
}

.record-link {
  border: none;
  background: transparent;
  color: #60a5fa;
  font-size: 13px;
  font-weight: 800;
  cursor: pointer;
  padding: 0;
  text-align: left;
}

.record-link:hover {
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
.device-tag{
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

.device-tag {
  background: #0f766e33;
  color: #5eead4;
}

.version-cell {
  overflow: hidden;
}

.software-tag,
.hardware-tag {
  display: inline-block;
  max-width: 230px;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  vertical-align: middle;
}

.software-tag {
  background: #16a34a33;
  color: #4ade80;
}

.hardware-tag {
  background: #9333ea33;
  color: #c084fc;
}
.sn-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.sn-tag {
  display: inline-flex;
  padding: 3px 8px;
  border-radius: 999px;
  background: #33415566;
  color: #cbd5e1;
  font-size: 12px;
  font-family: Consolas, Monaco, monospace;
}

.mac-text {
  color: #cbd5e1;
  font-size: 12px;
  word-break: break-all;
  font-family: Consolas, Monaco, monospace;
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
  width: 760px;
  max-width: 100%;
  background: #0f172a;
  border: 1px solid #334155;
  border-radius: 16px;
  color: #f8fafc;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.45);
}

.large-dialog {
  width: 920px;
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

  .table-card {
    overflow-x: auto;
  }

  .table-card table {
    min-width: 1500px;
  }

  .form-grid,
  .detail-card {
    grid-template-columns: 1fr;
  }

  .sn-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  max-width: 160px;
}

.sn-tag {
  display: inline-flex;
  max-width: 150px;
  padding: 3px 8px;
  border-radius: 999px;
  background: #33415566;
  color: #cbd5e1;
  font-size: 12px;
  font-family: Consolas, Monaco, monospace;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.table-card th:nth-child(1),
.table-card td:nth-child(1) {
  width: 160px;
}

.table-card th:nth-child(2),
.table-card td:nth-child(2) {
  width: 140px;
}

.table-card th:nth-child(3),
.table-card td:nth-child(3) {
  width: 130px;
}

.table-card th:nth-child(4),
.table-card td:nth-child(4) {
  width: 180px;
}

.table-card th:nth-child(5),
.table-card td:nth-child(5) {
  width: 150px;
}

.table-card th:nth-child(6),
.table-card td:nth-child(6) {
  width: 260px;
}

.table-card th:nth-child(7),
.table-card td:nth-child(7) {
  width: 260px;
}

.table-card th:nth-child(8),
.table-card td:nth-child(8) {
  width: 100px;
}

.table-card th:nth-child(9),
.table-card td:nth-child(9) {
  width: 120px;
}

.table-card th:nth-child(10),
.table-card td:nth-child(10) {
  width: 160px;
}

  
}
</style>