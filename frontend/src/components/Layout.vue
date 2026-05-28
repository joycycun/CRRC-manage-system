<template>
  <div class="app-layout">
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
import Sidebar from './components/Sidebar.vue'
import Header from './components/Header.vue'
</script>

<style scoped>
.app-layout {
  display: flex;
  height: 100vh;
  background: #020617;
  color: #f8fafc;
  overflow: hidden;
}

.main-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.page-content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}
</style>