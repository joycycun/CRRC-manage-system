<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>升级记录管理</h1>
      </div>

      <div class="header-actions">
        <button class="reset-btn" @click="exportUpgradeRecords">
          导出升级记录
        </button>

        <button class="primary-btn" @click="openCreateDialog">
          新增升级记录
        </button>
      </div>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索项目 / 终端 / SN / MAC / 版本 / 升级人"
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

      <select v-model="filters.upgradeStatus">
        <option value="">全部升级状态</option>
        <option value="success">升级成功</option>
        <option value="pending">待确认</option>
        <option value="failed">升级失败</option>
      </select>

      <button class="query-btn">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <!-- 升级记录表格 -->
    <div class="table-card">
      <div class="table-card-header">
        <div>
          <h3>软件升级记录列表</h3>
          <span>共 {{ filteredUpgradeList.length }} 条升级记录</span>
        </div>
      </div>

      <div class="table-wrapper">
        <table class="version-table">
          <thead>
            <tr>
              <th>项目</th>
              <th>终端类型</th>
              <th>SN序列号</th>
              <th>MAC地址</th>
              <th>升级前版本</th>
              <th>升级后版本</th>
              <th>升级状态</th>
              <th>升级人</th>
              <th>升级时间</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="item in filteredUpgradeList" :key="item.id">
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

              <td class="version-cell">
                <span class="old-version-tag" :title="item.oldSoftwareVersion">
                  {{ item.oldSoftwareVersion }}
                </span>
              </td>

              <td class="version-cell">
                <span class="software-tag" :title="item.newSoftwareVersion">
                  {{ item.newSoftwareVersion }}
                </span>
              </td>

              <td>
                <span class="status-tag" :class="item.upgradeStatus">
                  {{ getUpgradeStatusText(item.upgradeStatus) }}
                </span>
              </td>

              <td>{{ item.upgradeUser }}</td>

              <td class="muted">{{ item.upgradeTime }}</td>

              <td class="operation-col">
                <div class="action-group">
                  <button class="text-btn" @click="viewUpgrade(item)">
                    查看
                  </button>

                  <button class="text-btn blue" @click="openEditDialog(item)">
                    修改
                  </button>

                  <button
                    v-if="item.upgradeStatus === 'failed'"
                    class="text-btn red"
                    @click="deleteUpgrade(item)"
                  >
                    删除错误记录
                  </button>

                  <button
                    v-else
                    class="text-btn red"
                    @click="deleteUpgrade(item)"
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

    <!-- 新增 / 修改升级记录弹窗 -->
    <div v-if="showEditDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>{{ editMode === 'create' ? '新增升级记录' : '修改升级记录' }}</h3>
          <button @click="showEditDialog = false">×</button>
        </div>

        <div class="form-grid">
          <label>
            绑定项目
            <select v-model="upgradeForm.projectName">
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
            <select v-model="upgradeForm.deviceType">
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
              v-model="upgradeForm.sn"
              placeholder="例如：DCCU-202605100001"
            />
          </label>

          <label>
            MAC地址
            <input
              v-model="upgradeForm.macAddress"
              placeholder="例如：00:11:22:33:44:01"
            />
          </label>

          <label>
            升级前软件版本
            <input
              v-model="upgradeForm.oldSoftwareVersion"
              placeholder="例如：SW-DCCU.V1.1.0"
            />
          </label>

          <label>
            升级后软件版本
            <input
              v-model="upgradeForm.newSoftwareVersion"
              placeholder="例如：SW-DCCU.V1.2.0"
            />
          </label>

          <label>
            升级状态
            <select v-model="upgradeForm.upgradeStatus">
              <option value="success">升级成功</option>
              <option value="pending">待确认</option>
              <option value="failed">升级失败</option>
            </select>
          </label>

          <label>
            升级人
            <input
              v-model="upgradeForm.upgradeUser"
              placeholder="请输入升级人"
            />
          </label>

          <label>
            升级时间
            <input
              v-model="upgradeForm.upgradeTime"
              type="date"
            />
          </label>

          <label>
            升级方式
            <select v-model="upgradeForm.upgradeMethod">
              <option value="远程升级">远程升级</option>
              <option value="现场升级">现场升级</option>
              <option value="售后返修升级">售后返修升级</option>
              <option value="生产重新烧录">生产重新烧录</option>
            </select>
          </label>

          <label class="full-row">
            升级说明
            <textarea
              v-model="upgradeForm.remark"
              placeholder="例如：升级原因、升级过程、异常情况、验证结果等"
            ></textarea>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showEditDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="saveUpgrade">
            保存
          </button>
        </div>
      </div>
    </div>

    <!-- 查看详情弹窗 -->
    <div v-if="selectedUpgrade" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>升级记录详情</h3>
          <button @click="selectedUpgrade = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>绑定项目</span>
            <strong>{{ selectedUpgrade.projectName }}</strong>
          </div>

          <div>
            <span>终端类型</span>
            <strong>{{ selectedUpgrade.deviceType }}</strong>
          </div>

          <div>
            <span>SN序列号</span>
            <strong>{{ selectedUpgrade.sn }}</strong>
          </div>

          <div>
            <span>MAC地址</span>
            <strong>{{ selectedUpgrade.macAddress }}</strong>
          </div>

          <div>
            <span>升级前软件版本</span>
            <strong>{{ selectedUpgrade.oldSoftwareVersion }}</strong>
          </div>

          <div>
            <span>升级后软件版本</span>
            <strong>{{ selectedUpgrade.newSoftwareVersion }}</strong>
          </div>

          <div>
            <span>升级状态</span>
            <strong>{{ getUpgradeStatusText(selectedUpgrade.upgradeStatus) }}</strong>
          </div>

          <div>
            <span>升级方式</span>
            <strong>{{ selectedUpgrade.upgradeMethod }}</strong>
          </div>

          <div>
            <span>升级人</span>
            <strong>{{ selectedUpgrade.upgradeUser }}</strong>
          </div>

          <div>
            <span>升级时间</span>
            <strong>{{ selectedUpgrade.upgradeTime }}</strong>
          </div>

          <div>
            <span>最后修改时间</span>
            <strong>{{ selectedUpgrade.updateTime || '-' }}</strong>
          </div>
        </div>

        <div class="remark-card">
          <span>升级说明</span>
          <p>{{ selectedUpgrade.remark || '暂无说明' }}</p>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="openEditDialog(selectedUpgrade)">
            修改记录
          </button>

          <button class="primary-btn" @click="selectedUpgrade = null">
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
  upgradeStatus: ''
})

const showEditDialog = ref(false)
const selectedUpgrade = ref(null)
const currentEditUpgrade = ref(null)
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

const upgradeForm = reactive({
  projectName: '',
  deviceType: '',
  sn: '',
  macAddress: '',
  oldSoftwareVersion: '',
  newSoftwareVersion: '',
  upgradeStatus: 'success',
  upgradeUser: '',
  upgradeTime: '',
  upgradeMethod: '远程升级',
  remark: ''
})

const upgradeList = ref([
  {
    id: 1,
    projectName: '香港屯马项目',
    deviceType: '司机室控制盒',
    sn: 'DCCU-202605100001',
    macAddress: '00:11:22:33:44:01',
    oldSoftwareVersion: 'SW-DCCU.V1.1.0',
    newSoftwareVersion: 'SW-DCCU.V1.2.0',
    upgradeStatus: 'success',
    upgradeUser: '售后人员',
    upgradeTime: '2026-05-10',
    upgradeMethod: '远程升级',
    updateTime: '2026-05-10',
    remark: '远程升级成功，升级后广播、对讲、网络通信功能验证正常。'
  },
  {
    id: 2,
    projectName: '波尔图二期项目',
    deviceType: '解码板',
    sn: 'DEC-202605110001',
    macAddress: '00:11:22:55:66:01',
    oldSoftwareVersion: 'SW-DECODER.V1.0.2',
    newSoftwareVersion: 'SW-DECODER.V1.0.3',
    upgradeStatus: 'success',
    upgradeUser: '售后人员',
    upgradeTime: '2026-05-16',
    upgradeMethod: '现场升级',
    updateTime: '2026-05-16',
    remark: '现场升级完成，已验证客室广播解码功能。'
  },
  {
    id: 3,
    projectName: '阿根廷有轨项目',
    deviceType: '司机提醒单元',
    sn: 'DRU-202605120001',
    macAddress: '00:11:22:77:88:01',
    oldSoftwareVersion: 'SW-DRU.V0.9.0',
    newSoftwareVersion: 'SW-DRU.V1.0.0',
    upgradeStatus: 'pending',
    upgradeUser: '研发人员',
    upgradeTime: '2026-05-18',
    upgradeMethod: '生产重新烧录',
    updateTime: '2026-05-18',
    remark: '已完成软件升级，等待测试人员确认升级结果。'
  },
  {
    id: 4,
    projectName: '迪拜项目',
    deviceType: '编码板',
    sn: 'ENC-202605140001',
    macAddress: '00:11:22:BB:CC:01',
    oldSoftwareVersion: 'SW-ENCODER.V1.2.0',
    newSoftwareVersion: 'SW-ENCODER.V1.3.0',
    upgradeStatus: 'failed',
    upgradeUser: '售后人员',
    upgradeTime: '2026-05-20',
    upgradeMethod: '远程升级',
    updateTime: '2026-05-20',
    remark: '远程升级失败，升级过程中网络中断。该记录为异常记录，可删除后重新登记。'
  }
])

const filteredUpgradeList = computed(() => {
  return upgradeList.value.filter(item => {
    const keywordMatch =
      !filters.keyword ||
      item.projectName.includes(filters.keyword) ||
      item.deviceType.includes(filters.keyword) ||
      item.sn.includes(filters.keyword) ||
      item.macAddress.includes(filters.keyword) ||
      item.oldSoftwareVersion.includes(filters.keyword) ||
      item.newSoftwareVersion.includes(filters.keyword) ||
      item.upgradeUser.includes(filters.keyword)

    const projectMatch =
      !filters.projectName || item.projectName === filters.projectName

    const deviceTypeMatch =
      !filters.deviceType || item.deviceType === filters.deviceType

    const statusMatch =
      !filters.upgradeStatus || item.upgradeStatus === filters.upgradeStatus

    return keywordMatch && projectMatch && deviceTypeMatch && statusMatch
  })
})

function getUpgradeStatusText(status) {
  const map = {
    success: '升级成功',
    pending: '待确认',
    failed: '升级失败'
  }

  return map[status] || status
}

function resetFilters() {
  filters.keyword = ''
  filters.projectName = ''
  filters.deviceType = ''
  filters.upgradeStatus = ''
}

function openCreateDialog() {
  editMode.value = 'create'
  currentEditUpgrade.value = null

  upgradeForm.projectName = ''
  upgradeForm.deviceType = ''
  upgradeForm.sn = ''
  upgradeForm.macAddress = ''
  upgradeForm.oldSoftwareVersion = ''
  upgradeForm.newSoftwareVersion = ''
  upgradeForm.upgradeStatus = 'success'
  upgradeForm.upgradeUser = ''
  upgradeForm.upgradeTime = new Date().toISOString().slice(0, 10)
  upgradeForm.upgradeMethod = '远程升级'
  upgradeForm.remark = ''

  showEditDialog.value = true
}

function openEditDialog(item) {
  editMode.value = 'edit'
  currentEditUpgrade.value = item

  upgradeForm.projectName = item.projectName
  upgradeForm.deviceType = item.deviceType
  upgradeForm.sn = item.sn
  upgradeForm.macAddress = item.macAddress
  upgradeForm.oldSoftwareVersion = item.oldSoftwareVersion
  upgradeForm.newSoftwareVersion = item.newSoftwareVersion
  upgradeForm.upgradeStatus = item.upgradeStatus
  upgradeForm.upgradeUser = item.upgradeUser
  upgradeForm.upgradeTime = item.upgradeTime
  upgradeForm.upgradeMethod = item.upgradeMethod
  upgradeForm.remark = item.remark

  showEditDialog.value = true
}

function saveUpgrade() {
  if (!upgradeForm.projectName) {
    alert('请选择绑定项目')
    return
  }

  if (!upgradeForm.deviceType) {
    alert('请选择终端设备')
    return
  }

  if (!upgradeForm.sn) {
    alert('请输入SN序列号')
    return
  }

  if (!upgradeForm.macAddress) {
    alert('请输入MAC地址')
    return
  }

  if (!upgradeForm.oldSoftwareVersion) {
    alert('请输入升级前软件版本')
    return
  }

  if (!upgradeForm.newSoftwareVersion) {
    alert('请输入升级后软件版本')
    return
  }

  if (!upgradeForm.upgradeUser) {
    alert('请输入升级人')
    return
  }

  if (editMode.value === 'create') {
    upgradeList.value.unshift({
      id: Date.now(),
      projectName: upgradeForm.projectName,
      deviceType: upgradeForm.deviceType,
      sn: upgradeForm.sn,
      macAddress: upgradeForm.macAddress,
      oldSoftwareVersion: upgradeForm.oldSoftwareVersion,
      newSoftwareVersion: upgradeForm.newSoftwareVersion,
      upgradeStatus: upgradeForm.upgradeStatus,
      upgradeUser: upgradeForm.upgradeUser,
      upgradeTime: upgradeForm.upgradeTime,
      upgradeMethod: upgradeForm.upgradeMethod,
      updateTime: new Date().toISOString().slice(0, 10),
      remark: upgradeForm.remark
    })
  } else {
    currentEditUpgrade.value.projectName = upgradeForm.projectName
    currentEditUpgrade.value.deviceType = upgradeForm.deviceType
    currentEditUpgrade.value.sn = upgradeForm.sn
    currentEditUpgrade.value.macAddress = upgradeForm.macAddress
    currentEditUpgrade.value.oldSoftwareVersion = upgradeForm.oldSoftwareVersion
    currentEditUpgrade.value.newSoftwareVersion = upgradeForm.newSoftwareVersion
    currentEditUpgrade.value.upgradeStatus = upgradeForm.upgradeStatus
    currentEditUpgrade.value.upgradeUser = upgradeForm.upgradeUser
    currentEditUpgrade.value.upgradeTime = upgradeForm.upgradeTime
    currentEditUpgrade.value.upgradeMethod = upgradeForm.upgradeMethod
    currentEditUpgrade.value.updateTime = new Date().toISOString().slice(0, 10)
    currentEditUpgrade.value.remark = upgradeForm.remark
  }

  showEditDialog.value = false
}

function viewUpgrade(item) {
  selectedUpgrade.value = item
}

function deleteUpgrade(item) {
  const ok = confirm(`确认删除升级记录【${item.sn}】吗？`)
  if (!ok) return

  upgradeList.value = upgradeList.value.filter(record => record.id !== item.id)
}

function exportUpgradeRecords() {
  const header = [
    '绑定项目',
    '终端类型',
    'SN序列号',
    'MAC地址',
    '升级前软件版本',
    '升级后软件版本',
    '升级状态',
    '升级方式',
    '升级人',
    '升级时间',
    '最后修改时间',
    '升级说明'
  ]

  const rows = upgradeList.value.map(item => [
    item.projectName,
    item.deviceType,
    item.sn,
    item.macAddress,
    item.oldSoftwareVersion,
    item.newSoftwareVersion,
    getUpgradeStatusText(item.upgradeStatus),
    item.upgradeMethod,
    item.upgradeUser,
    item.upgradeTime,
    item.updateTime || '',
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
  link.download = '软件升级记录.csv'
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
  min-width: 1500px;
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

.version-cell {
  overflow: hidden;
}

.software-tag,
.old-version-tag {
  display: inline-block;
  max-width: 170px;
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

.old-version-tag {
  background: #47556933;
  color: #cbd5e1;
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

.status-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
}

.status-tag.success {
  background: #16a34a33;
  color: #4ade80;
}

.status-tag.pending {
  background: #d9770633;
  color: #fbbf24;
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
    min-width: 1500px;
  }

  .form-grid,
  .detail-card {
    grid-template-columns: 1fr;
  }
}
</style>