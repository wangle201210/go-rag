package v1

import (
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
)

type RetrieverReq struct {
	g.Meta        `path:"/v1/retriever" method:"post" tags:"rag"`
	Question      string  `json:"question" v:"required"`
	TopK          int     `json:"top_k"` // 默认为5
	Score         float64 `json:"score"` // 默认为0.2
	KnowledgeName string  `json:"knowledge_name" v:"required"`
}

type RetrieverRes struct {
	g.Meta   `mime:"application/json"`
	Document []*schema.Document `json:"document"`
}
