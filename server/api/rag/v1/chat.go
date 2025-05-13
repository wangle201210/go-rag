package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type ChatReq struct {
	g.Meta   `path:"/v1/chat" method:"post" tags:"rag"`
	ConvID   string  `json:"conv_id"` // 会话id
	Question string  `json:"question"`
	TopK     int     `json:"top_k"` // 如果需要检索文档则需要传入
	Score    float64 `json:"score"` // 如果需要检索文档则需要传入
}

type ChatRes struct {
	g.Meta `mime:"application/json"`
	Answer string `json:"answer"`
}
