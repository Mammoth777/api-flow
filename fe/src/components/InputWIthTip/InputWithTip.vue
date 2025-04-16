<script setup lang="ts">
import { inline, useFloating } from '@floating-ui/vue';
import { ref } from 'vue';

const reference = ref(null);
const floating = ref(null);
const { floatingStyles } = useFloating(reference, floating, {
  middleware: [inline()]
});
</script>

<template>
  <div class="input-width-tip-wrapper">
    <div contenteditable="true" class="input-with-tip"></div>
    <div ref="reference" class="cursor">c</div>
    <div ref="floating" :style="floatingStyles" class="floating-tip">Tooltip</div>
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

</style>
