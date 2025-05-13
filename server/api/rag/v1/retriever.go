package v1

import (
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
)

type RetrieverReq struct {
	g.Meta   `path:"/v1/retriever" method:"post" tags:"rag"`
	Question string  `json:"question"`
	TopK     int     `json:"top_k"`
	Score    float64 `json:"score"`
}

type RetrieverRes struct {
	g.Meta   `mime:"application/json"`
	Document []*schema.Document `json:"document"`
}
