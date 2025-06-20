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
	var knowledgeName string
	if value, ok := ctx.Value(common.KnowledgeName).(string); ok {
		knowledgeName = value
	} else {
		err = fmt.Errorf("必须提供知识库名称")
		return
	}
	wg := &sync.WaitGroup{}
	for _, doc := range docs {
		wg.Add(1)
		go func(doc *schema.Document) {
			defer wg.Done()
			qaContent, e := getQAContent(ctx, doc, knowledgeName)
			if e != nil {
				g.Log().Errorf(ctx, "getQAContent failed, err=%v", e)
				return
			}
			// 生成QA和内容放在一个chunk的不同字段
			doc.MetaData[common.FieldQAContent] = qaContent
		}(doc)
	}
	wg.Wait()
	return docs, nil
}

func getQAContent(ctx context.Context, doc *schema.Document, knowledgeName string) (qaContent string, err error) {
	// 已经有数据了就不要再生成了
	if s, ok := doc.MetaData[common.FieldQAContent].(string); ok && len(s) > 0 {
		return s, nil
	}
	cm, err := common.GetQAModel(ctx, nil)
	if err != nil {
		return
	}
	generate, err := cm.Generate(ctx, []*schema.Message{
		{
			Role: schema.System,
			Content: fmt.Sprintf("你是一个专业的问题生成助手，任务是从给定的文本中提取或生成可能的问题。你不需要回答这些问题，只需生成问题本身。\n"+
				"知识库名字是：《%s》\n\n"+
				"输出格式：\n"+
				"- 每个问题占一行\n"+
				"- 问题必须以问号结尾\n"+
				"- 避免重复或语义相似的问题\n\n"+
				"生成规则：\n"+
				"- 生成的问题必须严格基于文本内容，不能脱离文本虚构。\n"+
				"- 优先生成事实性问题（如谁、何时、何地、如何）。\n"+
				"- 对于复杂文本，可生成多层次问题（基础事实 + 推理问题）。\n"+
				"- 禁止生成主观或开放式问题（如“你认为...？”）。"+
				"- 数量控制在3-5个", knowledgeName),
		},
		{
			Role:    schema.User,
			Content: doc.Content,
		},
	})
	if err != nil {
		return
	}
	qaContent = generate.Content
	return
}
