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
const projects = [
  { id: 1, name: '香港屯马项目', owner: '研发：卢进', progress: 85 },
  { id: 2, name: '波尔图二期项目', owner: '测试：寸诗睿', progress: 62 },
  { id: 3, name: '香港东涌线', owner: '售后：卢进', progress: 92 },
  { id: 4, name: '阿根廷项目', owner: '硬件研发：王宇', progress: 45 },
  { id: 5, name: '波哥大有轨项目', owner: '需求：寸诗睿', progress: 15 }
]
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