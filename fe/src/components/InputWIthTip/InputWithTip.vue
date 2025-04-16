<script setup lang="ts">
import { inline, useFloating } from '@floating-ui/vue';
import { ref, onMounted } from 'vue';
import { get } from 'lodash-es'

// console.log('get', get({ a: { b: { c: 1 } } }, 'a.b.c', 0)) // 1

const reference = ref(null);
const floating = ref<HTMLDivElement | null>(null);
const inputarea = ref<HTMLDivElement | null>(null);
const { floatingStyles } = useFloating(reference, floating, {
  middleware: [inline()]
});

const dropdownList = ref([])

const floatTipStyle = ref({
  // position: 'absolute',
  top: '0px',
  left: '0px',
  display: 'none'
  // background: '#222',
  // color: 'white',
  // fontWeight: 'bold',
  // padding: '5px',
  // borderRadius: '4px',
  // fontSize: '90%',
})

let inputareaSelection: Selection | null = null;
let inputareaRange: Range | null = null;


function getCaretPosition(editableDiv: HTMLDivElement): number {
    let caretOffset = 0;
    const selection = window.getSelection();
    if (!selection) {
      console.error("Cannot get selection");
      return 0;
    }
    if (selection.rangeCount > 0) {
        const range = selection.getRangeAt(0);
        const preCaretRange = range.cloneRange();
        preCaretRange.selectNodeContents(editableDiv);
        preCaretRange.setEnd(range.startContainer, range.startOffset);
        caretOffset = preCaretRange.toString().length;

        inputareaSelection = selection;
        inputareaRange = range;
    }
    return caretOffset;
}

function getCaretPositionRelativeToElement(container: HTMLElement): { x: number; y: number } {
    const selection = inputareaSelection
    if (!selection?.rangeCount) {
        console.error("No selection found");
        return { x: 0, y: 0 };
    }

    const range = selection.getRangeAt(0).cloneRange();
    const caretRect = range.getBoundingClientRect(); // 光标矩形位置
    const containerRect = container.getBoundingClientRect(); // 容器矩形位置

    // 计算光标相对容器的偏移
    const offsetX = caretRect.left - containerRect.left;
    const offsetY = caretRect.top - containerRect.top;

    return { x: offsetX, y: offsetY };
}

function insertTextAtCursor(text: string) {
    const selection = inputareaSelection;
    const range = inputareaRange
    if (!selection) {
        console.error("No selection found");
        return;
    }
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

function getDropdownContent(expression: string) {
  console.log(expression, 'expression')
  const exp = expression.replace(/^\$/, '')
  const value = get(TipContent, exp, null)
  console.log(value, 'value')
  
  console.log(value, 'dropdown')
  return value
}

function getTextBeforeCursor(editableDiv: HTMLDivElement): string {
    const selection = window.getSelection();
    if (!selection) {
      console.error("No selection found");
      return ""
    };
    if (!selection.rangeCount) {
      console.error("No range found");
      return ""
    }

    const range = selection.getRangeAt(0).cloneRange(); // 克隆当前 Range
    range.selectNodeContents(editableDiv); // 将范围设置为整个内容区域
    range.setEnd(selection.anchorNode!, selection.anchorOffset); // 将结束点设置为光标位置

    return range.toString(); // 获取光标之前的文本
}

function calcTagContent(text: string) {
  const lastSign = text.lastIndexOf('$')
  const tagContent = text.substring(lastSign)
  return tagContent
}

onMounted(() => {
  const updatePosition = () => {
    const { x, y } = getCaretPositionRelativeToElement(inputarea.value!)
    floatTipStyle.value!.left = `${x + 3}px`;
    floatTipStyle.value!.top = `${y}px`;
    console.log("光标相对位置:", x, y);
  }

  const inputElm = inputarea.value!;
  inputElm.addEventListener('keyup', () => {
    console.log("光标位置:", getCaretPosition(inputElm));
  });
  inputElm.addEventListener('input', (e: Event) => {
    const inputEvt = e as InputEvent;
    if (inputEvt.data === '$') {
      const beforeText = getTextBeforeCursor(inputElm);
      const tagContent = calcTagContent(beforeText)
      console.log(beforeText, tagContent);
      const dropdownContent = getDropdownContent(tagContent)
      console.log('dropdownContent', dropdownContent);
      updatePosition();
      floatTipStyle.value!.display = 'block'
    } else if (inputEvt.data === '.') {
      const beforeText = getTextBeforeCursor(inputElm);
      const tagContent = calcTagContent(beforeText)
      console.log(beforeText, tagContent);
      const dropdownContent = getDropdownContent(tagContent)
      console.log('dropdownContent', dropdownContent);
      updatePosition();
      floatTipStyle.value!.display = 'block'
    } else {
      floatTipStyle.value!.display = 'none'
    }
  });

  const tipElm = floating.value!;
  tipElm.addEventListener('click', () => {
    setTimeout(() => {
      floatTipStyle.value!.display = 'none'
    }, 0)
  });
})
</script>

<template>
  <div class="input-width-tip-wrapper">
    <div contenteditable="true" class="input-with-tip" ref="inputarea"></div>
    <div ref="reference" class="cursor">c</div>
    <div ref="floating" :style="floatTipStyle" class="floating-tip">
      <span @click="insertTextAtCursor('hello')">hello</span>
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
.cursor{
  display: inline;
  posation: absolute;
}
.floating-tip {
  /* width: max-content; */
  position: absolute;
  top: 0;
  left: 0;
  background: #222;
  color: white;
  font-weight: bold;
  padding: 5px;
  border-radius: 4px;
  font-size: 90%;
}

.floating-tip span {
  display: block;
  cursor: pointer;
  padding: 2px 5px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

</style>
