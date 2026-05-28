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
const cards = [
  {
    title: '进行中项目',
    value: '18',
    desc: '较上月增加 3 个',
    icon: '📁'
  },
  {
    title: '待处理问题',
    value: '42',
    desc: '研发 / 测试 / 生产待闭环',
    icon: '⚠️'
  },
  {
    title: '本月发布版本',
    value: '12',
    desc: '软件版本发布记录',
    icon: '🚀'
  }
  // {
  //   title: '待量产终端类型',
  //   value: '6',
  //   desc: '待进入批量生产流程',
  //   icon: '🏭'
  // }
]
</script>

<style scoped>
.kpi-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
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