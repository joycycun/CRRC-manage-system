<template>
  <div class="card">
    <h3>近期软件发布版本</h3>

    <div class="release-list">
      <div class="release-item success">
        <p>u_dev_pacu_decode_bogota_v_s1.0.13-dec-bog-ts_d_20260416.img</p>
        <span>2026-05-18 14:30</span>
        <b>发布成功</b>
      </div>

      <div class="release-item running">
        <p>u_dev_pecu_dubai_v_s1.0.1-dubai_d_20260430.img</p>
        <span>正在部署...</span>
        <b>进行中</b>
      </div>

      <div class="release-item rollback">
        <p>u_dev_pacu_encode_hongkong_tm_v2_V10.01_d_20260513.img</p>
        <span>2026-05-15 09:12</span>
        <b>已回滚</b>
      </div>
    </div>
  </div>
</template>

<script setup>
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