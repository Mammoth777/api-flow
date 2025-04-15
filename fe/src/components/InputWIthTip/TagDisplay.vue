<template>
  <div class="tag-display">
    <template v-for="(part, index) in parsedContent" :key="index">
      <span 
        v-if="part.isTag" 
        class="display-tag"
        :class="{ [`tag-type-${part.type}`]: part.type }"
      >
        {{ part.text }}
      </span>
      <span v-else>{{ part }}</span>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface Props {
  content: string;
}

const props = defineProps<Props>();

// 解析内容，识别并处理标签格式 {{TAG:text:type}}
const parsedContent = computed(() => {
  const result = [];
  let remaining = props.content;
  const tagRegex = /\{\{TAG:(.*?):(.*?)\}\}/g;
  let match;
  let lastIndex = 0;
  
  // 查找所有标签
  while ((match = tagRegex.exec(props.content)) !== null) {
    // 添加标签前的普通文本
    if (match.index > lastIndex) {
      result.push(props.content.substring(lastIndex, match.index));
    }
    
    // 添加标签对象
    result.push({
      isTag: true,
      text: match[1],  // 标签文本
      type: match[2]   // 标签类型
    });
    
    // 更新处理位置
    lastIndex = match.index + match[0].length;
  }
  
  // 添加剩余的普通文本
  if (lastIndex < props.content.length) {
    result.push(props.content.substring(lastIndex));
  }
  
  return result;
});
</script>

<style scoped>
.tag-display {
  display: inline-flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 4px;
  line-height: 1.6;
}

.display-tag {
  display: inline-flex;
  align-items: center;
  background-color: #ecf5ff;
  color: #409eff;
  border-radius: 3px;
  padding: 2px 8px;
  font-size: 0.9em;
  border: 1px solid #d9ecff;
  white-space: nowrap;
  user-select: none;
}

/* 不同类型的标签样式 */
.tag-type-default {
  background-color: #ecf5ff;
  color: #409eff;
  border-color: #d9ecff;
}

.tag-type-customer {
  background-color: #f0f9eb;
  color: #67c23a;
  border-color: #e1f3d8;
}

.tag-type-settings {
  background-color: #fdf6ec;
  color: #e6a23c;
  border-color: #faecd8;
}

/* 复杂对象路径的标签样式 */
.tag-type-customer\.address {
  background-color: #f2f6fc;
  color: #909399;
  border-color: #ebeef5;
}
</style>