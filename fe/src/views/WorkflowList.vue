<template>
  <div class="workflow-list-container">
    <div class="workflow-list-header">
      <h1>工作流列表</h1>
      <button class="create-btn" @click="createWorkflow">新建</button>
    </div>
    
    <div v-if="isLoading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>
    
    <div v-else-if="error" class="error-state">
      <p>{{ error }}</p>
      <button @click="fetchWorkflows" class="retry-button">重试</button>
    </div>
    
    <div v-else class="workflow-table-container">
      <table class="workflow-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>名称</th>
            <th>描述</th>
            <th>状态</th>
            <th>更新时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="workflow in workflows" :key="workflow.id">
            <td>{{ workflow.id }}</td>
            <td>{{ workflow.name }}</td>
            <td class="description-cell">{{ workflow.description }}</td>
            <td>
              <span class="status-badge" :class="workflow.statusText">{{ workflow.statusText }}</span>
            </td>
            <td>{{ formatDate(workflow.updatedAt) }}</td>
            <td>
              <button class="edit-btn" @click="editWorkflow(workflow.id)">编辑</button>
              <button class="delete-btn" @click="confirmDelete(workflow.id)">删除</button>
              <!-- 根据状态显示不同按钮 -->
              <button v-if="workflow.status === 1" class="view-btn" @click="viewInvocation(workflow.id)">查看调用</button>
              <button v-else class="publish-btn" @click="publishIt(workflow.id)">发布</button>
            </td>
          </tr>
        </tbody>
      </table>
      
      <div class="pagination">
        <span>共 {{ total }} 条</span>
        <div class="pagination-controls">
          <button 
            class="page-btn" 
            :disabled="currentPage <= 1" 
            @click="changePage(currentPage - 1)"
          >
            上一页
          </button>
          
          <span class="page-info">{{ currentPage }} / {{ totalPages }}</span>
          
          <button 
            class="page-btn" 
            :disabled="currentPage >= totalPages" 
            @click="changePage(currentPage + 1)"
          >
            下一页
          </button>
          
          <select v-model="pageSize" class="page-size-select" @change="handlePageSizeChange">
            <option :value="10">10条/页</option>
            <option :value="20">20条/页</option>
            <option :value="50">50条/页</option>
            <option :value="100">100条/页</option>
          </select>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { workflowService } from '../services/workflow.service';
import { showConfirm, showSuccess, showError, closeLoading } from '../utils/alert';

const router = useRouter();
const workflows = ref<any[]>([]);
const isLoading = ref(true);
const error = ref<string | null>(null);
const total = ref(0);

// 分页相关状态
const currentPage = ref(1);
const pageSize = ref(10);
const totalPages = computed(() => Math.ceil(total.value / pageSize.value) || 1);

// 获取工作流列表
const fetchWorkflows = async () => {
  isLoading.value = true;
  error.value = null;
  
  try {
    // 添加分页参数
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value
    };
    
    const result = await workflowService.getWorkflows(params);
    result.data.forEach((workflow: any) => {
      workflow.statusText = getStatusText(workflow.status);
    })
    workflows.value = result.data || [];
    total.value = result.total || 0;
    
    console.log('工作流列表:', workflows.value);
  } catch (err: any) {
    console.error('获取工作流列表失败:', err);
    error.value = '获取工作流列表失败，请重试';
  } finally {
    isLoading.value = false;
  }
};

// 切换页码
const changePage = (newPage: number) => {
  if (newPage < 1 || newPage > totalPages.value) return;
  currentPage.value = newPage;
  fetchWorkflows();
};

// 修改每页数量
const handlePageSizeChange = () => {
  currentPage.value = 1; // 切换每页条数时，回到第一页
  fetchWorkflows();
};

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return '未知';
  const date = new Date(dateString);
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  });
};

// 获取状态文本
const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    0: 'ready',
    1: 'published'
  };
  
  return statusMap[status] || '未知';
};

// 新建工作流
const createWorkflow = () => {
  router.push('/workflow/create');
};

// 编辑工作流
const editWorkflow = (id: number) => {
  router.push(`/workflow/edit/${id}`);
};

// 确认删除
const confirmDelete = async (id: number) => {
  const result = await showConfirm('确认删除', '确定要删除该工作流吗？此操作无法撤销。', '删除');
  
  if (result.isConfirmed) {
    try {
      await workflowService.deleteWorkflow(id.toString());
      closeLoading();
      
      // 显示成功消息
      await showSuccess('删除成功！');
      
      // 删除成功后，重新获取列表
      fetchWorkflows();
    } catch (err: any) {
      closeLoading();
      showError('删除失败：' + (err.message || '未知错误'));
    }
  }
};

const publishIt = async (id: number) => {
  const result = await showConfirm('确认发布', '确定要发布该工作流吗？', '发布');
  
  if (result.isConfirmed) {
    try {
      await workflowService.publishWorkflow(id.toString());
      closeLoading();
      
      // 显示成功消息
      await showSuccess('发布成功！');
      
      // 发布成功后，重新获取列表
      fetchWorkflows();
    } catch (err: any) {
      closeLoading();
      showError('发布失败：' + (err.message || '未知错误'));
    }
  }
};

// 查看工作流调用统计
const viewInvocation = (id: number) => {
  router.push(`/workflow/statistics/${id}`);
};

// 组件挂载时获取工作流列表
onMounted(() => {
  fetchWorkflows();
});
</script>

<style scoped>
.workflow-list-container {
  max-width: 1200px;
  padding: 20px;
  margin: 0 auto;
}

.workflow-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.workflow-list-header h1 {
  font-size: 20px;
  margin: 0;
}

.create-btn {
  background-color: #1890ff;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 8px 16px;
  cursor: pointer;
}

.create-btn:hover {
  background-color: #40a9ff;
}

.workflow-table-container {
  background-color: white;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.workflow-table {
  width: 100%;
  border-collapse: collapse;
}

.workflow-table th,
.workflow-table td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #f0f0f0;
}

.workflow-table th {
  background-color: #fafafa;
  font-weight: 500;
}

.description-cell {
  max-width: 200px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.status-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 12px;
}

.status-badge.ready {
  background-color: #e6f7ff;
  color: #1890ff;
}

.status-badge.published {
  background-color: #e6fffb;
  color: #13c2c2;
}

.status-badge.error {
  background-color: #fff1f0;
  color: #f5222d;
}

.edit-btn, .delete-btn, .publish-btn, .view-btn {
  margin-right: 8px;
  padding: 4px 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.edit-btn {
  background-color: #e6f7ff;
  color: #1890ff;
}

.delete-btn {
  background-color: #fff1f0;
  color: #f5222d;
}

.publish-btn {
  background-color: #f6ffed;
  color: #52c41a;
}

.view-btn {
  background-color: #f9f0ff;
  color: #722ed1;
}

.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  color: #666;
  font-size: 14px;
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 8px;
}

.page-btn {
  padding: 4px 12px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  background-color: white;
  cursor: pointer;
  font-size: 14px;
  color: #333;
}

.page-btn:hover:not(:disabled) {
  color: #1890ff;
  border-color: #1890ff;
}

.page-btn:disabled {
  cursor: not-allowed;
  color: #d9d9d9;
  background-color: #f5f5f5;
}

.page-info {
  margin: 0 8px;
}

.page-size-select {
  margin-left: 16px;
  padding: 4px 8px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  background-color: white;
  cursor: pointer;
}

.loading-state, .error-state {
  text-align: center;
  padding: 40px;
  color: #666;
}

.loading-spinner {
  width: 30px;
  height: 30px;
  border: 3px solid rgba(24, 144, 255, 0.1);
  border-radius: 50%;
  border-top-color: #1890ff;
  margin: 0 auto 16px;
  animation: spin 1s linear infinite;
}

.retry-button {
  background-color: #1890ff;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 8px 16px;
  cursor: pointer;
  margin-top: 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>
