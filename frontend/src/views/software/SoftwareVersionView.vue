<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>软件版本管理</h1>
      </div>

      <button class="primary-btn" @click="openCreateDialog">
        新增软件版本
      </button>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索软件版本 / 项目 / 终端 / 业务说明"
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

      <select v-model="filters.status">
        <option value="">全部状态</option>
        <option value="developing">开发中</option>
        <option value="testing">测试中</option>
        <option value="released">已发布</option>
        <option value="deprecated">已废弃</option>
      </select>

      <button class="query-btn">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <!-- 数据表格 -->
    <div class="table-card">
      <table>
        <thead>
          <tr>
            <th>软件版本</th>
            <th>适配项目</th>
            <th>适配终端</th>
            <th>适配硬件版本</th>
            <th>版本状态</th>
            <th>负责人</th>
            <th>发布日期</th>
            <th class="operation-col">操作</th>
          </tr>
        </thead>

        <tbody>
          <tr v-for="item in filteredSoftwareList" :key="item.id">
            <td>
              <button class="version-link" @click="openDownloadPage(item)">
                {{ item.softwareVersion }}
              </button>
              <div class="business-desc">{{ item.businessDesc }}</div>
            </td>

            <td>
              <span class="project-tag">
                {{ item.projectName }}
              </span>
            </td>

            <td>
              <span class="device-tag">
                {{ item.deviceType }}
              </span>
            </td>

            <td>
              <span class="hardware-tag">
                {{ item.hardwareVersion }}
              </span>
            </td>

            <td>
              <span class="status-tag" :class="item.status">
                {{ getStatusText(item.status) }}
              </span>
            </td>

            <td>{{ item.owner }}</td>

            <td class="muted">
              {{ item.releaseDate }}
            </td>

            <td class="operation-col">
              <div class="action-group">
                <button class="text-btn" @click="viewSoftware(item)">
                  查看
                </button>

                <button class="text-btn blue" @click="openEditDialog(item)">
                  修改
                </button>

                <button class="text-btn green" @click="openDownloadPage(item)">
                  跳转下载
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <div class="table-footer">
        共 {{ filteredSoftwareList.length }} 条软件版本记录
      </div>
    </div>

    <!-- 新增 / 修改软件版本弹窗 -->
    <div v-if="showEditDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>{{ editMode === 'create' ? '新增软件版本' : '修改软件版本' }}</h3>
          <button @click="showEditDialog = false">×</button>
        </div>

        <div class="form-grid">
          <label>
            软件版本号
            <input
              v-model="softwareForm.softwareVersion"
              placeholder="例如：SW-CRRC-ARG-DACU.V1.0.0"
            />
          </label>

          <label>
            适配项目
            <select v-model="softwareForm.projectName" @change="onProjectChange">
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
            适配终端
            <select v-model="softwareForm.deviceType">
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

          <label>
            适配硬件版本
            <select v-model="softwareForm.hardwareVersion">
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
            版本状态
            <select v-model="softwareForm.status">
              <option value="developing">开发中</option>
              <option value="testing">测试中</option>
              <option value="released">已发布</option>
              <option value="deprecated">已废弃</option>
            </select>
          </label>

          <label>
            负责人
            <input
              v-model="softwareForm.owner"
              placeholder="请输入软件负责人"
            />
          </label>

          <label class="full-row">
            下载网页地址
            <input
              v-model="softwareForm.downloadUrl"
              placeholder="例如：http://bc.zycoo.com:8050/files/Argentina/DACU/"
            />
          </label>

          <label class="full-row">
            当前版本实现的业务
            <textarea
              v-model="softwareForm.businessDesc"
              placeholder="例如：实现阿根廷项目 DACU 广播控制盒的人工广播、OCC广播、PAD广播、紧急广播等业务功能"
            ></textarea>
          </label>

          <label class="full-row">
            版本说明
            <textarea
              v-model="softwareForm.description"
              placeholder="例如：修复音频播放异常、优化SIP注册流程、适配指定硬件版本"
            ></textarea>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showEditDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="saveSoftwareVersion">
            保存
          </button>
        </div>
      </div>
    </div>

    <!-- 查看详情弹窗 -->
    <div v-if="selectedSoftware" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>软件版本详情</h3>
          <button @click="selectedSoftware = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>软件版本</span>
            <strong>{{ selectedSoftware.softwareVersion }}</strong>
          </div>

          <div>
            <span>适配项目</span>
            <strong>{{ selectedSoftware.projectName }}</strong>
          </div>

          <div>
            <span>适配终端</span>
            <strong>{{ selectedSoftware.deviceType }}</strong>
          </div>

          <div>
            <span>适配硬件版本</span>
            <strong>{{ selectedSoftware.hardwareVersion }}</strong>
          </div>

          <div>
            <span>版本状态</span>
            <strong>{{ getStatusText(selectedSoftware.status) }}</strong>
          </div>

          <div>
            <span>负责人</span>
            <strong>{{ selectedSoftware.owner }}</strong>
          </div>

          <div>
            <span>发布日期</span>
            <strong>{{ selectedSoftware.releaseDate }}</strong>
          </div>

          <div>
            <span>下载网页</span>
            <button class="inline-link" @click="openDownloadPage(selectedSoftware)">
              打开下载网页
            </button>
          </div>
        </div>

        <div class="remark-card">
          <span>当前版本实现的业务</span>
          <p>{{ selectedSoftware.businessDesc || '暂无说明' }}</p>
        </div>

        <div class="remark-card">
          <span>版本说明</span>
          <p>{{ selectedSoftware.description || '暂无说明' }}</p>
        </div>

        <div class="dialog-footer">
          <button class="primary-btn" @click="selectedSoftware = null">
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
  deviceType: '',
  status: ''
})

const showEditDialog = ref(false)
const selectedSoftware = ref(null)
const currentEditSoftware = ref(null)
const editMode = ref('create')

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
    updateTime: '2026-05-10',
    zipFileName: '',
    zipFileUrl: '',
    description: '广播控制盒硬件版本，适配司机人工广播和乘客报警功能'
  },
  {
    id: 2,
    hardwareVersion: 'HD-CRRC-AGTB-04.T1.1.0',
    deviceType: '客室解码板',
    bindProjects: ['阿根廷有轨项目'],
    status: 'testing',
    owner: '王宇',
    updateTime: '2026-05-16',
    zipFileName: '',
    zipFileUrl: '',
    description: '客室解码板硬件测试版本，用于车厢广播解码'
  },
  {
    id: 3,
    hardwareVersion: 'HD-CRRC-BOGT-03.T1.1.0',
    deviceType: '乘客报警器',
    bindProjects: ['波哥大有轨项目'],
    status: 'released',
    owner: '郑宇',
    updateTime: '2026-05-18',
    zipFileName: '',
    zipFileUrl: '',
    description: '乘客报警器硬件冻结版本'
  },
  {
    id: 4,
    hardwareVersion: 'HD-CRRC-DUBAI-05.S1.1.0',
    deviceType: '编码板',
    bindProjects: ['迪拜项目'],
    status: 'draft',
    owner: '郑宇',
    updateTime: '2026-05-20',
    zipFileName: '',
    zipFileUrl: '',
    description: '编码板硬件草稿版本'
  }
])

const softwareForm = reactive({
  softwareVersion: '',
  projectName: '',
  deviceType: '',
  hardwareVersion: '',
  status: 'developing',
  owner: '',
  downloadUrl: '',
  businessDesc: '',
  description: ''
})

const softwareVersionList = ref([
  {
    id: 1,
    softwareVersion: 'SW-CRRC-ARG-DACU.V1.0.0',
    projectName: '阿根廷有轨项目',
    deviceType: '广播控制盒',
    hardwareVersion: 'HD-CRRC-HKTM.01.V1.1.0',
    status: 'released',
    owner: '卢进',
    releaseDate: '2026-05-10',
    downloadUrl: 'http://bc.zycoo.com:8050/files/Argentina/DACU/',
    businessDesc: '实现阿根廷项目 DACU 广播控制盒的人工广播、OCC广播、PAD广播、预录制紧急广播、特殊广播和站台广播业务。',
    description: '适配阿根廷项目 DACU 终端，完成广播业务主流程。'
  },
  {
    id: 2,
    softwareVersion: 'SW-CRRC-HKTM-PACU.V1.2.0',
    projectName: '香港屯马项目',
    deviceType: '广播控制盒',
    hardwareVersion: 'HD-CRRC-HKTM.01.V1.1.0',
    status: 'testing',
    owner: '卢进',
    releaseDate: '2026-05-16',
    downloadUrl: 'http://bc.zycoo.com:8050/files/HongKong/PACU/',
    businessDesc: '实现香港屯马项目 PACU 人工广播、乘客报警联动、司机室广播优先级处理。',
    description: '当前版本处于测试中，等待内部初始测试问题闭环。'
  },
  {
    id: 3,
    softwareVersion: 'SW-CRRC-BOGT-PECU.V1.0.1',
    projectName: '波哥大有轨项目',
    deviceType: '乘客报警器',
    hardwareVersion: 'HD-CRRC-BOGT-03.T1.1.0',
    status: 'developing',
    owner: '寸诗睿',
    releaseDate: '2026-05-18',
    downloadUrl: 'http://bc.zycoo.com:8050/files/Bogota/PECU/',
    businessDesc: '实现乘客报警器呼叫、报警上报、司机室接听、报警状态恢复等业务。',
    description: '开发阶段版本，尚未进入发布测试。'
  },
  {
    id: 4,
    softwareVersion: 'SW-CRRC-DUBAI-ECU.V0.9.0',
    projectName: '迪拜项目',
    deviceType: '编码板',
    hardwareVersion: 'HD-CRRC-DUBAI-05.S1.1.0',
    status: 'deprecated',
    owner: '寸诗睿',
    releaseDate: '2026-05-20',
    downloadUrl: 'http://bc.zycoo.com:8050/files/Dubai/ECU/',
    businessDesc: '实现迪拜项目编码板音频编码、码流发送、网络传输基础业务。',
    description: '旧测试版本，已废弃，不建议继续使用。'
  }
])

const availableHardwareVersions = computed(() => {
  if (!softwareForm.projectName && !softwareForm.deviceType) {
    return hardwareVersionList.value
  }

  return hardwareVersionList.value.filter(item => {
    const projectMatch =
      !softwareForm.projectName ||
      item.bindProjects.includes(softwareForm.projectName)

    const deviceMatch =
      !softwareForm.deviceType ||
      item.deviceType === softwareForm.deviceType

    return projectMatch && deviceMatch
  })
})

const availableDeviceTypes = computed(() => {
  if (!softwareForm.projectName) {
    return deviceTypeOptions
  }

  const types = hardwareVersionList.value
    .filter(item => item.bindProjects.includes(softwareForm.projectName))
    .map(item => item.deviceType)

  return [...new Set(types)]
})

const filteredSoftwareList = computed(() => {
  return softwareVersionList.value.filter(item => {
    const keywordMatch =
      !filters.keyword ||
      item.softwareVersion.includes(filters.keyword) ||
      item.projectName.includes(filters.keyword) ||
      item.deviceType.includes(filters.keyword) ||
      item.hardwareVersion.includes(filters.keyword) ||
      item.businessDesc.includes(filters.keyword)

    const projectMatch =
      !filters.projectName || item.projectName === filters.projectName

    const deviceTypeMatch =
      !filters.deviceType || item.deviceType === filters.deviceType

    const statusMatch =
      !filters.status || item.status === filters.status

    return keywordMatch && projectMatch && deviceTypeMatch && statusMatch
  })
})

function getStatusText(status) {
  const map = {
    developing: '开发中',
    testing: '测试中',
    released: '已发布',
    deprecated: '已废弃'
  }

  return map[status] || status
}

function resetFilters() {
  filters.keyword = ''
  filters.projectName = ''
  filters.deviceType = ''
  filters.status = ''
}

function onProjectChange() {
  softwareForm.deviceType = ''
  softwareForm.hardwareVersion = ''
}

function openCreateDialog() {
  editMode.value = 'create'
  currentEditSoftware.value = null

  softwareForm.softwareVersion = ''
  softwareForm.projectName = ''
  softwareForm.deviceType = ''
  softwareForm.hardwareVersion = ''
  softwareForm.status = 'developing'
  softwareForm.owner = ''
  softwareForm.downloadUrl = ''
  softwareForm.businessDesc = ''
  softwareForm.description = ''

  showEditDialog.value = true
}

function openEditDialog(item) {
  editMode.value = 'edit'
  currentEditSoftware.value = item

  softwareForm.softwareVersion = item.softwareVersion
  softwareForm.projectName = item.projectName
  softwareForm.deviceType = item.deviceType
  softwareForm.hardwareVersion = item.hardwareVersion
  softwareForm.status = item.status
  softwareForm.owner = item.owner
  softwareForm.downloadUrl = item.downloadUrl
  softwareForm.businessDesc = item.businessDesc
  softwareForm.description = item.description

  showEditDialog.value = true
}

function saveSoftwareVersion() {
  if (!softwareForm.softwareVersion) {
    alert('请输入软件版本号')
    return
  }

  if (!softwareForm.projectName) {
    alert('请选择适配项目')
    return
  }

  if (!softwareForm.deviceType) {
    alert('请选择适配终端')
    return
  }

  if (!softwareForm.hardwareVersion) {
    alert('请选择适配硬件版本')
    return
  }

  if (!softwareForm.businessDesc) {
    alert('请填写当前软件版本实现了哪些业务')
    return
  }

  if (!softwareForm.downloadUrl) {
    alert('请填写软件版本下载网页地址')
    return
  }

  if (editMode.value === 'create') {
    softwareVersionList.value.unshift({
      id: Date.now(),
      softwareVersion: softwareForm.softwareVersion,
      projectName: softwareForm.projectName,
      deviceType: softwareForm.deviceType,
      hardwareVersion: softwareForm.hardwareVersion,
      status: softwareForm.status,
      owner: softwareForm.owner || '未填写',
      releaseDate: new Date().toISOString().slice(0, 10),
      downloadUrl: softwareForm.downloadUrl,
      businessDesc: softwareForm.businessDesc,
      description: softwareForm.description
    })
  } else {
    currentEditSoftware.value.softwareVersion = softwareForm.softwareVersion
    currentEditSoftware.value.projectName = softwareForm.projectName
    currentEditSoftware.value.deviceType = softwareForm.deviceType
    currentEditSoftware.value.hardwareVersion = softwareForm.hardwareVersion
    currentEditSoftware.value.status = softwareForm.status
    currentEditSoftware.value.owner = softwareForm.owner || '未填写'
    currentEditSoftware.value.releaseDate = new Date().toISOString().slice(0, 10)
    currentEditSoftware.value.downloadUrl = softwareForm.downloadUrl
    currentEditSoftware.value.businessDesc = softwareForm.businessDesc
    currentEditSoftware.value.description = softwareForm.description
  }

  showEditDialog.value = false
}

function viewSoftware(item) {
  selectedSoftware.value = item
}

function openDownloadPage(item) {
  if (!item.downloadUrl) {
    alert('当前软件版本没有配置下载网页地址')
    return
  }

  window.open(item.downloadUrl, '_blank')
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
  grid-template-columns: 1.4fr 200px 180px 160px 90px 90px;
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

.version-link {
  border: none;
  background: transparent;
  color: #60a5fa;
  font-size: 13px;
  font-weight: 800;
  cursor: pointer;
  padding: 0;
  text-align: left;
}

.version-link:hover {
  color: #93c5fd;
  text-decoration: underline;
}

.business-desc {
  margin-top: 4px;
  color: #64748b;
  font-size: 12px;
  line-height: 1.5;
  word-break: break-all;
}

.project-tag,
.device-tag,
.hardware-tag {
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

.hardware-tag {
  background: #9333ea33;
  color: #c084fc;
}

.status-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
}

.status-tag.developing {
  background: #47556933;
  color: #94a3b8;
}

.status-tag.testing {
  background: #d9770633;
  color: #fbbf24;
}

.status-tag.released {
  background: #16a34a33;
  color: #4ade80;
}

.status-tag.deprecated {
  background: #dc262633;
  color: #f87171;
}

.muted {
  color: #94a3b8 !important;
}

.operation-col {
  width: 260px;
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
}

.inline-link {
  border: none;
  background: transparent;
  color: #60a5fa;
  cursor: pointer;
  padding: 0;
  font-size: 14px;
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
    min-width: 1200px;
  }

  .form-grid,
  .detail-card {
    grid-template-columns: 1fr;
  }
}
</style>