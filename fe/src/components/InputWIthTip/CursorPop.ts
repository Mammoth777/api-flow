export class CursorPop {
  inputDiv: HTMLDivElement;
  floatTipDiv: HTMLDivElement;
  isShown: boolean = false;
  private selection: Selection | null = null;
  constructor(editableDiv: HTMLDivElement, floatTipDiv: HTMLDivElement) {
    this.inputDiv = editableDiv;
    this.floatTipDiv = floatTipDiv;
    this.hide()
  }

  show() {
    this.floatTipDiv.style.display = "block";
    this.isShown = true;
  }

  hide() {
    this.floatTipDiv.style.display = "none";
    this.isShown = false
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
    return this.selection
  }

  getRange() {
    if (this.selection) {
      return this.selection.getRangeAt(0);
    } else {
      throw new Error("Selection is null");
    }
  }

  private getCaretPositionRelativeToElement(): { x: number; y: number } {
    const selection = this.getSelection()
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
    const { x, y } = this.getCaretPositionRelativeToElement()
    this.floatTipDiv.style.left = `${x + 3}px`;
    this.floatTipDiv.style.top = `${y}px`;
  }

  insertTextAtCursor(text: string) {
    const selection = this.getSelection();
    const range = this.getRange()
    if (!range) {
        console.error("No range found");
        return;
    }
    range.deleteContents(); // 可选，用于替换光标位置的内容

    // 创建一个文本节点
    const textNode = document.createTextNode(text);

    // 在光标位置插入内容
    range.insertNode(textNode);

    // 调整光标位置到文本节点之后
    range.setStartAfter(textNode);
    range.setEndAfter(textNode);
    selection.removeAllRanges();
    selection.addRange(range);
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
    const selection = this.getSelection()
    if (!selection) {
      console.error("No selection found");
      return ""
    };
    if (!selection.rangeCount) {
      console.error("No range found");
      return ""
    }

    const range = selection.getRangeAt(0).cloneRange(); // 克隆当前 Range
    range.selectNodeContents(this.inputDiv); // 将范围设置为整个内容区域
    range.setEnd(selection.anchorNode!, selection.anchorOffset); // 将结束点设置为光标位置

    return range.toString(); // 获取光标之前的文本
  }
}
