package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudwego/eino/components/model"
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
	idxer      compose.Runnable[any, []string]
	idxerAsync compose.Runnable[[]*schema.Document, []string]
	rtrvr      compose.Runnable[string, []*schema.Document]
	qaRtrvr    compose.Runnable[string, []*schema.Document]
	client     *elasticsearch.Client
	cm         model.BaseChatModel

	grader *grader.Grader // 暂时先弃用，使用 grader 会严重影响rag的速度
	conf   *config.Config
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
	buildIndexAsync, err := indexer.BuildIndexerAsync(ctx, conf)
	if err != nil {
		return nil, err
	}
	buildRetriever, err := retriever.BuildRetriever(ctx, conf)
	if err != nil {
		return nil, err
	}
	qaCtx := context.WithValue(ctx, common.RetrieverFieldKey, common.FieldQAContentVector)
	qaRetriever, err := retriever.BuildRetriever(qaCtx, conf)
	if err != nil {
		return nil, err
	}
	cm, err := common.GetChatModel(ctx, nil)
	if err != nil {
		g.Log().Error(ctx, "GetChatModel failed, err=%v", err)
		return nil, err
	}
	return &Rag{
		idxer:      buildIndex,
		idxerAsync: buildIndexAsync,
		rtrvr:      buildRetriever,
		qaRtrvr:    qaRetriever,
		client:     conf.Client,
		cm:         cm,
		conf:       conf,
		// grader:  grader.NewGrader(cm),
	}, nil
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
