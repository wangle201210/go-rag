# Makefile for go-rag project

.PHONY: build-fe build-server build docker-build clean

# 默认目标
all: build

# 构建前端
build-fe:
	cd fe && npm install && npm run build
	mkdir -p server/static/fe
	cp -r fe/dist/* server/static/fe/

# 构建后端
build-server:
	cd server && go mod tidy && go build -o go-rag-server main.go

# 构建整个项目
build: build-fe build-server

# 清理构建产物
clean:
	rm -rf fe/dist
	rm -rf server/static/fe
	rm -f server/go-rag-server

# 构建Docker镜像
docker-build: build
	docker build -t go-rag:latest -f Dockerfile .