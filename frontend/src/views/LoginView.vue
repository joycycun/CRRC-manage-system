<template>
  <div class="login-page">
    <section class="brand-panel">
      <div class="brand-mark">CR</div>

      <div class="brand-content">
        <span class="system-tag">CRRC Project Management</span>
        <h1>项目全生命周期管理平台</h1>
        <p>统一管理项目立项、研发测试、生产烧录、发货出库与售后闭环。</p>
      </div>

      <div class="module-grid">
        <div
          v-for="item in moduleItems"
          :key="item.title"
          class="module-item"
        >
          <strong>{{ item.title }}</strong>
          <span>{{ item.desc }}</span>
        </div>
      </div>
    </section>

    <section class="login-panel">
      <div class="login-card">
        <div class="login-header">
          <span>账号登录</span>
          <h2>欢迎回来</h2>
        </div>

        <form class="login-form" @submit.prevent="handleLogin">
          <label>
            <span>登录账号</span>
            <input
              v-model.trim="loginForm.username"
              type="text"
              autocomplete="username"
              placeholder="请输入用户名"
            />
          </label>

          <label>
            <span>登录密码</span>
            <input
              v-model="loginForm.password"
              type="password"
              autocomplete="current-password"
              placeholder="请输入密码"
            />
          </label>

          <div v-if="errorMessage" class="error-message">
            {{ errorMessage }}
          </div>

          <button class="login-btn" type="submit" :disabled="loading">
            {{ loading ? '登录中...' : '登录系统' }}
          </button>
        </form>

        <div class="role-preview">
          <div
            v-for="item in roleItems"
            :key="item.role"
            class="role-item"
          >
            <b>{{ item.role }}</b>
            <span>{{ item.scope }}</span>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { loginApi } from '@/api/auth'

const router = useRouter()

const loading = ref(false)
const errorMessage = ref('')

const loginForm = reactive({
  username: '',
  password: ''
})

const moduleItems = [
  { title: '项目立项', desc: '进度阶段与审核流转' },
  { title: '测试闭环', desc: '问题提交与责任人处理' },
  { title: '生产发货', desc: '烧录、测试、库存、出库' },
  { title: '售后管理', desc: '维修记录与故障分析' }
]

const roleItems = [
  { role: '领导', scope: '项目审核 / 看板统计' },
  { role: '研发', scope: '版本维护 / 问题闭环' },
  { role: '生产', scope: '烧录记录 / 出厂测试' },
  { role: '售后', scope: '维修记录 / 故障方案' }
]

async function handleLogin() {
  errorMessage.value = ''

  if (!loginForm.username) {
    errorMessage.value = '请输入登录账号'
    return
  }

  if (!loginForm.password) {
    errorMessage.value = '请输入登录密码'
    return
  }

  loading.value = true

  try {
    const res = await loginApi(loginForm.username, loginForm.password)
    const result = res?.data || res

    if (result.code !== 200 || !result.data?.token) {
      errorMessage.value = result.msg || '登录失败，请检查账号或密码'
      return
    }

    const user = result.data.user || {}
    const roles = result.data.roles || []
    const permissions = result.data.permissions || []

    localStorage.setItem('token', result.data.token)
    localStorage.setItem('user', JSON.stringify(user))
    localStorage.setItem('roles', JSON.stringify(roles))
    localStorage.setItem('permissions', JSON.stringify(permissions))
    localStorage.setItem('username', user.username || loginForm.username)
    localStorage.setItem('accountName', user.username || loginForm.username)
    localStorage.setItem('realName', user.realName || user.username || loginForm.username)
    localStorage.setItem('department', user.department || '')

    router.push('/dashboard')
  } catch (err) {
    console.error('登录失败：', err)
    errorMessage.value = err.response?.data?.msg || err.response?.data?.message || '登录失败，请检查网络或后端服务'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: grid;
  grid-template-columns: minmax(0, 1.08fr) minmax(420px, 0.92fr);
  background: #020617;
  color: #f8fafc;
}

.brand-panel {
  position: relative;
  padding: 56px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  overflow: hidden;
  border-right: 1px solid #1e293b;
  background:
    linear-gradient(135deg, rgba(37, 99, 235, 0.24), rgba(15, 23, 42, 0.08) 42%),
    #020617;
}

.brand-mark {
  width: 52px;
  height: 52px;
  border-radius: 12px;
  background: #2563eb;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 900;
  letter-spacing: 0;
  box-shadow: 0 18px 40px rgba(37, 99, 235, 0.32);
}

.brand-content {
  max-width: 640px;
}

.system-tag {
  display: inline-flex;
  color: #93c5fd;
  font-size: 13px;
  font-weight: 800;
  margin-bottom: 16px;
}

.brand-content h1 {
  margin: 0;
  font-size: 46px;
  line-height: 1.16;
  font-weight: 900;
}

.brand-content p {
  margin: 18px 0 0;
  max-width: 540px;
  color: #94a3b8;
  font-size: 15px;
  line-height: 1.8;
}

.module-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

.module-item {
  min-height: 86px;
  padding: 14px;
  border: 1px solid #1e293b;
  border-radius: 8px;
  background: rgba(15, 23, 42, 0.74);
}

.module-item strong {
  display: block;
  font-size: 13px;
  color: #e2e8f0;
}

.module-item span {
  display: block;
  margin-top: 8px;
  color: #64748b;
  font-size: 12px;
  line-height: 1.5;
}

.login-panel {
  padding: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-card {
  width: 100%;
  max-width: 430px;
  padding: 30px;
  border: 1px solid #1e293b;
  border-radius: 12px;
  background: #0f172a;
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.34);
}

.login-header {
  margin-bottom: 24px;
}

.login-header span {
  color: #60a5fa;
  font-size: 13px;
  font-weight: 800;
}

.login-header h2 {
  margin: 8px 0 0;
  color: #f8fafc;
  font-size: 28px;
  font-weight: 900;
}

.login-form {
  display: grid;
  gap: 16px;
}

.login-form label {
  display: grid;
  gap: 8px;
}

.login-form label span {
  color: #cbd5e1;
  font-size: 13px;
  font-weight: 700;
}

.login-form input {
  width: 100%;
  height: 42px;
  border: 1px solid #334155;
  border-radius: 8px;
  background: #020617;
  color: #e2e8f0;
  padding: 0 12px;
  outline: none;
  font-size: 14px;
  box-sizing: border-box;
}

.login-form input:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.2);
}

.login-form input::placeholder {
  color: #64748b;
}

.error-message {
  min-height: 36px;
  display: flex;
  align-items: center;
  padding: 0 12px;
  border: 1px solid rgba(239, 68, 68, 0.32);
  border-radius: 8px;
  background: rgba(127, 29, 29, 0.22);
  color: #fca5a5;
  font-size: 13px;
}

.login-btn {
  height: 42px;
  border: none;
  border-radius: 8px;
  background: #2563eb;
  color: #fff;
  font-size: 14px;
  font-weight: 800;
  cursor: pointer;
}

.login-btn:hover {
  background: #1d4ed8;
}

.login-btn:disabled {
  cursor: not-allowed;
  opacity: 0.72;
}

.role-preview {
  margin-top: 22px;
  padding-top: 18px;
  border-top: 1px solid #1e293b;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.role-item {
  padding: 10px;
  border-radius: 8px;
  background: #020617;
  border: 1px solid #1e293b;
}

.role-item b {
  display: block;
  color: #e2e8f0;
  font-size: 12px;
}

.role-item span {
  display: block;
  margin-top: 5px;
  color: #64748b;
  font-size: 11px;
}

@media (max-width: 1100px) {
  .login-page {
    grid-template-columns: 1fr;
  }

  .brand-panel {
    min-height: 420px;
    padding: 36px;
    border-right: none;
    border-bottom: 1px solid #1e293b;
  }
}

@media (max-width: 720px) {
  .brand-panel {
    min-height: auto;
    gap: 36px;
    padding: 26px 20px;
  }

  .brand-content h1 {
    font-size: 30px;
  }

  .module-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .login-panel {
    padding: 20px;
  }

  .login-card {
    padding: 22px;
  }
}
</style>
