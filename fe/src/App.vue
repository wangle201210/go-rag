<script setup>
// 导入全局样式
import './assets/styles/main.css';
</script>

<template>
  <div class="app-container">
    <header>
      <div class="header-content">
        <h1>GO-RAG</h1>
        <nav>
          <router-link to="/indexer">文档索引</router-link>
          <router-link to="/retriever">文档检索</router-link>
          <router-link to="/chat">智能问答</router-link>
        </nav>
      </div>
    </header>
    <main>
      <div class="page-container">
        <router-view />
      </div>
    </main>
    <footer>
      <div class="footer-content">
        <p>Go-RAG &copy; 2025</p>
      </div>
    </footer>
  </div>
</template>

<style>
.app-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

header {
  background-color: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  padding: 1rem 0;
}

.header-content {
  max-width: var(--container-width);
  margin: 0 auto;
  padding: 0 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

header h1 {
  font-size: 1.5rem;
  color: var(--primary-color);
  margin: 0;
}

nav {
  display: flex;
  gap: 1.5rem;
}

nav a {
  text-decoration: none;
  color: var(--text-color);
  padding: 0.5rem 0.75rem;
  border-radius: var(--border-radius);
  transition: all 0.3s ease;
  font-weight: 500;
}

nav a:hover {
  color: var(--primary-color);
}

nav a.router-link-active {
  background-color: var(--primary-color);
  color: #888888;
}

main {
  flex: 1;
  padding: 2rem 0;
}

footer {
  background-color: white;
  border-top: 1px solid var(--border-color);
  padding: 1.5rem 0;
}

.footer-content {
  max-width: var(--container-width);
  margin: 0 auto;
  padding: 0 2rem;
  text-align: center;
  color: var(--secondary-color);
  font-size: 0.9rem;
}
</style>
