<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue';

interface Props {
  tips: string[];
  position: { x: number; y: number };
  query: string; // 当前输入的查询文本
}

const props = defineProps<Props>();
const emit = defineEmits(['select-tip', 'close-tip']);

const activeIndex = ref(0);
const dropdownElement = ref<HTMLElement | null>(null);

// 当提示列表变化时，重置激活索引
watch(() => props.tips, () => {
  activeIndex.value = 0;
});

// 使激活项滚动到可视区域
const scrollActiveItemIntoView = () => {
  nextTick(() => {
    if (!dropdownElement.value) return;
    
    const activeItem = dropdownElement.value.querySelector('.tip-item.active');
    if (activeItem) {
      activeItem.scrollIntoView({ block: 'nearest', behavior: 'smooth' });
    }
  });
};

// 处理键盘导航
const handleKeydown = (e: KeyboardEvent) => {
  if (!props.tips.length) return;

  switch (e.key) {
    case 'ArrowDown':
      e.preventDefault();
      activeIndex.value = (activeIndex.value + 1) % props.tips.length;
      scrollActiveItemIntoView();
      break;
    case 'ArrowUp':
      e.preventDefault();
      activeIndex.value = (activeIndex.value - 1 + props.tips.length) % props.tips.length;
      scrollActiveItemIntoView();
      break;
    case 'Enter':
      e.preventDefault();
      if (props.tips[activeIndex.value]) {
        selectTip(props.tips[activeIndex.value]);
      }
      break;
    case 'Tab':
      e.preventDefault();
      if (props.tips[activeIndex.value]) {
        selectTip(props.tips[activeIndex.value]);
      }
      break;
    case 'Escape':
      e.preventDefault();
      // 通知父组件关闭提示
      emit('close-tip');
      break;
  }
};

// 选择提示项
const selectTip = (tip: string) => {
  emit('select-tip', tip);
};

// 点击外部关闭下拉菜单
const handleClickOutside = (e: MouseEvent) => {
  if (dropdownElement.value && !dropdownElement.value.contains(e.target as Node)) {
    emit('close-tip');
  }
};

// 添加和移除全局键盘事件监听
onMounted(() => {
  window.addEventListener('keydown', handleKeydown);
  window.addEventListener('mousedown', handleClickOutside);
});

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown);
  window.removeEventListener('mousedown', handleClickOutside);
});

// 高亮匹配的文本部分
const highlightMatch = (text: string) => {
  if (!props.query) return text;
  
  const index = text.toLowerCase().indexOf(props.query.toLowerCase());
  if (index === -1) return text;
  
  const beforeMatch = text.substring(0, index);
  const match = text.substring(index, index + props.query.length);
  const afterMatch = text.substring(index + props.query.length);
  
  return [
    beforeMatch,
    `<span class="highlight">${match}</span>`,
    afterMatch
  ].join('');
};

// 检测是否为对象类型的提示
const isObjectTip = (tip: string) => {
  return tip.includes(',');
};
</script>

<template>
  <div 
    ref="dropdownElement"
    class="tip-dropdown" 
    :style="{
      left: `${position.x}px`,
      top: `${position.y}px`
    }"
  >
    <div class="tips-container">
      <div
        v-for="(tip, index) in tips"
        :key="index"
        class="tip-item"
        :class="{ 'active': index === activeIndex }"
        @click="selectTip(tip)"
        @mouseover="activeIndex = index"
      >
        <div class="tip-content" v-html="highlightMatch(tip)"></div>
        <span v-if="isObjectTip(tip)" class="object-indicator">对象</span>
      </div>
      <div v-if="tips.length === 0" class="no-tips">
        无匹配提示
      </div>
    </div>
  </div>
</template>

<style scoped>
.tip-dropdown {
  position: absolute;
  background-color: white;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  z-index: 9999;
  max-width: 320px;
  max-height: 200px;
  overflow-y: auto;
}

.tips-container {
  padding: 5px 0;
}

.tip-item {
  padding: 8px 15px;
  cursor: pointer;
  font-size: 14px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tip-item:hover,
.tip-item.active {
  background-color: #f5f7fa;
  color: #409eff;
}

.tip-content {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
}

.object-indicator {
  font-size: 12px;
  color: #909399;
  background-color: #f0f2f5;
  padding: 2px 5px;
  border-radius: 3px;
  margin-left: 8px;
  flex-shrink: 0;
}

.no-tips {
  padding: 10px 15px;
  color: #909399;
  font-size: 14px;
  text-align: center;
}

:deep(.highlight) {
  color: #409eff;
  font-weight: bold;
}
</style>