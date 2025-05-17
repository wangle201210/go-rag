# go-rag api
rag api 项目

## api 运行

```bash
$ cd server
$ go mod tidy
$ go run main.go
```

## fe 运行

```bash
$ cd fe
$ npm install
$ npm run dev
```

## 使用Makefile构建

```bash
# 构建前端并将产物复制到server/static/fe目录
$ make build-fe

# 构建后端
$ make build-server

# 构建整个项目（前端+后端）
$ make build

# 清理构建产物
$ make clean
```

## Docker部署

### 构建Docker镜像

```bash
# 使用Makefile构建Docker镜像
$ make docker-build

# 或者直接使用docker命令
$ docker build -t go-rag:latest .
```

### 使用Docker Compose启动

```bash
$ docker-compose up -d
```


## indexer
解析文件并向量化到es ![](./static/indexer.png)
```bash
curl --request POST \
  --url http://localhost:8000/v1/indexer \
  --header 'content-type: multipart/form-data' \
  --form 'file=[object Object]'
```

## retriever
根据用户提问检索文档![](./static/retriever.png)
```bash
curl --request POST \
  --url http://localhost:8000/v1/retriever \
  --header 'Content-Type: application/json' \
  --data '{
    "question":"未来计划",
    "top_k":5,
    "score":0.2
}'
```

## chat
检索到文档后回答用户问题![](./static/chat.png)
```bash
curl --request POST \
  --url http://localhost:8000/v1/chat \
  --header 'Content-Type: application/json' \
  --data '{
    "question":"未来计划",
    "top_k":5,
    "score":0.2,
    "conv_id":"123-abc"
}'
```
conv_id 是会话id，尽量不相关的问题使用不同的会话id，否则上下文过长会导致大模型回答不准确  
可以用 uuid.New().String()