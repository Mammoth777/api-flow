<template>
  <div class="workflow-info">
    <div class="workflow-header">
      <div class="workflow-title-section">
        <input 
          type="text" 
          class="workflow-name-input" 
          v-model="workflowName" 
          placeholder="输入工作流名称" 
          @input="emitUpdate"
        />
        <div class="workflow-desc-wrapper">
          <textarea 
            class="workflow-desc-input" 
            v-model="workflowDesc" 
            placeholder="输入工作流描述（可选）"
            @input="handleInput"
            rows="1"
          ></textarea>
        </div>
      </div>
      <div class="workflow-actions">
        <button 
          class="action-button"
          :class="{'loading': isSaving}" 
          @click="$emit('save')"
          :disabled="isSaving"
        >
          {{ isSaving ? '保存中...' : '保存' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';

// 定义属性
const props = defineProps({
  name: {
    type: String,
    default: '新工作流'
  },
  description: {
    type: String,
    default: ''
  },
  isSaving: {
    type: Boolean,
    default: false
  }
});

// 定义事件
const emit = defineEmits(['update:name', 'update:description', 'save']);

// 本地状态
const workflowName = ref(props.name);
const workflowDesc = ref(props.description);

// 监听prop变化，同步到本地状态
watch(() => props.name, (newValue) => {
  workflowName.value = newValue;
});

watch(() => props.description, (newValue) => {
  workflowDesc.value = newValue;
});

// 自动调整文本域高度避免滚动条
const handleInput = (e: Event) => {
  const textarea = e.target as HTMLTextAreaElement;
  textarea.style.height = 'auto';
  textarea.style.height = textarea.scrollHeight + 'px';
  
  emitUpdate();
};

// 更新父组件数据
const emitUpdate = () => {
  emit('update:name', workflowName.value);
  emit('update:description', workflowDesc.value);
};
</script>

<style scoped>
.workflow-info {
  padding: 0;
  background-color: #fff;
  border-bottom: 1px solid #e8e8e8;
}

.workflow-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  background: linear-gradient(to right, #f8f8f8, #ffffff);
}

.workflow-title-section {
  flex: 1;
  padding-right: 20px;
}

.workflow-name-input {
  width: 100%;
  font-size: 16px;
  font-weight: 500;
  border: none;
  background: transparent;
  color: #333;
  padding: 4px 0;
  transition: all 0.3s;
  position: relative;
  box-shadow: 0 0 0 0 transparent;
}

.workflow-name-input:focus {
  outline: none;
  box-shadow: 0 1px 0 0 #1890ff;
}

.workflow-desc-wrapper {
  margin-top: 4px;
  position: relative;
}

.workflow-desc-input {
  width: 100%;
  border: none;
  background: transparent;
  color: #666;
  font-size: 13px;
  padding: 4px 0;
  resize: none;
  transition: all 0.3s;
  overflow: hidden;
  height: auto;
  max-height: 80px;
  min-height: 20px;
  line-height: 1.5;
  position: relative;
  box-shadow: 0 0 0 0 transparent;
}

.workflow-desc-input:focus {
  outline: none;
  box-shadow: 0 1px 0 0 #1890ff;
}

.workflow-actions {
  display: flex;
  gap: 8px;
}

.action-button {
  padding: 4px 12px;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s;
  min-width: 70px;
}

.action-button:hover:not(:disabled) {
  background: #40a9ff;
}

.action-button:disabled {
  background: #bae7ff;
  cursor: not-allowed;
}

.action-button.loading {
  position: relative;
  color: transparent;
}

.action-button.loading::after {
  content: "";
  position: absolute;
  width: 16px;
  height: 16px;
  top: 50%;
  left: 50%;
  margin-top: -8px;
  margin-left: -8px;
  border-radius: 50%;
  border: 2px solid rgba(255, 255, 255, 0.2);
  border-top-color: white;
  animation: spin 1s infinite linear;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>