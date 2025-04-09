import apiClient from './api-client';

// 工作流数据接口
export interface WorkflowData {
  name: string;
  description: string;
  status: string;
  id?: string;
  // 可以根据需要添加更多属性
}

// 工作流服务
export const workflowService = {
  // 保存工作流
  saveWorkflow(workflowData: WorkflowData) {
    return apiClient.post<WorkflowData>('/workflows/save', workflowData);
  },

  // 获取工作流列表
  getWorkflows() {
    return apiClient.get<WorkflowData[]>('/workflows');
  },

  // 获取单个工作流详情
  getWorkflowById(id: string) {
    return apiClient.get<WorkflowData>(`/workflows/${id}`);
  },

  // 删除工作流
  deleteWorkflow(id: string) {
    return apiClient.delete(`/workflows/${id}`);
  }
};