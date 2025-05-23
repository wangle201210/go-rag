// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package rag

import (
	"context"

	"github.com/wangle201210/go-rag/server/api/rag/v1"
)

type IRagV1 interface {
	Chat(ctx context.Context, req *v1.ChatReq) (res *v1.ChatRes, err error)
	ChatStream(ctx context.Context, req *v1.ChatStreamReq) (res *v1.ChatStreamRes, err error)
	Indexer(ctx context.Context, req *v1.IndexerReq) (res *v1.IndexerRes, err error)
	Retriever(ctx context.Context, req *v1.RetrieverReq) (res *v1.RetrieverRes, err error)
}
