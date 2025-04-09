<template>
  <div class="workflow-container">
    <WorkflowInfo
      v-model:name="workflowName"
      v-model:description="workflowDesc"
      :isSaving="isSaving"
      @save="handleSave"
    />
    <div class="workflow-editor">
      <div class="stencil-container">
        <div class="stencil-title">节点类型</div>
        <div ref="stencilContainer" class="stencil">
          <!-- 节点类型加载状态显示 -->
          <div v-if="isLoadingNodeTypes" class="loading-state">
            <div class="loading-spinner"></div>
            <p>加载节点类型...</p>
          </div>
          <div v-else-if="nodeTypeError" class="error-state">
            <p>{{ nodeTypeError }}</p>
            <button @click="fetchNodeTypes" class="retry-button">重试</button>
          </div>
        </div>
      </div>
      <div class="workflow-main">
        <div ref="graphContainer" class="workflow-canvas"></div>
        <!-- 修复节点检查器条件渲染，避免使用v-if导致DOM重建问题 -->
        <NodeInspector
          v-show="selectedNode"
          :selectedNode="selectedNode"
          :nodeTypes="nodeTypes"
          @close="clearSelectedNode"
          @nodeUpdated="handleNodeUpdated"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Graph, Shape, Cell } from '@antv/x6'
import { Dnd } from '@antv/x6-plugin-dnd'
import '@antv/x6-vue-shape'
import { onMounted, ref, onUnmounted } from 'vue'
import WorkflowInfo from './WorkflowInfo.vue'
import NodeInspector from './NodeInspector.vue'
import { workflowService } from '../services/workflow.service'
import { nodeTypeService } from '../services/node-type.service'
import type { NodeType } from '../services/node-type.service'

// 工作流属性
const workflowName = ref<string>('新工作流');
const workflowDesc = ref<string>('');
const isSaving = ref<boolean>(false);

// 节点类型状态
const nodeTypes = ref<NodeType[]>([]);
const isLoadingNodeTypes = ref<boolean>(false);
const nodeTypeError = ref<string | null>(null);

// 选中的节点
const selectedNode = ref<any>(null);
// 清除选中节点（安全处理）
const clearSelectedNode = () => {
  selectedNode.value = null;
};

// 处理节点更新
const handleNodeUpdated = (node: Cell) => {
  console.log('节点已更新:', node.getData());
};

// 保存工作流
const handleSave = async () => {
  if (isSaving.value) return;
  
  try {
    isSaving.value = true;
    
    const workflowData = {
      name: workflowName.value,
      description: workflowDesc.value,
      status: "ready"
    };
    
    const result = await workflowService.saveWorkflow(workflowData);
    console.log('工作流保存成功:', result);
    
    // 提示保存成功
    alert('工作流保存成功!');
    
  } catch (error: any) {
    console.error('保存工作流时出错:', error);
    // 显示错误信息
    const errorMessage = error.response?.data?.message || '保存失败';
    alert(`保存失败: ${errorMessage}`);
  } finally {
    isSaving.value = false;
  }
};

// 获取节点类型
const fetchNodeTypes = async () => {
  isLoadingNodeTypes.value = true;
  nodeTypeError.value = null;
  
  try {
    const result = await nodeTypeService.getNodeTypes();
    console.log('API返回结果:', result);
    
    // 根据API返回的结构获取数据
    nodeTypes.value = result.data || [];
    console.log('节点类型加载成功:', nodeTypes.value);
    
    // 为节点类型添加客户端属性
    if (Array.isArray(nodeTypes.value)) {
      nodeTypes.value = nodeTypes.value.map(nodeType => ({
        ...nodeType,
        // 根据节点类型代码分配颜色和类别
        color: getNodeColor(nodeType.code),
        category: getNodeCategory(nodeType.code)
      }));
    } else {
      console.error('API返回的节点类型不是数组:', nodeTypes.value);
      nodeTypes.value = [];
    }
    
    // 确保在初始化后再调用initStencil
    // 使用nextTick确保DOM更新后再操作
    if (graph) {
      // 使用setTimeout等待DOM更新
      setTimeout(() => {
        if (stencilContainer.value) {
          initStencil();
        }
      }, 0);
    }
  } catch (error: any) {
    console.error('获取节点类型失败:', error);
    nodeTypeError.value = '获取节点类型失败，请重试';
    nodeTypes.value = []; // 确保失败时节点类型为空数组，而不是undefined
  } finally {
    isLoadingNodeTypes.value = false;
  }
};

// 根据节点代码获取颜色
const getNodeColor = (code: string): string => {
  const colorMap: Record<string, string> = {
    'api': '#e6f7ff',   // 浅蓝色
    'text': '#f6ffed',  // 浅绿色
    'start': '#e6f7ff', // 浅蓝色
    'end': '#fff1f0',   // 浅红色
    'condition': '#f9f0ff', // 浅紫色
  };
  
  return colorMap[code] || '#ffffff'; // 默认白色
};

// 根据节点代码获取类别
const getNodeCategory = (code: string): string => {
  const categoryMap: Record<string, string> = {
    'api': '常用节点',
    'text': '常用节点',
    'start': '基本节点',
    'end': '基本节点',
    'condition': '高级节点',
  };
  
  return categoryMap[code] || '其他节点';
};

// 图容器的引用
const graphContainer = ref<HTMLElement | null>(null)
const stencilContainer = ref<HTMLElement | null>(null)
let graph: Graph | null = null
let dnd: Dnd | null = null

// 添加事件监听
const setupGraphEvents = () => {
  if (!graph) return;
  
  // 监听节点选择事件
  graph.on('node:click', ({ node }) => {
    selectedNode.value = node;
  });
  
  // 监听画布空白区域点击，取消选择
  graph.on('blank:click', () => {
    clearSelectedNode();
  });
};

// 注册自定义节点
const registerCustomNodes = () => {
  // 注册任务节点
  Graph.registerNode('task-node', {
    inherit: 'rect',
    width: 120,
    height: 60,
    attrs: {
      body: {
        fill: '#ffffff',
        stroke: '#5F95FF',
        strokeWidth: 1,
        rx: 6,
        ry: 6,
      },
      label: {
        text: '任务节点',
        fill: '#333',
        fontSize: 14,
        fontWeight: 'bold',
      },
    },
    ports: {
      groups: {
        top: {
          position: 'top',
          attrs: {
            circle: {
              r: 4,
              magnet: true,
              stroke: '#5F95FF',
              strokeWidth: 1,
              fill: '#fff',
            },
          },
        },
        right: {
          position: 'right',
          attrs: {
            circle: {
              r: 4,
              magnet: true,
              stroke: '#5F95FF',
              strokeWidth: 1,
              fill: '#fff',
            },
          },
        },
        bottom: {
          position: 'bottom',
          attrs: {
            circle: {
              r: 4,
              magnet: true,
              stroke: '#5F95FF',
              strokeWidth: 1,
              fill: '#fff',
            },
          },
        },
        left: {
          position: 'left',
          attrs: {
            circle: {
              r: 4,
              magnet: true,
              stroke: '#5F95FF',
              strokeWidth: 1,
              fill: '#fff',
            },
          },
        },
      },
      items: [
        { group: 'top' },
        { group: 'right' },
        { group: 'bottom' },
        { group: 'left' },
      ],
    },
  })
}

// 初始化 Stencil 面板
const initStencil = () => {
  if (!graph || !stencilContainer.value) {
    console.warn('无法初始化节点面板：图表或容器未就绪');
    return;
  }

  try {
    // 创建 DND 实例
    dnd = new Dnd({
      target: graph,
      scaled: false,
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
    
    // 如果没有从 API 获取到节点类型，使用默认节点
    if (nodeTypes.value.length === 0) {
      // 默认节点类型 (这里我们设置一些默认值以匹配新的接口)
      const defaultNodeTypes: NodeType[] = [
        { 
          id: 1, 
          createAt: '', 
          updateAt: '', 
          deleteAt: null, 
          code: 'api', 
          name: 'API节点', 
          description: '发送HTTP请求并处理响应的节点', 
          color: '#e6f7ff', 
          category: '常用节点' 
        },
        { 
          id: 2, 
          createAt: '', 
          updateAt: '', 
          deleteAt: null, 
          code: 'text', 
          name: '文本节点', 
          description: '直接返回配置的文本内容的节点', 
          color: '#f6ffed',
          category: '常用节点' 
        }
      ];
      
      // 按分类整理
      defaultNodeTypes.forEach(nodeType => {
        const category = nodeType.category || '其他节点';
        if (!nodeTypesByCategory[category]) {
          nodeTypesByCategory[category] = [];
        }
        nodeTypesByCategory[category].push(nodeType);
      });
    } else {
      // 使用从 API 获取的节点类型
      nodeTypes.value.forEach(nodeType => {
        const category = nodeType.category || '其他节点';
        if (!nodeTypesByCategory[category]) {
          nodeTypesByCategory[category] = [];
        }
        nodeTypesByCategory[category].push(nodeType);
      });
    }
    
    // 创建分类和节点
    Object.entries(nodeTypesByCategory).forEach(([category, types]) => {
      const itemsEl = createCategory(category);
      
      types.forEach(nodeType => {
        const nodeEl = createNodeTemplate(
          nodeType.code,
          nodeType.name,
          nodeType.color
        );
        itemsEl.appendChild(nodeEl);
      });
    });

    // 为所有节点模板绑定拖拽事件
    if (stencilContainer.value) {
      const nodeEls = stencilContainer.value.querySelectorAll('.stencil-node');
      nodeEls.forEach((nodeEl) => {
        const code = nodeEl.getAttribute('data-type') || 'task-node';
        const text = nodeEl.getAttribute('data-text') || '节点';
        const color = nodeEl.getAttribute('data-color');
        
        nodeEl.addEventListener('mousedown', (e) => {
          if (!graph || !dnd) return;
          
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
          
          const node = graph.createNode({
            shape: 'task-node', // 使用通用节点形状
            attrs: nodeAttrs,
            // 存储节点类型信息以便后续使用
            data: {
              nodeType: code
            }
          });
          
          dnd.start(node, e as any);
        });
      });
    }
  } catch (error) {
    console.error('初始化节点面板时出错:', error);
  }
}

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
}

// 创建节点模板元素
const createNodeTemplate = (type: string, text: string, backgroundColor?: string) => {
  const nodeEl = document.createElement('div')
  nodeEl.className = 'stencil-node'
  nodeEl.setAttribute('data-type', type)
  nodeEl.setAttribute('data-text', text)
  if (backgroundColor) {
    nodeEl.setAttribute('data-color', backgroundColor)
    nodeEl.style.backgroundColor = backgroundColor
  }
  
  // 更新节点内部结构，添加类型标签和名称标签
  const typeLabel = document.createElement('div')
  typeLabel.className = 'stencil-node-type'
  typeLabel.textContent = type

  const nameLabel = document.createElement('div')
  nameLabel.className = 'stencil-node-name'
  nameLabel.textContent = text
  
  nodeEl.appendChild(typeLabel)
  nodeEl.appendChild(nameLabel)
  
  return nodeEl
}

// 初始化画布
const initGraph = () => {
  if (graphContainer.value) {
    try {
      // 创建画布
      graph = new Graph({
        container: graphContainer.value,
        autoResize: true, // 自动调整画布大小以适应容器
        panning: {
          enabled: true,
          modifiers: null, // 不需要按键修饰符即可拖动画布
        },
        grid: {
          size: 10,
          visible: true,
          type: 'mesh',
          args: {
            color: '#cccccc',
            thickness: 1,
          },
        },
        connecting: {
          router: 'manhattan',
          connector: {
            name: 'rounded',
            args: {
              radius: 8,
            },
          },
          anchor: 'center',
          connectionPoint: 'boundary',
          allowBlank: false,
          snap: {
            radius: 20,
          },
          createEdge() {
            return new Shape.Edge({
              attrs: {
                line: {
                  stroke: '#5F95FF',
                  strokeWidth: 2,
                  targetMarker: {
                    name: 'block',
                    width: 12,
                    height: 8,
                  },
                },
              },
              zIndex: 0,
            })
          },
          validateConnection({ sourceView, targetView }) {
            if (sourceView === targetView) {
              return false
            }
            return true
          },
        },
        mousewheel: {
          enabled: true,
          modifiers: 'ctrl',
          factor: 1.1,
          maxScale: 1.5,
          minScale: 0.5,
        },
        // Removed selecting property as it is not recognized
      });

      // 添加事件监听
      setupGraphEvents();
      
      // 初始化 Stencil 面板
      initStencil();
      
      // 居中画布内容
      graph.centerContent();
      
      console.log('画布初始化完成');
    } catch (error) {
      console.error('初始化画布时出错:', error);
    }
  }
};

// 清理函数
const cleanup = () => {
  // 移除事件监听
  if (graph) {
    graph.off('node:click');
    graph.off('blank:click');
  }
  
  // 清理引用
  dnd = null;
  graph = null;
};

// 组件挂载和卸载
onMounted(() => {
  console.log('组件挂载中...');
  
  // 先注册自定义节点
  registerCustomNodes();
  
  // 等待DOM完全渲染后再初始化图表
  setTimeout(() => {
    console.log('初始化图表...');
    if (graphContainer.value) {
      initGraph();
      
      console.log('获取节点类型...');
      setTimeout(() => {
        fetchNodeTypes();
      }, 200);
    }
  }, 100);
});

// 组件卸载时清理资源
onUnmounted(() => {
  cleanup();
});
</script>

<style scoped>
.workflow-container {
  position: relative;
  width: 100%;
  height: 100vh;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.workflow-editor {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.workflow-main {
  flex: 1;
  position: relative;
  overflow: hidden;
}

.stencil-container {
  width: 180px;
  height: 100%;
  background-color: #f7f9fb;
  border-right: 1px solid #ddd;
  overflow-y: auto;
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE and Edge */
}

.stencil-container::-webkit-scrollbar {
  display: none; /* Chrome, Safari and Opera */
}

.stencil-title {
  padding: 8px 12px;
  font-size: 14px;
  font-weight: 500;
  border-bottom: 1px solid #e8e8e8;
  background-color: #fafafa;
}

.stencil {
  padding: 8px;
}

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
}

.stencil-node {
  width: 90px;
  height: 60px;
  border: 1px solid #5F95FF;
  border-radius: 4px;
  display: flex;
  flex-direction: column;
  background-color: #fff;
  font-size: 12px;
  cursor: move;
  transition: all 0.2s;
  overflow: hidden;
}

.stencil-node:hover {
  box-shadow: 0 0 6px rgba(95, 149, 255, 0.4);
}

.stencil-node-type {
  background-color: #f0f5ff;
  color: #5F95FF;
  padding: 2px 4px;
  font-size: 10px;
  width: 100%;
  text-align: center;
  border-bottom: 1px dashed #d9e6ff;
}

.stencil-node-name {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4px;
  font-weight: 500;
}

.workflow-canvas {
  width: 100%;
  height: 100%;
  background-color: #f5f5f5;
  position: relative;
  overflow: hidden; /* 防止画布出现滚动条 */
}

html, body {
  margin: 0;
  padding: 0;
  overflow: hidden; /* 防止全局滚动条 */
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
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>