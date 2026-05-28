<template>
  <div class="login-wrapper">
    <div class="login-box">
      <el-form
        class="login-form"
        :model="loginForm"
        ref="loginFormRef"
      >
        <div class="header-title">
          <h3 class="title">
            <center>{{ $t('login.model') }}</center>
          </h3>
        </div>

<el-form-item prop="username">
  <el-input
    v-model="loginForm.username"
    :placeholder="$t('login.username')"
  >
    <template #prefix>
      <el-icon><User /></el-icon>
    </template>
  </el-input>
</el-form-item>

<el-form-item prop="password">
  <el-input
    v-model="loginForm.password"
    type="password"
    :placeholder="$t('login.password')"
  >
    <template #prefix>
      <el-icon><Lock /></el-icon>
    </template>
  </el-input>
</el-form-item>


        <el-button
          :loading="loading"
          type="primary"
          style="width: 100%; margin-bottom: 30px"
          @click="handleLogin"
        >
          {{ $t('login.submit') }}
        </el-button>
      </el-form>
    </div>
  </div>
</template>

<script>
import { User, Lock } from '@element-plus/icons-vue'
import { loginApi } from "@/api/auth";

export default {
  components: { User, Lock },
  data() {
    return {
      loading: false,
      loginForm: {
        username: "",
        password: ""
      }
    };
  },
  methods: {
   async handleLogin() {
  this.loading = true;

  try {
    const res = await loginApi(
      this.loginForm.username,
      this.loginForm.password
    );
                               
    console.log("login res:",res)
    const token = res.data.data.token ;
    // 登录成功（后端返回 token）
    if (token) {
      localStorage.setItem("token", token);
      this.$router.push("/");
    } else {
      // 没有后端界面的时候 显示登陆失败
      this.$message.error("登陆失败：未取到token");
    }

  } catch (err) {
    /**
     * 这里专门处理“用户名或密码错误”
     * err.response 是 axios 的错误结构
     */
    if (err.response && err.response.data && err.response.data.message) {
      this.$message.error(err.response.data.message);
    } else {
      this.$message.error("登录失败，请检查网络");
    }
  } finally {
    this.loading = false;
  }
}

  }
};
</script>


<style scoped>
/* 页面居中 */
.login-wrapper {
  position: fixed;
  inset: 0;

  display: flex;
  align-items: center;
  justify-content: center;

  background: #2c2c2c;
}

/* 白色卡片 */
.login-box {
  width: 90%;
  max-width: 360px;

  background: #fff;
  padding: 32px 32px 20px;

  border-radius: 8px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
}

/* 表单本身必须 100% */
.login-form {
  width: 100%;
}

/* 标题 */
.header-title {
  margin-bottom: 24px;
  text-align: center;
}

.title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #333;
}

/*  重点：form-item */
.login-form :deep(.el-form-item) {
  width: 100%;
  display: block;
}

/*  重点：el-input */
.login-form :deep(.el-input) {
  width: 100%;
}

/* 重点：真正的输入框容器 */
.login-form :deep(.el-input__wrapper) {
  width: 100%;
  height: 42px;
  box-sizing: border-box;
}

/* 按钮 */
.login-form :deep(.el-button) {
  width: 100%;
  height: 42px;
  font-size: 16px;
}
</style>
