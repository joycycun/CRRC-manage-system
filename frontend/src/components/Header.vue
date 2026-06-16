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
      <div class="notify-wrap">
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
              <span class="notify-title">{{ item.title }}</span>
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
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import { getDashboardSummary, globalSearch, markNotificationRead } from '@/api/report'
import { confirmProductionRequest } from '@/api/shippingBatch'
import { confirmIssue } from '@/api/issue'
import { getCurrentUserParams } from '@/utils/currentUser'

const router = useRouter()
const keyword = ref('')
const showNotifications = ref(false)
const notifications = ref([])
const showSearchPanel = ref(false)
const searching = ref(false)
const hasSearched = ref(false)
const headerRef = ref(null)
const searchResults = reactive({
  projects: [],
  versions: [],
  devices: []
})

const notificationCount = computed(() => notifications.value.length)
const hasSearchResult = computed(() => {
  return searchResults.projects.length > 0 ||
    searchResults.versions.length > 0 ||
    searchResults.devices.length > 0
})

onMounted(() => {
  loadNotifications()
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

async function loadNotifications() {
  try {
    const res = await getDashboardSummary(getCurrentUserParams())
    const result = res?.data || res
    if (result.code !== 200) return
    notifications.value = [
      ...(result.data?.todos || []),
      ...(result.data?.notifications || [])
    ]
  } catch (err) {
    console.error('加载消息提醒失败：', err)
  }
}

async function toggleNotifications() {
  showNotifications.value = !showNotifications.value
  if (showNotifications.value) {
    await loadNotifications()
  }
}

async function goNotification(item) {
  showNotifications.value = false

  if (String(item.type || '').endsWith('AuditResult') || item.type === 'productionRequestConfirmed' || item.type === 'issueConfirmed') {
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

    if (searchResults.projects.length === 1 && !searchResults.versions.length && !searchResults.devices.length) {
      goSearchResult('project', searchResults.projects[0])
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
    router.push({ path: '/project/progress-report', query: { keyword: projectKeyword } })
    return
  }

  if (type === 'version') {
    const versionKeyword = item.software_version || item.hardware_version || text
    router.push({ path: '/version/matrix', query: { keyword: versionKeyword } })
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
  max-width: 100%;
  color: #f8fafc;
  font-size: 14px;
  font-weight: 600;
  line-height: 1.45;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
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
</style>
