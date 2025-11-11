package core

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/wangle201210/go-rag/server/core/common"
	"github.com/wangle201210/go-rag/server/core/config"
	"github.com/wangle201210/go-rag/server/core/grader"
	"github.com/wangle201210/go-rag/server/core/indexer"
	"github.com/wangle201210/go-rag/server/core/retriever"
	coretypes "github.com/wangle201210/go-rag/server/core/types"
)

const (
	scoreThreshold = 1.05 // 设置一个很小的阈值
	esTopK         = 50
	esTryFindDoc   = 10
)

type Rag struct {
	idxer      compose.Runnable[any, []string]
	idxerAsync compose.Runnable[[]*schema.Document, []string]
	rtrvr      compose.Runnable[string, []*schema.Document]
	qaRtrvr    compose.Runnable[string, []*schema.Document]
	client     *elasticsearch.Client // 保留用于兼容
	cm         model.BaseChatModel

	grader *grader.Grader // 暂时先弃用，使用 grader 会严重影响rag的速度
	conf   *config.Config
}

func New(ctx context.Context, conf *config.Config) (*Rag, error) {
	if len(conf.IndexName) == 0 {
		return nil, fmt.Errorf("indexName is empty")
	}
	// 确保 index 存在
	exists, err := conf.IndexExists(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = conf.CreateIndex(ctx)
		if err != nil {
			return nil, err
		}
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
	qaCtx := context.WithValue(ctx, coretypes.RetrieverFieldKey, coretypes.FieldQAContentVector)
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
	return x.conf.GetKnowledgeBaseList(ctx)
}
