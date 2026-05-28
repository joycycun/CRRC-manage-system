<template>
  <div class="coo-page-header">
    <!-- 左侧标题 -->
    <div class="coo-header-left">
      <router-link
        v-if="address"
        :to="address"
        class="title-link"
      >
        <h3 class="title">{{ title }}</h3>
      </router-link>

      <h3 v-else class="title">{{ title }}</h3>
    </div>

    <!-- 右侧子菜单 -->
    <div class="coo-header-right">
      <el-menu
        v-if="subMenu?.length"
        :default-active="route.path"
        class="coo-sub-menu"
        mode="horizontal"
        router
        @select="handleSelect"
      >
        <el-menu-item
          v-for="item in subMenu"
          :key="item.link"
          :index="item.link"
        >
          {{ item.name }}
        </el-menu-item>
      </el-menu>
    </div>
  </div>
</template>

<script setup>
import { useRoute } from "vue-router";

/* props（Vue3 标准写法） */
defineProps({
  title: {
    type: String,
    required: true
  },
  subMenu: {
    type: Array,
    default: () => []
  },
  address: {
    type: String,
    default: ""
  }
});

/* 当前路由 */
const route = useRoute();

/* 菜单选择 */
const handleSelect = (key, keyPath) => {
  console.log("select:", key, keyPath);
};
</script>

<style scoped>
.coo-page-header {

  justify-content: space-between;
  align-items: center;
  width: 100%;                      
  margin-bottom: 16px;

}

.title-link {
  color: #2c3e50;
  text-decoration: none;
}

.title {
  margin: 0;
}

.form-wrapper {
  max-width: 900px;
  margin: 0 ;
}
</style>


