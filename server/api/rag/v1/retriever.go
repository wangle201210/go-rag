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

type RetrieverDifyReq struct {
	g.Meta           `path:"/v1/dify/retrieval" method:"post" tags:"rag" no_wrap_resp:"true"`
	KnowledgeID      string            `json:"knowledge_id" v:"required"`
	Query            string            `json:"query" v:"required"`
	RetrievalSetting *RetrievalSetting `json:"retrieval_setting" v:"required"`
	// MetadataCondition map[string]interface{} `json:"metadata_condition"`
}

type RetrievalSetting struct {
	TopK           int     `json:"top_k"`
	ScoreThreshold float64 `json:"score_threshold"`
}
type RetrieverDifyRes struct {
	g.Meta  `mime:"application/json"`
	Records []*Record `json:"records"`
}

type Record struct {
	Metadata *Metadata `json:"metadata"`
	Score    float64   `json:"score"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
}

type Metadata struct {
	Path        string `json:"path"`
	Description string `json:"description"`
}
