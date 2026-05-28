<template>
  <div class="card">
    <div class="card-header">
      <h3>待办事项</h3>
      <span class="add-btn">+</span>
    </div>

    <div class="todo-list">
      <div class="todo-item" v-for="item in todos" :key="item.id">
        <input type="checkbox" />

        <div class="todo-content">
          <p>{{ item.title }}</p>
          <span>截止日期：{{ item.deadline }}</span>
        </div>

        <b class="level-tag" :class="getLevelClass(item.level)">
          {{ item.level }}
        </b>
      </div>
    </div>
  </div>
</template>

<script setup>
const todos = [
  { id: 1, title: '屯马项目RPAU升级故障', deadline: '05-20', level: '紧急' },
  { id: 2, title: '提交测试问题闭环', deadline: '05-22', level: '高' },
  { id: 3, title: '波尔图需求变更查看', deadline: '05-25', level: '中' },
  { id: 4, title: '查看新项目立项申请', deadline: '05-23', level: '低' }
]

function getLevelClass(level) {
  switch (level) {
    case '紧急':
      return 'level-urgent'
    case '高':
      return 'level-high'
    case '中':
      return 'level-middle'
    case '低':
      return 'level-low'
    default:
      return 'level-low'
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
  align-items: center;
  margin-bottom: 16px;
}

.card-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 700;
}

.add-btn {
  color: #94a3b8;
  cursor: pointer;
  font-size: 16px;
}

.add-btn:hover {
  color: #f8fafc;
}

.todo-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.todo-item {
  display: flex;
  gap: 10px;
  align-items: flex-start;
  background: #1e293b80;
  border: 1px solid #1e293b;
  border-radius: 8px;
  padding: 10px;
}

.todo-item input {
  margin-top: 3px;
}

.todo-content {
  flex: 1;
  min-width: 0;
}

.todo-item p {
  margin: 0;
  font-size: 13px;
  font-weight: 600;
  color: #f8fafc;
}

.todo-item span {
  font-size: 12px;
  color: #64748b;
}

.level-tag {
  margin-left: auto;
  font-size: 11px;
  font-weight: 700;
  padding: 2px 7px;
  border-radius: 999px;
  white-space: nowrap;
}

/* 紧急：红色 */
.level-urgent {
  color: #ef4444;
  background: rgba(239, 68, 68, 0.12);
}

/* 高：蓝色 */
.level-high {
  color: #3b82f6;
  background: rgba(59, 130, 246, 0.12);
}

/* 中：黄色 */
.level-middle {
  color: #f59e0b;
  background: rgba(245, 158, 11, 0.12);
}

/* 低：灰色 */
.level-low {
  color: #94a3b8;
  background: rgba(148, 163, 184, 0.12);
}
</style>