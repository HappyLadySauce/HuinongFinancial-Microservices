<template>
  <el-dialog
    v-model="visible"
    :title="dialogTitle"
    :width="width"
    :fullscreen="fullscreen"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <!-- 对话框内容 -->
    <div class="dialog-content">
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        :label-width="labelWidth"
        :label-position="labelPosition"
        @submit.prevent="handleSubmit"
      >
        <slot name="form" :form="formData" :rules="formRules" />
      </el-form>
    </div>

    <!-- 对话框底部 -->
    <template #footer>
      <div class="dialog-footer">
        <slot name="footer" :form="formData" :loading="loading">
          <el-button @click="handleCancel" :disabled="loading">
            {{ cancelText }}
          </el-button>
          <el-button 
            type="primary" 
            @click="handleSubmit"
            :loading="loading"
          >
            {{ confirmText }}
          </el-button>
        </slot>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, reactive } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'

interface Props {
  modelValue: boolean
  title?: string
  mode?: 'create' | 'edit'
  width?: string
  fullscreen?: boolean
  labelWidth?: string
  labelPosition?: 'left' | 'right' | 'top'
  initialData?: Record<string, any>
  rules?: FormRules
  confirmText?: string
  cancelText?: string
  createTitle?: string
  editTitle?: string
}

const props = withDefaults(defineProps<Props>(), {
  title: '',
  mode: 'create',
  width: '600px',
  fullscreen: false,
  labelWidth: '120px',
  labelPosition: 'right',
  initialData: () => ({}),
  rules: () => ({}),
  confirmText: '确定',
  cancelText: '取消',
  createTitle: '新增',
  editTitle: '编辑'
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  submit: [data: Record<string, any>, mode: 'create' | 'edit']
  cancel: []
  close: []
}>()

const formRef = ref<FormInstance>()
const loading = ref(false)
const formData = reactive<Record<string, any>>({})

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const dialogTitle = computed(() => {
  if (props.title) return props.title
  return props.mode === 'create' ? props.createTitle : props.editTitle
})

const formRules = computed(() => props.rules)

// 监听初始数据变化
watch(() => props.initialData, (newData) => {
  Object.assign(formData, newData)
}, { deep: true, immediate: true })

// 监听对话框打开状态
watch(visible, (isVisible) => {
  if (isVisible) {
    // 对话框打开时重置表单
    resetForm()
  }
})

const resetForm = () => {
  // 清空表单数据
  Object.keys(formData).forEach(key => {
    delete formData[key]
  })
  // 设置初始数据
  Object.assign(formData, props.initialData)
  
  // 清除表单验证
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    // 表单验证
    await formRef.value.validate()
    
    loading.value = true
    emit('submit', { ...formData }, props.mode)
    
  } catch (error) {
    console.error('表单验证失败:', error)
  } finally {
    loading.value = false
  }
}

const handleCancel = () => {
  emit('cancel')
  visible.value = false
}

const handleClose = () => {
  emit('close')
}

// 暴露给父组件的方法
const setLoading = (isLoading: boolean) => {
  loading.value = isLoading
}

const closeDialog = () => {
  visible.value = false
}

const validateForm = () => {
  return formRef.value?.validate()
}

const resetFields = () => {
  formRef.value?.resetFields()
}

const clearValidate = () => {
  formRef.value?.clearValidate()
}

defineExpose({
  setLoading,
  closeDialog,
  validateForm,
  resetFields,
  clearValidate,
  resetForm
})
</script>

<style scoped>
.dialog-content {
  max-height: 60vh;
  overflow-y: auto;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

:deep(.el-dialog__body) {
  padding: 20px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
}
</style> 