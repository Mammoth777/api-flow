<script setup lang="ts">
import { ref } from 'vue';
// import InputWithPrevNode from './InputWithPrevNode.vue';
// import TagDisplay from './TagDisplay.vue';
import InputWithTip from './InputWithTip.vue';

// 示例输入值
const inputValue = ref('');

// 示例提示数据
const suggestions = ref([
  'user',
  'product',
  'order',
  {
    customer: {
      name: 'John',
      age: 30,
      address: {
        city: 'Beijing',
        street: 'Main Street'
      }
    }
  },
  {
    settings: {
      theme: 'dark',
      language: 'zh-cn'
    }
  },
  'date',
  'status'
]);

// 处理输入变化
// const handleInput = (value: string) => {
//   console.log('输入值变化:', value);
//   inputValue.value = value;
// };
</script>

<template>
  <div class="example-container">
    <h2>带提示的标签化输入框组件示例 (所见即所得)</h2>
    
    <div class="description">
      <p>输入特殊字符 <code>$</code> 来触发提示，然后可以使用 <code>.</code> 访问对象属性</p>
      <p>例如: <code>$customer.address.city</code> 会形成一个带样式的标签</p>
      <p>标签和普通文本可以混合输入，互不影响</p>
      <p><strong>所见即所得:</strong> 输入时直接显示标签样式，无需额外预览</p>
    </div>
    
    <div class="input-wrapper">
      <label>带标签的输入框 (所见即所得):</label>
      <InputWithTip></InputWithTip>
    </div>
    
    <div class="current-value">
      <h3>当前输入的原始值:</h3>
      <pre>{{ inputValue }}</pre>
    </div>
    
    <div class="usage-examples">
      <h3>使用示例:</h3>
      <ul>
        <li>输入 <code>我想查询$customer.name的信息</code> 并选择提示</li>
        <li>输入 <code>设置系统$settings.theme为暗色</code> 并选择提示</li>
        <li>输入 <code>已完成$order状态修改</code> 并选择提示</li>
      </ul>
    </div>
    
    <div class="available-suggestions">
      <h3>可用的提示:</h3>
      <pre>{{ JSON.stringify(suggestions, null, 2) }}</pre>
    </div>
  </div>
</template>

<style scoped>
.example-container {
  margin: 20px;
  font-family: Arial, sans-serif;
}

.description {
  margin: 20px 0;
  padding: 15px;
  background-color: #f9f9f9;
  border-radius: 5px;
  border-left: 4px solid #409eff;
}

.input-wrapper {
  margin: 20px 0;
  display: flex;
  flex-direction: column;
}

.input-wrapper label {
  margin-bottom: 5px;
  font-weight: bold;
}

.current-value {
  margin: 20px 0;
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 5px;
}

.current-value pre,
.available-suggestions pre {
  background-color: #f0f2f5;
  padding: 10px;
  border-radius: 3px;
  overflow-x: auto;
}

.usage-examples {
  margin: 20px 0;
  padding: 15px;
  background-color: #e6f7ff;
  border-radius: 5px;
  border-left: 4px solid #1890ff;
}

.usage-examples ul {
  padding-left: 20px;
}

.usage-examples li {
  margin: 8px 0;
}

.available-suggestions {
  margin: 20px 0;
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 5px;
}

code {
  background-color: #e6e6e6;
  padding: 2px 5px;
  border-radius: 3px;
  font-family: monospace;
}
</style>