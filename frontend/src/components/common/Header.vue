<template>
  <el-header class="app-header">
    <div class="header-content">
      <div class="logo">
        <el-icon :size="28"><Document /></el-icon>
        <span>LaTeX PPT Generator</span>
      </div>
      <el-menu
        v-if="userStore.isLoggedIn()"
        :default-active="activeIndex"
        class="header-menu"
        mode="horizontal"
        router
      >
        <el-menu-item index="/">Home</el-menu-item>
        <el-menu-item index="/generate">Generate</el-menu-item>
        <el-menu-item index="/knowledge">Knowledge Base</el-menu-item>
        <el-menu-item index="/history">History</el-menu-item>
      </el-menu>
      <div class="header-actions">
        <template v-if="userStore.isLoggedIn()">
          <span class="username">{{ userStore.user?.username }}</span>
          <el-button @click="handleLogout">Logout</el-button>
        </template>
        <template v-else>
          <el-button @click="$router.push('/login')">Login</el-button>
        </template>
      </div>
    </div>
  </el-header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeIndex = computed(() => route.path)

const handleLogout = () => {
  userStore.logout()
  ElMessage.success('Logged out successfully')
  router.push('/')
}
</script>

<style scoped>
.app-header {
  background-color: #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 0 20px;
  height: 60px;
  display: flex;
  align-items: center;
}

.header-content {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 20px;
  font-weight: bold;
  color: #409eff;
}

.header-menu {
  flex: 1;
  margin: 0 40px;
  border-bottom: none;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 15px;
}

.username {
  color: #606266;
  font-weight: 500;
}
</style>
