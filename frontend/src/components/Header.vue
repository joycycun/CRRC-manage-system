<template>
  <header class="top-header">
    <div class="header-left">
      <button class="menu-btn" type="button">
        ☰
      </button>

      <div class="search-box">
        <span class="search-icon">⌕</span>
        <input
          v-model="keyword"
          type="text"
          placeholder="全局搜索项目、版本或SN..."
        />
      </div>
    </div>

    <div class="header-right">
      <button class="notify-btn" type="button">
        🔔
        <span class="notify-dot"></span>
      </button>

      <div class="divider"></div>

      <div class="date-box">
        <span>{{ todayText }}</span>
        <span class="calendar-icon">📅</span>
      </div>
    </div>
  </header>
</template>

<script setup>
import { computed, ref } from 'vue'

const keyword = ref('')

const todayText = computed(() => {
  const date = new Date()
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()
  return `${year}年${month}月${day}日`
})
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