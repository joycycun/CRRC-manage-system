<template>
  <aside class="sidebar">
    <!-- Logo -->
    <div class="sidebar-logo">
      <div class="logo-icon">▦</div>
      <span class="logo-text">PMS Dashboard</span>
    </div>

    <!-- Menu -->
    <nav class="menu">
      <!-- 项目看板 -->
      <router-link
        v-if="canAccess('/dashboard')"
        to="/dashboard"
        class="menu-item"
        :class="{ active: isActive('/dashboard') }"
      >
        <span class="menu-icon">⌁</span>
        <span>项目看板</span>
      </router-link>

      <!-- 项目立项管理 -->
      <router-link
        v-if="canAccess('/project/manage')"
        to="/project/manage"
        class="menu-item"
        :class="{ active: isActive('/project/manage') }"
      >
        <span class="menu-icon">▣</span>
        <span>项目立项管理</span>
      </router-link>

      <!-- 需求管理 -->
      <div v-if="canAccessGroup(['/requirement/book', '/requirement/change', '/requirement/customer-supplied'])" class="menu-group">
        <button class="menu-item menu-button" @click="toggleMenu('requirement')">
          <span class="menu-left">
            <span class="menu-icon">▤</span>
            <span>需求管理</span>
          </span>
          <span class="arrow" :class="{ open: openedMenus.requirement }">⌄</span>
        </button>

        <div v-show="openedMenus.requirement" class="submenu">
          <router-link
            v-if="canAccess('/requirement/book')"
            to="/requirement/book"
            class="submenu-item"
            :class="{ active: isActive('/requirement/book') }"
          >
            需求书
          </router-link>

          <router-link
            v-if="canAccess('/requirement/change')"
            to="/requirement/change"
            class="submenu-item"
            :class="{ active: isActive('/requirement/change') }"
          >
            需求变更
          </router-link>

          <router-link
            v-if="canAccess('/requirement/customer-supplied')"
            to="/requirement/customer-supplied"  
            class="submenu-item"
            :class="{ active: isActive('/requirement/customer-supplied') }"
          >
            客供资料
          </router-link>
        </div>
      </div>

      <!-- 硬件管理 -->
      <div v-if="canAccessGroup(['/hardware/version', '/hardware/test'])" class="menu-group">
        <button class="menu-item menu-button" @click="toggleMenu('hardware')">
          <span class="menu-left">
            <span class="menu-icon">⚙</span>
            <span>硬件管理</span>
          </span>
          <span class="arrow" :class="{ open: openedMenus.hardware }">⌄</span>
        </button>

        <div v-show="openedMenus.hardware" class="submenu">
          <router-link
            v-if="canAccess('/hardware/version')"
            to="/hardware/version"
            class="submenu-item"
            :class="{ active: isActive('/hardware/version') }"
          >
            硬件版本
          </router-link>

          <router-link
            v-if="canAccess('/hardware/test')"
            to="/hardware/test"
            class="submenu-item"
            :class="{ active: isActive('/hardware/test') }"
          >
            硬件测试
          </router-link>
        </div>
      </div>

      <!-- 软件管理 -->
      <div v-if="canAccessGroup(['/software/version', '/software/branch'])" class="menu-group">
        <button class="menu-item menu-button" @click="toggleMenu('software')">
          <span class="menu-left">
            <span class="menu-icon">&lt;/&gt;</span>
            <span>软件管理</span>
          </span>
          <span class="arrow" :class="{ open: openedMenus.software }">⌄</span>
        </button>

        <div v-show="openedMenus.software" class="submenu">
          <router-link
            v-if="canAccess('/software/version')"
            to="/software/version"
            class="submenu-item"
            :class="{ active: isActive('/software/version') }"
          >
            软件版本
          </router-link>

          <router-link
            v-if="canAccess('/software/branch')"
            to="/software/branch"
            class="submenu-item"
            :class="{ active: isActive('/software/branch') }"
          >
            项目分支
          </router-link>


          <!-- <router-link
            to="/software/project-config"
            class="submenu-item"
            :class="{ active: isActive('/software/project-config') }"
          >
            项目配置资料
          </router-link> -->
        </div>
      </div>

      <!-- 测试管理 -->
      <div v-if="canAccessGroup(['/test/case', '/test/issue'])" class="menu-group">
        <button class="menu-item menu-button" @click="toggleMenu('test')">
          <span class="menu-left">
            <span class="menu-icon">◇</span>
            <span>测试管理</span>
          </span>
          <span class="arrow" :class="{ open: openedMenus.test }">⌄</span>
        </button>

        <div v-show="openedMenus.test" class="submenu">
          <router-link
            v-if="canAccess('/test/case')"
            to="/test/case"
            class="submenu-item"
            :class="{ active: isActive('/test/case') }"
          >
            测试用例
          </router-link>

          <!-- <router-link
            to="/test/report"
            class="submenu-item"
            :class="{ active: isActive('/test/report') }"
          >
            测试报告
          </router-link> -->

          <router-link
            v-if="canAccess('/test/issue')"
            to="/test/issue"
            class="submenu-item"
            :class="{ active: isActive('/test/issue') }"
          >
            问题闭环
          </router-link>
        </div>
      </div>

      <!-- 生产管理 -->
      <div v-if="canAccessGroup(['/production/burn', '/production/factory-test', '/production/inventory'])" class="menu-group">
        <button class="menu-item menu-button" @click="toggleMenu('production')">
          <span class="menu-left">
            <span class="menu-icon">▥</span>
            <span>生产管理</span>
          </span>
          <span class="arrow" :class="{ open: openedMenus.production }">⌄</span>
        </button>

        <div v-show="openedMenus.production" class="submenu">
          <router-link
            v-if="canAccess('/production/burn')"
            to="/production/burn"
            class="submenu-item"
            :class="{ active: isActive('/production/burn') }"
          >
            生产烧录
          </router-link>

          <router-link
            v-if="canAccess('/production/factory-test')"
            to="/production/factory-test"
            class="submenu-item"
            :class="{ active: isActive('/production/factory-test') }"
          >
            出厂测试
          </router-link>

          <router-link
            v-if="canAccess('/production/inventory')"
            to="/production/inventory"
            class="submenu-item"
            :class="{ active: isActive('/production/inventory') }"
          >
            库存情况
          </router-link>
        </div>
      </div>

      <!-- 发货管理 -->
      <div v-if="canAccessGroup(['/shipping/out', '/shipping/batch'])" class="menu-group">
        <button class="menu-item menu-button" @click="toggleMenu('shipping')">
          <span class="menu-left">
            <span class="menu-icon">◈</span>
            <span>发货管理</span>
          </span>
          <span class="arrow" :class="{ open: openedMenus.shipping }">⌄</span>
        </button>

        <div v-show="openedMenus.shipping" class="submenu">
          <!-- <router-link
            to="/shipping/batch"
            class="submenu-item"
            :class="{ active: isActive('/shipping/batch') }"
          >
            发货批次
          </router-link> -->


          <!-- <router-link
            to="/shipping/receipt"
            class="submenu-item"
            :class="{ active: isActive('/shipping/receipt') }"
          >
            收货记录
          </router-link> -->

          <router-link
            v-if="canAccess('/shipping/out')"
            to="/shipping/out"
            class="submenu-item"
            :class="{ active: isActive('/shipping/out') }"
          >
            出库记录
          </router-link>
            <router-link
                  v-if="canAccess('/shipping/batch')"
                  to="/shipping/batch"
                  class="submenu-item"
                  :class="{ active: isActive('/shipping/batch') }"
                >
                  发货批次
                </router-link> 
              </div>
      </div>
 
      <!-- 售后管理 -->
      <div v-if="canAccessGroup(['/aftersales/repair', '/aftersales/fault-analysis'])" class="menu-group">
        <button class="menu-item menu-button" @click="toggleMenu('aftersales')">
          <span class="menu-left">
            <span class="menu-icon">⚒</span>
            <span>售后管理</span>
          </span>
          <span class="arrow" :class="{ open: openedMenus.aftersales }">⌄</span>
        </button>

        <div v-show="openedMenus.aftersales" class="submenu">
          <!-- <router-link
            to="/aftersales/upgrade"
            class="submenu-item"
            :class="{ active: isActive('/aftersales/upgrade') }"
          >
            升级记录
          </router-link> -->

          <router-link
            v-if="canAccess('/aftersales/repair')"
            to="/aftersales/repair"
            class="submenu-item"
            :class="{ active: isActive('/aftersales/repair') }"
          >
            维修记录
          </router-link>

          <router-link
            v-if="canAccess('/aftersales/fault-analysis')"
            to="/aftersales/fault-analysis"
            class="submenu-item"
            :class="{ active: isActive('/aftersales/fault-analysis') }"
          >
            故障分析方案
          </router-link>
        </div>
      </div>

      <!-- 统计报表 -->
      <div v-if="canAccessGroup(['/report/project-progress', '/report/version-matrix', '/report/issue-statistics'])" class="menu-group">
        <button class="menu-item menu-button" @click="toggleMenu('report')">
          <span class="menu-left">
            <span class="menu-icon">▥</span>
            <span>统计报表</span>
          </span>
          <span class="arrow" :class="{ open: openedMenus.report }">⌄</span>
        </button>

        <div v-show="openedMenus.report" class="submenu">
          <router-link
            v-if="canAccess('/report/project-progress')"
            to="/report/project-progress"
            class="submenu-item"
            :class="{ active: isActive('/report/project-progress') }"
          >
            项目进度
          </router-link>

          <router-link
            v-if="canAccess('/report/version-matrix')"
            to="/report/version-matrix"
            class="submenu-item"
            :class="{ active: isActive('/report/version-matrix') }"
          >
            版本矩阵
          </router-link>

          <router-link
            v-if="canAccess('/report/issue-statistics')"
            to="/report/issue-statistics"
            class="submenu-item"
            :class="{ active: isActive('/report/issue-statistics') }"
          >
            问题统计
          </router-link>
        </div>
      </div>
    </nav>

    <!-- User -->
    <div class="sidebar-user">
      <div class="avatar">{{ userInitial }}</div>
      <div class="user-info">
        <div class="user-name">{{ currentUserName }}</div>
        <div class="user-role">{{ currentRoleName }}</div>
      </div>
      <button class="logout" type="button" title="退出登录" @click="handleLogout">
        ⇥
      </button>
    </div>
  </aside>
</template>

<script setup>
import { computed, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { canAccessPage } from '@/utils/permission'

const route = useRoute()
const router = useRouter()

const storedUser = computed(() => {
  try {
    const text = localStorage.getItem('user')
    return text ? JSON.parse(text) : {}
  } catch (err) {
    console.warn('读取用户信息失败：', err)
    return {}
  }
})

const storedRoles = computed(() => {
  try {
    const text = localStorage.getItem('roles')
    return text ? JSON.parse(text) : []
  } catch (err) {
    console.warn('读取角色信息失败：', err)
    return []
  }
})

const currentUserName = computed(() => {
  return storedUser.value.realName || storedUser.value.username || localStorage.getItem('realName') || '当前用户'
})

const currentRoleName = computed(() => {
  const firstRole = storedRoles.value[0]
  return firstRole?.roleName || firstRole?.role_name || storedUser.value.department || localStorage.getItem('department') || '系统用户'
})

const userInitial = computed(() => {
  return currentUserName.value.slice(0, 1).toUpperCase()
})

const openedMenus = reactive({
  requirement: true,
  hardware: true,
  software: true,
  test: true,
  production: true,
  shipping: true,
  aftersales: true,
  report: true
})

function toggleMenu(name) {
  openedMenus[name] = !openedMenus[name]
}

function isActive(path) {
  return route.path === path
}

function canAccess(path) {
  return canAccessPage(path)
}

function canAccessGroup(paths) {
  return paths.some(path => canAccessPage(path))
}

function handleLogout() {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  localStorage.removeItem('roles')
  localStorage.removeItem('permissions')
  localStorage.removeItem('username')
  localStorage.removeItem('accountName')
  localStorage.removeItem('realName')
  localStorage.removeItem('department')
  router.replace('/login')
}
</script>

<style scoped>
.sidebar {
  width: 240px;
  height: 100vh;
  flex-shrink: 0;
  background: #0f172a;
  border-right: 1px solid #1e293b;
  color: #cbd5e1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.main-area {
  flex: 1;
  min-width: 0;
}

.sidebar-logo {
  height: 64px;
  padding: 0 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  border-bottom: 1px solid #1e293b;
}

.logo-icon {
  width: 32px;
  height: 32px;
  background: #2563eb;
  color: #fff;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 17px;
  font-weight: 700;
}

.logo-text {
  color: #f8fafc;
  font-size: 17px;
  font-weight: 700;
  white-space: nowrap;
}

.menu {
  flex: 1;
  padding: 14px 10px;
  overflow-y: auto;
}

.menu::-webkit-scrollbar {
  width: 4px;
}

.menu::-webkit-scrollbar-track {
  background: #0f172a;
}

.menu::-webkit-scrollbar-thumb {
  background: #334155;
  border-radius: 999px;
}

.menu-group {
  margin-top: 4px;
}

.menu-item {
  width: 100%;
  min-height: 38px;
  padding: 0 12px;
  border-radius: 8px;
  color: #94a3b8;
  font-size: 14px;
  font-weight: 500;
  text-decoration: none;
  display: flex;
  align-items: center;
  gap: 10px;
  transition: 0.2s;
}

.menu-button {
  border: none;
  background: transparent;
  cursor: pointer;
  justify-content: space-between;
}

.menu-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.menu-icon {
  width: 22px;
  color: #94a3b8;
  font-size: 15px;
  text-align: center;
  flex-shrink: 0;
}

.menu-item:hover {
  color: #f8fafc;
  background: #1e293b;
}

.menu-item.active {
  color: #3b82f6;
  background: #1e293b;
  border-left: 2px solid #3b82f6;
}

.arrow {
  font-size: 12px;
  color: #64748b;
  transition: transform 0.2s;
}

.arrow.open {
  transform: rotate(180deg);
}

.submenu {
  padding-left: 44px;
  margin: 4px 0 8px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.submenu-item {
  color: #64748b;
  font-size: 13px;
  text-decoration: none;
  padding: 7px 0;
  transition: 0.2s;
}

.submenu-item:hover {
  color: #60a5fa;
}

.submenu-item.active {
  color: #3b82f6;
  font-weight: 600;
}

.sidebar-user {
  height: 72px;
  border-top: 1px solid #1e293b;
  padding: 14px 16px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.avatar {
  width: 34px;
  height: 34px;
  border-radius: 999px;
  background: #1d4ed8;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name {
  color: #f8fafc;
  font-size: 14px;
  font-weight: 600;
}

.user-role {
  color: #64748b;
  font-size: 12px;
  margin-top: 2px;
}

.logout {
  border: none;
  background: transparent;
  color: #64748b;
  cursor: pointer;
  font-size: 18px;
  padding: 4px;
}

.logout:hover {
  color: #ef4444;
}
</style>
