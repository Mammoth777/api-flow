import Toastify from 'toastify-js';
import 'toastify-js/src/toastify.css';

/**
 * Toast通知工具类
 */
export class Toast {
  // 成功提示
  static success(message: string, duration = 3000) {
    return Toastify({
      text: message,
      duration,
      gravity: "top", 
      position: "center",
      style: {
        background: "linear-gradient(to right, #11998e, #38ef7d)", // 降低饱和度的绿色
        borderRadius: "4px",
        boxShadow: "0 3px 10px rgba(82, 196, 26, 0.15)",
        padding: "12px 16px"
      }
    }).showToast();
  }

  // 错误提示
  static error(message: string, duration = 4000) {
    return Toastify({
      text: message,
      duration,
      gravity: "top", 
      position: "center",
      style: {
        background: "linear-gradient(to right, #FF416C, #FF4B2B)", // 降低饱和度的红色
        borderRadius: "4px",
        boxShadow: "0 3px 10px rgba(255, 77, 79, 0.15)",
        padding: "12px 16px"
      }
    }).showToast();
  }

  // 警告提示
  static warning(message: string, duration = 3500) {
    return Toastify({
      text: message,
      duration,
      gravity: "top", 
      position: "center",
      style: {
        background: "linear-gradient(to right, #f7ff00, #db36a4)", // 降低饱和度的黄色
        borderRadius: "4px",
        boxShadow: "0 3px 10px rgba(250, 173, 20, 0.15)",
        padding: "12px 16px",
        color: "#5c4f17" // 深色文字以确保在浅色背景上的可读性
      }
    }).showToast();
  }

  // 信息提示
  static info(message: string, duration = 3000) {
    return Toastify({
      text: message,
      duration,
      gravity: "top", 
      position: "center",
      style: {
        background: "linear-gradient(to right, #bdc3c7, #2c3e50)", // 降低饱和度的蓝色
        borderRadius: "4px",
        boxShadow: "0 3px 10px rgba(24, 144, 255, 0.15)",
        padding: "12px 16px"
      }
    }).showToast();
  }

  // 加载提示实例
  private static loadingInstance: any = null;

  // 显示加载提示
  static showLoading(message = "加载中...") {
    // 先清除之前的加载提示
    this.closeLoading();
    
    this.loadingInstance = Toastify({
      text: message,
      duration: 0, // 持续显示，直到手动关闭
      gravity: "top",
      position: "center",
      style: {
        background: "linear-gradient(to right, #40a9ffaa, #bae7ffaa)", // 降低饱和度的蓝色
        borderRadius: "4px",
        boxShadow: "0 3px 10px rgba(24, 144, 255, 0.15)",
        padding: "12px 16px"
      }
    }).showToast();
    
    return this.loadingInstance;
  }
  
  // 关闭加载提示
  static closeLoading() {
    if (this.loadingInstance) {
      this.loadingInstance.hideToast();
      this.loadingInstance = null;
    }
  }
}
