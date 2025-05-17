# go-rag
基于eino实现知识库的rag

## 存储层
- [x] es8存储向量相关数据

## 功能列表
- [x] md、pdf、html 文档解析
- [x] 网页解析
- [x] 文档检索
- [x] 长文档自动切割(chunk)
- [x] rerank
- [x] 提供http接口 [rag-api](./server/README.md)

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
  
  
  docker run -d --name elasticsearch \
   -p 9200:9200 \
   -p 9300:9300 \
   -e "discovery.type=single-node" \
   elasticsearch:8.18.0
```
安装mysql
```bash
docker run -p 3306:3306 --name mysql \
    -v /Users/wanna/docker/mysql/log:/var/log/mysql \
    -v /Users/wanna/docker/mysql/data:/var/lib/mysql \
    --restart=always \
    -e MYSQL_ROOT_PASSWORD=123456 \
    -d mysql:8.0
```
初始化Rag对象
```go
    client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		log.Printf("NewClient of es8 failed, err=%v", err)
		return
	}
	ragSvr, err = New(context.Background(), &config.Config{
		Client:    client,
		IndexName: "rag",
		APIKey:    os.Getenv("OPENAI_API_KEY"),
		BaseURL:   os.Getenv("OPENAI_BASE_URL"),
		Model:     "text-embedding-3-large",
	})
	if err != nil {
		log.Printf("New of rag failed, err=%v", err)
		return
	}
```
加载各种数据源的数据，并将其向量化后存储进向量数据库。
```golang
func TestIndex(t *testing.T) {
	ctx := context.Background()
	uriList := []string{
		"./test_file/readme.md",
		"./test_file/readme2.md",
		"./test_file/readme.html",
		"./test_file/test.pdf",
		"https://deepchat.thinkinai.xyz/docs/guide/advanced-features/shortcuts.html",
	}
	for _, s := range uriList {
		req := &IndexReq{
			URI:           s,
			KnowledgeName: "wanna",
		}
		ids, err := ragSvr.Index(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		for _, id := range ids {
			t.Log(id)
		}
	}
}
```
检索
```go
func TestRetriever(t *testing.T) {
	ctx := context.Background()
	req := &RetrieveReq{
		Query:         "这里有很多内容",
		TopK:          5,
		Score:         1.2,
		KnowledgeName: "wanna",
	}
	msg, err := ragSvr.Retrieve(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range msg {
		t.Logf("content: %v, score: %v", m.Content, m.Score())
	}
}
```
详情可以参照[test文件](./rag_test.go)
