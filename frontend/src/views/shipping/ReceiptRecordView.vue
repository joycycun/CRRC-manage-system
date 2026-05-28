<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>收货记录管理</h1>
      </div>

      <div class="header-actions">
        <button class="reset-btn" @click="exportReceiptRecords">
          导出收货记录
        </button>

        <button class="primary-btn" @click="openCreateDialog">
          新增快递单号
        </button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="summary-panel">
      <div class="summary-card">
        <div class="summary-icon blue">收</div>
        <div>
          <span>收货记录总数</span>
          <strong>{{ totalCount }}</strong>
          <p>当前系统记录的快递收货单号</p>
        </div>
      </div>

      <div class="summary-card green">
        <div class="summary-icon green">签</div>
        <div>
          <span>已确认收货</span>
          <strong>{{ receivedCount }}</strong>
          <p>已通过物流状态自动确认收货</p>
        </div>
      </div>

      <div class="summary-card yellow">
        <div class="summary-icon yellow">运</div>
        <div>
          <span>运输中</span>
          <strong>{{ shippingCount }}</strong>
          <p>顺丰物流仍在运输中的快递</p>
        </div>
      </div>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索快递单号 / 收货人 / 上传人 / 备注"
      />

      <select v-model="filters.receiptStatus">
        <option value="">全部收货状态</option>
        <option value="shipping">运输中</option>
        <option value="received">已收货</option>
        <option value="exception">物流异常</option>
      </select>

      <button class="query-btn">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <!-- 收货记录表格 -->
    <div class="table-card">
      <div class="table-card-header">
        <div>
          <h3>收货记录列表</h3>
          <span>共 {{ filteredReceiptList.length }} 条收货记录</span>
        </div>
      </div>

      <div class="table-wrapper">
        <table class="version-table">
          <thead>
            <tr>
              <th>快递单号</th>
              <th>快递公司</th>
              <th>收货人</th>
              <th>收货状态</th>
              <th>物流状态</th>
              <th>新增人</th>
              <th>新增时间</th>
              <th>确认收货时间</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="item in filteredReceiptList" :key="item.id">
              <td>
                <button class="record-link" @click="viewReceipt(item)">
                  {{ item.expressNo }}
                </button>
              </td>

              <td>
                <span class="company-tag">
                  {{ item.expressCompany }}
                </span>
              </td>

              <td>{{ item.receiver }}</td>

              <td>
                <span class="status-tag" :class="item.receiptStatus">
                  {{ getReceiptStatusText(item.receiptStatus) }}
                </span>
              </td>

              <td>
                <span class="logistics-text" :title="item.logisticsStatus">
                  {{ item.logisticsStatus }}
                </span>
              </td>

              <td>{{ item.creator }}</td>

              <td class="muted">{{ item.createTime }}</td>

              <td class="muted">{{ item.receiptTime || '-' }}</td>

              <td class="operation-col">
                <div class="action-group">
                  <button class="text-btn" @click="viewReceipt(item)">
                    查看
                  </button>

                  <button
                    class="text-btn blue"
                    @click="querySfLogistics(item)"
                  >
                    查询物流
                  </button>

                  <button
                    v-if="item.receiptStatus !== 'received'"
                    class="text-btn green"
                    @click="confirmReceipt(item)"
                  >
                    确认收货
                  </button>

                  <button
                    class="text-btn red"
                    @click="deleteReceipt(item)"
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

    <!-- 新增快递单号弹窗 -->
    <div v-if="showCreateDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>新增快递单号</h3>
          <button @click="showCreateDialog = false">×</button>
        </div>

        <div class="form-grid">
          <label>
            快递单号
            <input
              v-model="receiptForm.expressNo"
              placeholder="请输入顺丰快递单号"
            />
          </label>

          <label>
            快递公司
            <select v-model="receiptForm.expressCompany">
              <option value="顺丰速运">顺丰速运</option>
            </select>
          </label>

          <label>
            收货人
            <input
              v-model="receiptForm.receiver"
              placeholder="请输入收货人"
            />
          </label>

          <label>
            新增人
            <input
              v-model="receiptForm.creator"
              placeholder="请输入新增人"
            />
          </label>

          <label class="full-row">
            备注
            <textarea
              v-model="receiptForm.remark"
              placeholder="例如：该快递为返修设备、样机设备、生产物料或客户回寄设备"
            ></textarea>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showCreateDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="createReceipt">
            保存
          </button>
        </div>
      </div>
    </div>

    <!-- 查看详情弹窗 -->
    <div v-if="selectedReceipt" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>收货记录详情</h3>
          <button @click="selectedReceipt = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>快递单号</span>
            <strong>{{ selectedReceipt.expressNo }}</strong>
          </div>

          <div>
            <span>快递公司</span>
            <strong>{{ selectedReceipt.expressCompany }}</strong>
          </div>

          <div>
            <span>收货人</span>
            <strong>{{ selectedReceipt.receiver }}</strong>
          </div>

          <div>
            <span>收货状态</span>
            <strong>{{ getReceiptStatusText(selectedReceipt.receiptStatus) }}</strong>
          </div>

          <div>
            <span>物流状态</span>
            <strong>{{ selectedReceipt.logisticsStatus }}</strong>
          </div>

          <div>
            <span>新增人</span>
            <strong>{{ selectedReceipt.creator }}</strong>
          </div>

          <div>
            <span>新增时间</span>
            <strong>{{ selectedReceipt.createTime }}</strong>
          </div>

          <div>
            <span>确认收货时间</span>
            <strong>{{ selectedReceipt.receiptTime || '-' }}</strong>
          </div>

          <div>
            <span>最后查询时间</span>
            <strong>{{ selectedReceipt.lastQueryTime || '-' }}</strong>
          </div>

          <div>
            <span>顺丰接口状态</span>
            <strong>{{ selectedReceipt.apiStatus || '未查询' }}</strong>
          </div>
        </div>

        <div class="remark-card">
          <span>备注</span>
          <p>{{ selectedReceipt.remark || '暂无备注' }}</p>
        </div>

        <div class="remark-card">
          <span>物流轨迹</span>
          <div
            v-if="selectedReceipt.traces.length === 0"
            class="empty-trace"
          >
            暂无物流轨迹，请点击“查询物流”
          </div>

          <div
            v-for="trace in selectedReceipt.traces"
            :key="trace.id"
            class="trace-item"
          >
            <strong>{{ trace.time }}</strong>
            <p>{{ trace.desc }}</p>
          </div>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="querySfLogistics(selectedReceipt)">
            查询物流
          </button>

          <button
            v-if="selectedReceipt.receiptStatus !== 'received'"
            class="green-btn"
            @click="confirmReceipt(selectedReceipt)"
          >
            确认收货
          </button>

          <button class="primary-btn" @click="selectedReceipt = null">
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
  receiptStatus: ''
})

const showCreateDialog = ref(false)
const selectedReceipt = ref(null)

const receiptForm = reactive({
  expressNo: '',
  expressCompany: '顺丰速运',
  receiver: '',
  creator: '',
  remark: ''
})

const receiptList = ref([
  {
    id: 1,
    expressNo: 'SF202605100001',
    expressCompany: '顺丰速运',
    receiver: '袁晓兰',
    receiptStatus: 'received',
    logisticsStatus: '已签收',
    creator: '售后人员',
    createTime: '2026-05-10',
    receiptTime: '2026-05-11',
    lastQueryTime: '2026-05-11',
    apiStatus: '查询成功',
    remark: '客户回寄返修设备，已确认收货。',
    traces: [
      {
        id: 101,
        time: '2026-05-10 14:20',
        desc: '顺丰已揽收'
      },
      {
        id: 102,
        time: '2026-05-11 10:15',
        desc: '快件已签收，签收人为袁晓兰'
      }
    ]
  },
  {
    id: 2,
    expressNo: 'SF202605160002',
    expressCompany: '顺丰速运',
    receiver: '生产人员',
    receiptStatus: 'shipping',
    logisticsStatus: '运输中',
    creator: '生产人员',
    createTime: '2026-05-16',
    receiptTime: '',
    lastQueryTime: '2026-05-16',
    apiStatus: '查询成功',
    remark: '供应商寄回的生产物料，等待到货。',
    traces: [
      {
        id: 201,
        time: '2026-05-16 09:30',
        desc: '快件已从深圳转运中心发出'
      }
    ]
  },
  {
    id: 3,
    expressNo: 'SF202605180003',
    expressCompany: '顺丰速运',
    receiver: '测试人员',
    receiptStatus: 'exception',
    logisticsStatus: '派送异常',
    creator: '测试人员',
    createTime: '2026-05-18',
    receiptTime: '',
    lastQueryTime: '2026-05-18',
    apiStatus: '查询成功',
    remark: '样机设备快递，当前派送异常，需要联系快递员。',
    traces: [
      {
        id: 301,
        time: '2026-05-18 16:40',
        desc: '派送异常，收件人电话暂时无法接通'
      }
    ]
  }
])

const filteredReceiptList = computed(() => {
  return receiptList.value.filter(item => {
    const keywordMatch =
      !filters.keyword ||
      item.expressNo.includes(filters.keyword) ||
      item.receiver.includes(filters.keyword) ||
      item.creator.includes(filters.keyword) ||
      item.remark.includes(filters.keyword)

    const statusMatch =
      !filters.receiptStatus || item.receiptStatus === filters.receiptStatus

    return keywordMatch && statusMatch
  })
})

const totalCount = computed(() => receiptList.value.length)

const receivedCount = computed(() => {
  return receiptList.value.filter(item => item.receiptStatus === 'received').length
})

const shippingCount = computed(() => {
  return receiptList.value.filter(item => item.receiptStatus === 'shipping').length
})

function getReceiptStatusText(status) {
  const map = {
    shipping: '运输中',
    received: '已收货',
    exception: '物流异常'
  }

  return map[status] || status
}

function resetFilters() {
  filters.keyword = ''
  filters.receiptStatus = ''
}

function openCreateDialog() {
  receiptForm.expressNo = ''
  receiptForm.expressCompany = '顺丰速运'
  receiptForm.receiver = ''
  receiptForm.creator = ''
  receiptForm.remark = ''

  showCreateDialog.value = true
}

function createReceipt() {
  if (!receiptForm.expressNo) {
    alert('请输入快递单号')
    return
  }

  if (!receiptForm.receiver) {
    alert('请输入收货人')
    return
  }

  if (!receiptForm.creator) {
    alert('请输入新增人')
    return
  }

  receiptList.value.unshift({
    id: Date.now(),
    expressNo: receiptForm.expressNo,
    expressCompany: receiptForm.expressCompany,
    receiver: receiptForm.receiver,
    receiptStatus: 'shipping',
    logisticsStatus: '待查询',
    creator: receiptForm.creator,
    createTime: new Date().toISOString().slice(0, 10),
    receiptTime: '',
    lastQueryTime: '',
    apiStatus: '未查询',
    remark: receiptForm.remark,
    traces: []
  })

  showCreateDialog.value = false
}

function viewReceipt(item) {
  selectedReceipt.value = item
}

function confirmReceipt(item) {
  item.receiptStatus = 'received'
  item.logisticsStatus = '已签收'
  item.receiptTime = new Date().toISOString().slice(0, 10)
  item.lastQueryTime = new Date().toISOString().slice(0, 10)

  item.traces.unshift({
    id: Date.now(),
    time: new Date().toISOString().slice(0, 16).replace('T', ' '),
    desc: '系统已确认收货'
  })

  if (selectedReceipt.value) {
    selectedReceipt.value = item
  }

  alert(`快递单号【${item.expressNo}】已确认收货`)
}

function querySfLogistics(item) {
  /**
   * 正式接入顺丰 API 时，不建议前端直接请求顺丰。
   * 正确方式：
   *
   * 前端调用你的 Go 后端：
   * GET /api/sf/track?expressNo=SF202605100001
   *
   * Go 后端再去请求顺丰接口：
   * 1. 保存顺丰 appKey/appSecret
   * 2. 生成签名
   * 3. 请求顺丰物流轨迹接口
   * 4. 返回统一格式给前端
   */

  const today = new Date().toISOString().slice(0, 10)
  const now = new Date().toISOString().slice(0, 16).replace('T', ' ')

  item.lastQueryTime = today
  item.apiStatus = '查询成功'

  if (item.expressNo.endsWith('1')) {
    item.receiptStatus = 'received'
    item.logisticsStatus = '已签收'
    item.receiptTime = today
    item.traces.unshift({
      id: Date.now(),
      time: now,
      desc: '顺丰接口返回：快件已签收，系统自动确认收货'
    })
  } else if (item.expressNo.endsWith('3')) {
    item.receiptStatus = 'exception'
    item.logisticsStatus = '派送异常'
    item.traces.unshift({
      id: Date.now(),
      time: now,
      desc: '顺丰接口返回：派送异常，请联系快递员'
    })
  } else {
    item.receiptStatus = 'shipping'
    item.logisticsStatus = '运输中'
    item.traces.unshift({
      id: Date.now(),
      time: now,
      desc: '顺丰接口返回：快件运输中'
    })
  }

  if (selectedReceipt.value) {
    selectedReceipt.value = item
  }

  alert(`快递单号【${item.expressNo}】物流状态已更新`)
}

function deleteReceipt(item) {
  const ok = confirm(`确认删除快递单号【${item.expressNo}】吗？`)
  if (!ok) return

  receiptList.value = receiptList.value.filter(record => record.id !== item.id)
}

function exportReceiptRecords() {
  const header = [
    '快递单号',
    '快递公司',
    '收货人',
    '收货状态',
    '物流状态',
    '新增人',
    '新增时间',
    '确认收货时间',
    '最后查询时间',
    '备注'
  ]

  const rows = receiptList.value.map(item => [
    item.expressNo,
    item.expressCompany,
    item.receiver,
    getReceiptStatusText(item.receiptStatus),
    item.logisticsStatus,
    item.creator,
    item.createTime,
    item.receiptTime || '',
    item.lastQueryTime || '',
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
  link.download = '收货记录.csv'
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
  min-width: 1150px;
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

.file-tag {
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

  .version-table {
    min-width: 1150px;
  }

  .form-grid,
  .detail-card {
    grid-template-columns: 1fr;
  }
}

.header-actions {
  display: flex;
  gap: 12px;
}

.summary-panel {
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 14px;
  padding: 16px;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 14px;
  margin-bottom: 20px;
}

.summary-card {
  background: #020617;
  border: 1px solid #1e293b;
  border-radius: 12px;
  padding: 16px;
  display: flex;
  align-items: flex-start;
  gap: 14px;
}

.summary-card span {
  color: #94a3b8;
  font-size: 13px;
}

.summary-card strong {
  display: block;
  margin-top: 8px;
  color: #f8fafc;
  font-size: 28px;
  font-weight: 800;
}

.summary-card p {
  margin: 6px 0 0;
  color: #64748b;
  font-size: 12px;
}

.summary-card.green strong {
  color: #4ade80;
}

.summary-card.yellow strong {
  color: #fbbf24;
}

.summary-icon {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 15px;
  font-weight: 800;
}

.summary-icon.blue {
  background: #1d4ed833;
  color: #60a5fa;
}

.summary-icon.green {
  background: #16a34a33;
  color: #4ade80;
}

.summary-icon.yellow {
  background: #d9770633;
  color: #fbbf24;
}
.company-tag {
  display: inline-flex;
  padding: 4px 9px;
  border-radius: 999px;
  background: #1d4ed833;
  color: #60a5fa;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
}

.logistics-text {
  display: inline-block;
  max-width: 160px;
  color: #cbd5e1;
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.status-tag.shipping {
  background: #d9770633;
  color: #fbbf24;
}

.status-tag.received {
  background: #16a34a33;
  color: #4ade80;
}

.status-tag.exception {
  background: #dc262633;
  color: #f87171;
}

.empty-trace {
  color: #64748b;
  font-size: 13px;
  padding: 8px 0;
}

.trace-item {
  border-top: 1px solid #1e293b;
  padding-top: 12px;
  margin-top: 12px;
}

.trace-item strong {
  color: #f8fafc;
  font-size: 13px;
}

.trace-item p {
  margin: 6px 0 0;
  color: #cbd5e1;
  font-size: 13px;
  line-height: 1.6;
}

@media (max-width: 960px) {
  .summary-grid {
    grid-template-columns: 1fr;
  }
}


</style>