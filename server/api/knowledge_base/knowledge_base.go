// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package knowledge_base

import (
	"context"

	"github.com/wangle201210/go-rag/server/api/knowledge_base/v1"
)

type IKnowledgeBaseV1 interface {
	KnowledgeBase(ctx context.Context, req *v1.KnowledgeBaseReq) (res *v1.KnowledgeBaseRes, err error)
	CreateKnowledgeBase(ctx context.Context, req *v1.CreateKnowledgeBaseReq) (res *v1.CreateKnowledgeBaseRes, err error)
	UpdateKnowledgeBase(ctx context.Context, req *v1.UpdateKnowledgeBaseReq) (res *v1.UpdateKnowledgeBaseRes, err error)
	DeleteKnowledgeBase(ctx context.Context, req *v1.DeleteKnowledgeBaseReq) (res *v1.DeleteKnowledgeBaseRes, err error)
}
