# go-rag
基于eino实现知识库的rag

## 存储层
- [x] es8存储向量相关数据

## 功能列表
- [x] md、pdf、html 文档解析
- [x] 网页解析
- [x] 文档检索

## 未来计划
- [ ] 使用mysql存储chunk和文档的映射关系，目前放在es的ext字段

## 使用
安装依赖
```bash
go get github.com/wangle201210/go-rag@latest
```
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
配置环境变量(也可以在[文件](./common/embedding.go)中直接配置)
```bash 
export OPENAI_API_KEY=YOUR_OPENAI_API_KEY
export OPENAI_BASE_URL=YOUR_OPENAI_BASE_URL
```
加载各种数据源的数据，并将其向量化后存储进向量数据库。
```golang
// Index
// uri: 文档地址，可以是文件路径（pdf，html，md等），也可以是网址
// ids: 文档id
func Index(uri string) (ids []string, err error) {
    buildIndex, err := indexer.BuildIndexer(context.Background())
    if err != nil {
        return
    }
    s := document.Source{
        URI: uri,
    }
	// 这个buildIndex可以复用
    ids, err = buildIndex.Invoke(context.Background(), s)
    if err != nil {
        return
    }
    return
}


```
检索
```go
// Retrieve
// input: 检索关键词
// score: 0-2, 0 完全相反，1 毫不相干，2 完全相同
func Retrieve(input string, score float64) (msg []*schema.Document, err error) {
	r, err := retriever.BuildRetriever(context.Background())
	if err != nil {
		return
	}
	// 这个r可以复用
	msg, err = r.Invoke(context.Background(), input,
		compose.WithRetrieverOption(
			er.WithScoreThreshold(score),
			er.WithTopK(5),
		),
	)
	if err != nil {
		return
	}
	return
}
```
详情可以参照[test文件](./rag_test.go)
