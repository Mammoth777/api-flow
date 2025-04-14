<template>
  <div class="workflow-stats-container">
    <div class="back-button">
      <button @click="goBack">
        <i class="arrow-left"></i> 返回列表
      </button>
    </div>
    
    <div class="stats-header">
      <h1>工作流调用统计</h1>
      <div class="workflow-info">
        <div class="info-item">
          <span class="label">ID:</span>
          <span class="value">{{ workflowId }}</span>
        </div>
        <div class="info-item">
          <span class="label">名称:</span>
          <span class="value">{{ workflowName }}</span>
        </div>
        <div class="info-item">
          <span class="label">状态:</span>
          <span class="value status-badge" :class="workflowStatus">{{ workflowStatus }}</span>
        </div>
      </div>
    </div>

    <div class="curl-section">
      <h2>API 调用方式</h2>
      <HighlightCode 
        :code="curlCommand" 
        language="bash" 
        :lineNumbers="false"
      />
    </div>

    <div v-if="isLoading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>加载数据中...</p>
    </div>
    
    <div v-else-if="error" class="error-state">
      <p>{{ error }}</p>
      <button @click="fetchStatistics" class="retry-button">重试</button>
    </div>

    <div v-else class="stats-content">
      <div class="stats-summary">
        <div class="stat-card">
          <div class="stat-number">{{ statistics.totalCalls || 0 }}</div>
          <div class="stat-label">总调用次数</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">{{ (statistics.successRate * 100).toFixed(0) || '0' }}%</div>
          <div class="stat-label">成功率</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">{{ statistics.avgDuration || '0' }}ms</div>
          <div class="stat-label">平均响应时间</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">{{ statistics.todayCalls || 0 }}</div>
          <div class="stat-label">今日调用</div>
        </div>
      </div>

      <div class="execution-history">
        <h2>调用历史记录</h2>
        <table class="history-table">
          <thead>
            <tr>
              <th>调用时间</th>
              <th>状态</th>
              <th>耗时</th>
              <th>输入参数</th>
              <th>返回结果</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="executions.length === 0">
              <td colspan="5" class="no-data">暂无调用记录</td>
            </tr>
            <tr v-for="(execution, index) in executions" :key="index">
              <td>{{ formatDate(execution.timestamp) }}</td>
              <td>
                <span class="status-indicator" :class="execution.status === 2 ? 'success' : 'error'">
                  <i class="status-icon"></i>
                  {{ execution.status === 2 ? '成功' : '失败' }}
                </span>
              </td>
              <td>{{ execution.duration }}ms</td>
              <td class="code-cell">
                <div class="code-preview">{{ formatJson(execution.input) }}</div>
                <button @click="showFullData(execution.input)" class="view-btn">查看</button>
              </td>
              <td class="code-cell">
                <div class="code-preview">{{ formatJson(execution.output) }}</div>
                <button @click="showFullData(execution.output)" class="view-btn">查看</button>
              </td>
            </tr>
          </tbody>
        </table>
        
        <div class="pagination">
          <span>共 {{ total }} 条记录</span>
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
            </select>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { workflowService } from '../services/workflow.service';
import { showComponentDialog } from '../utils/alert'; // 保留需要的对话框
import HighlightCode from '../components/HighlightCode.vue';
import { serverUrl } from '../env';

const route = useRoute();
const router = useRouter();
const workflowId = computed(() => Number(route.params.id));
const workflowName = ref('');
const workflowStatus = ref('published');
const isLoading = ref(true);
const error = ref<string | null>(null);

// 构建curl命令
const baseUrl = serverUrl;
const curlCommand = computed(() => `curl -X POST "${baseUrl}/api/workflows/execute" \\
  -H "Content-Type: application/json" \\
  -d '{
    "workflowId": ${workflowId.value},
    "sync": true,
    "inputs": {
        "content": "hello execute flow"
    }
  }'`);

// 统计数据
const statistics = ref<any>({
  totalCalls: 0,
  successRate: '--',
  avgDuration: '--',
  todayCalls: 0
});

// 执行历史相关状态
const executions = ref<any[]>([]);
const total = ref(0);
const currentPage = ref(1);
const pageSize = ref(10);
const totalPages = computed(() => Math.ceil(total.value / pageSize.value) || 1);

// 返回上一页
const goBack = () => {
  router.push('/workflows');
};

// 获取工作流统计数据
const fetchStatistics = async () => {
  isLoading.value = true;
  error.value = null;
  
  try {
    // 获取工作流详情
    const workflowDetail = await workflowService.getWorkflow(workflowId.value);
    workflowName.value = workflowDetail.name;
    workflowStatus.value = workflowDetail.status === 1 ? 'published' : 'ready';
    
    // 获取工作流执行历史，带分页参数
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value
    };
    
    const historyData = await workflowService.getWorkflowExecutionHistory(workflowId.value, params);
    
    // 转换执行历史数据格式
    executions.value = (historyData.data || []).map((item: any) => ({
      timestamp: item.createdAt,
      status: item.status,
      duration: item.duration,
      input: item.inputs || {},
      output: item.results || {}
    }));
    
    // 计算统计信息（可能需要单独的API调用）
    const statsData = historyData.statistics || {};
    // 设置总记录数
    statistics.value = {
      totalCalls: statsData.total,
      successRate: statsData.total ? statsData.successCount / statsData.total : 0, // %
      avgDuration: statsData.avgDuration, // ms
      todayCalls: statsData.todayCount
    };
    
    console.log('工作流统计:', statistics.value);
    console.log('执行历史:', executions.value);
  } catch (err: any) {
    console.error('获取工作流统计失败:', err);
    error.value = '获取统计数据失败，请重试';
    executions.value = [];
  } finally {
    isLoading.value = false;
  }
};

// 切换页码
const changePage = (newPage: number) => {
  if (newPage < 1 || newPage > totalPages.value) return;
  currentPage.value = newPage;
  fetchStatistics();
};

// 修改每页数量
const handlePageSizeChange = () => {
  currentPage.value = 1; // 切换每页条数时，回到第一页
  fetchStatistics();
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
    minute: '2-digit',
    second: '2-digit'
  });
};

// 格式化JSON显示
const formatJson = (obj: any) => {
  if (!obj) return '{}';
  try {
    const str = JSON.stringify(obj);
    return str.length > 30 ? str.substring(0, 30) + '...' : str;
  } catch (e) {
    return '{}';
  }
};

// 显示完整数据
const showFullData = (data: string) => {
  const obj = JSON.parse(data)
  const formattedData = JSON.stringify(obj, null, 2);
  
  // 使用新的支持组件渲染的对话框函数
  showComponentDialog(
    '数据详情',
    HighlightCode,
    {
      code: formattedData,
      language: 'json',
      lineNumbers: true,
      showHeader: true
    },
    {
      width: '800px',
    }
  );
};

// 组件挂载时获取数据
onMounted(() => {
  fetchStatistics();
});
</script>

<style scoped>
.workflow-stats-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px 20px;
}

.back-button {
  padding: 20px 0 10px;
  margin-bottom: 10px;
}

.stats-header {
  margin-bottom: 24px;
  overflow: visible;
}

.curl-section {
  margin-bottom: 30px;
  background-color: white;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 20px;
  overflow: visible;
}

.stats-content {
  display: flex;
  flex-direction: column;
  gap: 30px;
  overflow: visible;
}

/* 确保表格内容在需要时可滚动 */
.history-table {
  width: 100%;
  border-collapse: collapse;
  table-layout: fixed;
}

.code-cell {
  position: relative;
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 响应式设计调整 */
@media screen and (max-width: 768px) {
  .workflow-stats-container {
    padding: 0 10px 10px;
  }
  
  .stats-summary {
    grid-template-columns: repeat(2, 1fr); /* 在小屏幕上每行显示2个卡片 */
  }
}

.back-button button {
  display: flex;
  align-items: center;
  background: none;
  border: none;
  color: #1890ff;
  cursor: pointer;
  font-size: 14px;
  padding: 0;
}

.arrow-left {
  display: inline-block;
  width: 0;
  height: 0;
  border-top: 6px solid transparent;
  border-bottom: 6px solid transparent;
  border-right: 6px solid #1890ff;
  margin-right: 8px;
}

.stats-header h1 {
  font-size: 24px;
  margin-bottom: 16px;
}

.workflow-info {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.info-item {
  display: flex;
  align-items: center;
}

.info-item .label {
  font-weight: bold;
  margin-right: 8px;
  color: #666;
}

.status-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 12px;
}

.status-badge.published {
  background-color: #e6fffb;
  color: #13c2c2;
}

.status-badge.ready {
  background-color: #e6f7ff;
  color: #1890ff;
}

.curl-section h2 {
  font-size: 18px;
  margin-bottom: 16px;
  color: #333;
}

.stats-summary {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
}

.stat-card {
  background: white;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 20px;
  text-align: center;
}

.stat-number {
  font-size: 28px;
  font-weight: bold;
  color: #1890ff;
  margin-bottom: 10px;
}

.stat-label {
  color: #666;
  font-size: 14px;
}

.execution-history {
  background: white;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.execution-history h2 {
  font-size: 18px;
  margin-bottom: 16px;
  color: #333;
}

.history-table th, 
.history-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #f0f0f0;
}

.history-table th {
  background-color: #fafafa;
  font-weight: 500;
}

.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0 0;
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

.view-btn {
  position: absolute;
  right: 8px;
  top: 50%;
  transform: translateY(-50%);
  background-color: #f0f0f0;
  border: none;
  border-radius: 3px;
  padding: 2px 6px;
  font-size: 12px;
  cursor: pointer;
  color: #666;
}

.no-data {
  text-align: center;
  color: #999;
  padding: 20px 0;
}

.loading-state, .error-state {
  text-align: center;
  padding: 40px 0;
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

/* 代码高亮样式 */
:deep(.cmd-keyword) {
  color: #569cd6;
  font-weight: bold;
}

:deep(.cmd-param) {
  color: #9cdcfe;
}

:deep(.cmd-value) {
  color: #ce9178;
}

:deep(.cmd-string) {
  color: #6a9955;
}

:deep(.json-viewer) {
  background-color: #f5f5f5;
  padding: 16px;
  font-family: monospace;
  white-space: pre-wrap;
  word-break: break-word;
  max-height: 400px;
  overflow-y: auto;
  border-radius: 4px;
}

:deep(.json-detail-dialog) {
  width: 90%;
  max-width: 800px;
  max-height: 80vh;
  overflow: auto;
}

/* 修改根元素样式确保滚动正常工作 */
:deep(html), :deep(body) {
  height: 100%;
  overflow-y: auto;
}

/* 确保页面容器可以撑开并滚动 */
:deep(#app) {
  min-height: 100%;
  display: flex;
  flex-direction: column;
}

.status-indicator {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap;
}

.status-indicator.success {
  background-color: #f6ffed;
  color: #52c41a;
  border: 1px solid #b7eb8f;
}

.status-indicator.error {
  background-color: #fff2f0;
  color: #ff4d4f;
  border: 1px solid #ffccc7;
}

.status-icon {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin-right: 6px;
}

.status-indicator.success .status-icon {
  background-color: #52c41a;
}

.status-indicator.error .status-icon {
  background-color: #ff4d4f;
}
</style>
