<template>
  <div class="highlight-code-container" :class="{ 'theme-dark': isDarkTheme }">
    <div class="code-header" v-if="showHeader">
      <div class="language-tag" v-if="language">{{ language }}</div>
      <div class="actions">
        <button class="copy-btn" @click="copyCode" :title="copyBtnText">
          <svg v-if="!copied" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
            <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
          </svg>
          <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 6L9 17l-5-5"></path>
          </svg>
          {{ copyBtnText }}
        </button>
      </div>
    </div>
    <pre ref="codeBlock" class="code-block"><code :class="codeClass" v-html="highlightedCode"></code></pre>
    <div class="line-numbers" v-if="lineNumbers">
      <div v-for="i in lineCount" :key="i" class="line-number">{{ i }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import hljs from 'highlight.js';
import 'highlight.js/styles/github.css'; // 默认样式，可以根据需要导入不同主题
hljs.configure({ ignoreUnescapedHTML: true });

interface Props {
  code: string;
  language?: string;
  lineNumbers?: boolean;
  autoDetect?: boolean;
  showHeader?: boolean;
  isDarkTheme?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  code: '',
  language: '',
  lineNumbers: false,
  autoDetect: true,
  showHeader: true,
  isDarkTheme: false
});

const codeBlock = ref<HTMLElement | null>(null);
const copied = ref(false);
const copyBtnText = computed(() => copied.value ? '已复制' : '复制');
const lineCount = computed(() => props.code.split('\n').length);

const highlightedCode = computed(() => {
  if (!props.code) return '';
  
  try {
    if (props.language && hljs.getLanguage(props.language)) {
      return hljs.highlight(props.code, { language: props.language }).value;
    } else if (props.autoDetect) {
      return hljs.highlightAuto(props.code).value;
    } else {
      return escapeHtml(props.code);
    }
  } catch (e) {
    console.error('Highlighting error:', e);
    return escapeHtml(props.code);
  }
});

const codeClass = computed(() => {
  return props.language ? `language-${props.language}` : '';
});

const copyCode = async () => {
  try {
    await navigator.clipboard.writeText(props.code);
    copied.value = true;
    setTimeout(() => {
      copied.value = false;
    }, 2000);
  } catch (err) {
    console.error('Failed to copy text: ', err);
  }
};

const escapeHtml = (unsafe: string): string => {
  return unsafe
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
};

watch(() => props.code, () => {
  if (codeBlock.value) {
    hljs.highlightElement(codeBlock.value);
  }
});

onMounted(() => {
  if (codeBlock.value) {
    hljs.highlightElement(codeBlock.value);
  }
});
</script>

<style scoped>
.highlight-code-container {
  max-height: calc(100vh - 250px);
  position: relative;
  margin: 1rem 0;
  border-radius: 6px;
  overflow-y: auto;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  background-color: #f6f8fa;
  text-align: left;
}

.theme-dark {
  background-color: #282c34;
  color: #abb2bf;
}

.code-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 1rem;
  background-color: #e8eaed;
  border-bottom: 1px solid #d0d7de;
}

.theme-dark .code-header {
  background-color: #21252b;
  border-bottom: 1px solid #3e4451;
}

.language-tag {
  font-size: 0.8rem;
  font-weight: 500;
  color: #57606a;
}

.theme-dark .language-tag {
  color: #a0a8b7;
}

.actions {
  display: flex;
  gap: 0.5rem;
}

.copy-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  border: none;
  background: none;
  color: #57606a;
  font-size: 0.8rem;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: background-color 0.2s, color 0.2s;
}

.copy-btn:hover {
  background-color: #dde1e6;
  color: #24292e;
}

.theme-dark .copy-btn {
  color: #a0a8b7;
}

.theme-dark .copy-btn:hover {
  background-color: #3e4451;
  color: #efefef;
}

.code-block {
  margin: 0;
  padding: 1rem;
  overflow-x: auto;
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
  font-size: 0.85rem;
  line-height: 1.6;
  border-radius: 0;
}

.line-numbers {
  position: absolute;
  top: calc(2.5rem + 1px); /* 标题栏高度 + 边框 */
  left: 0;
  padding: 1rem 0.5rem 1rem 0.75rem;
  color: #6e7681;
  text-align: right;
  background-color: #f6f8fa;
  user-select: none;
  counter-reset: line;
}

.theme-dark .line-numbers {
  background-color: #282c34;
  color: #636d83;
}

.line-number {
  display: block;
  line-height: 1.6;
  font-size: 0.85rem;
}

:deep(code) {
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
}

/* 启用行号时的代码块左边距 */
.highlight-code-container:has(.line-numbers) .code-block {
  padding-left: 3.5rem;
}
</style>
