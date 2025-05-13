# Go-RAG 前端应用

这是为Go-RAG服务器项目提供的前端界面，使用Vue框架实现。前端提供三个主要功能页面：文档索引、文档检索和智能问答。

## 功能介绍

1. **文档索引 (Indexer)**
   - 上传文件并向量化到ES
   - 支持拖拽上传文件

2. **文档检索 (Retriever)**
   - 根据用户提问检索相关文档
   - 可配置Top K和最低相关度参数

3. **智能问答 (Chat)**
   - 检索文档后回答用户问题
   - 支持会话ID管理，可查看参考文档

## 项目设置

### 安装依赖
```bash
npm install
```

### 启动开发服务器
```bash
npm run dev
```

### 构建生产版本
```bash
npm run build
```

## 使用说明

1. 确保后端服务已启动（默认地址：http://localhost:8000）
2. 启动前端开发服务器
3. 访问前端应用（默认地址：http://localhost:5173）

## 技术栈

- Vue 3
- Vue Router
- Axios
- Vite
