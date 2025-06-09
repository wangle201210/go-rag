# go-rag
基于eino+gf+vue实现知识库的rag
1. 创建知识库
   ![](./server/static/kb.png)
2. 选择需要使用的知识库，上传文档
   ![](./server/static/kb-select.png)
   ![](./server/static/indexer.png)
3. 检索
   ![](./server/static/retriever.png)
4. 对话
   ![](./server/static/chat.png)
5. mcp （以集成到deepchat为例）
   ![](./server/static/mcp-cfg.png)
   ![](./server/static/mcp-use.png)


## 存储层
- [x] es8存储向量相关数据

## 功能列表
- [x] md、pdf、html 文档解析
- [x] 网页解析
- [x] 文档检索
- [x] 长文档自动切割(chunk)
- [x] 提供http接口 [rag-api](./server/README.md)
- [x] 提供 index、retrieve、chat 的前端界面
- [x] 多知识库支持


## 未来计划
- [ ] 使用mysql存储chunk和文档的映射关系，目前放在es的ext字段

## 使用
### clone项目
```bash
git clone https://github.com/wangle201210/go-rag.git
```

### 使用 Docker Compose 快速启动（推荐）
```bash
cp server/manifest/config/config.example.yaml server/manifest/config/config.yaml 
# 修改配置文件中的embedding、chat的配置
docker compose up -d
# 浏览器打开 http://localhost:8000
```

### 使用源码启动
如果有可用的es8和mysql,可以直接快速启动项目，否则需要先安装es8和mysql  
需要修改`config.yaml`文件的相关配置
```bash
cp server/manifest/config/config.example.yaml server/manifest/config/config.yaml 
# 修改配置文件中的embedding、chat、mysql、es的配置
make build
make run
# 浏览器打开 http://localhost:8000
````

### 安装依赖
*如果有可用的es8和mysql,可以不用安装*  
安装es8
```bash
docker run -d --name elasticsearch \
  -e "discovery.type=single-node" \
  -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" \
  -e "xpack.security.enabled=false" \
  -p 9200:9200 \
  -p 9300:9300 \
  elasticsearch:8.18.0
```
安装mysql
```bash
docker run -p 3306:3306 --name mysql \
    -v /Users/wanna/docker/mysql/log:/var/log/mysql \
    -v /Users/wanna/docker/mysql/data:/var/lib/mysql \
    --restart=always \
    -e MYSQL_ROOT_PASSWORD=123456 \
    -e MYSQL_DATABASE=go-rag \
    -d mysql:8.0
```

### 运行 api 项目

```bash
cd server
go mod tidy
go run main.go
```

### 运行前端项目

```bash
cd fe
npm install
npm run dev
```

## 使用Makefile构建

- 构建前端并将产物复制到server/static/fe目录 `make build-fe`

- 构建后端 `make build-server`

- 构建整个项目（前端+后端）`make build`

- 清理构建产物 `make clean`

