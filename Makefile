# Makefile for go-rag project

.PHONY: build-fe build-server build docker-build clean build-all release clean-release

# 默认目标
all: build

# 构建前端
build-fe:
	cd fe && pnpm install && pnpm run build
	mkdir -p server/static/fe
	cp -r fe/dist/* server/static/fe/

# 构建后端
build-server:
	cd server && go mod tidy && go build -o go-rag-server main.go

# 构建整个项目
build: build-fe build-server

# 运行
run:
	cd server && ./go-rag-server

# 清理构建产物
clean:
	rm -rf fe/dist
	rm -rf server/static/fe
	rm -f server/go-rag-server

# 构建Docker镜像
docker-build: build
	docker build -t go-rag:latest -f Dockerfile .

run-local:
	cd server && go mod tidy && go run .

build-linux:
	cd server && go mod tidy && GOOS=linux GOARCH=amd64 go build -o go-rag-server

run-by-docker:
	docker compose -f docker-compose.yml up -d

v := v0.0.3
buildx:
	docker buildx build \
		--platform linux/arm64,linux/amd64 \
		-t iwangle/go-rag:$(v) \
		--push \
		.

# 项目名称和版本
APP_NAME := go-rag

# 支持的平台
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# 多平台构建
build-all:
	@echo "Building go-rag for multiple platforms..."
	@mkdir -p releases
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d'/' -f1); \
		arch=$$(echo $$platform | cut -d'/' -f2); \
		output_name=$(APP_NAME)-$$os-$$arch; \
		if [ $$os = "windows" ]; then output_name=$$output_name.exe; fi; \
		echo "Building for $$os/$$arch..."; \
		(cd server && GOOS=$$os GOARCH=$$arch go build -o ../releases/$$output_name .); \
	done
	@echo "Build completed! Files are in releases/ directory"

# 发布版本（构建 + 压缩）
release: clean-release build-all
	@echo "Creating release archives..."
	@cd releases && \
	for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d'/' -f1); \
		arch=$$(echo $$platform | cut -d'/' -f2); \
		release_dir=$(APP_NAME); \
		output_name=$(APP_NAME); \
		if [ $$os = "windows" ]; then output_name=$$output_name.exe; fi; \
		exe_file=$(APP_NAME)-$$os-$$arch; \
		if [ $$os = "windows" ]; then exe_file=$$exe_file.exe; fi; \
		archive_name=$(APP_NAME)-$$os-$$arch; \
		\
		echo "Preparing $$archive_name..."; \
		mkdir -p $$release_dir; \
		\
		if [ -f $$exe_file ]; then \
			cp $$exe_file $$release_dir/$$output_name; \
			echo "Copied executable to $$release_dir/$$output_name"; \
		else \
			echo "Warning: $$exe_file not found, skipping $$platform"; \
			rm -rf $$release_dir; \
			continue; \
		fi; \
		\
		if [ -d ../server/static ]; then \
			cp -r ../server/static $$release_dir/; \
			echo "Copied static files to $$release_dir/static/"; \
		fi; \
		\
		if [ -f ../server/manifest/config/config_qd_demo.yaml ]; then \
			cp ../server/manifest/config/config_qd_demo.yaml $$release_dir/config.yaml; \
			echo "Copied config file to $$release_dir/config.yaml"; \
		fi; \
		\
		if [ $$os = "windows" ]; then \
			zip -q $$archive_name.zip -r $$release_dir; \
			echo "Created $$archive_name.zip"; \
		else \
			tar -czf $$archive_name.tar.gz $$release_dir; \
			echo "Created $$archive_name.tar.gz"; \
		fi; \
		rm -rf $$release_dir; \
	done
	@echo "Release archives created!"
	@echo "Files in releases/:"
	@ls -la releases/

# 清理发布文件
clean-release:
	@echo "Cleaning release files..."
	@rm -rf releases/
	@echo "Release files cleaned!"

