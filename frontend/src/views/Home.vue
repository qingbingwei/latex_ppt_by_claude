<template>
  <div class="home">
    <Header />
    <el-main class="home-content">
      <div class="hero">
        <h1>AI-Powered LaTeX PPT Generator</h1>
        <p>Generate professional LaTeX Beamer presentations with AI assistance</p>
        <div class="actions">
          <el-button type="primary" size="large" @click="$router.push('/generate')">
            Get Started
          </el-button>
          <el-button size="large" @click="$router.push('/knowledge')">
            Manage Knowledge Base
          </el-button>
        </div>
      </div>

      <div class="features">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-card shadow="hover">
              <template #header>
                <el-icon :size="40" color="#409eff"><Document /></el-icon>
              </template>
              <h3>RAG Support</h3>
              <p>Upload documents to build a knowledge base and generate PPTs with relevant context</p>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card shadow="hover">
              <template #header>
                <el-icon :size="40" color="#67c23a"><Edit /></el-icon>
              </template>
              <h3>AI Generation</h3>
              <p>Powered by OpenAI and Claude to generate high-quality LaTeX Beamer presentations</p>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card shadow="hover">
              <template #header>
                <el-icon :size="40" color="#e6a23c"><Files /></el-icon>
              </template>
              <h3>PDF Export</h3>
              <p>Automatically compile LaTeX to PDF for instant preview and download</p>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </el-main>
  </div>
</template>

<script setup lang="ts">
import Header from '@/components/common/Header.vue'
</script>

<style scoped>
.home {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.home-content {
  flex: 1;
  padding: 40px 20px;
}

.hero {
  text-align: center;
  padding: 60px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 10px;
  color: white;
  margin-bottom: 40px;
}

.hero h1 {
  font-size: 42px;
  margin-bottom: 20px;
}

.hero p {
  font-size: 20px;
  margin-bottom: 30px;
  opacity: 0.9;
}

.actions {
  display: flex;
  gap: 15px;
  justify-content: center;
}

.features {
  max-width: 1200px;
  margin: 0 auto;
}

.features h3 {
  font-size: 20px;
  margin: 15px 0 10px;
  color: #303133;
}

.features p {
  color: #606266;
  line-height: 1.6;
}

:deep(.el-card__header) {
  text-align: center;
  padding: 30px 20px 10px;
}

:deep(.el-card__body) {
  text-align: center;
}
</style>
