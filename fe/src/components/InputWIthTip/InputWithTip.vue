<script setup lang="ts">
import { ref, onMounted, watch, type PropType } from 'vue';
import { get } from 'lodash-es'
import { CursorPop } from './CursorPop'

const floating = ref<HTMLDivElement | null>(null);
const inputarea = ref<HTMLDivElement | null>(null);

type ItemTypes = 'string' | 'number' | 'boolean' | 'object' | 'array'
type ItemObjType = {
  [key: string]: ItemTypes | ItemObjType
}

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: ''
  },
  suggestions: {
    type: Object as PropType<ItemObjType>,
    default: () => ({
      fullname: {
        firstName: 'string',
        lastName: 'string'
      },
      age: 'number',
      address: {
        city: 'string',
        state: 'string'
      },
      hobbies: 'array',
      isStudent: 'boolean',
    })
  }
})

const emit = defineEmits(['update:modelValue']);

// 初始化输入值
watch(() => props.modelValue, (newValue) => {
  if (inputarea.value && newValue !== inputarea.value.textContent) {
    inputarea.value.textContent = newValue;
  }
}, { immediate: true });

const dropdownList = ref<[string, string][]>([])

dropdownList.value = getDropdownContent('')

/**
 * 一定会有 $ 符号开头
 * 1. 根据 "." 进行分割，获取 "." 前面的所有所有属性连接后的值。 
 * 1.1 如果"."后面有key， 则比对
 * 2. 如果没有 ".", 则返回顶层
 */
function getDropdownContent(expression: string): Array<[string, string]> {
  if (!expression.startsWith('$')) {
    return []
  }
  const exp = expression.replace(/^\$/, '')
  const expList = exp.split('.')
  const prevProps = expList.slice(0, expList.length - 1)
  const prevPropsStr = prevProps.join('.')
  const fullPropsValue = get(props.suggestions, exp)
  const prevPropsValue = get(props.suggestions, prevPropsStr)
  const value = fullPropsValue || prevPropsValue
  const lastProperty = expList[expList.length - 1]
  let list: [string, string][] = []
  const calcList = (obj: any) => {
    list = Object.entries(obj).map(item => {
      const key = item[0]
      const value = item[1]
      if (typeof value === 'object') {
        return [key, 'object']
      } else {
        return [key, value]
      }
    }) as [string, string][]
    if (!lastProperty) {
      return list
    }
    return list.filter(item => {
      return item[0].startsWith(lastProperty)
    })
  }
  if (prevProps.length === 0) {
    return calcList(props.suggestions)
  } else if (typeof value === 'object') {
    return calcList(value)
  } else {
    return []
  }
}

// 增强版的calcTagContent，支持处理标签内文本
function calcTagContent(text: string) {
  const lastSign = text.lastIndexOf('$');
  if (lastSign !== -1) {
    let tagContent = text.substring(lastSign);
    const tempDiv = document.createElement('div');
    tempDiv.innerHTML = text;
    const tagItems = tempDiv.querySelectorAll('.tag-item');
    if (tagItems.length > 0) {
      const lastTag = tagItems[tagItems.length - 1];
      const path = lastTag.getAttribute('data-path');
      const tagText = lastTag.textContent || '';
      const afterTagText = text.substring(text.indexOf(tagText) + tagText.length);
      if (afterTagText.includes('.')) {
        tagContent = (path || '') + afterTagText;
      } else if (path) {
        tagContent = path;
      }
    }
    return tagContent;
  }
  return '';
}

let popper!: CursorPop

// 更新 modelValue 的函数
function updateModelValue() {
  if (inputarea.value) {
    const text = inputarea.value.textContent || '';
    emit('update:modelValue', text);
  }
}

onMounted(() => {
  const inputElm = inputarea.value!;
  const floatElm = floating.value!;
  popper = new CursorPop(inputElm, floatElm)

  // 处理输入事件，显示提示
  inputElm.addEventListener('input', () => {
    // 每次输入后更新 modelValue
    updateModelValue();
    
    const beforeText = popper.getTextBeforeCursor();
    const tagContent = calcTagContent(beforeText)
    const lastChar = beforeText.charAt(beforeText.length - 1);
    if (lastChar === '.') {
      const dropdownContent = getDropdownContent(tagContent)
      dropdownList.value = dropdownContent
      if (dropdownContent.length > 0) {
        setTimeout(() => {
          popper.show()
          popper.updatePosition();
        }, 10)
      }
    } else {
      const dropdownContent = getDropdownContent(tagContent)
      dropdownList.value = dropdownContent
      if (dropdownContent.length === 0) {
        popper.hide()
      } else {
        setTimeout(() => {
          popper.show()
          popper.updatePosition();
        }, 10)
      }
    }
  });

  // 改用keydown事件以获得更好的按键捕获
  inputElm.addEventListener('keydown', (e: KeyboardEvent) => {
    const moveOutOfTagKeys = ['ArrowRight']
    const moveDropdownHighlightKeys = ['ArrowUp', 'ArrowDown']
    const confirmSelectKeys = ['Enter', 'Tab']
    const hideMenuKeys = ['Escape']
    if (([...moveOutOfTagKeys, ...moveDropdownHighlightKeys, ...confirmSelectKeys].includes(e.key)) && popper.isShown) {
      e.preventDefault();
      e.stopPropagation();
    }
    if (moveOutOfTagKeys.includes(e.key)) {
      popper.cursorRight()
    }
    if (confirmSelectKeys.includes(e.key)) {
      if (popper.isShown) {
        const selected = popper.confirmSelection();
        if (selected) {
          e.preventDefault();
          popper.hide();
          return;
        }
        if (popper.getItemCount() > 0 && popper.getActiveItemIndex() === -1) {
          popper.highlightItem(0);
          setTimeout(() => {
            popper.confirmSelection();
            popper.hide();
          }, 0);
          e.preventDefault();
        }
      }
    } else if (hideMenuKeys.includes(e.key)) {
      if (popper.isShown) {
        popper.hide();
      }
    } else if (e.key === '.') {
      setTimeout(() => {
        const beforeText = popper.getTextBeforeCursor();
        const tagContent = calcTagContent(beforeText);
        const dropdownContent = getDropdownContent(tagContent);
        if (dropdownContent.length > 0) {
          dropdownList.value = dropdownContent;
          popper.show();
          popper.updatePosition();
        }
      }, 10);
    }
  });

  // 使用keyup事件处理上下键导航
  inputElm.addEventListener('keyup', (e: KeyboardEvent) => {
    if (e.key === 'ArrowUp') {
      if (popper.isShown) {
        popper.selectPrevItem();
      }
    } else if (e.key === 'ArrowDown') {
      if (popper.isShown) {
        popper.selectNextItem();
      }
    } else if (e.key === 'Backspace') {
      console.log('Backspace pressed');
    }
  });

  // 监听失去焦点事件，确保值被正确更新
  inputElm.addEventListener('blur', () => {
    updateModelValue();
  });
})

// 鼠标悬停在下拉选项上时高亮显示
function highlightItem(index: number) {
  if (popper && popper.isShown) {
    popper.highlightItem(index);
    popper.cloneRange();
  }
}

function handleItemClick(e: Event, item: [string, string]) {
  e.stopPropagation();
  popper.confirmSelectionWithValue(item[0]);
  setTimeout(() => {
    updateModelValue();
  }, 0);
  popper.hide();
}

</script>

<template>
  <div class="input-width-tip-wrapper">
    <div contenteditable="true" class="input-with-tip" ref="inputarea" :data-placeholder="placeholder"></div>
    <div ref="floating" class="floating-tip">
      <span v-for="(item, index) in dropdownList" :key="index" @click="e => handleItemClick(e, item)"
        @mouseover="highlightItem(index)" class="dropdown-item" role="option" :aria-selected="false"
        :data-key="item[0]">
        <span class="item-key">{{ item[0] }}</span>
        <span class="item-type">{{ item[1] }}</span>
      </span>
    </div>
  </div>
</template>

<style scoped>
.input-width-tip-wrapper {
  position: relative;
}

.input-with-tip {
  border: 1px solid #ccc;
  padding: 10px;
  border-radius: 4px;
  width: 100%;
  min-height: 40px;
  outline: none;
  transition: border-color 0.3s;
}

/* 添加占位符样式 */
.input-with-tip:empty:before {
  content: attr(data-placeholder);
  color: #999;
  position: absolute;
  pointer-events: none;
}

.input-with-tip:focus {
  border-color: #0066cc;
  box-shadow: 0 0 0 2px rgba(0, 102, 204, 0.2);
}

.cursor {
  display: inline;
  position: absolute;
}

.floating-tip {
  position: absolute;
  top: 0;
  left: 0;
  background: rgba(250, 250, 250, 0.96);
  color: #333;
  padding: 6px 0;
  border-radius: 14px;
  font-size: 14px;
  min-width: 180px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.08), 0 2px 5px rgba(0, 0, 0, 0.05);
  backdrop-filter: blur(12px);
  border: none;
  overflow: hidden;
  transform-origin: top left;
  animation: dropdownFadeIn 0.2s ease-out;
}

@keyframes dropdownFadeIn {
  from {
    opacity: 0;
    transform: scale(0.98);
  }

  to {
    opacity: 1;
    transform: scale(1);
  }
}

.floating-tip .dropdown-item {
  display: flex;
  justify-content: space-between;
  cursor: pointer;
  padding: 10px 16px;
  transition: all 0.2s cubic-bezier(0.25, 0.1, 0.25, 1);
  margin: 0 4px;
  border-radius: 8px;
}

.floating-tip .dropdown-item:hover {
  background-color: rgba(0, 0, 0, 0.04);
}

.floating-tip .dropdown-item.active {
  background-color: rgba(81, 91, 212, 0.08);
  color: #515bd4;
  font-weight: 500;
}

.floating-tip .dropdown-item.active .item-key {
  color: #515bd4;
}

.floating-tip .dropdown-item.active .item-type {
  color: #515bd4;
  opacity: 0.75;
}

.floating-tip .item-key {
  color: #262626;
  margin-right: 8px;
  font-weight: 500;
}

.floating-tip .item-type {
  color: #737373;
  font-weight: normal;
  font-size: 13px;
}

:deep(.tag-item) {
  display: inline-flex;
  background-color: rgba(81, 91, 212, 0.08);
  color: #515bd4;
  padding: 2px 6px;
  margin: 0 2px;
  font-size: 14px;
  user-select: none;
  border-radius: 6px;
  border: none;
  position: relative;
  transition: all 0.2s ease;
  font-weight: 500;
  cursor: default;
}

:deep(.tag-item:hover) {
  background-color: rgba(81, 91, 212, 0.12);
}

:deep(.cursor-container) {
  display: inline;
  position: relative;
  margin-left: 1px;
  color: transparent;
  caret-color: black;
}

:deep(.tag-item + br) {
  display: none;
}
</style>
