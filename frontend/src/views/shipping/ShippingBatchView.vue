<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>发货批次管理</h1>
      </div>

      <button class="primary-btn" @click="openCreateDialog">
        新增发货批次
      </button>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索批次号 / SN序列号 / 上传人 / 发货单文件 / 快递单号"
      />

      <select v-model="filters.auditStatus">
        <option value="">全部审核状态</option>
        <option value="draft">草稿</option>
        <option value="submitted">待审核</option>
        <option value="approved">审核通过</option>
        <option value="rejected">审核驳回</option>
      </select>

      <button class="query-btn">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <!-- 待生成发货批次提示 -->
    <div v-if="pendingShippingDevices.length > 0" class="pending-card">
      <div>
        <strong>已从出库管理选择 {{ pendingShippingDevices.length }} 台设备</strong>
        <p>点击“新增发货批次”后，系统会自动带入这些设备的 SN 序列号和设备数量。</p>
      </div>

      <button class="primary-btn" @click="openCreateDialog">
        使用已选设备生成发货批次
      </button>
    </div>

    <!-- 发货批次表格 -->
    <div class="table-card">
      <div class="table-card-header">
        <div>
          <h3>发货批次列表</h3>
          <span>共 {{ filteredBatchList.length }} 条发货批次</span>
        </div>
      </div>

      <div class="table-wrapper">
        <table class="version-table">
          <thead>
            <tr>
              <th>发货批次</th>
              <th>发货单文件</th>
              <th>设备数量</th>
              <th>上传人</th>
              <th>上传时间</th>
              <th>审核状态</th>
              <th>审核人</th>
              <th class="operation-col">操作</th>
              <th>快递单号</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="item in filteredBatchList" :key="item.id">
              <td>
                <button class="record-link" @click="viewBatch(item)">
                  {{ item.batchNo }}
                </button>
              </td>

              <td>
                <span class="file-tag" :title="item.fileName">
                  {{ item.fileName }}
                </span>
              </td>

              <td>
                <span class="count-tag">
                  {{ item.deviceCount }} 台
                </span>
              </td>

              <td>{{ item.uploader }}</td>

              <td class="muted">{{ item.uploadTime }}</td>

              <td>
                <span class="status-tag" :class="item.auditStatus">
                  {{ getAuditStatusText(item.auditStatus) }}
                </span>
              </td>

              <td>{{ item.auditor || '-' }}</td>

              <td class="operation-col">
                <div class="action-group">
                  <button class="text-btn" @click="viewBatch(item)">
                    查看
                  </button>

                  <button class="text-btn blue" @click="downloadFile(item)">
                    下载
                  </button>

                  <button
                    v-if="item.auditStatus === 'draft' || item.auditStatus === 'rejected'"
                    class="text-btn blue"
                    @click="submitBatch(item)"
                  >
                    提交
                  </button>

                  <button
                    v-if="item.auditStatus === 'submitted'"
                    class="text-btn green"
                    @click="auditBatch(item)"
                  >
                    审核
                  </button>

                  <button
                    v-if="item.auditStatus === 'draft' || item.auditStatus === 'rejected'"
                    class="text-btn red"
                    @click="deleteBatch(item)"
                  >
                    删除
                  </button>
                </div>
              </td>

              <td>
                <span class="express-tag" :title="item.expressNo">
                  {{ item.expressNo || '-' }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="table-footer">
        支持按发货批次号、SN序列号、上传人、发货单文件、快递单号查询。
      </div>
    </div>

    <!-- 新增发货批次弹窗 -->
    <div v-if="showCreateDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>新增发货批次</h3>
          <button @click="showCreateDialog = false">×</button>
        </div>

        <div class="form-grid">
          <label>
            发货批次号
            <input
              v-model="batchForm.batchNo"
              placeholder="例如：0000013289"
            />
          </label>

          <label>
            设备数量
            <input
              v-model.number="batchForm.deviceCount"
              type="number"
              min="1"
              disabled
            />
          </label>

          <label>
            上传人
            <input
              v-model="batchForm.uploader"
              placeholder="请输入上传人"
            />
          </label>

          <label>
            快递单号
            <input
              v-model="batchForm.expressNo"
              placeholder="例如：SF1234567890"
            />
          </label>

          <label class="full-row">
            已选出库设备 SN
            <div class="selected-device-panel">
              <div v-if="batchForm.snList.length === 0" class="empty-device">
                暂无已选设备，请先到出库管理页面勾选设备并跳转到发货管理。
              </div>

              <div
                v-for="device in batchForm.deviceList"
                v-else
                :key="device.outboundId || device.id"
                class="selected-device-item"
              >
                <span>{{ device.deviceType }}</span>
                <strong>{{ device.sn }}</strong>
                <em>{{ device.macAddress }}</em>
              </div>
            </div>
          </label>

          <label class="full-row">
            发货单 Excel 文件
            <input
              type="file"
              accept=".xls,.xlsx"
              @change="handleFileChange"
            />
          </label>

          <label class="full-row">
            文件名称
            <input
              v-model="batchForm.fileName"
              placeholder="选择文件后自动填充"
              disabled
            />
          </label>

          <label class="full-row">
            发货说明
            <textarea
              v-model="batchForm.remark"
              placeholder="例如：本批次发货单包含SN、MAC、终端类型、软件版本、硬件版本等信息"
            ></textarea>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showCreateDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="createBatch">
            保存
          </button>
        </div>
      </div>
    </div>

    <!-- 查看详情 / 审核弹窗 -->
    <div v-if="selectedBatch" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>发货批次详情</h3>
          <button @click="selectedBatch = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>发货批次号</span>
            <strong>{{ selectedBatch.batchNo }}</strong>
          </div>

          <div>
            <span>快递单号</span>
            <strong>{{ selectedBatch.expressNo || '-' }}</strong>
          </div>

          <div>
            <span>发货单文件</span>
            <strong>{{ selectedBatch.fileName }}</strong>
          </div>

          <div>
            <span>设备数量</span>
            <strong>{{ selectedBatch.deviceCount }} 台</strong>
          </div>

          <div>
            <span>上传人</span>
            <strong>{{ selectedBatch.uploader }}</strong>
          </div>

          <div>
            <span>上传时间</span>
            <strong>{{ selectedBatch.uploadTime }}</strong>
          </div>

          <div>
            <span>审核状态</span>
            <strong>{{ getAuditStatusText(selectedBatch.auditStatus) }}</strong>
          </div>

          <div>
            <span>审核人</span>
            <strong>{{ selectedBatch.auditor || '-' }}</strong>
          </div>

          <div>
            <span>审核时间</span>
            <strong>{{ selectedBatch.auditTime || '-' }}</strong>
          </div>

          <div>
            <span>文件查看</span>
            <button class="inline-link" @click="openFile(selectedBatch)">
              点开查看文件
            </button>
          </div>
        </div>

        <div class="remark-card">
          <span>SN序列号</span>
          <div class="sn-list">
            <span
              v-for="sn in selectedBatch.snList"
              :key="sn"
              class="sn-tag"
              :title="sn"
            >
              {{ sn }}
            </span>
            <p v-if="!selectedBatch.snList || selectedBatch.snList.length === 0">
              暂无 SN 序列号
            </p>
          </div>
        </div>

        <div class="remark-card">
          <span>发货说明</span>
          <p>{{ selectedBatch.remark || '暂无说明' }}</p>
        </div>

        <div class="dialog-footer">
          <button
            v-if="selectedBatch.auditStatus === 'submitted'"
            class="green-btn"
            @click="approveBatch(selectedBatch)"
          >
            审核通过
          </button>

          <button
            v-if="selectedBatch.auditStatus === 'submitted'"
            class="red-btn"
            @click="rejectBatch(selectedBatch)"
          >
            审核驳回
          </button>

          <button class="reset-btn" @click="downloadFile(selectedBatch)">
            下载文件
          </button>

          <button class="primary-btn" @click="selectedBatch = null">
            关闭
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'

const STORAGE_PENDING_SHIPPING_KEY = 'pendingShippingDevices'

const filters = reactive({
  keyword: '',
  auditStatus: ''
})

const showCreateDialog = ref(false)
const selectedBatch = ref(null)

const pendingShippingDevices = ref(readStorageList(STORAGE_PENDING_SHIPPING_KEY, []))

const batchForm = reactive({
  batchNo: '',
  deviceCount: 0,
  uploader: '',
  expressNo: '',
  snList: [],
  deviceList: [],
  fileName: '',
  file: null,
  fileUrl: '',
  remark: ''
})

const batchList = ref([
  {
    id: 1,
    batchNo: '0000013289',
    deviceCount: 3,
    expressNo: 'SF202605100001',
    snList: [
      'DCCU-202605100001',
      'DCCU-202605100002',
      'DCCU-202605100003'
    ],
    fileName: '5月第一批司机室控制盒发货单.xlsx',
    fileUrl: '',
    uploader: '袁晓兰',
    uploadTime: '2026-05-10',
    auditStatus: 'approved',
    auditor: '领导',
    auditTime: '2026-05-11',
    remark: '本批次包含司机室控制盒，发货单记录SN、MAC、软件版本、硬件版本。'
  },
  {
    id: 2,
    batchNo: '0000013290',
    deviceCount: 3,
    expressNo: 'SF202605160001',
    snList: [
      'DEC-202605110001',
      'DEC-202605110002',
      'DEC-202605110003'
    ],
    fileName: '5月第二批解码板发货单.xlsx',
    fileUrl: '',
    uploader: '袁晓兰',
    uploadTime: '2026-05-16',
    auditStatus: 'submitted',
    auditor: '',
    auditTime: '',
    remark: '待领导审核确认。'
  }
])

const filteredBatchList = computed(() => {
  return batchList.value.filter(item => {
    const keyword = filters.keyword.trim()

    const snMatch =
      item.snList &&
      item.snList.some(sn => sn.includes(keyword))

    const keywordMatch =
      !keyword ||
      item.batchNo.includes(keyword) ||
      item.uploader.includes(keyword) ||
      item.fileName.includes(keyword) ||
      item.expressNo.includes(keyword) ||
      snMatch

    const auditStatusMatch =
      !filters.auditStatus || item.auditStatus === filters.auditStatus

    return keywordMatch && auditStatusMatch
  })
})

function readStorageList(key, fallback) {
  try {
    const raw = localStorage.getItem(key)
    if (!raw) return fallback
    return JSON.parse(raw)
  } catch (error) {
    return fallback
  }
}

function saveStorageList(key, value) {
  localStorage.setItem(key, JSON.stringify(value))
}

function getCurrentUserName() {
  return (
    localStorage.getItem('username') ||
    localStorage.getItem('accountName') ||
    localStorage.getItem('realName') ||
    '当前用户'
  )
}

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
  filters.auditStatus = ''
}

function openCreateDialog() {
  const pendingDevices = readStorageList(STORAGE_PENDING_SHIPPING_KEY, [])
  pendingShippingDevices.value = pendingDevices

  batchForm.batchNo = ''
  batchForm.deviceList = pendingDevices
  batchForm.snList = pendingDevices.map(item => item.sn).filter(Boolean)
  batchForm.deviceCount = batchForm.snList.length
  batchForm.uploader = getCurrentUserName()
  batchForm.expressNo = ''
  batchForm.fileName = ''
  batchForm.file = null
  batchForm.fileUrl = ''
  batchForm.remark = ''

  showCreateDialog.value = true
}

function handleFileChange(event) {
  const file = event.target.files[0]
  if (!file) return

  batchForm.file = file
  batchForm.fileName = file.name
  batchForm.fileUrl = URL.createObjectURL(file)
}

function createBatch() {
  if (!batchForm.batchNo) {
    alert('请输入发货批次号')
    return
  }

  if (!batchForm.deviceCount || batchForm.deviceCount <= 0) {
    alert('暂无已选设备，请先到出库管理页面勾选设备')
    return
  }

  if (!batchForm.uploader) {
    alert('上传人读取失败')
    return
  }

  if (!batchForm.file) {
    alert('请上传发货单 Excel 文件')
    return
  }

  batchList.value.unshift({
    id: Date.now(),
    batchNo: batchForm.batchNo,
    deviceCount: batchForm.deviceCount,
    expressNo: batchForm.expressNo,
    snList: [...batchForm.snList],
    deviceList: [...batchForm.deviceList],
    fileName: batchForm.fileName,
    fileUrl: batchForm.fileUrl,
    uploader: batchForm.uploader,
    uploadTime: new Date().toISOString().slice(0, 10),
    auditStatus: 'draft',
    auditor: '',
    auditTime: '',
    remark: batchForm.remark
  })

  localStorage.removeItem(STORAGE_PENDING_SHIPPING_KEY)
  pendingShippingDevices.value = []
  showCreateDialog.value = false

  alert('发货批次已生成，SN 已从出库管理页面自动带入')
}

function viewBatch(item) {
  selectedBatch.value = item
}

function auditBatch(item) {
  selectedBatch.value = item
}

function submitBatch(item) {
  item.auditStatus = 'submitted'
  item.auditor = ''
  item.auditTime = ''
  alert(`发货批次【${item.batchNo}】已提交领导审核`)
}

function approveBatch(item) {
  item.auditStatus = 'approved'
  item.auditor = '领导'
  item.auditTime = new Date().toISOString().slice(0, 10)
  selectedBatch.value = null
  alert(`发货批次【${item.batchNo}】审核通过`)
}

function rejectBatch(item) {
  item.auditStatus = 'rejected'
  item.auditor = '领导'
  item.auditTime = new Date().toISOString().slice(0, 10)
  selectedBatch.value = null
  alert(`发货批次【${item.batchNo}】已驳回`)
}

function openFile(item) {
  if (!item.fileUrl) {
    alert('当前是模拟数据，暂无可直接打开的原始文件')
    return
  }

  window.open(item.fileUrl, '_blank')
}

function downloadFile(item) {
  if (!item.fileUrl) {
    alert('当前是模拟数据，暂无可下载的原始文件')
    return
  }

  const link = document.createElement('a')
  link.href = item.fileUrl
  link.download = item.fileName || '发货单.xlsx'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

function deleteBatch(item) {
  const ok = confirm(`确认删除发货批次【${item.batchNo}】吗？`)
  if (!ok) return

  batchList.value = batchList.value.filter(record => record.id !== item.id)
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
  grid-template-columns: 1.4fr 180px 90px 90px;
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

.pending-card {
  background: #0f172a;
  border: 1px solid #1d4ed8;
  border-radius: 14px;
  padding: 16px;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.pending-card strong {
  color: #f8fafc;
  font-size: 15px;
}

.pending-card p {
  margin: 6px 0 0;
  color: #94a3b8;
  font-size: 12px;
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
  min-width: 1320px;
  border-collapse: collapse;
  table-layout: fixed;
}

.version-table thead {
  background: #020617;
}

.version-table th,
.version-table td {
  box-sizing: border-box;
  white-space: nowrap;
}

.version-table th {
  padding: 14px 16px;
  color: #94a3b8;
  font-size: 12px;
  font-weight: 600;
  text-align: left;
  border-bottom: 1px solid #1e293b;
}

.version-table td {
  padding: 15px 16px;
  font-size: 13px;
  color: #e2e8f0;
  border-bottom: 1px solid #1e293b;
  vertical-align: middle;
  overflow: hidden;
}

.record-link {
  display: inline-block;
  max-width: 130px;
  border: none;
  background: transparent;
  color: #60a5fa;
  font-size: 13px;
  font-weight: 800;
  cursor: pointer;
  padding: 0;
  text-align: left;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.record-link:hover {
  color: #93c5fd;
  text-decoration: underline;
}

.file-tag,
.express-tag {
  display: inline-block;
  width: 150px;
  max-width: 150px;
  padding: 4px 9px;
  border-radius: 999px;
  background: #33415566;
  color: #cbd5e1;
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;
}

.express-tag {
  background: #1d4ed833;
  color: #60a5fa;
  font-family: Consolas, Monaco, monospace;
}

.count-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  background: #0f766e33;
  color: #5eead4;
  font-size: 12px;
  font-weight: 700;
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
  width: 300px;
  text-align: right !important;
}

.action-group {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 12px;
  flex-wrap: nowrap;
  white-space: nowrap;
}

.text-btn {
  border: none;
  background: transparent;
  color: #cbd5e1;
  font-size: 13px;
  cursor: pointer;
  white-space: nowrap;
  padding: 0;
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
  max-height: 92vh;
  overflow-y: auto;
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

.selected-device-panel {
  min-height: 120px;
  max-height: 240px;
  overflow-y: auto;
  border: 1px solid #334155;
  border-radius: 10px;
  background: #020617;
  padding: 10px;
}

.selected-device-item {
  display: grid;
  grid-template-columns: 140px 1fr 180px;
  align-items: center;
  gap: 10px;
  padding: 8px;
  border-bottom: 1px solid #1e293b;
}

.selected-device-item:last-child {
  border-bottom: none;
}

.selected-device-item span {
  color: #5eead4;
  font-size: 12px;
}

.selected-device-item strong,
.selected-device-item em {
  color: #cbd5e1;
  font-size: 12px;
  font-family: Consolas, Monaco, monospace;
  font-style: normal;
}

.empty-device {
  color: #64748b;
  font-size: 13px;
  padding: 10px 0;
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

.sn-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.sn-tag {
  display: inline-block;
  max-width: 180px;
  padding: 4px 9px;
  border-radius: 999px;
  background: #33415566;
  color: #cbd5e1;
  font-size: 12px;
  font-family: Consolas, Monaco, monospace;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

@media (max-width: 960px) {
  .filter-card {
    grid-template-columns: 1fr;
  }

  .version-table {
    min-width: 1320px;
  }

  .form-grid,
  .detail-card {
    grid-template-columns: 1fr;
  }

  .selected-device-item {
    grid-template-columns: 1fr;
  }
}
</style>