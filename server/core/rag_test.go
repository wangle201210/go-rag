package core

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/elastic/go-elasticsearch/v8"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/wangle201210/go-rag/server/core/config"
)

var ragSvr = &Rag{}
var cfg = &config.Config{}

func _init() {
	ctx := context.Background()
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		log.Printf("NewClient of es8 failed, err=%v", err)
		return
	}
	cfg = &config.Config{
		Client:         client,
		IndexName:      "rag-test1",
		APIKey:         g.Cfg().MustGet(ctx, "embedding.apiKey").String(),
		BaseURL:        g.Cfg().MustGet(ctx, "embedding.baseURL").String(),
		EmbeddingModel: g.Cfg().MustGet(ctx, "embedding.model").String(),
		ChatModel:      g.Cfg().MustGet(ctx, "chat.model").String(),
	}
	ragSvr, err = New(context.Background(), cfg)
	if err != nil {
		log.Printf("New of rag failed, err=%v", err)
		return
	}
}
func TestIndex(t *testing.T) {
	_init()
	ctx := context.Background()
	uriList := []string{
		"./test_file/test.xlsx",
		// "./test_file/readme.md",
		// "./test_file/readme2.md",
		// "./test_file/readme.html",
		// "./test_file/test.pdf",
		// "https://deepchat.thinkinai.xyz/docs/guide/advanced-features/shortcuts.html",
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
		// QA 是异步的，不sleep后面会直接停掉
		time.Sleep(time.Second * 3)
	}
}

func TestIndexAsyncByDocsID(t *testing.T) {
	_init()
	ctx := gctx.New()
	var knowledgeName = "wanna"
	ids, err := ragSvr.indexAsyncByDocsID(ctx, &IndexAsyncByDocsIDReq{
		DocsIDs: []string{
			"8dd15a68-69f9-46e2-9045-f0995c4a0b3c",
		},
		KnowledgeName: knowledgeName,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
		return
	}
	t.Logf("ids: %v", ids)
}

func TestIndexAsync(t *testing.T) {
	_init()
	ctx := gctx.New()
	var knowledgeName = "wanna"
	ids, err := ragSvr.IndexAsync(ctx, &IndexAsyncReq{
		Docs: []*schema.Document{
			{
				ID:      "d5225c05-0536-4440-bc1b-edb8dd776de5",
				Content: "h1:这是一个readme文件，这里有很多内容 \\n# 这是一个readme文件，这里有很多内容",
				MetaData: map[string]any{
					"_extension": ".md",
					"_file_name": "readme.md",
					"_source":    "./test_file/readme.md",
					"h1":         "这是一个readme文件，这里有很多内容",
				},
			},
		},
		KnowledgeName: knowledgeName,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
		return
	}
	t.Logf("ids: %v", ids)
}

func TestIndexAsyncByRetrieve(t *testing.T) {
	_init()
	ctx := gctx.New()
	// 先检索点数据
	var knowledgeName = "deepchat使用文档"
	req := &RetrieveReq{
		Query:         "快捷键有哪些",
		TopK:          2,
		Score:         1.2,
		KnowledgeName: knowledgeName,
	}
	msg, err := ragSvr.Retrieve(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	var msgNew []*schema.Document
	for _, m := range msg {
		// Content + ID + MetaData 都是必传的
		// ID 一定要和原数据保持一致
		// MetaData.qa_content 可以置空也可以修改，置空后会在后续自动生成QA
		// MetaData 的其他字段不建议修改
		docParseExt(m)
		metaData := m.MetaData
		metaData["qa_content"] = "" // 置空,也可以是修改
		item := &schema.Document{
			Content:  "===================== 这个文档被我修改了 =====================\n" + m.Content,
			ID:       m.ID,
			MetaData: m.MetaData, // 修改后的 MetaData
		}
		msgNew = append(msgNew, item)
		t.Logf("content: %v, id: %v", item.Content, item.ID)
	}

	ids, err := ragSvr.IndexAsync(ctx, &IndexAsyncReq{
		Docs:          msgNew,
		KnowledgeName: knowledgeName,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
		return
	}
	t.Logf("ids: %v", ids)
}

func TestRetriever(t *testing.T) {
	_init()
	ctx := context.Background()
	req := &RetrieveReq{
		Query:         "dify知识库配置",
		TopK:          5,
		Score:         1.2,
		KnowledgeName: "deepchat使用文档",
	}
	msg, err := ragSvr.Retrieve(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range msg {
		t.Logf("content: %v, score: %v", m.Content, m.Score())
	}
}

func TestRag_GetKnowledgeList(t *testing.T) {
	ctx := context.Background()
	list, err := ragSvr.GetKnowledgeBaseList(ctx)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("list: %v", list)
}
