import axios from 'axios';

// 创建一个 axios 实例
const apiClient = axios.create({
  baseURL: '/api',  // API 的基础路径
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  },
  timeout: 10000 // 请求超时时间
});

// 请求拦截器
apiClient.interceptors.request.use(
  (config) => {
    // 在请求发送前可以做一些处理，如添加认证令牌等
    return config;
  },
  (error) => {
    // 请求错误处理
    return Promise.reject(error);
  }
);

// 响应拦截器
apiClient.interceptors.response.use(
  (response) => {
    // 2xx 范围内的状态码都会触发该函数
    return response.data;
  },
  (error) => {
    // 超出 2xx 范围的状态码都会触发该函数
    const errorMessage = error.response?.data?.message || '请求失败';
    console.error('请求错误:', errorMessage);
    return Promise.reject(error);
  }
);

export default apiClient;