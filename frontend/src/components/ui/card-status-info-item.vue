<template>
  <div class="info-item">
    <div class="label">{{ title }}</div>
    <div class="value">{{ value }}</div>
    <div class="status" :class="status">
      {{ statusText }}
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  title: String,
  value: String,
  status: String
})

const statusText = computed(() => {
  return props.status === 'inactive' ? '未激活' : '已注册'
})
</script>

<style scoped>
.info-item {
  display: flex;
  align-items: center;
  padding: 10px 16px;
  font-size: 13px;
  background: #fafafa;
  box-shadow: 0 1px 0 rgba(0, 0, 0, 0.04);

  margin-bottom: 6px;
  border-radius: 4px;;
}
.info-item:hover {
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.06);
}
.label {
  width: 140px;
  color: #606266;
  text-align: left;
}

.value {
  flex: 1;
  color: #303133;
}

.status {
  width: 80px;
  text-align: right;
  font-size: 12px;
}

.status.inactive {
  color: #909399;
}

.status.active {
  color: #67c23a;
}
</style>
