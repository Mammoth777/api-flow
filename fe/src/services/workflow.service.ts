import axios from 'axios';

// 定义工作流接口
export interface Workflow {
  id: string;
  name: string;
  description: string;
  status: string;
  createAt: string;
  updateAt: string;
  nodes: any[];
  edges: any[];
}

// 工作流服务
export const workflowService = {
  // 获取工作流列表
  getWorkflows: async (params?: { page?: number; pageSize?: number }) => {
    const queryParams = new URLSearchParams();
    
    if (params?.page) {
      queryParams.append('page', params.page.toString());
    }
    
    if (params?.pageSize) {
      queryParams.append('pageSize', params.pageSize.toString());
    }
    
    const queryString = queryParams.toString();
    const url = `/api/workflows${queryString ? '?' + queryString : ''}`;
    
    const response = await axios.get(url);
    return response.data;
  },
  
  // 获取单个工作流详情
  getWorkflow: async (id: number) => {
    const response = await axios.get(`/api/workflows/${id}`);
    return response.data;
  },
  
  // 保存/更新工作流
  saveWorkflow: async (data: any) => {
    const response = await axios.post('/api/workflows/save', data);
    return response.data;
  },
  
  // 删除工作流
  deleteWorkflow: async (id: string) => {
    const response = await axios.delete(`/api/workflows/${id}`);
    return response.data;
  },

  // 发布工作流
  publishWorkflow: async (id: string) => {
    const response = await axios.post(`/api/workflows/${id}/publish`);
    return response.data;
  },

  // 获取工作流执行历史记录
  getWorkflowExecutionHistory: async (workflowId: number, params: any = {}): Promise<any> => {
    try {
      const response = await axios.get(`/api/workflows/execute/${workflowId}/history`, { params });
      return response.data;
    } catch (error) {
      console.error('获取工作流执行历史失败:', error);
      throw error;
    }
  },
};