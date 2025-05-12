// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package retriever

import (
	"context"

	"github.com/wangle201210/go-rag/server/api/retriever/v1"
)

type IRetrieverV1 interface {
	Retriever(ctx context.Context, req *v1.RetrieverReq) (res *v1.RetrieverRes, err error)
}
