import { getNodeCategory, getNodeColor } from '../components/workflow';
import apiClient from './api-client';

type ParamDefine = {
  field: string;
  type: string;
  desc: string;
  default: any;
}

// 节点类型接口，匹配API返回的格式
export interface NodeType {
  id: number;
  createAt: string;
  updateAt: string;
  deleteAt: string | null;
  code: string;
  name: string;
  description: string;
  ui: {
    color: string;
    category: string;
  },
  input: ParamDefine[]
  output: ParamDefine[]
}

// 节点类型服务
export const nodeTypeService = {
  // 获取所有节点类型
  async getNodeTypes() {
    const res = await apiClient.get<NodeType[]>('/node-types');
    // 为节点类型添加客户端属性
    if (Array.isArray(res)) {
      return res.map(nodeType => ({
        ...nodeType,
        ui: {
          // 根据节点类型代码分配颜色和类别
          color: getNodeColor(nodeType.code),
          category: getNodeCategory(nodeType.code)
        }
      } as NodeType));
    } else {
      console.error('API返回的节点类型不是数组:', res);
      return [];
    }
  }
};