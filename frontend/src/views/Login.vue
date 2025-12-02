<template>
  <div class="login-page">
    <el-card class="login-card">
      <template #header>
        <h2>{{ isRegister ? 'Register' : 'Login' }}</h2>
      </template>
      
      <el-form :model="form" :rules="rules" ref="formRef" label-position="top">
        <el-form-item label="Username" prop="username">
          <el-input v-model="form.username" placeholder="Enter username" />
        </el-form-item>
        
        <el-form-item v-if="isRegister" label="Email" prop="email">
          <el-input v-model="form.email" type="email" placeholder="Enter email" />
        </el-form-item>
        
        <el-form-item label="Password" prop="password">
          <el-input v-model="form.password" type="password" placeholder="Enter password" show-password />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="loading" style="width: 100%">
            {{ isRegister ? 'Register' : 'Login' }}
          </el-button>
        </el-form-item>
      </el-form>
      
      <div class="toggle-mode">
        <el-link @click="toggleMode">
          {{ isRegister ? 'Already have an account? Login' : "Don't have an account? Register" }}
        </el-link>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const isRegister = ref(false)
const loading = ref(false)
const formRef = ref<FormInstance>()

const form = reactive({
  username: '',
  email: '',
  password: ''
})

const rules: FormRules = {
  username: [{ required: true, message: 'Please enter username', trigger: 'blur' }],
  email: [
    { required: true, message: 'Please enter email', trigger: 'blur' },
    { type: 'email', message: 'Please enter valid email', trigger: 'blur' }
  ],
  password: [
    { required: true, message: 'Please enter password', trigger: 'blur' },
    { min: 6, message: 'Password must be at least 6 characters', trigger: 'blur' }
  ]
}

const toggleMode = () => {
  isRegister.value = !isRegister.value
  form.email = ''
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    loading.value = true
    try {
      if (isRegister.value) {
        await userStore.register(form)
        ElMessage.success('Registered successfully')
      } else {
        await userStore.login({
          username: form.username,
          password: form.password
        })
        ElMessage.success('Logged in successfully')
      }
      
      const redirect = (route.query.redirect as string) || '/generate'
      router.push(redirect)
    } catch (error) {
      console.error('Auth error:', error)
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 100%;
  max-width: 400px;
}

:deep(.el-card__header) {
  text-align: center;
}

:deep(.el-card__header h2) {
  margin: 0;
  color: #303133;
}

.toggle-mode {
  text-align: center;
  margin-top: 10px;
}
</style>
