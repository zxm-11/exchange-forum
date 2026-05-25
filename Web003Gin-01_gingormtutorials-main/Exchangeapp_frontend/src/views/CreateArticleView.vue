<template>
  <el-container>
    <el-main>
      <el-card class="create-article-card">
        <template #header>
          <h2>发布文章</h2>
        </template>
        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-width="80px"
        >
          <el-form-item label="标题" prop="Title">
            <el-input v-model="form.Title" placeholder="请输入文章标题" />
          </el-form-item>
          <el-form-item label="预览" prop="Preview">
            <el-input
              v-model="form.Preview"
              type="textarea"
              :rows="3"
              placeholder="请输入文章预览摘要"
            />
          </el-form-item>
          <el-form-item label="正文" prop="Content">
            <el-input
              v-model="form.Content"
              type="textarea"
              :rows="10"
              placeholder="请输入文章正文内容"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="submitArticle" :loading="submitting">
              发布
            </el-button>
            <el-button @click="router.push({ name: 'News' })">取消</el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import type { FormInstance, FormRules } from "element-plus";
import axios from "../axios";
import { useAuthStore } from "../store/auth";

const router = useRouter();
const authStore = useAuthStore();
const formRef = ref<FormInstance>();
const submitting = ref(false);

const form = reactive({
  Title: "",
  Preview: "",
  Content: "",
});

const rules: FormRules = {
  Title: [{ required: true, message: "请输入文章标题", trigger: "blur" }],
  Preview: [{ required: true, message: "请输入文章预览", trigger: "blur" }],
  Content: [{ required: true, message: "请输入文章正文", trigger: "blur" }],
};

const submitArticle = async () => {
  if (!authStore.isAuthenticated) {
    ElMessage.error("请先登录后再发布文章");
    return;
  }

  const valid = await formRef.value?.validate().catch(() => false);
  if (!valid) return;

  submitting.value = true;
  try {
    await axios.post("/articles", form);
    ElMessage.success("文章发布成功");
    router.push({ name: "News" });
  } catch (error) {
    ElMessage.error("文章发布失败，请重试");
    console.error("Failed to create article:", error);
  } finally {
    submitting.value = false;
  }
};
</script>

<style scoped>
.create-article-card {
  max-width: 800px;
  margin: 0 auto;
}
</style>
