package indexer

import (
	"context"
	"fmt"
	"sync"

	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/wangle201210/go-rag/server/core/common"
)

func qa(ctx context.Context, docs []*schema.Document) (output []*schema.Document, err error) {
	wg := &sync.WaitGroup{}
	var knowledgeName string
	if value, ok := ctx.Value(common.KnowledgeName).(string); ok {
		knowledgeName = value
	} else {
		err = fmt.Errorf("必须提供知识库名称")
		return
	}
	for _, doc := range docs {
		wg.Add(1)
		go func(doc *schema.Document) {
			defer wg.Done()
			qaContent, e := getQAContent(ctx, nil, doc, knowledgeName) // 生成QA直接和内容放一起
			if e != nil {
				g.Log().Errorf(ctx, "getQAContent failed, err=%v", e)
				return
			}
			doc.MetaData[common.FieldQAContent] = qaContent
		}(doc)
	}
	wg.Wait()
	return docs, nil
}
