import { Graph } from '@antv/x6'

const ports = {
  groups: {
    right: {
      position: 'right',
      attrs: {
        circle: {
          r: 4,
          magnet: true,
          stroke: '#5F95FF',
          strokeWidth: 1,
          fill: '#fff',
          style: {
            visibility: 'hidden',
          },
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
          style: {
            visibility: 'hidden',
          },
        },
      },
    },
  },
  items: [
    {
      group: 'right',
    },
    {
      group: 'left',
    },
  ],
}

// 用于跟踪已注册的节点类型，避免重复注册
const registeredNodeTypes = new Set<string>();

// 注册自定义节点
export const registerCustomNodes = () => {
  // 注册任务节点前检查是否已注册
  if (!registeredNodeTypes.has('task')) {
    Graph.registerNode('task', {
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
      ports: { ...ports },
    });
    // 标记该节点类型已注册
    registeredNodeTypes.add('task');
  }

  if (!registeredNodeTypes.has('system')) {
    Graph.registerNode('system', {
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
          text: '系统节点',
          fill: '#333',
          fontSize: 14,
          fontWeight: 'bold',
        },
      },
      ports: { ...ports },
    });
    registeredNodeTypes.add('system');
  }
}

// 根据节点代码获取类别
export const getNodeCategory = (code: string): string => {
  const categoryMap: Record<string, string> = {
    'api': 'task',
    'text': 'task',
    'execInput': 'system'
  };

  if (!categoryMap[code]) {
    throw new Error(`节点代码 "${code}" 未定义, 请在 src/components/workflow.ts 中添加`);
  }

  return categoryMap[code];
};

// 根据节点代码获取颜色
export const getNodeColor = (code: string): string => {
  // 扩展颜色映射表，使用更加丰富的颜色
  const colorMap: Record<string, string> = {
    'api': '#91d5ff',      // 蓝色
    'text': '#b7eb8f',     // 绿色
    'start': '#87e8de',    // 青色
    'end': '#ffadd2',      // 粉色
    'condition': '#d3adf7', // 紫色
    'loop': '#ffd591',     // 橙色
    'timer': '#ffe58f',    // 黄色
    'data': '#87e8de',     // 青色
    'function': '#adc6ff', // 蓝紫色
    'notification': '#ffadd2', // 粉色
    'log': '#b5f5ec',      // 浅绿色
    'error': '#ffa39e',    // 红色
    'success': '#b7eb8f',  // 绿色
    'wait': '#91caff',     // 蓝色
    'request': '#95de64',  // 黄绿色
    'response': '#ff9c6e', // 橘色
    'transform': '#9254de', // 紫色
    'process': '#36cfc9',  // 青绿色
    'decision': '#ff7a45', // 橙红色
  };

  return colorMap[code.toLowerCase()] || '#f0f2f5'; // 默认浅灰色
};