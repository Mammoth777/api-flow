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
          :initialDisplayMode="inspectorMode"
          :getGraphBounding="getGraphBounding"
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
import { onMounted, ref, onUnmounted, watch } from 'vue'
import WorkflowInfo from './WorkflowInfo.vue'
import NodeInspector from './NodeInspector.vue'
import NodeTypePanel from './NodeTypePanel.vue'
import { workflowService } from '../services/workflow.service'
import type { NodeType } from '../services/node-type.service'
import { useRoute } from 'vue-router'
import { Toast } from '../utils/toast' // 导入新的Toast工具类
import { registerCustomNodes } from './workflow'
import router from '../router'

// 获取路由对象和当前路由名称
const route = useRoute();
const workflowId = ref<number>(Number(route.params.id) || 0);
const isCreateMode = ref(route.name === 'createWorkflow');

watch(() => route.params.id, (newId) => {
  if (newId && newId !== workflowId.value) {
    workflowId.value = Number(newId);
    console.log('路由参数变化，重新加载工作流:', newId);
    loadWorkflow(newId);
  }
});

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
// 选中的边
const selectedEdge = ref<any>(null);
// 节点检查器显示模式
const inspectorMode = ref<'compact' | 'full'>('compact');
// 单击计时器，用于区分单击和双击
const clickTimer = ref<number | null>(null);
// 单击延迟时间（毫秒）
const clickDelay = 300;

// 清除选中节点和边（安全处理）
const clearSelectedNode = () => {
  // 如果有选中的节点，先恢复其样式
  if (selectedNode.value) {
    // 移除高亮效果
    selectedNode.value.attr('body/strokeWidth', 1);
    selectedNode.value.attr('body/stroke', '#5F95FF');
  }
  selectedNode.value = null;

  // 清除选中的边
  clearSelectedEdge();
};

// 清除选中的边
const clearSelectedEdge = () => {
  if (selectedEdge.value) {
    // 恢复边的样式
    selectedEdge.value.attr('line/stroke', '#5F95FF');
    selectedEdge.value.attr('line/strokeWidth', 2);
  }
  selectedEdge.value = null;
};

// 删除节点
const handleDeleteNode = (node: Cell) => {
  if (!node) return;

  // 如果正在编辑节点， 则不需要删除
  if (inspectorMode.value === 'full') {
    return;
  }

  // 记录节点数据，以便可能需要的后续操作
  const nodeData = node.getData();
  console.log(`正在删除节点:`, nodeData);

  // 从图中移除节点
  node.remove();

  // 清除选中状态
  selectedNode.value = null;

  Toast.success('节点已删除');
};

// 删除边
const handleDeleteEdge = (edge: Cell) => {
  if (!edge) return;

  console.log(`正在删除边:`, edge.id);

  // 从图中移除边
  edge.remove();

  // 清除选中状态
  selectedEdge.value = null;

  Toast.success('连接已删除');
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
};

// 保存工作流
const handleSave = async () => {
  if (isSaving.value) return;

  try {
    isSaving.value = true;
    Toast.showLoading('正在保存...');

    // 收集节点信息
    const nodes: any[] = [];
    const graphNodes = graph?.getNodes() || [];

    for (const node of graphNodes) {
      const data = node.getData() || {};
      const pos = node.getPosition()
      console.log(pos, 'pos')
      console.log(data, 'data');
      nodes.push({
        id: data.id,
        nodeKey: data.nodeKey,
        nodeType: data.nodeType,
        name: data.name,
        status: data.status,
        description: data.description,
        config: data.config || {},
        ui: {
          ...data.ui,
          x: pos.x,  // 明确保存x, y坐标
          y: pos.y
        }
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
        console.log(edge, 'edge')
        edges.push({
          id: edge.id,
          sourceNodeKey: sourceData.nodeKey,
          targetNodeKey: targetData.nodeKey,
          config: {}
        });
      }
    }

    const workflowData = {
      id: workflowId.value || null,
      name: workflowName.value,
      description: workflowDesc.value,
      status: 0,
      nodes: nodes,
      edges: edges
    };

    console.log('准备保存的工作流数据:', workflowData);
    const result = await workflowService.saveWorkflow(workflowData);
    console.log('工作流保存成功:', result);

    // 提示保存成功
    Toast.success('工作流保存成功!');
    router.push({ path: `/workflow/edit/${result.id}`, params: { id: result.id } });
  } catch (error: any) {
    console.error('保存工作流时出错:', error);
    Toast.closeLoading();
    // 显示错误信息
    const errorMessage = error.response?.data?.message || '保存失败';
    Toast.error(`保存失败: ${errorMessage}`);
  } finally {
    Toast.closeLoading();
    isSaving.value = false;
  }
};

// 图容器的引用
const graphContainer = ref<HTMLElement | null>(null)
let graph: Graph | null = null

function getGraphBounding() {
  if (!graphContainer.value) return { x: 0, y: 0 }
  const rect = graphContainer.value.getBoundingClientRect()
  return {
    x: rect.left + window.scrollX,
    y: rect.top + window.scrollY,
    w: rect.width,
    h: rect.height,
  }
}

// 处理键盘删除事件
const handleKeyDown = (e: KeyboardEvent) => {
  // 检查是否按下了删除键或退格键
  if ((e.key === 'Delete' || e.key === 'Backspace') && graph) {
    if (selectedNode.value) {
      handleDeleteNode(selectedNode.value);
    } else if (selectedEdge.value) {
      handleDeleteEdge(selectedEdge.value);
    }
  }
};

// 添加事件监听
const setupGraphEvents = () => {
  if (!graph) return;

  // 监听节点单击和双击事件
  graph.on('node:click', ({ node }) => {
    // 清除选中的边
    clearSelectedEdge();
    
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
    
    // 使用定时器区分单击和双击
    if (clickTimer.value !== null) {
      // 这是双击事件，清除定时器
      clearTimeout(clickTimer.value);
      clickTimer.value = null;
      
      // 双击时显示完整详情模式
      inspectorMode.value = 'full';
      console.log('节点双击: 显示完整详情');
    } else {
      // 这是单击事件，设置定时器
      clickTimer.value = setTimeout(() => {
        // 单击时显示紧凑模式
        inspectorMode.value = 'compact';
        clickTimer.value = null;
      }, clickDelay) as unknown as number;
    }
  });

  // 监听边选择事件
  graph.on('edge:click', ({ edge }) => {
    // 清除选中的节点
    clearSelectedNode();

    // 如果之前有选中的边，恢复其样式
    if (selectedEdge.value && selectedEdge.value !== edge) {
      selectedEdge.value.attr('line/stroke', '#5F95FF');
      selectedEdge.value.attr('line/strokeWidth', 2);
    }

    // 设置新选中的边
    selectedEdge.value = edge;

    // 给选中的边添加高亮效果
    edge.attr('line/stroke', '#FF3366'); // 红色
    edge.attr('line/strokeWidth', 3);

    console.log('边被选中:', edge.id);
  });

  // 监听画布空白区域点击，取消选择
  graph.on('blank:click', () => {
    clearSelectedNode();
    clearSelectedEdge();
  });

  // 控制连接桩显示/隐藏
  const showPorts = (ports: NodeListOf<SVGElement>, show: boolean) => {
    for (let i = 0, len = ports.length; i < len; i += 1) {
      ports[i].style.visibility = show ? 'visible' : 'hidden'
    }
  }
  graph.on('node:mouseenter', () => {
    const container = graphContainer.value!
    const ports = container.querySelectorAll(
      '.x6-port-body',
    ) as NodeListOf<SVGElement>
    showPorts(ports, true)
  })
  graph.on('node:mouseleave', () => {
    const container = graphContainer.value!
    const ports = container.querySelectorAll(
      '.x6-port-body',
    ) as NodeListOf<SVGElement>
    showPorts(ports, false)
  })
};

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
const loadWorkflow = async (id: number) => {
  try {
    // isSaving.value = true;
    const result = await workflowService.getWorkflow(id);
    const workflowData = result

    // 更新工作流基本信息
    workflowName.value = workflowData.name || '未命名工作流';
    workflowDesc.value = workflowData.description || '';

    if (!graph) {
      console.error('画布未初始化');
      Toast.closeLoading();
      Toast.error('画布未初始化，无法加载工作流');
      return;
    }

    // 清空现有画布内容
    graph.clearCells();

    // 临时存储节点映射，用于创建边
    const nodeMap = new Map();

    // 如果有节点数据，加载节点
    if (workflowData.nodes && Array.isArray(workflowData.nodes) && workflowData.nodes.length > 0) {
      console.log('加载工作流节点:', workflowData.nodes.length);

      // 遍历节点数据并创建节点
      for (const nodeData of workflowData.nodes) {
        // 查找对应的节点类型定义
        // const nodeTypeInfo = nodeTypes.value.find(type => type.code === nodeData.nodeType);

        // 创建节点 - 使用间隔布局，后续可优化为保存的实际位置
        const node = graph.addNode({
          id: nodeData.id,
          shape: nodeData.ui.category || 'task', // 使用nodeType作为节点形状
          x: nodeData.ui.x, // 简单的网格布局
          y: nodeData.ui.y,
          attrs: {
            body: {
              fill: nodeData.ui.color,
              stroke: '#5F95FF',
              strokeWidth: 1,
            },
            label: {
              text: nodeData.name || '未命名节点',
              fontSize: 14,
            },
          },
          data: {
            id: nodeData.id,
            nodeKey: nodeData.nodeKey,
            nodeType: nodeData.nodeType,
            name: nodeData.name,
            description: nodeData.description,
            config: nodeData.config || {},
            status: nodeData.status,
            ui: nodeData.ui || {}  // 保存完整的UI配置信息
          },
        });

        // 将节点存入映射表，以便后续创建边
        nodeMap.set(nodeData.nodeKey, node);
      }

      // 如果有边数据，创建连接
      if (workflowData.edges && Array.isArray(workflowData.edges) && workflowData.edges.length > 0) {
        console.log('加载工作流连接:', workflowData.edges.length);

        for (const edgeData of workflowData.edges) {
          const sourceNode = nodeMap.get(edgeData.sourceNodeKey);
          const targetNode = nodeMap.get(edgeData.targetNodeKey);

          if (sourceNode && targetNode) {
            // 创建边
            graph.addEdge({
              id: edgeData.id,
              source: { cell: sourceNode },
              target: { cell: targetNode },
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
              data: edgeData.config || {},
            });
          }
        }
      }

      // 调整画布，使所有内容可见
      graph.centerContent();
    }
  } catch (error) {
    console.error('加载工作流失败:', error);
    Toast.error('加载工作流数据失败');
  } finally {
    // isSaving.value = false;

  }
};

// 清理函数
const cleanup = () => {
  // 移除事件监听
  if (graph) {
    graph.off('node:click');
    graph.off('edge:click');
    graph.off('blank:click');
  }

  // 移除键盘事件监听
  window.removeEventListener('keydown', handleKeyDown);

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

    // 添加键盘事件监听
    window.addEventListener('keydown', handleKeyDown);

    // 判断是创建模式还是编辑模式
    if (isCreateMode.value) {
      console.log('创建工作流模式，清空画布');
      // 对于创建模式，确保初始化一个空画布
      if (graph) {
        graph.clearCells();
        // 重置工作流名称和描述为默认值
        workflowName.value = '新工作流';
        workflowDesc.value = '';
      }
    } else if (workflowId.value) {
      // 编辑模式，加载现有工作流
      console.log('编辑工作流模式，加载工作流:', workflowId.value);
      loadWorkflow(workflowId.value);
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
  height: calc(100% - 1px);
  border-top: 1px solid #e8e8e8;
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