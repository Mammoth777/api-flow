<template>
  <!-- 编辑系统节点 -->
  <div v-if="systemNodeType" class="form-group input-config-section">
    <ExecInputNode v-if="nodeType === 'execInput'"></ExecInputNode>
    <div v-else>系统节点 {{ nodeType }} 暂不支持编辑</div>
  </div>
   <!-- 无inputs -->
  <div v-else-if="!inputFields || !inputFields.length">
    <div class="form-group input-config-section">
      <label>节点输入配置</label>
      <div class="empty-message">暂无输入参数定义</div>
    </div>
  </div>
  <!-- 编辑非系统节点 -->
  <div v-else class="form-group input-config-section">
    <div class="section-header">
      <label>节点输入配置</label>
      <button class="toggle-button" @click="toggleAdvancedMode">
        {{ showAdvancedConfig ? '标准模式' : '高级模式' }}
      </button>
    </div>
    <!-- 标准模式 - 使用表单界面配置 -->
    <div v-if="!showAdvancedConfig" class="standard-mode">
      <NodeInputConfig 
        :input-fields="inputFields" 
        :modelValue="nodeConfig" 
        @update:modelValue="updateConfig"
      />
    </div>
    <!-- 高级模式 - 使用JSON文本编辑 -->
    <div v-else class="advanced-mode">
      <div class="advanced-config-header">
        <label>高级配置(JSON)</label>
      </div>
      <textarea 
        v-model="localConfigStr" 
        @blur="updateAdvancedConfig" 
        class="config-editor" 
        placeholder="输入 JSON 格式配置" 
        rows="6"
      ></textarea>
      <div class="error-message" v-if="configError">
        {{ configError }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits, ref, watch, computed } from 'vue';
import { InputField } from './types'
import NodeInputConfig from './NodeInputConfig.vue';
import { isSystemNode } from './nodeTypes'
import ExecInputNode from './systemNode/ExecInputNode.vue';

const props = defineProps({
  nodeType: {
    type: String,
    default: ''
  },
  inputFields: {
    type: Array as () => InputField[],
    required: true
  },
  nodeConfig: {
    type: Object,
    default: () => ({})
  },
  // 高级模式相关props
  showAdvancedConfig: {
    type: Boolean,
    default: false
  },
  configStr: {
    type: String,
    default: '{}'
  },
  configError: {
    type: String,
    default: ''
  }
});

const emit = defineEmits([
  'update:config', 
  'update:config-str', 
  'toggle-advanced-mode',
  'update-advanced-config'
]);
console.log('NodeInputSection.vue', props);
const systemNodeType = computed(() => isSystemNode(props.nodeType));

// 本地保存configStr，避免直接修改props
const localConfigStr = ref(props.configStr);

// 监听configStr变化
watch(() => props.configStr, (newValue) => {
  localConfigStr.value = newValue;
});

const updateConfig = (newConfig: Record<string, any>) => {
  emit('update:config', newConfig);
};

const toggleAdvancedMode = () => {
  emit('toggle-advanced-mode');
};

const updateAdvancedConfig = () => {
  emit('update-advanced-config');
};
</script>

<style scoped>
.input-config-section {
  border: 1px solid #f0f0f0;
  border-radius: 4px;
  padding: 10px;
  background-color: #fafafa;
  margin-bottom: 16px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 12px;
  color: #666;
}

.toggle-button {
  background: #f0f0f0;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  padding: 4px 8px;
  font-size: 12px;
  cursor: pointer;
  color: #666;
  transition: all 0.3s;
}

.toggle-button:hover {
  color: #1890ff;
  border-color: #1890ff;
  background-color: #f0f7ff;
}

.advanced-config-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.config-editor {
  width: 100%;
  font-family: monospace;
  font-size: 12px;
  line-height: 1.5;
  resize: vertical;
  padding: 6px 8px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
}

.error-message {
  color: #f5222d;
  font-size: 12px;
  margin-top: 4px;
}
</style>
