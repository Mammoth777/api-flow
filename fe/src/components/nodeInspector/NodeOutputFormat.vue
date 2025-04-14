<script setup lang="ts">
interface OutputField {
  field: string;
  type: string;
  desc: string;
  default: any;
}

defineProps({
  outputs: {
    type: Array as () => OutputField[],
    required: true
  }
});
</script>

<template>
  <div class="output-format-container">
    <table class="output-table" v-if="outputs && outputs.length > 0">
      <thead>
        <tr>
          <th>字段名称</th>
          <th>类型</th>
          <th>描述</th>
          <th>默认值</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(output, index) in outputs" :key="index">
          <td class="field-name">{{ output.field }}</td>
          <td><span class="type-badge">{{ output.type }}</span></td>
          <td>{{ output.desc }}</td>
          <td>{{ output.default === null ? '-' : output.default }}</td>
        </tr>
      </tbody>
    </table>
    <div class="empty-message" v-else>暂无输出参数定义</div>
  </div>
</template>

<style scoped>
.output-format-container {
  margin-top: 8px;
}

.output-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
  text-align: left;
}

.output-table th {
  background-color: #f5f5f5;
  padding: 8px;
  border-bottom: 1px solid #e0e0e0;
  color: #333;
  font-weight: 500;
}

.output-table td {
  padding: 8px;
  border-bottom: 1px solid #e0e0e0;
  color: #555;
}

.field-name {
  font-weight: 500;
  color: #333;
}

.type-badge {
  background-color: #e8f0fe;
  color: #1a73e8;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
}

.empty-message {
  padding: 12px;
  color: #999;
  text-align: center;
  font-style: italic;
  background-color: #f9f9f9;
  border-radius: 4px;
}
</style>
