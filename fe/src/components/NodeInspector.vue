<template>
  <div class="node-inspector" v-show="selectedNode">
    <div class="inspector-header">
      <h3>节点属性</h3>
      <button class="close-button" @click="$emit('close')">&times;</button>
    </div>
    <div class="inspector-body">
      <div class="form-group">
        <label>节点类型</label>
        <select v-model="nodeData.nodeType" @change="updateNode">
          <option v-for="type in availableNodeTypes" :key="type.code" :value="type.code">
            {{ type.name }}
          </option>
        </select>
      </div>
      
      <div class="form-group">
        <label>节点名称</label>
        <input type="text" v-model="nodeData.name" @input="updateNode" placeholder="输入节点名称" />
      </div>
      
      <div class="form-group">
        <label>节点描述</label>
        <textarea v-model="nodeData.description" @input="updateNode" placeholder="输入节点描述" rows="2"></textarea>
      </div>
      
      <div class="form-group">
        <label>节点配置</label>
        <textarea 
          v-model="configStr" 
          @input="updateConfig" 
          class="config-editor" 
          placeholder="输入 JSON 格式配置" 
          rows="6"
        ></textarea>
        <div class="error-message" v-if="configError">
          {{ configError }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { Cell } from '@antv/x6';
import type { NodeType } from '../services/node-type.service';

// Props 定义，使用默认值防止空引用
const props = defineProps({
  selectedNode: {
    type: Object as () => Cell | null,
    default: null
  },
  nodeTypes: {
    type: Array as () => NodeType[],
    default: () => []
  }
});

const emit = defineEmits(['close', 'nodeUpdated']);

// 表单数据，使用安全的默认值
const nodeData = ref({
  nodeType: '',
  name: '',
  description: '',
  config: {} as Record<string, any>
});

// JSON 配置编辑器
const configStr = ref('{}');
const configError = ref('');

// 获取可用的节点类型，添加计算属性以处理空值
const availableNodeTypes = computed(() => {
  return props.nodeTypes || [];
});

// 重置表单为默认值 - 移到这里以便在 watch 中使用
const resetForm = () => {
  nodeData.value = {
    nodeType: '',
    name: '',
    description: '',
    config: {}
  };
  configStr.value = '{}';
  configError.value = '';
};

// 更新配置字符串显示，安全处理 JSON 转换
const updateConfigString = (config: Record<string, any>) => {
  try {
    configStr.value = JSON.stringify(config || {}, null, 2);
    configError.value = '';
  } catch (err) {
    configStr.value = '{}';
    configError.value = '无效的 JSON 配置';
  }
};

// 监听节点变化，安全处理节点数据提取
watch(() => props.selectedNode, (newNode) => {
  if (!newNode) {
    // 如果节点为空，重置表单数据
    resetForm();
    return;
  }
  
  try {
    // 获取节点数据，使用安全的数据访问方式
    const data = newNode.getData() || {};
    const attrs = newNode.getAttrs() || {};
    
    // 填充表单数据，确保每个字段都有默认值
    nodeData.value = {
      nodeType: data.nodeType || '',
      name: attrs.label?.text + '' || '未命名节点',
      description: data.description || '',
      config: data.config || {}
    };
    
    // 更新配置字符串
    updateConfigString(nodeData.value.config);
  } catch (err) {
    console.error('解析节点数据时出错:', err);
    resetForm();
  }
}, { immediate: true });

// 安全地更新节点数据
const updateNode = () => {
  if (!props.selectedNode) return;
  
  try {
    // 更新节点标签文本
    props.selectedNode.attr('label/text', nodeData.value.name);
    
    // 更新节点数据
    props.selectedNode.setData({
      ...(props.selectedNode.getData() || {}),
      nodeType: nodeData.value.nodeType,
      description: nodeData.value.description,
      config: nodeData.value.config
    });
    
    // 通知父组件节点已更新
    emit('nodeUpdated', props.selectedNode);
  } catch (err) {
    console.error('更新节点时出错:', err);
  }
};

// 更新配置
const updateConfig = () => {
  try {
    const newConfig = JSON.parse(configStr.value);
    nodeData.value.config = newConfig;
    configError.value = '';
    updateNode();
  } catch (err) {
    configError.value = '无效的 JSON 格式';
  }
};
</script>

<style scoped>
.node-inspector {
  position: absolute;
  top: 0;
  right: 0;
  width: 280px;
  height: 100%;
  background-color: #fff;
  border-left: 1px solid #e8e8e8;
  box-shadow: -2px 0 8px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  z-index: 10;
}

.inspector-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 16px;
  border-bottom: 1px solid #e8e8e8;
  background-color: #fafafa;
}

.inspector-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.close-button {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  color: #999;
}

.close-button:hover {
  color: #333;
}

.inspector-body {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 12px;
  color: #666;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 6px 8px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  font-size: 13px;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #40a9ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

.config-editor {
  font-family: monospace;
  font-size: 12px;
  line-height: 1.5;
  resize: vertical;
}

.error-message {
  color: #f5222d;
  font-size: 12px;
  margin-top: 4px;
}
</style>