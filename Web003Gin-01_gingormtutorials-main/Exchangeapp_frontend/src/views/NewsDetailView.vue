<template>
  <el-container>
    <el-main>
      <el-card v-if="article" class="article-detail">
        <h1>{{ article.Title }}</h1>
        <p>{{ article.Content }}</p>
        <div>
          <el-button type="primary" @click="likeArticle">点赞</el-button>
          <p>点赞数: {{ likes }}</p>
        </div>
      </el-card>
      <div v-else class="no-data">您必须登录/注册才可以阅读文章</div>

      <!-- 评论输入区 -->
      <el-card v-if="authStore.isAuthenticated" class="comment-form">
        <h3>发表评论</h3>
        <el-input
          v-model="newComment"
          type="textarea"
          placeholder="写下你的评论..."
          rows="3"
        />
        <el-button type="primary" @click="submitComment" style="margin-top: 10px">
          发表评论
        </el-button>
      </el-card>

      <!-- 评论列表 -->
      <el-card v-if="comments.length > 0" class="comment-list">
        <h3>评论 ({{ comments.length }})</h3>
        <div v-for="item in comments" :key="item.ID" class="comment-item">
          <div class="comment-header">
            <strong>{{ item.Username }}</strong>
            <span class="comment-time">{{ new Date(item.CreatedAt).toLocaleString() }}</span>
          </div>
          <p class="comment-content">{{ item.Content }}</p>
          <el-button
            v-if="item.Username === authStore.username"
            type="danger"
            size="small"
            @click="deleteComment(item.ID)"
          >
            删除
          </el-button>
        </div>
      </el-card>
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import axios from "../axios";
import { useAuthStore } from "../store/auth";
import type { Article, Like, Comment } from "../types/Article";

const article = ref<Article | null>(null);
const route = useRoute();
const likes = ref<number>(0);
const authStore = useAuthStore();

const comments = ref<Comment[]>([]);
const newComment = ref("");

const { id } = route.params;

const fetchArticle = async () => {
  try {
    const response = await axios.get<Article>(`/articles/${id}`);
    article.value = response.data;
  } catch (error) {
    console.error("Failed to load article:", error);
  }
};

const likeArticle = async () => {
  try {
    const res = await axios.post<Like>(`articles/${id}/like`)
    likes.value = res.data.likes
    await fetchLike()
  } catch (error) {
    console.log('Error Liking article:', error)
  }
};

const fetchLike = async () => {
  try {
    const res = await axios.get<Like>(`articles/${id}/like`)
    likes.value = res.data.likes
  } catch (error) {
    console.log('Error fetching likes:', error)
  }
}

const fetchComments = async () => {
  try {
    const res = await axios.get<Comment[]>(`/articles/${id}/comments`);
    comments.value = res.data;
  } catch (error) {
    console.log("Error fetching comments:", error);
  }
};

const submitComment = async () => {
  if (!newComment.value.trim()) return;
  try {
    await axios.post(`/articles/${id}/comments`, { content: newComment.value });
    newComment.value = "";
    await fetchComments();
  } catch (error) {
    console.log("Error submitting comment:", error);
  }
};

const deleteComment = async (commentId: number) => {
  try {
    await axios.delete(`/articles/${id}/comments/${commentId}`);
    await fetchComments();
  } catch (error) {
    console.log("Error deleting comment:", error);
  }
};

onMounted(() => {
  fetchArticle();
  fetchLike();
  fetchComments();
});
</script>

<style scoped>
.article-detail {
  margin: 20px 0;
}

.no-data {
  text-align: center;
  font-size: 1.2em;
  color: #999;
}

.comment-form {
  margin: 20px 0;
}

.comment-list {
  margin: 20px 0;
}

.comment-item {
  border-bottom: 1px solid #eee;
  padding: 12px 0;
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.comment-time {
  color: #999;
  font-size: 0.85em;
}

.comment-content {
  margin: 8px 0;
  line-height: 1.6;
}
</style>
