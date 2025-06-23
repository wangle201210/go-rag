package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/wangle201210/go-rag/server/internal/model/entity"
)

const (
	StatusPending  Status = 0
	StatusIndexing Status = 1
	StatusActive   Status = 2
	StatusFailed   Status = 3
)

type DocumentsListReq struct {
	g.Meta        `path:"/v1/documents" method:"get" tags:"rag"`
	KnowledgeName string `p:"knowledge_name" dc:"knowledge_name" v:"required|length:3,50"`
	Page          int    `p:"page" dc:"page" v:"required|min:1" d:"1"`
	Size          int    `p:"size" dc:"size" v:"required|min:1|max:100" d:"10"`
}

type DocumentsListRes struct {
	g.Meta `mime:"application/json"`
	Data   []entity.KnowledgeDocuments `json:"data"`
	Total  int                         `json:"total"`
	Page   int                         `json:"page"`
	Size   int                         `json:"size"`
}

type DocumentsDeleteReq struct {
	g.Meta     `path:"/v1/documents" method:"delete" tags:"rag" summary:"Delete a document and its chunks"`
	DocumentId int `p:"document_id" dc:"document_id" v:"required"`
}

type DocumentsDeleteRes struct {
	g.Meta `mime:"application/json"`
}
