<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { get } from 'lodash-es'
import { CursorPop } from './CursorPop'

// console.log('get', get({ a: { b: { c: 1 } } }, 'a.b.c', 0)) // 1

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

function getDropdownContent(expression: string): Array<[string, string]> {
  console.log(expression, 'expression')
  const exp = expression.replace(/^\$/, '').replace(/\.$/, '')
  const expList = exp.split('.').filter(Boolean)
  const lastProperty = expList[expList.length - 1]
  const value = get(TipContent, exp)
  const containsInCurrentDropdown = (t: string) => {
    const hit = dropdownList.value.some(item => item[0].includes(t))
    return hit
  }
  if (expression === '$') {
    // 1. 返回顶层
    return Object.entries(TipContent).map(item => {
      const key = item[0]
      const value = item[1]
      if (typeof value === 'object') {
        return [key, 'object']
      } else {
        return [key, value]
      }
    })
  } else if (expression.endsWith('.')) {
    // 3. 以.结尾
    const value = get(TipContent, exp)
    if (typeof value === 'object') {
      return Object.entries(value)
    } else {
      return []
    }
  } else if (expList.length > 0 && containsInCurrentDropdown(lastProperty)) {
    // 2. 最后一个属性已经在下拉列表中
    return dropdownList.value
  } else if (typeof value === 'string') {
    // 4. 值是字符串
    return []
  } else if (typeof value === 'object') {
    // 3. 值是对象
    return Object.entries(value)
  } else {
    // 5. 值是已知的类型字符串
    // 6. 其他
    // throw new Error('Invalid expression')
    console.warn('Invalid expression', expression)
    return []
  }
}

function calcTagContent(text: string) {
  const lastSign = text.lastIndexOf('$')
  const tagContent = text.substring(lastSign)
  return tagContent
}

let popper!: CursorPop

onMounted(() => {
  const inputElm = inputarea.value!;
  const floatElm = floating.value!;
  popper = new CursorPop(inputElm, floatElm)
  inputElm.addEventListener('keyup', () => {
    console.log("光标位置:", popper.getCaretPosition());
  });
  inputElm.addEventListener('input', (e: Event) => {
    const inputEvt = e as InputEvent;
    const beforeText = popper.getTextBeforeCursor();
    const tagContent = calcTagContent(beforeText)
    console.log(beforeText, tagContent);
    const dropdownContent = getDropdownContent(tagContent)
    console.log('dropdownContent', dropdownContent);
    if (inputEvt.data === '$') {
      dropdownList.value = dropdownContent
      setTimeout(() => {
        popper.show()
      }, 200)
    } else if (inputEvt.data === '.') {
      dropdownList.value = dropdownContent
      popper.show()
    } else {
      dropdownList.value = dropdownContent
      if (dropdownContent.length === 0) {
        popper.hide()
      } else {
        popper.show()
      }
    }
    popper.updatePosition();
  });

  const tipElm = floating.value!;
  tipElm.addEventListener('click', () => {
    setTimeout(() => {
      popper.hide()
    }, 0)
  });
})
</script>

<template>
  <div class="input-width-tip-wrapper">
    <div contenteditable="true" class="input-with-tip" ref="inputarea"></div>
    <div ref="floating" class="floating-tip">
      <span @click="popper.insertTextAtCursor(item[0])" v-for="item in dropdownList">
        {{ item[0] }}: {{ item[1] }}
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
