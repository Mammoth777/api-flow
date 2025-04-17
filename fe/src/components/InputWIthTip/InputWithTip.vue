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
    // console.log("光标位置:", popper.getCaretPosition());
  });
  inputElm.addEventListener('input', (e: Event) => {
    // 1. 生成dropdown
    const beforeText = popper.getTextBeforeCursor();
    const tagContent = calcTagContent(beforeText)
    const dropdownContent = getDropdownContent(tagContent)
    dropdownList.value = dropdownContent
    if (dropdownContent.length === 0) {
      popper.hide()
    } else {
      setTimeout(() => {
        popper.show()
      }, 100)
    }
    popper.updatePosition();

    // 2. 生成tag
    const inputEvt = e as InputEvent
    console.log(inputEvt.data, 'inputEvt')
  });

  inputElm.addEventListener('keyup', (e: KeyboardEvent) => {
    e.preventDefault()
    e.stopPropagation()
    const keyCode = e.keyCode
    if (keyCode === 13) {
      // Enter
      console.log('Enter pressed');
    } else if (keyCode === 8) {
      // Backspace
      console.log('Backspace pressed');
    }
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

.cursor {
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
