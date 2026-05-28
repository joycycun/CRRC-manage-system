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
        to="/dashboard"
        class="menu-item"
        :class="{ active: isActive('/dashboard') }"
      >
        <span class="menu-icon">⌁</span>
        <span>项目看板</span>
      </router-link>

      <!-- 项目立项管理 -->
      <router-link
        to="/project/manage"
        class="menu-item"
        :class="{ active: isActive('/project/manage') }"
      >
        <span class="menu-icon">▣</span>
        <span>项目立项管理</span>
      </router-link>

      <!-- 需求管理 -->
      <div class="menu-group">
        <button class="menu-item menu-button" @click="toggleMenu('requirement')">
          <span class="menu-left">
            <span class="menu-icon">▤</span>
            <span>需求管理</span>
          </span>
          <span class="arrow" :class="{ open: openedMenus.requirement }">⌄</span>
        </button>

        <div v-show="openedMenus.requirement" class="submenu">
          <router-link
            to="/requirement/book"
            class="submenu-item"
            :class="{ active: isActive('/requirement/book') }"
          >
            需求书
          </router-link>

          <router-link
            to="/requirement/change"
            class="submenu-item"
            :class="{ active: isActive('/requirement/change') }"
          >
            需求变更
          </router-link>
        </div>
      </div>

      <!-- 硬件管理 -->
      <div class="menu-group">
        <button class="menu-item menu-button" @click="toggleMenu('hardware')">
          <span class="menu-left">
            <span class="menu-icon">⚙</span>
            <span>硬件管理</span>
          </span>
          <span class="arrow" :class="{ open: openedMenus.hardware }">⌄</span>
        </button>

        <div v-show="openedMenus.hardware" class="submenu">
          <router-link
            to="/hardware/version"
            class="submenu-item"
            :class="{ active: isActive('/hardware/version') }"
          >
            硬件版本
          </router-link>

          <router-link
            to="/hardware/test"
            class="submenu-item"
            :class="{ active: isActive('/hardware/test') }"
          >
            硬件测试
          </router-link>
        </div>
      </div>

      <!-- 软件管理 -->
      <div class="menu-group">
        <button class="menu-item menu-button" @click="toggleMenu('software')">
          <span class="menu-left">
            <span class="menu-icon">&lt;/&gt;</span>
            <span>软件管理</span>
          </span>
          <span class="arrow" :class="{ open: openedMenus.software }">⌄</span>
        </button>

        <div v-show="openedMenus.software" class="submenu">
          <router-link
            to="/software/version"
            class="submenu-item"
            :class="{ active: isActive('/software/version') }"
          >
            软件版本
          </router-link>

          <router-link
            to="/software/branch"
            class="submenu-item"
            :class="{ active: isActive('/software/branch') }"
          >
            项目分支
          </router-link>


          <router-link
            to="/software/project-config"
            class="submenu-item"
            :class="{ active: isActive('/software/project-config') }"
          >
            项目配置资料
          </router-link>
        </div>
      </div>

      <!-- 测试管理 -->
      <div class="menu-group">
        <button class="menu-item menu-button" @click="toggleMenu('test')">
          <span class="menu-left">
            <span class="menu-icon">◇</span>
            <span>测试管理</span>
          </span>
          <span class="arrow" :class="{ open: openedMenus.test }">⌄</span>
        </button>

        <div v-show="openedMenus.test" class="submenu">
          <router-link
            to="/test/case"
            class="submenu-item"
            :class="{ active: isActive('/test/case') }"
          >
            测试用例
          </router-link>

          <router-link
            to="/test/report"
            class="submenu-item"
            :class="{ active: isActive('/test/report') }"
          >
            测试报告
          </router-link>

          <router-link
            to="/test/issue"
            class="submenu-item"
            :class="{ active: isActive('/test/issue') }"
          >
            问题闭环
          </router-link>
        </div>
      </div>

      <!-- 生产管理 -->
      <div class="menu-group">
        <button class="menu-item menu-button" @click="toggleMenu('production')">
          <span class="menu-left">
            <span class="menu-icon">▥</span>
            <span>生产管理</span>
          </span>
          <span class="arrow" :class="{ open: openedMenus.production }">⌄</span>
        </button>

        <div v-show="openedMenus.production" class="submenu">
          <router-link
            to="/production/burn"
            class="submenu-item"
            :class="{ active: isActive('/production/burn') }"
          >
            生产烧录记录
          </router-link>

          <router-link
            to="/production/factory-test"
            class="submenu-item"
            :class="{ active: isActive('/production/factory-test') }"
          >
            出厂测试
          </router-link>

          <router-link
            to="/production/inventory"
            class="submenu-item"
            :class="{ active: isActive('/production/inventory') }"
          >
            库存情况
          </router-link>
        </div>
      </div>

      <!-- 发货管理 -->
      <div class="menu-group">
        <button class="menu-item menu-button" @click="toggleMenu('shipping')">
          <span class="menu-left">
            <span class="menu-icon">◈</span>
            <span>发货管理</span>
          </span>
          <span class="arrow" :class="{ open: openedMenus.shipping }">⌄</span>
        </button>

        <div v-show="openedMenus.shipping" class="submenu">
          <router-link
            to="/shipping/batch"
            class="submenu-item"
            :class="{ active: isActive('/shipping/batch') }"
          >
            发货批次
          </router-link>

          <router-link
            to="/shipping/receipt"
            class="submenu-item"
            :class="{ active: isActive('/shipping/receipt') }"
          >
            收货记录
          </router-link>
        </div>
      </div>

      <!-- 售后管理 -->
      <div class="menu-group">
        <button class="menu-item menu-button" @click="toggleMenu('aftersales')">
          <span class="menu-left">
            <span class="menu-icon">⚒</span>
            <span>售后管理</span>
          </span>
          <span class="arrow" :class="{ open: openedMenus.aftersales }">⌄</span>
        </button>

        <div v-show="openedMenus.aftersales" class="submenu">
          <router-link
            to="/aftersales/upgrade"
            class="submenu-item"
            :class="{ active: isActive('/aftersales/upgrade') }"
          >
            升级记录
          </router-link>

          <router-link
            to="/aftersales/repair"
            class="submenu-item"
            :class="{ active: isActive('/aftersales/repair') }"
          >
            维修记录
          </router-link>

          <router-link
            to="/aftersales/fault-analysis"
            class="submenu-item"
            :class="{ active: isActive('/aftersales/fault-analysis') }"
          >
            故障分析方案
          </router-link>
        </div>
      </div>

      <!-- 统计报表 -->
      <div class="menu-group">
        <button class="menu-item menu-button" @click="toggleMenu('report')">
          <span class="menu-left">
            <span class="menu-icon">▥</span>
            <span>统计报表</span>
          </span>
          <span class="arrow" :class="{ open: openedMenus.report }">⌄</span>
        </button>

        <div v-show="openedMenus.report" class="submenu">
          <router-link
            to="/report/project-progress"
            class="submenu-item"
            :class="{ active: isActive('/report/project-progress') }"
          >
            项目进度
          </router-link>

          <router-link
            to="/report/version-matrix"
            class="submenu-item"
            :class="{ active: isActive('/report/version-matrix') }"
          >
            版本矩阵
          </router-link>

          <router-link
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
      <div class="avatar">丁</div>
      <div class="user-info">
        <div class="user-name">丁sir</div>
        <div class="user-role">研发主管</div>
      </div>
      <div class="logout">⇥</div>
    </div>
  </aside>
</template>

<script setup>
import { reactive } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

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
  color: #64748b;
  cursor: pointer;
}

.logout:hover {
  color: #ef4444;
}
</style>