<template>
  <el-container>
    <el-header>
      <el-menu
        :default-active="activeIndex"
        class="el-menu-demo"
        mode="horizontal"
        :ellipsis="true"
        @select="handleSelect"
      >
        <el-menu-item index="home">首页</el-menu-item>
        <el-menu-item index="currencyExchange">兑换货币</el-menu-item>
        <el-menu-item index="news">查看新闻</el-menu-item>
        <el-menu-item index="createArticle" v-if="authStore.isAuthenticated">发布文章</el-menu-item>
        <el-menu-item index="notifications" v-if="authStore.isAuthenticated">
          通知
          <el-badge :value="authStore.unreadCount" :hidden="authStore.unreadCount === 0" class="nav-badge" />
        </el-menu-item>
        <el-menu-item index="login" v-if="!authStore.isAuthenticated">登录</el-menu-item>
        <el-menu-item index="register" v-if="!authStore.isAuthenticated">注册</el-menu-item>
        <el-menu-item index="logout" v-if="authStore.isAuthenticated">退出</el-menu-item>
      </el-menu>
    </el-header>
    <el-main>
      <router-view></router-view>
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from './store/auth';

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();
const activeIndex = ref(route.name?.toString() || 'home');

watch(route, (newRoute) => {
  activeIndex.value = newRoute.name?.toString() || 'home';
});

let timer: ReturnType<typeof setInterval> | null = null;

onMounted(() => {
  if (authStore.isAuthenticated) {
    authStore.fetchUnreadCount();
    timer = setInterval(() => authStore.fetchUnreadCount(), 30000);
  }
});

watch(() => authStore.isAuthenticated, (val) => {
  if (val) {
    authStore.fetchUnreadCount();
    timer = setInterval(() => authStore.fetchUnreadCount(), 30000);
  } else {
    if (timer) clearInterval(timer);
  }
});

const handleSelect = (key: string) => {
  if (key === 'logout') {
    authStore.logout();
    router.push({ name: 'Home' });
  } else {
    router.push({ name: key.charAt(0).toUpperCase() + key.slice(1) });
  }
};
</script>

<style scoped>
.el-menu-demo {
  line-height: 60px;
}

.nav-badge {
  margin-left: 4px;
}
</style>
