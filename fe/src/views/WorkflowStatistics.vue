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
      <div class="curl-command">
        <pre v-html="formattedCurlCommand"></pre>
        <button @click="copyCommand" class="copy-btn">复制</button>
      </div>
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
          <div class="stat-number">{{ statistics.successRate || '0%' }}</div>
          <div class="stat-label">成功率</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">{{ statistics.avgDuration || '0ms' }}</div>
          <div class="stat-label">平均响应时间</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">{{ statistics.todayCalls || 0 }}</div>
          <div class="stat-label">今日调用</div>
        </div>
      </div>

      <div class="time-distribution">
        <h2>调用时间分布</h2>
        <div class="chart-container">
          <!-- 这里可以放置图表组件，如 ECharts -->
          <div class="placeholder-chart">时间分布图表将在这里显示</div>
        </div>
      </div>

      <div class="recent-calls">
        <h2>最近调用记录</h2>
        <table class="calls-table">
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
            <tr v-if="statistics.recentCalls && statistics.recentCalls.length === 0">
              <td colspan="5" class="no-data">暂无调用记录</td>
            </tr>
            <tr v-for="(call, index) in statistics.recentCalls" :key="index">
              <td>{{ formatDate(call.timestamp) }}</td>
              <td>
                <span class="status-indicator" :class="call.status === 'success' ? 'success' : 'error'">
                  {{ call.status === 'success' ? '成功' : '失败' }}
                </span>
              </td>
              <td>{{ call.duration }}ms</td>
              <td class="code-cell">
                <div class="code-preview">{{ formatJson(call.input) }}</div>
                <button @click="showFullData(call.input)" class="view-btn">查看</button>
              </td>
              <td class="code-cell">
                <div class="code-preview">{{ formatJson(call.output) }}</div>
                <button @click="showFullData(call.output)" class="view-btn">查看</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { workflowService } from '../services/workflow.service';
import { showSuccess, showError, showDialog } from '../utils/alert';

const route = useRoute();
const router = useRouter();
const workflowId = computed(() => route.params.id as string);
const workflowName = ref('');
const workflowStatus = ref('published');
const isLoading = ref(true);
const error = ref<string | null>(null);

// 构建curl命令
const baseUrl = 'http://example-server.com';
const curlCommand = computed(() => `curl -X POST "${baseUrl}/api/workflow/execute" -H "Content-Type: application/json" -d '{
    "workflowId": ${workflowId.value},
    "sync": true,
    "inputs": {
        "content": "hello execute flow"
    }
}'`);

// 带有高亮的格式化curl命令
const formattedCurlCommand = computed(() => {
  return `<span class="cmd-keyword">curl</span> <span class="cmd-param">-X</span> <span class="cmd-value">POST</span> <span class="cmd-string">"${baseUrl}/api/workflow/${workflowId.value}/invoke"</span> \\
  <span class="cmd-param">-H</span> <span class="cmd-string">"Content-Type: application/json"</span> \\
  <span class="cmd-param">-d</span> <span class="cmd-string">'{
    "workflowId": ${workflowId.value},
    "sync": true,
    "inputs": {
        "content": "hello execute flow"
    }
  }'</span>`;
});

// 统计数据
const statistics = ref<any>({
  totalCalls: 0,
  successRate: '0%',
  avgDuration: '0ms',
  todayCalls: 0,
  recentCalls: []
});

// 返回上一页
const goBack = () => {
  router.push('/workflow/list');
};

// 复制命令
const copyCommand = () => {
  navigator.clipboard.writeText(curlCommand.value)
    .then(() => showSuccess('命令已复制到剪贴板'))
    .catch(err => showError('复制失败: ' + err));
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
    
    // 这里应该调用获取统计数据的API，目前使用模拟数据
    // const statsData = await workflowService.getWorkflowStatistics(workflowId.value);
    
    // 模拟数据
    statistics.value = {
      totalCalls: 128,
      successRate: '92%',
      avgDuration: '245ms',
      todayCalls: 15,
      recentCalls: [
        {
          timestamp: new Date().toISOString(),
          status: 'success',
          duration: 230,
          input: { data: { query: "如何使用API工作流?" } },
          output: { result: "API工作流是一种自动化工具，可以帮助您...", status: "success" }
        },
        {
          timestamp: new Date(Date.now() - 3600000).toISOString(),
          status: 'error',
          duration: 456,
          input: { data: { query: "错误的查询" } },
          output: { error: "处理请求时出错", status: "error" }
        },
        {
          timestamp: new Date(Date.now() - 7200000).toISOString(),
          status: 'success',
          duration: 189,
          input: { data: { query: "工作流示例" } },
          output: { result: "以下是一些常见的工作流示例...", status: "success" }
        }
      ]
    };
    
    console.log('工作流统计:', statistics.value);
  } catch (err: any) {
    console.error('获取工作流统计失败:', err);
    error.value = '获取统计数据失败，请重试';
  } finally {
    isLoading.value = false;
  }
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
const showFullData = (data: any) => {
  const formattedData = JSON.stringify(data, null, 2);
  showDialog(
    '数据详情', 
    `<pre class="json-viewer">${formattedData}</pre>`, 
    '关闭'
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
  padding: 20px;
  margin: 0 auto;
}

.back-button {
  margin-bottom: 20px;
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

.stats-header {
  margin-bottom: 24px;
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

.curl-section {
  margin-bottom: 30px;
  background-color: white;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.curl-section h2 {
  font-size: 18px;
  margin-bottom: 16px;
  color: #333;
}

.curl-command {
  background-color: #282c34;
  border-radius: 6px;
  padding: 16px;
  position: relative;
  overflow-x: hidden;
}

.curl-command pre {
  white-space: pre-wrap;
  word-break: break-word;
  margin: 0;
  padding-bottom: 40px;
  font-family: 'Menlo', 'Monaco', 'Consolas', monospace;
  font-size: 14px;
  color: #d4d4d4;
  line-height: 1.5;
}

.copy-btn {
  position: absolute;
  right: 16px;
  bottom: 16px;
  background-color: #1890ff;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 6px 12px;
  cursor: pointer;
  font-size: 12px;
  transition: background-color 0.3s;
}

.copy-btn:hover {
  background-color: #40a9ff;
}

.stats-content {
  display: flex;
  flex-direction: column;
  gap: 30px;
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

.time-distribution, .recent-calls {
  background: white;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.time-distribution h2, .recent-calls h2 {
  font-size: 18px;
  margin-bottom: 16px;
  color: #333;
}

.chart-container {
  height: 300px;
}

.placeholder-chart {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f5f5;
  color: #999;
  border-radius: 4px;
}

.calls-table {
  width: 100%;
  border-collapse: collapse;
}

.calls-table th, 
.calls-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #f0f0f0;
}

.calls-table th {
  background-color: #fafafa;
  font-weight: 500;
}

.status-indicator {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 12px;
}

.status-indicator.success {
  background-color: #f6ffed;
  color: #52c41a;
}

.status-indicator.error {
  background-color: #fff1f0;
  color: #f5222d;
}

.code-cell {
  position: relative;
  max-width: 200px;
}

.code-preview {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-family: monospace;
  padding-right: 40px;
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
</style>
