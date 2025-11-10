# 向量存储抽象层

本模块提供了向量存储的抽象接口，支持多种向量数据库实现。

## 支持的向量存储

- **Elasticsearch (ES)**
- **Qdrant**

## 配置说明

### 使用 Elasticsearch

```yaml
vector:
  type: "es"  # 或 "elasticsearch"
  indexName: "rag-test"
  es:
    address: "http://elasticsearch:9200"
    username: "elastic"  # 可选
    password: "123456"   # 可选
```

### 使用 Qdrant

```yaml
vector:
  type: "qdrant"
  indexName: "rag-test"
  qdrant:
    address: "http://qdrant:6333"
    apiKey: ""  # 可选，如果需要认证
```

## 接口说明

### VectorStore 接口

```go
type VectorStore interface {
    // 创建索引/集合
    CreateIndex(ctx context.Context, indexName string) error
    
    // 检查索引/集合是否存在
    IndexExists(ctx context.Context, indexName string) (bool, error)
    
    // 删除文档
    DeleteDocument(ctx context.Context, indexName, documentID string) error
    
    // 获取知识库列表
    GetKnowledgeBaseList(ctx context.Context, indexName string) ([]string, error)
    
    // 搜索文档
    SearchDocuments(ctx context.Context, req *SearchRequest) (*SearchResponse, error)
    
    // 关闭连接
    Close() error
}
```

## 实现新的向量存储

要添加新的向量存储实现：

1. 在 `vector` 包中创建新文件（如 `milvus.go`）
2. 实现 `VectorStore` 接口
3. 在 `factory.go` 中添加新类型的支持
4. 更新配置结构

## 迁移说明

### 从旧配置迁移

旧配置格式：
```yaml
es:
  address: "http://elasticsearch:9200"
  indexName: "rag-test"
  username: "elastic"
  password: "123456"
```

新配置格式：
```yaml
vector:
  type: "es"
  indexName: "rag-test"
  es:
    address: "http://elasticsearch:9200"
    username: "elastic"
    password: "123456"
```

### 代码迁移

旧代码：
```go
err := common.CreateIndexIfNotExists(ctx, client, indexName)
err := common.DeleteDocument(ctx, client, documentID)
```

新代码：
```go
exists, err := vectorStore.IndexExists(ctx, indexName)
if !exists {
    err = vectorStore.CreateIndex(ctx, indexName)
}
err := vectorStore.DeleteDocument(ctx, indexName, documentID)
```

## TODO

- [ ] 添加单元测试
- [ ] 添加 Milvus 支持
- [ ] 添加 Pinecone 支持
