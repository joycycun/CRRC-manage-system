<template>
  <div class="dashboard-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div>
        <h1>项目管理看板</h1>
        <!-- <p>查看项目进度、待处理问题、软件版本发布和生产动态</p> -->
      </div>

      <!-- <button class="primary-btn">
        新建项目
      </button> -->
    </div>

    <!-- 第一行：总览卡片 -->
    <KpiCards />

    <!-- 第二行：项目进度 + 待办事项 + 近期版本 -->
    <div class="dashboard-grid three-columns">
      <ProjectProgressCard />
      <TodoCard />
      <RecentReleaseCard />
    </div>

    <!-- 第三行：问题图表 + 版本矩阵 -->
    <div class="dashboard-grid two-columns">
      <IssueChart />
      <VersionTable />
    </div>
  </div>
</template>

<script setup>
import KpiCards from '@/components/KpiCards.vue'
import ProjectProgressCard from '@/components/ProjectProgressCard.vue'
import TodoCard from '@/components/TodoCard.vue'
import RecentReleaseCard from '@/components/RecentReleaseCard.vue'
import IssueChart from '@/components/IssueChart.vue'
import VersionTable from '@/components/VersionTable.vue'
</script>

<style scoped>
.dashboard-page {
  width: 100%;
  min-height: 100%;
  color: #f8fafc;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
}

.page-header h1 {
  margin: 0;
  font-size: 28px;
  font-weight: 800;
  color: #f8fafc;
}

.page-header p {
  margin: 8px 0 0;
  color: #94a3b8;
  font-size: 14px;
}

.primary-btn {
  height: 38px;
  padding: 0 18px;
  border: none;
  border-radius: 8px;
  background: #2563eb;
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
}

.primary-btn:hover {
  background: #1d4ed8;
}

.dashboard-grid {
  display: grid;
  gap: 20px;
  margin-top: 20px;
}

.three-columns {
  grid-template-columns: 1.2fr 1fr 1.2fr;
}

.chart-row {
  grid-template-columns: 1fr;
  align-items: stretch;
}

.two-columns {
  grid-template-columns: 360px minmax(0, 1fr);
}

/* 中等屏幕 */
@media (max-width: 1280px) {
  .three-columns {
    grid-template-columns: 1fr;
  }

  .two-columns {
    grid-template-columns: 1fr;
  }
}

/* 小屏幕 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
  }

  .page-header h1 {
    font-size: 22px;
  }

  .primary-btn {
    width: 100%;
  }
}
</style>