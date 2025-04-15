<script setup lang="ts">
import { ref, watch, nextTick, onMounted, onBeforeUnmount } from 'vue';
import TipDropdown from './TipDrop.vue';

interface Tag {
  id: number;
  text: string; // 显示的文本
  path: string[]; // 对象路径，如 ['customer', 'address']
  type?: string; // 标签类型，用于样式区分
  fullPath: string; // 完整路径，如 customer.address.city
}

interface Props {
  modelValue: string;
  suggestions: (string | object)[];
  placeholder?: string;
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: '输入$触发提示...'
});

const emit = defineEmits(['update:modelValue']);

const editorRef = ref<HTMLDivElement | null>(null);
const showTips = ref(false);
const tipPosition = ref({ x: 0, y: 0 });
const currentQuery = ref('');
const currentObjectPath = ref<string[]>([]);
const filteredTips = ref<string[]>([]);
const tags = ref<Tag[]>([]);
const nextTagId = ref(1);
const lastCursorPosition = ref<Range | null>(null);
const isComposing = ref(false);
const isEmpty = ref(true);

// 初始化编辑器内容
const initializeContent = () => {
  if (!props.modelValue) {
    isEmpty.value = true;
    return;
  }
  
  // 解析初始内容中的标签
  const content = parseModelValue(props.modelValue);
  if (!editorRef.value) return;
  
  // 设置内容
  editorRef.value.innerHTML = content;
  isEmpty.value = content.trim() === '';
};

// 解析 modelValue 中的标签格式，转换为 HTML
const parseModelValue = (value: string): string => {
  if (!value) return '';
  
  // 解析 {{TAG:text:type}} 格式的标签
  const tagRegex = /\{\{TAG:(.*?):(.*?)\}\}/g;
  let result = value;
  let match;
  
  // 收集所有标签
  const newTags: Tag[] = [];
  
  // 替换所有标签为 HTML
  while ((match = tagRegex.exec(value)) !== null) {
    const text = match[1];
    const type = match[2];
    const fullTag = match[0];
    const pathParts = type.split('.');
    
    // 创建标签
    const tag: Tag = {
      id: nextTagId.value++,
      text,
      path: pathParts,
      type,
      fullPath: pathParts.join('.')
    };
    
    newTags.push(tag);
    
    // 创建标签的 HTML
    const tagHtml = createTagHtml(tag);
    
    // 替换原标签文本为 HTML
    result = result.replace(fullTag, tagHtml);
  }
  
  // 更新标签集合
  tags.value = newTags;
  
  return result;
};

// 创建标签的 HTML 结构
const createTagHtml = (tag: Tag): string => {
  return `<span class="editor-tag ${tag.type ? `tag-type-${tag.type}` : ''}" 
                data-id="${tag.id}" 
                data-text="${tag.text}"
                data-type="${tag.type || 'default'}"
                contenteditable="false"
                tabindex="-1">
            ${tag.text}
            <span class="remove-tag" data-id="${tag.id}" tabindex="-1">&times;</span>
          </span>`;
};

// 获取编辑器纯文本内容（不含标签HTML）
const getPlainTextContent = (): string => {
  if (!editorRef.value) return '';
  
  // 克隆节点以避免修改原始DOM
  const clone = editorRef.value.cloneNode(true) as HTMLElement;
  
  // 找到所有标签元素
  const tagElements = clone.querySelectorAll('.editor-tag');
  
  // 替换每个标签元素为占位符格式
  tagElements.forEach(el => {
    const text = el.getAttribute('data-text') || '';
    const type = el.getAttribute('data-type') || 'default';
    const placeholder = document.createTextNode(`{{TAG:${text}:${type}}}`);
    el.parentNode?.replaceChild(placeholder, el);
  });
  
  // 返回处理后的文本内容
  return clone.textContent || '';
};

// 监听 modelValue 变化
watch(() => props.modelValue, () => {
  // 只有当编辑器不是焦点时才更新内容
  if (document.activeElement !== editorRef.value) {
    initializeContent();
  }
});

// 更新模型值
const updateModelValue = () => {
  const plainText = getPlainTextContent();
  emit('update:modelValue', plainText);
};

// 焦点事件处理器
const handleFocus = () => {
  if (isEmpty.value && editorRef.value) {
    // 如果是空的，确保编辑器没有占位符内容
    if (editorRef.value.innerHTML === `<span class="placeholder">${props.placeholder}</span>`) {
      editorRef.value.innerHTML = '';
    }
    isEmpty.value = false;
  }
};

const handleBlur = () => {
  if (editorRef.value && (!editorRef.value.textContent || editorRef.value.textContent.trim() === '')) {
    isEmpty.value = true;
    editorRef.value.innerHTML = `<span class="placeholder">${props.placeholder}</span>`;
  }
  updateModelValue();
};

// 处理输入事件
const handleInput = () => {
  if (isComposing.value) return;
  
  if (editorRef.value) {
    isEmpty.value = !editorRef.value.textContent || editorRef.value.textContent.trim() === '';
  }
  
  // 保存当前光标位置
  saveCursorPosition();
  
  // 检查是否应该显示提示
  checkForTips();
  
  // 更新模型值
  updateModelValue();
};

// 保存当前光标位置
const saveCursorPosition = () => {
  const selection = window.getSelection();
  if (!selection || selection.rangeCount === 0) return;
  
  lastCursorPosition.value = selection.getRangeAt(0).cloneRange();
};

// 获取光标周围文本，用于检测 $ 触发
const getTextAroundCursor = (): { before: string, after: string } => {
  const selection = window.getSelection();
  if (!selection || selection.rangeCount === 0) return { before: '', after: '' };
  
  const range = selection.getRangeAt(0).cloneRange();
  
  // 获取当前节点
  const startNode = range.startContainer;
  
  // 如果是文本节点
  if (startNode.nodeType === Node.TEXT_NODE) {
    // 获取光标前的文本
    const beforeText = startNode.textContent?.substring(0, range.startOffset) || '';
    
    // 获取光标后的文本
    const afterText = startNode.textContent?.substring(range.startOffset) || '';
    
    return { before: beforeText, after: afterText };
  }
  
  return { before: '', after: '' };
};

// 处理组合输入（如中文输入法）
const handleCompositionStart = () => {
  isComposing.value = true;
};

const handleCompositionEnd = () => {
  isComposing.value = false;
  handleInput();
};

// 处理键盘事件
const handleKeyDown = (e: KeyboardEvent) => {
  // 处理删除编辑器标签
  if ((e.key === 'Backspace' || e.key === 'Delete') && window.getSelection) {
    const selection = window.getSelection();
    if (!selection || selection.rangeCount === 0) return;
    
    const range = selection.getRangeAt(0);
    
    // 检查选中的内容是否包含标签
    const selectedNodes = getSelectedNodes(range);
    const hasTag = selectedNodes.some(node => 
      node.nodeType === Node.ELEMENT_NODE && 
      (node as HTMLElement).classList?.contains('editor-tag')
    );
    
    if (hasTag) {
      e.preventDefault();
      
      // 删除选中的标签
      selectedNodes.forEach(node => {
        if (node.nodeType === Node.ELEMENT_NODE && 
            (node as HTMLElement).classList?.contains('editor-tag')) {
          const tagId = parseInt((node as HTMLElement).getAttribute('data-id') || '0');
          removeTag(tagId);
          node.parentNode?.removeChild(node);
        }
      });
      
      updateModelValue();
    }
  }
  
  // 在提示显示时，防止特定键默认行为
  if (showTips.value && ['ArrowUp', 'ArrowDown', 'Enter', 'Tab'].includes(e.key)) {
    e.preventDefault();
  }
};

// 获取选中的所有节点
const getSelectedNodes = (range: Range): Node[] => {
  const selectedNodes: Node[] = [];
  
  if (range.collapsed) {
    return selectedNodes;
  }
  
  const startNode = range.startContainer.nodeType === Node.TEXT_NODE ? 
    range.startContainer.parentNode : range.startContainer;
    
  const endNode = range.endContainer.nodeType === Node.TEXT_NODE ? 
    range.endContainer.parentNode : range.endContainer;
    
  if (!startNode || !endNode) return selectedNodes;
  
  if (startNode === endNode) {
    if (startNode.nodeType === Node.ELEMENT_NODE) {
      const element = startNode as HTMLElement;
      if (element.classList?.contains('editor-tag')) {
        selectedNodes.push(startNode);
      }
    }
    return selectedNodes;
  }
  
  // 收集选中范围内的所有节点
  const nodeIterator = document.createNodeIterator(
    editorRef.value!,
    NodeFilter.SHOW_ELEMENT,
    {
      acceptNode: (node: Node) => {
        if ((node as HTMLElement).classList?.contains('editor-tag')) {
          return NodeFilter.FILTER_ACCEPT;
        }
        return NodeFilter.FILTER_SKIP;
      }
    }
  );
  
  let currentNode: Node | null;
  let inSelection = false;
  
  while ((currentNode = nodeIterator.nextNode())) {
    if (currentNode === startNode) {
      inSelection = true;
    }
    
    if (inSelection) {
      selectedNodes.push(currentNode);
    }
    
    if (currentNode === endNode) {
      break;
    }
  }
  
  return selectedNodes;
};

// 检查是否应该显示提示
const checkForTips = () => {
  // 获取光标前的文本
  const { before } = getTextAroundCursor();
  
  // 查找最后一个 $ 符号
  const lastIndex = before.lastIndexOf('$');
  
  if (lastIndex === -1) {
    showTips.value = false;
    return;
  }
  
  // 获取 $ 后面的查询文本
  const query = before.substring(lastIndex + 1);
  
  // 检查是否有点号以处理对象属性
  const parts = query.split('.');
  currentObjectPath.value = parts.slice(0, -1);
  currentQuery.value = parts[parts.length - 1];
  
  // 过滤提示
  updateFilteredTips();
  
  if (filteredTips.value.length > 0) {
    showTips.value = true;
    
    // 计算提示框位置
    setTimeout(() => {
      if (!editorRef.value) return;
      
      const selection = window.getSelection();
      if (!selection || selection.rangeCount === 0) return;
      
      const range = selection.getRangeAt(0);
      const rect = range.getBoundingClientRect();
      
      // 设置下拉框位置
      tipPosition.value = {
        x: rect.left,
        y: rect.bottom + window.scrollY
      };
    }, 0);
  } else {
    showTips.value = false;
  }
};

// 获取当前层级的提示
const updateFilteredTips = () => {
  // 获取当前应该显示的提示数据
  let currentLevelSuggestions: (string | Record<string, any>)[] = props.suggestions;
  
  // 逐级查找对象属性
  if (currentObjectPath.value.length > 0) {
    for (const path of currentObjectPath.value) {
      let found = false;
      for (const suggestion of currentLevelSuggestions) {
        if (typeof suggestion === 'object' && suggestion !== null) {
          if (path in suggestion) {
            currentLevelSuggestions = Array.isArray(suggestion[path]) 
              ? suggestion[path] 
              : Object.keys(suggestion[path]).map(key => ({[key]: suggestion[path][key]}));
            found = true;
            break;
          }
        }
      }
      if (!found) {
        currentLevelSuggestions = [];
        break;
      }
    }
  }
  
  // 将当前层级的提示转换为字符串数组
  filteredTips.value = currentLevelSuggestions
    .filter(item => {
      if (typeof item === 'string') {
        return item.toLowerCase().includes(currentQuery.value.toLowerCase());
      } else {
        // 对于对象，使用其键作为提示
        return Object.keys(item).some(key => 
          key.toLowerCase().includes(currentQuery.value.toLowerCase())
        );
      }
    })
    .map(item => {
      if (typeof item === 'string') {
        return item;
      } else {
        // 对于对象，返回其所有键
        return Object.keys(item).join(', ');
      }
    });
};

// 选择一个提示
const selectTip = (tip: string) => {
  if (!tip) {
    showTips.value = false;
    return;
  }
  
  if (!editorRef.value) return;
  
  const selection = window.getSelection();
  if (!selection || selection.rangeCount === 0) return;
  
  const range = selection.getRangeAt(0).cloneRange();
  
  // 插入前确保恢复光标位置
  if (lastCursorPosition.value) {
    selection.removeAllRanges();
    selection.addRange(lastCursorPosition.value);
  }
  
  // 获取当前光标所在的文本节点
  const { before, after } = getTextAroundCursor();
  
  // 查找最后一个 $ 符号
  const lastIndex = before.lastIndexOf('$');
  if (lastIndex === -1) return;
  
  // 删除从 $ 开始到光标位置的文本
  const beforeTag = before.substring(0, lastIndex);
  
  // 创建新标签
  const newTag: Tag = {
    id: nextTagId.value++,
    text: tip,
    path: [...currentObjectPath.value],
    type: currentObjectPath.value.join('.') || 'default',
    fullPath: [...currentObjectPath.value, tip].join('.')
  };
  
  tags.value.push(newTag);
  
  // 创建标签HTML
  const tagHtml = createTagHtml(newTag);
  
  // 替换掉从 $ 到光标的文本
  const currentNode = range.startContainer;
  if (currentNode.nodeType === Node.TEXT_NODE) {
    // 创建文本节点
    const beforeTextNode = document.createTextNode(beforeTag);
    const afterTextNode = document.createTextNode(after);
    
    // 创建标签节点
    const tagContainer = document.createElement('div');
    tagContainer.innerHTML = tagHtml;
    const tagNode = tagContainer.firstChild as Node;
    
    // 替换当前节点
    const parentNode = currentNode.parentNode;
    if (parentNode) {
      parentNode.insertBefore(beforeTextNode, currentNode);
      parentNode.insertBefore(tagNode, currentNode);
      parentNode.insertBefore(afterTextNode, currentNode);
      parentNode.removeChild(currentNode);
      
      // 将光标移动到标签后面
      const newRange = document.createRange();
      newRange.setStartAfter(tagNode);
      newRange.collapse(true);
      
      selection.removeAllRanges();
      selection.addRange(newRange);
      lastCursorPosition.value = newRange.cloneRange();
    }
  }
  
  // 隐藏提示
  showTips.value = false;
  
  // 更新模型值
  updateModelValue();
  
  // 聚焦编辑器
  editorRef.value.focus();
};

// 关闭提示
const closeTip = () => {
  showTips.value = false;
};

// 移除标签
const removeTag = (tagId: number) => {
  const index = tags.value.findIndex(tag => tag.id === tagId);
  if (index !== -1) {
    tags.value.splice(index, 1);
  }
};

// 处理标签删除事件
const handleRemoveTagClick = (e: MouseEvent) => {
  const target = e.target as HTMLElement;
  if (target.classList.contains('remove-tag')) {
    const tagId = parseInt(target.getAttribute('data-id') || '0');
    
    // 找到包含此标签的父元素
    const tagElement = target.closest('.editor-tag') as HTMLElement;
    if (tagElement && tagElement.parentNode) {
      // 从DOM中移除标签
      tagElement.parentNode.removeChild(tagElement);
      
      // 从数据中移除标签
      removeTag(tagId);
      
      // 更新模型值
      updateModelValue();
      
      // 阻止事件冒泡
      e.stopPropagation();
    }
  }
};

// 编辑器初始化和事件绑定
onMounted(() => {
  if (editorRef.value) {
    editorRef.value.addEventListener('click', handleRemoveTagClick);
    
    // 初始化内容
    initializeContent();
    
    // 如果初始为空，显示占位符
    if (isEmpty.value) {
      editorRef.value.innerHTML = `<span class="placeholder">${props.placeholder}</span>`;
    }
  }
});

// 清理事件监听
onBeforeUnmount(() => {
  if (editorRef.value) {
    editorRef.value.removeEventListener('click', handleRemoveTagClick);
  }
});
</script>

<template>
  <div class="content-editable-container">
    <div
      ref="editorRef"
      class="content-editable"
      :class="{ 'is-empty': isEmpty }"
      contenteditable="true"
      @focus="handleFocus"
      @blur="handleBlur"
      @input="handleInput"
      @keydown="handleKeyDown"
      @compositionstart="handleCompositionStart"
      @compositionend="handleCompositionEnd"
    ></div>
    
    <TipDropdown
      v-if="showTips"
      :tips="filteredTips"
      :position="tipPosition"
      :query="currentQuery"
      @select-tip="selectTip"
      @close-tip="closeTip"
    />
  </div>
</template>

<style scoped>
.content-editable-container {
  position: relative;
  width: 100%;
}

.content-editable {
  width: 100%;
  min-height: 40px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 8px 12px;
  line-height: 1.5;
  font-size: 14px;
  color: #606266;
  background-color: #fff;
  outline: none;
  transition: border-color 0.2s;
  word-break: break-word;
  white-space: pre-wrap;
  cursor: text;
  /* 与 TagDisplay 组件保持一致的样式 */
  display: inline-flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 4px;
}

.content-editable:focus {
  border-color: #409eff;
}

.placeholder {
  color: #c0c4cc;
  pointer-events: none;
}

.content-editable.is-empty:not(:focus) {
  color: #c0c4cc;
}

.editor-tag {
  display: inline-flex;
  align-items: center;
  background-color: #ecf5ff;
  color: #409eff;
  border-radius: 3px;
  padding: 0 6px;
  height: 24px;
  margin: 0 2px;
  font-size: 0.9em;
  border: 1px solid #d9ecff;
  line-height: 22px;
  white-space: nowrap;
  vertical-align: middle;
  user-select: none;
  /* 与 TagDisplay 组件保持一致的样式 */
  box-sizing: border-box;
}

.remove-tag {
  margin-left: 4px;
  color: #909399;
  font-weight: bold;
  font-size: 14px;
  cursor: pointer;
  padding: 0 2px;
}

.remove-tag:hover {
  color: #f56c6c;
}

/* 不同类型的标签样式 - 与 TagDisplay 组件保持一致 */
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

.tag-type-customer\.address {
  background-color: #f2f6fc;
  color: #909399;
  border-color: #ebeef5;
}
</style>
