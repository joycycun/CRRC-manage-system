<template>
  <div class="kpi-grid">
    <div class="kpi-card" v-for="item in cards" :key="item.title">
      <div class="kpi-icon">
        {{ item.icon }}
      </div>

      <div class="kpi-content">
        <p>{{ item.title }}</p>
        <h2>{{ item.value }}</h2>
        <span>{{ item.desc }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { getDashboardSummary } from '@/api/report'
import { getCurrentUserParams } from '@/utils/currentUser'

const cards = ref([
  { title: '进行中项目', value: 0, desc: '较上月持平', icon: '项' },
  { title: '待处理问题', value: 0, desc: '当前负责人未闭环问题', icon: '问' },
  { title: '本月发布版本', value: 0, desc: '较上月持平', icon: '版' }
])

onMounted(() => {
  loadKpis()
})

async function loadKpis() {
  try {
    const res = await getDashboardSummary(getCurrentUserParams())
    const result = res?.data || res
    if (result.code !== 200) return
    cards.value = result.data?.kpis || cards.value
  } catch (err) {
    console.error('加载看板KPI失败：', err)
  }
}
</script>

<style scoped>
.kpi-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

.kpi-card {
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 14px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.kpi-icon {
  width: 46px;
  height: 46px;
  border-radius: 12px;
  background: #1e293b;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
}

.kpi-content p {
  margin: 0;
  color: #94a3b8;
  font-size: 13px;
}

.kpi-content h2 {
  margin: 6px 0;
  font-size: 28px;
  font-weight: 800;
  color: #f8fafc;
}

.kpi-content span {
  color: #64748b;
  font-size: 12px;
}

@media (max-width: 1280px) {
  .kpi-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .kpi-grid {
    grid-template-columns: 1fr;
  }
}
</style>
