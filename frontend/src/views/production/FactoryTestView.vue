<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>板卡出厂测试管理</h1>
      </div>

      <button class="primary-btn" @click="openUploadDialog">
        上传出厂测试文档
      </button>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索产品型号 / MAC地址 / SN序列号 / 测试文档 / 上传人"
      />

      <select v-model="filters.productModel">
        <option value="">全部产品型号</option>
        <option
          v-for="model in productModelOptions"
          :key="model"
          :value="model"
        >
          {{ model }}
        </option>
      </select>

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

    <!-- 按产品型号归类展示 -->
    <div class="table-card">
      <div class="table-wrapper">
        <table class="model-table">
          <thead>
            <tr>
              <th>产品型号</th>
              <th>MAC数量</th>
              <th>测试文档</th>
              <th>上传人</th>
              <th>上传时间</th>
              <th>审核状态</th>
              <th>审核人</th>
              <th class="operation-col">操作</th>
            </tr>
          </thead>

          <tbody>
            <tr
              v-for="group in filteredModelGroupList"
              :key="group.groupKey"
              class="model-row"
            >
              <td>
                <button
                  class="model-name-btn"
                  @click="openModelMacDialog(group)"
                >
                  {{ group.productModel }}
                </button>
              </td>

              <td>
                <span class="count-tag">
                  {{ group.records.length }} 个 MAC
                </span>
              </td>

              <td>
                <span class="file-text" :title="group.fileName">
                  {{ group.fileName }}
                </span>
              </td>

              <td>
                <span class="normal-text" :title="group.uploader">
                  {{ group.uploader }}
                </span>
              </td>

              <td class="muted nowrap">
                {{ group.uploadTime }}
              </td>

              <td>
                <span class="status-tag" :class="group.auditStatus">
                  {{ getAuditStatusText(group.auditStatus) }}
                </span>
              </td>

              <td>
                {{ group.auditor || '-' }}
              </td>

              <td class="operation-col">
                <div class="action-group">
                  <button
                    class="text-btn blue"
                    @click="openModelMacDialog(group)"
                  >
                    查看MAC
                  </button>

                  <button
                    class="text-btn blue"
                    @click="downloadModelGroup(group)"
                  >
                    下载
                  </button>

                  <button
                    v-if="group.auditStatus === 'draft' || group.auditStatus === 'rejected'"
                    class="text-btn blue"
                    @click="submitModelGroup(group)"
                  >
                    提交
                  </button>

                  <button
                    v-if="group.auditStatus === 'submitted'"
                    class="text-btn green"
                    @click="auditModelGroup(group)"
                  >
                    审核
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="table-footer">
        共 {{ filteredFactoryTestList.length }} 条 MAC 出厂测试记录，按 {{ filteredModelGroupList.length }} 个产品型号归类展示。
      </div>
    </div>

    <!-- 上传出厂测试记录弹窗 -->
    <div v-if="showUploadDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>上传出厂测试文档</h3>
          <button @click="showUploadDialog = false">×</button>
        </div>

        <div class="upload-tip">
          <strong>填写说明：</strong>
          <p>
            出厂测试文档按产品型号上传。同一个产品型号下的 MAC 使用同一份测试文档。
          </p>
          <p>
            选择产品型号后，只展示该型号下生产烧录完毕的 MAC 地址。
          </p>
          <p>
            MAC 和 SN 序列号均来自生产烧录记录页面解析出的生产数据。
          </p>
        </div>

        <div class="form-grid">
          <label>
            产品型号
            <select v-model="uploadForm.productModel" @change="onProductModelChange">
              <option value="">请选择产品型号</option>
              <option
                v-for="model in productionProductModelOptions"
                :key="model"
                :value="model"
              >
                {{ model }}
              </option>
            </select>
          </label>

          <label>
            上传人
            <input
              v-model="currentUserName"
              disabled
            />
          </label>

          <label class="full-row">
            测试文档名称
            <input
              v-model="uploadForm.fileName"
              placeholder="选择文件后自动读取文件名称"
              disabled
            />
          </label>

          <label class="full-row">
            关联 MAC 地址
            <div class="mac-select-panel">
              <div v-if="!uploadForm.productModel" class="empty-mac">
                请先选择产品型号
              </div>

              <template v-else>
                <input
                  v-model="uploadForm.macKeyword"
                  class="mac-search-input"
                  placeholder="输入 MAC / 产品名称 / SN序列号 / 产品编码进行筛选"
                />

                <div class="mac-panel-header">
                  <span>显示烧录成功MAC，勾选表示不上传</span>
                  <strong>
                    需上传 {{ finalUploadMacList.length }} 个 / 已排除 {{ uploadForm.excludedMacs.length }} 个
                  </strong>
                </div>

                <div v-if="unpassedMacList.length === 0" class="empty-mac">
                  当前产品型号暂无烧录完毕的MAC。
                </div>

                <div v-else-if="availableMacList.length === 0" class="empty-mac">
                  没有匹配到 MAC，请换一个关键词。
                </div>

                <label
                  v-for="mac in availableMacList"
                  v-else
                  :key="mac.macAddress"
                  class="mac-check-item"
                >
                  <input
                    v-model="uploadForm.excludedMacs"
                    type="checkbox"
                    :value="mac.macAddress"
                  />
                  <span>{{ mac.macAddress }}</span>
                  <em>{{ mac.serialNumber || '-' }}</em>
                  <em>{{ mac.productName || '-' }}</em>
                    <b
                      v-if="uploadForm.excludedMacs.includes(mac.macAddress)"
                      class="exclude-tag"
                    >
                      不上传
                    </b>
                </label>
              </template>
            </div>
          </label>

          <label class="full-row">
            出厂测试文档
            <input
              type="file"
              accept=".doc,.docx,.xls,.xlsx,.pdf,.txt,.csv,.zip"
              @change="handleFileChange"
            />
          </label>

          <label class="full-row">
            出厂测试说明
            <textarea
              v-model="uploadForm.remark"
              placeholder="例如：该型号板卡出厂测试通过，测试项包括网络通信、音频、按键、接口、电源等"
            ></textarea>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showUploadDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="uploadFactoryTest">
            保存上传
          </button>
        </div>
      </div>
    </div>

    <!-- 产品型号 MAC 查看弹窗 -->
    <div v-if="selectedModelGroup" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>MAC 地址明细</h3>
          <button @click="selectedModelGroup = null">×</button>
        </div>

        <div class="model-summary-card">
          <div>
            <span>产品型号</span>
            <strong>{{ selectedModelGroup.productModel }}</strong>
          </div>

          <div>
            <span>测试文档</span>
            <strong>{{ selectedModelGroup.fileName }}</strong>
          </div>

          <div>
            <span>MAC 数量</span>
            <strong>{{ selectedModelMacList.length }} 个</strong>
          </div>

          <div>
            <span>数据来源</span>
            <strong>生产烧录记录</strong>
          </div>
        </div>

        <div class="mac-dialog-table-wrapper">
          <table class="mac-dialog-table">
            <thead>
              <tr>
                <th>序号</th>
                <th>产品名称</th>
                <th>MAC地址</th>
                <th>SN序列号</th>
                <th>产品编码</th>
              </tr>
            </thead>

            <tbody>
              <tr
                v-for="(item, index) in selectedModelMacList"
                :key="item.macAddress"
              >
                <td>{{ index + 1 }}</td>

                <td>{{ item.productName || '-' }}</td>

                <td>
                  <span class="mac-tag" :title="item.macAddress">
                    {{ item.macAddress }}
                  </span>
                </td>

                <td>
                  <span class="sn-tag" :title="item.serialNumber">
                    {{ item.serialNumber || '-' }}
                  </span>
                </td>

                <td>
                  <span class="code-text" :title="item.productCode">
                    {{ item.productCode || '-' }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div v-if="selectedModelMacList.length === 0" class="empty-dialog-data">
          当前型号下暂无 MAC 明细。
        </div>
      </div>
    </div>

    <!-- 查看详情 / 审核弹窗 -->
    <div v-if="selectedFactoryTest" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>出厂测试记录详情</h3>
          <button @click="selectedFactoryTest = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>产品型号</span>
            <strong>{{ selectedFactoryTest.productModel }}</strong>
          </div>

          <div>
            <span>MAC地址</span>
            <strong>{{ selectedFactoryTest.macAddress }}</strong>
          </div>

          <div>
            <span>SN序列号</span>
            <strong>{{ getBurnInfoByMac(selectedFactoryTest.macAddress).serialNumber || '-' }}</strong>
          </div>

          <div>
            <span>测试文档</span>
            <strong>{{ selectedFactoryTest.fileName }}</strong>
          </div>

          <div>
            <span>上传人</span>
            <strong>{{ selectedFactoryTest.uploader }}</strong>
          </div>

          <div>
            <span>上传时间</span>
            <strong>{{ selectedFactoryTest.uploadTime }}</strong>
          </div>

          <div>
            <span>审核状态</span>
            <strong>{{ getAuditStatusText(selectedFactoryTest.auditStatus) }}</strong>
          </div>

          <div>
            <span>审核人</span>
            <strong>{{ selectedFactoryTest.auditor || '-' }}</strong>
          </div>

          <div>
            <span>审核时间</span>
            <strong>{{ selectedFactoryTest.auditTime || '-' }}</strong>
          </div>

          <div>
            <span>文件查看</span>
            <button class="inline-link" @click="openFactoryTestFile(selectedFactoryTest)">
              点开查看文件
            </button>
          </div>
        </div>

        <div class="remark-card">
          <span>出厂测试说明</span>
          <p>{{ selectedFactoryTest.remark || '暂无说明' }}</p>
        </div>

        <div class="dialog-footer">
          <button
            v-if="selectedFactoryTest.auditStatus === 'submitted'"
            class="green-btn"
            @click="approveFactoryTest(selectedFactoryTest)"
          >
            审核通过
          </button>

          <button
            v-if="selectedFactoryTest.auditStatus === 'submitted'"
            class="red-btn"
            @click="rejectFactoryTest(selectedFactoryTest)"
          >
            审核驳回
          </button>

          <button class="reset-btn" @click="downloadFactoryTest(selectedFactoryTest)">
            下载文件
          </button>

          <button class="primary-btn" @click="selectedFactoryTest = null">
            关闭
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'

const currentUserName = ref(
  localStorage.getItem('username') ||
  localStorage.getItem('accountName') ||
  localStorage.getItem('realName') ||
  '当前用户'
)

const filters = reactive({
  keyword: '',
  productModel: '',
  auditStatus: ''
})

const showUploadDialog = ref(false)
const selectedFactoryTest = ref(null)
const selectedModelGroup = ref(null)

const uploadForm = reactive({
  productModel: '',
  excludedMacs: [],
  macKeyword: '',
  fileName: '',
  file: null,
  fileUrl: '',
  remark: ''
})

/**
 * 模拟“生产烧录记录页面”解析 Excel 后得到的数据。
 * 实际项目中这里应由后端接口返回，例如：
 * GET /api/production/burn-records
 *
 * factoryTestStatus:
 * passed：已通过，不再允许勾选
 * failed：未通过，允许勾选上传出厂测试文档
 */
const productionBurnRecordList = ref([
  {
    id: 1,
    productName: '控制盒',
    productModel: 'DACU-SIP-Porto',
    productCode: 'M00004665176',
    serialNumber: 'ZDAPO2300001',
    macAddress: '68:69:2E:DD:17:27',
    hardwareVersion: '-',
    softwareVersion: '-',
    factoryTestStatus: 'passed'
  },
  {
    id: 2,
    productName: '控制盒',
    productModel: 'DACU-SIP-Porto',
    productCode: 'M00004665176',
    serialNumber: 'ZDAPO2300002',
    macAddress: '68:69:2E:DD:1A:E6',
    hardwareVersion: '-',
    softwareVersion: '-',
    factoryTestStatus: 'passed'
  },
  {
    id: 3,
    productName: '控制盒',
    productModel: 'DACU-SIP-Porto',
    productCode: 'M00004665176',
    serialNumber: 'ZDAPO2300003',
    macAddress: '68:69:2E:DD:1A:E7',
    hardwareVersion: '-',
    softwareVersion: '-',
    factoryTestStatus: 'failed'
  },
  {
    id: 4,
    productName: '乘客报警器',
    productModel: 'PECU-BOG-Metro-PCBA',
    productCode: 'M00004665176',
    serialNumber: 'ZPEBOG230000X',
    macAddress: '68:69:2E:DD:22:01',
    hardwareVersion: '-',
    softwareVersion: '-',
    factoryTestStatus: 'failed'
  },
  {
    id: 5,
    productName: '编码板',
    productModel: 'ACSU-BOG-Encode-PCBA',
    productCode: 'M00004665172',
    serialNumber: 'ZENBOG230000X',
    macAddress: '68:69:2E:DD:33:01',
    hardwareVersion: '-',
    softwareVersion: '-',
    factoryTestStatus: 'failed'
  }
])

const factoryTestList = ref([
  {
    id: 1,
    productModel: 'DACU-SIP-Porto',
    macAddress: '68:69:2E:DD:17:27',
    recordName: 'DACU-SIP-Porto_出厂测试记录_V1.0.xlsx',
    fileName: 'DACU-SIP-Porto_出厂测试记录_V1.0.xlsx',
    fileUrl: '',
    uploader: '生产人员',
    uploadTime: '2026-05-10',
    auditStatus: 'approved',
    auditor: '领导',
    auditTime: '2026-05-11',
    remark: '该型号板卡出厂测试通过。'
  },
  {
    id: 2,
    productModel: 'DACU-SIP-Porto',
    macAddress: '68:69:2E:DD:1A:E6',
    recordName: 'DACU-SIP-Porto_出厂测试记录_V1.0.xlsx',
    fileName: 'DACU-SIP-Porto_出厂测试记录_V1.0.xlsx',
    fileUrl: '',
    uploader: '生产人员',
    uploadTime: '2026-05-10',
    auditStatus: 'approved',
    auditor: '领导',
    auditTime: '2026-05-11',
    remark: '同型号板卡使用同一份测试文档。'
  },
  {
    id: 3,
    productModel: 'PECU-BOG-Metro-PCBA',
    macAddress: '68:69:2E:DD:22:01',
    recordName: 'PECU-BOG-Metro-PCBA_出厂测试记录_V1.0.xlsx',
    fileName: 'PECU-BOG-Metro-PCBA_出厂测试记录_V1.0.xlsx',
    fileUrl: '',
    uploader: '生产人员',
    uploadTime: '2026-05-16',
    auditStatus: 'submitted',
    auditor: '',
    auditTime: '',
    remark: '待领导审核确认。'
  }
])

const productionProductModelOptions = computed(() => {
  const models = productionBurnRecordList.value
    .map(item => item.productModel)
    .filter(Boolean)

  return [...new Set(models)]
})

const productModelOptions = computed(() => {
  const models = factoryTestList.value
    .map(item => item.productModel)
    .filter(Boolean)

  return [...new Set(models)]
})

const unpassedMacList = computed(() => {
  if (!uploadForm.productModel) {
    return []
  }

  return productionBurnRecordList.value.filter(item => {
    return (
      item.productModel === uploadForm.productModel &&
      item.macAddress &&
      item.factoryTestStatus === 'failed'
    )
  })
})

const availableMacList = computed(() => {
  const keyword = uploadForm.macKeyword.trim().toLowerCase()

  if (!keyword) {
    return unpassedMacList.value
  }

  return unpassedMacList.value.filter(item => {
    return (
      item.macAddress.toLowerCase().includes(keyword) ||
      item.productName.toLowerCase().includes(keyword) ||
      item.productCode.toLowerCase().includes(keyword) ||
      item.serialNumber.toLowerCase().includes(keyword)
    )
  })
})

const finalUploadMacList = computed(() => {
  return unpassedMacList.value.filter(item => {
    return !uploadForm.excludedMacs.includes(item.macAddress)
  })
})

const filteredFactoryTestList = computed(() => {
  return factoryTestList.value.filter(item => {
    const keyword = filters.keyword.trim()
    const burnInfo = getBurnInfoByMac(item.macAddress)

    const keywordMatch =
      !keyword ||
      item.productModel.includes(keyword) ||
      item.macAddress.includes(keyword) ||
      item.recordName.includes(keyword) ||
      item.uploader.includes(keyword) ||
      item.fileName.includes(keyword) ||
      burnInfo.serialNumber.includes(keyword) ||
      burnInfo.productCode.includes(keyword) ||
      burnInfo.productName.includes(keyword)

    const modelMatch =
      !filters.productModel || item.productModel === filters.productModel

    const auditStatusMatch =
      !filters.auditStatus || item.auditStatus === filters.auditStatus

    return keywordMatch && modelMatch && auditStatusMatch
  })
})

const filteredModelGroupList = computed(() => {
  const map = new Map()

  filteredFactoryTestList.value.forEach(item => {
    const groupKey = [
      item.productModel || '未填写型号',
      item.fileName || '未上传文档',
      item.uploader || '未知上传人',
      item.uploadTime || '未知上传时间',
      item.auditStatus || 'draft'
    ].join('__')

    if (!map.has(groupKey)) {
      map.set(groupKey, {
        groupKey,
        productModel: item.productModel || '未填写型号',
        records: []
      })
    }

    map.get(groupKey).records.push(item)
  })

  return Array.from(map.values()).map(group => {
    const sortedRecords = [...group.records].sort((a, b) => {
      return String(a.macAddress).localeCompare(String(b.macAddress))
    })

    const first = sortedRecords[0] || {}

    return {
      groupKey: group.groupKey,
      productModel: group.productModel,
      records: sortedRecords,
      fileName: first.fileName || '-',
      fileUrl: first.fileUrl || '',
      uploader: first.uploader || '-',
      uploadTime: first.uploadTime || '-',
      auditStatus: first.auditStatus || 'draft',
      auditor: first.auditor || '',
      auditTime: first.auditTime || ''
    }
  })
})

const selectedModelMacList = computed(() => {
  if (!selectedModelGroup.value) {
    return []
  }

  return selectedModelGroup.value.records.map(record => {
    const burnInfo = getBurnInfoByMac(record.macAddress)

    return {
      macAddress: record.macAddress,
      serialNumber: burnInfo.serialNumber,
      productName: burnInfo.productName,
      productCode: burnInfo.productCode
    }
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

function getBurnInfoByMac(macAddress) {
  return (
    productionBurnRecordList.value.find(item => item.macAddress === macAddress) ||
    {
      productName: '',
      productModel: '',
      productCode: '',
      serialNumber: '',
      macAddress: '',
      hardwareVersion: '',
      softwareVersion: ''
    }
  )
}

function resetFilters() {
  filters.keyword = ''
  filters.productModel = ''
  filters.auditStatus = ''
}

function openModelMacDialog(group) {
  selectedModelGroup.value = group
}

function openUploadDialog() {
  uploadForm.productModel = ''
  uploadForm.excludedMacs = []
  uploadForm.macKeyword = ''
  uploadForm.fileName = ''
  uploadForm.file = null
  uploadForm.fileUrl = ''
  uploadForm.remark = ''

  showUploadDialog.value = true
}

function onProductModelChange() {
  uploadForm.excludedMacs = []
  uploadForm.macKeyword = ''
}

function handleFileChange(event) {
  const file = event.target.files[0]
  if (!file) return

  uploadForm.file = file
  uploadForm.fileName = file.name
  uploadForm.fileUrl = URL.createObjectURL(file)
}

function uploadFactoryTest() {
  if (!uploadForm.productModel) {
    alert('请选择产品型号')
    return
  }

  if (finalUploadMacList.value.length === 0) {
    alert('当前没有需要上传的 MAC，请至少保留一个未排除的 MAC')
    return
  }

  if (!uploadForm.file) {
    alert('请上传出厂测试文档')
    return
  }

  const today = new Date().toISOString().slice(0, 10)
  const uploader = currentUserName.value
  const recordName = uploadForm.fileName

  const newRecords = finalUploadMacList.value.map((macItem, index) => {
    return {
      id: Date.now() + index,
      productModel: uploadForm.productModel,
      macAddress: macItem.macAddress,
      recordName,
      fileName: uploadForm.fileName,
      fileUrl: uploadForm.fileUrl,
      uploader,
      uploadTime: today,
      auditStatus: 'draft',
      auditor: '',
      auditTime: '',
      remark: uploadForm.remark
    }
  })
  factoryTestList.value.unshift(...newRecords)

  showUploadDialog.value = false
}

function openFactoryTestFile(item) {
  if (!item.fileUrl) {
    alert('当前是模拟数据，暂无可直接打开的原始文件')
    return
  }

  window.open(item.fileUrl, '_blank')
}

function downloadFactoryTest(item) {
  if (!item.fileUrl) {
    alert('当前是模拟数据，暂无可下载的原始文件')
    return
  }

  const link = document.createElement('a')
  link.href = item.fileUrl
  link.download = item.fileName || '出厂测试记录文件'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

function downloadModelGroup(group) {
  const first = group.records[0]

  if (!first) {
    alert('当前产品型号暂无测试记录')
    return
  }

  downloadFactoryTest(first)
}

function submitModelGroup(group) {
  group.records.forEach(item => {
    item.auditStatus = 'submitted'
    item.auditor = ''
    item.auditTime = ''
  })

  alert(`产品型号【${group.productModel}】的出厂测试文档已提交领导审核`)
}

function auditModelGroup(group) {
  const first = group.records[0]

  if (!first) return

  selectedFactoryTest.value = first
}

function approveFactoryTest(item) {
  updateSameGroupAuditStatus(item, 'approved')

  selectedFactoryTest.value = null
  alert(`出厂测试文档【${item.fileName}】审核通过`)
}

function rejectFactoryTest(item) {
  updateSameGroupAuditStatus(item, 'rejected')

  selectedFactoryTest.value = null
  alert(`出厂测试文档【${item.fileName}】已驳回`)
}

function updateSameGroupAuditStatus(target, status) {
  factoryTestList.value.forEach(item => {
    const sameGroup =
      item.productModel === target.productModel &&
      item.fileName === target.fileName &&
      item.uploader === target.uploader &&
      item.uploadTime === target.uploadTime

    if (sameGroup) {
      item.auditStatus = status
      item.auditor = '领导'
      item.auditTime = new Date().toISOString().slice(0, 10)
    }
  })
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
  grid-template-columns: 1.4fr 220px 180px 90px 90px;
  gap: 12px;
  margin-bottom: 20px;
}

.filter-card input,
.filter-card select,
.form-grid input,
.form-grid select,
.form-grid textarea,
.mac-search-input {
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
.form-grid select,
.mac-search-input {
  height: 36px;
}

.form-grid textarea {
  min-height: 90px;
  padding: 10px 12px;
  resize: vertical;
}

.filter-card input::placeholder,
.form-grid input::placeholder,
.form-grid textarea::placeholder,
.mac-search-input::placeholder {
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

.table-wrapper,
.mac-dialog-table-wrapper {
  width: 100%;
  overflow-x: auto;
  overflow-y: hidden;
}

.table-wrapper::-webkit-scrollbar,
.mac-dialog-table-wrapper::-webkit-scrollbar {
  height: 10px;
}

.table-wrapper::-webkit-scrollbar-track,
.mac-dialog-table-wrapper::-webkit-scrollbar-track {
  background: #020617;
  border-radius: 999px;
}

.table-wrapper::-webkit-scrollbar-thumb,
.mac-dialog-table-wrapper::-webkit-scrollbar-thumb {
  background: #334155;
  border-radius: 999px;
  border: 2px solid #020617;
}

.table-wrapper::-webkit-scrollbar-thumb:hover,
.mac-dialog-table-wrapper::-webkit-scrollbar-thumb:hover {
  background: #475569;
}

.table-wrapper::-webkit-scrollbar-button,
.mac-dialog-table-wrapper::-webkit-scrollbar-button {
  display: none;
}

.table-wrapper,
.mac-dialog-table-wrapper {
  scrollbar-width: thin;
  scrollbar-color: #334155 #020617;
}

.model-table {
  width: 100%;
  min-width: 1280px;
  border-collapse: collapse;
  table-layout: fixed;
}

.model-table thead {
  background: #020617;
}

.model-table th,
.model-table td {
  box-sizing: border-box;
  white-space: nowrap;
}

.model-table th {
  padding: 14px 16px;
  color: #94a3b8;
  font-size: 12px;
  font-weight: 600;
  text-align: left;
  border-bottom: 1px solid #1e293b;
}

.model-table td {
  padding: 15px 16px;
  font-size: 13px;
  color: #e2e8f0;
  border-bottom: 1px solid #1e293b;
  vertical-align: middle;
  overflow: hidden;
}

.model-table th:nth-child(1),
.model-table td:nth-child(1) {
  width: 220px;
}

.model-table th:nth-child(2),
.model-table td:nth-child(2) {
  width: 120px;
}

.model-table th:nth-child(3),
.model-table td:nth-child(3) {
  width: 280px;
}

.model-table th:nth-child(4),
.model-table td:nth-child(4) {
  width: 120px;
}

.model-table th:nth-child(5),
.model-table td:nth-child(5) {
  width: 140px;
}

.model-table th:nth-child(6),
.model-table td:nth-child(6) {
  width: 130px;
}

.model-table th:nth-child(7),
.model-table td:nth-child(7) {
  width: 110px;
}

.model-table th:nth-child(8),
.model-table td:nth-child(8) {
  width: 160px;
}

.model-row {
  background: #0f172a;
}

.model-row:hover {
  background: #1e293b80;
}

.model-name-btn {
  display: inline-block;
  max-width: 190px;
  border: none;
  background: transparent;
  color: #60a5fa;
  font-size: 14px;
  font-weight: 800;
  cursor: pointer;
  padding: 0;
  text-align: left;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;
}

.model-name-btn:hover {
  color: #93c5fd;
  text-decoration: underline;
}

.count-tag {
  display: inline-block;
  padding: 4px 9px;
  border-radius: 999px;
  background: #33415566;
  color: #cbd5e1;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
}

.mac-tag,
.sn-tag,
.code-text,
.file-text,
.normal-text {
  display: inline-block;
  max-width: 100%;
  color: #cbd5e1;
  font-size: 12px;
  font-family: Consolas, Monaco, monospace;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;
}

.mac-tag,
.sn-tag {
  padding: 3px 8px;
  border-radius: 999px;
  background: #33415566;
}

.normal-text {
  font-family: inherit;
  font-size: 13px;
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

.nowrap {
  white-space: nowrap !important;
}

.operation-col {
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

.upload-tip {
  margin: 18px 20px 0;
  padding: 12px;
  background: #020617;
  border: 1px solid #1e293b;
  border-radius: 10px;
}

.upload-tip strong {
  color: #f8fafc;
  font-size: 13px;
}

.upload-tip p {
  margin: 6px 0 0;
  color: #94a3b8;
  font-size: 12px;
  line-height: 1.6;
}

.mac-select-panel {
  min-height: 120px;
  max-height: 260px;
  overflow-y: auto;
  border: 1px solid #334155;
  border-radius: 10px;
  background: #020617;
  padding: 10px;
}

.mac-select-panel::-webkit-scrollbar {
  width: 8px;
}

.mac-select-panel::-webkit-scrollbar-track {
  background: #020617;
}

.mac-select-panel::-webkit-scrollbar-thumb {
  background: #334155;
  border-radius: 999px;
}

.mac-search-input {
  width: 100%;
  box-sizing: border-box;
  margin-bottom: 10px;
}

.mac-panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  color: #64748b;
  font-size: 12px;
}

.mac-panel-header strong {
  color: #5eead4;
  font-size: 12px;
}

.exclude-tag {
  display: inline-flex;
  justify-content: center;
  align-items: center;
  height: 22px;
  padding: 0 8px;
  border-radius: 999px;
  background: #7f1d1d55;
  color: #fca5a5;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
}

.empty-mac,
.empty-dialog-data {
  color: #64748b;
  font-size: 13px;
  padding: 8px 0;
}

.empty-dialog-data {
  padding: 0 20px 20px;
}

.mac-check-item {
  display: grid !important;
  grid-template-columns: 18px 180px 130px 1fr 70px;
  align-items: center;
  gap: 8px !important;
  padding: 8px;
  border-bottom: 1px solid #1e293b;
  color: #cbd5e1 !important;
}

.mac-check-item:last-child {
  border-bottom: none;
}

.mac-check-item input {
  width: 14px;
  height: 14px;
}

.mac-check-item span {
  font-family: Consolas, Monaco, monospace;
  color: #e2e8f0;
  font-size: 12px;
}

.mac-check-item em {
  color: #64748b;
  font-size: 12px;
  font-style: normal;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* MAC 弹窗 */
.model-summary-card {
  padding: 20px;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 14px;
}

.model-summary-card div {
  background: #020617;
  border: 1px solid #1e293b;
  border-radius: 10px;
  padding: 12px;
}

.model-summary-card span {
  display: block;
  color: #64748b;
  font-size: 12px;
  margin-bottom: 6px;
}

.model-summary-card strong {
  color: #f8fafc;
  font-size: 14px;
  word-break: break-all;
}

.mac-dialog-table-wrapper {
  padding: 0 20px 20px;
  box-sizing: border-box;
}

.mac-dialog-table {
  width: 100%;
  min-width: 760px;
  border-collapse: collapse;
  table-layout: fixed;
  border: 1px solid #1e293b;
  border-radius: 10px;
  overflow: hidden;
}

.mac-dialog-table thead {
  background: #020617;
}

.mac-dialog-table th,
.mac-dialog-table td {
  padding: 13px 14px;
  border-bottom: 1px solid #1e293b;
  color: #e2e8f0;
  font-size: 13px;
  text-align: left;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.mac-dialog-table th {
  color: #94a3b8;
  font-size: 12px;
  font-weight: 600;
}

.mac-dialog-table th:nth-child(1),
.mac-dialog-table td:nth-child(1) {
  width: 70px;
}

.mac-dialog-table th:nth-child(2),
.mac-dialog-table td:nth-child(2) {
  width: 160px;
}

.mac-dialog-table th:nth-child(3),
.mac-dialog-table td:nth-child(3) {
  width: 190px;
}

.mac-dialog-table th:nth-child(4),
.mac-dialog-table td:nth-child(4) {
  width: 180px;
}

.mac-dialog-table th:nth-child(5),
.mac-dialog-table td:nth-child(5) {
  width: 180px;
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
  max-height: 92vh;
  overflow-y: auto;
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

@media (max-width: 960px) {
  .filter-card {
    grid-template-columns: 1fr;
  }

  .model-table {
    min-width: 1280px;
  }

  .form-grid,
  .detail-card,
  .model-summary-card {
    grid-template-columns: 1fr;
  }

  .mac-check-item {
    grid-template-columns: 18px 1fr;
  }

  .mac-check-item em {
    grid-column: 2;
  }

  .mac-dialog-table {
    min-width: 760px;
  }
}
</style>