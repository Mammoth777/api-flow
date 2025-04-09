<template>
  <div class="workflow-container">
    <WorkflowInfo v-model:name="workflowName" v-model:description="workflowDesc" :isSaving="isSaving"
      @save="handleSave" />
    <div class="workflow-editor">
      <!-- 使用新的节点类型面板组件 -->
      <NodeTypePanel :graph="graph" @node-types-loaded="handleNodeTypesLoaded" />
      
      <div class="workflow-main">
        <div ref="graphContainer" class="workflow-canvas"></div>
        <NodeInspector 
          v-show="selectedNode" 
          :selectedNode="selectedNode" 
          :nodeTypes="nodeTypes"
          @close="clearSelectedNode" 
          @nodeUpdated="handleNodeUpdated" 
          @deleteNode="handleDeleteNode" 
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Graph, Shape, Cell } from '@antv/x6'
import '@antv/x6-vue-shape'
import { onMounted, ref, onUnmounted } from 'vue'
import WorkflowInfo from './WorkflowInfo.vue'
import NodeInspector from './NodeInspector.vue'
import NodeTypePanel from './NodeTypePanel.vue'
import { workflowService } from '../services/workflow.service'
import type { NodeType } from '../services/node-type.service'
import { useRoute } from 'vue-router'
import { showConfirm, showSuccess, showError, closeLoading } from '../utils/alert'

// 获取路由参数
const route = useRoute();
const workflowId = route.params.id;

// 工作流属性
const workflowName = ref<string>('新工作流');
const workflowDesc = ref<string>('');
const isSaving = ref<boolean>(false);

// 节点类型状态 - 现在从NodeTypePanel组件获取
const nodeTypes = ref<NodeType[]>([]);

// 处理NodeTypePanel组件加载的节点类型
const handleNodeTypesLoaded = (types: NodeType[]) => {
  nodeTypes.value = types;
};

// 选中的节点
const selectedNode = ref<any>(null);

// 清除选中节点（安全处理）
const clearSelectedNode = () => {
  // 如果有选中的节点，先恢复其样式
  if (selectedNode.value) {
    // 移除高亮效果
    selectedNode.value.attr('body/strokeWidth', 1);
    selectedNode.value.attr('body/stroke', '#5F95FF');
  }
  selectedNode.value = null;
};

// 删除节点
const handleDeleteNode = (node: Cell) => {
  if (!node) return;
  
  showConfirm('删除节点', '确定要删除此节点吗？', '删除').then((result) => {
    if (result.isConfirmed) {
      // 记录节点数据，以便可能需要的后续操作
      const nodeData = node.getData();
      console.log(`正在删除节点:`, nodeData);
      
      // 从图中移除节点
      node.remove();
      
      // 清除选中状态
      selectedNode.value = null;
      
      showSuccess('节点已删除');
    }
  });
};

// 设置键盘事件监听，用于删除节点
const setupKeyboardEvents = () => {
  // 获取画布容器元素
  if (!graphContainer.value) return;

};

// 处理节点更新
const handleNodeUpdated = (updatedData: any) => {
  if (!selectedNode.value) return;
  
  const node = selectedNode.value;
  const currentData = node.getData() || {};
  
  // 创建更新数据副本，排除 nodeType
  const { nodeType: updatedNodeType, ...dataToUpdate } = updatedData;
  
  // 如果尝试更新节点类型且与当前类型不同，发出警告
  if (updatedNodeType && updatedNodeType !== currentData.nodeType) {
    console.warn('节点类型不可更改');
    // 可选：提醒用户
    // alert('节点类型不可更改');
  }
  
  // 更新节点数据，保持原始 nodeType
  const newData = {
    ...currentData,
    ...dataToUpdate
  };
  
  // 将更新后的数据设置回节点
  node.setData(newData);
  
  // 如果名称发生变化，同时更新节点上显示的文本
  if (dataToUpdate.name) {
    node.attr('label/text', dataToUpdate.name);
  }
  
  console.log('节点已更新:', newData);
};

// 保存工作流
const handleSave = async () => {
  if (isSaving.value) return;

  try {
    isSaving.value = true;

    // 收集节点信息
    const nodes: any[] = [];
    const graphNodes = graph?.getNodes() || [];
    
    for (const node of graphNodes) {
      const data = node.getData() || {};
      nodes.push({
        nodeKey: data.nodeKey,
        nodeType: data.nodeType || 'unknown',
        name: data.name || node.attr('label/text') || '未命名节点',
        description: data.description || '',
        config: data.config || {},
        status: 'ready'
      });
    }

    // 收集连线信息
    const edges: any[] = [];
    const graphEdges = graph?.getEdges() || [];
    
    for (const edge of graphEdges) {
      const source = edge.getSourceNode();
      const target = edge.getTargetNode();
      
      if (source && target) {
        const sourceData = source.getData() || {};
        const targetData = target.getData() || {};
        
        edges.push({
          sourceNodeKey: sourceData.nodeKey,
          targetNodeKey: targetData.nodeKey,
          config: {}
        });
      }
    }

    const workflowData = {
      name: workflowName.value,
      description: workflowDesc.value,
      status: "ready",
      nodes: nodes,
      edges: edges
    };

    const result = await workflowService.saveWorkflow(workflowData);
    console.log('工作流保存成功:', result);

    closeLoading();
    // 提示保存成功
    await showSuccess('工作流保存成功!');

  } catch (error: any) {
    console.error('保存工作流时出错:', error);
    closeLoading();
    // 显示错误信息
    const errorMessage = error.response?.data?.message || '保存失败';
    showError(`保存失败: ${errorMessage}`);
  } finally {
    isSaving.value = false;
  }
};

// 图容器的引用
const graphContainer = ref<HTMLElement | null>(null)
let graph: Graph | null = null

// 添加事件监听
const setupGraphEvents = () => {
  if (!graph) return;

  // 监听节点选择事件
  graph.on('node:click', ({ node }) => {
    // 如果之前有选中的节点，恢复其样式
    if (selectedNode.value && selectedNode.value !== node) {
      selectedNode.value.attr('body/strokeWidth', 1);
      selectedNode.value.attr('body/stroke', '#5F95FF');
    }
    
    // 设置新选中的节点
    selectedNode.value = node;
    
    // 给选中的节点添加高亮效果
    node.attr('body/strokeWidth', 2);
    node.attr('body/stroke', '#ff7f0e'); // 橙色边框
    
    // 可选：让选中的节点在视图中居中
    // graph?.centerCell(node);
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

      // 居中画布内容
      graph.centerContent();

      console.log('画布初始化完成');
    } catch (error) {
      console.error('初始化画布时出错:', error);
    }
  }
};

// 加载工作流数据
const loadWorkflow = async (id: string) => {
  try {
    isSaving.value = true;
    const result = await workflowService.getWorkflow(id);
    const workflowData = result
    
    // 更新工作流基本信息
    workflowName.value = workflowData.name || '未命名工作流';
    workflowDesc.value = workflowData.description || '';
    
    // 如果有节点和边的数据，加载它们
    if (workflowData.nodes && Array.isArray(workflowData.nodes)) {
      console.log('加载工作流节点:', workflowData.nodes.length);
      // 加载节点和连线的逻辑...
      // 这里需要实现加载节点到画布的具体逻辑
    }
    
    console.log('工作流加载成功:', workflowData);
  } catch (error) {
    console.error('加载工作流失败:', error);
    alert('加载工作流数据失败');
  } finally {
    isSaving.value = false;
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
  graph = null;
};

// 组件挂载和卸载
onMounted(() => {
  console.log('组件挂载中...');

  // 先注册自定义节点
  registerCustomNodes();

  // 等待DOM完全渲染后再初始化图表
  console.log('初始化图表...');
  if (graphContainer.value) {
    initGraph();
    setupKeyboardEvents(); // 设置键盘事件
    
    // 如果有ID参数，加载对应的工作流
    if (workflowId) {
      console.log('加载工作流:', workflowId);
      loadWorkflow(workflowId as string);
    }
  }
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

.workflow-canvas {
  width: 100%;
  height: 100%;
  background-color: #f5f5f5;
  position: relative;
  overflow: hidden;
}

html,
body {
  margin: 0;
  padding: 0;
  overflow: hidden;
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