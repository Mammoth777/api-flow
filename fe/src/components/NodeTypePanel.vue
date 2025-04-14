<template>
  <div class="stencil-container">
    <div class="stencil-title">节点类型</div>
    <div ref="stencilContainer" class="stencil"></div>
    <!-- 节点类型加载状态显示 -->
    <div v-if="isLoading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>加载节点类型...</p>
    </div>
    <div v-else-if="error" class="error-state">
      <p>{{ error }}</p>
      <button @click="fetchNodeTypes" class="retry-button">重试</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { Graph } from '@antv/x6';
import { Dnd } from '@antv/x6-plugin-dnd';
import { nodeTypeService } from '../services/node-type.service';
import type { NodeType } from '../services/node-type.service';
import { getNodeCategory, getNodeColor } from './workflow'

// 接收的属性
const props = defineProps({
  graph: {
    type: Object as () => Graph | null,
    required: true,
  }
});

const emit = defineEmits(['nodeTypesLoaded']);

// 节点类型状态
const nodeTypes = ref<NodeType[]>([]);
const isLoading = ref<boolean>(false);
const error = ref<string | null>(null);
const stencilContainer = ref<HTMLElement | null>(null);
let dnd: Dnd | null = null;

// 节点唯一标识计数器
const generateNodeKey = () => `node-${Date.now()}-${Math.floor(Math.random() * 10000)}`;

// 获取节点类型
const fetchNodeTypes = async () => {
  isLoading.value = true;
  error.value = null;

  try {
    const result = await nodeTypeService.getNodeTypes();

    // 根据API返回的结构获取数据
    nodeTypes.value = result;
    
    // 向父组件emit节点类型
    emit('nodeTypesLoaded', nodeTypes.value);

    // 初始化面板
    if (props.graph) {
      // 使用setTimeout等待DOM更新
      setTimeout(() => {
        if (stencilContainer.value) {
          initStencil();
        }
      }, 0);
    }
  } catch (error: any) {
    console.error('获取节点类型失败:', error);
    error.value = '获取节点类型失败，请重试';
    nodeTypes.value = [];
  } finally {
    isLoading.value = false;
  }
};

// 监听graph变化，重新初始化Stencil
watch(() => props.graph, (newGraph) => {
  if (newGraph && stencilContainer.value) {
    initStencil();
  }
});

// 初始化 Stencil 面板
const initStencil = () => {
  if (!props.graph || !stencilContainer.value) {
    console.warn('无法初始化节点面板：图表或容器未就绪');
    return;
  }

  try {
    // 创建 DND 实例
    dnd = new Dnd({
      target: props.graph,
      scaled: false,
      getDragNode: (node) => node.clone({ keepId: true }),
      getDropNode: (node) => node.clone({ keepId: true }),
    });

    // 清空现有的 stencil 内容
    if (stencilContainer.value) {
      stencilContainer.value.innerHTML = '';
    } else {
      console.warn('stencilContainer不存在');
      return;
    }

    // 按分类整理节点类型
    const nodeTypesByCategory: Record<string, NodeType[]> = {};

    // 使用从 API 获取的节点类型
    nodeTypes.value.forEach(nodeType => {
      const category = nodeType.ui.category || '其他节点';
      if (!nodeTypesByCategory[category]) {
        nodeTypesByCategory[category] = [];
      }
      nodeTypesByCategory[category].push(nodeType);
    });

    // 创建分类和节点
    Object.entries(nodeTypesByCategory).forEach(([category, types]) => {
      const itemsEl = createCategory(category);

      types.forEach(nodeType => {
        const nodeEl = createNodeTemplate(
          nodeType.code,
          nodeType.name,
          getNodeColor(nodeType.code) // 使用节点颜色
        );
        itemsEl.appendChild(nodeEl);
      });
    });

    // 为所有节点模板绑定拖拽事件
    if (stencilContainer.value) {
      const nodeEls = stencilContainer.value.querySelectorAll('.stencil-node');
      nodeEls.forEach((nodeEl) => {
        const code = nodeEl.getAttribute('data-type') || 'task';
        const text = nodeEl.getAttribute('data-text') || '节点';
        const color = nodeEl.getAttribute('data-color');

        nodeEl.addEventListener('mousedown', (e) => {
          if (!props.graph || !dnd) return;

          const nodeAttrs: any = {
            label: {
              text: text,
            },
          };

          if (color) {
            nodeAttrs.body = {
              fill: color,
            };
          }

          const nodeKey = generateNodeKey();

          const category = getNodeCategory(code)
          const node = props.graph.createNode({
            id: nodeKey, // 似乎没用呢？
            shape: category, // 使用通用节点形状
            attrs: nodeAttrs,
            data: {
              id: nodeKey,
              nodeType: code,
              nodeKey: nodeKey,
              name: text,
              description: '',
              config: {},
              status: 'ready',
              ui: {
                category,
                color: color,
              },
            }
          });
          console.log('create node', node)

          dnd.start(node, e as any);
        });
      });
    }
  } catch (error) {
    console.error('初始化节点面板时出错:', error);
  }
};

// 创建分类标题
const createCategory = (title: string) => {
  if (!stencilContainer.value) return document.createElement('div');

  const categoryEl = document.createElement('div');
  categoryEl.className = 'stencil-category';

  const titleEl = document.createElement('div');
  titleEl.className = 'stencil-category-title';
  titleEl.textContent = title;

  const itemsEl = document.createElement('div');
  itemsEl.className = 'stencil-category-items';

  categoryEl.appendChild(titleEl);
  categoryEl.appendChild(itemsEl);
  stencilContainer.value.appendChild(categoryEl);

  return itemsEl;
};

// 创建节点模板元素
const createNodeTemplate = (type: string, text: string, backgroundColor?: string) => {
  const nodeEl = document.createElement('div');
  nodeEl.className = 'stencil-node';
  nodeEl.setAttribute('data-type', type);
  nodeEl.setAttribute('data-text', text);
  
  // 设置背景色 - 使用纯色背景代替渐变效果
  const bgColor = backgroundColor || getNodeColor(type);
  nodeEl.setAttribute('data-color', bgColor);
  
  // 扁平设计 - 使用单一背景色
  nodeEl.style.backgroundColor = bgColor;
  
  // 计算边框颜色 - 比背景色稍深
  const borderColor = adjustColor(bgColor, -10);
  nodeEl.style.borderColor = borderColor;
  
  // 更新节点内部结构
  const headerEl = document.createElement('div');
  headerEl.className = 'stencil-node-header';
  
  const typeLabel = document.createElement('div');
  typeLabel.className = 'stencil-node-type';
  typeLabel.textContent = type;
  
  headerEl.appendChild(typeLabel);

  const nameLabel = document.createElement('div');
  nameLabel.className = 'stencil-node-name';
  nameLabel.textContent = text;

  nodeEl.appendChild(headerEl);
  nodeEl.appendChild(nameLabel);

  return nodeEl;
};

// 颜色调整辅助函数
const adjustColor = (color: string, percent: number): string => {
  // 去掉#号并解析为RGB
  const hex = color.replace('#', '');
  let r = parseInt(hex.substring(0, 2), 16);
  let g = parseInt(hex.substring(2, 4), 16);
  let b = parseInt(hex.substring(4, 6), 16);

  // 调整颜色值
  r = Math.min(255, Math.max(0, Math.round(r * (1 + percent / 100))));
  g = Math.min(255, Math.max(0, Math.round(g * (1 + percent / 100))));
  b = Math.min(255, Math.max(0, Math.round(b * (1 + percent / 100))));

  // 转换回十六进制并添加#号
  const toHex = (c: number): string => {
    const hex = c.toString(16);
    return hex.length === 1 ? '0' + hex : hex;
  };

  return `#${toHex(r)}${toHex(g)}${toHex(b)}`;
};

// 组件挂载时获取节点类型
onMounted(() => {
  fetchNodeTypes();
});

// 对外暴露方法
defineExpose({
  fetchNodeTypes,
  nodeTypes
});
</script>

<style scoped>
.stencil-container {
  width: 120px;
  height: 100%;
  background-color: #f7f9fb;
  border-right: 1px solid #ddd;
  overflow-y: auto;
  scrollbar-width: none;
  /* Firefox */
  -ms-overflow-style: none;
  /* IE and Edge */
}

.stencil-container::-webkit-scrollbar {
  display: none;
  /* Chrome, Safari and Opera */
}

.stencil-title {
  padding: 8px 12px;
  font-size: 14px;
  font-weight: 500;
  border-bottom: 1px solid #e8e8e8;
  background-color: #fafafa;
}

::v-deep .stencil {
  padding: 8px;

  .stencil-category {
    margin-bottom: 12px;
  }
  .stencil-category-title {
    font-size: 12px;
    color: #666;
    margin-bottom: 8px;
    padding-left: 4px;
    font-weight: 500;
  }
  .stencil-category-items {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    justify-content: center;
  }
  
  .stencil-node {
    width: 100%;
    height: 60px;
    border: 1px solid;
    border-radius: 6px;
    display: flex;
    flex-direction: column;
    font-size: 12px;
    cursor: move;
    transition: all 0.2s;
    overflow: hidden;
    /* 扁平化设计 - 减小阴影 */
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
    margin-bottom: 10px;
  }
  
  .stencil-node:hover {
    /* 扁平化设计 - 简化悬停效果 */
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
    transform: translateY(-1px);
  }
  
  .stencil-node-header {
    display: flex;
    align-items: center;
    padding: 2px 6px;
    /* 扁平化设计 - 实线而非虚线 */
    border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  }
  
  .stencil-node-icon {
    width: 16px;
    height: 16px;
    border-radius: 3px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 10px;
    margin-right: 4px;
    color: white;
    font-weight: bold;
  }
  
  .stencil-node-type {
    color: rgba(0, 0, 0, 0.65);
    padding: 2px 0;
    font-size: 10px;
    font-weight: 500;
    flex: 1;
  }
  
  .stencil-node-name {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 4px;
    font-weight: 500;
    color: rgba(0, 0, 0, 0.85);
  }
}




.loading-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px 0;
  color: #666;
  font-size: 13px;
}

.loading-spinner {
  width: 24px;
  height: 24px;
  border: 2px solid rgba(24, 144, 255, 0.1);
  border-left-color: #1890ff;
  border-radius: 50%;
  margin-bottom: 8px;
  animation: spin 1s linear infinite;
}

.retry-button {
  margin-top: 8px;
  padding: 4px 12px;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }

  100% {
    transform: rotate(360deg);
  }
}
</style>
