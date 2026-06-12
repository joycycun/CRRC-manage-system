<template>
  <div class="card">
    <div class="card-header">
      <h3>项目进度概览</h3>
        <button @click="$router.push('/project/progress-report')">
          查看全部
        </button>
    </div>

    <div class="project-list">
      <div class="project-item" v-for="item in projects" :key="item.id">
        <div class="project-title">
          <span>{{ item.name }}</span>
          <span>{{ item.owner }} | <b>{{ item.progress }}%</b></span>
        </div>
        <div class="progress-bg">
          <div class="progress-bar" :style="{ width: item.progress + '%' }"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { getProjectProgressReport } from '@/api/report'

const projects = ref([])

onMounted(() => {
  loadProjects()
})

async function loadProjects() {
  try {
    const res = await getProjectProgressReport()
    const result = res?.data || res
    if (result.code !== 200) return
    projects.value = (result.data || []).slice(0, 5).map(item => ({
      id: item.id || item.projectId,
      name: item.name || item.projectName || '',
      owner: item.owner || '-',
      progress: Math.round(Number(item.progress || 0))
    }))
  } catch (err) {
    console.error('加载项目进度概览失败：', err)
  }
}
</script>

<style scoped>
.card {
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 12px;
  padding: 20px;
  color: #f8fafc;
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
}

.card-header button {
  background: transparent;
  border: none;
  color: #3b82f6;
  cursor: pointer;
}

.project-list {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.project-title {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
  margin-bottom: 8px;
}

.project-title b {
  color: #3b82f6;
}

.progress-bg {
  width: 100%;
  height: 8px;
  background: #1e293b;
  border-radius: 999px;
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  background: #3b82f6;
  border-radius: 999px;
}
</style>
