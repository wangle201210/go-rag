package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type IndexerReq struct {
	g.Meta        `path:"/v1/indexer" method:"post" mime:"multipart/form-data" tags:"rag"`
	File          *ghttp.UploadFile `p:"file" type:"file" dc:"如果是本地文件，怎上传文件"`
	URL           string            `p:"url" dc:"如果是网络文件则直接输入url即可"`
	KnowledgeName string            `p:"knowledge_name" dc:"知识库名称" v:"required"`
}

type IndexerRes struct {
	g.Meta `mime:"application/json"`
	DocIDs []string `json:"doc_ids"`
}
