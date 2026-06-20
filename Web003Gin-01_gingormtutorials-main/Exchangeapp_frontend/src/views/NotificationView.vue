<template>
  <el-container>
    <el-main class="notification-container">
      <h2>消息通知</h2>
      <div v-if="notifications.length === 0" class="no-data">暂无通知</div>
      <el-card
        v-for="item in notifications"
        :key="item.ID"
        :class="['notification-item', { unread: !item.is_read }]"
        @click="goToArticle(item)"
      >
        <div class="notification-header">
          <strong>{{ item.title }}</strong>
          <span class="notification-time">
            {{ new Date(item.CreatedAt).toLocaleString() }}
          </span>
        </div>
        <p class="notification-content">{{ item.content }}</p>
        <el-tag v-if="!item.is_read" type="danger" size="small">未读</el-tag>
        <el-tag v-else type="info" size="small">已读</el-tag>
      </el-card>
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import axios from "../axios";
import { useAuthStore } from "../store/auth";
import type { Notification } from "../types/Article";

const router = useRouter();
const authStore = useAuthStore();
const notifications = ref<Notification[]>([]);

const fetchNotifications = async () => {
  try {
    const res = await axios.get<Notification[]>("/notifications");
    notifications.value = res.data;
  } catch (error) {
    console.error("Error fetching notifications:", error);
  }
};

const goToArticle = async (item: Notification) => {
  if (!item.is_read) {
    try {
      await axios.put(`/notifications/${item.ID}/read`);
      item.is_read = true;
      await authStore.fetchUnreadCount();
    } catch (error) {
      console.error("Error marking notification as read:", error);
    }
  }
  // 从 link 中提取文章ID，跳转到文章详情
  const articleId = item.link.split("/").pop();
  if (articleId) {
    router.push({ name: "NewsDetail", params: { id: articleId } });
  }
};

onMounted(fetchNotifications);
</script>

<style scoped>
.notification-container {
  max-width: 700px;
  margin: 0 auto;
}

.no-data {
  text-align: center;
  color: #999;
  margin-top: 40px;
}

.notification-item {
  margin: 12px 0;
  cursor: pointer;
}

.notification-item.unread {
  background-color: #f0f9ff;
  border-left: 3px solid #409eff;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.notification-time {
  color: #999;
  font-size: 0.85em;
}

.notification-content {
  margin: 8px 0;
  color: #666;
}
</style>
