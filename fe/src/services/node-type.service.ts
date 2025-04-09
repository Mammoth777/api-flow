import apiClient from './api-client';

// 节点类型接口，匹配API返回的格式
export interface NodeType {
  id: number;
  createAt: string;
  updateAt: string;
  deleteAt: string | null;
  code: string;
  name: string;
  description: string;
  // 客户端额外属性，根据需要添加
  color?: string;
  category?: string;
}

// 节点类型服务
export const nodeTypeService = {
  // 获取所有节点类型
  getNodeTypes() {
    return apiClient.get<NodeType[]>('/node-types');
  }
};