<template>
  <div class="card">
    <h3>近期软件发布版本</h3>

    <div class="release-list">
      <div
        class="release-item"
        v-for="item in releaseList"
        :key="item.id"
        :class="item.statusType"
      >
        <p>{{ item.version }}</p>
        <span>{{ item.releaseTime || '暂无发布时间' }}</span>
        <b>{{ item.status }}</b>
      </div>

      <div v-if="releaseList.length === 0" class="empty-state">
        暂无软件发布记录
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { getDashboardSummary } from '@/api/report'
import { getCurrentUserParams } from '@/utils/currentUser'

const releaseList = ref([])

onMounted(() => {
  loadRecentReleases()
})

async function loadRecentReleases() {
  try {
    const res = await getDashboardSummary(getCurrentUserParams())
    const result = res?.data || res
    if (result.code !== 200) return
    releaseList.value = result.data?.recentReleases || []
  } catch (err) {
    console.error('加载近期软件发布版本失败：', err)
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

h3 {
  margin: 0 0 16px;
  font-size: 16px;
}

.release-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.release-item {
  padding: 10px 12px;
  border-left: 4px solid #334155;
  background: #1e293b40;
}

.release-item p {
  margin: 0;
  font-size: 13px;
  font-weight: 700;
  word-break: break-all;
}

.release-item span {
  display: block;
  margin-top: 4px;
  color: #64748b;
  font-size: 12px;
}

.release-item b {
  display: inline-block;
  margin-top: 6px;
  font-size: 11px;
}

.empty-state {
  padding: 24px 0;
  color: #64748b;
  font-size: 13px;
  text-align: center;
}

.success {
  border-left-color: #10b981;
}

.success b {
  color: #10b981;
}

.running {
  border-left-color: #3b82f6;
}

.running b {
  color: #3b82f6;
}

.rollback {
  border-left-color: #ef4444;
}

.rollback b {
  color: #ef4444;
}
</style>
