<template>
  <header ref="headerRef" class="top-header">
    <div class="header-left">
      <button class="menu-btn" type="button">
        ☰
      </button>

      <div class="search-wrap">
        <form class="search-box" @submit.prevent="handleSearch">
        <span class="search-icon">⌕</span>
        <input
          v-model="keyword"
          type="text"
          placeholder="全局搜索项目、版本或SN..."
          @focus="openSearchPanel"
        />
        </form>

        <div v-if="showSearchPanel" class="search-popover">
          <div v-if="searching" class="search-empty">搜索中...</div>

          <template v-else>
            <div v-if="searchResults.projects.length" class="search-group">
              <strong>项目</strong>
              <button
                v-for="item in searchResults.projects"
                :key="`project-${item.id}`"
                type="button"
                class="search-item"
                @click="goSearchResult('project', item)"
              >
                <span>{{ item.project_name || item.projectName }}</span>
                <em>{{ item.project_code || item.projectCode || '项目看板' }}</em>
              </button>
            </div>

            <div v-if="searchResults.versions.length" class="search-group">
              <strong>版本</strong>
              <button
                v-for="item in searchResults.versions"
                :key="`version-${item.version_type}-${item.id}`"
                type="button"
                class="search-item"
                @click="goSearchResult('version', item)"
              >
                <span>{{ item.software_version || item.hardware_version || '未命名版本' }}</span>
                <em>{{ item.project_name || item.device_type || '版本矩阵' }}</em>
              </button>
            </div>

            <div v-if="searchResults.devices.length" class="search-group">
              <strong>终端</strong>
              <button
                v-for="item in searchResults.devices"
                :key="`device-${item.location}-${item.id}`"
                type="button"
                class="search-item"
                @click="goSearchResult('device', item)"
              >
                <span>{{ item.sn || item.mac_address || '未填写SN/MAC' }}</span>
                <em>{{ item.location === 'outbound' ? '出库记录' : '库存情况' }} · {{ item.status || '-' }}</em>
              </button>
            </div>

            <div v-if="hasSearched && !hasSearchResult" class="search-empty">
              没有找到匹配结果
            </div>
          </template>
        </div>
      </div>
    </div>

    <div class="header-right">
      <button v-if="isSystemAdmin" class="password-btn" type="button" @click="openUserDialog">
        新增用户
      </button>

      <button class="password-btn" type="button" @click="openPasswordDialog">
        修改密码
      </button>

      <div v-if="!isSystemAdmin" class="notify-wrap">
        <button class="notify-btn" type="button" @click="toggleNotifications">
          🔔
          <span v-if="notificationCount > 0" class="notify-dot"></span>
          <span v-if="notificationCount > 0" class="notify-count">
            {{ notificationCount > 99 ? '99+' : notificationCount }}
          </span>
        </button>

        <div v-if="showNotifications" class="notify-popover">
          <div class="notify-header">
            <strong>消息提醒</strong>
            <button type="button" @click="loadNotifications">刷新</button>
          </div>

          <div class="notify-list">
            <div
              v-for="item in notifications"
              :key="`${item.type || 'todo'}-${item.id}`"
              class="notify-item"
              @click="goNotification(item)"
            >
              <button
                class="notify-title"
                :class="{ expanded: expandedNotificationIds.includes(item.id) }"
                type="button"
                :title="item.title"
                @click.stop="toggleNotificationExpand(item)"
              >
                {{ item.title }}
              </button>
              <button
                v-if="isLongNotification(item)"
                class="notify-expand"
                type="button"
                @click.stop="toggleNotificationExpand(item)"
              >
                {{ expandedNotificationIds.includes(item.id) ? '收起' : '展开' }}
              </button>
              <span class="notify-meta">
                {{ item.deadline || '暂无截止时间' }} · {{ item.level || '普通' }}
              </span>
              <button
                v-if="item.type === 'productionRequest'"
                class="notify-action"
                type="button"
                @click.stop="confirmProductionRequestTodo(item)"
              >
                确认
              </button>
              <button
                v-if="item.type === 'issue'"
                class="notify-action"
                type="button"
                @click.stop="confirmIssueTodo(item)"
              >
                确认
              </button>
              <button
                v-if="item.type === 'requirementChangeReceipt'"
                class="notify-action"
                type="button"
                @click.stop="confirmRequirementChangeTodo(item)"
              >
                确认
              </button>
            </div>

            <div v-if="notifications.length === 0" class="notify-empty">
              暂无待办消息
            </div>
          </div>
        </div>
      </div>

      <div class="divider"></div>

      <div class="date-box">
        <span>{{ todayText }}</span>
        <span class="calendar-icon">📅</span>
      </div>
    </div>
  </header>

  <div v-if="showPasswordDialog" class="dialog-mask">
    <div class="password-dialog">
      <div class="dialog-header">
        <h3>修改密码</h3>
        <button type="button" @click="closePasswordDialog">×</button>
      </div>

      <div class="password-form">
        <label>
          原密码
          <input
            v-model="passwordForm.oldPassword"
            type="password"
            autocomplete="current-password"
            placeholder="请输入当前密码"
          />
        </label>

        <label>
          新密码
          <input
            v-model="passwordForm.newPassword"
            type="password"
            autocomplete="new-password"
            placeholder="至少6位"
          />
        </label>

        <label>
          确认新密码
          <input
            v-model="passwordForm.confirmPassword"
            type="password"
            autocomplete="new-password"
            placeholder="再次输入新密码"
          />
        </label>

        <p v-if="passwordError" class="password-error">{{ passwordError }}</p>
      </div>

      <div class="dialog-footer">
        <button class="cancel-btn" type="button" @click="closePasswordDialog">
          取消
        </button>
        <button class="save-btn" type="button" :disabled="passwordSaving" @click="submitPasswordChange">
          {{ passwordSaving ? '保存中...' : '保存' }}
        </button>
      </div>
    </div>
  </div>

  <div v-if="showUserDialog" class="dialog-mask">
    <div class="user-dialog">
      <div class="dialog-header">
        <h3>新增用户</h3>
        <button type="button" @click="closeUserDialog">×</button>
      </div>

      <div class="user-form">
        <label>
          登录账号
          <input v-model.trim="userForm.username" placeholder="请输入登录账号" />
        </label>

        <label>
          真实姓名
          <input v-model.trim="userForm.realName" placeholder="请输入真实姓名" />
        </label>

        <label>
          初始密码
          <input v-model="userForm.password" type="password" placeholder="默认 123456" />
        </label>

        <label>
          部门
          <select v-model="userForm.department">
            <option value="">请选择部门</option>
            <option
              v-for="department in departmentOptions"
              :key="department"
              :value="department"
            >
              {{ department }}
            </option>
          </select>
        </label>

        <label>
          邮箱
          <input v-model.trim="userForm.email" placeholder="可选" />
        </label>

        <label>
          手机/微信
          <input v-model.trim="userForm.phone" placeholder="可选" />
        </label>

        <label>
          状态
          <select v-model="userForm.status">
            <option value="启用">启用</option>
            <option value="禁用">禁用</option>
          </select>
        </label>

        <div class="role-field">
          <span>角色</span>
          <div class="role-grid">
            <label
              v-for="role in roleOptions"
              :key="role.roleCode"
              class="role-option"
            >
              <input
                v-model="userForm.roleCodes"
                type="checkbox"
                :value="role.roleCode"
              />
              <span>{{ role.roleName }}</span>
            </label>
          </div>
        </div>

        <p v-if="userError" class="password-error">{{ userError }}</p>
      </div>

      <div class="account-section">
        <div class="account-section-header">
          <strong>已有账户</strong>
          <button type="button" @click="loadUsers">刷新</button>
        </div>

        <div class="account-list">
          <div v-for="user in userList" :key="user.id" class="account-item">
            <div>
              <strong>{{ user.realName || user.username }}</strong>
              <span>{{ user.username }} · {{ user.department || '未填写部门' }} · {{ user.roles || '未绑定角色' }}</span>
            </div>
            <button
              type="button"
              class="permission-account-btn"
              @click="openPermissionDialog(user)"
            >
              权限
            </button>
            <button
              type="button"
              class="delete-account-btn"
              :disabled="user.username === 'admin' || deletingUserId === user.id"
              @click="deleteUser(user)"
            >
              {{ deletingUserId === user.id ? '删除中' : '删除' }}
            </button>
          </div>

          <div v-if="userList.length === 0" class="account-empty">
            暂无账户数据
          </div>
        </div>
      </div>

      <div class="dialog-footer">
        <button class="cancel-btn" type="button" @click="closeUserDialog">
          取消
        </button>
        <button class="save-btn" type="button" :disabled="userSaving" @click="submitCreateUser">
          {{ userSaving ? '保存中...' : '保存' }}
        </button>
      </div>
    </div>
  </div>

  <div v-if="showPermissionDialog" class="dialog-mask">
    <div class="account-dialog permission-dialog">
      <div class="dialog-header">
        <div>
          <h3>配置用户功能权限</h3>
          <p>{{ permissionUser?.realName || permissionUser?.username }}</p>
        </div>
        <button type="button" @click="closePermissionDialog">×</button>
      </div>

      <div class="permission-section">
        <div
          v-for="group in permissionGroups"
          :key="group.module"
          class="permission-group"
        >
          <strong>{{ group.module }}</strong>
          <div class="permission-grid">
            <label
              v-for="permission in group.permissions"
              :key="permission.code"
              class="role-option"
            >
              <input
                v-model="permissionForm.permissions"
                type="checkbox"
                :value="permission.code"
              />
              <span>{{ permission.name }}</span>
            </label>
          </div>
        </div>
      </div>

      <p v-if="permissionError" class="password-error">{{ permissionError }}</p>

      <div class="dialog-footer">
        <button class="cancel-btn" type="button" @click="closePermissionDialog">
          取消
        </button>
        <button class="save-btn" type="button" :disabled="permissionSaving" @click="savePermissionConfig">
          {{ permissionSaving ? '保存中...' : '保存权限' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import { getDashboardSummary, globalSearch, markNotificationRead } from '@/api/report'
import { confirmProductionRequest } from '@/api/shippingBatch'
import { confirmIssue } from '@/api/issue'
import { confirmRequirementChange } from '@/api/requirement'
import { changePasswordApi, createUserApi, deleteUserApi, getRoleOptionsApi, getUserPermissionsApi, getUsersApi, saveUserPermissionsApi } from '@/api/auth'
import { getCurrentUserParams } from '@/utils/currentUser'
import { getStoredRoles } from '@/utils/permission'

const router = useRouter()
const keyword = ref('')
const showNotifications = ref(false)
const notifications = ref([])
const showSearchPanel = ref(false)
const searching = ref(false)
const hasSearched = ref(false)
const headerRef = ref(null)
const expandedNotificationIds = ref([])
const showPasswordDialog = ref(false)
const showUserDialog = ref(false)
const passwordSaving = ref(false)
const userSaving = ref(false)
const permissionSaving = ref(false)
const deletingUserId = ref(0)
const passwordError = ref('')
const userError = ref('')
const permissionError = ref('')
const roleOptions = ref([])
const userList = ref([])
const showPermissionDialog = ref(false)
const permissionUser = ref(null)
const permissionForm = reactive({
  permissions: []
})
const permissionGroups = [
  {
    module: '项目管理',
    permissions: [
      { code: 'project:view', name: '查看' },
      { code: 'project:create', name: '新增' },
      { code: 'project:update', name: '修改' },
      { code: 'project:delete', name: '删除' },
      { code: 'project:submit', name: '提交审核' },
      { code: 'project:audit', name: '审核' },
      { code: 'project:close', name: '关闭项目' },
      { code: 'project:reopen', name: '重新打开' },
      { code: 'project:archive', name: '归档' }
    ]
  },
  {
    module: '需求管理',
    permissions: [
      { code: 'requirement:view', name: '需求查看' },
      { code: 'requirement:download', name: '需求下载' },
      { code: 'requirement:upload', name: '需求上传' },
      { code: 'requirement:submit', name: '需求提交' },
      { code: 'requirement:delete', name: '需求删除' },
      { code: 'requirement:close', name: '需求关闭' },
      { code: 'requirement:audit', name: '需求审核' },
      { code: 'customer:view', name: '客供资料查看' },
      { code: 'customer:upload', name: '客供资料上传' },
      { code: 'customer:download', name: '客供资料下载' },
      { code: 'customer:delete', name: '客供资料删除' }
    ]
  },
  {
    module: '硬件管理',
    permissions: [
      { code: 'hardware:view', name: '硬件查看' },
      { code: 'hardware:create', name: '硬件版本新增' },
      { code: 'hardware:update', name: '硬件版本修改' },
      { code: 'hardware:upload', name: '硬件上传' },
      { code: 'hardware:download', name: '硬件下载' },
      { code: 'hardware:delete', name: '硬件删除' },
      { code: 'hardware:submit', name: '硬件提交' },
      { code: 'hardware:audit', name: '硬件审核' },
      { code: 'hardware-test:view', name: '硬件测试查看' },
      { code: 'hardware-dev-doc:view', name: '开发文档查看' },
      { code: 'hardware-dev-doc:upload', name: '开发文档上传' },
      { code: 'hardware-dev-doc:download', name: '开发文档下载' },
      { code: 'hardware-dev-doc:delete', name: '开发文档删除' },
      { code: 'board-composition:view', name: '板卡组成查看' },
      { code: 'board-composition:upload', name: '板卡组成上传' },
      { code: 'board-composition:create', name: '板卡组成新增' },
      { code: 'board-composition:delete', name: '板卡组成删除' }
    ]
  },
  {
    module: '软件管理',
    permissions: [
      { code: 'software:view', name: '软件版本查看' },
      { code: 'software:create', name: '软件版本新增' },
      { code: 'software:update', name: '软件版本修改' },
      { code: 'software:delete', name: '软件版本删除' },
      { code: 'software:download', name: '软件版本下载' },
      { code: 'software:release', name: '软件版本发布' },
      { code: 'branch:view', name: '项目分支查看' },
      { code: 'branch:create', name: '项目分支新增' },
      { code: 'branch:update', name: '项目分支修改' },
      { code: 'branch:delete', name: '项目分支删除' },
      { code: 'branch:download', name: '项目分支下载' }
    ]
  },
  {
    module: '测试管理',
    permissions: [
      { code: 'testcase:view', name: '测试用例查看' },
      { code: 'testcase:upload', name: '测试用例上传' },
      { code: 'testcase:uploadReport', name: '测试报告上传' },
      { code: 'testcase:download', name: '测试资料下载' },
      { code: 'testcase:submit', name: '测试用例提交' },
      { code: 'testcase:delete', name: '测试用例删除' },
      { code: 'testcase:audit', name: '测试用例审核' },
      { code: 'issue:view', name: '问题查看' },
      { code: 'issue:create', name: '问题新增' },
      { code: 'issue:update', name: '问题修改' },
      { code: 'issue:reply', name: '问题回复' },
      { code: 'issue:close', name: '问题关闭' },
      { code: 'issue:reopen', name: '问题重开' },
      { code: 'issue:export', name: '问题导出' }
    ]
  },
  {
    module: '生产管理',
    permissions: [
      { code: 'production:view', name: '查看' },
      { code: 'production:create', name: '新增' },
      { code: 'production:update', name: '修改' },
      { code: 'production:delete', name: '删除' },
      { code: 'production:audit', name: '审核' },
      { code: 'burn:deleteBatch', name: '生产烧录删除批次' },
      { code: 'production:outline:view', name: '测试大纲查看' },
      { code: 'production:outline:upload', name: '测试大纲上传' },
      { code: 'production:outline:update', name: '测试大纲修改' },
      { code: 'production:outline:delete', name: '测试大纲删除' },
      { code: 'board-inbound:view', name: '板卡入库查看' },
      { code: 'board-inbound:import', name: '板卡入库上传' },
      { code: 'board-inbound:delete', name: '板卡入库删除' },
      { code: 'inventory:view', name: '库存情况查看' }
    ]
  },
  {
    module: '发货管理',
    permissions: [
      { code: 'shipping:view', name: '查看' },
      { code: 'shipping:manage', name: '业务管理' },
      { code: 'shipping:create', name: '新增' },
      { code: 'shipping:update', name: '修改' },
      { code: 'shipping:delete', name: '删除' },
      { code: 'shipping:submit', name: '提交审核' },
      { code: 'shipping:audit', name: '审核' }
    ]
  },
  {
    module: '售后管理',
    permissions: [
      { code: 'aftersales:view', name: '查看' },
      { code: 'aftersales:create', name: '新增' },
      { code: 'aftersales:update', name: '修改' },
      { code: 'aftersales:delete', name: '删除' },
      { code: 'aftersales:upload', name: '上传' },
      { code: 'aftersales:download', name: '下载' },
      { code: 'aftersales:submit', name: '提交审核' },
      { code: 'aftersales:export', name: '导出' },
      { code: 'aftersales:audit', name: '审核' }
    ]
  },
  {
    module: '统计报表',
    permissions: [
      { code: 'report:view', name: '查看' },
      { code: 'report:manage', name: '管理' }
    ]
  }
]
const departmentOptions = [
  '管理部',
  '项目助理',
  '软件研发',
  '硬件研发',
  '生产',
  '发货',
  '发货审核',
  '售后',
  '质量检查',
  '领导'
]
const searchResults = reactive({
  projects: [],
  versions: [],
  devices: []
})
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})
const userForm = reactive({
  username: '',
  realName: '',
  password: '123456',
  email: '',
  phone: '',
  department: '',
  status: '启用',
  roleCodes: []
})

const notificationCount = computed(() => notifications.value.length)
const isSystemAdmin = computed(() => getStoredRoles().includes('system_admin'))
const hasSearchResult = computed(() => {
  return searchResults.projects.length > 0 ||
    searchResults.versions.length > 0 ||
    searchResults.devices.length > 0
})

onMounted(() => {
  if (isSystemAdmin.value) {
    clearAdminNotifications()
  } else {
    loadNotifications()
  }
  document.addEventListener('click', handleDocumentClick)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleDocumentClick)
})

const todayText = computed(() => {
  const date = new Date()
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()
  return `${year}年${month}月${day}日`
})

function clearAdminNotifications() {
  notifications.value = []
  expandedNotificationIds.value = []
  showNotifications.value = false
}

async function loadNotifications() {
  if (isSystemAdmin.value) {
    clearAdminNotifications()
    return
  }

  try {
    const res = await getDashboardSummary(getCurrentUserParams())
    const result = res?.data || res
    if (result.code !== 200) return
    notifications.value = [
      ...(result.data?.todos || []),
      ...(result.data?.notifications || [])
    ]
    const currentIds = new Set(notifications.value.map(item => item.id))
    expandedNotificationIds.value = expandedNotificationIds.value.filter(id => currentIds.has(id))
  } catch (err) {
    console.error('加载消息提醒失败：', err)
  }
}

function isLongNotification(item) {
  return String(item?.title || '').length > 24
}

function toggleNotificationExpand(item) {
  if (!isLongNotification(item)) return
  const id = item.id
  if (expandedNotificationIds.value.includes(id)) {
    expandedNotificationIds.value = expandedNotificationIds.value.filter(itemId => itemId !== id)
  } else {
    expandedNotificationIds.value = [...expandedNotificationIds.value, id]
  }
}

async function toggleNotifications() {
  if (isSystemAdmin.value) {
    clearAdminNotifications()
    return
  }

  showNotifications.value = !showNotifications.value
  if (showNotifications.value) {
    await loadNotifications()
  }
}

async function goNotification(item) {
  showNotifications.value = false

  if (String(item.type || '').endsWith('AuditResult') || item.type === 'productionRequestConfirmed' || item.type === 'issueConfirmed' || item.type === 'requirementChangeConfirmed') {
    await markAuditResultRead(item)
    notifications.value = notifications.value.filter(notification => notification.id !== item.id)
  }

  if (item.link) {
    router.push(item.link)
  }
}

async function confirmProductionRequestTodo(item) {
  const id = Number(String(item.id || '').replace('production-request-', ''))
  if (!id) return

  try {
    const user = getCurrentUserParams()
    const res = await confirmProductionRequest(id, {
      confirmerId: user.userId,
      confirmerName: user.realName || user.username || '生产人员'
    })
    const result = res?.data || res
    if (result.code !== 200) {
      alert(result.msg || '确认失败')
      return
    }
    notifications.value = notifications.value.filter(notification => notification.id !== item.id)
    await loadNotifications()
    alert('已确认生产请求，发货人员将收到通知')
  } catch (err) {
    console.error('确认生产请求失败：', err)
    alert(err.response?.data || '确认失败，请检查后端接口')
  }
}

async function confirmIssueTodo(item) {
  const id = Number(String(item.id || '').replace('issue-', ''))
  if (!id) return

  try {
    const user = getCurrentUserParams()
    const res = await confirmIssue(id, {
      confirmUserId: user.userId,
      confirmUserName: user.realName || user.username || '负责人'
    })
    const result = res?.data || res
    if (result.code !== 200) {
      alert(result.msg || '确认失败')
      return
    }
    notifications.value = notifications.value.filter(notification => notification.id !== item.id)
    await loadNotifications()
    alert('已确认问题，项目助理将收到通知')
  } catch (err) {
    console.error('确认问题失败：', err)
    alert(err.response?.data || '确认失败，请检查后端接口')
  }
}

async function confirmRequirementChangeTodo(item) {
  const id = Number(String(item.id || '').replace('requirement-change-receipt-', ''))
  if (!id) return

  try {
    const user = getCurrentUserParams()
    const res = await confirmRequirementChange(id, {
      confirmUserId: user.userId,
      confirmUserName: user.realName || user.username || '软件负责人'
    })
    const result = res?.data || res
    if (result.code !== 200) {
      alert(result.msg || '确认失败')
      return
    }
    notifications.value = notifications.value.filter(notification => notification.id !== item.id)
    await loadNotifications()
    alert('已确认需求变更，上传人将收到通知')
  } catch (err) {
    console.error('确认需求变更失败：', err)
    alert(err.response?.data || '确认失败，请检查后端接口')
  }
}

async function markAuditResultRead(item) {
  try {
    const user = getCurrentUserParams()
    await markNotificationRead({
      userId: user.userId,
      username: user.username,
      notificationId: item.id
    })
  } catch (err) {
    console.error('标记消息已读失败：', err)
  }
}

function handleDocumentClick(event) {
  if (!headerRef.value?.contains(event.target)) {
    showNotifications.value = false
    showSearchPanel.value = false
  }
}

function openPasswordDialog() {
  passwordError.value = ''
  passwordForm.oldPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
  showPasswordDialog.value = true
}

function closePasswordDialog() {
  if (passwordSaving.value) return
  showPasswordDialog.value = false
}

async function openUserDialog() {
  userError.value = ''
  userForm.username = ''
  userForm.realName = ''
  userForm.password = '123456'
  userForm.email = ''
  userForm.phone = ''
  userForm.department = ''
  userForm.status = '启用'
  userForm.roleCodes = []
  showUserDialog.value = true
  await Promise.all([loadRoleOptions(), loadUsers()])
}

function closeUserDialog() {
  if (userSaving.value) return
  showUserDialog.value = false
}

async function loadRoleOptions() {
  if (roleOptions.value.length > 0) return
  try {
    const res = await getRoleOptionsApi()
    const result = res?.data || res
    if (result.code !== 200) {
      userError.value = result.msg || '加载角色失败'
      return
    }
    roleOptions.value = result.data || []
  } catch (err) {
    console.error('加载角色失败：', err)
    userError.value = err.response?.data?.msg || err.response?.data || '加载角色失败'
  }
}

async function loadUsers() {
  try {
    const res = await getUsersApi()
    const result = res?.data || res
    if (result.code !== 200) {
      userError.value = result.msg || '加载用户失败'
      return
    }
    userList.value = result.data || []
  } catch (err) {
    console.error('加载用户失败：', err)
    userError.value = err.response?.data?.msg || err.response?.data || '加载用户失败'
  }
}

async function deleteUser(user) {
  if (!user?.id || user.username === 'admin') return
  if (!confirm(`确认删除账户【${user.realName || user.username}】吗？`)) return

  try {
    deletingUserId.value = user.id
    const res = await deleteUserApi(user.id)
    const result = res?.data || res
    if (result.code !== 200) {
      userError.value = result.msg || '删除账户失败'
      return
    }
    alert('删除账户成功')
    await loadUsers()
  } catch (err) {
    console.error('删除账户失败：', err)
    userError.value = err.response?.data?.msg || err.response?.data || '删除账户失败，请检查后端接口'
  } finally {
    deletingUserId.value = 0
  }
}

async function openPermissionDialog(user) {
  permissionUser.value = user
  permissionError.value = ''
  permissionForm.permissions = []
  showPermissionDialog.value = true

  try {
    const res = await getUserPermissionsApi(user.id)
    const result = res?.data || res
    if (result.code !== 200) {
      permissionError.value = result.msg || '加载用户权限失败'
      return
    }
    permissionForm.permissions = result.data || []
  } catch (err) {
    console.error('加载用户权限失败：', err)
    permissionError.value = err.response?.data?.msg || err.response?.data || '加载用户权限失败'
  }
}

function closePermissionDialog() {
  if (permissionSaving.value) return
  showPermissionDialog.value = false
  permissionUser.value = null
}

async function savePermissionConfig() {
  if (!permissionUser.value?.id) return
  try {
    permissionSaving.value = true
    const res = await saveUserPermissionsApi(permissionUser.value.id, permissionForm.permissions)
    const result = res?.data || res
    if (result.code !== 200) {
      permissionError.value = result.msg || '保存权限失败'
      return
    }
    const currentUser = JSON.parse(localStorage.getItem('user') || '{}')
    if (Number(currentUser.id) === Number(permissionUser.value.id)) {
      localStorage.setItem('permissions', JSON.stringify(result.data || permissionForm.permissions))
    }
    alert('保存权限成功')
    closePermissionDialog()
  } catch (err) {
    console.error('保存用户权限失败：', err)
    permissionError.value = err.response?.data?.msg || err.response?.data || '保存用户权限失败'
  } finally {
    permissionSaving.value = false
  }
}

async function submitCreateUser() {
  userError.value = ''

  if (!userForm.username) {
    userError.value = '请输入登录账号'
    return
  }

  if ((userForm.password || '').length < 6) {
    userError.value = '初始密码至少需要6位'
    return
  }

  if (userForm.roleCodes.length === 0) {
    userError.value = '请至少选择一个角色'
    return
  }

  try {
    userSaving.value = true
    const res = await createUserApi({
      username: userForm.username,
      password: userForm.password || '123456',
      realName: userForm.realName || userForm.username,
      email: userForm.email,
      phone: userForm.phone,
      department: userForm.department,
      status: userForm.status,
      roleCodes: userForm.roleCodes
    })
    const result = res?.data || res
    if (result.code !== 200) {
      userError.value = result.msg || '新增用户失败'
      return
    }
    alert('新增用户成功')
    await loadUsers()
    closeUserDialog()
  } catch (err) {
    console.error('新增用户失败：', err)
    userError.value = err.response?.data?.msg || err.response?.data || '新增用户失败，请检查后端接口'
  } finally {
    userSaving.value = false
  }
}

async function submitPasswordChange() {
  passwordError.value = ''

  if (!passwordForm.oldPassword) {
    passwordError.value = '请输入原密码'
    return
  }

  if (passwordForm.newPassword.length < 6) {
    passwordError.value = '新密码至少需要6位'
    return
  }

  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    passwordError.value = '两次输入的新密码不一致'
    return
  }

  try {
    passwordSaving.value = true
    const user = getCurrentUserParams()
    const res = await changePasswordApi({
      userId: user.userId,
      username: user.username,
      oldPassword: passwordForm.oldPassword,
      newPassword: passwordForm.newPassword
    })
    const result = res?.data || res

    if (result.code !== 200) {
      passwordError.value = result.msg || '修改密码失败'
      return
    }

    alert(result.msg || '密码修改成功，请重新登录')
    localStorage.clear()
    router.push('/login')
  } catch (err) {
    console.error('修改密码失败：', err)
    passwordError.value = err.response?.data?.msg || err.response?.data || '修改密码失败，请检查后端接口'
  } finally {
    passwordSaving.value = false
  }
}

function openSearchPanel() {
  if (hasSearched.value || keyword.value.trim()) {
    showSearchPanel.value = true
  }
}

async function handleSearch() {
  const text = keyword.value.trim()
  if (!text) return

  showSearchPanel.value = true
  hasSearched.value = true
  searching.value = true

  try {
    const res = await globalSearch({ keyword: text })
    const result = res?.data || res
    if (result.code !== 200) return

    searchResults.projects = result.data?.projects || []
    searchResults.versions = result.data?.versions || []
    searchResults.devices = result.data?.devices || []

    const total =
      searchResults.projects.length +
      searchResults.versions.length +
      searchResults.devices.length

    if (total === 1) {
      if (searchResults.projects.length === 1) {
        goSearchResult('project', searchResults.projects[0])
      } else if (searchResults.versions.length === 1) {
        goSearchResult('version', searchResults.versions[0])
      } else if (searchResults.devices.length === 1) {
        goSearchResult('device', searchResults.devices[0])
      }
    }
  } catch (err) {
    console.error('全局搜索失败：', err)
  } finally {
    searching.value = false
  }
}

function goSearchResult(type, item) {
  showSearchPanel.value = false
  const text = keyword.value.trim()

  if (type === 'project') {
    const projectKeyword = item.project_name || item.projectName || text
    router.push({ path: '/dashboard', query: { keyword: projectKeyword, projectId: item.id } })
    return
  }

  if (type === 'version') {
    const versionKeyword = item.software_version || item.hardware_version || text
    router.push({ path: '/version/matrix', query: { keyword: versionKeyword, type: item.version_type || '' } })
    return
  }

  const snKeyword = item.sn || item.mac_address || text
  router.push({
    path: item.location === 'outbound' ? '/shipping/out' : '/production/inventory',
    query: { keyword: snKeyword }
  })
}
</script>

<style scoped>
.top-header {
  height: 64px;
  flex-shrink: 0;
  background: #0f172a;
  border-bottom: 1px solid #1e293b;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 28px;
  color: #f8fafc;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.menu-btn {
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 8px;
  background: transparent;
  color: #94a3b8;
  font-size: 18px;
  cursor: pointer;
}

.menu-btn:hover {
  background: #1e293b;
  color: #f8fafc;
}

.search-wrap {
  position: relative;
}

.search-box {
  width: 340px;
  height: 38px;
  background: #020617;
  border: 1px solid #1e293b;
  border-radius: 999px;
  display: flex;
  align-items: center;
  padding: 0 14px;
}

.search-popover {
  position: absolute;
  top: 46px;
  left: 0;
  width: min(430px, calc(100vw - 32px));
  max-height: 500px;
  border: 1px solid #263244;
  border-radius: 8px;
  background: #0b1220;
  box-shadow: 0 18px 45px rgba(0, 0, 0, 0.35);
  overflow-y: auto;
  z-index: 20;
}

.search-group {
  padding: 10px 0;
  border-bottom: 1px solid #1e293b;
}

.search-group strong {
  display: block;
  padding: 0 14px 8px;
  color: #94a3b8;
  font-size: 12px;
  font-weight: 700;
}

.search-item {
  width: 100%;
  border: none;
  background: transparent;
  color: #e2e8f0;
  text-align: left;
  padding: 9px 14px;
  cursor: pointer;
}

.search-item:hover {
  background: #111c30;
}

.search-item span,
.search-item em {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.search-item span {
  color: #f8fafc;
  font-size: 14px;
  font-style: normal;
  font-weight: 600;
}

.search-item em {
  margin-top: 3px;
  color: #94a3b8;
  font-size: 12px;
  font-style: normal;
}

.search-empty {
  padding: 22px 14px;
  color: #94a3b8;
  text-align: center;
  font-size: 14px;
}

.search-icon {
  color: #64748b;
  margin-right: 8px;
  font-size: 16px;
}

.search-box input {
  flex: 1;
  border: none;
  outline: none;
  background: transparent;
  color: #e2e8f0;
  font-size: 14px;
}

.search-box input::placeholder {
  color: #64748b;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.password-btn {
  height: 34px;
  border: 1px solid #334155;
  border-radius: 6px;
  background: #111827;
  color: #cbd5e1;
  cursor: pointer;
  font-size: 13px;
  font-weight: 700;
  padding: 0 12px;
}

.password-btn:hover {
  border-color: #38bdf8;
  color: #f8fafc;
}

.notify-btn {
  position: relative;
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 999px;
  background: transparent;
  color: #94a3b8;
  cursor: pointer;
  font-size: 16px;
}

.notify-wrap {
  position: relative;
}

.notify-btn:hover {
  background: #1e293b;
  color: #f8fafc;
}

.notify-dot {
  position: absolute;
  top: 9px;
  right: 9px;
  width: 8px;
  height: 8px;
  background: #ef4444;
  border-radius: 999px;
  border: 2px solid #0f172a;
}

.notify-count {
  position: absolute;
  top: -5px;
  right: -8px;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  border-radius: 999px;
  background: #ef4444;
  border: 1px solid #0f172a;
  color: #ffffff;
  font-size: 11px;
  line-height: 16px;
  font-weight: 700;
}

.notify-popover {
  position: absolute;
  top: 46px;
  right: 0;
  width: min(360px, calc(100vw - 32px));
  max-height: 440px;
  border: 1px solid #263244;
  border-radius: 8px;
  background: #0b1220;
  box-shadow: 0 18px 45px rgba(0, 0, 0, 0.35);
  overflow: hidden;
  z-index: 20;
}

.notify-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  border-bottom: 1px solid #1e293b;
}

.notify-header strong {
  color: #f8fafc;
  font-size: 15px;
}

.notify-header button {
  border: none;
  background: transparent;
  color: #38bdf8;
  cursor: pointer;
  font-size: 13px;
}

.notify-list {
  max-height: 376px;
  overflow-y: auto;
}

.notify-item {
  width: 100%;
  border: none;
  border-bottom: 1px solid #1e293b;
  background: transparent;
  color: #e2e8f0;
  text-align: left;
  padding: 12px 16px;
  cursor: pointer;
}

.notify-item:hover {
  background: #111c30;
}

.notify-title {
  display: block;
  width: 100%;
  max-width: 100%;
  border: none;
  background: transparent;
  color: #f8fafc;
  font-size: 14px;
  font-weight: 600;
  line-height: 1.45;
  text-align: left;
  padding: 0;
  cursor: pointer;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.notify-title.expanded {
  overflow: visible;
  text-overflow: clip;
  white-space: normal;
  word-break: break-word;
}

.notify-expand {
  margin-top: 5px;
  border: none;
  background: transparent;
  color: #38bdf8;
  cursor: pointer;
  font-size: 12px;
  padding: 0;
}

.notify-expand:hover {
  color: #7dd3fc;
  text-decoration: underline;
}

.notify-meta {
  display: block;
  margin-top: 4px;
  color: #94a3b8;
  font-size: 12px;
}

.notify-action {
  margin-top: 8px;
  border: 1px solid rgba(56, 189, 248, 0.45);
  border-radius: 6px;
  background: rgba(14, 165, 233, 0.16);
  color: #7dd3fc;
  cursor: pointer;
  font-size: 12px;
  font-weight: 700;
  padding: 5px 10px;
}

.notify-action:hover {
  background: rgba(14, 165, 233, 0.25);
  color: #e0f2fe;
}

.notify-empty {
  padding: 28px 16px;
  color: #94a3b8;
  text-align: center;
  font-size: 14px;
}

.divider {
  width: 1px;
  height: 28px;
  background: #1e293b;
}

.date-box {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #cbd5e1;
  font-size: 14px;
  font-weight: 500;
}

.calendar-icon {
  color: #64748b;
  font-size: 15px;
}

.dialog-mask {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(2, 6, 23, 0.68);
  padding: 24px;
}

.password-dialog,
.user-dialog {
  width: min(420px, calc(100vw - 32px));
  border: 1px solid #263244;
  border-radius: 8px;
  background: #0b1220;
  box-shadow: 0 24px 70px rgba(0, 0, 0, 0.38);
  color: #e2e8f0;
}

.user-dialog {
  width: min(640px, calc(100vw - 32px));
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 18px;
  border-bottom: 1px solid #1e293b;
}

.dialog-header h3 {
  margin: 0;
  color: #f8fafc;
  font-size: 17px;
}

.dialog-header button {
  width: 30px;
  height: 30px;
  border: none;
  border-radius: 6px;
  background: transparent;
  color: #94a3b8;
  cursor: pointer;
  font-size: 20px;
}

.dialog-header button:hover {
  background: #1e293b;
  color: #f8fafc;
}

.password-form,
.user-form {
  display: grid;
  gap: 14px;
  padding: 18px;
}

.user-form {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.password-form label,
.user-form label {
  display: grid;
  gap: 7px;
  color: #cbd5e1;
  font-size: 13px;
  font-weight: 700;
}

.password-form input,
.user-form input,
.user-form select {
  width: 100%;
  height: 38px;
  border: 1px solid #334155;
  border-radius: 6px;
  background: #020617;
  color: #f8fafc;
  outline: none;
  padding: 0 11px;
}

.password-form input:focus,
.user-form input:focus,
.user-form select:focus {
  border-color: #38bdf8;
}

.role-field {
  grid-column: 1 / -1;
  display: grid;
  gap: 9px;
}

.role-field > span {
  color: #cbd5e1;
  font-size: 13px;
  font-weight: 700;
}

.role-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
}

.role-option {
  min-height: 36px;
  border: 1px solid #263244;
  border-radius: 6px;
  background: #020617;
  display: flex !important;
  grid-template-columns: none !important;
  align-items: center;
  gap: 8px !important;
  padding: 0 10px;
  cursor: pointer;
}

.role-option input {
  width: 14px;
  height: 14px;
  padding: 0;
}

.role-option span {
  color: #e2e8f0;
  font-size: 13px;
  font-weight: 600;
}

.password-error {
  grid-column: 1 / -1;
  margin: 0;
  border: 1px solid rgba(248, 113, 113, 0.35);
  border-radius: 6px;
  background: rgba(127, 29, 29, 0.24);
  color: #fecaca;
  font-size: 13px;
  padding: 9px 10px;
}

.account-section {
  margin: 0 18px 18px;
  border: 1px solid #263244;
  border-radius: 8px;
  background: #020617;
  overflow: hidden;
}

.account-section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-bottom: 1px solid #1e293b;
}

.account-section-header strong {
  color: #f8fafc;
  font-size: 13px;
}

.account-section-header button,
.delete-account-btn,
.permission-account-btn {
  border: none;
  background: transparent;
  color: #38bdf8;
  cursor: pointer;
  font-size: 12px;
  font-weight: 700;
}

.account-list {
  max-height: 220px;
  overflow-y: auto;
}

.account-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 10px 12px;
  border-bottom: 1px solid #111827;
}

.account-item div {
  min-width: 0;
}

.account-item strong,
.account-item span {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.account-item strong {
  color: #f8fafc;
  font-size: 13px;
}

.account-item span {
  margin-top: 3px;
  color: #94a3b8;
  font-size: 12px;
}

.delete-account-btn {
  flex: 0 0 auto;
  color: #f87171;
}

.permission-account-btn {
  flex: 0 0 auto;
  color: #60a5fa;
}

.delete-account-btn:disabled {
  color: #64748b;
  cursor: not-allowed;
}

.permission-dialog {
  width: min(860px, calc(100vw - 32px));
}

.permission-section {
  padding: 18px;
  display: grid;
  gap: 14px;
  max-height: 58vh;
  overflow-y: auto;
}

.permission-group {
  border: 1px solid #263244;
  border-radius: 8px;
  background: #020617;
  padding: 12px;
}

.permission-group > strong {
  display: block;
  margin-bottom: 10px;
  color: #f8fafc;
  font-size: 13px;
}

.permission-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 8px;
}

.account-empty {
  padding: 18px 12px;
  color: #64748b;
  font-size: 13px;
  text-align: center;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 14px 18px 18px;
  border-top: 1px solid #1e293b;
}

.cancel-btn,
.save-btn {
  height: 34px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 700;
  padding: 0 16px;
}

.cancel-btn {
  border: 1px solid #334155;
  background: transparent;
  color: #cbd5e1;
}

.save-btn {
  border: 1px solid #0ea5e9;
  background: #0ea5e9;
  color: #ffffff;
}

.save-btn:disabled {
  cursor: not-allowed;
  opacity: 0.65;
}

@media (max-width: 720px) {
  .user-form,
  .role-grid {
    grid-template-columns: 1fr;
  }
}
</style>
