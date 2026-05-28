<template>
  <div class="version-card">
    <div class="card-header">
      <div>
        <h3>版本状态矩阵</h3>
        <!-- <p>展示项目对应的硬件版本、软件版本、状态和更新时间</p> -->
      </div>

      <div class="header-actions">
        <button class="secondary-btn">导出 Excel</button>
        <button class="primary-btn" @click="$router.push('/version/matrix')">
          筛选
        </button>
      </div>
    </div>

    <div class="table-wrapper">
      <table class="version-table">
        <thead>
          <tr>
            <th>项目名称</th>
            <th>终端类型</th>
            <th>硬件版本</th>
            <th>软件版本</th>
            <th>当前状态</th>
            <th>最近更新</th>
            <th>负责人</th>
            <th class="operation-col">操作</th>
          </tr>
        </thead>

        <tbody>
          <tr v-for="item in versionList" :key="item.id">
            <td class="project-name">
              {{ item.projectName }}
            </td>

            <td>
              {{ item.deviceType }}
            </td>

            <td>
              <span class="version-tag hardware">
                {{ item.hardwareVersion }}
              </span>
            </td>

            <td>
              <span class="version-tag software">
                {{ item.softwareVersion }}
              </span>
            </td>

            <td>
              <span class="status-tag" :class="item.statusType">
                {{ item.status }}
              </span>
            </td>

            <td class="muted">
              {{ item.updateTime }}
            </td>

            <td>
              {{ item.owner }}
            </td>

            <td class="operation-col">
              <button class="text-btn" @click="$router.push('/version/matrix')">
                查看
              </button>
              <button class="more-btn">...</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="table-footer">
      <span>共 {{ versionList.length }} 条记录</span>
      <!-- <span>数据来源：软件版本、硬件版本、生产烧录记录、发货记录</span> -->
    </div>
  </div>
</template>

<script setup>
const versionList = [
  {
    id: 1,
    projectName: '波尔图二期项目',
    deviceType: '广播控制盒',
    hardwareVersion: 'HW_V2.1.0',
    softwareVersion: 'SW_V1.2.4-Release',
    status: '硬件开发中',
    statusType: 'developing',
    updateTime: '2026-05-18',
    owner: '卢进'
  },
  {
    id: 2,
    projectName: '阿根廷有轨项目',
    deviceType: '客室解码板',
    hardwareVersion: 'HW_V1.0.2',
    softwareVersion: 'SW_V2.0.1-Beta',
    status: '集成测试',
    statusType: 'testing',
    updateTime: '2026-05-19',
    owner: '寸诗睿'
  },
  {
    id: 3,
    projectName: '屯马报警器项目',
    deviceType: '乘客报警器',
    hardwareVersion: 'HW_V3.0.0',
    softwareVersion: 'SW_V3.0.0-RC1',
    status: '版本冻结',
    statusType: 'frozen',
    updateTime: '2026-05-12',
    owner: '王宇'
  },
  {
    id: 4,
    projectName: '波哥大有轨项目',
    deviceType: '编码板',
    hardwareVersion: 'HW_V1.4.0',
    softwareVersion: 'SW_V1.0.13',
    status: '待量产',
    statusType: 'production',
    updateTime: '2026-05-16',
    owner: '丁sir'
  }
]
</script>

<style scoped>
.version-card {
  width: 100%;
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 14px;
  color: #f8fafc;
  overflow: hidden;
}

.card-header {
  padding: 20px 22px;
  border-bottom: 1px solid #1e293b;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
}

.card-header p {
  margin: 6px 0 0;
  font-size: 13px;
  color: #94a3b8;
}

.header-actions {
  display: flex;
  gap: 10px;
  flex-shrink: 0;
}

.primary-btn,
.secondary-btn {
  height: 32px;
  padding: 0 12px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
}

.primary-btn {
  border: none;
  background: #2563eb;
  color: #fff;
}

.secondary-btn {
  border: 1px solid #334155;
  background: #1e293b;
  color: #cbd5e1;
}

.table-wrapper {
  width: 100%;
  overflow-x: auto;
  overflow-y: hidden;
}

/* Chrome / Edge / Safari 滚动条整体 */
.table-wrapper::-webkit-scrollbar {
  height: 10px;
}

/* 滚动条轨道，也就是底色 */
.table-wrapper::-webkit-scrollbar-track {
  background: #020617;
  border-radius: 999px;
}

/* 可以拖动的滑块 */
.table-wrapper::-webkit-scrollbar-thumb {
  background: #334155;
  border-radius: 999px;
  border: 2px solid #020617;
}

/* 鼠标放上去时，滑块稍微亮一点 */
.table-wrapper::-webkit-scrollbar-thumb:hover {
  background: #475569;
}

/* 左右两边的小箭头区域 */
.table-wrapper::-webkit-scrollbar-button {
  display: none;
}

.version-table {
  width: 100%;
  min-width: 980px;
  border-collapse: collapse;
  table-layout: fixed;
}

.version-table thead {
  background: #020617;
}

.version-table th {
  padding: 14px 16px;
  color: #94a3b8;
  font-size: 12px;
  font-weight: 600;
  text-align: left;
  border-bottom: 1px solid #1e293b;
  white-space: nowrap;
}

.version-table td {
  padding: 15px 16px;
  font-size: 13px;
  color: #e2e8f0;
  border-bottom: 1px solid #1e293b;
  vertical-align: middle;
  white-space: nowrap;
}

.version-table tbody tr {
  transition: background 0.2s;
}

.version-table tbody tr:hover {
  background: #1e293b80;
}

.project-name {
  color: #f8fafc;
  font-weight: 700;
}

.version-tag {
  display: inline-flex;
  align-items: center;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
}

.version-tag.hardware {
  background: #312e81;
  color: #c4b5fd;
}

.version-tag.software {
  background: #0f766e33;
  color: #5eead4;
}

.status-tag {
  display: inline-flex;
  align-items: center;
  padding: 4px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
}

.status-tag.developing {
  background: #1d4ed833;
  color: #60a5fa;
}

.status-tag.testing {
  background: #d9770633;
  color: #fbbf24;
}

.status-tag.frozen {
  background: #16a34a33;
  color: #4ade80;
}

.status-tag.production {
  background: #9333ea33;
  color: #c084fc;
}

.muted {
  color: #94a3b8 !important;
}

.operation-col {
  width: 120px;
  text-align: right !important;
}

.text-btn {
  border: none;
  background: transparent;
  color: #60a5fa;
  font-size: 13px;
  cursor: pointer;
  margin-right: 10px;
}

.more-btn {
  width: 26px;
  height: 26px;
  border: 1px solid #334155;
  background: #1e293b;
  color: #cbd5e1;
  border-radius: 6px;
  cursor: pointer;
}

.table-footer {
  padding: 12px 22px;
  border-top: 1px solid #1e293b;
  display: flex;
  justify-content: space-between;
  color: #64748b;
  font-size: 12px;
}

/* 小屏幕适配 */
@media (max-width: 768px) {
  .card-header {
    flex-direction: column;
  }

  .header-actions {
    width: 100%;
  }

  .primary-btn,
  .secondary-btn {
    flex: 1;
  }

  .table-footer {
    flex-direction: column;
    gap: 6px;
  }
}
</style>