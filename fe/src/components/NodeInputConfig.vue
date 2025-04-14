<template>
  <div class="node-input-config">
    <div class="input-field" v-for="(field, index) in inputFields" :key="index">
      <div class="field-header">
        <label>{{ field.desc }} ({{ field.field }})</label>
        <span class="field-type">{{ field.type }}</span>
      </div>
      
      <!-- 根据字段类型显示不同的输入控件 -->
      <!-- 字符串类型输入 -->
      <input 
        v-if="field.type === 'string'" 
        type="text" 
        v-model="nodeConfig[field.field]" 
        @input="debouncedUpdate"
        :placeholder="field.desc" 
      />
      
      <!-- 数字类型输入 -->
      <input 
        v-else-if="field.type === 'number'" 
        type="number" 
        v-model.number="nodeConfig[field.field]" 
        @input="debouncedUpdate"
        :placeholder="field.desc" 
      />
      
      <!-- 选项类型输入 -->
      <select 
        v-else-if="field.type === 'options'" 
        v-model="nodeConfig[field.field]" 
        @change="updateConfig"
      >
        <option 
          v-for="(option, optIdx) in field.options" 
          :key="optIdx" 
          :value="option"
        >
          {{ option }}
        </option>
      </select>
      
      <!-- 对象类型输入 -->
      <div v-else-if="field.type === 'object'" class="object-input-container">
        <textarea 
          v-model="objectValues[field.field]" 
          @input="debouncedUpdateObject(field.field)"
          placeholder="输入 JSON 对象，例如: {&quot;key&quot;: &quot;value&quot;}"
          rows="3"
        ></textarea>
        <div class="error-message" v-if="objectErrors[field.field]">
          {{ objectErrors[field.field] }}
        </div>
      </div>
      
      <!-- 布尔类型输入 -->
      <div v-else-if="field.type === 'boolean'" class="boolean-input">
        <label class="switch">
          <input 
            type="checkbox" 
            v-model="nodeConfig[field.field]" 
            @change="updateConfig"
          >
          <span class="slider round"></span>
        </label>
        <span class="boolean-value">{{ nodeConfig[field.field] ? '是' : '否' }}</span>
      </div>
      
      <!-- 默认输入框（未知类型） -->
      <input 
        v-else 
        type="text" 
        v-model="nodeConfig[field.field]" 
        @input="debouncedUpdate"
        :placeholder="field.desc" 
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';

// 定义节点输入字段的类型
interface InputField {
  field: string;
  type: 'string' | 'number' | 'object' | 'options' | 'boolean';
  desc: string;
  default: any;
  // 仅在 type 为 'options' 时使用
  options?: string[] | { label: string, value: number | string}[];
}

// 接收的属性
const props = defineProps({
  // 节点的输入字段定义
  inputFields: {
    type: Array as () => InputField[],
    required: true
  },
  // 节点配置，用于双向绑定
  modelValue: {
    type: Object,
    default: () => ({})
  }
});

// 定义事件
const emit = defineEmits(['update:modelValue']);

// 局部存储配置对象
const nodeConfig = ref<Record<string, any>>({});

// 对象类型字段的临时文本值
const objectValues = ref<Record<string, string>>({});
const objectErrors = ref<Record<string, string>>({});

// 防抖定时器
const debounceTimerId = ref<number | null>(null);
const objectDebounceTimers = ref<Record<string, number>>({});

// 防抖函数，减少更新频率
const debouncedUpdate = () => {
  if (debounceTimerId.value) {
    clearTimeout(debounceTimerId.value);
  }
  
  debounceTimerId.value = setTimeout(() => {
    updateConfig();
    debounceTimerId.value = null;
  }, 300) as unknown as number;
};

// 对象字段的防抖更新
const debouncedUpdateObject = (fieldName: string) => {
  if (objectDebounceTimers.value[fieldName]) {
    clearTimeout(objectDebounceTimers.value[fieldName]);
  }
  
  objectDebounceTimers.value[fieldName] = setTimeout(() => {
    updateObjectField(fieldName);
    objectDebounceTimers.value[fieldName] = 0;
  }, 500) as unknown as number;
};

// 优化监听逻辑，减少不必要的深度比较
watch(() => props.modelValue, (newConfig) => {
  // 避免不必要的深拷贝，只更新值
  for (const key in newConfig) {
    if (nodeConfig.value[key] !== newConfig[key]) {
      nodeConfig.value[key] = newConfig[key];
    }
  }
  
  // 初始化对象类型字段的文本表示
  props.inputFields.forEach(field => {
    if (field.type === 'object' && nodeConfig.value[field.field] !== undefined) {
      try {
        const jsonString = JSON.stringify(nodeConfig.value[field.field], null, 2);
        // 只有当值变化时才更新，避免不必要的DOM更新
        if (objectValues.value[field.field] !== jsonString) {
          objectValues.value[field.field] = jsonString;
          objectErrors.value[field.field] = '';
        }
      } catch (err) {
        if (!objectValues.value[field.field]) {
          objectValues.value[field.field] = '{}';
          objectErrors.value[field.field] = '';
        }
      }
    }
  });
}, { immediate: true });

// 组件挂载时初始化配置
onMounted(() => {
  initializeConfig();
});

// 初始化配置，确保每个字段都有默认值
const initializeConfig = () => {
  props.inputFields.forEach(field => {
    // 如果配置中没有该字段的值，使用默认值
    if (nodeConfig.value[field.field] === undefined) {
      nodeConfig.value[field.field] = field.default;
      
      // 如果是对象类型，初始化其文本表示
      if (field.type === 'object') {
        try {
          objectValues.value[field.field] = JSON.stringify(field.default || {}, null, 2);
          objectErrors.value[field.field] = '';
        } catch (err) {
          objectValues.value[field.field] = '{}';
          objectErrors.value[field.field] = '';
        }
      }
    }
  });
  
  updateConfig();
};

// 更新配置并触发事件
const updateConfig = () => {
  emit('update:modelValue', nodeConfig.value);
};

// 更新对象类型字段 - 添加错误处理
const updateObjectField = (fieldName: string) => {
  try {
    const parsedValue = JSON.parse(objectValues.value[fieldName]);
    nodeConfig.value[fieldName] = parsedValue;
    objectErrors.value[fieldName] = '';
    updateConfig();
  } catch (err) {
    objectErrors.value[fieldName] = '无效的 JSON 格式';
    // 解析失败时不触发更新
  }
};
</script>

<style scoped>
.node-input-config {
  margin-bottom: 16px;
}

.input-field {
  margin-bottom: 16px;
}

.field-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.field-header label {
  font-size: 12px;
  color: #666;
  font-weight: 500;
}

.field-type {
  font-size: 10px;
  color: #888;
  background-color: #f0f0f0;
  padding: 2px 6px;
  border-radius: 10px;
}

input, select, textarea {
  width: 100%;
  padding: 6px 8px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  font-size: 13px;
}

input:focus, select:focus, textarea:focus {
  outline: none;
  border-color: #40a9ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

.object-input-container textarea {
  font-family: monospace;
  font-size: 12px;
  resize: vertical;
}

.error-message {
  color: #f5222d;
  font-size: 12px;
  margin-top: 4px;
}

/* 布尔选项开关样式 */
.boolean-input {
  display: flex;
  align-items: center;
  gap: 10px;
}

.switch {
  position: relative;
  display: inline-block;
  width: 40px;
  height: 20px;
}

.switch input { 
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  transition: .4s;
}

.slider:before {
  position: absolute;
  content: "";
  height: 16px;
  width: 16px;
  left: 2px;
  bottom: 2px;
  background-color: white;
  transition: .4s;
}

input:checked + .slider {
  background-color: #1890ff;
}

input:focus + .slider {
  box-shadow: 0 0 1px #1890ff;
}

input:checked + .slider:before {
  transform: translateX(20px);
}

.slider.round {
  border-radius: 34px;
}

.slider.round:before {
  border-radius: 50%;
}

.boolean-value {
  font-size: 12px;
  color: #666;
}
</style>