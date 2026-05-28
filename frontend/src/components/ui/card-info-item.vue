<template>
  <div class="info-item">
    <div class="label">{{ title }}</div>
    
    <div class="value">
       
      <slot>
        {{ value }}
      </slot>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: String,
  value: String
})
</script>

<style scoped>
.info-item {
  display: flex;
  align-items: flex;
  padding: 10px 16px;
  font-size: 13px;
  background: #fafafa;
  box-shadow: 0 1px 0 rgba(0, 0, 0, 0.04);

  margin-bottom: 6px;
  border-radius: 4px;
}
.info-item:hover {
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.06);
}
.info-item:last-child {
  border-bottom: none;
}

.label {
  width: 140px;
  text-align: left;
  color: #606266;
}

.value {
  flex: 1;
  color: #303133;
}
</style>
