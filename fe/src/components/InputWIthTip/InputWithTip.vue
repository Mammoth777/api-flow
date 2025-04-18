<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { get } from 'lodash-es'
import { CursorPop } from './CursorPop'


const floating = ref<HTMLDivElement | null>(null);
const inputarea = ref<HTMLDivElement | null>(null);

const TipContent = {
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
}

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
  // console.log(expression, 'expression')
  const exp = expression.replace(/^\$/, '')
  const expList = exp.split('.')
  // 已输入完成的 properties, 对应的真实值
  const prevProps = expList.slice(0, expList.length - 1)
  const prevPropsStr = prevProps.join('.')
  const fullPropsValue = get(TipContent, exp)
  const prevPropsValue = get(TipContent, prevPropsStr)
  const value = fullPropsValue || prevPropsValue
  // 正在输入中的 property
  const lastProperty = expList[expList.length - 1]
  console.log({ exp, prevProps, prevPropsStr, value, lastProperty })
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
    // 1. 返回顶层
    return calcList(TipContent)
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
    // 先检查这是普通文本还是已经在标签内
    // 获取从 $ 开始的部分
    let tagContent = text.substring(lastSign);
    
    // 处理标签中的路径
    // 如果在文本中找到数据路径属性标签，则使用其data-path属性值
    const tempDiv = document.createElement('div');
    tempDiv.innerHTML = text;
    const tagItems = tempDiv.querySelectorAll('.tag-item');
    
    // 如果有标签，可能需要获取最后一个标签的路径
    if (tagItems.length > 0) {
      const lastTag = tagItems[tagItems.length - 1];
      const path = lastTag.getAttribute('data-path');
      
      // 检查标签后是否有点号
      const tagText = lastTag.textContent || '';
      const afterTagText = text.substring(text.indexOf(tagText) + tagText.length);
      
      if (afterTagText.includes('.')) {
        // 如果标签后有点号，则使用完整路径
        tagContent = (path || '') + afterTagText;
      } else if (path) {
        // 如果标签后没有点号，就只用标签路径
        tagContent = path;
      }
    }
    
    return tagContent;
  }
  return '';
}


let popper!: CursorPop

onMounted(() => {
  const inputElm = inputarea.value!;
  const floatElm = floating.value!;
  popper = new CursorPop(inputElm, floatElm)
  
  // 处理输入事件，显示提示
  inputElm.addEventListener('input', (e: Event) => {
    // 1. 生成dropdown
    const beforeText = popper.getTextBeforeCursor();
    const tagContent = calcTagContent(beforeText)
    
    // 检查是否输入了点号，如果是则显示提示
    const lastChar = beforeText.charAt(beforeText.length - 1);
    if (lastChar === '.') {
      // 用户输入了点号，应该显示提示
      const dropdownContent = getDropdownContent(tagContent)
      dropdownList.value = dropdownContent
      
      if (dropdownContent.length > 0) {
        setTimeout(() => {
          popper.show()
          popper.updatePosition();
        }, 10)
      }
    } else {
      // 常规情况
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

    // 2. 生成tag
    const inputEvt = e as InputEvent
    console.log(inputEvt.data, 'inputEvt')
  });

  // 处理标签点击事件，确保标签的正确交互
  inputElm.addEventListener('click', (e: MouseEvent) => {
    // 获取点击位置
    const target = e.target as HTMLElement;
    
    // 如果点击了标签，浏览器会自动处理光标位置
    if (target === inputElm) {
      // 找到点击位置最近的文本位置
      const selection = window.getSelection();
      if (selection) {
        // 获取点击的准确位置
        const range = document.caretRangeFromPoint(e.clientX, e.clientY);
        if (range) {
          // 应用新的选择范围
          selection.removeAllRanges();
          selection.addRange(range);
          e.preventDefault(); // 防止默认的点击行为干扰我们的位置设置
        }
      }
    }
  });

  // 改用keydown事件以获得更好的按键捕获
  inputElm.addEventListener('keydown', (e: KeyboardEvent) => {
    // 对于上下键，需要阻止默认行为，防止光标移动
    if ((e.key === 'ArrowUp' || e.key === 'ArrowDown' || e.key === 'Enter' || e.key === 'Tab') && popper.isShown) {
      e.preventDefault();
      e.stopPropagation();
    }
    
    // 处理右箭头键，使光标能够从标签内移出
    if (e.key === 'ArrowRight') {
      console.log('right arrow')
      popper.cursorRight()
    }

    if (e.key === 'Enter' || e.key === 'Tab') {
      // Enter或Tab - 选中当前高亮项
      if (popper.isShown) {
        const selected = popper.confirmSelection();
        if (selected) {
          e.preventDefault();
          // 插入文本后隐藏下拉框
          popper.hide();
          return;
        }
        
        // 如果没有选中的项但有显示的选项，则选择第一个
        if (popper.getItemCount() > 0 && popper.getActiveItemIndex() === -1) {
          popper.highlightItem(0);
          setTimeout(() => {
            popper.confirmSelection();
            popper.hide();
          }, 0);
          e.preventDefault();
        }
      }
    } else if (e.key === 'Escape') {
      // ESC - 隐藏下拉菜单
      if (popper.isShown) {
        popper.hide();
      }
    } else if (e.key === '.') {
      // 检测点号输入，可能需要触发提示
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
      // 上箭头 - 选择上一项
      if (popper.isShown) {
        popper.selectPrevItem();
      }
    } else if (e.key === 'ArrowDown') {
      // 下箭头 - 选择下一项
      if (popper.isShown) {
        popper.selectNextItem();
      }
    } else if (e.key === 'Backspace') {
      // Backspace
      console.log('Backspace pressed');
    }
  });

  // const tipElm = floating.value!;
  // tipElm.addEventListener('click', () => {
  //   setTimeout(() => {
  //     popper.hide()
  //   }, 0)
  // });
})

// 鼠标悬停在下拉选项上时高亮显示
function highlightItem(index: number) {
  if (popper && popper.isShown) {
    popper.highlightItem(index);
  }
}

function handleItemClick(item: [string, string]) {
  // 记录当前选中项的索引
  const index = dropdownList.value.findIndex(i => i[0] === item[0]);
  if (index >= 0) {
    popper.highlightItem(index);
  }
  
  // 确保输入区域获得焦点
  inputarea.value?.focus();
  
  // 使用requestAnimationFrame确保DOM更新后再插入文本
  popper.confirmSelectionWithValue(item[0]);
  
  // 隐藏下拉菜单
  popper.hide();
}

</script>

<template>
  <div class="input-width-tip-wrapper">
    <div contenteditable="true" class="input-with-tip" ref="inputarea"></div>
    <div ref="floating" class="floating-tip">
      <span 
        v-for="(item, index) in dropdownList" 
        :key="index"
        @click="handleItemClick(item)"
        @mouseover="highlightItem(index)"
        class="dropdown-item"
        role="option"
        :aria-selected="false"
        :data-key="item[0]"
      >
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
  height: 40px;
  outline: none;
  transition: border-color 0.3s;
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
  background: rgba(250, 250, 250, 0.96); /* 更柔和的背景色 */
  color: #333;
  padding: 6px 0;
  border-radius: 14px; /* 更圆润的边角 */
  font-size: 14px;
  min-width: 180px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.08), 0 2px 5px rgba(0, 0, 0, 0.05); /* 更自然的阴影 */
  backdrop-filter: blur(12px); /* 增强模糊效果 */
  border: none; /* 移除边框 */
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

/* 高亮项样式 */
.floating-tip .dropdown-item.active {
  background-color: rgba(81, 91, 212, 0.08); /* 更现代的主题紫色，透明度低 */
  color: #515bd4; /* Instagram 风格紫色调 */
  font-weight: 500;
}

.floating-tip .dropdown-item.active .item-key {
  color: #515bd4; /* 统一主题色 */
}

.floating-tip .dropdown-item.active .item-type {
  color: #515bd4;
  opacity: 0.75;
}

.floating-tip .item-key {
  color: #262626; /* 更深的文本颜色，提高可读性 */
  margin-right: 8px;
  font-weight: 500;
}

.floating-tip .item-type {
  color: #737373; /* 更柔和的灰色 */
  font-weight: normal;
  font-size: 13px;
}

/* 标签统一样式 */
:deep(.tag-item) {
  display: inline-flex;
  background-color: rgba(81, 91, 212, 0.08); /* 匹配主题色 */
  color: #515bd4; /* Instagram 风格紫色 */
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

/* 光标容器样式 */
:deep(.cursor-container) {
  display: inline;
  position: relative;
  margin-left: 1px;
  color: transparent;
  caret-color: black; /* 使光标可见 */
}

/* 确保标签后面的点号操作正常 */
:deep(.tag-item + br) {
  display: none;
}
</style>
