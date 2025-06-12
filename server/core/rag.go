package core

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/cloudwego/eino-ext/components/retriever/es8"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/model"
	er "github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/wangle201210/go-rag/server/core/common"
	"github.com/wangle201210/go-rag/server/core/config"
	"github.com/wangle201210/go-rag/server/core/grader"
	"github.com/wangle201210/go-rag/server/core/indexer"
	"github.com/wangle201210/go-rag/server/core/retriever"
)

const (
	scoreThreshold = 1.05 // 设置一个很小的阈值
	esTopK         = 10
)

type Rag struct {
	idxer   compose.Runnable[any, []string]
	rtrvr   compose.Runnable[string, []*schema.Document]
	qaRtrvr compose.Runnable[string, []*schema.Document]
	client  *elasticsearch.Client
	cm      model.BaseChatModel
	grader  *grader.Grader
}

func New(ctx context.Context, conf *config.Config) (*Rag, error) {
	if len(conf.IndexName) == 0 {
		return nil, fmt.Errorf("indexName is empty")
	}
	// 确保es index存在
	err := common.CreateIndexIfNotExists(ctx, conf.Client, conf.IndexName)
	if err != nil {
		return nil, err
	}
	buildIndex, err := indexer.BuildIndexer(ctx, conf)
	if err != nil {
		return nil, err
	}
	buildRetriever, err := retriever.BuildRetriever(ctx, conf)
	if err != nil {
		return nil, err
	}
	QActx := context.WithValue(ctx, common.RetrieverFieldKey, common.FieldQAContentVector)
	qaRetriever, err := retriever.BuildRetriever(QActx, conf)
	if err != nil {
		return nil, err
	}
	cm, err := common.GetChatModel(ctx, conf.GetChatModelConfig())
	if err != nil {
		g.Log().Error(ctx, "GetChatModel failed, err=%v", err)
		return nil, err
	}
	// err = indexer.InitQaIndexer(ctx, conf)
	// if err != nil {
	// 	g.Log().Error(ctx, "InitQaIndexer failed, err=%v", err)
	// 	return nil, err
	// }
	return &Rag{
		idxer:   buildIndex,
		rtrvr:   buildRetriever,
		qaRtrvr: qaRetriever,
		client:  conf.Client,
		cm:      cm,
		grader:  grader.NewGrader(cm),
	}, nil
}

type IndexReq struct {
	URI           string // 文档地址，可以是文件路径（pdf，html，md等），也可以是网址
	KnowledgeName string // 知识库名称
}

// Index
// uri:
// ids: 文档id
func (x *Rag) Index(ctx context.Context, req *IndexReq) (ids []string, err error) {
	s := document.Source{
		URI: req.URI,
	}
	ctx = context.WithValue(ctx, common.KnowledgeName, req.KnowledgeName)
	ids, err = x.idxer.Invoke(ctx, s)
	if err != nil {
		return
	}
	return
}

type RetrieveReq struct {
	Query         string   // 检索关键词
	TopK          int      // 检索结果数量
	Score         float64  //  分数阀值(0-2, 0 完全相反，1 毫不相干，2 完全相同,一般需要传入一个大于1的数字，如1.5)
	KnowledgeName string   // 知识库名字
	optQuery      string   // 优化后的检索关键词
	excludeIDs    []string // 要排除的 _id 列表
	rankScore     float64  // 排名分数
}

// Retrieve 检索
func (x *Rag) Retrieve(ctx context.Context, req *RetrieveReq) (msg []*schema.Document, err error) {
	var (
		used        = ""
		relatedDocs = &sync.Map{}
		docNum      = 0
	)
	req.rankScore = req.Score
	// 大于1的需要-1
	if req.rankScore >= 1 {
		req.rankScore -= 1
	}
	// 最多尝试N次,后续做成可配置
	for i := 0; i < 3; i++ {
		question := req.Query
		var (
			messages []*schema.Message
			generate *schema.Message
			docs     []*schema.Document
			qaDocs   []*schema.Document
			pass     bool
		)
		messages, err = getOptimizedQueryMessages(used, question, req.KnowledgeName)
		if err != nil {
			return
		}
		generate, err = x.cm.Generate(ctx, messages)
		if err != nil {
			return
		}
		optimizedQuery := generate.Content
		used += optimizedQuery + " "
		req.optQuery = optimizedQuery
		docs, err = x.retrieve(ctx, req, false)
		if err != nil {
			return
		}
		qaDocs, err = x.retrieve(ctx, req, true)
		if err != nil {
			return
		}
		docs = append(docs, qaDocs...)
		docs, err = retriever.NewRerank(ctx, req.optQuery, docs, req.TopK)
		if err != nil {
			return
		}
		wg := &sync.WaitGroup{}
		for _, doc := range docs {
			// 分数不够的直接不管
			if doc.Score() < req.rankScore {
				continue
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				// 检查下检索到的结果是否和用户问题相关
				pass, err = x.grader.Related(ctx, doc, req.Query)
				if err != nil {
					return
				}
				if pass {
					req.excludeIDs = append(req.excludeIDs, doc.ID) // 后续不要检索这个_id对应的数据
					relatedDocs.Store(doc.ID, doc)
					docNum++
				} else {
					g.Log().Infof(ctx, "not doc score: %v, related: %v", doc.Score(), doc.Content)
				}
			}()
		}
		wg.Wait()
		// 数量不够就再次检索
		if docNum < req.TopK {
			continue
		}
		// 数量够了，就直接返回
		rDocs := make([]*schema.Document, 0, req.TopK)
		relatedDocs.Range(func(key, value any) bool {
			rDocs = append(rDocs, value.(*schema.Document))
			return true
		})
		pass, err = x.grader.Retriever(ctx, rDocs, req.Query)
		if err != nil {
			return
		}
		if pass {
			return rDocs, nil
		}
	}
	// 最后数量不够，就返回所有数据
	relatedDocs.Range(func(key, value any) bool {
		msg = append(msg, value.(*schema.Document))
		return true
	})
	return
}

func (x *Rag) retrieve(ctx context.Context, req *RetrieveReq, qa bool) (msg []*schema.Document, err error) {
	g.Log().Infof(ctx, "query: %v", req.optQuery)
	esQuery := []types.Query{
		{
			Bool: &types.BoolQuery{
				Must: []types.Query{{Match: map[string]types.MatchQuery{common.KnowledgeName: {Query: req.KnowledgeName}}}},
			},
		},
	}
	if len(req.excludeIDs) > 0 {
		esQuery[0].Bool.MustNot = []types.Query{
			{
				Terms: &types.TermsQuery{
					TermsQuery: map[string]types.TermsQueryField{
						"_id": req.excludeIDs,
					},
				},
			},
		}
	}
	r := x.rtrvr
	if qa {
		r = x.qaRtrvr
	}
	msg, err = r.Invoke(ctx, req.optQuery,
		compose.WithRetrieverOption(
			// er.WithScoreThreshold(req.Score),
			er.WithTopK(esTopK),
			es8.WithFilters(esQuery),
		),
	)
	for _, s := range msg {
		s.WithScore(s.Score() - 1)
	}
	if err != nil {
		return
	}
	return
}

// GetKnowledgeBaseList 获取知识库列表
func (x *Rag) GetKnowledgeBaseList(ctx context.Context) (list []string, err error) {
	names := "distinct_knowledge_names"
	query := search.NewRequest()
	query.Size = common.Of(0) // 不返回原始文档
	query.Aggregations = map[string]types.Aggregations{
		names: {
			Terms: &types.TermsAggregation{
				Field: common.Of(common.KnowledgeName),
				Size:  common.Of(10000),
			},
		},
	}
	res, err := search.NewSearchFunc(x.client)().
		Request(query).
		Do(ctx)
	if err != nil {
		return
	}
	if res.Aggregations == nil {
		g.Log().Infof(ctx, "No aggregations found")
		return
	}
	termsAgg, ok := res.Aggregations[names].(*types.StringTermsAggregate)
	if !ok || termsAgg == nil {
		err = errors.New("failed to parse terms aggregation")
		return
	}
	for _, bucket := range termsAgg.Buckets.([]types.StringTermsBucket) {
		list = append(list, bucket.Key.(string))
	}
	return
}
