<template>
  <div class="app-container">
    <header class="app-header">
      <h1>Geminate</h1>
      <nav class="nav-links">
        <router-link to="/workflows">工作流列表</router-link>
        <router-link to="/workflow/create">创建工作流</router-link>
      </nav>
    </header>
    <main class="app-main">
      <router-view />
    </main>
  </div>
</template>

<style>
body {
  margin: 0;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  overflow: hidden; /* 防止 body 出现滚动条 */
}

.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden; /* 防止容器出现滚动条 */
}

.app-header {
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
  padding: 0 20px;
  height: 42px; /* 减少高度 */
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.app-header h1 {
  font-size: 16px; /* 减小标题字号 */
  color: #333;
  margin: 0;
  font-weight: 500; /* 调整字重 */
}

.nav-links {
  display: flex;
  gap: 16px;
}

.nav-links a {
  color: #1890ff;
  text-decoration: none;
  font-size: 14px;
}

.nav-links a:hover {
  color: #40a9ff;
}

.nav-links a.router-link-active {
  color: #096dd9;
  font-weight: 500;
}

.app-main {
  flex: 1;
  padding: 0; /* 去掉内边距以最大化工作区域 */
  overflow: hidden;
  background-color: #f0f2f5;
}
</style>
