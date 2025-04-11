import Swal from 'sweetalert2';
import { createApp, defineComponent, h } from 'vue';
import { Toast } from './toast';

// 成功提示 - 改为使用Toast
export const showSuccess = (message: string) => {
  return Toast.success(message);
};

// 错误提示 - 改为使用Toast
export const showError = (message: string) => {
  return Toast.error(message);
};

// 警告提示 - 改为使用Toast
export const showWarning = (message: string) => {
  return Toast.warning(message);
};

// 普通信息提示 - 改为使用Toast
export const showInfo = (message: string) => {
  return Toast.info(message);
};

// 确认对话框 - 保留SweetAlert2实现
export const showConfirm = (title: string, text: string, confirmButtonText = '确定') => {
  return Swal.fire({
    title,
    text,
    icon: 'question',
    showCancelButton: true,
    confirmButtonColor: '#1890ff',
    cancelButtonColor: '#d9d9d9',
    confirmButtonText,
    cancelButtonText: '取消'
  });
};

// 
export const showDialog = (title: string, html: string, confirmButtonText = '确定') => {
  return Swal.fire({
    title,
    html,
    // showCancelButton: true,
    confirmButtonColor: '#1890ff',
    cancelButtonColor: '#d9d9d9',
    confirmButtonText,
  });
};

// 加载中 - 改为使用Toast
export const showLoading = (title = '处理中...') => {
  return Toast.showLoading(title);
};

// 关闭加载 - 改为使用Toast
export const closeLoading = () => {
  Toast.closeLoading();
};

// 带导航功能的对话框
export const showNavigateDialog = (title: string, text: string, url: string, buttonText = '前往查看') => {
  return Swal.fire({
    title,
    text,
    icon: 'info',
    showCancelButton: true,
    confirmButtonColor: '#1890ff',
    cancelButtonColor: '#d9d9d9',
    confirmButtonText: buttonText,
    cancelButtonText: '取消'
  }).then((result) => {
    if (result.isConfirmed) {
      window.location.href = url;
    }
    return result;
  });
};

// 带有 Vue 组件渲染能力的对话框
export const showComponentDialog = (title: string, component: any, props: any = {}, options: any = {}) => {
  // 创建挂载点
  const mountPoint = document.createElement('div');
  document.body.appendChild(mountPoint);

  // 创建包装组件
  const wrapper = defineComponent({
    render() {
      return h(component, props);
    }
  });

  // 创建应用并挂载
  const app = createApp(wrapper);
  app.mount(mountPoint);

  // 配置 SweetAlert2 选项
  const swalOptions = {
    title,
    html: mountPoint,
    showConfirmButton: true,
    confirmButtonText: options.confirmButtonText || '关闭',
    confirmButtonColor: options.confirmButtonColor || '#1890ff',
    width: options.width || '600px',
    ...options,
    didClose: () => {
      // 在对话框关闭时卸载组件并清理
      app.unmount();
      if (options.didClose) options.didClose();
    }
  };

  // 显示对话框
  return Swal.fire(swalOptions);
};
