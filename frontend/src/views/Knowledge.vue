<template>
  <div class="knowledge-page">
    <Header />
    <el-main class="content">
      <el-card>
        <template #header>
          <div class="header-actions">
            <h3>Knowledge Base</h3>
            <el-upload
              :show-file-list="false"
              :before-upload="handleUpload"
              accept=".pdf,.docx,.txt,.md"
            >
              <el-button type="primary">
                <el-icon><Upload /></el-icon> Upload Document
              </el-button>
            </el-upload>
          </div>
        </template>
        
        <el-table :data="documents" v-loading="loading" style="width: 100%">
          <el-table-column prop="filename" label="Filename" />
          <el-table-column prop="file_type" label="Type" width="100" />
          <el-table-column prop="file_size" label="Size" width="120">
            <template #default="{ row }">
              {{ formatFileSize(row.file_size) }}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="Status" width="120">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="chunk_count" label="Chunks" width="100" />
          <el-table-column prop="created_at" label="Upload Date" width="180">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="Actions" width="150">
            <template #default="{ row }">
              <el-button
                type="danger"
                size="small"
                @click="handleDelete(row.id)"
              >
                Delete
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </el-main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Header from '@/components/common/Header.vue'
import { getDocuments, uploadDocument, deleteDocument } from '@/api/knowledge'
import { ElMessage, ElMessageBox, ElLoading } from 'element-plus'
import type { Document } from '@/types'

const documents = ref<Document[]>([])
const loading = ref(false)

const loadDocuments = async () => {
  loading.value = true
  try {
    documents.value = await getDocuments()
  } catch (error) {
    console.error('Load documents error:', error)
  } finally {
    loading.value = false
  }
}

const handleUpload = async (file: File) => {
  const loadingInstance = ElLoading.service({ fullscreen: true, text: 'Uploading...' })
  try {
    await uploadDocument(file)
    ElMessage.success('Document uploaded successfully')
    loadDocuments()
  } catch (error) {
    console.error('Upload error:', error)
  } finally {
    loadingInstance.close()
  }
  return false
}

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('Are you sure to delete this document?', 'Warning', {
      type: 'warning'
    })
    
    await deleteDocument(id)
    ElMessage.success('Document deleted')
    loadDocuments()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Delete error:', error)
    }
  }
}

const formatFileSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleString()
}

const getStatusType = (status: string) => {
  const types: Record<string, any> = {
    pending: 'info',
    processing: 'warning',
    completed: 'success',
    failed: 'danger'
  }
  return types[status] || 'info'
}

onMounted(() => {
  loadDocuments()
})
</script>

<style scoped>
.knowledge-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.content {
  flex: 1;
  padding: 20px;
}

.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions h3 {
  margin: 0;
}
</style>
