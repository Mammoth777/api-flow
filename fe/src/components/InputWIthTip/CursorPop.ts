import { TagItem } from "./TagItem";

export class CursorPop {
  inputDiv: HTMLDivElement;
  floatTipDiv: HTMLDivElement;
  isShown: boolean = false;
  private selection: Selection | null = null;
  private activeItemIndex: number = -1; // 改为公共属性，以便外部访问

  getActiveItemIndex() {
    return this.activeItemIndex;
  }

  setActiveItemIndex(index: number) {
    this.activeItemIndex = index;
  }

  constructor(editableDiv: HTMLDivElement, floatTipDiv: HTMLDivElement) {
    this.inputDiv = editableDiv;
    this.floatTipDiv = floatTipDiv;
    this.hide();
    this.setupStyles();
  }

  // 设置基本样式
  setupStyles() {
    // 浮动提示框基础样式
    Object.assign(this.floatTipDiv.style, {
      maxHeight: '200px',
      overflowY: 'auto',
      zIndex: '1000'
    });
  }

  show() {
    this.floatTipDiv.style.display = "block";
    this.isShown = true;
    this.activeItemIndex = -1; // 重置活动项索引
  }

  hide() {
    this.floatTipDiv.style.display = "none";
    this.isShown = false;
  }

  // 高亮选中项
  highlightItem(index: number) {
    if (!this.isShown) return;

    // 移除所有现有高亮
    const items = this.floatTipDiv.querySelectorAll('.dropdown-item');
    items.forEach((item, idx) => {
      item.classList.remove('active');
      if (idx === this.activeItemIndex) {
        item.setAttribute('aria-selected', 'false');
      }
    });

    // 设置新的高亮项
    if (index >= 0 && index < items.length) {
      this.activeItemIndex = index;
      const activeItem = items[index] as HTMLElement;
      activeItem.classList.add('active');
      activeItem.setAttribute('aria-selected', 'true');

      // 确保选中项在视口内
      this.scrollItemIntoView(activeItem);
    }
  }

  // 确保选中项在滚动视图中可见
  scrollItemIntoView(element: HTMLElement) {
    const container = this.floatTipDiv;
    const containerHeight = container.clientHeight;
    const itemTop = element.offsetTop;
    const itemHeight = element.offsetHeight;

    // 如果选中项在可视区域上方，滚动到选中项
    if (itemTop < container.scrollTop) {
      container.scrollTop = itemTop;
    }
    // 如果选中项在可视区域下方，滚动以显示选中项
    else if (itemTop + itemHeight > container.scrollTop + containerHeight) {
      container.scrollTop = itemTop + itemHeight - containerHeight;
    }
  }

  // 获取下一个项目
  selectNextItem() {
    if (!this.isShown) return;

    const items = this.floatTipDiv.querySelectorAll('.dropdown-item');
    if (items.length === 0) return;

    let nextIndex = this.activeItemIndex + 1;
    if (nextIndex >= items.length) nextIndex = 0;

    this.highlightItem(nextIndex);
    return nextIndex;
  }

  // 获取上一个项目
  selectPrevItem() {
    if (!this.isShown) return;

    const items = this.floatTipDiv.querySelectorAll('.dropdown-item');
    if (items.length === 0) return;

    let prevIndex = this.activeItemIndex - 1;
    if (prevIndex < 0) prevIndex = items.length - 1;

    this.highlightItem(prevIndex);
    return prevIndex;
  }

  confirmSelectionWithValue(value: string): boolean {
    
    // 确保输入区域获取焦点
    this.inputDiv.focus();
    
    // 等待DOM更新完成
    setTimeout(() => {
      this.insertTextAtCursor(value);
    }, 0);
    
    return true;

  }

  // 确认选择当前高亮项
  confirmSelection(): boolean {
    if (!this.isShown || this.activeItemIndex < 0) return false;

    const items = this.floatTipDiv.querySelectorAll('.dropdown-item');
    if (this.activeItemIndex < items.length) {
      const activeItem = items[this.activeItemIndex] as HTMLElement;
      
      // 确保输入区域获取焦点
      this.inputDiv.focus();
      
      // 如果选择状态丢失，尝试恢复最后一个已知的选择
      if (!window.getSelection()?.rangeCount) {
        console.warn("Selection lost, attempting to restore");
        // 如果没有有效的选择，将光标放在内容末尾
        const selection = window.getSelection();
        if (selection) {
          const range = document.createRange();
          if (this.inputDiv.lastChild) {
            range.setStartAfter(this.inputDiv.lastChild);
            range.collapse(true);
          } else {
            range.setStart(this.inputDiv, 0);
          }
          selection.removeAllRanges();
          selection.addRange(range);
          this.selection = selection;
        }
      } else {
        // 保存当前的选择状态
        this.selection = window.getSelection();
      }
      
      // 等待DOM更新完成
      setTimeout(() => {
        const key = activeItem.dataset.key as string;
        this.insertTextAtCursor(key);
      }, 0);
      
      return true;
    }

    return false;
  }

  // 获取当前选中的项目文本
  getSelectedItemText(): string | null {
    if (!this.isShown || this.activeItemIndex < 0) return null;

    const items = this.floatTipDiv.querySelectorAll('.dropdown-item');
    if (this.activeItemIndex < items.length) {
      const itemKey = items[this.activeItemIndex].querySelector('.item-key');
      if (itemKey) {
        // 移除冒号，只返回键值
        return itemKey.textContent?.replace(':', '') || null;
      }
    }

    return null;
  }

  // 获取所有可选项的数量
  getItemCount(): number {
    return this.floatTipDiv.querySelectorAll('.dropdown-item').length;
  }

  getSelection() {
    if (this.selection) {
      return this.selection;
    }
    const selection = window.getSelection();
    if (!selection) {
      throw new Error("Selection is null");
    }
    this.selection = selection;
    return this.selection;
  }

  getRange() {
    if (this.selection) {
      return this.selection.getRangeAt(0);
    } else {
      throw new Error("Selection is null");
    }
  }

  private getCaretPositionRelativeToElement(): { x: number; y: number } {
    const selection = this.getSelection();
    if (!selection) {
      console.error("No selection found");
      return { x: 0, y: 0 };
    }
    if (!selection.rangeCount) {
      console.error("No range found");
      return { x: 0, y: 0 };
    }

    const range = selection.getRangeAt(0).cloneRange();
    const caretRect = range.getBoundingClientRect(); // 光标矩形位置
    const containerRect = this.inputDiv.getBoundingClientRect(); // 容器矩形位置

    // 计算光标相对容器的偏移
    const offsetX = caretRect.left - containerRect.left;
    const offsetY = caretRect.top - containerRect.top;

    return { x: offsetX, y: offsetY };
  }

  updatePosition() {
    const { x, y } = this.getCaretPositionRelativeToElement();
    this.floatTipDiv.style.left = `${x + 3}px`;
    this.floatTipDiv.style.top = `${y + 20}px`; // 略微下移，避免遮挡光标

    // 确保下拉框不超出视口边界
    const rect = this.floatTipDiv.getBoundingClientRect();
    const viewportHeight = window.innerHeight;
    const viewportWidth = window.innerWidth;

    // 如果下拉框底部超出视口
    if (rect.bottom > viewportHeight) {
      this.floatTipDiv.style.top = `${y - rect.height - 5}px`; // 向上显示
    }

    // 如果下拉框右侧超出视口
    if (rect.right > viewportWidth) {
      this.floatTipDiv.style.left = `${viewportWidth - rect.width - 10}px`;
    }
  }

  // 修改插入文本方法，确保点号后的内容不成为独立标签
  insertTextAtCursor(text: string) {
    const selection = this.getSelection();
    if (!selection || !selection.rangeCount) {
      console.error("Cannot get selection for text insertion");
      return;
    }
    
    // 使用当前实时的选择状态，而不是缓存的
    const range = selection.getRangeAt(0);
    if (!range) {
      console.error("No range found");
      return;
    }

    // 获取当前光标前的文本
    const beforeText = this.getTextBeforeCursor();

    // 查找最后一个 $ 符号的位置
    const lastDollarIndex = beforeText.lastIndexOf('$');
    if (lastDollarIndex !== -1) {
      console.log('1. 找到 $ 符号，进行补全');
      // 获取从 $ 到光标位置的文本
      const partialInput = beforeText.substring(lastDollarIndex);

      // 检查是否是点号后面的情况（如 $fullname.last）
      const lastDotIndex = partialInput.lastIndexOf('.');

      if (lastDotIndex !== -1) {
        console.log('1.1 这是点号后面的补全情况');
        // 这是点号后面的补全情况
        const textAfterDot = partialInput.substring(lastDotIndex + 1); // 点号后的文本

        // 如果选项以用户输入的部分开头，需要删除用户输入的部分
        if (text.startsWith(textAfterDot)) {
          console.log('1.1.1 这是点号后面的补全情况，删除已输入的部分');
          // 创建一个新范围来选择从点号后到当前光标的文本
          const deleteRange = document.createRange();
          let currentNode = range.startContainer;
          let currentOffset = range.startOffset;

          // 检查当前节点是否是文本节点
          if (currentNode.nodeType === Node.TEXT_NODE) {
            // 计算要删除的字符数
            const charsToDelete = textAfterDot.length;

            // 设置范围起始点为当前光标位置减去要删除的字符数
            deleteRange.setStart(currentNode, currentOffset - charsToDelete);
            deleteRange.setEnd(currentNode, currentOffset);

            // 删除已输入的部分
            deleteRange.deleteContents();
          }

          // 直接插入文本，而不是标签
          const textNode = document.createTextNode(text);
          range.insertNode(textNode);

          // 将光标设置在文本节点之后
          range.setStartAfter(textNode);
          range.setEndAfter(textNode);
          selection.removeAllRanges();
          selection.addRange(range);
        }
      } else {
        console.log('1.2 这是 $ 符号后的补全情况');

        // 创建一个新范围来选择从 $ 符号（包含）到当前光标的文本
        const deleteRange = document.createRange();
        let currentNode = range.startContainer;
        let currentOffset = range.startOffset;

        if (currentNode.nodeType === Node.TEXT_NODE) {
          // 计算要删除的字符数（包括$符号）
          const charsToDelete = partialInput.length;

          // 设置范围起始点为当前光标位置减去要删除的字符数
          deleteRange.setStart(currentNode, currentOffset - charsToDelete);
          deleteRange.setEnd(currentNode, currentOffset);

          // 删除包括 $ 在内的所有已输入部分
          deleteRange.deleteContents();

          const tagItem = new TagItem('$' + text);
          tagItem.insertByRange(range);
          // 将光标设置在零宽度空格之后
          selection.removeAllRanges();
          selection.addRange(range);
        }
      }
    } else {
      console.log('2. 没有找到 $ 符号，直接插入文本');
      const tagItem = new TagItem(text);
      tagItem.insertByRange(range);

      selection.removeAllRanges();
      selection.addRange(range);
    }
    
    // 更新内部保存的选择
    this.selection = selection;
  }

  getCaretPosition(): number {
    let caretOffset = 0;
    const selection = window.getSelection();
    if (!selection) {
      console.error("Cannot get selection");
      return 0;
    }
    if (selection.rangeCount > 0) {
      const range = selection.getRangeAt(0);
      const preCaretRange = range.cloneRange();
      preCaretRange.selectNodeContents(this.inputDiv);
      preCaretRange.setEnd(range.startContainer, range.startOffset);
      caretOffset = preCaretRange.toString().length;
    }
    return caretOffset;
  }

  getTextBeforeCursor(): string {
    // 使用当前实时的选择状态，而不是可能过期的缓存
    const selection = window.getSelection();
    if (!selection) {
      console.error("No selection found");
      return "";
    }
    if (!selection.rangeCount) {
      console.error("No range found");
      return "";
    }

    const range = selection.getRangeAt(0).cloneRange(); // 克隆当前 Range
    range.selectNodeContents(this.inputDiv); // 将范围设置为整个内容区域
    
    // 添加安全检查，确保 anchorNode 不为 null
    if (!selection.anchorNode) {
      console.error("Selection anchorNode is null");
      return "";
    }
    
    range.setEnd(selection.anchorNode, selection.anchorOffset); // 将结束点设置为光标位置

    return range.toString(); // 获取光标之前的文本
  }

  // 光标向右移动
  cursorRight() {
    const selection = this.selection;
    if (!selection || selection.rangeCount === 0) return;

    const range = selection.getRangeAt(0);
    const node = range.startContainer;

    // 检查当前光标是否在标签内
    let isAtTagEnd = false;
    let tagElement = null;

    // 情况1: 光标在标签内的文本节点末尾
    if (node.nodeType === Node.TEXT_NODE && node.parentElement?.classList.contains('tag-item')) {
      tagElement = node.parentElement;
      const textContent = node.textContent || '';
      isAtTagEnd = range.startOffset >= textContent.length;
    }
    // 情况2: 光标直接在标签元素内(没有文本子节点的情况)
    else if (node.nodeType === Node.ELEMENT_NODE && (node as Element).classList.contains('tag-item')) {
      tagElement = node as HTMLElement;
      isAtTagEnd = true;
    }

    // 如果光标在标签末尾，移动到标签后面
    if (isAtTagEnd && tagElement) {
      // 创建新范围，位于标签之后
      const newRange = document.createRange();
      newRange.setStartAfter(tagElement);
      newRange.collapse(true); // 折叠到起点

      // 应用新范围
      selection.removeAllRanges();
      selection.addRange(newRange);
      return true; // 表示光标已移动
    }

    return false; // 表示没有特殊处理，让浏览器默认行为生效
  }
}
