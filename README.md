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
    ids, err := ragSvr.Index("./test_file/readme.md")
    if err != nil {
        t.Fatal(err)
    }
    for _, id := range ids {
        t.Log(id)
    }
    ragSvr.Index("./test_file/readme2.md")
    ragSvr.Index("./test_file/readme.html")
    ragSvr.Index("./test_file/test.pdf")
    ragSvr.Index("https://deepchat.thinkinai.xyz/docs/guide/advanced-features/shortcuts.html")
    ... ...
```
检索
```go
    msg, err := ragSvr.Retrieve("这里有很多内容", 1.5, 5)
    if err != nil {
        t.Fatal(err)
    }
    for _, m := range msg {
        t.Logf("content: %v, score: %v", m.Content, m.Score())
    }
    msg, err = ragSvr.Retrieve("代码解析", 1.5, 5)
    ... ...
```
详情可以参照[test文件](./rag_test.go)
