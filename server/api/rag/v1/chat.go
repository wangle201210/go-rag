package v1

import (
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
)

type ChatReq struct {
	g.Meta        `path:"/v1/chat" method:"post" tags:"rag"`
	ConvID        string  `json:"conv_id" v:"required"` // 会话id
	Question      string  `json:"question" v:"required"`
	KnowledgeName string  `json:"knowledge_name" v:"required"`
	TopK          int     `json:"top_k"` // 默认为5
	Score         float64 `json:"score"` // 默认为0.2
}

type ChatRes struct {
	g.Meta     `mime:"application/json"`
	Answer     string             `json:"answer"`
	References []*schema.Document `json:"references"`
}

// ChatStreamReq 流式输出请求
type ChatStreamReq struct {
	g.Meta        `path:"/v1/chat/stream" method:"post" tags:"rag"`
	ConvID        string  `json:"conv_id" v:"required"` // 会话id
	Question      string  `json:"question" v:"required"`
	KnowledgeName string  `json:"knowledge_name" v:"required"`
	TopK          int     `json:"top_k"` // 默认为5
	Score         float64 `json:"score"` // 默认为0.2
}

// ChatStreamRes 流式输出响应
type ChatStreamRes struct {
	g.Meta `mime:"text/event-stream"`
	// 流式输出不需要返回具体内容，内容通过HTTP响应流返回
}
