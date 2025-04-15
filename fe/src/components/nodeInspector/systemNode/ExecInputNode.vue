<script setup lang="ts">
import { ref } from 'vue';

const inputParams = ref([
  {
    name: '',
    type: 'string', // 默认类型
    defaultValue: '',
    options: [] as string[], // 如果类型是 options，可以定义选项
  },
]);

const dataTypes = [
  { value: 'string', label: '文本' },
  { value: 'number', label: '数字' },
  { value: 'boolean', label: '布尔' },
  { value: 'array', label: '数组' },
  { value: 'options', label: '选项' },
  { value: 'object', label: '对象' },
  // { value: 'null', label: '空值' },
  { value: 'any', label: '任意' },
];

function addParam() {
  inputParams.value.push({
    name: '',
    type: 'string',
    defaultValue: '',
    options: [],
  });
}

function removeParam(index: number) {
  inputParams.value.splice(index, 1);
}

function handleOptionsInput(param: any, event: Event) {
  const target = event.target as HTMLInputElement | null;
  if (target) {
    param.options = target.value.split(',').map((item) => item.trim());
  }
}

function handleBooleanDefault(param: any, event: Event) {
  const target = event.target as HTMLInputElement | null;
  if (target) {
    param.defaultValue = target.checked;
  }
}

function handleTypeChange(param: any) {
  param.defaultValue = param.type === 'boolean' ? false : param.type === 'options' ? [] : '';
  if (param.type !== 'options') {
    param.options = [];
  }
}
</script>

<template>
  <div>
    <h3>定义工作流入参</h3>
    <table class="editable-table">
      <thead>
        <tr>
          <th>Name</th>
          <th>Type</th>
          <th>Default</th>
          <th class="delete-column">Action</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(param, index) in inputParams" :key="index">
          <td>
            <input v-model="param.name" placeholder="参数名称" />
          </td>
          <td>
            <select v-model="param.type" @change="handleTypeChange(param)">
              <option v-for="type in dataTypes" :key="type.value" :value="type.value">
                {{ type.label }}
              </option>
            </select>
          </td>
          <td>
            <template v-if="param.type === 'boolean'">
              <input
                type="checkbox"
                :checked="param.defaultValue as unknown as boolean"
                @change="handleBooleanDefault(param, $event)"
              />
            </template>
            <template v-else-if="param.type === 'options'">
              <div class="options-editor">
                <label>选项:</label>
                <textarea
                  :value="param.options.join(', ')"
                  placeholder="用逗号分隔选项，例如: A,B,C"
                  rows="3"
                  style="width: 100%; resize: vertical;"
                  @input="handleOptionsInput(param, $event)"
                ></textarea>
              </div>
            </template>
            <template v-else-if="param.type === 'string' || param.type === 'array'">
              <textarea
                v-model="param.defaultValue"
                rows="3"
                style="width: 100%; resize: vertical;"
              ></textarea>
            </template>
            <template v-else>
              <input v-model="param.defaultValue" />
            </template>
          </td>
          <td style="text-align: center">
            <button @click="removeParam(index)" class="delete-btn">删除</button>
          </td>
        </tr>
      </tbody>
    </table>
    <button @click="addParam">添加参数</button>
  </div>
</template>

<style scoped>
table.editable-table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 1rem;
  font-size: 14px;
  color: #333;
}

table.editable-table th, table.editable-table td {
  border: 1px solid #e0e0e0;
  padding: 0;
  text-align: left;
}

table.editable-table th {
  background-color: #f9f9f9;
  font-weight: 600;
  padding: 8px;
}

table.editable-table th.delete-column {
  width: 50px; /* 调整删除列的宽度 */
  text-align: center;
}

table.editable-table td input, table.editable-table td select, table.editable-table td textarea {
  width: 100%;
  height: 100%; /* 确保高度占满父元素 */
  padding: 6px;
  border: none;
  font-size: 14px;
  color: #333;
  background-color: transparent; /* 改为透明背景 */
  outline: none;
  box-sizing: border-box;
  transition: background-color 0.2s;
}

table.editable-table td input:focus, table.editable-table td select:focus, table.editable-table td textarea:focus {
  outline: none;
}

table.editable-table td textarea {
  width: 100%;
  padding: 6px;
  border: none; /* 与其他输入框保持一致 */
  font-size: 14px;
  color: #333;
  outline: none;
  box-sizing: border-box;
  resize: vertical;
  transition: background-color 0.2s;
}

table.editable-table td textarea:focus {
  background-color: #f0f8ff; /* 与其他输入框的焦点样式一致 */
  outline: none;
}

table.editable-table td:first-child {
  background-color: #f9f9f9; /* 确保 Name 列背景色统一 */
}

table.editable-table td:focus-within {
  background-color: #f0f8ff; /* 激活时改变单元格背景色 */
}

.options-editor {
  margin-top: 4px;
}

.options-editor label {
  font-size: 12px;
  color: #666;
  margin-right: 4px;
}

button {
  padding: 6px 12px;
  font-size: 14px;
  color: #007bff;
  background-color: transparent;
  border: none;
  cursor: pointer;
  transition: color 0.2s;
}

button:hover {
  color: #0056b3;
}

button:disabled {
  color: #ccc;
  cursor: not-allowed;
}

.delete-btn {
  padding: 0;
  font-size: 12px; /* 调整字体大小 */
  color: #999; /* 使用灰色淡化按钮 */
  background-color: transparent;
  border: none;
  cursor: pointer;
  transition: color 0.2s;
}

delete-btn:hover {
  color: #666; /* 悬停时稍微加深颜色 */
}
</style>
