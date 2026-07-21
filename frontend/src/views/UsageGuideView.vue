<template>
  <div class="guide-page">
    <div class="page-header">
      <div>
        <h1>使用说明</h1>
        <p>集中查看和下载系统操作文档</p>
      </div>
      <button v-if="isSystemAdmin" class="primary-btn" type="button" @click="openCreateDialog">新增说明</button>
    </div>

    <div class="filter-bar">
      <input v-model.trim="keyword" placeholder="搜索说明名称或内容" />
      <span>共 {{ filteredGuides.length }} 条</span>
    </div>

    <div class="guide-list">
      <section v-for="item in filteredGuides" :key="item.id" class="guide-item">
        <button class="guide-summary" type="button" @click="toggleGuide(item.id)">
          <span class="expand-icon">{{ expandedGuideIds.includes(item.id) ? '⌄' : '›' }}</span>
          <span class="summary-main">
            <strong>{{ item.title }}</strong>
            <small>{{ item.files.length }} 个文件 · {{ item.createdBy || '管理员' }} · {{ item.createdAt }}</small>
          </span>
          <span class="summary-description">{{ item.description || '暂无补充说明' }}</span>
        </button>

        <div v-if="expandedGuideIds.includes(item.id)" class="guide-content">
          <p v-if="item.description" class="description">{{ item.description }}</p>
          <div class="file-list">
            <div v-for="(file, index) in item.files" :key="file.id || file.fileId" class="file-row">
              <span class="file-index">{{ index + 1 }}</span>
              <span class="file-info">
                <strong :title="file.fileName">{{ file.fileName }}</strong>
                <small>{{ fileTypeLabel(file) }}</small>
              </span>
              <span class="file-actions">
                <button type="button" @click="openLocalFilePreview(file)">查看</button>
                <button type="button" @click="downloadLocalFile(file, file.fileName)">下载</button>
              </span>
            </div>
          </div>
          <div v-if="isSystemAdmin" class="guide-actions">
            <button class="delete-btn" type="button" @click="removeGuide(item)">删除说明</button>
          </div>
        </div>
      </section>

      <div v-if="filteredGuides.length === 0" class="empty-state">暂无使用说明</div>
    </div>

    <div v-if="showCreateDialog" class="dialog-mask">
      <div class="dialog">
        <div class="dialog-header">
          <h3>新增使用说明</h3>
          <button type="button" title="关闭" @click="closeCreateDialog">×</button>
        </div>

        <div class="form-body">
          <label>
            说明名称
            <input v-model.trim="form.title" maxlength="128" placeholder="例如：新增硬件版本" />
          </label>
          <label>
            补充说明
            <textarea v-model.trim="form.description" rows="3" placeholder="简要说明适用场景或注意事项"></textarea>
          </label>
          <label>
            说明文件
            <span class="upload-row">
              <input ref="fileInput" class="native-file-input" type="file" accept=".doc,.docx,.pdf,.txt,.xls,.xlsx,.ppt,.pptx,.png,.jpg,.jpeg,.gif,.webp" multiple @change="handleFileChange" />
              <button class="choose-btn" type="button" @click="fileInput?.click()">选择文件</button>
              <span>{{ form.files.length ? `已选择 ${form.files.length} 个` : '支持文档、PDF、表格、演示文稿和图片' }}</span>
            </span>
          </label>

          <div v-if="form.files.length" class="selected-files">
            <div v-for="(file, index) in form.files" :key="file.fileId" class="selected-file">
              <span class="file-index">{{ index + 1 }}</span>
              <span :title="file.fileName">{{ file.fileName }}</span>
              <button type="button" title="移除文件" @click="removeSelectedFile(index)">×</button>
            </div>
          </div>

          <p v-if="formError" class="form-error">{{ formError }}</p>
        </div>

        <div class="dialog-footer">
          <button class="cancel-btn" type="button" @click="closeCreateDialog">取消</button>
          <button class="primary-btn" type="button" :disabled="saving" @click="saveGuide">
            {{ saving ? '保存中...' : '保存说明' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { createUsageGuide, deleteUsageGuide, getUsageGuides } from '@/api/usageGuide'
import { buildUploadFilePayload, downloadLocalFile, openLocalFilePreview } from '@/utils/filePreview'
import { hasRole } from '@/utils/permission'

const guides = ref([])
const keyword = ref('')
const expandedGuideIds = ref([])
const showCreateDialog = ref(false)
const saving = ref(false)
const formError = ref('')
const fileInput = ref(null)
const form = reactive({ title: '', description: '', files: [] })
const isSystemAdmin = computed(() => hasRole('system_admin'))

const filteredGuides = computed(() => {
  const value = keyword.value.toLowerCase()
  if (!value) return guides.value
  return guides.value.filter(item =>
    item.title.toLowerCase().includes(value) || item.description.toLowerCase().includes(value)
  )
})

onMounted(loadGuides)

async function loadGuides() {
  try {
    const result = (await getUsageGuides())?.data
    if (result.code !== 200) throw new Error(result.msg || '查询失败')
    guides.value = (result.data || []).map(item => ({
      ...item,
      title: item.title || '',
      description: item.description || '',
      files: item.files || item.images || []
    }))
  } catch (err) {
    alert(err.response?.data || err.message || '加载使用说明失败')
  }
}

function toggleGuide(id) {
  expandedGuideIds.value = expandedGuideIds.value.includes(id)
    ? expandedGuideIds.value.filter(itemId => itemId !== id)
    : [...expandedGuideIds.value, id]
}

function openCreateDialog() {
  form.title = ''
  form.description = ''
  form.files = []
  formError.value = ''
  showCreateDialog.value = true
}

function closeCreateDialog() {
  if (!saving.value) showCreateDialog.value = false
}

async function handleFileChange(event) {
  const files = Array.from(event.target.files || [])
  if (!files.length) return
  try {
    const payloads = await Promise.all(files.map(file => buildUploadFilePayload(file)))
    form.files = [...form.files, ...payloads]
  } catch (err) {
    formError.value = '读取文件失败，请重新选择'
  } finally {
    event.target.value = ''
  }
}

function removeSelectedFile(index) {
  form.files.splice(index, 1)
}

function fileTypeLabel(file) {
  const name = String(file.fileName || '')
  const extension = name.includes('.') ? name.split('.').pop().toUpperCase() : '文件'
  return `${extension} 文件`
}

async function saveGuide() {
  formError.value = ''
  if (!form.title) {
    formError.value = '请输入说明名称'
    return
  }
  if (!form.files.length) {
    formError.value = '请至少选择一个说明文件'
    return
  }
  try {
    saving.value = true
    const result = (await createUsageGuide({
      title: form.title,
      description: form.description,
      files: form.files
    }))?.data
    if (result.code !== 200) throw new Error(result.msg || '保存失败')
    showCreateDialog.value = false
    await loadGuides()
  } catch (err) {
    formError.value = err.response?.data || err.message || '保存使用说明失败'
  } finally {
    saving.value = false
  }
}

async function removeGuide(item) {
  if (!confirm(`确认删除使用说明【${item.title}】吗？`)) return
  try {
    const result = (await deleteUsageGuide(item.id))?.data
    if (result.code !== 200) throw new Error(result.msg || '删除失败')
    await loadGuides()
  } catch (err) {
    alert(err.response?.data || err.message || '删除使用说明失败')
  }
}
</script>

<style scoped>
.guide-page { display: flex; flex-direction: column; gap: 18px; color: #e5e7eb; }
.page-header { display: flex; align-items: center; justify-content: space-between; gap: 16px; }
.page-header h1 { margin: 0 0 6px; font-size: 25px; letter-spacing: 0; }
.page-header p { margin: 0; color: #94a3b8; font-size: 14px; }
.primary-btn, .choose-btn, .cancel-btn, .delete-btn { min-height: 38px; padding: 0 16px; border: 0; border-radius: 6px; cursor: pointer; }
.primary-btn { background: #2563eb; color: #fff; }
.primary-btn:disabled { opacity: .6; cursor: wait; }
.filter-bar { display: flex; align-items: center; justify-content: space-between; gap: 16px; padding: 14px 16px; background: #0f172a; border: 1px solid #1e293b; border-radius: 8px; }
.filter-bar input { width: min(420px, 100%); }
.filter-bar span { color: #94a3b8; white-space: nowrap; }
input, textarea { box-sizing: border-box; width: 100%; color: #e5e7eb; background: #111827; border: 1px solid #334155; border-radius: 6px; padding: 10px 12px; outline: none; }
input:focus, textarea:focus { border-color: #3b82f6; }
.guide-list { display: flex; flex-direction: column; gap: 10px; }
.guide-item { overflow: hidden; background: #0f172a; border: 1px solid #1e293b; border-radius: 8px; }
.guide-summary { width: 100%; min-height: 76px; display: grid; grid-template-columns: 26px minmax(220px, 1fr) minmax(220px, 1.5fr); align-items: center; gap: 12px; padding: 14px 18px; color: inherit; background: transparent; border: 0; text-align: left; cursor: pointer; }
.guide-summary:hover { background: #111c31; }
.expand-icon { color: #60a5fa; font-size: 24px; line-height: 1; }
.summary-main { display: flex; flex-direction: column; gap: 6px; min-width: 0; }
.summary-main strong { font-size: 15px; }
.summary-main small, .summary-description { color: #94a3b8; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.guide-content { padding: 20px 24px 24px 56px; border-top: 1px solid #1e293b; background: #0b1220; }
.description { margin: 0 0 18px; color: #cbd5e1; white-space: pre-wrap; }
.file-list { overflow: hidden; border: 1px solid #263449; border-radius: 6px; }
.file-row { min-height: 68px; display: grid; grid-template-columns: 32px minmax(0, 1fr) auto; align-items: center; gap: 12px; padding: 10px 14px; background: #0f172a; border-bottom: 1px solid #263449; }
.file-row:last-child { border-bottom: 0; }
.file-index { display: grid; place-items: center; width: 28px; height: 28px; color: #bfdbfe; background: #1e3a8a; border-radius: 50%; }
.file-info { min-width: 0; display: flex; flex-direction: column; gap: 5px; }
.file-info strong { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.file-info small { color: #64748b; }
.file-actions { display: flex; gap: 8px; }
.file-actions button { min-height: 34px; padding: 0 13px; color: #dbeafe; background: #1e3a8a; border: 0; border-radius: 5px; cursor: pointer; }
.guide-actions { display: flex; justify-content: flex-end; margin-top: 18px; }
.delete-btn { color: #fecaca; background: #7f1d1d; }
.empty-state { padding: 54px; color: #64748b; text-align: center; background: #0f172a; border: 1px dashed #334155; border-radius: 8px; }
.dialog-mask { position: fixed; inset: 0; z-index: 1000; display: grid; place-items: center; padding: 20px; background: rgba(2, 6, 23, .72); }
.dialog { width: min(720px, 100%); max-height: calc(100vh - 40px); overflow-y: auto; background: #0f172a; border: 1px solid #334155; border-radius: 8px; box-shadow: 0 20px 60px rgba(0, 0, 0, .45); }
.dialog-header, .dialog-footer { display: flex; align-items: center; justify-content: space-between; padding: 16px 20px; }
.dialog-header { border-bottom: 1px solid #1e293b; }
.dialog-header h3 { margin: 0; }
.dialog-header button { width: 34px; height: 34px; color: #cbd5e1; background: transparent; border: 0; font-size: 24px; cursor: pointer; }
.form-body { display: flex; flex-direction: column; gap: 16px; padding: 20px; }
.form-body label { display: flex; flex-direction: column; gap: 8px; color: #cbd5e1; }
.upload-row { display: flex; align-items: center; gap: 12px; color: #94a3b8; }
.native-file-input { display: none; }
.choose-btn { color: #dbeafe; background: #1d4ed8; white-space: nowrap; }
.selected-files { display: flex; flex-direction: column; gap: 8px; }
.selected-file { min-width: 0; display: grid; grid-template-columns: 32px minmax(0, 1fr) 30px; align-items: center; gap: 10px; padding: 9px 10px; background: #111827; border: 1px solid #263449; border-radius: 6px; }
.selected-file > span:nth-child(2) { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.selected-file button { color: #fca5a5; background: transparent; border: 0; font-size: 20px; cursor: pointer; }
.form-error { margin: 0; color: #fca5a5; }
.dialog-footer { justify-content: flex-end; gap: 10px; border-top: 1px solid #1e293b; }
.cancel-btn { color: #cbd5e1; background: #334155; }
@media (max-width: 760px) {
  .guide-summary { grid-template-columns: 24px minmax(0, 1fr); }
  .summary-description { display: none; }
  .guide-content { padding: 16px; }
  .file-row { grid-template-columns: 32px minmax(0, 1fr); }
  .file-actions { grid-column: 2; }
}
</style>
