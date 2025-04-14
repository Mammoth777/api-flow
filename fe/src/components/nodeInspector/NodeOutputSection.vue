<template>
  <div class="form-group output-format-section">
    <label>输出格式</label>
    <NodeOutputFormat :outputs="outputs" />
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue';
import NodeOutputFormat from './NodeOutputFormat.vue';

defineProps({
  outputs: {
    type: Array as () => Array<{ field: string; type: string; desc: string; default: any }>,
    required: true
  }
});
</script>

<style scoped>
.output-format-section {
  margin-top: 16px;
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 12px;
  color: #666;
}
</style>
