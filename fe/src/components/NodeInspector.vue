<template>
  <div class="node-inspector" v-show="selectedNode" :class="{'compact-mode': isCompactMode, 'dragging': isDragging}" 
       :style="compactPositionStyle" ref="inspectorEl">
    <div class="inspector-header" @mousedown="startDrag">
      <h3>{{ nodeData.name }}</h3>
      <div class="header-actions">
        <!-- 添加切换按钮 -->
        <button class="toggle-mode-button" @click="toggleMode">
          {{ isCompactMode ? '编辑' : '收起详情' }}
        </button>
        <button class="close-button" @click="$emit('close')">&times;</button>
      </div>
    </div>
    <div class="inspector-body">
      <!-- 基本信息区域（在简洁模式下显示） -->
      <div class="basic-info">
        <!-- 显示节点类型（只读） -->
        <div class="form-group">
          <label>节点类型</label>
          <input type="text" :value="nodeData.nodeType" disabled class="readonly-input" />
        </div>
      </div>
      
      <!-- 详细信息区域（在完整模式下显示） -->
      <div class="detailed-info" v-if="!isCompactMode">
        <div class="form-group">
          <label>节点名称</label>
          <input type="text" v-model="nodeData.name" @blur="debouncedUpdateNode" placeholder="输入节点名称" />
        </div>

        <div class="form-group">
          <label>节点描述</label>
          <textarea v-model="nodeData.description" @blur="debouncedUpdateNode" placeholder="输入节点描述" rows="2"></textarea>
        </div>
        
        <!-- 新增：节点输入配置区域 -->
        <div class="form-group input-config-section" v-if="nodeInputs && nodeInputs.length > 0">
          <label>节点输入配置</label>
          <NodeInputConfig 
            :input-fields="nodeInputs" 
            v-model="nodeData.config" 
            @update:modelValue="onInputConfigUpdate"
          />
        </div>

        <!-- 保留原有的节点配置编辑器，作为高级模式 -->
        <div class="form-group advanced-config" v-if="showAdvancedConfig">
          <div class="advanced-config-header">
            <label>高级配置（JSON）</label>
            <button class="toggle-button" @click="showAdvancedConfig = false">隐藏</button>
          </div>
          <textarea 
            v-model="configStr" 
            @blur="updateConfig" 
            class="config-editor" 
            placeholder="输入 JSON 格式配置" 
            rows="6"
          ></textarea>
          <div class="error-message" v-if="configError">
            {{ configError }}
          </div>
        </div>
        <div class="form-group" v-else>
          <button class="toggle-button show-advanced" @click="showAdvancedConfig = true">
            显示高级配置
          </button>
        </div>
        
        <!-- 节点输出格式区域 -->
        <div class="form-group output-format-section" v-if="nodeOutputs && nodeOutputs.length > 0">
          <label>输出格式</label>
          <div class="output-format-container">
            <table class="output-format-table">
              <thead>
                <tr>
                  <th>参数名</th>
                  <th>数据类型</th>
                  <th>描述</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(output, idx) in nodeOutputs" :key="idx">
                  <td class="param-name">{{ output.field }}</td>
                  <td class="param-type">{{ output.type }}</td>
                  <td class="param-desc">{{ output.desc }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
      
      <!-- 删除节点按钮，始终显示 -->
      <div class="node-actions">
        <button class="delete-node-btn" @click="confirmDeleteNode">删除节点</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed, onMounted, onUnmounted, nextTick } from 'vue';
import { Cell } from '@antv/x6';
import type { NodeType } from '../services/node-type.service';
import NodeInputConfig from './NodeInputConfig.vue';

// Props 定义，使用默认值防止空引用
const props = defineProps({
  selectedNode: {
    type: Object as () => Cell | null,
    default: null
  },
  nodeTypes: {
    type: Array as () => NodeType[],
    default: () => []
  },
  // 增加新的props来控制初始展开模式
  inspectorMode: {
    type: String,
    default: 'compact' // 'compact' 或 'full'
  },
  // 画布相对屏幕左上角偏移量
  getGraphBounding: {
    type: Function,
    default: () => ({ x: 0, y: 0, w: 0, h: 0 })
  }
});

const emit = defineEmits(['close', 'nodeUpdated', 'deleteNode', 'update:inspectorMode']);

// 表单数据，使用安全的默认值
const nodeData = ref({
  nodeType: '',
  name: '',
  description: '',
  config: {} as Record<string, any>
});

// 添加展示模式状态
const isCompactMode = ref(props.inspectorMode === 'compact');

// 切换展示模式
const toggleMode = () => {
  isCompactMode.value = !isCompactMode.value;
  emit('update:inspectorMode', isCompactMode.value ? 'compact' : 'full');
};

// 显示高级配置的状态
const showAdvancedConfig = ref(false);

// 防抖定时器
const debounceTimerId = ref<number | null>(null);

// 防抖更新函数
const debouncedUpdateNode = () => {
  if (debounceTimerId.value) {
    clearTimeout(debounceTimerId.value);
  }
  
  debounceTimerId.value = setTimeout(() => {
    updateNode();
    debounceTimerId.value = null;
  }, 300) as number;
};

// 计算属性：找到当前节点类型的输出格式
const nodeOutputs = computed(() => {
  const nodeType = props.nodeTypes.find(nodeType => {
    return nodeType.code === nodeData.value.nodeType;
  });
  return nodeType?.output || [];
});

// 计算属性：找到当前节点类型的输入字段定义
const nodeInputs = computed(() => {
  const nodeType = props.nodeTypes.find(nodeType => {
    return nodeType.code === nodeData.value.nodeType;
  });
  return (nodeType?.input || []).map(input => ({
    ...input,
    type: input.type as "string" | "number" | "boolean" | "object" | "options"
  }));
});

// JSON 配置编辑器
const configStr = ref('{}');
const configError = ref('');

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
  showAdvancedConfig.value = false;
  // 重置为默认显示模式
  isCompactMode.value = props.inspectorMode === 'compact';
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
    
    // 确保位置已初始化
    if (!positionInitialized.value) {
      // 使用 nextTick 确保 DOM 已更新
      nextTick(() => {
        initPosition();
      });
    }
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
      name: nodeData.value.name,
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

// 高级配置更新
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

// 处理输入配置更新
const onInputConfigUpdate = (newConfig: Record<string, any>) => {
  nodeData.value.config = newConfig;
  updateConfigString(newConfig);
  // 不直接调用updateNode，而是采用防抖方式
  debouncedUpdateNode();
};

// 确认并删除节点
const confirmDeleteNode = () => {
  if (!props.selectedNode) return;
  emit('deleteNode', props.selectedNode);
  emit('close');
};

// 拖动功能相关变量
const isDragging = ref(false);
const position = ref({ x: 50, y: 50 }); // 初始位置
const dragOffset = ref({ x: 0, y: 0 });
const inspectorEl = ref<HTMLElement | null>(null);
const positionInitialized = ref(false); // 添加一个标记，表示位置是否已初始化

// 计算样式
const compactPositionStyle = computed(() => {
  if (!isCompactMode.value) return {};
  return {
    top: `${position.value.y}px`,
    right: 'auto',
    left: `${position.value.x}px`,
  };
});

// 初始化位置，仅在首次显示时调用一次
const initPosition = () => {
  if (positionInitialized.value || !inspectorEl.value) return;
  
  const bound = props.getGraphBounding()
  const margin = 10
  position.value = {
    x: bound.w - 300 - margin,
    y: margin
  };
  
  positionInitialized.value = true;
};

// 开始拖动
const startDrag = (event: MouseEvent) => {
  if (!isCompactMode.value || !inspectorEl.value) return;
  
  // 防止文本选择
  event.preventDefault();
  
  // 确保位置已初始化
  if (!positionInitialized.value) {
    initPosition();
  }
  
  isDragging.value = true;
  const offset = props.getGraphBounding()
  
  // 计算鼠标在元素内的偏移
  const rect = inspectorEl.value.getBoundingClientRect();
  dragOffset.value = {
    x: event.clientX - rect.left + offset.x,
    y: event.clientY - rect.top + offset.y
  };
  
  // 添加移动和释放事件监听器
  document.addEventListener('mousemove', onDrag);
  document.addEventListener('mouseup', stopDrag);
};

// 拖动中
const onDrag = (event: MouseEvent) => {
  if (!isDragging.value) return;
  
  // 更新位置
  position.value = {
    x: event.clientX - dragOffset.value.x,
    y: event.clientY - dragOffset.value.y
  };
};

// 停止拖动
const stopDrag = () => {
  isDragging.value = false;
  
  // 移除事件监听器
  document.removeEventListener('mousemove', onDrag);
  document.removeEventListener('mouseup', stopDrag);
};

// 组件挂载和卸载时的清理
onMounted(() => {
  // 确保在组件卸载时清理事件监听器
  window.addEventListener('resize', onWindowResize);
  // 添加ESC键监听
  document.addEventListener('keydown', handleKeyDown);
});

onUnmounted(() => {
  document.removeEventListener('mousemove', onDrag);
  document.removeEventListener('mouseup', stopDrag);
  window.removeEventListener('resize', onWindowResize);
  // 移除ESC键监听
  document.removeEventListener('keydown', handleKeyDown);
});

// 处理键盘按下事件，关注ESC键
const handleKeyDown = (event: KeyboardEvent) => {
  // 检测ESC键
  if (event.key === 'Escape' || event.key === 'Esc') {
    // 触发关闭事件
    emit('close');
  }
};

// 窗口大小变化时重新定位，确保不超出边界
const onWindowResize = () => {
  if (!inspectorEl.value || !isCompactMode.value) return;
  
  const rect = inspectorEl.value.getBoundingClientRect();
  const windowWidth = window.innerWidth;
  const windowHeight = window.innerHeight;
  
  // 确保不超出窗口边界
  if (position.value.x + rect.width > windowWidth) {
    position.value.x = windowWidth - rect.width - 10;
  }
  
  if (position.value.y + rect.height > windowHeight) {
    position.value.y = windowHeight - rect.height - 10;
  }
  
  if (position.value.x < 0) {
    position.value.x = 10;
  }
  
  if (position.value.y < 0) {
    position.value.y = 10;
  }
};
</script>

<style scoped>
.node-inspector {
  position: absolute;
  top: 0;
  right: 0;
  width: 500px;
  height: 100%;
  background-color: #fff;
  border-left: 1px solid #e8e8e8;
  box-shadow: -2px 0 8px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  z-index: 10;
  transition: width 0.3s, height 0.3s;
}

/* 紧凑模式样式 */
.node-inspector.compact-mode {
  width: 300px;
  height: 220px;
  top: 50px;
  right: auto;
  left: 50px;
  border-radius: 6px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
  border: 1px solid #d9d9d9;
  transition: none; /* 拖动时不需要过渡动画 */
}

/* 拖动中的样式 */
.node-inspector.dragging {
  opacity: 0.95;
  cursor: grabbing;
}

.inspector-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 16px;
  border-bottom: 1px solid #e8e8e8;
  background-color: #fafafa;
  cursor: grab; /* 显示可拖动的鼠标指针 */
}

.node-inspector.dragging .inspector-header {
  cursor: grabbing; /* 拖动时更改鼠标指针 */
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
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

.toggle-mode-button {
  background: #f0f0f0;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  padding: 4px 8px;
  font-size: 12px;
  cursor: pointer;
  color: #666;
  transition: all 0.3s;
}

.toggle-mode-button:hover {
  color: #1890ff;
  border-color: #1890ff;
  background-color: #f0f7ff;
}

.inspector-body {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
}

.basic-info {
  margin-bottom: 16px;
}

.detailed-info {
  padding-top: 10px;
  border-top: 1px dashed #e8e8e8;
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

.readonly-input {
  background-color: #f5f5f5;
  color: #666;
  cursor: not-allowed;
}

.input-config-section {
  border: 1px solid #f0f0f0;
  border-radius: 4px;
  padding: 10px;
  background-color: #fafafa;
}

.advanced-config-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.toggle-button {
  background: none;
  border: none;
  font-size: 12px;
  color: #1890ff;
  cursor: pointer;
  padding: 2px 6px;
}

.toggle-button.show-advanced {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  padding: 6px 12px;
  background-color: #f7f7f7;
  width: 100%;
}

.toggle-button:hover {
  color: #40a9ff;
  background-color: #f0f7ff;
}

.node-actions {
  margin-top: 24px;
  display: flex;
  justify-content: flex-end;
}

.delete-node-btn {
  background-color: #ff4d4f;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 8px 16px;
  cursor: pointer;
  font-size: 13px;
  transition: all 0.3s;
}

.delete-node-btn:hover {
  background-color: #ff7875;
}

/* 输出格式区域样式 */
.output-format-section {
  margin-top: 16px;
  margin-bottom: 16px;
}

.output-format-container {
  background-color: #fafafa;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  max-height: 200px;
  overflow-y: auto;
  font-size: 12px;
}

.output-format-table {
  width: 100%;
  border-collapse: collapse;
}

.output-format-table th,
.output-format-table td {
  padding: 6px 8px;
  text-align: left;
  border-bottom: 1px solid #e8e8e8;
  font-size: 12px;
}

.output-format-table th {
  background-color: #f0f0f0;
  font-weight: 500;
  color: #666;
}

.param-name {
  font-weight: 500;
  color: #1890ff;
}

.param-type {
  color: #722ed1;
  font-family: monospace;
}

.param-desc {
  color: #666;
}
</style>