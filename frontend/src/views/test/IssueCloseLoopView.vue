<template>
  <div class="page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div>
        <h1>问题闭环管理</h1>
      </div>

      <div class="header-actions">
        <button class="secondary-btn" @click="exportIssueRecords">
          导出问题记录
        </button>

        <button class="primary-btn" @click="openCreateDialog">
          新增问题
        </button>
      </div>
    </div>

    <!-- 查询条件 -->
    <div class="filter-card">
      <input
        v-model="filters.keyword"
        placeholder="搜索项目名称 / 问题标题 / 负责人 / 提出人"
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

      <select v-model="filters.issueSource">
        <option value="">全部来源</option>
        <option value="研发">研发</option>
        <option value="测试">测试</option>
        <option value="生产">生产</option>
        <option value="售后">售后</option>
      </select>

      <select v-model="filters.closeStatus">
        <option value="">全部状态</option>
        <option value="open">未关闭</option>
        <option value="closed">已关闭</option>
      </select>

      <button class="query-btn">查询</button>
      <button class="reset-btn" @click="resetFilters">重置</button>
    </div>

    <!-- 数据表格 -->
    <div class="table-card">
      <table>
        <thead>
          <tr>
            <th>问题标题</th>
            <th>绑定项目</th>
            <th>终端类型</th>
            <th>问题来源</th>
            <th>严重等级</th>
            <th>负责人</th>
            <th>关闭状态</th>
            <th>提出时间</th>
            <th class="operation-col">操作</th>
          </tr>
        </thead>

        <tbody>
          <tr v-for="item in filteredIssueList" :key="item.id">
            <td>
              <button class="issue-link" @click="viewIssue(item)">
                {{ item.issueTitle }}
              </button>
              <div class="issue-desc">{{ item.issueDesc }}</div>
            </td>

            <td>
              <span class="project-tag">{{ item.projectName }}</span>
            </td>

            <td>
              <span class="device-tag">{{ item.deviceType }}</span>
            </td>

            <td>
              <span class="source-tag" :class="getSourceClass(item.issueSource)">
                {{ item.issueSource }}
              </span>
            </td>

            <td>
              <span class="level-tag" :class="getLevelClass(item.level)">
                {{ item.level }}
              </span>
            </td>

            <td>{{ item.owner }}</td>

            <td>
              <span class="close-tag" :class="item.closeStatus">
                {{ getCloseStatusText(item.closeStatus) }}
              </span>
            </td>

            <td class="muted">{{ item.createTime }}</td>

            <td class="operation-col">
              <div class="action-group">
                <button class="text-btn" @click="viewIssue(item)">
                  查看
                </button>

                <button
                  v-if="item.closeStatus === 'open'"
                  class="text-btn blue"
                  @click="openEditDialog(item)"
                >
                  修改
                </button>

                <button
                  v-if="item.closeStatus === 'open'"
                  class="text-btn green"
                  @click="openReplyDialog(item)"
                >
                  回复
                </button>

                <button
                  v-if="item.closeStatus === 'open'"
                  class="text-btn red"
                  @click="closeIssue(item)"
                >
                  关闭
                </button>

                <button
                  v-if="item.closeStatus === 'closed'"
                  class="text-btn yellow"
                  @click="reopenIssue(item)"
                >
                  重新打开
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <div class="table-footer">
        共 {{ filteredIssueList.length }} 条问题记录
      </div>
    </div>

    <!-- 新增 / 修改问题弹窗 -->
    <div v-if="showEditDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>{{ editMode === 'create' ? '新增问题' : '修改问题' }}</h3>
          <button @click="showEditDialog = false">×</button>
        </div>

        <div class="form-grid">
          <label>
            绑定项目
            <select v-model="issueForm.projectName">
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
            <select v-model="issueForm.deviceType">
              <option value="">请选择终端</option>
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
            问题来源
            <select v-model="issueForm.issueSource">
              <option value="">请选择来源</option>
              <option value="研发">研发</option>
              <option value="测试">测试</option>
              <option value="生产">生产</option>
              <option value="售后">售后</option>
            </select>
          </label>

          <label>
            严重等级
            <select v-model="issueForm.level">
              <option value="紧急">紧急</option>
              <option value="高">高</option>
              <option value="中">中</option>
              <option value="低">低</option>
            </select>
          </label>

          <label>
            问题标题
            <input
              v-model="issueForm.issueTitle"
              placeholder="例如：广播控制盒上电后无声音输出"
            />
          </label>

          <label>
            负责人
            <input
              v-model="issueForm.owner"
              placeholder="请输入问题负责人"
            />
          </label>

          <label>
            提出人
            <input
              v-model="issueForm.creator"
              placeholder="请输入问题提出人"
            />
          </label>

          <label>
            计划关闭时间
            <input
              v-model="issueForm.planCloseTime"
              type="date"
            />
          </label>

          <label class="full-row">
            问题描述
            <textarea
              v-model="issueForm.issueDesc"
              placeholder="请输入问题现象、复现步骤、影响范围、现场条件等"
            ></textarea>
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showEditDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="saveIssue">
            保存
          </button>
        </div>
      </div>
    </div>

    <!-- 回复问题弹窗 -->
    <div v-if="showReplyDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>回复问题</h3>
          <button @click="showReplyDialog = false">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>问题标题</span>
            <strong>{{ currentReplyIssue?.issueTitle }}</strong>
          </div>

          <div>
            <span>绑定项目</span>
            <strong>{{ currentReplyIssue?.projectName }}</strong>
          </div>

          <div>
            <span>终端类型</span>
            <strong>{{ currentReplyIssue?.deviceType }}</strong>
          </div>

          <div>
            <span>负责人</span>
            <strong>{{ currentReplyIssue?.owner }}</strong>
          </div>
        </div>

        <div class="form-grid reply-form">
          <label class="full-row">
            回复内容
            <textarea
              v-model="replyForm.content"
              placeholder="请输入问题分析、处理方案、验证结果或下一步计划"
            ></textarea>
          </label>

          <label>
            回复人
            <input
              v-model="replyForm.replyUser"
              placeholder="请输入回复人"
            />
          </label>
        </div>

        <div class="dialog-footer">
          <button class="reset-btn" @click="showReplyDialog = false">
            取消
          </button>

          <button class="primary-btn" @click="saveReply">
            保存回复
          </button>
        </div>
      </div>
    </div>

    <!-- 查看问题详情弹窗 -->
    <div v-if="selectedIssue" class="dialog-mask">
      <div class="dialog large-dialog">
        <div class="dialog-header">
          <h3>问题详情</h3>
          <button @click="selectedIssue = null">×</button>
        </div>

        <div class="detail-card">
          <div>
            <span>问题标题</span>
            <strong>{{ selectedIssue.issueTitle }}</strong>
          </div>

          <div>
            <span>绑定项目</span>
            <strong>{{ selectedIssue.projectName }}</strong>
          </div>

          <div>
            <span>终端类型</span>
            <strong>{{ selectedIssue.deviceType }}</strong>
          </div>

          <div>
            <span>问题来源</span>
            <strong>{{ selectedIssue.issueSource }}</strong>
          </div>

          <div>
            <span>严重等级</span>
            <strong>{{ selectedIssue.level }}</strong>
          </div>

          <div>
            <span>负责人</span>
            <strong>{{ selectedIssue.owner }}</strong>
          </div>

          <div>
            <span>提出人</span>
            <strong>{{ selectedIssue.creator }}</strong>
          </div>

          <div>
            <span>提出时间</span>
            <strong>{{ selectedIssue.createTime }}</strong>
          </div>

          <div>
            <span>计划关闭时间</span>
            <strong>{{ selectedIssue.planCloseTime || '-' }}</strong>
          </div>

          <div>
            <span>实际关闭时间</span>
            <strong>{{ selectedIssue.realCloseTime || '-' }}</strong>
          </div>

          <div>
            <span>关闭状态</span>
            <strong>{{ getCloseStatusText(selectedIssue.closeStatus) }}</strong>
          </div>

          <div>
            <span>关闭人</span>
            <strong>{{ selectedIssue.closeUser || '-' }}</strong>
          </div>
        </div>

        <div class="remark-card">
          <span>问题描述</span>
          <p>{{ selectedIssue.issueDesc || '暂无描述' }}</p>
        </div>

        <div class="reply-card">
          <span>问题回复记录</span>

          <div
            v-if="selectedIssue.replies.length === 0"
            class="empty-reply"
          >
            暂无回复记录
          </div>

          <div
            v-for="reply in selectedIssue.replies"
            :key="reply.id"
            class="reply-item"
          >
            <div class="reply-header">
              <strong>{{ reply.replyUser }}</strong>
              <span>{{ reply.replyTime }}</span>
            </div>
            <p>{{ reply.content }}</p>
          </div>
        </div>

        <div class="dialog-footer">
          <button
            v-if="selectedIssue.closeStatus === 'open'"
            class="reset-btn"
            @click="openReplyDialog(selectedIssue)"
          >
            回复问题
          </button>

          <button
            v-if="selectedIssue.closeStatus === 'open'"
            class="red-btn"
            @click="closeIssue(selectedIssue)"
          >
            关闭问题
          </button>

          <button
            v-if="selectedIssue.closeStatus === 'closed'"
            class="green-btn"
            @click="reopenIssue(selectedIssue)"
          >
            重新打开
          </button>

          <button class="primary-btn" @click="selectedIssue = null">
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
  issueSource: '',
  closeStatus: ''
})

const showEditDialog = ref(false)
const showReplyDialog = ref(false)

const selectedIssue = ref(null)
const currentEditIssue = ref(null)
const currentReplyIssue = ref(null)
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

const issueForm = reactive({
  projectName: '',
  deviceType: '',
  issueSource: '',
  level: '中',
  issueTitle: '',
  owner: '',
  creator: '',
  planCloseTime: '',
  issueDesc: ''
})

const replyForm = reactive({
  replyUser: '',
  content: ''
})

const issueList = ref([
  {
    id: 1,
    projectName: '香港屯马项目',
    deviceType: '广播控制盒',
    issueSource: '测试',
    level: '高',
    issueTitle: '广播控制盒上电后偶发无声音输出',
    owner: '卢进',
    creator: '寸诗睿',
    createTime: '2026-05-10',
    planCloseTime: '2026-05-15',
    realCloseTime: '',
    closeStatus: 'open',
    closeUser: '',
    issueDesc: '内部测试时发现广播控制盒上电后偶发无声音输出，重启后恢复，需要研发分析音频初始化流程。',
    replies: [
      {
        id: 101,
        replyUser: '卢进',
        replyTime: '2026-05-11',
        content: '已初步定位为音频服务启动时序问题，后续增加初始化完成检测。'
      }
    ]
  },
  {
    id: 2,
    projectName: '波尔图二期项目',
    deviceType: '乘客报警器',
    issueSource: '生产',
    level: '中',
    issueTitle: '生产烧录后部分终端无法注册SIP',
    owner: '王宇',
    creator: '生产人员',
    createTime: '2026-05-16',
    planCloseTime: '2026-05-20',
    realCloseTime: '',
    closeStatus: 'open',
    closeUser: '',
    issueDesc: '生产烧录后发现部分乘客报警器无法注册SIP账号，需要确认配置文件和烧录流程。',
    replies: []
  },
  {
    id: 3,
    projectName: '阿根廷有轨项目',
    deviceType: '客室解码板',
    issueSource: '研发',
    level: '低',
    issueTitle: '客室解码板日志打印过多',
    owner: '郑宇',
    creator: '研发人员',
    createTime: '2026-05-18',
    planCloseTime: '2026-05-22',
    realCloseTime: '2026-05-19',
    closeStatus: 'closed',
    closeUser: '郑宇',
    issueDesc: '调试版本日志打印过多，影响现场查看关键错误日志。',
    replies: [
      {
        id: 102,
        replyUser: '郑宇',
        replyTime: '2026-05-19',
        content: '已降低日志等级，保留关键错误日志。'
      }
    ]
  },
  {
    id: 4,
    projectName: '波哥大有轨项目',
    deviceType: '编码板',
    issueSource: '售后',
    level: '紧急',
    issueTitle: '现场反馈编码板音频中断',
    owner: '卢进',
    creator: '售后人员',
    createTime: '2026-05-20',
    planCloseTime: '2026-05-21',
    realCloseTime: '',
    closeStatus: 'open',
    closeUser: '',
    issueDesc: '现场反馈编码板运行一段时间后音频中断，需要排查网络、buffer和音视频同步机制。',
    replies: [
      {
        id: 103,
        replyUser: '卢进',
        replyTime: '2026-05-20',
        content: '需要现场提供码流、日志、RTP时间戳和中断时刻抓包。'
      }
    ]
  }
])

const filteredIssueList = computed(() => {
  return issueList.value.filter(item => {
    const keywordMatch =
      !filters.keyword ||
      item.projectName.includes(filters.keyword) ||
      item.deviceType.includes(filters.keyword) ||
      item.issueTitle.includes(filters.keyword) ||
      item.owner.includes(filters.keyword) ||
      item.creator.includes(filters.keyword) ||
      item.issueDesc.includes(filters.keyword)

    const projectMatch =
      !filters.projectName || item.projectName === filters.projectName

    const deviceTypeMatch =
      !filters.deviceType || item.deviceType === filters.deviceType

    const sourceMatch =
      !filters.issueSource || item.issueSource === filters.issueSource

    const closeStatusMatch =
      !filters.closeStatus || item.closeStatus === filters.closeStatus

    return keywordMatch && projectMatch && deviceTypeMatch && sourceMatch && closeStatusMatch
  })
})

function getCloseStatusText(status) {
  const map = {
    open: '未关闭',
    closed: '已关闭'
  }

  return map[status] || status
}

function getSourceClass(source) {
  const map = {
    研发: 'dev',
    测试: 'test',
    生产: 'production',
    售后: 'aftersales'
  }

  return map[source] || 'other'
}

function getLevelClass(level) {
  const map = {
    紧急: 'urgent',
    高: 'high',
    中: 'middle',
    低: 'low'
  }

  return map[level] || 'low'
}

function resetFilters() {
  filters.keyword = ''
  filters.projectName = ''
  filters.deviceType = ''
  filters.issueSource = ''
  filters.closeStatus = ''
}

function openCreateDialog() {
  editMode.value = 'create'
  currentEditIssue.value = null

  issueForm.projectName = ''
  issueForm.deviceType = ''
  issueForm.issueSource = ''
  issueForm.level = '中'
  issueForm.issueTitle = ''
  issueForm.owner = ''
  issueForm.creator = ''
  issueForm.planCloseTime = ''
  issueForm.issueDesc = ''

  showEditDialog.value = true
}

function openEditDialog(item) {
  editMode.value = 'edit'
  currentEditIssue.value = item

  issueForm.projectName = item.projectName
  issueForm.deviceType = item.deviceType
  issueForm.issueSource = item.issueSource
  issueForm.level = item.level
  issueForm.issueTitle = item.issueTitle
  issueForm.owner = item.owner
  issueForm.creator = item.creator
  issueForm.planCloseTime = item.planCloseTime
  issueForm.issueDesc = item.issueDesc

  showEditDialog.value = true
}

function saveIssue() {
  if (!issueForm.projectName) {
    alert('请选择问题绑定项目')
    return
  }

  if (!issueForm.deviceType) {
    alert('请选择问题绑定终端')
    return
  }

  if (!issueForm.issueSource) {
    alert('请选择问题来源')
    return
  }

  if (!issueForm.issueTitle) {
    alert('请输入问题标题')
    return
  }

  if (!issueForm.owner) {
    alert('请输入问题负责人')
    return
  }

  if (!issueForm.creator) {
    alert('请输入问题提出人')
    return
  }

  if (editMode.value === 'create') {
    issueList.value.unshift({
      id: Date.now(),
      projectName: issueForm.projectName,
      deviceType: issueForm.deviceType,
      issueSource: issueForm.issueSource,
      level: issueForm.level,
      issueTitle: issueForm.issueTitle,
      owner: issueForm.owner,
      creator: issueForm.creator,
      createTime: new Date().toISOString().slice(0, 10),
      planCloseTime: issueForm.planCloseTime,
      realCloseTime: '',
      closeStatus: 'open',
      closeUser: '',
      issueDesc: issueForm.issueDesc,
      replies: []
    })
  } else {
    currentEditIssue.value.projectName = issueForm.projectName
    currentEditIssue.value.deviceType = issueForm.deviceType
    currentEditIssue.value.issueSource = issueForm.issueSource
    currentEditIssue.value.level = issueForm.level
    currentEditIssue.value.issueTitle = issueForm.issueTitle
    currentEditIssue.value.owner = issueForm.owner
    currentEditIssue.value.creator = issueForm.creator
    currentEditIssue.value.planCloseTime = issueForm.planCloseTime
    currentEditIssue.value.issueDesc = issueForm.issueDesc
  }

  showEditDialog.value = false
}

function viewIssue(item) {
  selectedIssue.value = item
}

function openReplyDialog(item) {
  currentReplyIssue.value = item
  replyForm.replyUser = ''
  replyForm.content = ''
  showReplyDialog.value = true
}

function saveReply() {
  if (!currentReplyIssue.value) return

  if (!replyForm.replyUser) {
    alert('请输入回复人')
    return
  }

  if (!replyForm.content) {
    alert('请输入回复内容')
    return
  }

  currentReplyIssue.value.replies.push({
    id: Date.now(),
    replyUser: replyForm.replyUser,
    replyTime: new Date().toISOString().slice(0, 10),
    content: replyForm.content
  })

  showReplyDialog.value = false
  alert('问题回复已保存')
}

function closeIssue(item) {
  const ok = confirm(`确认关闭问题【${item.issueTitle}】吗？`)
  if (!ok) return

  item.closeStatus = 'closed'
  item.closeUser = '当前用户'
  item.realCloseTime = new Date().toISOString().slice(0, 10)

  selectedIssue.value = null
  alert(`问题【${item.issueTitle}】已关闭`)
}

function reopenIssue(item) {
  const ok = confirm(`确认重新打开问题【${item.issueTitle}】吗？`)
  if (!ok) return

  item.closeStatus = 'open'
  item.closeUser = ''
  item.realCloseTime = ''

  selectedIssue.value = null
  alert(`问题【${item.issueTitle}】已重新打开`)
}

function exportIssueRecords() {
  const header = [
    '问题标题',
    '绑定项目',
    '终端类型',
    '问题来源',
    '严重等级',
    '负责人',
    '提出人',
    '提出时间',
    '计划关闭时间',
    '实际关闭时间',
    '关闭状态',
    '关闭人',
    '问题描述',
    '回复记录'
  ]

  const rows = issueList.value.map(item => [
    item.issueTitle,
    item.projectName,
    item.deviceType,
    item.issueSource,
    item.level,
    item.owner,
    item.creator,
    item.createTime,
    item.planCloseTime || '',
    item.realCloseTime || '',
    getCloseStatusText(item.closeStatus),
    item.closeUser || '',
    item.issueDesc || '',
    item.replies.map(reply => `${reply.replyTime} ${reply.replyUser}: ${reply.content}`).join('；')
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
  link.download = '问题闭环记录.csv'
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

.header-actions {
  display: flex;
  gap: 12px;
}

.primary-btn,
.secondary-btn,
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

.secondary-btn,
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
  grid-template-columns: 1.4fr 180px 160px 140px 140px 90px 90px;
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

.issue-link {
  border: none;
  background: transparent;
  color: #60a5fa;
  font-size: 13px;
  font-weight: 800;
  cursor: pointer;
  padding: 0;
  text-align: left;
}

.issue-link:hover {
  color: #93c5fd;
  text-decoration: underline;
}

.issue-desc {
  margin-top: 4px;
  color: #64748b;
  font-size: 12px;
  line-height: 1.5;
  word-break: break-all;
}

.project-tag,
.device-tag,
.source-tag,
.level-tag,
.close-tag {
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

.source-tag.dev {
  background: #9333ea33;
  color: #c084fc;
}

.source-tag.test {
  background: #1d4ed833;
  color: #60a5fa;
}

.source-tag.production {
  background: #d9770633;
  color: #fbbf24;
}

.source-tag.aftersales {
  background: #be123c33;
  color: #fb7185;
}

.source-tag.other {
  background: #47556933;
  color: #94a3b8;
}

.level-tag.urgent {
  background: #dc262633;
  color: #f87171;
}

.level-tag.high {
  background: #ea580c33;
  color: #fb923c;
}

.level-tag.middle {
  background: #d9770633;
  color: #fbbf24;
}

.level-tag.low {
  background: #47556933;
  color: #94a3b8;
}

.close-tag.open {
  background: #16a34a33;
  color: #4ade80;
}

.close-tag.closed {
  background: #64748b33;
  color: #cbd5e1;
}

.muted {
  color: #94a3b8 !important;
}

.operation-col {
  width: 330px;
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

.text-btn.yellow {
  color: #fbbf24;
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

.reply-form {
  padding-top: 0;
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
.remark-card,
.reply-card {
  background: #020617;
  border: 1px solid #1e293b;
  border-radius: 10px;
  padding: 12px;
}

.detail-card span,
.remark-card span,
.reply-card > span {
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

.remark-card,
.reply-card {
  margin: 0 20px 20px;
}

.remark-card p {
  margin: 0;
  color: #cbd5e1;
  font-size: 13px;
  line-height: 1.6;
}

.empty-reply {
  color: #64748b;
  font-size: 13px;
  padding: 8px 0;
}

.reply-item {
  border-top: 1px solid #1e293b;
  padding-top: 12px;
  margin-top: 12px;
}

.reply-header {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 6px;
}

.reply-header strong {
  color: #f8fafc;
  font-size: 13px;
}

.reply-header span {
  color: #64748b;
  font-size: 12px;
}

.reply-item p {
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
    min-width: 1300px;
  }

  .form-grid,
  .detail-card {
    grid-template-columns: 1fr;
  }
}
</style>