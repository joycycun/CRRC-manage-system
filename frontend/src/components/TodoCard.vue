<template>
  <div class="card">
    <div class="card-header">
      <h3>待办事项</h3>
    </div>

    <div class="todo-list">
      <button
        class="todo-item"
        v-for="item in todos"
        :key="item.id"
        type="button"
        @click="goTodo(item)"
      >
        <div class="todo-content">
          <p>{{ item.title }}</p>
          <span>截止日期：{{ item.deadline }}</span>
        </div>

        <b class="level-tag" :class="getLevelClass(item.level)">
          {{ item.level }}
        </b>
      </button>

      <div v-if="todos.length === 0" class="empty-state">
        暂无待办事项
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getDashboardSummary } from '@/api/report'
import { getCurrentUserParams } from '@/utils/currentUser'

const router = useRouter()
const todos = ref([])

onMounted(() => {
  loadTodos()
})

async function loadTodos() {
  try {
    const res = await getDashboardSummary(getCurrentUserParams())
    const result = res?.data || res
    if (result.code !== 200) return
    todos.value = result.data?.todos || []
  } catch (err) {
    console.error('加载待办事项失败：', err)
  }
}

function goTodo(item) {
  if (item.link) {
    router.push(item.link)
  }
}

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

.todo-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.todo-item {
  width: 100%;
  display: flex;
  gap: 10px;
  align-items: flex-start;
  background: #1e293b80;
  border: 1px solid #1e293b;
  border-radius: 8px;
  padding: 10px;
  text-align: left;
  cursor: pointer;
}

.todo-item:hover {
  border-color: #334155;
  background: #1e293bcc;
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

.empty-state {
  padding: 24px 0;
  color: #64748b;
  font-size: 13px;
  text-align: center;
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
