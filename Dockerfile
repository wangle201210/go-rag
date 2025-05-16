# 多阶段构建
# 阶段2: 构建后端
FROM golang:1.23.9-alpine AS server-builder
WORKDIR /app
# 复制整个项目代码
COPY . .
# 复制前端构建产物到server目录
# 构建后端
RUN cd server && go mod tidy && go build -o go-rag-server main.go

# 阶段3: 最终镜像
FROM alpine:latest
WORKDIR /app
# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata
# 设置时区
ENV TZ=Asia/Shanghai
# 复制后端构建产物
COPY --from=server-builder /app/server/go-rag-server /app/
COPY --from=server-builder /app/server/static/ /app/static/
COPY --from=server-builder /app/server/manifest/ /app/manifest/

# 暴露端口
EXPOSE 8000

# 启动命令
CMD ["/app/go-rag-server"]