import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { PPTRecord } from '@/types'

export const usePPTStore = defineStore('ppt', () => {
  const currentPPT = ref<PPTRecord | null>(null)
  const pptHistory = ref<PPTRecord[]>([])
  const templates = ref<string[]>([])

  function setCurrentPPT(ppt: PPTRecord) {
    currentPPT.value = ppt
  }

  function setPPTHistory(history: PPTRecord[]) {
    pptHistory.value = history
  }

  function setTemplates(templateList: string[]) {
    templates.value = templateList
  }

  function clearCurrentPPT() {
    currentPPT.value = null
  }

  return {
    currentPPT,
    pptHistory,
    templates,
    setCurrentPPT,
    setPPTHistory,
    setTemplates,
    clearCurrentPPT
  }
})
