package grader

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
)

type Grader struct {
	cm model.BaseChatModel
}

func NewGrader(cm model.BaseChatModel) *Grader {
	return &Grader{
		cm: cm,
	}
}

// Retriever 检查下检索到的结果是否能够回答当前问题
func (x *Grader) Retriever(ctx context.Context, docs []*schema.Document, question string) (pass bool, err error) {
	messages, err := retrieverMessages(docs, question)
	if err != nil {
		return
	}
	result, err := x.cm.Generate(ctx, messages)
	if err != nil {
		return false, fmt.Errorf("检查下检索到的结果是否能够回答当前问题失败: %v", err)
	}
	pass = isPass(result.Content)
	return
}

func (x *Grader) Related(ctx context.Context, doc *schema.Document, question string) (pass bool, err error) {
	messages, err := docRelatedMessages(doc, question)
	if err != nil {
		return
	}
	result, err := x.cm.Generate(ctx, messages)
	if err != nil {
		return false, fmt.Errorf("检查下检索到的结果是否和用户问题相关失败: %v", err)
	}
	pass = isPass(result.Content)
	return
}

func isPass(msg string) bool {
	g.Log().Infof(context.Background(), "isPass: %s", msg)
	msg = strings.ToLower(msg)
	return strings.Contains(msg, "yes")
}
