<template>
  <div class="generate-page">
    <Header />
    <el-container class="main-container">
      <el-aside width="400px" class="sidebar">
        <el-card>
          <template #header>
            <h3>Generate PPT</h3>
          </template>
          
          <el-form :model="form" label-position="top">
            <el-form-item label="Title">
              <el-input v-model="form.title" placeholder="Enter PPT title" />
            </el-form-item>
            
            <el-form-item label="Requirements">
              <el-input
                v-model="form.prompt"
                type="textarea"
                :rows="6"
                placeholder="Describe your PPT requirements..."
              />
            </el-form-item>
            
            <el-form-item label="Template">
              <el-select v-model="form.template" style="width: 100%">
                <el-option label="Default" value="default" />
                <el-option label="Madrid" value="madrid" />
                <el-option label="Modern" value="modern" />
              </el-select>
            </el-form-item>
            
            <el-form-item label="Use Knowledge Base">
              <el-switch v-model="useKnowledge" />
            </el-form-item>
            
            <el-form-item label="AI Model">
              <el-radio-group v-model="form.use_openai">
                <el-radio :label="true">OpenAI</el-radio>
                <el-radio :label="false">Claude</el-radio>
              </el-radio-group>
            </el-form-item>
            
            <el-button
              type="primary"
              @click="handleGenerate"
              :loading="generating"
              style="width: 100%"
            >
              Generate PPT
            </el-button>
          </el-form>
        </el-card>
      </el-aside>
      
      <el-main class="content-area">
        <el-tabs v-model="activeTab">
          <el-tab-pane label="LaTeX Code" name="latex">
            <el-input
              v-model="latexContent"
              type="textarea"
              :rows="25"
              placeholder="LaTeX code will appear here..."
            />
          </el-tab-pane>
          <el-tab-pane label="PDF Preview" name="pdf">
            <div v-if="pdfUrl" class="pdf-preview">
              <iframe :src="pdfUrl" width="100%" height="600px"></iframe>
            </div>
            <el-empty v-else description="No PDF generated yet" />
          </el-tab-pane>
        </el-tabs>
        
        <div v-if="currentPPT" class="actions">
          <el-button @click="handleCompile">Compile LaTeX</el-button>
          <el-button v-if="currentPPT.pdf_path" @click="handleDownload">
            <el-icon><Download /></el-icon> Download PDF
          </el-button>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import Header from '@/components/common/Header.vue'
import { generatePPT, compileLaTeX, downloadPPT } from '@/api/ppt'
import { usePPTStore } from '@/store/ppt'
import { ElMessage } from 'element-plus'

const pptStore = usePPTStore()

const form = reactive({
  title: '',
  prompt: '',
  template: 'default',
  use_openai: true
})

const useKnowledge = ref(false)
const generating = ref(false)
const activeTab = ref('latex')
const latexContent = ref('')
const pdfUrl = ref('')
const currentPPT = ref<any>(null)

const handleGenerate = async () => {
  if (!form.title || !form.prompt) {
    ElMessage.warning('Please fill in all required fields')
    return
  }
  
  generating.value = true
  try {
    const result = await generatePPT(form)
    currentPPT.value = result
    latexContent.value = result.latex_content || ''
    
    if (result.pdf_path) {
      pdfUrl.value = downloadPPT(result.id)
    }
    
    ElMessage.success('PPT generated successfully')
    pptStore.setCurrentPPT(result)
  } catch (error) {
    console.error('Generation error:', error)
  } finally {
    generating.value = false
  }
}

const handleCompile = async () => {
  if (!latexContent.value) {
    ElMessage.warning('No LaTeX content to compile')
    return
  }
  
  try {
    const result = await compileLaTeX(latexContent.value)
    currentPPT.value = result
    pdfUrl.value = downloadPPT(result.id)
    ElMessage.success('Compiled successfully')
  } catch (error) {
    console.error('Compilation error:', error)
  }
}

const handleDownload = () => {
  if (currentPPT.value) {
    window.open(downloadPPT(currentPPT.value.id), '_blank')
  }
}
</script>

<style scoped>
.generate-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-container {
  flex: 1;
  background-color: #f5f7fa;
}

.sidebar {
  background-color: #fff;
  padding: 20px;
  overflow-y: auto;
}

.content-area {
  padding: 20px;
}

.pdf-preview {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
}

.actions {
  margin-top: 20px;
  display: flex;
  gap: 10px;
}

:deep(.el-textarea__inner) {
  font-family: 'Courier New', monospace;
}
</style>
