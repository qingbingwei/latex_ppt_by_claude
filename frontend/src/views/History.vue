<template>
  <div class="history-page">
    <Header />
    <el-main class="content">
      <el-card>
        <template #header>
          <h3>PPT Generation History</h3>
        </template>
        
        <el-table :data="history" v-loading="loading" style="width: 100%">
          <el-table-column prop="title" label="Title" />
          <el-table-column prop="template" label="Template" width="120" />
          <el-table-column prop="status" label="Status" width="120">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="Created" width="180">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="Actions" width="250">
            <template #default="{ row }">
              <el-button size="small" @click="handleView(row)">View</el-button>
              <el-button
                v-if="row.pdf_path"
                size="small"
                type="success"
                @click="handleDownload(row.id)"
              >
                Download
              </el-button>
              <el-button
                size="small"
                type="danger"
                @click="handleDelete(row.id)"
              >
                Delete
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
      
      <el-dialog v-model="dialogVisible" title="PPT Details" width="80%">
        <div v-if="selectedPPT">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="Title">{{ selectedPPT.title }}</el-descriptions-item>
            <el-descriptions-item label="Template">{{ selectedPPT.template }}</el-descriptions-item>
            <el-descriptions-item label="Status">
              <el-tag :type="getStatusType(selectedPPT.status)">{{ selectedPPT.status }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="Created">{{ formatDate(selectedPPT.created_at) }}</el-descriptions-item>
          </el-descriptions>
          
          <h4 style="margin-top: 20px">Prompt:</h4>
          <el-input
            v-model="selectedPPT.prompt"
            type="textarea"
            :rows="4"
            readonly
          />
          
          <h4 style="margin-top: 20px">LaTeX Content:</h4>
          <el-input
            v-model="selectedPPT.latex_content"
            type="textarea"
            :rows="15"
            readonly
          />
        </div>
      </el-dialog>
    </el-main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Header from '@/components/common/Header.vue'
import { getPPTHistory, deletePPT, downloadPPT } from '@/api/ppt'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { PPTRecord } from '@/types'

const history = ref<PPTRecord[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const selectedPPT = ref<PPTRecord | null>(null)

const loadHistory = async () => {
  loading.value = true
  try {
    history.value = await getPPTHistory()
  } catch (error) {
    console.error('Load history error:', error)
  } finally {
    loading.value = false
  }
}

const handleView = (ppt: PPTRecord) => {
  selectedPPT.value = ppt
  dialogVisible.value = true
}

const handleDownload = async (id: number) => {
  try {
    await downloadPPT(id)
  } catch (error) {
    console.error('Download error:', error)
    ElMessage.error('Download failed')
  }
}

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('Are you sure to delete this PPT?', 'Warning', {
      type: 'warning'
    })
    
    await deletePPT(id)
    ElMessage.success('PPT deleted')
    loadHistory()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Delete error:', error)
    }
  }
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleString()
}

const getStatusType = (status: string) => {
  const types: Record<string, any> = {
    pending: 'info',
    generating: 'warning',
    completed: 'success',
    failed: 'danger'
  }
  return types[status] || 'info'
}

onMounted(() => {
  loadHistory()
})
</script>

<style scoped>
.history-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.content {
  flex: 1;
  padding: 20px;
}

:deep(.el-textarea__inner) {
  font-family: 'Courier New', monospace;
}
</style>
