<template>
  <router-view v-if="isLoginPage" />

  <div v-else class="app-layout">
    <Sidebar />

    <div class="main-area">
      <Header />

      <main class="page-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import Sidebar from './components/Sidebar.vue'
import Header from './components/Header.vue'

const route = useRoute()

const isLoginPage = computed(() => route.path === '/login')
</script>

<style scoped>
.app-layout {
  width: 100vw;
  min-height: 100vh;
  display: flex;
  background: #020617;
  color: #f8fafc;
  overflow: hidden;
}

.main-area {
  flex: 1;
  min-width: 0;
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.page-content {
  flex: 1;
  min-width: 0;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 24px;
  background: #020617;
}

/* 平板和小屏幕 */
@media (max-width: 1024px) {
  .page-content {
    padding: 16px;
  }
}

/* 手机屏幕 */
@media (max-width: 768px) {
  .app-layout {
    flex-direction: column;
  }

  .main-area {
    height: auto;
    min-height: 100vh;
  }

  .page-content {
    padding: 12px;
  }
}
</style>
