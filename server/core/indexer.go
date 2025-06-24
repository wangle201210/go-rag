package core

import (
	"context"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	v1 "github.com/wangle201210/go-rag/server/api/rag/v1"
	"github.com/wangle201210/go-rag/server/core/common"
	"github.com/wangle201210/go-rag/server/core/retriever"
	"github.com/wangle201210/go-rag/server/internal/logic/knowledge"
	"github.com/wangle201210/go-rag/server/internal/model/entity"
)

type IndexReq struct {
	URI           string // 文档地址，可以是文件路径（pdf，html，md等），也可以是网址
	KnowledgeName string // 知识库名称
	DocumentsId   int64  // 文档ID
}

type IndexAsyncReq struct {
	Docs          []*schema.Document
	KnowledgeName string // 知识库名称
	DocumentsId   int64  // 文档ID
}

type IndexAsyncByDocsIDReq struct {
	DocsIDs       []string
	KnowledgeName string // 知识库名称
	DocumentsId   int64  // 文档ID
}

// Index
// 这里处理文档的读取、分割、合并和存储
// 真正的embedding 和 QA 异步执行
func (x *Rag) Index(ctx context.Context, req *IndexReq) (ids []string, err error) {
	s := document.Source{
		URI: req.URI,
	}
	ctx = context.WithValue(ctx, common.KnowledgeName, req.KnowledgeName)
	ids, err = x.idxer.Invoke(ctx, s)
	if err != nil {
		return
	}
	go func() {
		// 测试下来这里必须 sleep 一段时间，否则下面的 indexAsyncByDocsID 在es里面搜索不到该条数据，应该是es本身会有一定延迟
		// 这里会有一定隐患，刚提交index后项目就崩了，可能会有几条chunk没有生成QA
		// 但是这个场景几乎不会出现，且不影响用户使用，可以忽略
		time.Sleep(time.Second)
		ctxN := gctx.New()
		defer func() {
			if e := recover(); e != nil {
				g.Log().Errorf(ctxN, "recover indexAsyncByDocsID failed, err=%v", e)
			}
		}()
		_, err = x.indexAsyncByDocsID(ctxN, &IndexAsyncByDocsIDReq{
			DocsIDs:       ids,
			KnowledgeName: req.KnowledgeName,
			DocumentsId:   req.DocumentsId,
		})
		if err != nil {
			g.Log().Errorf(ctxN, "indexAsyncByDocsID failed, err=%v", err)
		}
	}()
	return
}

// IndexAsync
// 通过 schema.Document 异步 生成QA&embedding
func (x *Rag) IndexAsync(ctx context.Context, req *IndexAsyncReq) (ids []string, err error) {
	ctx = context.WithValue(ctx, common.KnowledgeName, req.KnowledgeName)
	ids, err = x.idxerAsync.Invoke(ctx, req.Docs)
	if err != nil {
		return
	}

	return
}

// 通过docIDs 异步 生成QA&embedding
// 这个方法不用暴露出去
func (x *Rag) indexAsyncByDocsID(ctx context.Context, req *IndexAsyncByDocsIDReq) (ids []string, err error) {
	esQuery := &types.Query{
		Bool: &types.BoolQuery{
			Must: []types.Query{
				{Match: map[string]types.MatchQuery{common.KnowledgeName: {Query: req.KnowledgeName}}},
				{Terms: &types.TermsQuery{TermsQuery: map[string]types.TermsQueryField{"_id": req.DocsIDs}}},
			},
		},
	}

	sreq := search.NewRequest()
	sreq.Query = esQuery
	resp, err := search.NewSearchFunc(x.client)().
		Index(x.conf.IndexName).
		Request(sreq).
		Do(ctx)
	if err != nil {
		return
	}
	var docs []*schema.Document
	var chunks []entity.KnowledgeChunks
	for _, hit := range resp.Hits.Hits {
		doc := &schema.Document{}
		doc, err = retriever.EsHit2Document(ctx, hit)
		if err != nil {
			g.Log().Errorf(ctx, "EsHit2Document failed, err=%v", err)
			return
		}
		docParseExt(doc)
		docs = append(docs, doc)

		ext, err := sonic.Marshal(doc.MetaData)
		if err != nil {
			g.Log().Errorf(ctx, "sonic.Marshal failed, err=%v", err)
			continue
		}
		chunks = append(chunks, entity.KnowledgeChunks{
			KnowledgeDocId: req.DocumentsId,
			ChunkId:        doc.ID,
			Content:        doc.Content,
			Ext:            string(ext),
		})
	}
	if err = knowledge.SaveChunksData(ctx, req.DocumentsId, chunks); err != nil {
		// 这里不返回err，不影响用户使用
		g.Log().Errorf(ctx, "indexAsyncByDocsID insert chunks failed, err=%v", err)
	}

	asyncReq := &IndexAsyncReq{
		Docs:          docs,
		KnowledgeName: req.KnowledgeName,
		DocumentsId:   req.DocumentsId,
	}
	ids, err = x.IndexAsync(ctx, asyncReq)
	if err != nil {
		return
	}
	knowledge.UpdateDocumentsStatus(ctx, req.DocumentsId, int(v1.StatusActive))
	return
}

// 检索会把原来的 MetaData 放到 MetaData.ext 中，这里需要把原来的 MetaData 恢复
func docParseExt(doc *schema.Document) {
	if ext, ok := doc.MetaData[common.FieldExtra].(string); ok && len(ext) > 0 {
		extData := map[string]any{}
		if err := sonic.Unmarshal([]byte(doc.MetaData[common.FieldExtra].(string)), &extData); err != nil {
			// 忽略err
			g.Log().Errorf(gctx.New(), "documentParseExt unmarshal failed, err=%v", err)
			return
		}
		doc.MetaData = extData
	}
}

func (x *Rag) DeleteDocument(ctx context.Context, documentID string) error {
	return common.DeleteDocument(ctx, x.conf.Client, documentID)
}
