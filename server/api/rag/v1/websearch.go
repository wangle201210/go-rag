package v1

import (
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
)

type WebSearchReq struct {
	g.Meta   `path:"/v1/websearch" method:"post" tags:"rag"`
	Question string `json:"question" v:"required"`
	TopK     int    `json:"top_k"` // 默认为5
}

type WebSearchRes struct {
	g.Meta   `mime:"application/json"`
	Document []*schema.Document `json:"document"`
}
