<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>维修记录管理</h1>
      </div>

      <div class="header-actions">
        <button class="reset-btn" @click="exportRepairRecords">
          导出维修记录
        </button>

        <button class="primary-btn" @click="openCreateDialog">
          新增维修记录
        </button>
      </div>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索项目 / 终端 / SN / MAC / 维修人 / 故障现象"
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
        <option value="">全部终端类型</option>
        <option
          v-for="type in deviceTypeOptions"
          :key="type"
          :value="type"
        >
          {{ type }}
        </option>
      </select>

      <select v-model="filters.repairStatus">
        <option value="">全部维修状态</option>
        <option value="pending">待维修</option>
        <option value="repairing">维修中</option>
        <option value="finished">维修完成</option>
        <option value="failed">维修失败</option>
      </select>

      <button class="query-btn">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <!-- 维修记录表格 -->
    <div class="table-card">
      <div class="table-card-header">
        <div>
          <h3>维修记录列表</h3>
          <span>共 {{ filteredRepairList.length }} 条维修记录</span>
        </div>
      </div>

      <div class="table-wrapper">
        <table class="version-table">
          <thead>
            <tr>
              <th>绑定项目</th>
              <th>终端类型</th>
              <th>SN序列号</th>
              <th>MAC地址</th>
              <th>故障现象</th>
              <th>维修状态</th>
              <th>维修人</th>
              <th>维修时间</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="item in filteredRepairList" :key="item.id">
              <td>
                <span class="project-tag">{{ item.projectName }}</span>
              </td>

              <td>
                <span class="device-tag">{{ item.deviceType }}</span>
              </td>

              <td>
                <span class="sn-tag" :title="item.sn">
                  {{ item.sn }}
                </span>
              </td>

              <td>
                <span class="mac-text" :title="item.macAddress">
                  {{ item.macAddress }}
                </span>
              </td>

              <td>
                <span class="fault-text" :title="item.faultDesc">
                  {{ item.faultDesc }}
                </span>
              </td>

              <td>
                <span class="status-tag" :class="item.repairStatus">
                  {{ getRepairStatusText(item.repairStatus) }}
                </span>
              </td>

              <td>{{ item.repairUser }}</td>

              <td class="muted">{{ item.repairTime }}</td>

              <td class="operation-col">
                <div class="action-group">
                  <button class="text-btn" @click="viewRepair(item)">
                    查看
                  </button>

                  <button class="text-btn blue" @click="openEditDialog(item)">
                    修改
                  </button>

                  <button
                    v-if="item.repairStatus === 'failed'"
                    class="text-btn red"
                    @click="deleteRepair(item)"
                  >
                    删除错误记录
                  </button>

                  <button
                    v-else
                    class="text-btn red"
                    @click="deleteRepair(item)"
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
      </div>
    </div>

    <!-- 新增 / 修改维修记录弹窗 -->
    <div v-if="showEditDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>{{ editMode === 'create' ? '新增维修记录' : '修改维修记录' }}</h3>
          <button @click="showEditDialog = false">×</button>
        </div>

        <div class="form-grid">
          <label>
            绑定项目
            <select v-model="repairForm.projectName">
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
            <select v-model="repairForm.deviceType">
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
            SN序列号
            <input
              v-model="repairForm.sn"
              placeholder="例如：DCCU-202605100001"
            />
          </label>

          <label>
            MAC地址
            <input
              v-model="repairForm.macAddress"
              placeholder="例如：00:11:22:33:44:01"
            />
          </label>

          <label>
            维修状态
            <select v-model="repairForm.repairStatus">
              <option value="pending">待维修</option>
              <option value="repairing">维修中</option>
              <option value="finished">维修完成</option>
              <option value="failed">维修失败</option>
            </select>
          </label>

          <label>
            维修人
            <input
              v-model="repairForm.repairUser"
              placeholder="请输入维修人"
            />
          </label>

          <label>
            维修时间
            <input
              v-model="repairForm.repairTime"
              type="date"
            />
          </label>

          <label>
            维修方式
            <select v-model="repairForm.repairMethod">
              <option value="现场维修">现场维修</option>
              <option value="返厂维修">返厂维修</option>
              <option value="远程指导维修">远程指导维修</option>
              <option value="更换板卡">更换板卡</option>
            </select>
          </label>

          <label class="full-row">
            故障现象
            <textarea
              v-model="repairForm.faultDesc"
              placeholder="例如：设备上电无响应、音频无输出、SIP注册失败、按键无效等"
            ></textarea>
          </label>

          <label class="full-row">
            维修处理过程
            <textarea
              v-model="repairForm.repairProcess"
              placeholder="例如：检查电源、重新烧录、替换板卡、复测功能等"
            ></textarea>
          </label>

          <label class="full-row">
            维修结果 / 备注
            <textarea
              v-model="repairForm.remark"
              placeholder="例如：维修完成，功能验证正常；或维修失败，建议更换整机"
            ></textarea>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showEditDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="saveRepair">
            保存
          </button>
        </div>
      </div>
    </div>

    <!-- 查看详情弹窗 -->
    <div v-if="selectedRepair" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>维修记录详情</h3>
          <button @click="selectedRepair = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>绑定项目</span>
            <strong>{{ selectedRepair.projectName }}</strong>
          </div>

          <div>
            <span>终端类型</span>
            <strong>{{ selectedRepair.deviceType }}</strong>
          </div>

          <div>
            <span>SN序列号</span>
            <strong>{{ selectedRepair.sn }}</strong>
          </div>

          <div>
            <span>MAC地址</span>
            <strong>{{ selectedRepair.macAddress }}</strong>
          </div>

          <div>
            <span>维修状态</span>
            <strong>{{ getRepairStatusText(selectedRepair.repairStatus) }}</strong>
          </div>

          <div>
            <span>维修方式</span>
            <strong>{{ selectedRepair.repairMethod }}</strong>
          </div>

          <div>
            <span>维修人</span>
            <strong>{{ selectedRepair.repairUser }}</strong>
          </div>

          <div>
            <span>维修时间</span>
            <strong>{{ selectedRepair.repairTime }}</strong>
          </div>

          <div>
            <span>最后修改时间</span>
            <strong>{{ selectedRepair.updateTime || '-' }}</strong>
          </div>
        </div>

        <div class="remark-card">
          <span>故障现象</span>
          <p>{{ selectedRepair.faultDesc || '暂无故障现象' }}</p>
        </div>

        <div class="remark-card">
          <span>维修处理过程</span>
          <p>{{ selectedRepair.repairProcess || '暂无维修过程' }}</p>
        </div>

        <div class="remark-card">
          <span>维修结果 / 备注</span>
          <p>{{ selectedRepair.remark || '暂无备注' }}</p>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="openEditDialog(selectedRepair)">
            修改记录
          </button>

          <button class="primary-btn" @click="selectedRepair = null">
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
  repairStatus: ''
})

const showEditDialog = ref(false)
const selectedRepair = ref(null)
const currentEditRepair = ref(null)
const editMode = ref('create')

const projectOptions = [
  '香港屯马项目',
  '波尔图二期项目',
  '阿根廷有轨项目',
  '波哥大有轨项目',
  '迪拜项目'
]

const deviceTypeOptions = [
  '司机室控制盒',
  '解码板',
  '司机提醒单元',
  '噪声检测',
  '编码板',
  '功放板'
]

const repairForm = reactive({
  projectName: '',
  deviceType: '',
  sn: '',
  macAddress: '',
  faultDesc: '',
  repairStatus: 'pending',
  repairUser: '',
  repairTime: '',
  repairMethod: '返厂维修',
  repairProcess: '',
  remark: ''
})

const repairList = ref([
  {
    id: 1,
    projectName: '香港屯马项目',
    deviceType: '司机室控制盒',
    sn: 'DCCU-202605100001',
    macAddress: '00:11:22:33:44:01',
    faultDesc: '设备上电后无音频输出',
    repairStatus: 'finished',
    repairUser: '售后人员',
    repairTime: '2026-05-10',
    repairMethod: '返厂维修',
    repairProcess: '检查音频功放输出，重新烧录软件后复测正常。',
    updateTime: '2026-05-10',
    remark: '维修完成，广播和对讲功能验证正常。'
  },
  {
    id: 2,
    projectName: '波尔图二期项目',
    deviceType: '解码板',
    sn: 'DEC-202605110001',
    macAddress: '00:11:22:55:66:01',
    faultDesc: '客室广播偶发解码失败',
    repairStatus: 'repairing',
    repairUser: '维修人员',
    repairTime: '2026-05-16',
    repairMethod: '现场维修',
    repairProcess: '正在检查网络模块和音频解码流程。',
    updateTime: '2026-05-16',
    remark: '问题仍在排查。'
  },
  {
    id: 3,
    projectName: '阿根廷有轨项目',
    deviceType: '司机提醒单元',
    sn: 'DRU-202605120001',
    macAddress: '00:11:22:77:88:01',
    faultDesc: '提醒音播放异常',
    repairStatus: 'pending',
    repairUser: '售后人员',
    repairTime: '2026-05-18',
    repairMethod: '远程指导维修',
    repairProcess: '等待现场提供日志和设备状态。',
    updateTime: '2026-05-18',
    remark: '待维修确认。'
  },
  {
    id: 4,
    projectName: '迪拜项目',
    deviceType: '编码板',
    sn: 'ENC-202605140001',
    macAddress: '00:11:22:BB:CC:01',
    faultDesc: '远程维修失败，设备无法连接',
    repairStatus: 'failed',
    repairUser: '售后人员',
    repairTime: '2026-05-20',
    repairMethod: '远程指导维修',
    repairProcess: '远程连接失败，设备无响应。',
    updateTime: '2026-05-20',
    remark: '错误维修记录，可删除后重新登记。'
  }
])

const filteredRepairList = computed(() => {
  return repairList.value.filter(item => {
    const keywordMatch =
      !filters.keyword ||
      item.projectName.includes(filters.keyword) ||
      item.deviceType.includes(filters.keyword) ||
      item.sn.includes(filters.keyword) ||
      item.macAddress.includes(filters.keyword) ||
      item.faultDesc.includes(filters.keyword) ||
      item.repairUser.includes(filters.keyword) ||
      item.repairProcess.includes(filters.keyword)

    const projectMatch =
      !filters.projectName || item.projectName === filters.projectName

    const deviceTypeMatch =
      !filters.deviceType || item.deviceType === filters.deviceType

    const statusMatch =
      !filters.repairStatus || item.repairStatus === filters.repairStatus

    return keywordMatch && projectMatch && deviceTypeMatch && statusMatch
  })
})

function getRepairStatusText(status) {
  const map = {
    pending: '待维修',
    repairing: '维修中',
    finished: '维修完成',
    failed: '维修失败'
  }

  return map[status] || status
}

function resetFilters() {
  filters.keyword = ''
  filters.projectName = ''
  filters.deviceType = ''
  filters.repairStatus = ''
}

function openCreateDialog() {
  editMode.value = 'create'
  currentEditRepair.value = null

  repairForm.projectName = ''
  repairForm.deviceType = ''
  repairForm.sn = ''
  repairForm.macAddress = ''
  repairForm.faultDesc = ''
  repairForm.repairStatus = 'pending'
  repairForm.repairUser = ''
  repairForm.repairTime = new Date().toISOString().slice(0, 10)
  repairForm.repairMethod = '返厂维修'
  repairForm.repairProcess = ''
  repairForm.remark = ''

  showEditDialog.value = true
}

function openEditDialog(item) {
  editMode.value = 'edit'
  currentEditRepair.value = item

  repairForm.projectName = item.projectName
  repairForm.deviceType = item.deviceType
  repairForm.sn = item.sn
  repairForm.macAddress = item.macAddress
  repairForm.faultDesc = item.faultDesc
  repairForm.repairStatus = item.repairStatus
  repairForm.repairUser = item.repairUser
  repairForm.repairTime = item.repairTime
  repairForm.repairMethod = item.repairMethod
  repairForm.repairProcess = item.repairProcess
  repairForm.remark = item.remark

  showEditDialog.value = true
}

function saveRepair() {
  if (!repairForm.projectName) {
    alert('请选择绑定项目')
    return
  }

  if (!repairForm.deviceType) {
    alert('请选择终端设备')
    return
  }

  if (!repairForm.sn) {
    alert('请输入SN序列号')
    return
  }

  if (!repairForm.macAddress) {
    alert('请输入MAC地址')
    return
  }

  if (!repairForm.faultDesc) {
    alert('请输入故障现象')
    return
  }

  if (!repairForm.repairUser) {
    alert('请输入维修人')
    return
  }

  if (editMode.value === 'create') {
    repairList.value.unshift({
      id: Date.now(),
      projectName: repairForm.projectName,
      deviceType: repairForm.deviceType,
      sn: repairForm.sn,
      macAddress: repairForm.macAddress,
      faultDesc: repairForm.faultDesc,
      repairStatus: repairForm.repairStatus,
      repairUser: repairForm.repairUser,
      repairTime: repairForm.repairTime,
      repairMethod: repairForm.repairMethod,
      repairProcess: repairForm.repairProcess,
      updateTime: new Date().toISOString().slice(0, 10),
      remark: repairForm.remark
    })
  } else {
    currentEditRepair.value.projectName = repairForm.projectName
    currentEditRepair.value.deviceType = repairForm.deviceType
    currentEditRepair.value.sn = repairForm.sn
    currentEditRepair.value.macAddress = repairForm.macAddress
    currentEditRepair.value.faultDesc = repairForm.faultDesc
    currentEditRepair.value.repairStatus = repairForm.repairStatus
    currentEditRepair.value.repairUser = repairForm.repairUser
    currentEditRepair.value.repairTime = repairForm.repairTime
    currentEditRepair.value.repairMethod = repairForm.repairMethod
    currentEditRepair.value.repairProcess = repairForm.repairProcess
    currentEditRepair.value.updateTime = new Date().toISOString().slice(0, 10)
    currentEditRepair.value.remark = repairForm.remark
  }

  showEditDialog.value = false
}

function viewRepair(item) {
  selectedRepair.value = item
}

function deleteRepair(item) {
  const ok = confirm(`确认删除维修记录【${item.sn}】吗？`)
  if (!ok) return

  repairList.value = repairList.value.filter(record => record.id !== item.id)
}

function exportRepairRecords() {
  const header = [
    '绑定项目',
    '终端类型',
    'SN序列号',
    'MAC地址',
    '故障现象',
    '维修状态',
    '维修方式',
    '维修人',
    '维修时间',
    '最后修改时间',
    '维修过程',
    '维修结果备注'
  ]

  const rows = repairList.value.map(item => [
    item.projectName,
    item.deviceType,
    item.sn,
    item.macAddress,
    item.faultDesc,
    getRepairStatusText(item.repairStatus),
    item.repairMethod,
    item.repairUser,
    item.repairTime,
    item.updateTime || '',
    item.repairProcess || '',
    item.remark || ''
  ])

  const csvContent = [header, ...rows]
    .map(row => row.map(cell => `"${String(cell).replace(/"/g, '""')}"`).join(','))
    .join('\n')

  const blob = new Blob(['\uFEFF' + csvContent], {
    type: 'text/csv;charset=utf-8;'
  })

  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = '维修记录.csv'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
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

.header-actions {
  display: flex;
  gap: 12px;
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
  grid-template-columns: 1.4fr 180px 180px 180px 90px 90px;
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

.table-card-header {
  padding: 14px 16px;
  border-bottom: 1px solid #1e293b;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.table-card-header h3 {
  margin: 0;
  color: #f8fafc;
  font-size: 15px;
}

.table-card-header span {
  color: #64748b;
  font-size: 12px;
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

.version-table {
  width: 100%;
  min-width: 1450px;
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

.project-tag,
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

.device-tag {
  background: #0f766e33;
  color: #5eead4;
}

.sn-tag {
  display: inline-flex;
  max-width: 170px;
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

.mac-text {
  display: inline-block;
  max-width: 130px;
  color: #cbd5e1;
  font-size: 12px;
  word-break: break-all;
  font-family: Consolas, Monaco, monospace;
}

.fault-text {
  display: inline-block;
  max-width: 180px;
  color: #cbd5e1;
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.status-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
}

.status-tag.pending {
  background: #d9770633;
  color: #fbbf24;
}

.status-tag.repairing {
  background: #1d4ed833;
  color: #60a5fa;
}

.status-tag.finished {
  background: #16a34a33;
  color: #4ade80;
}

.status-tag.failed {
  background: #dc262633;
  color: #f87171;
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

  .version-table {
    min-width: 1450px;
  }

  .form-grid,
  .detail-card {
    grid-template-columns: 1fr;
  }
}
</style>